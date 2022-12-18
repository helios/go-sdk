package heliosmux

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	expectedStatusCode = 200
	expectedBody       = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\"}"
)

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

func TestNewRouterInstrumentation(t *testing.T) {
	spanRecorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	router := NewRouter()
	router.HandleFunc("/users", http.HandlerFunc(postUser))
	http.Handle("/", router)
	go func() { http.ListenAndServe(":8000", nil) }()

	user, _ := json.Marshal(User{Id: 123, Name: "Lior Govrin", Role: "Software Engineer"})
	response, _ := http.Post("http://localhost:8000/users", "application/json", bytes.NewBuffer(user))
	statusCode := response.StatusCode
	body, _ := io.ReadAll(response.Body)

	assert.Equal(t, expectedStatusCode, statusCode)
	assert.Equal(t, expectedBody, strings.TrimSpace(string(body)))
	spanRecorder.ForceFlush(context.Background())
	spans := spanRecorder.Ended()
	assert.Equal(t, 1, len(spans))
	span := spans[0]
	assert.False(t, span.Parent().HasTraceID())
	assert.False(t, span.Parent().HasSpanID())
	assert.True(t, span.SpanContext().HasTraceID())
	assert.True(t, span.SpanContext().HasSpanID())
	assert.Equal(t, "/users", span.Name())
	assert.Equal(t, trace.SpanKindServer, span.SpanKind())

	for _, attribute := range span.Attributes() {
		key := attribute.Key
		value := attribute.Value.AsString()

		switch key {
		case semconv.HTTPHostKey:
			assert.Equal(t, "localhost:8000", value)
		case semconv.HTTPMethodKey:
			assert.Equal(t, "POST", value)
		case semconv.HTTPRouteKey:
			assert.Equal(t, "/users", value)
		case semconv.HTTPSchemeKey:
			assert.Equal(t, "http", value)
		case semconv.HTTPServerNameKey:
			assert.Equal(t, "opentelemetry-middleware", value)
		case semconv.HTTPTargetKey:
			assert.Equal(t, "/users", value)
		}
	}
}

func TestInterfaceMatch(t *testing.T) {
	// Get original mux exports.
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/gorilla/mux", "v1.8.0")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "mux")
	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i int, j int) bool {
		return originalExports[i].Name < originalExports[j].Name
	})

	// Get Helios mux exports.
	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliosmux")
	sort.Slice(heliosExports, func(i int, j int) bool {
		return heliosExports[i].Name < heliosExports[j].Name
	})

	assert.Equal(t, len(originalExports), len(heliosExports))
	assert.EqualValues(t, originalExports, heliosExports)
}
