package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lucasmcclean/url-shortener/config"
	"github.com/lucasmcclean/url-shortener/postgres"
	"github.com/lucasmcclean/url-shortener/server"
)

// FIX: Setup logging package

func main() {
	ctx := context.Background()
	ctx, cancelCtx := signal.NotifyContext(ctx, os.Interrupt)

	dbCfg := config.GetDB()
	srvCfg := config.GetServer()

	repo, err := postgres.New(dbCfg)
	if err != nil {
		// log.Printf("error connecting to database: %s", err)
		os.Exit(1)
	}

	server := server.New(ctx, srvCfg, repo)
	go func() {
		// log.Printf("listening and serving on: %s", server.Addr)
		// TODO: Get certs for TLS
		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil && err != http.ErrServerClosed {
			// log.Printf("error listening and serving: %s", err)
			cancelCtx() // Server has failed so begin shutdown
		}
	}()

	// TODO: Add a server to redirect http requests to https

	var wg sync.WaitGroup
	wg.Add(1)
	successfulShutdown := true
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancelCtx := context.WithTimeout(shutdownCtx, time.Second*10)
		defer cancelCtx()

		err = server.Shutdown(shutdownCtx)
		if err != nil {
			// log.Printf("error shutting down HTTP server: %s", err)
			successfulShutdown = false
		}

		err = repo.Close()
		if err != nil {
			// log.Printf("error closing repository: %s", err)
			successfulShutdown = false
		}
	}()
	wg.Wait()

	if !successfulShutdown {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
