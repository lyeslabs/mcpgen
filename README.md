# mcpgen: Seamlessly Transform OpenAPI APIs into AI Agent Tools

**mcpgen** is a command-line tool that seamlessly generates production-ready Model Context Protocol (MCP) server boilerplate from your OpenAPI specifications, enabling you to easily expose your existing APIs as powerful tools for AI agents.


## Key Features

-   **OpenAPI Compatibility:** Reads and processes OpenAPI specifications in YAML or JSON format, supporting versions **3.1, 3.0**.
-   **Comprehensive MCP Server Generation:** Generates the full Go boilerplate required to set up an MCP server, including server initialization, tool registration, and handler skeletons.
-   **Accurate Schema Translation:** Automatically translates OpenAPI schema definitions into the necessary **JSON Schemas** for tool inputs (compatible with MCP) and generates detailed markdown-based **Response Templates (Prompts)** for various status codes and content types, providing rich context for AI models.
-   **Advanced Schema Support:** Handles complex OpenAPI schema constructs, including:
    *   Arrays and nested structures
    *   Type unions and `oneOf`/`anyOf`/`allOf` combinators
    *   Recursive type definitions
    *   Validation constraints (e.g., `minimum`, `maximum`, `maxLength`, `pattern`)
-   **Multiple Content Type Handling:** Correctly processes and generates templates for operations defining multiple request and response content types.
-   **Generated Code Quality:** Produces well-structured, idiomatic Go code with embedded schemas and prompts, leveraging constants for clarity and maintainability.
-   **Optional Client & Types Generation:** Can optionally generate a Go HTTP client and corresponding Go types based on your OpenAPI schema, simplifying the implementation logic within the generated MCP handlers.
-   **Developer Friendly:** Provides clear handler function skeletons with guidance on where to integrate your core logic to connect to the actual backend API.
-   **Battle-Tested Foundation:** The code generation logic is backed by extensive testing, ensuring high reliability (as demonstrated by 94% coverage and 200% test volume).

## Installation

```sh
go install github.com/lyeslabs/mcpgen/cmd/mcpgen@latest
```

By default, the binary is installed to `$HOME/go/bin` (or `%USERPROFILE%\go\bin` on Windows).
Make sure this directory is in your `PATH`.

## Usage

```sh
mcpgen --input openapi.yaml --output generated-server
```

### Required flags

-   `--input`
    Path to your OpenAPI specification file (YAML or JSON).

-   `--output`
    Output directory for the generated MCP server boilerplate.

### Optional flags

-   `--validation`
    Enable OpenAPI validation (default: `false`).

-   `--package`
    Name for the generated Go package (default: `mcpgen`).

-   `--includes`
    Comma-separated list of additional includes for the generated code. Use `httpclient,types` to generate the HTTP client and types.

### Example

```sh
mcpgen --input api/openapi.yaml --output ./generated-server --validation --package myserver --includes=httpclient,types
```

## How It Works

`mcpgen` acts as a bridge between your declarative OpenAPI specification and the programmatic Go code required for an MCP server. It reads your OpenAPI definition and automatically generates the necessary boilerplate, including the structured schemas and prompts essential for effective AI agent interaction.

Let's illustrate this with an example of a moderately complex endpoint defined in OpenAPI:

```yaml
# This is a snippet from your OpenAPI specification
/todos:
    get:
      tags:
        - Todos
      summary: List all todo items
      description: Retrieves a list of todo items, optionally filtered by status.
      operationId: listTodos
      parameters:
        - name: status
          in: query
          description: Filter todos by status (e.g., "pending", "completed")
          required: false
          schema:
            type: string
            enum: [pending, completed, in-progress]
        - name: token
          in: cookie
          description: Token for authentication
          required: false
          schema:
            type: integer
            format: int32
            minimum: 1
            default: 20
        - name: limit
          in: query
          description: Maximum number of todos to return
          required: false
          schema:
            type: integer
            format: int32
            minimum: 1
            default: 20
        - name: offset
          in: query
          description: Number of todos to skip for pagination
          required: false
          schema:
            type: integer
            format: int32
            minimum: 0
            default: 0
      responses:
        '200':
          description: A list of todo items.
          content:
            application/json:
              schema:
                type: array
                items:
                  oneOf:
                    - $ref: '#/components/schemas/Todo'
                    - $ref: '#/components/schemas/NewTodo'
        '400':
          $ref: '#/components/responses/BadRequest'
        '500':
          $ref: '#/components/responses/InternalServerError'

# ... (The full spec would include components for #/components/schemas/Todo,
# #/components/schemas/NewTodo, #/components/responses/BadRequest, etc.)
```

When you run `mcpgen --input your_openapi.yaml --output generated-server` (and optionally `--includes=httpclient,types`), `mcpgen` analyzes this operation (`operationId: listTodos`) and generates Go code. This includes:

1.  **JSON Schema for Input:** A constant string containing the JSON Schema representing the required and optional parameters for the `listTodos` tool:

    ```go
    // Input Schema for the ListTodos tool
    const listTodosInputSchema = `{
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
    ```

2.  **Markdown Templates for Responses:** Constant strings containing detailed markdown prompts for each potential response (based on status codes and content types), describing the structure and meaning of the data. This is crucial for LLMs to understand the tool's output:

    ```go
    // Response Template for the ListTodos tool (Status: 200, Content-Type: application/json)
    const ListTodosResponseTemplate_A = `# API Response Information
    ... (detailed markdown describing the 200 response structure including the oneOf combining Todo and NewTodo schemas) ...
    `

    // Response Template for the ListTodos tool (Status: 400, Content-Type: application/json)
    const ListTodosResponseTemplate_B = `# API Response Information
    ... (detailed markdown describing the 400 error response structure) ...
    `

    // Response Template for the ListTodos tool (Status: 500, Content-Type: application/json)
    const ListTodosResponseTemplate_C = `# API Response Information
    ... (detailed markdown describing the 500 error response structure) ...
    `
    ```
    *Note: The full content of the markdown templates is extensive and generated based on the OpenAPI response schemas.*

3.  **MCP Tool Registration:** A function to create and configure the `mcp.Tool` instance, embedding the operation's description and the generated input schema:

    ```go
    // NewListTodosMCPTool creates the MCP Tool instance for ListTodos
    func NewListTodosMCPTool() mcp.Tool {
    	return mcp.NewToolWithRawSchema(
    		"ListTodos", // Operation ID becomes the Tool Name
    		"List all todo items - Retrieves a list of todo items, optionally filtered by status.", // Summary + Description
    		[]byte(listTodosInputSchema), // Embedded Input Schema
    	)
    }
    ```

4.  **Handler Function Skeleton:** A placeholder function where you will write the code to handle the tool call. This function receives the `mcp.CallToolRequest` (containing the input payload as JSON) and is where you will integrate with your actual backend API:

    ```go
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
    	return nil, fmt.Errorf("ListTodos handler not implemented") // Placeholder until you add your logic
    }
    ```

By generating all this structured boilerplate code, `mcpgen` allows you to focus solely on implementing the core integration logic within the generated handler functions – parsing the input payload (potentially simplified by generated types), calling your existing backend API (potentially simplified by a generated client), and mapping the backend response to the expected MCP `CallToolResult` format.


## License

This project is licensed under the [MIT License](LICENSE).
