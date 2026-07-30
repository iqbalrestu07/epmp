# EPMP Backend

Enterprise Property Management Platform — Backend API server.

---

## Tech Stack

| Concern         | Library / Tool               | Notes                              |
| --------------- | ---------------------------- | ---------------------------------- |
| Language        | Go 1.26                      |                                    |
| HTTP Framework  | [Echo v4](https://echo.labstack.com) | Router, middleware, binding |
| Database Driver | pgx/v5                       | pgxpool for connection pooling     |
| Query Builder   | sqlc                         | Type-safe SQL queries              |
| Logging         | [zerolog](https://github.com/rs/zerolog) | Structured JSON in prod, coloured console in dev |
| Validation      | go-playground/validator v10  | Struct-tag based validation        |
| Migration       | golang-migrate               | SQL migration files in `migrations/` |
| Testing         | testing + testify            |                                    |
| Architecture    | Clean Architecture + DDD     |                                    |

---

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Composition root — DI, Echo setup, graceful shutdown
├── internal/
│   ├── shared/               # Cross-domain utilities
│   │   ├── errors.go         # DomainError types + sentinel errors
│   │   ├── health.go         # /health helper
│   │   ├── types.go          # Shared value types
│   │   ├── logger/           # zerolog setup (New() zerolog.Logger)
│   │   ├── middleware/       # RequestLogger, Recover (zerolog-aware)
│   │   └── response/         # JSON envelope helpers (OK, Created, BadRequest …)
│   └── modules/              # Bounded contexts — one folder per domain
│       ├── property/         # Property module
│       ├── tenant/           # Tenant module
│       └── room/             # Room module
├── configs/                  # Config files (koanf / env)
├── migrations/               # SQL migration files (golang-migrate)
└── go.mod
```

---

## Router Architecture

Echo is the **only** HTTP router. There is no `net/http.Server` — Echo's `e.Start()` is used directly with graceful shutdown.

### Route hierarchy

```
GET  /health                          # unauthenticated health check

POST   /api/v1/properties             # create
GET    /api/v1/properties             # list (paginated)
GET    /api/v1/properties/:id         # get by ID
PUT    /api/v1/properties/:id         # update
DELETE /api/v1/properties/:id         # delete

POST   /api/v1/tenants
GET    /api/v1/tenants
GET    /api/v1/tenants/:id
PUT    /api/v1/tenants/:id
DELETE /api/v1/tenants/:id

POST   /api/v1/rooms
GET    /api/v1/rooms
GET    /api/v1/rooms/:id
PUT    /api/v1/rooms/:id
DELETE /api/v1/rooms/:id
```

### Module registration pattern

Every bounded context implements `RegisterRoutes(v1 *echo.Group)`.  
`main.go` creates the shared `/api/v1` group and passes it to each module:

```go
v1 := e.Group("/api/v1")

property.NewModule(db, log).RegisterRoutes(v1)
tenant.NewModule(db, log).RegisterRoutes(v1)
room.NewModule(db, log).RegisterRoutes(v1)
```

---

## Bounded Context Structure

Each module follows the same layout:

```
{module}/
├── module.go                 # DI wiring; exposes NewModule() and RegisterRoutes()
├── MODULE.md                 # Module contract (fields, endpoints, behaviors)
├── application/
│   ├── dto/                  # Request / response DTOs
│   └── service/              # Use-case orchestration
├── domain/
│   ├── entity/               # Domain entity structs
│   └── repository/           # Repository interface (port)
├── infrastructure/
│   └── repository/           # PostgreSQL adapter (sqlc-backed)
└── interfaces/
    └── http/                 # Echo handler + route registration
```

---

## Middleware

Middleware is registered globally in `main.go`.

| Middleware       | Package                          | Behaviour                                   |
| ---------------- | -------------------------------- | ------------------------------------------- |
| `RequestLogger`  | `internal/shared/middleware`     | Logs method, path, status, latency via zerolog |
| `Recover`        | `internal/shared/middleware`     | Catches panics, logs via zerolog, returns 500 |

---

## Response Envelope

All endpoints return a consistent JSON shape:

**Success**
```json
{ "success": true, "data": { … } }
```

**Error**
```json
{ "success": false, "error": { "code": "NOT_FOUND", "message": "…" } }
```

---

## Environment Variables

Copy `.env.example` and adjust for your environment:

```bash
cp .env.example .env
```

| Variable       | Default                                                      | Description                  |
| -------------- | ------------------------------------------------------------ | ---------------------------- |
| `APP_ENV`      | `development`                                                | `production` → JSON logging  |
| `PORT`         | `8080`                                                       | HTTP listen port             |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable` | PostgreSQL DSN      |

---

## Running

```bash
# install dependencies
go mod tidy

# run development server
go run ./cmd/server/...

# build binary
go build -o bin/server ./cmd/server/...

# run binary
./bin/server
```

---

## Database Migrations

```bash
# apply all pending migrations
migrate -path ./migrations -database "$DATABASE_URL" up

# rollback last migration
migrate -path ./migrations -database "$DATABASE_URL" down 1
```

---

## Testing

```bash
go test ./... -v -count=1
```

---

## Code Generation

Modules are scaffolded by the EPMP SDK code generator.

```bash
cd tools/epmp-sdk/be/codegen
go build -o epmp-codegen ./cmd/codegen/

# generate a new module from a schema
./epmp-codegen generate --schema examples/property.yaml

# dry-run (print files without writing)
./epmp-codegen generate --schema examples/property.yaml --dry-run
```

See [`tools/epmp-sdk/be/codegen/README.md`](../tools/epmp-sdk/be/codegen/README.md) for full documentation.
