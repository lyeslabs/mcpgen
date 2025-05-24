package converter

import (
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestCreateResponseTemplates(t *testing.T) {
	c := &Converter{}

	// Build a fake operation with two responses, each with two content types
	op := &openapi3.Operation{
		Responses: openapi3.NewResponses(),
	}

	// 200 response with application/json and text/plain (both with schemas)
	schemaJSON := &openapi3.Schema{Description: "JSON schema"}
	schemaPlain := &openapi3.Schema{Description: "Plain schema"}
	op.Responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: func(s string) *string { return &s }("OK"),
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{Value: schemaJSON},
				},
				"text/plain": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{Value: schemaPlain},
				},
			},
		},
	})

	// 404 response with only application/json (with schema)
	schema404 := &openapi3.Schema{Description: "Not found schema"}
	op.Responses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: func(s string) *string { return &s }("Not found"),
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{Value: schema404},
				},
			},
		},
	})

	// 500 response with no schema (should be skipped)
	op.Responses.Set("500", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: func(s string) *string { return &s }("Error"),
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					// No schema
				},
			},
		},
	})

	templates, err := c.createResponseTemplates(op)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 3 templates: 200/json, 200/plain, 404/json
	if len(templates) != 3 {
		t.Fatalf("expected 3 templates, got %d", len(templates))
	}

	// Check that all expected combinations are present
	found := map[string]bool{}
	for _, tpl := range templates {
		key := tpl.StatusCode
		ct := tpl.ContentType
		body := tpl.PrependBody
		switch {
		case key == 200 && ct == "application/json":
			if !strings.Contains(body, "JSON schema") {
				t.Errorf("expected JSON schema in body, got: %q", body)
			}
			found["200-json"] = true
		case key == 200 && ct == "text/plain":
			if !strings.Contains(body, "Plain schema") {
				t.Errorf("expected Plain schema in body, got: %q", body)
			}
			found["200-plain"] = true
		case key == 404 && ct == "application/json":
			if !strings.Contains(body, "Not found schema") {
				t.Errorf("expected Not found schema in body, got: %q", body)
			}
			found["404-json"] = true
		default:
			t.Errorf("unexpected template: status=%d, contentType=%s", key, ct)
		}
	}

	for _, k := range []string{"200-json", "200-plain", "404-json"} {
		if !found[k] {
			t.Errorf("missing template for %s", k)
		}
	}
}
