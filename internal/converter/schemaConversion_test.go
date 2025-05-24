package converter

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGenerateJSONSchemaDraft7(t *testing.T) {
	// Minimal schemas for testing
	s1 := &Schema{
		Types: []string{"string"},
		Title: "First",
	}
	s2 := &Schema{
		Types: []string{"integer"},
		Title: "Second",
	}

	args := []Arg{
		{
			Name:     "foo",
			Required: true,
			Schema:   s1,
		},
		{
			Name:     "bar",
			Required: false,
			Schema:   s2,
		},
		{
			Name:     "skipme",
			Required: true,
			Schema:   nil, // Should be skipped
		},
	}

	got, err := GenerateJSONSchemaDraft7(args)
	if err != nil {
		t.Fatalf("GenerateJSONSchemaDraft7() error = %v", err)
	}

	// Unmarshal the result for easier assertions
	var gotMap map[string]interface{}
	if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}

	// Check root type
	if gotMap["type"] != "object" {
		t.Errorf("root type = %v, want object", gotMap["type"])
	}

	// Check properties
	props, ok := gotMap["properties"].(map[string]interface{})
	if !ok {
		t.Fatalf("properties missing or wrong type: %v", gotMap["properties"])
	}
	if _, ok := props["foo"]; !ok {
		t.Errorf("missing property foo")
	}
	if _, ok := props["bar"]; !ok {
		t.Errorf("missing property bar")
	}
	if _, ok := props["skipme"]; ok {
		t.Errorf("property skipme should be skipped")
	}

	// Check required
	req, ok := gotMap["required"].([]interface{})
	if !ok {
		t.Fatalf("required missing or wrong type: %v", gotMap["required"])
	}
	wantRequired := []interface{}{"foo"}
	if !reflect.DeepEqual(req, wantRequired) {
		t.Errorf("required = %v, want %v", req, wantRequired)
	}

	// Check property types
	foo := props["foo"].(map[string]interface{})
	if foo["type"] != "string" {
		t.Errorf("foo.type = %v, want string", foo["type"])
	}
	bar := props["bar"].(map[string]interface{})
	if bar["type"] != "integer" {
		t.Errorf("bar.type = %v, want integer", bar["type"])
	}
}
