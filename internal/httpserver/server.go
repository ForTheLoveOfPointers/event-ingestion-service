package httpserver

import (
	"context"
	"event-ingestion-service/internal/httpserver/handlers"
	"event-ingestion-service/internal/httpserver/middleware"
	"event-ingestion-service/internal/ingest"
	"event-ingestion-service/internal/metrics"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	ingestor   *ingest.Ingest
}

func New(s string, bufferSize int, ctx context.Context, wp *ingest.WorkerPool) *Server {
	mux := http.NewServeMux()
	ingestor := ingest.New(bufferSize)
	m := metrics.NewMetric()

	batcher := ingest.Make(100, time.Second)
	batcher.Start(ingestor, ctx)

	wp.Start(batcher, ctx)

	mux.HandleFunc("/healthz", handlers.Healthz())
	mux.HandleFunc("/api-key", handlers.APIKeyGen())
	mux.HandleFunc("/events", middleware.APIKeyMiddleware(
		handlers.EventHandler(ingestor, m),
	),
	)

	mux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

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
