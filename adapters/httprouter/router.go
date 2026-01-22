package httprouter

import (
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
)

// Router is the default net/http router implementation.
//
// Today it re-exports the chi-based implementation from the openapi package.
// Later we can make openapi core framework-agnostic and keep the net/http router here.
type Router = openapi.Router

func New() *Router { return openapi.NewRouter() }

// HandlerOption Re-export route options for convenience.
type HandlerOption = openapi.HandlerOption

var (
	WithRequestSchema  = openapi.WithRequestSchema
	WithResponseSchema = openapi.WithResponseSchema
	WithSecurity       = openapi.WithSecurity
	WithTags           = openapi.WithTags
	WithResponses      = openapi.WithResponses
	WithQueryParams    = openapi.WithQueryParams
	JSONRoute          = openapi.JSONRoute
)

// Re-export helpers.
var (
	Register  = openapi.Register
	Bind      = openapi.Bind
	JSON      = openapi.JSON
	PathValue = openapi.PathValue
	BuildSpec = openapi.BuildSpec
	_         = http.MethodGet
)
