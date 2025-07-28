package builder

import (
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
	"github.com/janmarkuslanger/ssgo/writer"
)

type Builder struct {
	OutputDir string
	Pages     []page.Generator
	Writer    writer.Writer
	Renderer  rendering.Renderer
}
