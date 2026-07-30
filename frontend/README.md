# EPMP Frontend

Enterprise Property Management Platform — Frontend.

## Tech Stack

| Area         | Standard        |
| ------------ | --------------- |
| Framework    | React           |
| Language     | TypeScript      |
| Build Tool   | Vite            |
| Router       | React Router    |
| Server State | TanStack Query  |
| Forms        | React Hook Form |
| Validation   | Zod             |
| Styling      | Tailwind CSS    |
| UI Component | shadcn/ui       |
| Icons        | Lucide          |
| Tables       | TanStack Table  |
| Charts       | Recharts        |

## Structure

```
frontend/
  src/
    features/    # bounded contexts (feature-first)
      property/  # generated: Property module
      tenant/    # generated: Tenant module
      room/      # generated: Room module
    shared/      # cross-feature utilities
    layouts/     # layout components
    pages/       # route-level pages
    router/      # route configuration
    hooks/       # global hooks
    services/    # API client and external services
    styles/      # global styles
    assets/      # static assets
```

## Generated Modules

Modules are generated using the EPMP SDK frontend code generator.

| Module   | Package  | REST Base Path  | Files |
| -------- | -------- | --------------- | ----- |
| Property | property | /api/properties | 10    |
| Tenant   | tenant   | /api/tenants    | 10    |
| Room     | room     | /api/rooms      | 10    |

Each generated feature includes: types, Zod schema, API client, TanStack Query hooks, form component, table component, and list/create/edit/detail pages.

### Regenerating Modules

```bash
cd tools/epmp-sdk/fe/codegen
go build -o epmp-fe-codegen ./cmd/codegen/

./epmp-fe-codegen --config examples/property.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config examples/tenant.yaml --output ../../../../frontend/src
./epmp-fe-codegen --config examples/room.yaml --output ../../../../frontend/src
```

See `tools/epmp-sdk/fe/codegen/README.md` for full documentation.

## Feature Structure

Each feature follows:

```
features/{domain}/
  pages/
  components/
  forms/
  hooks/
  api/
  types/
```

## Running

```bash
npm install
npm run dev
```
