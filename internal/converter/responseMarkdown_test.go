package converter

import (
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestWriteSchemaDetails_StringValidations(t *testing.T) {
	c := &Converter{}
	maxLen := uint64(8)
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type:      &schemaType,
		MinLength: 2,
		MaxLength: &maxLen,
		Pattern:   "foo`bar",
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "Min Length: 2") {
		t.Errorf("expected Min Length, got: %q", out)
	}
	if !strings.Contains(out, "Max Length: 8") {
		t.Errorf("expected Max Length, got: %q", out)
	}
	if !strings.Contains(out, "Pattern: 'foo'bar'") {
		t.Errorf("expected Pattern with backtick replaced, got: %q", out)
	}
}

func TestWriteSchemaDetails_NumericValidations(t *testing.T) {
	c := &Converter{}
	min := 1.5
	max := 10.0
	mult := 2.0
	schemaType := openapi3.Types{"number"}
	schema := &openapi3.Schema{
		Type:         &schemaType,
		Min:          &min,
		Max:          &max,
		ExclusiveMin: true,
		ExclusiveMax: true,
		MultipleOf:   &mult,
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 1)
	out := b.String()
	for _, want := range []string{
		"Minimum: 1.5",
		"Maximum: 10",
		"Exclusive Minimum: true",
		"Exclusive Maximum: true",
		"Multiple Of: 2",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestWriteSchemaDetails_ArrayValidations(t *testing.T) {
	c := &Converter{}
	maxItems := uint64(5)
	schemaType := openapi3.Types{"array"}
	schema := &openapi3.Schema{
		Type:        &schemaType,
		MinItems:    1,
		MaxItems:    &maxItems,
		UniqueItems: true,
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 2)
	out := b.String()
	for _, want := range []string{
		"Min Items: 1",
		"Max Items: 5",
		"Unique Items: true",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestWriteSchemaDetails_Nullable(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		Nullable: true,
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "Nullable: true") {
		t.Errorf("expected Nullable: true, got: %q", out)
	}
}

func TestWriteSchemaDetails_DefaultAndExample_StringType(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type:    &schemaType,
		Default: "foo",
		Example: "bar",
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "Default: 'foo'") {
		t.Errorf("expected Default: 'foo', got: %q", out)
	}
	if !strings.Contains(out, "Example: 'bar'") {
		t.Errorf("expected Example: 'bar', got: %q", out)
	}
}

func TestWriteSchemaDetails_DefaultAndExample_NonStringType(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"integer"}
	schema := &openapi3.Schema{
		Type:    &schemaType,
		Default: 42,
		Example: 99,
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "Default: '42'") {
		t.Errorf("expected Default: '42', got: %q", out)
	}
	if !strings.Contains(out, "Example: '99'") {
		t.Errorf("expected Example: '99', got: %q", out)
	}
}

func TestWriteSchemaDetails_Enum_StringType(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type: &schemaType,
		Enum: []interface{}{"a", "b", "foo`bar"},
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "Enum: ['a', 'b', 'foo'bar']") {
		t.Errorf("expected Enum: ['a', 'b', 'foo'bar'], got: %q", out)
	}
}

func TestWriteSchemaDetails_Enum_NonStringType(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"object"}
	schema := &openapi3.Schema{
		Type: &schemaType,
		Enum: []interface{}{
			map[string]interface{}{"a": 1},
			map[string]interface{}{"b": 2},
		},
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, `Enum: ['{"a":1}', '{"b":2}']`) {
		t.Errorf("expected Enum: ['{\"a\":1}', '{\"b\":2}'], got: %q", out)
	}
}

func TestWriteSchemaDetails_EmptySchema(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 0)
	out := b.String()
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}

func TestWriteSchemaDetails_Indentation(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type:      &schemaType,
		MinLength: 1,
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 2)
	out := b.String()
	// Should be 4 spaces for indent=2, plus 2 more for details
	if !strings.HasPrefix(out, "      - ") {
		t.Errorf("expected 6 spaces indentation, got: %q", out)
	}
}

