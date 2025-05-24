package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to create a temporary go.mod file
func createTempGoMod(t *testing.T, dir string, content string) string {
	t.Helper()
	path := filepath.Join(dir, goModFile)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp go.mod: %v", err)
	}
	return path
}

func TestParseModuleName(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name          string
		content       string
		expectedName  string
		expectError   bool
		errorContains string
	}{
		{"simple", "module example.com/project\n\ngo 1.20", "example.com/project", false, ""},
		{"quoted", "module \"example.com/project/v2\"\n", "example.com/project/v2", false, ""},
		{"single quoted", "module 'example.com/project/v3'\n", "example.com/project/v3", false, ""},
		{"with comment", "module example.com/project // my module\n", "example.com/project", false, ""},
		{"quoted with comment", "module \"example.com/project\" // my module\n", "example.com/project", false, ""},
		{"extra spaces", "  module    example.com/another   \n", "example.com/another", false, ""},
		{"no module line", "go 1.20\nrequire other/thing v1.0.0", "", true, "module declaration not found"},
		{"empty file", "", "", true, "module declaration not found"},
		{"unclosed quote", "module \"example.com/project\n", "", true, "unclosed quote"},
		{"empty module name", "module \n", "", true, "empty module name after 'module' keyword"},
		{"module keyword only", "module\n", "", true, "empty module name after 'module' keyword"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			goModPath := createTempGoMod(t, tempDir, tc.content)
			name, err := parseModuleName(goModPath)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tc.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if name != tc.expectedName {
					t.Errorf("Expected module name %q, got %q", tc.expectedName, name)
				}
			}
		})
	}

	t.Run("file not found", func(t *testing.T) {
		_, err := parseModuleName(filepath.Join(tempDir, "nonexistent.mod"))
		if err == nil {
			t.Errorf("Expected error for non-existent file, got nil")
		} else if !strings.Contains(err.Error(), "failed to open") {
			t.Errorf("Error message %q does not contain 'failed to open'", err.Error())
		}
	})
}

// Helper to set up a test module directory structure
func setupTestModuleStructure(t *testing.T, modulePathInTemp string, moduleName string, subDirsToCreate ...string) (moduleRootDir string, deepestDir string) {
	t.Helper()
	baseTempDir := t.TempDir()
	moduleRootDir = filepath.Join(baseTempDir, modulePathInTemp)

	if err := os.MkdirAll(moduleRootDir, 0755); err != nil {
		t.Fatalf("Failed to create module root dir %s: %v", moduleRootDir, err)
	}
	createTempGoMod(t, moduleRootDir, fmt.Sprintf("module %s\n\ngo 1.20", moduleName))

	currentPath := moduleRootDir
	for _, subDir := range subDirsToCreate {
		currentPath = filepath.Join(currentPath, subDir)
		if err := os.MkdirAll(currentPath, 0755); err != nil {
			t.Fatalf("Failed to create subdir %s: %v", currentPath, err)
		}
	}
	return moduleRootDir, currentPath
}

func TestFindModulePath(t *testing.T) {
	testCases := []struct {
		name             string
		modulePathInTemp string // e.g., "projectA"
		moduleName       string
		startPathSubDirs []string // subdirs from modulePathInTemp to start search
		expectedModName  string
		expectError      bool
		errorContains    string
		setupFunc        func(t *testing.T, moduleRootDir string) // <-- Add this line
	}{
		{
			name:             "gomod in startDir",
			modulePathInTemp: "proj1",
			moduleName:       "example.com/proj1",
			startPathSubDirs: []string{},
			expectedModName:  "example.com/proj1",
		},
		{
			name:             "gomod one level up",
			modulePathInTemp: "proj2",
			moduleName:       "example.com/proj2",
			startPathSubDirs: []string{"cmd"},
			expectedModName:  "example.com/proj2",
		},
		{
			name:             "gomod multiple levels up",
			modulePathInTemp: "proj3",
			moduleName:       "example.com/proj3",
			startPathSubDirs: []string{"internal", "service", "impl"},
			expectedModName:  "example.com/proj3",
		},
		{
			name:             "no gomod found (search depth)",
			modulePathInTemp: "proj4_no_mod_setup", // We won't create go.mod for this one
			moduleName:       "",                   // No module name
			startPathSubDirs: []string{"a", "b", "c", "d", "e"},
			expectError:      true,
			errorContains:    fmt.Sprintf("no go.mod found in %d levels", maxSearchDepth),
		},
		{
			name:             "malformed gomod",
			modulePathInTemp: "proj5_malformed",
			moduleName:       "example.com/proj5", // This name won't be used due to malformed content
			startPathSubDirs: []string{},
			expectError:      true,
			errorContains:    "module declaration not found", // Error from parseModuleName
			// Custom setup for this case:
			setupFunc: func(t *testing.T, moduleRootDir string) {
				createTempGoMod(t, moduleRootDir, "this is not a valid go.mod")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var moduleRoot, startDir string
			if tc.modulePathInTemp == "proj4_no_mod_setup" {
				// Special case: create directory structure but no go.mod
				baseTempDir := t.TempDir()
				moduleRoot = filepath.Join(baseTempDir, tc.modulePathInTemp) // This dir won't have go.mod
				currentPath := moduleRoot
				for _, subDir := range tc.startPathSubDirs {
					currentPath = filepath.Join(currentPath, subDir)
				}
				if err := os.MkdirAll(currentPath, 0755); err != nil {
					t.Fatalf("Failed to create dirs for no_mod_setup: %v", err)
				}
				startDir = currentPath
			} else {
				moduleRoot, startDir = setupTestModuleStructure(t, tc.modulePathInTemp, tc.moduleName, tc.startPathSubDirs...)
				if tc.setupFunc != nil {
					// Allow custom setup, e.g., for malformed go.mod
					type testCaseWithSetup struct { // Define struct locally for setupFunc
						setupFunc func(t *testing.T, moduleRootDir string)
					}
					// This is a bit of a workaround to access tc.setupFunc if it exists
					// A better way would be to add setupFunc to the main testCases struct definition
					// For now, this cast will work if tc has a setupFunc field.
					// Let's assume tc has setupFunc for the "malformed gomod" case.
					if s, ok := interface{}(tc).(testCaseWithSetup); ok {
						s.setupFunc(t, moduleRoot)
					} else if tc.name == "malformed gomod" { // specific check for this test case
						createTempGoMod(t, moduleRoot, "this is not a valid go.mod")
					}
				}
			}

			modName, modRootPath, err := findModulePath(startDir)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tc.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if modName != tc.expectedModName {
					t.Errorf("Expected module name %q, got %q", tc.expectedModName, modName)
				}
				// Normalize paths for comparison as moduleRoot from setupTestModuleStructure is absolute
				expectedModRoot := filepath.Clean(moduleRoot)
				if filepath.Clean(modRootPath) != expectedModRoot {
					t.Errorf("Expected module root %q, got %q", expectedModRoot, modRootPath)
				}
			}
		})
	}
}

