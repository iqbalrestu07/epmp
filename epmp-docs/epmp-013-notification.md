# EPMP-013

# Notification & Messaging Service Architecture

```text
Document ID    : EPMP-013
Document Name  : Notification & Messaging Service Architecture
Version        : 1.1.0
Status         : Implemented & Active
Owner          : Software Architecture & Engineering Team
Dependencies   : EPMP-003, EPMP-006, EPMP-007
Referenced By  : Communication Module, Frontend Messaging Settings, Blast Messaging
```

---

# 1. Purpose

Dokumen ini mendefinisikan arsitektur dan spesifikasi untuk **Notification & Messaging Service** dalam ekosistem Enterprise Property Management Platform (EPMP). Layanan ini bertanggung jawab atas penyampaian informasi real-time kepada pengguna internal (Web Dashboard) maupun pengguna eksternal (Tenant/Penghuni).

Sistem komunikasi ini dibangun dengan kapabilitas multi-channel:
- **WhatsApp Multi-Device Gateway**: Integrasi native Go menggunakan pustaka `go.mau.fi/whatsmeow` untuk komunikasi langsung dua arah dan broadcast ke penghuni.
- **In-App Real-time Notification**: Menggunakan *WebSocket* untuk *real-time push notification* ke pengguna aplikasi.

---

# 2. Key Capabilities

1. **WhatsApp Gateway Mandiri (Whatsmeow Native Go)**
   - Tidak memerlukan service pihak ketiga berbayar atau server Node.js eksternal.
   - Mendukung banyak nomor WhatsApp (*multi-device*) per organisasi (misal: satu untuk Divisi Keuangan, satu untuk Operasional/CS).
   - Pengguna cukup men-scan QR code resmi dari WhatsApp melalui dashboard EPMP.

2. **Deteksi Otomatis Nomor HP via Callback Scan QR**
   - User **tidak perlu menginput nomor handphone** secara manual saat menambah perangkat.
   - Nomor HP langsung dideteksi secara otomatis dari handshake callback server WhatsApp (`client.Store.ID.User`) begitu scan QR berhasil.

3. **Verifikasi Nomor WhatsApp Terdaftar (`IsOnWhatsApp`) Sebelum Kirim**
   - Untuk mencegah kegagalan massal dan proteksi nomor dari pemblokiran (rate limiting/spam penalty), sistem **memeriksa apakah nomor tujuan terdaftar di WhatsApp** terlebih dahulu via `client.IsOnWhatsApp()`.
   - Nomor berformat lokal (`08...`) otomatis dinormalisasi ke format internasional (`628...`).
   - Hanya nomor yang terbukti aktif di WhatsApp yang akan dikirimkan pesan ke JID terverifikasinya. Nomor yang tidak terdaftar akan dicatat di log pengiriman dengan status `failed` dan keterangan yang transparan.

4. **Blast Message (Broadcast) dengan Dynamic Templates**
   - Admin dapat mengirim pesan massal berdasarkan segmentasi audiens:
     - **Semua Tenant Aktif** (`all_tenants`)
     - **Tenant dengan Tunggakan / Belum Lunas** (`overdue`)
     - **Manual Selection** (`manual`)
   - Mendukung variabel dinamis seperti `{{tenant_name}}`, `{{room_name}}`, `{{invoice_amount}}`, `{{rent_amount}}`, `{{due_date}}`, dan `{{days_left}}`.
   - **Strict Real Delivery**: Sistem mewajibkan adanya perangkat pengirim yang terhubung. Tidak ada *mock simulation* terselubung.

5. **Penyimpanan Sesi & Auto-Reconnect**
   - Kunci enkripsi sesi disimpan di database PostgreSQL (`whatsmeow_device` dan `wa_devices.session_data`).
   - Ketika server backend restart, seluruh sesi perangkat yang berstatus *connected* akan **otomatis terhubung kembali** (*auto-reconnect*) tanpa mewajibkan scan QR ulang.

---

# 3. Arsitektur Teknis

