---
trigger: always_on
description: Setup awal untuk anggota tim baru yang baru join ke proyek ini
---

# Onboarding — Setup Lingkungan Pengembangan

Panduan ini wajib dibaca dan diikuti oleh setiap anggota tim baru sebelum mulai berkontribusi ke repositori `kes-ref-screening-directus`.

## Prasyarat

- [Docker Desktop](https://www.docker.com/) sudah terpasang dan berjalan
- Node.js >= 18 sudah terpasang
- Akses ke repositori GitLab sudah diberikan oleh lead / admin tim

## Urutan Setup (Jalankan Sekali)

```bash
# 1. Clone repositori
git clone <gitlab-repo-url>
cd kes-ref-screening-directus

# 2. Salin file environment
cp .env.example .env
# Isi nilai credential di .env (minta dari lead tim):
# DIRECTUS_URL, ADMIN_EMAIL, ADMIN_PASSWORD

# 3. Jalankan Directus via Docker
docker compose up -d

# 4. Jalankan migration skema (tunggu sampai Directus siap ~30 detik)
node scripts/run-migrations.js

# 5. Seeding data sampel form (opsional, untuk development)
node --env-file=.env scripts/seed-samples.js

# 6. Perbaiki tampilan UI relasi di panel admin
node --env-file=.env scripts/fix-display-templates.js
```

Setelah selesai, buka: **http://localhost:8055**

## Konvensi Penting

| Aspek | Aturan |
|---|---|
| **Zero Data-Loss** | Script deployment dilarang keras melakukan DROP collection atau DELETE/TRUNCATE data tanpa otorisasi eksplisit (`--allow-delete`). Partial deploy wajib steril & aditif. |
| **Schema DB** | Jangan ubah lewat GUI Directus Admin. Gunakan `scripts/migrations/` |
| **Kolom `code`** | Field wajib di tabel `forms` — slug unik untuk identifikasi form |
| **Script CLI** | Selalu gunakan `node --env-file=.env scripts/<nama>.js` |
| **Commit format** | Ikuti Conventional Commits: `feat(scope): description` |
| **Branch kerja** | Buat dari `development`: `git checkout -b feature/nama-fitur development` |
| **Alur merge** | `feature/*` → `development` → `staging` → `master` |

## Referensi Dokumen

| Dokumen | Lokasi | Keterangan |
|---|---|---|
| PRD | `docs/PRD/` | Kebutuhan fungsional & bisnis |
| TSD | `docs/TSD/` | Spesifikasi teknis & arsitektur |
| Schema DB | `docs/schema.md` | Semua field & relasi koleksi Directus |
| Audit Log | `docs/audits/` | Hasil review kode |

## Cara Menggunakan AI Agent

AI Coding Assistant yang Anda gunakan sudah dikonfigurasi dengan konteks penuh proyek ini. Gunakan slash commands untuk menugaskan pekerjaan:

```
@[/1-research]   ← Planning sebelum coding
@[/2-implement]  ← TDD cycle
@[/4-verify]     ← Validasi sebelum commit
@[/5-commit]     ← Commit dengan format standar
@[/audit]        ← Review kode
@[/quick-fix]    ← Hotfix cepat
@[/create-spec]  ← Dokumentasi fitur baru (out-of-scope)
```

> ⚠️ Format yang benar: `@[/nama-command]` atau `/nama-command`
> ❌ Jangan tulis hanya `nama-command` tanpa `/` — AI tidak akan mengenalinya.

## Pertanyaan Umum

**Q: Script migration gagal dengan error "collection already exists"?**
A: Normal — migration runner bersifat idempotent dan akan skip collection yang sudah ada.

**Q: Endpoint `/survey-forms` tidak bisa diakses?**
A: Pastikan Docker sudah jalan (`docker compose ps`) dan extension sudah ter-load. Cek logs: `docker compose logs -f directus`.

**Q: Ingin menambahkan field baru ke schema?**
A: Jangan lewat GUI. Buat migration baru di `scripts/migrations/` dengan nama `002-nama-fitur.js`, lalu jalankan `node scripts/run-migrations.js`.
