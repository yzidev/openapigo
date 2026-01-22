package openapi

import (
	"net/http"
	"path"
	"strings"
)

// Group allows applying shared options (e.g., WithTags) and a common path prefix
// to multiple routes.
//
// Example:
//
//	api := r.Group("", WithTags("Users"))
//	api.GET("/users", ...)
//
// Options provided to the group are applied to every route in that group,
// and can be overridden/extended by per-route options.
//
// Path joining is kept simple and consistent with common router behavior.
// We avoid cleaning ":" or "{}" segments; NormalizePath handles OpenAPI normalization.
type Group struct {
	prefix string
	opts   []HandlerOption
	route  func(method, path string, h http.HandlerFunc, opts ...HandlerOption)
}

func (r *Router) Group(prefix string, opts ...HandlerOption) *Group {
	return &Group{
		prefix: prefix,
		opts:   opts,
		route:  r.Handle,
	}
}

func (g *Group) join(p string) string {
	if g.prefix == "" {
		return p
	}
	if p == "" {
		return g.prefix
	}
	// path.Join cleans slashes; keep trailing slash behavior minimal.
	j := path.Join(g.prefix, p)
	if strings.HasSuffix(p, "/") && !strings.HasSuffix(j, "/") {
		j += "/"
	}
	if !strings.HasPrefix(j, "/") {
		j = "/" + j
	}
	return j
}

func (g *Group) Handle(method, p string, h http.HandlerFunc, opts ...HandlerOption) {
	all := make([]HandlerOption, 0, len(g.opts)+len(opts))
	all = append(all, g.opts...)
	all = append(all, opts...)
	g.route(method, g.join(p), h, all...)
}

func (g *Group) GET(p string, h http.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodGet, p, h, opts...)
}
func (g *Group) POST(p string, h http.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPost, p, h, opts...)
}
func (g *Group) PUT(p string, h http.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPut, p, h, opts...)
}
func (g *Group) PATCH(p string, h http.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPatch, p, h, opts...)
}
func (g *Group) DELETE(p string, h http.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodDelete, p, h, opts...)
}
