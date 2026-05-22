//go:build gin && !typed && !security

package main

import (
	"log"

	ginlib "github.com/gin-gonic/gin"
	"github.com/yzidev/goas/adapters/ginadapter"
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

	registerSystemRoutes(engine)
	registerUserRoutes(engine)

	ginadapter.Docs(engine, openAPICfg())
	log.Fatal(engine.Run(":8080"))
}
