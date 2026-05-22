package muxadapter

import (
	"net/http"

	"github.com/yzidev/goas/openapi"
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
	return New(mux...)
}

// New creates the default net/http Goas router adapter.
func New(mux ...*http.ServeMux) *Router {
	r := openapi.NewRouter()
	if len(mux) > 0 && mux[0] != nil {
		mux[0].Handle("/", r)
	}
	return r
}

// Docs mounts OpenAPI JSON and Swagger UI for a net/http Goas router.
func Docs(r *Router, cfg openapi.Config) {
	if r == nil {
		return
	}
	r.Docs(cfg)
}

// AutoDocs is an alias for Docs.
func AutoDocs(r *Router, cfg openapi.Config) {
	Docs(r, cfg)
}

// Mount creates an Goas router, mounts it on mux, and registers docs.
// Use this for a Springdoc-like setup with net/http:
//
//	r := muxadapter.Mount(mux, openapi.Config{Title: "API", Version: "1.0.0"})
//	r.GET("/users", listUsers)
func Mount(mux *http.ServeMux, cfg openapi.Config) *Router {
	r := New(mux)
	r.Docs(cfg)
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
	Req                = openapi.Req
	MultipartUpload    = openapi.MultipartUpload
	Res                = openapi.Res
	Tags               = openapi.Tags
	Security           = openapi.Security
	Query              = openapi.Query
	Headers            = openapi.Headers
	Status             = openapi.Status
	Created            = openapi.Created
	NoContent          = openapi.NoContent
	Responses          = openapi.Responses
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
