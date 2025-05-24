package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
)

// GenerateHelpers creates a helpers.go file with utility functions for MCP tools
func (g *Generator) GenerateHelpers() error {
	helpersTemplate, err := templatesFS.ReadFile("templates/helpers.templ")
	if err != nil {
		return fmt.Errorf("failed to read helpers template file: %w", err)
	}

	tmpl, err := template.New("helpers").Parse(string(helpersTemplate))
	if err != nil {
		return fmt.Errorf("failed to parse helpers template: %w", err)
	}

	data := struct {
		PackageName string
	}{
		PackageName: g.PackageName,
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return fmt.Errorf("failed to execute helpers template: %w", err)
	}

	formattedCode, err := format.Source(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated helpers code: %w", err)
	}

	err = writeFileContent(g.outputDir + "/helpers", "params.go", func() ([]byte, error) {
		return formattedCode, nil
	})
	if err != nil {
		return fmt.Errorf("failed to write helpers.go file: %w", err)
	}

	return nil
}
