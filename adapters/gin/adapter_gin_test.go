package gin

import (
	"testing"

	"github.com/aizacoders/openapigo/openapi"
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
	openapiCfg := openapi.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}
