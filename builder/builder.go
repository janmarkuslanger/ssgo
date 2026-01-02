package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
	"github.com/janmarkuslanger/ssgo/task"
	"github.com/janmarkuslanger/ssgo/writer"
)

type Builder struct {
	OutputDir   string
	Generators  []page.Generator
	Writer      writer.Writer
	Renderer    rendering.Renderer
	BeforeTasks []task.Task
	AfterTasks  []task.Task
}

func (b Builder) RunTasks(tasks []task.Task) error {
	for _, t := range tasks {
		err := t.Run(task.TaskContext{
			OutputDir: b.OutputDir,
		})

		if err != nil && t.IsCritical() {
			return fmt.Errorf("failed to run tasks: %w", err)
		}

		if err != nil {
			fmt.Printf("warning: task failed to run: %v", err)
		}
	}

	return nil
}

func (b Builder) Build() error {
	if err := b.RunTasks(b.BeforeTasks); err != nil {
		return err
	}

	for _, g := range b.Generators {
		pages, err := g.GeneratePageInstances()
		if err != nil {
			return fmt.Errorf("failed to generate pages: %w", err)
		}

		for _, p := range pages {
			cleanPath := filepath.Clean(p.Path)
			if cleanPath == "." {
				return fmt.Errorf("page path must not be empty")
			}
			if filepath.IsAbs(cleanPath) {
				return fmt.Errorf("page path must be relative, got %q", p.Path)
			}
			if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
				return fmt.Errorf("page path must not traverse outside output dir: %q", p.Path)
			}

			content, err := p.Render()
			if err != nil {
				// TODO: make configurable if it should continue if single page fails
				return fmt.Errorf("failed to render page %s: %w", p.Path, err)
			}

			fullPath := filepath.Join(b.OutputDir, cleanPath)
			if err := b.Writer.Write(fullPath, content); err != nil {
				return fmt.Errorf("failed to write page %s: %w", p.Path, err)
			}
		}
	}

	if err := b.RunTasks(b.AfterTasks); err != nil {
		return err
	}

	return nil
}
