package converter

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// createRequestTemplate creates an MCP request template from an OpenAPI operation
func (c *Converter) createRequestTemplate(path, method string, operation *openapi3.Operation) (*RequestTemplate, error) {
	// Get the server URL from the OpenAPI specification
	var serverURL string
	if servers := c.parser.GetDocument().Servers; len(servers) > 0 {
		serverURL = servers[0].URL
	}

	// Remove trailing slash from server URL if present
	serverURL = strings.TrimSuffix(serverURL, "/")

	// Create the request template
	template := &RequestTemplate{
		URL:     serverURL + path,
		Method:  strings.ToUpper(method),
		Headers: []Header{},
	}

	// Add Content-Type header based on request body content type
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		for contentType := range operation.RequestBody.Value.Content {
			// Add the Content-Type header
			template.Headers = append(template.Headers, Header{
				Key:   "Content-Type",
				Value: contentType,
			})
			break // Just use the first content type
		}
	}

	return template, nil
}


