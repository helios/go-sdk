package datautils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

var blocklistRules = []string{
	"$.metadata.*",
	"$.collection",
	"$.details[*].name",
	"$.nestedDetails",
	"$.topic",
	"$.information[*].age",
	"$..information[?(@.address=='Unclassified')].address"}

func initHelperObfuscator(samplingRatio float64) {
	rulesAsJsonString, _ := json.Marshal(blocklistRules)
	os.Setenv(hsDataObfuscationBlocklistEnvVar, string(rulesAsJsonString))
	os.Setenv(hsDatahMacKeyEnvVar, "12345")
	getObfuscationConfig()
}

func init() {
	initHelperObfuscator(1)
}

func validateObfusactedAttribute(t *testing.T, keyValueAttr attribute.KeyValue, obfuscatedDataExpectedValue string) {
	attrs := []attribute.KeyValue{keyValueAttr}
	obfuscatedData := ObfuscateAttributeValue(attrs[0])
	assert.Equal(t, obfuscatedDataExpectedValue, obfuscatedData.Value.AsString())
}

func TestObfuscationBlocklistDbStatement(t *testing.T) {
	stringAttr := "{\"collection\":\"spec\",\"details\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}]}"
	obfuscatedDataExpectedValue := "{\"collection\":\"d3ae0dfc\",\"details\":[{\"address\":\"New York\",\"age\":35,\"male\":true,\"name\":\"dac02c19\",\"null\":null},{\"address\":\"Jerusalem\",\"age\":42,\"extra\":\"field\",\"male\":false,\"name\":\"f175ac0e\"}]}"
	keyValueAttr := attribute.KeyValue{
		Key:   "db.statement",
		Value: attribute.StringValue(stringAttr),
	}
	validateObfusactedAttribute(t, keyValueAttr, obfuscatedDataExpectedValue)
}

func TestNestedObfuscationBlocklistDbStatement(t *testing.T) {
	stringAttr := "{\"collection\":\"spec\",\"details\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}], \"nestedDetails\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}]}"
	obfuscatedDataExpectedValue := "{\"collection\":\"d3ae0dfc\",\"details\":[{\"address\":\"New York\",\"age\":35,\"male\":true,\"name\":\"dac02c19\",\"null\":null},{\"address\":\"Jerusalem\",\"age\":42,\"extra\":\"field\",\"male\":false,\"name\":\"f175ac0e\"}],\"nestedDetails\":[{\"address\":\"7d639f21\",\"age\":\"2df9a61a\",\"male\":true,\"name\":\"dac02c19\",\"null\":null},{\"address\":\"7d7ae621\",\"age\":\"44343c77\",\"extra\":\"f08b0238\",\"male\":false,\"name\":\"f175ac0e\"}]}"
	keyValueAttr := attribute.KeyValue{
		Key:   "db.statement",
		Value: attribute.StringValue(stringAttr),
	}
	validateObfusactedAttribute(t, keyValueAttr, obfuscatedDataExpectedValue)
}

func TestObfuscationBlocklistHttpRequestBody(t *testing.T) {
	stringAttr := "{\"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null,\"metadata\":{\"date\":\"2022-04-01T00:00:00.000Z\",\"count\":5}}"
	obfuscatedDataExpectedValue := "{\"address\":\"New York\",\"age\":35,\"male\":true,\"metadata\":{\"count\":\"07eb9d8b\",\"date\":\"c6e6d6c3\"},\"name\":\"Lior Govrin\",\"null\":null}"
	keyValueAttr := attribute.KeyValue{
		Key:   "http.request.body",
		Value: attribute.StringValue(stringAttr),
	}
	validateObfusactedAttribute(t, keyValueAttr, obfuscatedDataExpectedValue)
}

func TestObfuscationBlocklistMessagingPayload(t *testing.T) {
	stringAttr := "{\"topic\":\"test\",\"information\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"},{\"name\":\"Bob Wilson\",\"male\":true,\"age\":100,\"address\":\"Unclassified\",\"extra\":\"field\"}]}"
	obfuscatedDataExpectedValue := "{\"information\":[{\"address\":\"New York\",\"age\":\"2df9a61a\",\"male\":true,\"name\":\"Lior Govrin\",\"null\":null},{\"address\":\"Jerusalem\",\"age\":\"44343c77\",\"extra\":\"field\",\"male\":false,\"name\":\"Alice Smith\"},{\"address\":\"119b419b\",\"age\":\"960e14e3\",\"extra\":\"field\",\"male\":true,\"name\":\"Bob Wilson\"}],\"topic\":\"e031ba1c\"}"
	keyValueAttr := attribute.KeyValue{
		Key:   "messaging.payload",
		Value: attribute.StringValue(stringAttr),
	}
	validateObfusactedAttribute(t, keyValueAttr, obfuscatedDataExpectedValue)
}

func TestObfuscationBlocklistNonJsonVal(t *testing.T) {
	stringAttr := "test"
	obfuscatedDataExpectedValue := "e031ba1c"
	keyValueAttr := attribute.KeyValue{
		Key:   "faas.event",
		Value: attribute.StringValue(stringAttr),
	}
	validateObfusactedAttribute(t, keyValueAttr, obfuscatedDataExpectedValue)
}

func TestObfuscationBlocklistDontObfuscateNonRelevantKey(t *testing.T) {
	stringAttr := "{\"collection\":\"spec\",\"details\":[{ \"name\":\"Lior Govrin\",\"male\":true,\"age\":35,\"address\":\"New York\",\"null\":null},{\"name\":\"Alice Smith\",\"male\":false,\"age\":42,\"address\":\"Jerusalem\",\"extra\":\"field\"}]}"

	keyValueAttr := attribute.KeyValue{
		Key:   "span.name",
		Value: attribute.StringValue(stringAttr),
	}
	validateObfusactedAttribute(t, keyValueAttr, stringAttr)
}
