# EPMP-006

# Domain Map & Context Boundaries

```text
Document ID    : EPMP-006
Document Name  : Domain Map & Context Boundaries
Version        : 1.0.0
Status         : Draft
Owner          : Software Architecture Team
Dependencies   : EPMP-001, EPMP-002, EPMP-003, EPMP-004, EPMP-005
Referenced By  : Seluruh Domain Specification, Module Specification, API, Database
```

---

# 1. Purpose

Dokumen ini mendefinisikan **Bounded Context**, hubungan antar domain, aturan komunikasi, dan batas tanggung jawab masing-masing domain dalam Enterprise Property Management Platform (EPMP).

Tujuannya adalah memastikan:

- Domain tetap independen.
- Coupling antar modul rendah.
- Business Rules tidak tercampur.
- Pengembangan dapat dilakukan oleh beberapa tim secara paralel.
- AI Agent memiliki batas konteks yang jelas ketika menghasilkan kode.

---

# 2. What is a Bounded Context?

Dalam EPMP, **Bounded Context** adalah batas logis tempat sebuah model bisnis berlaku.

Di dalam sebuah context:

- istilah memiliki satu arti,
- business rule konsisten,
- entity dimiliki oleh context tersebut.

Di luar context tersebut, entity hanya boleh diakses melalui kontrak yang telah ditentukan.

---

# 3. Core Context Map

```text
                 +----------------------+
                 |     Organization     |
                 +----------+-----------+
                            |
        +-------------------+-------------------+
        |                                       |
+-------v--------+                     +--------v-------+
| Configuration  |                     | Identity & IAM |
+-------+--------+                     +--------+-------+
        |                                       |
        +-------------------+-------------------+
                            |
                    +-------v--------+
                    | Property Mgmt  |
                    +-------+--------+
                            |
             +--------------+--------------+
             |                             |
     +-------v-------+             +-------v-------+
     | Reservation   |             | Asset Mgmt    |
     +-------+-------+             +-------+-------+
             |                             |
     +-------v-------+             +-------v-------+
     | Contract Mgmt |             | Maintenance   |
     +-------+-------+             +---------------+
             |
     +-------v-------+
     | Occupancy     |
     +-------+-------+
             |
     +-------v-------+
     | Billing       |
     +-------+-------+
             |
     +-------v-------+
     | Reporting     |
     +---------------+
```

Diagram ini menggambarkan hubungan konseptual antar bounded context, bukan urutan eksekusi proses bisnis.

---

# 4. Bounded Context List

| Context                | Purpose                               | Owns                                       |
| ---------------------- | ------------------------------------- | ------------------------------------------ |
| Organization           | Mengelola organisasi penyedia layanan | Organization                               |
| Configuration          | Mengelola konfigurasi sistem          | Settings, Master Data                      |
| Identity & IAM         | Mengelola identitas dan hak akses     | User, Role, Permission                     |
| Property Management    | Mengelola struktur fisik properti     | Property, Building, Floor, Zone, Room, Bed |
| Tenant Management      | Mengelola data penyewa                | Tenant, Documents, Contacts                |
| Reservation Management | Mengelola pemesanan                   | Reservation, Booking Fee                   |
| Contract Management    | Mengelola perjanjian sewa             | Contract                                   |
| Occupancy Management   | Mengelola status penggunaan unit      | Occupancy                                  |
| Billing & Finance      | Mengelola transaksi keuangan          | Invoice, Payment, Deposit                  |
| Asset Management       | Mengelola aset                        | Asset                                      |
| Maintenance Management | Mengelola pemeliharaan                | Work Order                                 |
| Reporting              | Menyediakan data analitik             | Read Models                                |

---

# 5. Domain Ownership Rules

Sebuah entity hanya boleh dimiliki oleh satu bounded context.

Contoh:

| Entity   | Owner Context       |
| -------- | ------------------- |
| Room     | Property Management |
| Contract | Contract Management |
| Invoice  | Billing & Finance   |
| Tenant   | Tenant Management   |
| Asset    | Asset Management    |

Context lain **tidak boleh** mengubah entity tersebut secara langsung.

---

# 6. Allowed Communication

Komunikasi antar context dilakukan melalui:

1. Application Service
2. Domain Event
3. Query Interface (Read Only)
4. Published API

Komunikasi langsung ke database context lain tidak diperbolehkan.

---

# 7. Communication Patterns

### Command

Digunakan ketika satu context meminta context lain melakukan aksi.

Contoh:

