T# EPMP-005

# Core Domain Model

```text
Document ID    : EPMP-005
Document Name  : Core Domain Model
Version        : 1.0.0
Status         : Draft
Owner          : Software Architecture Team
Dependencies   : EPMP-001, EPMP-002, EPMP-003, EPMP-004
Referenced By  : Seluruh Domain Specification, Module Specification, API, Database
```

---

# 1. Purpose

Dokumen ini mendefinisikan struktur domain utama Enterprise Property Management Platform (EPMP).

Dokumen ini **tidak membahas implementasi teknis**, melainkan menjelaskan bagaimana bisnis dipisahkan menjadi domain-domain yang independen namun saling berkolaborasi.

Core Domain Model menjadi fondasi bagi seluruh desain entity, database, API, dan modul pada EPMP.

---

# 2. Domain-Driven Design (DDD)

EPMP mengadopsi pendekatan **Domain-Driven Design (DDD)**.

Artinya, desain sistem dibangun berdasarkan model bisnis, bukan berdasarkan tampilan (UI), struktur database, ataupun framework.

Setiap domain memiliki:

- Tanggung jawab yang jelas.
- Batas (Boundary) yang tegas.
- Bahasa (Ubiquitous Language) yang konsisten.
- Aturan bisnis (Business Rules) yang spesifik.
- Siklus hidup (Lifecycle) masing-masing.

---

# 3. Domain Classification

Domain dalam EPMP dibagi menjadi tiga kategori.

## 3.1 Core Domain

Core Domain adalah domain yang memberikan nilai bisnis utama dan menjadi pembeda EPMP.

- Property Management
- Reservation Management
- Contract Management
- Occupancy Management
- Billing & Finance

Perubahan pada Core Domain harus melalui analisis arsitektur yang mendalam.

---

## 3.2 Supporting Domain

Supporting Domain mendukung operasional Core Domain.

- Tenant Management
- Asset Management
- Maintenance Management
- Notification
- Reporting

Supporting Domain dapat berkembang secara independen selama tidak melanggar kontrak dengan Core Domain.

---

## 3.3 Generic Domain

Generic Domain adalah kemampuan umum yang dapat digunakan lintas aplikasi.

- Authentication
- Authorization
- File Storage
- Audit Log
- Configuration
- Search
- Logging
- Scheduler
- Integration

Generic Domain tidak mengandung aturan bisnis khusus properti.

---

# 4. Domain Landscape

Diagram berikut menunjukkan hubungan antar domain utama.

```text
Organization
│
├── Property Management
│   ├── Property
│   ├── Building
│   ├── Floor
│   ├── Zone
│   ├── Room
│   └── Bed
│
├── Tenant Management
│   ├── Tenant
│   ├── Identity
│   ├── Contact
│   └── Documents
│
├── Reservation Management
│
├── Contract Management
│
├── Occupancy Management
│
├── Billing & Finance
│
├── Asset Management
│
├── Maintenance Management
│
├── Reporting
│
└── Notification
```

---

# 5. Domain Relationships

Hubungan antar domain harus mengikuti aturan berikut.

- Property Management **tidak mengetahui** detail pembayaran.
- Finance **tidak mengubah** status Room secara langsung.
- Reservation **tidak mengaktifkan** Contract secara langsung.
- Contract **tidak mengubah** Asset.
- Maintenance **tidak mengubah** Invoice.

Komunikasi dilakukan melalui Application Layer dan Domain Events.

---

# 6. Domain Ownership

Setiap domain memiliki kepemilikan yang jelas.

| Domain                 | Owner            |
| ---------------------- | ---------------- |
| Property Management    | Property Team    |
| Tenant Management      | Tenant Team      |
| Reservation Management | Reservation Team |
| Contract Management    | Contract Team    |
| Billing & Finance      | Finance Team     |
| Asset Management       | Asset Team       |
| Maintenance            | Maintenance Team |
| Reporting              | Analytics Team   |
| Configuration          | Platform Team    |

Ownership ini penting untuk menjaga batas tanggung jawab saat tim berkembang.

---

# 7. Core Domain Detail

## 7.1 Property Management

Bertanggung jawab atas struktur fisik properti.

Mencakup:

- Property
- Building
- Floor
- Zone
- Room
- Bed
- Facility
- Room Type

Tidak mengelola penghuni maupun pembayaran.

