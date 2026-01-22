//go:build gin && security && !typed

package main

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

func openAPICfgSecurity() (openapi.Config, *openapi3.SecurityRequirement, *openapi3.SecurityRequirement) {
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
	return cfg, &bearer, &apiKey
}

func registerSecureRoutes(r *simple.GinRouter, bearer, apiKey *openapi3.SecurityRequirement) {
	b := simple.NewSpec()
	b.GroupTags("", []string{"Secure Users"}, func(s *simple.SpecBuilder) {
		s.GET("/secure/healthz").Security(bearer).Res(map[string]string{}).OK()
		s.GET("/secure/users").Security(bearer).Res([]SecUser{}).OK()
		s.POST("/secure/users").Security(apiKey).Res(struct{}{}).Created()

		s.GET("/secure/demo-errors").Security(bearer).Res(map[string]string{}).OK().Responses(
			openapi.ResponseSpec{Status: 400, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: 401, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: 500, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: 503, Schema: openapi.ErrorResponse{}},
		)
	})
	r.Spec = b.Spec()

	r.GET("/secure/healthz", handleSecureHealthz)

	secure := r.Group("", gin.WithTags("Secure Users"))
	secure.GET("/secure/users", handleSecureListUsers)
	secure.POST("/secure/users", handleSecureCreateUser)
	secure.GET("/secure/demo-errors", handleSecureDemoErrors)
}
