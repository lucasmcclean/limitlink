package server

import (
	"net/http"
	"time"

	"github.com/lucasmcclean/limitlink/link"
)

func New(repo link.Repository) *http.Server {
	mux := http.NewServeMux()

	registerRoutes(mux, repo)

	var handler = maxBodySizeMiddleware(mux)

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
