package taskutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/janmarkuslanger/ssgo/task"
)

func NewCopyTask(sourceDir string, outputSubDir string, pathResolver PathResolver) CopyTask {
	if pathResolver == nil {
		pathResolver = defaultPathResolver{}
	}

	return CopyTask{
		SourceDir:    sourceDir,
		OutputSubDir: outputSubDir,
		pathResolver: pathResolver,
	}
}

type PathResolver interface {
	Abs(path string) (string, error)
	Rel(basepath string, targpath string) (string, error)
}

type defaultPathResolver struct{}

func (p defaultPathResolver) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (p defaultPathResolver) Rel(basepath string, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

type CopyTask struct {
	SourceDir    string
	OutputSubDir string
	pathResolver PathResolver
}

func (c *CopyTask) Run(ctx task.TaskContext) error {
	if c.pathResolver == nil {
		return fmt.Errorf("pathresolver not defined")
	}

	srcDirAbs, err := c.pathResolver.Abs(c.SourceDir)
	if err != nil {
		return fmt.Errorf("failed to resolve source dir: %w", err)
	}

	outDir := ctx.OutputDir
	if c.OutputSubDir != "" {
		outDir = filepath.Join(outDir, c.OutputSubDir)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return filepath.Walk(srcDirAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}
		relPath, err := c.pathResolver.Rel(srcDirAbs, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path from %s to %s: %w", srcDirAbs, path, err)
		}

		targetPath := filepath.Join(outDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return CopyFile(path, targetPath, info.Mode())
	})
}

func (c *CopyTask) IsCritical() bool {
	return false
}

type CopyFunc func(dst io.Writer, src io.Reader) (int64, error)

func CopyFile(src, dest string, mode os.FileMode) error {
	return CopyFileWith(src, dest, mode, io.Copy)
}

func CopyFileWith(src, dest string, mode os.FileMode, copier CopyFunc) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src error: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create dest error: %w", err)
	}
	defer destFile.Close()

	if _, err := copier(destFile, srcFile); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return os.Chmod(dest, mode)
}
