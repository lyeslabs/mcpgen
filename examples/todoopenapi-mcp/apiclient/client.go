// Package apiclient provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/lyeslabs/mcpgen version v0.0.0-20250525151702-6ab30df4cde6 DO NOT EDIT.
package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// Defines values for NewTodoStatus.
const (
	NewTodoStatusCompleted  NewTodoStatus = "completed"
	NewTodoStatusInProgress NewTodoStatus = "in-progress"
	NewTodoStatusPending    NewTodoStatus = "pending"
)

// Defines values for TodoStatus.
const (
	TodoStatusCompleted  TodoStatus = "completed"
	TodoStatusInProgress TodoStatus = "in-progress"
	TodoStatusPending    TodoStatus = "pending"
)

// Defines values for UpdateTodoStatus.
const (
	UpdateTodoStatusCompleted  UpdateTodoStatus = "completed"
	UpdateTodoStatusInProgress UpdateTodoStatus = "in-progress"
	UpdateTodoStatusPending    UpdateTodoStatus = "pending"
)

// Defines values for ListTodosParamsStatus.
const (
	ListTodosParamsStatusCompleted  ListTodosParamsStatus = "completed"
	ListTodosParamsStatusInProgress ListTodosParamsStatus = "in-progress"
	ListTodosParamsStatusPending    ListTodosParamsStatus = "pending"
)

// Defines values for CreateTodoJSONBodyPriority.
const (
	High   CreateTodoJSONBodyPriority = "high"
	Low    CreateTodoJSONBodyPriority = "low"
	Medium CreateTodoJSONBodyPriority = "medium"
)

// Error defines model for Error.
type Error struct {
	// Code An application-specific error code.
	Code int32 `json:"code"`

	// Details Optional array of specific field validation errors.
	Details *[]struct {
		Field *string `json:"field,omitempty"`
		Issue *string `json:"issue,omitempty"`
	} `json:"details,omitempty"`

	// Message A human-readable description of the error.
	Message string `json:"message"`
}

// NewTodo defines model for NewTodo.
type NewTodo struct {
	// Description Optional detailed description of the todo item.
	Description *string `json:"description"`

	// DueDate Optional due date for the todo item.
	DueDate *openapi_types.Date `json:"dueDate"`

	// Status Current status of the todo item.
	Status *NewTodoStatus `json:"status,omitempty"`

	// Title The main content of the todo item.
	Title string `json:"title"`
}

// NewTodoStatus Current status of the todo item.
type NewTodoStatus string

