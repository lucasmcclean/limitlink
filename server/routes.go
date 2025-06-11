package server

import (
	"net/http"

	"github.com/lucasmcclean/limitlink/handlers"
	"github.com/lucasmcclean/limitlink/link"
)

func registerRoutes(mux *http.ServeMux, repo link.Repository) {
	mux.Handle("/{$}", handlers.Index())
	mux.Handle("/links", handlers.Links(repo))
	mux.Handle("/static/", http.StripPrefix("/static/", handlers.Static()))
	mux.Handle("/admin/", handlers.Admin(repo))
	mux.Handle("/", handlers.Redirect(repo))
}
