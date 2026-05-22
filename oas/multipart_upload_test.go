package oas

import (
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yzidev/goas"
)

func TestMultipartUploadHelperProducesMultipartFormData(t *testing.T) {
	base := goas.NewRouter()
	b := NewSpec()
	b.Group("", func(s *SpecBuilder) {
		s.POST("/upload").MultipartUpload("file", goas.MultipartField{Name: "note", Type: goas.ParamString}).Res(map[string]string{}).OK()
	})

	r := NewHttpRouter(base, b.Spec())
	r.POST("/upload", func(w http.ResponseWriter, r *http.Request) {})

	doc := goas.BuildSpec(base.Routes(), goas.Config{Title: "T", Version: "1"})
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
