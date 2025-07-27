package page_test

import (
	"testing"

	"github.com/janmarkuslanger/ssgo/page"
	"github.com/stretchr/testify/assert"
)

func TestExtractPattern(t *testing.T) {
	assert.Equal(t, page.ExtractPattern("/hello/:id", "/hello/123"), map[string]string{
		"id": "123",
	})

	assert.Equal(t, page.ExtractPattern("/:foo/:id", "/hello/123"), map[string]string{
		"foo": "hello",
		"id":  "123",
	})
}