```text
Create Contract
Assign Asset
Issue Invoice
```

Command bersifat eksplisit dan ditujukan kepada satu context.

---

### Query

Digunakan untuk membaca data tanpa mengubah state.

Contoh:

```text
Get Available Rooms
Get Active Contracts
Get Tenant Profile
```

Query tidak boleh memiliki efek samping.

---

### Domain Event

Digunakan untuk memberi tahu context lain bahwa sesuatu telah terjadi.

Contoh:

```text
ContractActivated
InvoicePaid
RoomVacated
ReservationExpired
```

Event tidak memerintahkan context lain, tetapi menginformasikan fakta yang telah terjadi.

---

# 8. Upstream & Downstream Relationships

Beberapa context menjadi sumber data bagi context lain.

| Upstream    | Downstream  |
| ----------- | ----------- |
| Property    | Reservation |
| Reservation | Contract    |
| Contract    | Occupancy   |
| Occupancy   | Billing     |
| Billing     | Reporting   |

Downstream boleh menggunakan informasi dari upstream, tetapi tidak boleh mengubah data upstream.

---

# 9. Anti-Corruption Layer (ACL)

Jika di masa depan EPMP terintegrasi dengan sistem eksternal (ERP, PMS lain, Payment Gateway, IoT), setiap integrasi harus melalui **Anti-Corruption Layer**.

Tujuannya adalah:

- Melindungi model domain EPMP.
- Mencegah istilah atau struktur sistem eksternal mencemari model internal.
- Memudahkan penggantian penyedia layanan eksternal.

---

# 10. Shared Kernel

Beberapa konsep digunakan lintas context dan harus konsisten.

Contohnya:

- Money
- Address
- Date Range
- Identifier
- Attachment
- Audit Information
- Status Reference

Shared Kernel harus dijaga seminimal mungkin agar tidak menciptakan coupling yang tinggi.

---

# 11. Context Independence

Setiap bounded context harus dapat berkembang secara mandiri.

Sebagai contoh:

- Property Management dapat menambahkan konsep `Tower` tanpa memengaruhi Billing.
- Billing dapat mendukung metode pembayaran baru tanpa memengaruhi Property.
- Asset dapat menambahkan kategori aset tanpa mengubah Contract.

---

# 12. Context Lifecycle

Setiap context memiliki siklus pengembangan sendiri.

```text
Planning
↓

Development
↓

Testing
↓

Release
↓

Maintenance
↓

Evolution
```

Perubahan dalam satu context tidak boleh mengharuskan perubahan pada seluruh sistem.

---

# 13. AI Development Context

Dokumentasi ini juga menjadi dasar pembagian konteks untuk AI Coding Agent.

Idealnya satu AI Agent hanya menerima spesifikasi untuk satu bounded context dalam satu sesi.

Contoh:

- AI #1 → Property Management
- AI #2 → Contract Management
- AI #3 → Billing & Finance

Pendekatan ini membantu AI menghasilkan implementasi yang lebih fokus, mengurangi konteks yang harus diproses, dan meminimalkan risiko mencampurkan business rule antar domain.

---

# 14. Future Context Expansion

Arsitektur memungkinkan penambahan bounded context baru tanpa mengubah context yang telah ada.

Contoh:

- Visitor Management
- Parking Management
- Housekeeping
- Smart Lock
- Utility Meter
- CRM
- Loyalty
- AI Recommendation
- Workflow Automation

Context baru harus:

1. Memiliki owner yang jelas.
2. Memiliki model domain sendiri.
3. Berkomunikasi melalui kontrak yang terdokumentasi.
4. Tidak mengakses database context lain secara langsung.

---

# 15. Context Evolution Principles

Ketika domain berkembang, perubahan harus mengikuti prinsip berikut:

- **Tambah sebelum mengubah**: lebih baik menambahkan extension point daripada mengubah perilaku inti.
- **Jaga backward compatibility** untuk API dan event jika memungkinkan.
- **Pisahkan model tulis dan model baca** bila kebutuhan analitik mulai kompleks.
- **Refactor di dalam context**, bukan dengan memindahkan tanggung jawab ke context lain tanpa alasan bisnis yang kuat.

---

# Closing Statement

Bounded Context adalah batas organisasi pengetahuan dalam EPMP. Dengan batas yang jelas, setiap domain dapat berkembang secara independen, tim dapat bekerja secara paralel, dan AI Coding Agent dapat menghasilkan implementasi yang konsisten tanpa mencampurkan aturan bisnis dari domain lain.

---
