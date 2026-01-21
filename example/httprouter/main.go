//go:build !security

package main

import (
	"encoding/json"
	"net/http"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/getkin/kin-openapi/openapi3"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

func main() {
	r := openapi.NewRouter()

	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode([]User{{ID: "1", Name: "Alice"}})
	}, openapi.WithTags("Users"), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: []User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.GET("/search", func(w http.ResponseWriter, req *http.Request) {
		_, _, _ = openapi.QueryValue[int](req, "limit")
		w.WriteHeader(http.StatusOK)
	}, openapi.WithTags("Users"), openapi.WithQueryParams(
		openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
		openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
	), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: struct{}{}, Description: "OK"},
	))

	r.POST("/users", func(w http.ResponseWriter, req *http.Request) {
		var in CreateUser
		if err := openapi.Bind(req, &in); err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		w.WriteHeader(http.StatusCreated)
	}, openapi.WithTags("Users"), openapi.WithRequestSchema(CreateUser{}), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusCreated, Schema: struct{}{}, Description: "Created"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	// GET /users/{id}
	r.GET("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: "Alice"})
	}, openapi.WithTags("Users"), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	// PUT /users/{id}
	r.PUT("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		var in UpdateUser
		if err := openapi.Bind(req, &in); err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: in.Name})
	}, openapi.WithTags("Users"), openapi.WithRequestSchema(UpdateUser{}), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	// PATCH /users/{id}
	r.PATCH("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		var in UpdateUser
		if err := openapi.Bind(req, &in); err != nil {
			openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		openapi.JSON(w, http.StatusOK, User{ID: id, Name: in.Name})
	}, openapi.WithTags("Users"), openapi.WithRequestSchema(UpdateUser{}), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	// DELETE /users/{id}
	r.DELETE("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			openapi.JSON(w, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}, openapi.WithTags("Users"), openapi.WithResponses(
		openapi.ResponseSpec{Status: http.StatusNoContent, Schema: struct{}{}, Description: "No Content"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	openapi.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	_ = http.ListenAndServe(":8080", r)
}
