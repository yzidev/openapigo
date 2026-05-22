package echoadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"
	"github.com/yzidev/goas"
)

func TestEchoNewAndWrap(t *testing.T) {
	r := New()
	if r == nil || r.Echo == nil {
		t.Fatalf("New() returned nil")
	}
	r2 := NewEchoAdapters(nil)
	if r2 == nil || r2.Echo == nil {
		t.Fatalf("NewEchoAdapters(nil) returned nil")
	}
	openapiCfg := goas.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}

func TestEchoDocsDiscoversNativeRoutes(t *testing.T) {
	e := echolib.New()
	e.GET("/native/users/:id", func(c echolib.Context) error {
		return c.NoContent(http.StatusOK)
	})

	Docs(e, goas.Config{Title: "native", Version: "1"})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	e.ServeHTTP(rec, req)

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
