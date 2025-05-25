package mcptools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

// Input Schema for the GetTodoById tool
const GetTodoByIdInputSchema = `{
  "properties": {
    "todoId": {
      "description": "ID of the todo item to retrieve.",
      "format": "uuid",
      "type": "string"
    }
  },
  "required": [
    "todoId"
  ],
  "type": "object"
}`

// Response Template for the GetTodoById tool (Status: 200, Content-Type: application/json)
const GetTodoByIdResponseTemplate_A = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 200

**Content-Type:** application/json

> The requested todo item.

## Response Structure

- Structure (Type: object):
  - **status**: Current status of the todo item. (Type: string):
      - Default: 'pending'
      - Example: 'pending'
      - Enum: ['pending', 'in-progress', 'completed']
  - **title**: The main content of the todo item. (Type: string):
      - Example: 'Buy groceries'
  - **updatedAt**: Timestamp of when the todo item was last updated. (Type: string, date-time):
      - Example: '2025-05-10T10:00:00Z'
  - **createdAt**: Timestamp of when the todo item was created. (Type: string, date-time):
      - Example: '2025-05-09T18:12:54Z'
  - **id**: Unique identifier for the todo item. (Type: string, uuid):
      - Example: 'd290f1ee-6c54-4b01-90e6-d701748f0851'
`

// Response Template for the GetTodoById tool (Status: 404, Content-Type: application/json)
const GetTodoByIdResponseTemplate_B = `# API Response Information

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

// Response Template for the GetTodoById tool (Status: 500, Content-Type: application/json)
const GetTodoByIdResponseTemplate_C = `# API Response Information

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

// NewGetTodoByIdMCPTool creates the MCP Tool instance for GetTodoById
func NewGetTodoByIdMCPTool() mcp.Tool {
	return mcp.NewToolWithRawSchema(
		"GetTodoById",
		"Get a specific todo item - Retrieves a single todo item by its ID.",
		[]byte(GetTodoByIdInputSchema),
	)
}

// GetTodoByIdHandler is the handler function for the GetTodoById tool.
// This function is automatically generated. Users should implement the actual
// logic within this function body to integrate with backend APIs.
// You can generate types, http client and helpers for parsing request params to facilitate the implementation.
func GetTodoByIdHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	// IMPORTANT: Replace the following placeholder implementation with your actual logic.
	// Use the 'request' parameter to access tool call arguments.
	// Make HTTP calls or interact with services as needed.
	// Return an *mcp.CallToolResult with the response payload, or an error.

	// Example placeholder implementation:
	// Extract the parameters from the request and parse them.
	// Call your backend API or perform the necessary operations using 'params'.
	// Handle the response and errors accordingly.
	return nil, fmt.Errorf("%s not implemented", "GetTodoById")
}
