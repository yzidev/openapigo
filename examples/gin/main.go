//go:build gin && !typed && !security

package main

import (
	"github.com/aizacoders/openapigo/adapters/ginadapter"
	ginlib "github.com/gin-gonic/gin"
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
	engine := ginlib.New()

	r := ginadapter.Wrap(engine)

	registerSystemRoutes(r)
	registerUserRoutes(r)

	r.Docs(openAPICfg())
	_ = r.Engine.Run(":8080")
}
