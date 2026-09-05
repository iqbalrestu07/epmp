# EPMP-009

# Repository & Solution Structure

---

```text
Document ID    : EPMP-009
Document Name  : Repository & Solution Structure
Version        : 1.1.0
Status         : Implemented & Active
Owner          : Software Architecture Team
Dependencies   : EPMP-001 ~ EPMP-008
Referenced By  : Backend, Frontend, DevOps, CI/CD, Testing
```

---

# 1. Purpose

Dokumen ini mendefinisikan struktur repository dan organisasi source code EPMP aktual.

Tujuannya adalah:

- menjaga konsistensi struktur proyek,
- memudahkan onboarding developer,
- memudahkan AI Coding Agent memahami konteks,
- memisahkan business domain dari implementation detail,
- memungkinkan pertumbuhan proyek tanpa reorganisasi besar.

Dokumen ini tidak membahas implementasi detail setiap modul, tetapi menetapkan aturan bagaimana kode diorganisasi.

---

# 2. Repository Strategy

EPMP versi pertama menggunakan **Monorepo**.

Repository menyimpan:

- Backend (Go + Echo + whatsmeow)
- Frontend (React + Vite + Three.js WebGL Canvas)
- Documentation (`epmp-docs/`, `README.md`, `PROGRESS.md`)
- Database Migration (`backend/migrations/*.sql`)
- End-to-End Test Suite (`frontend/tests/e2e/*.spec.ts`)
- Infrastructure & Docker (`docker-compose.yml`)

dalam satu repository Git.

### Alasan memilih Monorepo

- Konsistensi versi antar aplikasi.
- Dokumentasi dan kode berada pada satu sumber kebenaran.
- Pull Request dapat mencakup perubahan lintas layer.
- Mempermudah AI Coding Agent memahami keseluruhan konteks.
- Sederhana untuk tim kecil hingga menengah.

---

# 3. High-Level Repository Layout

```text
epmp/
├── epmp-docs/           # Dokumentasi arsitektur & domain komprehensif
├── backend/             # Monolith modular Go (Echo framework)
│   ├── cmd/server/      # Entrypoint aplikasi utama
│   ├── internal/modules/# Bounded contexts domain
│   ├── migrations/      # Migrasi database PostgreSQL (termasuk trigger auto-settlement)
│   └── pkg/             # Shared utilities
├── frontend/            # Single Page Application React 19 + TypeScript + Three.js
│   ├── src/features/    # Bounded contexts frontend (pages, components, api, hooks)
│   └── tests/e2e/       # Playwright E2E test suite (40+ dashboard routes)
├── docker-compose.yml   # Kontainer PostgreSQL lokal
├── Makefile             # Automation runner (test-e2e, dev, migrate)
├── PROGRESS.md          # Log status implementasi & changelog fitur
└── README.md            # Dokumentasi arsitektur sistem level tinggi
```

---

# 4. Documentation Structure

Seluruh dokumentasi mengikuti struktur yang telah kita bangun.

```text
epmp-docs/
├── epmp-001.md ~ epmp-012.md   # Core Foundation, Business, & Architecture
└── epmp-013-notification.md     # Communication & WhatsApp Gateway Specification
```

Dokumentasi merupakan bagian dari source code dan wajib diperbarui bersama perubahan implementasi.

---

# 5. Backend Structure

Backend diorganisasi berdasarkan **bounded context**, bukan berdasarkan layer teknis semata.

```text
backend/
├── cmd/
│   └── server/          # main.go, server initialization & dependency injection
├── internal/
│   ├── app/             # Application container & module registration
│   ├── middleware/      # Auth JWT, CORS, Logger, Error recovery
│   └── modules/         # Bounded Contexts
│       ├── property/    # Property, Building, Floor, Room, Bed Template
│       ├── tenant/      # Tenant directory & KYC
│       ├── reservation/ # Unit booking & reservations
│       ├── contract/    # Lease contracts & terms
│       ├── occupancy/   # Check-in, check-out, occupancy log
│       ├── billing/     # Invoicing, multi-currency fees
│       ├── payment/     # Payments & auto-paid settlement
│       ├── deposit/     # Security deposit management
│       ├── communication/# whatsmeow WhatsApp gateway, blast, audit logs
│       ├── asset/       # Inventory items, asset assignments, inspections
│       └── maintenance/ # Work orders & technician management
├── migrations/          # PostgreSQL schema migrations (000001 - 000038)
└── pkg/                 # Database connection, JWT helper, validator
```

