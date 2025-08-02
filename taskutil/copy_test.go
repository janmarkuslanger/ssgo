package taskutil_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/janmarkuslanger/ssgo/task"
	"github.com/janmarkuslanger/ssgo/taskutil"
)

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	fullPath := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		t.Fatalf("mkdir error: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}
	return fullPath
}

type PathResolverFail struct{}

func (p PathResolverFail) Abs(path string) (string, error) {
	return "", errors.New("Fail")
}

func TestCopyTask_New(t *testing.T) {
	sDir := "src"
	dirName := "name"
	ct := taskutil.NewCopyTask(sDir, dirName, nil)

	if ct.SourceDir != sDir {
		t.Fatalf("expected source dir to be %v but got %q", sDir, ct.SourceDir)
	}

	if ct.OutputSubDir != dirName {
		t.Fatalf("expected output dir name to be %v but got %q", dirName, ct.OutputSubDir)
	}
}

func TestCopyTask_CustomPathResolver(t *testing.T) {
	sDir := "src"
	dirName := "name"
	taskutil.NewCopyTask(sDir, dirName, PathResolverFail{})
}

func TestCopyTask_Run(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	createTempFile(t, srcDir, "a.txt", "Hello A")
	createTempFile(t, srcDir, "sub/b.txt", "Hello B")

	copyTask := taskutil.NewCopyTask(srcDir, "nested", nil)

	err := copyTask.Run(task.TaskContext{
		OutputDir: outDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "nested", "a.txt")); err != nil {
		t.Errorf("file not copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "nested", "sub", "b.txt")); err != nil {
		t.Errorf("nested file not copied: %v", err)
	}
}

func TestCopyTask_Run_NoSubdir(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	createTempFile(t, srcDir, "x.txt", "Data X")

	copyTask := taskutil.NewCopyTask(srcDir, "", nil)

	err := copyTask.Run(task.TaskContext{OutputDir: outDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "x.txt")); err != nil {
		t.Errorf("file not copied: %v", err)
	}
}

func TestCopyTask_IsCritical(t *testing.T) {
	copyTask := &taskutil.CopyTask{}
	if copyTask.IsCritical() != false {
		t.Errorf("expected IsCritical to return false")
	}
}

func TestCopyTask_Run_ErrorOnAbsPath(t *testing.T) {
	srcDir := t.TempDir()
	copyTask := taskutil.NewCopyTask(srcDir, "", PathResolverFail{})
	err := copyTask.Run(task.TaskContext{OutputDir: t.TempDir()})
	if err == nil {
		t.Errorf("expected error on invalid source path")
	}

	msg := err.Error()
	expected := "failed to resolve source dir:"
	if !strings.Contains(msg, expected) {
		t.Errorf("expected err msg %v but got %q", expected, msg)
	}
}

func TestCopyTask_Run_ErrorOnMkdirOutput(t *testing.T) {
	srcDir := t.TempDir()
	badOutput := filepath.Join(srcDir, "nonWritable")

	if err := os.MkdirAll(badOutput, 0500); err != nil {
		t.Fatalf("setup error: %v", err)
	}
	defer os.Chmod(badOutput, 0755)

	copyTask := taskutil.NewCopyTask(srcDir, "x", nil)

	err := copyTask.Run(task.TaskContext{OutputDir: badOutput})
	if err == nil {
		t.Errorf("expected mkdir error")
	}
}

func TestCopyTask_Run_ErrorOnWalk(t *testing.T) {
	srcDir := t.TempDir()
	if err := os.RemoveAll(srcDir); err != nil {
		t.Fatalf("failed to remove source dir")
	}

	copyTask := taskutil.NewCopyTask(srcDir, "", nil)

	err := copyTask.Run(task.TaskContext{OutputDir: t.TempDir()})
	if err == nil {
		t.Errorf("expected walk error")
	}
}

func TestCopyTask_Run_ErrorOnFileOpen(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	file := createTempFile(t, srcDir, "x.txt", "content")
	if err := os.Chmod(file, 0000); err != nil {
		t.Fatalf("chmod failed")
	}
	defer os.Chmod(file, 0644)

	copyTask := &taskutil.CopyTask{SourceDir: srcDir}
	err := copyTask.Run(task.TaskContext{OutputDir: outDir})
	if err == nil {
		t.Errorf("expected error on file open")
	}
}

func TestCopyTask_Run_ErrorOnCreateDest(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	createTempFile(t, srcDir, "x.txt", "abc")
	targetFile := filepath.Join(outDir, "x.txt")
	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		t.Fatalf("mkdir error")
	}
	if err := os.WriteFile(targetFile, []byte("x"), 0444); err != nil {
		t.Fatalf("file create error")
	}
	if err := os.Chmod(filepath.Dir(targetFile), 0400); err != nil {
		t.Fatalf("chmod error")
	}
	defer os.Chmod(filepath.Dir(targetFile), 0755)

	copyTask := &taskutil.CopyTask{SourceDir: srcDir}
	err := copyTask.Run(task.TaskContext{OutputDir: outDir})
	if err == nil {
		t.Errorf("expected error on create dest")
	}
}

func TestCopyTask_Run_ErrorOnCopy(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	createTempFile(t, srcDir, "fail.txt", "abc")
	copyTask := &taskutil.CopyTask{SourceDir: srcDir}

	destDir := filepath.Join(outDir, "fail.txt")
	if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
		t.Fatalf("mkdir error")
	}
	if err := os.Mkdir(destDir, 0444); err != nil {
		t.Fatalf("dir create error")
	}
	defer os.RemoveAll(destDir)

	err := copyTask.Run(task.TaskContext{OutputDir: outDir})
	if err == nil {
		t.Errorf("expected create dest error")
	}
}

func TestCopyTask_Run_ErrorOnChmod(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	copyTask := taskutil.NewCopyTask(srcDir, "", nil)

	destPath := filepath.Join(outDir, "chmod.txt")
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		t.Fatalf("mkdir error")
	}
	if err := os.WriteFile(destPath, []byte{}, 0444); err != nil {
		t.Fatalf("write error")
	}
	defer os.Remove(destPath)

	err := copyTask.Run(task.TaskContext{OutputDir: outDir})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCopyFile_OpenSrcFails(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")

	if err := os.WriteFile(srcFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(srcFile, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(srcFile, 0644)

	destFile := filepath.Join(tmpDir, "dest.txt")

	err := taskutil.CopyFile(srcFile, destFile, 0644)
	if err == nil || !strings.Contains(err.Error(), "open src error") {
		t.Errorf("expected open src error, got: %v", err)
	}
}

func TestCopyFile_CreateDestFails(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	destDir := filepath.Join(tmpDir, "no-write")
	destFile := filepath.Join(destDir, "dest.txt")

	if err := os.WriteFile(srcFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(destDir, 0500); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(destDir, 0755)

	err := taskutil.CopyFile(srcFile, destFile, 0644)
	if err == nil || !strings.Contains(err.Error(), "create dest error") {
		t.Errorf("expected create dest error, got: %v", err)
	}
}

func TestCopyFile_CopyFails(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	destFile := filepath.Join(tmpDir, "dest.txt")

	if err := os.WriteFile(srcFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	failingCopier := func(dst io.Writer, src io.Reader) (int64, error) {
		return 0, errors.New("simulated copy failure")
	}

	err := taskutil.CopyFileWith(srcFile, destFile, 0644, failingCopier)
	if err == nil || !strings.Contains(err.Error(), "copy error") {
		t.Errorf("expected copy error, got: %v", err)
	}
}
