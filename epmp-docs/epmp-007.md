# EPMP-007

# Core Business Processes & Business Event Flow

```text
Document ID    : EPMP-007
Document Name  : Core Business Processes & Business Event Flow
Version        : 1.1.0
Status         : Implemented & Active
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

Fokus dokumen ini adalah **proses bisnis**, model event-driven, dan implementasi lifecycle aktif di platform EPMP.

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

| Process            | Domain Utama         | Deskripsi Implementasi |
| ------------------ | -------------------- | ---------------------- |
| Property Setup     | Property Management  | Registrasi Property, Building, Floor, Room, Bed Template & Bed assignment, visualisasi 3D WebGL. |
| Room Preparation   | Property / Asset     | Setup aset kamar, inspeksi, dan penentuan ketersediaan awal. |
| Reservation        | Reservation          | Pemesanan kamar/bed, down payment/booking fee, auto-hold status. |
| Check-In           | Contract + Occupancy | Aktivasi kontrak, penyerahan kunci/fasilitas, perubahan status kamar menjadi Occupied. |
| Occupancy          | Occupancy            | Manajemen penghuni, log pergerakan, pelaporan fasilitas. |
| Billing Cycle      | Billing & Finance    | Penerbitan Invoice tagihan sewa/utilitas berulang berbasis multi-currency (IDR, USD, EUR, SGD, MYR). |
| Payment Collection | Billing & Finance    | Pencatatan pembayaran, auto-reconciliation & trigger auto-paid pada invoice saat lunas. |
| Communication      | Communication        | Pairing gateway WhatsApp (whatsmeow), validasi nomor WhatsApp, blast pengumuman & tagihan otomatis. |
| Maintenance        | Maintenance          | Work order perbaikan fasilitas kamar/gedung, penugasan teknisi, isolasi kamar ke Maintenance state. |
| Contract Renewal   | Contract             | Perpanjangan masa sewa, penyesuaian tarif rental, re-kontrak. |
| Check-Out          | Contract + Occupancy | Inspeksi akhir, serah terima aset, pengosongan kamar (Vacant / Available). |
| Deposit Settlement | Finance              | Pengembalian deposit setelah pemotongan denda/kerusakan aset (refund / adjustment / penalty). |
| Reporting          | Reporting            | Analitik okupansi, laporan arus kas (cash flow), utilisasi unit 3D. |

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
| Room Lifecycle   | RoomCreated, RoomAvailabilityChanged (Available, Reserved, Occupied, Maintenance), BedAssigned |
| Billing          | InvoiceIssued, InvoiceOverdue, InvoiceSettledViaTrigger                            |
| Payment          | PaymentReceived, PaymentFailed, PaymentRefunded                                    |
| Deposit          | DepositCollected, DepositReturned, DepositForfeited                                |
| Communication    | WhatsAppDeviceConnected, WhatsAppDeviceDisconnected, WhatsAppMessageSent, WhatsAppMessageFailed, BlastCompleted |
| Maintenance      | MaintenanceRequested, WorkOrderCreated, MaintenanceCompleted                       |
| Asset            | AssetAssigned, AssetReturned, AssetInspected                                       |

Business event ini menjadi kontrak komunikasi antar bounded context.

---

# 6. High-Level Workflow per Domain

### Reservation
1. Pilih Room/Bed yang tersedia (Visual 3D atau list unit).
2. Buat Reservation dan kunci unit (status berubah ke Reserved).
3. Validasi aturan bisnis (masa berlaku, booking fee, dsb.).
4. Konfirmasi atau batalkan (jika batal, status kembali ke Available).

### Contract & Check-In
1. Buat Draft Contract berbasis penyewa & unit.
2. Review klausul & termin pembayaran (multi-currency: IDR, USD, EUR, SGD, MYR).
3. Aktivasi kontrak.
4. Check-In: unit otomatis berpindah ke status `Occupied` (Merah di kanvas 3D).

### Billing & Payment Auto-Settlement
1. Generate Invoice (tagihan sewa, deposit, utilitas).
2. Kirim notifikasi tagihan melalui WhatsApp Gateway.
3. Catat penerimaan pembayaran (Payment).
4. **Auto-Settlement Trigger (`migration 000038`)**: Database trigger menghitung total pembayaran yang sukses (`completed`) untuk invoice terkait. Jika total bayar >= total tagihan, status invoice seketika diubah menjadi `Paid` secara atomik.

### Communication & WhatsApp Gateway
1. **Device Pairing**: Scan QR code resmi WhatsApp (protokol whatsmeow multi-device) tanpa input manual nomor telepon.
2. **Auto-Detection**: Server otomatis menerima event JID login WhatsApp dan menyimpan nomor pemilik perangkat.
3. **Pre-Verification (`IsOnWhatsApp`)**: Setiap pengiriman pesan massal (Blast) melakukan verifikasi nomor ke server WhatsApp. Nomor yang tidak terdaftar otomatis ditandai `failed` tanpa membuang kuota koneksi.
4. **Message Dispatch & Templating**: Merender variabel dinamis (`{{tenant_name}}`, `{{room_name}}`, `{{invoice_amount}}`, dsb.) dan mengirim pesan teks/dokumen secara real-time.
5. **Delivery Audit Log**: Setiap pesan tercatat dalam `wa_message_logs` untuk rekam jejak pengiriman.

### Maintenance
1. Buat permintaan tiket maintenance (kamar atau fasilitas umum).
2. Buat Work Order dan tugaskan teknisi.
3. Unit kamar ditandai sebagai `Maintenance` (Biru di kanvas 3D).
4. Selesaikan perbaikan, verifikasi hasil, dan kembalikan unit ke `Available`.

---

# 7. Cross-Domain Business Rules

Beberapa aturan lintas domain yang menjadi dasar implementasi:

- **Room Availability Integrity**: Room/Bed tidak dapat di-_check-in_ tanpa Contract yang aktif.
- **Visual Spatial Sync**: Status ketersediaan kamar (`Available`, `Reserved`, `Occupied`, `Maintenance`) secara otomatis memantulkan kode warna standar pada visualisasi 3D WebGL (Hijau, Oranye, Merah, Biru).
- **Payment-Invoice Synchronization**: Pembayaran yang berstatus `completed` langsung memicu auto-update status invoice menjadi `Paid` lewat trigger basis data.
- **WhatsApp Phone Normalization**: Semua nomor tujuan lokal (format `08...`) dinormalisasi secara otomatis menjadi standar internasional `628...` sebelum validasi whatsmeow JID.
- **Contract Room Validity**: Contract tidak dapat diaktifkan tanpa Room/Bed yang valid dan tersedia.
- **Invoice Pre-requisite**: Invoice sewa tidak boleh diterbitkan untuk Contract yang belum aktif.
- **Maintenance Lock**: Room dalam status Maintenance tidak dapat dipesan atau dialokasikan untuk reservasi baru.
- **Deposit Refund Requirement**: Deposit hanya dapat dikembalikan setelah proses inspeksi aset saat check-out selesai.

---

# 8. State Transition Principle

Setiap proses bisnis harus mengikuti **state transition** yang eksplisit.

### Room Availability State Machine
```text
               ┌───────────────┐
               │   AVAILABLE   │ (Hijau)
               └───────┬───────┘
                       │
       ┌───────────────┴───────────────┐
       │ (Reservation)                 │ (Maintenance Request)
       ▼                               ▼
