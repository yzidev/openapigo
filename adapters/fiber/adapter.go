//go:build fiber

package fiber

import (
	"net/http"
	"reflect"

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
	swagPath := cfg.SwaggerPath
	if swagPath == "" {
		swagPath = "/swagger"
	}

	r.App.Get(specPath, func(c *fiberlib.Ctx) error {
		return c.Status(200).JSON(doc)
	})

	r.App.Get(swagPath, func(c *fiberlib.Ctx) error {
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
}

func Bind(c *fiberlib.Ctx, v interface{}) error           { return c.BodyParser(v) }
func JSON(c *fiberlib.Ctx, code int, v interface{}) error { return c.Status(code).JSON(v) }

type SecurityRequirement = openapi3.SecurityRequirement

// Typed handler support (full-auto schema)
type TypedHandler[TReq any, TRes any] func(c *fiberlib.Ctx, req TReq) (res TRes, status int, err error)

func isZeroStructType[T any]() bool {
	var zero T
	t := reflect.TypeOf(zero)
	return t != nil && t.Kind() == reflect.Struct && t.NumField() == 0
}

func typedOptions[TReq any, TRes any]() (reqOpt, resOpt HandlerOption) {
	var reqZero TReq
	var resZero TRes
	if !isZeroStructType[TReq]() {
		reqOpt = WithRequestSchema(reqZero)
	}
	if !isZeroStructType[TRes]() {
		resOpt = WithResponseSchema(resZero)
	}
	return reqOpt, resOpt
}

func mergeOpts(base []HandlerOption, add ...HandlerOption) []HandlerOption {
	out := make([]HandlerOption, 0, len(base)+len(add))
	out = append(out, base...)
	out = append(out, add...)
	return out
}

func wrapTyped[TReq any, TRes any](h TypedHandler[TReq, TRes]) fiberlib.Handler {
	return func(c *fiberlib.Ctx) error {
		var reqVal TReq
		if !isZeroStructType[TReq]() {
			_ = Bind(c, &reqVal)
		}

		res, code, err := h(c, reqVal)
		if err != nil {
			return JSON(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if code == 0 {
			code = http.StatusOK
		}
		if isZeroStructType[TRes]() {
			return c.SendStatus(code)
		}
		return JSON(c, code, res)
	}
}

func GETT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.Handle(http.MethodGet, path, wrapTyped(h), mergeOpts(base, opts...)...)
}

func POSTT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.Handle(http.MethodPost, path, wrapTyped(h), mergeOpts(base, opts...)...)
}

func PUTT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.Handle(http.MethodPut, path, wrapTyped(h), mergeOpts(base, opts...)...)
}

func PATCHT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.Handle(http.MethodPatch, path, wrapTyped(h), mergeOpts(base, opts...)...)
}

func DELETET[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.Handle(http.MethodDelete, path, wrapTyped(h), mergeOpts(base, opts...)...)
}
