package converter

import (
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
)

// createResponseTemplates generates Markdown response templates for an operation.
func (c *Converter) createResponseTemplates(
	operation *openapi3.Operation,
) ([]ResponseTemplate, error) {
	if operation == nil || operation.Responses == nil {
		return nil, nil
	}

	sortedCodes := sortedResponseCodes(operation.Responses)
	var templates []ResponseTemplate

	for _, code := range sortedCodes {
		responseRef := operation.Responses.Map()[code]
		if responseRef == nil || responseRef.Value == nil {
			continue
		}
		statusCode, _ := strconv.Atoi(code)
		contentTypes := sortedContentTypes(responseRef.Value.Content)

		for _, contentType := range contentTypes {
			mediaType := responseRef.Value.Content[contentType]
			if !hasSchema(mediaType) {
				continue
			}
			schema := mediaType.Schema.Value
			markdown := c.buildResponseMarkdown(code, contentType, responseRef, schema)
			templates = append(templates, ResponseTemplate{
				PrependBody: markdown,
				StatusCode:  statusCode,
				ContentType: contentType,
			})
		}
	}
	return assignSuffixes(templates), nil
}
