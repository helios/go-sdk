package helioshttp

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/helios/helios-go-instrumenter/exports_extractor"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const responseBody = "hello1234"

func getHello(responseWriter ResponseWriter, request *Request) {
	io.WriteString(responseWriter, responseBody)
}

func validateAttributes(attrs []attribute.KeyValue, t *testing.T) {
	for _, value := range attrs {
		key := value.Key
		if key == semconv.HTTPMethodKey {
			assert.Equal(t, "GET", value.Value.AsString())
		} else if key == semconv.HTTPTargetKey {
			assert.Equal(t, "/test", value.Value.AsString())
		} else if key == semconv.HTTPStatusCodeKey {
			assert.Equal(t, 200, int(value.Value.AsInt64()))
		}
	}
}

func TestServerInstrumentation(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	Handle("/test", HandlerFunc(getHello))
	go func() {
		ListenAndServe(":8000", nil)
	}()

	res, _ := Get("http://localhost:8000/test")
	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, responseBody, string(body))
	sr.ForceFlush(context.Background())
	spans := sr.Ended()
	assert.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	assert.Equal(t, trace.SpanKind(2), serverSpan.SpanKind())
	validateAttributes(serverSpan.Attributes(), t)
	clientSpan := spans[1]
	assert.Equal(t, trace.SpanKind(3), clientSpan.SpanKind())
	assert.False(t, clientSpan.Parent().HasTraceID())
	validateAttributes(clientSpan.Attributes(), t)
	assert.Equal(t, serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
}

func TestInterfaceMatch(t *testing.T) {
	mainRepoFolder := exports_extractor.CloneGitRepository("https://github.com/golang/go", "go1.18.8")
	netHttpPackageName := "http"
	packagePath := filepath.Join(mainRepoFolder, "/src/net/http")
	netHttpExports := exports_extractor.ExtractExports(packagePath, netHttpPackageName)
	os.RemoveAll(mainRepoFolder)
	sort.Slice(netHttpExports, func(i, j int) bool {
		return netHttpExports[i].Name < netHttpExports[j].Name
	})

	heliosHttpRoot, _ := filepath.Abs(".")
	heliosHttpPackageName := "helioshttp"
	heliosHttpExports := exports_extractor.ExtractExports(heliosHttpRoot, heliosHttpPackageName)
	sort.Slice(heliosHttpExports, func(i, j int) bool {
		return heliosHttpExports[i].Name < heliosHttpExports[j].Name
	})

	assert.Equal(t, len(netHttpExports), len(heliosHttpExports))
	assert.EqualValues(t, netHttpExports, heliosHttpExports)
}
