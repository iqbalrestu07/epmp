---
name: patch-cve
description: Triage CVE findings dari audit.sh, review audit-log, dan patch package.json overrides untuk memperbaiki vulnerable dependencies. Gunakan setelah audit.sh gagal atau setelah /audit menemukan dependency vulnerabilities.
---

# Patch CVE Skill

## Purpose
Menjalankan security audit, membaca hasil `audit-log`, dan secara otomatis memperbaiki
vulnerable dependencies melalui `package.json` overrides — tanpa memerlukan upgrade Directus.

## When to Invoke
- Setelah `./audit.sh` gagal dengan exit code 1
- Setelah `/audit` workflow menemukan HIGH/CRITICAL CVE pada dependency
- Saat ada CVE baru dilaporkan pada library yang digunakan Directus

## When NOT to Use
- Ketika masalah bukan dependency (gunakan `/quick-fix` atau `/refactor`)
- Ketika CVE hanya ada di dev tooling lokal (bukan di Docker image)

---

## Steps

### Step 1: Jalankan audit.sh (atau baca audit-log yang ada)

Jika `audit-log` sudah ada dan fresh (dibuat < 30 menit), skip ke Step 2.
Jika tidak, jalankan:

```bash
./audit.sh
```

> Audit akan memakan waktu 5-15 menit karena melakukan `docker compose build`.
> Jika ingin skip build dan hanya scan image yang sudah ada, jalankan langsung:
> ```bash
> docker compose --profile audit run --rm trivy image \
>   --exit-code 0 --severity CRITICAL,HIGH,MEDIUM,LOW \
>   --scanners vuln --pkg-types library \
>   --format json --skip-version-check \
>   directus:latest 2>/dev/null | node -e "..." gate
> ```

### Step 2: Baca dan parse audit-log

Baca file `audit-log` di root project. Format yang diharapkan:

```
PACKAGE                       VERSION        SEVERITY    CVE                         FIX
-----------------------------------------------------------------------------------------------
basic-ftp                     5.2.2          HIGH        GHSA-rp42-5vxx-qpwr         5.3.0
liquidjs                      10.25.3        HIGH        CVE-2026-41311              10.25.7
...
```

Untuk setiap baris, ekstrak:
- `PACKAGE` — nama npm package
- `VERSION` — versi yang terinstall
- `SEVERITY` — CRITICAL / HIGH / MEDIUM / LOW
- `CVE` — ID CVE atau GHSA
- `FIX` — versi yang sudah di-patch (`no fix` jika belum ada)

### Step 3: Triage findings

Klasifikasikan setiap CVE:

| Kondisi | Aksi |
|---|---|
| Package adalah **build tool / dev dependency** | **HAPUS DARI DOCKERFILE** (Step 3.5) |
| `FIX` tersedia + SEVERITY CRITICAL/HIGH | **CEK DIRECTUS** dulu (Step 3.6), lalu PATCH jika perlu |
| `FIX = "no fix"` + SEVERITY CRITICAL/HIGH | **IGNORE** — tambah ke `.trivyignore` dengan justifikasi |
| `FIX` tersedia tapi **tidak bisa di-apply** + SEVERITY CRITICAL/HIGH | **IGNORE** — tambah ke `.trivyignore` dengan justifikasi (lihat catatan di bawah) |
| SEVERITY MEDIUM | **SKIP** — hanya dokumentasikan di findings, tidak ada aksi |
| SEVERITY LOW | **SKIP** — hanya dokumentasikan di findings, tidak ada aksi |

> **Alasan:** Security gate di `audit.sh` hanya memblokir HIGH/CRITICAL.
> MEDIUM dan LOW bersifat informatif dan tidak memblokir commit.

> **Kapan CVE dengan fix "tidak bisa di-apply"?**
> - **OS-level package** (Alpine apk) — Dockerfile sudah punya `apk upgrade --no-cache`
>   tapi base image belum publish versi terbaru. Ini di luar kontrol kita; tunggu
>   Alpine update dan re-build. Sementara itu, boleh di-ignore di `.trivyignore`.
> - **Major version jump tanpa safe path** — Fix mengharuskan upgrade major (misal
>   `tar@6→7`) dan override menyebabkan runtime error karena breaking API change
>   yang tidak kompatibel dengan consumer. Dokumentasikan alasan di `.trivyignore`.
> - **Package hanya ada di build-time** — Vulnerability hanya exists di dependency
>   yang digunakan saat build (dev dependency), bukan runtime. Attack vector tidak
>   applicable di production image.

