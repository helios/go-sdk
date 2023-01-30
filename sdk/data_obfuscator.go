package sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/martian/v3/log"
	"github.com/ohler55/ojg/gen"
	"github.com/ohler55/ojg/jp"
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

func obfuscateDataHelper(value attribute.Value) string {
	var attrValueAsJson map[string]interface{}
	heliosConfig := getHeliosConfig()
	if heliosConfig == nil {
		log.Errorf("Can't apply obfuscation before configuration initialized")
		return ""
	}
	if value.Type() == attribute.STRING {
		err := json.Unmarshal([]byte(value.AsString()), &attrValueAsJson)
		if err != nil {
			log.Errorf("Failed parsing attribute value to json")
			return ""
		}
		var result any
		switch heliosConfig.obfuscationConfig.obfuscationMode {
		case "blocklist":
			for _, rule := range heliosConfig.obfuscationConfig.obfuscationRules {
				ruleAsExpr, err := jp.ParseString(rule)
				if err != nil {
					log.Errorf("Failed parsing obfuscation rule")
				}
				result, err = ruleAsExpr.Modify(attrValueAsJson, modifyElement)
				if err != nil {
					log.Errorf("Failed applying obfuscation in blocklist mode")
				}
			}
			data, _ := json.Marshal(result)
			return string(data)
		case "allowlist":
			allowlistedNodes := make(map[string][]gen.Node)
			p := gen.Parser{}
			// 1. get the value by the jsonpath expressions
			// 2. assign everyNodeRule := every node
			// 3. everyNodeRule.Modify(attrValueAsJson, modifyElement)
			// 4. for each jsonpath expression: modify to the original value

			// ruleAsExpr.Modify(attrValueAsJson, modifyElement)
			for _, rule := range heliosConfig.obfuscationConfig.obfuscationRules {
				ruleAsExpr, err := jp.ParseString(rule)
				if err != nil {
					log.Errorf("Failed parsing obfuscation rule")
					continue
				}
				obj, err := p.Parse([]byte(value.AsString()))
				nodes := ruleAsExpr.GetNodes(obj)
				nodeTest := ruleAsExpr.String()
				if len(nodes) > 0 {
					allowlistedNodes[rule] = nodes
					fmt.Printf("node: %v\n", nodeTest)
				}
			}
			allNodesRule := ".*.*"
			ruleAsExprAllData, err := jp.ParseString(allNodesRule)
			if err != nil {
				log.Errorf("Failed parsing all nodes expression")
			}
			result, err = ruleAsExprAllData.Modify(attrValueAsJson, modifyElement)

			for key, value := range allowlistedNodes {
				allowlistRuleAsExpr, err := jp.ParseString(key)
				if err != nil {
					log.Errorf("Failed parsing obfuscation allowlist rule")
				}
				allowlistRuleAsExpr.Modify(attrValueAsJson, func(element any) (any, bool) {
					return value, true
				})
			}
		}
	}
	return ""
}

// TODO - this method will be called from the reflected code in the processor
// func obfuscateData(attributes []attribute.KeyValue) []attribute.KeyValue {
// 	for _, attr := range attributes {
// 		if slices.Contains(DATA_TO_OBFUSCATE, string(attr.Key)) {
// 			obfuscatedVal := obfuscateDataHelper(attr.Value)
// 			if obfuscatedVal != "" {
// 				attributes[i] = obfuscatedVal
// 			}
// 		}
// 	}
// }

func obfuscateAttribute(attribute attribute.KeyValue) string {
	if slices.Contains(DATA_TO_OBFUSCATE, string(attribute.Key)) {
		return obfuscateDataHelper(attribute.Value)
	}
	return ""
}

func gethMacKey() []byte {
	heliosConfig := getHeliosConfig()
	if hMacKey == nil {
		hMacKey = make([]byte, heliosConfig.obfuscationConfig.obfuscationhmacKey)
	}
	return hMacKey
}
