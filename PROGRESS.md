# EPMP Module Integration Progress

## Hierarchy

```
Organization → Property → Building → Floor → Room
```

## Module Status

| #   | Module       | Backend Tests    | Backend Fixes                                  | Frontend Integration              | Status |
| --- | ------------ | ---------------- | ---------------------------------------------- | --------------------------------- | ------ |
| 1   | Property     | ✅ 8 tests pass  | ✅ timestamps, search, pagination, org-scoped  | ✅ real API, 3D view              | DONE   |
| 2   | Organization | ✅ 8 tests pass  | ✅ timestamps, search, pagination, org_members | ✅ real API, routes, UI, switcher | DONE   |
| 3   | Building     | ✅ 10 tests pass | ✅ timestamps, search, pagination, org-scoped  | ✅ real API, routes, UI           | DONE   |
| 4   | Floor        | ✅ 10 tests pass | ✅ timestamps, search, pagination, org-scoped  | ✅ real API, routes, UI           | DONE   |
| 5   | Room         | ✅ 10 tests pass | ✅ timestamps, search, pagination, org-scoped  | ✅ real API, routes, UI           | DONE   |
| 6   | Tenant       | ✅ 10 tests pass | ✅ timestamps, search, pagination, org-scoped  | TODO                              | DONE   |

## Org-Scoped Filtering Architecture

All core modules now enforce organization-scoped filtering:

- **Header**: `X-Organization-ID` extracted by `mw.GetOrgID(c)` middleware
- **Repository**: `FindByID(id, orgID)`, `FindAll(limit, offset, search, parentId, orgID)`, `Count(search, parentId, orgID)`
- **Service**: All methods accept `orgID` param, pass to repository
- **Handler**: Extracts `orgID` from context, validates header on Create
- **Request DTOs**: `organization_id` removed from Create/Update bodies (derived from context)
- **Frontend**: `api.ts` injects `X-Organization-ID` header from localStorage; `OrgContext` manages org selection; org switcher dropdown in `MainLayout`

### Org Switcher (Frontend)

- `OrgContext.tsx` — fetches user's orgs from `/organizations/mine`, manages current org, persists to localStorage
- `MainLayout.tsx` — dropdown in header showing org name, list of orgs with checkmark, "Create Organization" link
- `api.ts` — `getStoredOrgId()`/`setStoredOrgId()` injects header on every API call

## Per-Module Checklist

### Organization

- [x] Backend: Fix created_at/updated_at in response
- [x] Backend: Add search/filter to list API
- [x] Backend: Add Count method for pagination
- [x] Backend: Automated API integration tests
- [x] Backend: org_members pivot table + ListMine endpoint
- [x] Frontend: Connect to real API (remove mock/SDK defaults)
- [x] Frontend: Improve form (checkbox for is_active, placeholders)
- [x] Frontend: Improve table (hide raw IDs, status badges)
- [x] Frontend: Fix all routes to /dashboard/organizations
- [x] Frontend: Detail page with properties link
- [x] Frontend: Org switcher in header
- [x] Verify: tsc + build + backend tests pass

### Property

- [x] Backend: Fix created_at/updated_at in response
- [x] Backend: Add search/filter to list API
- [x] Backend: Add Count method for pagination
- [x] Backend: Automated API integration tests
- [x] Backend: Org-scoped filtering (FindByID, FindAll, Count)
- [x] Frontend: Connect to real API (remove mock/SDK defaults)
- [x] Frontend: Improve form (checkbox for is_active, placeholders)
- [x] Frontend: Improve table (hide raw IDs, status badges)
- [x] Frontend: Fix all routes to /dashboard/properties
- [x] Frontend: Detail page with buildings link
- [x] Verify: tsc + build + backend tests pass

### Building

