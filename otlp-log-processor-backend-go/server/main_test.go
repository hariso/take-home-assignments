package server

import (
	"context"
	"fmt"
	"os"
	"testing"

	"dash0.com/otlp-log-processor-backend/telemetry"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	shutdownFn, err := telemetry.Setup(ctx, "test")
	if err != nil {
		fmt.Printf("error setting up OpenTelemetry: %v", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	err = shutdownFn(ctx)
	if err != nil {
		fmt.Printf("error shutting down OpenTelemetry: %v", err)
	}
	os.Exit(exitCode)
}
