package helioshttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const responseBody = "hello1234"

func getHello(responseWriter ResponseWriter, request *Request) {
	io.WriteString(responseWriter, responseBody)
}

func validateAttributes(attrs []attribute.KeyValue, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "GET", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/test", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		}
	}
}

func validateServerAttributes(attrs []attribute.KeyValue, metadataOnly bool, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == "http.response.body" {
			assert.False(t, metadataOnly)
			assert.Equal(t, responseBody, value.Value.AsString())
		} else if key == "http.request.headers" {
			assert.False(t, metadataOnly)
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			assert.Equal(t, "gzip", headers["Accept-Encoding"][0])
		}
	}
}

func testHelper(t *testing.T, port int, metadataOnly bool) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	Handle("/test", HandlerFunc(getHello))
	go func() {
		ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	res, _ := Get(fmt.Sprintf("http://localhost:%d/test", port))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), t)
	validateServerAttributes(serverSpan.Attributes(), metadataOnly, t)
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), t)
	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
}

func TestServerInstrumentation(t *testing.T) {
	testHelper(t, 8000, false)
}

func TestServerInstrumentationMetadataOnly(t *testing.T) {
	os.Setenv("HS_METADATA_ONLY", "true")
	testHelper(t, 8001, true)
}
