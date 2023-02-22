package heliosmux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

const (
	expectedStatusCode = 200
	expectedBody       = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\"}"
)

const requestResponseBody = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\"}"
const obfuscatedRequestResponseBody = "{\"id\":123,\"name\":\"dac02c19\",\"role\":\"Software Engineer\"}"

func init() {
	blocklistRules, _ := json.Marshal([]string{"$.name"})
	os.Setenv("HS_DATA_OBFUSCATION_HMAC_KEY", "12345")
	os.Setenv("HS_DATA_OBFUSCATION_BLOCKLIST", string(blocklistRules))
}

func validateAttributes(attrs []attribute.KeyValue, path string, metadataOnly bool, t *testing.T) {
	requestBodyFound := false
	requestHeadersFound := false
	responseBodyFound := false
	for _, attribute := range attrs {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.HTTPHostKey:
			assert.Equal(t, "localhost:8000", value)
		case semconv.HTTPMethodKey:
			assert.Equal(t, "POST", value)
		case semconv.HTTPRouteKey:
			assert.Equal(t, path, value)
		case semconv.HTTPSchemeKey:
			assert.Equal(t, "http", value)
		case semconv.HTTPServerNameKey:
			assert.Equal(t, "opentelemetry-middleware", value)
		case semconv.HTTPTargetKey:
			assert.Equal(t, path, value)
		case "http.request.body":
			requestBodyFound = true
			assert.Equal(t, obfuscatedRequestResponseBody, value)
		case "http.request.headers":
			requestHeadersFound = true
			headers := map[string][]string{}
			json.Unmarshal([]byte(value), &headers)
			assert.Equal(t, "application/json", headers["Content-Type"][0])
		case "http.response.body":
			responseBodyFound = true
			assert.Equal(t, obfuscatedRequestResponseBody, value)
		}
	}

	assert.Equal(t, metadataOnly, !requestBodyFound)
	assert.Equal(t, metadataOnly, !requestHeadersFound)
	assert.Equal(t, metadataOnly, !responseBodyFound)
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func postUser(responseWriter http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var user User
	decoder.Decode(&user)
	responseWriter.Header().Set("content-type", "application/json")
	json.NewEncoder(responseWriter).Encode(user)
}

func testHelper(t *testing.T, metadataOnly bool, path string) {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	router := NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s", path), http.HandlerFunc(postUser))
	http.Handle(fmt.Sprintf("/%s", path), router)
	go func() { http.ListenAndServe(":8000", nil) }()

	response, _ := http.Post(fmt.Sprintf("http://localhost:8000/%s", path), "application/json", bytes.NewBuffer([]byte(requestResponseBody)))
	statusCode := response.StatusCode
	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, expectedStatusCode, statusCode)
	assert.Equal(t, requestResponseBody, strings.TrimSpace(string(body)))
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 1, len(spans))
	span := spans[0]
	assert.False(t, span.Parent().HasTraceID())
	assert.False(t, span.Parent().HasSpanID())
	assert.True(t, span.SpanContext().HasTraceID())
	assert.True(t, span.SpanContext().HasSpanID())
	assert.Equal(t, fmt.Sprintf("/%s", path), span.Name())
	assert.Equal(t, trace.SpanKindServer, span.SpanKind())

	validateAttributes(span.Attributes(), fmt.Sprintf("/%s", path), metadataOnly, t)
	assert.Equal(t, response.Header.Get("traceresponse"), fmt.Sprintf("00-%s-%s-01", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String()))
}

func TestNewRouterInstrumentation(t *testing.T) {
	testHelper(t, false, "allData")
}

func TestNewRouterInstrumentationMetadataOnly(t *testing.T) {
	os.Setenv("HS_METADATA_ONLY", "true")
	testHelper(t, true, "metadataOnly")
}
