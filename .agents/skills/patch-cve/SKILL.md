---
name: patch-cve
description: Triage dan patch vulnerable dependencies EPMP — Go modules (govulncheck), npm (npm audit), dan Docker base image. Gunakan setelah /audit menemukan HIGH/CRITICAL CVE atau saat ada advisory baru pada library yang dipakai backend/frontend.
---

# Patch CVE Skill

## Purpose

Menjalankan dependency audit, membaca hasilnya, dan memperbaiki vulnerable dependencies
di tiga permukaan EPMP:

| Surface        | Manifest                                              | Scanner                                                |
| -------------- | ----------------------------------------------------- | ------------------------------------------------------ |
| Go backend     | `backend/go.mod`, `backend/go.sum`                    | `govulncheck`                                          |
| React frontend | `frontend/package.json`, `frontend/package-lock.json` | `npm audit`                                            |
| Docker images  | `backend/Dockerfile`, `frontend/Dockerfile`           | base image tag review (+ `docker scout` bila tersedia) |

## When to Invoke

- Setelah `/audit` workflow menemukan HIGH/CRITICAL CVE pada dependency
- Saat ada advisory baru pada library yang dipakai (Echo, pgx, jwt, whatsmeow, React, Vite, dll)
- Sebelum release / merge ke `main` jika audit belum dijalankan > 30 hari

## When NOT to Use

- Masalah bukan dependency (gunakan `/quick-fix` atau `/refactor`)
- CVE hanya ada di dev tooling lokal yang tidak masuk image produksi (dokumentasikan, tidak perlu patch)

---

## Steps

### Step 1: Scan

```bash
# Go — hanya melaporkan vuln yang benar-benar reachable dari kode kita
cd backend && go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# npm — production deps saja (dev deps tidak masuk nginx runtime image)
cd frontend && npm audit --omit=dev --json > /tmp/npm-audit.json; npm audit --omit=dev
```

Cek base image:

```bash
grep -n '^FROM' backend/Dockerfile frontend/Dockerfile
```

> `npm audit` tanpa `--omit=dev` boleh dijalankan sebagai informasi, tetapi hanya
> production findings yang memblokir.

### Step 2: Triage

Klasifikasikan setiap temuan:

| Kondisi                                                                       | Aksi                                                                                                              |
| ----------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| Go vuln **reachable** (govulncheck melaporkan call stack) + fix tersedia      | **PATCH** — `go get <module>@<fixed>` (Step 3a)                                                                   |
| Go vuln di module yang di-import tapi symbol tidak dipanggil                  | **DOKUMENTASIKAN** — govulncheck sudah memfilternya; catat di audit log saja                                      |
| npm prod dep, CRITICAL/HIGH, fix tersedia tanpa major bump                    | **PATCH** — `npm audit fix` atau bump versi (Step 3b)                                                             |
| npm prod dep, fix hanya lewat major bump                                      | **OVERRIDE** transitive via `overrides` di `package.json` jika aman; kalau tidak, ADR + ignore dengan justifikasi |
| npm dev-only dep (eslint, vite plugin, playwright, dll)                       | **SKIP** — tidak ada di runtime image; catat                                                                      |
| `FIX = none` + CRITICAL/HIGH                                                  | **IGNORE dengan justifikasi** — catat di `epmp-docs/audits/` (Step 4)                                             |
| SEVERITY MEDIUM/LOW                                                           | **SKIP** — dokumentasikan saja                                                                                    |
| Base image (`golang:`, `alpine:`, `node:`, `nginx:`) tertinggal patch release | **BUMP** tag ke patch terbaru dalam minor yang sama                                                               |

### Step 3a: Patch Go module

```bash
cd backend
go get github.com/<owner>/<module>@v<fixed>
go mod tidy
go build ./... && go test ./... -count=1
```

**Rules:**

