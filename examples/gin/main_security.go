//go:build gin && security && !typed

package main

import (
	"github.com/aizacoders/openapigo/adapters/ginadapter"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := ginadapter.New()

	cfg, bearer, apiKey := openAPICfgSecurity()

	registerSecureRoutes(r, bearer, apiKey)

	r.Docs(cfg)
	_ = r.Engine.Run(":8080")
}
