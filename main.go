package main

import (
	"context"
	"event-ingestion-service/internal/httpserver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer close()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment variables")
	}

	var bufferSz int

	bufferSz, err := strconv.Atoi(os.Getenv("INGEST_BUFFER_SIZE"))
	if err != nil {
		bufferSz = 50
	}

	server := httpserver.New(os.Getenv("PORT"), bufferSz, ctx)
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed starting the server: %s\r\n", err.Error())
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("\nGraceful shutdown failed: %v\r\n", err)
	}

	log.Printf("\nGraceful shutdown successful\r\n")
}
