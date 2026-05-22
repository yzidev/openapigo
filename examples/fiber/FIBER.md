# Fiber example (Goas)

Fiber example uses native Fiber routes plus one Goas docs call.

## Quick start

Install Fiber if you don't have it:

```bash
go get github.com/gofiber/fiber/v2@latest
```

Run the example:

```bash
go run ./examples/fiber
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./examples/fiber
```

Open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:

- http://localhost:8080/openapi.json

---

## Implementation details (step-by-step)

1) Imports

```go
import (
    fiberlib "github.com/gofiber/fiber/v2"
    fiberadapter "github.com/yzidev/goas/adapters/fiberadapter"
    "github.com/yzidev/goas/openapi"
)
```

2) Create Fiber app

```go
app := fiberlib.New()
```

3) Register handlers with plain Fiber

```go
users := app.Group("")
users.Get("/users", func(c *fiberlib.Ctx) error {
    return fiberadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
})

users.Post("/users", createUser)
users.Post("/users/upload", uploadUserFile)
```

4) Mount OpenAPI and run

```go
fiberadapter.Docs(app, openapi.Config{Title: "User API", Version: "1.0.0"})
app.Listen(":8080")
```

5) Notes

- Auto-docs can discover paths, methods, and path params from native Fiber routes.
- Use `fiberadapter.Wrap` with `Req`, `Res`, `Tags`, or `MultipartUpload` when you want explicit body schemas or richer per-route docs.

### Note about core router

The Goas core router is a lightweight net/http-backed mux. Adapter packages (including Fiber) integrate with the core behavior and continue to work as before. If you use the `httprouter` adapter you can optionally mount the router automatically onto a `*http.ServeMux` by calling `muxadapter.Mount(mux, cfg)`.
