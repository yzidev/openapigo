//go:build gin

package simple

import (
	"net/http"
	"strings"

	ginadapter "github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	ginlib "github.com/gin-gonic/gin"
)

// GinRouter wraps the gin adapter Router and injects options from Spec automatically.
type GinRouter struct {
	Base *ginadapter.Router
	Spec Spec
}

// GinGroup provides grouping with prefix + shared options, while preserving Spec injection.
type GinGroup struct {
	prefix string
	opts   []ginadapter.HandlerOption
	r      *GinRouter
}

func NewGin(base *ginadapter.Router, spec Spec) *GinRouter {
	return &GinRouter{Base: base, Spec: spec}
}

func (r *GinRouter) Routes() []openapi.RouteMeta { return r.Base.Routes() }

func (r *GinRouter) Group(prefix string, opts ...ginadapter.HandlerOption) *GinGroup {
	return &GinGroup{prefix: prefix, opts: opts, r: r}
}

func joinGin(prefix, p string) string {
	if prefix == "" {
		return p
	}
	if p == "" {
		return prefix
	}
	// keep it simple: mimic adapter join behavior
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return prefix + p
}

func (g *GinGroup) Handle(method, path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	all := make([]ginadapter.HandlerOption, 0, len(g.opts)+len(opts))
	all = append(all, g.opts...)
	all = append(all, opts...)
	g.r.Handle(method, joinGin(g.prefix, path), h, all...)
}

func (g *GinGroup) GET(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	g.Handle(http.MethodGet, path, h, opts...)
}
func (g *GinGroup) POST(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	g.Handle(http.MethodPost, path, h, opts...)
}
func (g *GinGroup) PUT(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	g.Handle(http.MethodPut, path, h, opts...)
}
func (g *GinGroup) PATCH(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	g.Handle(http.MethodPatch, path, h, opts...)
}
func (g *GinGroup) DELETE(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	g.Handle(http.MethodDelete, path, h, opts...)
}

func (r *GinRouter) Handle(method, path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	all := make([]openapi.HandlerOption, 0, len(opts))
	for _, o := range opts {
		all = append(all, o)
	}
	if def, ok := r.Spec[Key(method, path)]; ok {
		all = Inject(all, def)
	}
	// Convert back to adapter options (same underlying type)
	out := make([]ginadapter.HandlerOption, 0, len(all))
	for _, o := range all {
		out = append(out, o)
	}
	r.Base.Handle(method, path, h, out...)
}

func (r *GinRouter) GET(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *GinRouter) POST(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *GinRouter) PUT(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *GinRouter) PATCH(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *GinRouter) DELETE(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
