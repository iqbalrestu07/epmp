# EPMP SDK — Backend Code Generator

Code generation toolkit for EPMP backend modules.

## Overview

The codegen tool generates Clean Architecture boilerplate for bounded contexts from a single YAML config file.

```
epmp-sdk/be/codegen/
  cmd/codegen/          # CLI entry point
  config/               # Config schema, loader, validator
  internal/
    renderer/           # Template rendering (text/template)
    filesystem/         # File I/O operations
    generator/          # Module generator (uses renderer + filesystem)
      templates/        # Embedded Go templates
  examples/             # Sample config files
```

## Architecture

```
config.yaml
    │
    ▼
config.Load() ──► config.Validate()
    │
    ▼
generator.Generate()
    │
    ├── renderer.Render(template, data)  →  []byte
    │
    └── filesystem.WriteFile(path, content)
```

- **Renderer** loads and renders Go `text/template` templates. Returns bytes. Does not write files.
- **Filesystem** handles directory creation and file writing. Does not know about templates.
- **Generator** orchestrates both. Maps config to template data, renders each artifact, writes to disk.

## Installation

```bash
cd tools/epmp-sdk/be/codegen
go build -o epmp-codegen ./cmd/codegen/
```

## Usage

### Generate a Module

```bash
./epmp-codegen --config examples/property.yaml --output ../../../../backend/internal
```

The generator creates files under `{output}/modules/{package}/`. So with `--output backend/internal`, files go to `backend/internal/modules/property/`.

### Dry Run (no files written)

```bash
./epmp-codegen --config examples/property.yaml --output ./output --dry-run
```

### Override Module Path

```bash
./epmp-codegen --config examples/property.yaml --output ./output --module github.com/myorg/myapp
```

### Flags

| Flag        | Description                         | Required |
| ----------- | ----------------------------------- | -------- |
| `--config`  | Path to the module config YAML file | Yes      |
| `--output`  | Output root directory               | No\*     |
| `--module`  | Go module path                      | No\*     |
| `--dry-run` | Render without writing files        | No       |

\* Falls back to `backend.output_root` and `backend.module_path` from config.

## Config File Format

```yaml
version: "1.0"
dry_run: false

backend:
  output_root: "./internal"
  module_path: "github.com/epmp/backend"
  artifacts:
    - dto
    - repository
    - rest
    - test

domain:
  name: Property # PascalCase entity name
  package: property # lowercase Go package name
  table: properties # database table name
  bounded_context: property

  fields:
    - name: id
      type: uuid
      primary_key: true
      auto: true

    - name: name
      type: string
      max_length: 255
      searchable: true

    - name: description
      type: text
      nullable: true

  behaviors:
    soft_delete: true
    audit_trail: true
    pagination: true
    search: true
    filter_by:
      - is_active
    sort_by:
      - name

  rest:
    base_path: /api/properties
    auth_required: true
    operations:
      - create
      - read
      - update
      - delete
      - list

  events:
    - PropertyCreated
    - PropertyUpdated
```

### Field Types

| Type        | Go Type     |
| ----------- | ----------- |
| `uuid`      | `string`    |
| `string`    | `string`    |
| `int`       | `int`       |
| `bool`      | `bool`      |
| `text`      | `string`    |
| `enum`      | `string`    |
| `timestamp` | `time.Time` |
| `decimal`   | `float64`   |

### Field Options

| Option        | Type     | Description                          |
| ------------- | -------- | ------------------------------------ |
| `name`        | string   | Field name (snake_case)              |
| `type`        | string   | One of the field types above         |
| `nullable`    | bool     | Field can be null                    |
| `primary_key` | bool     | Field is primary key                 |
| `unique`      | bool     | Field has unique constraint          |
| `max_length`  | int      | Max string length                    |
| `default`     | any      | Default value                        |
| `searchable`  | bool     | Field is searchable                  |
| `auto`        | bool     | Auto-generated (e.g. auto-increment) |
| `enum_values` | []string | Required for enum type (min 2)       |
| `foreign_key` | object   | Reference to another domain          |

### REST Operations

| Operation | HTTP Method | Route Pattern      |
| --------- | ----------- | ------------------ |
| `create`  | POST        | `{base_path}`      |
| `read`    | GET         | `{base_path}/{id}` |
| `update`  | PUT         | `{base_path}/{id}` |
| `delete`  | DELETE      | `{base_path}/{id}` |
| `list`    | GET         | `{base_path}`      |

