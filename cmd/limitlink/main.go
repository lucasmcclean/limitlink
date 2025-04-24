package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lucasmcclean/url-shortener/internal/config"
	"github.com/lucasmcclean/url-shortener/internal/link"
	"github.com/lucasmcclean/url-shortener/internal/logger"
	"github.com/lucasmcclean/url-shortener/internal/router"
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	log := logger.NewStdLogger(logger.DEBUG, os.Stdout)
	config.InitApp(log)
	if config.App.IsProd() {
		log = logger.NewStdLogger(logger.INFO, os.Stderr)
	}

	dbCfg := config.GetDB(log)
	srvCfg := config.GetServer(log)

	repo, err := link.NewRepository(dbCfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	log.Info("database connection established")

	err = repo.Migrate(log)
	if err != nil {
		log.Error("failed to run database migrations", "error", err)
		repo.Close()
		os.Exit(1)
	}
	log.Info("fatabase migrations ran successfully")

  rtr := router.New(repo)

	srv := newServer(rtr, srvCfg)
	go func() {
		log.Info("server listening and serving", "address", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err)
			cancelCtx()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	successfulShutdown := true
	go func() {
		defer wg.Done()

		<-ctx.Done()

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelShutdown()

		err = srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Error("failed to shut down server", "error", err)
			successfulShutdown = false
		}

		err = repo.Close()
		if err != nil {
			log.Error("failed to close database connection", "error", err)
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

func newServer(handler http.Handler, srvCfg *config.Server) *http.Server {
	return &http.Server{
		Addr:              srvCfg.Port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
