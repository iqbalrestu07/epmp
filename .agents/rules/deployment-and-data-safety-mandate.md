# Deployment and Data Safety Mandate

### Core Principle: Zero Data-Loss Architecture
Deployments and migrations between environments (`dev` ➔ `stg` ➔ `prod`) MUST NEVER destroy, drop, or truncate collections, fields, or existing database records unless explicitly requested with destructive flags and interactive confirmation.

---

### 1. Schema Migration Safety (`/schema/diff` & `/schema/apply`)

#### The Directus Partial Snapshot Trap
Directus native schema engine (`/schema/diff`) treats the incoming snapshot as the *absolute desired state*. If a snapshot only contains a subset of collections (e.g. during partial deployment with `--collection=ref_tl_*`), Directus will automatically generate `DROP` actions (`kind: "D"`) for every single collection and field absent from the snapshot.

#### Mandatory Safeguards:
1. **Never Apply Raw Diffs**:
   All diff payloads from Directus `/schema/diff` MUST be sanitized before being sent to `/schema/apply`:
   - Filter out all collections, fields, and relations whose change consists entirely of deletions (`kind: "D"`), unless `--allow-delete` or `--destructive` is explicitly provided.
   - For partial deployments (`--collection=...`), strictly discard any diff entries belonging to collections outside the target scope.
2. **Logged & Blocked Drops**:
   Any blocked drop instruction must be logged to stderr/stdout with `🛡️ KEAMANAN: X instruksi DROP/DELETE skema diblokir otomatis`.
3. **No Direct Admin GUI Schema Alterations**:
   Schema definitions must be codified in `scripts/migrations/` and deployed via `scripts/deploy-collection.js` or `scripts/run-migrations.js`.

---

### 2. Data Synchronization Safety

1. **Additive & Idempotent by Default**:
   - The default data mode for all transfer scripts MUST be `skip-existing` (insert only records whose primary key does not yet exist on the target).
   - Deletion of target data (`--replace`) is strictly opt-in and requires explicit confirmation in production.
2. **Virtual Relational Alias Sanitization**:
   - Virtual M2M/O2M alias fields (fields with `type: "alias"` or `schema: null`, such as `age_groups`, `topic_code`, `layanan_code`) MUST be stripped from insert payloads.
   - Junction tables (e.g., `ref_cluster_age_groups`, `ref_topic_ref_service`) are distinct collections and must be seeded independently in strict dependency order.
3. **Dynamic Primary Key Resolution**:
   - Do not assume `id` is the primary key for all tables. Collections like `ref_topic` use `code` as their primary key. Scripts must inspect or fallback between `id` and `code`.
4. **Token Expiry Resilience**:
   - Batch insert operations for large collections (e.g. `ref_zscore_reference` > 10,000 rows) must proactively refresh admin session tokens and implement automatic re-login on HTTP 401 (`TOKEN_EXPIRED`).

---

### 3. Environment Roles & Truth Hierarchy

- **DEV**: Active development and prototyping space (CMS UI experiments, form testing). Data may undergo temporary mutations.
- **STG**: Pre-production staging environment. Skema and reference data must mirror production stability.
- **PROD**: Golden source of truth for business and reference data.
- **Safe Partial Deployment**: When shipping new feature collections (e.g. Tatalaksana V3 `ref_tl_*`) from DEV to STG, only the target feature collections are allowed to be transferred. STG baseline master data (collections 1–29) must remain 100% untouched.

### Enforcement Checklist
Before executing any cross-environment script:
- [ ] Are schema drops (`kind: "D"`) blocked by default?
- [ ] Is data transfer using `skip-existing` without destructive truncates?
- [ ] Are virtual relational alias fields stripped prior to row insertion?
- [ ] Is token refresh handling active for long-running batches?
