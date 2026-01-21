//go:build !security

package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/adapters/httprouter"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

func main() {
	r := httprouter.New()

	usersOpts := []httprouter.HandlerOption{httprouter.WithTags("Users")}

	postOpts := append(append([]httprouter.HandlerOption{}, usersOpts...), httprouter.JSONRoute(CreateUser{}, User{}, http.StatusCreated)...)
	httprouter.POSTT[CreateUser, User](r, "/users", func(w http.ResponseWriter, req *http.Request, in CreateUser) (User, int, error) {
		_ = req
		return User{ID: "1", Name: in.Name}, http.StatusCreated, nil
	}, postOpts...)

	getListOpts := append(append([]httprouter.HandlerOption{}, usersOpts...), httprouter.JSONRoute(struct{}{}, []User{}, http.StatusOK)...)
	httprouter.GETT[struct{}, []User](r, "/users", func(w http.ResponseWriter, req *http.Request, _ struct{}) ([]User, int, error) {
		_ = w
		_ = req
		return []User{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	}, getListOpts...)

	getByIDOpts := append(append([]httprouter.HandlerOption{}, usersOpts...), httprouter.JSONRoute(struct{}{}, User{}, http.StatusOK)...)
	httprouter.GETT[struct{}, User](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, _ struct{}) (User, int, error) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			return User{}, http.StatusNotFound, nil
		}
		return User{ID: id, Name: "Alice"}, http.StatusOK, nil
	}, getByIDOpts...)

	putOpts := append(append([]httprouter.HandlerOption{}, usersOpts...), httprouter.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)
	httprouter.PUTT[UpdateUser, User](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, in UpdateUser) (User, int, error) {
		id := openapi.PathValue(req, "id")
		if in.Name == "" {
			return User{}, http.StatusBadRequest, nil
		}
		if id == "404" {
			return User{}, http.StatusNotFound, nil
		}
		return User{ID: id, Name: in.Name}, http.StatusOK, nil
	}, putOpts...)

	patchOpts := append(append([]httprouter.HandlerOption{}, usersOpts...), httprouter.JSONRoute(UpdateUser{}, User{}, http.StatusOK)...)
	httprouter.PATCHT[UpdateUser, User](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, in UpdateUser) (User, int, error) {
		id := openapi.PathValue(req, "id")
		if in.Name == "" {
			return User{}, http.StatusBadRequest, nil
		}
		if id == "404" {
			return User{}, http.StatusNotFound, nil
		}
		return User{ID: id, Name: in.Name}, http.StatusOK, nil
	}, patchOpts...)

	deleteOpts := append(append([]httprouter.HandlerOption{}, usersOpts...), httprouter.JSONRoute(struct{}{}, struct{}{}, http.StatusNoContent)...)
	httprouter.DELETET[struct{}, struct{}](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, _ struct{}) (struct{}, int, error) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			return struct{}{}, http.StatusNotFound, nil
		}
		return struct{}{}, http.StatusNoContent, nil
	}, deleteOpts...)

	openapi.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = http.ListenAndServe(":8080", r)
}
