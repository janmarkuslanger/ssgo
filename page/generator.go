package page

import (
	"errors"
	"fmt"
	"strings"
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
	// MaxWorkers controls parallel page generation. Values <= 1 run sequentially.
	MaxWorkers int
	Renderer   rendering.Renderer
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
	pages := make([]Page, len(paths))
	if len(paths) == 0 {
		return pages, nil
	}
	if err := g.validatePaths(paths); err != nil {
		return nil, err
	}

	workers := g.Config.MaxWorkers
	if workers <= 1 {
		for i, path := range paths {
			pages[i] = g.GeneratePageInstance(path)
		}
		return pages, nil
	}
	if len(paths) < workers {
		workers = len(paths)
	}

	type job struct {
		index int
		path  string
	}

	jobs := make(chan job)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				pages[j.index] = g.GeneratePageInstance(j.path)
			}
		}()
	}

	for i, path := range paths {
		jobs <- job{index: i, path: path}
	}
	close(jobs)
	wg.Wait()

	return pages, nil
}

func (g Generator) validatePaths(paths []string) error {
	for i, path := range paths {
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("path at index %d is empty", i)
		}
		if err := validatePatternMatch(g.Config.Pattern, path); err != nil {
			return fmt.Errorf("path at index %d: %w", i, err)
		}
	}
	return nil
}
