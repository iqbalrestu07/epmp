# MODULE.md — Reservation

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Reservation            |
| Package          | reservation         |
| Table            | reservations           |
| Bounded Context  | reservation  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/reservations  |
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
| TenantId | string | false | false | true |
| PropertyId | string | false | false | true |
| RoomId | string | false | false | true |
| Status | string | false | false | false |
| CheckInDate | time.Time | false | false | false |
| CheckOutDate | time.Time | false | true | false |
| BookingFee | float64 | false | false | false |
| Notes | string | false | true | false |

## Structure

```
internal/modules/reservation/
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
