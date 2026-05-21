//go:build gin && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/yzidev/openapigo/adapters/ginadapter"
	"github.com/yzidev/openapigo/openapi"
)

// registerRoutes wires the endpoints in a readable and grouped way.
// (Non-typed, non-security variant.)

func registerSystemRoutes(r *ginadapter.Router) {
	r.GET("/healthz", handleHealthz, ginadapter.Res(map[string]string{}), ginadapter.Tags("System"))
}

func registerUserRoutes(r *ginadapter.Router) {
	users := r.Group("", ginadapter.Tags("Users"))

	users.GET("/users", handleListUsers, ginadapter.Res([]User{}))
	users.GET("/search", handleSearchUsers,
		ginadapter.Query(
			openapi.QueryParam{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
			openapi.QueryParam{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
		),
		ginadapter.Res(struct{}{}),
	)
	users.POST("/users", handleCreateUser, ginadapter.Req(CreateUser{}), ginadapter.Created())
	users.POST("/users/upload", handleUploadUserFile,
		ginadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
		ginadapter.Res(map[string]string{}),
	)
	users.GET("/users/demo-errors", handleDemoErrors,
		ginadapter.Headers(openapi.HeaderParam{Name: "X-Demo-Fail", Type: openapi.ParamString, Required: false, Description: "Set to 400/401/500/503 to simulate an error"}),
		ginadapter.Res(map[string]string{}),
		ginadapter.Responses(
			openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}},
			openapi.ResponseSpec{Status: http.StatusUnauthorized, Schema: ErrorResponse{}},
			openapi.ResponseSpec{Status: http.StatusInternalServerError, Schema: ErrorResponse{}},
			openapi.ResponseSpec{Status: http.StatusServiceUnavailable, Schema: ErrorResponse{}},
		),
	)
	users.GET("/users/:id", handleGetUser, ginadapter.Res(User{}))
	users.PUT("/users/:id", handlePutUser, ginadapter.Req(UpdateUser{}), ginadapter.Res(User{}))
	users.PATCH("/users/:id", handlePatchUser, ginadapter.Req(UpdateUser{}), ginadapter.Res(User{}))
	users.DELETE("/users/:id", handleDeleteUser, ginadapter.NoContent())
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
