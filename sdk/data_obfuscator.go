package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ohler55/ojg/jp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
	"golang.org/x/exp/slices"
)
const lengthToObfuscatedByteArray = 8
var DATA_TO_OBFUSCATE = []string{"http.request.body", "http.response.body", "db.query_result", string(semconv.DBStatementKey), "messaging.payload", "faas.event", "faas.res"}
var hMacKey []byte

func parseMap(aMap map[string]interface{}) (any,bool) {
	var result any
	var changed bool
	for _, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			result, changed = parseMap(val.(map[string]interface{}))
		case []interface{}:
			result, changed = parseArray(val.([]interface{}))
		default:
			result, changed = obfuscatePrimitives(val)
		}
	}
	return result, changed
}

func parseArray(anArray []interface{}) (any, bool) {
	var result any
	var changed bool
	for _, val := range anArray {
		switch val.(type) {
		case map[string]interface{}:
			result, changed = parseMap(val.(map[string]interface{}))
		case []interface{}:
			result, changed = parseArray(val.([]interface{}))
		default:
			result, changed = obfuscatePrimitives(val)
		}
	}
	return result, changed
}

func modifyElement(element any) (any, bool) {
	var result any
	var changed bool
	switch element.(type) {
	case map[string]interface{}:
		result, changed = parseMap(element.(map[string]interface{}))
	case []interface{}:
		result, changed = parseArray(element.([]interface{}))
	default:
		result, changed = obfuscatePrimitives(element)
	}

	return result, changed
}

func obfuscatePrimitives(val any) (any, bool) {
	var bs []byte
	if val != nil {
		h := hmac.New(sha256.New, gethMacKey())
		switch val.(type) {
		case string:
			h.Write([]byte(val.(string)))
			bs = h.Sum(nil)
		case int:
			h.Write([]byte(strconv.Itoa(val.(int))))
			bs = h.Sum(nil)
		case float32:
			h.Write([]byte(fmt.Sprintf("%f", val.(float32))))
			bs = h.Sum(nil)
		case float64:
			h.Write([]byte(fmt.Sprintf("%f", val.(float64))))
			bs = h.Sum(nil)
		}
		if bs != nil {
			return hex.EncodeToString(bs)[:lengthToObfuscatedByteArray], true
		}
	}
	return val, false
}

func obfuscateDataHelper(value attribute.Value, obfuscationMode string, obfuscationRules []jp.Expr) attribute.Value {
	var attrValueAsJson map[string]interface{}
	if value.Type() == attribute.STRING {
		err := json.Unmarshal([]byte(value.AsString()), &attrValueAsJson)
		if err != nil {
			obfuscatedPrimitiveValue, obfuscated := obfuscatePrimitives(value.AsString())
			if obfuscated  {
				return attribute.StringValue(obfuscatedPrimitiveValue.(string))
			} else {
				return value
			}
		}
		var result any
		switch obfuscationMode {
		case "blocklist":
			for _, rule := range obfuscationRules {
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
	heliosConfig := getHeliosConfig()
	if heliosConfig != nil && heliosConfig.obfuscationConfig.obfuscationEnabled && slices.Contains(DATA_TO_OBFUSCATE, string(attribute.Key)) {
		return obfuscateDataHelper(attribute.Value, heliosConfig.obfuscationConfig.obfuscationMode, heliosConfig.obfuscationConfig.obfuscationRules)
	}
	return attribute.Value
}

func gethMacKey() []byte {
	heliosConfig := getHeliosConfig()
	if hMacKey == nil {
		hMacKey = []byte(heliosConfig.obfuscationConfig.obfuscationhmacKey)
	}
	return hMacKey
}
