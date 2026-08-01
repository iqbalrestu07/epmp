# MODULE.md — TenantDocument

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | TenantDocument            |
| Package          | tenantdocument         |
| Table            | tenant_documents           |
| Bounded Context  | tenant  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/tenant_documents  |
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
| DocumentType | string | false | false | false |
| FileUrl | string | false | false | false |

## Structure

```
internal/modules/tenantdocument/
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
