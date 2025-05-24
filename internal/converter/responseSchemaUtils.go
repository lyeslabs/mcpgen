package converter

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// hasSchema checks if a media type has a valid schema.
func hasSchema(mediaType *openapi3.MediaType) bool {
	return mediaType != nil && mediaType.Schema != nil && mediaType.Schema.Value != nil
}

// schemaTypeDescription returns a string describing the schema's type and format.
func schemaTypeDescription(schema *openapi3.Schema) string {
	var typeStrs []string
	if schema.Type != nil {
		typeStrs = *schema.Type
	}
	if schema.Format != "" {
		typeStrs = append(typeStrs, schema.Format)
	}
	if schema.Nullable {
		typeStrs = append(typeStrs, "nullable")
	}
	if len(typeStrs) == 0 {
		if len(schema.OneOf) > 0 || len(schema.AnyOf) > 0 || len(schema.AllOf) > 0 || schema.Not != nil {
			return "Combinator"
		}
		return "unknown"
	}
	return strings.Join(typeStrs, ", ")
}

// isArray checks if a schema represents an array.
func isArray(schema *openapi3.Schema) bool {
	return schema != nil && schema.Type != nil &&
		len(*schema.Type) > 0 && (*schema.Type)[0] == "array"
}

// isObject checks if a schema represents an object.
func isObject(schema *openapi3.Schema) bool {
	return schema != nil && schema.Type != nil &&
		len(*schema.Type) > 0 && (*schema.Type)[0] == "object"
}



// Type checking helpers - now with nil checking
func hasStringType(schema *openapi3.Schema) bool {
	return schema != nil && schema.Type != nil && contains(*schema.Type, "string")
}

func hasNumericType(schema *openapi3.Schema) bool {
	return schema != nil && schema.Type != nil &&
		(contains(*schema.Type, "number") || contains(*schema.Type, "integer"))
}

func hasArrayType(schema *openapi3.Schema) bool {
	return schema != nil && schema.Type != nil && contains(*schema.Type, "array")
}

func hasObjectType(schema *openapi3.Schema) bool {
	return schema != nil && schema.Type != nil && contains(*schema.Type, "object")
}
