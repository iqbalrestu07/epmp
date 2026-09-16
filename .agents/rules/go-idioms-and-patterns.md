---
trigger: model_decision
description: When writing or reviewing Go code in backend/ or tools/ — idioms, error handling, context, testing, and tooling for this project
---

## Go Idioms and Patterns

> Applies to `backend/` (module `github.com/epmp/backend`, Go 1.26) and `tools/epmp-sdk/*/codegen`.
> Naming rules come from `tools/epmp-ai/CONVENTIONS.md`; layering rules from `epmp-docs/epmp-011.md`.

### Formatting & Tooling

| Tool | Purpose | Command |
|---|---|---|
| `gofmt` | Canonical formatting (non-negotiable) | `gofmt -l -w .` |
| `go vet` | Static correctness | `go vet ./...` |
| `staticcheck` | Lint (installed in `~/go/bin`) | `staticcheck ./...` |
| `gosec` | Security lint (installed in `~/go/bin`) | `gosec -quiet ./...` |
| `govulncheck` | Dependency CVE scan | `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` |
| `go test` | Tests | `go test ./... -count=1` (add `-race` before shipping) |

All commands run from `backend/`. `golangci-lint` is referenced by `backend/Makefile` but is **not** guaranteed to be installed — fall back to `go vet` + `staticcheck` + `gosec`.

### Errors

- Return errors, don't panic. `panic` only in `main`/bootstrap for unrecoverable startup failures.
- Wrap with context and `%w`: `fmt.Errorf("property service: create: %w", err)`. Prefix = `<package> <layer>: <operation>`.
- Business failures are `*errs.DomainError`; check with `errors.As` / `errs.IsDomainError(err, "NOT_FOUND")`, never string comparison.
- Sentinel errors are exported `Err*` vars describing the business fact: `ErrRoomAlreadyOccupied`, not `ErrBad`.
- Never swallow errors with `_ =` on operations whose failure matters (DB writes, scans that feed responses).

### Context

- First parameter of every service/repository method is `ctx context.Context`; pass `c.Request().Context()` from handlers.
- Never store `context.Context` in a struct; never pass `nil` — use `context.Background()` only at the composition root or in tests.

### Structs, Interfaces, Constructors

- Interfaces are defined by the **consumer** (repository interface lives in the module's `repository/` package, consumed by `service/`). Keep them small.
- Constructors `NewX(deps...) *X`; accept interfaces, return concrete types.
- Prefer value receivers for small immutable types, pointer receivers for entities and anything with mutex/pool.
- Exported identifiers need doc comments starting with the identifier name.
- DTO JSON tags are `snake_case` and mirror the frontend Zod schema field names.

### Database (pgx)

- Use `pgxpool.Pool` from `internal/database/postgres`; never open connections in modules.
- Parameterised queries only (`$1, $2`); build dynamic filters by appending `AND col = $%d` with an arg index, never by string-interpolating values.
- Always filter `deleted_at IS NULL` and `organization_id` (see @multi-tenancy-boundary.md).
- `pgx.ErrNoRows` → translate to `errs.ErrNotFound` at the repository boundary.
- Multi-statement writes go in a transaction (`pool.Begin(ctx)` + `defer tx.Rollback(ctx)`), owned by the service/use case, not the repository.

### HTTP (Echo)

- Handlers return `error`; use `internal/pkg/response` helpers for every response.
- Bind with `c.Bind(&req)`, validate required fields explicitly, then call the service. Do not put business rules in handlers.
- Path params via `c.Param("id")`; pagination via `page`/`per_page` with defaults `1`/`20`.
- Middleware order: `AuthRequired` → `AuditLog` → (`RequirePermission`) → handler.

### Concurrency

- Goroutines must have a clear owner and shutdown path (context cancellation or `sync.WaitGroup`). Background work (WhatsApp client, websocket hub) lives in `internal/pkg/*`, not in request handlers.
- Protect shared maps with `sync.RWMutex`; prefer channels for hand-off, mutexes for state.

### Testing

- Table-driven tests with `t.Run(name, ...)`; name tests `Test<Func>_<Scenario>` (e.g. `TestCancelReservation_RejectsPaidReservation`).
- Unit tests: hand-written fake implementing the repository interface (no mocking framework in `go.mod`; do not add one without an ADR).
- API integration tests: `internal/testutil` — `env := testutil.Setup(t)` boots the full router on `httptest` against real Postgres and skips when the DB is unreachable. Use `env.Create`, `env.List`, `env.RoomChain`, `env.ContractChain`, `testutil.AssertCRUD`, `testutil.AssertUnauthorized`.
- Use `t.Helper()` in helpers, `t.Cleanup` for teardown, unique fixture names via `time.Now().UnixNano()`.
- Never delete or `t.Skip` a failing test to make the build green.

### Migrations (golang-migrate)

- `cd backend && go run ./cmd/migrate create <snake_case_name>` generates `NNNNNN_<name>.up.sql` + `.down.sql`.
- `.down.sql` must fully reverse `.up.sql`. Use `IF NOT EXISTS` / `IF EXISTS` for idempotency.
- Never edit an already-applied migration; add a new one.

### Anti-Patterns

❌ `response.InternalError(c, err.Error())` — leaks SQL/internal details to clients
❌ `fmt.Sprintf("... WHERE name = '%s'", name)` — SQL injection
❌ `init()` functions with side effects, package-level mutable state
❌ Returning `interface{}`/`any` from services — return typed DTOs
❌ `time.Now()` deep inside domain logic without an injectable clock when the value affects business rules

### Related Principles
- EPMP Module Patterns @epmp-module-patterns.md
- Error Handling Principles @error-handling-principles.md
- Testing Strategy @testing-strategy.md
- Code Completion Mandate @code-completion-mandate.md
