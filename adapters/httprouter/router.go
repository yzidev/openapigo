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
)

// Re-export helpers.
var (
	Register   = openapi.Register
	Bind       = openapi.Bind
	JSON       = openapi.JSON
	PathValue  = openapi.PathValue
	BuildSpec  = openapi.BuildSpec
	HTTPMethod = http.MethodGet
)

// TypedHandler enables full-auto schema via type parameters.
type TypedHandler[TReq any, TRes any] = openapi.TypedHandler[TReq, TRes]

// Full auto schema helpers (generic). These attach request/response schema automatically.
func GETT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	openapi.GETT[TReq, TRes](r, path, h, opts...)
}

func POSTT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	openapi.POSTT[TReq, TRes](r, path, h, opts...)
}

func PUTT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	openapi.PUTT[TReq, TRes](r, path, h, opts...)
}

func PATCHT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	openapi.PATCHT[TReq, TRes](r, path, h, opts...)
}

func DELETET[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	openapi.DELETET[TReq, TRes](r, path, h, opts...)
}
