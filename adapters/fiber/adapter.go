//go:build fiber

package fiber

import (
	"net/http"
	"strings"

	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/getkin/kin-openapi/openapi3"
)

type Router struct {
	App    *fiberlib.App
	routes []openapi.RouteMeta
}

func New() *Router {
	return &Router{App: fiberlib.New()}
}

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

func (r *Router) Handle(method, path string, h fiberlib.Handler, opts ...HandlerOption) {
	meta := openapi.RouteMeta{Method: method, Path: path}
	for _, opt := range opts {
		opt(&meta)
	}
	r.routes = append(r.routes, meta)

	r.App.Add(method, path, h)
}

func (r *Router) GET(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *Router) POST(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *Router) PUT(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *Router) DELETE(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
func (r *Router) PATCH(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *Router) HEAD(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodHead, path, h, opts...)
}
func (r *Router) OPTIONS(path string, h fiberlib.Handler, opts ...HandlerOption) {
	r.Handle(http.MethodOptions, path, h, opts...)
}

func (r *Router) Routes() []openapi.RouteMeta { return r.routes }

func Register(r *Router, cfg openapi.Config) {
	doc := openapi.BuildSpec(r.routes, cfg)

	specPath := cfg.SpecPath
	if specPath == "" {
		specPath = "/openapi.json"
	}
	mount := cfg.SwaggerPath
	if mount == "" {
		mount = "/swagger-ui"
	}
	mount = strings.TrimSuffix(mount, "/")
	indexPath := mount + "/index.html"

	r.App.Get(specPath, func(c *fiberlib.Ctx) error {
		return c.Status(200).JSON(doc)
	})

	redirect := func(c *fiberlib.Ctx) error {
		return c.Redirect(indexPath+"#/", http.StatusFound)
	}

	r.App.Get(mount, redirect)
	r.App.Get(mount+"/", redirect)
	r.App.Get(indexPath, func(c *fiberlib.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.Status(200).SendString(`<!DOCTYPE html>
<html>
<head>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script>
SwaggerUIBundle({
  url: '` + specPath + `',
  dom_id: '#swagger-ui'
});
</script>
</body>
</html>`)
	})

	// Legacy /swagger redirect
	r.App.Get("/swagger", redirect)
	r.App.Get("/swagger/", redirect)
}

func Bind(c *fiberlib.Ctx, v interface{}) error           { return c.BodyParser(v) }
func JSON(c *fiberlib.Ctx, code int, v interface{}) error { return c.Status(code).JSON(v) }

type SecurityRequirement = openapi3.SecurityRequirement

// NOTE: Typed (generic) handler helpers were removed to keep the API simple.

// Group allows applying shared options (e.g., tags/security) and a common path prefix.
type Group struct {
	prefix string
	opts   []HandlerOption
	r      *Router
}

func (r *Router) Group(prefix string, opts ...HandlerOption) *Group {
	return &Group{prefix: prefix, opts: opts, r: r}
}

func (g *Group) Handle(method, p string, h fiberlib.Handler, opts ...HandlerOption) {
	all := make([]HandlerOption, 0, len(g.opts)+len(opts))
	all = append(all, g.opts...)
	all = append(all, opts...)
	g.r.Handle(method, joinPaths(g.prefix, p), h, all...)
}

func (g *Group) GET(p string, h fiberlib.Handler, opts ...HandlerOption) {
	g.Handle(http.MethodGet, p, h, opts...)
}
func (g *Group) POST(p string, h fiberlib.Handler, opts ...HandlerOption) {
	g.Handle(http.MethodPost, p, h, opts...)
}
func (g *Group) PUT(p string, h fiberlib.Handler, opts ...HandlerOption) {
	g.Handle(http.MethodPut, p, h, opts...)
}
func (g *Group) PATCH(p string, h fiberlib.Handler, opts ...HandlerOption) {
	g.Handle(http.MethodPatch, p, h, opts...)
}
func (g *Group) DELETE(p string, h fiberlib.Handler, opts ...HandlerOption) {
	g.Handle(http.MethodDelete, p, h, opts...)
}

func joinPaths(prefix, p string) string {
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
