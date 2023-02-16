package helioslogrus

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestLogInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	tracer := tp.Tracer("TestLogrus")
	ctx, span := tracer.Start(context.Background(), "hello-span")

	l := logrus.New()
	b := &bytes.Buffer{}
	l.Formatter = &logrus.JSONFormatter{}
	l.Out = b
	l.AddHook(NewHook(WithLevels(logrus.WarnLevel)))
	l.WithContext(ctx).Warn("test")
	data := map[string]string{}
	json.Unmarshal(b.Bytes(), &data)

	span.End()
	sr.ForceFlush(ctx)
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	value, ok := data["go_to_helios"]
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("https://app.gethelios.dev?actionTraceId=%s&spanId=%s", span.SpanContext().TraceID(), span.SpanContext().SpanID()), value)
}
