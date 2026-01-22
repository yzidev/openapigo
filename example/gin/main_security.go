//go:build gin && security && !typed

package main

import (
	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	cfg, bearer, apiKey := openAPICfgSecurity()

	sr := simple.NewGin(r, simple.Spec{})
	registerSecureRoutes(sr, bearer, apiKey)

	gin.Register(r, cfg)
	_ = r.Engine.Run(":8080")
}
