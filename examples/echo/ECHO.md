# Echo example (Goas)

This example uses native Echo routes plus one Goas docs call.

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
    echoadapter "github.com/yzidev/goas/adapters/echoadapter"
    "github.com/yzidev/goas/openapi"
)
```

2) Create Echo instance

```go
base := echolib.New()
```

3) Register handlers with plain Echo

```go
users := base.Group("")
users.GET("/users", func(c echolib.Context) error {
    return echoadapter.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
})

users.POST("/users", createUser)
users.POST("/users/upload", uploadUserFile)
```

4) Mount OpenAPI and run

```go
echoadapter.Docs(base, openapi.Config{Title: "User API", Version: "1.0.0"})
base.Start(":8080")
```

5) Notes

- Auto-docs can discover paths, methods, and path params from native Echo routes.
- Use `echoadapter.Wrap` with `Req`, `Res`, `Tags`, or `MultipartUpload` when you want explicit body schemas or richer per-route docs.

### Note about core router

The Goas core router is a lightweight net/http-backed mux. Adapter packages (including Echo) integrate with this core behavior and continue to work as before. If you use the `httprouter` adapter you can optionally mount the router automatically onto a `*http.ServeMux` by calling `muxadapter.Mount(mux, cfg)`.
