package builder

import (
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/writer"
)

type Builder struct {
	OutputDir string
	Pages     []page.Config
	Writer    writer.Writer
}

func (b *Builder) Build() {

}
