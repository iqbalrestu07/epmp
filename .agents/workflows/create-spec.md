---
description: Requirement gathering & spec update untuk request yang di luar scope saat ini
---

# /create-spec Workflow

Saat user menjalankan `/create-spec`, ikuti langkah berikut untuk mendraft dan mengintegrasikan permintaan fitur baru ke dokumen spesifikasi EPMP. **Jangan menulis kode aplikasi** selama workflow ini.

## Step 1: Gather Requirements

1. Minta user mendeskripsikan fitur, user stories, dan business rules secara lengkap. Tanyakan property model mana yang terpengaruh (boarding house, apartment, co-living, office, dll).
2. Tanyakan constraint UI/UX, API, dan apakah fitur bersifat per-organization (multi-tenant) atau global.
3. Stop dan **tunggu jawaban user**.

## Step 2: Analyze Impact

1. Dampak ke schema PostgreSQL: migration baru di `backend/migrations/` dan `MODULE.md` modul terkait.
2. Dampak ke backend: modul mana di `backend/internal/modules/` yang berubah / modul baru (lihat `epmp-module-patterns.md`). Cross-module dependency harus lewat service interface, bukan repository.
3. Dampak ke frontend: feature folder `frontend/src/features/`, routes di `frontend/src/App.tsx`, dan `PAGES_TO_TEST` di `frontend/run_e2e_tests.cjs`.
4. Dampak ke dokumen EPMP AI (`tools/epmp-ai/CONTEXT.md`, `RULES.md`) jika ada aturan bisnis baru.
5. Usulkan daftar Functional Requirements (FR) dan Non-Functional Requirements (NFR) dengan ID (`FR-xx`, `NFR-xx`).

## Step 3: Propose Plan

1. Buat `implementation_plan.md` yang merinci dokumen apa saja yang akan diubah dan bagaimana.
2. **Tunggu approval user.**

## Step 4: Update Documentation

Setelah disetujui:

1. Tulis/update spec di `epmp-docs/epmp-XXX-<slug>.md` (lanjutkan penomoran yang ada; lihat `epmp-docs/epmp-013-notification.md` sebagai contoh format).
2. Jika ada keputusan arsitektur, buat ADR di `epmp-docs/adr/` dengan **ADR Skill**.
3. Update `PROGRESS.md` bila fitur masuk roadmap.
4. Ingatkan user bahwa spec kini "In Scope" dan lanjut dengan `/orchestrator` atau `/2-implement`.
