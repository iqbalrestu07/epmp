# MODULE.md — Tenant

| Metadata        | Value                        |
| --------------- | ---------------------------- |
| Module          | Tenant                       |
| Package         | tenant                       |
| Table           | tenants                      |
| Bounded Context | tenant                       |
| Module Path     | github.com/epmp/backend      |

## REST API

| Property   | Value                              |
| ---------- | ---------------------------------- |
| Base Path  | `/api/v1/tenants`                  |
| Operations | create, read, update, delete, list |
| Auth       | Required (JWT middleware on v1 group) |

### Endpoints

| Method | Path                 | Handler | Description      |
| ------ | -------------------- | ------- | ---------------- |
| POST   | /api/v1/tenants      | Create  | Create a tenant  |
| GET    | /api/v1/tenants      | List    | List tenants     |
| GET    | /api/v1/tenants/:id  | GetByID | Get one by ID    |
| PUT    | /api/v1/tenants/:id  | Update  | Update a tenant  |
| DELETE | /api/v1/tenants/:id  | Delete  | Delete a tenant  |

## Behaviors

| Behavior    | Enabled |
| ----------- | ------- |
| Soft Delete | true    |
| Pagination  | true    |
| Search      | true    |

## Fields

| Name           | Type   | PK    | Nullable | Searchable |
| -------------- | ------ | ----- | -------- | ---------- |
| Id             | string | true  | false    | false      |
| FullName       | string | false | false    | true       |
| Email          | string | false | false    | false      |
| Phone          | string | false | true     | false      |
| IdentityNumber | string | false | true     | false      |
| IsActive       | bool   | false | false    | false      |

## Wiring

```go
// cmd/server/main.go
v1 := e.Group("/api/v1")
tenant.NewModule(db, log).RegisterRoutes(v1)
```

## Structure

```
internal/modules/tenant/
  application/
    dto/        # CreateTenantRequest, UpdateTenantRequest, TenantResponse
    service/    # TenantService (use-case orchestration)
  domain/
    entity/     # Tenant entity
    repository/ # TenantRepository interface
  infrastructure/
    repository/ # PostgreSQL implementation
  interfaces/
    http/       # TenantHandler, RegisterTenantRoutes
  module.go     # DI wiring + RegisterRoutes(v1 *echo.Group)
  MODULE.md
```
