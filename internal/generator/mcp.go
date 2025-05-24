package generator

import (
	"embed"
	"fmt"

	"github.com/lyeslabs/mcpgen/internal/converter"
)

//go:embed templates/*.templ
var templatesFS embed.FS

// ToolTemplateData holds the data to pass to the template for a single tool
type ToolTemplateData struct {
	ToolNameOriginal      string
	ToolNameGo            string
	ToolHandlerName       string
	ToolDescription       string
	RawInputSchema        string
	ResponseTemplate      []converter.ResponseTemplate
	InputSchemaConst      string
	ResponseTemplateConst string
}

// GenerateMCP generates the MCP tool files while preserving existing handler implementations and imports
func (g *Generator) GenerateMCP() error {
	config, err := g.converter.Convert()
	if err != nil {
		return fmt.Errorf("failed at converting OpenAPI schema into MCP code %w", err)
	}

	if err := g.GenerateServerFile(config); err != nil {
		return fmt.Errorf("failed to generate server file: %w", err)
	}

	if err := g.GenerateToolFiles(config); err != nil {
		return fmt.Errorf("failed to generate tool files: %w", err)
	}

	if err := g.GenerateHelpers(); err != nil {
		return fmt.Errorf("failed to generate helpers: %w", err)
	}

	return nil
}
