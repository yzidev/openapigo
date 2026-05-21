//go:build !security

package main

import (
	"encoding/json"
	"net/http"

	"github.com/aizacoders/openapigo/adapters/muxadapter"
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

// --- Handlers for v1 (users group)
func listUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode([]User{{ID: "1", Name: "Alice"}})
}

func searchUsers(w http.ResponseWriter, req *http.Request) {
	_, _, _ = openapi.QueryValue[int](req, "limit")
	w.WriteHeader(http.StatusOK)
}

func createUser(w http.ResponseWriter, req *http.Request) {
	var in CreateUser
	if err := openapi.Bind(req, &in); err != nil || in.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid body"})
		return
	}
	created := User{ID: "2", Name: in.Name}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}

func uploadUserFile(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid multipart"})
		return
	}
	f, fh, err := req.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "missing file"})
		return
	}
	_ = f.Close()
	note := req.FormValue("note")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"filename": fh.Filename, "note": note})
}

func getUser(w http.ResponseWriter, req *http.Request) {
	id := openapi.PathValue(req, "id")
	if id == "404" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(User{ID: id, Name: "Alice"})
}

func putUser(w http.ResponseWriter, req *http.Request) {
	id := openapi.PathValue(req, "id")
	var in UpdateUser
	if err := openapi.Bind(req, &in); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid body"})
		return
	}
	if id == "404" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(User{ID: id, Name: in.Name})
}

func patchUser(w http.ResponseWriter, req *http.Request) {
	id := openapi.PathValue(req, "id")
	var in UpdateUser
	if err := openapi.Bind(req, &in); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid body"})
		return
	}
	if id == "404" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(User{ID: id, Name: in.Name})
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	id := openapi.PathValue(req, "id")
	if id == "404" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- API2 handlers (moved into main.go as requested)
func api2Hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello from api2"})
}

func api2CreateUser(w http.ResponseWriter, req *http.Request) {
	var in CreateUser
	if err := openapi.Bind(req, &in); err != nil || in.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid body"})
		return
	}
	created := User{ID: "100", Name: in.Name}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}

func api2GetUser(w http.ResponseWriter, req *http.Request) {
	id := openapi.PathValue(req, "id")
	if id == "404" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "user not found (api2)"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(User{ID: id, Name: "API2 User"})
}

func main() {
	mux := http.NewServeMux()

	r := muxadapter.New(mux)

	users := r.Group("/api/v1.0", openapi.Tags("Users"))
	users.GET("/users", listUsers, openapi.Res([]User{}))
	users.GET("/search", searchUsers,
		openapi.Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		),
		openapi.Res(struct{}{}),
	)
	users.POST("/users", createUser, openapi.Req(CreateUser{}), openapi.Res(User{}), openapi.Created())
	users.POST("/users/upload", uploadUserFile,
		openapi.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
		openapi.Res(map[string]string{}),
	)
	users.GET("/users/{id}", getUser, openapi.Res(User{}))
	users.PUT("/users/{id}", putUser, openapi.Req(UpdateUser{}), openapi.Res(User{}))
	users.PATCH("/users/{id}", patchUser, openapi.Req(UpdateUser{}), openapi.Res(User{}))
	users.DELETE("/users/{id}", deleteUser, openapi.NoContent())

	api2 := r.Group("/api/v2.0", openapi.Tags("Users V2.0"))
	api2.GET("/hello", api2Hello, openapi.Res(map[string]string{}))
	api2.POST("/users", api2CreateUser, openapi.Req(CreateUser{}), openapi.Res(User{}), openapi.Created())
	api2.GET("/users/{id}", api2GetUser, openapi.Res(User{}))

	r.Docs(openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	_ = http.ListenAndServe(":8080", mux)
}
