# EPMP-001

# Project Overview

```text
Document ID    : EPMP-001
Document Name  : Project Overview
Version        : 1.2.0
Status         : Implemented & Active
Owner          : Product & Architecture Team
Reviewer       : -
Dependencies   : README.md
Referenced By  : Seluruh dokumentasi EPMP
```

---

# 1. Introduction

Enterprise Property Management Platform (EPMP) adalah sebuah platform enterprise yang dirancang untuk mengelola berbagai jenis properti sewa melalui satu sistem yang fleksibel, modular, scalable, dan dapat dikonfigurasi.

EPMP tidak dibangun untuk satu model bisnis tertentu, melainkan sebagai platform yang mampu beradaptasi dengan berbagai kebutuhan operasional properti melalui konfigurasi, tanpa memerlukan perubahan pada source code inti.

Platform ini menjadi fondasi untuk pengelolaan:

- Boarding House (Kost)
- Apartment
- Dormitory
- Student Housing
- Co-Living
- Guest House
- Villa
- Warehouse
- Commercial Building
- Office Rental
- Storage
- Mixed Property

dan jenis properti lainnya.

---

# 2. Background

Sebagian besar software property management yang tersedia saat ini dibuat untuk satu jenis bisnis tertentu.

Contohnya:

- Software khusus hotel
- Software khusus apartemen
- Software khusus kost
- Software khusus gudang

Pendekatan tersebut menyebabkan keterbatasan ketika kebutuhan bisnis berkembang.

Sebagai contoh:

Hari ini sebuah bisnis hanya memiliki satu gedung kost.

Beberapa tahun kemudian bisnis tersebut berkembang menjadi:

- beberapa gedung
- berbagai tipe properti
- berbagai jenis kamar
- berbagai aturan deposit
- berbagai metode pembayaran
- berbagai konfigurasi kontrak

Sebagian besar aplikasi harus dimodifikasi secara besar-besaran untuk mengakomodasi perubahan tersebut.

EPMP dibangun untuk menghilangkan keterbatasan tersebut.

---

# 3. Problem Statement

EPMP dikembangkan untuk menyelesaikan beberapa permasalahan utama dalam industri property management.

## 3.1 Hardcoded Business Rules

Sebagian besar aplikasi memiliki aturan bisnis yang ditanam langsung di dalam kode aplikasi.

Contoh:

- Deposit selalu Rp1.000.000
- Booking Fee wajib
- Termin hanya 3 bulan, 6 bulan, dan 12 bulan
- Hanya ada satu gedung
- Hanya ada satu lantai
- Hanya ada AC dan Non-AC

Pendekatan tersebut membuat aplikasi sulit berkembang.

---

## 3.2 Poor Scalability

Banyak aplikasi dirancang hanya untuk satu properti.

Ketika bisnis berkembang menjadi multi-property, arsitektur harus diubah secara besar.

---

## 3.3 Poor Extensibility

Penambahan fitur baru sering kali memerlukan perubahan pada fitur lama.

Hal ini meningkatkan risiko bug dan biaya pengembangan.

---

## 3.4 Fragmented Data

Data penghuni, pembayaran, aset, maintenance, dan laporan sering berada pada sistem yang berbeda.

Akibatnya:

- data sulit dianalisis
- laporan tidak real-time
- audit sulit dilakukan

---

## 3.5 Limited Business Adaptability

Setiap bisnis memiliki aturan yang berbeda.

Misalnya:

- kontrak
- deposit
- denda
- booking
- maintenance

Software harus mampu mengikuti bisnis.

Bukan bisnis yang mengikuti software.

---

# 4. Product Vision

Menjadi platform property management enterprise yang fleksibel, modular, dan scalable sehingga dapat digunakan oleh berbagai jenis bisnis properti tanpa perubahan pada arsitektur inti sistem.

---

# 5. Product Mission

Menyediakan platform property management modern yang:

- mudah dikembangkan
- mudah dikonfigurasi
- mudah diintegrasikan
- mudah dipelihara
- ramah terhadap AI-assisted development

---

# 6. Product Objectives

EPMP memiliki beberapa tujuan utama.

## Objective 1

Menyediakan platform yang dapat digunakan untuk berbagai jenis properti.

---

## Objective 2

Menghilangkan business rule yang bersifat hardcoded.

---

## Objective 3

Membangun arsitektur modular yang mudah dikembangkan.

---

## Objective 4

Memungkinkan seluruh proses operasional properti berjalan melalui satu platform.

---

## Objective 5

Menjadi fondasi jangka panjang untuk pengembangan produk enterprise.

---

# 7. Product Scope

Versi implementasi saat ini mencakup domain-domain fungsional berikut:

## Property Management & Spatial Visualization
- **Property, Building, Floor, Zone, Room**
- **Bed Management & Templates**: Dukungan sewa kamar privat (*Whole-Room Rental*) maupun sewa tempat tidur (*Bed/Coliving*), dilengkapi template tempat tidur (Single, Queen, King, Bunk).
- **Interactive Spatial Views**: Visualisasi 3D WebGL interaktif gedung dan denah lantai (*Three.js Building Drop*) dengan indikator status ketersediaan warna real-time (Hijau=Available/Vacant, Oranye=Reserved, Merah=Occupied, Biru=Maintenance).

---

## Tenant Management
- Tenant Profiles & Multi-tenant isolation (Org-Scoped)
- Tenant Identity, Documents, dan Emergency Contacts
- Tenant Rental & Payment History

---

## Reservation & Booking
- Pemesanan kamar/bed sebelum tanggal sewa
- Validasi ketersediaan dan status transisi otomatis

---

