// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// RegisterResponse defines model for RegisterResponse.
type RegisterResponse struct {
	Id int64 `json:"id"`
}

// UpdateUserRequest defines model for UpdateUserRequest.
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// UpdateUserResponse defines model for UpdateUserResponse.
type UpdateUserResponse struct {
	Id int64 `json:"id"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// GetUserJSONRequestBody defines body for GetUser for application/json ContentType.
type GetUserJSONRequestBody = RegisterRequest

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserJSONRequestBody = UpdateUserRequest

// UserRegisterJSONRequestBody defines body for UserRegister for application/json ContentType.
type UserRegisterJSONRequestBody = RegisterRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Login user.
	// (POST /login)
	Login(ctx echo.Context) error
	// Get user data.
	// (GET /users)
	GetUser(ctx echo.Context) error
	// Update user data.
	// (PATCH /users)
	UpdateUser(ctx echo.Context) error
	// Register a new user.
	// (POST /users/register)
	UserRegister(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Login(ctx)
	return err
}

// GetUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUser(ctx)
	return err
}

// UpdateUser converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateUser(ctx)
	return err
}

// UserRegister converts echo context to params.
func (w *ServerInterfaceWrapper) UserRegister(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UserRegister(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/login", wrapper.Login)
	router.GET(baseURL+"/users", wrapper.GetUser)
	router.PATCH(baseURL+"/users", wrapper.UpdateUser)
	router.POST(baseURL+"/users/register", wrapper.UserRegister)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xWwW7bOhD8FWHfOwqW2wQ96BigDQK0l6Q5BTkw1FpiKpHMcuXAMPTvBUnbqiPZLlLb",
	"KHqTxfXO7Ax37CVI01ijUbODfAlOVtiI8PiZyNAtOmu0Q//CkrFIrDAcN+icKMMBLyxCDo5J6RK6LgXC",
	"l1YRFpA/bAof03WheXpGydCl8NWUSt/iS4uOhxBWOPdqqBjBSMFWRv8GeixL+157aOwaVQUGM0ONYMhB",
	"af50CZs2SjOWSL4Pmx+oD3NSBaxrx9jcYqkcI+3URYsGxzU5omAB5IBuPdM/lG4o0BjcvS0E4717lzTv",
	"mP4Qh/MMvRfqxPP6cqVnxjeqlcQVh4gK326+h1uvuPYfPdPkDmmupG85R3LKaMjhw2Q6mfpKY1ELqyCH",
	"i/DKXy+uwiBZ7XcwDGiisX5MwcromwLyuKIQyaPjK1MsfJE0mlGHemFtrWT4RvbsjO7jzD/9TziDHP7L",
	"+rzLVmGXbaVQty0RU4vhRXQgcP04nR4be+VvAC/QSVKWo3h3rZToXBL16VK4PCL6dsaPoF+JIumVScG1",
	"TSNosTYkaR3SJJxk/jHIU+KIf9fI/n6cyMG3kXlmE7d2dI+HJXJQLCkEi+jlxfm8/GLoSRUF6jdOXv/K",
	"ahJ/R1hWQw/77DuRjcOAP7eRw3TfY2df85c4GelvmbnZzIxWO7I7Y+Pcq6p/clEHf1j2uCsJvZgaX4Og",
	"Z4/eGz0XtdoVv+tJErFhOIldHNI8JPHDElqqIYeK2eZZVhsp6so73z12PwMAAP//ASahZfkLAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
