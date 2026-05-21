//go:build gin && security && !typed

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/ginadapter"
	"github.com/aizacoders/openapigo/openapi"
)

func openAPICfgSecurity() (openapi.Config, *openapi3.SecurityRequirement, *openapi3.SecurityRequirement) {
	cfg := openapi.Config{
		Title:       "User API (Gin + Security)",
		Version:     "1.0.0",
		Description: "An examples API with secured endpoints using Gin and OpenAPIGO",
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

func registerSecureRoutes(r *ginadapter.Router, bearer, apiKey *openapi3.SecurityRequirement) {
	r.GET("/secure/healthz", handleSecureHealthz,
		ginadapter.Tags("Secure Users"),
		ginadapter.Security(bearer),
		ginadapter.Res(map[string]string{}),
	)

	secure := r.Group("", ginadapter.Tags("Secure Users"))
	secure.GET("/secure/users", handleSecureListUsers, ginadapter.Security(bearer), ginadapter.Res([]SecUser{}))
	secure.POST("/secure/users", handleSecureCreateUser, ginadapter.Security(apiKey), ginadapter.Created())
	secure.POST("/secure/users/upload", handleSecureUploadUserFile,
		ginadapter.Security(apiKey),
		ginadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
		ginadapter.Res(map[string]string{}),
	)
	secure.GET("/secure/demo-errors", handleSecureDemoErrors,
		ginadapter.Security(bearer),
		ginadapter.Res(map[string]string{}),
		ginadapter.Responses(
			openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: http.StatusUnauthorized, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: openapi.ErrorResponse{}},
			openapi.ResponseSpec{Status: http.StatusServiceUnavailable, Schema: openapi.ErrorResponse{}},
		),
	)
}
