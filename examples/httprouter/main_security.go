//go:build security

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yzidev/goas"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	cfg := goas.Config{
		Title:   "User API (Security)",
		Version: "1.0.0",
		Tags: goas.DocumentTags{
			{Name: "Secure Users", Description: "Secured endpoints (Bearer / X-API-Key)"},
		},
		SecuritySchemes: map[string]*goas.SecuritySchemeRef{
			"bearerAuth": {Value: &goas.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &goas.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := goas.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := goas.NewSecurityRequirement().Authenticate("apiKeyAuth")
	cfg.Security = goas.SecurityRequirements{bearer, apiKey}

	r := goas.New(cfg)
	secure := r.Group("", goas.Tags("Secure Users"))

	secure.GET("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode([]SecUser{{ID: "1", Name: "Alice"}})
	}, goas.Security(&bearer), goas.Res([]SecUser{}))

	secure.POST("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-Key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}, goas.Security(&apiKey), goas.Res(struct{}{}), goas.Created())

	secure.POST("/secure/users/upload", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-Key") == "" {
			goas.JSON(w, http.StatusUnauthorized, goas.ErrorResponse{Error: "unauthorized"})
			return
		}
		if err := req.ParseMultipartForm(10 << 20); err != nil {
			goas.JSON(w, http.StatusBadRequest, goas.ErrorResponse{Error: "invalid multipart"})
			return
		}
		f, fh, err := req.FormFile("file")
		if err != nil {
			goas.JSON(w, http.StatusBadRequest, goas.ErrorResponse{Error: "missing file"})
			return
		}
		_ = f.Close()
		note := req.FormValue("note")
		goas.JSON(w, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	},
		goas.Security(&apiKey),
		goas.MultipartUpload("file", goas.MultipartField{Name: "note", Type: goas.ParamString}),
		goas.Res(map[string]string{}),
	)

	secure.GET("/secure/demo-errors", func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			goas.JSON(w, http.StatusUnauthorized, goas.ErrorResponse{Error: "unauthorized"})
			return
		}
		switch req.URL.Query().Get("code") {
		case "400":
			goas.JSON(w, http.StatusBadRequest, goas.ErrorResponse{Error: "bad request"})
			return
		case "500":
			goas.JSON(w, http.StatusInternalServerError, goas.ErrorResponse{Error: "internal error"})
			return
		case "503":
			goas.JSON(w, http.StatusServiceUnavailable, goas.ErrorResponse{Error: "service unavailable"})
			return
		default:
			goas.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
		}
	}, goas.Security(&bearer), goas.Res(map[string]string{}))

	r.Docs()
	log.Fatal(http.ListenAndServe(":8080", r))
}
