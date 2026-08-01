# MODULE.md — AssetInspection

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | AssetInspection            |
| Package          | assetinspection         |
| Table            | asset_inspections           |
| Bounded Context  | asset  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/asset_inspections  |
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
| InspectionDate | time.Time | false | false | false |
| Condition | string | false | false | false |
| Notes | string | false | true | false |

## Structure

```
internal/modules/assetinspection/
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
