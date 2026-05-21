# Gin example (OpenAPIGO)

This example uses route-level OpenAPI options so routes, handlers, and docs stay in one place.

## Quick start

Install (if you don't already have Gin):

```bash
go get github.com/gin-gonic/gin@latest
```

Run the example:

```bash
go run ./examples/gin
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./examples/gin
```

Open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:

- http://localhost:8080/openapi.json

---

## Implementation details (step-by-step)

This section shows how to wire Gin with OpenAPIGO in your own project.

1) Imports

```go
import (
    ginlib "github.com/gin-gonic/gin"
    ginadapter "github.com/yzidev/openapigo/adapters/ginadapter"
    "github.com/yzidev/openapigo/openapi"
)
```

2) Create your Gin engine (you can customize middleware, logger, etc.)

```go
engine := ginlib.Default()      // or ginlib.New()
```

3) Wrap the engine with the adapter so OpenAPIGO can capture route metadata

```go
r := ginadapter.Wrap(engine)
```

4) Register handlers with short OpenAPI options

```go
users := r.Group("", ginadapter.Tags("Users"))

users.GET("/users", func(c *ginlib.Context) {
    ginadapter.JSON(c, 200, []User{{ID: "1", Name: "Alice"}})
}, ginadapter.Res([]User{}))

users.POST("/users", createUser,
    ginadapter.Req(CreateUser{}),
    ginadapter.Res(User{}),
    ginadapter.Created(),
)

users.POST("/users/upload", uploadUserFile,
    ginadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
    ginadapter.Res(map[string]string{}),
)

users.GET("/users/demo-errors", demoErrors,
    ginadapter.Res(map[string]string{}),
    ginadapter.Responses(
        openapi.ResponseSpec{Status: 400, Schema: openapi.ErrorResponse{}},
        openapi.ResponseSpec{Status: 401, Schema: openapi.ErrorResponse{}},
        openapi.ResponseSpec{Status: 500, Schema: openapi.ErrorResponse{}},
    ),
)
```

5) Mount OpenAPI JSON + Swagger UI and run

```go
r.Docs(openapi.Config{Title: "User API", Version: "1.0.0"})
r.Engine.Run(":8080")
```

6) Security (optional)

- Define schemes in `openapi.Config.SecuritySchemes` and attach per-route via `ginadapter.Security(...)`.

Example:

```go
bearer := &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}
apiKey := &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}
cfg := openapi.Config{Title: "API", Version: "1.0.0", SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
    "bearer": {Value: bearer},
    "xapikey": {Value: apiKey},
}}
// Attach per-route in builder:
// s.GET("/secure").Security(&openapi3.SecurityRequirement{"bearer": {}}).Res(...).OK()
```

9) Troubleshooting

- If you get type errors around constructors: make sure you import Gin framework (github.com/gin-gonic/gin) and the adapter package separately (use aliases to avoid name collisions: `ginlib` vs `ginadapter`).
- If Swagger UI doesn't show request/response schemas: ensure you declared `Req`/`Res` on the route.

---

## What to inspect in this repo

- `example/gin/main.go` — demonstrates wrapping an existing engine and registering OpenAPI
- `example/gin/routes.go` — shows clean route declarations
- `openapi/spec` — optional config-first builder for teams that prefer central route metadata

### Note about core router

The OpenAPIGO core router is a lightweight net/http-backed mux. The Gin adapter continues to work unchanged. For the net/http example you can mount the router on a ServeMux easily (the `httprouter` adapter supports `httprouter.New(mux)` to auto-mount).
