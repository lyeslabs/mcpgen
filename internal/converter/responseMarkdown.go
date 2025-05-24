package converter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// buildResponseMarkdown builds the Markdown documentation for a response.
func (c *Converter) buildResponseMarkdown(
	code, contentType string,
	responseRef *openapi3.ResponseRef,
	schema *openapi3.Schema,
) string {
	var b strings.Builder
	b.WriteString("# API Response Information\n\n")
	b.WriteString("Below is the response template for this API endpoint.\n\n")
	b.WriteString("The template shows a possible response, including its status code and content type, to help you understand and generate correct outputs.\n\n")
	b.WriteString(fmt.Sprintf("**Status Code:** %s\n\n", code))
	b.WriteString(fmt.Sprintf("**Content-Type:** %s\n\n", contentType))
	if desc := getResponseDescription(responseRef); desc != "" {
		b.WriteString(fmt.Sprintf("> %s\n\n", desc))
	}
	b.WriteString("## Response Structure\n\n")
	c.writeSchemaMarkdown(&b, schema, 0, "")
	return b.String()
}

// writeSchemaMarkdown documents a schema in Markdown, recursively.
func (c *Converter) writeSchemaMarkdown(
	b *strings.Builder,
	schema *openapi3.Schema,
	indent int,
	fieldName string,
) {
	if schema == nil {
		return
	}
	ind := strings.Repeat("  ", indent)
	typeDesc := schemaTypeDescription(schema)
	description := schema.Description

	// Print the field or root schema line
	if fieldName != "" {
		if description == "" {
			b.WriteString(fmt.Sprintf("%s- **%s** (Type: %s):\n", ind, fieldName, typeDesc))
		} else {
			b.WriteString(fmt.Sprintf("%s- **%s**: %s (Type: %s):\n", ind, fieldName, description, typeDesc))
		}
	} else {
		if description == "" {
			b.WriteString(fmt.Sprintf("%s- Structure (Type: %s):\n", ind, typeDesc))
		} else {
			b.WriteString(fmt.Sprintf("%s- %s (Type: %s):\n", ind, description, typeDesc))
		}
	}

	c.writeSchemaDetails(b, schema, indent+1)
	c.writeSchemaProperties(b, schema, indent)
	c.writeSchemaCombinators(b, schema, indent)
	c.writeAdditionalProperties(b, schema, indent)
}

// writeSchemaProperties documents object properties and array items.
func (c *Converter) writeSchemaProperties(
	b *strings.Builder,
	schema *openapi3.Schema,
	indent int,
) {
	// Object properties
	if isObject(schema) && len(schema.Properties) > 0 {
		for propName, propRef := range schema.Properties {
			if propRef != nil && propRef.Value != nil {
				c.writeSchemaMarkdown(b, propRef.Value, indent+1, propName)
			}
		}
	}
	// Array items
	if isArray(schema) && schema.Items != nil && schema.Items.Value != nil {
		c.writeSchemaMarkdown(b, schema.Items.Value, indent+1, "Items")
	}
}

// writeSchemaCombinators documents oneOf, anyOf, allOf, and not combinators.
func (c *Converter) writeSchemaCombinators(
	b *strings.Builder,
	schema *openapi3.Schema,
	indent int,
) {
	ind := strings.Repeat("  ", indent)
	if len(schema.OneOf) > 0 {
		b.WriteString(fmt.Sprintf("%s  - **One Of the following structures**:\n", ind))
		for i, sub := range schema.OneOf {
			c.writeSchemaMarkdown(b, sub.Value, indent+2, fmt.Sprintf("Option %d", i+1))
		}
	}
	if len(schema.AnyOf) > 0 {
		b.WriteString(fmt.Sprintf("%s  - **Any Of the following structures**:\n", ind))
		for i, sub := range schema.AnyOf {
			c.writeSchemaMarkdown(b, sub.Value, indent+2, fmt.Sprintf("Option %d", i+1))
		}
	}
	if len(schema.AllOf) > 0 {
		b.WriteString(fmt.Sprintf("%s  - **Combines All Of the following structures**:\n", ind))
		for i, sub := range schema.AllOf {
			c.writeSchemaMarkdown(b, sub.Value, indent+2, fmt.Sprintf("Part %d", i+1))
		}
	}
	if schema.Not != nil && schema.Not.Value != nil {
		b.WriteString(fmt.Sprintf("%s  - **Not**: Cannot be the following structure:\n", ind))
		c.writeSchemaMarkdown(b, schema.Not.Value, indent+2, "Forbidden Structure")
	}
}

