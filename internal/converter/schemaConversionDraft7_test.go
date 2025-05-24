package converter


import (
    "reflect"
    "testing"
)

func TestSchemaToDraft7Map_Basic(t *testing.T) {
    min := 1.0
    max := 5.0
    mult := 2.0
    var maxLen uint64 = 10
    var minLen uint64 = 2
    var maxItems uint64 = 3
    var minItems uint64 = 1
    var maxProps uint64 = 2
    var minProps uint64 = 1

    schema := &Schema{
        Title:       "Test Title",
        Description: "Test Description",
        Format:      "date-time",
        Default:     "default",
        Example:     "example",
        Enum:        []interface{}{"A", "B"},
        ReadOnly:    true,
        WriteOnly:   false,
        Types:       []string{"string"},
        String: &StringValidation{
            MinLength: minLen,
            MaxLength: &maxLen,
            Pattern:   "^[a-z]+$",
        },
        Number: &NumberValidation{
            Minimum:          &min,
            ExclusiveMinimum: false,
            Maximum:          &max,
            ExclusiveMaximum: true,
            MultipleOf:       &mult,
        },
        Array: &ArrayValidation{
            Items: &Schema{
                Types: []string{"integer"},
            },
            MinItems:    minItems,
            MaxItems:    &maxItems,
            UniqueItems: true,
        },
        Object: &ObjectValidation{
            Properties: map[string]*Schema{
                "foo": {Types: []string{"string"}},
            },
            Required:                     []string{"foo"},
            MinProperties:                minProps,
            MaxProperties:                &maxProps,
            DisallowAdditionalProperties: true,
        },
    }

    got, err := schemaToDraft7Map(schema)
    if err != nil {
        t.Fatalf("schemaToDraft7Map() error = %v", err)
    }

    want := map[string]interface{}{
        "title":            "Test Title",
        "description":      "Test Description",
        "format":           "date-time",
        "default":          "default",
        "example":          "example",
        "enum":             []interface{}{"A", "B"},
        "readOnly":         true,
        "type":             "string",
        "minLength":        minLen,
        "maxLength":        maxLen,
        "pattern":          "^[a-z]+$",
        "minimum":          min,
        "exclusiveMaximum": max,
        "multipleOf":       mult,
        "items": map[string]interface{}{
            "type": "integer",
        },
        "minItems":    minItems,
        "maxItems":    maxItems,
        "uniqueItems": true,
        "properties": map[string]interface{}{
            "foo": map[string]interface{}{
                "type": "string",
            },
        },
        "required":             []string{"foo"},
        "minProperties":        minProps,
        "maxProperties":        maxProps,
        "additionalProperties": false,
    }

    // Only check for keys we expect (since the function may add more)
    for k, v := range want {
        if !reflect.DeepEqual(got[k], v) {
            t.Errorf("schemaToDraft7Map()[%q] = %v, want %v", k, got[k], v)
        }
    }
}


func TestAddBasicMetadata(t *testing.T) {
    s := &Schema{
        Title:       "My Title",
        Description: "My Description",
        Format:      "date",
        Default:     42,
        Example:     "foo",
        Enum:        []interface{}{"A", "B"},
        ReadOnly:    true,
        WriteOnly:   true,
    }
    result := make(map[string]interface{})
    addBasicMetadata(result, s)

    want := map[string]interface{}{
        "title":       "My Title",
        "description": "My Description",
        "format":      "date",
        "default":     42,
        "example":     "foo",
        "enum":        []interface{}{"A", "B"},
        "readOnly":    true,
        "writeOnly":   true,
    }
    if !reflect.DeepEqual(result, want) {
        t.Errorf("addBasicMetadata() = %v, want %v", result, want)
    }
}


func TestAddCombinators(t *testing.T) {
    s := &Schema{
        OneOf: []*Schema{
            {Title: "A"},
            {Title: "B"},
        },
        AnyOf: []*Schema{
            {Title: "C"},
        },
        AllOf: []*Schema{
            {Title: "D"},
        },
        Not: &Schema{Title: "E"},
    }
    result := make(map[string]interface{})
    err := addCombinators(result, s)
    if err != nil {
        t.Fatalf("addCombinators() error = %v", err)
    }
    // Only check keys exist and are correct type
    if _, ok := result["oneOf"]; !ok {
        t.Error("addCombinators() missing oneOf")
    }
    if _, ok := result["anyOf"]; !ok {
        t.Error("addCombinators() missing anyOf")
    }
    if _, ok := result["allOf"]; !ok {
        t.Error("addCombinators() missing allOf")
    }
    if _, ok := result["not"]; !ok {
        t.Error("addCombinators() missing not")
    }
}


