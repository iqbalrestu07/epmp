---
trigger: model_decision
description: When writing tests, organizing test files, implementing test doubles, or setting up testing infrastructure
---

## Testing Strategy

### Test Pyramid

**Unit Tests (70% of tests):**

- **What:** Test domain logic in isolation with mocked dependencies
- **Speed:** Fast (<100ms per test)
- **Scope:** Single function, class, or module
- **Dependencies:** All external dependencies mocked (repositories, APIs, time, random)
- **Coverage Goal:** >85% of domain logic

**Integration Tests (20% of tests):**

- **What:** Test adapters against real infrastructure
- **Speed:** Medium (100ms-5s per test)
- **Scope:** Component interaction with infrastructure (database, cache, message queue)
- **Dependencies:** Real infrastructure (local Docker PostgreSQL via `internal/testutil`)
- **Coverage Goal:** All adapter implementations, critical integration points

**End-to-End Tests (10% of tests):**

- **What:** Test complete user journeys through all layers
- **Speed:** Slow (5s-30s per test)
- **Scope:** Full system from HTTP request to database and back
- **Dependencies:** Entire system running (or close approximation)
- **Coverage Goal:** Happy paths, critical business flows

### Test-Driven Development (TDD)

**Red-Green-Refactor Cycle:**

1. **Red:** Write a failing test for next bit of functionality
2. **Green:** Write minimal code to make test pass
3. **Refactor:** Clean up code while keeping tests green
4. **Repeat:** Next test

**Benefits:**

- Tests written first ensure testable design
- Comprehensive test coverage (code without test doesn't exist)
- Faster development (catch bugs immediately, not in QA)
- Better design (forces thinking about interfaces before implementation)

### Test Doubles Strategy

**Unit Tests:** Use mocks/stubs for all driven ports

- Mock repositories return pre-defined data
- Mock external APIs return successful responses
- Mock time/random for deterministic tests
- Control test environment completely

**Integration Tests:** Use real infrastructure

- Real PostgreSQL 16 from `docker compose up -d postgres`, migrated with `go run ./cmd/migrate up`
- Test actual database queries, connection handling, transactions
- Verify adapter implementations work with real services

> **EPMP override:** Integration tests do **not** use Testcontainers. `backend/internal/testutil` (`testutil.Setup(t)`) boots the full Echo router on `httptest` against `DATABASE_URL` (default `postgres://postgres:postgres@localhost:5432/epmp`) and **skips** when the DB is unreachable. Each module ships `<name>_api_test.go` using `testutil.AssertCRUD` (CRUD + cross-org isolation) and `testutil.AssertUnauthorized`. A skipped integration test is not a pass — start Postgres before verifying adapter changes.

**Best Practice:**

- Generate at least 2 implementations per driven port:
  1. Production adapter (PostgreSQL, GCP GCS, etc.)
  2. Test adapter (in-memory, fake implementation)

### Test Organization

**Universal Rule: Co-locate implementation tests; Separate system tests.**

**1. Unit & Integration Tests (Co-located)**

- **Rule:** Place tests **next to the file** they test.
- **Why:** Keeps tests visible, encourages maintenance, and supports refactoring (moving a file moves its tests).
- **Naming Convention (this project):**
  - **Go:** `<file>_test.go` (unit, same package or `_test` package), `<name>_api_test.go` (API integration via `testutil`)
  - **TS:** no unit runner configured yet; if added (ADR required) use `*.test.ts` consistently
    > You must strictly follow the convention for the target language. Do not mix `test` and `spec` suffixes in the same application context.
    > See `architectural-pattern.md` § Test co-location and `epmp-module-patterns.md` for the authoritative rule.

**2. End-to-End Tests (Separate)**

- **Rule:** E2E lives outside feature folders. In EPMP the Playwright runner is `frontend/run_e2e_tests.cjs`, executed via `make test-e2e` (headless) or `make test-e2e-gui` (headed). It uses the system Google Chrome (`channel: 'chrome'`) — no `npx playwright install`.
- **Registration:** every new page route MUST be appended to `PAGES_TO_TEST` in `run_e2e_tests.cjs`; business-flow steps extend the existing flow section.
- **Preconditions:** Postgres up + migrated, backend on `:8080`, frontend on `:3000`.
- **Pass criterion:** the run ends with `0 Errors Found` (console errors, failed requests, and uncaught exceptions all count).
- Details and policy: `E2E_TESTING.md`.

**Using a browser interactively during development/verification:**

When a UI change needs manual confirmation before/after the automated run, use whatever browser tooling your agent environment provides (Playwright MCP, `browser_preview`, or a headed run with `make test-e2e-gui`) against `http://localhost:3000`. Capture a screenshot/snapshot at each major step as proof in the task summary.

**E2E Test Requirements:**

- Every new/changed route appears in `PAGES_TO_TEST`
- Test happy path AND at least one error path for changed business flows
- Use unique identifiers (`Date.now()`) for created records; never depend on pre-existing data
- Never weaken the runner's error detection to make a run pass

**Key Principles:**

- **Unit/Integration tests**: Co-located with implementation
- **E2E tests**: Separate runner (crosses boundaries)
- **Test doubles**: Co-located with interface (hand-written fake in the test file)
- **Pattern consistency**: All features follow same structure

### Test Quality Standards

**AAA Pattern (Arrange-Act-Assert):**

```
// Arrange: Set up test data and mocks
const user = { id: '123', email: 'test@example.com' };
const mockRepo = createMockRepository();

// Act: Execute the code under test
const result = await userService.createUser(user);

// Assert: Verify expected outcome
expect(result.id).toBe('123');
expect(mockRepo.save).toHaveBeenCalledWith(user);
```

**Test Naming:**

- Descriptive: `should [expected behavior] when [condition]`
- Examples:
  - `should return 404 when user not found`
  - `should hash password before saving to database`
  - `should reject email with invalid format`

**Coverage Requirements:**

- Unit tests: >85% code coverage
- Integration tests: All adapter implementations
- E2E tests: Critical user journeys

### Related Principles

- Architectural Patterns — Testability-First Design @architectural-pattern.md
- Error Handling Principles @error-handling-principles.md
- Project Structure @project-structure.md
