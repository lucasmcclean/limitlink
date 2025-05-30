package server

import (
	"context"
	"net/http"

	"github.com/lucasmcclean/limitlink/handlers"
	"github.com/lucasmcclean/limitlink/link"
)

func registerRoutes(mux *http.ServeMux, ctx context.Context, repo link.Repository) {
	mux.Handle("/{$}", handlers.Index())
	mux.Handle("/links", handlers.Links(ctx, repo))
	mux.Handle("/static/", http.StripPrefix("/static/", handlers.Static()))
	mux.Handle("/admin/links/", handlers.AdminLinks(ctx, repo))
	// mux.Handle("/", handlers.Redirect())
}