### Step 3.5: Evaluasi Hapus Build Tool dari Runtime Image (Paling Direkomendasikan)

Sebelum melakukan override, cek apakah package yang terkena CVE hanyalah **build tool** atau **dev dependency** (misalnya `vite`, `esbuild`, `npm`, `npx`) yang tidak sengaja terbawa ke dalam *runtime image* produksi.

Jika package tersebut tidak dipanggil sama sekali saat aplikasi berjalan (*runtime*), maka **cara terbaik dan terbersih** adalah menghapusnya langsung dari `Dockerfile` di *stage* akhir (*runtime stage*).

**Cara eksekusi:**
Tambahkan baris berikut di `Dockerfile` pada *stage* runtime:
```dockerfile
# Hapus <nama-package> dari runtime (hanya dibutuhkan saat build-time, memicu CVE)
RUN find /directus/node_modules -name "<nama-package>" -type d -exec rm -rf {} + 2>/dev/null || true
```
Setelah itu, jalankan rebuild, commit dengan pesan `fix(security): remove <package> from runtime image to fix CVE...`. Jika berhasil, Anda bisa skip langkah selanjutnya.

---

### Step 3.6: Cek apakah Directus minor terbaru sudah memfix CVE

Sebelum menulis override, cek apakah upgrade `DIRECTUS_VERSION` di Dockerfile
sudah cukup untuk menghilangkan CVE tersebut.

**Cara cek:**

1. Baca versi Directus yang sedang dipakai dari `Dockerfile`:
   ```bash
   grep DIRECTUS_VERSION Dockerfile
   # Contoh: ENV DIRECTUS_VERSION=v11.17.0
   ```

2. Cari versi minor terbaru di GitHub releases (same major, latest minor):
   ```
   https://github.com/directus/directus/releases
   ```
   Atau via API:
   ```bash
   curl -s https://api.github.com/repos/directus/directus/releases/latest | grep '"tag_name"'
   ```

3. Untuk setiap package CRITICAL/HIGH yang ditemukan, cek apakah sudah di-bump
   di versi Directus terbaru dengan melihat `pnpm-lock.yaml` di branch tersebut:
   ```
   https://raw.githubusercontent.com/directus/directus/<latest-tag>/pnpm-lock.yaml
   ```
   Cari nama package dan bandingkan versi yang terpasang.

**Keputusan:**

| Hasil cek | Aksi |
|---|---|
| Directus terbaru sudah bump package ke versi ≥ fix version | **UPGRADE** `DIRECTUS_VERSION` di Dockerfile (lebih bersih dari override) |
| Directus terbaru belum fix (masih pakai versi vulnerable) | **OVERRIDE** — lanjut ke Step 4 |
| Versi terbaru adalah major baru (misal v11 → v12) | **SKIP upgrade** — jangan upgrade major; lanjut ke Step 4 |

> **Preferensi:** Upgrade Directus > override package.
> Override hanya sebagai fallback jika Directus upstream belum fix.
> Semakin sedikit override, semakin mudah maintenance jangka panjang.

**Jika upgrade Directus dipilih:**
```bash
# Update Dockerfile
# ENV DIRECTUS_VERSION=v11.17.0  →  ENV DIRECTUS_VERSION=v11.XX.0

# Commit terpisah dari override agar mudah di-revert
git add Dockerfile
git commit -m "chore(docker): upgrade Directus to vXX.XX.XX (fixes CVE-XXXX)"
```
Setelah itu jalankan `./audit.sh` untuk konfirmasi CVE hilang.
Jika masih ada CVE lain yang belum terfix, lanjut ke Step 4 untuk sisanya.

### Step 3.7: Audit override yang sudah tidak relevan (wajib setelah upgrade Directus)

Setiap kali `DIRECTUS_VERSION` di-bump, beberapa override di `package.json` mungkin
sudah **tidak diperlukan** karena Directus versi baru sudah menyertakan versi aman
secara native. Override yang tidak perlu harus dihapus — membiarkannya berisiko
memaksa Directus menggunakan versi yang tidak diuji upstream-nya.

**Cara cek setiap override:**

1. Baca semua entry di `pnpm.overrides` pada `package.json`
2. Untuk setiap override, cek versi yang dipakai Directus versi baru di `pnpm-lock.yaml`:
   ```
   https://raw.githubusercontent.com/directus/directus/<new-tag>/pnpm-lock.yaml
   ```
   Cari nama package dan catat versi resolusinya.

3. Bandingkan:

