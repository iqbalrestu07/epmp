# EPMP-007

# Core Business Processes & Business Event Flow

```text
Document ID    : EPMP-007
Document Name  : Core Business Processes & Business Event Flow
Version        : 1.0.0
Status         : Draft
Owner          : Product & Architecture
Dependencies   : EPMP-001 s.d. EPMP-006
Referenced By  : Seluruh Domain Specification, Workflow, Automation, API, UI
```

---

# 1. Purpose

Dokumen ini mendefinisikan **alur bisnis utama** (Core Business Processes) dalam EPMP beserta **Business Event** yang dihasilkan pada setiap tahapan.

Dokumen ini menjadi acuan untuk:

- Business Workflow
- State Machine
- Automation
- Domain Events
- Use Case
- API Design
- UI Flow

Fokus dokumen ini adalah **proses bisnis**, bukan implementasi teknis.

---

# 2. Business Process Principles

Setiap proses bisnis di EPMP harus memenuhi prinsip berikut:

- Memiliki titik awal dan akhir yang jelas.
- Menghasilkan perubahan state yang terdokumentasi.
- Menghasilkan Business Event ketika terjadi perubahan penting.
- Tidak melakukan perubahan pada domain lain secara langsung.
- Dapat diaudit (Audit Trail).

---

# 3. Core Business Processes

EPMP memiliki proses bisnis inti berikut:

| Process            | Domain Utama         |
| ------------------ | -------------------- |
| Property Setup     | Property Management  |
| Room Preparation   | Property / Asset     |
| Reservation        | Reservation          |
| Check-In           | Contract + Occupancy |
| Occupancy          | Occupancy            |
| Billing Cycle      | Billing & Finance    |
| Payment Collection | Billing & Finance    |
| Maintenance        | Maintenance          |
| Contract Renewal   | Contract             |
| Check-Out          | Contract + Occupancy |
| Deposit Settlement | Finance              |
| Reporting          | Reporting            |

---

# 4. End-to-End Tenant Journey

Alur utama penyewa digambarkan sebagai berikut:

```text
Property Ready
        │
        ▼
Room Available
        │
        ▼
Reservation Created
        │
        ▼
Reservation Confirmed
        │
        ▼
Contract Created
        │
        ▼
Contract Activated
        │
        ▼
Check-In
        │
        ▼
Occupied
        │
        ▼
Recurring Billing
        │
        ▼
Payment Received
        │
        ▼
Contract Renewal
        │
     atau
        │
        ▼
Check-Out
        │
        ▼
Inspection
        │
        ▼
Deposit Settlement
        │
        ▼
Room Cleaning
        │
        ▼
Room Available
```

Siklus ini merupakan **happy path**. Variasi seperti pembatalan reservasi, gagal bayar, atau terminasi dini akan didefinisikan pada spesifikasi domain terkait.

---

# 5. Business Event Flow

Peristiwa utama yang menjadi penghubung antar domain:

| Business Process | Event                                                                              |
| ---------------- | ---------------------------------------------------------------------------------- |
| Reservation      | ReservationCreated, ReservationConfirmed, ReservationCancelled, ReservationExpired |
| Contract         | ContractCreated, ContractActivated, ContractRenewed, ContractTerminated            |
| Occupancy        | TenantCheckedIn, TenantCheckedOut, RoomOccupied, RoomVacated                       |
| Billing          | InvoiceIssued, InvoiceOverdue                                                      |
| Payment          | PaymentReceived, PaymentFailed, PaymentRefunded                                    |
| Deposit          | DepositCollected, DepositReturned, DepositForfeited                                |
| Maintenance      | MaintenanceRequested, WorkOrderCreated, MaintenanceCompleted                       |
| Asset            | AssetAssigned, AssetReturned, AssetInspected                                       |

Business event ini menjadi kontrak komunikasi antar bounded context.

---

# 6. High-Level Workflow per Domain

### Reservation

1. Pilih Room/Bed yang tersedia.
2. Buat Reservation.
3. Validasi aturan bisnis (masa berlaku, booking fee, dsb.).
4. Konfirmasi atau batalkan.

### Contract

1. Buat Draft Contract.
2. Review.
3. Aktivasi.
4. Kelola perubahan (renewal, extension, termination).

### Occupancy

1. Check-In.
2. Ubah status Room/Bed menjadi Occupied.
3. Pantau status hingga Check-Out.

### Billing

1. Generate Invoice.
2. Kirim tagihan.
3. Catat pembayaran.
4. Tangani keterlambatan, denda, dan penyesuaian.

### Maintenance

1. Buat permintaan.
2. Buat Work Order.
3. Kerjakan.
4. Inspeksi.
5. Selesaikan.

---

# 7. Cross-Domain Business Rules

Beberapa aturan lintas domain yang menjadi dasar implementasi:

- Room tidak dapat di-_check-in_ tanpa Contract yang aktif.
- Contract tidak dapat diaktifkan tanpa Room/Bed yang valid.
- Invoice tidak boleh diterbitkan untuk Contract yang belum aktif.
- Room dalam status Maintenance tidak dapat dipesan.
- Deposit hanya dapat dikembalikan setelah proses inspeksi selesai.

Detail aturan akan dibahas pada spesifikasi masing-masing domain.

---

# 8. State Transition Principle

Setiap proses bisnis harus mengikuti **state transition** yang eksplisit.

Contoh untuk Reservation:

```text
Draft
  │
  ▼
Pending
  │
  ├──► Confirmed
  │
  ├──► Cancelled
  │
  └──► Expired
```

Tidak diperbolehkan melakukan perpindahan state yang tidak didefinisikan.

---

# 9. Automation Opportunities

EPMP dirancang agar proses tertentu dapat diotomatisasi, misalnya:

- Mengubah Reservation menjadi Expired ketika melewati batas waktu.
- Membuat Invoice bulanan secara otomatis.
- Mengirim pengingat sebelum jatuh tempo.
- Mengubah status Contract menjadi Expired setelah masa berlaku berakhir.
- Menjadwalkan inspeksi saat Check-Out selesai.

Dokumen ini hanya mengidentifikasi peluang otomasi. Mekanisme implementasi akan dijelaskan pada dokumen Automation di fase engineering.

---

# 10. Process Ownership

Setiap proses memiliki domain owner yang bertanggung jawab atas aturan bisnisnya.

| Process              | Owner                  |
| -------------------- | ---------------------- |
| Property Setup       | Property Management    |
| Reservation          | Reservation Management |
| Contract Lifecycle   | Contract Management    |
| Occupancy Lifecycle  | Occupancy Management   |
| Billing Cycle        | Billing & Finance      |
| Asset Assignment     | Asset Management       |
| Maintenance Workflow | Maintenance Management |

Domain lain dapat berpartisipasi melalui event, tetapi tidak mengambil alih kepemilikan proses.

---

# Closing Statement

Core Business Processes & Business Event Flow adalah jembatan antara model bisnis dan implementasi teknis. Dokumen ini memastikan setiap alur operasional memiliki state yang jelas, menghasilkan business event yang konsisten, dan dapat diimplementasikan pada arsitektur Go + React berbasis REST API tanpa mengorbankan modularitas domain.

---
