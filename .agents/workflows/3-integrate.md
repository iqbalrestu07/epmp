---
description: Integrate phase - test adapters with real infrastructure
---

# Phase 3: Integrate

## Purpose

Test adapter implementations (database, external APIs) with real infrastructure via the `backend/internal/testutil` harness against local PostgreSQL.

## Prerequisites

- Phase 2 (Implement) completed
- Unit tests passing
- Local PostgreSQL running: `docker compose up -d postgres` and migrated: `cd backend && go run ./cmd/migrate up`

## Applicability (Default: REQUIRED)

This phase is REQUIRED by default. You may ONLY skip it if ALL of the following are true:

- [ ] No storage/repository files were modified or created (`backend/internal/modules/*/repository/*_impl.go`)
- [ ] No external API client files were modified or created
- [ ] No database queries or schemas were changed (`backend/internal/database/**`, `backend/migrations/*.sql`)
- [ ] No message queue, cache, or I/O adapter code was touched

If ANY adapter was modified, you MUST write integration tests.

External adapters (`backend/internal/modules/communication/**`, e.g. WhatsApp via whatsmeow) must be faked in tests — never hit real WhatsApp.

## If This Phase Fails

If integration tests fail:

1. Check PostgreSQL logs for errors (`docker compose logs postgres`)
2. Verify schema matches expectations (migrations applied)
3. Fix adapter implementation
4. Re-run tests before proceeding

## Steps

### 1. Setup the testutil Harness

`testutil.Setup(t)` boots the full Echo router on `httptest` against `DATABASE_URL` and **skips** the test when the DB is unreachable — a skipped test is not a pass.

**Go Example:**

```go
func TestThingsAPI_CRUD(t *testing.T) {
    env := testutil.Setup(t) // skips if DATABASE_URL unreachable; env already has a user + org

    testutil.AssertCRUD(t, env, "/api/v1/things",
        `{"name":"Thing A"}`,   // create body
        `{"name":"Thing A2"}`,  // update body
    )
}
```

### 2. Write Integration Tests

Test file naming: co-located `*_test.go` (each module ships `<name>_api_test.go`)

```go
func TestThingsAPI_Unauthorized(t *testing.T) {
    env := testutil.Setup(t)
    testutil.AssertUnauthorized(t, env, "/api/v1/things")
}
```

Use `env.Do(t, method, path, body, token, orgID)` for raw requests, `env.Create(t, path, body)` to seed a record (returns its ID), `env.CreateOrg(t)` to create a second org for isolation checks, and the `*Chain` helpers (`RoomChain`, `ContractChain`, `InvoiceChain`, `PaymentChain`) for prerequisite fixtures.

### 3. Run Integration Tests

```bash
# Go - run all tests including integration
cd backend && go test -v ./...
```

### 4. Manual Check (Optional)

- If UI involved, launch a browser to verify the basic flow.
- If API involved, use `curl` or `client` to hit the endpoint.

## Completion Criteria

- [ ] Integration tests written for all adapters
- [ ] Tests pass with real infrastructure (`internal/testutil` + local PostgreSQL)
- [ ] Database queries verified against real PostgreSQL

## Next Phase

Proceed to **Phase 4: Verify** (`/4-verify`)
