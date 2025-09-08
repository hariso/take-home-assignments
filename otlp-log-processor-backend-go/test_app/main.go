package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
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
	attribute.KeyValue{
		Key:   "resource-foo",
		Value: attribute.StringValue("resource-hello"),
	},
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
	logger := otelslog.NewLogger("foo-service")

	for range time.Tick(100 * time.Millisecond) {
		if rand.Float64() < 0.5 {
			logger.Info("logging with key foo", slog.String("foo", fmt.Sprintf("foo-value-%d", rand.IntN(3))))
		} else {
			logger.Info("logging with key bar", slog.String("bar", fmt.Sprintf("bar-value-%d", rand.IntN(3))))
		}
	}
}
