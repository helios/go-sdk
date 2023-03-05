package helioslogrus

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"encoding/json"

	"github.com/sirupsen/logrus"
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
	tracer := tp.Tracer("TestLogrus")
	ctx, span := tracer.Start(context.Background(), "hello-span")

	l := New()
	b := &bytes.Buffer{}
	l.Formatter = &logrus.JSONFormatter{}
	l.Out = b
	l.WithContext(ctx).Warn("test")
	data := map[string]string{}
	json.Unmarshal(b.Bytes(), &data)

	span.End()
	sr.ForceFlush(ctx)
	return data
}

func TestLogInstrumentation(t *testing.T) {
	sr, tp := initTracing(t)
	
	data := registerLoggerAndLogMessage(t, sr, tp)

	spans := sr.Ended()
	resultedSpan := spans[0]
	assert.Equal(t, 1, len(spans))
	value, ok := data["go_to_helios"]
	assert.True(t, ok)
	traceId := resultedSpan.SpanContext().TraceID().String()
	spanId := resultedSpan.SpanContext().SpanID().String()
	expectedGoToHeliosVal := fmt.Sprintf("https://app.gethelios.dev?actionTraceId=%s&spanId=%s&source=logrus", traceId, spanId)
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

	_, exists := data["go_to_helios"]
	assert.False(t, exists)
}
