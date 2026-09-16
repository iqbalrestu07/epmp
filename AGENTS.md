# EPMP — Agent Entry Point

EPMP (Enterprise Property Management Platform): monorepo Go (Echo v4) + PostgreSQL 16 backend,
React 18 + TypeScript + Vite frontend, Playwright E2E, Docker Compose. Multi-tenant per
`organization_id`; mendukung banyak property model (boarding house, apartment, co-living, office, ...).

## Wajib dibaca sebelum mengubah kode

1. `.agents/rules/project-structure.md` — single source of truth untuk layout & command
2. `.agents/rules/epmp-module-patterns.md` — anatomi modul backend & feature frontend
3. `.agents/rules/multi-tenancy-boundary.md` — isolasi organization di setiap query/endpoint
4. `.agents/rules/rule-priority.md` — urutan prioritas jika rules berkonflik
5. `tools/epmp-ai/RULES.md`, `CONVENTIONS.md`, `CONTEXT.md` — EPMP AI Operating System (konteks bisnis otoritatif)

Semua file `.agents/rules/*.md` bersifat mengikat. Bahasa file campuran ID/EN — ikuti gaya file yang diedit.

## Layout singkat

| Area                                                      | Path                                                                            |
| --------------------------------------------------------- | ------------------------------------------------------------------------------- |
| Backend module (bounded context)                          | `backend/internal/modules/<module>/` (+ `MODULE.md`)                            |
| Composition root / routes `/api/v1`                       | `backend/internal/modules/modules.go`                                           |
| Shared pkg (errs, response, middleware, logger, whatsapp) | `backend/internal/pkg/`                                                         |
| SQL migrations (golang-migrate)                           | `backend/migrations/NNNNNN_name.{up,down}.sql`                                  |
| API integration harness                                   | `backend/internal/testutil/apitest.go`                                          |
| Frontend feature                                          | `frontend/src/features/<feature>/{api,hooks,types,schema,components,pages}`     |
| Route table                                               | `frontend/src/App.tsx`                                                          |
| API client (JWT + `X-Organization-ID`)                    | `frontend/src/services/api.ts`                                                  |
| E2E runner (system Chrome)                                | `frontend/run_e2e_tests.cjs` → `PAGES_TO_TEST`                                  |
| Codegen SDK                                               | `tools/epmp-sdk/{schemas,be/codegen,fe/codegen}`                                |
| Specs / ADR / audits / research logs                      | `epmp-docs/`, `epmp-docs/adr/`, `epmp-docs/audits/`, `epmp-docs/research-logs/` |

## Commands

```bash
docker compose up -d postgres                 # DB lokal (5432)
cd backend && go run ./cmd/migrate up         # migrasi
make dev-backend                              # Echo :8080  (go run ./cmd/server/...  — bukan main.go saja)
make dev-frontend                             # Vite  :3000
cd backend && gofmt -l . && go vet ./... && go build ./... && go test ./... -count=1
cd frontend && npm run lint && npm run build  # build = tsc -b && vite build (typecheck)
make test-e2e                                 # Playwright headless, harus "0 Errors Found"
make test                                     # backend + frontend + e2e
```

Integration test (`*_api_test.go`) **skip** jika Postgres tidak jalan — skip bukan pass.

## Aturan inti (ringkas; detail di `.agents/rules/`)

- Modul memiliki datanya sendiri. Cross-module lewat service interface, **tidak pernah** query tabel modul lain.
- Setiap query & endpoint difilter `organization_id` dari context (`X-Organization-ID`), tidak dari body/query param.
- Schema berubah hanya via migration up+down; tidak ada DDL manual, tidak ada SQL destruktif tanpa permintaan eksplisit.
- Logic bisnis baru wajib punya test; bug fix wajib punya regression test.
- Secrets dari env/config; jangan hardcode, jangan log data sensitif (token, NIK, nomor WA).
- Jangan mengarang aturan bisnis; jika tidak ada di spec/`tools/epmp-ai`, tanya.
- Route baru → tambahkan ke `App.tsx` **dan** `PAGES_TO_TEST`.
- Conventional Commits; jangan push tanpa diminta.

## Workflows & skills

Workflow fase (`.agents/workflows/`) dipanggil lewat skill `workflow`:
`orchestrator` → `1-research` → `2-implement` → `3-integrate` → `4-verify` → `5-commit`;
juga `quick-fix`, `refactor`, `audit`, `create-spec`.
Skill lain: `guardrails`, `code-review`, `debugging-protocol`, `adr`, `patch-cve`, `sequential-thinking`.
