package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateHelpers_NoErrorOnSuccess(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "generator_helpers_test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	// Instantiate the Generator with the temporary output directory
	g := &Generator{
		PackageName: "mytools", // Use a specific package name for the test
		outputDir:   tmpDir,
	}

	// Call the function under test
	err = g.GenerateHelpers()

	// Assert that no error occurred
	if err != nil {
		t.Errorf("GenerateHelpers returned an unexpected error: %v", err)
	}

	// Optional: You could still check for the *existence* of the file
	// to ensure the writeFileContent call was at least attempted.
	expectedFilePath := filepath.Join(tmpDir, "helpers", "params.go")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Errorf("expected generated file %s to exist, but it does not", expectedFilePath)
	}
}
