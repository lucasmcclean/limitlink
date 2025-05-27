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

	"github.com/lucasmcclean/limitlink/link"
	"github.com/lucasmcclean/limitlink/server"
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	repo, err := link.NewMongo("temp", "temp", "temp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to database: %s\n", err)
		cancelCtx()
		os.Exit(1)
	}

	srv := server.New(repo)
	go func() {
		log.Printf("listening and serving on: %s\n", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
			cancelCtx()
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		log.Println("shutting down http srv...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelShutdown()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http srv: %s\n", err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}
