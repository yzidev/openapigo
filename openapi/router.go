package openapi

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
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

// Router uses a small net/http-backed mux that understands oas path params
// of the form /users/{id}. This avoids depending on chi while preserving
// behavior needed by the OpenAPI builder (path param extraction via context).
type Router struct {
	Mux    *httpMux
	config Config
	routes []RouteMeta
}

// Get registers a plain GET handler (used by Swagger UI mount helpers).
func (r *Router) Get(path string, h http.HandlerFunc) {
	r.Mux.MethodFunc(http.MethodGet, path, h)
}

// New creates an OpenAPI-aware net/http router.
//
// Typical usage:
//
//	r := openapi.New(openapi.Config{Title: "User API", Version: "1.0.0"})
//	r.GET("/users", listUsers, openapi.Res([]User{}), openapi.Tags("Users"))
//	r.Docs()
func New(cfg ...Config) *Router {
	r := &Router{
		Mux: newHTTPMux(),
	}
	if len(cfg) > 0 {
		r.config = cfg[0]
	}
	return r
}

func NewRouter() *Router {
	return New()
}

// Docs registers the OpenAPI JSON endpoint and Swagger UI using the router config.
// Pass a Config here to override or set the config after New().
func (r *Router) Docs(cfg ...Config) {
	if len(cfg) > 0 {
		r.config = cfg[0]
	}
	Register(r, r.config)
}

// HandlerOption configures RouteMeta.
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

	// Register handler on underlying mux. We wrap the provided handler so the
	// request context contains the path params under pathParamsKey, just like
	// previous chi-based implementation expected.
	r.Mux.MethodFunc(method, path, func(w http.ResponseWriter, req *http.Request) {
		// extract params and set into context
		params := extractPathParams(path, req.URL.Path)
		ctx := context.WithValue(req.Context(), pathParamsKey, params)
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

// extractPathParams extracts params from pattern like /users/{id} and actual
// path like /users/123 returning map[string]string{"id":"123"}.
func extractPathParams(pattern, actual string) map[string]string {
	pParts := splitPath(pattern)
	aParts := splitPath(actual)
	out := map[string]string{}
	if len(pParts) != len(aParts) {
		return out
	}
	for i := range pParts {
		pp := pParts[i]
		ap := aParts[i]
		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			name := strings.Trim(pp, "{}")
			if name != "" {
				out[name] = ap
			}
		}
	}
	return out
}

func splitPath(p string) []string {
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return []string{}
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	parts := strings.Split(p, "/")
	// drop leading empty from split
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	return parts
}

// httpMux is a tiny net/http-based router supporting MethodFunc and path
// patterns with {param} placeholders.
type httpMux struct {
	mu     sync.RWMutex
	routes []httpRoute
}

type httpRoute struct {
	method  string
	pattern string
	h       http.HandlerFunc
}

func newHTTPMux() *httpMux { return &httpMux{} }

// MethodFunc registers a handler for a method+pattern.
func (m *httpMux) MethodFunc(method, pattern string, h http.HandlerFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.routes = append(m.routes, httpRoute{method: method, pattern: pattern, h: h})
}

// ServeHTTP matches first registered route with same method and compatible pattern.
func (m *httpMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.mu.RLock()
	routes := append([]httpRoute(nil), m.routes...)
	m.mu.RUnlock()
	path := req.URL.Path
	for _, rt := range routes {
		if rt.method != req.Method {
			continue
		}
		if matchPattern(rt.pattern, path) {
			rt.h(w, req)
			return
		}
	}
	// default: 404
	http.NotFound(w, req)
}

func matchPattern(pattern, actual string) bool {
	pParts := splitPath(pattern)
	aParts := splitPath(actual)
	if len(pParts) != len(aParts) {
		return false
	}
	for i := range pParts {
		pp := pParts[i]
		ap := aParts[i]
		if strings.HasPrefix(pp, "{") && strings.HasSuffix(pp, "}") {
			continue
		}
		if pp != ap {
			return false
		}
	}
	return true
}

func getChiURLParam(r *http.Request) map[string]string {
	// Backwards-compatible helper: attempt to read params placed into context by our mux.
	if v := r.Context().Value(pathParamsKey); v != nil {
		if m, ok := v.(map[string]string); ok {
			return m
		}
	}
	return nil
}