func TestWriteSchemaDetails_AllCombined(t *testing.T) {
	c := &Converter{}
	maxLen := uint64(10)
	min := 1.0
	max := 5.0
	mult := 2.0
	maxItems := uint64(3)
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type:         &schemaType,
		MinLength:    2,
		MaxLength:    &maxLen,
		Pattern:      "abc",
		Min:          &min,
		Max:          &max,
		ExclusiveMin: true,
		ExclusiveMax: true,
		MultipleOf:   &mult,
		MinItems:     1,
		MaxItems:     &maxItems,
		UniqueItems:  true,
		Nullable:     true,
		Default:      "foo",
		Example:      "bar",
		Enum:         []interface{}{"a", "b"},
	}
	var b strings.Builder
	c.writeSchemaDetails(&b, schema, 1)
	out := b.String()
	expected := []string{
		"Min Length: 2",
		"Max Length: 10",
		"Pattern: 'abc'",
		"Minimum: 1",
		"Maximum: 5",
		"Exclusive Minimum: true",
		"Exclusive Maximum: true",
		"Multiple Of: 2",
		"Min Items: 1",
		"Max Items: 3",
		"Unique Items: true",
		"Nullable: true",
		"Default: 'foo'",
		"Example: 'bar'",
		"Enum: ['a', 'b']",
	}
	for _, want := range expected {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestWriteAdditionalProperties_WithSchema(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"object"}
	propSchema := &openapi3.Schema{
		Type:        &openapi3.Types{"string"},
		Description: "Extra property value",
	}
	schema := &openapi3.Schema{
		Type: &schemaType,
		AdditionalProperties: openapi3.AdditionalProperties{
			Schema: &openapi3.SchemaRef{Value: propSchema},
		},
	}

	var b strings.Builder
	c.writeAdditionalProperties(&b, schema, 1)
	out := b.String()
	if !strings.Contains(out, "  - **Additional Properties**:") {
		t.Errorf("expected Additional Properties header, got: %q", out)
	}
	if !strings.Contains(out, "Extra property value") {
		t.Errorf("expected property value schema to be documented, got: %q", out)
	}
}

func TestWriteAdditionalProperties_WithHasTrue(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"object"}
	has := true
	schema := &openapi3.Schema{
		Type: &schemaType,
		AdditionalProperties: openapi3.AdditionalProperties{
			Has: &has,
		},
	}
	var b strings.Builder
	c.writeAdditionalProperties(&b, schema, 2)
	out := b.String()
	if !strings.Contains(out, "    - **Allows Additional Properties**") {
		t.Errorf("expected Allows Additional Properties, got: %q", out)
	}
}

func TestWriteAdditionalProperties_NonObject(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type: &schemaType,
	}
	var b strings.Builder
	c.writeAdditionalProperties(&b, schema, 0)
	out := b.String()
	if out != "" {
		t.Errorf("expected no output for non-object, got: %q", out)
	}
}

func TestWriteAdditionalProperties_NilAndFalse(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"object"}
	has := false
	schema := &openapi3.Schema{
		Type: &schemaType,
		AdditionalProperties: openapi3.AdditionalProperties{
			Has: &has,
		},
	}
	var b strings.Builder
	c.writeAdditionalProperties(&b, schema, 0)
	out := b.String()
	if out != "" {
		t.Errorf("expected no output for Has=false, got: %q", out)
	}

	// Also test with zero-value AdditionalProperties
	schema2 := &openapi3.Schema{
		Type: &schemaType,
	}
	var b2 strings.Builder
	c.writeAdditionalProperties(&b2, schema2, 0)
	out2 := b2.String()
	if out2 != "" {
		t.Errorf("expected no output for zero-value AdditionalProperties, got: %q", out2)
	}
}


func TestWriteSchemaCombinators_OneOf(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		OneOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Description: "Option A"}},
			{Value: &openapi3.Schema{Description: "Option B"}},
		},
	}
	var b strings.Builder
	c.writeSchemaCombinators(&b, schema, 1)
	out := b.String()
	if !strings.Contains(out, "  - **One Of the following structures**:") {
		t.Errorf("expected One Of header, got: %q", out)
	}
	if !strings.Contains(out, "Option A") || !strings.Contains(out, "Option B") {
		t.Errorf("expected both options' descriptions, got: %q", out)
	}
}

