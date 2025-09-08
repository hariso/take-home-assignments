package main

import (
	"context"
	"log/slog"
	"time"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	v1 "go.opentelemetry.io/proto/otlp/logs/v1"
)

//go:generate mockgen -typed -source=logs_service.go -destination=logs_counter_mock.go -package=main -mock_names=logsCounter=LogsCounter . logsCounter
type logsCounter interface {
	// count counts the logs in the given ResourceLogs.
	count(context.Context, []*v1.ResourceLogs)
	// getAndReset returns the current counts and resets the internal counter to 0.
	getAndReset() map[string]int64
}

//go:generate mockgen -typed -source=logs_service.go -destination=logs_printer_mock.go -package=main -mock_names=countPrinter=CountPrinter . countPrinter
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
	slog.Debug("starting printer")
	// todo cancel when service shuts down
	for range s.printTicker.C {
		slog.Debug("printing counts")
		counts := s.counter.getAndReset()
		s.printer.print(counts)
	}
}