| Kondisi | Aksi |
|---|---|
| Directus baru sudah pakai versi ≥ override version | **HAPUS** override (sudah tidak diperlukan) |
| Directus baru masih pakai versi < override version | **PERTAHANKAN** override (masih melindungi) |
| Override terkait workaround build (bukan CVE, misal `rollup`, `tsconfig`) | **PERTAHANKAN** selalu — ini bukan CVE override |

> **Contoh:**
> Override `"basic-ftp": "5.3.0"` ditambahkan karena CVE-XXXX.
> Directus v11.18.0 sudah memakai `basic-ftp@5.3.1` secara native.
> → Hapus override `basic-ftp` dari `package.json`.

**Setelah cleanup:**
```bash
# Validasi JSON tetap valid
node -e "JSON.parse(require('fs').readFileSync('package.json','utf8')); console.log('valid')"

git add package.json
git commit -m "chore(deps): remove stale overrides after Directus upgrade to vXX.XX"
```

### Step 4: Terapkan patch di package.json


Baca `package.json` dan update bagian `pnpm.overrides`:

```json
{
  "pnpm": {
    "overrides": {
      "basic-ftp": "5.3.0",
      "liquidjs": "10.25.7"
    }
  }
}
```

**Rules:**
- Selalu set ke versi **minimum fixed version** dari kolom FIX di audit-log
- Jika package sudah ada di overrides, **bump** ke versi fix (jangan downgrade)
- Jika package punya multiple CVE dengan fix version berbeda, gunakan **versi tertinggi**
- Simpan komentar yang ada di package.json
- Jangan hapus override yang tidak terkait CVE ini

### Step 5: Untuk CVE tanpa fix — tambah ke .trivyignore

Format `.trivyignore`:

```
# CVE-XXXX-XXXXX — <package>@<version>
# Justifikasi: <alasan mengapa ini acceptable>
# Ditambahkan: <tanggal>
CVE-XXXX-XXXXX
```

**Panduan justifikasi yang valid:**
- Package hanya digunakan di build-time, tidak di runtime
- Vector attack tidak applicable (misal: vulnerability memerlukan akses lokal)
- Package sudah akan di-remove di Directus versi berikutnya (link issue)

**JANGAN tambahkan ke .trivyignore jika tidak ada justifikasi valid.**

### Step 6: Validasi

Jalankan syntax check terlebih dahulu:

```bash
node --check extensions/survey-forms/index.js
```

Kemudian verifikasi `package.json` valid JSON:

```bash
node -e "JSON.parse(require('fs').readFileSync('package.json','utf8')); console.log('valid')"
```

> **Catatan:** Jalankan `./audit.sh` ulang hanya jika user meminta verifikasi penuh.
> Build ulang memakan waktu lama. Cukup validasi JSON dan syntax.

### Step 7: Ringkasan perubahan

Setelah semua perubahan dibuat, tampilkan ringkasan:

```
## CVE Patch Summary

### Patched via package.json overrides
| Package | Dari | Ke | CVE |
|---|---|---|---|
| basic-ftp | 5.2.2 | 5.3.0 | GHSA-rp42-5vxx-qpwr |
| liquidjs | 10.25.3 | 10.25.7 | CVE-2026-41311 |

### Added to .trivyignore
| CVE | Package | Justifikasi |
|---|---|---|
| CVE-XXXX | pm2 | No fix available; build-only tool |

### Skipped (no fix, LOW severity)
- CVE-YYYY: nodemailer@7.0.11 — no fix available
```

### Step 8: Commit

Gunakan conventional commit:

```bash
git add package.json .trivyignore
git commit -m "fix(deps): patch HIGH/CRITICAL CVEs via pnpm overrides

- basic-ftp: 5.2.2 → 5.3.0 (GHSA-rp42-5vxx-qpwr)
- liquidjs: 10.25.3 → 10.25.7 (CVE-2026-41311)"
```

---

## Important Constraints

1. **Jangan upgrade ke major version berbeda** tanpa memahami breaking changes.
   - `uuid@14.0.0` adalah major upgrade — check jika Directus kompatibel sebelum override.
   - Prefer minimum fix version (`5.3.0`) daripada latest (`6.0.0`).

2. **Overrides bersifat transitive** — `"basic-ftp": "5.3.0"` akan mempengaruhi semua
   packages yang depend pada `basic-ftp`, termasuk versi nested.

3. **Setelah rebuild, re-run `./audit.sh`** untuk konfirmasi CVE sudah hilang.
   Beberapa CVE mungkin masih muncul jika ada multiple paths ke library yang sama.

4. **jangan hapus override yang sudah ada** kecuali ada alasan eksplisit —
   override lama mungkin masih memproteksi dari CVE lain.
