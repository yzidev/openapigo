# Goas 

[![CI](https://github.com/yzidev/goas/actions/workflows/ci.yml/badge.svg)](https://github.com/yzidev/goas/actions/workflows/ci.yml)

Auto-generate **OpenAPI 3.x** from your Go route registrations.

The goal is to keep your routing code **clean** (plain `GET/POST/PUT/PATCH/DELETE`) while still producing a good OpenAPI spec + Swagger UI.

---

Background and motivation

Creating OpenAPI (OpenAPI 3.x) documentation for Go projects is often tedious and error-prone. Most common workflows require hand-maintaining large YAML or JSON files that declare the entire API surface — types, request/response schemas, parameters, security schemes, and more. For medium to large APIs this quickly becomes unmanageable: teams may end up with thousands of lines of YAML (10k+ lines is not unusual) that must be edited and kept in sync with code changes.

Every change to a handler, request/response type, or parameter often means manually editing the documentation files. This duplication increases the risk of inconsistencies, stale docs, and significant maintenance overhead. Compared to frameworks like Spring Boot or FastAPI — which offer more integrated or declarative approaches for keeping API docs close to code — the Go ecosystem has historically lacked a lightweight, ergonomic solution for automatic OpenAPI generation.

Goas was created to bridge that gap. Instead of writing a YAML entry for every endpoint, Goas captures route registrations and a small, config-first specification to generate a complete OpenAPI document and Swagger UI automatically. The goals are:

- Eliminate the need to maintain huge, hand-written OpenAPI YAML files.
- Keep handlers idiomatic and minimal while centralizing schema metadata in a compact config.
- Reduce duplication and human error by generating docs from the same source of truth as your routes.
- Provide practical features teams need (multipart uploads, security schemes, grouped tags) so large APIs remain maintainable and well-documented.

This README continues with the features you get and examples on how to use the library.

---

## What you get

- `GET /openapi.json` (generated OpenAPI document)
- Swagger UI mounted at:
  - `http://localhost:8080/swagger-ui/index.html#/`
  - `/swagger` is kept as a legacy redirect

---

## Key concepts

### 1) Base router (net/http)

Use the built-in router:

- `openapi.New(...)` → returns an `http.Handler` (lightweight net/http-backed router)
- register routes with `GET/POST/PUT/PATCH/DELETE`
- call `Docs()` once after registering your routes to mount `/openapi.json` and Swagger UI

Note: the default router implementation used to be chi-backed; it now uses a small net/http-based mux compatible with the project's needs. Adapters for Gin, Echo and Fiber remain available.

### 2) Config-first spec (SpringBoot-like)

Go handlers don’t expose schema information automatically.
So Goas uses a **config-first** approach:

- put route schemas/tags/security/query/header params in one place using `openapi/spec`
- keep your handlers clean and readable

### 3) Multipart upload support

Use `MultipartUpload(...)` to get `multipart/form-data` request bodies and a file upload field in Swagger UI.

---

## Installation

```bash
go get github.com/yzidev/goas@latest
```

---

## Minimal example (net/http)

```go
package main

import (
	"net/http"

	"github.com/yzidev/goas/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := openapi.New(openapi.Config{Title: "User API", Version: "1.0.0"})

	r.GET("/users", func(w http.ResponseWriter, _ *http.Request) {
		openapi.JSON(w, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	}, openapi.Res([]User{}), openapi.Tags("Users"))

	r.Docs()
	_ = http.ListenAndServe(":8080", r)
}
```

Prefer grouped config instead of per-route options? `openapi/spec` is still available as an advanced config-first layer:

```go
b := spec.New()
b.GroupTags("", []string{"Users"}, func(s *spec.SpecBuilder) {
	s.GET("/users").Res([]User{}).OK()
})

base := openapi.New(openapi.Config{Title: "User API", Version: "1.0.0"})
r := spec.HTTP(base, b.Spec())
r.GET("/users", listUsers)
base.Docs()
```

---

## Springdoc-like mode

For Gin, Echo, and Fiber, you can register routes with the framework directly and
add Goas with one docs call. Goas discovers method/path information from
the framework router and mounts Swagger UI.

```go
engine := gin.Default()

engine.GET("/users/:id", getUser)
engine.POST("/users", createUser)

ginadapter.Docs(engine, openapi.Config{
	Title:   "User API",
	Version: "1.0.0",
})

_ = engine.Run(":8080")
```

This zero-config mode documents paths, methods, path params, and default
responses. Go handlers do not expose request/response body types at runtime, so
body schemas are intentionally opt-in:

```go
r := ginadapter.Wrap(engine)
r.POST("/users", createUser,
	ginadapter.Req(CreateUser{}),
	ginadapter.Res(User{}),
	ginadapter.Created(),
)
r.Docs(openapi.Config{Title: "User API", Version: "1.0.0"})
```

For the built-in net/http router or `muxadapter`, use `openapi.New(...)` or
`muxadapter.Mount(...)`:

```go
mux := http.NewServeMux()
r := muxadapter.Mount(mux, openapi.Config{
	Title:   "User API",
	Version: "1.0.0",
})

r.GET("/users/{id}", getUser, openapi.Res(User{}), openapi.Tags("Users"))

_ = http.ListenAndServe(":8080", mux)
```

Note: the standard `http.ServeMux` does not expose registered route metadata, so
Goas can auto-document routes registered through `openapi.Router` or
`muxadapter.Router`, but not arbitrary handlers registered directly on
`http.ServeMux`.

---

## Multipart upload example

On a route:

```go
r.POST("/users/upload", uploadUserFile,
	openapi.MultipartUpload(
		"file",
		openapi.MultipartField{Name: "note", Type: openapi.ParamString},
	),
	openapi.Res(map[string]string{}),
)
```

In `openapi/spec`, the equivalent is:

```go
s.POST("/users/upload").MultipartUpload(
	"file",
	openapi.MultipartField{Name: "note", Type: openapi.ParamString},
).Res(map[string]string{}).OK()
```

In Swagger UI this will show:
- `file` as file chooser
- `note` as text input
- requestBody content type: `multipart/form-data`

---

## Security

You can provide security schemes via `openapi.Config.SecuritySchemes`.
For Springdoc-like auto-docs, set `openapi.Config.Security` as a global requirement.
For route-specific security, attach requirements per route with `openapi.Security(...)` or adapter aliases like `ginadapter.Security(...)`.
Examples include two schemes:

- **Bearer** JWT (`Authorization: Bearer <token>`)
- **API key** (`X-API-Key: <key>`)

---

## Examples (recommended)

Run examples and open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

### Default (net/http)

- Docs: [`EXAMPLE_HTTPROUTER.md`](examples/httprouter/EXAMPLE_HTTPROUTER.mdMPLE_HTTPROUTER.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

### Gin

- Docs: [`EXAMPLE_GIN.md`](examples/gin/EXAMPLE_GIN.mdMPLE_GIN.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

### Echo

- Docs: [`EXAMPLE_ECHO.md`](examples/echo/EXAMPLE_ECHO.mdMPLE_ECHO.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

### Fiber

- Docs: [`EXAMPLE_FIBER.md`](examples/fiber/EXAMPLE_FIBER.mdMPLE_FIBER.md)
  (See the doc above for run commands, endpoints, security, and upload sample.)

---

## Current support (today)

Goas is currently focused on **4 frameworks/router setups**:

1. **net/http (built-in `openapi.Router` based on chi)**
2. **Gin**
3. **Echo**
4. **Fiber**

Notes:
- Other frameworks may be added later, but the repo intentionally stays small and dependency-light.
- Adapters are provided as packages under `adapters/*` so you can use them when needed.
  They are compiled by default and no special build tags are required to use them.
  If you prefer to keep adapter dependencies optional for your project, consider
  shipping adapters as separate modules (e.g. `github.com/yzidev/goas-adapter-gin`) so downstream projects opt-in.

---

## Roadmap / future updates

The direction going forward:

- **Keep the public API simple**:
  - common HTTP methods only: `GET/POST/PUT/PATCH/DELETE`
  - grouping via `Group(...)`
  - OpenAPI metadata via route options or config-first spec (`openapi/spec`)

- **Improve schema inference gradually**:
  - better tag support (`omitempty`, pointer handling)
  - better nested struct handling
  - better multipart documentation

- **Better DX in Swagger UI**:
  - theming improvements
  - cleaner auth UX
  - consistent error schemas

- **Adapter expansion (optional)**:
  - If more frameworks are added, they will follow the same pattern:
    - keep handlers/framework usage idiomatic
    - keep Goas integration minimal
    - keep core library independent of adapter dependencies

### Update policy / compatibility

- The project is evolving quickly.
- We aim to keep the **core API stable** (`openapi.New`, `openapi.Router`, `openapi.Register`, and `openapi/spec`).
- Adapter APIs may change as we simplify integration and keep parity across frameworks.

### Framework support timeline

For now Goas only ships examples + adapters for:
- `net/http` (built-in router)
- Gin
- Echo
- Fiber

Additional frameworks are considered **future work** (optional adapters behind build tags).

### How to add another framework (adapter concept)

If you want to support another framework, the recommended approach is:

- Create a new adapter package under `adapters/<framework>`.
- Guard it with a build tag (so the dependency stays optional).
- The adapter should expose a router wrapper similar to the existing ones:
  - register `GET/POST/PUT/PATCH/DELETE`
  - keep grouping if the framework supports groups
  - call `openapi.Router.Handle(...)` / attach `HandlerOption`s in the same way.

For a starting point, check:
- `adapters/gin`
- `adapters/echo`
- `adapters/fiber`

---

## Adapters (how to use with frameworks)

Goas provides lightweight adapters for multiple frameworks so you can keep your
handler code clean while still generating OpenAPI and mounting Swagger UI.

Pattern (recommended):

1. Create your framework engine/app (e.g., `gin`, `echo`, `fiber`).
2. For zero-config route discovery, call `ginadapter.Docs(engine, cfg)`, `echoadapter.Docs(e, cfg)`, or `fiberadapter.Docs(app, cfg)`. For net/http, use `openapi.New(...)` or `muxadapter.Mount(...)`.
3. If you want body schemas/tags/security per route, wrap with `Wrap`/`New` and register routes with short options like `Res`, `Req`, `Tags`, and `Created`.
4. Run the engine/app.

Examples:

- Gin

```go
import (
    ginlib "github.com/gin-gonic/gin"
    "github.com/yzidev/goas/openapi"
    "github.com/yzidev/goas/adapters/ginadapter"
)

engine := ginlib.New()
engine.GET("/users", listUsers)
ginadapter.Docs(engine, openapi.Config{Title: "My API", Version: "0.1.0"})
engine.Run(":8080")
```

- Echo

```go
import (
    echolib "github.com/labstack/echo/v4"
    "github.com/yzidev/goas/openapi"
    "github.com/yzidev/goas/adapters/echoadapter"
)

base := echolib.New()
base.GET("/users", listUsers)
echoadapter.Docs(base, openapi.Config{Title: "My API", Version: "0.1.0"})
base.Start(":8080")
```

- Fiber

```go
import (
    fiberlib "github.com/gofiber/fiber/v2"
    "github.com/yzidev/goas/openapi"
    "github.com/yzidev/goas/adapters/fiberadapter"
)

app := fiberlib.New()
app.Get("/users", listUsers)
fiberadapter.Docs(app, openapi.Config{Title: "My API", Version: "0.1.0"})
app.Listen(":8080")
```

Notes:
- The `Docs` helpers discover native framework routes. The `Wrap` helpers add richer per-route metadata when you need schemas, tags, security, or custom responses.
- If you previously built with `-tags`, adapters are now compiled by default — no need to use build tags to get adapter implementations.

---

## License

MIT. See [`LICENSE`](./LICENSE).
