package heliosmacaron

import (
	"context"
	"net/http"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/stretchr/testify/assert"
)

func validateAttributes(attrs []attribute.KeyValue, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "GET", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/users/abcd1234", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		} else if key == semconv.HTTPRouteKey {
			assert.Equal(t, "/users/:id", value.Value.AsString())
		}
	}
}



func initTracing(t *testing.T) *tracetest.SpanRecorder {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return sr
}

func registerMacaronAndPerformCall(t *testing.T) {
	m := Classic()
	m.Get("/users/:id", func(ctx *Context) string {
		id := ctx.Params("id")
		return id
	})

	go func() {
		m.Run()
	}()

	http.Get("http://localhost:4000/users/abcd1234")
}

func TestInstrumentation(t *testing.T) {
	sr := initTracing(t)
	registerMacaronAndPerformCall(t)
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	serverSpan := spans[0]
	validateAttributes(serverSpan.Attributes(), t)
}

func TestDisableInstrumentation(t *testing.T) {
	os.Setenv("HS_DISABLED", "true")
	defer os.Setenv("HS_DISABLED", "")
	sr := initTracing(t)
	registerMacaronAndPerformCall(t)
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 0, len(spans))
}
