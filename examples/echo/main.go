//go:build echo && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"
	"github.com/yzidev/goas/adapters/echoadapter"

	"github.com/yzidev/goas/openapi"
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

	users := base.Group("")
	users.GET("/users", func(c echolib.Context) error {
		return echoadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.GET("/search", func(c echolib.Context) error {
		_ = c.QueryParam("q")
		return c.NoContent(http.StatusOK)
	})

	users.POST("/users", func(c echolib.Context) error {
		var in CreateUser
		if err := echoadapter.Bind(c, &in); err != nil || in.Name == "" {
			_ = echoadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.NoContent(http.StatusCreated)
	})

	users.POST("/users/upload", func(c echolib.Context) error {
		f, err := c.FormFile("file")
		if err != nil {
			_ = echoadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
			return nil
		}
		note := c.FormValue("note")
		return echoadapter.JSON(c, http.StatusOK, map[string]string{"filename": f.Filename, "note": note})
	})

	users.GET("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			return echoadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echoadapter.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

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
	})

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
	})

	users.DELETE("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			_ = echoadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.NoContent(http.StatusNoContent)
	})

	echoadapter.Docs(base, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = base.Start(":8080")
}
