package main

import (
	"context"
	"event-ingestion-service/internal/httpserver"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer close()

	server := httpserver.New(os.Getenv("PORT"))

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed starting the server")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}
}
