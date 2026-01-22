//go:build gin

package simple

import (
	"net/http"
	"testing"

	ginadapter "github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	ginlib "github.com/gin-gonic/gin"
)

type testReq struct {
	Name string `json:"name"`
}

type testRes struct {
	ID string `json:"id"`
}

func TestGinGroupKeepsSpecInjection(t *testing.T) {
	base := ginadapter.New()
	g := ginlib.New()
	base.Engine = g

	spec := Spec{
		Key(http.MethodPost, "/users"): {
			Tags:      []string{"Users"},
			ReqSchema: testReq{},
			ResSchema: testRes{},
			Status:    http.StatusCreated,
		},
	}

	r := NewGin(base, spec)
	grp := r.Group("", ginadapter.WithTags("Users"))
	grp.POST("/users", func(c *ginlib.Context) {})

	routes := r.Routes()
	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}
	if routes[0].RequestSchema == nil || routes[0].ResponseSchema == nil {
		t.Fatalf("expected request+response schema to be injected, got req=%T res=%T", routes[0].RequestSchema, routes[0].ResponseSchema)
	}

	doc := openapi.BuildSpec(routes, openapi.Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/users")
	if p == nil || p.Post == nil {
		t.Fatalf("expected POST /users")
	}
	if p.Post.RequestBody == nil {
		t.Fatalf("expected requestBody in spec")
	}
	if p.Post.Responses == nil || p.Post.Responses.Value("201") == nil {
		t.Fatalf("expected 201 response")
	}
	if doc.Components == nil || len(doc.Components.Schemas) == 0 {
		t.Fatalf("expected component schemas to be generated")
	}
}
