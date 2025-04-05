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
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	log := logger.NewStdLogger(logger.DEBUG, os.Stderr)
	config.InitApp(log)

	if config.App.IsProd() {
		log = logger.NewStdLogger(logger.INFO, os.Stderr)
	}

	dbCfg := config.GetDB(log)
	srvCfg := config.GetServer(log)

	repo, err := postgres.New(dbCfg)
	if err != nil {
		log.Error("couldn't connect to database", "error", err)
		os.Exit(1)
	}

	srv := server.New(ctx, srvCfg, repo)
	go func() {
		log.Info("listening and serving", "address", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("failed to listen and serve", "error", err)
			cancelCtx() // Trigger shutdown if server fails
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	successfulShutdown := true
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		// Give 10 seconds until forced shutdown
		shutdownCtx, cancelCtx := context.WithTimeout(shutdownCtx, time.Second*10)
		defer cancelCtx()

		err = srv.Shutdown(shutdownCtx)
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
