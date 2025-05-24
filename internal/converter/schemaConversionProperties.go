package converter

import "fmt"

// buildPropertySchema builds the JSON Schema property for a given Arg.
// Returns nil if the property should be skipped.
func buildPropertySchema(arg Arg) (map[string]interface{}, error) {
	var propSchema map[string]interface{}
	var err error

	switch arg.Source {
	case "body":
		propSchema, err = buildBodySchema(arg)
	default:
		if arg.Schema == nil {
			return nil, nil
		}
		propSchema, err = schemaToDraft7Map(arg.Schema)
	}
	if err != nil || propSchema == nil {
		return propSchema, err
	}

	if arg.Description != "" {
		if _, hasDesc := propSchema["description"]; !hasDesc {
			propSchema["description"] = arg.Description
		} else if propSchema["description"] == "" {
			propSchema["description"] = arg.Description
		}
	}

	return propSchema, nil
}

// buildBodySchema handles the "body" source, including multiple content types.
func buildBodySchema(arg Arg) (map[string]interface{}, error) {
	if len(arg.ContentTypes) == 0 {
		return nil, nil
	}
	if len(arg.ContentTypes) == 1 {
		// Only one content type, use its schema directly
		for _, schema := range arg.ContentTypes {
			return schemaToDraft7Map(schema)
		}
	}

	// Multiple content types: use oneOf
	oneOfSchemas := []map[string]interface{}{}
	for contentType, schema := range arg.ContentTypes {
		branchSchema, err := schemaToDraft7Map(schema)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to convert body schema branch for content type '%s': %w",
				contentType, err,
			)
		}
		if branchSchema != nil {
			// Add content type info to title/description
			addContentTypeInfo(branchSchema, contentType)
			oneOfSchemas = append(oneOfSchemas, branchSchema)
		}
	}
	if len(oneOfSchemas) == 0 {
		return nil, nil
	}
	return map[string]interface{}{
		"oneOf": oneOfSchemas,
	}, nil
}

// addContentTypeInfo adds content type info to the schema's title or description.
func addContentTypeInfo(schema map[string]interface{}, contentType string) {
	if desc, ok := schema["description"].(string); ok {
		schema["description"] = fmt.Sprintf("[%s] %s", contentType, desc)
	} else if title, ok := schema["title"].(string); ok {
		schema["title"] = fmt.Sprintf("[%s] %s", contentType, title)
	} else {
		schema["title"] = fmt.Sprintf("Schema for %s", contentType)
	}
}
