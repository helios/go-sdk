package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var exporterObfuscatorTest *tracetest.InMemoryExporter
var providerObfuscatorTest *trace.TracerProvider

var blocklistRules = []string{
	"$.metadata.*",
	"$.collection",
	"$.details[*].name",
	"$.topic",
	"$.information[*].age",
	"$..information[?(@.address=='Unclassified')].address"}

const testServiceNameObfuscator = "test_service_obfuscator"


func initHelperObfuscator(samplingRatio float64) {
	if provider != nil {
		provider.Shutdown(context.Background())
	}
	providerSingelton = nil
	provider, _ = Initialize(serviceName, "abcd1234", WithCollectorEndpoint(""), WithSamplingRatio(samplingRatio), WithObfuscationBlocklistRules(blocklistRules), WithhmacKey("12345"))
	exporter = tracetest.NewInMemoryExporter()
	provider.RegisterSpanProcessor(trace.NewSimpleSpanProcessor(exporter))
}

func init() {
	initHelperObfuscator(1)
}


func TestObfuscationBlocklistDbStatement(t *testing.T) {
	stringAttr := "{\"collection\":\"spec\",\"details\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}]}"
	obfuscatedDataExpectedValue := "{\"collection\":\"d3ae0dfc\",\"details\":[{\"address\":\"New York\",\"age\":35,\"male\":true,\"name\":\"dac02c19\",\"null\":null},{\"address\":\"Jerusalem\",\"age\":42,\"extra\":\"field\",\"male\":false,\"name\":\"f175ac0e\"}]}"
	keyValueAttr := attribute.KeyValue{
		Key:   "db.statement",
		Value: attribute.StringValue(stringAttr),
	}
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := obfuscateAttributeValue(attrs[0]).AsString()
	assert.Equal(t, obfuscatedDataExpectedValue, obfuscatedData)
}

func TestObfuscationBlocklistHttpRequestBody(t *testing.T) {
	stringAttr := "{\"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null,\"metadata\":{\"date\":\"2022-04-01T00:00:00.000Z\",\"count\":5}}"
	obfuscatedDataExpectedValue := "{\"address\":\"New York\",\"age\":35,\"male\":true,\"metadata\":{\"count\":\"7337c795\",\"date\":\"c6e6d6c3\"},\"name\":\"Lior Govrin\",\"null\":null}"
	keyValueAttr := attribute.KeyValue{
		Key:   "http.request.body",
		Value: attribute.StringValue(stringAttr),
	}
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := obfuscateAttributeValue(attrs[0]).AsString()
	assert.Equal(t, obfuscatedDataExpectedValue, obfuscatedData)
}

func TestObfuscationBlocklistMessagingPayload(t *testing.T) {
	stringAttr := "{\"topic\":\"test\",\"information\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"},{\"name\":\"Bob Wilson\",\"male\":true,\"age\":100,\"address\":\"Unclassified\",\"extra\":\"field\"}]}"
	obfuscatedDataExpectedValue := "{\"information\":[{\"address\":\"New York\",\"age\":\"8baca100\",\"male\":true,\"name\":\"Lior Govrin\",\"null\":null},{\"address\":\"Jerusalem\",\"age\":\"99f900ae\",\"extra\":\"field\",\"male\":false,\"name\":\"Alice Smith\"},{\"address\":\"119b419b\",\"age\":\"fde61109\",\"extra\":\"field\",\"male\":true,\"name\":\"Bob Wilson\"}],\"topic\":\"e031ba1c\"}"
	keyValueAttr := attribute.KeyValue{
		Key:   "messaging.payload",
		Value: attribute.StringValue(stringAttr),
	}
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := obfuscateAttributeValue(attrs[0]).AsString()
	assert.Equal(t, obfuscatedDataExpectedValue, obfuscatedData)
}

func TestObfuscationBlocklistNonJsonVal(t *testing.T) {
	stringAttr := "test"
	obfuscatedDataExpectedValue := "e031ba1c"
	keyValueAttr := attribute.KeyValue{
		Key:   "faas.event",
		Value: attribute.StringValue(stringAttr),
	}
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := obfuscateAttributeValue(attrs[0]).AsString()
	assert.Equal(t, obfuscatedDataExpectedValue, obfuscatedData)
}

func TestObfuscationBlocklistDontObfuscateNonRelevantKey(t *testing.T) {
	stringAttr := "{\"collection\":\"spec\",\"details\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}]}"

	keyValueAttr := attribute.KeyValue{
		Key:   "span.name",
		Value: attribute.StringValue(stringAttr),
	}
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := obfuscateAttributeValue(attrs[0]).AsString()
	assert.Equal(t, stringAttr, obfuscatedData)
}
