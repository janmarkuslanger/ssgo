package builder

import (
	"fmt"
	"path/filepath"

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

func (b Builder) Build() error {
	for _, g := range b.Pages {
		pages, err := g.GeneratePageInstances()
		if err != nil {
			return fmt.Errorf("failed to generate pages: %w", err)
		}

		for _, p := range pages {
			content, err := p.Render()
			if err != nil {
				// TODO: make configurable if it should continue if single page faileds
				return fmt.Errorf("failed to render page %s: %w", p.Path, err)
			}

			fullPath := filepath.Join(b.OutputDir, p.Path)
			if err := b.Writer.Write(fullPath, content); err != nil {
				return fmt.Errorf("failed to write page %s: %w", p.Path, err)
			}
		}
	}

	return nil
}
