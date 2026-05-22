//go:build gin && !typed && !security

package main

import (
	"github.com/getkin/kin-openapi/openapi3"
	ginlib "github.com/gin-gonic/gin"

	"github.com/yzidev/goas/openapi"
)

// registerRoutes wires the endpoints in a readable and grouped way.
// (Non-typed, non-security variant.)

func registerSystemRoutes(r *ginlib.Engine) {
	r.GET("/healthz", handleHealthz)
}

func registerUserRoutes(r *ginlib.Engine) {
	users := r.Group("")

	users.GET("/users", handleListUsers)
	users.GET("/search", handleSearchUsers)
	users.POST("/users", handleCreateUser)
	users.POST("/users/upload", handleUploadUserFile)
	users.GET("/users/demo-errors", handleDemoErrors)
	users.GET("/users/:id", handleGetUser)
	users.PUT("/users/:id", handlePutUser)
	users.PATCH("/users/:id", handlePatchUser)
	users.DELETE("/users/:id", handleDeleteUser)
}

func openAPICfg() openapi.Config {
	return openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	}
}
