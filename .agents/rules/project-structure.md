---
trigger: always_on
---

> **This file is the SINGLE SOURCE OF TRUTH for project organization.**
> All other rules and workflows that reference paths should defer to this file.

## Project Structure

**Project Type:** EPMP (Enterprise Property Management Platform) — monorepo with a Go (Echo v4) REST API, a React + Vite SPA, PostgreSQL 16, and a Go-based code generator SDK. Docker Compose for local infrastructure.

### Directory Layout

```
epmp/
├── .agents/                    # AI agent configuration (rules, skills, workflows)
├── .devin/                     # Devin CLI config (config.json, rules/, permissions)
├── AGENTS.md                   # Always-on entry point for AI agents (concise; points to .agents/rules)
├── Makefile                    # Root commands: test, test-e2e, dev-backend, dev-frontend, build
├── docker-compose.yml          # postgres, migrate (one-shot), backend, frontend (nginx)
├── .env.example                # Root env template (DB_*, BACKEND_PORT, FRONTEND_PORT, VITE_API_BASE_URL)
├── E2E_TESTING.md              # Playwright E2E policy & how-to
├── PROGRESS.md                 # Feature progress tracker
│
├── backend/                    # Go module: github.com/epmp/backend
│   ├── cmd/
│   │   ├── server/             # main.go, server.go (Echo start/graceful shutdown), bootstrap.go (DI)
│   │   └── migrate/            # golang-migrate runner: up | down | version | force | create | watch
│   ├── configs/                # config.go + config.yaml (env-driven runtime config)
│   ├── internal/
│   │   ├── database/postgres/  # pgxpool connection
│   │   ├── modules/            # Bounded contexts — one folder per domain (see below)
│   │   │   └── modules.go      # Register(): composition root, mounts /api/v1 + auth middleware
│   │   ├── pkg/                # Shared packages: errs, response, middleware, logger, health,
│   │   │                       #   types, uid, email, websocket, whatsapp
│   │   └── testutil/           # apitest.go — API integration harness (real Postgres, httptest)
│   ├── migrations/             # golang-migrate SQL: NNNNNN_name.up.sql + .down.sql (+ README.md)
│   ├── Makefile                # run, build, test, lint, migrate-*
│   └── README.md
│
├── frontend/                   # React 18 + TypeScript + Vite + Tailwind (npm)
│   ├── src/
│   │   ├── App.tsx             # Route table (every page route lives here)
│   │   ├── components/         # ErrorBoundary, ProtectedRoute, PermissionGuard, ui/*
│   │   ├── features/           # Feature modules — one folder per domain (see below)
│   │   ├── layouts/            # AuthLayout, MainLayout
│   │   ├── pages/              # DashboardPage, LandingPage
│   │   ├── services/api.ts     # Fetch wrapper: BASE_URL=/api/v1, JWT + X-Organization-ID headers
│   │   ├── lib/, utils/, styles/
│   │   └── main.tsx
│   ├── run_e2e_tests.cjs       # Playwright runner (system Chrome) — PAGES_TO_TEST list
│   ├── vite.config.ts, tsconfig.json, tailwind.config.js, nginx.conf, Dockerfile
│   └── package.json            # scripts: dev, build (tsc -b && vite build), lint, test:e2e[:gui]
│
├── tools/
│   ├── epmp-sdk/
│   │   ├── schemas/*.yaml      # Domain definitions — input for BOTH generators
│   │   ├── be/codegen/         # Go backend module generator (epmp-codegen)
│   │   └── fe/codegen/         # React feature generator (epmp-fe-codegen)
│   └── epmp-ai/                # EPMP AI Operating System docs (RULES, CONVENTIONS, WORKFLOW, ...)
│
└── epmp-docs/                  # Platform documentation (EPMP-001 … EPMP-013)
    ├── epmp-003.md             # Architecture Overview
    ├── epmp-010.md             # Engineering Platform Standard (stack, branches, commits)
    ├── epmp-011.md             # Backend Architecture Standard (layers, modules, errors)
    ├── epmp-012.md             # AI Engineering Standard
    ├── epmp-013-notification.md
    ├── research-logs/          # Agent research logs (created by /1-research)
    ├── audits/                 # Code review findings (created by /audit)
    └── adr/                    # Architecture Decision Records (created by ADR skill)
```

