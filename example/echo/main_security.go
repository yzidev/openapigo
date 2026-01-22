//go:build echo && !typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := echo.New()

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

	b := simple.NewSpec()
	b.GroupTags("", []string{"Secure Users"}, func(s *simple.SpecBuilder) {
		s.GET("/secure/users").Security(&bearer).Res([]SecUser{}).OK()
		s.POST("/secure/users").Security(&apiKey).Res(struct{}{}).Created()

		s.GET("/secure/demo-errors").Security(&bearer).Res(map[string]string{}).OK().Responses(
			openapi.ResponseSpec{Status: 400, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: 401, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: 500, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: 503, Schema: openapi.ErrorResponse{}},
		)
	})

	spec := b.Spec()

	r := simple.NewEcho(base, spec)
	secure := r.Group("", echo.WithTags("Secure Users"))

	secure.GET("/secure/users", func(c echolib.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.NoContent(http.StatusUnauthorized)
		}
		return echo.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	})

	secure.POST("/secure/users", func(c echolib.Context) error {
		if c.Request().Header.Get("X-API-Key") == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusCreated)
	})

	secure.GET("/secure/demo-errors", func(c echolib.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return echo.JSON(c, http.StatusUnauthorized, openapi.ErrorResponse{Error: "unauthorized"})
		}
		switch c.QueryParam("code") {
		case "400":
			return echo.JSON(c, http.StatusBadRequest, openapi.ErrorResponse{Error: "bad request"})
		case "500":
			return echo.JSON(c, http.StatusInternalServerError, openapi.ErrorResponse{Error: "internal error"})
		case "503":
			return echo.JSON(c, http.StatusServiceUnavailable, openapi.ErrorResponse{Error: "service unavailable"})
		default:
			return echo.JSON(c, http.StatusOK, map[string]string{"status": "ok"})
		}
	})

	echo.Register(base, cfg)
	_ = base.Echo.Start(":8080")
}
