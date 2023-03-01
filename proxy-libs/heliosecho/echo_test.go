package heliosecho

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const requestBody = "{\"id\":123,\"name\":\"Lior Govrin\",\"role\":\"Software Engineer\""
const responseBody = "{\"ID\":\"abcd1234\",\"Name\":\"Random\"}"

func validateAttributes(t *testing.T, attrs []attribute.KeyValue) {
	foundHeaders := false
	foundReqBody := false
	foundResBody := false
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "POST", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/users/abcd1234", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == semconv.HTTPRouteKey {
			assert.Equal(t, "/users/:id", value.Value.AsString())
		} else if key == "http.request.body" {
			foundReqBody = true
			assert.Equal(t, requestBody, value.Value.AsString())
		} else if key == "http.response.body" {
			foundResBody = true
			assert.Equal(t, responseBody, strings.Trim(value.Value.AsString(), "\n"))
		} else if key == "http.request.headers" {
			foundHeaders = true
			headers := map[string][]string{}
			json.Unmarshal([]byte(value.Value.AsString()), &headers)
			assert.Equal(t, "application/json", headers["Content-Type"][0])
		}
	}

	assert.True(t, foundHeaders)
	assert.True(t, foundReqBody)
	assert.True(t, foundResBody)
}

func TestInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	r := New()

	r.POST("/users/:id", func(c Context) error {
		print(io.ReadAll(c.Request().Body))
		id := c.Param("id")
		name := "Random"
		return c.JSON(http.StatusOK, struct {
			ID   string
			Name string
		}{
			ID:   id,
			Name: name,
		})
	})

	go func() {
		_ = r.Start(":8091")
	}()

	url := "http://localhost:8091/users/abcd1234"
	res, _ := http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, strings.Trim(string(body), "\n"))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	serverSpan := spans[0]
	validateAttributes(t, serverSpan.Attributes())

	// Send again
	http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	sr.ForceFlush(context.Background())
	serverSpan = sr.Ended()[1]
	validateAttributes(t, serverSpan.Attributes())
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	r := New()

	r.POST("/users/:id", func(c Context) error {
		print(io.ReadAll(c.Request().Body))
		id := c.Param("id")
		name := "Random"
		return c.JSON(http.StatusOK, struct {
			ID   string
			Name string
		}{
			ID:   id,
			Name: name,
		})
	})

	go func() {
		_ = r.Start(":8092")
	}()

	url := "http://localhost:8092/users/abcd1234"
	res, _ := http.Post(url, "application/json", bytes.NewBuffer([]byte(requestBody)))
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, strings.Trim(string(body), "\n"))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 0, len(spans))
}
