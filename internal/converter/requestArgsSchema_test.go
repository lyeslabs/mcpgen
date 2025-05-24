package converter

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestCreateObjectValidation_NoProperties(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatal("expected non-nil ObjectValidation")
	}
	if obj.Properties != nil {
		t.Errorf("expected nil Properties, got: %+v", obj.Properties)
	}
}

func TestCreateObjectValidation_WithProperties(t *testing.T) {
	c := &Converter{}
	propSchema := &openapi3.Schema{Title: "Prop1"}
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: propSchema},
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.Properties == nil || obj.Properties["foo"] == nil {
		t.Fatalf("expected property 'foo' to be present")
	}
	if obj.Properties["foo"].Title != "Prop1" {
		t.Errorf("expected property title 'Prop1', got %q", obj.Properties["foo"].Title)
	}
}

func TestCreateObjectValidation_RequiredAndMinMax(t *testing.T) {
	c := &Converter{}
	min := uint64(1)
	max := uint64(5)
	schema := &openapi3.Schema{
		Required: []string{"foo", "bar"},
		MinProps: min,
		MaxProps: &max,
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: &openapi3.Schema{}},
			"bar": {Value: &openapi3.Schema{}},
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(obj.Required) != 2 || obj.Required[0] != "foo" || obj.Required[1] != "bar" {
		t.Errorf("expected required [foo bar], got %+v", obj.Required)
	}
	if obj.MinProperties != min {
		t.Errorf("expected MinProperties %d, got %d", min, obj.MinProperties)
	}
	if obj.MaxProperties == nil || *obj.MaxProperties != max {
		t.Errorf("expected MaxProperties %d, got %+v", max, obj.MaxProperties)
	}
}

func TestCreateObjectValidation_PropertyNilSchemaRef(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"foo": nil,
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not panic, should not add property
	if obj.Properties != nil && obj.Properties["foo"] != nil {
		t.Errorf("expected nil or empty property for 'foo', got %+v", obj.Properties["foo"])
	}
}

func TestCreateObjectValidation_PropertyNilValue(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: nil},
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should map to empty schema
	if obj.Properties == nil || obj.Properties["foo"] == nil {
		t.Errorf("expected empty schema for 'foo', got nil")
	}
}

