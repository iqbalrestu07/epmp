# MODULE.md — TenantIdentity

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | TenantIdentity            |
| Package          | tenantidentity         |
| Table            | tenant_identities           |
| Bounded Context  | tenant  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/tenant_identities  |
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
| TenantId | string | false | false | true |
| IdentityType | string | false | false | false |
| IdentityNumber | string | false | false | true |
| FileUrl | string | false | true | false |

## Structure

```
internal/modules/tenantidentity/
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
