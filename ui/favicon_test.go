package ui

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testMux struct{ m map[string]http.HandlerFunc }

func (t *testMux) Get(path string, h http.HandlerFunc) {
	if t.m == nil {
		t.m = map[string]http.HandlerFunc{}
	}
	t.m[path] = h
}

func TestFaviconServed(t *testing.T) {
	mux := &testMux{}
	RegisterSwaggerUI(mux, SwaggerUIConfig{MountPath: "/swagger-ui"})

	h, ok := mux.m["/favicon.ico"]
	if !ok {
		t.Fatalf("expected /favicon.ico route")
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	h(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected image/png, got %q", ct)
	}
	b := rr.Body.Bytes()
	// PNG files start with 89 50 4E 47
	if len(b) < 4 {
		t.Fatalf("expected non-empty PNG body, got len=%d", len(b))
	}
	if b[0] != 0x89 || b[1] != 0x50 || b[2] != 0x4e || b[3] != 0x47 {
		t.Fatalf("expected PNG signature, got %v", b[:4])
	}
}
