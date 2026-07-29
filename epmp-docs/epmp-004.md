Mantap.

Saya justru sangat senang kita memutuskan membuat **EPMP-004** sekarang.

Karena menurut pengalaman saya, salah satu penyebab proyek enterprise menjadi berantakan bukan karena coding, tetapi karena **setiap orang menggunakan istilah yang berbeda untuk objek yang sama**.

Contohnya saja:

Developer A mengatakan **Room**.

Developer B mengatakan **Unit**.

Designer menulis **Kamar**.

Database memakai **Rooms**.

API memakai **Units**.

AI nanti akan mulai bingung.

Karena itu kita akan membuat **bahasa resmi EPMP**.

---

---

# EPMP-004

# Ubiquitous Language & Business Glossary

```text
Document ID    : EPMP-004
Document Name  : Ubiquitous Language & Business Glossary
Version        : 1.0.0
Status         : Draft
Owner          : Product & Architecture
Dependencies   : EPMP-001, EPMP-002, EPMP-003
Referenced By  : Seluruh dokumentasi EPMP
```

---

# 1. Purpose

Dokumen ini mendefinisikan bahasa resmi (Ubiquitous Language) yang digunakan di seluruh Enterprise Property Management Platform (EPMP).

Seluruh stakeholder, termasuk:

- Product Owner
- Software Architect
- UI/UX Designer
- Backend Developer
- Frontend Developer
- QA Engineer
- DevOps Engineer
- AI Coding Agent

wajib menggunakan istilah yang didefinisikan dalam dokumen ini.

Tujuannya adalah menghilangkan ambiguitas dan memastikan setiap istilah memiliki satu arti yang konsisten.

---

# 2. Language Principles

EPMP menggunakan beberapa prinsip dasar dalam penamaan.

## 2.1 One Concept, One Name

Satu konsep hanya boleh memiliki satu nama resmi.

Contoh:

❌ Customer

❌ Resident

❌ Occupant

✅ Tenant

Mulai sekarang seluruh sistem menggunakan istilah **Tenant**.

---

## 2.2 One Name, One Meaning

Satu istilah hanya boleh memiliki satu arti.

Misalnya:

Room tidak boleh berarti:

- kamar
- bangunan
- properti

Room hanya berarti satu ruang yang dapat disewakan.

---

## 2.3 Business Language First

Nama entity mengikuti istilah yang dipakai bisnis.

Bukan mengikuti nama tabel database.

---

# 3. Organization Domain

## Organization

Entity tertinggi dalam EPMP.

Mewakili perusahaan atau organisasi yang mengelola properti.

Contoh:

- PT ABC Property
- Flosse Group
- XYZ Residence

Satu Organization dapat memiliki banyak Property.

---

# 4. Property Domain

## Property

Lokasi utama yang dikelola.

Property bukan gedung.

Property adalah satu kawasan atau lokasi bisnis.

Contoh:

- Flosse House
- Flosse Residence
- ABC Apartment

Satu Property dapat memiliki banyak Building.

---

## Building

Gedung fisik di dalam Property.

Contoh:

Building A

Building B

Building Timur

Tower North

---

## Floor

Lantai di dalam Building.

Contoh:

Ground Floor

Level 2

Floor 15

---

## Zone

Area logis dalam satu lantai.

Zone bersifat opsional.

Contoh:

North Wing

South Wing

West Area

VIP Area

---

## Room

Unit fisik yang dapat disewakan.

Contoh:

101

102

305

A-12

WH-05

Room selalu berada di bawah:

Property

↓

Building

↓

Floor

↓

Zone (optional)

↓

Room

---

## Bed

Sub-unit dari Room.

Digunakan apabila satu Room memiliki lebih dari satu penyewa.

Contoh:

Room 201

↓

Bed A

↓

Bed B

↓

Bed C

↓

Bed D

---

# 5. Tenant Domain

## Tenant

Individu yang menyewa Room atau Bed.

Tenant adalah istilah resmi.

Bukan:

- Customer
- Guest
- Resident
- User

---

## Occupancy

Status apakah Room sedang ditempati.

Occupancy bukan berarti kontrak aktif.

