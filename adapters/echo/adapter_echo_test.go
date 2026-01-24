package echo

import (
	"testing"

	"github.com/aizacoders/openapigo/openapi"
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
	openapiCfg := openapi.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}
