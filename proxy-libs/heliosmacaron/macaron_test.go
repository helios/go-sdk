package heliosmacaron

import (
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

	exportsExtractor "github.com/helios/helios-go-instrumenter/exports_extractor"
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
	// r := New()
	// tmplName := "user"
	// tmplStr := "user {{ .name }} (id {{ .id }})\n"
	// tmpl := template.Must(template.New(tmplName).Parse(tmplStr))
	// r.SetHTMLTemplate(tmpl)
	// r.GET("/users/:id", func(c *Context) {
	// 	id := c.Param("id")
	// 	c.HTML(http.StatusOK, tmplName, H{
	// 		"name": "whatever",
	// 		"id":   id,
	// 	})
	// })

	// go func() {
	// 	_ = r.Run(":8090")
	// }()

	// http.Get("http://localhost:8090/users/abcd1234")
	// sr.ForceFlush(context.Background())
	// spans := sr.Ended()
	// assert.Equal(t, 1, len(spans))
	// serverSpan := spans[0]
	// validateAttributes(serverSpan.Attributes(), t)
}

func TestInterfaceMatch(t *testing.T) {
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/go-macaron/macaron", "v1.4.0")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "macaron")
	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i, j int) bool {
		return originalExports[i].Name < originalExports[j].Name
	})

	srcDir, _ := filepath.Abs(".")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliosmacaron")
	sort.Slice(heliosExports, func(i, j int) bool {
		return heliosExports[i].Name < heliosExports[j].Name
	})

	// Compare.
	assert.EqualValues(t, originalExports, heliosExports)
}
