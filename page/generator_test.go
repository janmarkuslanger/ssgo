package page_test

import (
	"testing"

	"github.com/janmarkuslanger/ssgo/page"
)

func TestGeneratorGeneratePages_missinggetpaths(t *testing.T) {
	g := page.Generator{}
	p, err := g.GeneratePageInstances()

	if err.Error() != "GetPaths is not defined in Config" {
		t.Error("should throw an error")
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
		t.Error("there should be an error")
	}

	if len(p) != 2 {
		t.Error("there should be 2 pages")
	}

	if p[0].Path != "hello/world" {
		t.Errorf("page should not be: %q", p[0].Path)
	}

	if p[1].Path != "foo/bar" {
		t.Errorf("page should not be: %q", p[1].Path)
	}
}
