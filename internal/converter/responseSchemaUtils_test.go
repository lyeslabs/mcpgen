package converter


import (
    "testing"

    "github.com/getkin/kin-openapi/openapi3"
)

func TestHasSchema(t *testing.T) {
    s := &openapi3.Schema{}
    ref := &openapi3.SchemaRef{Value: s}
    mt := &openapi3.MediaType{Schema: ref}

    if !hasSchema(mt) {
        t.Error("hasSchema should return true for valid schema")
    }
    if hasSchema(&openapi3.MediaType{Schema: nil}) {
        t.Error("hasSchema should return false if Schema is nil")
    }
    if hasSchema(&openapi3.MediaType{Schema: &openapi3.SchemaRef{Value: nil}}) {
        t.Error("hasSchema should return false if Schema.Value is nil")
    }
    if hasSchema(nil) {
        t.Error("hasSchema should return false if MediaType is nil")
    }
}

func TestSchemaTypeDescription(t *testing.T) {
    strType := openapi3.Types{"string"}
    arrType := openapi3.Types{"array"}
    objType := openapi3.Types{"object"}
    numType := openapi3.Types{"number"}
    intType := openapi3.Types{"integer"}

    cases := []struct {
        name   string
        schema *openapi3.Schema
        want   string
    }{
        {"string", &openapi3.Schema{Type: &strType}, "string"},
        {"string+format", &openapi3.Schema{Type: &strType, Format: "date"}, "string, date"},
        {"nullable", &openapi3.Schema{Type: &strType, Nullable: true}, "string, nullable"},
        {"array", &openapi3.Schema{Type: &arrType}, "array"},
        {"object", &openapi3.Schema{Type: &objType}, "object"},
        {"number", &openapi3.Schema{Type: &numType}, "number"},
        {"integer", &openapi3.Schema{Type: &intType}, "integer"},
        {"combinator", &openapi3.Schema{OneOf: []*openapi3.SchemaRef{{}}}, "Combinator"},
        {"unknown", &openapi3.Schema{}, "unknown"},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := schemaTypeDescription(c.schema)
            if got != c.want {
                t.Errorf("schemaTypeDescription() = %q, want %q", got, c.want)
            }
        })
    }
}

func TestIsArray(t *testing.T) {
    arrType := openapi3.Types{"array"}
    strType := openapi3.Types{"string"}
    if !isArray(&openapi3.Schema{Type: &arrType}) {
        t.Error("isArray should return true for array type")
    }
    if isArray(&openapi3.Schema{Type: &strType}) {
        t.Error("isArray should return false for non-array type")
    }
    if isArray(nil) {
        t.Error("isArray should return false for nil schema")
    }
}

func TestIsObject(t *testing.T) {
    objType := openapi3.Types{"object"}
    strType := openapi3.Types{"string"}
    if !isObject(&openapi3.Schema{Type: &objType}) {
        t.Error("isObject should return true for object type")
    }
    if isObject(&openapi3.Schema{Type: &strType}) {
        t.Error("isObject should return false for non-object type")
    }
    if isObject(nil) {
        t.Error("isObject should return false for nil schema")
    }
}

func TestHasStringType(t *testing.T) {
    strType := openapi3.Types{"string"}
    arrType := openapi3.Types{"array"}
    if !hasStringType(&openapi3.Schema{Type: &strType}) {
        t.Error("hasStringType should return true for string type")
    }
    if hasStringType(&openapi3.Schema{Type: &arrType}) {
        t.Error("hasStringType should return false for non-string type")
    }
    if hasStringType(nil) {
        t.Error("hasStringType should return false for nil schema")
    }
}

func TestHasNumericType(t *testing.T) {
    numType := openapi3.Types{"number"}
    intType := openapi3.Types{"integer"}
    strType := openapi3.Types{"string"}
    if !hasNumericType(&openapi3.Schema{Type: &numType}) {
        t.Error("hasNumericType should return true for number type")
    }
    if !hasNumericType(&openapi3.Schema{Type: &intType}) {
        t.Error("hasNumericType should return true for integer type")
    }
    if hasNumericType(&openapi3.Schema{Type: &strType}) {
        t.Error("hasNumericType should return false for non-numeric type")
    }
    if hasNumericType(nil) {
        t.Error("hasNumericType should return false for nil schema")
    }
}

func TestHasArrayType(t *testing.T) {
    arrType := openapi3.Types{"array"}
    strType := openapi3.Types{"string"}
    if !hasArrayType(&openapi3.Schema{Type: &arrType}) {
        t.Error("hasArrayType should return true for array type")
    }
    if hasArrayType(&openapi3.Schema{Type: &strType}) {
        t.Error("hasArrayType should return false for non-array type")
    }
    if hasArrayType(nil) {
        t.Error("hasArrayType should return false for nil schema")
    }
}

func TestHasObjectType(t *testing.T) {
    objType := openapi3.Types{"object"}
    strType := openapi3.Types{"string"}
    if !hasObjectType(&openapi3.Schema{Type: &objType}) {
        t.Error("hasObjectType should return true for object type")
    }
    if hasObjectType(&openapi3.Schema{Type: &strType}) {
        t.Error("hasObjectType should return false for non-object type")
    }
    if hasObjectType(nil) {
        t.Error("hasObjectType should return false for nil schema")
    }
}
