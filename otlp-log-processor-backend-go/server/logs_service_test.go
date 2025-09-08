package server

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"go.uber.org/mock/gomock"
)

// Given an empty request, the server should call the counter,
// print the results after the given window, and return an empty response.
func TestLogsServiceServer_Export_CountPrint_EmptyRequest(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	testCountWindow := 100 * time.Millisecond

	counter := NewLogsCounter(ctrl)
	printer := NewCountPrinter(ctrl)
	underTest := newLogsServiceServer(testCountWindow, counter, printer)

	req := &collogspb.ExportLogsServiceRequest{}
	counter.EXPECT().count(ctx, req.ResourceLogs).Return()
	counts := map[string]int64{
		"something": 123,
	}
	counter.EXPECT().getAndReset().Return(counts).AnyTimes()
	printer.EXPECT().print(counts).AnyTimes()

	got, err := underTest.Export(ctx, req)

	time.Sleep(testCountWindow + 10*time.Millisecond)

	is.NoErr(err)
	is.Equal(got, &collogspb.ExportLogsServiceResponse{})
}
