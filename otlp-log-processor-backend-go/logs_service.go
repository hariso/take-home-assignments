package main

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/metric"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

var (
	logsReceivedCounter metric.Int64Counter
)

func init() {
	var err error
	logsReceivedCounter, err = meter.Int64Counter("com.dash0.homeexercise.logs.received",
		metric.WithDescription("The number of logs received by otlp-log-processor-backend"),
		metric.WithUnit("{log}"))
	if err != nil {
		panic(err)
	}
}

type dash0LogsServiceServer struct {
	// addr removed because it was unused
	attributeKey string
	countWindow  time.Duration

	collogspb.UnimplementedLogsServiceServer
}

func newServer(attributeKey string, countWindow time.Duration) collogspb.LogsServiceServer {
	s := &dash0LogsServiceServer{attributeKey: attributeKey, countWindow: countWindow}
	return s
}

func (l *dash0LogsServiceServer) Export(ctx context.Context, request *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	slog.DebugContext(ctx, "Received ExportLogsServiceRequest")
	logsReceivedCounter.Add(ctx, 1)

	// Do something with the logs

	return &collogspb.ExportLogsServiceResponse{}, nil
}
