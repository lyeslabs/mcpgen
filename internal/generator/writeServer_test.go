package generator

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func Test_RenderAndWriteServerTemplate(t *testing.T) {
	// Load the real template file from disk (relative to test file)
	templatePath := filepath.Join("templates", "server.templ")
	templateBytes, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read template file: %v", err)
	}

	// Create ToolTemplateData slice
	tools := []ToolTemplateData{
		{
			ToolNameOriginal: "Echo",
			ToolNameGo:       "Echo",
			ToolHandlerName:  "EchoHandler",
			ToolDescription:  "Echoes input",
		},
		{
			ToolNameOriginal: "Reverse",
			ToolNameGo:       "Reverse",
			ToolHandlerName:  "ReverseHandler",
			ToolDescription:  "Reverses input",
		},
	}

	// Prepare the data struct as GenerateServerFile would
	data := struct {
		PackageName        string
		MCPToolsImportPath string
		Tools              []ToolTemplateData
	}{
		PackageName:        "mytools",
		MCPToolsImportPath: "github.com/example/project/mcptools",
		Tools:              tools,
	}

	// Parse and render the template
	tmpl, err := template.New("server.templ").Parse(string(templateBytes))
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	// Format the code
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		t.Fatalf("failed to format code: %v", err)
	}

	// Write to a temp dir
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "server.go")
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// Read back and check content
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	strContent := string(content)

	// Check for package declaration and handler names
	if !strings.Contains(strContent, "package mytools") {
		t.Errorf("Generated file missing package declaration")
	}
	if !strings.Contains(strContent, "EchoHandler") || !strings.Contains(strContent, "ReverseHandler") {
		t.Errorf("Generated file missing expected handler names")
	}
}
