//go:build gin && typed && security

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

type SecCreateUser struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SecureHealth struct {
	Status string `json:"status"`
}

func main() {
	r := gin.New()

	cfg := openapi.Config{
		Title:   "User API (Gin + Security)",
		Version: "1.0.0",
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	// Sample route WITHOUT group: typed + security + direct registration.
	secureHealthOpts := append(
		[]gin.HandlerOption{gin.WithTags("System"), gin.WithSecurity(&bearer)},
		gin.JSONRoute(struct{}{}, SecureHealth{}, http.StatusOK)...,
	)
	gin.GETT[struct{}, SecureHealth](r, "/secure/typed/healthz", func(c *ginlib.Context, _ struct{}) (SecureHealth, int, error) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return SecureHealth{}, http.StatusUnauthorized, nil
		}
		return SecureHealth{Status: "ok"}, http.StatusOK, nil
	}, secureHealthOpts...)

	secureOpts := []gin.HandlerOption{gin.WithTags("Secure Users")}

	postOpts := append(
		append([]gin.HandlerOption{}, secureOpts...),
		gin.WithSecurity(&bearer),
	)
	postOpts = append(postOpts, gin.JSONRoute(SecCreateUser{}, SecUser{}, http.StatusCreated)...)
	gin.POSTT[SecCreateUser, SecUser](r, "/secure/users", func(c *ginlib.Context, in SecCreateUser) (SecUser, int, error) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return SecUser{}, http.StatusUnauthorized, nil
		}
		return SecUser{ID: "1", Name: in.Name}, http.StatusCreated, nil
	}, postOpts...)

	getOpts := append(
		append([]gin.HandlerOption{}, secureOpts...),
		gin.WithSecurity(&apiKey),
	)
	getOpts = append(getOpts, gin.JSONRoute(struct{}{}, []SecUser{}, http.StatusOK)...)
	gin.GETT[struct{}, []SecUser](r, "/secure/users", func(c *ginlib.Context, _ struct{}) ([]SecUser, int, error) {
		if c.GetHeader("X-API-Key") == "" {
			return nil, http.StatusUnauthorized, nil
		}
		return []SecUser{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	}, getOpts...)

	gin.Register(r, cfg)
	_ = r.Engine.Run(":8080")
}
