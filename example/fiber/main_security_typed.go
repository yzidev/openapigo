//go:build fiber && typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiber"
	"github.com/aizacoders/openapigo/openapi"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SecCreateUser struct {
	Name string `json:"name"`
}

func main() {
	r := fiber.New()

	cfg := openapi.Config{
		Title:   "User API (Fiber + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}
	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	secureOpts := []fiber.HandlerOption{fiber.WithTags("Secure Users")}

	postOpts := openapi.MergeOptionSlices(
		secureOpts,
		[]fiber.HandlerOption{fiber.WithSecurity(&bearer)},
		fiber.JSONRoute(SecCreateUser{}, SecUser{}, http.StatusCreated),
	)
	fiber.POSTT[SecCreateUser, SecUser](r, "/secure/users", func(c *fiberlib.Ctx, in SecCreateUser) (SecUser, int, error) {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return SecUser{}, http.StatusUnauthorized, nil
		}
		return SecUser{ID: "1", Name: in.Name}, http.StatusCreated, nil
	}, postOpts...)

	getOpts := openapi.MergeOptionSlices(
		secureOpts,
		[]fiber.HandlerOption{fiber.WithSecurity(&apiKey)},
		fiber.JSONRoute(struct{}{}, []SecUser{}, http.StatusOK),
	)
	fiber.GETT[struct{}, []SecUser](r, "/secure/users", func(c *fiberlib.Ctx, _ struct{}) ([]SecUser, int, error) {
		if c.Get("X-API-Key") == "" {
			return nil, http.StatusUnauthorized, nil
		}
		return []SecUser{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	}, getOpts...)

	fiber.Register(r, cfg)
	_ = r.App.Listen(":8080")
}
