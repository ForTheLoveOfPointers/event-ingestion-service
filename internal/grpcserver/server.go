package grpcserver

import (
	"context"
	"encoding/json"
	"event-ingestion-service/internal/httpserver/types"
	"event-ingestion-service/internal/ingest"
	"event-ingestion-service/internal/metrics"
	"event-ingestion-service/proto/rpcserver"
	"fmt"
)

type ServergRPC struct {
	rpcserver.EventIngestionServer
	Ingestor *ingest.Ingest
	Metric   *metrics.Metric
}

func NewServergRPC(ingestor *ingest.Ingest, metric *metrics.Metric) *ServergRPC {
	return &ServergRPC{Ingestor: ingestor, Metric: metric}
}

func (s *ServergRPC) CreateEvent(ctx context.Context, req *rpcserver.EventRequest) (*rpcserver.IngestResponse, error) {

	var event types.Event

	raw, err := json.Marshal(req)

	if err != nil {
		return &rpcserver.IngestResponse{Status: metrics.StatusFailure}, err
	}

	if err := json.Unmarshal(raw, &event); err != nil {
		return &rpcserver.IngestResponse{Status: metrics.StatusFailure}, err
	}

	if ok := s.Ingestor.Enqueue(event); ok != true {
		err := fmt.Errorf("Too many requests")
		return &rpcserver.IngestResponse{Status: metrics.StatusFailure}, err
	}

	return &rpcserver.IngestResponse{Status: metrics.StatusSuccess}, nil
}
