package converter

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestSortedResponseCodes(t *testing.T) {
	responses := &openapi3.Responses{}
	responses.Set("404", &openapi3.ResponseRef{})
	responses.Set("2XX", &openapi3.ResponseRef{})
	responses.Set("200", &openapi3.ResponseRef{})
	responses.Set("default", &openapi3.ResponseRef{})
	responses.Set("500", &openapi3.ResponseRef{})

	got := sortedResponseCodes(responses)
	want := []string{"200", "404", "500", "2XX", "default"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortedResponseCodes() = %v, want %v", got, want)
	}
}

func TestSortedContentTypes(t *testing.T) {
	content := openapi3.Content{
		"application/json":         &openapi3.MediaType{},
		"application/xml":          &openapi3.MediaType{},
		"text/plain":               &openapi3.MediaType{},
		"application/octet-stream": &openapi3.MediaType{},
	}

	got := sortedContentTypes(content)
	want := []string{
		"application/json",
		"application/octet-stream",
		"application/xml",
		"text/plain",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortedContentTypes() = %v, want %v", got, want)
	}
}

func TestToAlphaSuffix(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
		{701, "ZZ"},
		{702, "AAA"},
	}

	for _, c := range cases {
		got := toAlphaSuffix(c.n)
		if got != c.want {
			t.Errorf("toAlphaSuffix(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestAssignSuffixes(t *testing.T) {
	responses := make([]ResponseTemplate, 5)
	wantSuffixes := []string{"A", "B", "C", "D", "E"}

	got := assignSuffixes(responses)
	for i, resp := range got {
		if resp.Suffix != wantSuffixes[i] {
			t.Errorf("assignSuffixes: response %d has Suffix %q, want %q", i, resp.Suffix, wantSuffixes[i])
		}
	}

	// Test with more than 26 to check double letters
	responses = make([]ResponseTemplate, 28)
	got = assignSuffixes(responses)
	if got[26].Suffix != "AA" || got[27].Suffix != "AB" {
		t.Errorf("assignSuffixes: got[26]=%q, got[27]=%q, want AA, AB", got[26].Suffix, got[27].Suffix)
	}
}

func TestFormatForGoRawString(t *testing.T) {
	strType := openapi3.Types{"string"}
	intType := openapi3.Types{"integer"}
	boolType := openapi3.Types{"boolean"}

	cases := []struct {
		name   string
		schema *openapi3.Schema
		value  interface{}
		want   string
	}{
		{
			name:   "simple string",
			schema: &openapi3.Schema{Type: &strType},
			value:  "hello",
			want:   "hello",
		},
		{
			name:   "quoted string",
			schema: &openapi3.Schema{Type: &strType},
			value:  `"hello"`,
			want:   "hello",
		},
		{
			name:   "string with backtick",
			schema: &openapi3.Schema{Type: &strType},
			value:  "foo`bar",
			want:   "foo'bar",
		},
		{
			name:   "integer value",
			schema: &openapi3.Schema{Type: &intType},
			value:  42,
			want:   "42",
		},
		{
			name:   "boolean value",
			schema: &openapi3.Schema{Type: &boolType},
			value:  true,
			want:   "true",
		},
		{
			name:   "array value",
			schema: &openapi3.Schema{Type: &strType},
			value:  []string{"a", "b"},
			want:   `["a","b"]`,
		},
		{
			name:   "map value",
			schema: &openapi3.Schema{Type: &strType},
			value:  map[string]int{"a": 1, "b": 2},
			want:   `{"a":1,"b":2}`,
		},
		{
			name:   "string with quotes",
			schema: &openapi3.Schema{Type: &strType},
			value:  `"foo"`,
			want:   "foo",
		},
		{
			name:   "string with backtick and quotes",
			schema: &openapi3.Schema{Type: &strType},
			value:  "`foo`",
			want:   "'foo'",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := formatForGoRawString(c.schema, c.value)
			if got != c.want {
				t.Errorf("formatForGoRawString(%v, %v) = %q, want %q", c.schema.Type, c.value, got, c.want)
			}
		})
	}
}

func TestGetResponseDescription(t *testing.T) {
	desc := "A description"
	cases := []struct {
		name     string
		response *openapi3.ResponseRef
		want     string
	}{
		{
			name: "has description",
			response: &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: &desc,
				},
			},
			want: "A description",
		},
		{
			name: "nil description",
			response: &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: nil,
				},
			},
			want: "",
		},
		{
			name: "nil value",
			response: &openapi3.ResponseRef{
				Value: nil,
			},
			want: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ""
			if c.response != nil && c.response.Value != nil {
				got = getResponseDescription(c.response)
			}
			if got != c.want {
				t.Errorf("getResponseDescription(%v) = %q, want %q", c.response, got, c.want)
			}
		})
	}
}

func TestGetDescription(t *testing.T) {
	cases := []struct {
		name      string
		operation *openapi3.Operation
		want      string
	}{
		{
			name: "summary and description",
			operation: &openapi3.Operation{
				Summary:     "Short",
				Description: "Long",
			},
			want: "Short - Long",
		},
		{
			name: "summary only",
			operation: &openapi3.Operation{
				Summary: "Short",
			},
			want: "Short",
		},
		{
			name: "description only",
			operation: &openapi3.Operation{
				Description: "Long",
			},
			want: "Long",
		},
		{
			name:      "neither",
			operation: &openapi3.Operation{},
			want:      "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := getDescription(c.operation)
			if got != c.want {
				t.Errorf("getDescription(%v) = %q, want %q", c.operation, got, c.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		name  string
		slice []string
		str   string
		want  bool
	}{
		{"found", []string{"a", "b", "c"}, "b", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty slice", []string{}, "a", false},
		{"empty string in slice", []string{"", "b"}, "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := contains(c.slice, c.str)
			if got != c.want {
				t.Errorf("contains(%v, %q) = %v, want %v", c.slice, c.str, got, c.want)
			}
		})
	}
}
