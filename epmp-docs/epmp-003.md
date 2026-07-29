# EPMP-003 — Architecture Overview

| Metadata      | Value                             |
| ------------- | --------------------------------- |
| Document ID   | EPMP-003                          |
| Version       | 2.0.0                             |
| Status        | Approved                          |
| Owner         | Software Architecture Team        |
| Depends On    | EPMP-001, EPMP-002                |
| Referenced By | Seluruh Engineering Documentation |

---

# 1. Purpose

Dokumen ini mendefinisikan arsitektur tingkat tinggi (**High-Level Architecture**) Enterprise Property Management Platform (EPMP).

Dokumen ini menjelaskan bagaimana seluruh komponen sistem saling berhubungan tanpa membahas detail implementasi teknis. Tujuan utama dokumen ini adalah memberikan pemahaman bersama mengenai arah arsitektur platform sehingga seluruh keputusan desain, implementasi, dan pengembangan berikutnya tetap konsisten.

Dokumen ini menjadi jembatan antara Product Architecture dan Engineering Architecture.

---

# 2. Architecture Vision

EPMP dibangun sebagai **Enterprise Property Management Platform** yang:

- Modular
- Configurable
- API First
- Domain Driven
- Cloud Ready
- AI Friendly
- Enterprise Grade

Seluruh arsitektur dirancang agar mampu berkembang dari aplikasi web sederhana menjadi platform enterprise tanpa perubahan fundamental pada model domain.

---

# 3. Architecture Principles

Seluruh keputusan arsitektur harus mengikuti prinsip berikut.

## 3.1 Domain First

Business Domain merupakan inti sistem.

Database, framework, UI, dan library merupakan detail implementasi.

---

## 3.2 Configuration Over Customization

Perubahan perilaku bisnis sebisa mungkin dilakukan melalui konfigurasi, bukan perubahan kode.

---

## 3.3 API First

Semua fitur harus dapat diakses melalui API.

Frontend hanyalah salah satu konsumen API.

---

## 3.4 Modular by Design

Setiap bounded context merupakan modul independen.

---

## 3.5 Replaceable Infrastructure

Database, cache, storage, maupun external service harus dapat diganti tanpa mengubah Domain Layer.

---

## 3.6 AI-Friendly Architecture

Struktur sistem harus mudah dipahami oleh AI maupun developer.

Dokumentasi menjadi sumber kebenaran utama.

---

# 4. System Context

Versi pertama EPMP difokuskan pada aplikasi web.

```
┌────────────────────┐
│    Web Browser     │
└─────────┬──────────┘
          │ HTTPS
          ▼
┌────────────────────┐
│ React Web Frontend │
└─────────┬──────────┘
          │ REST API
          ▼
┌────────────────────┐
│   Golang Backend   │
└─────────┬──────────┘
          │
          ▼
┌────────────────────┐
│    PostgreSQL      │
└────────────────────┘
```

Pada fase berikutnya arsitektur dapat berkembang dengan:

- Mobile Application
- Public API
- External Integration
- Message Broker
- Cache Layer
- Object Storage

tanpa mengubah Domain Model.

---

# 5. Architecture Style

EPMP menggunakan kombinasi beberapa pendekatan arsitektur.

## Clean Architecture

Memisahkan:

- Interface
- Application
- Domain
- Infrastructure

Business Rule tidak boleh bergantung pada framework.

---

## Domain Driven Design

Setiap module dibangun berdasarkan bounded context.

Contoh:

- Organization
- Property
- Reservation
- Contract
- Finance
- Maintenance
- Asset

Bukan berdasarkan tabel database.

---

## Modular Monolith

Versi pertama menggunakan Modular Monolith.

Keuntungan:

- deployment sederhana
- debugging mudah
- transaction sederhana
- maintainability tinggi

Setiap modul tetap memiliki batas yang jelas sehingga dapat dipisahkan menjadi microservice bila diperlukan.

---

## Event Ready

Semua perubahan penting menghasilkan Domain Event.

Pada MVP event diproses secara in-process.

Arsitektur tetap memungkinkan penggunaan:

- Kafka
- RabbitMQ
- NATS

di masa depan.

---

# 6. Logical Layer

Backend terdiri dari empat lapisan.

```
Presentation Layer

↓

Application Layer

↓

Domain Layer

↓

Infrastructure Layer
```

## Presentation Layer

Tanggung jawab:

- HTTP
- Authentication
- Request
- Response
- Middleware

Tidak boleh memiliki business rule.

---

## Application Layer

Mengelola:

