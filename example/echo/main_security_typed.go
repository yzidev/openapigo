//go:build echo && typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
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
	r := echo.New()

	cfg := openapi.Config{
		Title:   "User API (Echo + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}
	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	secureOpts := []echo.HandlerOption{echo.WithTags("Secure Users")}

	postOpts := append(append([]echo.HandlerOption{}, secureOpts...), echo.WithSecurity(&bearer))
	postOpts = append(postOpts, echo.JSONRoute(SecCreateUser{}, SecUser{}, http.StatusCreated)...)
	echo.POSTT[SecCreateUser, SecUser](r, "/secure/users", func(c echolib.Context, in SecCreateUser) (SecUser, int, error) {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return SecUser{}, http.StatusUnauthorized, nil
		}
		return SecUser{ID: "1", Name: in.Name}, http.StatusCreated, nil
	}, postOpts...)

	getOpts := append(append([]echo.HandlerOption{}, secureOpts...), echo.WithSecurity(&apiKey))
	getOpts = append(getOpts, echo.JSONRoute(struct{}{}, []SecUser{}, http.StatusOK)...)
	echo.GETT[struct{}, []SecUser](r, "/secure/users", func(c echolib.Context, _ struct{}) ([]SecUser, int, error) {
		if c.Request().Header.Get("X-API-Key") == "" {
			return nil, http.StatusUnauthorized, nil
		}
		return []SecUser{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	}, getOpts...)

	echo.Register(r, cfg)
	_ = r.Echo.Start(":8080")
}
