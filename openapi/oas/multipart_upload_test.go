package oas

import (
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yzidev/openapigo/openapi"
)

func TestMultipartUploadHelperProducesMultipartFormData(t *testing.T) {
	base := openapi.NewRouter()
	b := NewSpec()
	b.Group("", func(s *SpecBuilder) {
		s.POST("/upload").MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()
	})

	r := NewHttpRouter(base, b.Spec())
	r.POST("/upload", func(w http.ResponseWriter, r *http.Request) {})

	doc := openapi.BuildSpec(base.Routes(), openapi.Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/upload")
	if p == nil || p.Post == nil || p.Post.RequestBody == nil || p.Post.RequestBody.Value == nil {
		t.Fatalf("expected requestBody")
	}
	if _, ok := p.Post.RequestBody.Value.Content["multipart/form-data"]; !ok {
		t.Fatalf("expected multipart/form-data content, got keys=%v", keys(p.Post.RequestBody.Value.Content))
	}
}

func keys(m map[string]*openapi3.MediaType) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