- Use Case
- Transaction
- Authorization
- Orchestration

Application Layer mengatur proses bisnis tetapi tidak menyimpan business rule inti.

---

## Domain Layer

Merupakan inti sistem.

Berisi:

- Entity
- Aggregate
- Value Object
- Domain Service
- Repository Interface
- Domain Event

Tidak mengetahui:

- SQL
- HTTP
- JSON
- Framework
- Logger

---

## Infrastructure Layer

Berisi implementasi teknis.

Contoh:

- PostgreSQL
- Redis
- SMTP
- Storage
- Queue
- External API

---

# 7. Core Domains

Core Domain EPMP terdiri dari:

- Organization
- Identity & Access Management
- Property
- Reservation
- Contract
- Occupancy
- Tenant
- Finance
- Maintenance
- Asset
- Configuration
- Reporting

Masing-masing domain memiliki bounded context yang independen.

---

# 8. Cross-Cutting Components

Beberapa komponen digunakan lintas domain.

## Authentication

Mengidentifikasi pengguna.

---

## Authorization

Mengatur Role dan Permission.

---

## Audit Trail

Mencatat seluruh aktivitas penting.

---

## Notification

Mengirim Email, WhatsApp, Push Notification, atau media lain.

---

## Configuration

Menyimpan Business Configuration.

---

## Logging

Structured Logging.

---

## Event Bus

Media komunikasi antar bounded context.

---

# 9. Technology Strategy

## Initial Platform

Versi pertama hanya mendukung:

- Web Admin
- Web Dashboard

Belum mencakup:

- Android
- iOS
- Desktop

Namun seluruh API harus dirancang agar platform-platform tersebut dapat ditambahkan tanpa perubahan besar.

---

## Backend

Bahasa:

- Go (Golang)

Pendekatan:

- Clean Architecture
- DDD
- Modular Monolith

---

## Frontend

Framework:

- React

Karakteristik:

- SPA
- Component Based
- Feature Based

---

## API

REST API menjadi standar utama.

GraphQL atau gRPC dapat dipertimbangkan apabila terdapat kebutuhan bisnis yang jelas.

---

## Database

Versi pertama menggunakan relational database.

Pemilihan vendor dibahas pada dokumen Engineering Platform Standard.

---

# 10. Scalability Strategy

EPMP dirancang untuk berkembang melalui tahapan berikut.

## Phase 1

Single Instance

```
Browser

↓

React

↓

Go

↓

PostgreSQL
```

---

## Phase 2

Horizontal Scaling

```
Browser

↓

Load Balancer

↓

Backend x N

↓

PostgreSQL
```

---

## Phase 3

Event Driven

```
Backend

↓

Event Bus

↓

Background Worker
```

---

## Phase 4

Selective Microservices

Hanya bounded context yang membutuhkan skalabilitas tinggi yang dipisahkan menjadi service tersendiri.

---

# 11. Architectural Decision

EPMP secara resmi menetapkan keputusan berikut.

| Area          | Decision                  |
| ------------- | ------------------------- |
| Architecture  | Modular Monolith          |
| API           | REST                      |
| Frontend      | React                     |
| Backend       | Golang                    |
| Design        | Domain Driven Design      |
| Pattern       | Clean Architecture        |
| Communication | Interface + Domain Event  |
| Documentation | Markdown                  |
| Development   | Human-Led, AI-Accelerated |

Perubahan terhadap keputusan ini harus melalui Engineering Decision Record (EDR).

---

# 12. Non Goals

Versi pertama tidak menargetkan:

- Full Microservices
- CQRS
- Event Sourcing
- Multi Region Deployment
- Multi Database
- Offline First

Arsitektur tetap membuka kemungkinan evolusi menuju kebutuhan tersebut apabila dibutuhkan.

---

# 13. Future Evolution

Roadmap evolusi arsitektur:

```
Single Application

↓

Modular Monolith

↓

Event Driven

↓

Service Extraction

↓

Distributed Platform
```

Evolusi dilakukan berdasarkan kebutuhan bisnis, bukan tren teknologi.

---

# Closing Statement

Architecture Overview menjadi fondasi teknis seluruh Enterprise Property Management Platform.

Semua keputusan implementasi harus tetap konsisten terhadap prinsip-prinsip yang dijelaskan dalam dokumen ini. Framework, library, maupun teknologi dapat berubah seiring waktu, tetapi Domain Model, Architecture Principles, dan Modular Design harus tetap menjadi inti dari platform.

Dokumen ini harus dibaca sebelum mempelajari dokumen engineering maupun implementation lainnya.
