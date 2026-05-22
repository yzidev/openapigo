package goas

import (
	"encoding/json"
	"net/http"

	"github.com/yzidev/goas/ui"

	"github.com/getkin/kin-openapi/openapi3"
)

type contextKey int

const (
	pathParamsKey contextKey = iota
)

// Re-export common OpenAPI config types so applications can configure Goas
// without importing the underlying kin-openapi package directly.
type (
	SecurityRequirement  = openapi3.SecurityRequirement
	SecurityRequirements = openapi3.SecurityRequirements
	SecurityScheme       = openapi3.SecurityScheme
	SecuritySchemeRef    = openapi3.SecuritySchemeRef
	Tag                  = openapi3.Tag
	DocumentTags         = openapi3.Tags
)

func NewSecurityRequirement() SecurityRequirement {
	return openapi3.NewSecurityRequirement()
}

func PathValue(r *http.Request, key string) string {
	if p, ok := r.Context().Value(pathParamsKey).(map[string]string); ok {
		return p[key]
	}
	return ""
}

type Config struct {
	Title           string
	Version         string
	Description     string
	SecuritySchemes map[string]*SecuritySchemeRef
	Security        SecurityRequirements
	Tags            DocumentTags
	SpecPath        string
	SwaggerPath     string

	// Schemas registers component schemas by name without attaching them to a route.
	// Useful when you want config-only schema registration.
	Schemas SchemaRegistry

	// DefaultErrorResponses controls which standard error responses are automatically
	// added to every operation (if not already declared).
	//
	// If nil, a sensible default set is used.
	// If empty (len==0), automatic error responses are disabled.
	DefaultErrorResponses []int

	// DefaultErrorSchema is the schema used for DefaultErrorResponses.
	// If nil, goas.ErrorResponse{} is used.
	DefaultErrorSchema any
}

func Register(r *Router, cfg Config) {
	if cfg.SpecPath == "" {
		cfg.SpecPath = "/openapi.json"
	}
	if cfg.SwaggerPath == "" {
		cfg.SwaggerPath = "/swagger-ui"
	}

	// Serve OpenAPI JSON
	// Use Router.Get so we don't depend on underlying mux implementation details.
	r.Get(cfg.SpecPath, func(w http.ResponseWriter, _ *http.Request) {
		doc := BuildSpec(r.Routes(), cfg)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(doc)
	})

	// Serve Swagger UI + favicon on the Router (not the raw chi mux)
	ui.RegisterSwaggerUI(r, ui.SwaggerUIConfig{MountPath: cfg.SwaggerPath, SpecURLPath: cfg.SpecPath, Version: cfg.Version})
}

func ptr[T any](v T) *T {
	return &v
}

func Bind(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func JSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
