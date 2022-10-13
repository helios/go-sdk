package sdk

import (
	"os"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type HeliosConfig struct {
	serviceName       string
	apiToken          string
	sampler           trace.Sampler
	collectorInsecure bool
	collectorEndpoint string
	collectorPath     string
	environment       string
	commitHash        string
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

// Default values
const defaultCollectorInsecure = "false"
const defaultCollectorEndpoint = "collector.heliosphere.io:443"
const defaultCollectorPath = "traces"

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

		return trace.TraceIDRatioBased(res)
	}

	if samplerConfig.Key == "" {
		return trace.AlwaysSample()
	}

	return trace.TraceIDRatioBased(samplerConfig.Value.AsFloat64())
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

func isCollectorInsecure(attrs []attribute.KeyValue) bool {
	collectorInsecureConfig := getConfigByKey(collectorInsecureKey, attrs)
	bool, _ := strconv.ParseBool(getStringConfig(collectorInsecureEnvVar, defaultCollectorInsecure, collectorInsecureConfig))
	return bool
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

func getHeliosConfig(serviceName string, apiToken string, attrs ...attribute.KeyValue) HeliosConfig {
	sampler := getSampler(attrs)
	collectorInsecure := isCollectorInsecure(attrs)
	collectorEndpoint := getCollectorEndpoint(attrs)
	collectorPath := getCollectorPath(attrs)
	environment := getEnvironment(attrs)
	commitHash := getCommitHash(attrs)
	return HeliosConfig{serviceName, apiToken, sampler, collectorInsecure, collectorEndpoint, collectorPath, environment, commitHash}
}
