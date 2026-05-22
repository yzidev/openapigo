//go:build gin && security && !typed

package main

import (
	ginlib "github.com/gin-gonic/gin"

	"github.com/yzidev/goas"
)

func openAPICfgSecurity() goas.Config {
	bearer := goas.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := goas.NewSecurityRequirement().Authenticate("apiKeyAuth")

	return goas.Config{
		Title:       "User API (Gin + Security)",
		Version:     "1.0.0",
		Description: "An examples API with secured endpoints using Gin and Goas",
		Tags: goas.DocumentTags{
			{Name: "Secure Users", Description: "Secured endpoints (Bearer / X-API-Key)"},
		},
		SecuritySchemes: map[string]*goas.SecuritySchemeRef{
			"bearerAuth": {Value: &goas.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &goas.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
		Security: goas.SecurityRequirements{bearer, apiKey},
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
