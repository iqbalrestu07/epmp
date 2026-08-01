# MODULE.md — Technician

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Technician            |
| Package          | technician         |
| Table            | technicians           |
| Bounded Context  | maintenance  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/technicians  |
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
| Phone | string | false | false | false |
| Specialty | string | false | false | false |

## Structure

```
internal/modules/technician/
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
