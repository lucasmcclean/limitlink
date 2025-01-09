package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/lucasmcclean/url-shortener/db"
)

func main() {
	ctx := context.Background()

	err := run(ctx, os.Getenv)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func run(ctx context.Context, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	dbConn, err := db.New(getenv)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	// TODO: Implement NewServer()
	httpServer := http.Server{} // NewServer(ctx, getenv, dbConn)

	// Used to force early shutdown if http server errs
	httpServerErr := make(chan error)

	go func() {
		log.Printf("listening and serving on: %s", httpServer.Addr)
		// TODO: Setup TLS
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("error listening and serving: %s", err)
			httpServerErr <- err
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	var shutdownErr error = nil
	go func() {
		defer wg.Done()

		select {
		case err = <-httpServerErr:
			log.Printf("%s", err)
			shutdownErr = err

		case <-ctx.Done():
			shutdownCtx := context.Background()
			shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Second*10)
			defer cancel()

			err = httpServer.Shutdown(shutdownCtx)
			if err != nil {
				log.Printf("error shutting down HTTP server: %s", err)
				shutdownErr = err
			}
		}

		err = dbConn.Close()
		if err != nil {
			log.Printf("error closing databse connection: %s", err)
			if shutdownErr == nil {
				shutdownErr = err
			} else {
				shutdownErr = fmt.Errorf("%s & %s ", err, shutdownErr)
			}
		}
	}()

	wg.Wait()
	return shutdownErr
}
