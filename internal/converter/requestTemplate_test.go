package converter

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestCreateRequestTemplate_Basic(t *testing.T) {
	// Setup a fake OpenAPI doc with a server
	doc := &openapi3.T{
		Servers: openapi3.Servers{
			&openapi3.Server{URL: "http://api.example.com/"},
		},
	}
	parser := &Parser{doc: doc}
	c := &Converter{parser: parser}

	// Operation with a request body and two content types
	op := &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{},
					"text/plain":       &openapi3.MediaType{},
				},
			},
		},
	}

	template, err := c.createRequestTemplate("/v1/hello", "post", op)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// URL should not have double slash
	if template.URL != "http://api.example.com/v1/hello" {
		t.Errorf("expected URL 'http://api.example.com/v1/hello', got %q", template.URL)
	}
	// Method should be uppercase
	if template.Method != "POST" {
		t.Errorf("expected method POST, got %q", template.Method)
	}
	// Should have a Content-Type header (the first one in the map)
	if len(template.Headers) == 0 || template.Headers[0].Key != "Content-Type" {
		t.Errorf("expected Content-Type header, got %+v", template.Headers)
	}
	// Should be one of the content types
	if template.Headers[0].Value != "application/json" && template.Headers[0].Value != "text/plain" {
		t.Errorf("unexpected Content-Type value: %q", template.Headers[0].Value)
	}
}

func TestCreateRequestTemplate_NoRequestBody(t *testing.T) {
	doc := &openapi3.T{
		Servers: openapi3.Servers{
			&openapi3.Server{URL: "http://api.example.com"},
		},
	}
	parser := &Parser{doc: doc}
	c := &Converter{parser: parser}

	op := &openapi3.Operation{} // No request body

	template, err := c.createRequestTemplate("/v1/hello", "get", op)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.URL != "http://api.example.com/v1/hello" {
		t.Errorf("expected URL 'http://api.example.com/v1/hello', got %q", template.URL)
	}
	if template.Method != "GET" {
		t.Errorf("expected method GET, got %q", template.Method)
	}
	if len(template.Headers) != 0 {
		t.Errorf("expected no headers, got %+v", template.Headers)
	}
}

func TestCreateRequestTemplate_ServerNoTrailingSlash(t *testing.T) {
	doc := &openapi3.T{
		Servers: openapi3.Servers{
			&openapi3.Server{URL: "http://api.example.com"},
		},
	}
	parser := &Parser{doc: doc}
	c := &Converter{parser: parser}

	op := &openapi3.Operation{
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{},
				},
			},
		},
	}

	template, err := c.createRequestTemplate("/v1/hello", "put", op)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if template.URL != "http://api.example.com/v1/hello" {
		t.Errorf("expected URL 'http://api.example.com/v1/hello', got %q", template.URL)
	}
	if template.Method != "PUT" {
		t.Errorf("expected method PUT, got %q", template.Method)
	}
	if len(template.Headers) != 1 || template.Headers[0].Key != "Content-Type" || template.Headers[0].Value != "application/json" {
		t.Errorf("expected Content-Type header with application/json, got %+v", template.Headers)
	}
}
