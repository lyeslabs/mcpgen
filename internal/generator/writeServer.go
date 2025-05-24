package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/lyeslabs/mcpgen/internal/converter"
)

// GenerateServerFile creates a server.go file in the same package as the tools
func (g *Generator) GenerateServerFile(config *converter.MCPConfig) error {
	serverTemplateContent, err := templatesFS.ReadFile("templates/server.templ")
	if err != nil {
		return fmt.Errorf("failed to read server template file: %w", err)
	}

	tmpl, err := template.New("server.templ").Parse(string(serverTemplateContent))
	if err != nil {
		return fmt.Errorf("failed to parse server template: %w", err)
	}

	importPath, err := BuildImportPath(g.outputDir)
	if err != nil {
		return fmt.Errorf("failed to build import path: %w", err)
	}

	data := struct {
		PackageName        string
		MCPToolsImportPath string
		Tools              []ToolTemplateData
	}{
		PackageName:        g.PackageName,
		Tools:              make([]ToolTemplateData, 0, len(config.Tools)),
		MCPToolsImportPath: importPath,
	}

	for _, tool := range config.Tools {
		capitalizedName := capitalizeFirstLetter(tool.Name)

		data.Tools = append(data.Tools, ToolTemplateData{
			ToolNameOriginal: capitalizedName,
			ToolNameGo:       capitalizedName,
			ToolHandlerName:  capitalizedName + "Handler",
			ToolDescription:  tool.Description,
		})
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to render server template: %w", err)
	}

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated server.go: %w", err)
	}

	if err := writeFileContent(g.outputDir, "server.go", func() ([]byte, error) {
		return formattedCode, nil
	}); err != nil {
		return fmt.Errorf("failed to write server.go file: %w", err)
	}

	return nil
}