### Backend Module Convention

Each bounded context lives under `backend/internal/modules/<package>/` and is wired in `modules.go`:

```
internal/modules/<package>/
├── MODULE.md                       # Module metadata: table, base path, fields, behaviors
├── module.go                       # NewModule(db, log) *Module + RegisterRoutes(*echo.Group)
├── entity/<name>.go                # Domain entity (struct + invariants)
├── dto/<name>_dto.go               # Create/Update Request, Response, ListResponse
├── repository/<name>_repository.go       # Interface (contract) — org-scoped signatures
├── repository/<name>_repository_impl.go  # PostgreSQL impl (pgx)
├── service/<name>_service.go       # Application service / use cases
├── delivery/http/<name>_handler.go # Echo handlers (bind → validate → service → response)
├── delivery/http/<name>_routes.go  # Register<Name>Routes(g, h)
└── <name>_api_test.go              # API integration test (uses internal/testutil)
```

- Package names are singular, lowercase, no underscores (`roomtype`, `tenantcontact`).
- Routes are mounted under `/api/v1/<plural-resource>` by `modules.Register`.
- Modules are frequently **generated** by `tools/epmp-sdk/be/codegen` from `tools/epmp-sdk/schemas/<name>.yaml`. Files with the header `// Code generated by epmp-sdk. DO NOT EDIT.` should be regenerated, not hand-edited, unless the generator cannot express the change (then remove the header).

### Frontend Feature Convention

Each domain lives under `frontend/src/features/<feature>/`:

```
features/<feature>/
├── api/index.ts          # fetch<Name>s, fetch<Name>ById, create/update/delete<Name>
├── hooks/index.ts        # TanStack Query hooks: use<Name>s, use<Name>, useCreate<Name>, ...
├── types/index.ts        # Entity, ListResponse, Create/Update Request, QueryParams
├── schema/index.ts       # Zod schemas (create/update) + inferred form types
├── components/           # <Name>Table.tsx, <Name>Form.tsx, ...
└── pages/                # <Name>ListPage, <Name>CreatePage, <Name>EditPage, <Name>DetailPage
```

- Import alias `@/` → `frontend/src/`.
- New page routes go in `frontend/src/App.tsx` **and** in `PAGES_TO_TEST` in `frontend/run_e2e_tests.cjs`.

### Key Conventions

1. **Monorepo, two runtimes** — Go for `backend/` and `tools/`, Node (npm, `package-lock.json`) for `frontend/`. Never introduce pnpm/yarn.
2. **Schema-as-code** — every DB change is a `backend/migrations/NNNNNN_*.up.sql` + `.down.sql` pair, created via `cd backend && go run ./cmd/migrate create <name>`. Never alter the schema by hand.
3. **Generator first** — CRUD boilerplate MUST come from `tools/epmp-sdk` when the schema can express it (EPMP-010 §5). Do not hand-write what the generator produces.
4. **Docker-first infra** — PostgreSQL runs via `docker compose up -d postgres`; backend/frontend run natively in dev (`make dev-backend`, `make dev-frontend`) or via compose.
5. **Ports** — backend `8080`, frontend `3000`, Postgres `5432`. Health check: `GET /health`.
6. **Root Makefile is the entry point** for tests: `make test-backend`, `make test-frontend`, `make test-e2e`, `make test`.
7. **Documentation lives in `epmp-docs/`** — no root `docs/` folder. Agent artefacts (research logs, audits, ADRs) go under `epmp-docs/`.

### Related Principles

- EPMP Module Patterns @epmp-module-patterns.md
- Multi-Tenancy Boundary @multi-tenancy-boundary.md
- Code Organization Principles @code-organization-principles.md