### Penjelasan

- `cmd/server/` → entry point aplikasi.
- `internal/modules/` → implementasi domain modular yang tidak diekspor.
- `communication/` → integrasi native whatsmeow multi-device WhatsApp, container SQLite session, template resolver, dan auto phone check.
- `migrations/` → file SQL up/down untuk skema tabel dan database triggers.

---

# 6. Internal Module Structure

Setiap bounded context memiliki pola yang sama.

Contoh untuk `property`:

```text
property/

application/
domain/
infrastructure/
interfaces/

README.md
```

## application/

Berisi:

- Use Case
- Command
- Query
- Handler
- DTO
- Mapper

## domain/

Berisi:

- Entity
- Aggregate
- Repository Interface
- Value Object
- Domain Service
- Domain Event

## infrastructure/

Berisi implementasi teknis:

- PostgreSQL Repository
- External API
- Cache
- Storage

## interfaces/

Berisi:

- HTTP Controller
- REST Endpoint
- Request
- Response

Pendekatan ini menjaga agar Domain tetap bersih dari detail implementasi.

---

# 7. Frontend Structure

Frontend menggunakan **feature-first architecture** agar selaras dengan bounded context backend.

```text
frontend/

src/

features/
shared/
layouts/
pages/
router/
hooks/
services/
styles/
assets/
```

Setiap fitur memiliki struktur sendiri:

```text
features/property/

pages/
components/
forms/
hooks/
api/
types/
```

Dengan demikian, developer frontend dapat bekerja pada satu domain tanpa bergantung pada domain lain.

---

# 7. Frontend Structure

Frontend menggunakan **feature-first architecture** agar selaras dengan bounded context backend.

```text
frontend/
├── src/
│   ├── components/      # Shared UI primitives (Shadcn/Tailwind, Modals, Buttons)
│   ├── features/        # Feature modules
│   │   ├── property/    # Properties, Buildings, Floors, Rooms, Bed Templates, 3D WebGL Canvas
│   │   ├── tenant/      # Tenant directory & onboarding
│   │   ├── reservation/ # Unit booking & calendar
│   │   ├── contract/    # Contract agreements & renewals
│   │   ├── occupancy/   # Check-in, check-out, living logs
│   │   ├── invoice/     # Billing statements, invoice generation
│   │   ├── payment/     # Payment recording & settlement
│   │   ├── deposit/     # Security deposit & deductions
│   │   ├── communication/# WhatsApp devices (QR pairing), blast campaigns, template editor
│   │   ├── asset/       # Inventory items, inspections
│   │   ├── workorder/   # Maintenance tickets & technician assignments
│   │   └── dashboard/   # Executive summary & quick action shortcuts
│   ├── hooks/           # Custom React hooks (TanStack query, auth, currency)
│   ├── lib/             # API client (Axios/fetch with JWT interceptor), utils, currency formatting
│   ├── types/           # Global TypeScript type declarations
│   └── App.tsx          # Master routing & Layout providers
└── tests/
    └── e2e/             # Playwright E2E tests for all 40 dashboard routes
```

Setiap fitur memiliki struktur modular:
```text
features/<feature-name>/
├── api/                 # API calls ke endpoint backend
├── components/          # Komponen UI spesifik (misal: WhatsAppQRModal, BuildingCanvas3D)
├── hooks/               # React Query hooks (useCreate, useUpdate, useList)
├── pages/               # ListPage, CreatePage, EditPage, DetailPage
└── types/               # TypeScript interfaces & DTOs
```

---

# 8. Shared Components Policy

Komponen bersama hanya dibuat apabila benar-benar digunakan oleh lebih dari satu fitur.

Contoh:
- `Button`, `Modal`, `Dialog`, `Input`, `Select`, `Badge`
- `CurrencySelector` & formatters (IDR, USD, EUR, SGD, MYR)
- `StatusBadge` (Available: Hijau, Reserved: Oranye, Occupied: Merah, Maintenance: Biru)

Komponen yang hanya digunakan oleh satu fitur (misal: `Canvas3D`, `QRCodeDisplay`) tetap berada di dalam folder fitur tersebut.

---

# 9. API Specification & Integration

REST API terpusat di backend Go pada prefix `/api/v1/...`.
Frontend berinteraksi melalui API client seragam dengan penanganan token JWT otomatis, normalisasi tanggal, dan human-readable identifiers.

---

# 10. Database Migrations

Skema database dikelola melalui migration files terurut di `backend/migrations/`:

