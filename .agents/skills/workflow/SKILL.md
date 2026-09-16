---
name: workflow
description: Menjalankan salah satu workflow fase EPMP dari .agents/workflows (orchestrator, 1-research, 2-implement, 3-integrate, 4-verify, 5-commit, quick-fix, refactor, audit, create-spec). Gunakan saat user menyebut "/orchestrator", "/quick-fix", "jalankan workflow X", atau memulai fitur/perbaikan yang butuh fase terstruktur.
---

# Workflow Skill

Wrapper tipis agar workflow di `.agents/workflows/` bisa dipanggil dari Devin. Skill ini **tidak** berisi
langkah sendiri — isi otoritatif ada di file workflow.

## Cara pakai

1. Tentukan workflow dari permintaan user:

| Permintaan | File |
|---|---|
| Fitur baru end-to-end / `/orchestrator` | `.agents/workflows/orchestrator.md` (mengarahkan ke fase 1–5) |
| Riset / scoping saja | `.agents/workflows/1-research.md` |
| Implementasi (TDD) | `.agents/workflows/2-implement.md` |
| Test adapter / repository / migration | `.agents/workflows/3-integrate.md` |
| Verifikasi penuh sebelum ship | `.agents/workflows/4-verify.md` |
| Commit | `.agents/workflows/5-commit.md` |
| Bug kecil, ≤ 3 file, tanpa perubahan schema | `.agents/workflows/quick-fix.md` |
| Refactor tanpa perubahan perilaku | `.agents/workflows/refactor.md` |
| Audit kualitas/keamanan/dependency | `.agents/workflows/audit.md` |
| Request di luar scope spec | `.agents/workflows/create-spec.md` |

2. **Baca file workflow tersebut secara penuh**, lalu baca rule yang disebut sebagai *Mandatory Rules* di dalamnya (`.agents/rules/`).
3. Ikuti langkahnya secara berurutan. Untuk `orchestrator`, jangan lompat fase; setiap fase punya *Completion Criteria* yang harus terpenuhi sebelum lanjut.
4. Artefak ditulis ke lokasi EPMP: `epmp-docs/research-logs/`, `epmp-docs/audits/`, `epmp-docs/adr/`.

## Skill pendamping

- `guardrails` — pre-flight sebelum menulis kode & self-review sesudahnya
- `code-review`, `debugging-protocol`, `adr`, `patch-cve`, `sequential-thinking`
