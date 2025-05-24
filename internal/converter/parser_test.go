package converter

import (
	"os"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)


func getTestSpecPath(t *testing.T) string {
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Fatalf("Test setup error: fixture file %s does not exist. Please create it.", specPath)
	}
	return specPath
}

func getTestSpecBytes(t *testing.T) []byte {
	specPath := getTestSpecPath(t)
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("Test setup error: could not read %s: %v", specPath, err)
	}
	return data
}

func TestNewParser(t *testing.T) {
	p := NewParser(true)
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
	if !p.ValidateDocument {
		t.Error("expected ValidateDocument to be true")
	}
}

func TestParser_Parse_ValidYAML(t *testing.T) {
	p := NewParser(false)
	data := getTestSpecBytes(t)
	err := p.Parse(data)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	doc := p.GetDocument()
	if doc == nil {
		t.Fatal("expected parsed document")
	}
	if doc.Info == nil || doc.Info.Title == "" {
		t.Errorf("expected non-empty info title")
	}
}

func TestParser_ParseFile(t *testing.T) {
	p := NewParser(false)
	specPath := getTestSpecPath(t)
	err := p.ParseFile(specPath)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	doc := p.GetDocument()
	if doc == nil {
		t.Fatal("expected parsed document from file")
	}
}

func TestParser_GetDocument(t *testing.T) {
	p := NewParser(false)
	_ = p.Parse(getTestSpecBytes(t))
	doc := p.GetDocument()
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
}

func TestParser_GetPaths(t *testing.T) {
	p := NewParser(false)
	_ = p.Parse(getTestSpecBytes(t))
	paths := p.GetPaths()
	if paths == nil {
		t.Fatal("expected non-nil paths")
	}
	// Optionally check for a known path in your simple_openapi.yaml
	// if _, ok := paths["/hello"]; !ok {
	// 	t.Errorf("expected /hello path")
	// }
}

func TestParser_GetServers(t *testing.T) {
	p := NewParser(false)
	_ = p.Parse(getTestSpecBytes(t))
	servers := p.GetServers()
	if len(servers) == 0 {
		t.Fatal("expected at least one server")
	}
}

func TestParser_GetInfo(t *testing.T) {
	p := NewParser(false)
	_ = p.Parse(getTestSpecBytes(t))
	info := p.GetInfo()
	if info == nil {
		t.Fatal("expected non-nil info")
	}
}

func TestParser_GetOperationID_WithExplicitID(t *testing.T) {
	p := NewParser(false)
	op := &openapi3.Operation{OperationID: "explicitID"}
	id := p.GetOperationID("/foo", "GET", op)
	if id != "explicitID" {
		t.Errorf("expected explicitID, got %q", id)
	}
}

func TestParser_GetOperationID_Generated(t *testing.T) {
	p := NewParser(false)
	op := &openapi3.Operation{}
	id := p.GetOperationID("/foo/{bar}", "POST", op)
	if !strings.HasPrefix(id, "post_foo_bar") {
		t.Errorf("expected generated id to start with post_foo_bar, got %q", id)
	}
}
