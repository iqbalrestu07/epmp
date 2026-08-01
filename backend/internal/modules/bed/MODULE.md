# MODULE.md — Bed

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Bed            |
| Package          | bed         |
| Table            | beds           |
| Bounded Context  | property  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/beds  |
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
| RoomId | string | false | false | true |
| Name | string | false | false | true |
| Status | string | false | false | false |

## Structure

```
internal/modules/bed/
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
