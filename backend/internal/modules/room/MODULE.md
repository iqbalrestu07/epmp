# MODULE.md — Room

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Room            |
| Package          | room         |
| Table            | rooms           |
| Bounded Context  | property  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/rooms  |
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
| Floor | int | false | false | false |
| Capacity | int | false | false | false |
| Price | float64 | false | false | false |
| IsAvailable | bool | false | false | false |
| PropertyId | string | false | false | false |

## Structure

```
internal/modules/room/
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
