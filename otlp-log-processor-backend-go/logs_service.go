package main

import (
	"context"
	"log/slog"
	"time"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	v1 "go.opentelemetry.io/proto/otlp/logs/v1"
)

type logsCounter interface {
	count([]*v1.ResourceLogs)
}

type countPrinter interface {
	print(map[any]int64)
}

type dash0LogsServiceServer struct {
	// addr removed because it was unused
	attributeKey string
	countWindow  time.Duration
	counter      logsCounter
	printer      countPrinter

	collogspb.UnimplementedLogsServiceServer
}

func newServer(attributeKey string, countWindow time.Duration) collogspb.LogsServiceServer {
	s := &dash0LogsServiceServer{
		attributeKey: attributeKey,
		countWindow:  countWindow,
		// todo inject as argument to newServer()
		counter: newInMemoryCounter(),
	}

	return s
}

func (s *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest", "resourceLogs.size", len(request.ResourceLogs))

	s.counter.count(request.ResourceLogs)

	return &collogspb.ExportLogsServiceResponse{}, nil
}
