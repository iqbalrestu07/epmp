# MODULE.md — Refund

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Refund            |
| Package          | refund         |
| Table            | refunds           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/refunds  |
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
| PaymentId | string | false | false | true |
| TenantId | string | false | false | true |
| Amount | float64 | false | false | false |
| Status | string | false | false | false |
| RefundDate | time.Time | false | false | false |
| Reason | string | false | false | false |

## Structure

```
internal/modules/refund/
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
