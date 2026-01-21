//go:build gin

package gin

import (
	"net/http"
	"reflect"

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
)

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

// Register mounts /openapi.json and /swagger and uses captured routes.
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

	r.Engine.GET(specPath, func(c *ginlib.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, doc)
	})

	// Minimal swagger UI (same html as openapi/ui)
	r.Engine.GET(swagPath, func(c *ginlib.Context) {
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
}

// Helpers for gin
func Bind(c *ginlib.Context, v interface{}) error     { return c.ShouldBindJSON(v) }
func JSON(c *ginlib.Context, code int, v interface{}) { c.JSON(code, v) }

// Security helper alias.
type SecurityRequirement = openapi3.SecurityRequirement

type TypedHandler[TReq any, TRes any] func(c *ginlib.Context, req TReq) (res TRes, status int, err error)

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

func wrapTyped[TReq any, TRes any](h TypedHandler[TReq, TRes]) ginlib.HandlerFunc {
	return func(c *ginlib.Context) {
		var reqVal TReq
		if !isZeroStructType[TReq]() {
			_ = Bind(c, &reqVal)
		}

		res, code, err := h(c, reqVal)
		if err != nil {
			JSON(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if code == 0 {
			code = http.StatusOK
		}
		if isZeroStructType[TRes]() {
			c.Status(code)
			return
		}
		JSON(c, code, res)
	}
}

// GETT registers a typed GET handler with full-auto schema.
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

// POSTT registers a typed POST handler with full-auto schema.
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
