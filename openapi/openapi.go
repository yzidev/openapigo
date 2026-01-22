package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/aizacoders/openapigo/openapi/ui"

	"github.com/getkin/kin-openapi/openapi3"
)

type contextKey int

const (
	pathParamsKey contextKey = iota
)

func PathValue(r *http.Request, key string) string {
	if p, ok := r.Context().Value(pathParamsKey).(map[string]string); ok {
		return p[key]
	}
	return ""
}

type Config struct {
	Title           string
	Version         string
	SecuritySchemes map[string]*openapi3.SecuritySchemeRef
	Tags            openapi3.Tags
	SpecPath        string
	SwaggerPath     string

	// Schemas registers component schemas by name without attaching them to a route.
	// Useful when you want config-only schema registration.
	Schemas SchemaRegistry
}

func Register(r *Router, cfg Config) {
	if cfg.SpecPath == "" {
		cfg.SpecPath = "/openapi.json"
	}
	if cfg.SwaggerPath == "" {
		cfg.SwaggerPath = "/swagger-ui"
	}

	doc := BuildSpec(r.Routes(), cfg)

	r.Mux.Get(cfg.SpecPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(doc)
	})

	ui.RegisterSwaggerUI(r.Mux, ui.SwaggerUIConfig{MountPath: cfg.SwaggerPath, SpecURLPath: cfg.SpecPath})
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
