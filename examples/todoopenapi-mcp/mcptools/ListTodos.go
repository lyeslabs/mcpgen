package mcptools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

// Input Schema for the ListTodos tool
const ListTodosInputSchema = `{
  "properties": {
    "limit": {
      "default": 20,
      "description": "Maximum number of todos to return",
      "format": "int32",
      "minimum": 1,
      "type": "integer"
    },
    "offset": {
      "default": 0,
      "description": "Number of todos to skip for pagination",
      "format": "int32",
      "minimum": 0,
      "type": "integer"
    },
    "status": {
      "description": "Filter todos by status (e.g., \"pending\", \"completed\")",
      "enum": [
        "pending",
        "completed",
        "in-progress"
      ],
      "type": "string"
    },
    "token": {
      "default": 20,
      "description": "Token for authentication",
      "format": "int32",
      "minimum": 1,
      "type": "integer"
    }
  },
  "type": "object"
}`

// Response Template for the ListTodos tool (Status: 200, Content-Type: application/json)
const ListTodosResponseTemplate_A = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 200

**Content-Type:** application/json

> A list of todo items.

## Response Structure

- Structure (Type: array):
  - **Items** (Type: Combinator):
    - **One Of the following structures**:
      - **Option 1** (Type: object):
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
      - **Option 2** (Type: object):
        - **title**: The main content of the todo item. (Type: string):
            - Example: 'Plan weekend trip'
        - **description**: Optional detailed description of the todo item. (Type: string, nullable):
            - Nullable: true
            - Example: 'Research destinations and book accommodation.'
        - **dueDate**: Optional due date for the todo item. (Type: string, date, nullable):
            - Nullable: true
            - Example: '2025-06-15'
        - **status**: Current status of the todo item. (Type: string):
            - Default: 'pending'
            - Example: 'pending'
            - Enum: ['pending', 'in-progress', 'completed']
`

// Response Template for the ListTodos tool (Status: 400, Content-Type: application/json)
const ListTodosResponseTemplate_B = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 400

**Content-Type:** application/json

> The request was malformed or invalid.

## Response Structure

- Structure (Type: object):
  - **code**: An application-specific error code. (Type: integer, int32):
  - **details**: Optional array of specific field validation errors. (Type: array):
    - **Items** (Type: object):
      - **field** (Type: string):
      - **issue** (Type: string):
  - **message**: A human-readable description of the error. (Type: string):
`

// Response Template for the ListTodos tool (Status: 500, Content-Type: application/json)
const ListTodosResponseTemplate_C = `# API Response Information

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

// NewListTodosMCPTool creates the MCP Tool instance for ListTodos
func NewListTodosMCPTool() mcp.Tool {
	return mcp.NewToolWithRawSchema(
		"ListTodos",
		"List all todo items - Retrieves a list of todo items, optionally filtered by status.",
		[]byte(ListTodosInputSchema),
	)
}

// ListTodosHandler is the handler function for the ListTodos tool.
// This function is automatically generated. Users should implement the actual
// logic within this function body to integrate with backend APIs.
// You can generate types, http client and helpers for parsing request params to facilitate the implementation.
func ListTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	// IMPORTANT: Replace the following placeholder implementation with your actual logic.
	// Use the 'request' parameter to access tool call arguments.
	// Make HTTP calls or interact with services as needed.
	// Return an *mcp.CallToolResult with the response payload, or an error.

	// Example placeholder implementation:
	// Extract the parameters from the request and parse them.
	// Call your backend API or perform the necessary operations using 'params'.
	// Handle the response and errors accordingly.
	return nil, fmt.Errorf("%s not implemented", "ListTodos")
}
