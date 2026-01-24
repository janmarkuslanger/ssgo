<p align="center"><img src="/logo.svg" alt="Logo" width="200" /></p>

<h1 align="center">SSGO</h1>
<p align="center">SSGO is a minimal static site generator written in Go. It is designed for <strong>clarity, explicit APIs, and flexibility</strong>.</p>

<p align="center">
  <a href="https://codecov.io/gh/janmarkuslanger/ssgo"><img src="https://codecov.io/gh/janmarkuslanger/ssgo/graph/badge.svg?token=XUZ7Y1VN3T" alt="Code coverage"></a>
  <a href="https://goreportcard.com/report/github.com/janmarkuslanger/ssgo"><img src="https://goreportcard.com/badge/github.com/janmarkuslanger/ssgo" alt="Go Report"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/releases"><img src="https://img.shields.io/github/release/janmarkuslanger/ssgo" alt="Latest Release"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/actions"><img src="https://github.com/janmarkuslanger/ssgo/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
  <a href="https://github.com/janmarkuslanger/ssgo/archive/refs/heads/main.zip"><img src="https://img.shields.io/badge/Download-ZIP-blue" alt="Download ZIP"></a>
</p>

---

## Features

- **Explicit API**: You control paths, data, templates, and output.  
- **Pluggable renderers**: Ships with `rendering.HTMLRenderer` (`html/template`); you can implement your own.  
- **Flexible writers**: Implement the `Writer` interface (disk, memory, S3, etc.).  
- **Tasks**: Run hooks before/after the build (asset copying, cleanup, etc.).  
- **Dev server**: Serve your build locally during development.  

---

## Quickstart

### 1) Install

```bash
go get github.com/janmarkuslanger/ssgo@latest
```

### 2) Create templates

```bash
mkdir -p templates
cat > templates/layout.html <<'EOF'
{{define "root"}}<html><body>{{template "content" .}}</body></html>{{end}}
EOF
cat > templates/page.html <<'EOF'
{{define "content"}}<h1>{{.Title}}</h1><p>{{.Content}}</p>{{end}}
EOF
```

### 3) Create `main.go`

```go
package main

import (
	"flag"
	"log"

	"github.com/janmarkuslanger/ssgo/builder"
	"github.com/janmarkuslanger/ssgo/dev"
	"github.com/janmarkuslanger/ssgo/page"
	"github.com/janmarkuslanger/ssgo/rendering"
	"github.com/janmarkuslanger/ssgo/writer"
)

func main() {
	devMode := flag.Bool("dev", false, "start dev server")
	flag.Parse()

	renderer := rendering.HTMLRenderer{
		Layout: []string{"templates/layout.html"},
	}

	blog := page.Generator{
		Config: page.Config{
			Template: "templates/page.html",
			Pattern:  "/blog/:slug",
			GetPaths: func() []string { return []string{"/blog/hello", "/blog/ssgo"} },
			GetData: func(p page.PagePayload) map[string]any {
				return map[string]any{
					"Title":   p.Params["slug"],
					"Content": "Hello from " + p.Params["slug"],
				}
			},
			Renderer: renderer,
		},
	}

	b := builder.Builder{
		OutputDir:  "dist",
		Writer:     writer.NewFileWriter(),
		Generators: []page.Generator{blog},
	}

	if *devMode {
		dev.StartServer(b)
		return
	}

	if err := b.Build(); err != nil {
		log.Fatal(err)
	}
}
```

### 4) Run

```bash
go run . --dev
# open http://localhost:8080/blog/hello
```

```bash
go run .
# outputs dist/blog/hello.html and dist/blog/ssgo.html
```

---

## Core Concepts & API

### Builder

The `builder.Builder` orchestrates everything.

```go
type Builder struct {
    OutputDir   string
    Generators  []page.Generator
    Writer      writer.Writer
    Renderer    rendering.Renderer
    BeforeTasks []task.Task
    AfterTasks  []task.Task
}

func (b Builder) RunTasks(tasks []task.Task) error
func (b Builder) Build() error
```

- **`OutputDir`** – where generated files go.  
- **`Writer`** – implements how files are written (e.g. `writer.NewFileWriter()`).  
- **`Generators`** – list of page generators.  
- **`Renderer`** – currently unused by `Builder`; renderers are set per generator.  
- **`BeforeTasks` / `AfterTasks`** – tasks to run before/after the build.  
- **`RunTasks(tasks)`** – runs a task list and stops on critical failures.  
- **`Build()`** – executes the full build.  

---

### Pages & Generators

A generator is the unit the builder executes. It holds the config for how pages are created.

#### Generator

```go
type Generator struct {
    Config Config
}

func (g Generator) GeneratePageInstance(path string) Page
func (g Generator) GeneratePageInstances() ([]Page, error)
```

