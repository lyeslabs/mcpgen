package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lyeslabs/mcpgen/internal/converter"
)

// testConverter implements a minimal converter.Converter interface for testing.
type testConverter struct {
	config *converter.MCPConfig
}

func (tc *testConverter) Convert() (*converter.MCPConfig, error) {
	return tc.config, nil
}

func TestGenerateMCP(t *testing.T) {
	tmpDir := t.TempDir()

	// Prepare a minimal MCPConfig with one tool
	config := &converter.MCPConfig{
		Tools: []converter.Tool{
			{
				Name:           "echo",
				Description:    "Echoes input",
				RawInputSchema: `{"type":"object","properties":{"msg":{"type":"string"}}}`,
				Responses: []converter.ResponseTemplate{
					{PrependBody: "// response", StatusCode: 200, ContentType: "application/json", Suffix: "default"},
				},
				RequestTemplate: converter.RequestTemplate{
					URL:    "/echo",
					Method: "POST",
				},
			},
		},
	}

	// Use the test converter
	g := &Generator{
		PackageName: "mytools",
		outputDir:   tmpDir,
		converter:   &testConverter{config: config},
	}

	// Call GenerateMCP
	if err := g.GenerateMCP(); err != nil {
		t.Fatalf("GenerateMCP failed: %v", err)
	}

	// Check that server.go exists
	serverGoPath := filepath.Join(tmpDir, "server.go")
	if _, err := os.Stat(serverGoPath); err != nil {
		t.Errorf("Expected server.go to be generated, but it does not exist")
	}

	// Check that helpers.go exists
	helpersGoPath := filepath.Join(tmpDir, "helpers", "params.go")
	if _, err := os.Stat(helpersGoPath); err != nil {
		t.Errorf("Expected helpers.go to be generated, but it does not exist")
	}

	// Check that the tool file exists
	toolFilePath := filepath.Join(tmpDir, "mcptools", "Echo.go")
	data, err := os.ReadFile(toolFilePath)
	if err != nil {
		t.Fatalf("Failed to read generated tool file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "EchoHandler") {
		t.Errorf("Generated tool file missing handler name")
	}
	if !strings.Contains(content, "Echoes input") {
		t.Errorf("Generated tool file missing tool description")
	}
	if !strings.Contains(content, "package mcptools") {
		t.Errorf("Generated tool file missing package declaration")
	}
}
