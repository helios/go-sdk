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
	config := getHeliosConfig(serviceName, token)
	assert.Equal(t, config.serviceName, serviceName)
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.collectorEndpoint, defaultCollectorEndpoint)
	assert.Equal(t, config.collectorPath, defaultCollectorPath)
	assert.Equal(t, config.commitHash, "")
	assert.Equal(t, config.environment, "")
	assert.Equal(t, config.sampler.Description(), "AlwaysOnSampler")
}

func TestConfigWithOptions(t *testing.T) {
	testCollectorEndpoint := "aaa.bbb.com:1234"
	testCollectorPath := "/sababa"
	testSamplingRatio := 0.1234
	config := getHeliosConfig(serviceName, token, WithCollectorEndpoint(testCollectorEndpoint), WithCollectorPath(testCollectorPath), WithSamplingRatio(testSamplingRatio))
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.collectorEndpoint, testCollectorEndpoint)
	assert.Equal(t, config.collectorPath, testCollectorPath)
	assert.Equal(t, config.sampler.Description(), fmt.Sprintf("TraceIDRatioBased{%.4f}", testSamplingRatio))
}

func TestConfigWithEnvVars(t *testing.T) {
	testCollectorEndpoint := "aaa.bbb.com:1234"
	testCollectorPath := "/sababa"
	testSamplingRatio := 0.1234
	os.Setenv(collectorEndpointEnvVar, testCollectorEndpoint)
	os.Setenv(collectorPathEnvVar, testCollectorPath)
	os.Setenv(samplingRatioEnvVar, fmt.Sprintf("%.4f", testSamplingRatio))

	config := getHeliosConfig(serviceName, token)
	assert.Equal(t, config.apiToken, token)
	assert.Equal(t, config.collectorEndpoint, testCollectorEndpoint)
	assert.Equal(t, config.collectorPath, testCollectorPath)
	assert.Equal(t, config.sampler.Description(), fmt.Sprintf("TraceIDRatioBased{%.4f}", testSamplingRatio))
}
