//go:build !security

package main

import (
	"encoding/json"
	"net/http"

	muxadapter "github.com/aizacoders/openapigo/adapters/mux"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
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

// openapiSpec builds the OpenAPI spec for the example (keeps main() clean),
// includes both API v1.0 and v2.0 groups to mirror the routes in main().
func openapiSpec() simple.Spec {
	b := simple.NewSpec()

	// API v1.0 (Users)
	b.GroupTags("/api/v1.0", []string{"Users"}, func(s *simple.SpecBuilder) {
		s.GET("/users").Res([]User{}).OK()
		s.GET("/search").Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		).Res(struct{}{}).OK()
		// Create user returns created resource (status 201)
		s.POST("/users").Req(CreateUser{}).Res(User{}).Created()
		// Upload user file: multipart/form-data.
		s.POST("/users/upload").MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}).Res(map[string]string{}).OK()
		s.GET("/users/{id}").Res(User{}).OK()
		s.PUT("/users/{id}").Req(UpdateUser{}).Res(User{}).OK()
		s.PATCH("/users/{id}").Req(UpdateUser{}).Res(User{}).OK()
		s.DELETE("/users/{id}").Res(struct{}{}).NoContent()
	})

	// API v2.0 (Users V2.0) - smaller surface per main.go: hello + users create/get
	b.GroupTags("/api/v2.0", []string{"Users V2.0"}, func(s *simple.SpecBuilder) {
		s.GET("/hello").Res(map[string]string{}).OK()
		// v2 create returns the created user object
		s.POST("/users").Req(CreateUser{}).Res(User{}).Created()
		s.GET("/users/{id}").Res(User{}).OK()
	})

	return b.Spec()
}

func main() {
	mux := http.NewServeMux()

	// Router setup by openapigo/httprouter adapter
	base := muxadapter.NewHttpAdapters(mux)
	r := simple.NewHttpRouter(base, openapiSpec())

	// Clean routes: just HTTP methods + handlers.
	users := r.Group("/api/v1.0", openapi.WithTags("Users"))
	users.GET("/users", listUsers)
	users.GET("/search", searchUsers)
	users.POST("/users", createUser)
	users.POST("/users/upload", uploadUserFile)
	users.GET("/users/{id}", getUser)
	users.PUT("/users/{id}", putUser)
	users.PATCH("/users/{id}", patchUser)
	users.DELETE("/users/{id}", deleteUser)

	// --- API2: separate group that uses default writer (encoding/json) instead of openapi.JSON
	api2 := r.Group("/api/v2.0", openapi.WithTags("Users V2.0"))
	api2.GET("/hello", api2Hello)
	api2.POST("/users", api2CreateUser)
	api2.GET("/users/{id}", api2GetUser)

	muxadapter.Register(base, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	// Server: ServeMux already has the router mounted by httprouter.New(mux)

	_ = http.ListenAndServe(":8080", mux)
}
