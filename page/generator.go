package page

import (
	"errors"

	"github.com/janmarkuslanger/ssgo/rendering"
)

type PagePayload struct {
	Params map[string]string
	Path   string
}

type Config struct {
	Template string
	Pattern  string
	GetData  func(payload PagePayload) map[string]any
	GetPaths func() []string
	Renderer rendering.Renderer
}

type Generator struct {
	Config Config
}

func (g Generator) GeneratePageInstance(path string) Page {
	data := make(map[string]any)
	params := ExtractParams(g.Config.Pattern, path)

	if g.Config.GetData != nil {
		data = g.Config.GetData(PagePayload{
			Path:   path,
			Params: params,
		})
	}

	return Page{
		Path:     path,
		Params:   params,
		Data:     data,
		Template: g.Config.Template,
		Renderer: g.Config.Renderer,
	}
}

func (g Generator) GeneratePageInstances() ([]Page, error) {
	pages := []Page{}

	if g.Config.GetPaths == nil {
		return pages, errors.New("GetPaths is not defined in Config")
	}

	for _, path := range g.Config.GetPaths() {
		p := g.GeneratePageInstance(path)
		pages = append(pages, p)
	}

	return pages, nil
}
