# MODULE.md — Organization

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Organization            |
| Package          | organization         |
| Table            | organizations           |
| Bounded Context  | organization  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/organizations  |
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
| Id | string | true | false | false |
| Name | string | false | false | true |
| Domain | string | false | false | true |
| IsActive | bool | false | false | false |

## Structure

```
internal/modules/organization/
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
