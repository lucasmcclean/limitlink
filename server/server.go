package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/lucasmcclean/url-shortener/config"
	"github.com/lucasmcclean/url-shortener/repository"
)

func New(ctx context.Context, srvCfg *config.Server, repo repository.Repository) *http.Server {
	mux := http.NewServeMux()

	addRoutes(mux)

	var handler http.Handler = mux

	// apply middleware

	// TODO: Use config
	server := &http.Server{
		Addr:         ":443",
		Handler:      handler,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	return server
}

// TODO: Setup tls
func getTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{}
	return tlsConfig
}
