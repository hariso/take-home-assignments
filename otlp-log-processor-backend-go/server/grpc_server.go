package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"

	"dash0.com/otlp-log-processor-backend/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const name = "dash0.com/otlp-log-processor-backend"

// Run runs the gRPC server.
func Run(ctx context.Context) error {
	otelShutdown, err := telemetry.Setup(ctx, name)
	if err != nil {
		return fmt.Errorf("error setting up OpenTelemetry: %w", err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	cfg, err := parseConfig()
	if err != nil {
		slog.ErrorContext(ctx, "error parsing config", "err", err)
		printHelp()
		os.Exit(1)
	}

	slog.Debug("Starting listener", slog.String("listenAddr", cfg.listenAddr))
	listener, err := net.Listen("tcp", cfg.listenAddr)
	if err != nil {
		return fmt.Errorf("error listening at %v: %w", cfg.listenAddr, err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.MaxRecvMsgSize(cfg.maxReceiveMessageSize),
		grpc.Creds(insecure.NewCredentials()),
	)
	collogspb.RegisterLogsServiceServer(
		grpcServer,
		newLogsServiceServer(
			cfg.countWindow,
			newInMemoryCounter(cfg.attributeKey),
			newStdoutPrinter(),
		),
	)

	slog.Debug("Starting gRPC server")

	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("error serving gRPC server: %w", err)
	}

	return nil
}

func printHelp() {
	fmt.Println(
		`Usage:
  otlp-log-processor-backend [flags]

Flags:
  -listenAddr string
        The address to listen on for gRPC requests (default "localhost:4317").
  -maxReceiveMessageSize int
        The max message size in bytes the server can receive (default 16777216).
  -attributeKey string
        (REQUIRED) Attribute key for which the numbers of distinct values should be tracked.
  -countWindow duration
        Duration of the time window after which the number of distinct values of attributeKey will be printed (default 10s)`,
	)
}