- [x] Backend: Fix created_at/updated_at in response
- [x] Backend: Add search/filter to list API
- [x] Backend: Add Count method for pagination
- [x] Backend: Automated API integration tests
- [x] Backend: Org-scoped filtering + property_id filter
- [x] Frontend: Connect to real API
- [x] Frontend: Improve form (select for property_id, validation)
- [x] Frontend: Improve table (hide raw IDs, floor badges)
- [x] Frontend: Fix all routes to /dashboard/buildings
- [x] Frontend: Detail page with floors link
- [x] Verify: tsc + build + backend tests pass

### Floor

- [x] Backend: Fix created_at/updated_at in response
- [x] Backend: Add search/filter to list API
- [x] Backend: Add Count method for pagination
- [x] Backend: Automated API integration tests (10 tests)
- [x] Backend: Org-scoped filtering + building_id filter
- [x] Frontend: Connect to real API
- [x] Frontend: Improve form (select for building_id, checkbox for is_active)
- [x] Frontend: Improve table (icons, status badges, action buttons)
- [x] Frontend: Fix all routes to /dashboard/floors
- [x] Frontend: Add create/edit/detail pages with breadcrumbs
- [x] Frontend: Detail page with rooms link
- [x] Verify: tsc + build + backend tests pass

### Room

- [x] Backend: Fix created_at/updated_at in response
- [x] Backend: Add search/filter to list API
- [x] Backend: Add Count method for pagination
- [x] Backend: Automated API integration tests (10 tests)
- [x] Backend: Org-scoped filtering + floor_id filter
- [x] Backend: floor→floor_id migration (ULID FK)
- [x] Frontend: Connect to real API
- [x] Frontend: Improve form (select for property/floor_id, checkbox for is_available)
- [x] Frontend: Improve table (icons, price/capacity, status badges, action buttons)
- [x] Frontend: Fix all routes to /dashboard/rooms
- [x] Frontend: Add create/edit/detail pages with breadcrumbs
- [x] Frontend: Detail page with related floor link
- [x] Frontend: Fix Floor3D + PropertyInteractiveView (floor→floor_id)
- [x] Verify: tsc + build + backend tests pass

## Remaining Work

### Backend — Org-Scoped Filtering for Remaining Modules

| Module          | Status  |
| --------------- | ------- |
| Tenant          | ✅ DONE |
| Reservation     | ✅ DONE |
| Contract        | ✅ DONE |
| Occupancy       | ✅ DONE |
| Invoice/Billing | ✅ DONE |
| Payment         | ✅ DONE |
| Deposit         | ✅ DONE |
| Charge          | ✅ DONE |
| Refund          | ✅ DONE |
| Adjustment      | ✅ DONE |
| Penalty         | ✅ DONE |
| Asset           | ✅ DONE |
| AssetAssignment | ✅ DONE |
| AssetInspection | ✅ DONE |
| WorkOrder       | ✅ DONE |
| Technician      | ✅ DONE |
| Supplier        | ✅ DONE |
| Zone            | ✅ DONE |
| Bed             | ✅ DONE |
| Facility        | ✅ DONE |
| RoomType        | ✅ DONE |
| TenantContact   | ✅ DONE |
| TenantDocument  | ✅ DONE |
| TenantIdentity  | ✅ DONE |

### Frontend — Remaining

| Task                                      | Status  |
| ----------------------------------------- | ------- |
| Floor form: building select from API      | DONE    |
| Room form: property/floor select from API | DONE    |
| All remaining modules UI                  | PENDING |

## Notes

- Backend uses Echo framework with JWT auth, RBAC, PostgreSQL
- Frontend uses React + React Query + TypeScript + Zod + TailwindCSS
- API base: /api/v1
- Frontend routes: /dashboard/{module}
- JWT test secret: "test-jwt-secret-for-integration-tests"
- DB: postgres://postgres:postgres@localhost:5432/epmp?sslmode=disable
- All IDs are ULID (TEXT), not UUID
- 56 backend tests total (Property: 8, Organization: 8, Building: 10, Floor: 10, Room: 10, Tenant: 10)
- Soft deletes via `deleted_at` column
- `organization_id` always from context (X-Organization-ID header), never from request body