- **`GeneratePageInstances()`** – uses `GetPaths()` and errors if it is nil.  
- **`GeneratePageInstance(path)`** – extracts params via `Pattern` and calls `GetData` if set.  

#### Config

```go
type Config struct {
    Template string
    Pattern  string
    GetPaths func() []string
    GetData  func(PagePayload) map[string]any
    MaxWorkers int
    Renderer rendering.Renderer
}

type PagePayload struct {
    Path   string
    Params map[string]string
}
```

- **`Template`** – path to the template file.  
- **`Pattern`** – route pattern used for param extraction only (supports params, e.g. `/blog/:slug`).  
- **`GetPaths()`** – returns all paths to generate (required for `GeneratePageInstances`).  
- **`GetPaths()` values** – normalized before build output; leading `/` is accepted and stripped, traversal paths are rejected, and empty/root (`/`) paths are invalid for build output.  
- **`GetData(payload)`** – returns data for each path.  
- **`MaxWorkers`** – max parallel page generation; values <= 1 run sequentially, values > 1 run concurrently; **order is always preserved regardless of the value**, but for values > 1 `GetData` must be concurrency-safe.  
- **`Renderer`** – responsible for rendering (must be set, e.g. `rendering.HTMLRenderer`).  

#### Page

```go
type Page struct {
    Path     string
    Params   map[string]string
    Data     map[string]any
    Template string
    Renderer rendering.Renderer
}

func (p Page) Render() (string, error)
```

- **`Render()`** – errors if no renderer is set and renders with `Template` + `Data`.  

#### Path helpers

```go
func ExtractParams(pattern, path string) map[string]string
func BuildPath(pattern string, params map[string]string) (string, error)
func NormalizePagePath(path string) (string, error)
func NormalizeRoutePath(path string) (string, error)
```

- **`BuildPath`** – returns an error if a required param is missing.  
- **`ExtractParams`** – does not validate segment counts; ensure pattern and path match.  
- **`NormalizePagePath`** – strips leading `/`, cleans path, rejects traversal, and errors on empty/root paths.  
- **`NormalizeRoutePath`** – normalizes to a leading `/` and rejects traversal; `/` is allowed.  

---

### Rendering

Abstracted by the `rendering.Renderer` interface:

```go
type RenderContext struct {
    Data     map[string]any
    Template string
}
```

```go
type Renderer interface {
    Render(RenderContext) (string, error)
}
```

#### HTMLRenderer

```go
type HTMLRenderer struct {
    CustomFuncs template.FuncMap
    Layout      []string
}
```

- **Layouts** – must define `{{ define "root" }}`.  
- **Content templates** – must define `{{ define "content" }}`.  
- **CustomFuncs** – inject helper functions.  
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

func NewFileWriter() *FileWriter
func (w *FileWriter) Write(path, content string) error
```

Writes files to disk (mkdir + write) and appends `.html` when missing.  

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
func NewCopyTask(sourceDir, outputSubDir string, resolver PathResolver) *CopyTask
```

Note: it returns a `*CopyTask`, which implements `task.Task`.  

---

### Dev Server

Run a simple dev server for local development.
Tasks are also executed on each page request.

```go
b := builder.Builder{...}
dev.StartServer(b)
```

`dev.NewServer` returns an `http.Handler`; `dev.StartServer` listens on `:8080`.

---

## Example

A minimal blog generator:

```go
package main

import (
    "html/template"
    "github.com/janmarkuslanger/ssgo/builder"
    "github.com/janmarkuslanger/ssgo/page"
    "github.com/janmarkuslanger/ssgo/rendering"
    "github.com/janmarkuslanger/ssgo/task"
    "github.com/janmarkuslanger/ssgo/taskutil"
    "github.com/janmarkuslanger/ssgo/writer"
    "strings"
)

var posts = map[string]map[string]any{
    "hello-world": {"title": "Hello World", "content": "Welcome to my blog!"},
    "second-post": {"title": "Second Post", "content": "More content here..."},
}

func main() {
    gen := page.Generator{
        Config: page.Config{
            Template: "templates/blog.html",
            Pattern:  "blog/:slug",
            GetPaths: func() []string {
                return []string{"blog/hello-world", "blog/second-post"}
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

    copyTask := taskutil.NewCopyTask("static", "assets", nil)

    b := builder.Builder{
        OutputDir:  "public",
        Writer:     writer.NewFileWriter(),
        Generators: []page.Generator{gen},
        BeforeTasks: []task.Task{
            copyTask,
        },
    }

    if err := b.Build(); err != nil {
        panic(err)
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

## Output Example

```
public/
  blog/
    hello-world.html
    second-post.html
  assets/
    style.css
```

---

## License

MIT © [Jan Markus Langer](https://github.com/janmarkuslanger)

## Showcases

Websites built with SSGO:

- https://www.yoga-by-julia.de/
