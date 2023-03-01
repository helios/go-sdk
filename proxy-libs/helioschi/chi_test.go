package helioschi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
)

const requestBody = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\""
const responseBody = "user test (id abcd1234)"

func validateAttributes(attrs []attribute.KeyValue, t *testing.T, metadataOnly bool) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "POST", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/users/abcd1234", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == semconv.HTTPRouteKey {
			assert.Equal(t, "/users/{id}", value.Value.AsString())
		} else if key == "http.request.body" {
			assert.False(t, metadataOnly)
			assert.Equal(t, requestBody, value.Value.AsString())
		} else if key == "http.response.body" {
			assert.False(t, metadataOnly)
			assert.Equal(t, responseBody, value.Value.AsString())
		} else if key == "http.request.headers" {
			assert.False(t, metadataOnly)
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			assert.Equal(t, "application/json", headers["Content-Type"][0])
		}
	}
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func runTests(t *testing.T, port string, metadataOnly bool) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	r := NewRouter()

	r.HandleFunc("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user User
		decoder.Decode(&user)

		id := URLParam(r, "id")
		name := "test"
		reply := fmt.Sprintf("user %s (id %s)", name, id)
		w.Write(([]byte)(reply))
	}))

	go func() {
		http.ListenAndServe(":"+port, r)
	}()

	url := fmt.Sprintf("http://localhost:%s/users/abcd1234", port)
	res, _ := http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, "/users/{id}", serverSpan.Name())
	assert.Equal(t, serverSpan.SpanKind(), trace.SpanKindServer)
	validateAttributes(serverSpan.Attributes(), t, metadataOnly)
	assert.Equal(t, res.Header.Get("traceresponse"), fmt.Sprintf("00-%s-%s-01", serverSpan.SpanContext().TraceID().String(), serverSpan.SpanContext().SpanID().String()))

	// Send again
	http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	spans = sr.Ended()
	serverSpan = spans[1]
	validateAttributes(serverSpan.Attributes(), t, metadataOnly)
}

func TestInstrumentation(t *testing.T) {
	runTests(t, "3333", false)
}

func TestInstrumentationInMetadataOnlyMode(t *testing.T) {
	os.Setenv("HS_METADATA_ONLY", "true")
	runTests(t, "3334", true)
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	port := "3335"
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	r := NewRouter()

	r.HandleFunc("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user User
		decoder.Decode(&user)

		id := URLParam(r, "id")
		name := "test"
		reply := fmt.Sprintf("user %s (id %s)", name, id)
		w.Write(([]byte)(reply))
	}))

	go func() {
		http.ListenAndServe(":" + port, r)
	}()

	url := fmt.Sprintf("http://localhost:%s/users/abcd1234", port)
	res, _ := http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 0, len(spans))
}