func TestCreateObjectValidation_AdditionalProperties_HasFalse(t *testing.T) {
	c := &Converter{}
	has := false
	schema := &openapi3.Schema{
		AdditionalProperties: openapi3.AdditionalProperties{
			Has: &has,
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !obj.DisallowAdditionalProperties {
		t.Errorf("expected DisallowAdditionalProperties true, got false")
	}
	if obj.AdditionalProperties != nil {
		t.Errorf("expected nil AdditionalProperties, got %+v", obj.AdditionalProperties)
	}
}

func TestCreateObjectValidation_AdditionalProperties_HasTrue(t *testing.T) {
	c := &Converter{}
	has := true
	schema := &openapi3.Schema{
		AdditionalProperties: openapi3.AdditionalProperties{
			Has: &has,
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.DisallowAdditionalProperties {
		t.Errorf("expected DisallowAdditionalProperties false, got true")
	}
	if obj.AdditionalProperties == nil {
		t.Errorf("expected non-nil AdditionalProperties, got nil")
	}
}

func TestCreateObjectValidation_AdditionalProperties_Schema(t *testing.T) {
	c := &Converter{}
	propSchema := &openapi3.Schema{Title: "Extra"}
	schema := &openapi3.Schema{
		AdditionalProperties: openapi3.AdditionalProperties{
			Schema: &openapi3.SchemaRef{Value: propSchema},
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.AdditionalProperties == nil {
		t.Fatalf("expected non-nil AdditionalProperties")
	}
	if obj.AdditionalProperties.Title != "Extra" {
		t.Errorf("expected AdditionalProperties title 'Extra', got %q", obj.AdditionalProperties.Title)
	}
}

func TestCreateObjectValidation_AdditionalProperties_SchemaNilValue(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		AdditionalProperties: openapi3.AdditionalProperties{
			Schema: &openapi3.SchemaRef{Value: nil},
		},
	}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.AdditionalProperties == nil {
		t.Fatalf("expected non-nil AdditionalProperties")
	}
}

func TestCreateObjectValidation_AdditionalProperties_ZeroValue(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{}
	obj, err := c.createObjectValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj.AdditionalProperties != nil {
		t.Errorf("expected nil AdditionalProperties, got %+v", obj.AdditionalProperties)
	}
	if obj.DisallowAdditionalProperties {
		t.Errorf("expected DisallowAdditionalProperties false, got true")
	}
}

func TestCreateStringValidation_Nil(t *testing.T) {
	c := &Converter{}
	if c.createStringValidation(nil) != nil {
		t.Error("expected nil for nil schema")
	}
}

func TestCreateStringValidation_Fields(t *testing.T) {
	c := &Converter{}
	maxLen := uint64(10)
	schema := &openapi3.Schema{
		MinLength: 2,
		MaxLength: &maxLen,
		Pattern:   "abc.*",
	}
	sv := c.createStringValidation(schema)
	if sv == nil {
		t.Fatal("expected non-nil StringValidation")
	}
	if sv.MinLength != 2 {
		t.Errorf("expected MinLength 2, got %d", sv.MinLength)
	}
	if sv.MaxLength == nil || *sv.MaxLength != 10 {
		t.Errorf("expected MaxLength 10, got %+v", sv.MaxLength)
	}
	if sv.Pattern != "abc.*" {
		t.Errorf("expected Pattern 'abc.*', got %q", sv.Pattern)
	}
}

func TestCreateNumberValidation_Nil(t *testing.T) {
	c := &Converter{}
	if c.createNumberValidation(nil) != nil {
		t.Error("expected nil for nil schema")
	}
}

func TestCreateNumberValidation_Fields(t *testing.T) {
	c := &Converter{}
	min := 1.5
	max := 10.0
	mult := 2.0
	schema := &openapi3.Schema{
		Min:          &min,
		Max:          &max,
		MultipleOf:   &mult,
		ExclusiveMin: true,
		ExclusiveMax: true,
	}
	nv := c.createNumberValidation(schema)
	if nv == nil {
		t.Fatal("expected non-nil NumberValidation")
	}
	if nv.Minimum == nil || *nv.Minimum != 1.5 {
		t.Errorf("expected Minimum 1.5, got %+v", nv.Minimum)
	}
	if nv.Maximum == nil || *nv.Maximum != 10.0 {
		t.Errorf("expected Maximum 10.0, got %+v", nv.Maximum)
	}
	if nv.MultipleOf == nil || *nv.MultipleOf != 2.0 {
		t.Errorf("expected MultipleOf 2.0, got %+v", nv.MultipleOf)
	}
	if !nv.ExclusiveMinimum {
		t.Errorf("expected ExclusiveMinimum true, got false")
	}
	if !nv.ExclusiveMaximum {
		t.Errorf("expected ExclusiveMaximum true, got false")
	}
}

func TestCreateArrayValidation_Nil(t *testing.T) {
	c := &Converter{}
	arr, err := c.createArrayValidation(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arr != nil {
		t.Error("expected nil for nil schema")
	}
}
func TestCreateArrayValidation_Fields(t *testing.T) {
	c := &Converter{}
	maxItems := uint64(5)
	itemSchema := &openapi3.Schema{Title: "Item"}
	schema := &openapi3.Schema{
		MinItems:    1,
		MaxItems:    &maxItems,
		UniqueItems: true,
		Items:       &openapi3.SchemaRef{Value: itemSchema},
	}
	arr, err := c.createArrayValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arr == nil {
		t.Fatal("expected non-nil ArrayValidation")
	}
	if arr.MinItems != 1 {
		t.Errorf("expected MinItems 1, got %d", arr.MinItems)
	}
	if arr.MaxItems == nil || *arr.MaxItems != 5 {
		t.Errorf("expected MaxItems 5, got %+v", arr.MaxItems)
	}
	if !arr.UniqueItems {
		t.Errorf("expected UniqueItems true, got false")
	}
	if arr.Items == nil {
		t.Errorf("expected non-nil Items, got nil")
	} else if arr.Items.Title != "Item" {
		t.Errorf("expected Items.Title 'Item', got %q", arr.Items.Title)
	}
}

func TestCreateArrayValidation_ItemsNil(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		Items: nil,
	}
	arr, err := c.createArrayValidation(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if arr == nil {
		t.Fatal("expected non-nil ArrayValidation")
	}
	if arr.Items != nil {
		t.Errorf("expected nil Items, got %+v", arr.Items)
	}
}

func TestApplySchema_NilInput(t *testing.T) {
	c := &Converter{}
	_, err := c.applySchema(nil)
	if err == nil {
		t.Fatal("expected error for nil schema")
	}
}

func TestApplySchema_MetadataFields(t *testing.T) {
	c := &Converter{}
	schema := &openapi3.Schema{
		Title:       "TestTitle",
		Description: "TestDesc",
		Format:      "uuid",
		Enum:        []interface{}{"a", "b"},
		Default:     "foo",
		Example:     "bar",
		ReadOnly:    true,
		WriteOnly:   true,
		Type:        &openapi3.Types{"string"},
	}
	result, err := c.applySchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Title != "TestTitle" ||
		result.Description != "TestDesc" ||
		result.Format != "uuid" ||
		!reflect.DeepEqual(result.Enum, []interface{}{"a", "b"}) ||
		result.Default != "foo" ||
		result.Example != "bar" ||
		!result.ReadOnly ||
		!result.WriteOnly {
		t.Errorf("metadata fields not copied correctly: %+v", result)
	}
}

func TestApplySchema_TypesAndNullable(t *testing.T) {
	c := &Converter{}
	// Type present, nullable false
	schema := &openapi3.Schema{
		Type: &openapi3.Types{"string"},
	}
	result, err := c.applySchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result.Types, []string{"string"}) {
		t.Errorf("expected Types [string], got %+v", result.Types)
	}

	// Type present, nullable true
	schema2 := &openapi3.Schema{
		Type:     &openapi3.Types{"string"},
		Nullable: true,
	}
	result2, err := c.applySchema(schema2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result2.Types, []string{"string", "null"}) {
		t.Errorf("expected Types [string null], got %+v", result2.Types)
	}

	// Type already includes null
	schema3 := &openapi3.Schema{
		Type:     &openapi3.Types{"string", "null"},
		Nullable: true,
	}
	result3, err := c.applySchema(schema3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result3.Types, []string{"string", "null"}) {
		t.Errorf("expected Types [string null], got %+v", result3.Types)
	}

	// No type, nullable true
	schema4 := &openapi3.Schema{
		Nullable: true,
	}
	result4, err := c.applySchema(schema4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result4.Types, []string{"null"}) && !reflect.DeepEqual(result4.Types, []string{"string", "null"}) {
		t.Errorf("expected Types to include null, got %+v", result4.Types)
	}
}

func TestApplySchema_StringNumberArrayObject(t *testing.T) {
	c := &Converter{}
	// String
	maxLen := uint64(10)
	schema := &openapi3.Schema{
		Type:      &openapi3.Types{"string"},
		MinLength: 2,
		MaxLength: &maxLen,
		Pattern:   "abc",
	}
	result, err := c.applySchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String == nil || result.String.MinLength != 2 || result.String.MaxLength == nil || *result.String.MaxLength != 10 || result.String.Pattern != "abc" {
		t.Errorf("string validation not set correctly: %+v", result.String)
	}

	// Number
	min := 1.0
	max := 5.0
	mult := 2.0
	schema2 := &openapi3.Schema{
		Type:         &openapi3.Types{"number"},
		Min:          &min,
		Max:          &max,
		MultipleOf:   &mult,
		ExclusiveMin: true,
		ExclusiveMax: true,
	}
	result2, err := c.applySchema(schema2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result2.Number == nil || result2.Number.Minimum == nil || *result2.Number.Minimum != 1.0 ||
		result2.Number.Maximum == nil || *result2.Number.Maximum != 5.0 ||
		result2.Number.MultipleOf == nil || *result2.Number.MultipleOf != 2.0 ||
		!result2.Number.ExclusiveMinimum || !result2.Number.ExclusiveMaximum {
		t.Errorf("number validation not set correctly: %+v", result2.Number)
	}

	// Array
	maxItems := uint64(3)
	itemSchema := &openapi3.Schema{Title: "Item"}
	schema3 := &openapi3.Schema{
		Type:        &openapi3.Types{"array"},
		MinItems:    1,
		MaxItems:    &maxItems,
		UniqueItems: true,
		Items:       &openapi3.SchemaRef{Value: itemSchema},
	}
	result3, err := c.applySchema(schema3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result3.Array == nil || result3.Array.MinItems != 1 || result3.Array.MaxItems == nil || *result3.Array.MaxItems != 3 || !result3.Array.UniqueItems {
		t.Errorf("array validation not set correctly: %+v", result3.Array)
	}
	if result3.Array.Items == nil || result3.Array.Items.Title != "Item" {
		t.Errorf("array item schema not set correctly: %+v", result3.Array.Items)
	}

	// Object
	propSchema := &openapi3.Schema{Title: "Prop"}
	schema4 := &openapi3.Schema{
		Type: &openapi3.Types{"object"},
		Properties: map[string]*openapi3.SchemaRef{
			"foo": {Value: propSchema},
		},
	}
	result4, err := c.applySchema(schema4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result4.Object == nil || result4.Object.Properties == nil || result4.Object.Properties["foo"] == nil || result4.Object.Properties["foo"].Title != "Prop" {
		t.Errorf("object validation not set correctly: %+v", result4.Object)
	}
}

func TestApplySchema_OneOf_AnyOf_AllOf_Not(t *testing.T) {
	c := &Converter{}
	// OneOf
	schema := &openapi3.Schema{
		OneOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Title: "OneA"}},
			{Value: &openapi3.Schema{Title: "OneB"}},
		},
	}
	result, err := c.applySchema(schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.OneOf) != 2 || result.OneOf[0].Title != "OneA" || result.OneOf[1].Title != "OneB" {
		t.Errorf("OneOf not set correctly: %+v", result.OneOf)
	}

	// AnyOf
	schema2 := &openapi3.Schema{
		AnyOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Title: "AnyA"}},
			{Value: &openapi3.Schema{Title: "AnyB"}},
		},
	}
	result2, err := c.applySchema(schema2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result2.AnyOf) != 2 || result2.AnyOf[0].Title != "AnyA" || result2.AnyOf[1].Title != "AnyB" {
		t.Errorf("AnyOf not set correctly: %+v", result2.AnyOf)
	}

	// AllOf
	schema3 := &openapi3.Schema{
		AllOf: []*openapi3.SchemaRef{
			{Value: &openapi3.Schema{Title: "AllA"}},
			{Value: &openapi3.Schema{Title: "AllB"}},
		},
	}
	result3, err := c.applySchema(schema3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result3.AllOf) != 2 || result3.AllOf[0].Title != "AllA" || result3.AllOf[1].Title != "AllB" {
		t.Errorf("AllOf not set correctly: %+v", result3.AllOf)
	}

	// Not
	schema4 := &openapi3.Schema{
		Not: &openapi3.SchemaRef{Value: &openapi3.Schema{Title: "NotA"}},
	}
	result4, err := c.applySchema(schema4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result4.Not == nil || result4.Not.Title != "NotA" {
		t.Errorf("Not not set correctly: %+v", result4.Not)
	}
}
