package data-obfuscator

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/ohler55/ojg/jp"
)


const hsDatahMacKeyEnvVar = "HS_DATA_OBFUSCATION_HMAC_KEY"
const hsDataObfuscationBlocklistEnvVar = "HS_DATA_OBFUSCATION_BLOCKLIST"
const hsDataObfuscationAllowlistEnvVAr = "HS_DATA_OBFUSCATION_ALLOWLIST"


type HeliosObfuscationConfig struct {
	obfuscationEnabled bool
	obfuscationMode    string
	obfuscationRules   []jp.Expr
	obfuscationhmacKey string
}

var obfuscationConfigSingelton HeliosObfuscationConfig

func getStringSliceConfig(envVar string, defaultValue []string) []string {
	envVarValue := os.Getenv(envVar)
	var returnVal []string
	if envVarValue != "" {
		json.Unmarshal([]byte(envVarValue), &returnVal)
		return returnVal
	}
	return defaultValue
}

func parseObfuscationRules(rules []string) []jp.Expr {
	parsedRules := []jp.Expr{}

	for _, rule := range rules {
		ruleAsExpr, err := jp.ParseString(rule)
		if err != nil {
			log.Printf("Failed parsing obfuscation rule %s", rule)
			continue
		}
		parsedRules = append(parsedRules, ruleAsExpr)
	}
	return parsedRules
}


func createObfuscationConfig()  HeliosObfuscationConfig{
	hsDataObfuscationBlocklist := getStringSliceConfig(hsDataObfuscationBlocklistEnvVar, []string{})
	hsDataObfuscationAllowlist := getStringSliceConfig(hsDataObfuscationAllowlistEnvVAr, []string{})
	hsDatahMacKey := os.Getenv(hsDatahMacKeyEnvVar)
	if hsDatahMacKey != "" {
		if len(hsDataObfuscationBlocklist) > 0 {
			return HeliosObfuscationConfig{true, "blocklist", parseObfuscationRules(hsDataObfuscationBlocklist), hsDatahMacKey}
		} else if len(hsDataObfuscationAllowlist) > 0 {
			return HeliosObfuscationConfig{true, "allowlist", parseObfuscationRules(hsDataObfuscationAllowlist), hsDatahMacKey}
		}
	}
	return HeliosObfuscationConfig{false, "", []jp.Expr{}, "0"}
}

func getObfuscationConfig() HeliosObfuscationConfig {
	if reflect.ValueOf(obfuscationConfigSingelton).IsZero() {
		obfuscationConfigSingelton = createObfuscationConfig()
	}
	return obfuscationConfigSingelton
}