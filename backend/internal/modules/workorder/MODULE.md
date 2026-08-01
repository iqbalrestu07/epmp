# MODULE.md — WorkOrder

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | WorkOrder            |
| Package          | workorder         |
| Table            | work_orders           |
| Bounded Context  | maintenance  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/work_orders  |
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
| RoomId | string | false | true | true |
| Description | string | false | false | false |
| Status | string | false | false | false |
| Priority | string | false | false | false |

## Structure

```
internal/modules/workorder/
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
