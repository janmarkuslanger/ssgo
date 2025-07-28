package writer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janmarkuslanger/ssgo/writer"
)

func TestFileWriter_New(t *testing.T) {
	var fw any = writer.NewFileWriter()

	if _, ok := fw.(*writer.FileWriter); !ok {
		t.Fatalf("doesnt create correct pointer")
	}
}

func TestFileWriter_Write(t *testing.T) {
	tmpDir := t.TempDir()

	writer := writer.FileWriter{}
	path := filepath.Join(tmpDir, "test", "index.html")
	content := "<h1>Hello</h1>"

	err := writer.Write(path, content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file was not created: %v", err)
	}

	if string(data) != content {
		t.Errorf("unexpected content: got %q, want %q", data, content)
	}
}

func TestFileWriter_Write_err(t *testing.T) {
	tmp := t.TempDir()

	conflictPath := filepath.Join(tmp, "foo")
	err := os.WriteFile(conflictPath, []byte("I am a file"), 0644)
	if err != nil {
		t.Fatalf("unexpected error creating file: %v", err)
	}

	writer := writer.FileWriter{}
	targetPath := filepath.Join(conflictPath, "bar", "index.html")
	err = writer.Write(targetPath, "test content")

	if err == nil {
		t.Fatal("expected mkdir to fail, but got no error")
	}
	t.Logf("got expected error: %v", err)
}
