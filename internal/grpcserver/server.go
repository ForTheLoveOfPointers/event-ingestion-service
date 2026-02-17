package grpcserver

import (
	"context"
	"encoding/json"
	"event-ingestion-service/internal/httpserver/types"
	"event-ingestion-service/internal/ingest"
	"event-ingestion-service/internal/metrics"
	"event-ingestion-service/proto/rpcserver"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ServergRPC struct {
	rpcserver.EventIngestionServer
	Ingestor   *ingest.Ingest
	Metric     *metrics.Metric
	grpcServer *grpc.Server
}

func NewServergRPC(ingestor *ingest.Ingest, metric *metrics.Metric) *ServergRPC {
	grpcServer := grpc.NewServer()
	return &ServergRPC{Ingestor: ingestor, Metric: metric, grpcServer: grpcServer}
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

func (s *ServergRPC) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	rpcserver.RegisterEventIngestionServer(s.grpcServer, s)

	log.Printf("gRPC server listening on %s\n", port)

	return s.grpcServer.Serve(lis)
}

func (s *ServergRPC) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}

}