```text
+-------------------------------------------------------------+
|                      React Web Dashboard                    |
|                                                             |
|  [Pengaturan Gateway]                [Blast Message (WA)]   |
|  - Render Live QR (react-qr-code)    - Filter Target Audiens|
|  - Real-time Polling Status          - Dynamic Template Vars|
|  - Connect / Disconnect / Delete     - WA Phone Mockup      |
+-------------------------------------------------------------+
                                │ REST API
                                ▼
+-------------------------------------------------------------+
|                      Go Backend (Echo)                      |
|                                                             |
|  internal/modules/communication/                            |
|  ├── module.go        (REST Handler: devices, blast, tmpl)  |
|  └── wa_manager.go    (Whatsmeow Manager, Lifecycle,        |
|                        IsOnWhatsApp, SendMessage, Reconnect)|
+-------------------------------------------------------------+
               │                                │
      PostgreSQL (Database)              WebSocket (WSS)
   - wa_devices                          wss://web.whatsapp.com
   - blast_messages                             │
   - blast_message_logs                         ▼
   - message_templates                 +----------------------+
   - whatsmeow_device (Keys)           |   WhatsApp Servers   |
                                       +----------------------+
```

---

# 4. WhatsApp Multi-Device Lifecycle

### 4.1. Alur Penambahan & Pairing Device via QR Code
1. User membuka halaman **Pengaturan Gateway** (`/dashboard/messaging/devices`) dan menekan tombol **Tambah Perangkat**.
2. User hanya mengisi **Label Perangkat** (contoh: *Divisi Keuangan*). Nomor HP tidak diminta.
3. Backend membuat record baru di tabel `wa_devices` dengan status `disconnected`.
4. User mengklik tombol **Hubungkan** pada perangkat:
   - Backend memanggil `m.waMgr.GetQR(ctx, deviceID)`.
   - `whatsmeow` membuat *device store* baru di PostgreSQL, meminta QR channel via `client.GetQRChannel(ctx)`, dan melakukan dial WebSocket ke server WhatsApp.
   - WhatsApp mengembalikan kode pairing resmi berupa raw string (contoh: `https://wa.me/settings/linked_devices#2@...`).
   - Frontend menerima raw string tersebut dan merendernya sebagai QR Code asli via komponen `react-qr-code`.
5. User membuka WhatsApp di HP (**Perangkat Tertaut** → **Tautkan Perangkat**) dan memindai QR Code di layar.
6. Server WhatsApp menyelesaikan *Noise Handshake*, lalu memancarkan event `success` ke goroutine backend:
   - Backend membaca `client.Store.ID.User` (nomor HP terhubung).
   - Backend memperbarui tabel `wa_devices`: `status = 'connected'`, `phone = <nomor_hp>`, `session_data = <jid>`.
7. Frontend yang melakukan polling interval ringan (2 detik) langsung mendeteksi perubahan status dan menampilkan layar sukses **"WhatsApp Berhasil Terhubung!"** lengkap dengan nomor HP yang baru dipasangkan.

### 4.2. Alur Pengiriman Pesan & Verifikasi Nomor
```text
[Tenant Phone: "085157943349"]
               │
               ▼
   [Normalisasi Nomor: "6285157943349"]
               │
               ▼
   [client.IsOnWhatsApp(ctx, ["+6285157943349"])]
               │
      ┌────────┴────────┐
      │                 │
[IsIn == true]    [IsIn == false]
      │                 │
      ▼                 ▼
[Kirim ke JID]    [Skip Pengiriman]
      │                 │
      ▼                 ▼
[Log: "sent"]     [Log: "failed" ("Nomor tidak terdaftar di WhatsApp")]
```

---

# 5. Spesifikasi Database (PostgreSQL)

