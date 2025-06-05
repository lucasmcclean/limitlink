package server

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/lucasmcclean/limitlink/handlers"
	"github.com/lucasmcclean/limitlink/link"
)

func registerRoutes(
	mux *http.ServeMux,
	ctx context.Context,
	repo link.Repository,
	staticFS fs.FS,
	templatesFS fs.FS,
) {
	mux.Handle("/{$}", handlers.Index(staticFS))
	mux.Handle("/links", handlers.Links(ctx, repo, templatesFS))
	mux.Handle("/static/", http.StripPrefix("/static/", handlers.Static(staticFS)))
	mux.Handle("/admin/", handlers.AdminLinks(ctx, repo, templatesFS))
	mux.Handle("/", handlers.Redirect(ctx, repo, templatesFS))
}