## Contract & Occupancy
- Manajemen Kontrak Sewa (Durasi, Nilai Sewa, Tanggal Jatuh Tempo)
- Siklus Occupancy: Check-In, Active Staying, Check-Out
- Hubungan terintegrasi: Tenant → Room/Bed → Contract → Occupancy

---

## Financial, Billing & Multi-Currency
- **Multi-Currency Support**: Dukungan pemilihan mata uang (IDR, USD, EUR, SGD, MYR) untuk fleksibilitas multi-country.
- **Invoices**: Penagihan otomatis berbasis kontrak & tenant dengan shortcut pembuatan langsung dari detail kontrak.
- **Payments**: Pencatatan pembayaran dengan **Auto-Sync Status Invoice** (ketika pembayaran lunas, status invoice otomatis berubah menjadi 'Paid' secara real-time via database trigger).
- **Deposits**: Pengelolaan uang jaminan sewa (Security Deposit).
- **Charges, Adjustments, Penalties, & Refunds**: Manajemen biaya tambahan, penyesuaian/diskon, denda keterlambatan, dan pengembalian dana.
- **Human-Readable Presentation**: Penyajian nama referensi eksplisit (Nama Tenant, Nomor Kontrak, Nama Kamar) dan format tanggal lokal yang mudah dibaca (misal: 15 Sep 2026).

---

## Communication & WhatsApp Multi-Device Gateway
- **Native WhatsApp Multi-Device Gateway (`whatsmeow`)**: Terhubung langsung ke protokol resmi WhatsApp tanpa service perantara berbayar.
- **Live QR Code Pairing**: Menampilkan QR code noise handshake resmi WhatsApp secara real-time menggunakan `react-qr-code`.
- **Zero-Manual Phone Input**: Mendeteksi nomor handphone secara otomatis melalui callback server WhatsApp saat QR code berhasil dipindai oleh HP.
- **Verifikasi Nomor WhatsApp Terdaftar (`IsOnWhatsApp`)**: Memvalidasi nomor sebelum pengiriman dan menormalisasi awalan lokal `08...` ke format internasional `628...` guna mencegah pemblokiran atau kegagalan kirim.
- **Blast Message (Broadcast)**: Pengiriman pesan massal berbasis target audiens (Semua Tenant, Tenant Menunggak) dengan variabel dinamis (`{{tenant_name}}`, `{{invoice_amount}}`, dll.) dan log audit per penerima yang transparan.

---

## Asset & Inventory Management
- Master data aset properti
- Asset Assignment (Penempatan aset ke kamar/gedung)
- Asset Inspection (Pemeriksaan kondisi aset saat check-in/check-out)

---

## Maintenance & Work Orders
- Pelaporan masalah & perbaikan
- Work Order lifecycle (Draft, Assigned, In Progress, Completed)
- Manajemen Teknisi & Supplier pihak ketiga

---

## Dashboard & Analytics
- Metrik real-time: Tingkat hunian (Occupancy Rate), Total Pendapatan, Properti Aktif, Tenant Aktif
- Visualisasi ringkasan status sewa dan ketersediaan kamar

---

## Quality Assurance & E2E Testing Standard
- Pengujian otomatis berbasis Playwright (`make test-e2e`) mencakup 40+ rute dashboard dan visualisasi 3D Three.js.

---

# 8. Out of Scope

Beberapa fitur berikut tidak menjadi bagian dari versi awal.

- Accounting ERP
- Payroll
- HR Management
- CRM Marketing
- Marketplace
- Hotel Booking Engine
- POS System

Namun arsitektur harus memungkinkan integrasi di masa depan.

---

# 9. Target Users

EPMP dirancang untuk berbagai jenis pengguna.

- Organization Owner
- Property Owner
- Property Manager
- Operational Staff
- Finance Staff
- Maintenance Staff
- Receptionist
- Tenant
- Vendor

---

# 10. Business Model

EPMP mendukung berbagai model bisnis.

- Single Property
- Multi Property
- Multi Organization
- Franchise
- Enterprise

---

# 11. Core Capabilities

Platform harus mampu:

- mengelola properti
- mengelola tenant
- mengelola kontrak
- mengelola pembayaran
- mengelola aset
- mengelola maintenance
- menghasilkan laporan
- melakukan audit
- melakukan automasi

---

# 12. Success Metrics

Keberhasilan EPMP diukur melalui beberapa indikator.

### Functional

Seluruh modul utama dapat berjalan secara independen.

---

### Technical

Arsitektur mampu menangani penambahan fitur tanpa refactor besar.

---

### Business

Platform dapat digunakan untuk lebih dari satu jenis bisnis properti.

---

### Operational

Seluruh proses operasional dapat dilakukan melalui satu platform.

---

# 13. Long-Term Vision

Dalam lima tahun ke depan EPMP diharapkan berkembang menjadi platform enterprise yang mendukung:

- SaaS Multi Tenant
- Public API
- Marketplace
- Plugin System
- Mobile Apps
- Smart Lock
- IoT
- AI Analytics
- Predictive Maintenance
- Dynamic Pricing

---

# 14. Guiding Philosophy

EPMP dibangun berdasarkan satu filosofi utama.

> **Everything should be configurable. Nothing should be hardcoded.**

Seluruh keputusan desain, arsitektur, dan implementasi harus mengacu pada prinsip tersebut.

---

# 15. References

Dokumen ini menjadi referensi utama bagi seluruh dokumentasi EPMP selanjutnya.

Dokumen yang bergantung pada EPMP-001 meliputi:

- EPMP-002 Product Principles
- EPMP-003 Architecture Overview
- EPMP-004 Domain Model
- Seluruh Module Specification
- Seluruh API Specification
- Seluruh Database Specification

---
