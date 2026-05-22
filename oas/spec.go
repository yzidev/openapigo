package oas

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yzidev/goas"
)

// RouteDef is a config-only description of an endpoint.
//
// This is the "SpringBoot-like" mode: keep handlers/routes clean, put OpenAPI
// request/response schema + tags/security/query definition in one config place.
//
// Note: request/response schema inference is not possible from plain Go handlers,
// so you still need to declare them here (or omit them).
//
// If ReqSchema/ResSchema is nil, schema isn't declared.
type RouteDef struct {
	Tags         []string
	Security     *openapi3.SecurityRequirement
	QueryParams  []goas.QueryParam
	HeaderParams []goas.HeaderParam

	ReqSchema any
	ResSchema any
	Status    int

	// Optional extra responses (errors, alternate status codes)
	Responses []goas.ResponseSpec
}

// Spec maps method+path to its OpenAPI definition.
// Key format: "METHOD /path".
type Spec map[string]RouteDef

// Key builds the Spec key.
func Key(method, path string) string {
	return strings.ToUpper(method) + " " + path
}

// Inject converts a RouteDef into route options.
func Inject(opts []goas.HandlerOption, def RouteDef) []goas.HandlerOption {
	out := make([]goas.HandlerOption, 0, len(opts)+8)
	out = append(out, opts...)
	if len(def.Tags) > 0 {
		out = append(out, goas.WithTags(def.Tags...))
	}
	if def.Security != nil {
		out = append(out, goas.WithSecurity(def.Security))
	}
	if len(def.QueryParams) > 0 {
		out = append(out, goas.WithQueryParams(def.QueryParams...))
	}
	if len(def.HeaderParams) > 0 {
		out = append(out, goas.WithHeaderParams(def.HeaderParams...))
	}
	if def.ReqSchema != nil || def.ResSchema != nil || def.Status != 0 {
		if def.ReqSchema != nil {
			out = append(out, goas.WithRequestSchema(def.ReqSchema))
		}
		if def.ResSchema != nil {
			out = append(out, goas.WithResponseSchema(def.ResSchema))
		}
		if def.ResSchema != nil || def.Status != 0 {
			status := def.Status
			if status == 0 {
				status = 200
			}
			out = append(out, goas.WithResponses(goas.ResponseSpec{Status: status, Schema: def.ResSchema}))
		}
	}
	if len(def.Responses) > 0 {
		out = append(out, goas.WithResponses(def.Responses...))
	}
	return out
}
