package converter

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// applySchemaMetadata applies basic schema metadata to create a Schema
func (c *Converter) applySchema(schema *openapi3.Schema) (*Schema, error) {
	if schema == nil {
		return nil, fmt.Errorf("cannot apply metadata to nil schema")
	}

	// Create a new Schema
	result := &Schema{
		Title:       schema.Title,
		Description: schema.Description,
		Format:      schema.Format,
		Enum:        schema.Enum,
		Default:     schema.Default,
		Example:     schema.Example,
		ReadOnly:    schema.ReadOnly,
		WriteOnly:   schema.WriteOnly,
	}

	// Handle types, including nullable
	if schema.Type != nil {
		result.Types = *schema.Type
	} else {
		result.Types = []string{}
	}

	if schema.Nullable {
		isNullableAlreadyPresent := false
		if result.Types != nil {
			for _, t := range result.Types {
				if t == "null" {
					isNullableAlreadyPresent = true
					break
				}
			}
		} else {
			// If no types are set, default to string
			result.Types = append(result.Types, "string")
		}
		if !isNullableAlreadyPresent {
			result.Types = append(result.Types, "null")
		}
	}

	var err error

	if hasStringType(schema) {
		result.String = c.createStringValidation(schema)
	}

	if hasNumericType(schema) {
		result.Number = c.createNumberValidation(schema)
	}

	if hasArrayType(schema) {
		result.Array, err = c.createArrayValidation(schema)
		if err != nil {
			return nil, fmt.Errorf("error creating array validation: %w", err)
		}
	}

	if hasObjectType(schema) {
		result.Object, err = c.createObjectValidation(schema)
		if err != nil {
			return nil, fmt.Errorf("error creating object validation: %w", err)
		}
	}

	// Handle OneOf
	if len(schema.OneOf) > 0 {
		result.OneOf = make([]*Schema, len(schema.OneOf))
		for i, subSchemaRef := range schema.OneOf {
			if subSchemaRef == nil || subSchemaRef.Value == nil {
				return nil, fmt.Errorf("oneOf contains a nil schema reference or value at index %d", i)
			}
			subSchema, err := c.applySchema(subSchemaRef.Value) // Recursive call
			if err != nil {
				return nil, fmt.Errorf("error processing oneOf sub-schema at index %d: %w", i, err)
			}
			if subSchema != nil {
				result.OneOf[i] = subSchema
			} else {
				// Handle case where recursive call returned nil schema but no error?
				// For now, return error as a nil schema here is likely unexpected.
				return nil, fmt.Errorf("oneOf sub-schema at index %d resulted in a nil schema", i)
			}
		}
	}

	// Handle AnyOf
	if len(schema.AnyOf) > 0 {
		result.AnyOf = make([]*Schema, len(schema.AnyOf))
		for i, subSchemaRef := range schema.AnyOf {
			if subSchemaRef == nil || subSchemaRef.Value == nil {
				return nil, fmt.Errorf("anyOf contains a nil schema reference or value at index %d", i)
			}
			subSchema, err := c.applySchema(subSchemaRef.Value)
			if err != nil {
				return nil, fmt.Errorf("error processing anyOf sub-schema at index %d: %w", i, err)
			}
			if subSchema != nil {
				result.AnyOf[i] = subSchema
			} else {
				return nil, fmt.Errorf("anyOf sub-schema at index %d resulted in a nil schema", i)
			}
		}
	}

	// Handle AllOf
	if len(schema.AllOf) > 0 {
		result.AllOf = make([]*Schema, len(schema.AllOf))
		for i, subSchemaRef := range schema.AllOf {
			if subSchemaRef == nil || subSchemaRef.Value == nil {
				return nil, fmt.Errorf("allOf contains a nil schema reference or value at index %d", i)
			}
			subSchema, err := c.applySchema(subSchemaRef.Value)
			if err != nil {
				return nil, fmt.Errorf("error processing allOf sub-schema at index %d: %w", i, err)
			}
			if subSchema != nil {
				result.AllOf[i] = subSchema
			} else {
				return nil, fmt.Errorf("allOf sub-schema at index %d resulted in a nil schema", i)
			}
		}
	}

	// Handle Not
	if schema.Not != nil && schema.Not.Value != nil {
		notSchema, err := c.applySchema(schema.Not.Value)
		if err != nil {
			return nil, fmt.Errorf("error processing not sub-schema: %w", err)
		}
		if notSchema != nil {
			result.Not = notSchema
		} else {
			return nil, fmt.Errorf("not sub-schema resulted in a nil schema")
		}
	}

	return result, nil
}

