//go:build fiber && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/yzidev/goas/adapters/fiberadapter"
	"github.com/yzidev/goas/openapi"
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

	users := base.Group("")
	users.Get("/users", func(c *fiberlib.Ctx) error {
		return fiberadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.Get("/search", func(c *fiberlib.Ctx) error {
		_ = c.Query("q")
		return c.SendStatus(http.StatusOK)
	})

	users.Post("/users", func(c *fiberlib.Ctx) error {
		var in CreateUser
		if err := fiberadapter.Bind(c, &in); err != nil || in.Name == "" {
			_ = fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.SendStatus(http.StatusCreated)
	})

	users.Post("/users/upload", func(c *fiberlib.Ctx) error {
		fh, err := c.FormFile("file")
		if err != nil {
			return fiberadapter.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "missing file"})
		}
		note := c.FormValue("note")
		return fiberadapter.JSON(c, http.StatusOK, map[string]string{"filename": fh.Filename, "note": note})
	})

	users.Get("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			return fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return fiberadapter.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	users.Put("/users/:id", func(c *fiberlib.Ctx) error {
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
	})

	users.Patch("/users/:id", func(c *fiberlib.Ctx) error {
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
	})

	users.Delete("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			_ = fiberadapter.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.SendStatus(http.StatusNoContent)
	})

	fiberadapter.Docs(base, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = base.Listen(":8080")
}
