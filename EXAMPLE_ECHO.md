# Echo example (OpenAPIGO)

This example uses route-level OpenAPI options so Echo routes stay direct and easy to scan.

## Quick start

Install Echo if you don't have it:

```bash
go get github.com/labstack/echo/v4@latest
```

Run the example:

```bash
go run ./examples/echo
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./examples/echo
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
    echolib "github.com/labstack/echo/v4"
    echoadapter "github.com/aizacoders/openapigo/adapters/echoadapter"
    "github.com/aizacoders/openapigo/openapi"
)
```

2) Create Echo instance and wrap with adapter

```go
base := echolib.New()
r := echoadapter.Wrap(base)
```

3) Register handlers with short OpenAPI options

```go
users := r.Group("", echoadapter.Tags("Users"))
users.GET("/users", func(c echolib.Context) error {
    return echoadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
}, echoadapter.Res([]User{}))

users.POST("/users", createUser,
    echoadapter.Req(CreateUser{}),
    echoadapter.Res(User{}),
    echoadapter.Created(),
)

users.POST("/users/upload", uploadUserFile,
    echoadapter.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
    echoadapter.Res(map[string]string{}),
)
```

4) Mount OpenAPI and run

```go
r.Docs(openapi.Config{Title: "User API", Version: "1.0.0"})
r.Echo.Start(":8080")
```

5) Notes

- `Wrap` lets you create middleware and configure the Echo instance before wrapping it with the adapter.
- Use `echoadapter.MultipartUpload` to expose file upload inputs in Swagger UI.

### Note about core router

The OpenAPIGO core router is a lightweight net/http-backed mux. Adapter packages (including Echo) integrate with this core behavior and continue to work as before. If you use the `httprouter` adapter you can optionally mount the router automatically onto a `*http.ServeMux` by calling `httprouter.New(mux)`.
