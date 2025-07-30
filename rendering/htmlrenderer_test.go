package rendering_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janmarkuslanger/ssgo/rendering"
)

func TestHTMLRenderer_Render_Success(t *testing.T) {
	tmp := t.TempDir()

	layoutPath := filepath.Join(tmp, "layout.html")
	err := os.WriteFile(layoutPath, []byte(`<html><body>{{ template "content" . }}</body></html>`), 0644)
	if err != nil {
		t.Fatalf("could not write layout: %v", err)
	}

	templatePath := filepath.Join(tmp, "index.html")
	err = os.WriteFile(templatePath, []byte(`{{ define "content" }}<h1>Hello {{ .Name }}</h1>{{ end }}`), 0644)
	if err != nil {
		t.Fatalf("could not write template: %v", err)
	}

	renderer := rendering.HTMLRenderer{
		Layout: []string{layoutPath},
	}
	out, err := renderer.Render(rendering.RenderContext{
		Data:     map[string]any{"Name": "Jan"},
		Template: templatePath,
	})
	if err != nil {
		t.Fatalf("rendering failed: %v", err)
	}

	want := `<html><body><h1>Hello Jan</h1></body></html>`
	if out != want {
		t.Errorf("unexpected output:\n%s\nexpected:\n%s", out, want)
	}
}

func TestHTMLRenderer_Render_TemplateNotFound(t *testing.T) {
	renderer := rendering.HTMLRenderer{}
	_, err := renderer.Render(rendering.RenderContext{
		Data:     nil,
		Template: "not-exist.html",
	})
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestHTMLRenderer_Render_UndefinedTemplate(t *testing.T) {
	tmp := t.TempDir()

	layoutPath := filepath.Join(tmp, "layout.html")
	err := os.WriteFile(layoutPath, []byte(`<html><body>{{ template "somethingthatisnotthere" . }}</body></html>`), 0644)
	if err != nil {
		t.Fatalf("could not write layout: %v", err)
	}

	templatePath := filepath.Join(tmp, "index.html")
	err = os.WriteFile(templatePath, []byte(`{{ define "content" }}<h1>Hello {{ .Name }}</h1>{{ end }}`), 0644)
	if err != nil {
		t.Fatalf("could not write template: %v", err)
	}

	renderer := rendering.HTMLRenderer{
		Layout: []string{layoutPath},
	}
	_, err = renderer.Render(rendering.RenderContext{
		Data:     map[string]any{"Name": "Jan"},
		Template: templatePath,
	})
	if err == nil {
		t.Fatal("expected error for undefined content block")
	}
}
