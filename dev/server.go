package dev

import (
	"net/http"

	"github.com/janmarkuslanger/ssgo/builder"
)

func StartServer(builder builder.Builder) {
	mux := http.NewServeMux()

	for _, g := range builder.Pages {
		for _, path := range g.Config.GetPaths() {
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				p := g.GeneratePageInstance(path)
				c, _ := p.Render()
				w.Write([]byte(c))
			})
		}
	}

	http.ListenAndServe(":8080", mux)
}
