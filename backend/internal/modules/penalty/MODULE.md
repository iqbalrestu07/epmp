# MODULE.md — Penalty

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Penalty            |
| Package          | penalty         |
| Table            | penalties           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/penalties  |
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
| Amount | float64 | false | false | false |
| Status | string | false | false | false |
| PenaltyDate | time.Time | false | false | false |
| Description | string | false | false | false |

## Structure

```
internal/modules/penalty/
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