- Pakai **minimum fixed version**, bukan `@latest`, kecuali `latest` adalah patch release yang sama minor-nya
- Untuk indirect dependency, `go get` tetap berlaku — Go akan menaikkan versi di `go.mod` sebagai `// indirect`
- Jangan bump major (`/v2` → `/v3`) dalam skill ini — itu pekerjaan `/refactor` dengan ADR
- Jalankan ulang `govulncheck` untuk konfirmasi

### Step 3b: Patch npm dependency

```bash
cd frontend
npm audit fix --omit=dev          # non-breaking saja; JANGAN pakai --force
npm run lint && npm run build     # tsc -b + vite build harus tetap hijau
```

Jika `npm audit fix` tidak bisa (fix ada di transitive dep yang di-pin parent-nya), tambahkan
override **minimal**:

```json
{
  "overrides": {
    "<package>": "<minimum-fixed-version>"
  }
}
```

**Rules:**

- Set ke **minimum fixed version**; jika ada beberapa CVE, gunakan versi tertinggi di antara fix-nya
- Jangan hapus override yang sudah ada tanpa alasan eksplisit
- Setelah override, `npm install` agar `package-lock.json` ikut berubah, lalu `npm run build`
- Pilih versi yang sudah dipublikasikan ≥ 7 hari (hindari versi baru yang belum ter-vetting)

### Step 3c: Bump Docker base image

Ubah tag `FROM` ke patch release terbaru di minor yang sama (mis. `golang:1.26.0-alpine` → `golang:1.26.1-alpine`).
Lalu:

```bash
docker compose build backend frontend
```

Go version di `backend/Dockerfile` harus konsisten dengan direktif `go` di `backend/go.mod`.

### Step 4: Dokumentasikan yang di-ignore

Untuk CVE tanpa fix atau dengan fix yang tidak bisa diterapkan, tulis di
`epmp-docs/audits/{YYYY-MM-DD}-dependency-cve.md`:

```
| CVE | Package@version | Severity | Keputusan | Justifikasi | Review ulang |
|---|---|---|---|---|---|
| GO-2026-XXXX | golang.org/x/net@v0.x | HIGH | IGNORE | symbol tidak reachable (govulncheck) | 2026-XX-XX |
```

**Justifikasi valid:** tidak reachable / build-time only / attack vector memerlukan kondisi
yang tidak ada di deployment kita / menunggu upstream (sertakan link issue).
**JANGAN ignore tanpa justifikasi.**

### Step 5: Validasi

```bash
cd backend && go build ./... && go test ./... -count=1
cd frontend && npm run lint && npm run build
```

Jika integration test skip karena DB tidak jalan, nyalakan `docker compose up -d postgres`
dan ulangi — skip bukan pass.

### Step 6: Ringkasan

```
## CVE Patch Summary

### Go modules
| Module | Dari | Ke | Advisory |
|---|---|---|---|

### npm (production)
| Package | Dari | Ke | Advisory | Cara (audit fix / override) |
|---|---|---|---|---|

### Docker base images
| File | Dari | Ke |
|---|---|---|

### Ignored (dengan justifikasi → epmp-docs/audits/...)
### Skipped (dev-only / MEDIUM / LOW)
```

### Step 7: Commit

Satu commit per surface agar mudah di-revert:

```bash
git add backend/go.mod backend/go.sum
git commit -m "fix(deps): bump <module> to v<fixed> (GO-2026-XXXX)"

git add frontend/package.json frontend/package-lock.json
git commit -m "fix(deps): patch <package> HIGH CVE via npm audit fix"

git add backend/Dockerfile frontend/Dockerfile
git commit -m "chore(docker): bump base images to latest patch release"
```

---

## Important Constraints

1. **Jangan bump major version** dalam skill ini. Major bump = `/refactor` + ADR.
2. **`npm audit fix --force` dilarang** — bisa mem-bump major secara diam-diam.
3. **Override bersifat transitive** — satu override mempengaruhi semua consumer package tsb.
4. **Lockfile harus ikut di-commit** (`go.sum`, `package-lock.json`); jangan pernah edit lockfile manual.
5. **Selalu re-scan** setelah patch untuk konfirmasi advisory hilang.
