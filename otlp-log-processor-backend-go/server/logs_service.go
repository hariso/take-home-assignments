package server

import (
	"context"
	"log/slog"
	"time"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	v1 "go.opentelemetry.io/proto/otlp/logs/v1"
)

// logsCounter counts the occurrences of an attribute's value in a set of ResourceLogs.
//
//go:generate mockgen -typed -source=logs_service.go -destination=counter_mock.go -package=server -mock_names=logsCounter=LogsCounter . logsCounter
type logsCounter interface {
	// count counts the logs in the given ResourceLogs.
	count(context.Context, []*v1.ResourceLogs)
	// getAndReset returns the current counts and resets the internal counter to 0.
	getAndReset() map[string]int64
}

// countPrinter prints the counts (occurrences of an attribute's value) to stdout.'
//
//go:generate mockgen -typed -source=logs_service.go -destination=printer_mock.go -package=server -mock_names=countPrinter=CountPrinter . countPrinter
type countPrinter interface {
	print(map[string]int64)
}

// dash0LogsServiceServer is an implementation of v1.LogsServiceServer that
// counts the occurrences of an attribute's value in a set of ResourceLogs,
// and periodically prints the counts using the configured printer.
type dash0LogsServiceServer struct {
	counter     logsCounter
	printTicker *time.Ticker
	printer     countPrinter

	collogspb.UnimplementedLogsServiceServer
}

func newLogsServiceServer(countWindow time.Duration, counter logsCounter, printer countPrinter) collogspb.LogsServiceServer {
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
	// todo this runs in a goroutine, cancel when service shuts down
	for range s.printTicker.C {
		slog.Debug("printing counts")
		counts := s.counter.getAndReset()
		s.printer.print(counts)
	}
}
