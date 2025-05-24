package converter

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ConvertRequestBody converts an OpenAPI request body to our Arg structures
func (c *Converter) convertRequestBody(requestBodyRef *openapi3.RequestBodyRef) (*Arg, error) {
	if requestBodyRef == nil || requestBodyRef.Value == nil {
		return nil, nil
	}

	requestBody := requestBodyRef.Value
	Arg := Arg{
		Name:         "body",
		Source:       "body",
		Description:  requestBody.Description,
		Required:     requestBody.Required,
		ContentTypes: make(map[string]*Schema),
	}

	// Process each content type
	validContent := false
	for contentType, mediaType := range requestBody.Content {
		if mediaType == nil || mediaType.Schema == nil || mediaType.Schema.Value == nil {
			continue
		}

		schema, err := c.applySchema(mediaType.Schema.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert schema for content type %s: %w", contentType, err)
		}

		if schema != nil {
			Arg.ContentTypes[contentType] = schema
			validContent = true
		}
	}

	if validContent {
		return &Arg, nil
	}

	return nil, nil
}

// ConvertParameters converts OpenAPI parameters to our Arg structures
func (c *Converter) convertParameters(parameters openapi3.Parameters) ([]Arg, error) {
	args := []Arg{}

	for i, paramRef := range parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value

		// Skip invalid parameters
		if param.Schema == nil || param.Schema.Value == nil {
			continue
		}

		// Convert the schema using our new function
		schema, err := c.applySchema(param.Schema.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert schema for parameter %s (index %d): %w",
				param.Name, i, err)
		}

		// Create an arg for this parameter
		arg := Arg{
			Name:        param.Name,
			Description: param.Description,
			Source:      param.In,
			Required:    param.Required,
			Schema:      schema,
			Deprecated:  param.Deprecated,
		}

		args = append(args, arg)
	}

	return args, nil
}
