package helioslambda

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type mockRequest struct {
	Headers map[string]string
}

var (
	traceId       = "83d8d6c5347593d092e9409f4978bd51"
	parentSpanId  = "6f2a23d2d1e9159c"
	tracingHeader = "00" + "-" + traceId + "-" + parentSpanId + "-" + "01"
	testEvent     = mockRequest{Headers: map[string]string{"traceparent": tracingHeader}}
)

func TestInstrumentHandlerTracingWithMockPropagator(t *testing.T) {
	ctx := context.Background()
	exporter := tracetest.NewInMemoryExporter()
	otel.SetTracerProvider(trace.NewTracerProvider(trace.WithBatcher(exporter)))

	customerHandler := func(event mockRequest) (string, error) {
		return "hello world", nil
	}

	wrapped := instrumentHandler(customerHandler)

	wrappedCallable := reflect.ValueOf(wrapped)
	resp := wrappedCallable.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(testEvent)})
	assert.Len(t, resp, 2)
	assert.Equal(t, "hello world", resp[0].Interface())
	assert.Nil(t, resp[1].Interface())

	spans := exporter.GetSpans()
	assert.Len(t, spans, 1)
	span := spans[0]
	assert.Equal(t, traceId, span.SpanContext.TraceID().String())
	assert.Equal(t, parentSpanId, span.Parent.SpanID().String())
}
