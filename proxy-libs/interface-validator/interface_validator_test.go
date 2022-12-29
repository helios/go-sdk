package interfacevalidator

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	exportsExtractor "github.com/helios/go-instrumentor/exports_extractor"
	"github.com/stretchr/testify/assert"
)

func TestSaramaInterfaceMatch(t *testing.T) {
	delete := func(exports []exportsExtractor.ExtractedObject, name string) []exportsExtractor.ExtractedObject {
		for i, export := range exports {
			if export.Name == name {
				return append(exports[:i], exports[i+1:]...)
			}
		}

		return exports
	}

	// Get original sarama exports.
	originalRepository := exportsExtractor.CloneGitRepository("https://github.com/Shopify/sarama", "v1.37.2")
	originalExports := exportsExtractor.ExtractExports(originalRepository, "sarama")
	os.RemoveAll(originalRepository)
	sort.Slice(originalExports, func(i int, j int) bool { return originalExports[i].Name < originalExports[j].Name })

	// Get Helios sarama exports.
	srcDir, _ := filepath.Abs("../heliossarama")
	heliosExports := exportsExtractor.ExtractExports(srcDir, "heliossarama")
	sort.Slice(heliosExports, func(i int, j int) bool { return heliosExports[i].Name < heliosExports[j].Name })

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
