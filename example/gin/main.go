//go:build gin && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
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
	r := gin.New()

	r.GET("/users", func(c *ginlib.Context) {
		gin.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, gin.WithTags("Users"), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: []User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.GET("/search", func(c *ginlib.Context) {
		_ = c.Query("q")
		c.Status(http.StatusOK)
	}, gin.WithTags("Users"), gin.WithQueryParams(
		openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
		openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
	), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: struct{}{}, Description: "OK"},
	))

	r.POST("/users", func(c *ginlib.Context) {
		var in CreateUser
		if err := c.ShouldBindJSON(&in); err != nil {
			gin.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		c.Status(http.StatusCreated)
	}, gin.WithTags("Users"), gin.WithRequestSchema(CreateUser{}), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusCreated, Schema: struct{}{}, Description: "Created"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.GET("/users/:id", func(c *ginlib.Context) {
		id := c.Param("id")
		if id == "404" {
			gin.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		gin.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	}, gin.WithTags("Users"), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.PUT("/users/:id", func(c *ginlib.Context) {
		id := c.Param("id")
		var in UpdateUser
		if err := c.ShouldBindJSON(&in); err != nil {
			gin.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		if id == "404" {
			gin.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		gin.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, gin.WithTags("Users"), gin.WithRequestSchema(UpdateUser{}), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.PATCH("/users/:id", func(c *ginlib.Context) {
		id := c.Param("id")
		var in UpdateUser
		if err := c.ShouldBindJSON(&in); err != nil {
			gin.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return
		}
		if id == "404" {
			gin.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		gin.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	}, gin.WithTags("Users"), gin.WithRequestSchema(UpdateUser{}), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	r.DELETE("/users/:id", func(c *ginlib.Context) {
		id := c.Param("id")
		if id == "404" {
			gin.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}, gin.WithTags("Users"), gin.WithResponses(
		openapi.ResponseSpec{Status: http.StatusNoContent, Schema: struct{}{}, Description: "No Content"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
		openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}, Description: "Internal Server Error"},
	))

	gin.Register(r, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	})
	_ = r.Engine.Run(":8080")
}
