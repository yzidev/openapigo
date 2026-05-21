//go:build fiber && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/yzidev/openapigo/adapters/fiberadapter"
	"github.com/yzidev/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	base := fiberlib.New()
	r := fiberadapter.Wrap(base)

	users := r.Group("", fiberadapter.Tags("Users"))
	users.GET("/users", func(c *fiberlib.Ctx) error {
		return fiberadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, fiberadapter.Res([]User{}))

	users.GET("/search", func(c *fiberlib.Ctx) error {
		_ = c.Query("q")
		return c.SendStatus(http.StatusOK)
	},
		fiberadapter.Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		),
		fiberadapter.Res(struct{}{}),
	)

	users.POST("/users", func(c *fiberlib.Ctx) error {
		var in CreateUser
		if err := fiberadapter.Bind(c, &in); err != nil || in.Name == "" {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.SendStatus(http.StatusCreated)
	}, fiberadapter.Req(CreateUser{}), fiberadapter.Created())

	users.POST("/users/upload", func(c *fiberlib.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
			return fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
		}
		note := c.FormValue("note")
		return fiberadapter.JSON(c, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	},
		fiberadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
		fiberadapter.Res(map[string]string{}),
	)

	users.GET("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			return fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	}, fiberadapter.Res(User{}))

	users.PUT("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiberadapter.Bind(c, &in); err != nil {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, fiberadapter.Req(UpdateUser{}), fiberadapter.Res(User{}))

	users.PATCH("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiberadapter.Bind(c, &in); err != nil {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, fiberadapter.Req(UpdateUser{}), fiberadapter.Res(User{}))

	users.DELETE("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.SendStatus(http.StatusNoContent)
	}, fiberadapter.NoContent())

	r.Docs(openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = r.App.Listen(":8080")
}
