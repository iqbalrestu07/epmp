# MODULE.md — TenantContact

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | TenantContact            |
| Package          | tenantcontact         |
| Table            | tenant_contacts           |
| Bounded Context  | tenant  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/tenant_contacts  |
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
| ContactType | string | false | false | false |
| ContactValue | string | false | false | true |
| IsPrimary | bool | false | false | false |

## Structure

```
internal/modules/tenantcontact/
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
