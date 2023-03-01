package helioshttptest

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

const response = "abcd1234"

func assertSpan(t *testing.T, span sdktrace.ReadOnlySpan, hostname string) {
	attributes := span.Attributes()
	assert.Contains(t, attributes, attribute.String("http.response.body", response))
	assert.Contains(t, attributes, attribute.String("http.method", "GET"))
	assert.Contains(t, attributes, attribute.String("http.host", hostname[7:]))
}

func TestHeliosHttpTest(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	ts := NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		require.NoError(t, err)
	}))

	client := ts.Client()
	res, err := client.Get(ts.URL)
	resBytes, _ := io.ReadAll(res.Body)
	assert.Equal(t, response, string(resBytes))
	assert.Nil(t, err)
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKindServer, serverSpan.SpanKind())
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKindClient, clientSpan.SpanKind())
	assertSpan(t, serverSpan, ts.URL)
	assertSpan(t, clientSpan, ts.URL)
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")

	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	ts := NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		require.NoError(t, err)
	}))

	client := ts.Client()
	res, err := client.Get(ts.URL)
	resBytes, _ := io.ReadAll(res.Body)
	assert.Equal(t, response, string(resBytes))
	assert.Nil(t, err)
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 0, len(spans))
}
