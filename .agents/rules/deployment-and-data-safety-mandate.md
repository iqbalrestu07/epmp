---
trigger: always_on
---

# Deployment and Data Safety Mandate

### Core Principle: Zero Data-Loss Architecture

Deployments and migrations between environments (`dev` ➔ `staging` ➔ `production`) MUST NEVER drop tables, drop columns, truncate, or bulk-delete existing records unless explicitly requested by a human with an interactive confirmation for that specific action.

---

### 1. Schema Migration Safety (golang-migrate)

Migrations live in `backend/migrations/` and run via `cd backend && go run ./cmd/migrate <up|down|version|force|create>` or the one-shot `migrate` service in `docker-compose.yml`.

#### Mandatory Safeguards

1. **Additive by default** — new tables/columns/indexes are the norm. Destructive DDL (`DROP TABLE`, `DROP COLUMN`, `ALTER ... TYPE` with data loss, `TRUNCATE`) in an `.up.sql` requires:
   - explicit user request in the current task, and
   - a preceding data-preservation step (backfill / copy to new column) in an earlier migration, and
   - a note in the commit body: `DESTRUCTIVE: <what and why>`.
2. **Down migrations are for local rollback only.** Never run `migrate down` against staging/production data without a verified backup. `.down.sql` still MUST exist and reverse the `.up.sql`.
3. **Never use `migrate force`** to paper over a failed migration in shared environments without first inspecting the partial state and documenting it.
4. **Never edit an applied migration.** Add a new sequence number.
5. **Column removal is two-phase**: (a) stop reading/writing the column in code and ship, (b) drop it in a later release once no running version depends on it.

---

### 2. Data Safety

1. **No manual production data edits by AI** (`tools/epmp-ai/RULES.md` §8). Data fixes are migrations or reviewed scripts, run by a human.
2. **Soft delete only** — application code uses `deleted_at = now()`; physical `DELETE` on tenant data is prohibited outside purge jobs explicitly designed for it.
3. **Seed/fixture scripts must be idempotent** and target only `APP_ENV=development`. Guard with an environment check before any write.
4. **Backfills** are separate migrations with `WHERE <col> IS NULL` guards and, for large tables, batched updates.

---

### 3. Docker & Local Environment

- `docker compose down -v` **destroys the `postgres_data` volume**. Never run it (or suggest it) as a routine "restart"; use `docker compose down` / `docker compose restart postgres`.
- `docker compose up -d` runs the `migrate` service automatically before `backend` starts; failed migrations block the backend by design — fix the migration, do not bypass the dependency.
- Local DB URL default: `postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable` (also the default in `internal/testutil`). Integration tests write to this DB — never point `DATABASE_URL` for tests at a shared environment.

---

### 4. Environment Roles & Truth Hierarchy

- **DEV**: local Docker Postgres; disposable, may be reset by the developer (not by AI without asking).
- **STAGING**: mirrors production schema; data may be anonymised copies. Migrations are rehearsed here first.
- **PRODUCTION**: golden source of truth for business data. Only additive migrations by default; backups verified before any release containing schema changes.

---

### Enforcement Checklist

Before proposing or running anything that touches a database:

- [ ] Is every DDL statement additive? If not, was destruction explicitly requested and documented?
- [ ] Do `.up.sql` and `.down.sql` both exist and reverse each other?
- [ ] Am I targeting the local dev database (`localhost:5432/epmp`) and not a shared one?
- [ ] Have I avoided `docker compose down -v`, `migrate force`, and `migrate down` on shared data?
- [ ] For column/table removal: is this the second phase, after code stopped using it?

### Related Principles

- Schema Documentation Mandate @schema-documentation-mandate.md
- Database Design Principles @database-design-principles.md
- Security Mandate @security-mandate.md
