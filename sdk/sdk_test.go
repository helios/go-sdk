package sdk

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var exporter *tracetest.InMemoryExporter
var provider *trace.TracerProvider

const testServiceName = "test_service"

func initHelper(samplingRatio float64) {
	os.Setenv("HS_DISABLED", "")
	if provider != nil {
		provider.ForceFlush(context.Background())
		provider.Shutdown(context.Background())
	}
	providerSingleton = nil
	heliosConfigSingleton = nil
	provider, _ = Initialize(serviceName, "abcd1234", WithCollectorEndpoint(""), WithSamplingRatio(samplingRatio))
	otel.SetTracerProvider(provider)
	exporter = tracetest.NewInMemoryExporter()
	provider.RegisterSpanProcessor(trace.NewSimpleSpanProcessor(exporter))
}

func init() {
	initHelper(1)
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

func TestPropagateTestContext(t *testing.T) {
	exporter.Reset()
	ctx := context.Background()
	bMember, _ := baggage.NewMember(HELIOS_TEST_TRIGGERED_TRACE, "true")
	b, _ := baggage.New(bMember)
	ctx = baggage.ContextWithBaggage(ctx, b)

	tracer := provider.Tracer("helios")
	spanName := "test_test"
	_, span := tracer.Start(ctx, spanName)
	span.End()

	exported := exporter.GetSpans()
	testSpan := exported[0]
	attrs := []attribute.KeyValue{attribute.String(HELIOS_TEST_TRIGGERED_TRACE, "true")}
	assert.Equal(t, testSpan.Attributes, attrs)
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
	var foundUserAttr bool = false
	var foundCustomSpanAttr bool = false
	for i := range customSpan.Attributes {
		if string(customSpan.Attributes[i].Key) == attrKey {
			assert.Equal(t, customSpan.Attributes[i].Value.AsString(), attrValue)
			foundUserAttr = true
		}

		if string(customSpan.Attributes[i].Key) == customSpanAttr {
			assert.Equal(t, customSpan.Attributes[i].Value.AsString(), "true")
			foundCustomSpanAttr = true
		}
	}
	assert.True(t, foundUserAttr)
	assert.True(t, foundCustomSpanAttr)
}

func TestSamplerNoSampling(t *testing.T) {
	exporter.Reset()
	sampledCtx := CreateCustomSpan(context.Background(), "sampled1", []attribute.KeyValue{}, nil)
	exported := exporter.GetSpans()
	assert.Equal(t, len(exported), 1)
	assert.Equal(t, exported[0].Name, "sampled1")
	exporter.Reset()
	initHelper(0)
	CreateCustomSpan(sampledCtx, "sampled2", []attribute.KeyValue{}, nil)
	CreateCustomSpan(context.Background(), "not_sampled", []attribute.KeyValue{}, nil)
	exported = exporter.GetSpans()
	assert.Equal(t, 1, len(exported))
	assert.Equal(t, exported[0].Name, "sampled2")
}

func TestDisabledSDK(t *testing.T) {
	if provider != nil {
		provider.ForceFlush(context.Background())
		provider.Shutdown(context.Background())
	}
	providerSingleton = nil
	heliosConfigSingleton = nil
	_, err := Initialize(serviceName, "abcd1234", WithCollectorEndpoint(""), WithSamplingRatio(1), WithInstrumentationDisabled())
	assert.NotNil(t, err)	
}
