# MODULE.md — Invoice

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Invoice            |
| Package          | billing         |
| Table            | invoices           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/invoices  |
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
| DueDate | time.Time | false | false | false |
| PaidDate | time.Time | false | true | false |
| PaymentMethod | string | false | true | false |
| Notes | string | false | true | false |

## Structure

```
internal/modules/billing/
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
