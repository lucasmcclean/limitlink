package server

import (
	"context"
	"io/fs"
	"net/http"
	"time"

	"github.com/lucasmcclean/limitlink/link"
)

func New(
	ctx context.Context,
	repo link.Repository,
	staticFS fs.FS,
	templatesFS fs.FS,
) *http.Server {
	mux := http.NewServeMux()

	registerRoutes(mux, ctx, repo, staticFS, templatesFS)

	var handler http.Handler = mux

	return &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
