package generator

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create a temporary spec file
func createTempSpecFileWithContent(t *testing.T, content string) string {
	t.Helper()
	tempFile, err := os.CreateTemp(t.TempDir(), "testspec-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp spec file: %v", err)
	}
	if _, err := tempFile.WriteString(content); err != nil {
		tempFile.Close()
		t.Fatalf("Failed to write to temp spec file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp spec file: %v", err)
	}
	return tempFile.Name()
}
func TestNewGenerator_Success_WithValidFile(t *testing.T) {
	// Assumes testdata/valid_openapi.yaml exists relative to this test file
	specPath := filepath.Join("../..", "testdata", "simple_openapi.yaml")

	// Check if the fixture file exists
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Fatalf("Test setup error: fixture file %s does not exist. Please create it.", specPath)
	}

	packageName := "testpkgfromfile"
	outputDir := t.TempDir()

	testCases := []struct {
		name       string
		validation bool
	}{
		{"with validation", true},
		{"without validation", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := NewGenerator(specPath, tc.validation, packageName, outputDir)

			if err != nil {
				t.Fatalf("NewGenerator() with file %s error = %v, wantErr nil", specPath, err)
			}
			if gen == nil {
				t.Fatal("NewGenerator() returned nil Generator, want non-nil")
			}

			if gen.specPath != specPath {
				t.Errorf("gen.specPath = %q, want %q", gen.specPath, specPath)
			}
			if gen.PackageName != packageName {
				t.Errorf("gen.PackageName = %q, want %q", gen.PackageName, packageName)
			}
			if gen.outputDir != outputDir {
				t.Errorf("gen.outputDir = %q, want %q", gen.outputDir, outputDir)
			}

			if gen.converter == nil {
				t.Error("gen.converter is nil, want non-nil")
			}
			if gen.spec == nil {
				t.Error("gen.spec is nil, want non-nil")
			} else {
				// Check a value from the testdata/valid_openapi.yaml file
				if gen.spec.Info == nil {
					t.Errorf("gen.spec.Info.Title not parsed correctly from file, got '%s', want 'My Valid API From File'", gen.spec.Info.Title)
				}
			}
		})
	}
}

func TestNewGenerator_Error_ParsingFailed(t *testing.T) {
	packageName := "testpkgerror"
	outputDir := t.TempDir()

	testCases := []struct {
		name          string
		specPathSetup func(t *testing.T) string // Returns path to use
		validation    bool
		expectedError string // The top-level error string we wrap with
	}{
		{
			name: "non-existent file",
			specPathSetup: func(t *testing.T) string {
				// Use a path within the temp directory that won't exist
				return filepath.Join(t.TempDir(), "this_file_does_not_exist.yaml")
			},
			validation:    false,
			expectedError: "error parsing OpenAPI specification",
		},
		{
			name: "invalid spec content in file",
			specPathSetup: func(t *testing.T) string {
				return createTempSpecFileWithContent(t, "this is not: valid: yaml { content")
			},
			validation:    false,
			expectedError: "error parsing OpenAPI specification",
		},
		{
			name: "invalid openapi structure with validation",
			specPathSetup: func(t *testing.T) string {
				// Valid YAML, but invalid OpenAPI (missing info.version)
				return createTempSpecFileWithContent(t, "openapi: 3.0.0\ninfo:\n  title: API Without Version")
			},
			validation:    true,
			expectedError: "error parsing OpenAPI specification",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			specPath := tc.specPathSetup(t)
			gen, err := NewGenerator(specPath, tc.validation, packageName, outputDir)

			if err == nil {
				t.Fatal("NewGenerator() error = nil, wantErr non-nil")
			}
			if gen != nil {
				t.Errorf("NewGenerator() returned non-nil Generator (%+v), want nil on error", gen)
			}

			if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("NewGenerator() error message = %q, want substring %q", err.Error(), tc.expectedError)
			}

			// This assertion REQUIRES NewGenerator to use fmt.Errorf("...: %w", err)
			// If you haven't made that change in generator.go, this will fail.
			unwrappedErr := errors.Unwrap(err)
			if unwrappedErr == nil {
				t.Errorf("Expected a wrapped error (NewGenerator should use %%w to wrap), but unwrapping returned nil. Original error: %v", err)
			} else {
				t.Logf("Successfully unwrapped error: %v", unwrappedErr)
			}
		})
	}
}
