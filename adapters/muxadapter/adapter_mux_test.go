package muxadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yzidev/goas/openapi"
)

func TestHTTPRouterNew(t *testing.T) {
	r := NewHttpAdapters()
	if r == nil {
		t.Fatalf("New() returned nil")
	}
	openapiCfg := openapi.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}

func TestMuxMountAllowsRoutesAfterDocs(t *testing.T) {
	mux := http.NewServeMux()
	r := Mount(mux, openapi.Config{Title: "native", Version: "1"})
	r.GET("/users/{id}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, Res(map[string]string{}), Tags("Users"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var doc openapi3.T
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid OpenAPI JSON: %v", err)
	}
	p := doc.Paths.Find("/users/{id}")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET /users/{id}")
	}
	if len(p.Get.Tags) != 1 || p.Get.Tags[0] != "Users" {
		t.Fatalf("expected Users tag, got %#v", p.Get.Tags)
	}
}
