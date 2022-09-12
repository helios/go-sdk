package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var exporter *tracetest.InMemoryExporter
var provider *trace.TracerProvider

const testServiceName = "test_service"

func init() {
	provider, _ = Initialize(serviceName, "abcd1234", WithCollectorEndpoint(""))
	exporter = tracetest.NewInMemoryExporter()
	provider.RegisterSpanProcessor(trace.NewSimpleSpanProcessor(exporter))
}

func TestCreateCustomSpanNoCallback(t *testing.T) {
	exporter.Reset()
	spanName := "abcd1234"
	CreateCustomSpan(context.Background(), spanName, []attribute.KeyValue{}, nil)
	exported := exporter.GetSpans()
	customSpan := exported[0]
	assert.Equal(t, customSpan.Name, spanName)
	serviceName := customSpan.Resource.Attributes()[0]
	assert.Equal(t, serviceName.Value.AsString(), testServiceName)
}

func TestCreateCustomSpanWithCallback(t *testing.T) {
	exporter.Reset()
	spanName := "abcd1234"
	var value = 1
	attrKey := "abcd"
	attrValue := "1234"
	keyValueAttr := attribute.KeyValue{
		Key:   attribute.Key(attrKey),
		Value: attribute.StringValue(attrValue),
	}
	CreateCustomSpan(context.Background(), spanName, []attribute.KeyValue{keyValueAttr}, func() {
		value = 2
	})
	exported := exporter.GetSpans()
	assert.Equal(t, value, 2)
	customSpan := exported[0]
	assert.Equal(t, customSpan.Name, spanName)
	var found bool = false
	for i := range customSpan.Attributes {
		if string(customSpan.Attributes[i].Key) == attrKey {
			assert.Equal(t, customSpan.Attributes[i].Value.AsString(), attrValue)
			found = true
		}
	}
	assert.True(t, found)
}
