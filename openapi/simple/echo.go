//go:build echo

package simple

import (
	"net/http"
	"strings"

	echoadapter "github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
	echolib "github.com/labstack/echo/v4"
)

// EchoRouter wraps the echo adapter Router and injects options from Spec automatically.
type EchoRouter struct {
	Base *echoadapter.Router
	Spec Spec
}

type EchoGroup struct {
	prefix string
	opts   []echoadapter.HandlerOption
	r      *EchoRouter
}

func NewEcho(base *echoadapter.Router, spec Spec) *EchoRouter {
	return &EchoRouter{Base: base, Spec: spec}
}

func (r *EchoRouter) Routes() []openapi.RouteMeta { return r.Base.Routes() }

func (r *EchoRouter) Group(prefix string, opts ...echoadapter.HandlerOption) *EchoGroup {
	return &EchoGroup{prefix: prefix, opts: opts, r: r}
}

func joinEcho(prefix, p string) string {
	if prefix == "" {
		return p
	}
	if p == "" {
		return prefix
	}
	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return prefix + p
}

func (g *EchoGroup) Handle(method, path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	all := make([]echoadapter.HandlerOption, 0, len(g.opts)+len(opts))
	all = append(all, g.opts...)
	all = append(all, opts...)
	g.r.Handle(method, joinEcho(g.prefix, path), h, all...)
}

func (g *EchoGroup) GET(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	g.Handle(http.MethodGet, path, h, opts...)
}
func (g *EchoGroup) POST(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	g.Handle(http.MethodPost, path, h, opts...)
}
func (g *EchoGroup) PUT(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	g.Handle(http.MethodPut, path, h, opts...)
}
func (g *EchoGroup) PATCH(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	g.Handle(http.MethodPatch, path, h, opts...)
}
func (g *EchoGroup) DELETE(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	g.Handle(http.MethodDelete, path, h, opts...)
}

func (r *EchoRouter) Handle(method, path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	all := make([]openapi.HandlerOption, 0, len(opts))
	for _, o := range opts {
		all = append(all, o)
	}
	if def, ok := r.Spec[Key(method, path)]; ok {
		all = Inject(all, def)
	}
	out := make([]echoadapter.HandlerOption, 0, len(all))
	for _, o := range all {
		out = append(out, o)
	}
	r.Base.Handle(method, path, h, out...)
}

func (r *EchoRouter) GET(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *EchoRouter) POST(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *EchoRouter) PUT(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *EchoRouter) PATCH(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *EchoRouter) DELETE(path string, h echolib.HandlerFunc, opts ...echoadapter.HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
