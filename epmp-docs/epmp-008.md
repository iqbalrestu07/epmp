# EPMP-008

# Solution Architecture (Go + React)

```text
Document ID    : EPMP-008
Document Name  : Solution Architecture
Version        : 1.1.0
Status         : Implemented & Active
Owner          : Software Architecture Team
Dependencies   : EPMP-001 ~ EPMP-007
Referenced By  : Repository Structure, Backend Architecture, Frontend Architecture, API, Database
```

---

# 1. Purpose

Dokumen ini mendefinisikan **arsitektur implementasi nyata** Enterprise Property Management Platform (EPMP).

Berbeda dengan EPMP-003 yang menjelaskan arsitektur secara konseptual, dokumen ini menjelaskan bagaimana platform dibangun dan dioperasikan menggunakan stack teknologi aktif: Go (Echo), PostgreSQL, React (Vite, Tailwind), Three.js (3D WebGL Canvas), dan native Go WhatsApp Gateway (`whatsmeow`).

Dokumen ini menjadi jembatan antara desain bisnis dan implementasi perangkat lunak.

---

# 2. Technology Stack

## Backend

- **Language & Runtime**: Go (Golang 1.24+)
- **HTTP Framework**: Echo v4 (REST API routing, middleware, CORS, request validation)
- **Database Engine**: PostgreSQL 16 (Relational DB, ACID transactions, database triggers)
- **Database Migrations**: `golang-migrate/migrate/v4`
- **WhatsApp Gateway**: `go.mau.fi/whatsmeow` (Multi-device WhatsApp protocol, session management via SQLite3/PostgreSQL container store)
- **Authentication**: JWT (JSON Web Tokens) dengan bcrypt password hashing dan RBAC middleware.

Alasan pemilihan:
- Performa tinggi dengan binary footprint kecil.
- Concurrency goroutine yang sangat tangguh untuk menangani koneksi websocket WhatsApp dan broadcast massal.
- Strong typing mencegah runtime bugs pada model finansial dan status ketersediaan.

---

## Frontend

- **Core Library**: React 19 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS + Shadcn UI primitives
- **Spatial 3D Engine**: Three.js (`@react-three/fiber` & `@react-three/drei`) untuk visualisasi 3D WebGL bangunan dan status kamar real-time.
- **Routing**: React Router v7 (Feature-driven routing)
- **Icons & QR**: `lucide-react`, `react-qr-code` (render QR pairing resmi WhatsApp)
- **State & Data Fetching**: TanStack React Query / Custom Hook Service Architecture

Alasan pemilihan:
- SPA responsif dengan pemisahan bounded context berbasis fitur (`src/features/*`).
- Visualisasi 3D langsung di browser tanpa plugin eksternal.
- Dynamic theme dan multi-currency formatting yang fleksibel (IDR, USD, EUR, SGD, MYR).

---

## Testing & Quality Assurance

- **End-to-End (E2E)**: Playwright (`make test-e2e` mencakup 40+ rute dashboard dan visual 3D canvas).
- **Backend Unit/Integration Tests**: Go standard test suite (`go test ./...`).

---

# 3. Architectural Style

EPMP mengadopsi kombinasi beberapa pendekatan arsitektur:

## Clean Architecture
Memisahkan:
- Business Rules (Domain Entities & Value Objects)
- Application Logic (Use Cases & Handlers)
- Infrastructure (Repositories, whatsmeow client, DB drivers)
- Framework (Echo HTTP router)

## Domain Driven Design (DDD)
Setiap modul backend dan fitur frontend dibangun berdasarkan Bounded Context domain:
```text
backend/internal/modules/
├── property/          (Property, Building, Floor, Room, Bed, 3D spatial layout)
├── tenant/            (Tenant directory & profiles)
├── reservation/       (Booking, holding, expiration)
├── contract/          (Lease agreement, renewal, terms)
├── occupancy/         (Check-in, check-out, living records)
├── billing/           (Invoice, recurring fees, multi-currency)
├── payment/           (Payment collection, auto-settlement trigger)
├── deposit/           (Security deposit, refund, deductions)
├── communication/     (WhatsApp devices, pairing, blast campaigns, audit logs)
├── asset/             (Inventory items, assignments, inspections)
└── maintenance/       (Work orders, technician dispatch, maintenance locks)
```

