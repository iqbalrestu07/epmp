# EPMP-010 — Engineering Platform Standard

| Metadata      | Value                             |
| ------------- | --------------------------------- |
| Document ID   | EPMP-010                          |
| Version       | 2.0.0                             |
| Status        | Approved                          |
| Owner         | Software Architecture Team        |
| Depends On    | EPMP-001 ~ EPMP-009               |
| Referenced By | Seluruh Engineering Documentation |

---

# 1. Purpose

Dokumen ini menetapkan standar engineering resmi yang digunakan dalam pengembangan Enterprise Property Management Platform (EPMP).

Dokumen ini merupakan **single source of truth** untuk seluruh keputusan teknologi, toolchain, library, workflow, dan engineering practices yang digunakan selama siklus hidup proyek.

Seluruh implementasi backend, frontend, database, testing, deployment, dan AI-assisted development harus mengacu pada dokumen ini.

---

# 2. Scope

Dokumen ini mencakup:

- Technology Baseline
- Backend Platform
- Frontend Platform
- Database Platform
- Development Workflow
- Code Quality
- Security Baseline
- Observability
- Tooling
- AI Development Support

Dokumen ini **tidak** membahas implementasi detail modul tertentu.

---

# 3. Engineering Principles

Seluruh keputusan engineering wajib mengikuti prinsip berikut.

## EP-01 — Simplicity First

Engineer **MUST** memilih solusi paling sederhana yang mampu memenuhi kebutuhan bisnis.

### Rationale

Kompleksitas merupakan biaya jangka panjang.

---

## EP-02 — Long-Term Maintainability

Seluruh keputusan **SHOULD** mempertimbangkan maintainability dibanding optimasi prematur.

---

## EP-03 — Convention Over Configuration

Apabila terdapat standar internal, implementasi **MUST** mengikuti standar tersebut.

---

## EP-04 — Open Source Preferred

Library open source yang matang **SHOULD** diprioritaskan dibanding solusi proprietary apabila memenuhi kebutuhan proyek.

---

## EP-05 — Explicit Over Magic

Framework atau library yang menyembunyikan terlalu banyak proses internal **SHOULD NOT** digunakan apabila mengurangi keterbacaan kode.

---

## EP-06 — AI-Friendly Development

Source code, struktur folder, dan dokumentasi **MUST** mudah dipahami oleh manusia maupun AI Coding Agent.

---

# 4. Technology Baseline

## Backend

| Area            | Standard                |
| --------------- | ----------------------- |
| Language        | Go                      |
| HTTP            | net/http                |
| Router          | echo                    |
| Architecture    | Clean Architecture      |
| Design          | Domain Driven Design    |
| Database Access | pgx + sqlc              |
| Logging         | zerolog                 |
| Validation      | go-playground/validator |
| Migration       | golang-migrate          |
| Configuration   | koanf                   |
| Testing         | testing + testify       |

---

## Frontend

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

---

## Database

Primary Database

- PostgreSQL

Optional Components

- Redis
- Object Storage (future)

---

# 5. Official Development Toolkit

EPMP menggunakan toolkit resmi untuk mempercepat implementasi.

```
helper-package/

├── be/

└── fe/
```

Toolkit ini merupakan bagian resmi dari platform.

CRUD boilerplate, template, dan generator **MUST** berasal dari toolkit ini apabila tersedia.

Developer maupun AI **MUST NOT** membuat ulang boilerplate yang telah tersedia.

---

## Backend Toolkit

Folder:

```
helper-package/be/
```

Digunakan untuk menghasilkan:

- CRUD
- Repository
- DTO
- Handler
- Migration
- Validation
- Test Skeleton

---

## Frontend Toolkit

Folder:

```
helper-package/fe/
```

Digunakan untuk menghasilkan:

- CRUD Page
- Form
- Table
- API Client
- Hooks
- Types
- Layout

---

## Rationale

Generator resmi memberikan:

- konsistensi struktur
- pengurangan boilerplate
- efisiensi token AI
- kemudahan maintenance

---

# 6. Development Model

EPMP menggunakan model:

```
Human-Led

↓

AI-Accelerated
```

Human bertanggung jawab terhadap:

- Requirement
- Architecture
- Design
- Review
- Release

AI membantu:

