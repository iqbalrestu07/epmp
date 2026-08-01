# MODULE.md — Supplier

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Supplier            |
| Package          | supplier         |
| Table            | vendors           |
| Bounded Context  | maintenance  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/vendors  |
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
| Name | string | false | false | true |
| ContactPerson | string | false | false | false |
| Phone | string | false | false | false |
| ServiceType | string | false | false | false |

## Structure

```
internal/modules/supplier/
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
