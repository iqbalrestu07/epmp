# MODULE.md — Adjustment

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Adjustment            |
| Package          | adjustment         |
| Table            | adjustments           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/adjustments  |
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
| InvoiceId | string | false | false | true |
| AdjustmentType | string | false | false | false |
| Amount | float64 | false | false | false |
| AdjustmentDate | time.Time | false | false | false |
| Reason | string | false | false | false |

## Structure

```
internal/modules/adjustment/
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
