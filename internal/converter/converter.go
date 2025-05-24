package converter

import (
	"fmt"
	"sort"
)

// Converter represents an OpenAPI to MCP converter
type Converter struct {
	parser  *Parser
	options ConvertOptions
}


type ConverterInterface interface {
    Convert() (*MCPConfig, error)
}

// NewConverter creates a new OpenAPI to MCP converter
func NewConverter(parser *Parser) *Converter {
	return &Converter{
		parser: parser,
		options: ConvertOptions{
			ServerConfig: make(map[string]interface{}),
		},
	}
}


// Convert converts an OpenAPI document to an MCP configuration
func (c *Converter) Convert() (*MCPConfig, error) {
	if c.parser.GetDocument() == nil {
		return nil, fmt.Errorf("no OpenAPI document loaded")
	}

	// Create the MCP configuration
	config := &MCPConfig{
		Server: ServerConfig{
			Config: c.options.ServerConfig,
		},
		Tools: []Tool{},
	}

	// Process each path and operation
	for path, pathItem := range c.parser.GetPaths() {
		operations := getOperations(pathItem)
		for method, operation := range operations {
			tool, err := c.convertOperation(path, method, operation)
			if err != nil {
				return nil, fmt.Errorf("failed to convert operation %s %s: %w", method, path, err)
			}
			config.Tools = append(config.Tools, *tool)
		}
	}

	// Sort tools by name for consistent output
	sort.Slice(config.Tools, func(i, j int) bool {
		return config.Tools[i].Name < config.Tools[j].Name
	})

	return config, nil
}
