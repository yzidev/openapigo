package mux

import (
	"testing"

	"github.com/aizacoders/openapigo/openapi"
)

func TestHTTPRouterNew(t *testing.T) {
	r := NewHttpAdapters()
	if r == nil {
		t.Fatalf("New() returned nil")
	}
	openapiCfg := openapi.Config{Title: "smoke", Version: "0"}
	Register(r, openapiCfg)
}
