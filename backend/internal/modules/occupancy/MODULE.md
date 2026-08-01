# MODULE.md — Occupancy

| Metadata         | Value                |
| ---------------- | -------------------- |
| Module           | Occupancy            |
| Package          | occupancy         |
| Table            | occupancies           |
| Bounded Context  | occupancy  |
| Module Path      | github.com/epmp/backend      |

## REST API

| Property     | Value              |
| ------------ | ------------------ |
| Base Path    | /api/occupancies  |
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
| RoomId | string | false | false | true |
| TenantId | string | false | false | true |
| Status | string | false | false | false |
| CheckInTime | time.Time | false | false | false |
| CheckOutTime | time.Time | false | true | false |
| Notes | string | false | true | false |

## Structure

```
internal/modules/occupancy/
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