func TestWriteSchemaCombinators_AnyOf(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		AnyOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Description: "Any A"}},
			{Value: &openapi3.Schema{Description: "Any B"}},
		},
	}
	var b strings.Builder
	c.writeSchemaCombinators(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "- **Any Of the following structures**:") {
		t.Errorf("expected Any Of header, got: %q", out)
	}
	if !strings.Contains(out, "Any A") || !strings.Contains(out, "Any B") {
		t.Errorf("expected both anyOf descriptions, got: %q", out)
	}
}

func TestWriteSchemaCombinators_AllOf(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		AllOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Description: "All A"}},
			{Value: &openapi3.Schema{Description: "All B"}},
		},
	}
	var b strings.Builder
	c.writeSchemaCombinators(&b, schema, 2)
	out := b.String()
	if !strings.Contains(out, "    - **Combines All Of the following structures**:") {
		t.Errorf("expected All Of header, got: %q", out)
	}
	if !strings.Contains(out, "All A") || !strings.Contains(out, "All B") {
		t.Errorf("expected both allOf descriptions, got: %q", out)
	}
}

func TestWriteSchemaCombinators_Not(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		Not: &openapi3.SchemaRef{
			Value: &openapi3.Schema{Description: "Forbidden!"},
		},
	}
	var b strings.Builder
	c.writeSchemaCombinators(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "- **Not**: Cannot be the following structure:") {
		t.Errorf("expected Not header, got: %q", out)
	}
	if !strings.Contains(out, "Forbidden!") {
		t.Errorf("expected Not description, got: %q", out)
	}
}

func TestWriteSchemaCombinators_AllCombinators(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		OneOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Description: "OneOf"}},
		},
		AnyOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Description: "AnyOf"}},
		},
		AllOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Description: "AllOf"}},
		},
		Not: &openapi3.SchemaRef{
			Value: &openapi3.Schema{Description: "Not"},
		},
	}
	var b strings.Builder
	c.writeSchemaCombinators(&b, schema, 0)
	out := b.String()
	for _, want := range []string{
		"**One Of the following structures**",
		"OneOf",
		"**Any Of the following structures**",
		"AnyOf",
		"**Combines All Of the following structures**",
		"AllOf",
		"**Not**: Cannot be the following structure:",
		"Not",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestWriteSchemaCombinators_None(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{}
	var b strings.Builder
	c.writeSchemaCombinators(&b, schema, 0)
	out := b.String()
	if out != "" {
		t.Errorf("expected no output for empty combinators, got: %q", out)
	}
}


func TestWriteSchemaProperties_ObjectProperties(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"object"}
	schema := &openapi3.Schema{
		Type: &schemaType,
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: &openapi3.Schema{Description: "Foo property"}},
			"bar": {Value: &openapi3.Schema{Description: "Bar property"}},
		},
	}
	var b strings.Builder
	c.writeSchemaProperties(&b, schema, 0)
	out := b.String()
	if !strings.Contains(out, "Foo property") {
		t.Errorf("expected Foo property description, got: %q", out)
	}
	if !strings.Contains(out, "Bar property") {
		t.Errorf("expected Bar property description, got: %q", out)
	}
}

func TestWriteSchemaProperties_ArrayItems(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"array"}
	schema := &openapi3.Schema{
		Type:  &schemaType,
		Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Description: "Array item schema"}},
	}
	var b strings.Builder
	c.writeSchemaProperties(&b, schema, 1)
	out := b.String()
	if !strings.Contains(out, "Array item schema") {
		t.Errorf("expected array item schema description, got: %q", out)
	}
}

func TestWriteSchemaProperties_ObjectAndArray(t *testing.T) {
	c := &Converter{}
	objType := openapi3.Types{"object"}
	arrType := openapi3.Types{"array"}
	// Object with properties and array items
	objSchema := &openapi3.Schema{
		Type: &objType,
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: &openapi3.Schema{Description: "Foo property"}},
		},
		Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Description: "Should not appear"}},
	}
	arrSchema := &openapi3.Schema{
		Type:  &arrType,
		Items: &openapi3.SchemaRef{Value: &openapi3.Schema{Description: "Array item schema"}},
		Properties: map[string]*openapi3.SchemaRef{
			"bar": {Value: &openapi3.Schema{Description: "Bar property"}},
		},
	}
	// Test object schema (should only print properties, not items)
	var b1 strings.Builder
	c.writeSchemaProperties(&b1, objSchema, 0)
	out1 := b1.String()
	if !strings.Contains(out1, "Foo property") {
		t.Errorf("expected Foo property description, got: %q", out1)
	}
	if strings.Contains(out1, "Should not appear") {
		t.Errorf("did not expect array item schema in object, got: %q", out1)
	}
	// Test array schema (should only print items, not properties)
	var b2 strings.Builder
	c.writeSchemaProperties(&b2, arrSchema, 0)
	out2 := b2.String()
	if !strings.Contains(out2, "Array item schema") {
		t.Errorf("expected array item schema description, got: %q", out2)
	}
	if strings.Contains(out2, "Bar property") {
		t.Errorf("did not expect object property in array, got: %q", out2)
	}
}

