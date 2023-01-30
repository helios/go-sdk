package sdk

import (
	"os"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

var heliosConfigSingletone *HeliosConfig

type HeliosObfuscationConfig struct {
	obfuscationEnabled bool
	obfuscationMode    string
	obfuscationRules   []string
	obfuscationhmacKey int
}

type HeliosConfig struct {
	serviceName       string
	apiToken          string
	sampler           trace.Sampler
	collectorInsecure bool
	collectorEndpoint string
	collectorPath     string
	environment       string
	commitHash        string
	debug             bool
	metadataOnly      bool
	obfuscationConfig HeliosObfuscationConfig
}

// Keys and their matching env vars
const samplingRatioKey = "samplingRatio"
const samplingRatioEnvVar = "HS_SAMPLING_RATIO"
const environmentKey = "environment"
const environmentEnvVar = "HS_ENVIRONMENT"
const collectorInsecureKey = "collectorInsecure"
const collectorInsecureEnvVar = "HS_COLLECTOR_INSECURE"
const collectorEndpointKey = "collectorEndpoint"
const collectorEndpointEnvVar = "HS_COLLECTOR_ENDPOINT"
const collectorPathKey = "collectorPath"
const collectorPathEnvVar = "HS_COLLECTOR_PATH"
const commitHashKey = "commitHash"
const commitHashEnvVar = "HS_COMMIT_HASH"
const debugKey = "debug"
const debugEnvVar = "HS_DEBUG"
const metadataOnlyKey = "metadataOnly"
const metadataOnlyEnvVar = "HS_METADATA_ONLY"
const hsDataObfuscationAllowlistEnvVAr = "HS_DATA_OBFUSCATION_ALLOWLIST"
const hsDataObfuscationAllowlistKey = "dataObfuscationAllowlist"
const hsDataObfuscationBlocklistEnvVar = "HS_DATA_OBFUSCATION_BLOCKLIST"
const hsDataObfuscationBlocklistKey = "dataObfuscationBlocklist"
const hsDatahMacKeyEnvVar = "HS_DATA_OBFUSCATION_HMAC_KEY"
const hsDatahMacKey = "dataObfuscationhMacKey"

// Default values
const defaultCollectorInsecure = false
const defaultCollectorEndpoint = "collector.heliosphere.io:443"
const defaultCollectorPath = "traces"
const defaultDebug = false
const defaultMetadataOnly = false

func getConfigByKey(key string, attrs []attribute.KeyValue) attribute.KeyValue {
	for i := range attrs {
		if string(attrs[i].Key) == key {
			return attrs[i]
		}
	}

	return attribute.KeyValue{Key: "", Value: attribute.Value{}}
}

func getSampler(attrs []attribute.KeyValue) trace.Sampler {
	samplerConfig := getConfigByKey(samplingRatioKey, attrs)
	samplingRatio := os.Getenv(samplingRatioEnvVar)
	if samplingRatio != "" {
		res, err := strconv.ParseFloat(samplingRatio, 64)
		if err != nil {
			return trace.AlwaysSample()
		}

		return NewHeliosSampler(res)
	}

	if samplerConfig.Key == "" {
		return trace.AlwaysSample()
	}

	return NewHeliosSampler(samplerConfig.Value.AsFloat64())
}

func getStringConfig(envVar string, defaultValue string, config attribute.KeyValue) string {
	envVarValue := os.Getenv(envVar)
	if envVarValue != "" {
		return envVarValue
	}

	if config.Key == "" {
		return defaultValue
	}

	return config.Value.AsString()
}

func getBoolConfig(envVar string, defaultValue bool, config attribute.KeyValue) bool {
	result, err := strconv.ParseBool(getStringConfig(envVar, strconv.FormatBool(defaultValue), config))
	if err != nil {
		return defaultValue
	}

	return result
}

func isCollectorInsecure(attrs []attribute.KeyValue) bool {
	collectorInsecureConfig := getConfigByKey(collectorInsecureKey, attrs)
	return getBoolConfig(collectorInsecureEnvVar, defaultCollectorInsecure, collectorInsecureConfig)
}

func isDebugMode(attrs []attribute.KeyValue) bool {
	debugConfig := getConfigByKey(debugKey, attrs)
	return getBoolConfig(debugEnvVar, defaultDebug, debugConfig)
}

func isMetadataOnlyMode(attrs []attribute.KeyValue) bool {
	metadataOnlyConfig := getConfigByKey(metadataOnlyKey, attrs)
	return getBoolConfig(metadataOnlyEnvVar, defaultMetadataOnly, metadataOnlyConfig)
}

func getCollectorEndpoint(attrs []attribute.KeyValue) string {
	collectorEndpointConfig := getConfigByKey(collectorEndpointKey, attrs)
	return getStringConfig(collectorEndpointEnvVar, defaultCollectorEndpoint, collectorEndpointConfig)
}

func getCollectorPath(attrs []attribute.KeyValue) string {
	collectorPathConfig := getConfigByKey(collectorPathKey, attrs)
	return getStringConfig(collectorPathEnvVar, defaultCollectorPath, collectorPathConfig)
}

func getEnvironment(attrs []attribute.KeyValue) string {
	environmentConfig := getConfigByKey(environmentKey, attrs)
	return getStringConfig(environmentEnvVar, "", environmentConfig)
}

func getCommitHash(attrs []attribute.KeyValue) string {
	commitHashConfig := getConfigByKey(commitHashKey, attrs)
	return getStringConfig(commitHashEnvVar, "", commitHashConfig)
}

func getObfuscationDetails(attrs []attribute.KeyValue) HeliosObfuscationConfig {
	hsDataObfuscationBlocklist := []string{}
	hsDataObfuscationAllowlist := []string{
		"$.metadata.*",
		"$.collection",
		"$.details[*].name",
		"$.topic",
		"$.information[*].age",
		"$..information[?(@.address=='Unclassified')].address"}
	hsDatahMacKey := "1234"
	if hsDatahMacKey != "" {
		hsDatahMacKeyAsInt, err := strconv.Atoi(hsDatahMacKey)
		if err == nil {
		if len(hsDataObfuscationBlocklist) > 0 {
			return HeliosObfuscationConfig{true, "blocklist", hsDataObfuscationBlocklist, hsDatahMacKeyAsInt}
		} else if len(hsDataObfuscationAllowlist) > 0 {
			return HeliosObfuscationConfig{true, "allowlist", hsDataObfuscationAllowlist, hsDatahMacKeyAsInt}
		}
	}
	}
	return HeliosObfuscationConfig{false, "", []string{}, 0}
}

func getOrCreateHeliosConfig(serviceName string, apiToken string, attrs ...attribute.KeyValue) *HeliosConfig {
	if heliosConfigSingletone != nil {
		return heliosConfigSingletone
	} else {
		sampler := getSampler(attrs)
		collectorInsecure := isCollectorInsecure(attrs)
		collectorEndpoint := getCollectorEndpoint(attrs)
		collectorPath := getCollectorPath(attrs)
		environment := getEnvironment(attrs)
		commitHash := getCommitHash(attrs)
		debug := isDebugMode(attrs)
		metadataOnly := isMetadataOnlyMode(attrs)
		obfuscationConfig := getObfuscationDetails(attrs)
		heliosConfigSingletone = &HeliosConfig{serviceName, apiToken, sampler, collectorInsecure, collectorEndpoint, collectorPath, environment, commitHash, debug, metadataOnly, obfuscationConfig}
		return heliosConfigSingletone
	}
}

func getHeliosConfig() *HeliosConfig {
	if heliosConfigSingletone != nil {
		return heliosConfigSingletone
	}
	return  nil
}
