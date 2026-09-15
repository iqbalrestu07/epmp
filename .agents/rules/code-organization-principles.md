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
> **Note:** The generic filenames below are illustrative. Always use language-specific naming conventions from the relevant `project-structure-*` file (e.g., `storage.go` in Go, `task.api.ts` in TypeScript, `task_repository.dart` in Dart).

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

**Directus Extension Example:**

For this project, extensions are thin API adapters — each is a single `index.js`:

```
extensions/survey-forms/
├── index.js                    # Single entry point
└── package.json                # Extension manifest
```

For larger extensions with reusable logic, extract pure functions:

```
extensions/survey-forms/
├── index.js                    # Router & handlers (I/O layer)
├── transform.js                # Pure data transformation (business logic)
├── transform.spec.js           # Unit tests for pure functions
└── package.json
```

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