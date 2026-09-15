

## Data Source Boundary — `ref_*` vs `sspkg_*`

### Core Principle

**Revamped endpoints MUST source data exclusively from `ref_*` collections. Never use `sspkg_*` as a data source for API responses.**

### Collection Classification

| Prefix | Purpose | Usage | Mutable? |
|---|---|---|---|
| `ref_*` | **Revamped reference data** — denormalized, simplified, with Directus relations | ✅ Data source for all revamped endpoints | Yes (via migration scripts) |
| `sspkg_*` | **Staging snapshot** — as-is mirror of Odoo tables, reference only | ❌ NEVER use as data source for API responses | No (read-only snapshot) |

### Data Pipeline

```
  Odoo (Source of Truth)
      │
      │ Snapshot Script (as-is, no transformation)
      ▼
  sspkg_* (Staging — ONLY for reference & seeder input)
      │
      │ Seeder Script (transform + denormalize)
      ▼
  ref_* (Production — used by endpoints, builder, fetcher)
      │
      │ Application Code (fetcher → builder)
      ▼
  API Response (SurveyJS JSON, etc.)
```

### `sspkg_*` Schema Rules — As-Is Mirror

**`sspkg_*` harus merupakan mirror persis dari tabel Odoo:**

1. **Field names as-is** — gunakan nama field yang sama dengan Odoo model (e.g. `prosedur_tatalaksana_id`, `product_template_id`, `kfa_code`)
2. **Flat tables** — semua kolom ID (`prosedur_tatalaksana_id`, `diagnosis_id`, dll.) disimpan sebagai **integer biasa** di Directus, **BUKAN** M2O relation
3. **No Directus relations** — tidak ada foreign key constraint di Directus. Relasi antar sspkg_* hanya referensi logis
4. **Stored related fields boleh disertakan** — jika Odoo menyimpan related field (e.g. `kfa_code` = related dari `product_template_id.kfa_code`), boleh di-snapshot juga
5. **Subset diizinkan** — untuk tabel master yang sangat besar (e.g. `product.product`), boleh snapshot hanya subset yang direferensi
6. **Tambah `snapshot_at` timestamp** — satu-satunya field tambahan yang tidak ada di Odoo
7. **Tidak ada transformasi** — data masuk apa adanya dari Odoo REST API

> **Mengapa flat?** Agar terlihat sederhana di Directus Admin sebagai tabel referensi. Admin tidak perlu navigate relasi Odoo yang kompleks.

### `ref_*` Schema Rules — Revamp/Simplifikasi

**`ref_*` adalah hasil pengembangan/simplifikasi dari `sspkg_*`:**

1. **Denormalized** — gabungkan data dari beberapa `sspkg_*` menjadi 1 tabel yang self-contained. Contoh: `ref_tatalaksana_prescription` = gabungan `sspkg_prescription_dosage_line` + `sspkg_product_product`
2. **Dengan Directus relations** — gunakan M2O FK, proper interfaces, editable di Directus Admin
3. **UUID primary key** — mengikuti standar Directus (bukan integer Odoo ID)
4. **`code` field unique** — setiap `ref_*` harus punya business identifier `code` yang unique
5. **`status` field** — `published` / `draft` / `archived` (standar Directus)
6. **Admin-friendly** — admin bisa langsung edit data di Directus GUI tanpa harus paham struktur Odoo

### Why This Matters

1. **`sspkg_*` is staging** — it's a frozen copy of Odoo data. It exists solely as input for seeder scripts that populate `ref_*`.

2. **`ref_*` is the source of truth** — all revamped, normalized data lives here. Schema is designed for SurveyJS compatibility, proper relations, and maintainability.

3. **Mixing sources creates coupling** — if an endpoint reads from both `ref_*` and `sspkg_*`, it couples the revamped system to legacy data structures. When `sspkg_*` data is removed or updated, the endpoint breaks.

4. **Data format mismatch** — `sspkg_*` uses Odoo conventions (integer IDs, numeric values). `ref_*` uses Directus conventions (UUIDs, code-based values). Mixing them causes value format conflicts in SurveyJS forms.

### Allowed Uses of `sspkg_*`

| Use Case | Allowed? |
|---|---|
| **Seeder scripts** — populating `ref_*` from `sspkg_*` snapshot | ✅ Yes |
| **Verification scripts** — comparing `ref_*` output with `sspkg_*` for correctness | ✅ Yes |
| **API endpoints returning data to clients** | ❌ Never |
| **Building SurveyJS form JSON** | ❌ Never |
| **Choices / visibleIf / validators** | ❌ Never |
| **Fetcher / builder logic** | ❌ Never |

### Migration Path

When a feature requires data that currently only exists in `sspkg_*`:

1. **Snapshot the Odoo table** — create `sspkg_*` as-is mirror (flat, no relations)
2. **Create the `ref_*` equivalent** — design the denormalized/simplified schema
3. **Write a migration** — `scripts/migrations/xxx-create-ref-*.js`
4. **Write a seeder** — populate `ref_*` from `sspkg_*` snapshot (transform + denormalize)
5. **Update the endpoint** — source from `ref_*` only
6. **Document in `docs/schema.md`** — add the new collections/fields

> **Rule:** If you find yourself writing `await svc('sspkg_*')` in an endpoint handler, STOP. The data must come from `ref_*`. If the `ref_*` equivalent doesn't exist yet, create it first.

### Related Principles
- Architectural Patterns @/Users/panjiadhiemusthofa/project/kesprimkom/services/reference/kes-ref-screening-directus/.agents/rules/architectural-pattern.md
- Schema Documentation Mandate @/Users/panjiadhiemusthofa/project/kesprimkom/services/reference/kes-ref-screening-directus/.agents/rules/schema-documentation-mandate.md
- Database Design Principles @/Users/panjiadhiemusthofa/project/kesprimkom/services/reference/kes-ref-screening-directus/.agents/rules/database-design-principles.md
