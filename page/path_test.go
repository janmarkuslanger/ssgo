package page_test

import (
	"testing"

	"github.com/janmarkuslanger/ssgo/page"
	"github.com/stretchr/testify/assert"
)

func TestExtractPattern(t *testing.T) {
	assert.Equal(t, page.ExtractParams("/hello/:id", "/hello/123"), map[string]string{
		"id": "123",
	})

	assert.Equal(t, page.ExtractParams("/:foo/:id", "/hello/123"), map[string]string{
		"foo": "hello",
		"id":  "123",
	})
}

func TestBuildPath(t *testing.T) {
	assert.Equal(t, page.BuildPath("/hello/:id", map[string]string{
		"id": "123",
	}), "/hello/123")

	assert.Equal(t, page.BuildPath("/:hello/:id", map[string]string{
		"id":    "123",
		"hello": "world",
	}), "/world/123")
}
