package heliosmacaron

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
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
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

func TestInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	m := Classic()
	m.Get("/users/:id", func(ctx *Context) string {
		id := ctx.Params("id")
		return id
	})

	go func() {
		m.Run()
	}()

	http.Get("http://localhost:4000/users/abcd1234")
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 1, len(spans))
	serverSpan := spans[0]
	validateAttributes(serverSpan.Attributes(), t)
}

func removeDuplicateValues(slice []exportsExtractor.ExtractedObject) []exportsExtractor.ExtractedObject {
	keys := make(map[string]bool)
	list := []exportsExtractor.ExtractedObject{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry.Name]; !value {
			keys[entry.Name] = true
			list = append(list, entry)
		}
	}
	return list
}

func TestInterfaceMatch(t *testing.T) {
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/go-macaron/macaron", "v1.4.0")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "macaron")
	// Macaron has separate implementations for PathUnescape for Go 1.17 and 1.18 - until the extractor properly
	// handles that we're forced to remove duplicates
	originalExports = removeDuplicateValues(originalExports)

	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i, j int) bool {
		return originalExports[i].Name < originalExports[j].Name
	})

	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliosmacaron")
	sort.Slice(heliosExports, func(i, j int) bool {
		return heliosExports[i].Name < heliosExports[j].Name
	})

	for index, value := range originalExports {
		heliosVal := heliosExports[index]
		if value.Name == "NewRouteMap" {
			// The return value can't be used by the proxy lib as its not exported by the original package
			assert.Equal(t, value.FunctionReturnValues[0].AttributeType, "routeMap")
			assert.Equal(t, heliosVal.FunctionReturnValues[0].AttributeType, "interface{}")
		} else {
			assert.Equal(t, value, heliosVal)
		}
	}
}
