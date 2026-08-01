# MODULE.md — Zone

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Zone            |
| Package          | zone         |
| Table            | zones           |
| Bounded Context  | property  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/zones  |
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
| BuildingId | string | false | false | true |
| Floor | int | false | false | false |
| Name | string | false | false | true |

## Structure

```
internal/modules/zone/
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
