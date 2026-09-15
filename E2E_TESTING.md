# End-to-End (E2E) Testing Guide — EPMP

Dokumen ini menjelaskan standar, arsitektur, dan cara menjalankan **End-to-End (E2E) Testing** pada platform EPMP.

---

## 🌟 1. Kebijakan Kualitas (Mandatory Quality Gate)

> **ATURAN WAJIB PENGEMBANGAN:**  
> Setiap kali ada penambahan fitur baru, perubahan alur bisnis, atau refactoring komponen UI di frontend/backend, **wajib** menjalankan dan memperbarui pengujian End-to-End (E2E) sebelum kode dianggap selesai (*Definition of Done*).

### Standar Checklist Fitur Baru:
- [ ] Backend handler & route sudah terdaftar dan teruji (*unit/integration test*).
- [ ] Rute frontend terdaftar di `frontend/src/App.tsx`.
- [ ] Rute halaman ditambahkan ke daftar `PAGES_TO_TEST` di `frontend/run_e2e_tests.cjs`.
- [ ] Pengujian E2E dieksekusi dengan hasil **0 Errors Found** (`make test-e2e`).

---

## 🚀 2. Arsitektur & Teknologi E2E

Sistem E2E EPMP menggunakan **Playwright** dengan peramban **Google Chrome Resmi (Official System Build)**:

```
[Playwright Runner] ───(channel: 'chrome')───> [Google Chrome Sistem macOS/Linux]
                                                       │
                                   ┌───────────────────┴───────────────────┐
                                   ▼                                       ▼
                       [Frontend: React Vite]                  [Backend: Go Echo]
                        http://localhost:3000                  http://localhost:8080
                                                                           │
                                                                           ▼
                                                                  [PostgreSQL Database]
```

### Mengapa Menggunakan Chrome Sistem (`channel: 'chrome'`)?
1. **Zero Download Overhead:** Tidak perlu mendownload binary Chromium Playwright yang berukuran ratusan megabyte (`npx playwright install` bisa dilewati).
2. **Dukungan 3D / WebGL Native:** Komponen Spatial 3D (`@react-three/fiber`, Three.js, model GLB) berjalan di atas engine hardware acceleration yang sama persis dengan browser pengguna harian.
3. **Validasi Lingkungan Nyata:** Menguji rendering CSS, local storage, cookie, dan response browser komersial sesungguhnya.

---

## 💻 3. Cara Menjalankan E2E Test

### Prasyarat
Pastikan service backend dan frontend sudah aktif:
1. PostgreSQL di Docker aktif (`docker-compose up -d postgres redis`)
2. Backend berjalan di port `8080` (`cd backend && go run ./cmd/server/...`)
3. Frontend berjalan di port `3000` (`cd frontend && npm run dev`)

### Eksekusi Cepat (Dari Root Directory)

```bash
# 1. Mode Headless (Cepat, di latar belakang tanpa jendela browser)
make test-e2e
# atau
make e2e

# 2. Mode GUI / Headed (Membuka jendela browser Google Chrome visual secara real-time)
make test-e2e-gui
# atau
make e2e-gui
```

Atau bisa juga langsung dari direktori `frontend`:
```bash
cd frontend

# Mode Headless
npm run test:e2e

# Mode GUI (Headed Google Chrome dengan visual browser)
npm run test:e2e:gui
```

---

## 📋 4. Cakupan Pengujian (Test Coverage)

File pengujian utama berada di:
`frontend/run_e2e_tests.cjs`

Pengujian menguji alur berikut secara menyeluruh:

1. **Autentikasi & Sesi Enterprise:**
   - Login menggunakan akun tester ke `/auth/signin`.
   - Inisialisasi token JWT, refresh token, dan context organisasi (`X-Organization-ID`).
2. **Pemeriksaan 64 Halaman & Form:**
   - **Dashboard & Core:** Overview (`/dashboard`), Spatial Explorer (`/dashboard/explorer`).
   - **Struktur Properti:** Organizations, Properties, Buildings, Floors, Rooms, Room Types, Beds, Zones, Facilities (baik halaman List maupun New Form).
   - **Operasional:** Tenants, Reservations, Contracts, Occupancies.
   - **Finansial:** Invoices, Payments, Deposits, Charges, Refunds, Adjustments, Penalties.
   - **Aset & Pemeliharaan:** Assets, Asset Assignments, Asset Inspections, Work Orders, Technicians, Suppliers.
   - **Komunikasi:** Device Management, Blast Messages.
3. **Siklus Penuh CRUD (Create -> Read -> Update -> Delete):**
   - **Property CRUD:** Pembuatan properti baru dengan form lengkap, submit, verifikasi kemunculan instan pada tabel Properti.
   - **Building CRUD:** Pembuatan gedung baru, relasi dengan properti, navigasi form Basic View vs 3D, verifikasi daftar gedung.
   - **Tenant FULL CRUD:**
     - *Create*: Pendaftaran penyewa baru melalui form validasi Zod.
     - *Read / Detail*: Klik baris tabel dan verifikasi detail penyewa (`/dashboard/tenants/:id`).
     - *Update*: Pengeditan nama/kontak melalui halaman edit (`/dashboard/tenants/:id/edit`), simpan dan verifikasi perubahan data.
     - *Delete*: Penghapusan data dengan konfirmasi dialog otomatis, verifikasi data bersih dari database dan UI tabel.
   - **Asset CRUD:** Pendaftaran aset operasional (kategori elektronik, harga perolehan, status ketersediaan) ke basis data.
4. **Interaktivitas Khusus & 3D:**
   - **Building Creation:** Pergantian tab *Basic Form* vs *Interactive 3D View* (drop model `building.glb` di grid 3D).
   - **Property 3D Map:** Pemuatan scene Three.js digital twin drill-down gedung dan ruangan.
5. **Deteksi Error Otomatis:**
   - Setiap `pageerror` (unhandled exception / crash).
   - Setiap `console.error` di konsol browser.
   - Setiap HTTP response failure (status 4xx / 5xx).
   - Deteksi halaman kosong (*blank white page*).

---

## ➕ 5. Cara Menambahkan Rute Baru ke Test Suite

Jika Anda membuat modul atau halaman baru:
1. Buka file `frontend/run_e2e_tests.cjs`.
2. Tambahkan URL rute ke array `PAGES_TO_TEST`:
   ```javascript
   const PAGES_TO_TEST = [
     // ... rute yang sudah ada ...
     '/dashboard/nama-modul-baru',
     '/dashboard/nama-modul-baru/new',
   ];
   ```
3. Jika modul memiliki interaksi khusus (misal: klik tombol kalkulator, drag-and-drop, atau modal popup), tambahkan blok assertion di bawah bagian pengujian interaktif:
   ```javascript
   console.log('Testing Fitur Baru...');
   await page.goto(`${BASE_URL}/dashboard/nama-modul-baru`);
   await page.click('button:has-text("Aksi")');
   // Verifikasi hasil
   ```
4. Jalankan `make test-e2e` untuk memverifikasi.
