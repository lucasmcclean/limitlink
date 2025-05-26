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
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	server := newServer()
	go func() {
		log.Printf("listening and serving on: %s\n", server.Addr)
		err := server.ListenAndServe()
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

		log.Println("shutting down http server...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelShutdown()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}
