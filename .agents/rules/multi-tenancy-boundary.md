---
trigger: always_on
---

## Multi-Tenancy Boundary — Organization Scope

### Core Principle

**Every tenant-owned record belongs to exactly one `organization_id`, and every read/write MUST be scoped to the caller's active organization.** A request from organization A must never see, mutate, or infer the existence of organization B's data.

### How Scope Flows

```
Frontend (localStorage epmp_org_id)
    │  Authorization: Bearer <jwt>
    │  X-Organization-ID: <org uuid>
    ▼
mw.AuthRequired  ── sets ContextKeyOrgID from header
    ▼
Handler          ── orgID := mw.GetOrgID(c)   (400 if required and empty)
    ▼
Service          ── passes orgID through every use case
    ▼
Repository       ── WHERE organization_id = $n AND deleted_at IS NULL
```

### Rules

1. **Handlers** read the org with `mw.GetOrgID(c)` and pass it explicitly. Create operations MUST reject an empty org (`response.BadRequest(c, "X-Organization-ID header is required")`).
2. **Repository interfaces** take `orgID string` on every method that touches tenant data (`FindByID(ctx, id, orgID)`, `Delete(ctx, id, orgID)`, ...). Do not add unscoped variants.
3. **SQL** always includes `organization_id = $n` for tenant tables. Cross-org lookups return `errs.ErrNotFound` (404), not 403 — do not leak existence.
4. **Soft delete** is org-scoped too: `UPDATE ... SET deleted_at = now() WHERE id = $1 AND organization_id = $2`.
5. **Do not trust client-supplied org identifiers in the body or query string** for authorization decisions. The header set by `AuthRequired` is the source of truth. If an endpoint accepts `organization_id` as a filter (e.g. list endpoints), it MUST be validated against the caller's memberships before use, or removed.
6. **Global tables** (`users`, `roles`, `organizations`, `permissions`) are the only exception; membership tables (`property_user_roles`, org memberships) bridge them to tenant scope.
7. **Migrations** for new tenant tables MUST include `organization_id UUID NOT NULL` with an index (see @schema-documentation-mandate.md).
8. **Frontend** never hardcodes an org id; it reads `getStoredOrgId()` from `@/services/api` and the wrapper injects the header.

### Test Requirement

Every module's `<name>_api_test.go` MUST prove isolation with `testutil.AssertCRUD`, which creates a second organization and asserts `GET/PUT/DELETE` on the first org's entity return non-2xx (expected 404). Do not weaken this assertion.

### Violation Examples

❌ `SELECT * FROM rooms WHERE id = $1` (no org filter)
❌ `orgID := c.QueryParam("organization_id")` used for authorization without membership check
❌ Returning `403 Forbidden` for another org's record (confirms existence)

✅ `SELECT ... FROM rooms WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL`
✅ `orgID := mw.GetOrgID(c)` passed down to repository
✅ `errs.ErrNotFound` → `response.NotFound`

### Related Principles
- EPMP Module Patterns @epmp-module-patterns.md
- Security Mandate @security-mandate.md
- Database Design Principles @database-design-principles.md
- API Design Principles @api-design-principles.md
