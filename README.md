<p align="center"><img src="/logo.svg" alt="Logo" width="200" /></p>
<h1 align="center">SSGO</h1>
<p align="center">Simple, fast and extendable static site generator.</p>

<p align="center">
  <a href="https://codecov.io/gh/janmarkuslanger/ssgo"><img src="https://codecov.io/gh/janmarkuslanger/ssgo/graph/badge.svg?token=XUZ7Y1VN3T" alt="Code coverage"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/releases"><img src="https://img.shields.io/github/release/janmarkuslanger/ssgo" alt="Latest Release"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/actions"><img src="https://github.com/janmarkuslanger/ssgo/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/archive/refs/heads/main.zip"><img src="https://img.shields.io/badge/Download-ZIP-blue" alt="Download ZIP"></a>
</p>

---

## ✨ What is SSGO?

**SSGO** is a minimal, Go-native static site generator focused on:

- 🧩 *Composability* – every page is generated via a clear data + template pipeline
- ⚡ *Simplicity* – no magic or custom config formats
- 🔧 *Extensibility* – plug in your own rendering or writing logic

---

## 📦 Installation

Install via:

```bash
go get github.com/janmarkuslanger/ssgo
```

Or download directly:  
[⬇ Download latest ZIP](https://github.com/janmarkuslanger/ssgo/archive/refs/heads/main.zip)

---

## 🚀 Minimal Example

In this example we use the default implementation of the renderer (https://pkg.go.dev/html/template) and the writer. 

Create the following structure:

```
your-project/
├── main.go
├── templates/
│   ├── layout.html
│   └── blog.html
└── public/ (generated after running)
```

### main.go

```go
package main

import (
	"log"

	"github.com/janmarkuslanger/ssgo/builder"
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
)

func main() {
	renderer := rendering.HTMLRenderer{
		Layout: []string{"templates/layout.html"},
	}

	posts := map[string]map[string]any{
		"hello-world": {
			"Title":   "Hello World",
			"Content": "Welcome to my blog!",
		},
		"second-post": {
			"Title":   "Second Post",
			"Content": "Another blog entry.",
		},
	}

	generator := page.Generator{
		Config: page.Config{
			Pattern:  "/blog/:slug",
			Template: "templates/blog.html",
			GetPaths: func() []string {
				return []string{"/blog/hello-world", "/blog/second-post"}
			},
			GetData: func(p page.PagePayload) map[string]any {
				return posts[p.Params["slug"]]
			},
			Renderer: renderer,
		},
	}

	b := builder.Builder{
		OutputDir: "public",
		Writer:    builder.FileWriter{},
		Generators: []page.Generator{
			generator,
		},
	}

	if err := b.Build(); err != nil {
		log.Fatal(err)
	}
}
```

### templates/layout.html

```html
<!DOCTYPE html>
<html>
  <head><title>{{ .Title }}</title></head>
  <body>
    {{ template "content" . }}
  </body>
</html>
```

### templates/blog.html

```html
{{ define "content" }}
<h1>{{ .Title }}</h1>
<p>{{ .Content }}</p>
{{ end }}
```

### Run it

```bash
go run main.go
```

→ Two files will be generated:

```
public/
└── blog/
    ├── hello-world
    └── second-post
```

---

## 🧱 Concepts

### 🔨 Builder

Orchestrates the generation of pages and writes them to disk:

```go
type Builder struct {
	OutputDir  string
	Writer     Writer
	Generators []page.Generator
}
```

### 📄 Generator / Config

Each `page.Generator` is driven by a `Config`:

```go
type Config struct {
	Template string
	Pattern  string
	GetPaths func() []string
	GetData  func(PagePayload) map[string]any
	Renderer rendering.Renderer
}
```

This allows dynamic paths with params like `/blog/:slug`.

### 📦 PagePayload

Passed to `GetData` so you can access dynamic URL parameters:

```go
type PagePayload struct {
	Path   string
	Params map[string]string
}
```

---

## 🖋 Writer

The `Writer` interface is used to persist rendered output:

```go
type Writer interface {
	Write(path string, content string) error
}
```

Default implementation:

```go
type FileWriter struct{}

func (FileWriter) Write(path string, content string) error {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	return os.WriteFile(path, []byte(content), 0644)
}
```

Swap this out to write to memory, S3, etc.

---

## 🖌️ Rendering

Rendering is abstracted via this interface:

```go
type Renderer interface {
	Render(RenderContext) (string, error)
}
```

The built-in `HTMLRenderer` supports:

- Go templates (`html/template`)
- Layouts (via `[]string`)
- Custom data injection

---

## 📖 License

MIT © [Jan Markus Langer](https://github.com/janmarkuslanger)
