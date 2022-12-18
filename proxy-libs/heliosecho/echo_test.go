package heliosecho

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

func TestInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	r := New()

	r.GET("/users/:id", func(c Context) error {
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

	http.Get("http://localhost:8091/users/abcd1234")
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	serverSpan := spans[0]
	validateAttributes(serverSpan.Attributes(), t)
}

func TestInterfaceMatch(t *testing.T) {
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/labstack/echo", "v4.9.1")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "echo")
	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i, j int) bool {
		return originalExports[i].Name < originalExports[j].Name
	})

	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliosecho")
	sort.Slice(heliosExports, func(i, j int) bool {
		return heliosExports[i].Name < heliosExports[j].Name
	})

	// Compare.
	assert.EqualValues(t, originalExports, heliosExports)
}
