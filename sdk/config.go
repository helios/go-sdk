package sdk

import (
	"os"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

var heliosConfigSingleton *HeliosConfig

type HeliosConfig struct {
	instrumentationDisabled  bool
	serviceName              string
	apiToken                 string
	sampler                  trace.Sampler
	collectorInsecure        bool
	collectorEndpoint        string
	collectorPath            string
	collectorMetricsPath     string
	environment              string
	serviceNamespace         string
	commitHash               string
	debug                    bool
	metadataOnly             bool
	disableMetricsCollection bool
}

// Keys and their matching env vars
const instrumentationDisabledKey = "disabled"
const instrumentationDisabledEnvVar = "HS_DISABLED"
const samplingRatioKey = "samplingRatio"
const samplingRatioEnvVar = "HS_SAMPLING_RATIO"
const environmentKey = "environment"
const environmentEnvVar = "HS_ENVIRONMENT"
const serviceNamespaceKey = "serviceNamespace"
const serviceNamespaceEnvVar = "HS_SERVICE_NAMESPACE"
const collectorInsecureKey = "collectorInsecure"
const collectorInsecureEnvVar = "HS_COLLECTOR_INSECURE"
const collectorEndpointKey = "collectorEndpoint"
const collectorEndpointEnvVar = "HS_COLLECTOR_ENDPOINT"
const collectorPathKey = "collectorPath"
const collectorPathEnvVar = "HS_COLLECTOR_PATH"
const collectorMetricsPathKey = "collectorMetricsPath"
const collectorMetricsPathEnvVar = "HS_COLLECTOR_METRICS_PATH"
const commitHashKey = "commitHash"
const commitHashEnvVar = "HS_COMMIT_HASH"
const debugKey = "debug"
const debugEnvVar = "HS_DEBUG"
const metadataOnlyKey = "metadataOnly"
const metadataOnlyEnvVar = "HS_METADATA_ONLY"
const disableMetricsCollectionKey = "disableMetricsCollection"
const disableMetricsCollectionEnvVar = "HS_DISABLE_METRICS_COLLECTION"
const hsDataObfuscationAllowlistKey = "dataObfuscationAllowlist"
const hsDataObfuscationBlocklistEnvVar = "HS_DATA_OBFUSCATION_BLOCKLIST"
const hsDataObfuscationBlocklistKey = "dataObfuscationBlocklist"
const hsDatahMacKeyEnvVar = "HS_DATA_OBFUSCATION_HMAC_KEY"
const hsDatahMacKey = "dataObfuscationhMacKey"

// Default values
const defaultInstrumentationDisabled = false
const defaultCollectorInsecure = false
const defaultCollectorEndpoint = "collector.gethelios.dev:443"
const defaultCollectorPath = "/v1/traces"
const defaultCollectorMetricsPath = "/v1/metrics"
const defaultDebug = false
const defaultMetadataOnly = false
const defaultDisableMetricsCollection = false

func getConfigByKey(key string, attrs []attribute.KeyValue) attribute.KeyValue {
	for i := range attrs {
		if string(attrs[i].Key) == key {
			return attrs[i]
		}
	}

	return attribute.KeyValue{Key: "", Value: attribute.Value{}}
}

func isInstrumentationDisabled(attrs []attribute.KeyValue) bool {
	instrumentationDisabledConfig := getConfigByKey(instrumentationDisabledKey, attrs)
	return getBoolConfig(instrumentationDisabledEnvVar, defaultInstrumentationDisabled, instrumentationDisabledConfig)
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

func isDisableMetricsCollection(attrs []attribute.KeyValue) bool {
	disableMetricsCollectionConfig := getConfigByKey(disableMetricsCollectionKey, attrs)
	return getBoolConfig(disableMetricsCollectionEnvVar, defaultDisableMetricsCollection, disableMetricsCollectionConfig)
}

func getCollectorEndpoint(attrs []attribute.KeyValue) string {
	collectorEndpointConfig := getConfigByKey(collectorEndpointKey, attrs)
	return getStringConfig(collectorEndpointEnvVar, defaultCollectorEndpoint, collectorEndpointConfig)
}

func getCollectorPath(attrs []attribute.KeyValue) string {
	collectorPathConfig := getConfigByKey(collectorPathKey, attrs)
	return getStringConfig(collectorPathEnvVar, defaultCollectorPath, collectorPathConfig)
}

func getCollectorMetricsPath(attrs []attribute.KeyValue) string {
	collectorMetricsPathConfig := getConfigByKey(collectorMetricsPathKey, attrs)
	return getStringConfig(collectorMetricsPathEnvVar, defaultCollectorMetricsPath, collectorMetricsPathConfig)
}

func getEnvironment(attrs []attribute.KeyValue) string {
	environmentConfig := getConfigByKey(environmentKey, attrs)
	return getStringConfig(environmentEnvVar, "", environmentConfig)
}

func getServiceNamespace(attrs []attribute.KeyValue) string {
	serviceNamespaceConfig := getConfigByKey(serviceNamespaceKey, attrs)
	return getStringConfig(serviceNamespaceEnvVar, "", serviceNamespaceConfig)
}

func getCommitHash(attrs []attribute.KeyValue) string {
	commitHashConfig := getConfigByKey(commitHashKey, attrs)
	return getStringConfig(commitHashEnvVar, "", commitHashConfig)
}

func createHeliosConfig(serviceName string, apiToken string, attrs ...attribute.KeyValue) *HeliosConfig {
	if heliosConfigSingleton != nil {
		return heliosConfigSingleton
	} else {
		instrumentationDisabled := isInstrumentationDisabled(attrs)
		sampler := getSampler(attrs)
		collectorInsecure := isCollectorInsecure(attrs)
		collectorEndpoint := getCollectorEndpoint(attrs)
		collectorPath := getCollectorPath(attrs)
		collectorMetricsPath := getCollectorMetricsPath(attrs)
		environment := getEnvironment(attrs)
		serviceNamespace := getServiceNamespace(attrs)
		commitHash := getCommitHash(attrs)
		debug := isDebugMode(attrs)
		metadataOnly := isMetadataOnlyMode(attrs)
		disableMetricsCollection := isDisableMetricsCollection(attrs)
		heliosConfigSingleton = &HeliosConfig{instrumentationDisabled, serviceName, apiToken, sampler, collectorInsecure, collectorEndpoint, collectorPath, collectorMetricsPath, environment, serviceNamespace, commitHash, debug, metadataOnly, disableMetricsCollection}
		return heliosConfigSingleton
	}
}
