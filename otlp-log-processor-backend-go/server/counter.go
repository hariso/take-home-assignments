package server

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"dash0.com/otlp-log-processor-backend/telemetry"
	"go.opentelemetry.io/otel/metric"
	v2 "go.opentelemetry.io/proto/otlp/common/v1"
	"go.opentelemetry.io/proto/otlp/logs/v1"
)

var unknownAttrKey = &v2.AnyValue{
	Value: &v2.AnyValue_StringValue{StringValue: "unknown"},
}

// inMemoryCounter is a logs' counter that stores the counts in memory.
type inMemoryCounter struct {
	// attrKey is the attribute key to use for counting.
	attrKey string
	// counts stores the counts for each resource/scope/attribute value.
	// The key is the stringified attribute value, the value is the count.
	counts      map[string]int64
	countsMutex sync.Mutex

	// TODO We can add more metrics, e.g. for the number of logs received by resource
	// (so that we can have a better idea of the load on the system).

	logsReceivedCounter metric.Int64Counter
}

func newInMemoryCounter(attrKey string) *inMemoryCounter {
	logsReceivedCounter, err := telemetry.Meter.Int64Counter(
		"com.dash0.homeexercise.logs.received",
		metric.WithDescription("The number of logs received by otlp-log-processor-backend"),
		metric.WithUnit("{log}"),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create logs received counter: %w", err))
	}

	return &inMemoryCounter{
		attrKey:             attrKey,
		counts:              make(map[string]int64, 1000),
		logsReceivedCounter: logsReceivedCounter,
	}
}

func (c *inMemoryCounter) count(ctx context.Context, resLogs []*v1.ResourceLogs) {
	c.countsMutex.Lock()
	defer c.countsMutex.Unlock()

	slog.DebugContext(ctx, "counting logs", "resourceLogs.size", len(resLogs))
	for _, resLog := range resLogs {
		c.countInResource(ctx, resLog)
	}
}

func (c *inMemoryCounter) getAndReset() map[string]int64 {
	c.countsMutex.Lock()
	defer c.countsMutex.Unlock()

	slog.Debug("getting and resetting counts")
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
	slog.Debug("counting logs in resource", "resource.name", resLog.Resource.String())

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
		slog.DebugContext(ctx, "attribute found on resource, counted log records",
			"resource", telemetry.Resource.String(),
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
	slog.Debug("counting logs in scope", "scope.name", scopeLog.Scope.String())

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
		slog.DebugContext(ctx, "attribute found on scope, counted log records",
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

	c.logsReceivedCounter.Add(ctx, int64(len(scopeLog.LogRecords)))
}
