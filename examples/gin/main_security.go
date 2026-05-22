//go:build gin && security && !typed

package main

import (
	ginlib "github.com/gin-gonic/gin"
	"github.com/yzidev/goas/adapters/ginadapter"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	engine := ginlib.New()

	cfg := openAPICfgSecurity()

	registerSecureRoutes(engine)

	ginadapter.Docs(engine, cfg)
	_ = engine.Run(":8080")
}
