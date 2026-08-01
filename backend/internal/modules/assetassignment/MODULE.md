# MODULE.md — AssetAssignment

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | AssetAssignment            |
| Package          | assetassignment         |
| Table            | asset_assignments           |
| Bounded Context  | asset  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/asset_assignments  |
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
| AssetId | string | false | false | true |
| RoomId | string | false | false | true |
| AssignedDate | time.Time | false | false | false |

## Structure

```
internal/modules/assetassignment/
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
