---
trigger: always_on
---

## Documentation Principles

### Self-Documenting Code

**Clear naming reduces need for comments:**

- Code shows WHAT is happening
- Comments explain WHY it's done this way

**When to comment:**

- Complex business logic deserves explanation
- Non-obvious algorithms (explain approach)
- Workarounds for bugs (link to issue tracker)
- Performance optimizations (explain trade-offs)

### Documentation Levels

1. **Inline comments:** Explain WHY for complex code
2. **Function/method docs:** API contract (parameters, returns, errors)
3. **Module/package docs:** High-level purpose and usage
4. **README:** Setup, usage, examples
5. **Architecture docs:** System design, component interactions, use valid mermaid diagram

### Spec, MODULE.md & PROGRESS.md Update Mandate

**When implementing features that fulfill or change planned milestones, the corresponding spec (`epmp-docs/epmp-XXX-*.md`), the module's `MODULE.md`, and `PROGRESS.md` MUST be updated as part of the same work.**

**Triggers for update:**

- New table or schema change → new migration in `backend/migrations/` + update `MODULE.md` (entities, migrations)
- Milestone status change (Planned → In Progress → Done) → update `PROGRESS.md`
- New endpoint or API change → update `MODULE.md` (API table) and the feature spec if the contract changes
- New FR (functional requirement) completed → mark done in the `epmp-docs` spec

**What to update:**
| Document | When | What |
|---|---|---|
| **PROGRESS.md** | Milestone/feature status changes | Update the progress tracker |
| **`epmp-docs/epmp-XXX-*.md`** | Feature spec changes | Update scope, FR/NFR IDs, mark items done |
| **`MODULE.md`** (per module) | New entity / endpoint / schema change | Update entities, API table, and migration list |
| **`epmp-docs/adr/`** | Architecture decision | Add ADR via the ADR skill |

> **Rule:** Never ship a feature that changes the schema or completes a milestone without updating `MODULE.md` and `PROGRESS.md`. Treat documentation as part of the definition of done.

### Related Principles

- Core Design Principles @core-design-principles.md
- Code Organization Principles @code-organization-principles.md
- Schema Documentation Mandate @schema-documentation-mandate.md