## Modular Monolith
Seluruh sistem berjalan sebagai satu service backend Go terpadu dengan modul internal independen yang berkomunikasi melalui Go interface dan database event triggers.

## Event Driven & Triggers
- **Database Trigger (`migration 000038`)**: Otomatisasi sinkronisasi status Invoice menjadi `Paid` ketika akumulasi pembayaran berstatus `completed` memenuhi nilai tagihan.
- **Whatsmeow Event Handlers**: Menangkap event live login, session disconnect, dan delivery confirmation dari WhatsApp Web protocol.

---

# 4. System Context

```text
       ┌─────────────────────────────────────────────────────────┐
       │                   Web Browser (Client)                 │
       │  React 19 SPA + Tailwind CSS + Three.js 3D WebGL Canvas │
       └────────────────────────────┬────────────────────────────┘
                                    │ HTTPS (REST API + JSON)
                                    ▼
       ┌─────────────────────────────────────────────────────────┐
       │                    Golang Backend (Echo)                │
       │                                                         │
       │  ┌─────────────────┐             ┌───────────────────┐  │
       │  │ Domain Handlers │             │ Communication Svc │  │
       │  │ & Use Cases     │             │ (whatsmeow engine)│  │
       │  └────────┬────────┘             └─────────┬─────────┘  │
       └───────────┼────────────────────────────────┼────────────┘
                   │                                │
                   ▼                                ▼
       ┌───────────────────────┐        ┌────────────────────────┐
       │     PostgreSQL 16     │        │  WhatsApp Web Servers  │
       │  - Core App Database  │        │  (Direct TLS Pairing)  │
       │  - Auto-Paid Triggers │        └────────────────────────┘
       └───────────────────────┘
```

---

# 5. Backend Architecture

Backend dibagi menjadi beberapa lapisan.

```
HTTP Layer
        │
        ▼
Application Layer
        │
        ▼
Domain Layer
        │
        ▼
Infrastructure Layer
```

## HTTP Layer

Tanggung jawab:

- Routing
- Request Parsing
- Response Formatting
- Middleware

Tidak boleh berisi business rule.

---

## Application Layer

Tanggung jawab:

- Use Case
- Orchestration
- Transaction Boundary
- Authorization
- Validation

Application Layer mengoordinasikan domain tanpa menyimpan aturan bisnis inti.

---

## Domain Layer

Lapisan paling penting.

Berisi:

- Entity
- Value Object
- Aggregate
- Domain Service
- Repository Interface
- Domain Event
- Business Rules

Lapisan ini tidak mengetahui framework, database, atau HTTP.

---

## Infrastructure Layer

Implementasi teknis.

Contoh:

- PostgreSQL Repository
- File Storage
- Email
- Cache
- Logger
- Queue
- External API

Infrastructure hanya memenuhi kontrak yang didefinisikan oleh Domain atau Application Layer.

---

# 6. Frontend Architecture

Frontend menggunakan pendekatan berbasis fitur (feature-based architecture), bukan berdasarkan jenis file.

Contoh:

```
features/
├── property/
├── reservation/
├── contract/
├── finance/
├── maintenance/
```

Setiap fitur memiliki:

- halaman
- komponen
- hooks
- service API
- validasi
- state lokal

Pendekatan ini menjaga konsistensi dengan bounded context di backend.

---

# 7. API Strategy

REST API menjadi standar utama.

Prinsip:

- Resource-oriented.
- Stateless.
- JSON.
- Versioning (`/api/v1/...`).
- Idempotent untuk operasi yang sesuai.
- Konsisten dalam format request dan response.

