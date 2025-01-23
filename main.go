package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lucasmcclean/url-shortener/config"
	"github.com/lucasmcclean/url-shortener/logger"
	"github.com/lucasmcclean/url-shortener/repository/postgres"
	"github.com/lucasmcclean/url-shortener/server"
)

func main() {
	ctx := context.Background()
	ctx, cancelCtx := signal.NotifyContext(ctx, os.Interrupt)

	dbCfg := config.GetDB()
	srvCfg := config.GetServer()

	log := logger.NewStdLogger(logger.DEBUG, os.Stderr)

	repo, err := postgres.New(dbCfg)
	if err != nil {
		log.Error("couldn't connect to database", "error", err)
		os.Exit(1)
	}

	server, err := server.New(ctx, srvCfg, repo)
	if err != nil {
		log.Error("failed to setup HTTPS server", "error", err)
		err = repo.Close()
		if err != nil {
			log.Error("failed to close repository", "error", err)
		}
		os.Exit(1)
	}

	go func() {
		log.Info("listening and serving", "address", server.Addr)
		err := server.ListenAndServeTLS(srvCfg.CertPath+"/server.crt", srvCfg.CertPath+"/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Error("failed to listen and serve", "error", err)
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
			log.Error("failed to shutdown server", "error", err)
			successfulShutdown = false
		}

		err = repo.Close()
		if err != nil {
			log.Error("failed to close repository", "error", err)
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
