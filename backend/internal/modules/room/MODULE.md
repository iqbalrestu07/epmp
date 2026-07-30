# MODULE.md — Room

| Metadata        | Value                        |
| --------------- | ---------------------------- |
| Module          | Room                         |
| Package         | room                         |
| Table           | rooms                        |
| Bounded Context | property                     |
| Module Path     | github.com/epmp/backend      |

## REST API

| Property   | Value                              |
| ---------- | ---------------------------------- |
| Base Path  | `/api/v1/rooms`                    |
| Operations | create, read, update, delete, list |
| Auth       | Required (JWT middleware on v1 group) |

### Endpoints

| Method | Path              | Handler | Description   |
| ------ | ----------------- | ------- | ------------- |
| POST   | /api/v1/rooms     | Create  | Create a room |
| GET    | /api/v1/rooms     | List    | List rooms    |
| GET    | /api/v1/rooms/:id | GetByID | Get one by ID |
| PUT    | /api/v1/rooms/:id | Update  | Update a room |
| DELETE | /api/v1/rooms/:id | Delete  | Delete a room |

## Behaviors

| Behavior    | Enabled |
| ----------- | ------- |
| Soft Delete | true    |
| Pagination  | true    |
| Search      | true    |

## Fields

| Name        | Type    | PK    | Nullable | Searchable |
| ----------- | ------- | ----- | -------- | ---------- |
| Id          | string  | true  | false    | false      |
| Name        | string  | false | false    | true       |
| Floor       | int     | false | false    | false      |
| Capacity    | int     | false | false    | false      |
| Price       | float64 | false | false    | false      |
| IsAvailable | bool    | false | false    | false      |
| PropertyId  | string  | false | false    | false      |

## Wiring

```go
// cmd/server/main.go
v1 := e.Group("/api/v1")
room.NewModule(db, log).RegisterRoutes(v1)
```

## Structure

```
internal/modules/room/
  application/
    dto/        # CreateRoomRequest, UpdateRoomRequest, RoomResponse
    service/    # RoomService (use-case orchestration)
  domain/
    entity/     # Room entity
    repository/ # RoomRepository interface
  infrastructure/
    repository/ # PostgreSQL implementation
  interfaces/
    http/       # RoomHandler, RegisterRoomRoutes
  module.go     # DI wiring + RegisterRoutes(v1 *echo.Group)
  MODULE.md
```
