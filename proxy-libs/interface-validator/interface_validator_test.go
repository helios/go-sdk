package interfacevalidator

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	"github.com/stretchr/testify/assert"
)

func sortExports(exports []exportsExtractor.ExtractedObject) {
	sort.Slice(exports, func(i int, j int) bool { return exports[i].Name < exports[j].Name })
}

func cloneRepositoryAndExtractExports(repoUrl string, tag string, moduleName string, modulePath string) []exportsExtractor.ExtractedObject {
	originalRepository := exportsExtractor.CloneGitRepository(repoUrl, tag)
	defer os.RemoveAll(originalRepository)
	originalExports := exportsExtractor.ExtractExports(originalRepository+modulePath, moduleName)
	sortExports(originalExports)
	return originalExports
}

func extractProxyLibExports(libName string) []exportsExtractor.ExtractedObject {
	srcDir, _ := filepath.Abs("../" + libName)
	heliosExports := exportsExtractor.ExtractExports(srcDir, libName)
	heliosExports = deleteByName(heliosExports, "InstrumentedSymbols")
	sortExports(heliosExports)
	return heliosExports
}

func deleteByName(exports []exportsExtractor.ExtractedObject, name string) []exportsExtractor.ExtractedObject {
	for i, export := range exports {
		if export.Name == name {
			return append(exports[:i], exports[i+1:]...)
		}
	}

	return exports
}

func TestHttpInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/golang/go", "go1.19", "http", "/src/net/http")
	heliosExports := extractProxyLibExports("helioshttp")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestGrpcInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/grpc/grpc-go", "v1.50.1", "grpc", "")
	heliosExports := extractProxyLibExports("heliosgrpc")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestMongoInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/mongodb/mongo-go-driver", "v1.11.0", "mongo", "/mongo")
	heliosExports := extractProxyLibExports("heliosmongo")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestMuxInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/gorilla/mux", "v1.8.0", "mux", "")
	heliosExports := extractProxyLibExports("heliosmux")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestEchoInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/labstack/echo", "v4.9.1", "echo", "")
	heliosExports := extractProxyLibExports("heliosecho")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestS3InterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/aws/aws-sdk-go-v2", "service/s3/v1.30.0", "s3", "/service/s3")
	heliosExports := extractProxyLibExports("helioss3")
	originalExports = deleteByName(originalExports,"NewDefaultEndpointResolver")
	heliosExports = deleteByName(heliosExports,"NewDefaultEndpointResolver")
	assert.EqualValues(t, originalExports, heliosExports)
}

func removeDuplicateValues(slice []exportsExtractor.ExtractedObject) []exportsExtractor.ExtractedObject {
	keys := make(map[string]bool)
	list := []exportsExtractor.ExtractedObject{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry.Name]; !value {
			keys[entry.Name] = true
			list = append(list, entry)
		}
	}
	return list
}

func TestMacaronInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/go-macaron/macaron", "v1.4.0", "macaron", "")
	heliosExports := extractProxyLibExports("heliosmacaron")
	// Macaron has separate implementations for PathUnescape for Go 1.17 and 1.18 - until the extractor properly
	// handles that we're forced to remove duplicates
	originalExports = removeDuplicateValues(originalExports)

	for index, value := range originalExports {
		heliosVal := heliosExports[index]
		if value.Name == "NewRouteMap" {
			// The return value can't be used by the proxy lib as its not exported by the original package
			assert.Equal(t, value.FunctionReturnValues[0].AttributeType, "routeMap")
			assert.Equal(t, heliosVal.FunctionReturnValues[0].AttributeType, "interface{}")
		} else {
			assert.Equal(t, value, heliosVal)
		}
	}
}

func TestGinInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/gin-gonic/gin", "v1.8.1", "gin", "")
	heliosExports := extractProxyLibExports("heliosgin")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/go-chi/chi", "v5.0.8", "chi", "")
	heliosExports := extractProxyLibExports("helioschi")
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestSaramaInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/Shopify/sarama", "v1.37.2", "sarama", "")
	heliosExports := extractProxyLibExports("heliossarama")

	// "NewMockWrapper" cannot be wrapped because its parameter's type is private - Remove it from the expected list.
	originalExports = deleteByName(originalExports, "NewMockWrapper")
	// The signature of "Wrap" was changed because the original return type is private - Remove it from both lists.
	originalExports = deleteByName(originalExports, "Wrap")
	heliosExports = deleteByName(heliosExports, "Wrap")
	// A helper method we've added to improve context propagation
	heliosExports = deleteByName(heliosExports, "InjectContextToMessage")

	assert.Equal(t, len(originalExports), len(heliosExports))
	assert.EqualValues(t, originalExports, heliosExports)
}

func TestAwsLambdaInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/aws/aws-lambda-go", "v1.37.0", "lambda", "/lambda")
	heliosExports := extractProxyLibExports("helioslambda")

	// Generics, not supported
	originalExports = deleteByName(originalExports, "HandlerFunc")
	originalExports = deleteByName(originalExports, "StartHandlerFunc")

	assert.EqualValues(t, originalExports, heliosExports)
}
