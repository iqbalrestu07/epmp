---
trigger: always_on
---

> **This file is the SINGLE SOURCE OF TRUTH for project organization.**
> All other rules and workflows that reference paths should defer to this file.

## Project Structure

**Project Type:** Directus Headless CMS — single-app, Docker-based deployment with custom endpoint extensions.

### Directory Layout

```
kes-ref-screening-directus/
├── .agents/                    # AI agent configuration (rules, skills, workflows)
├── docs/
│   ├── PRD/                    # Product Requirements Document
│   ├── TSD/                    # Technical Specification Document
│   ├── spec-document/          # Flow & API Specifications (flow-openapi-*.md)
│   ├── audits/                 # Laporan hasil code audit
│   ├── research_logs/          # Catatan penelitian & keputusan teknis
│   └── schema.md               # Blueprint & dokumentasi skema database (SurveyJS)
├── extensions/                 # Directus custom extensions
│   ├── .gitkeep
│   └── survey-forms/           # Endpoint extension: SurveyJS JSON transformer
│       ├── index.js            # REST API endpoint (ESM JavaScript)
│       └── package.json        # Extension manifest (directus:extension)
├── scripts/                    # Utilitas CLI & manajemen skema
│   ├── migrations/             # Sequential migration files (dieksekusi oleh run-migrations.js)
│   │   └── 001-initialization.js  # Inisialisasi 6 koleksi awal
│   ├── fix-display-templates.js   # Perbaikan UI relasi di Directus Admin
│   ├── run-migrations.js          # Migration runner (idempotent, safe to re-run)
│   └── seed-samples.js            # Dynamic JSON seeder untuk data sampel form
├── database/                   # SQLite database files (gitignored)
├── uploads/                    # Uploaded assets (gitignored)
├── docker-compose.yml          # Docker services (Directus + Trivy scanner)
├── Dockerfile                  # Custom Directus image with pnpm overrides
├── package.json                # Root package.json (pnpm overrides for security patches)
├── pnpm-lock.yaml
└── readme.md
```

### Extension Structure Convention

Each Directus extension lives under `extensions/<extension-name>/`:

```
extensions/<extension-name>/
├── package.json                # Must declare directus:extension type & host
└── index.js                    # Single entry point (ESM default export)
```

**Extension package.json format:**
```json
{
  "name": "<extension-name>",
  "version": "1.0.0",
  "type": "module",
  "directus:extension": {
    "type": "endpoint",
    "path": "index.js",
    "source": "index.js",
    "host": "^11.0.0"
  }
}
```

### Key Conventions

1. **No build step** — extensions use plain ESM JavaScript (no TypeScript, no bundler)
2. **No `apps/` directory** — this is a single Directus instance, not a monorepo
3. **No feature-based vertical slicing** — Directus handles CRUD/routing/auth; extensions are thin API adapters
4. **Schema-as-code** — `scripts/migrations/` berisi sequential migration files, dijalankan via `scripts/run-migrations.js` (bukan GUI)
5. **Docker-first** — all infrastructure defined in `docker-compose.yml`
6. **Scripts are run with `--env-file=.env`** — semua script CLI membutuhkan env file: `node --env-file=.env scripts/<name>.js`

### Related Principles
- Directus Extension Patterns @directus-extension-patterns.md
- Code Organization Principles @code-organization-principles.md
