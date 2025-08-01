package taskutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/janmarkuslanger/ssgo/task"
)

type FileSystem interface {
	Open(string) (*os.File, error)
	Create(string) (*os.File, error)
	MkdirAll(string, os.FileMode) error
}

type Copier interface {
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

type CopyTask struct {
	sourceDir    string
	outputSubdir string
	fs           FileSystem
	copier       Copier
}

func NewCopyTask(sourceDir, outputSubdir string, fs FileSystem, copier Copier) CopyTask {
	if fs == nil {
		fs = defaultFS{}
	}
	if copier == nil {
		copier = defaultCopier{}
	}
	return CopyTask{
		sourceDir:    sourceDir,
		outputSubdir: outputSubdir,
		fs:           fs,
		copier:       copier,
	}
}

type defaultFS struct{}

func (defaultFS) Open(name string) (*os.File, error)           { return os.Open(name) }
func (defaultFS) Create(name string) (*os.File, error)         { return os.Create(name) }
func (defaultFS) MkdirAll(path string, perm os.FileMode) error { return os.MkdirAll(path, perm) }

type defaultCopier struct{}

func (defaultCopier) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

func (t CopyTask) Run(ctx task.TaskContext) error {
	return filepath.Walk(t.sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(t.sourceDir, path)
		if err != nil {
			return err
		}

		dest := filepath.Join(ctx.OutputDir, t.outputSubdir, relPath)

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = t.copier.Copy(dstFile, srcFile)
		return err
	})
}

func (t CopyTask) IsCritical() bool {
	return true
}
