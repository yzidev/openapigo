//go:build gin && !typed && security

package main

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	cfg := openapi.Config{
		Title:   "User API (Gin + Security)",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Secure Users", Description: "Secured endpoints (Bearer / X-API-Key)"},
		},
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	// Sample route WITHOUT group: secured + direct registration.
	secureHealthzOpts := append(
		[]gin.HandlerOption{gin.WithTags("System"), gin.WithSecurity(&bearer)},
		gin.JSONRoute(struct{}{}, map[string]string{}, http.StatusOK)...,
	)
	r.GET("/secure/healthz", func(c *ginlib.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.Status(http.StatusUnauthorized)
			return
		}
		gin.JSON(c, http.StatusOK, map[string]string{"status": "ok"})
	}, secureHealthzOpts...)

	secure := r.Group("", gin.WithTags("Secure Users"))

	secure.GET("/secure/users", func(c *ginlib.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.Status(http.StatusUnauthorized)
			return
		}
		gin.JSON(c, http.StatusOK, []SecUser{{ID: "1", Name: "Alice"}})
	}, gin.WithSecurity(&bearer))

	postOpts := append(
		[]gin.HandlerOption{gin.WithSecurity(&apiKey)},
		gin.JSONRoute(nil, struct{}{}, http.StatusCreated)...,
	)
	secure.POST("/secure/users", func(c *ginlib.Context) {
		if c.GetHeader("X-API-Key") == "" {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusCreated)
	}, postOpts...)

	gin.Register(r, cfg)
	_ = r.Engine.Run(":8080")
}
