package writer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janmarkuslanger/ssgo/writer"
)

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
