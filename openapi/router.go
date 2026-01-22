package openapi

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

type RouteMeta struct {
	Method         string
	Path           string
	Handler        http.HandlerFunc
	Summary        string
	Description    string
	Tags           []string
	RequestSchema  interface{}
	ResponseSchema interface{}
	Responses      []ResponseSpec
	Security       *openapi3.SecurityRequirement
	QueryParams    []QueryParam
	HeaderParams   []HeaderParam
	PathParams     []PathParamSpec
}

type Router struct {
	Mux    *chi.Mux
	routes []RouteMeta
}

func NewRouter() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

type HandlerOption func(*RouteMeta)

func WithRequestSchema(schema interface{}) HandlerOption {
	return func(meta *RouteMeta) {
		meta.RequestSchema = schema
	}
}

func WithResponseSchema(schema interface{}) HandlerOption {
	return func(meta *RouteMeta) {
		meta.ResponseSchema = schema
	}
}

func WithSecurity(security *openapi3.SecurityRequirement) HandlerOption {
	return func(meta *RouteMeta) {
		meta.Security = security
	}
}

func WithTags(tags ...string) HandlerOption {
	return func(meta *RouteMeta) {
		meta.Tags = append(meta.Tags, tags...)
	}
}

func (r *Router) Handle(method, path string, h http.HandlerFunc, opts ...HandlerOption) {
	meta := RouteMeta{
		Method:  method,
		Path:    path,
		Handler: h,
	}
	for _, opt := range opts {
		opt(&meta)
	}
	r.routes = append(r.routes, meta)

	r.Mux.MethodFunc(method, path, func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), pathParamsKey, getChiURLParam(req))
		h(w, req.WithContext(ctx))
	})
}

func (r *Router) GET(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}

func (r *Router) POST(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}

func (r *Router) PUT(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}

func (r *Router) DELETE(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}

func (r *Router) PATCH(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}

func (r *Router) HEAD(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodHead, path, h, opts...)
}

func (r *Router) OPTIONS(path string, h http.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodOptions, path, h, opts...)
}

func (r *Router) Routes() []RouteMeta {
	return r.routes
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Mux.ServeHTTP(w, req)
}

func getChiURLParam(r *http.Request) map[string]string {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return nil
	}
	params := make(map[string]string)
	for i, key := range rctx.URLParams.Keys {
		params[key] = rctx.URLParams.Values[i]
	}
	return params
}
