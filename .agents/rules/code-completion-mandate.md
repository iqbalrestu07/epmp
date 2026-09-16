---
trigger: always_on
---

## Code Completion Mandate

### Universal Requirement

**Before marking any code task as complete, you MUST run automated quality checks and remediate all issues.**

This is NOT OPTIONAL. Delivering code without validation violates the Rugged Software Constitution @rugged-software-constitution.md.

### The Completion Checklist

Every code generation task follows this workflow:

1. **Generate** - Write the code based on requirements
2. **Validate** - Run language-appropriate quality checks (see below)
3. **Remediate** - Fix all detected issues
4. **Verify** - Re-run checks to confirm fixes
5. **Deliver** - Mark task complete only after all checks pass

**Never skip validation "to save time." Validation IS the work.**

### Quality Commands for This Project

Run the set that matches the code you touched. Paths are relative to the repo root.

| Scope                               | Check              | Command                                                                                           | Notes                                                                                                                        |
| ----------------------------------- | ------------------ | ------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Backend (`backend/`)                | Format             | `cd backend && gofmt -l .`                                                                        | Must print nothing                                                                                                           |
|                                     | Vet                | `cd backend && go vet ./...`                                                                      |                                                                                                                              |
|                                     | Lint               | `cd backend && staticcheck ./... && gosec -quiet ./...`                                           | Both installed in `~/go/bin`; `golangci-lint` optional                                                                       |
|                                     | Build              | `cd backend && go build ./...`                                                                    |                                                                                                                              |
|                                     | Tests              | `make test-backend` (= `go test ./... -v -count=1`)                                               | Integration tests need Postgres up + migrated; they `SKIP` otherwise — a skip is **not** a pass for adapter changes          |
|                                     | Race (before ship) | `cd backend && go test -race ./...`                                                               |                                                                                                                              |
| Frontend (`frontend/`)              | Typecheck + build  | `make test-frontend` (= `npm run build` → `tsc -b && vite build`)                                 | Zero TS errors                                                                                                               |
|                                     | Lint               | `cd frontend && npm run lint`                                                                     | `eslint .` — **no ESLint config is committed yet**; if it errors on missing config, report it, do not fabricate one silently |
|                                     | Dependency audit   | `cd frontend && npm audit --audit-level=high`                                                     |                                                                                                                              |
|                                     | E2E (UI touched)   | `make test-e2e`                                                                                   | Requires backend `:8080` + frontend `:3000` running; must end with 0 errors                                                  |
| Tools (`tools/epmp-sdk/*/codegen`)  | Tests              | `cd tools/epmp-sdk/be/codegen && go test ./...` (same for `fe/codegen`)                           | Run when templates or config change                                                                                          |
| Migrations                          | Round-trip         | `cd backend && go run ./cmd/migrate up && go run ./cmd/migrate down 1 && go run ./cmd/migrate up` | Local DB only                                                                                                                |
| Docker (compose/Dockerfile touched) | Build              | `docker compose build`                                                                            |                                                                                                                              |

### Failure Protocol

**If any quality check fails:**

1. Read the error output completely
2. Fix the identified issues in the code
3. Re-run the failing command
4. Do not proceed until all checks pass

> Never disable a lint rule, add `//nolint`, `// eslint-disable`, `@ts-ignore`, or `t.Skip` to make checks pass. Fix the root cause. Never delete a failing test.

### Related Principles

- Rugged Software Constitution @rugged-software-constitution.md
- Code Idioms and Conventions @code-idioms-and-conventions.md
- Go Idioms and Patterns @go-idioms-and-patterns.md
- TypeScript Idioms and Patterns @typescript-idioms-and-patterns.md
- EPMP Module Patterns @epmp-module-patterns.md
