---
trigger: always_on
---

## Rule Priority (When Rules Conflict)

When two rules pull in opposite directions, use this priority to decide:

### Priority Order (Highest to Lowest)

1. **Security Mandate** — always wins. Never compromise security for velocity, simplicity, or convenience.
2. **Rugged Software Constitution** — foundational philosophy. Code must be defensible.
3. **Code Completion Mandate**, **Logging and Observability Mandate**, **Schema Documentation Mandate**, **Deployment and Data Safety Mandate**, and **Multi-Tenancy Boundary** — all are always-on enforcement rules. Validation, instrumentation, schema documentation, zero-data-loss safeguards, and organization isolation are non-negotiable; none can be skipped to ship faster.
4. **Testability-First Design** — maintainability enables future improvements.
5. **Feature-specific principles** — context-dependent guidance for the task at hand. This includes `epmp-module-patterns.md`, `go-idioms-and-patterns.md`, and `typescript-idioms-and-patterns.md`. When an idiom conflicts with a higher-priority rule (e.g., security or testability), the higher-priority rule always wins.
6. **Spec-gated principles** — only apply when `epmp-docs/` or a `MODULE.md` _explicitly requires_ the capability. Must not be applied based on speculation; agent must confirm requirement before activating.
7. **YAGNI / KISS** — only when no security, reliability, or maintainability trade-off exists.

### Common Conflict Resolutions

| Conflict                                                              | Resolution                                                                                 |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| YAGNI vs Security ("don't add input validation, it's not needed yet") | **Security wins.** Input validation is always needed.                                      |
| KISS vs Testability ("adding an interface makes it more complex")     | **Testability wins.** Interfaces enable testing, which enables maintainability.            |
| Performance vs YAGNI ("should I optimize this now?")                  | **Measure first.** Only optimize after profiling shows a real bottleneck.                  |
| DRY vs Clarity ("should I abstract this into a shared utility?")      | **Clarity wins** until duplication reaches 3+ instances (Rule of Three).                   |
| Speed vs Logging ("skip logging to ship faster")                      | **Logging wins.** Silent failures are the enemy.                                           |
| YAGNI vs spec-gated requirement ("feature flags aren't needed")       | **Spec wins.** If `epmp-docs` explicitly requires it, YAGNI cannot override it.            |
| "it's just a small column" vs Schema Documentation                    | **Documentation wins.** Every column appears in `MODULE.md` + migration README.            |
| Speed vs Data Safety ("just drop the column in this migration")       | **Safety wins.** Destructive DDL only on explicit request, two-phase removal.              |
| Convenience vs Org Scope ("skip the organization_id filter here")     | **Isolation wins.** Every tenant query is org-scoped; 404 on cross-org access.             |
| Hand-edit vs Generator ("faster to edit the generated file")          | **Generator wins.** Change the schema YAML and regenerate, or drop the DO NOT EDIT header. |

### Guiding Principle

When in doubt, ask: _"Which choice makes the code more defensible and maintainable?"_

If both options are equally defensible, choose the simpler one (KISS).

### Related Principles

- Security Mandate @security-mandate.md
- Rugged Software Constitution @rugged-software-constitution.md
- Code Completion Mandate @code-completion-mandate.md
- Logging and Observability Mandate @logging-and-observability-mandate.md
- Schema Documentation Mandate @schema-documentation-mandate.md
- Deployment and Data Safety Mandate @deployment-and-data-safety-mandate.md
- Multi-Tenancy Boundary @multi-tenancy-boundary.md
- Architectural Patterns — Testability-First Design @architectural-pattern.md
