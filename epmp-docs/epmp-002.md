# EPMP-002

# Core Design Principles

```text
Document ID    : EPMP-002
Document Name  : Core Design Principles
Version        : 1.0.0
Status         : Draft
Owner          : Product & Architecture
Dependencies   : EPMP-001 Project Overview
Referenced By  : Seluruh dokumentasi EPMP
```

---

# 1. Purpose

Dokumen ini mendefinisikan prinsip desain yang menjadi fondasi seluruh Enterprise Property Management Platform (EPMP).

Semua keputusan terkait:

- Business Rule
- Database
- UI
- API
- Backend
- Frontend
- Automation
- AI Development

harus mengikuti prinsip-prinsip yang dijelaskan dalam dokumen ini.

Apabila terdapat konflik antara implementasi dan dokumen ini, maka dokumen ini menjadi acuan utama.

---

# 2. Core Philosophy

## Everything Should Be Configurable

Ini adalah filosofi utama EPMP.

Sistem tidak boleh mengandung business rule yang bersifat tetap apabila aturan tersebut dapat direpresentasikan sebagai konfigurasi.

Contoh yang **tidak diperbolehkan**:

- Deposit selalu Rp1.000.000
- Booking Fee selalu wajib
- Hanya ada 3 termin pembayaran
- Hanya ada AC dan Non-AC
- Hanya ada satu gedung

Sebaliknya, semua aturan tersebut harus dapat dikonfigurasi oleh pengguna yang memiliki hak akses.

---

# 3. Configuration Over Customization

EPMP lebih mengutamakan konfigurasi dibandingkan modifikasi kode.

Jika suatu kebutuhan dapat diselesaikan melalui konfigurasi, maka konfigurasi harus menjadi pilihan utama.

Dengan pendekatan ini, satu platform dapat melayani berbagai model bisnis tanpa percabangan kode yang kompleks.

---

# 4. Domain First

Arsitektur sistem harus mengikuti domain bisnis.

Urutan desain yang benar adalah:

```text
Business Domain
    ↓
Business Rules
    ↓
Domain Model
    ↓
Application Services
    ↓
API
    ↓
UI
```

UI tidak boleh menjadi dasar desain domain.

---

# 5. API First

Semua kemampuan sistem harus tersedia melalui API.

UI hanyalah salah satu client dari API tersebut.

Dengan pendekatan ini, EPMP akan lebih mudah dikembangkan menjadi:

- Web Application
- Mobile Application
- Public API
- Integrasi pihak ketiga
- Automation Service

---

# 6. Modular Architecture

Setiap modul harus memiliki satu tanggung jawab utama.

Contoh:

- Property Module
- Tenant Module
- Asset Module
- Contract Module

Setiap modul harus dapat dikembangkan secara independen selama kontraknya tidak berubah.

---

# 7. Single Source of Truth

Satu informasi hanya boleh memiliki satu sumber utama.

Contoh:

Status kamar berasal dari Room Module.

Status kontrak berasal dari Contract Module.

Status pembayaran berasal dari Finance Module.

Modul lain hanya mengonsumsi informasi tersebut.

---

# 8. Event Driven Mindset

Peristiwa penting dalam sistem harus menghasilkan Business Event.

Contoh:

```
ReservationCreated
ReservationExpired

TenantCheckedIn
TenantCheckedOut

ContractActivated
ContractExpired

InvoiceIssued
InvoicePaid

DepositReturned

AssetAssigned

MaintenanceRequested
MaintenanceCompleted
```

Business Event akan menjadi dasar automation di masa depan.

---

# 9. Audit Everything

Seluruh perubahan penting harus memiliki histori.

Audit minimal mencatat:

- waktu
- pengguna
- aksi
- data sebelum
- data sesudah

Tidak diperbolehkan melakukan perubahan data tanpa jejak.

---

# 10. Open for Extension

Sistem harus mudah ditambah.

Namun tidak mudah diubah.

Prinsip ini memungkinkan modul baru ditambahkan tanpa merusak modul yang telah stabil.

---

# 11. AI-First Documentation

Dokumentasi harus mudah dipahami oleh:

- manusia
- AI Assistant
- AI Coding Agent

Karena itu:

- istilah harus konsisten
- satu dokumen satu topik
- tidak ada informasi yang bertentangan

---

# 12. Security by Design

Keamanan bukan fitur tambahan.

Keamanan merupakan bagian dari desain.

Setiap modul wajib mempertimbangkan:

- Authentication
- Authorization
- Validation
- Audit
- Encryption

---

# 13. Scalability by Default

Setiap keputusan arsitektur harus mengasumsikan bahwa sistem akan berkembang.

Walaupun MVP hanya memiliki:

```
1 Property
1 Building
38 Rooms
```

Desain harus mampu berkembang menjadi:

```
100 Organizations

1000 Properties

5000 Buildings

100.000 Rooms
```

tanpa perubahan arsitektur.

---

# 14. Business Rules Must Not Live in UI

Business Rule hanya boleh berada pada Domain Layer atau Application Layer.

UI hanya bertugas menampilkan data.

---

# 15. Everything Has Lifecycle

Setiap entity memiliki lifecycle.

Misalnya:

Room

```
Available

Reserved

Occupied

Maintenance

Cleaning

Available
```

Contract

```
Draft

Pending

Active

Expired

Renewed

Completed
```

Asset

```
Available

Assigned

Maintenance

Disposed
```

Lifecycle akan dijelaskan lebih detail pada masing-masing Domain Specification.

---

# 16. Prefer Composition Over Inheritance

Entity tidak boleh saling bergantung melalui inheritance yang kompleks.

Hubungan antar entity lebih diutamakan menggunakan composition.

Pendekatan ini membuat domain lebih fleksibel.

---

# 17. Consistency Over Convenience

Konsistensi lebih penting dibanding kemudahan implementasi.

Jika terdapat dua pendekatan:

- satu cepat tetapi tidak konsisten
- satu sedikit lebih kompleks tetapi konsisten

maka pendekatan kedua dipilih.

---

# 18. Guiding Question

Sebelum menambahkan fitur baru, selalu jawab pertanyaan berikut:

> Apakah fitur ini dapat dibuat sebagai konfigurasi?

Jika jawabannya "Ya",

maka implementasi tidak boleh di-hardcode.

---
