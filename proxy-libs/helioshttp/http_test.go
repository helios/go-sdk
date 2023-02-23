package helioshttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/helios/opentelemetry-go-contrib/instrumentation/net/http/otelhttp"
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
const obfuscatedRequestBody = "{\"id\":123,\"name\":\"dac02c19\",\"role\":\"Software Engineer\"}"
const responseBody = "hello1234"
const obfuscatedResponseBody = "87468e56"

func init() {
	blocklistRules, _ := json.Marshal([]string{"$.name"})
	os.Setenv("HS_DATA_OBFUSCATION_HMAC_KEY", "12345")
	os.Setenv("HS_DATA_OBFUSCATION_BLOCKLIST", string(blocklistRules))
}

func getHello(responseWriter ResponseWriter, request *Request) {
	body, _ := io.ReadAll(request.Body)
	if string(body) != requestBody {
		log.Fatal("Invalid request body")
	}
	io.WriteString(responseWriter, responseBody)
}

func validateAttributes(attrs []attribute.KeyValue, path string, metadataOnly bool, t *testing.T) {
	requestBodyFound := false
	requestHeadersFound := false
	responseBodyFound := false
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "POST", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/"+path, value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == "http.response.body" {
			responseBodyFound = true
			assert.Equal(t, obfuscatedResponseBody, value.Value.AsString())
		} else if key == "http.request.headers" {
			requestHeadersFound = true
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			assert.Equal(t, "application/json", headers["Content-Type"][0])
		} else if key == "http.request.body" {
			requestBodyFound = true
			assert.Equal(t, obfuscatedRequestBody, value.Value.AsString())
		}
	}

	assert.Equal(t, metadataOnly, !requestBodyFound)
	assert.Equal(t, metadataOnly, !requestHeadersFound)
	assert.Equal(t, metadataOnly, !responseBodyFound)
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
	assert.Equal(t, res.Header.Get("traceresponse"), fmt.Sprintf("00-%s-%s-01", serverSpan.SpanContext().TraceID().String(), serverSpan.SpanContext().SpanID().String()))

	// Send again
	res, _ = Post(fmt.Sprintf("http://localhost:%d/%s", port, path), "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ = io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(context.Background())
	spans = sr.Ended()
	serverSpan = spans[2]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), path, metadataOnly, t)
	clientSpan = spans[3]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), path, metadataOnly, t)
}

func TestServerInstrumentation(t *testing.T) {
	testHelper(t, 8000, "test1", false)
}

func TestServerInstrumentationMetadataOnly(t *testing.T) {
	os.Setenv("HS_METADATA_ONLY", "true")
	// Reset the client so that metadaaonly mode canbe properly applied
	otelhttp.DefaultClient = &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	DefaultClient = &Client{}
	testHelper(t, 8001, "test2", true)
}
