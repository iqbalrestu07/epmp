# MODULE.md — RoomType

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | RoomType            |
| Package          | roomtype         |
| Table            | room_types           |
| Bounded Context  | property  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/room_types  |
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
| Description | string | false | true | false |
| BasePrice | float64 | false | false | false |

## Structure

```
internal/modules/roomtype/
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
