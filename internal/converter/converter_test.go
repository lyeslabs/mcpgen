package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConverter(t *testing.T) {
	parser := NewParser(false)
	c := NewConverter(parser)
	if c == nil {
		t.Fatal("expected non-nil Converter")
	}
	if c.parser != parser {
		t.Error("expected parser to be set")
	}
	if c.options.ServerConfig == nil {
		t.Error("expected ServerConfig to be initialized")
	}
}

func TestConverter_Convert(t *testing.T) {
	// Load a real OpenAPI spec
	specPath := filepath.Join("..", "testdata", "simple_openapi.yaml")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Fatalf("Test setup error: fixture file %s does not exist. Please create it.", specPath)
	}
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("could not read %s: %v", specPath, err)
	}
	parser := NewParser(false)
	if err := parser.Parse(data); err != nil {
		t.Fatalf("failed to parse OpenAPI: %v", err)
	}

	c := NewConverter(parser)
	config, err := c.Convert()
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	if config == nil {
		t.Fatal("expected non-nil MCPConfig")
	}
	if config.Server.Config == nil {
		t.Error("expected Server.Config to be set")
	}
	if len(config.Tools) == 0 {
		t.Error("expected at least one tool in Tools")
	}
	// Check that tools are sorted by name
	for i := 1; i < len(config.Tools); i++ {
		if config.Tools[i-1].Name > config.Tools[i].Name {
			t.Errorf("tools not sorted by name: %q > %q", config.Tools[i-1].Name, config.Tools[i].Name)
		}
	}
}

func TestConverter_Convert_NoDocument(t *testing.T) {
	parser := NewParser(false)
	c := NewConverter(parser)
	_, err := c.Convert()
	if err == nil {
		t.Fatal("expected error when no OpenAPI document is loaded")
	}
}
