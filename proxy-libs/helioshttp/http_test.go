package helioshttp

import (
	"bytes"
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

const requestBody = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\"}"
const responseBody = "hello1234"

func getHello(responseWriter ResponseWriter, request *Request) {
	io.WriteString(responseWriter, responseBody)
}

func validateAttributes(attrs []attribute.KeyValue, path string, metadataOnly bool, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "POST", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/"+path, value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == "http.response.body" {
			assert.False(t, metadataOnly)
			assert.Equal(t, responseBody, value.Value.AsString())
		} else if key == "http.request.headers" {
			assert.False(t, metadataOnly)
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			assert.Equal(t, "application/json", headers["Content-Type"][0])
		} else if key == "http.request.body" {
			assert.False(t, metadataOnly)
			assert.Equal(t, requestBody, value.Value.AsString())
		}
	}
}

func testHelper(t *testing.T, port int, path string, metadataOnly bool) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	Handle("/"+path, HandlerFunc(getHello))
	go func() {
		ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	res, _ := Post(fmt.Sprintf("http://localhost:%d/%s", port, path), "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), path, metadataOnly, t)
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), path, metadataOnly, t)
	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
}

func TestServerInstrumentation(t *testing.T) {
	testHelper(t, 8000, "test1", false)
}

func TestServerInstrumentationMetadataOnly(t *testing.T) {
	os.Setenv("HS_METADATA_ONLY", "true")
	testHelper(t, 8001, "test2", true)
}
