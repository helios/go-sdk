package helioszerolog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	// "os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestZerolog(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	tracer := tp.Tracer("TestZerolog")
	ctx, span := tracer.Start(context.Background(), "hello-span")

	out := &bytes.Buffer{}
	newCtx := New(out).WithContext(ctx)
	newLogger := Ctx(newCtx)
	newLogger.Info().Msg("test")
	data := map[string]string{}
	json.Unmarshal(out.Bytes(), &data)

	span.End()
	sr.ForceFlush(ctx)
	spans := sr.Ended()
	resultedSpan := spans[0]
	assert.Equal(t, 1, len(spans))
	value, ok := data["go_to_helios"]
	assert.True(t, ok)
	traceId := resultedSpan.SpanContext().TraceID().String()
	spanId := resultedSpan.SpanContext().SpanID().String()
	testedVal := fmt.Sprintf("https://app.gethelios.dev?actionTraceId=%s&spanId=%s&source=zerolog", traceId, spanId)
	assert.True(t, strings.HasPrefix(value, testedVal))
	value, ok = data["traceId"]
	assert.True(t, ok)
	assert.Equal(t, traceId, value)
	value, ok = data["spanId"]
	assert.True(t, ok)
	assert.Equal(t, spanId, value)

}
