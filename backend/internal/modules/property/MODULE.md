# MODULE.md — Property

| Metadata        | Value                        |
| --------------- | ---------------------------- |
| Module          | Property                     |
| Package         | property                     |
| Table           | properties                   |
| Bounded Context | property                     |
| Module Path     | github.com/epmp/backend      |

## REST API

| Property   | Value                          |
| ---------- | ------------------------------ |
| Base Path  | `/api/v1/properties`           |
| Operations | create, read, update, delete, list |
| Auth       | Required (JWT middleware on v1 group) |

### Endpoints

| Method | Path                    | Handler        | Description         |
| ------ | ----------------------- | -------------- | ------------------- |
| POST   | /api/v1/properties      | Create         | Create a property   |
| GET    | /api/v1/properties      | List           | List properties     |
| GET    | /api/v1/properties/:id  | GetByID        | Get one by ID       |
| PUT    | /api/v1/properties/:id  | Update         | Update a property   |
| DELETE | /api/v1/properties/:id  | Delete         | Delete a property   |

## Behaviors

| Behavior    | Enabled |
| ----------- | ------- |
| Soft Delete | true    |
| Pagination  | true    |
| Search      | true    |

## Fields

| Name         | Type   | PK    | Nullable | Searchable |
| ------------ | ------ | ----- | -------- | ---------- |
| Id           | string | true  | false    | false      |
| Name         | string | false | false    | true       |
| Description  | string | false | true     | false      |
| Address      | string | false | true     | false      |
| PropertyType | string | false | false    | false      |
| IsActive     | bool   | false | false    | false      |

## Wiring

```go
// cmd/server/main.go
v1 := e.Group("/api/v1")
property.NewModule(db, log).RegisterRoutes(v1)
```

## Structure

```
internal/modules/property/
  application/
    dto/        # CreatePropertyRequest, UpdatePropertyRequest, PropertyResponse
    service/    # PropertyService (use-case orchestration)
  domain/
    entity/     # Property entity
    repository/ # PropertyRepository interface
  infrastructure/
    repository/ # PostgreSQL implementation
  interfaces/
    http/       # PropertyHandler, RegisterPropertyRoutes
  module.go     # DI wiring + RegisterRoutes(v1 *echo.Group)
  MODULE.md
```
