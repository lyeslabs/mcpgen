package converter

import (
	"reflect"
	"testing"
)

func TestBuildPropertySchema(t *testing.T) {
	// Minimal schema for testing
	s := &Schema{
		Title:       "Test",
		Description: "A test schema",
		Types:       []string{"string"},
	}

	t.Run("default with schema", func(t *testing.T) {
		arg := Arg{
			Source:      "query",
			Schema:      s,
			Description: "desc",
		}
		got, err := buildPropertySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatalf("expected non-nil schema")
		}
		if got["title"] != "Test" {
			t.Errorf("title = %v, want Test", got["title"])
		}
		// Should not overwrite existing description
		if got["description"] != "A test schema" {
			t.Errorf("description = %v, want 'A test schema'", got["description"])
		}
	})

	t.Run("default with schema and empty description", func(t *testing.T) {
		s2 := &Schema{
			Title:       "Test2",
			Description: "",
			Types:       []string{"string"},
		}
		arg := Arg{
			Source:      "query",
			Schema:      s2,
			Description: "desc2",
		}
		got, err := buildPropertySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["description"] != "desc2" {
			t.Errorf("description = %v, want 'desc2'", got["description"])
		}
	})

	t.Run("default with nil schema", func(t *testing.T) {
		arg := Arg{
			Source: "query",
			Schema: nil,
		}
		got, err := buildPropertySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("body with one content type", func(t *testing.T) {
		arg := Arg{
			Source: "body",
			ContentTypes: map[string]*Schema{
				"application/json": s,
			},
		}
		got, err := buildPropertySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["title"] != "Test" {
			t.Errorf("title = %v, want Test", got["title"])
		}
	})

	t.Run("body with multiple content types", func(t *testing.T) {
		s2 := &Schema{Title: "Other", Types: []string{"integer"}}
		arg := Arg{
			Source: "body",
			ContentTypes: map[string]*Schema{
				"application/json": s,
				"application/xml":  s2,
			},
		}
		got, err := buildPropertySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		oneOf, ok := got["oneOf"].([]map[string]interface{})
		if !ok || len(oneOf) != 2 {
			t.Fatalf("expected oneOf with 2 schemas, got %v", got["oneOf"])
		}
		// Check content type info is added
		found := false
		for _, sch := range oneOf {
			if title, ok := sch["title"].(string); ok && title == "[application/xml] Other" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected content type info in title for xml branch, got %v", got)
		}
	})

	t.Run("body with no content types", func(t *testing.T) {
		arg := Arg{
			Source:       "body",
			ContentTypes: map[string]*Schema{},
		}
		got, err := buildPropertySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestBuildBodySchema(t *testing.T) {
	s := &Schema{
		Title: "Test",
		Types: []string{"string"},
	}
	s2 := &Schema{
		Title: "Other",
		Types: []string{"integer"},
	}

	t.Run("no content types", func(t *testing.T) {
		arg := Arg{ContentTypes: map[string]*Schema{}}
		got, err := buildBodySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	t.Run("one content type", func(t *testing.T) {
		arg := Arg{ContentTypes: map[string]*Schema{
			"application/json": s,
		}}
		got, err := buildBodySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got["title"] != "Test" {
			t.Errorf("title = %v, want Test", got["title"])
		}
		if got["type"] != "string" {
			t.Errorf("type = %v, want string", got["type"])
		}
	})

	t.Run("multiple content types", func(t *testing.T) {
		arg := Arg{ContentTypes: map[string]*Schema{
			"application/json": s,
			"application/xml":  s2,
		}}
		got, err := buildBodySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		oneOf, ok := got["oneOf"].([]map[string]interface{})
		if !ok || len(oneOf) != 2 {
			t.Fatalf("expected oneOf with 2 schemas, got %v", got["oneOf"])
		}
		found := false
		for _, sch := range oneOf {
			if title, ok := sch["title"].(string); ok && title == "[application/xml] Other" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected content type info in title for xml branch, got %v", got)
		}
	})

	t.Run("multiple content types with nil schema", func(t *testing.T) {
		arg := Arg{ContentTypes: map[string]*Schema{
			"application/json": s,
			"application/xml":  nil,
		}}
		got, err := buildBodySchema(arg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		oneOf, ok := got["oneOf"].([]map[string]interface{})
		if !ok || len(oneOf) != 1 {
			t.Fatalf("expected oneOf with 1 schema, got %v", got["oneOf"])
		}
		if oneOf[0]["title"] != "[application/json] Test" {
			t.Errorf("title = %v, want [application/json] Test", oneOf[0]["title"])
		}
	})
}

func TestAddContentTypeInfo(t *testing.T) {
	t.Run("adds to description", func(t *testing.T) {
		schema := map[string]interface{}{
			"description": "A description",
		}
		addContentTypeInfo(schema, "application/json")
		want := map[string]interface{}{
			"description": "[application/json] A description",
		}
		if !reflect.DeepEqual(schema, want) {
			t.Errorf("got %v, want %v", schema, want)
		}
	})

	t.Run("adds to title if no description", func(t *testing.T) {
		schema := map[string]interface{}{
			"title": "A title",
		}
		addContentTypeInfo(schema, "application/xml")
		want := map[string]interface{}{
			"title": "[application/xml] A title",
		}
		if !reflect.DeepEqual(schema, want) {
			t.Errorf("got %v, want %v", schema, want)
		}
	})

	t.Run("sets title if neither present", func(t *testing.T) {
		schema := map[string]interface{}{}
		addContentTypeInfo(schema, "text/plain")
		want := map[string]interface{}{
			"title": "Schema for text/plain",
		}
		if !reflect.DeepEqual(schema, want) {
			t.Errorf("got %v, want %v", schema, want)
		}
	})
}
