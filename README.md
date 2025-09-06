<p align="center"><img src="/logo.svg" alt="Logo" width="200" /></p>

<h1 align="center">SSGO</h1>
<p align="center">SSGO is a minimal **static site generator** written in Go. It is designed for **clarity, explicit APIs, and flexibility**.</p>

<p align="center">
  <a href="https://codecov.io/gh/janmarkuslanger/ssgo"><img src="https://codecov.io/gh/janmarkuslanger/ssgo/graph/badge.svg?token=XUZ7Y1VN3T" alt="Code coverage"></a>
  <a href="https://goreportcard.com/report/github.com/janmarkuslanger/ssgo"><img src="https://goreportcard.com/badge/github.com/janmarkuslanger/ssgo" alt="Go Report"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/releases"><img src="https://img.shields.io/github/release/janmarkuslanger/ssgo" alt="Latest Release"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/actions"><img src="https://github.com/janmarkuslanger/ssgo/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/archive/refs/heads/main.zip"><img src="https://img.shields.io/badge/Download-ZIP-blue" alt="Download ZIP"></a>
</p>

---

## ✨ Features

- **Explicit API**: You control paths, data, templates, and output.  
- **Pluggable renderers**: Default is `html/template`, but you can implement your own.  
- **Flexible writers**: Write to disk, memory, S3, or anywhere.  
- **Tasks**: Run hooks before/after the build (asset copying, cleanup, etc.).  
- **Dev server**: Serve your build locally during development.  

---

## 📦 Installation

```bash
go get github.com/janmarkuslanger/ssgo@latest
```

---

## 🔑 Core Concepts & API

### Builder

The `builder.Builder` orchestrates everything.

```go
type Builder struct {
    OutputDir   string
    Writer      writer.Writer
    Generators  []page.Generator
    BeforeTasks []task.Task
    AfterTasks  []task.Task
}

func (b *Builder) Build() error
```

- **`OutputDir`** – where generated files go.  
- **`Writer`** – implements how files are written (default: `writer.FileWriter`).  
- **`Generators`** – list of page generators.  
- **`BeforeTasks` / `AfterTasks`** – tasks to run before/after the build.  
- **`Build()`** – executes the full build.  

---

### Pages & Generators

Define what content to generate with a `page.Config`.  
A generator combines a pattern, template, and data.

```go
type Config struct {
    Template string
    Pattern  string
    GetPaths func() []string
    GetData  func(PagePayload) map[string]any
    Renderer rendering.Renderer
}

type PagePayload struct {
    Path   string
    Params map[string]string
}
```

- **`Template`** – path to the template file.  
- **`Pattern`** – route pattern (supports params, e.g. `/blog/:slug`).  
- **`GetPaths()`** – returns all paths to generate.  
- **`GetData(payload)`** – returns data for each path.  
- **`Renderer`** – responsible for rendering (default: `HTMLRenderer`).  

---

### Rendering

Abstracted by the `rendering.Renderer` interface:

```go
type Renderer interface {
    Render(RenderContext) (string, error)
}
```

#### HTMLRenderer (default)

```go
type HTMLRenderer struct {
    Layout      []string
    CustomFuncs template.FuncMap
    ExtraData   map[string]any
}
```

- **Layouts** – must define `{{ define "root" }}`.  
- **Content templates** – must define `{{ define "content" }}`.  
- **CustomFuncs** – inject helper functions.  
- **ExtraData** – pass additional values to all templates.  

---

### Writer

Defines how output is written.

```go
type Writer interface {
    Write(path string, content string) error
}
```

Default implementation:

```go
type FileWriter struct{}

func (FileWriter) Write(path, content string) error
```

Writes files to disk (mkdir + write).  

---

### Tasks

Hook into the build with before/after tasks.

```go
type Task interface {
    Run(ctx TaskContext) error
    IsCritical() bool
}

type TaskContext struct {
    OutputDir string
}
```

- **Critical tasks** – stop the build on failure.  
- **Non-critical tasks** – log and continue.  

#### CopyTask (built-in)

Copy static assets into the build output.

```go
func NewCopyTask(sourceDir, outputSubDir string, resolver PathResolver) CopyTask
```

---

### Dev Server

Run a simple dev server for local development.
Tasks will ne also executed on each "refresh".

```go
b := builder.Builder{...}
dev.StartServer(b)
```

---

## 🚀 Example

A minimal blog generator:

```go
package main

import (
    "html/template"
    "github.com/janmarkuslanger/ssgo/builder"
    "github.com/janmarkuslanger/ssgo/dev"
    "github.com/janmarkuslanger/ssgo/page"
    "github.com/janmarkuslanger/ssgo/rendering"
    "github.com/janmarkuslanger/ssgo/task"
    "github.com/janmarkuslanger/ssgo/taskutil"
    "github.com/janmarkuslanger/ssgo/writer"
    "strings"
	"os"
)

var posts = map[string]map[string]any{
    "hello-world": {"title": "Hello World", "content": "Welcome to my blog!"},
    "second-post": {"title": "Second Post", "content": "More content here..."},
}

func main() {
    gen := page.Generator{
        Config: page.Config{
            Template: "templates/blog.html",
            Pattern:  "/blog/:slug",
            GetPaths: func() []string {
                return []string{"/blog/hello-world", "/blog/second-post"}
            },
            GetData: func(p page.PagePayload) map[string]any {
                return posts[p.Params["slug"]]
            },
            Renderer: rendering.HTMLRenderer{
                Layout: []string{"templates/layout.html"},
                CustomFuncs: template.FuncMap{
                    "upper": strings.ToUpper,
                },
            },
        },
    }

    b := builder.Builder{
        OutputDir:  "public",
        Writer:     writer.FileWriter{},
        Generators: []page.Generator{gen},
        BeforeTasks: []task.Task{
            taskutil.NewCopyTask("static", "assets", nil),
        },
    }

    if err := b.Build(); err != nil {
        panic(err)
    }

	if os.Getenv("ENV") == "test" {
		dev.StartServer(b)
	}
}
```

Folder structure:

```
templates/
  layout.html
  blog.html
static/
  style.css
```

---

## 🗂️ Output Example

```
public/
  blog/
    hello-world
    second-post
  assets/
    style.css
```

---

✅ With SSGO you control **exactly what gets built** — no hidden magic, just Go code.

## 📖 License

MIT © [Jan Markus Langer](https://github.com/janmarkuslanger)

## Showcases

Websites built with SSGO:

- https://www.yoga-by-julia.de/
