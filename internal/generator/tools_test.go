package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/lyeslabs/mcpgen/internal/converter"
)

func Test_capitalizeFirstLetter(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"h", "H"},
		{"", ""},
		{"123abc", "123abc"},
		{"éclair", "Éclair"},
		{"aBC", "ABC"},
		{"A", "A"},
		{"!bang", "!bang"},
	}

	for _, tt := range tests {
		got := capitalizeFirstLetter(tt.in)
		if got != tt.want {
			t.Errorf("capitalizeFirstLetter(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func Test_extractImports(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected []string
	}{
		{
			name: "single import",
			src: `
				package main
				import "fmt"
			`,
			expected: []string{`"fmt"`},
		},
		{
			name: "multiple imports",
			src: `
				package main
				import (
					"fmt"
					"os"
				)
			`,
			expected: []string{`"fmt"`, `"os"`},
		},
		{
			name: "named import",
			src: `
				package main
				import f "fmt"
			`,
			expected: []string{`f "fmt"`},
		},
		{
			name: "dot import",
			src: `
				package main
				import . "math"
			`,
			expected: []string{`. "math"`},
		},
		{
			name: "underscore import",
			src: `
				package main
				import _ "net/http/pprof"
			`,
			expected: []string{`_ "net/http/pprof"`},
		},
		{
			name: "mixed imports",
			src: `
				package main
				import (
					"fmt"
					. "math"
					_ "net/http/pprof"
					myjson "encoding/json"
				)
			`,
			expected: []string{`"fmt"`, `. "math"`, `_ "net/http/pprof"`, `myjson "encoding/json"`},
		},
		{
			name: "no imports",
			src: `
				package main
				func main() {}
			`,
			expected: []string{},
		},
		{
			name:     "invalid Go code",
			src:      `not a go file`,
			expected: []string{},
		},
		{
			name: "import with comment",
			src: `
				package main
				import (
					"fmt" // standard fmt
					_ "net/http/pprof" // pprof
				)
			`,
			expected: []string{`"fmt"`, `_ "net/http/pprof"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractImports(tt.src)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("extractImports() = %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func Test_extractHandlerImplementation(t *testing.T) {
	const correctFunc = `
package main

import (
	"context"
	"fmt"
	"mcp"
)

// CreateTodoHandler is the handler function for the CreateTodo tool.
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// User logic here
	return nil, fmt.Errorf("%s not implemented", "CreateTodo")
}
`
	const multipleHandlers = `
package main
import (
	"context"
	"fmt"
	"mcp"
)
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return nil, fmt.Errorf("one")
}
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return nil, fmt.Errorf("two")
}
`
	const wrongNameFunc = `
package main
func OtherHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return nil, nil
}
`
	const wrongSignatureFunc = `
package main
func CreateTodoHandler(ctx context.Context) error {
	return nil
}
`
	const noFuncs = `
package main
var x = 42
`
	const invalidGo = `not a go file`

	tests := []struct {
		name         string
		src          string
		handlerName  string
		wantContains string // substring that must be in the extracted body
		wantEmpty    bool
		wantErr      bool
		errContains  string
	}{
		{
			name:         "correct handler",
			src:          correctFunc,
			handlerName:  "CreateTodoHandler",
			wantContains: "return nil, fmt.Errorf",
		},
		{
			name:        "multiple handlers",
			src:         multipleHandlers,
			handlerName: "CreateTodoHandler",
			wantEmpty:   true,
			wantErr:     true,
			errContains: "multiple handlers",
		},
		{
			name:        "wrong name",
			src:         wrongNameFunc,
			handlerName: "CreateTodoHandler",
			wantEmpty:   true,
		},
		{
			name:        "wrong signature",
			src:         wrongSignatureFunc,
			handlerName: "CreateTodoHandler",
			wantEmpty:   true,
		},
		{
			name:        "no functions",
			src:         noFuncs,
			handlerName: "CreateTodoHandler",
			wantEmpty:   true,
		},
		{
			name:        "invalid Go code",
			src:         invalidGo,
			handlerName: "CreateTodoHandler",
			wantEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := extractHandlerImplementation(tt.src, tt.handlerName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				if body != "" {
					t.Errorf("expected empty body, got: %q", body)
				}
			} else if tt.wantEmpty {
				if body != "" {
					t.Errorf("expected empty string, got: %q", body)
				}
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if !strings.Contains(body, tt.wantContains) {
					t.Errorf("expected body to contain %q, got: %q", tt.wantContains, body)
				}
				if !strings.HasPrefix(body, "{") {
					t.Errorf("expected body to start with '{', got: %q", body)
				}
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			}
		})
	}
}

func Test_replaceHandlerImplementation(t *testing.T) {
	const origFunc = `
package main

import (
	"context"
	"fmt"
	"mcp"
)

// CreateTodoHandler is the handler function for the CreateTodo tool.
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// User logic here
	return nil, fmt.Errorf("%s not implemented", "CreateTodo")
}
`

	const noHandler = `
