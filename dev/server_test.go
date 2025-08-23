package dev_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/janmarkuslanger/ssgo/builder"
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
	"github.com/janmarkuslanger/ssgo/writer"

	"github.com/janmarkuslanger/ssgo/dev"
)

func makeTempTemplates(t *testing.T) (layoutPath, tplPath string) {
	t.Helper()
	dir := t.TempDir()
	layout := `{{define "root"}}{{template "content" .}}{{end}}`
	tpl := `{{define "content"}}{{.Content}}{{end}}`
	lp := filepath.Join(dir, "layout.html")
	tp := filepath.Join(dir, "page.html")
	if err := os.WriteFile(lp, []byte(layout), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tp, []byte(tpl), 0o600); err != nil {
		t.Fatal(err)
	}
	return lp, tp
}

func makeTestBuilder(t *testing.T) builder.Builder {
	t.Helper()
	layout, tpl := makeTempTemplates(t)
	renderer := rendering.HTMLRenderer{Layout: []string{layout}}

	gen := func(path, body string) page.Generator {
		posts := map[string]map[string]any{path: {"Content": body}}
		return page.Generator{
			Config: page.Config{
				Pattern:  path,
				Template: tpl,
				GetPaths: func() []string { return []string{path} },
				GetData:  func(p page.PagePayload) map[string]any { return posts[path] },
				Renderer: renderer,
			},
		}
	}

	return builder.Builder{
		OutputDir: t.TempDir(),
		Writer:    &writer.FileWriter{},
		Pages: []page.Generator{
			gen("/", "home"),
			gen("/about", "about"),
		},
	}
}

func makeBrokenBuilder(t *testing.T) builder.Builder {
	t.Helper()
	layout, _ := makeTempTemplates(t)
	renderer := rendering.HTMLRenderer{Layout: []string{layout}}

	gen := func(path string) page.Generator {
		return page.Generator{
			Config: page.Config{
				Pattern:  path,
				Template: "does-not-exist.html",
				GetPaths: func() []string { return []string{path} },
				GetData:  func(p page.PagePayload) map[string]any { return map[string]any{} },
				Renderer: renderer,
			},
		}
	}

	return builder.Builder{
		OutputDir: t.TempDir(),
		Writer:    &writer.FileWriter{},
		Pages:     []page.Generator{gen("/broken")},
	}
}

func TestNewServer_RendersRegisteredPaths(t *testing.T) {
	b := makeTestBuilder(t)
	mux := dev.NewServer(b)

	cases := []struct {
		path     string
		expected string
	}{
		{path: "/", expected: "home"},
		{path: "/about", expected: "about"},
	}

	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, c.path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: expected 200, got %d", c.path, rec.Code)
		}
		if rec.Body.String() != c.expected {
			t.Fatalf("GET %s: expected %q, got %q", c.path, c.expected, rec.Body.String())
		}
	}
}

func TestNewServer_PanicsOnRenderError(t *testing.T) {
	b := makeBrokenBuilder(t)
	mux := dev.NewServer(b)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on render error")
		}
	}()
	req := httptest.NewRequest(http.MethodGet, "/broken", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
}

func TestStartServer_PanicsOnListenError(t *testing.T) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Skip("port :8080 not available for test")
	}
	defer ln.Close()

	done := make(chan any, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- r
			} else {
				done <- "no panic"
			}
		}()
		dev.StartServer(makeTestBuilder(t))
	}()

	select {
	case v := <-done:
		if v == "no panic" {
			t.Fatalf("expected panic when ListenAndServe fails")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for StartServer panic")
	}
}
