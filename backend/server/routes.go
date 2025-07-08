package server

import (
	"net/http"

	"github.com/lucasmcclean/limitlink/handler"
	"github.com/lucasmcclean/limitlink/link"
)

func registerRoutes(mux *http.ServeMux, repo link.Repository) {
	mux.Handle("/link", handler.Link(repo))
}
