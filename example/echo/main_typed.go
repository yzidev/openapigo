//go:build echo && typed && !security

package main

import (
	"net/http"

	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
)

type UserTyped struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserTyped struct {
	Name string `json:"name"`
}

func main() {
	r := echo.New()

	typedOpts := []echo.HandlerOption{echo.WithTags("Typed")}

	postOpts := append(append([]echo.HandlerOption{}, typedOpts...), echo.JSONRoute(CreateUserTyped{}, UserTyped{}, http.StatusCreated)...)
	echo.POSTT[CreateUserTyped, UserTyped](r, "/typed/users", func(c echolib.Context, in CreateUserTyped) (UserTyped, int, error) {
		return UserTyped{ID: "1", Name: in.Name}, http.StatusCreated, nil
	}, postOpts...)

	getOpts := append(append([]echo.HandlerOption{}, typedOpts...), echo.JSONRoute(struct{}{}, []UserTyped{}, http.StatusOK)...)
	echo.GETT[struct{}, []UserTyped](r, "/typed/users", func(c echolib.Context, _ struct{}) ([]UserTyped, int, error) {
		return []UserTyped{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	}, getOpts...)

	echo.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = r.Echo.Start(":8080")
}
