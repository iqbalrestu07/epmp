# MODULE.md — Tenant

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Tenant            |
| Package          | tenant         |
| Table            | tenants           |
| Bounded Context  | tenant  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/tenants  |
| Operations   | create read update delete list |

## Behaviors

| Behavior     | Enabled |
| ------------ | ------- |
| Soft Delete  | true  |
| Pagination   | true  |
| Search       | true      |

## Fields

| Name | Type | Primary Key | Nullable | Searchable |
| ---- | ---- | ----------- | -------- | ---------- |
| OrganizationId | string | false | false | true |
| Id | string | true | false | false |
| FullName | string | false | false | true |
| Email | string | false | false | false |
| Phone | string | false | true | false |
| IdentityNumber | string | false | true | false |
| IsActive | bool | false | false | false |

## Structure

```
internal/modules/tenant/
  application/
    dto/
    service/
  domain/
    entity/
    repository/
  infrastructure/
    repository/
  interfaces/
    http/
  module.go
  MODULE.md
```
