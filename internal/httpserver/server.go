package httpserver

import (
	"context"
	"event-ingestion-service/internal/httpserver/handlers"
	"event-ingestion-service/internal/httpserver/middleware"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(s string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Healthz())
	mux.HandleFunc("POST /events", handlers.EventHandler())

	return &Server{
		httpServer: &http.Server{
			Addr:              s,
			Handler:           middleware.Limiter(mux),
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
