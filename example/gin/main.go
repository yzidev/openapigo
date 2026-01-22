//go:build gin && !typed && !security

package main

import (
	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi/simple"
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

	sr := simple.NewGin(r, springSpec())

	registerSystemRoutes(sr)
	registerUserRoutes(sr)

	gin.Register(r, openAPICfg())
	_ = r.Engine.Run(":8080")
}
