# MODULE.md — Payment

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Payment            |
| Package          | payment         |
| Table            | payments           |
| Bounded Context  | billing  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/payments  |
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
| TenantId | string | false | false | true |
| Amount | float64 | false | false | false |
| PaymentDate | time.Time | false | false | false |
| PaymentMethod | string | false | false | false |
| Status | string | false | false | false |
| ReferenceNumber | string | false | true | false |

## Structure

```
internal/modules/payment/
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
