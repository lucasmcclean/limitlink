package server

import (
	"context"
	"net/http"
	"time"

	"github.com/lucasmcclean/url-shortener/config"
	"github.com/lucasmcclean/url-shortener/repository"
)

func New(ctx context.Context, srvCfg *config.Server, repo repository.Repository) *http.Server {
	mux := http.NewServeMux()

	registerRoutes(mux)

	var handler http.Handler = mux

	// apply middleware

	server := &http.Server{
		Addr:              srvCfg.Port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	return server
}
