# EPMP-011 — Backend Architecture Standard

| Metadata      | Value                                                                 |
| ------------- | --------------------------------------------------------------------- |
| Document ID   | EPMP-011                                                              |
| Version       | 2.0.0                                                                 |
| Status        | Approved                                                              |
| Owner         | Software Architecture Team                                            |
| Depends On    | EPMP-001 ~ EPMP-010                                                   |
| Referenced By | Backend Modules, API Standards, Database Standards, Testing Standards |

---

# 1. Purpose

Dokumen ini mendefinisikan standar arsitektur backend resmi untuk Enterprise Property Management Platform (EPMP).

Dokumen ini menjadi pedoman utama dalam membangun seluruh bounded context menggunakan Golang dengan pendekatan:

- Domain Driven Design (DDD)
- Clean Architecture
- Modular Monolith
- API First
- Human-Led, AI-Accelerated Development

Seluruh implementasi backend **MUST** mengikuti standar pada dokumen ini.

---

# 2. Architectural Vision

Backend EPMP dibangun dengan prinsip:

> **Business Domain is the Core. Everything Else is Infrastructure.**

Artinya:

- Framework dapat diganti.
- Database dapat diganti.
- Router dapat diganti.
- Logger dapat diganti.

Namun Domain Model harus tetap stabil.

---

# 3. Backend Architecture

```

Presentation Layer

↓

Application Layer

↓

Domain Layer

↓

Infrastructure Layer

```

Dependency Rule:

```

Presentation
↓

Application
↓

Domain

Infrastructure
↑

```

Domain tidak memiliki dependency ke layer lain.

---

# 4. Layer Responsibilities

## 4.1 Presentation Layer

Presentation Layer bertanggung jawab terhadap:

- HTTP Endpoint
- Request Validation
- Response Mapping
- Middleware
- Authentication
- Serialization

Layer ini **MUST NOT** mengandung business rule.

---

## 4.2 Application Layer

Application Layer bertanggung jawab terhadap:

- Use Case
- Command
- Query
- Transaction Boundary
- Authorization
- Orchestration
- Event Publishing

Layer ini mengatur jalannya proses bisnis.

Layer ini **MUST NOT** menjadi tempat penyimpanan business rule inti.

---

## 4.3 Domain Layer

Domain Layer merupakan inti sistem.

Layer ini berisi:

- Entity
- Aggregate Root
- Value Object
- Domain Service
- Repository Interface
- Domain Event
- Business Rule
- Domain Policy
- Specification (optional)

Domain Layer:

- MUST bebas framework
- MUST bebas SQL
- MUST bebas HTTP
- MUST bebas Logger
- MUST bebas JSON

---

## 4.4 Infrastructure Layer

Infrastructure mengimplementasikan kebutuhan teknis.

Contoh:

- PostgreSQL
- Redis
- SMTP
- Queue
- File Storage
- External API
- Scheduler

Infrastructure hanya memenuhi kontrak yang dibuat Domain atau Application.

---

# 5. Module Structure

Setiap bounded context menggunakan struktur yang identik.

```

property/

README.md

application/
domain/
infrastructure/
interfaces/

```

Tidak diperbolehkan membuat struktur berbeda tanpa Architecture Review.

---

# 6. Internal Package Structure

## application/

```

application/

command/
query/
dto/
mapper/
service/

```

---

## domain/

```

domain/

entity/
aggregate/
repository/
event/
service/
policy/
valueobject/
errors/

```

---

## infrastructure/

```

infrastructure/

postgres/
cache/
storage/
email/
external/

```

---

## interfaces/

```

interfaces/

http/
middleware/
request/
response/

```

---

# 7. Aggregate Design

Setiap Aggregate:

- memiliki satu Aggregate Root
- menjaga konsistensi data
- menjaga invariant
- menjadi batas transaksi bisnis

Aggregate **MUST NOT** bergantung pada Aggregate lain.

Komunikasi dilakukan melalui:

- Application Layer
- Domain Event

---

# 8. Repository Standard

Repository merupakan kontrak.

Repository berada pada Domain Layer.

Contoh:

```go
type PropertyRepository interface {
    FindByID(...)
    Save(...)
    Delete(...)
}
```

Implementasi berada pada Infrastructure.

Repository **MUST NOT** berisi business rule.

