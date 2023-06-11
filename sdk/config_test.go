package sdk

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const serviceName = "test_service"
const token = "abcd1234"

func TestBasicConfig(t *testing.T) {
	heliosConfigSingleton = nil
	config := createHeliosConfig(serviceName, token)
	assert.Equal(t, config.serviceName, serviceName)
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.collectorEndpoint, defaultCollectorEndpoint)
	assert.Equal(t, config.collectorPath, defaultCollectorPath)
	assert.Equal(t, config.collectorMetricsPath, defaultCollectorMetricsPath)
	assert.Equal(t, config.commitHash, "")
	assert.Equal(t, config.environment, "")
	assert.Equal(t, config.serviceNamespace, "")
	assert.Equal(t, config.debug, false)
	assert.Equal(t, config.metadataOnly, false)
	assert.Equal(t, config.disableMetricsCollection, false)
	assert.Equal(t, config.sampler.Description(), "AlwaysOnSampler")
}

func TestConfigWithOptions(t *testing.T) {
	heliosConfigSingleton = nil
	testCollectorEndpoint := "aaa.bbb.com:1234"
	testCollectorPath := "/sababa"
	testSamplingRatio := 0.1234
	metricsPath := "/v2/metrics"
	serviceNamespace := "ns1"
	config := createHeliosConfig(serviceName, token, WithCollectorInsecure(), WithCollectorEndpoint(testCollectorEndpoint), WithCollectorPath(testCollectorPath), WithSamplingRatio(testSamplingRatio), WithDebugMode(), WithMetadataOnlyMode(), WithDisableMetricsCollection(), WithCollectorMetricsPath(metricsPath), WithServiceNamespace(serviceNamespace))
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.collectorInsecure, true)
	assert.Equal(t, config.collectorEndpoint, testCollectorEndpoint)
	assert.Equal(t, config.collectorPath, testCollectorPath)
	assert.Equal(t, config.sampler.Description(), fmt.Sprintf("HeliosSampler(%.4f)", testSamplingRatio))
	assert.Equal(t, config.debug, true)
	assert.Equal(t, config.metadataOnly, true)
	assert.Equal(t, config.collectorMetricsPath, metricsPath)
	assert.Equal(t, config.serviceNamespace, serviceNamespace)
}

func TestConfigWithDisabledInstrumentationOption(t *testing.T) {
	heliosConfigSingleton = nil
	config := createHeliosConfig(serviceName, token, WithInstrumentationDisabled())
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.instrumentationDisabled, true)
}

func TestConfigWithEnvVars(t *testing.T) {
	heliosConfigSingleton = nil
	testCollectorEndpoint := "aaa.bbb.com:1234"
	testCollectorPath := "/sababa"
	testSamplingRatio := 0.1234
	metricsPath := "/v2/metrics"
	serviceNamespace := "ns1"
	os.Setenv(collectorInsecureEnvVar, "true")
	os.Setenv(collectorEndpointEnvVar, testCollectorEndpoint)
	os.Setenv(collectorPathEnvVar, testCollectorPath)
	os.Setenv(samplingRatioEnvVar, fmt.Sprintf("%.4f", testSamplingRatio))
	os.Setenv(debugEnvVar, "true")
	os.Setenv(metadataOnlyEnvVar, "true")
	os.Setenv(disableMetricsCollectionEnvVar, "true")
	os.Setenv(collectorMetricsPathEnvVar, metricsPath)
	os.Setenv(serviceNamespaceEnvVar, serviceNamespace)

	config := createHeliosConfig(serviceName, token)
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.collectorInsecure, true)
	assert.Equal(t, config.collectorEndpoint, testCollectorEndpoint)
	assert.Equal(t, config.collectorPath, testCollectorPath)
	assert.Equal(t, config.sampler.Description(), fmt.Sprintf("HeliosSampler(%.4f)", testSamplingRatio))
	assert.Equal(t, config.debug, true)
	assert.Equal(t, config.metadataOnly, true)
	assert.Equal(t, config.disableMetricsCollection, true)
	assert.Equal(t, config.collectorMetricsPath, metricsPath)
	assert.Equal(t, config.serviceNamespace, serviceNamespace)
}

func TestConfigWithDisabledInstrumentationEnvVar(t *testing.T) {
	heliosConfigSingleton = nil
	os.Setenv(instrumentationDisabledEnvVar, "true")
	defer os.Setenv(instrumentationDisabledEnvVar, "")

	config := createHeliosConfig(serviceName, token)
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.instrumentationDisabled, true)
}
