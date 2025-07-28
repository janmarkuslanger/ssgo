package page_test

import (
	"testing"

	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
)

type MockRenderer struct{}

func (r MockRenderer) Render(ctx rendering.RenderContext) (string, error) {
	return "hello world", nil
}

func TestPage_Render_norenderer(t *testing.T) {
	p := page.Page{}
	_, err := p.Render()

	if err.Error() != "no renderer set" {
		t.Error("should throw an error")
	}
}

func TestPage_Render_simplerenderer(t *testing.T) {
	p := page.Page{
		Renderer: MockRenderer{},
	}

	out, err := p.Render()

	if err != nil {
		t.Errorf("should not return an error but got %q", err)
	}

	if out != "hello world" {
		t.Errorf("expected other output but got %q", out)
	}

}
