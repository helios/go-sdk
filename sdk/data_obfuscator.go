package sdk

import (
	"crypto/sha256"
	"encoding/json"
	"strconv"

	"github.com/google/martian/v3/log"
	"github.com/ohler55/ojg/jp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
)

var DATA_TO_OBFUSCATE = [5]string{"http.reqeust.body", "http.response.body", "db.query_result", string(semconv.DBSystemKey), "messaging.payload" }

func obfuscate_data(value attribute.Value, spanId string) any {
	var attrValueAsJson = interface{}(nil)
	rules := [6]string{
	"$.metadata.*",
	"$.collection",
	"$.details[*].name",
	"$.topic",
	"$.information[*].age",
	"$..information[?(@.address=='Unclassified')].address"}
	if( value.Type() == attribute.STRING ) {
		err := json.Unmarshal([]byte(value.AsString()), &attrValueAsJson)
	
		if err != nil {
			log.Errorf("Failed parsing attribute value to json")
		}
		var result any
		for _, rule := range rules {
			x, err := jp.ParseString(rule)

			if err != nil {
				log.Errorf("Failed parsing obfuscation rule")
			}
			result, err = x.Modify(attrValueAsJson, func(element any) (any, bool) {
				h := sha256.New()
				if a, ok := (element.(string)); ok {
					h.Write([]byte(a))
				} else if a, ok := (element.(int)); ok {
					h.Write([]byte(strconv.Itoa(a)))
				}
				bs := h.Sum(nil)
				return bs, true
				
			})

			if err != nil {
				log.Errorf("Failed parsing attribute value to json")
			}

		}
		data, _ := json.Marshal(result)
		return string(data)
	}
	return ""
}