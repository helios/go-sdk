package interfacevalidator

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func assertExportsEquality(t *testing.T, proxyExports []exportsExtractor.ExtractedObject, originalExports []exportsExtractor.ExtractedObject) {
	for _, proxyExport := range proxyExports {
		if proxyExport.Name == "InstrumentedSymbols" {
			continue
		}

		var found bool = false
		for _, originalExport := range originalExports {
			if proxyExport.Name == originalExport.Name {
				found = true
				assertExtractedObjectEquality(t, proxyExport, originalExport)
			}
		}

		assert.True(t, found)
	}
}

func assertFunctionParameterEquality(t *testing.T, heliosAttribute exportsExtractor.ObjectAttribute, originalAttribute exportsExtractor.ObjectAttribute) {
	assert.True(t, strings.HasSuffix(heliosAttribute.AttributeType, originalAttribute.AttributeType))
	assert.True(t, strings.HasSuffix(heliosAttribute.AttributeTypeKey, originalAttribute.AttributeTypeKey))
}

func assertExtractedObjectEquality(t *testing.T, heliosObject exportsExtractor.ExtractedObject, originalObject exportsExtractor.ExtractedObject) {
	assert.Equal(t, heliosObject.Name, originalObject.Name)
	assert.Equal(t, heliosObject.PackageAttributeType, originalObject.PackageAttributeType)
	for index, proxyInput := range heliosObject.FunctionAttributeInput {
		originalInput := originalObject.FunctionAttributeInput[index]
		assertFunctionParameterEquality(t, proxyInput, originalInput)
	}

	for index, proxyInput := range heliosObject.FunctionReturnValues {
		originalInput := originalObject.FunctionReturnValues[index]
		assertFunctionParameterEquality(t, proxyInput, originalInput)
	}
}

func deleteExportedMember(exports []exportsExtractor.ExtractedObject, name string) []exportsExtractor.ExtractedObject {
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
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestGrpcInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/grpc/grpc-go", "v1.50.1", "grpc", "")
	heliosExports := extractProxyLibExports("heliosgrpc")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestMongoInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/mongodb/mongo-go-driver", "v1.11.0", "mongo", "/mongo")
	heliosExports := extractProxyLibExports("heliosmongo")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestMuxInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/gorilla/mux", "v1.8.0", "mux", "")
	heliosExports := extractProxyLibExports("heliosmux")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestEchoInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/labstack/echo", "v4.9.1", "echo", "")
	heliosExports := extractProxyLibExports("heliosecho")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestMacaronInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/go-macaron/macaron", "v1.4.0", "macaron", "")
	heliosExports := extractProxyLibExports("heliosmacaron")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestGinInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/gin-gonic/gin", "v1.8.1", "gin", "")
	heliosExports := extractProxyLibExports("heliosgin")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestChiInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/go-chi/chi", "v5.0.8", "chi", "")
	heliosExports := extractProxyLibExports("helioschi")
	assertExportsEquality(t, heliosExports, originalExports)
}

func TestSaramaInterfaceMatch(t *testing.T) {
	originalExports := cloneRepositoryAndExtractExports("https://github.com/Shopify/sarama", "v1.37.2", "sarama", "")
	heliosExports := deleteExportedMember(extractProxyLibExports("heliossarama"), "InjectContextToMessage")
	assertExportsEquality(t, heliosExports, originalExports)
}
