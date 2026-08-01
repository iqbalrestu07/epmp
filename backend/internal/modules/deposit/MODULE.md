# MODULE.md — Deposit

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Deposit            |
| Package          | deposit         |
| Table            | deposits           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/deposits  |
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
| ContractId | string | false | false | true |
| TenantId | string | false | false | true |
| Amount | float64 | false | false | false |
| Status | string | false | false | false |
| CollectionDate | time.Time | false | false | false |
| RefundDate | time.Time | false | true | false |
| Notes | string | false | true | false |

## Structure

```
internal/modules/deposit/
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