### 5.1. `wa_devices`
Menyimpan metadata perangkat WhatsApp gateway yang didaftarkan.
```sql
CREATE TABLE wa_devices (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id       UUID NOT NULL,
    label        VARCHAR(100) NOT NULL,
    phone        VARCHAR(30),
    status       VARCHAR(20) NOT NULL DEFAULT 'disconnected', -- connected | disconnected | qr_pending
    session_data TEXT,                                       -- WhatsApp JID string untuk reconnect
    last_seen    TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 5.2. `message_templates`
Koleksi template pesan broadcast dengan parameter dinamis.
```sql
CREATE TABLE message_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL,
    name        VARCHAR(100) NOT NULL,
    content     TEXT NOT NULL,
    variables   TEXT[],                                       -- Array variabel e.g. {tenant_name, invoice_amount}
    category    VARCHAR(50) DEFAULT 'general',                -- general | invoice | reminder | announcement
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 5.3. `blast_messages`
Catatan kampanye pengiriman pesan massal.
```sql
CREATE TABLE blast_messages (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id           UUID NOT NULL,
    device_id        UUID REFERENCES wa_devices(id) ON DELETE SET NULL,
    title            VARCHAR(200),
    template         TEXT NOT NULL,
    target_type      VARCHAR(50) NOT NULL DEFAULT 'all_tenants', -- all_tenants | overdue | manual
    target_filter    JSONB,
    status           VARCHAR(20) NOT NULL DEFAULT 'draft',       -- draft | sending | done | failed
    total_recipients INT NOT NULL DEFAULT 0,
    sent_count       INT NOT NULL DEFAULT 0,
    failed_count     INT NOT NULL DEFAULT 0,
    scheduled_at     TIMESTAMPTZ,
    started_at       TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ,
    created_by       UUID,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 5.4. `blast_message_logs`
Catatan transmisi per penerima individual untuk transparansi dan audit.
```sql
CREATE TABLE blast_message_logs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blast_id       UUID NOT NULL REFERENCES blast_messages(id) ON DELETE CASCADE,
    tenant_id      UUID,
    phone          VARCHAR(30) NOT NULL,
    recipient_name VARCHAR(200),
    message        TEXT NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending | sent | failed | read
    error_message  TEXT,                                  -- Alasan kegagalan (misal: tidak terdaftar di WA)
    sent_at        TIMESTAMPTZ,
    read_at        TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

# 6. Daftar Endpoint API

Semua rute berada di bawah proteksi autentikasi JWT (`/api/v1/communication`):

| Method | Path | Deskripsi |
|---|---|---|
| `GET` | `/communication/devices` | Mendapatkan daftar perangkat WA |
| `POST` | `/communication/devices` | Mendaftarkan perangkat baru (hanya label) |
| `GET` | `/communication/devices/:id/qr` | Meminta string QR pairing WhatsApp aktif |
| `PUT` | `/communication/devices/:id/status`| Memperbarui status / memutuskan koneksi |
| `DELETE`| `/communication/devices/:id` | Menghapus perangkat dan sesi whatsmeow |
| `GET` | `/communication/templates` | Mendapatkan daftar template pesan aktif |
| `POST` | `/communication/templates` | Membuat template pesan baru |
| `PUT` | `/communication/templates/:id` | Mengedit template pesan |
| `DELETE`| `/communication/templates/:id` | Menghapus (soft delete) template pesan |
| `GET` | `/communication/blast` | Mendapatkan riwayat blast messages |
| `POST` | `/communication/blast` | Membuat draft / kampanye blast baru |
| `GET` | `/communication/blast/:id` | Mendapatkan detail status blast |
| `GET` | `/communication/blast/:id/logs` | Mendapatkan log detail per penerima |
| `POST` | `/communication/blast/:id/send` | Menjalankan pengiriman pesan blast via WhatsApp |
| `GET` | `/communication/recipients` | Preview daftar penerima berdasarkan target type |

---

# 7. Frontend Features

1. **Halaman Pengaturan Gateway (`/dashboard/messaging/devices`)**:
   - Tab **Perangkat WA**: List perangkat, indikator online/offline, tombol Hubungkan (QR Modal), Putuskan, Hapus.
   - Modal QR: Rendering QR otomatis via `react-qr-code`, countdown timer kadaluarsa 60 detik, tombol perbarui QR, dan auto-detection saat HP selesai scan.
   - Tab **Template Pesan**: Manajemen template lengkap dengan tag pills variabel (`{{tenant_name}}`, dll.) yang bisa diklik untuk disisipkan.

2. **Halaman Blast Message (`/dashboard/messaging/blast`)**:
   - Tab **Compose**:
     - Auto-select perangkat WhatsApp yang sedang terhubung.
     - Pilihan target penerima dengan live counter dan daftar kontak penerima.
     - Dropdown template & variable pills.
     - Tampilan simulasi mockup pesan WhatsApp di sisi kanan secara real-time.
     - Tombol "Kirim Sekarang" otomatis *disabled* jika belum ada perangkat WA yang aktif terhubung.
   - Tab **Riwayat**:
     - Tabel kampanye blast lengkap dengan status badge (*Selesai*, *Mengirim*, *Gagal*), metrik jumlah berhasil/gagal.
     - Modal detail berisi log per kontak beserta alasan kegagalan jika nomor tidak terdaftar.
