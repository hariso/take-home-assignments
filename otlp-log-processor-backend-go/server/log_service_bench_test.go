package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"dash0.com/otlp-log-processor-backend/telemetry"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logsv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

func BenchmarkLogServiceServer_Export_LargeRequest(b *testing.B) {
	ctx := context.Background()
	shutdown, err := telemetry.Setup(ctx, b.Name())
	if err != nil {
		b.Fatalf("error setting up OpenTelemetry: %v", err)
	}
	b.Cleanup(func() {
		shutdownErr := shutdown(ctx)
		if err != nil {
			b.Logf("error shutting down OpenTelemetry: %v", shutdownErr)
		}
	})

	underTest := newLogsServiceServer( // adapt if it's not exported
		time.Millisecond,
		newInMemoryCounter("foo"),
		newStdoutPrinter(),
	)

	// Create a representative request with some logs
	req := makeExportLogsServiceRequest(1_000_000)

	// Reset timer so setup is not measured
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := underTest.Export(ctx, req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func makeExportLogsServiceRequest(recs int) *collogspb.ExportLogsServiceRequest {
	req := &collogspb.ExportLogsServiceRequest{
		ResourceLogs: []*logsv1.ResourceLogs{
			{
				Resource: &resourcev1.Resource{
					Attributes: []*commonv1.KeyValue{
						{
							Key: "service.name",
							Value: &commonv1.AnyValue{
								Value: &commonv1.AnyValue_StringValue{StringValue: "benchmark-service"},
							},
						},
					},
				},
				ScopeLogs: []*logsv1.ScopeLogs{
					{
						Scope: &commonv1.InstrumentationScope{
							Name:    "benchmark-scope",
							Version: "1.0.0",
						},
					},
				},
			},
		},
	}

	logRecords := make([]*logsv1.LogRecord, recs)
	for i := 0; i < recs; i++ {
		logRecords[i] = &logsv1.LogRecord{
			Body: &commonv1.AnyValue{Value: &commonv1.AnyValue_StringValue{
				StringValue: fmt.Sprintf("log-record-%d", i),
			}},
			Attributes: []*commonv1.KeyValue{
				{
					Key: "foo",
					Value: &commonv1.AnyValue{
						Value: &commonv1.AnyValue_StringValue{StringValue: "bar"},
					},
				},
			},
		}
	}

	req.ResourceLogs[0].ScopeLogs[0].LogRecords = logRecords

	return req
}
