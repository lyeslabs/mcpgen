package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// File utility functions
func writeFileContent(outputDir, fileName string, generateContent func() ([]byte, error)) error {
	if err := ensureOutputDir(outputDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(outputDir, fileName)
	existingContent, readErr := os.ReadFile(filePath)

	newContent, err := generateContent()
	if err != nil {
		return err
	}

	contentEqual := readErr == nil && bytes.Equal(existingContent, newContent)

	if !contentEqual {
		if err := os.WriteFile(filePath, newContent, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", fileName, err)
		}
	}

	return nil
}


func ensureOutputDir(dir string) error {
    info, err := os.Stat(dir)
    if err == nil {
        // Path exists
        if info.IsDir() {
            // It's a directory, all good.
            return nil
        }
        // It exists but is not a directory (it's a file).
        return fmt.Errorf("path '%s' exists and is not a directory", dir)
    }

    // Path does not exist (or some other error occurred during Stat)
    if os.IsNotExist(err) {
        // The path or a prefix of it does not exist. Attempt to create.
        // This will also handle cases where parent directories need to be created.
        if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
            // If MkdirAll fails, it might be because a parent component is a file.
            return fmt.Errorf("failed to create output directory '%s': %w", dir, mkdirErr)
        }
        return nil
    }

    if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
        return fmt.Errorf("failed to create output directory '%s' (underlying stat error: %v): %w", dir, err, mkdirErr)
    }

    return nil
}
