# EPMP-008

# Solution Architecture (Go + React)

```text
Document ID    : EPMP-008
Document Name  : Solution Architecture
Version        : 1.0.0
Status         : Draft
Owner          : Software Architecture Team
Dependencies   : EPMP-001 ~ EPMP-007
Referenced By  : Repository Structure, Backend Architecture, Frontend Architecture, API, Database
```

---

# 1. Purpose

Dokumen ini mendefinisikan **arsitektur implementasi** Enterprise Property Management Platform (EPMP).

Berbeda dengan EPMP-003 yang menjelaskan arsitektur secara konseptual, dokumen ini menjelaskan bagaimana platform akan dibangun menggunakan teknologi yang telah dipilih.

Dokumen ini menjadi jembatan antara desain bisnis dan implementasi perangkat lunak.

---

# 2. Technology Stack

## Backend

- Go (Golang)

Alasan pemilihan:

- Performa tinggi.
- Binary deployment yang sederhana.
- Concurrency yang sangat baik.
- Strong typing.
- Mudah dipelihara dalam jangka panjang.
- Sangat cocok untuk REST API dan service modular.

---

## Frontend

- React

Alasan pemilihan:

- Component-based.
- Ekosistem yang matang.
- Mudah membangun dashboard kompleks.
- Reusable UI.
- Cocok untuk SPA (Single Page Application).

---

## Platform

Versi pertama EPMP hanya mendukung:

- Web Admin Portal
- Web Management Portal

Versi mobile akan menjadi proyek terpisah di masa depan dan menggunakan API yang sama.

---

# 3. Architectural Style

EPMP mengadopsi kombinasi beberapa pendekatan arsitektur.

## Clean Architecture

Memisahkan:

- Business Rules
- Application Logic
- Infrastructure
- Framework

Framework tidak boleh memengaruhi Domain.

---

## Domain Driven Design

Setiap module dibangun berdasarkan Domain.

Contoh:

```
Property

Contract

Reservation

Finance

Asset

Maintenance
```

Bukan berdasarkan tabel database.

---

## Modular Monolith

Seluruh sistem berjalan sebagai satu aplikasi.

Namun secara internal terdiri dari module yang independen.

```
EPMP

├── Property Module

├── Reservation Module

├── Contract Module

├── Finance Module

├── Asset Module

└── Maintenance Module
```

Module hanya berkomunikasi melalui interface dan event.

---

## Event Driven Ready

Walaupun MVP menggunakan pemanggilan langsung antar module melalui interface aplikasi, setiap perubahan penting harus menghasilkan **Domain Event**.

Contoh:

```
ReservationConfirmed

↓

ContractCreated

↓

InvoiceIssued
```

Dengan pendekatan ini, ketika nanti dibutuhkan message broker (Kafka, NATS, RabbitMQ, dsb.), implementasi dapat ditambahkan tanpa mengubah model domain.

---

# 4. System Context

```
Browser
     │
     ▼
React Web Application
     │
 REST API (HTTPS)
     │
     ▼
Golang Backend
     │
     ▼
Database
```

Komunikasi dilakukan menggunakan REST API berbasis JSON pada versi pertama.

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
