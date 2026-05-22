// Package spec provides the config-first OpenAPI route specification builder.
//
// It is the recommended import path for keeping handler registration clean while
// declaring request/response schemas, tags, security, and parameters in one place.
package spec

import (
	"github.com/yzidev/goas"
	"github.com/yzidev/goas/adapters/echoadapter"
	"github.com/yzidev/goas/adapters/fiberadapter"
	"github.com/yzidev/goas/adapters/ginadapter"
	"github.com/yzidev/goas/oas"
)

type (
	Builder      = oas.Builder
	SpecBuilder  = oas.SpecBuilder
	RouteBuilder = oas.RouteBuilder
	RouteDef     = oas.RouteDef
	Spec         = oas.Spec

	HTTPRouter  = oas.Router
	GinRouter   = oas.GinRouter
	EchoRouter  = oas.EchoRouter
	FiberRouter = oas.FiberRouter
)

func New() *Builder {
	return oas.NewSpec()
}

func NewSpec() *Builder {
	return oas.NewSpec()
}

func Key(method, path string) string {
	return oas.Key(method, path)
}

func Inject(opts []goas.HandlerOption, def RouteDef) []goas.HandlerOption {
	return oas.Inject(opts, def)
}

func HTTP(base *goas.Router, s Spec) *HTTPRouter {
	return oas.NewHttpRouter(base, s)
}

func NewHTTPRouter(base *goas.Router, s Spec) *HTTPRouter {
	return oas.NewHttpRouter(base, s)
}

func Gin(base *ginadapter.Router, s Spec) *GinRouter {
	return oas.NewGinRouter(base, s)
}

func NewGinRouter(base *ginadapter.Router, s Spec) *GinRouter {
	return oas.NewGinRouter(base, s)
}

func Echo(base *echoadapter.Router, s Spec) *EchoRouter {
	return oas.NewEchoRouter(base, s)
}

func NewEchoRouter(base *echoadapter.Router, s Spec) *EchoRouter {
	return oas.NewEchoRouter(base, s)
}

func Fiber(base *fiberadapter.Router, s Spec) *FiberRouter {
	return oas.NewFiberRouter(base, s)
}

func NewFiberRouter(base *fiberadapter.Router, s Spec) *FiberRouter {
	return oas.NewFiberRouter(base, s)
}
