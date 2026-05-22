//go:build fiber && !typed && security

package main

import (
	"log"
	"net/http"
	"strings"

	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/yzidev/goas"
	"github.com/yzidev/goas/adapters/fiberadapter"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := fiberlib.New()
	bearer := goas.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := goas.NewSecurityRequirement().Authenticate("apiKeyAuth")

	cfg := goas.Config{
		Title:   "User API (Fiber + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*goas.SecuritySchemeRef{
			"bearerAuth": {Value: &goas.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &goas.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
		Security: goas.SecurityRequirements{bearer, apiKey},
	}

	secure := base.Group("")

	secure.Get("/secure/users", func(c *fiberlib.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.SendStatus(http.StatusUnauthorized)
		}
		return fiberadapter.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	})

	secure.Post("/secure/users", func(c *fiberlib.Ctx) error {
		if c.Get("X-API-Key") == "" {
			return c.SendStatus(http.StatusUnauthorized)
		}
		return c.SendStatus(http.StatusCreated)
	})

	secure.Post("/secure/users/upload", func(c *fiberlib.Ctx) error {
		if c.Get("X-API-Key") == "" {
			return fiberadapter.JSON(c, http.StatusUnauthorized, goas.ErrorResponse{Error: "unauthorized"})
		}
		fh, err := c.FormFile("file")
		if err != nil {
			return fiberadapter.JSON(c, http.StatusBadRequest, goas.ErrorResponse{Error: "missing file"})
		}
		note := c.FormValue("note")
		return fiberadapter.JSON(c, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	})

	secure.Get("/secure/demo-errors", func(c *fiberlib.Ctx) error {
		if !strings.HasPrefix(c.Get("Authorization"), "Bearer ") {
			return c.SendStatus(http.StatusUnauthorized)
		}
		switch c.Query("code") {
		case "400":
			return fiberadapter.JSON(c, http.StatusBadRequest, goas.ErrorResponse{Error: "bad request"})
		case "500":
			return fiberadapter.JSON(c, http.StatusInternalServerError, goas.ErrorResponse{Error: "internal error"})
		case "503":
			return fiberadapter.JSON(c, http.StatusServiceUnavailable, goas.ErrorResponse{Error: "service unavailable"})
		default:
			return fiberadapter.JSON(c, http.StatusOK, map[string]string{"status": "ok"})
		}
	})

	fiberadapter.Docs(base, cfg)
	log.Fatal(base.Listen(":8080"))
}
