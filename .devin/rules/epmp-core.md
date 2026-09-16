# EPMP core rules (Devin)

Rule lengkap ada di `.agents/rules/` dan `AGENTS.md`; file ini hanya menegaskan yang paling sering dilanggar.

## Sebelum edit

- Baca `MODULE.md` modul yang disentuh, atau `backend/internal/modules/property/MODULE.md` sebagai referensi bila membuat modul baru.
- Modul/feature baru: pertimbangkan generator `tools/epmp-sdk` (`be/codegen`, `fe/codegen`) dari `schemas/*.yaml` sebelum menulis manual.
- Cek `frontend/src/services/api.ts` sebelum menambah call API — jangan buat fetch wrapper baru.

## Saat edit

- Repository interface di `repository/`, implementasi pgx di `repository/*_impl.go`; handler Echo tidak pernah memegang `*pgxpool.Pool`.
- Error lewat `internal/pkg/errs` + `internal/pkg/response`; jangan `echo.NewHTTPError` langsung di service.
- Setiap SQL memakai `organization_id` dari context. Tidak ada `SELECT *`, tidak ada string concat SQL.
- Migration: `cd backend && go run ./cmd/migrate create <name>` lalu isi `.up.sql` **dan** `.down.sql`.
- Frontend: Zod schema di `schema/`, TanStack Query hooks di `hooks/`, tidak ada `any`, tidak ada `fetch` langsung di komponen.
- WhatsApp (`internal/pkg/whatsapp`, modul `communication`): di test selalu di-fake; jangan pernah kirim ke nomor asli.

## Sebelum selesai

- `cd backend && gofmt -l . && go vet ./... && go test ./... -count=1` hijau, dan integration test **tidak** ter-skip (Postgres harus nyala).
- `cd frontend && npm run lint && npm run build` hijau.
- UI berubah → route ada di `PAGES_TO_TEST` dan `make test-e2e` berakhir `0 Errors Found`.
- Tidak ada secret, tidak ada `console.log`/`fmt.Println` debug tersisa, `MODULE.md` diperbarui bila kontrak modul berubah.
- Jangan commit/push kecuali diminta.
