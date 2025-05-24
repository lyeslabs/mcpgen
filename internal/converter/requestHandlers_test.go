package converter

import (
	"os"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGetOperations_AllMethods(t *testing.T) {
	op := &openapi3.Operation{Summary: "test"}
	pathItem := &openapi3.PathItem{
		Get:     op,
		Post:    op,
		Put:     op,
		Delete:  op,
		Options: op,
		Head:    op,
		Patch:   op,
		Trace:   op,
	}
	ops := getOperations(pathItem)
	expected := []string{"get", "post", "put", "delete", "options", "head", "patch", "trace"}
	for _, method := range expected {
		if ops[method] != op {
			t.Errorf("expected %s operation to be set", method)
		}
	}
	if len(ops) != len(expected) {
		t.Errorf("expected %d operations, got %d", len(expected), len(ops))
	}
}

func TestGetOperations_None(t *testing.T) {
	ops := getOperations(&openapi3.PathItem{})
	if len(ops) != 0 {
		t.Errorf("expected 0 operations, got %d", len(ops))
	}
}

func TestConvertOperation_RealData(t *testing.T) {
	// Load a real OpenAPI spec (use your tested Parser)
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Fatalf("Test setup error: fixture file %s does not exist. Please create it.", specPath)
	}
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("could not read %s: %v", specPath, err)
	}
	parser := NewParser(false)
	if err := parser.Parse(data); err != nil {
		t.Fatalf("failed to parse OpenAPI: %v", err)
	}

	c := &Converter{parser: parser}

	// Pick the first path and method from the spec
	paths := parser.GetPaths()
	if len(paths) == 0 {
		t.Fatal("no paths found in OpenAPI spec")
	}
	var path string
	var pathItem *openapi3.PathItem
	for p, pi := range paths {
		path = p
		pathItem = pi
		break
	}
	ops := getOperations(pathItem)
	if len(ops) == 0 {
		t.Fatal("no operations found in path item")
	}
	var method string
	var operation *openapi3.Operation
	for m, op := range ops {
		method = m
		operation = op
		break
	}

	tool, err := c.convertOperation(path, method, operation)
	if err != nil {
		t.Fatalf("convertOperation failed: %v", err)
	}
	if tool == nil {
		t.Fatal("expected non-nil tool")
	}
	if tool.Name == "" {
		t.Error("expected tool name to be set")
	}
	if !strings.Contains(strings.ToLower(tool.Description), "test") && tool.Description == "" {
		t.Error("expected tool description to be set")
	}
	if tool.RequestTemplate.URL == "" {
		t.Error("expected request template URL to be set")
	}
	if len(tool.Responses) == 0 {
		t.Error("expected at least one response template")
	}
	// Args and RawInputSchema are optional, but you can check them if you want
}
