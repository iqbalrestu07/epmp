# EPMP Backend

Enterprise Property Management Platform — Backend API server.

---

## Tech Stack

| Concern         | Library / Tool                           | Notes                                            |
| --------------- | ---------------------------------------- | ------------------------------------------------ |
| Language        | Go 1.26                                  |                                                  |
| HTTP Framework  | [Echo v4](https://echo.labstack.com)     | Router, middleware, binding                      |
| Database Driver | pgx/v5                                   | pgxpool for connection pooling                   |
| Query Builder   | sqlc                                     | Type-safe SQL queries                            |
| Logging         | [zerolog](https://github.com/rs/zerolog) | Structured JSON in prod, coloured console in dev |
| Validation      | go-playground/validator v10              | Struct-tag based validation                      |
| Migration       | golang-migrate                           | SQL migration files in `migrations/`             |
| Testing         | testing + testify                        |                                                  |
| Architecture    | Clean Architecture + DDD                 |                                                  |

---

## Project Structure

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go           # Composition root — DI, Echo setup, graceful shutdown
│   └── migrate/
│       └── main.go           # Migration runner (up/down/version/force/create/watch)
├── internal/
│   ├── pkg/                # Shared packages across domains
│       ├── helper/               # Helper functions
│       ├── error/               # Error handling
│           ├── errors.go         # DomainError types + sentinel errors
│       ├── health/               # Health check
│           ├── health.go         # /health helper
│       ├── types/                # Shared value types
│       ├── logger/           # zerolog setup (New() zerolog.Logger)
│       ├── middleware/       # RequestLogger, Recover (zerolog-aware)
│       └── response/         # JSON envelope helpers (OK, Created, BadRequest …)
│   └── modules/              # Bounded contexts — one folder per domain
│       ├── property/         # Property module
│       ├── tenant/           # Tenant module
│       └── room/             # Room module
├── configs/                  # Config files (koanf / env)
├── migrations/               # SQL migration files (golang-migrate)
│   ├── 000001_create_properties.up.sql
│   ├── 000001_create_properties.down.sql
│   ├── 000002_create_tenants.up.sql
│   ├── 000002_create_tenants.down.sql
│   ├── 000003_create_rooms.up.sql
│   ├── 000003_create_rooms.down.sql
│   └── README.md             # Migration guide
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
├── entity/                   # Domain entity structs
├── dto/                      # Request / response DTOs
├── repository/               # PostgreSQL adapter (sqlc-backed)
│   └── repository.go         # Repository interface
│   └── repository_implementation.go # Repository implementation
├── service/
│   └── service.go            # Use-case orchestration interface
│   └── service_implementation.go # Use-case implementation
└── delivery/
    └── http/                 # Echo handler REST API
    └── route/                # Route registration
```

---

## Middleware

Middleware is registered globally in `main.go`.

| Middleware      | Package                      | Behaviour                                      |
| --------------- | ---------------------------- | ---------------------------------------------- |
| `RequestLogger` | `internal/shared/middleware` | Logs method, path, status, latency via zerolog |
| `Recover`       | `internal/shared/middleware` | Catches panics, logs via zerolog, returns 500  |

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

| Variable       | Default                                                            | Description                 |
| -------------- | ------------------------------------------------------------------ | --------------------------- |
| `APP_ENV`      | `development`                                                      | `production` → JSON logging |
| `PORT`         | `8080`                                                             | HTTP listen port            |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable` | PostgreSQL DSN              |

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

Migrations are managed by `golang-migrate` and run via `cmd/migrate`.

### Apply / rollback

```bash
# Apply all pending migrations
go run ./cmd/migrate up

# Rollback last migration
go run ./cmd/migrate down

# Rollback N migrations
go run ./cmd/migrate down 2

# Show current schema version
go run ./cmd/migrate version

# Fix dirty state after a failed migration
go run ./cmd/migrate force <version>
```

### Add a new migration

```bash
# Scaffold a numbered migration pair
go run ./cmd/migrate create add_status_to_properties
# → migrations/000004_add_status_to_properties.up.sql
# → migrations/000004_add_status_to_properties.down.sql

# Fill in the SQL, then apply
go run ./cmd/migrate up
```

### Watch mode (auto-apply)

```bash
go run ./cmd/migrate watch
```

Polls `migrations/` every 3 seconds. When new `.up.sql` files appear (e.g. after running the EPMP codegen), it automatically applies all pending migrations. Stop with `Ctrl+C`.

### Makefile shortcuts

```bash
make migrate-up
make migrate-down
make migrate-create name=add_column_to_tenants
```

See [`migrations/README.md`](migrations/README.md) for the full migration guide.

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
