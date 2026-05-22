package ginadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"
	"github.com/yzidev/goas"
)

func TestGinNewAndWrap(t *testing.T) {
	r := New()
	if r == nil || r.Engine == nil {
		t.Fatalf("New() returned nil")
	}
	// wrap nil engine -> should create a non-nil router
	r2 := NewGinAdapters(nil)
	if r2 == nil || r2.Engine == nil {
		t.Fatalf("NewGinOAS(nil) returned nil")
	}
	// Register should not panic with minimal config
	openapiCfg := goas.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}

func TestGinDocsDiscoversNativeRoutes(t *testing.T) {
	engine := ginlib.New()
	engine.GET("/native/users/:id", func(c *ginlib.Context) {
		c.Status(http.StatusOK)
	})

	Docs(engine, goas.Config{Title: "native", Version: "1"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var doc openapi3.T
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid OpenAPI JSON: %v", err)
	}
	if p := doc.Paths.Find("/native/users/{id}"); p == nil || p.Get == nil {
		t.Fatalf("expected discovered GET /native/users/{id}")
	}
}
