# EPMP-009

# Repository & Solution Structure

---

```text
Document ID    : EPMP-009
Document Name  : Repository & Solution Structure
Version        : 1.0.0
Status         : Draft
Owner          : Software Architecture Team
Dependencies   : EPMP-001 ~ EPMP-008
Referenced By  : Backend, Frontend, DevOps, CI/CD, Testing
```

---

# 1. Purpose

Dokumen ini mendefinisikan struktur repository dan organisasi source code EPMP.

Tujuannya adalah:

- menjaga konsistensi struktur proyek,
- memudahkan onboarding developer,
- memudahkan AI Coding Agent memahami konteks,
- memisahkan business domain dari implementation detail,
- memungkinkan pertumbuhan proyek tanpa reorganisasi besar.

Dokumen ini tidak membahas implementasi detail setiap modul, tetapi menetapkan aturan bagaimana kode harus diorganisasi.

---

# 2. Repository Strategy

EPMP versi pertama menggunakan **Monorepo**.

Repository akan menyimpan:

- Backend (Go)
- Frontend (React)
- Documentation
- Database Migration
- Infrastructure Configuration
- API Specification
- Automation Scripts

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

├── docs/
│
├── backend/
│
├── frontend/
│
├── database/
│
├── infrastructure/
│
├── scripts/
│
├── tools/
│
├── .github/
│
├── docker/
│
└── README.md
```

---

# 4. Documentation Structure

Seluruh dokumentasi mengikuti struktur yang telah kita bangun.

```text
docs/

README.md

01-foundation/

02-business/

03-architecture/

04-domain/

05-modules/

06-api/

07-database/

08-ui/

09-engineering/

10-devops/

11-testing/

12-adr/
```

Dokumentasi merupakan bagian dari source code dan wajib diperbarui bersama perubahan implementasi.

---

# 5. Backend Structure

Backend diorganisasi berdasarkan **bounded context**, bukan berdasarkan layer teknis semata.

Contoh:

```text
backend/

cmd/
internal/

property/
reservation/
contract/
finance/
asset/
maintenance/
tenant/
configuration/
shared/

pkg/

configs/

migrations/

test/
```

### Penjelasan

- `cmd/` → entry point aplikasi.
- `internal/` → implementasi utama yang tidak diekspor.
- Setiap folder domain memiliki struktur internal sendiri (application, domain, infrastructure, interface).
- `shared/` berisi komponen lintas domain yang telah disetujui.
- `pkg/` hanya untuk utilitas generik yang benar-benar reusable.

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

# 8. Shared Components Policy

Komponen bersama hanya dibuat apabila benar-benar digunakan oleh lebih dari satu fitur.

Contoh:

- Button
- Modal
- Table
- Input
- Date Picker
- Dialog

Komponen yang hanya digunakan oleh satu fitur tetap berada di dalam folder fitur tersebut.

---

# 9. API Specification Location

Seluruh spesifikasi API ditempatkan di:

```text
docs/06-api/
```

Implementasi backend dan frontend harus mengacu pada dokumen tersebut.

Jika API berubah, dokumentasi harus diperbarui pada commit yang sama.

---

# 10. Database Structure

Folder `database/` berisi artefak yang berkaitan dengan penyimpanan data.

```text
database/

migrations/
seeds/
schema/
views/
functions/
```

Model domain **bukan** berasal dari skema database. Skema database diturunkan dari model domain yang telah didefinisikan.

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
