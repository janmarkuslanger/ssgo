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

func TestBuildPath_success(t *testing.T) {
	p, err := page.BuildPath("/wow/:hello/:id/test", map[string]string{
		"id":    "123",
		"hello": "world",
	})
	assert.Equal(t, p, "/wow/world/123/test")
	assert.Equal(t, err, nil)
}

func TestBuildPath_missingdata(t *testing.T) {
	p, err := page.BuildPath("/hello/:id", map[string]string{})
	assert.Equal(t, p, "")
	assert.EqualError(t, err, "could not replace url param: :id")
}
