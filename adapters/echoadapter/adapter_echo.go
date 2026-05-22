package echoadapter

import (
	"net/http"
	"strings"

	echolib "github.com/labstack/echo/v4"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yzidev/goas"
	"github.com/yzidev/goas/ui"
)

type Router struct {
	Echo   *echolib.Echo
	routes []goas.RouteMeta
}

func New(e ...*echolib.Echo) *Router {
	if len(e) > 0 {
		return Wrap(e[0])
	}
	return &Router{Echo: echolib.New()}
}

// NewEchoAdapters wraps an existing *echo.Echo into the adapter Router so callers
// who create their own echo server (e.g., echo.New() or echo.Default()) can use the adapter.
func NewEchoAdapters(e *echolib.Echo) *Router {
	return Wrap(e)
}

// Wrap converts an existing Echo instance into an Goas router adapter.
func Wrap(e *echolib.Echo) *Router {
	if e == nil {
		e = echolib.New()
	}
	return &Router{Echo: e}
}

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

func (r *Router) Handle(method, path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	meta := goas.RouteMeta{Method: method, Path: path}
	for _, opt := range opts {
		opt(&meta)
	}
	r.routes = append(r.routes, meta)

	r.Echo.Add(method, path, h)
}

func (r *Router) GET(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *Router) POST(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *Router) PUT(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *Router) DELETE(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
func (r *Router) PATCH(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *Router) HEAD(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodHead, path, h, opts...)
}
func (r *Router) OPTIONS(path string, h echolib.HandlerFunc, opts ...HandlerOption) {
	r.Handle(http.MethodOptions, path, h, opts...)
}

func (r *Router) Routes() []goas.RouteMeta { return r.routes }

// Docs mounts the generated OpenAPI JSON document and Swagger UI.
func (r *Router) Docs(cfg goas.Config) {
	Register(r, cfg)
}

// Docs mounts OpenAPI JSON and Swagger UI for a native Echo instance.
// It discovers routes registered directly on Echo, so you can use plain Echo
// routing and add Goas with a single call.
func Docs(e *echolib.Echo, cfg goas.Config) {
	Wrap(e).Docs(cfg)
}

// AutoDocs is an alias for Docs.
func AutoDocs(e *echolib.Echo, cfg goas.Config) {
	Docs(e, cfg)
}

func Register(r *Router, cfg goas.Config) {
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

	r.Echo.GET(specPath, func(c echolib.Context) error {
		doc := goas.BuildSpec(r.discoveredRoutes(specPath, mount, indexPath), cfg)
		return c.JSON(200, doc)
	})

	redirect := func(c echolib.Context) error {
		return c.Redirect(http.StatusFound, indexPath+"#/")
	}

	r.Echo.GET(mount, redirect)
	r.Echo.GET(mount+"/", redirect)
	r.Echo.GET(indexPath, func(c echolib.Context) error {
		c.Response().Header().Set("Content-Type", "text/html")
		ui.WriteSwaggerUIHTML(c.Response().Writer, ui.SwaggerUIConfig{SpecURLPath: specPath})
		return nil
	})

	// Legacy /swagger redirect
	r.Echo.GET("/swagger", redirect)
	r.Echo.GET("/swagger/", redirect)
}

func (r *Router) discoveredRoutes(specPath, mount, indexPath string) []goas.RouteMeta {
	routes := append([]goas.RouteMeta(nil), r.routes...)
	seen := map[string]bool{}
	for _, route := range routes {
		seen[route.Method+" "+route.Path] = true
	}
	for _, route := range r.Echo.Routes() {
		if skipDocsRoute(route.Path, specPath, mount, indexPath) {
			continue
		}
		key := route.Method + " " + route.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		routes = append(routes, goas.RouteMeta{Method: route.Method, Path: route.Path})
	}
	return routes
}

func skipDocsRoute(routePath, specPath, mount, indexPath string) bool {
	switch routePath {
	case specPath, mount, mount + "/", indexPath, "/swagger", "/swagger/":
		return true
	default:
		return false
	}
}

func Bind(c echolib.Context, v interface{}) error           { return c.Bind(v) }
func JSON(c echolib.Context, code int, v interface{}) error { return c.JSON(code, v) }

type SecurityRequirement = openapi3.SecurityRequirement

// NOTE: Typed (generic) handler helpers were removed to keep the API oas.

// Group allows applying shared options (e.g., tags/security) and a common path prefix.
type Group struct {
	prefix string
	opts   []HandlerOption
	r      *Router
}

func (r *Router) Group(prefix string, opts ...HandlerOption) *Group {
	return &Group{prefix: prefix, opts: opts, r: r}
}

func (g *Group) Handle(method, p string, h echolib.HandlerFunc, opts ...HandlerOption) {
	all := make([]HandlerOption, 0, len(g.opts)+len(opts))
	all = append(all, g.opts...)
	all = append(all, opts...)
	g.r.Handle(method, joinPaths(g.prefix, p), h, all...)
}

func (g *Group) GET(p string, h echolib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodGet, p, h, opts...)
}
func (g *Group) POST(p string, h echolib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPost, p, h, opts...)
}
func (g *Group) PUT(p string, h echolib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPut, p, h, opts...)
}
func (g *Group) PATCH(p string, h echolib.HandlerFunc, opts ...HandlerOption) {
	g.Handle(http.MethodPatch, p, h, opts...)
}
func (g *Group) DELETE(p string, h echolib.HandlerFunc, opts ...HandlerOption) {
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