---

## 7.2 Tenant Management

Mengelola identitas dan informasi penyewa.

Mencakup:

- Biodata
- Dokumen
- Kontak
- Riwayat
- Emergency Contact

Tidak mengelola kontrak maupun tagihan.

---

## 7.3 Reservation Management

Mengelola proses pemesanan sebelum penyewa resmi menempati unit.

Mencakup:

- Reservation
- Booking Fee
- Expiry
- Approval

---

## 7.4 Contract Management

Mengelola hubungan hukum antara Tenant dan Organization.

Mencakup:

- Draft
- Active
- Renewal
- Extension
- Termination
- Renewal History

---

## 7.5 Occupancy Management

Mengelola kondisi penggunaan fisik Room atau Bed.

Status umum:

- Available
- Reserved
- Occupied
- Cleaning
- Maintenance
- Blocked

Occupancy adalah domain tersendiri karena status fisik tidak selalu sama dengan status kontrak.

---

## 7.6 Billing & Finance

Mengelola seluruh transaksi finansial.

Mencakup:

- Charge
- Invoice
- Payment
- Deposit
- Refund
- Adjustment
- Penalty

---

## 7.7 Asset Management

Mengelola aset yang dimiliki organisasi.

Mencakup:

- Asset Registry
- Assignment
- Inspection
- Disposal

---

## 7.8 Maintenance Management

Mengelola proses pemeliharaan.

Mencakup:

- Work Order
- Technician
- Vendor
- Schedule
- Inspection
- Completion

---

## 7.9 Reporting

Menyediakan data agregasi dan analitik.

Contoh:

- Occupancy Rate
- Revenue
- Outstanding Invoice
- Asset Condition
- Maintenance Summary

Reporting tidak mengubah data operasional.

---

# 8. Domain Dependency Rules

Agar sistem tetap modular, ketergantungan antar domain diatur sebagai berikut.

```text
Property
    ↓
Reservation
    ↓
Contract
    ↓
Occupancy
    ↓
Finance
```

Domain di sebelah kanan boleh menggunakan informasi dari domain di sebelah kiri melalui kontrak yang telah ditentukan, namun tidak boleh memodifikasi state internal domain lain secara langsung.

---

# 9. Cross-Cutting Domains

Beberapa domain mendukung seluruh sistem.

- Authentication
- Authorization
- Configuration
- Notification
- File Management
- Audit Log
- Scheduler
- Search
- Integration

Cross-Cutting Domain tidak menyimpan aturan bisnis inti.

---

# 10. Domain Events

Semua domain menghasilkan event yang dapat dikonsumsi domain lain.

Contoh:

| Domain      | Event                                                        |
| ----------- | ------------------------------------------------------------ |
| Reservation | ReservationCreated, ReservationConfirmed, ReservationExpired |
| Contract    | ContractActivated, ContractRenewed, ContractTerminated       |
| Occupancy   | TenantCheckedIn, TenantCheckedOut, RoomVacated               |
| Finance     | InvoiceIssued, PaymentReceived, DepositReturned              |
| Asset       | AssetAssigned, AssetReturned                                 |
| Maintenance | MaintenanceRequested, MaintenanceCompleted                   |

Event menjadi mekanisme komunikasi utama antar domain.

---

# 11. Future Domains

Arsitektur disiapkan agar dapat menambahkan domain baru tanpa mengubah Core Domain.

Contoh domain masa depan:

- Visitor Management
- Parking Management
- Utility Meter Management
- Housekeeping
- Dynamic Pricing
- CRM
- Loyalty Program
- AI Recommendation
- IoT Device Management
- Smart Lock Management

Domain baru harus mematuhi prinsip-prinsip pada EPMP-002 dan berintegrasi melalui kontrak yang jelas.

---

# 12. Domain Evolution Policy

Core Domain harus stabil. Perubahan besar hanya dilakukan jika benar-benar mengubah model bisnis.

Supporting Domain dapat berkembang mengikuti kebutuhan operasional.

Generic Domain dapat diganti atau ditingkatkan selama antarmukanya tetap kompatibel.

---

# Closing Statement

Core Domain Model adalah peta utama EPMP. Semua spesifikasi domain, modul, database, dan API harus mengacu pada model ini untuk menjaga konsistensi, modularitas, dan kemampuan sistem berkembang dalam jangka panjang.

---
