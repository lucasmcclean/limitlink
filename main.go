package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lucasmcclean/limitlink/link"
	"github.com/lucasmcclean/limitlink/server"
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	log.Println("starting limitlink...")

	repo, err := link.NewMongo()
	if err != nil {
		log.Printf("error connecting to database: %s\n", err)
		os.Exit(1)
	}

	srv := server.New(repo)
	go func() {
		log.Printf("listening and serving on: %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error listening and serving: %s\n", err)
			cancelCtx()
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	log.Println("starting shutdown...")

	success := true

	if err = repo.Close(shutdownCtx); err != nil {
		log.Printf("error closing database connection: %s\n", err)
		success = false
	} else {
		log.Println("database connection closed successfully")
	}

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error shutting down http server: %s\n", err)
		success = false
	} else {
		log.Println("http server shut down successfully")
	}

	if !success {
		os.Exit(1)
	}
}
