package mux

import (
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
)

// Router is the default net/http router implementation.
//
// Today it re-exports the chi-based implementation from the openapi package.
// Later we can make openapi core framework-agnostic and keep the net/http router here.
type Router = openapi.Router

// NewHttpAdapters creates a new adapter Router. If an *http.ServeMux is provided it will
// automatically mount the created router on the mux under the root path "/",
// so callers can call `base := httprouter.NewHttpAdapters(mux)` and skip `mux.Handle("/", r)`.
func NewHttpAdapters(mux ...*http.ServeMux) *Router {
	r := openapi.NewRouter()
	if len(mux) > 0 && mux[0] != nil {
		mux[0].Handle("/", r)
	}
	return r
}

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
