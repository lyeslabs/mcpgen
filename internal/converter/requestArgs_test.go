package converter

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestConvertRequestBody_NilAndEmpty(t *testing.T) {
	c := &Converter{}
	arg, err := c.convertRequestBody(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arg != nil {
		t.Errorf("expected nil for nil requestBodyRef, got %+v", arg)
	}

	arg, err = c.convertRequestBody(&openapi3.RequestBodyRef{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arg != nil {
		t.Errorf("expected nil for empty requestBodyRef, got %+v", arg)
	}
}

func TestConvertRequestBody_ValidSingleContentType(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{Title: "BodySchema"}
	body := &openapi3.RequestBody{
		Description: "A body",
		Required:    true,
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{Value: schema},
			},
		},
	}
	arg, err := c.convertRequestBody(&openapi3.RequestBodyRef{Value: body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arg == nil {
		t.Fatal("expected non-nil Arg")
	}
	if arg.Name != "body" || arg.Source != "body" || arg.Description != "A body" || !arg.Required {
		t.Errorf("unexpected Arg fields: %+v", arg)
	}
	if len(arg.ContentTypes) != 1 {
		t.Errorf("expected 1 content type, got %d", len(arg.ContentTypes))
	}
	if arg.ContentTypes["application/json"] == nil || arg.ContentTypes["application/json"].Title != "BodySchema" {
		t.Errorf("expected schema with title 'BodySchema', got %+v", arg.ContentTypes["application/json"])
	}
}

func TestConvertRequestBody_MultipleContentTypes(t *testing.T) {
	c := &Converter{}
	schema1 := &openapi3.Schema{Title: "Schema1"}
	schema2 := &openapi3.Schema{Title: "Schema2"}
	body := &openapi3.RequestBody{
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{Value: schema1},
			},
			"text/plain": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{Value: schema2},
			},
		},
	}
	arg, err := c.convertRequestBody(&openapi3.RequestBodyRef{Value: body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arg == nil {
		t.Fatal("expected non-nil Arg")
	}
	if len(arg.ContentTypes) != 2 {
		t.Errorf("expected 2 content types, got %d", len(arg.ContentTypes))
	}
	if arg.ContentTypes["application/json"] == nil || arg.ContentTypes["application/json"].Title != "Schema1" {
		t.Errorf("expected schema with title 'Schema1', got %+v", arg.ContentTypes["application/json"])
	}
	if arg.ContentTypes["text/plain"] == nil || arg.ContentTypes["text/plain"].Title != "Schema2" {
		t.Errorf("expected schema with title 'Schema2', got %+v", arg.ContentTypes["text/plain"])
	}
}

func TestConvertRequestBody_InvalidContentType(t *testing.T) {
	c := &Converter{}
	body := &openapi3.RequestBody{
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: nil, // Should be skipped
			},
		},
	}
	arg, err := c.convertRequestBody(&openapi3.RequestBodyRef{Value: body})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arg != nil {
		t.Errorf("expected nil Arg for invalid content type, got %+v", arg)
	}
}

func TestConvertParameters_Empty(t *testing.T) {
	c := &Converter{}
	args, err := c.convertParameters(openapi3.Parameters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestConvertParameters_Valid(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{Title: "ParamSchema"}
	param := &openapi3.Parameter{
		Name:        "id",
		Description: "The ID",
		In:          "query",
		Required:    true,
		Deprecated:  true,
		Schema:      &openapi3.SchemaRef{Value: schema},
	}
	args, err := c.convertParameters(openapi3.Parameters{
		&openapi3.ParameterRef{Value: param},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
	a := args[0]
	if a.Name != "id" || a.Description != "The ID" || a.Source != "query" || !a.Required || !a.Deprecated {
		t.Errorf("unexpected Arg fields: %+v", a)
	}
	if a.Schema == nil || a.Schema.Title != "ParamSchema" {
		t.Errorf("expected schema with title 'ParamSchema', got %+v", a.Schema)
	}
}

func TestConvertParameters_SkipInvalid(t *testing.T) {
	c := &Converter{}
	// Nil paramRef, nil Value, nil Schema, nil Schema.Value
	params := openapi3.Parameters{
		nil,
		&openapi3.ParameterRef{},
		&openapi3.ParameterRef{Value: &openapi3.Parameter{}},
		&openapi3.ParameterRef{Value: &openapi3.Parameter{Schema: &openapi3.SchemaRef{}}},
	}
	args, err := c.convertParameters(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}
