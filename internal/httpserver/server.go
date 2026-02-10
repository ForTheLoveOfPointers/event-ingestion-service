package httpserver

import (
	"context"
	"event-ingestion-service/internal/httpserver/handlers"
	"event-ingestion-service/internal/httpserver/middleware"
	"event-ingestion-service/internal/ingest"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	ingestor   *ingest.Ingest
}

func New(s string, bufferSize int) *Server {
	mux := http.NewServeMux()
	ingestor := ingest.New(bufferSize)

	mux.HandleFunc("/healthz", handlers.Healthz())
	mux.HandleFunc("/events", handlers.EventHandler(ingestor))

	return &Server{
		httpServer: &http.Server{
			Addr:              s,
			Handler:           middleware.Limiter(mux),
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
		ingestor: ingestor,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
