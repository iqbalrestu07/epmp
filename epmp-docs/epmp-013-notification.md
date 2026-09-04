# EPMP-013

# Notification & Messaging Service Architecture

```text
Document ID    : EPMP-013
Document Name  : Notification & Messaging Service Architecture
Version        : 1.0.0
Status         : Draft
Owner          : Software Architecture Team
Dependencies   : EPMP-003, EPMP-006, EPMP-007
Referenced By  : Notification Module, Frontend Messaging Settings
```

---

# 1. Purpose

Dokumen ini mendefinisikan arsitektur dan spesifikasi untuk **Notification & Messaging Service** dalam ekosistem Enterprise Property Management Platform (EPMP). Layanan ini bertanggung jawab atas penyampaian informasi real-time kepada pengguna internal (Web Dashboard) maupun pengguna eksternal (Tenant/Penghuni).

Sistem pengiriman pesan dibangun dengan kapabilitas multi-channel, dengan fokus utama pada:
- **Web Dashboard Notification**: Menggunakan *WebSocket* untuk *real-time push notification* ke pengguna aplikasi.
- **WhatsApp Integration**: Menggunakan *Go WhatsApp Client (`whatsmeow`)* untuk komunikasi langsung dengan penghuni.

---

# 2. Key Capabilities

1. **In-App Real-time Notification (WebSocket)**
   Memberikan notifikasi langsung pada UI ketika suatu *Business Event* terjadi (misal: Invoice Dibayar, Work Order Selesai).

2. **WhatsApp Multi-Device Management**
   - Pengguna (Organization) dapat menghubungkan akun WhatsApp (Multi-Device) ke dalam EPMP.
   - Autentikasi dilakukan secara mandiri oleh pengguna melalui fitur *Scan QR Code* langsung dari dashboard EPMP.
   - Sistem mendukung koneksi ke lebih dari satu nomor WhatsApp (Multiple Devices) dalam satu organisasi, sehingga admin dapat memilih "Nomor Pengirim" yang spesifik.

3. **Direct & Blast Messaging**
   - **Direct Message**: Sistem secara otomatis mengirim pesan WhatsApp ke penghuni berdasarkan event (misal: Tagihan Jatuh Tempo, Notifikasi Check-in).
   - **Blast Message (Broadcast)**: Admin dapat mengirim pengumuman atau pesan promosi ke banyak penghuni sekaligus, dengan memilih *Device/Nomor WhatsApp* yang akan digunakan sebagai pengirim.

---

# 3. High-Level Architecture

```text
+-------------------+       +-----------------------+       +-------------------+
|   React Web App   |       |   Go Backend (EPMP)   |       | External Services |
|                   |       |                       |       |                   |
| 1. Scan QR Code   | <---> | 1. WhatsApp Handler   | <---> | WhatsApp Servers  |
| 2. Blast Message  | REST  |    (uses whatsmeow)   |       |                   |
|                   |       |                       |       |                   |
| 3. In-App Bell    | <---> | 2. WebSocket Hub      |       |                   |
+-------------------+  WS   +-----------------------+       +-------------------+
                                       ^
                                       | (Consumes Event)
                               +-------------------+
                               |   Event Bus       |
                               | (Local / Kafka)   |
                               +-------------------+
```

---

# 4. WhatsApp Integration (`whatsmeow`)

Kita menggunakan `go.mau.fi/whatsmeow` sebagai library inti.

### 4.1. Entity Model

- **`WhatsAppDevice`**: Menyimpan state sesi dari sebuah akun WhatsApp yang terhubung.
  - `id`, `organization_id`, `phone_number`, `session_data`, `status` (Connected, Disconnected).
- **`MessageLog`**: Menyimpan histori pesan keluar/masuk.
  - `id`, `organization_id`, `tenant_id`, `device_id`, `message_body`, `status` (Sent, Failed, Delivered).

### 4.2. QR Login Flow
1. User masuk ke halaman **Messaging Settings** di UI.
2. User menekan tombol "Add New Device".
3. Frontend menembak API `POST /api/v1/messaging/devices/qr`.
4. Backend `whatsmeow` membuat *session* baru dan mengembalikan gambar QR Code (dalam format Base64) ke Frontend.
5. User memindai QR Code melalui aplikasi WhatsApp di ponselnya.
6. Backend menerima callback keberhasilan login dari `whatsmeow`, menyimpan `session_data` ke database, dan meng-*update* status device menjadi `Connected`.

---

# 5. In-App Notification (WebSocket)

### 5.1. WebSocket Hub
- Backend mengelola sebuah `Hub` yang menyimpan daftar koneksi TCP/WebSocket aktif dari setiap pengguna (berdasarkan `user_id` atau `organization_id`).
- Setiap kali pengguna berhasil login ke aplikasi React, aplikasi akan membuka koneksi WebSocket (`ws://host/api/v1/ws`).

### 5.2. Event Driven Push
- Saat domain bisnis lain menerbitkan *Domain Event* (misal: `InvoicePaid` dari domain Finance), *Notification Service* mendengarkan event tersebut.
- Notification Service membangun struktur notifikasi, menyimpannya di tabel `notifications`, dan mem- *push* data JSON ke klien WebSocket yang bersangkutan.

---

# 6. Messaging Features

### 6.1. Blast Message (Broadcast)
- **Targeting**: Admin memilih filter (misal: Semua Penghuni di Gedung A, Semua Penghuni yang Menunggak).
- **Sender Selection**: Admin dapat memilih `WhatsAppDevice` mana yang akan menjadi pengirim pesan blast ini (berguna jika organisasi memiliki nomor khusus keuangan dan nomor khusus operasional).
- **Rate Limiting & Queue**: Backend *wajib* memasukkan proses *blast* ke dalam *Background Worker Queue* (misalnya menggunakan Asynq atau Redis PubSub) dengan jeda acak (delay) antar pesan, guna mencegah nomor WhatsApp terblokir (Anti-Spam).

### 6.2. Template Engine
- Pesan otomatis harus menggunakan template yang mendukung variabel dinamis (misal: `Halo {{tenant_name}}, tagihan Anda sebesar {{amount}}...`).

---

# 7. Proposed Frontend Pages

1. **Messaging Settings (`/dashboard/messaging/devices`)**
   - Menampilkan daftar perangkat WhatsApp yang terhubung beserta statusnya (Hijau=Online, Merah=Offline).
   - Modal QR Code untuk menghubungkan perangkat baru.
   - Pengaturan *Default Sender Device* untuk berbagai jenis event otomatis (Invoice, Maintenance, dll).

2. **Blast Message (`/dashboard/messaging/blast`)**
   - Editor pesan untuk *broadcasting*.
   - Pemilihan target audiens (Berdasarkan Property/Building/Status).
   - Pemilihan perangkat pengirim (Sender ID).
   - Histori *blast message* dan metrik keberhasilan kirim.

---

# 8. Proposed Backend Packages

- `internal/pkg/whatsapp/`: Abstraksi pembungkus (wrapper) untuk `whatsmeow`, menangani *connection pool* untuk *multiple devices*.
- `internal/pkg/websocket/`: Sistem *pub-sub* WebSocket untuk *in-app notification*.
- `internal/modules/communication/`: Modul bisnis yang mengelola rute API untuk CRUD Device, generate QR, Blast Message, dan menyimpan histori komunikasi.
