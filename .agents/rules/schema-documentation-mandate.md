---
trigger: always_on
---

## Schema Documentation Mandate

> **This rule has the same enforcement weight as `code-completion-mandate.md`.**
> Delivering a migration or schema change without updating its documentation is INCOMPLETE work.

### Core Requirement

**Every change to a PostgreSQL table or column MUST be delivered in the same commit as:**

1. A `backend/migrations/NNNNNN_<name>.up.sql` **and** matching `.down.sql`
2. An updated row in the migration table of `backend/migrations/README.md`
3. The affected module's `MODULE.md` (Fields / Behaviors tables)
4. The module's schema YAML in `tools/epmp-sdk/schemas/<name>.yaml` (when the module is generator-managed)
5. The frontend `types/` + `schema/` (Zod) of the corresponding feature, when the column is exposed through the API

This applies to ALL operations: new table, new column, rename/drop column, type change, new index/constraint, new enum value.

No exceptions. Not "will document later." Not "minor column."

---

### When This Rule Applies

This rule activates whenever you:

1. Create or modify a file in `backend/migrations/`
2. Change a repository SQL statement to read/write a column that did not exist before
3. Change a field in `tools/epmp-sdk/schemas/*.yaml`
4. Discover an undocumented column during development

---

### Migration Standards

```bash
cd backend && go run ./cmd/migrate create add_currency_to_properties
# → migrations/0000NN_add_currency_to_properties.up.sql / .down.sql
```

- Sequence is a zero-padded 6-digit integer; one logical change per migration.
- `.up.sql` and `.down.sql` are both mandatory; `.down.sql` must fully reverse `.up.sql`.
- Tenant-owned tables MUST have: `id UUID PRIMARY KEY`, `organization_id UUID NOT NULL` (+ index), `created_at`, `updated_at`, `deleted_at TIMESTAMPTZ NULL`.
- Use `IF NOT EXISTS` / `IF EXISTS` so re-runs are safe.
- Never edit a migration that may already be applied in any environment — add a new one.

---

### Documentation Format

**`backend/migrations/README.md`** — one row per migration:

```markdown
| Seq | Name | Domain | Table | Depends On |
| 0000NN | add_currency_to_properties | property | properties | 000001 |
```

**`MODULE.md`** — keep the `Fields` table a 1-to-1 mirror of the live table (name, type, PK, nullable, searchable).

**Behaviour changes** (new business rule, new state) — also update the relevant `epmp-docs/epmp-0xx.md` module spec.

---

### Strict 1-to-1 Synchronization Rule

`MODULE.md` and `tools/epmp-sdk/schemas/<name>.yaml` must mirror the live schema. If a column exists in Postgres but not in these files, documentation is **out of sync** and must be fixed immediately.

**Enforcement:**

- Before writing a migration: read the module's `MODULE.md` and the latest `*.up.sql` touching that table
- After writing a migration: update all artefacts listed under _Core Requirement_
- If you discover an undocumented column during any task: document it before continuing

---

### Discovery Protocol for Undocumented Columns

1. **Stop current task**
2. Investigate: `grep -rn "<column>" backend/migrations/` to find when it was added
3. Document it in `MODULE.md` (and schema YAML if applicable)
4. Note in the commit body: `[Retroactive documentation] <table>.<column>`
5. Resume current task

---

### AI Agent Checklist (Pre-flight for Any Migration Task)

Before writing any migration file, verify:

- [ ] Have I read the module's `MODULE.md` and existing migrations for this table?
- [ ] Does the documentation reflect the current live schema?
- [ ] Is there an undocumented column I need to backfill first?

After writing any migration file, verify:

- [ ] `.up.sql` and `.down.sql` both exist and `go run ./cmd/migrate up` then `down` succeed locally
- [ ] `backend/migrations/README.md` row added
- [ ] `MODULE.md` Fields table updated
- [ ] `tools/epmp-sdk/schemas/<name>.yaml` updated (if generator-managed) and module regenerated or hand-edit header removed
- [ ] Frontend `types/` + `schema/` updated when the column is API-visible

---

### Violation Examples

❌ **Wrong:** Add `properties.currency` in a migration and only update the repository SQL.

✅ **Correct:** Migration pair + README row + `MODULE.md` field + `schemas/property.yaml` + frontend `propertySchema` (`currency: z.string()`).

❌ **Wrong:** Notice `rooms.price` is missing from `MODULE.md` and ignore it to finish faster.

✅ **Correct:** Document `rooms.price` immediately, then continue.

---

### Related Principles

- Database Design Principles @database-design-principles.md
- EPMP Module Patterns @epmp-module-patterns.md
- Code Completion Mandate @code-completion-mandate.md
- Project Structure @project-structure.md
- Deployment and Data Safety Mandate @deployment-and-data-safety-mandate.md
