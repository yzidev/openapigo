//go:build gin

package oas

import (
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"
	"github.com/yzidev/goas/adapters/ginadapter"
	"github.com/yzidev/goas/openapi"
)

type errSchemaDefaults struct {
	Error string `json:"error"`
}

func TestDefaultErrorResponsesIncludedEvenWithCustomSuccessStatuses(t *testing.T) {
	base := ginadapter.New()
	base.Engine = ginlib.New()

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")

	b := NewSpec()
	b.GroupTags("", []string{"T"}, func(s *SpecBuilder) {
		// custom success status will populate route.Responses
		// We still expect default errors to be added.
		s.POST("/secure/users").Security(&bearer).Res(struct{}{}).Created()
	})

	r := NewGinRouter(base, b.Spec())
	r.POST("/secure/users", func(c *ginlib.Context) {})

	doc := openapi.BuildSpec(r.Routes(), openapi.Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/secure/users")
	if p == nil || p.Post == nil {
		t.Fatalf("expected POST /secure/users")
	}

	// success
	if p.Post.Responses.Value("201") == nil {
		t.Fatalf("expected 201 response")
	}

	// defaults
	if p.Post.Responses.Value("400") == nil {
		t.Fatalf("expected default 400 response")
	}
	if p.Post.Responses.Value("401") == nil {
		t.Fatalf("expected default 401 response")
	}
	if p.Post.Responses.Value("500") == nil {
		t.Fatalf("expected default 500 response")
	}

	// Sanity: schema should have json content
	ref := p.Post.Responses.Value("500")
	if ref == nil || ref.Value == nil || ref.Value.Content == nil {
		t.Fatalf("expected 500 response content")
	}

	_ = http.StatusInternalServerError
}
