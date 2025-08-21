package server

import (
	"net/http"

	"github.com/lucasmcclean/limitlink/link"
)

func registerRoutes(mux *http.ServeMux, repo link.Repository) {
	mux.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			LinkHandler(repo)(w, r)
			return
		}
		http.NotFound(w, r)
	})
	mux.Handle("/links/", LinkHandler(repo))
	mux.Handle("/", RedirectHandler(repo))
}
