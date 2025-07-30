package page_test

import (
	"testing"

	"github.com/janmarkuslanger/ssgo/page"
)

func TestExtractPattern(t *testing.T) {
	got := page.ExtractParams("/hello/:id", "/hello/123")
	expected := map[string]string{"id": "123"}

	if len(got) != len(expected) {
		t.Fatalf("unexpected number of params: got %d, want %d", len(got), len(expected))
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("unexpected param value for key %q: got %q, want %q", k, got[k], v)
		}
	}

	got = page.ExtractParams("/:foo/:id", "/hello/123")
	expected = map[string]string{"foo": "hello", "id": "123"}

	if len(got) != len(expected) {
		t.Fatalf("unexpected number of params: got %d, want %d", len(got), len(expected))
	}
	for k, v := range expected {
		if got[k] != v {
			t.Errorf("unexpected param value for key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestBuildPath_success(t *testing.T) {
	got, err := page.BuildPath("/wow/:hello/:id/test", map[string]string{
		"id":    "123",
		"hello": "world",
	})
	expected := "/wow/world/123/test"

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expected {
		t.Errorf("unexpected path: got %q, want %q", got, expected)
	}
}

func TestBuildPath_missingdata(t *testing.T) {
	got, err := page.BuildPath("/hello/:id", map[string]string{})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
	expectedErr := "could not replace url param: :id"
	if err.Error() != expectedErr {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), expectedErr)
	}
	if got != "" {
		t.Errorf("expected empty path, got %q", got)
	}
}

func TestBuildPath_noparamneeded(t *testing.T) {
	got, err := page.BuildPath("/hello/world", map[string]string{})
	if err != nil {
		t.Fatal("expected no error but got one")
	}

	if got != "/hello/world" {
		t.Error("expected a valid path")
	}
}
