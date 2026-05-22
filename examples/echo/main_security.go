//go:build echo && !typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"
	"github.com/yzidev/goas/adapters/echoadapter"

	"github.com/yzidev/goas"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := echolib.New()
	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	cfg := goas.Config{
		Title:   "User API (Echo + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
		Security: openapi3.SecurityRequirements{bearer, apiKey},
	}

	secure := base.Group("")

	secure.GET("/secure/users", func(c echolib.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.NoContent(http.StatusUnauthorized)
		}
		return echoadapter.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	})

	secure.POST("/secure/users", func(c echolib.Context) error {
		if c.Request().Header.Get("X-API-Key") == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusCreated)
	})

	secure.POST("/secure/users/upload", func(c echolib.Context) error {
		if c.Request().Header.Get("X-API-Key") == "" {
			return echoadapter.JSON(c, http.StatusUnauthorized, goas.ErrorResponse{Error: "unauthorized"})
		}
		f, err := c.FormFile("file")
		if err != nil {
			return echoadapter.JSON(c, http.StatusBadRequest, goas.ErrorResponse{Error: "missing file"})
		}
		note := c.FormValue("note")
		return echoadapter.JSON(c, http.StatusOK, map[string]string{"filename": f.Filename, "note": note})
	})

	secure.GET("/secure/demo-errors", func(c echolib.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return echoadapter.JSON(c, http.StatusUnauthorized, goas.ErrorResponse{Error: "unauthorized"})
		}
		switch c.QueryParam("code") {
		case "400":
			return echoadapter.JSON(c, http.StatusBadRequest, goas.ErrorResponse{Error: "bad request"})
		case "500":
			return echoadapter.JSON(c, http.StatusInternalServerError, goas.ErrorResponse{Error: "internal error"})
		case "503":
			return echoadapter.JSON(c, http.StatusServiceUnavailable, goas.ErrorResponse{Error: "service unavailable"})
		default:
			return echoadapter.JSON(c, http.StatusOK, map[string]string{"status": "ok"})
		}
	})

	echoadapter.Docs(base, cfg)
	_ = base.Start(":8080")
}