Semua endpoint akan dirancang berdasarkan domain, bukan tabel database.

---

# 8. Authentication & Authorization

Versi pertama mendukung:

- Login berbasis email/username.
- Password yang di-hash.
- JWT untuk autentikasi.
- Role-Based Access Control (RBAC).

Authorization dilakukan di Application Layer berdasarkan Role dan Permission.

---

# 9. Data Persistence

Versi pertama menggunakan **satu database relasional**.

Alasan:

- Konsistensi transaksi.
- Kemudahan deployment.
- Kompleksitas operasional lebih rendah.
- Sesuai dengan pendekatan Modular Monolith.

Pemilihan vendor database akan dibahas pada dokumen Database Architecture. (Secara pribadi saya cenderung memilih PostgreSQL karena fitur dan stabilitasnya, tetapi keputusan final akan kita dokumentasikan di dokumen database agar tetap terpisah dari desain domain.)

---

# 10. Configuration Management

Konfigurasi dibedakan menjadi dua kategori:

### System Configuration

Dikelola melalui environment atau deployment.

Contoh:

- Database Connection
- JWT Secret
- SMTP
- Storage
- Logging

### Business Configuration

Dikelola melalui aplikasi.

Contoh:

- Room Type
- Charge Type
- Deposit Policy
- Payment Term
- Reservation Expiry
- Late Fee Policy

Perubahan Business Configuration tidak memerlukan deployment ulang.

---

# 11. Error Handling Strategy

Semua error harus:

- Konsisten.
- Memiliki kode.
- Memiliki pesan yang dapat dipahami pengguna.
- Dapat dicatat untuk audit dan debugging.

Business error dipisahkan dari system error.

---

# 12. Logging & Observability

Sistem harus mendukung:

- Structured Logging.
- Audit Trail.
- Request ID.
- Correlation ID.
- Performance Metrics.
- Health Check.

Hal ini memudahkan troubleshooting dan pemantauan aplikasi di lingkungan produksi.

---

# 13. Deployment Strategy

Versi pertama di-deploy sebagai satu aplikasi backend dan satu aplikasi frontend.

```
React Web

↓

Nginx

↓

Go Backend

↓

PostgreSQL
```

Arsitektur ini sederhana namun tetap memberikan ruang untuk berkembang.

---

# 14. Scalability Strategy

EPMP dirancang agar dapat berkembang tanpa perubahan arsitektur besar.

Tahapan evolusi yang direncanakan:

1. **Single Instance** – MVP.
2. **Modular Monolith** – penambahan modul dan pengguna.
3. **Horizontal Scaling** – beberapa instance backend di belakang load balancer.
4. **Event Bus** – komunikasi asinkron untuk proses tertentu.
5. **Selective Microservices** – hanya domain yang benar-benar membutuhkan skalabilitas atau isolasi tinggi yang dipisahkan menjadi layanan mandiri.

Pendekatan ini menghindari kompleksitas dini sekaligus menjaga fleksibilitas jangka panjang.

---

# 15. AI-Assisted Development Strategy

Dokumentasi EPMP disusun agar dapat digunakan langsung oleh AI Coding Agent.

Prinsipnya:

- Satu bounded context = satu ruang kerja AI.
- Satu modul = satu spesifikasi utama.
- Semua istilah mengacu pada EPMP-004.
- Semua implementasi mengikuti EPMP-002.

Dengan demikian, AI dapat menghasilkan kode yang konsisten antar modul tanpa harus memahami keseluruhan sistem sekaligus.

---

# 16. Non-Goals

Versi pertama **tidak** menargetkan:

- Microservice penuh.
- Multi-database.
- Multi-region deployment.
- Event sourcing.
- CQRS penuh.
- Offline-first.

Arsitektur tetap memungkinkan evolusi ke arah tersebut bila ada kebutuhan bisnis yang nyata.

---