## Generated Artifacts

For each module, the generator produces 10 files following Clean Architecture:

```
internal/modules/{package}/
  MODULE.md                                    # Module documentation
  module.go                                    # Module init & DI wiring
  domain/
    entity/
      {name}.go                                # Domain entity struct (lowercase)
    repository/
      {name}_repository.go                     # Repository interface (lowercase)
  application/
    dto/
      {name}_dto.go                            # Request/Response DTOs (lowercase)
    service/
      {name}_service.go                        # Application service / use case (lowercase)
  infrastructure/
    repository/
      {name}_repository_impl.go                # PostgreSQL repository impl (lowercase)
  interfaces/
    http/
      {name}_handler.go                        # HTTP handler (lowercase)
      {name}_routes.go                         # Route registration (lowercase)
  {name}_test.go                               # Test skeleton (lowercase)
```

Import paths use `{module_path}/internal/modules/{package}/...`.

Each module is self-contained and isolated:

- `module.go` exposes `NewModule(db, log) *Module` and `RegisterRoutes(e)`
- The composition root in `cmd/server/main.go` wires modules together
- Modules accept `*pgxpool.Pool` and `zerolog.Logger` as dependencies
- Middleware can be injected dynamically per module via Echo groups

### Example: Property Module

Generated from `examples/property.yaml`:

- `modules/property/MODULE.md` — module metadata, fields, behaviors, REST config
- `modules/property/module.go` — `NewModule(db, log)` + `RegisterRoutes(e)`
- `modules/property/domain/entity/property.go` — `Property` struct with fields + soft delete
- `modules/property/domain/repository/property_repository.go` — `PropertyRepository` interface (Save, FindByID, FindAll, Delete)
- `modules/property/application/dto/property_dto.go` — `CreatePropertyRequest`, `UpdatePropertyRequest`, `PropertyResponse`, `PropertyListResponse`
- `modules/property/application/service/property_service.go` — `PropertyService` with Create, GetByID, List, Update, Delete
- `modules/property/infrastructure/repository/property_repository_impl.go` — `PropertyRepositoryImpl` (PostgreSQL, TODO: sqlc queries)
- `modules/property/interfaces/http/property_handler.go` — `PropertyHandler` with Echo handlers for all CRUD operations
- `modules/property/interfaces/http/property_routes.go` — `RegisterPropertyRoutes(g, h)` for route registration
- `modules/property/property_test.go` — test skeleton with suggested test cases

## Generated Modules

The following modules have been generated in `backend/internal/modules/`:

| Module   | Package  | Table      | REST Base Path  |
| -------- | -------- | ---------- | --------------- |
| Property | property | properties | /api/properties |
| Tenant   | tenant   | tenants    | /api/tenants    |
| Room     | room     | rooms      | /api/rooms      |

## Using the Generator Programmatically

```go
package main

import (
    "github.com/epmp/sdk/codegen/config"
    "github.com/epmp/sdk/codegen/internal/filesystem"
    "github.com/epmp/sdk/codegen/internal/generator"
)

func main() {
    cfg, _ := config.Load("property.yaml")

    fs := filesystem.New()
    gen := generator.New(generator.DefaultRenderer(), fs)

    gen.Generate(&generator.GenerateRequest{
        Config:     cfg,
        OutputRoot: "./internal",
        ModulePath: "github.com/epmp/backend",
    })
    // Files are generated at: ./internal/modules/{package}/...
}
```

## Validation Rules

The config validator checks:

- `version` is required
- `backend.output_root` is required
- `backend.module_path` is required
- `domain.name` must be PascalCase (e.g. `Property`)
- `domain.package` is required and must be lowercase
- `domain.table` is required
- `domain.rest.base_path` must start with `/api/` if provided
- Enum fields must have at least 2 `enum_values`

## Testing

```bash
cd tools/epmp-sdk/be/codegen
go test ./... -v
```

## Design Principles

- **Small interfaces** — `Renderer`, `Filesystem`, `TemplateLoader` are minimal
- **No third-party dependencies** — only stdlib + `gopkg.in/yaml.v3` for config
- **Renderer does not write files** — pure render → bytes
- **Filesystem does not know about templates** — pure I/O
- **Generator orchestrates** — depends on both, knows about neither's internals
- **Idempotent** — running twice produces the same output
- **Embedded templates** — no external template files needed at runtime
