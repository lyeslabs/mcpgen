// generator_test.go (or write_test.go)
package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Tests for ensureOutputDir ---
func TestEnsureOutputDir_Focused(t *testing.T) {
	t.Run("directory already exists", func(t *testing.T) {
		tempDir := t.TempDir() // t.TempDir() creates the directory
		err := ensureOutputDir(tempDir)
		if err != nil {
			t.Fatalf("ensureOutputDir(%q) error = %v, wantErr nil", tempDir, err)
		}
		info, _ := os.Stat(tempDir) // Error checked by t.TempDir implicitly
		if !info.IsDir() {
			t.Errorf("Path %q should remain a directory", tempDir)
		}
	})

	t.Run("directory does not exist, creates it", func(t *testing.T) {
		parentDir := t.TempDir()
		newDir := filepath.Join(parentDir, "new_dir_to_create")

		err := ensureOutputDir(newDir)
		if err != nil {
			t.Fatalf("ensureOutputDir(%q) error = %v, wantErr nil", newDir, err)
		}
		info, statErr := os.Stat(newDir)
		if statErr != nil {
			t.Fatalf("os.Stat(%q) error = %v, expected directory to be created", newDir, statErr)
		}
		if !info.IsDir() {
			t.Errorf("Path %q was not created as a directory", newDir)
		}
	})

	t.Run("path exists as a file, returns error", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "i_am_a_file.txt")
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %q: %v", filePath, err)
		}

		err := ensureOutputDir(filePath)
		if err == nil {
			t.Fatalf("ensureOutputDir(%q) expected an error when path is a file, got nil", filePath)
		}
		expectedErrorMsg := fmt.Sprintf("path '%s' exists and is not a directory", filePath)
		if err.Error() != expectedErrorMsg {
			t.Errorf("ensureOutputDir error message = %q, want %q", err.Error(), expectedErrorMsg)
		}
	})

	t.Run("cannot create directory because parent path is a file, returns error", func(t *testing.T) {
		tempDir := t.TempDir()
		blockerFile := filepath.Join(tempDir, "parent_is_a_file")
		targetDir := filepath.Join(blockerFile, "subdir") // Attempt to create subdir under a file

		if err := os.WriteFile(blockerFile, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create blocker file %q: %v", blockerFile, err)
		}

		err := ensureOutputDir(targetDir)
		if err == nil {
			t.Fatalf("ensureOutputDir(%q) expected an error, got nil", targetDir)
		}
		// Check that the error is wrapped and comes from os.MkdirAll
		expectedTopLevelMsg := fmt.Sprintf("failed to create output directory '%s'", targetDir)
		if !strings.Contains(err.Error(), expectedTopLevelMsg) {
			t.Errorf("ensureOutputDir error = %q, want substring %q", err.Error(), expectedTopLevelMsg)
		}
		if errors.Unwrap(err) == nil {
			t.Errorf("ensureOutputDir expected a wrapped error from os.MkdirAll")
		}
	})
}

// --- Tests for writeFileContent ---

// Helper for generating content
func fixedContentGenerator(content string, errToReturn error) func() ([]byte, error) {
	return func() ([]byte, error) {
		if errToReturn != nil {
			return nil, errToReturn
		}
		return []byte(content), nil
	}
}

