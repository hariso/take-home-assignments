package main

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/proto/otlp/logs/v1"
)

type inMemoryCounter struct {
}

func newInMemoryCounter() *inMemoryCounter {
	logsReceivedCounter, err := meter.Int64Counter(
		"com.dash0.homeexercise.logs.received",
		metric.WithDescription("The number of logs received by otlp-log-processor-backend"),
		metric.WithUnit("{log}"),
	)
	if err != nil {
		panic(err)
	}
	return &inMemoryCounter{}
}

func (c *inMemoryCounter) count(logs []*v1.ResourceLogs) {
	//TODO implement me
	panic("implement me")
}