package main
func OtherHandler() {}
`
	const noBody = `
package main
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
`
	const multipleFuncs = `
package main
func Foo() {}
func CreateTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// old
}
func Bar() {}
`

	tests := []struct {
		name           string
		src            string
		handlerName    string
		implementation string
		wantContains   string
		wantEqual      string
	}{
		{
			name:        "replace handler body",
			src:         origFunc,
			handlerName: "CreateTodoHandler",
			implementation: `{
	// NEW IMPLEMENTATION
	return &mcp.CallToolResult{}, nil
}`,
			wantContains: "// NEW IMPLEMENTATION",
		},
		{
			name:           "handler not found",
			src:            noHandler,
			handlerName:    "CreateTodoHandler",
			implementation: "// should not appear",
			wantEqual:      noHandler,
		},
		{
			name:           "handler with no body",
			src:            noBody,
			handlerName:    "CreateTodoHandler",
			implementation: "// should not appear",
			wantEqual:      noBody,
		},
		{
			name:           "multiple functions, only first replaced",
			src:            multipleFuncs,
			handlerName:    "CreateTodoHandler",
			implementation: "{ /* replaced */ }",
			wantContains:   "/* replaced */",
		},
		{
			name:           "invalid Go code",
			src:            "not a go file",
			handlerName:    "CreateTodoHandler",
			implementation: "// should not appear",
			wantEqual:      "not a go file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := replaceHandlerImplementation(tt.src, tt.handlerName, tt.implementation)
			if tt.wantEqual != "" {
				if out != tt.wantEqual {
					t.Errorf("expected output to be unchanged, got:\n%s", out)
				}
			} else if tt.wantContains != "" {
				if !strings.Contains(out, tt.wantContains) {
					t.Errorf("expected output to contain %q, got:\n%s", tt.wantContains, out)
				}
			}
		})
	}
}

func TestGenerateToolFiles(t *testing.T) {
	tmpDir := t.TempDir()
	toolsDir := filepath.Join(tmpDir, "mcptools")

	// Only populate the fields needed for template rendering
	config := &converter.MCPConfig{
		Tools: []converter.Tool{
			{
				Name:           "echo",
				Description:    "Echoes input",
				RawInputSchema: `{"type":"object","properties":{"msg":{"type":"string"}}}`,
				Responses: []converter.ResponseTemplate{
					{PrependBody: "// response", StatusCode: 200, ContentType: "application/json", Suffix: "// end"},
				},
				RequestTemplate: converter.RequestTemplate{
					URL:    "/echo",
					Method: "POST",
				},
			},
			{
				Name:           "reverse",
				Description:    "Reverses input",
				RawInputSchema: `{"type":"object","properties":{"msg":{"type":"string"}}}`,
				Responses: []converter.ResponseTemplate{
					{PrependBody: "// response", StatusCode: 200, ContentType: "application/json", Suffix: "// end"},
				},
				RequestTemplate: converter.RequestTemplate{
					URL:    "/reverse",
					Method: "POST",
				},
			},
		},
	}

	g := &Generator{
		PackageName: "mytools",
		outputDir:   tmpDir,
	}

	err := g.GenerateToolFiles(config)
	if err != nil {
		t.Fatalf("GenerateToolFiles failed: %v", err)
	}

	for _, tool := range config.Tools {
		fileName := capitalizeFirstLetter(tool.Name) + ".go"
		filePath := filepath.Join(toolsDir, fileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated file %s: %v", filePath, err)
		}
		content := string(data)

		// Check for package declaration
		if !strings.Contains(content, "package mcptools") {
			t.Errorf("Generated file %s missing package declaration", fileName)
		}
		// Check for handler function
		handlerName := capitalizeFirstLetter(tool.Name) + "Handler"
		if !strings.Contains(content, handlerName) {
			t.Errorf("Generated file %s missing handler %s", fileName, handlerName)
		}
		// Check for tool description
		if !strings.Contains(content, tool.Description) {
			t.Errorf("Generated file %s missing tool description", fileName)
		}
		// Check for input schema
		if tool.RawInputSchema != "" && !strings.Contains(content, tool.RawInputSchema) {
			t.Errorf("Generated file %s missing input schema", fileName)
		}
		// Check for response template content
		if len(tool.Responses) > 0 && !strings.Contains(content, tool.Responses[0].PrependBody) {
			t.Errorf("Generated file %s missing response template", fileName)
		}
	}
}

func TestGenerateToolFilesWithHandlerBodyImplemented(t *testing.T) {
	tmpDir := t.TempDir()
	toolsDir := filepath.Join(tmpDir, "mcptools")

	// Initial config with two tools
	config := &converter.MCPConfig{
		Tools: []converter.Tool{
			{
				Name:           "echo",
				Description:    "Echoes input",
				RawInputSchema: `{"type":"object","properties":{"msg":{"type":"string"}}}`,
				Responses: []converter.ResponseTemplate{
					{PrependBody: "// response", StatusCode: 200, ContentType: "application/json", Suffix: "// end"},
				},
				RequestTemplate: converter.RequestTemplate{
					URL:    "/echo",
					Method: "POST",
				},
			},
			{
				Name:           "reverse",
				Description:    "Reverses input",
				RawInputSchema: `{"type":"object","properties":{"msg":{"type":"string"}}}`,
				Responses: []converter.ResponseTemplate{
					{PrependBody: "// response", StatusCode: 200, ContentType: "application/json", Suffix: "// end"},
				},
				RequestTemplate: converter.RequestTemplate{
					URL:    "/reverse",
					Method: "POST",
				},
			},
		},
	}

	g := &Generator{
		PackageName: "mytools",
		outputDir:   tmpDir,
	}

	// 1. Generate tool files for the first time
	if err := g.GenerateToolFiles(config); err != nil {
		t.Fatalf("GenerateToolFiles failed: %v", err)
	}

	// 2. Overwrite the Echo.go file with a custom handler implementation
	echoFile := filepath.Join(toolsDir, "Echo.go")
	origContent, err := os.ReadFile(echoFile)
	if err != nil {
		t.Fatalf("Failed to read generated Echo.go: %v", err)
	}
	customHandler := `
func EchoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// CUSTOM USER LOGIC
	return &mcp.CallToolResult{Payload: []byte("custom")}, nil
}
`

	// Replace the handler body in the file
	modified := replaceHandlerImplementation(string(origContent), "EchoHandler", customHandler)

	if err := os.WriteFile(echoFile, []byte(modified), 0644); err != nil {
		t.Fatalf("Failed to write custom Echo.go: %v", err)
	}

	// 3. Regenerate tool files (should preserve custom handler)
	if err := g.GenerateToolFiles(config); err != nil {
		t.Fatalf("GenerateToolFiles (second run) failed: %v", err)
	}

	// 4. Check that Echo.go still contains the custom handler
	data, err := os.ReadFile(echoFile)
	if err != nil {
		t.Fatalf("Failed to read Echo.go after regeneration: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "// CUSTOM USER LOGIC") {
		t.Errorf("Custom handler implementation was not preserved in Echo.go")
	}
}
