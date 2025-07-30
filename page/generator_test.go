package page_test

import (
	"strconv"
	"testing"

	"github.com/janmarkuslanger/ssgo/page"
)

func TestGeneratorGeneratePages_missinggetpaths(t *testing.T) {
	g := page.Generator{}
	p, err := g.GeneratePageInstances()

	if err == nil {
		t.Fatal("expected an error but got nil")
	}

	expectedErr := "GetPaths is not defined in Config"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), expectedErr)
	}

	if len(p) != 0 {
		t.Error("pages should be empty")
	}
}

func TestGeneratorGeneratePages_simple(t *testing.T) {
	c := page.Config{
		GetPaths: func() []string {
			return []string{
				"hello/world",
				"foo/bar",
			}
		},
	}
	g := page.Generator{
		Config: c,
	}
	p, err := g.GeneratePageInstances()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(p) != 2 {
		t.Errorf("expected 2 pages, got %d", len(p))
	}

	if p[0].Path != "hello/world" {
		t.Errorf("page should not be: %q", p[0].Path)
	}

	if p[1].Path != "foo/bar" {
		t.Errorf("page should not be: %q", p[1].Path)
	}
}

func TestGeneratorGeneratePages_withdata(t *testing.T) {
	c := page.Config{
		Pattern: "/add/:number",
		GetPaths: func() []string {
			return []string{
				"/add/1",
				"/add/2",
				"/add/666",
			}
		},
		GetData: func(payload page.PagePayload) map[string]any {
			n, ok := payload.Params["number"]

			if !ok {
				return nil
			}

			num, err := strconv.Atoi(n)
			if err != nil {
				return nil
			}

			return map[string]any{
				"newnumber": num + 1,
			}
		},
	}
	g := page.Generator{
		Config: c,
	}
	p, err := g.GeneratePageInstances()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(p) != 3 {
		t.Errorf("expected 3 pages, got %d", len(p))
	}

	if p[0].Path != "/add/1" {
		t.Errorf("page should not be: %q", p[0].Path)
	}

	if n := p[0].Data["newnumber"]; n != 2 {
		t.Errorf("unexpected newnumber: got %q, want %q", n, 2)
	}

	if p[2].Path != "/add/666" {
		t.Errorf("page should not be: %q", p[2].Path)
	}

	if n := p[2].Data["newnumber"]; n != 667 {
		t.Errorf("unexpected newnumber: got %q, want %q", n, 667)
	}
}
