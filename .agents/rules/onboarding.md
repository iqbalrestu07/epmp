---
trigger: manual
description: Setup awal untuk anggota tim baru yang baru join ke proyek EPMP
---

# Onboarding — Setup Lingkungan Pengembangan

Panduan ini wajib dibaca dan diikuti oleh setiap anggota tim baru sebelum mulai berkontribusi ke repositori `epmp`.

## Prasyarat

- [Docker Desktop](https://www.docker.com/) sudah terpasang dan berjalan (untuk PostgreSQL 16)
- Go **1.26+** (`go version`)
- Node.js **22+** dan npm (`node --version`) — gunakan **npm**, bukan pnpm/yarn
- Google Chrome terpasang (dipakai Playwright E2E via `channel: 'chrome'`, tidak perlu `npx playwright install`)
- (Opsional) `staticcheck`, `gosec` di `~/go/bin` untuk lint backend

## Urutan Setup (Jalankan Sekali)

```bash
# 1. Clone repositori
git clone <repo-url>
cd epmp

# 2. Salin file environment
cp .env.example .env
# Nilai default sudah cocok untuk lokal (DB postgres/postgres, port 8080/3000/5432)

# 3. Jalankan PostgreSQL
docker compose up -d postgres

# 4. Jalankan migrasi skema
cd backend && go run ./cmd/migrate up && cd ..

# 5. Install dependency frontend
cd frontend && npm install && cd ..

# 6. Jalankan backend & frontend (dua terminal)
make dev-backend      # http://localhost:8080  (health: /health)
make dev-frontend     # http://localhost:3000
```

Alternatif full-Docker: `docker compose up -d --build` (postgres → migrate → backend → frontend/nginx).

## Verifikasi Setup

```bash
make test-backend     # go test ./... (integration test skip otomatis jika DB tidak aktif)
make test-frontend    # tsc -b && vite build
make test-e2e         # Playwright, butuh backend :8080 dan frontend :3000 aktif
```

## Konvensi Penting

| Aspek                  | Aturan                                                                                                                                               |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Schema DB**          | Hanya lewat `backend/migrations/` (`go run ./cmd/migrate create <nama>`), selalu ada `.up.sql` + `.down.sql`. Lihat @schema-documentation-mandate.md |
| **Zero Data-Loss**     | Tidak ada DDL destruktif tanpa permintaan eksplisit; jangan `docker compose down -v`. Lihat @deployment-and-data-safety-mandate.md                   |
| **Multi-tenant**       | Semua query tenant data wajib filter `organization_id` dari header `X-Organization-ID`. Lihat @multi-tenancy-boundary.md                             |
| **Boilerplate**        | CRUD module/feature dibuat dari `tools/epmp-sdk` (schemas YAML → codegen), bukan ditulis manual                                                      |
| **Route baru**         | Daftarkan di `frontend/src/App.tsx` **dan** `PAGES_TO_TEST` di `frontend/run_e2e_tests.cjs`                                                          |
| **Commit format**      | Conventional Commits: `feat(property): add room availability service` (EPMP-010 §12)                                                                 |
| **Branch kerja**       | `feature/{module}/{feature}`, `fix/{issue}` dari `develop`; `main` untuk rilis (EPMP-010 §10)                                                        |
| **Definition of Done** | Build + lint + unit/integration test + E2E (jika sentuh UI) + dokumentasi diperbarui                                                                 |

## Referensi Dokumen

| Dokumen              | Lokasi                                                          | Keterangan                                           |
| -------------------- | --------------------------------------------------------------- | ---------------------------------------------------- |
| Arsitektur           | `epmp-docs/epmp-003.md`, `epmp-docs/epmp-011.md`                | Overview & Backend Architecture Standard             |
| Engineering Standard | `epmp-docs/epmp-010.md`                                         | Stack, branch, commit, code quality                  |
| AI Standard          | `epmp-docs/epmp-012.md`, `tools/epmp-ai/*`                      | Aturan kerja AI agent (RULES, CONVENTIONS, WORKFLOW) |
| Spesifikasi modul    | `epmp-docs/epmp-0xx.md`, `backend/internal/modules/*/MODULE.md` | Domain & field per modul                             |
| E2E                  | `E2E_TESTING.md`                                                | Kebijakan & cara jalankan Playwright                 |
| Progress             | `PROGRESS.md`                                                   | Status fitur                                         |
| Audit Log            | `epmp-docs/audits/`                                             | Hasil review kode (`/audit`)                         |
| Research Log         | `epmp-docs/research-logs/`                                      | Catatan riset per fitur (`/1-research`)              |
| ADR                  | `epmp-docs/adr/`                                                | Architecture Decision Records                        |

## Cara Menggunakan AI Agent

AI Coding Assistant sudah dikonfigurasi dengan konteks proyek ini (`AGENTS.md`, `.agents/rules/`, `.agents/skills/`). Gunakan slash commands:

```
/orchestrator   ← Fitur baru end-to-end (research → implement → integrate → verify → commit)
/1-research     ← Planning sebelum coding
/2-implement    ← TDD cycle
/3-integrate    ← Integration test dengan Postgres nyata (testutil)
/4-verify       ← Validasi penuh sebelum commit
/5-commit       ← Commit dengan format standar
/audit          ← Review kode
/quick-fix      ← Hotfix kecil
/refactor       ← Refactor terarah
/create-spec    ← Dokumentasi fitur yang belum ada di epmp-docs (out-of-scope)
/patch-cve      ← Triage & patch dependency vulnerability
```

> Format yang benar: `/nama-command`. Di Windsurf/Antigravity juga bisa `@[/nama-command]`.

## Pertanyaan Umum

**Q: `make test-backend` menampilkan banyak `SKIP`?**
A: Integration test skip otomatis bila Postgres tidak bisa dihubungi. Jalankan `docker compose up -d postgres` lalu `go run ./cmd/migrate up`.

**Q: E2E gagal semua dengan connection refused?**
A: Backend (`:8080`) dan frontend (`:3000`) harus sudah berjalan sebelum `make test-e2e`.

**Q: Ingin menambah kolom baru?**
A: `cd backend && go run ./cmd/migrate create add_<kolom>_to_<tabel>`, isi `.up.sql`/`.down.sql`, lalu update `MODULE.md`, `backend/migrations/README.md`, `tools/epmp-sdk/schemas/<modul>.yaml`, dan Zod schema frontend.

**Q: Migrasi gagal dan status `dirty`?**
A: Periksa state DB manual, perbaiki migrasi, lalu `go run ./cmd/migrate force <version>` **hanya di lokal**.
