package mcptools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

// Input Schema for the CreateTodo tool
const CreateTodoInputSchema = `{
  "properties": {
    "body": {
      "description": "Todo item to create.",
      "properties": {
        "priority": {
          "enum": [
            "low",
            "medium",
            "high"
          ],
          "type": "string"
        },
        "title": {
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "title"
      ],
      "type": "object"
    }
  },
  "required": [
    "body"
  ],
  "type": "object"
}`

// Response Template for the CreateTodo tool (Status: 201, Content-Type: application/json)
const CreateTodoResponseTemplate_A = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 201

**Content-Type:** application/json

> Todo item created successfully.

## Response Structure

- Structure (Type: object):
  - **completed** (Type: boolean):
  - **id** (Type: integer):
  - **title** (Type: string):
`

// Response Template for the CreateTodo tool (Status: 201, Content-Type: application/xml)
const CreateTodoResponseTemplate_B = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 201

**Content-Type:** application/xml

> Todo item created successfully.

## Response Structure

- Structure (Type: object):
  - **completed** (Type: boolean):
  - **id** (Type: integer):
  - **title** (Type: string):
`

// Response Template for the CreateTodo tool (Status: 201, Content-Type: text/plain)
const CreateTodoResponseTemplate_C = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 201

**Content-Type:** text/plain

> Todo item created successfully.

## Response Structure

- Structure (Type: string):
    - Example: 'Created'
`

// Response Template for the CreateTodo tool (Status: 207, Content-Type: application/json)
const CreateTodoResponseTemplate_D = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 207

**Content-Type:** application/json

> Multi-status response for batch creation.

## Response Structure

- Structure (Type: array):
  - **Items** (Type: Combinator):
    - **One Of the following structures**:
      - **Option 1** (Type: object):
        - **id** (Type: integer):
        - **title** (Type: string):
      - **Option 2** (Type: object):
        - **error** (Type: string):
`

// Response Template for the CreateTodo tool (Status: 400, Content-Type: application/json)
const CreateTodoResponseTemplate_E = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 400

**Content-Type:** application/json

> Bad request due to invalid input.

## Response Structure

- Structure (Type: Combinator):
  - **Combines All Of the following structures**:
    - **Part 1** (Type: object):
      - **message** (Type: string):
    - **Part 2** (Type: object):
      - **details** (Type: array):
        - **Items** (Type: string):
`

// Response Template for the CreateTodo tool (Status: 422, Content-Type: application/json)
const CreateTodoResponseTemplate_F = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 422

**Content-Type:** application/json

> Unprocessable entity due to validation errors.

## Response Structure

- Structure (Type: object):
  - **errors** (Type: array):
    - **Items** (Type: Combinator):
      - **Any Of the following structures**:
        - **Option 1** (Type: object):
          - **error** (Type: string):
          - **field** (Type: string):
        - **Option 2** (Type: string):
`

// Response Template for the CreateTodo tool (Status: 500, Content-Type: application/json)
const CreateTodoResponseTemplate_G = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 500

**Content-Type:** application/json

> Internal server error.

## Response Structure

- Structure (Type: object):
  - **message** (Type: string):
  - **traceId** (Type: string):
`

// Response Template for the CreateTodo tool (Status: 500, Content-Type: text/plain)
const CreateTodoResponseTemplate_H = `# API Response Information

Below is the response template for this API endpoint.

The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.

**Status Code:** 500

**Content-Type:** text/plain

> Internal server error.

## Response Structure

- Structure (Type: string):
    - Example: 'Internal Server Error'
`

// NewCreateTodoMCPTool creates the MCP Tool instance for CreateTodo
func NewCreateTodoMCPTool() mcp.Tool {
	return mcp.NewToolWithRawSchema(
		"CreateTodo",
		"Create a new todo item - Adds a new item to the todo list.",
		[]byte(CreateTodoInputSchema),
	)
}

// CreateTodoHandler is the handler function for the CreateTodo tool.
// This function is automatically generated. Users should implement the actual
// logic within this function body to integrate with backend APIs.
// You can generate types, http client and helpers for parsing request params to facilitate the implementation.
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	// IMPORTANT: Replace the following placeholder implementation with your actual logic.
	// Use the 'request' parameter to access tool call arguments.
	// Make HTTP calls or interact with services as needed.
	// Return an *mcp.CallToolResult with the response payload, or an error.

	// Example placeholder implementation:
	// Extract the parameters from the request and parse them.
	// Call your backend API or perform the necessary operations using 'params'.
	// Handle the response and errors accordingly.
	return nil, fmt.Errorf("%s not implemented", "CreateTodo")
}