- Implementation
- Boilerplate
- Test
- Documentation
- Refactoring

---

# 7. Coding Standard

Seluruh source code **MUST** mengikuti dokumen:

- Backend Architecture Standard
- Go Coding Standard
- Frontend Standard

Apabila terjadi konflik:

Engineering Standard memiliki prioritas lebih tinggi.

---

# 8. Documentation Standard

Seluruh dokumentasi menggunakan Markdown.

Setiap module **MUST** memiliki README.

Setiap perubahan perilaku sistem **MUST** disertai pembaruan dokumentasi.

---

# 9. Repository Standard

Repository menggunakan Monorepo.

```
docs/

backend/

frontend/

database/

helper-package/

scripts/

docker/
```

Tidak diperbolehkan membuat struktur folder baru tanpa alasan yang terdokumentasi.

---

# 10. Branch Strategy

Official Branch

```
main
```

Development Branch

```
develop
```

Feature

```
feature/{module}/{feature}
```

Bugfix

```
fix/{issue}
```

Release

```
release/{version}
```

Hotfix

```
hotfix/{version}
```

---

# 11. Versioning

EPMP menggunakan Semantic Versioning.

```
MAJOR.MINOR.PATCH
```

Contoh:

```
1.4.2
```

---

# 12. Commit Convention

Format:

```
type(scope): description
```

Contoh:

```
feat(property): add room availability service

fix(contract): prevent duplicate activation

docs(epmp): update architecture overview

refactor(finance): simplify invoice calculation
```

Jenis commit:

- feat
- fix
- docs
- test
- refactor
- perf
- chore
- ci
- build

---

# 13. Code Quality Requirements

Seluruh Pull Request wajib memenuhi:

- Build berhasil
- Lint berhasil
- Unit Test berhasil
- Integration Test berhasil (jika relevan)
- Tidak menurunkan code coverage
- Dokumentasi diperbarui

Reviewer **MUST** menolak Pull Request yang tidak memenuhi persyaratan tersebut.

---

# 14. Security Baseline

Seluruh endpoint wajib:

- tervalidasi
- diautentikasi
- diautorisasi
- memiliki audit trail

Rahasia aplikasi **MUST NOT** disimpan dalam repository.

Password **MUST** di-hash menggunakan algoritma yang sesuai standar keamanan proyek.

---

# 15. Observability

Platform **MUST** menyediakan:

- Structured Logging
- Request ID
- Correlation ID
- Health Check
- Metrics

Platform **SHOULD** mendukung Distributed Tracing pada fase berikutnya.

---

# 16. Configuration Policy

Konfigurasi dibagi menjadi dua kategori.

## Runtime Configuration

Contoh:

- Database
- SMTP
- JWT
- Storage

Disimpan di environment.

---

## Business Configuration

Contoh:

- Deposit Policy
- Payment Term
- Reservation Expiry

Disimpan di database.

Business Configuration **MUST NOT** di-hardcode.

---

# 17. AI Development Support

Seluruh AI Coding Agent **MUST** mengikuti:

- EPMP-012 AI Development Protocol
- helper-package
- Prompt Library
- Work Package

AI **MUST NOT** membuat struktur proyek sendiri.

AI **MUST** menggunakan generator resmi apabila tersedia.

---

# 18. Decision Records

Perubahan terhadap teknologi berikut **WAJIB** melalui Engineering Decision Record (EDR):

- Bahasa pemrograman
- Framework utama
- Database
- Toolkit resmi
- Architecture Style
- Repository Strategy

---

# 19. Future Consideration

Beberapa teknologi dapat dipertimbangkan di masa depan apabila terdapat kebutuhan nyata.

Contoh:

- Redis Cluster
- Kafka
- OpenTelemetry
- GraphQL
- gRPC
- Object Storage
- Kubernetes

Dokumen ini tidak mengharuskan implementasi teknologi tersebut pada MVP.

---

# Closing Statement

Engineering Platform Standard merupakan kontrak teknis resmi EPMP.

Seluruh engineer, AI Coding Agent, reviewer, dan maintainer wajib mengacu pada dokumen ini dalam setiap aktivitas pengembangan.

Konsistensi engineering lebih penting daripada preferensi pribadi terhadap framework, library, atau gaya implementasi.
