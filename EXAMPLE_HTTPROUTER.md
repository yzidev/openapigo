# net/http (default) router example (OpenAPIGO)

The "default" router in this repo is `openapi.Router` (a lightweight net/http-based mux).

## Quick start

Run the example:

```bash
go run ./examples/httprouter
```

Use `-tags "security"` only when running the security variant:

```bash
go run -tags "security" ./examples/httprouter
```

Open Swagger UI:
- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:
- http://localhost:8080/openapi.json

---

## Implementation details (step-by-step)

This section shows how to wire the default HTTP router with OpenAPIGO in your project.

1) Imports

```go
import (
    "net/http"

    "github.com/aizacoders/openapigo/adapters/muxadapter"
    "github.com/aizacoders/openapigo/openapi"
)
```

2) Create the router

```go
mux := http.NewServeMux()
r := muxadapter.New(mux)
```

3) Register handlers using route options

```go
r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
    openapi.JSON(w, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
}, openapi.Res([]User{}), openapi.Tags("Users"))

r.POST("/users", func(w http.ResponseWriter, req *http.Request) {
    var in CreateUser
    if err := openapi.Bind(req, &in); err != nil || in.Name == "" {
        openapi.JSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
        return
    }
    w.WriteHeader(http.StatusCreated)
}, openapi.Req(CreateUser{}), openapi.Res(User{}), openapi.Created(), openapi.Tags("Users"))

r.POST("/users/upload", uploadUserFile,
    openapi.MultipartUpload("file", openapi.MultipartField{Name: "note", Type: openapi.ParamString}),
    openapi.Res(map[string]string{}),
    openapi.Tags("Users"),
)
```

4) Mount OpenAPI JSON + Swagger UI on the router and run

```go
r.Docs(openapi.Config{Title: "User API", Version: "1.0.0"})
_ = http.ListenAndServe(":8080", mux)
```

5) Security (optional)

- Define schemes in `openapi.Config.SecuritySchemes` and attach per-route via `openapi.Security(...)`.

6) Multipart uploads

- Use `openapi.MultipartUpload(...)` to declare a `multipart/form-data` body with a file field; the Swagger UI will render a file chooser and corresponding fields.

---

## Notes

- Examples follow the pattern: build a router/engine, register routes with short OpenAPI options, call `Docs(...)`, then run.
