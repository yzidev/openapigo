# Gin example (Goas)

This example uses native Gin routes plus one Goas docs call, similar to a Springdoc-style setup.

## Quick start

Install (if you don't already have Gin):

```bash
go get github.com/gin-gonic/gin@latest
```

Run the example:

```bash
go run -tags "gin" ./examples/gin
```

Use `-tags "gin,security"` when running the security variant:

```bash
go run -tags "gin,security" ./examples/gin
```

Open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:

- http://localhost:8080/openapi.json

---

## Implementation details (step-by-step)

This section shows how to wire Gin with Goas in your own project.

1) Imports

```go
import (
    ginlib "github.com/gin-gonic/gin"
    ginadapter "github.com/yzidev/goas/adapters/ginadapter"
    "github.com/yzidev/goas"
)
```

2) Create your Gin engine (you can customize middleware, logger, etc.)

```go
engine := ginlib.Default()      // or ginlib.New()
```

3) Register handlers with plain Gin

```go
users := engine.Group("")

users.GET("/users", func(c *ginlib.Context) {
    ginadapter.JSON(c, 200, []User{{ID: "1", Name: "Alice"}})
})

users.POST("/users", createUser)
users.POST("/users/upload", uploadUserFile)
users.GET("/users/demo-errors", demoErrors)
```

4) Mount OpenAPI JSON + Swagger UI and run

```go
ginadapter.Docs(engine, goas.Config{Title: "User API", Version: "1.0.0"})
engine.Run(":8080")
```

5) Add richer schemas only when you need them

```go
r := ginadapter.Wrap(engine)
r.POST("/users", createUser,
    ginadapter.Req(CreateUser{}),
    ginadapter.Res(User{}),
    ginadapter.Created(),
)
```

6) Security (optional)

- Define schemes in `goas.Config.SecuritySchemes` and attach per-route via `ginadapter.Security(...)`.

Example:

```go
bearer := &goas.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}
apiKey := &goas.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}
cfg := goas.Config{Title: "API", Version: "1.0.0", SecuritySchemes: map[string]*goas.SecuritySchemeRef{
    "bearer": {Value: bearer},
    "xapikey": {Value: apiKey},
}}
bearerReq := goas.NewSecurityRequirement().Authenticate("bearer")
cfg.Security = goas.SecurityRequirements{bearerReq}
```

9) Troubleshooting

- If you get type errors around constructors: make sure you import Gin framework (github.com/gin-gonic/gin) and the adapter package separately (use aliases to avoid name collisions: `ginlib` vs `ginadapter`).
- Auto-docs can discover paths, methods, and path params from native Gin routes.
- If Swagger UI doesn't show request/response schemas: add `Req`/`Res` on routes that need explicit body schemas.

---

## What to inspect in this repo

- `examples/gin/main.go` — demonstrates native Gin routes with one docs call
- `examples/gin/routes.go` — shows clean route declarations
- `spec` — optional config-first builder for teams that prefer central route metadata

### Note about core router

The Goas core router is a lightweight net/http-backed mux. The Gin adapter continues to work unchanged. For the net/http example you can mount the router on a ServeMux easily with `muxadapter.Mount(mux, cfg)`.