func TestConvertSubSchemas(t *testing.T) {
    subs := []*Schema{
        {Title: "A"},
        {Title: "B"},
    }
    got, err := convertSubSchemas(subs)
    if err != nil {
        t.Fatalf("convertSubSchemas() error = %v", err)
    }
    if len(got) != 2 {
        t.Errorf("convertSubSchemas() length = %d, want 2", len(got))
    }
    if got[0]["title"] != "A" || got[1]["title"] != "B" {
        t.Errorf("convertSubSchemas() = %v", got)
    }
}

func TestAddType(t *testing.T) {
    s := &Schema{Types: []string{"string"}}
    result := make(map[string]interface{})
    addType(result, s)
    if result["type"] != "string" {
        t.Errorf("addType() = %v, want 'string'", result["type"])
    }

    s2 := &Schema{Types: []string{"string", "null"}}
    result2 := make(map[string]interface{})
    addType(result2, s2)
    if !reflect.DeepEqual(result2["type"], []string{"string", "null"}) {
        t.Errorf("addType() = %v, want [string null]", result2["type"])
    }
}


func TestAddStringValidation(t *testing.T) {
    var maxLen uint64 = 10
    s := &Schema{
        String: &StringValidation{
            MinLength: 2,
            MaxLength: &maxLen,
            Pattern:   "abc",
        },
    }
    result := make(map[string]interface{})
    addStringValidation(result, s)
    want := map[string]interface{}{
        "minLength": uint64(2),
        "maxLength": maxLen,
        "pattern":   "abc",
    }
    if !reflect.DeepEqual(result, want) {
        t.Errorf("addStringValidation() = %v, want %v", result, want)
    }
}


func TestAddNumberValidation(t *testing.T) {
    min := 1.0
    max := 5.0
    mult := 2.0
    s := &Schema{
        Number: &NumberValidation{
            Minimum:          &min,
            ExclusiveMinimum: false,
            Maximum:          &max,
            ExclusiveMaximum: true,
            MultipleOf:       &mult,
        },
    }
    result := make(map[string]interface{})
    addNumberValidation(result, s)
    want := map[string]interface{}{
        "minimum":          min,
        "exclusiveMaximum": max,
        "multipleOf":       mult,
    }
    if !reflect.DeepEqual(result, want) {
        t.Errorf("addNumberValidation() = %v, want %v", result, want)
    }
}


func TestAddArrayValidation(t *testing.T) {
    var maxItems uint64 = 3
    s := &Schema{
        Array: &ArrayValidation{
            Items: &Schema{Types: []string{"integer"}},
            MinItems: 1,
            MaxItems: &maxItems,
            UniqueItems: true,
        },
    }
    result := make(map[string]interface{})
    err := addArrayValidation(result, s)
    if err != nil {
        t.Fatalf("addArrayValidation() error = %v", err)
    }
    if result["minItems"] != uint64(1) {
        t.Errorf("addArrayValidation() minItems = %v, want 1", result["minItems"])
    }
    if result["maxItems"] != maxItems {
        t.Errorf("addArrayValidation() maxItems = %v, want %v", result["maxItems"], maxItems)
    }
    if result["uniqueItems"] != true {
        t.Errorf("addArrayValidation() uniqueItems = %v, want true", result["uniqueItems"])
    }
    items, ok := result["items"].(map[string]interface{})
    if !ok || items["type"] != "integer" {
        t.Errorf("addArrayValidation() items = %v, want type=integer", result["items"])
    }
}

func TestAddObjectValidation(t *testing.T) {
    var maxProps uint64 = 2
    s := &Schema{
        Object: &ObjectValidation{
            Properties: map[string]*Schema{
                "foo": {Types: []string{"string"}},
            },
            Required: []string{"foo"},
            MinProperties: 1,
            MaxProperties: &maxProps,
            DisallowAdditionalProperties: true,
        },
    }
    result := make(map[string]interface{})
    err := addObjectValidation(result, s)
    if err != nil {
        t.Fatalf("addObjectValidation() error = %v", err)
    }
    if result["minProperties"] != uint64(1) {
        t.Errorf("addObjectValidation() minProperties = %v, want 1", result["minProperties"])
    }
    if result["maxProperties"] != maxProps {
        t.Errorf("addObjectValidation() maxProperties = %v, want %v", result["maxProperties"], maxProps)
    }
    if result["additionalProperties"] != false {
        t.Errorf("addObjectValidation() additionalProperties = %v, want false", result["additionalProperties"])
    }
    props, ok := result["properties"].(map[string]interface{})
    if !ok || props["foo"] == nil {
        t.Errorf("addObjectValidation() properties = %v, want foo", result["properties"])
    }
}
