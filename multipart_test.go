package goas

import (
	"net/http"
	"testing"
)

type uploadReq struct {
	File MultipartFile `json:"file"`
	Note string        `json:"note"`
}

func TestMultipartRequestBodyUsesFormData(t *testing.T) {
	r := NewRouter()
	r.POST("/upload", func(w http.ResponseWriter, r *http.Request) {}, WithRequestSchema(uploadReq{}))

	doc := BuildSpec(r.Routes(), Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/upload")
	if p == nil || p.Post == nil || p.Post.RequestBody == nil || p.Post.RequestBody.Value == nil {
		t.Fatalf("expected requestBody")
	}
	if _, ok := p.Post.RequestBody.Value.Content["multipart/form-data"]; !ok {
		t.Fatalf("expected multipart/form-data content")
	}
}
