package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lyeslabs/mcpgen/internal/converter"
)

// GenerateToolFiles generates individual tool files while preserving existing handler implementations
func (g *Generator) GenerateToolFiles(config *converter.MCPConfig) error {
	toolTemplateContent, err := templatesFS.ReadFile("templates/tool.templ")
	if err != nil {
		return fmt.Errorf("failed to read tool template file: %w", err)
	}

	tmpl, err := template.New("tool.templ").Parse(string(toolTemplateContent))
	if err != nil {
		return fmt.Errorf("failed to parse tool template: %w", err)
	}

	for _, tool := range config.Tools {
		capitalizedName := capitalizeFirstLetter(tool.Name)
		data := struct {
			ToolTemplateData
			URL     string
			Method  string
			Headers []converter.Header
		}{
			ToolTemplateData: ToolTemplateData{
				ToolNameOriginal:      capitalizedName,
				ToolNameGo:            capitalizedName,
				ToolHandlerName:       capitalizedName + "Handler",
				ToolDescription:       tool.Description,
				RawInputSchema:        tool.RawInputSchema,
				ResponseTemplate:      tool.Responses,
				InputSchemaConst:      fmt.Sprintf("%sInputSchema", tool.Name),
				ResponseTemplateConst: fmt.Sprintf("%sResponseTemplate", tool.Name),
			},
			URL:     tool.RequestTemplate.URL,
			Method:  tool.RequestTemplate.Method,
			Headers: tool.RequestTemplate.Headers,
		}

		outputFileName := capitalizedName + ".go"
		outputFilePath := filepath.Join(g.outputDir+"/mcptools", outputFileName)

		// Check if file already exists and extract handler implementation if it does
		existingImplementation := ""
		existingImports := []string{}

		if _, err := os.Stat(outputFilePath); err == nil {
			existingContent, err := os.ReadFile(outputFilePath)
			if err == nil {
				existingImplementation, err = extractHandlerImplementation(string(existingContent), data.ToolHandlerName)
				if err != nil {
					return err
				}
				// Extract existing imports
				existingImports = extractImports(string(existingContent))
			}
		}

		// Generate code for this tool
		var toolBuf bytes.Buffer

		// Write package declaration
		fmt.Fprintf(&toolBuf, "package mcptools\n\n")

		// Merge imports
		requiredImports := []string{
			"context",
			"fmt",
			"github.com/mark3labs/mcp-go/mcp",
		}

		if len(existingImports) > 0 {
			fmt.Fprintf(&toolBuf, "import (\n")
			for _, imp := range existingImports {
				fmt.Fprintf(&toolBuf, "\t%s\n", imp)
			}
			fmt.Fprintf(&toolBuf, ")\n\n")
		} else {
			fmt.Fprintf(&toolBuf, "import (\n")
			for _, imp := range requiredImports {
				fmt.Fprintf(&toolBuf, "\t\"%s\"\n", imp)
			}
			fmt.Fprintf(&toolBuf, ")\n\n")
		}

		// Execute template to get the boilerplate
		if err := tmpl.Execute(&toolBuf, data); err != nil {
			return fmt.Errorf("failed to render template for tool %s: %w", tool.Name, err)
		}

		// If we have an existing implementation, replace the default one
		if existingImplementation != "" {
			toolContent := toolBuf.String()
			toolContent = replaceHandlerImplementation(toolContent, data.ToolHandlerName, existingImplementation)
			toolBuf.Reset()
			toolBuf.WriteString(toolContent)
		}

		// Format the generated code
		formattedCode, err := format.Source(toolBuf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to format generated code for %s: %w", outputFileName, err)
		}

		err = writeFileContent(g.outputDir+"/mcptools", outputFileName, func() ([]byte, error) {
			return formattedCode, nil
		})

		if err != nil {
			return fmt.Errorf("failed to write %s: %w", outputFileName, err)
		}
	}

	return nil
}
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}

func extractImports(fileContent string) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fileContent, parser.ImportsOnly)
	if err != nil {
		return []string{}
	}
	imports := make([]string, 0, len(f.Imports))
	for _, imp := range f.Imports {
		importLine := ""
		if imp.Name != nil {
			importLine += imp.Name.Name + " "
		}
		importLine += imp.Path.Value // includes quotes
		imports = append(imports, importLine)
	}
	return imports
}

func extractHandlerImplementation(fileContent, handlerName string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fileContent, parser.ParseComments)
	if err != nil {
		return "", nil
	}

	var foundBodies []string

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != handlerName {
			continue
		}
		// Check signature: (ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
		if len(fn.Type.Params.List) != 2 || len(fn.Type.Results.List) != 2 {
			continue
		}
		param2 := fn.Type.Params.List[1]
		result1 := fn.Type.Results.List[0]

		if exprToString(param2.Type) != "mcp.CallToolRequest" {
			continue
		}
		if exprToString(result1.Type) != "*mcp.CallToolResult" {
			continue
		}

		if fn.Body != nil {
			start := fset.Position(fn.Body.Lbrace).Offset
			end := fset.Position(fn.Body.Rbrace).Offset
			if start < end && end < len(fileContent) {
				body := fileContent[start : end+1]
				if !strings.HasSuffix(body, "\n") {
					body += "\n"
				}
				foundBodies = append(foundBodies, body)
			}
		}
	}

	if len(foundBodies) == 0 {
		return "", nil
	}
	if len(foundBodies) > 1 {
		return "", fmt.Errorf("multiple handlers named %s with the same signature found in file", handlerName)
	}
	return foundBodies[0], nil
}

func replaceHandlerImplementation(fileContent, handlerName, implementation string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fileContent, parser.ParseComments)
	if err != nil {
		return fileContent
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != handlerName {
			continue
		}
		if fn.Body == nil {
			continue
		}
		start := fset.Position(fn.Body.Lbrace).Offset
		end := fset.Position(fn.Body.Rbrace).Offset
		if start < end && end < len(fileContent) {
			var buf bytes.Buffer
			buf.WriteString(fileContent[:start])
			impl := implementation
			if !strings.HasSuffix(impl, "\n") {
				impl += "\n"
			}
			buf.WriteString(impl)
			buf.WriteString(fileContent[end+1:])
			return buf.String()
		}
	}
	return fileContent
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), expr)
	return buf.String()
}
