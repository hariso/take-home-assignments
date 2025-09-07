package main

import (
	"context"
	"log/slog"
	"sync"

	"go.opentelemetry.io/otel/metric"
	v2 "go.opentelemetry.io/proto/otlp/common/v1"
	"go.opentelemetry.io/proto/otlp/logs/v1"
)

var (
	logsReceivedCounter metric.Int64Counter
	unknownAttrKey      = &v2.AnyValue{
		Value: &v2.AnyValue_StringValue{StringValue: "unknown"},
	}
)

func init() {
	var err error
	logsReceivedCounter, err = meter.Int64Counter(
		"com.dash0.homeexercise.logs.received",
		metric.WithDescription("The number of logs received by otlp-log-processor-backend"),
		metric.WithUnit("{log}"),
	)
	if err != nil {
		panic(err)
	}
}

type inMemoryCounter struct {
	attrKey string
	counts  map[string]int64
	m       sync.Mutex
}

func newInMemoryCounter(attrKey string) *inMemoryCounter {
	return &inMemoryCounter{
		attrKey: attrKey,
		counts:  make(map[string]int64, 1000),
	}
}

func (c *inMemoryCounter) count(ctx context.Context, resLogs []*v1.ResourceLogs) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, resLog := range resLogs {
		c.countInResource(ctx, resLog)
	}
}

func (c *inMemoryCounter) getAndReset() map[string]int64 {
	c.m.Lock()
	defer c.m.Unlock()

	counts := c.counts
	c.counts = make(map[string]int64)

	return counts
}

func (c *inMemoryCounter) countAllLogRecords(slogs []*v1.ScopeLogs) int {
	var count int
	for _, l := range slogs {
		count += len(l.LogRecords)
	}

	return count
}

func (c *inMemoryCounter) countInResource(ctx context.Context, resLog *v1.ResourceLogs) {
	var resValue *v2.AnyValue
	// Get the attribute value from the resource, if present.
	for _, attr := range resLog.Resource.Attributes {
		if attr.Key == c.attrKey {
			resValue = attr.Value
		}
	}

	// Attribute found on resource.
	// We consider that the value is "inherited" to all the logs coming from this resource.
	// We do not handle the case where the attribute is overridden in scope attributes or
	// log record attributes, as it's not a usual thing to do (e.g. it's unusual for a log
	// record to override the service.name attribute from the resource).
	if resValue != nil {
		count := c.countAllLogRecords(resLog.ScopeLogs)
		c.counts[resValue.GetStringValue()] += int64(count)
		slog.DebugContext(ctx, "counted log records",
			"resource", res.String(),
			"attribute.value", resValue.GetStringValue(),
			"logs.count", count,
		)
		return
	}

	// check the scopes
	for _, scopeLog := range resLog.ScopeLogs {
		c.countInScope(ctx, scopeLog)
	}
}

func (c *inMemoryCounter) countInScope(ctx context.Context, scopeLog *v1.ScopeLogs) {
	// Get the attribute value from the scope, if present.
	// If not present, we take the attribute's value from the resource
	var scopeVal *v2.AnyValue
	for _, attr := range scopeLog.Scope.Attributes {
		if attr.Key == c.attrKey {
			scopeVal = attr.Value
		}
	}

	// Attribute found on scope.
	// We consider that the value is "inherited" to all the logs in this scope.
	// We do not handle the case where the attribute is overridden in log record attributes,
	// as it's not a usual thing to do (e.g. it's unusual for a log
	// record to override the otel.scope.name attribute from the scope).
	if scopeVal != nil {
		count := len(scopeLog.LogRecords)
		c.counts[scopeVal.GetStringValue()] += int64(count)
		slog.DebugContext(ctx, "counted log records",
			"scope", scopeLog.Scope.String(),
			"attribute.value", scopeVal.GetStringValue(),
			"logs.count", count,
		)

		return
	}

	for _, logRec := range scopeLog.LogRecords {
		attrVal := unknownAttrKey
		for _, attr := range logRec.Attributes {
			if attr.Key == c.attrKey {
				attrVal = attr.Value
			}
		}
		c.counts[attrVal.GetStringValue()]++
	}

	logsReceivedCounter.Add(ctx, int64(len(scopeLog.LogRecords)))
}
