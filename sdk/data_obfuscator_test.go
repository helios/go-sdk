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
	provider, _ = Initialize(serviceName, "abcd1234", WithCollectorEndpoint(""), WithSamplingRatio(samplingRatio), WithObfuscationBlocklistRules(blocklistRules), WithhmacKey("1234"))
	exporter = tracetest.NewInMemoryExporter()
	provider.RegisterSpanProcessor(trace.NewSimpleSpanProcessor(exporter))
}

func init() {
	initHelperObfuscator(1)
}


func TestObfuscationBlocklist(t *testing.T) {
	stringAttr := "{\"collection\":\"spec\",\"details\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}]}"
	obfuscatedDataExpectedValue := "{\"collection\":\"XHF3xRtbaOzWm4lxXtvBUi9HTArz+dw1Q7yxr5G7E0k=\",\"details\":[{\"address\":\"New York\",\"age\":35,\"male\":true,\"name\":\"n1IkECL7qb8VEcv9/7NUwFtabP6Gs8aW7fK7codetqE=\",\"null\":null},{\"address\":\"Jerusalem\",\"age\":42,\"extra\":\"field\",\"male\":false,\"name\":\"wEHZ1CbCELc6+Tv5RJgHetjYFeSZPvIh8KNIZcQZx7E=\"}]}"
	keyValueAttr := attribute.KeyValue{
		Key:   "db.statement",
		Value: attribute.StringValue(stringAttr),
	}
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := obfuscateAttributeValue(attrs[0]).AsString()
	assert.Equal(t, obfuscatedDataExpectedValue, obfuscatedData)
}
