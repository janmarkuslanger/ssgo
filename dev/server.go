package dev

import (
	"net/http"
	"path/filepath"

	"github.com/janmarkuslanger/ssgo/builder"
)

func NewServer(builder builder.Builder) http.Handler {
	mux := http.NewServeMux()

	pagePaths := make(map[string]struct{})
	for _, g := range builder.Pages {
		for _, path := range g.Config.GetPaths() {
			pagePaths[path] = struct{}{}
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				p := g.GeneratePageInstance(path)
				c, err := p.Render()
				if err != nil {
					panic(err)
				}

				w.Write([]byte(c))
			})
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := pagePaths[r.URL.Path]; ok {
			mux.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, filepath.Join(builder.OutputDir, filepath.Clean(r.URL.Path)))
	})
}

func StartServer(builder builder.Builder) {
	mux := NewServer(builder)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
