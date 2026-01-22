package openapi

import (
	"net/http"
	"testing"
)

type hdrRes struct {
	OK bool `json:"ok"`
}

func TestHeaderParamsAppearInOperation(t *testing.T) {
	routes := []RouteMeta{
		{
			Method: http.MethodGet,
			Path:   "/demo",
			HeaderParams: []HeaderParam{
				{Name: "X-Demo-Fail", Type: ParamString, Required: false, Description: "demo"},
			},
			ResponseSchema: hdrRes{},
		},
	}
	doc := BuildSpec(routes, Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/demo")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET /demo")
	}
	found := false
	for _, prm := range p.Get.Parameters {
		if prm.Value != nil && prm.Value.In == "header" && prm.Value.Name == "X-Demo-Fail" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected header param X-Demo-Fail")
	}
}