---

# 9. Use Case Standard

Setiap Use Case memiliki satu tanggung jawab.

Contoh:

```

CreateReservation

ActivateContract

CheckInTenant

GenerateInvoice

```

Tidak diperbolehkan membuat Use Case seperti:

```

ManageProperty

```

Karena terlalu luas.

---

# 10. Transaction Boundary

Transaction dikelola pada Application Layer.

Repository **MUST NOT** membuka transaction sendiri.

Semua perubahan Aggregate yang terkait dilakukan dalam satu transaction yang dikelola oleh Use Case.

---

# 11. Domain Events

Semua perubahan penting menghasilkan Domain Event.

Contoh:

- ReservationCreated
- ReservationConfirmed
- ContractActivated
- InvoiceIssued
- PaymentReceived

Domain Event merupakan mekanisme komunikasi resmi antar bounded context.

---

# 12. Dependency Injection

EPMP menggunakan Manual Dependency Injection.

Composition Root berada pada:

```

cmd/

```

Seluruh dependency dibangun di sana.

Framework Dependency Injection **MUST NOT** digunakan tanpa Architecture Decision Record.

---

# 13. Error Handling

Error dibagi menjadi tiga kategori.

## Domain Error

Contoh:

```

RoomAlreadyOccupied

```

---

## Application Error

Contoh:

```

PermissionDenied

```

---

## Infrastructure Error

Contoh:

```

DatabaseUnavailable

```

Mapping ke HTTP dilakukan pada Presentation Layer.

---

# 14. Cross Module Communication

Urutan komunikasi resmi:

```

Application Service

↓

Domain Event

```

Repository antar bounded context tidak boleh saling diakses.

---

# 15. Logging Policy

Logging hanya dilakukan pada:

- Presentation Layer
- Infrastructure Layer

Domain Layer **MUST NOT** melakukan logging.

---

# 16. Configuration Policy

Runtime Configuration

Contoh:

- Database
- SMTP
- JWT
- Storage

Business Configuration

Contoh:

- Deposit Policy
- Reservation Expiry
- Payment Term

Business Configuration **MUST** berada di database.

---

# 17. Generator Policy

Seluruh implementasi CRUD **MUST** menggunakan helper-package apabila generator tersedia.

```

helper-package/

├── be/

└── fe/

```

Developer maupun AI **MUST NOT** membuat ulang CRUD boilerplate.

Apabila generator belum mendukung kebutuhan tertentu, implementasi manual diperbolehkan dengan tetap mengikuti standar arsitektur.

---

# 18. AI Collaboration Rules

AI bertugas:

- Generate CRUD
- Generate DTO
- Generate Repository
- Generate Migration
- Generate Test

Human bertugas:

- Domain Design
- Architecture
- Review
- Business Rule
- Merge

---

# 19. Architecture Constraints

Backend EPMP:

MUST:

- menggunakan Clean Architecture
- menggunakan DDD
- menggunakan Repository Pattern
- menggunakan Manual DI
- menggunakan REST API

MUST NOT:

- SQL di Domain
- HTTP di Domain
- JSON di Domain
- Framework di Domain
- Business Rule di Controller
- Business Rule di Repository

---

# 20. Backend Review Checklist

Sebelum Pull Request diterima:

- Domain bebas framework.
- Business Rule berada di Domain.
- Repository hanya kontrak.
- Transaction di Application.
- Test tersedia.
- Dokumentasi diperbarui.
- Menggunakan helper-package bila tersedia.
- Tidak melanggar bounded context.

Semua checklist harus bernilai **PASS**.

---

# 21. Future Evolution

Arsitektur backend dirancang berkembang melalui tahapan berikut:

```

Modular Monolith

↓

Event Driven

↓

Shared Infrastructure

↓

Selective Microservices

```

Tidak diperbolehkan melakukan migrasi ke microservice hanya berdasarkan tren teknologi.

Migrasi harus memiliki alasan bisnis dan didokumentasikan melalui Engineering Decision Record (EDR).

---

# Closing Statement

Backend Architecture Standard merupakan konstitusi implementasi backend EPMP.

Seluruh engineer dan AI Coding Agent wajib mengikuti dokumen ini agar seluruh bounded context memiliki struktur, pola, dan kualitas yang konsisten sepanjang siklus hidup platform.
