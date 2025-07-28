package writer

import (
	"os"
	"path/filepath"
)

const (
	FilePerm = 0644
	DirPerm  = 0755
)

type FileWriter struct{}

func (w *FileWriter) Write(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), DirPerm); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), FilePerm)
}
