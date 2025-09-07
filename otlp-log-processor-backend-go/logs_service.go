package main

import (
	"context"
	"log/slog"
	"time"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	v1 "go.opentelemetry.io/proto/otlp/logs/v1"
)

type logsCounter interface {
	count(context.Context, []*v1.ResourceLogs)
	getAndReset() map[string]int64
}

type countPrinter interface {
	print(map[string]int64)
}

type dash0LogsServiceServer struct {
	counter     logsCounter
	printTicker *time.Ticker
	printer     countPrinter

	collogspb.UnimplementedLogsServiceServer
}

func newServer(countWindow time.Duration, counter logsCounter, printer countPrinter) collogspb.LogsServiceServer {
	s := &dash0LogsServiceServer{
		printTicker: time.NewTicker(countWindow),
		counter:     counter,
		printer:     printer,
	}

	go s.startPrinter()

	return s
}

func (s *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest", "resourceLogs.size", len(request.ResourceLogs))

	s.counter.count(ctx, request.ResourceLogs)

	return &collogspb.ExportLogsServiceResponse{}, nil
}

func (s *dash0LogsServiceServer) startPrinter() {
	// todo cancel when service shuts down
	for range s.printTicker.C {
		counts := s.counter.getAndReset()
		s.printer.print(counts)
	}
}
