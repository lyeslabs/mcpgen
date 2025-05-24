package converter

import (
	"fmt"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
)

// getOperations returns a map of HTTP method to operation
func getOperations(pathItem *openapi3.PathItem) map[string]*openapi3.Operation {
	operations := make(map[string]*openapi3.Operation)

	if pathItem.Get != nil {
		operations["get"] = pathItem.Get
	}
	if pathItem.Post != nil {
		operations["post"] = pathItem.Post
	}
	if pathItem.Put != nil {
		operations["put"] = pathItem.Put
	}
	if pathItem.Delete != nil {
		operations["delete"] = pathItem.Delete
	}
	if pathItem.Options != nil {
		operations["options"] = pathItem.Options
	}
	if pathItem.Head != nil {
		operations["head"] = pathItem.Head
	}
	if pathItem.Patch != nil {
		operations["patch"] = pathItem.Patch
	}
	if pathItem.Trace != nil {
		operations["trace"] = pathItem.Trace
	}

	return operations
}

// convertOperation converts an OpenAPI operation to an MCP tool
func (c *Converter) convertOperation(path, method string, operation *openapi3.Operation) (*Tool, error) {
	// Generate a tool name
	toolName := c.parser.GetOperationID(path, method, operation)

	// Create the tool
	tool := &Tool{
		Name:        toolName,
		Description: getDescription(operation),
		Args:        []Arg{},
	}

	// Convert parameters to arguments
	args, err := c.convertParameters(operation.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to convert parameters: %w", err)
	}
	if len(args) > 0 {
		tool.Args = append(tool.Args, args...)
	}

	// Convert request body to arguments
	bodyArgs, err := c.convertRequestBody(operation.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request body: %w", err)
	}
	if bodyArgs != nil {
		tool.Args = append(tool.Args, *bodyArgs)
	}

	rawInputSchema, err := GenerateJSONSchemaDraft7(tool.Args)
	if err != nil {
		return nil, fmt.Errorf("failed creating raw input schema for the %s tool input", toolName)
	}

	tool.RawInputSchema = rawInputSchema

	// Sort arguments by name for consistent output
	sort.Slice(tool.Args, func(i, j int) bool {
		return tool.Args[i].Name < tool.Args[j].Name
	})

	// Create request template
	requestTemplate, err := c.createRequestTemplate(path, method, operation)
	if err != nil {
		return nil, fmt.Errorf("failed to create request template: %w", err)
	}
	tool.RequestTemplate = *requestTemplate

	// Create response template
	responseTemplate, err := c.createResponseTemplates(operation)
	if err != nil {
		return nil, fmt.Errorf("failed to create response template: %w", err)
	}
	tool.Responses = responseTemplate

	return tool, nil
}
