package fiberadapter

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"
	"github.com/yzidev/goas/openapi"
)

func TestFiberNewAndWrap(t *testing.T) {
	r := New()
	if r == nil || r.App == nil {
		t.Fatalf("New() returned nil")
	}
	r2 := NewFiberAdapters(nil)
	if r2 == nil || r2.App == nil {
		t.Fatalf("NewFiberAdapters(nil) returned nil")
	}
	openapiCfg := openapi.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}

func TestFiberDocsDiscoversNativeRoutes(t *testing.T) {
	app := fiberlib.New()
	app.Get("/native/users/:id", func(c *fiberlib.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	Docs(app, openapi.Config{Title: "native", Version: "1"})

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	var doc openapi3.T
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatalf("invalid OpenAPI JSON: %v", err)
	}
	if p := doc.Paths.Find("/native/users/{id}"); p == nil || p.Get == nil {
		t.Fatalf("expected discovered GET /native/users/{id}")
	}
}
