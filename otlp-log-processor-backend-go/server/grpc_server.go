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

func Run(ctx context.Context) error {
	// Set up OpenTelemetry.
	otelShutdown, err := telemetry.Setup(ctx, name)
	if err != nil {
		return fmt.Errorf("error setting up OpenTelemetry: %w", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	cfg, err := parseConfig()
	if err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		printHelp()
		os.Exit(1)
	}

	slog.Debug("Starting listener", slog.String("listenAddr", cfg.listenAddr))
	listener, err := net.Listen("tcp", cfg.listenAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.MaxRecvMsgSize(cfg.maxReceiveMessageSize),
		grpc.Creds(insecure.NewCredentials()),
	)
	collogspb.RegisterLogsServiceServer(
		grpcServer,
		newServer(
			cfg.countWindow,
			newInMemoryCounter(cfg.attributeKey),
			newStdoutPrinter(),
		),
	)

	slog.Debug("Starting gRPC server")

	return grpcServer.Serve(listener)
}

func printHelp() {
	// todo pretty print the usage instructions
	fmt.Println("Usage: TODO")
}
