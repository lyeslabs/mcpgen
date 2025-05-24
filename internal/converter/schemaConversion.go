package converter

import (
	"encoding/json"
	"fmt"
)

// GenerateJSONSchemaDraft7 converts a slice of Arg structs into a JSON Schema Draft 7 string.
// It creates a root object schema with properties for each argument.
func GenerateJSONSchemaDraft7(args []Arg) (string, error) {
	rootSchema := map[string]interface{}{
		"type": "object",
	}

	properties := make(map[string]interface{})
	requiredProperties := []string{}

	for _, arg := range args {
		propSchema, err := buildPropertySchema(arg)
		if err != nil {
			return "", err
		}
		if propSchema == nil {
			continue
		}

		properties[arg.Name] = propSchema

		if arg.Required {
			requiredProperties = append(requiredProperties, arg.Name)
		}
	}

	if len(properties) > 0 {
		rootSchema["properties"] = properties
	}
	if len(requiredProperties) > 0 {
		rootSchema["required"] = requiredProperties
	}

	schemaBytes, err := json.MarshalIndent(rootSchema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON schema: %w", err)
	}

	return string(schemaBytes), nil
}


