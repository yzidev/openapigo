package ui

import (
	"io"
	"strings"
)

// WriteSwaggerUIHTML renders the embedded Swagger UI template into the provided writer.
//
// It intentionally does not register any routes; use RegisterSwaggerUI for that.
func WriteSwaggerUIHTML(w io.Writer, cfg SwaggerUIConfig) {
	spec := cfg.SpecURLPath
	if spec == "" {
		spec = "/openapi.json"
	}
	mount := cfg.MountPath
	if mount == "" {
		mount = "/swagger-ui"
	}
	if !strings.HasPrefix(mount, "/") {
		mount = "/" + mount
	}
	mount = strings.TrimSuffix(mount, "/")
	ver := cfg.Version
	if ver == "" {
		ver = "1"
	}

	_ = swaggerUITpl.Execute(w, map[string]any{"SpecURL": spec, "MountPath": mount, "Version": ver})
}
