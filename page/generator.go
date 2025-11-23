package page

import (
	"errors"
	"sync"

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
	if g.Config.GetPaths == nil {
		return nil, errors.New("GetPaths is not defined in Config")
	}

	paths := g.Config.GetPaths()
	if len(paths) == 0 {
		return nil, nil
	}

	const maxWorkers = 8
	workers := maxWorkers
	if len(paths) < workers {
		workers = len(paths)
	}

	jobs := make(chan string)
	pages := make([]Page, 0, len(paths))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				p := g.GeneratePageInstance(path)

				mu.Lock()
				pages = append(pages, p)
				mu.Unlock()
			}
		}()
	}

	for _, path := range paths {
		jobs <- path
	}
	close(jobs)
	wg.Wait()

	return pages, nil
}