// createStringValidation creates string-specific validations
func (c *Converter) createStringValidation(schema *openapi3.Schema) *StringValidation {
	if schema == nil {
		return nil
	}
	return &StringValidation{
		MinLength: schema.MinLength,
		MaxLength: schema.MaxLength,
		Pattern:   schema.Pattern,
	}
}

// createNumberValidation creates number-specific validations
func (c *Converter) createNumberValidation(schema *openapi3.Schema) *NumberValidation {
	if schema == nil {
		return nil
	}
	return &NumberValidation{
		Minimum:          schema.Min,
		Maximum:          schema.Max,
		MultipleOf:       schema.MultipleOf,
		ExclusiveMinimum: schema.ExclusiveMin,
		ExclusiveMaximum: schema.ExclusiveMax,
	}
}

// createArrayValidation creates array-specific validations
func (c *Converter) createArrayValidation(schema *openapi3.Schema) (*ArrayValidation, error) {
	if schema == nil {
		return nil, nil
	}
	result := &ArrayValidation{
		MinItems:    schema.MinItems,
		MaxItems:    schema.MaxItems,
		UniqueItems: schema.UniqueItems,
	}

	if schema.Items != nil && schema.Items.Value != nil {
		itemsSchema, err := c.applySchema(schema.Items.Value)
		if err != nil {
			return nil, fmt.Errorf("error processing array items schema: %w", err)
		}
		result.Items = itemsSchema
	}

	return result, nil
}

// createObjectValidation creates object-specific validations
func (c *Converter) createObjectValidation(schema *openapi3.Schema) (*ObjectValidation, error) {
	if schema == nil {
		return nil, nil
	}

	result := &ObjectValidation{
		Required:      schema.Required,
		MinProperties: schema.MinProps,
		MaxProperties: schema.MaxProps,
	}

	if len(schema.Properties) > 0 {
		result.Properties = make(map[string]*Schema)
		for propName, propSchemaRef := range schema.Properties {
			if propSchemaRef != nil {
				if propSchemaRef.Value != nil {
					propSchema, err := c.applySchema(propSchemaRef.Value)
					if err != nil {
						return nil, fmt.Errorf("error processing property '%s': %w", propName, err)
					}
					if propSchema != nil {
						result.Properties[propName] = propSchema
					}
				} else {
					result.Properties[propName] = &Schema{}
					fmt.Printf("Warning: Property '%s' has a non-nil SchemaRef but nil Value. Mapping to empty schema.\n", propName)
				}
			} else {
				fmt.Printf("Warning: Property '%s' has a nil schema reference in the OpenAPI spec.\n", propName)
			}
		}
		if len(result.Properties) == 0 {
			result.Properties = nil
		}
	}

	// --- Handle additionalProperties ---
	if schema.AdditionalProperties.Has != nil {
		// Case 1: additionalProperties is explicitly true or false
		if !*schema.AdditionalProperties.Has {
			result.DisallowAdditionalProperties = true
		} else {
			result.AdditionalProperties = &Schema{} // Represents allowing any additional properties ({})
		}
	} else if schema.AdditionalProperties.Schema != nil {
		// Case 2: additionalProperties is a schema object (or meant to be)
		if schema.AdditionalProperties.Schema.Value != nil {
			addPropSchema, err := c.applySchema(schema.AdditionalProperties.Schema.Value)
			if err != nil {
				return nil, fmt.Errorf("error processing additionalProperties schema: %w", err)
			}
			if addPropSchema != nil {
				result.AdditionalProperties = addPropSchema
			} else {
				// SchemaRef.Value was non-nil, but applySchema returned nil. Map to {}.
				result.AdditionalProperties = &Schema{}
			}
		} else {
			result.AdditionalProperties = &Schema{}
			fmt.Printf("Warning: AdditionalProperties SchemaRef has a nil Value. Mapping to empty schema {}.\n")
		}
	}

	return result, nil
}
