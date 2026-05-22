package ui

import (
	_ "embed"
	"net/http"
	"strings"
	"text/template"
)

//go:embed templates/swagger-ui.html
var swaggerUITemplate string

//go:embed templates/favicon.ico
var openAPIFaviconPNG []byte

var swaggerUITpl = template.Must(template.New("swagger-ui.html").Parse(swaggerUITemplate))

type SwaggerUIConfig struct {
	MountPath   string // default: /swagger-ui
	SpecURLPath string // default: /openapi.json
	Version     string // optional: used for cache busting assets like favicon
}

func RegisterSwaggerUI(mux interface {
	Get(string, http.HandlerFunc)
}, cfg SwaggerUIConfig) {
	mount := cfg.MountPath
	if mount == "" {
		mount = "/swagger-ui"
	}
	if !strings.HasPrefix(mount, "/") {
		mount = "/" + mount
	}
	mount = strings.TrimSuffix(mount, "/")

	// Standard favicon.ico (most reliable across browsers)
	mux.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openAPIFaviconPNG)
	})
	// Also serve under mount path for setups that expect it there.
	mux.Get(mount+"/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openAPIFaviconPNG)
	})

	spec := cfg.SpecURLPath
	if spec == "" {
		spec = "/openapi.json"
	}

	indexPath := mount + "/index.html"
	redirectHTML := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, indexPath+"#/", http.StatusFound)
	}

	ver := cfg.Version
	if ver == "" {
		ver = "1"
	}

	// New canonical paths
	mux.Get(mount, redirectHTML)
	mux.Get(mount+"/", redirectHTML)
	mux.Get(indexPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_ = swaggerUITpl.Execute(w, map[string]any{"SpecURL": spec, "MountPath": mount, "Version": ver})
	})

	// Legacy: /swagger should redirect to new canonical UI.
	if mount != "/swagger" {
		mux.Get("/swagger", redirectHTML)
		mux.Get("/swagger/", redirectHTML)
	}
}
