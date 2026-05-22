package fiberadapter

import (
	"net/http"
	"strings"

	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yzidev/goas"
	"github.com/yzidev/goas/ui"
)

type Router struct {
	App    *fiberlib.App
	routes []goas.RouteMeta
}

func New(app ...*fiberlib.App) *Router {
	if len(app) > 0 {
		return Wrap(app[0])
	}
	return &Router{App: fiberlib.New()}
}

// NewFiberAdapters wraps an existing *fiber.App into the adapter Router so callers
// who create their own app (e.g., fiber.New()) can still use the adapter.
func NewFiberAdapters(app *fiberlib.App) *Router {
	return Wrap(app)
}

// Wrap converts an existing Fiber app into an Goas router adapter.
func Wrap(app *fiberlib.App) *Router {
	if app == nil {
		app = fiberlib.New()
	}
	return &Router{App: app}
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

func (r *Router) Handle(method, path string, h fiberlib.Handler, opts ...HandlerOption) {
	meta := goas.RouteMeta{Method: method, Path: path}
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

func (r *Router) Routes() []goas.RouteMeta { return r.routes }

// Docs mounts the generated OpenAPI JSON document and Swagger UI.
func (r *Router) Docs(cfg goas.Config) {
	Register(r, cfg)
}

// Docs mounts OpenAPI JSON and Swagger UI for a native Fiber app.
// It discovers routes registered directly on Fiber, so you can use plain Fiber
// routing and add Goas with a single call.
func Docs(app *fiberlib.App, cfg goas.Config) {
	Wrap(app).Docs(cfg)
}

// AutoDocs is an alias for Docs.
func AutoDocs(app *fiberlib.App, cfg goas.Config) {
	Docs(app, cfg)
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

	r.App.Get(specPath, func(c *fiberlib.Ctx) error {
		doc := goas.BuildSpec(r.discoveredRoutes(specPath, mount, indexPath), cfg)
		return c.Status(200).JSON(doc)
	})

	redirect := func(c *fiberlib.Ctx) error {
		return c.Redirect(indexPath+"#/", http.StatusFound)
	}

	r.App.Get(mount, redirect)
	r.App.Get(mount+"/", redirect)
	r.App.Get(indexPath, func(c *fiberlib.Ctx) error {
		c.Set("Content-Type", "text/html")
		ui.WriteSwaggerUIHTML(c.Context().Response.BodyWriter(), ui.SwaggerUIConfig{SpecURLPath: specPath})
		return nil
	})

	// Legacy /swagger redirect
	r.App.Get("/swagger", redirect)
	r.App.Get("/swagger/", redirect)
}

func (r *Router) discoveredRoutes(specPath, mount, indexPath string) []goas.RouteMeta {
	routes := append([]goas.RouteMeta(nil), r.routes...)
	seen := map[string]bool{}
	for _, route := range routes {
		seen[route.Method+" "+route.Path] = true
	}
	for _, route := range r.App.GetRoutes(true) {
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

func Bind(c *fiberlib.Ctx, v interface{}) error           { return c.BodyParser(v) }
func JSON(c *fiberlib.Ctx, code int, v interface{}) error { return c.Status(code).JSON(v) }

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
