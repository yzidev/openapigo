package muxadapter

import (
	"net/http"

	"github.com/yzidev/goas"
)

// Router is the default net/http router implementation.
//
// Today it re-exports the net/http implementation from the goas package.
// Later we can make goas core framework-agnostic and keep the net/http router here.
type Router = goas.Router

// NewHttpAdapters creates a new adapter Router. If an *http.ServeMux is provided it will
// automatically mount the created router on the mux under the root path "/",
// so callers can call `base := httprouter.NewHttpAdapters(mux)` and skip `mux.Handle("/", r)`.
func NewHttpAdapters(mux ...*http.ServeMux) *Router {
	return New(mux...)
}

// New creates the default net/http Goas router adapter.
func New(mux ...*http.ServeMux) *Router {
	r := goas.NewRouter()
	if len(mux) > 0 && mux[0] != nil {
		mux[0].Handle("/", r)
	}
	return r
}

// Docs mounts OpenAPI JSON and Swagger UI for a net/http Goas router.
func Docs(r *Router, cfg goas.Config) {
	if r == nil {
		return
	}
	r.Docs(cfg)
}

// AutoDocs is an alias for Docs.
func AutoDocs(r *Router, cfg goas.Config) {
	Docs(r, cfg)
}

// Mount creates an Goas router, mounts it on mux, and registers docs.
// Use this for a Springdoc-like setup with net/http:
//
//	r := muxadapter.Mount(mux, goas.Config{Title: "API", Version: "1.0.0"})
//	r.GET("/users", listUsers)
func Mount(mux *http.ServeMux, cfg goas.Config) *Router {
	r := New(mux)
	r.Docs(cfg)
	return r
}

// HandlerOption Re-export route options for convenience.
type HandlerOption = goas.HandlerOption

var (
	WithRequestSchema  = goas.WithRequestSchema
	WithResponseSchema = goas.WithResponseSchema
	WithSecurity       = goas.WithSecurity
	WithTags           = goas.WithTags
	WithResponses      = goas.WithResponses
	WithQueryParams    = goas.WithQueryParams
	Req                = goas.Req
	MultipartUpload    = goas.MultipartUpload
	Res                = goas.Res
	Tags               = goas.Tags
	Security           = goas.Security
	Query              = goas.Query
	Headers            = goas.Headers
	Status             = goas.Status
	Created            = goas.Created
	NoContent          = goas.NoContent
	Responses          = goas.Responses
	JSONRoute          = goas.JSONRoute
)

// Re-export helpers.
var (
	Register  = goas.Register
	Bind      = goas.Bind
	JSON      = goas.JSON
	PathValue = goas.PathValue
	BuildSpec = goas.BuildSpec
	_         = http.MethodGet
)