// writeAdditionalProperties documents additionalProperties for objects.
func (c *Converter) writeAdditionalProperties(
	b *strings.Builder,
	schema *openapi3.Schema,
	indent int,
) {
	ind := strings.Repeat("  ", indent)
	if isObject(schema) && schema.AdditionalProperties.Schema != nil && schema.AdditionalProperties.Schema.Value != nil {
		b.WriteString(fmt.Sprintf("%s  - **Additional Properties**:\n", ind))
		c.writeSchemaMarkdown(b, schema.AdditionalProperties.Schema.Value, indent+2, "property value")
	} else if isObject(schema) && schema.AdditionalProperties.Has != nil && *schema.AdditionalProperties.Has {
		b.WriteString(fmt.Sprintf("%s  - **Allows Additional Properties**\n", ind))
	}
}

// writeSchemaDetails adds validation rules, examples, and default values in Markdown.
func (c *Converter) writeSchemaDetails(
	b *strings.Builder,
	schema *openapi3.Schema,
	indent int,
) {
	ind := strings.Repeat("  ", indent)
	var details []string

	// String validations
	if schema.MinLength > 0 {
		details = append(details, fmt.Sprintf("Min Length: %d", schema.MinLength))
	}
	if schema.MaxLength != nil && *schema.MaxLength > 0 {
		details = append(details, fmt.Sprintf("Max Length: %d", *schema.MaxLength))
	}
	if schema.Pattern != "" {
		details = append(details, fmt.Sprintf("Pattern: '%s'", strings.ReplaceAll(schema.Pattern, "`", "'")))
	}

	// Numeric validations
	if schema.Min != nil {
		details = append(details, fmt.Sprintf("Minimum: %v", *schema.Min))
	}
	if schema.Max != nil {
		details = append(details, fmt.Sprintf("Maximum: %v", *schema.Max))
	}
	if schema.ExclusiveMin {
		details = append(details, "Exclusive Minimum: true")
	}
	if schema.ExclusiveMax {
		details = append(details, "Exclusive Maximum: true")
	}
	if schema.MultipleOf != nil {
		details = append(details, fmt.Sprintf("Multiple Of: %v", *schema.MultipleOf))
	}

	// Array validations
	if schema.MinItems > 0 {
		details = append(details, fmt.Sprintf("Min Items: %d", schema.MinItems))
	}
	if schema.MaxItems != nil && *schema.MaxItems > 0 {
		details = append(details, fmt.Sprintf("Max Items: %d", *schema.MaxItems))
	}
	if schema.UniqueItems {
		details = append(details, "Unique Items: true")
	}

	// Nullable
	if schema.Nullable {
		details = append(details, "Nullable: true")
	}

	// Default/Example handling
	if schema.Default != nil {
		details = append(details, fmt.Sprintf("Default: '%s'", formatForGoRawString(schema, schema.Default)))
	}
	if schema.Example != nil {
		details = append(details, fmt.Sprintf("Example: '%s'", formatForGoRawString(schema, schema.Example)))
	}

	// Enums
	if len(schema.Enum) > 0 {
		var enumStrings []string
		for _, e := range schema.Enum {
			enumStrings = append(enumStrings, fmt.Sprintf("'%s'", formatForGoRawString(schema, e)))
		}
		details = append(details, fmt.Sprintf("Enum: [%s]", strings.Join(enumStrings, ", ")))
	}

	// Print details
	if len(details) > 0 {
		detailIndent := ind + "  "
		for _, detail := range details {
			b.WriteString(fmt.Sprintf("%s- %s\n", detailIndent, detail))
		}
	}
}
