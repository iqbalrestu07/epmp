# MODULE.md — Asset

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Asset            |
| Package          | asset         |
| Table            | assets           |
| Bounded Context  | asset  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/assets  |
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
| PropertyId | string | false | false | true |
| Name | string | false | false | true |
| Category | string | false | false | false |
| Status | string | false | false | false |
| PurchasePrice | float64 | false | false | false |

## Structure

```
internal/modules/asset/
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
