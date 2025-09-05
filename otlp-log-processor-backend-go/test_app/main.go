package main

import (
	"context"
	"log"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var res = resource.NewWithAttributes(
	semconv.SchemaURL,
	semconv.ServiceNameKey.String("test-app"),
	semconv.ServiceNamespaceKey.String("dash0-exercise"),
	semconv.ServiceVersionKey.String("1.0.0"),
)

func main() {
	ctx := context.Background()

	// Create OTLP log exporter (to local OTel Collector on 4317)
	exp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		log.Fatalf("failed to create log exporter: %v", err)
	}

	// Create LoggerProvider with exporter
	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
	)
	defer func() { _ = lp.Shutdown(ctx) }()

	// Register as global provider
	global.SetLoggerProvider(lp)

	// Create slog.Logger that routes logs to OTel
	logger := otelslog.NewLogger("my-service")

	// Optionally set as default slog logger
	slog.SetDefault(logger)

	// Both work and go to OTLP collector
	logger.Info("hello from otelslog.NewLogger")
	slog.Info("hello from slog default")
}
