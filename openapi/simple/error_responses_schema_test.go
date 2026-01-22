//go:build gin

package simple

import (
	"net/http"
	"testing"

	ginadapter "github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	ginlib "github.com/gin-gonic/gin"
)

type errSchema struct {
	Error string `json:"error"`
}

func TestErrorResponseSchemasAppearInComponents(t *testing.T) {
	base := ginadapter.New()
	base.Engine = ginlib.New()

	spec := Spec{
		Key(http.MethodGet, "/users/demo-errors"): {
			Tags:      []string{"Users"},
			ResSchema: map[string]string{},
			Status:    http.StatusOK,
			Responses: []openapi.ResponseSpec{
				{Status: 400, Schema: errSchema{}},
				{Status: 401, Schema: errSchema{}},
				{Status: 500, Schema: errSchema{}},
				{Status: 503, Schema: errSchema{}},
			},
		},
	}

	r := NewGin(base, spec)
	r.GET("/users/demo-errors", func(c *ginlib.Context) {})

	doc := openapi.BuildSpec(r.Routes(), openapi.Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/users/demo-errors")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET /users/demo-errors")
	}

	for _, code := range []string{"400", "401", "500", "503"} {
		if p.Get.Responses.Value(code) == nil {
			t.Fatalf("expected %s response", code)
		}
	}

	if doc.Components == nil || len(doc.Components.Schemas) == 0 {
		t.Fatalf("expected component schemas to be generated for error responses")
	}
}
