# MODULE.md — Contract

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Contract            |
| Package          | contract         |
| Table            | contracts           |
| Bounded Context  | contract  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/contracts  |
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
| ReservationId | string | false | true | true |
| TenantId | string | false | false | true |
| PropertyId | string | false | false | true |
| RoomId | string | false | false | true |
| Status | string | false | false | false |
| StartDate | time.Time | false | false | false |
| EndDate | time.Time | false | false | false |
| MonthlyRent | float64 | false | false | false |
| DepositAmount | float64 | false | false | false |
| Terms | string | false | true | false |

## Structure

```
internal/modules/contract/
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
