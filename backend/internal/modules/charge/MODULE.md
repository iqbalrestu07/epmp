# MODULE.md — Charge

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Charge            |
| Package          | charge         |
| Table            | charges           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/charges  |
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
| InvoiceId | string | false | true | true |
| ChargeType | string | false | false | false |
| Amount | float64 | false | false | false |
| Status | string | false | false | false |
| ChargeDate | time.Time | false | false | false |
| Notes | string | false | true | false |

## Structure

```
internal/modules/charge/
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
