package main

import (
	"net/http"
	"time"
)

func newServer() *http.Server {
	mux := http.NewServeMux()
	registerRoutes(mux)

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
