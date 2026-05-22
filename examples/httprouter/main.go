//go:build !security

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yzidev/goas"
	"github.com/yzidev/goas/adapters/muxadapter"
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
	_, _, _ = goas.QueryValue[int](req, "limit")
	w.WriteHeader(http.StatusOK)
}

func createUser(w http.ResponseWriter, req *http.Request) {
	var in CreateUser
	if err := goas.Bind(req, &in); err != nil || in.Name == "" {
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
	id := goas.PathValue(req, "id")
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
	id := goas.PathValue(req, "id")
	var in UpdateUser
	if err := goas.Bind(req, &in); err != nil {
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
	id := goas.PathValue(req, "id")
	var in UpdateUser
	if err := goas.Bind(req, &in); err != nil {
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
	id := goas.PathValue(req, "id")
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
	if err := goas.Bind(req, &in); err != nil || in.Name == "" {
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
	id := goas.PathValue(req, "id")
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

	r := muxadapter.Mount(mux, goas.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: goas.DocumentTags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})

	users := r.Group("/api/v1.0", goas.Tags("Users"))
	users.GET("/users", listUsers, goas.Res([]User{}))
	users.GET("/search", searchUsers,
		goas.Query(
			goas.QueryParam{Name: "q", Type: goas.ParamString, Required: true, Description: "Search term"},
			goas.QueryParam{Name: "limit", Type: goas.ParamInteger, Required: false, Description: "Max results"},
		),
		goas.Res(struct{}{}),
	)
	users.POST("/users", createUser, goas.Req(CreateUser{}), goas.Res(User{}), goas.Created())
	users.POST("/users/upload", uploadUserFile,
		goas.MultipartUpload("file", goas.MultipartField{Name: "note", Type: goas.ParamString}),
		goas.Res(map[string]string{}),
	)
	users.GET("/users/{id}", getUser, goas.Res(User{}))
	users.PUT("/users/{id}", putUser, goas.Req(UpdateUser{}), goas.Res(User{}))
	users.PATCH("/users/{id}", patchUser, goas.Req(UpdateUser{}), goas.Res(User{}))
	users.DELETE("/users/{id}", deleteUser, goas.NoContent())

	api2 := r.Group("/api/v2.0", goas.Tags("Users V2.0"))
	api2.GET("/hello", api2Hello, goas.Res(map[string]string{}))
	api2.POST("/users", api2CreateUser, goas.Req(CreateUser{}), goas.Res(User{}), goas.Created())
	api2.GET("/users/{id}", api2GetUser, goas.Res(User{}))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
