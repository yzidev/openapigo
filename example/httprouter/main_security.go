//go:build security

package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := openapi.NewRouter()

	cfg := openapi.Config{
		Title:   "User API (Security)",
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

	r := simple.New(base, spec)
	secure := r.Group("", openapi.WithTags("Secure Users"))

	// Bearer-protected endpoint
	secure.GET("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode([]SecUser{{ID: "1", Name: "Alice"}})
	})

	// API-key-protected endpoint
	secure.POST("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-Key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	secure.GET("/secure/demo-errors", func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			openapi.JSON(w, http.StatusUnauthorized, openapi.ErrorResponse{Error: "unauthorized"})
			return
		}
		switch req.URL.Query().Get("code") {
		case "400":
			openapi.JSON(w, http.StatusBadRequest, openapi.ErrorResponse{Error: "bad request"})
			return
		case "500":
			openapi.JSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Error: "internal error"})
			return
		case "503":
			openapi.JSON(w, http.StatusServiceUnavailable, openapi.ErrorResponse{Error: "service unavailable"})
			return
		default:
			openapi.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		}
	})

	openapi.Register(base, cfg)
	_ = http.ListenAndServe(":8080", r)
}
