package simple

import (
	"net/http"
	"path"
	"strings"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/getkin/kin-openapi/openapi3"
)

// Builder helps you build a Spec without manually typing keys like "GET /users/:id".
//
// Example:
//
//	spec := simple.NewSpec().
//		Group("/", func(s *simple.SpecBuilder) {
//			s.GET("/users").Tags("Users").Res([]User{}).OK()
//			s.POST("/users").Tags("Users").Req(CreateUser{}).Res(struct{}{}).Status(http.StatusCreated)
//		})
//		Spec()
//
// NOTE: This builder does not try to infer schemas from handlers.
// You still declare Req/Res schema samples here.
type Builder struct {
	spec Spec
}

func NewSpec() *Builder {
	return &Builder{spec: make(Spec)}
}

// With preloads an existing spec (merged; def from b wins on conflict).
func (b *Builder) With(s Spec) *Builder {
	for k, v := range s {
		b.spec[k] = v
	}
	return b
}

// Group helps apply a path prefix to multiple routes.
func (b *Builder) Group(prefix string, fn func(s *SpecBuilder)) *Builder {
	sb := &SpecBuilder{b: b, prefix: prefix}
	fn(sb)
	return b
}

// GroupTags is like Group, but also applies default tags to every route defined in the group.
func (b *Builder) GroupTags(prefix string, tags []string, fn func(s *SpecBuilder)) *Builder {
	sb := &SpecBuilder{b: b, prefix: prefix, tags: append([]string(nil), tags...)}
	fn(sb)
	return b
}

// Spec returns the built Spec.
func (b *Builder) Spec() Spec { return b.spec }

// SpecBuilder is a scoped builder that can apply a prefix.
type SpecBuilder struct {
	b      *Builder
	prefix string
	tags   []string
}

// WithTags sets default tags for all routes built from this builder.
func (s *SpecBuilder) WithTags(tags ...string) *SpecBuilder {
	s.tags = append(s.tags, tags...)
	return s
}

func (s *SpecBuilder) join(p string) string {
	if s.prefix == "" {
		return p
	}
	if p == "" {
		return s.prefix
	}
	j := path.Join(s.prefix, p)
	// Keep trailing slash behavior similar to router group join.
	if strings.HasSuffix(p, "/") && !strings.HasSuffix(j, "/") {
		j += "/"
	}
	if !strings.HasPrefix(j, "/") {
		j = "/" + j
	}
	return j
}

func (s *SpecBuilder) route(method, p string) *RouteBuilder {
	full := s.join(p)
	k := Key(method, full)
	def := RouteDef{}
	if len(s.tags) > 0 {
		def.Tags = append(def.Tags, s.tags...)
	}
	return &RouteBuilder{sb: s, method: method, path: full, key: k, def: def}
}

func (s *SpecBuilder) GET(p string) *RouteBuilder    { return s.route(http.MethodGet, p) }
func (s *SpecBuilder) POST(p string) *RouteBuilder   { return s.route(http.MethodPost, p) }
func (s *SpecBuilder) PUT(p string) *RouteBuilder    { return s.route(http.MethodPut, p) }
func (s *SpecBuilder) PATCH(p string) *RouteBuilder  { return s.route(http.MethodPatch, p) }
func (s *SpecBuilder) DELETE(p string) *RouteBuilder { return s.route(http.MethodDelete, p) }

// RouteBuilder builds a single route definition.
type RouteBuilder struct {
	sb     *SpecBuilder
	method string
	path   string
	key    string
	def    RouteDef
}

func (r *RouteBuilder) Tags(tags ...string) *RouteBuilder {
	r.def.Tags = append(r.def.Tags, tags...)
	return r
}

func (r *RouteBuilder) Security(sec *openapi3.SecurityRequirement) *RouteBuilder {
	r.def.Security = sec
	return r
}

func (r *RouteBuilder) Query(params ...openapi.QueryParam) *RouteBuilder {
	r.def.QueryParams = append(r.def.QueryParams, params...)
	return r
}

func (r *RouteBuilder) Headers(params ...openapi.HeaderParam) *RouteBuilder {
	r.def.HeaderParams = append(r.def.HeaderParams, params...)
	return r
}

func (r *RouteBuilder) Req(schema any) *RouteBuilder {
	r.def.ReqSchema = schema
	return r
}

func (r *RouteBuilder) Res(schema any) *RouteBuilder {
	r.def.ResSchema = schema
	return r
}

func (r *RouteBuilder) Status(status int) *RouteBuilder {
	r.def.Status = status
	// commit
	r.sb.b.spec[r.key] = r.def
	return r
}

// OK is shorthand for Status(http.StatusOK).
func (r *RouteBuilder) OK() *RouteBuilder {
	return r.Status(http.StatusOK)
}

// Created is shorthand for Status(http.StatusCreated).
func (r *RouteBuilder) Created() *RouteBuilder {
	return r.Status(http.StatusCreated)
}

// NoContent is shorthand for Status(http.StatusNoContent).
func (r *RouteBuilder) NoContent() *RouteBuilder {
	return r.Status(http.StatusNoContent)
}

// Responses appends additional response specs (like 400/500 error shape).
func (r *RouteBuilder) Responses(specs ...openapi.ResponseSpec) *RouteBuilder {
	r.def.Responses = append(r.def.Responses, specs...)
	// commit
	r.sb.b.spec[r.key] = r.def
	return r
}

// Done commits the current definition without changing status.
// Use this if you only want tags/security/query and no schema.
func (r *RouteBuilder) Done() *RouteBuilder {
	r.sb.b.spec[r.key] = r.def
	return r
}