┌──────────────┐               ┌──────────────┐
│   RESERVED   │ (Oranye)      │ MAINTENANCE  │ (Biru)
└──────┬───────┘               └──────┬───────┘
       │                              │
       │ (Check-In)                   │ (Resolved)
       ▼                              │
┌──────────────┐                      │
│   OCCUPIED   │ (Merah)              │
└──────┬───────┘                      │
       │                              │
       │ (Check-Out)                  │
       └───────────────┬──────────────┘
                       ▼
               ┌───────────────┐
               │   AVAILABLE   │
               └───────────────┘
```

### WhatsApp Device State Machine
```text
┌─────────────────┐
│  DISCONNECTED   │
└────────┬────────┘
         │ (Initiate Pairing)
         ▼
┌─────────────────┐
│   CONNECTING    │ ───► (Generate QR Channel)
└────────┬────────┘
         │ (Successful Scan & Handshake)
         ▼
┌─────────────────┐
│    CONNECTED    │ (Device JID Stored, Ready to Send)
└─────────────────┘
```

---

# 9. Automation Opportunities & Active Triggers

EPMP mengimplementasikan otomasi langsung di tingkat aplikasi dan basis data:

- **Trigger Auto-Settlement (`000038_sync_invoice_status.up.sql`)**: Sinkronisasi instan status tagihan ke `Paid` saat pembayaran lunas tanpa menunggu cron job batch.
- **Whatsmeow WhatsApp Verification**: Pencegahan blast ke nomor palsu menggunakan pengecekan aktif `IsOnWhatsApp` langsung ke protokol WhatsApp.
- **Normalization Engine**: Auto-formatting nomor telepon Indonesia (`08...` -> `628...`).
- **Recurring Invoice Generation**: Penerbitan invoice bulanan berkala untuk kontrak aktif.
- **Auto-Expiration**: Mengubah Reservation menjadi Expired saat melewati batas waktu pembayaran booking fee.

---

# 10. Process Ownership

Setiap proses memiliki domain owner yang bertanggung jawab atas aturan bisnisnya:

| Process              | Owner                  |
| -------------------- | ---------------------- |
| Property Setup       | Property Management    |
| Reservation          | Reservation Management |
| Contract Lifecycle   | Contract Management    |
| Occupancy Lifecycle  | Occupancy Management   |
| Billing Cycle        | Billing & Finance      |
| Payment Collection   | Billing & Finance      |
| Communication        | Communication Gateway  |
| Asset Assignment     | Asset Management       |
| Maintenance Workflow | Maintenance Management |
| System Audit & Logs  | Platform Engineering   |

Domain lain dapat berpartisipasi melalui event, tetapi tidak mengambil alih kepemilikan proses.

---

# Closing Statement

Core Business Processes & Business Event Flow adalah jembatan antara model bisnis dan implementasi teknis. Dokumen ini memastikan setiap alur operasional memiliki state yang jelas, menghasilkan business event yang konsisten, dan diimplementasikan secara solid pada arsitektur Go + Echo + whatsmeow + React + Three.js tanpa mengorbankan integritas data dan modularitas domain.

---

