package fiber

import (
	"testing"

	"github.com/aizacoders/openapigo/openapi"
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
