---
trigger: always_on
---

## Code Organization Principles

- Generate small, focused functions with clear single purposes (typically 10-50 lines)
- Keep cognitive complexity low (cyclomatic complexity < 10 for most functions)
- Maintain clear boundaries between different layers (presentation, business logic, data access)
- Design for testability from the start, avoiding tight coupling that prevents testing
- Apply consistent naming conventions that reveal intent without requiring comments

### Module Boundaries

**Problem:** Cross-module coupling makes changes ripple across codebase.

**Solution:** Feature-based organization with clear public interfaces:

- One feature = one directory
- Each module exposes a public API (exported functions/classes)
- Internal implementation details are private
- Cross-module calls only through public API

**Directory Structure (Language-Agnostic):**

> Paths below are illustrative examples following `project-structure.md` — the single source of truth for project layout.
> **Note:** The generic filenames below are illustrative. Always use the concrete EPMP layout from `project-structure.md` and `epmp-module-patterns.md`.

```
/task

- public_api.{ext}      # Exported interface
- business.{ext}        # Pure logic
- store.{ext}           # I/O abstraction (interface)
- postgres.{ext}        # I/O implementation
- mock.{ext}            # Test implementation
- test.{ext}            # Unit tests (mocked I/O)
- integration.test.{ext} # Integration tests (real I/O)
```

**EPMP Backend Module Example** (`backend/internal/modules/reservation/`):

```
reservation/
├── module.go                               # Public API: NewModule(db, log), RegisterRoutes(g)
├── entity/reservation.go                   # Pure domain logic & invariants
├── dto/reservation_dto.go                  # Request/Response contracts
├── repository/reservation_repository.go    # I/O abstraction (interface)
├── repository/reservation_repository_impl.go # I/O implementation (pgx)
├── service/reservation_service.go          # Use cases (orchestrates repo + entity)
├── delivery/http/reservation_handler.go    # Echo handlers
├── delivery/http/reservation_routes.go     # Route table
├── reservation_api_test.go                 # Integration test (testutil, real Postgres)
└── MODULE.md
```

**EPMP Frontend Feature Example** (`frontend/src/features/reservation/`):

```
reservation/
├── api/index.ts          # I/O: calls @/services/api
├── hooks/index.ts        # TanStack Query hooks (public API of the feature)
├── types/index.ts        # Contracts
├── schema/index.ts       # Zod validation (pure)
├── components/           # ReservationTable.tsx, ReservationForm.tsx
└── pages/                # ReservationListPage.tsx, ReservationCreatePage.tsx, ...
```

Other features import only from `hooks/`, `types/`, and `components/` — never from another feature's `api/` directly.

### Avoid Circular Dependencies

**Problem:** Module A imports B, B imports A

- Causes build failures, initialization issues
- Indicates poor module boundaries

**Solution:**

- Extract shared code to third module
- Restructure dependencies (A→C, B→C)
- Use dependency injection

### Related Principles

- Core Design Principles @core-design-principles.md
- Project Structure @project-structure.md
- Architectural Patterns — Testability-First Design @architectural-pattern.md
