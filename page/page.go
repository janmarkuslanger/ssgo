package page

import (
	"errors"

	"github.com/janmarkuslanger/ssgo/rendering"
)

type Page struct {
	Path     string
	Params   map[string]string
	Data     map[string]any
	Template string
	Renderer rendering.Renderer
}

func (p Page) Render() (string, error) {
	if p.Renderer == nil {
		return "", errors.New("no renderer set")
	}

	return p.Renderer.Render(rendering.RenderContext{
		Data:     p.Data,
		Template: p.Template,
	})
}
