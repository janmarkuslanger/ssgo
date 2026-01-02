package dev

import (
	"net/http"

	"github.com/janmarkuslanger/ssgo/builder"
)

func NewServer(builder builder.Builder) http.Handler {
	mux := http.NewServeMux()

	pagePaths := make(map[string]struct{})
	for _, g := range builder.Generators {
		if g.Config.GetPaths == nil {
			continue
		}
		for _, path := range g.Config.GetPaths() {
			pagePaths[path] = struct{}{}
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				if err := builder.RunTasks(builder.BeforeTasks); err != nil {
					panic(err)
				}

				p := g.GeneratePageInstance(path)
				c, err := p.Render()

				if err != nil {
					panic(err)
				}

				if err := builder.RunTasks(builder.AfterTasks); err != nil {
					panic(err)
				}

				w.Write([]byte(c))
			})
		}
	}

	fs := http.FileServer(http.Dir(builder.OutputDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := pagePaths[r.URL.Path]; ok {
			mux.ServeHTTP(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	})
}

func StartServer(builder builder.Builder) {
	mux := NewServer(builder)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
