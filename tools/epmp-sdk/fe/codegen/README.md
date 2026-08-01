# EPMP SDK — Frontend Code Generator

Code generation toolkit for EPMP frontend modules.

## Overview

The frontend codegen tool generates React + TypeScript feature modules from a single YAML config file. It produces types, Zod schemas, API clients, TanStack Query hooks, form/table components, and CRUD pages following the EPMP-010 frontend tech stack.

```
epmp-sdk/fe/codegen/
  cmd/codegen/          # CLI entry point
  config/               # Config schema, loader, validator
  internal/
    renderer/           # Template rendering (text/template)
    filesystem/         # File I/O operations
    generator/          # Frontend module generator
      templates/        # Embedded TypeScript/React templates
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

Same architecture as the backend generator:
- **Renderer** loads and renders Go `text/template` templates. Returns bytes. Does not write files.
- **Filesystem** handles directory creation and file writing. Does not know about templates.
- **Generator** orchestrates both.

## Installation

```bash
cd tools/epmp-sdk/fe/codegen
go build -o epmp-fe-codegen ./cmd/codegen/
```

## Usage

### Generate a Module

```bash
./epmp-fe-codegen --config ../../schemas/property.yaml --output ../../../../frontend/src
```

### Dry Run (no files written)

```bash
./epmp-fe-codegen --config ../../schemas/property.yaml --output ./output --dry-run
```

### Flags

| Flag        | Description                          | Required |
| ----------- | ------------------------------------ | -------- |
| `--config`  | Path to the module config YAML file  | Yes      |
| `--output`  | Output root directory                | No*      |
| `--base-url`| API base URL                         | No*      |
| `--dry-run` | Render without writing files         | No       |

\* Falls back to `frontend.output_root` and `frontend.base_url` from config.

## Config File Format

```yaml
version: "1.0"
dry_run: false

frontend:
  output_root: "./src"
  base_url: "/api"

domain:
  name: Property          # PascalCase entity name
  package: property       # lowercase feature directory name
  table: properties       # database table name
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
    pagination: true
    search: true

  rest:
    base_path: /api/properties
    auth_required: true
    operations:
      - create
      - read
      - update
      - delete
      - list
```

### Field Types → TypeScript Mapping

| Config Type  | TypeScript Type |
| ------------ | --------------- |
| `uuid`       | `string`        |
| `string`     | `string`        |
| `int`        | `number`        |
| `bool`       | `boolean`       |
| `text`       | `string`        |
| `enum`       | `string`        |
| `timestamp`  | `string`        |
| `decimal`    | `number`        |

## Generated Artifacts

For each module, the generator produces 10 files:

```
features/{package}/
  types/index.ts                          # TypeScript interfaces
  schema/index.ts                         # Zod validation schemas
  api/index.ts                            # API client functions
  hooks/index.ts                          # TanStack Query hooks
  components/
    {Name}Form.tsx                        # React Hook Form + Zod form
    {Name}Table.tsx                       # TanStack Table data table
  pages/
    {Name}ListPage.tsx                    # List page with search + pagination
    {Name}CreatePage.tsx                  # Create page
    {Name}EditPage.tsx                    # Edit page
    {Name}DetailPage.tsx                  # Detail page with delete
```

### Artifact Details

| Artifact | Tech | Description |
| -------- | ---- | ----------- |
| **types** | TypeScript | Entity, request, response, and query param interfaces |
| **schema** | Zod | Validation schemas for create/update forms |
| **api** | fetch | CRUD API client functions using `@/services/api` |
| **hooks** | TanStack Query | `use{Name}s`, `use{Name}`, `useCreate{Name}`, `useUpdate{Name}`, `useDelete{Name}` |
| **form** | React Hook Form + Zod | Reusable form component with validation |
| **table** | TanStack Table | Sortable data table with row click handler |
| **list page** | React Router | List page with search, pagination, and navigation |
| **create page** | React Router | Create form page |
| **edit page** | React Router | Edit form page with data loading |
| **detail page** | React Router | Detail view with edit/delete actions |

## Generated Modules

The following modules have been generated in `frontend/src/features/`:

| Module   | Package  | REST Base Path      | Files |
| -------- | -------- | ------------------- | ----- |
| Property | property | /api/properties     | 10    |
| Tenant   | tenant   | /api/tenants        | 10    |
| Room     | room     | /api/rooms          | 10    |

### Regenerating Modules

```bash
cd tools/epmp-sdk/fe/codegen
go build -o epmp-fe-codegen ./cmd/codegen/

./epmp-fe-codegen --config ../../schemas/property.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/tenant.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config ../../schemas/room.yaml --output ../../../../frontend/src
```

## Using the Generator Programmatically

```go
package main

import (
    "github.com/epmp/sdk/fe-codegen/config"
    "github.com/epmp/sdk/fe-codegen/internal/filesystem"
    "github.com/epmp/sdk/fe-codegen/internal/generator"
)

func main() {
    cfg, _ := config.Load("../../schemas/property.yaml")

    fs := filesystem.New()
    gen := generator.New(generator.DefaultRenderer(), fs)

    gen.Generate(&generator.GenerateRequest{
        Config:     cfg,
        OutputRoot: "./src",
        BaseURL:    "/api",
    })
}
```

## Feature Directory Structure

Generated modules follow the EPMP-009 feature-first architecture:

```
src/features/{package}/
  types/       # TypeScript type definitions
  schema/      # Zod validation schemas
  api/         # API client functions
  hooks/       # TanStack Query hooks
  components/  # Form and Table components
  pages/       # Route-level pages (List, Create, Edit, Detail)
```

## Testing

```bash
cd tools/epmp-sdk/fe/codegen
go test ./... -v
```

## Design Principles

- Same as backend: small interfaces, no third-party deps (except yaml.v3 for config)
- Renderer does not write files — pure render → bytes
- Filesystem does not know about templates — pure I/O
- Generator orchestrates — depends on both, knows about neither's internals
- Idempotent — running twice produces the same output
- Embedded templates — no external template files needed at runtime
- Follows EPMP-010 frontend tech stack exactly
