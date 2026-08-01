# MODULE.md — Property

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Property            |
| Package          | property         |
| Table            | properties           |
| Bounded Context  | property  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/properties  |
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
| Address | string | false | true | false |
| PropertyType | string | false | false | false |
| IsActive | bool | false | false | false |

## Structure

```
internal/modules/property/
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
