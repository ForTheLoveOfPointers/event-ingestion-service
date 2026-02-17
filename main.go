package main

import (
	"context"
	"event-ingestion-service/internal/grpcserver"
	"event-ingestion-service/internal/httpserver"
	"event-ingestion-service/internal/ingest"
	"event-ingestion-service/internal/metrics"
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

	workerPoolSz, err := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
	if err != nil {
		workerPoolSz = 20
	}

	wp := ingest.NewWorkerPool(workerPoolSz)
	defer wp.Wait()

	ingestor := ingest.New(bufferSz)
	m := metrics.NewMetric()

	server := httpserver.New(os.Getenv("PORT"), ctx, wp, ingestor, m)
	grpcServer := grpcserver.NewServergRPC(ingestor, m)
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed starting the server: %s\r\n", err.Error())
		}
	}()

	go func() {
		if err := grpcServer.Start("9001"); err != nil {
			log.Fatalf("Failed starting the grpc server: %s\r\n", err.Error())
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("\nGraceful shutdown failed for http server: %v\r\n", err)
	}
	if err := grpcServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("\nGraceful shutdown failed for grpc server: %v\r\n", err)
	}

	log.Printf("\nShutdown completed\r\n")
}
