//go:build security

package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/yzidev/openapigo/openapi"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
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

	r := openapi.New(cfg)
	secure := r.Group("", openapi.Tags("Secure Users"))

	secure.GET("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode([]SecUser{{ID: "1", Name: "Alice"}})
	}, openapi.Security(&bearer), openapi.Res([]SecUser{}))

	secure.POST("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-Key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}, openapi.Security(&apiKey), openapi.Res(struct{}{}), openapi.Created())

	secure.POST("/secure/users/upload", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-Key") == "" {
			openapi.JSON(w, http.StatusUnauthorized, openapi.ErrorResponse{Error: "unauthorized"})
			return
		}
		if err := req.ParseMultipartForm(10 << 20); err != nil {
			openapi.JSON(w, http.StatusBadRequest, openapi.ErrorResponse{Error: "invalid multipart"})
			return
		}
		f, fh, err := req.FormFile("file")
		if err != nil {
			openapi.JSON(w, http.StatusBadRequest, openapi.ErrorResponse{Error: "missing file"})
			return
		}
		_ = f.Close()
		note := req.FormValue("note")
		openapi.JSON(w, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	},
		openapi.Security(&apiKey),
		openapi.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
		openapi.Res(map[string]string{}),
	)

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
	}, openapi.Security(&bearer), openapi.Res(map[string]string{}))

	r.Docs()
	_ = http.ListenAndServe(":8080", r)
}
