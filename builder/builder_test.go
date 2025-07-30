package builder_test

import (
	"errors"
	"testing"

	"github.com/janmarkuslanger/ssgo/builder"
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
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

func TestBuilder_Build_success(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Pages: []page.Generator{
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

func TestBuilder_Build_failingGeneratePageInstances(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Pages: []page.Generator{
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

func TestBuilder_Build_failingwriter(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriterFail{},
		Pages: []page.Generator{
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

func TestBuilder_Build_failingrenderer(t *testing.T) {
	b := builder.Builder{
		OutputDir: "/test",
		Writer:    MockWriter{},
		Pages: []page.Generator{
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
