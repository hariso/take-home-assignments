package server

import (
	"context"
	"testing"

	"github.com/matryer/is"
	v1 "go.opentelemetry.io/proto/otlp/common/v1"
	logsv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

// TODO More tests are needed for various cases:
// 1. Attribute not found anywhere (counts go towards the "unknown" key).
// 2. Attr found on resource or scope, but not overridden.
// 3. Attr found on resource or scope, but overridden on a lower level (scope or log).
// 4. etc.

func TestInMemoryCounter_Count_ResourceAttr(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	resLog := []*logsv1.ResourceLogs{
		{
			Resource: &resourcev1.Resource{
				Attributes: []*v1.KeyValue{
					makeKeyValue("foo", "bar"),
				},
			},
			ScopeLogs: []*logsv1.ScopeLogs{
				{
					Scope: &v1.InstrumentationScope{
						Attributes: nil,
					},
					LogRecords: []*logsv1.LogRecord{
						{
							Attributes: nil,
						},
					},
				},
			},
		},
	}
	underTest := newInMemoryCounter("foo")
	underTest.count(ctx, resLog)
	counts := underTest.getAndReset()
	is.Equal(len(counts), 1)
	count, ok := counts["bar"]
	is.True(ok) // expected value not found
	is.Equal(count, int64(1))
}

func makeKeyValue(key string, value string) *v1.KeyValue {
	return &v1.KeyValue{
		Key:   key,
		Value: &v1.AnyValue{Value: &v1.AnyValue_StringValue{StringValue: value}},
	}
}
