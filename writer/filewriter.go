package writer

import (
	"os"
	"path/filepath"
)

const (
	FilePerm = 0644
	DirPerm  = 0755
)

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

type FileWriter struct{}

func (w *FileWriter) Write(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), DirPerm); err != nil {
		return err
	}

	return os.WriteFile(path+".html", []byte(content), FilePerm)
}
