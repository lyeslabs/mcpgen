package mcptools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

// Input Schema for the DeleteTodoById tool
const DeleteTodoByIdInputSchema = `{
  "properties": {
    "todoId": {
      "description": "ID of the todo item to delete.",
      "format": "uuid",
      "type": "string"
    }
  },
  "required": [
    "todoId"
  ],
  "type": "object"
}`

// Response Template for the DeleteTodoById tool (Status: 404, Content-Type: application/json)
const DeleteTodoByIdResponseTemplate_A = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 404

**Content-Type:** application/json

> The specified resource was not found.

## Response Structure

- Structure (Type: object):
  - **code**: An application-specific error code. (Type: integer, int32):
  - **details**: Optional array of specific field validation errors. (Type: array):
    - **Items** (Type: object):
      - **field** (Type: string):
      - **issue** (Type: string):
  - **message**: A human-readable description of the error. (Type: string):
`

// Response Template for the DeleteTodoById tool (Status: 500, Content-Type: application/json)
const DeleteTodoByIdResponseTemplate_B = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 500

**Content-Type:** application/json

> An unexpected error occurred on the server.

## Response Structure

- Structure (Type: object):
  - **code**: An application-specific error code. (Type: integer, int32):
  - **details**: Optional array of specific field validation errors. (Type: array):
    - **Items** (Type: object):
      - **field** (Type: string):
      - **issue** (Type: string):
  - **message**: A human-readable description of the error. (Type: string):
`

// NewDeleteTodoByIdMCPTool creates the MCP Tool instance for DeleteTodoById
func NewDeleteTodoByIdMCPTool() mcp.Tool {
	return mcp.NewToolWithRawSchema(
		"DeleteTodoById",
		"Delete a todo item - Removes a todo item by its ID.",
		[]byte(DeleteTodoByIdInputSchema),
	)
}

// DeleteTodoByIdHandler is the handler function for the DeleteTodoById tool.
// This function is automatically generated. Users should implement the actual
// logic within this function body to integrate with backend APIs.
// You can generate types, http client and helpers for parsing request params to facilitate the implementation.
func DeleteTodoByIdHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	// IMPORTANT: Replace the following placeholder implementation with your actual logic.
	// Use the 'request' parameter to access tool call arguments.
	// Make HTTP calls or interact with services as needed.
	// Return an *mcp.CallToolResult with the response payload, or an error.

	// Example placeholder implementation:
	// Extract the parameters from the request and parse them.
	// Call your backend API or perform the necessary operations using 'params'.
	// Handle the response and errors accordingly.
	return nil, fmt.Errorf("%s not implemented", "DeleteTodoById")
}
