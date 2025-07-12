package server

import (
	"net/http"

	"github.com/lucasmcclean/limitlink/link"
)

func registerRoutes(mux *http.ServeMux, repo link.Repository) {
	mux.Handle("/links", LinkHandler(repo))
	mux.Handle("/", RedirectHandler(repo))
}
