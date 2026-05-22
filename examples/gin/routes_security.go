//go:build gin && security && !typed

package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"

	"github.com/yzidev/goas/openapi"
)

func openAPICfgSecurity() openapi.Config {
	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	return openapi.Config{
		Title:       "User API (Gin + Security)",
		Version:     "1.0.0",
		Description: "An examples API with secured endpoints using Gin and Goas",
		Tags: openapi3.Tags{
			{Name: "Secure Users", Description: "Secured endpoints (Bearer / X-API-Key)"},
		},
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
		Security: openapi3.SecurityRequirements{bearer, apiKey},
	}
}

func registerSecureRoutes(r *ginlib.Engine) {
	r.GET("/secure/healthz", handleSecureHealthz)

	secure := r.Group("")
	secure.GET("/secure/users", handleSecureListUsers)
	secure.POST("/secure/users", handleSecureCreateUser)
	secure.POST("/secure/users/upload", handleSecureUploadUserFile)
	secure.GET("/secure/demo-errors", handleSecureDemoErrors)
}
