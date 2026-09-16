---
description: Verify phase - run full validation suite
---

# Phase 4: Verify

## Purpose

Run all linters, static analysis, and tests to ensure code quality.

> **Note:** Paths below follow the project structure defined in `project-structure.md`.

### Phase 3 Gate

- [ ] Were any storage/database adapter files modified? List them: \_\_\_
- [ ] Were any external API adapters modified? List them: \_\_\_
- [ ] If YES to either: are integration tests written and passing?
- [ ] If NO to both: state why Phase 3 was skipped

> CRITICAL: If you modified `backend/internal/modules/*/repository/*_impl.go`, `backend/internal/database/**`, `backend/migrations/*.sql`, or any SQL query, Phase 3 is MANDATORY. Do not proceed without integration tests.

### Phase 3.5 Gate

- [ ] Were any UI changes or modified? List them: \_\_\_
- [ ] If YES: are E2E tests written and passing?
- [ ] If NO to both: state why Phase 3.5 was skipped

> CRITICAL: If you modified `frontend/*` or any frontend component and integration point, Phase 3.5 is MANDATORY. Do not proceed without E2E tests.

## If This Phase Fails

If lint/test/build fails:

1. **Do not proceed** to Ship
2. Fix the issue (go back to Phase 2 or 3 as needed)
3. Re-run full verification
4. Only proceed when ALL checks pass

## Steps

**Set Mode:** Use `task_boundary` to set mode to **VERIFICATION**.

### 1. Backend Validation

Run the FULL validation suite for the backend path as defined in `project-structure.md`:

```bash
cd backend && gofmt -l . && go vet ./... && go build ./... && go test ./... -count=1
```

### 2. Frontend Validation

```bash
cd frontend && npm run lint && npm run build
```

(`npm run build` runs `tsc -b && vite build` — that is the typecheck. No frontend unit test runner exists.)

### 3. Build Check

```bash
# Backend
cd backend && go build ./...

# Frontend
cd frontend && npm run build
```

### 4. Check Coverage

Report actual coverage in task summary.

**Go:**

```bash
cd backend && go test -cover ./internal/modules/...
```

**Frontend:** no unit test runner exists — coverage is not collected; validate via `npm run lint && npm run build` and E2E.

## Completion Criteria

- [ ] All lint checks pass
- [ ] All tests pass
- [ ] Build succeeds
- [ ] Coverage reported (target >85% on domain logic)

## On Success

Mark task as `[x]` in task.md (verification passed = task complete).

## Next Phase

Proceed to **Phase 5: Ship** (`/5-commit`)
