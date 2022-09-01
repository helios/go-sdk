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
const collectorEndpointKey = "collectorEndpoint"
const collectorEndpointEnvVar = "HS_COLLECTOR_ENDPOINT"
const collectorPathKey = "collectorPath"
const collectorPathEnvVar = "HS_COLLECTOR_PATH"
const commitHashKey = "commitHash"
const commitHashEnvVar = "HS_COMMIT_HASH"

// Default values
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

func getSampler(samplerConfig attribute.KeyValue) trace.Sampler {
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

func getCollectorEndpoint(collectorEndpointConfig attribute.KeyValue) string {
	collectorEndpoint := os.Getenv(collectorEndpointEnvVar)
	if collectorEndpoint != "" {
		return collectorEndpoint
	}

	if collectorEndpointConfig.Key == "" {
		return defaultCollectorEndpoint
	}

	return collectorEndpointConfig.Value.AsString()
}

func getCollectorPath(collectorPathConfig attribute.KeyValue) string {
	collectorPath := os.Getenv(collectorPathEnvVar)
	if collectorPath != "" {
		return collectorPath
	}

	if collectorPathConfig.Key == "" {
		return defaultCollectorPath
	}

	return collectorPathConfig.Value.AsString()
}

func getEnvironment(environmentConfig attribute.KeyValue) string {
	environment := os.Getenv(environmentEnvVar)
	if environment != "" {
		return environment
	}

	if environmentConfig.Key == "" {
		return ""
	}

	return environmentConfig.Value.AsString()
}

func getCommitHash(commitHashConfig attribute.KeyValue) string {
	commitHash := os.Getenv(commitHashEnvVar)
	if commitHash != "" {
		return commitHash
	}

	if commitHashConfig.Key == "" {
		return ""
	}

	return commitHashConfig.Value.AsString()
}

func getHeliosConfig(serviceName string, apiToken string, attrs ...attribute.KeyValue) HeliosConfig {
	sampler := getSampler(getConfigByKey(samplingRatioKey, attrs))
	collectorEndpoint := getCollectorEndpoint(getConfigByKey(collectorEndpointKey, attrs))
	collectorPath := getCollectorPath(getConfigByKey(collectorPathKey, attrs))
	environment := getEnvironment(getConfigByKey(environmentKey, attrs))
	commitHash := getCommitHash(getConfigByKey(commitHashKey, attrs))
	return HeliosConfig{serviceName, apiToken, sampler, collectorEndpoint, collectorPath, environment, commitHash}
}
