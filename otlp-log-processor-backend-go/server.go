package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const name = "dash0.com/otlp-log-processor-backend"

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
	logger = otelslog.NewLogger(name)
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	slog.SetDefault(logger)
	logger.Info("Starting application")

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(context.Background())
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	cfg, err := parseConfig()
	if err != nil {
		printHelp()
		return nil
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
	collogspb.RegisterLogsServiceServer(grpcServer, newServer(cfg.attributeKey, cfg.countWindow))

	slog.Debug("Starting gRPC server")

	return grpcServer.Serve(listener)
}

func printHelp() {
	// todo pretty print the usage instructions
	fmt.Println("Usage: ")
}
