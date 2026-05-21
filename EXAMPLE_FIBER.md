# Fiber example (OpenAPIGO)

Fiber example uses route-level OpenAPI options for the shortest setup.

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
    fiberadapter "github.com/aizacoders/openapigo/adapters/fiberadapter"
    "github.com/aizacoders/openapigo/openapi"
)
```

2) Create Fiber app and wrap with adapter

```go
app := fiberlib.New()
r := fiberadapter.Wrap(app)
```

3) Register handlers with short OpenAPI options

```go
users := r.Group("", fiberadapter.Tags("Users"))
users.GET("/users", func(c *fiberlib.Ctx) error {
    return fiberadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
}, fiberadapter.Res([]User{}))

users.POST("/users", createUser,
    fiberadapter.Req(CreateUser{}),
    fiberadapter.Res(User{}),
    fiberadapter.Created(),
)

users.POST("/users/upload", uploadUserFile,
    fiberadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
    fiberadapter.Res(map[string]string{}),
)
```

4) Mount OpenAPI and run

```go
r.Docs(openapi.Config{Title: "User API", Version: "1.0.0"})
r.App.Listen(":8080")
```

5) Notes

- `Wrap` lets you configure middleware and settings on the Fiber app before wrapping it with the adapter.
- Use `fiberadapter.MultipartUpload` to expose file upload in Swagger UI.

### Note about core router

The OpenAPIGO core router is a lightweight net/http-backed mux. Adapter packages (including Fiber) integrate with the core behavior and continue to work as before. If you use the `httprouter` adapter you can optionally mount the router automatically onto a `*http.ServeMux` by calling `httprouter.New(mux)`.