func TestBuildImportPath(t *testing.T) {
	testCases := []struct {
		name               string
		modulePathInTemp   string
		moduleName         string
		cwdSubDirs         []string // Subdirs from modulePathInTemp to set as CWD
		outputDirRelToCwd  string   // outputDir parameter for BuildImportPath
		expectedImportPath string
		expectError        bool
		errorContains      string
	}{
		{
			name:               "output in cwd, cwd is module root",
			modulePathInTemp:   "app1",
			moduleName:         "example.com/app1",
			cwdSubDirs:         []string{},
			outputDirRelToCwd:  ".", // mcptools will be in example.com/app1/mcptools
			expectedImportPath: "example.com/app1/mcptools",
		},
		{
			name:               "output in subdir of cwd, cwd is module root",
			modulePathInTemp:   "app2",
			moduleName:         "example.com/app2",
			cwdSubDirs:         []string{},
			outputDirRelToCwd:  "generated", // mcptools will be in example.com/app2/generated/mcptools
			expectedImportPath: "example.com/app2/generated/mcptools",
		},
		{
			name:               "cwd is subdir, output relative to cwd",
			modulePathInTemp:   "app3",
			moduleName:         "example.com/app3",
			cwdSubDirs:         []string{"cmd", "mytool"},
			outputDirRelToCwd:  "pkg", // mcptools will be in example.com/app3/cmd/mytool/pkg/mcptools
			expectedImportPath: "example.com/app3/cmd/mytool/pkg/mcptools",
		},
		{
			name:               "output dir uses .. to go up from cwd",
			modulePathInTemp:   "app4",
			moduleName:         "example.com/app4",
			cwdSubDirs:         []string{"cmd", "mytool"},
			outputDirRelToCwd:  "../../generated", // mcptools will be in example.com/app4/generated/mcptools
			expectedImportPath: "example.com/app4/generated/mcptools",
		},
		{
			name:              "no go.mod found",
			modulePathInTemp:  "app5_no_mod", // No go.mod will be created here
			moduleName:        "",
			cwdSubDirs:        []string{"some", "dir"},
			outputDirRelToCwd: ".",
			expectError:       true,
			errorContains:     "failed to find module",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var moduleRoot, cwdSimulated string
			if tc.modulePathInTemp == "app5_no_mod" {
				baseTempDir := t.TempDir()
				moduleRoot = filepath.Join(baseTempDir, tc.modulePathInTemp) // This dir won't have go.mod
				currentPath := moduleRoot
				for _, subDir := range tc.cwdSubDirs {
					currentPath = filepath.Join(currentPath, subDir)
				}
				if err := os.MkdirAll(currentPath, 0755); err != nil {
					t.Fatalf("Failed to create dirs for no_mod_setup: %v", err)
				}
				cwdSimulated = currentPath
			} else {
				moduleRoot, cwdSimulated = setupTestModuleStructure(t, tc.modulePathInTemp, tc.moduleName, tc.cwdSubDirs...)
			}

			originalCwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get original CWD: %v", err)
			}
			if err := os.Chdir(cwdSimulated); err != nil {
				t.Fatalf("Failed to change CWD to %s: %v", cwdSimulated, err)
			}
			defer func() {
				if err := os.Chdir(originalCwd); err != nil {
					// Log error, but don't fail test here as main test logic is done
					t.Logf("Warning: failed to restore original CWD %s: %v", originalCwd, err)
				}
			}()

			importPath, err := BuildImportPath(tc.outputDirRelToCwd)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tc.errorContains)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if importPath != tc.expectedImportPath {
					t.Errorf("Expected import path %q, got %q", tc.expectedImportPath, importPath)
				}
			}
		})
	}
}
