package helioszerolog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	// "os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func initTracing(t *testing.T) (*tracetest.SpanRecorder, *sdktrace.TracerProvider) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return sr, tp
}

func registerLoggerAndLogMessage(t *testing.T, sr *tracetest.SpanRecorder, tp *sdktrace.TracerProvider) map[string]string {
	tracer := tp.Tracer("TestZerolog")
	ctx, span := tracer.Start(context.Background(), "hello-span")

	out := &bytes.Buffer{}
	newLogger := NewWithContext(out, ctx)
	newLogger.Info().Msg("test")
	data := map[string]string{}
	json.Unmarshal(out.Bytes(), &data)

	span.End()
	sr.ForceFlush(ctx)
	return data
}

func TestZerolog(t *testing.T) {
	sr, tp := initTracing(t)

	data := registerLoggerAndLogMessage(t, sr, tp)

	spans := sr.Ended()
	resultedSpan := spans[0]
	assert.Equal(t, 1, len(spans))
	value, ok := data["go_to_helios"]
	assert.True(t, ok)
	traceId := resultedSpan.SpanContext().TraceID().String()
	spanId := resultedSpan.SpanContext().SpanID().String()
	expectedGoToHeliosVal := fmt.Sprintf("https://app.gethelios.dev?actionTraceId=%s&spanId=%s&source=zerolog", traceId, spanId)
	assert.True(t, strings.HasPrefix(value, expectedGoToHeliosVal))
	value, ok = data["traceId"]
	assert.True(t, ok)
	assert.Equal(t, traceId, value)
	value, ok = data["spanId"]
	assert.True(t, ok)
	assert.Equal(t, spanId, value)
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	sr, tp := initTracing(t)

	data := registerLoggerAndLogMessage(t, sr, tp)

	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	_, exists := data["go_to_helios"]
	assert.False(t, exists)
}
