package converter

import "fmt"

func schemaToDraft7Map(s *Schema) (map[string]interface{}, error) {
	if s == nil {
		return nil, nil
	}

	result := make(map[string]interface{})

	addBasicMetadata(result, s)
	addType(result, s)
	addStringValidation(result, s)
	addNumberValidation(result, s)

	if err := addCombinators(result, s); err != nil {
		return nil, err
	}
	if err := addArrayValidation(result, s); err != nil {
		return nil, err
	}
	if err := addObjectValidation(result, s); err != nil {
		return nil, err
	}

	return result, nil
}

func addBasicMetadata(result map[string]interface{}, s *Schema) {
	if s.Title != "" {
		result["title"] = s.Title
	}
	if s.Description != "" {
		result["description"] = s.Description
	}
	if s.Format != "" {
		result["format"] = s.Format
	}
	if s.Default != nil {
		result["default"] = s.Default
	}
	if s.Example != nil {
		result["example"] = s.Example
	}
	if len(s.Enum) > 0 {
		result["enum"] = s.Enum
	}
	if s.ReadOnly {
		result["readOnly"] = true
	}
	if s.WriteOnly {
		result["writeOnly"] = true
	}
}

func addCombinators(result map[string]interface{}, s *Schema) error {
	if len(s.OneOf) > 0 {
		oneOfSchemas, err := convertSubSchemas(s.OneOf)
		if err != nil {
			return fmt.Errorf("failed to convert oneOf: %w", err)
		}
		result["oneOf"] = oneOfSchemas
	}
	if len(s.AnyOf) > 0 {
		anyOfSchemas, err := convertSubSchemas(s.AnyOf)
		if err != nil {
			return fmt.Errorf("failed to convert anyOf: %w", err)
		}
		result["anyOf"] = anyOfSchemas
	}
	if len(s.AllOf) > 0 {
		allOfSchemas, err := convertSubSchemas(s.AllOf)
		if err != nil {
			return fmt.Errorf("failed to convert allOf: %w", err)
		}
		result["allOf"] = allOfSchemas
	}
	if s.Not != nil {
		notSchemaMap, err := schemaToDraft7Map(s.Not)
		if err != nil {
			return fmt.Errorf("failed to convert not sub-schema: %w", err)
		}
		if notSchemaMap == nil {
			return fmt.Errorf("not sub-schema resulted in a nil schema map")
		}
		result["not"] = notSchemaMap
	}
	return nil
}

func convertSubSchemas(subSchemas []*Schema) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(subSchemas))
	for i, subSchema := range subSchemas {
		subSchemaMap, err := schemaToDraft7Map(subSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to convert sub-schema at index %d: %w", i, err)
		}
		if subSchemaMap == nil {
			return nil, fmt.Errorf("sub-schema at index %d resulted in a nil schema map", i)
		}
		result[i] = subSchemaMap
	}
	return result, nil
}

func addType(result map[string]interface{}, s *Schema) {
	if len(s.Types) == 1 {
		result["type"] = s.Types[0]
	} else if len(s.Types) > 1 {
		result["type"] = s.Types
	}
}

func addStringValidation(result map[string]interface{}, s *Schema) {
	if s.String == nil {
		return
	}
	if s.String.MinLength > 0 {
		result["minLength"] = s.String.MinLength
	}
	if s.String.MaxLength != nil {
		result["maxLength"] = *s.String.MaxLength
	}
	if s.String.Pattern != "" {
		result["pattern"] = s.String.Pattern
	}
}

func addNumberValidation(result map[string]interface{}, s *Schema) {
	if s.Number == nil {
		return
	}
	if s.Number.Minimum != nil {
		if s.Number.ExclusiveMinimum {
			result["exclusiveMinimum"] = *s.Number.Minimum
		} else {
			result["minimum"] = *s.Number.Minimum
		}
	}
	if s.Number.Maximum != nil {
		if s.Number.ExclusiveMaximum {
			result["exclusiveMaximum"] = *s.Number.Maximum
		} else {
			result["maximum"] = *s.Number.Maximum
		}
	}
	if s.Number.MultipleOf != nil {
		result["multipleOf"] = *s.Number.MultipleOf
	}
}

func addArrayValidation(result map[string]interface{}, s *Schema) error {
	if s.Array == nil {
		return nil
	}
	if s.Array.Items != nil {
		itemsSchemaMap, err := schemaToDraft7Map(s.Array.Items)
		if err != nil {
			return fmt.Errorf("failed to convert array items schema: %w", err)
		}
		if itemsSchemaMap != nil {
			result["items"] = itemsSchemaMap
		}
	}
	if s.Array.MinItems > 0 {
		result["minItems"] = s.Array.MinItems
	}
	if s.Array.MaxItems != nil {
		result["maxItems"] = *s.Array.MaxItems
	}
	if s.Array.UniqueItems {
		result["uniqueItems"] = true
	}
	return nil
}

func addObjectValidation(result map[string]interface{}, s *Schema) error {
	if s.Object == nil {
		return nil
	}
	if len(s.Object.Properties) > 0 {
		propertiesMap := make(map[string]interface{})
		for propName, propSchema := range s.Object.Properties {
			propSchemaMap, err := schemaToDraft7Map(propSchema)
			if err != nil {
				return fmt.Errorf("failed to convert property '%s': %w", propName, err)
			}
			if propSchemaMap != nil {
				propertiesMap[propName] = propSchemaMap
			}
		}
		if len(propertiesMap) > 0 {
			result["properties"] = propertiesMap
		}
	}
	if len(s.Object.Required) > 0 {
		result["required"] = s.Object.Required
	}
	if s.Object.MinProperties > 0 {
		result["minProperties"] = s.Object.MinProperties
	}
	if s.Object.MaxProperties != nil {
		result["maxProperties"] = *s.Object.MaxProperties
	}

	// Handle additionalProperties mapping
	if s.Object.DisallowAdditionalProperties {
		result["additionalProperties"] = false
	} else if s.Object.AdditionalProperties != nil {
		addPropSchemaMap, err := schemaToDraft7Map(s.Object.AdditionalProperties)
		if err != nil {
			return fmt.Errorf("failed to convert additionalProperties schema: %w", err)
		}
		if addPropSchemaMap != nil {
			result["additionalProperties"] = addPropSchemaMap
		} else {
			result["additionalProperties"] = true
		}

	}
	return nil
}
