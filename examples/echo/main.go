//go:build echo && !typed && !security

package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/adapters/echoadapter"
	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/openapi"
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
	base := echolib.New()
	r := echoadapter.Wrap(base)

	users := r.Group("", echoadapter.Tags("Users"))
	users.GET("/users", func(c echolib.Context) error {
		return echoadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, echoadapter.Res([]User{}))

	users.GET("/search", func(c echolib.Context) error {
		_ = c.QueryParam("q")
		return c.NoContent(http.StatusOK)
	},
		echoadapter.Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		),
		echoadapter.Res(struct{}{}),
	)

	users.POST("/users", func(c echolib.Context) error {
		var in CreateUser
		if err := echoadapter.Bind(c, &in); err != nil || in.Name == "" {
			_ = echoadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.NoContent(http.StatusCreated)
	}, echoadapter.Req(CreateUser{}), echoadapter.Created())

	users.POST("/users/upload", func(c echolib.Context) error {
		f, err := c.FormFile("file")
		if err != nil {
			_ = echoadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
			return nil
		}
		note := c.FormValue("note")
		return echoadapter.JSON(c, http.StatusOK, map[string]string{"filename": f.Filename, "note": note})
	},
		echoadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
		echoadapter.Res(map[string]string{}),
	)

	users.GET("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			return echoadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echoadapter.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	}, echoadapter.Res(User{}))

	users.PUT("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := echoadapter.Bind(c, &in); err != nil {
			_ = echoadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = echoadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return echoadapter.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, echoadapter.Req(UpdateUser{}), echoadapter.Res(User{}))

	users.PATCH("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := echoadapter.Bind(c, &in); err != nil {
			_ = echoadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = echoadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return echoadapter.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, echoadapter.Req(UpdateUser{}), echoadapter.Res(User{}))

	users.DELETE("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			_ = echoadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.NoContent(http.StatusNoContent)
	}, echoadapter.NoContent())

	r.Docs(openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = r.Echo.Start(":8080")
}
