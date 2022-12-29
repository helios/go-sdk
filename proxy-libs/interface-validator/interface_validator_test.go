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
	sort.Slice(heliosExports, func(i int, j int) bool { return heliosExports[i].Name < heliosExports[j].Name })
	return heliosExports
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

func TestSaramaInterfaceMatch(t *testing.T) {
	delete := func(exports []exportsExtractor.ExtractedObject, name string) []exportsExtractor.ExtractedObject {
		for i, export := range exports {
			if export.Name == name {
				return append(exports[:i], exports[i+1:]...)
			}
		}

		return exports
	}

	originalExports := cloneRepositoryAndExtractExports("https://github.com/Shopify/sarama", "v1.37.2", "sarama", "")
	heliosExports := extractProxyLibExports("heliossarama")

	// "NewMockWrapper" cannot be wrapped because its parameter's type is private - Remove it from the expected list.
	originalExports = delete(originalExports, "NewMockWrapper")
	// The signature of "Wrap" was changed because the original return type is private - Remove it from both lists.
	originalExports = delete(originalExports, "Wrap")
	heliosExports = delete(heliosExports, "Wrap")
	// A helper method we've added to improve context propagation
	heliosExports = delete(heliosExports, "InjectContextToMessage")

	assert.Equal(t, len(originalExports), len(heliosExports))
	assert.EqualValues(t, originalExports, heliosExports)
}