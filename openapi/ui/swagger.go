package ui

import (
	"net/http"
	"strings"
)

type SwaggerUIConfig struct {
	MountPath   string // default: /swagger-ui
	SpecURLPath string // default: /openapi.json
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

	spec := cfg.SpecURLPath
	if spec == "" {
		spec = "/openapi.json"
	}

	indexPath := mount + "/index.html"
	redirectHTML := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, indexPath+"#/", http.StatusFound)
	}

	// New canonical paths
	mux.Get(mount, redirectHTML)
	mux.Get(mount+"/", redirectHTML)
	mux.Get(indexPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
<!DOCTYPE html>
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
  url: '` + spec + `',
  dom_id: '#swagger-ui'
});
</script>
</body>
</html>
`))
	})

	// Legacy: /swagger should redirect to new canonical UI.
	if mount != "/swagger" {
		mux.Get("/swagger", redirectHTML)
		mux.Get("/swagger/", redirectHTML)
	}
}