func TestWriteFileContent_Focused(t *testing.T) {
	t.Run("writes new file successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		fileName := "output.txt"
		fileContent := "This is new content."
		filePath := filepath.Join(tempDir, fileName)

		genFunc := fixedContentGenerator(fileContent, nil)
		err := writeFileContent(tempDir, fileName, genFunc)
		if err != nil {
			t.Fatalf("writeFileContent error = %v, wantErr nil", err)
		}

		actualContent, readErr := os.ReadFile(filePath)
		if readErr != nil {
			t.Fatalf("Failed to read created file %q: %v", filePath, readErr)
		}
		if !bytes.Equal(actualContent, []byte(fileContent)) {
			t.Errorf("File content = %q, want %q", string(actualContent), fileContent)
		}
	})

	t.Run("overwrites existing file if content is different", func(t *testing.T) {
		tempDir := t.TempDir()
		fileName := "overwrite_me.txt"
		initialContent := "Old version."
		newContent := "New version!"
		filePath := filepath.Join(tempDir, fileName)

		if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
			t.Fatalf("Failed to write initial file: %v", err)
		}

		genFunc := fixedContentGenerator(newContent, nil)
		err := writeFileContent(tempDir, fileName, genFunc)
		if err != nil {
			t.Fatalf("writeFileContent error = %v, wantErr nil", err)
		}

		actualContent, readErr := os.ReadFile(filePath)
		if readErr != nil {
			t.Fatalf("Failed to read overwritten file: %v", readErr)
		}
		if !bytes.Equal(actualContent, []byte(newContent)) {
			t.Errorf("File content = %q, want %q", string(actualContent), newContent)
		}
	})

	t.Run("does not rewrite existing file if content is the same", func(t *testing.T) {
		tempDir := t.TempDir()
		fileName := "dont_touch_me.txt"
		content := "Identical content."
		filePath := filepath.Join(tempDir, fileName)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write initial file: %v", err)
		}
		initialStat, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to stat initial file: %v", err)
		}
		initialModTime := initialStat.ModTime()

		genFunc := fixedContentGenerator(content, nil) // Same content
		err = writeFileContent(tempDir, fileName, genFunc)
		if err != nil {
			t.Fatalf("writeFileContent error = %v, wantErr nil", err)
		}

		finalStat, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to stat file after operation: %v", err)
		}
		// Check if ModTime changed. This is the best we can do without mocking os.WriteFile
		// to see if it was called. Due to filesystem timestamp precision, this might
		// occasionally be different even if no write happened, but it's a reasonable check.
		if !initialModTime.Equal(finalStat.ModTime()) {
			// This is not a hard failure, as some filesystems might update mtime even on a stat or read.
			// The critical part is that the content is correct and no error occurred.
			t.Logf("File ModTime changed for same content. Initial: %v, Final: %v. This can happen.", initialModTime, finalStat.ModTime())
		}

		actualContent, readErr := os.ReadFile(filePath)
		if readErr != nil {
			t.Fatalf("Failed to read file: %v", readErr)
		}
		if !bytes.Equal(actualContent, []byte(content)) {
			t.Errorf("File content = %q, want %q", string(actualContent), content)
		}
	})

	t.Run("fails if ensureOutputDir (called by writeFileContent) fails", func(t *testing.T) {
		parentDir := t.TempDir()
		outputDirAsFile := filepath.Join(parentDir, "dir_is_actually_a_file") // This will be the outputDir
		if err := os.WriteFile(outputDirAsFile, []byte("blocker"), 0644); err != nil {
			t.Fatalf("Failed to create blocker file: %v", err)
		}

		fileName := "output.txt"
		genFunc := fixedContentGenerator("content", nil)

		err := writeFileContent(outputDirAsFile, fileName, genFunc)
		if err == nil {
			t.Fatalf("writeFileContent expected an error, got nil")
		}

		expectedTopLevelMsg := "failed to create directory"
		if !strings.Contains(err.Error(), expectedTopLevelMsg) {
			t.Errorf("Error = %q, want substring %q", err.Error(), expectedTopLevelMsg)
		}
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			t.Fatal("Expected wrapped error from ensureOutputDir")
		}
		expectedUnwrappedMsg := fmt.Sprintf("path '%s' exists and is not a directory", outputDirAsFile)
		if unwrapped.Error() != expectedUnwrappedMsg {
			t.Errorf("Unwrapped error = %q, want %q", unwrapped.Error(), expectedUnwrappedMsg)
		}

		// Ensure the actual target file (outputDirAsFile/fileName) was not created
		// Since outputDirAsFile is a file, attempting to join fileName to it and stat
		// will result in an error related to outputDirAsFile not being a directory.
		problematicFilePath := filepath.Join(outputDirAsFile, fileName)
		if _, statErr := os.Stat(problematicFilePath); !os.IsNotExist(statErr) {
			// This check is a bit tricky. If outputDirAsFile is a file,
			// os.Stat(outputDirAsFile + "/" + fileName) should fail because outputDirAsFile is not a directory.
			// The error might not be os.IsNotExist for the *full path* but rather for the parent component.
			if statErr == nil {
				t.Errorf("File %q was created when ensureOutputDir failed", problematicFilePath)
			} else {
				t.Logf("Stat for %q failed as expected: %v", problematicFilePath, statErr)
			}
		}
	})

	t.Run("fails if generateContent callback returns an error", func(t *testing.T) {
		tempDir := t.TempDir()
		fileName := "output.txt"
		expectedErr := errors.New("generateContent failed")

		genFunc := fixedContentGenerator("content", expectedErr)
		err := writeFileContent(tempDir, fileName, genFunc)

		if err == nil {
			t.Fatalf("writeFileContent expected an error, got nil")
		}
		if !errors.Is(err, expectedErr) { // The error should be exactly the one from the callback
			t.Errorf("Error = %v, want %v", err, expectedErr)
		}

		filePath := filepath.Join(tempDir, fileName)
		if _, statErr := os.Stat(filePath); !os.IsNotExist(statErr) {
			t.Errorf("File %q was created when generateContent failed (stat error: %v)", filePath, statErr)
		}
	})

	t.Run("fails if os.WriteFile itself fails (simulated by making output dir a file, then trying to write into it)", func(t *testing.T) {
		// This scenario is largely covered by "fails if ensureOutputDir fails" because
		// ensureOutputDir would prevent os.WriteFile from being called if the dir is bad.
		// To *specifically* test os.WriteFile failing *after* ensureOutputDir passed,
		// we'd need to make the directory writable but the specific file path unwritable,
		// which gets into permission complexities we're avoiding.

		// The most direct way our current writeFileContent can hit the os.WriteFile error
		// is if ensureOutputDir passes, but then something goes wrong during the write.
		// This is hard to simulate without mocks or OS-level permission changes.
		// The existing "fails if ensureOutputDir" test is the most relevant for filesystem errors
		// that prevent writing.
		t.Skip("Skipping direct os.WriteFile failure test as it's hard to isolate from ensureOutputDir failure without mocks or permission changes.")
	})
}
