package dataobfuscator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ohler55/ojg/jp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
	"golang.org/x/exp/slices"
)

const lengthToObfuscatedByteArray = 8

var DATA_TO_OBFUSCATE = []string{"http.request.body", "http.response.body", "db.query_result", string(semconv.DBStatementKey), "messaging.payload", "faas.event", "faas.res"}
var hMacKey []byte

func obfuscateMap(aMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(aMap))
	for key, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			result[key] = obfuscateMap(val.(map[string]interface{}))
		case []interface{}:
			result[key] = obfuscateArray(val.([]interface{}))
		default:
			result[key] = obfuscatePrimitives(val)
		}
	}
	return result
}

func obfuscateArray(anArray []interface{}) []interface{} {
	result := make([]interface{}, len(anArray))
	for i, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			result[i] = obfuscateMap(val.(map[string]interface{}))
		case []interface{}:
			result[i] = obfuscateArray(val.([]interface{}))
		default:
			result[i] = obfuscatePrimitives(val)
		}
	}
	return result
}

func modifyElement(element any) (any, bool) {
	var result any
	switch element.(type) {
	case map[string]interface{}:
		result = obfuscateMap(element.(map[string]interface{}))
	case []interface{}:
		result = obfuscateArray(element.([]interface{}))
	default:
		result = obfuscatePrimitives(element)
	}

	return result, true
}

func obfuscatePrimitives(val any) any {
	_, isbool := val.(bool)
	if val == nil || isbool {
		return val
	}
	var bs []byte
	h := hmac.New(sha256.New, hMacKey)
	switch val.(type) {
	case string:
		h.Write([]byte(val.(string)))
		bs = h.Sum(nil)
	default:
		h.Write([]byte(fmt.Sprintf("%v", val)))
		bs = h.Sum(nil)
	}

	return hex.EncodeToString(bs)[:lengthToObfuscatedByteArray]
}

func obfuscateDataHelper(value attribute.Value, obfuscationMode string, obfuscationRules []jp.Expr) attribute.Value {
	var attrValueAsJson map[string]interface{}
	if value.Type() == attribute.STRING {
		err := json.Unmarshal([]byte(value.AsString()), &attrValueAsJson)
		if err != nil {
			obfuscatedPrimitiveValue := obfuscatePrimitives(value.AsString())
			return attribute.StringValue(obfuscatedPrimitiveValue.(string))

		}
		var result any
		switch obfuscationMode {
		case "blocklist":
			for _, rule := range obfuscationRules {
				result, err = rule.Modify(attrValueAsJson, modifyElement)
				if err != nil {
					log.Printf("Failed applying obfuscation in blocklist mode")
					return value
				}
			}
			data, _ := json.Marshal(result)
			return attribute.StringValue(string(data))
		}
	}
	return value
}

func obfuscateAttributeValue(attribute attribute.KeyValue) attribute.Value {
	obfuscationConfig := getObfuscationConfig()
	hMacKey = []byte(obfuscationConfig.obfuscationhmacKey)
	if obfuscationConfig.obfuscationEnabled && slices.Contains(DATA_TO_OBFUSCATE, string(attribute.Key)) {
		return obfuscateDataHelper(attribute.Value, obfuscationConfig.obfuscationMode, obfuscationConfig.obfuscationRules)
	}
	return attribute.Value
}
