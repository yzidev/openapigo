//go:build gin

package gin

import (
	"net/http"
	"path"
	"strings"

	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/getkin/kin-openapi/openapi3"
)

// Router wraps gin.Engine and captures route metadata for OpenAPI generation.
//
// This adapter is intentionally minimal: it captures method/path and allows you
// to provide request/response schema samples via options.
type Router struct {
	Engine *ginlib.Engine
	routes []openapi.RouteMeta
}

func New() *Router {
	return &Router{Engine: ginlib.New()}
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

// Group allows applying shared options (e.g., WithTags) and a common path prefix
// to multiple routes when using the Gin adapter.
//
// Example:
//
//	api := r.Group("", WithTags("Users"))
//	api.GET("/users", ...)
//
// Options provided to the group are applied to every route in that group.
type Group struct {
	prefix string
	opts   []HandlerOption
	route  func(method, path string, h ginlib.HandlerFunc, opts ...HandlerOption)
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
	j := path.Join(g.prefix, p)
	if strings.HasSuffix(p, "/") && !strings.HasSuffix(j, "/") {
		j += "/"
	}
	if !strings.HasPrefix(j, "/") {
		j = "/" + j
	}
	return j
}

func (g *Group) Handle(method, p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	all := make([]HandlerOption, 0, len(g.opts)+len(opts))
	all = append(all, g.opts...)
	all = append(all, opts...)
	g.route(method, g.join(p), h, all...)
}

func (g *Group) GET(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodGet, p, h, opts...)
}
func (g *Group) POST(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPost, p, h, opts...)
}
func (g *Group) PUT(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPut, p, h, opts...)
}
func (g *Group) DELETE(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodDelete, p, h, opts...)
}
func (g *Group) PATCH(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPatch, p, h, opts...)
}
func (g *Group) HEAD(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodHead, p, h, opts...)
}
func (g *Group) OPTIONS(p string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodOptions, p, h, opts...)
}

func (r *Router) Handle(method, path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	meta := openapi.RouteMeta{Method: method, Path: path}
	for _, opt := range opts {
		opt(&meta)
	}
	r.routes = append(r.routes, meta)

	r.Engine.Handle(method, path, h)
}

func (r *Router) GET(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *Router) POST(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *Router) PUT(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *Router) DELETE(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
func (r *Router) PATCH(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *Router) HEAD(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodHead, path, h, opts...)
}
func (r *Router) OPTIONS(path string, h ginlib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodOptions, path, h, opts...)
}

func (r *Router) Routes() []openapi.RouteMeta { return r.routes }

// Register mounts /openapi.json and Swagger UI and uses captured routes.
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

	r.Engine.GET(specPath, func(c *ginlib.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, doc)
	})

	redirect := func(c *ginlib.Context) {
		c.Redirect(http.StatusFound, indexPath+"#/")
	}

	// Canonical swagger UI paths
	r.Engine.GET(mount, redirect)
	r.Engine.GET(mount+"/", redirect)
	r.Engine.GET(indexPath, func(c *ginlib.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `<!DOCTYPE html>
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
  url: '`+specPath+`',
  dom_id: '#swagger-ui'
});
</script>
</body>
</html>`)
	})

	// Legacy /swagger redirect
	r.Engine.GET("/swagger", redirect)
	r.Engine.GET("/swagger/", redirect)
}

// Helpers for gin
func Bind(c *ginlib.Context, v interface{}) error     { return c.ShouldBindJSON(v) }
func JSON(c *ginlib.Context, code int, v interface{}) { c.JSON(code, v) }

// Security helper alias.
type SecurityRequirement = openapi3.SecurityRequirement

// NOTE: Typed (generic) handler helpers were removed to keep the API simple.
