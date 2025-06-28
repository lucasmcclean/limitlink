package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lucasmcclean/limitlink/link"
	"github.com/lucasmcclean/limitlink/mongo"
	"github.com/lucasmcclean/limitlink/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Println("starting limitlink...")

	repo, err := mongo.New(ctx)
	if err != nil {
		log.Fatalf("error connecting to the database: %v\n", err)
	}

	srv := server.New(repo)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening and serving on: %s\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("received shutdown signal")

	case err := <-serverErr:
		ctx, _ := context.WithTimeout(ctx, 1*time.Second)

		closeErr := repo.Close(ctx)
		if closeErr != nil {
			log.Printf("error closing repo after server failure: %v", closeErr)
		}

		log.Fatalf("error listening and serving: %v\n", err)
	}

	log.Println("starting shutdown...")

	if !shutdown(srv, repo) {
		os.Exit(1)
	}
}

func shutdown(srv *http.Server, repo link.Repository) (ok bool) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := repo.Close(shutdownCtx)
	if err != nil {
		log.Printf("error closing database connection: %v\n", err)
		ok = false
	} else {
		log.Println("database connection closed successfully")
	}

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error shutting down http server: %v\n", err)
		ok = false
	} else {
		log.Println("http server shut down successfully")
	}

	return ok
}
