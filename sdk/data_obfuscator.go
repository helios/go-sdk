package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"log"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
	"golang.org/x/exp/slices"
)

var DATA_TO_OBFUSCATE = []string{"http.reqeust.body", "http.response.body", "db.query_result", string(semconv.DBStatementKey), "messaging.payload"}
var hMacKey []byte

func parseMap(aMap map[string]interface{}) any {
	var result any
	for _, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			result = parseMap(val.(map[string]interface{}))
		case []interface{}:
			result = parseArray(val.([]interface{}))
		default:
			result = obfuscateStringOrInt(val)
		}
	}
	return result
}

func parseArray(anArray []interface{}) any {
	var result any
	for _, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			result = parseMap(val.(map[string]interface{}))
		case []interface{}:
			result = parseArray(val.([]interface{}))
		default:
			result = obfuscateStringOrInt(val)
		}
	}
	return result
}

func obfuscateStringOrInt(val any) any {
	var bs any
	if val != nil {
		h := hmac.New(sha256.New, gethMacKey())
		if a, ok := (val.(string)); ok {
			h.Write([]byte(a))
			bs = h.Sum(nil)
		} else if a, ok := (val.(int)); ok {
			h.Write([]byte(strconv.Itoa(a)))
			bs = h.Sum(nil)
		}
		if bs != nil {
			return bs
		}
	}
	return val
}

func modifyElement(element any) (any, bool) {
	var result any
	switch element.(type) {
	case map[string]interface{}:
		result = parseMap(element.(map[string]interface{}))
	default:
		result = obfuscateStringOrInt(element)
	}
	return result, true
}

func obfuscateDataHelper(value attribute.Value) attribute.Value {
	var attrValueAsJson map[string]interface{}
	heliosConfig := getHeliosConfig()
	if heliosConfig == nil {
		log.Printf("Can't apply obfuscation before configuration initialized")
		return value
	}
	if value.Type() == attribute.STRING {
		err := json.Unmarshal([]byte(value.AsString()), &attrValueAsJson)
		if err != nil {
			log.Printf("Failed parsing attribute value to json")
			return value
		}
		var result any
		switch heliosConfig.obfuscationConfig.obfuscationMode {
		case "blocklist":
			for _, rule := range heliosConfig.obfuscationConfig.obfuscationRules {
				result, err = rule.Modify(attrValueAsJson, modifyElement)
				if err != nil {
					log.Printf("Failed applying obfuscation in blocklist mode")
				}
			}
			data, _ := json.Marshal(result)
			return attribute.StringValue(string(data))
		}
	}
	return value
}

func obfuscateAttributeValue(attribute attribute.KeyValue) attribute.Value {
	if slices.Contains(DATA_TO_OBFUSCATE, string(attribute.Key)) {
		return obfuscateDataHelper(attribute.Value)
	}
	return attribute.Value
}

func gethMacKey() []byte {
	heliosConfig := getHeliosConfig()
	if hMacKey == nil {
		hMacKey = make([]byte, heliosConfig.obfuscationConfig.obfuscationhmacKey)
	}
	return hMacKey
}