```text
backend/migrations/
├── 000001_create_initial_schema.up.sql
├── ...
├── 000035_add_communication_schema.up.sql
├── 000036_add_bed_templates.up.sql
├── 000037_add_currency_support.up.sql
└── 000038_sync_invoice_status.up.sql (Trigger auto-paid pada settlement payment)
```

Model domain diimplementasikan di layer Go, sedangkan database migration menjamin integritas relasional, foreign keys, indeks performa, dan trigger database.

---

# 11. Infrastructure

Folder `infrastructure/` berisi konfigurasi operasional.

Contoh:

```text
docker/
nginx/
compose/
deployment/
monitoring/
```

Folder ini tidak boleh berisi business logic.

---

# 12. Scripts

Folder `scripts/` berisi utilitas otomatisasi.

Contoh:

- bootstrap environment
- generate code
- lint
- test
- build
- release

Semua script harus idempotent dan terdokumentasi.

---

# 13. Testing Structure

Struktur pengujian mengikuti struktur modul.

```text
test/

property/
contract/
reservation/
finance/
```

Jenis pengujian yang direncanakan:

- Unit Test
- Integration Test
- API Test
- End-to-End Test

---

# 14. Dependency Rules

Agar arsitektur tetap bersih, berlaku aturan berikut:

- Domain tidak boleh bergantung pada Infrastructure.
- Application tidak boleh mengakses database secara langsung.
- UI tidak boleh memuat business rule.
- Shared module tidak boleh menjadi tempat "menumpuk" logika yang tidak jelas kepemilikannya.
- Setiap bounded context hanya mengekspos antarmuka yang memang diperlukan oleh context lain.

---

# 15. Naming Convention

### Folder

Gunakan huruf kecil dengan pemisah yang konsisten.

```text
property
reservation
contract
```

### Package Go

Menggunakan nama domain, bukan singkatan yang sulit dipahami.

### React Component

PascalCase.

```text
PropertyList.tsx
TenantDetail.tsx
InvoiceTable.tsx
```

### File

Gunakan nama yang mendeskripsikan isi dan tanggung jawabnya.

---

# 16. AI-Friendly Repository Rules

Repository dirancang agar mudah dipahami AI.

Aturannya:

- Satu folder = satu konteks bisnis.
- README lokal pada setiap modul menjelaskan tujuan modul.
- Hindari folder "misc", "helper", atau "utils" sebagai tempat serba ada.
- Dokumentasi domain ditempatkan sedekat mungkin dengan implementasinya bila relevan.
- Struktur konsisten di semua bounded context.

---

# 17. Evolution Strategy

Repository harus mampu berkembang tanpa reorganisasi besar.

Tahapan evolusi:

1. Modular Monolith.
2. Penambahan bounded context.
3. Pemisahan service bila diperlukan.
4. Pemisahan repository hanya jika ada alasan operasional yang kuat.

Dengan strategi ini, perubahan arsitektur di masa depan tidak mengharuskan penulisan ulang seluruh kode.

---

# Closing Statement

Repository & Solution Structure adalah kontrak organisasi source code EPMP. Konsistensi struktur akan memudahkan kolaborasi antar developer, mempercepat onboarding, dan meningkatkan efektivitas AI Coding Agent dalam menghasilkan implementasi yang sesuai dengan domain bisnis.

---

# 📌 Catatan sebagai Chief Software Architect

Saya ingin memberikan satu rekomendasi strategis yang menurut saya akan menjadi investasi terbesar untuk kualitas proyek.

## Kita perlu menetapkan **Technology Baseline** sebelum mulai menulis satu baris kode.

Artinya, sebelum membuat dokumen seperti API Guidelines atau Coding Standards, kita menyepakati seluruh teknologi inti yang akan digunakan. Contohnya:

### Backend

- Go 1.25 (atau versi stabil yang kita sepakati saat implementasi dimulai)
- Router (misalnya `echo` atau `gin`)
- ORM atau SQL builder (misalnya `sqlc`, `ent`, `gorm`, atau kombinasi `pgx` + query builder)
- Dependency Injection (manual atau `wire`)
- Logger
- Validator
- Configuration library
- Migration tool

### Frontend

- React
- Vite
- TypeScript
- React Router
- State management (jika diperlukan)
- Data fetching (misalnya TanStack Query)
- Form library
- Validation library
- UI component library atau design system

### Database

- PostgreSQL
- Redis (opsional pada MVP)
- Object Storage (fase berikutnya)
