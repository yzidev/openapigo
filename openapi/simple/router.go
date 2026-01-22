package simple

import (
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
)

// Router wraps an openapi.Router and injects options from Spec automatically.
// Your route registrations can stay as plain GET/POST/... without JSONRoute/With... per route.
type Router struct {
	Base *openapi.Router
	Spec Spec
}

func New(base *openapi.Router, spec Spec) *Router {
	return &Router{Base: base, Spec: spec}
}

func (r *Router) Routes() []openapi.RouteMeta                        { return r.Base.Routes() }
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) { r.Base.ServeHTTP(w, req) }

func (r *Router) Group(prefix string, opts ...openapi.HandlerOption) *openapi.Group {
	return r.Base.Group(prefix, opts...)
}

func (r *Router) Handle(method, path string, h http.HandlerFunc, opts ...openapi.HandlerOption) {
	all := opts
	if def, ok := r.Spec[Key(method, path)]; ok {
		all = Inject(opts, def)
	}
	// Note: calling Base.Handle keeps chi path param capturing behavior.
	r.Base.Handle(method, path, h, all...)
}

func (r *Router) GET(path string, h http.HandlerFunc, opts ...openapi.HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *Router) POST(path string, h http.HandlerFunc, opts ...openapi.HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *Router) PUT(path string, h http.HandlerFunc, opts ...openapi.HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *Router) PATCH(path string, h http.HandlerFunc, opts ...openapi.HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *Router) DELETE(path string, h http.HandlerFunc, opts ...openapi.HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
