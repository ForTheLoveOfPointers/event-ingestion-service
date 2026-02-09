package httpserver

import (
	"context"
	"event-ingestion-service/internal/httpserver/handlers"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(s string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlers.Healthz())

	return &Server{
		httpServer: &http.Server{
			Addr:              s,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
