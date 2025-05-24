package converter

// MCPConfig represents the top-level MCP server configuration
type MCPConfig struct {
	Server ServerConfig
	Tools  []Tool
}

// ServerConfig represents the MCP server configuration
type ServerConfig struct {
	Config          map[string]interface{}
	SecuritySchemes []SecurityScheme
}

// SecurityScheme defines a security scheme that can be used by the tools.
type SecurityScheme struct {
	ID                string
	Type              string // e.g., "http", "apiKey", "oauth2", "openIdConnect"
	Scheme            string // e.g., "basic", "bearer" for "http" type
	In                string // e.g., "header", "query", "cookie" for "apiKey" type
	Name              string // Name of the header, query parameter or cookie for "apiKey" type
	DefaultCredential string
}

// Tool represents an MCP tool configuration
type Tool struct {
	Name            string
	Description     string
	Args            []Arg
	RequestTemplate RequestTemplate
	Responses       []ResponseTemplate
	RawInputSchema  string
}

// RequestTemplate represents the MCP request template
type RequestTemplate struct {
	URL            string
	Method         string
	Headers        []Header
	Body           string
	ArgsToJsonBody bool
	ArgsToUrlParam bool
	ArgsToFormBody bool
	Security       []ToolSecurityRequirement
}

// ToolSecurityRequirement specifies a security scheme requirement for a tool.
type ToolSecurityRequirement struct {
	ID string
}

// Header represents an HTTP header
type Header struct {
	Key   string
	Value string
}

// ResponseTemplate represents the MCP response template
type ResponseTemplate struct {
	PrependBody string
	StatusCode  int
	ContentType string
	Suffix       string 
}

// ConvertOptions represents options for the conversion process
type ConvertOptions struct {
	ServerConfig map[string]interface{}
}

// ToolTemplate represents a template for applying to all tools
type ToolTemplate struct {
	RequestTemplate  *RequestTemplate
	ResponseTemplate *ResponseTemplate
}

// MCPConfigTemplate represents a template for patching the generated config
type MCPConfigTemplate struct {
	Server ServerConfig
	Tools  ToolTemplate
}

// Arg represents an argument in an API, which can come from path, query, or body
type Arg struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Source      string  `json:"source"` // "path", "query", "header", "body"
	Required    bool    `json:"required"`
	Deprecated  bool    `json:"deprecated,omitempty"`
	Schema      *Schema `json:"schema"`
	// For request bodies with multiple content types
	ContentTypes map[string]*Schema `json:"contentTypes,omitempty"`
}

// Schema represents the structure and validation rules for data
type Schema struct {
	Types       []string          `json:"types"`
	OneOf       []*Schema         `json:"oneOf,omitempty"`
	AnyOf       []*Schema         `json:"anyOf,omitempty"`
	AllOf       []*Schema         `json:"allOf,omitempty"`
	Not         *Schema           `json:"not,omitempty"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Format      string            `json:"format,omitempty"`
	Default     interface{}       `json:"default,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Enum        []interface{}     `json:"enum,omitempty"`
	ReadOnly    bool              `json:"readOnly,omitempty"`
	WriteOnly   bool              `json:"writeOnly,omitempty"`
	String      *StringValidation `json:"string,omitempty"`
	Number      *NumberValidation `json:"number,omitempty"`
	Array       *ArrayValidation  `json:"array,omitempty"`
	Object      *ObjectValidation `json:"object,omitempty"`
}

// StringValidation contains validation rules specific to string types
type StringValidation struct {
	MinLength uint64  `json:"minLength,omitempty"`
	MaxLength *uint64 `json:"maxLength,omitempty"`
	Pattern   string  `json:"pattern,omitempty"`
}

// NumberValidation contains validation rules specific to number types
type NumberValidation struct {
	Minimum          *float64 `json:"minimum,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	ExclusiveMinimum bool     `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum bool     `json:"exclusiveMaximum,omitempty"`
}

// ArrayValidation contains validation rules specific to array types
type ArrayValidation struct {
	Items       *Schema `json:"items,omitempty"`
	MinItems    uint64  `json:"minItems,omitempty"`
	MaxItems    *uint64 `json:"maxItems,omitempty"`
	UniqueItems bool    `json:"uniqueItems,omitempty"`
}

// ObjectValidation contains validation rules specific to object types
type ObjectValidation struct {
	Properties                   map[string]*Schema `json:"properties,omitempty"`
	AdditionalProperties         *Schema            `json:"additionalProperties,omitempty"`         // Schema if specified or represents {} for `true`
	DisallowAdditionalProperties bool               `json:"disallowAdditionalProperties,omitempty"` // True if additionalProperties: false
	Required                     []string           `json:"required,omitempty"`
	MinProperties                uint64             `json:"minProperties,omitempty"`
	MaxProperties                *uint64            `json:"maxProperties,omitempty"`
}
