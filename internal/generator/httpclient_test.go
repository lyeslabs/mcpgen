package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This assumes you have NewGenerator as in your previous code.

func TestGenerateHTTPClient_WithFixture(t *testing.T) {
	// Path to your OpenAPI fixture
	specPath := filepath.Join("../..", "testdata", "simple_openapi.yaml")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Fatalf("Fixture file %s does not exist. Please create it.", specPath)
	}

	packageName := "testpkg"
	outputDir := t.TempDir()

	// Create a real Generator using the fixture
	gen, err := NewGenerator(specPath, true, packageName, outputDir)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	tests := []struct {
		name           string
		includes       []string
		expectError    bool
		expectFile     bool
		errorSubstring string
	}{
		{
			name:       "types only",
			includes:   []string{"types"},
			expectFile: true,
		},
		{
			name:       "httpclient only",
			includes:   []string{"httpclient"},
			expectFile: true,
		},
		{
			name:       "both types and httpclient",
			includes:   []string{"types", "httpclient"},
			expectFile: true,
		},
		{
			name:           "invalid include",
			includes:       []string{"foo"},
			expectError:    true,
			errorSubstring: "no valid includes specified",
		},
		{
			name:           "empty includes",
			includes:       []string{},
			expectError:    true,
			errorSubstring: "no valid includes specified",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up output dir before each run
			_ = os.RemoveAll(filepath.Join(outputDir, "apiclient"))

			err := gen.GenerateHTTPClient(tc.includes)
			outputFile := filepath.Join(outputDir, "apiclient", "client.go")

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.errorSubstring != "" && !strings.Contains(err.Error(), tc.errorSubstring) {
					t.Errorf("Error message %q does not contain expected substring %q", err.Error(), tc.errorSubstring)
				}
				if _, statErr := os.Stat(outputFile); !os.IsNotExist(statErr) {
					t.Errorf("Expected no output file, but file exists: %s", outputFile)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				// Check that the file was created and is not empty
				data, readErr := os.ReadFile(outputFile)
				if readErr != nil {
					t.Fatalf("Expected output file, but got error: %v", readErr)
				}
				if len(data) == 0 {
					t.Errorf("Generated file is empty")
				}
				// Optionally, check for some Go code markers
				if !strings.Contains(string(data), "package apiclient") {
					t.Errorf("Generated file does not contain expected package declaration")
				}
			}
		})
	}
}

// Optionally, test the nil spec error path
func TestGenerateHTTPClient_NilSpec(t *testing.T) {
	outputDir := t.TempDir()
	gen := &Generator{
		spec:      nil,
		outputDir: outputDir,
	}
	err := gen.GenerateHTTPClient([]string{"types"})
	if err == nil || !strings.Contains(err.Error(), "OpenAPI spec is nil") {
		t.Errorf("Expected error about nil spec, got: %v", err)
	}
}
