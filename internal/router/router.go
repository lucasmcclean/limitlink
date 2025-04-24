package router

import (
	"github.com/lucasmcclean/url-shortener/internal/link"
	"net/http"
)

// New creates and returns an http.Handler for the application with the routes configured.
func New(repo link.LinkRepository) http.Handler {
	mux := http.NewServeMux()

	h := link.NewHandler(link.NewService(repo), nil)

	registerRoutes(mux, h)

	return mux
}

// registerRoutes registers all routes with the provided mux
func registerRoutes(mux *http.ServeMux, h *link.Handler) {
	staticFS := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	mux.HandleFunc("/links", h.Create)
	mux.HandleFunc("/v/", h.Visit)
	mux.HandleFunc("/admin/", h.AdminView)
	mux.HandleFunc("/delete/", h.Delete)
}