Occupancy berarti kondisi fisik penggunaan ruang.

---

## Check-In

Proses Tenant mulai menempati Room.

---

## Check-Out

Proses Tenant meninggalkan Room.

---

## Move-In

Aktivitas fisik masuk ke Room.

Biasanya bersamaan dengan Check-In.

---

## Move-Out

Aktivitas fisik keluar dari Room.

---

# 6. Reservation Domain

## Reservation

Proses memesan Room sebelum kontrak dibuat.

Reservation dapat:

- Pending
- Confirmed
- Cancelled
- Expired

---

## Booking Fee

Biaya yang dibayarkan untuk mengamankan Reservation.

Booking Fee belum tentu wajib.

Konfigurasinya tergantung Property.

---

## Reservation Expiry

Batas waktu Reservation.

Apabila lewat maka Room kembali tersedia.

---

# 7. Contract Domain

## Contract

Perjanjian sewa antara Tenant dan Organization.

Contract mengatur:

- durasi
- harga
- aturan
- pembayaran
- deposit

---

## Lease

Sinonim Contract.

Namun dalam EPMP istilah resmi adalah:

Contract.

---

## Renewal

Perpanjangan Contract.

---

## Extension

Penambahan durasi Contract sebelum Contract berakhir.

---

## Termination

Pengakhiran Contract.

---

# 8. Finance Domain

## Invoice

Dokumen tagihan.

Invoice belum tentu sudah dibayar.

---

## Payment

Pencatatan pembayaran.

---

## Charge

Komponen biaya.

Contoh:

Rental

Water

Electricity

Cleaning

Parking

Internet

---

## Deposit

Dana jaminan.

Deposit dapat:

- refundable
- non-refundable
- partial refundable

---

## Refund

Pengembalian dana.

---

## Adjustment

Perubahan nominal tagihan.

---

# 9. Asset Domain

## Asset

Barang milik Organization.

Contoh:

AC

Bed

Table

Chair

Key

TV

Water Heater

---

## Asset Assignment

Penempatan Asset pada Room.

---

## Asset Inspection

Pemeriksaan kondisi Asset.

---

# 10. Maintenance Domain

## Maintenance

Perawatan Asset atau Room.

---

## Work Order

Perintah kerja Maintenance.

---

## Technician

Petugas yang mengerjakan Maintenance.

---

## Inspection

Pemeriksaan kondisi Property.

---

# 11. Status Terminology

## Available

Dapat digunakan.

---

## Occupied

Sedang digunakan.

---

## Reserved

Sudah dipesan.

---

## Maintenance

Sedang diperbaiki.

---

## Cleaning

Sedang dibersihkan.

---

## Disabled

Tidak dapat digunakan.

---

# 12. Financial Status

## Draft

Belum aktif.

---

## Issued

Sudah diterbitkan.

---

## Paid

Sudah dibayar.

---

## Partial

Dibayar sebagian.

---

## Overdue

Lewat jatuh tempo.

---

## Cancelled

Dibatalkan.

---

# 13. Business Events

Event resmi EPMP.

ReservationCreated

ReservationConfirmed

ReservationExpired

TenantCheckedIn

TenantCheckedOut

ContractActivated

ContractRenewed

InvoiceGenerated

PaymentReceived

DepositReturned

MaintenanceRequested

MaintenanceCompleted

AssetAssigned

AssetReturned

---

# 14. Naming Convention

Nama Entity:

Singular

✅ Property

❌ Properties

---

API:

Plural

/property

/properties

/tenants

/contracts

---

Database Table

Plural

properties

rooms

contracts

payments

---

Class

Singular

Property

Room

Contract

Payment

---

# 15. Forbidden Terminology

Istilah berikut tidak boleh digunakan.

Customer

Guest

Member

Client

Resident

House

Boarding

Kos

Karena semuanya sudah diwakili oleh entity resmi.

---

# 16. Official Vocabulary

Mulai saat ini seluruh dokumentasi EPMP menggunakan istilah yang ada di dokumen ini.

Apabila muncul istilah baru, maka dokumen ini harus diperbarui terlebih dahulu sebelum istilah tersebut digunakan di dokumen lain.

---
