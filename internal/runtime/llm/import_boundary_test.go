package llm

import (
	"os"
	"strings"
	"testing"
)

func TestRuntimeLLM_ShouldNotImportProviderDomain(t *testing.T) {
	for _, file := range []string{"factory.go", "executor.go", "model_config.go"} {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		if strings.Contains(string(content), "internal/modules/provider/domain") {
			t.Fatalf("%s imports provider domain", file)
		}
	}
}
