# MODULE.md — Floor

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Floor            |
| Package          | floor         |
| Table            | floors           |
| Bounded Context  | property  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/floors  |
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
| OrganizationId | string | false | false | true |
| BuildingId | string | false | false | true |
| Name | string | false | false | true |
| FloorNumber | int | false | false | true |
| IsActive | bool | false | false | false |

## Structure

```
internal/modules/floor/
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
