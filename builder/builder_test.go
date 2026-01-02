package builder_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/janmarkuslanger/ssgo/builder"
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
	"github.com/janmarkuslanger/ssgo/task"
)

type MockWriter struct{}

func (w MockWriter) Write(filepath string, content string) (err error) {
	return nil
}

type MockWriterFail struct{}

func (w MockWriterFail) Write(filepath string, content string) (err error) {
	return errors.New("something went wrong")
}

type MockRenderer struct{}

func (r MockRenderer) Render(ctx rendering.RenderContext) (output string, err error) {
	return "hello world", nil
}

type MockRendererFail struct{}

func (r MockRendererFail) Render(ctx rendering.RenderContext) (output string, err error) {
	return "", errors.New("something went wrong")
}

type MockTask struct{}

func (t MockTask) Run(ctx task.TaskContext) error {
	return nil
}

func (t MockTask) IsCritical() bool {
	return true
}

type MockTaskFail struct{}

func (t MockTaskFail) Run(ctx task.TaskContext) error {
	return errors.New("task fails!")
}

func (t MockTaskFail) IsCritical() bool {
	return false
}

type MockTaskFailCitical struct{}

func (t MockTaskFailCitical) Run(ctx task.TaskContext) error {
	return errors.New("task fails!")
}

func (t MockTaskFailCitical) IsCritical() bool {
	return true
}

func TestBuilder_Build_Success(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
	}

	err := b.Build()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuilder_Build_Tasks_NoErr(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
		BeforeTasks: []task.Task{
			MockTask{},
		},
		AfterTasks: []task.Task{
			MockTask{},
		},
	}

	err := b.Build()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuilder_Build_Tasks_Err(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
		BeforeTasks: []task.Task{
			MockTask{},
			MockTaskFail{},
		},
		AfterTasks: []task.Task{
			MockTask{},
			MockTaskFail{},
		},
	}

	err := b.Build()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBuilder_Build_TasksBefore_FatalErr(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
		BeforeTasks: []task.Task{
			MockTask{},
			MockTaskFail{},
			MockTaskFailCitical{},
		},
		AfterTasks: []task.Task{
			MockTask{},
			MockTaskFail{},
		},
	}

	err := b.Build()

	if err == nil {
		t.Error("expected error")
	}

	expected := "failed to run tasks:"
	if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("expected error %v but got %q", expected, err.Error())
	}
}

func TestBuilder_Build_TasksAfter_FatalErr(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
		BeforeTasks: []task.Task{
			MockTask{},
			MockTaskFail{},
		},
		AfterTasks: []task.Task{
			MockTask{},
			MockTaskFail{},
			MockTaskFailCitical{},
		},
	}

	err := b.Build()

	if err == nil {
		t.Error("expected error")
	}

	expected := "failed to run tasks:"
	if !strings.HasPrefix(err.Error(), expected) {
		t.Errorf("expected error %v but got %q", expected, err.Error())
	}
}

func TestBuilder_Build_RejectsAbsolutePath(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"/absolute"}
					},
				},
			},
		},
	}

	err := b.Build()
	if err == nil {
		t.Fatal("expected error for absolute path")
	}
	if !strings.Contains(err.Error(), "page path must be relative") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuilder_Build_RejectsTraversalPath(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"../outside"}
					},
				},
			},
		},
	}

	err := b.Build()
	if err == nil {
		t.Fatal("expected error for traversal path")
	}
	if !strings.Contains(err.Error(), "must not traverse outside output dir") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuilder_Build_RejectsEmptyPath(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{""}
					},
				},
			},
		},
	}

	err := b.Build()
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if !strings.Contains(err.Error(), "page path must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuilder_Build_FailingGeneratePageInstances(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{},
			},
		},
	}

	err := b.Build()

	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestBuilder_Build_FailingWriter(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriterFail{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRenderer{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
	}

	err := b.Build()

	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestBuilder_Build_FailingRenderer(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Generators: []page.Generator{
			{
				Config: page.Config{
					Renderer: MockRendererFail{},
					GetPaths: func() []string {
						return []string{"a", "b"}
					},
				},
			},
		},
	}

	err := b.Build()

	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}
