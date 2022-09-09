package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var provider, _ = Initialize("test_service", "abcd1234", WithCollectorEndpoint(""))
var exporter = tracetest.NewInMemoryExporter()
var tracer = provider.Tracer("test")

func init() {
	provider.RegisterSpanProcessor(trace.NewSimpleSpanProcessor(exporter))
}

func TestCreateSpan(t *testing.T) {
	spanName := "abcd1234"
	_, span := tracer.Start(context.Background(), spanName)
	span.End()
	exported := exporter.GetSpans()
	assert.Equal(t, exported[0].Name, spanName)
}
