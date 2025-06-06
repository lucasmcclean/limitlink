package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lucasmcclean/limitlink/mongo"
	"github.com/lucasmcclean/limitlink/server"
)

//go:embed static/*
var staticEmbed embed.FS
var staticFS, _ = fs.Sub(staticEmbed, "static")

//go:embed templates/*
var templatesEmbed embed.FS
var templatesFS, _ = fs.Sub(templatesEmbed, "templates")

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	log.Println("starting limitlink...")

	repo, err := mongo.New(ctx)
	if err != nil {
		log.Printf("error connecting to database: %v\n", err)
		os.Exit(1)
	}

	srv := server.New(ctx, repo, staticFS, templatesFS)
	go func() {
		log.Printf("listening and serving on: %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error listening and serving: %v\n", err)
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
		log.Printf("error closing database connection: %v\n", err)
		success = false
	} else {
		log.Println("database connection closed successfully")
	}

	if err = srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error shutting down http server: %v\n", err)
		success = false
	} else {
		log.Println("http server shut down successfully")
	}

	if !success {
		os.Exit(1)
	}
}
