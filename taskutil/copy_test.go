package taskutil_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/janmarkuslanger/ssgo/task"
	"github.com/janmarkuslanger/ssgo/taskutil"
)

type mockFile struct {
	io.Reader
	io.Writer
	closed bool
}

func (f *mockFile) Close() error {
	f.closed = true
	return nil
}

type mockFS struct {
	openErr    error
	createErr  error
	mkdirErr   error
	openFile   *mockFile
	createFile *mockFile
}

func (fs *mockFS) Open(name string) (*os.File, error) {
	if fs.openErr != nil {
		return nil, fs.openErr
	}
	return (*os.File)(nil), nil
}

func (fs *mockFS) Create(name string) (*os.File, error) {
	if fs.createErr != nil {
		return nil, fs.createErr
	}
	return (*os.File)(nil), nil
}

func (fs *mockFS) MkdirAll(path string, perm os.FileMode) error {
	return fs.mkdirErr
}

type mockCopier struct {
	err error
}

func (c *mockCopier) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return 0, c.err
}

func createTempDir(t *testing.T) string {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	err := os.MkdirAll(filepath.Join(src, "subdir"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(src, "subdir", "file.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestCopyTask_Run_Success(t *testing.T) {
	dir := createTempDir(t)
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")

	copyTask := taskutil.NewCopyTask(src, "assets", nil, nil)

	err := copyTask.Run(task.TaskContext{OutputDir: dst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outFile := filepath.Join(dst, "assets", "subdir", "file.txt")
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("file not copied: %v", err)
	}
	if string(data) != "content" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestCopyTask_Run_OpenError(t *testing.T) {
	fs := &mockFS{openErr: errors.New("open fail")}
	copyTask := taskutil.NewCopyTask(".", "", fs, nil)
	err := copyTask.Run(task.TaskContext{OutputDir: "."})
	if err == nil || err.Error() != "open fail" {
		t.Errorf("expected open fail, got %v", err)
	}
}

func TestCopyTask_Run_CreateError(t *testing.T) {
	fs := &mockFS{createErr: errors.New("create fail")}
	copyTask := taskutil.NewCopyTask(".", "", fs, nil)
	err := copyTask.Run(task.TaskContext{OutputDir: "."})
	if err == nil || err.Error() != "create fail" {
		t.Errorf("expected create fail, got %v", err)
	}
}

func TestCopyTask_Run_MkdirError(t *testing.T) {
	fs := &mockFS{mkdirErr: errors.New("mkdir fail")}
	copyTask := taskutil.NewCopyTask(".", "", fs, nil)
	err := copyTask.Run(task.TaskContext{OutputDir: "."})
	if err == nil || err.Error() != "mkdir fail" {
		t.Errorf("expected mkdir fail, got %v", err)
	}
}

func TestCopyTask_Run_CopyError(t *testing.T) {
	dir := createTempDir(t)
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")

	copyTask := taskutil.NewCopyTask(src, "assets", nil, &mockCopier{err: errors.New("copy fail")})
	err := copyTask.Run(task.TaskContext{OutputDir: dst})
	if err == nil || err.Error() != "copy fail" {
		t.Errorf("expected copy fail, got %v", err)
	}
}

func TestCopyTask_IsCritical(t *testing.T) {
	task := taskutil.NewCopyTask("a", "b", nil, nil)
	if !task.IsCritical() {
		t.Errorf("expected true")
	}
}