// Todo defines model for Todo.
type Todo struct {
	// CreatedAt Timestamp of when the todo item was created.
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// Id Unique identifier for the todo item.
	Id openapi_types.UUID `json:"id"`

	// Status Current status of the todo item.
	Status TodoStatus `json:"status"`

	// Title The main content of the todo item.
	Title string `json:"title"`

	// UpdatedAt Timestamp of when the todo item was last updated.
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// TodoStatus Current status of the todo item.
type TodoStatus string

// UpdateTodo defines model for UpdateTodo.
type UpdateTodo struct {
	// Description Optional detailed description of the todo item.
	Description *string `json:"description"`

	// DueDate Optional due date for the todo item.
	DueDate *openapi_types.Date `json:"dueDate"`

	// Status Current status of the todo item.
	Status *UpdateTodoStatus `json:"status,omitempty"`

	// Title The main content of the todo item.
	Title *string `json:"title,omitempty"`
}

// UpdateTodoStatus Current status of the todo item.
type UpdateTodoStatus string

// BadRequest defines model for BadRequest.
type BadRequest = Error

// InternalServerError defines model for InternalServerError.
type InternalServerError = Error

// NotFound defines model for NotFound.
type NotFound = Error

// UnprocessableEntity defines model for UnprocessableEntity.
type UnprocessableEntity = Error

// ListTodosParams defines parameters for ListTodos.
type ListTodosParams struct {
	// Status Filter todos by status (e.g., "pending", "completed")
	Status *ListTodosParamsStatus `form:"status,omitempty" json:"status,omitempty"`

	// Limit Maximum number of todos to return
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`

	// Offset Number of todos to skip for pagination
	Offset *int32 `form:"offset,omitempty" json:"offset,omitempty"`

	// Token Token for authentication
	Token *int32 `form:"token,omitempty" json:"token,omitempty"`
}

// ListTodosParamsStatus defines parameters for ListTodos.
type ListTodosParamsStatus string

// CreateTodoJSONBody defines parameters for CreateTodo.
type CreateTodoJSONBody struct {
	Priority *CreateTodoJSONBodyPriority `json:"priority,omitempty"`
	Title    string                      `json:"title"`
}

// CreateTodoJSONBodyPriority defines parameters for CreateTodo.
type CreateTodoJSONBodyPriority string

// CreateTodoJSONRequestBody defines body for CreateTodo for application/json ContentType.
type CreateTodoJSONRequestBody CreateTodoJSONBody

// UpdateTodoByIdJSONRequestBody defines body for UpdateTodoById for application/json ContentType.
type UpdateTodoByIdJSONRequestBody = UpdateTodo

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ListTodos request
	ListTodos(ctx context.Context, params *ListTodosParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateTodoWithBody request with any body
	CreateTodoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateTodo(ctx context.Context, body CreateTodoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteTodoById request
	DeleteTodoById(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetTodoById request
	GetTodoById(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateTodoByIdWithBody request with any body
	UpdateTodoByIdWithBody(ctx context.Context, todoId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateTodoById(ctx context.Context, todoId openapi_types.UUID, body UpdateTodoByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListTodos(ctx context.Context, params *ListTodosParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListTodosRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTodoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTodoRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateTodo(ctx context.Context, body CreateTodoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateTodoRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteTodoById(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteTodoByIdRequest(c.Server, todoId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetTodoById(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetTodoByIdRequest(c.Server, todoId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateTodoByIdWithBody(ctx context.Context, todoId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateTodoByIdRequestWithBody(c.Server, todoId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateTodoById(ctx context.Context, todoId openapi_types.UUID, body UpdateTodoByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateTodoByIdRequest(c.Server, todoId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListTodosRequest generates requests for ListTodos
func NewListTodosRequest(server string, params *ListTodosParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/todos")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Status != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "status", runtime.ParamLocationQuery, *params.Status); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Offset != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "offset", runtime.ParamLocationQuery, *params.Offset); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {

		if params.Token != nil {
			var cookieParam0 string

			cookieParam0, err = runtime.StyleParamWithLocation("simple", true, "token", runtime.ParamLocationCookie, *params.Token)
			if err != nil {
				return nil, err
			}

			cookie0 := &http.Cookie{
				Name:  "token",
				Value: cookieParam0,
			}
			req.AddCookie(cookie0)
		}
	}
	return req, nil
}

// NewCreateTodoRequest calls the generic CreateTodo builder with application/json body
func NewCreateTodoRequest(server string, body CreateTodoJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateTodoRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateTodoRequestWithBody generates requests for CreateTodo with any type of body
func NewCreateTodoRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/todos")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteTodoByIdRequest generates requests for DeleteTodoById
func NewDeleteTodoByIdRequest(server string, todoId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "todoId", runtime.ParamLocationPath, todoId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/todos/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetTodoByIdRequest generates requests for GetTodoById
func NewGetTodoByIdRequest(server string, todoId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "todoId", runtime.ParamLocationPath, todoId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/todos/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateTodoByIdRequest calls the generic UpdateTodoById builder with application/json body
func NewUpdateTodoByIdRequest(server string, todoId openapi_types.UUID, body UpdateTodoByIdJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateTodoByIdRequestWithBody(server, todoId, "application/json", bodyReader)
}

// NewUpdateTodoByIdRequestWithBody generates requests for UpdateTodoById with any type of body
func NewUpdateTodoByIdRequestWithBody(server string, todoId openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "todoId", runtime.ParamLocationPath, todoId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/todos/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ListTodosWithResponse request
	ListTodosWithResponse(ctx context.Context, params *ListTodosParams, reqEditors ...RequestEditorFn) (*ListTodosResponse, error)

	// CreateTodoWithBodyWithResponse request with any body
	CreateTodoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTodoResponse, error)

	CreateTodoWithResponse(ctx context.Context, body CreateTodoJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTodoResponse, error)

	// DeleteTodoByIdWithResponse request
	DeleteTodoByIdWithResponse(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteTodoByIdResponse, error)

	// GetTodoByIdWithResponse request
	GetTodoByIdWithResponse(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetTodoByIdResponse, error)

	// UpdateTodoByIdWithBodyWithResponse request with any body
	UpdateTodoByIdWithBodyWithResponse(ctx context.Context, todoId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTodoByIdResponse, error)

	UpdateTodoByIdWithResponse(ctx context.Context, todoId openapi_types.UUID, body UpdateTodoByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTodoByIdResponse, error)
}

type ListTodosResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]ListTodos_200_Item
	JSON400      *BadRequest
	JSON500      *InternalServerError
}
type ListTodos_200_Item struct {
	union json.RawMessage
}

// Status returns HTTPResponse.Status
func (r ListTodosResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListTodosResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateTodoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *struct {
		Completed *bool   `json:"completed,omitempty"`
		Id        *int    `json:"id,omitempty"`
		Title     *string `json:"title,omitempty"`
	}
	XML201 *struct {
		Completed *bool   `json:"completed,omitempty"`
		Id        *int    `json:"id,omitempty"`
		Title     *string `json:"title,omitempty"`
	}
	JSON207 *[]CreateTodo_207_Item
	JSON400 *struct {
		Details *[]string `json:"details,omitempty"`
		Message *string   `json:"message,omitempty"`
	}
	JSON422 *struct {
		Errors *[]CreateTodo_422_Errors_Item `json:"errors,omitempty"`
	}
	JSON500 *struct {
		Message *string `json:"message,omitempty"`
		TraceId *string `json:"traceId,omitempty"`
	}
}
type CreateTodo2070 struct {
	Id    *int    `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
}
type CreateTodo2071 struct {
	Error *string `json:"error,omitempty"`
}
type CreateTodo_207_Item struct {
	union json.RawMessage
}
type CreateTodo422Errors0 struct {
	Error *string `json:"error,omitempty"`
	Field *string `json:"field,omitempty"`
}
type CreateTodo422Errors1 = string
type CreateTodo_422_Errors_Item struct {
	union json.RawMessage
}

// Status returns HTTPResponse.Status
func (r CreateTodoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateTodoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteTodoByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON404      *NotFound
	JSON500      *InternalServerError
}

// Status returns HTTPResponse.Status
func (r DeleteTodoByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteTodoByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetTodoByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Todo
	JSON404      *NotFound
	JSON500      *InternalServerError
}

// Status returns HTTPResponse.Status
func (r GetTodoByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetTodoByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateTodoByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Todo
	JSON400      *BadRequest
	JSON404      *NotFound
	JSON422      *UnprocessableEntity
	JSON500      *InternalServerError
}

// Status returns HTTPResponse.Status
func (r UpdateTodoByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateTodoByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListTodosWithResponse request returning *ListTodosResponse
func (c *ClientWithResponses) ListTodosWithResponse(ctx context.Context, params *ListTodosParams, reqEditors ...RequestEditorFn) (*ListTodosResponse, error) {
	rsp, err := c.ListTodos(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListTodosResponse(rsp)
}

// CreateTodoWithBodyWithResponse request with arbitrary body returning *CreateTodoResponse
func (c *ClientWithResponses) CreateTodoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateTodoResponse, error) {
	rsp, err := c.CreateTodoWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTodoResponse(rsp)
}

func (c *ClientWithResponses) CreateTodoWithResponse(ctx context.Context, body CreateTodoJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateTodoResponse, error) {
	rsp, err := c.CreateTodo(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateTodoResponse(rsp)
}

// DeleteTodoByIdWithResponse request returning *DeleteTodoByIdResponse
func (c *ClientWithResponses) DeleteTodoByIdWithResponse(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteTodoByIdResponse, error) {
	rsp, err := c.DeleteTodoById(ctx, todoId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteTodoByIdResponse(rsp)
}

// GetTodoByIdWithResponse request returning *GetTodoByIdResponse
func (c *ClientWithResponses) GetTodoByIdWithResponse(ctx context.Context, todoId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetTodoByIdResponse, error) {
	rsp, err := c.GetTodoById(ctx, todoId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetTodoByIdResponse(rsp)
}

// UpdateTodoByIdWithBodyWithResponse request with arbitrary body returning *UpdateTodoByIdResponse
func (c *ClientWithResponses) UpdateTodoByIdWithBodyWithResponse(ctx context.Context, todoId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateTodoByIdResponse, error) {
	rsp, err := c.UpdateTodoByIdWithBody(ctx, todoId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTodoByIdResponse(rsp)
}

func (c *ClientWithResponses) UpdateTodoByIdWithResponse(ctx context.Context, todoId openapi_types.UUID, body UpdateTodoByIdJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateTodoByIdResponse, error) {
	rsp, err := c.UpdateTodoById(ctx, todoId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateTodoByIdResponse(rsp)
}

// ParseListTodosResponse parses an HTTP response from a ListTodosWithResponse call
func ParseListTodosResponse(rsp *http.Response) (*ListTodosResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListTodosResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []ListTodos_200_Item
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseCreateTodoResponse parses an HTTP response from a CreateTodoWithResponse call
func ParseCreateTodoResponse(rsp *http.Response) (*CreateTodoResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateTodoResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest struct {
			Completed *bool   `json:"completed,omitempty"`
			Id        *int    `json:"id,omitempty"`
			Title     *string `json:"title,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 207:
		var dest []CreateTodo_207_Item
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON207 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Details *[]string `json:"details,omitempty"`
			Message *string   `json:"message,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest struct {
			Errors *[]CreateTodo_422_Errors_Item `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON422 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Message *string `json:"message,omitempty"`
			TraceId *string `json:"traceId,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "xml") && rsp.StatusCode == 201:
		var dest struct {
			Completed *bool   `json:"completed,omitempty"`
			Id        *int    `json:"id,omitempty"`
			Title     *string `json:"title,omitempty"`
		}
		if err := xml.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.XML201 = &dest

	case rsp.StatusCode == 201:
	// Content-type (text/plain) unsupported

	case rsp.StatusCode == 500:
		// Content-type (text/plain) unsupported

	}

	return response, nil
}

// ParseDeleteTodoByIdResponse parses an HTTP response from a DeleteTodoByIdWithResponse call
func ParseDeleteTodoByIdResponse(rsp *http.Response) (*DeleteTodoByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteTodoByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetTodoByIdResponse parses an HTTP response from a GetTodoByIdWithResponse call
func ParseGetTodoByIdResponse(rsp *http.Response) (*GetTodoByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetTodoByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Todo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseUpdateTodoByIdResponse parses an HTTP response from a UpdateTodoByIdWithResponse call
func ParseUpdateTodoByIdResponse(rsp *http.Response) (*UpdateTodoByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateTodoByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Todo
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest UnprocessableEntity
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON422 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
