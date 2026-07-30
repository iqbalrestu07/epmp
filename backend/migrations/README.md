# Migrations

SQL migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

## Naming Convention

```
{sequence}_{description}.up.sql    — forward migration
{sequence}_{description}.down.sql  — rollback
```

Sequence is a zero-padded 6-digit integer (`000001`, `000002`, …).  
Each migration **must** have both an `.up.sql` and a `.down.sql`.

## Current Migrations

| Seq    | Name                    | Domain   | Table       | Depends On |
| ------ | ----------------------- | -------- | ----------- | ---------- |
| 000001 | create_properties       | property | properties  | —          |
| 000002 | create_tenants          | tenant   | tenants     | —          |
| 000003 | create_rooms            | room     | rooms       | 000001     |

## Running Migrations

> All commands must be run from the `backend/` directory.

```bash
# Apply all pending migrations
go run ./cmd/migrate up

# Rollback the last migration
go run ./cmd/migrate down

# Rollback N migrations
go run ./cmd/migrate down 2

# Show current version
go run ./cmd/migrate version

# Fix dirty state (use only after a failed migration)
go run ./cmd/migrate force <version>
```

Or via Makefile:

```bash
make migrate-up
make migrate-down
make migrate-create name=add_column_to_tenants
```

## Adding a New Migration

### Option 1 — via the runner (recommended)

```bash
go run ./cmd/migrate create add_status_to_properties
```

This generates a numbered pair:

```
migrations/
  000004_add_status_to_properties.up.sql
  000004_add_status_to_properties.down.sql
```

Fill in the SQL, then apply:

```bash
go run ./cmd/migrate up
```

### Option 2 — manually

1. Find the next sequence number from existing files.
2. Create both files following the naming convention.
3. Write idempotent SQL (`IF NOT EXISTS`, `IF EXISTS`).
4. Run `go run ./cmd/migrate up`.

## Auto-Watch Mode

The watch command polls `migrations/` every 3 seconds. When it detects new `.up.sql` files, it automatically applies all pending migrations.

```bash
go run ./cmd/migrate watch
```

Useful during active development when generating new modules with the EPMP SDK codegen — the watcher picks up migration files the moment the generator writes them.

Stop with `Ctrl+C`.

## Writing Good Migrations

- Always pair `.up.sql` with `.down.sql`.
- Use `IF NOT EXISTS` / `IF EXISTS` for safety.
- Prefer `ALTER TABLE` for schema changes (never drop and re-create).
- Keep each migration focused on a single concern.
- Foreign keys: reference the parent table's migration sequence in the comment.
- Triggers: define `set_updated_at()` only once (migration `000001`), reuse in all others.