func TestWriteSchemaProperties_Neither(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{}
	var b strings.Builder
	c.writeSchemaProperties(&b, schema, 0)
	out := b.String()
	if out != "" {
		t.Errorf("expected no output for schema with no properties or items, got: %q", out)
	}
}

func TestBuildResponseMarkdown_Basic(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"object"}
	schema := &openapi3.Schema{
		Type:        &schemaType,
		Description: "A test object",
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: &openapi3.Schema{Description: "Foo property", Type: &openapi3.Types{"string"}}},
		},
	}
	resp := &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: func(s string) *string { return &s }("A test response"),
		},
	}
	md := c.buildResponseMarkdown("200", "application/json", resp, schema)
	if !strings.Contains(md, "# API Response Information") {
		t.Errorf("expected header, got: %q", md)
	}
	if !strings.Contains(md, "**Status Code:** 200") {
		t.Errorf("expected status code, got: %q", md)
	}
	if !strings.Contains(md, "**Content-Type:** application/json") {
		t.Errorf("expected content type, got: %q", md)
	}
	if !strings.Contains(md, "> A test response") {
		t.Errorf("expected response description, got: %q", md)
	}
	if !strings.Contains(md, "## Response Structure") {
		t.Errorf("expected response structure header, got: %q", md)
	}
	if !strings.Contains(md, "- A test object (Type: object):") {
		t.Errorf("expected schema description, got: %q", md)
	}
	if !strings.Contains(md, "Foo property") {
		t.Errorf("expected property description, got: %q", md)
	}
}

func TestBuildResponseMarkdown_NoDescription(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"string"}
	schema := &openapi3.Schema{
		Type: &schemaType,
	}
	resp := &openapi3.ResponseRef{
		Value: &openapi3.Response{},
	}
	md := c.buildResponseMarkdown("404", "text/plain", resp, schema)
	if !strings.Contains(md, "**Status Code:** 404") {
		t.Errorf("expected status code, got: %q", md)
	}
	if !strings.Contains(md, "**Content-Type:** text/plain") {
		t.Errorf("expected content type, got: %q", md)
	}
	if !strings.Contains(md, "- Structure (Type: string):") {
		t.Errorf("expected root schema line, got: %q", md)
	}
}

func TestWriteSchemaMarkdown_FieldNameAndDescription(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"integer"}
	schema := &openapi3.Schema{
		Type:        &schemaType,
		Description: "A number",
	}
	var b strings.Builder
	c.writeSchemaMarkdown(&b, schema, 1, "count")
	out := b.String()
	if !strings.Contains(out, "- **count**: A number (Type: integer):") {
		t.Errorf("expected field name and description, got: %q", out)
	}
}

func TestWriteSchemaMarkdown_FieldNameNoDescription(t *testing.T) {
	c := &Converter{}
	schemaType := openapi3.Types{"boolean"}
	schema := &openapi3.Schema{
		Type: &schemaType,
	}
	var b strings.Builder
	c.writeSchemaMarkdown(&b, schema, 2, "flag")
	out := b.String()
	if !strings.Contains(out, "- **flag** (Type: boolean):") {
		t.Errorf("expected field name and type, got: %q", out)
	}
}

func TestWriteSchemaMarkdown_NilSchema(t *testing.T) {
	c := &Converter{}
	var b strings.Builder
	c.writeSchemaMarkdown(&b, nil, 0, "")
	out := b.String()
	if out != "" {
		t.Errorf("expected no output for nil schema, got: %q", out)
	}
}