package builder

import "github.com/janmarkuslanger/ssgo/page"

type Builder struct {
	OutputDir string
	Pages     []page.Config
}
