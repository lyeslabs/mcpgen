package converter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// sortedResponseCodes returns response codes sorted numerically, then alphabetically.
func sortedResponseCodes(responses *openapi3.Responses) []string {
	var codes []string
	for code := range responses.Map() {
		codes = append(codes, code)
	}
	sort.SliceStable(codes, func(i, j int) bool {
		codeI, errI := strconv.Atoi(codes[i])
		codeJ, errJ := strconv.Atoi(codes[j])
		switch {
		case errI == nil && errJ == nil:
			return codeI < codeJ
		case errI != nil && errJ != nil:
			return codes[i] < codes[j]
		default:
			return errI == nil
		}
	})
	return codes
}

// sortedContentTypes returns sorted content types for consistent output.
func sortedContentTypes(content openapi3.Content) []string {
	var types []string
	for ct := range content {
		types = append(types, ct)
	}
	sort.Strings(types)
	return types
}

// assignSuffixes adds a unique letter suffix (_A, _B, ..., _Z, _AA, _AB, ...) to each response template.
func assignSuffixes(responses []ResponseTemplate) []ResponseTemplate {
	for i := range responses {
		responses[i].Suffix = toAlphaSuffix(i)
	}
	return responses
}

// toAlphaSuffix converts an integer to a base-26 alphabetic suffix (A, B, ..., Z, AA, AB, ...).
func toAlphaSuffix(n int) string {
    var b strings.Builder
    n++ // 1-based
    for n > 0 {
        n-- // 0-based for this digit
        b.WriteByte(byte('A' + (n % 26)))
        n /= 26
    }
    // The result is reversed, so reverse it
    s := b.String()
    // Reverse the string
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}


// formatForGoRawString safely formats values for Markdown.
func formatForGoRawString(schema *openapi3.Schema, value interface{}) string {
    str := fmt.Sprintf("%v", value)
    if schema.Type != nil && len(*schema.Type) > 0 && (*schema.Type)[0] == "string" {
        // If value is a string, use it directly
        if s, ok := value.(string); ok {
            str = s
        } else if bts, err := json.Marshal(value); err == nil {
            str = string(bts)
            if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
                str = strings.Trim(str, `"`)
            }
        }
        if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
            str = strings.Trim(str, `"`)
        }
    } else {
        // For non-string types, use JSON
        if bts, err := json.Marshal(value); err == nil {
            str = string(bts)
        }
    }
    str = strings.ReplaceAll(str, "`", "'")
    return str
}


// getResponseDescription safely returns the response description.
func getResponseDescription(responseRef *openapi3.ResponseRef) string {
	if responseRef.Value.Description != nil {
		return *responseRef.Value.Description
	}
	return ""
}

// getDescription returns a description for an operation.
func getDescription(operation *openapi3.Operation) string {
	if operation.Summary != "" {
		if operation.Description != "" {
			return fmt.Sprintf("%s - %s", operation.Summary, operation.Description)
		}
		return operation.Summary
	}
	return operation.Description
}

// contains checks if a string slice contains a string.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
