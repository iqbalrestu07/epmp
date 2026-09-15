## Schema Documentation Mandate

> **This rule has the same enforcement weight as `code-completion-mandate.md`.**
> Delivering a migration or schema change without updating `docs/schema.md` is INCOMPLETE work.

### Core Requirement

**Every change to a Directus collection or field MUST be documented in `docs/schema.md` in the same commit.**

This applies to ALL operations:
- ✅ Creating a new collection
- ✅ Adding a field to an existing collection
- ✅ Renaming or removing a field
- ✅ Changing a field's type, interface, or relationship
- ✅ Adding or changing display templates, conditions, or metadata

No exceptions. Not "will document later." Not "minor field."

---

### When This Rule Applies

This rule activates whenever you:
1. Write or modify a file in `scripts/migrations/`
2. Call any Directus REST API that touches `/collections` or `/fields`
3. Apply any schema change via `run-migrations.js`
4. Discover an undocumented field during development

---

### Documentation Format

Each collection in `docs/schema.md` must have a table with these columns:

```markdown
| Field Name | Type & Interface | Keterangan / Kegunaan |
| :--- | :--- | :--- |
| `field_name` | Type (Interface) | Clear, human-readable purpose |
```

**Every single field that exists in the live Directus instance must appear in this table.**
This includes:
- Custom fields you added
- Standard/system fields (`id`, `date_created`, `date_updated`, `user_created`, `user_updated`, `sort`, `status`)
- Alias and relational fields (O2M, M2O markers)

**For new fields, add to the correct collection section**, and add a **Changelog entry**:

```markdown
### [X.Y.Z] — YYYY-MM-DD
#### Added
- `collection_name.field_name` — brief description of why it was added
```

---

### Strict 1-to-1 Field Synchronization Rule

`docs/schema.md` must be a **mirror** of the live database schema. If a field exists in Directus but not in `docs/schema.md`, the schema documentation is **out of sync** and must be fixed immediately.

**Enforcement:**
- Before writing any migration: read the relevant collection section in `docs/schema.md`
- After writing any migration: update `docs/schema.md` to reflect all changes
- If you discover an undocumented field during any task: document it immediately before continuing

---

### Discovery Protocol for Undocumented Fields

If you encounter a field in Directus that is missing from `docs/schema.md`:

1. **Stop current task**
2. Investigate: check `scripts/migrations/` for when the field was added
3. Document the field in `docs/schema.md` under the correct collection
4. Add a Changelog entry with an estimated date and note: `[Retroactive documentation]`
5. Resume current task

---

### AI Agent Checklist (Pre-flight for Any Migration Task)

Before writing any migration file, verify:
- [ ] Have I read the relevant collection section in `docs/schema.md`?
- [ ] Does `docs/schema.md` accurately reflect the current live schema?
- [ ] Is there any undocumented field I need to backfill documentation for first?

After writing any migration file, verify:
- [ ] Have I added all new fields to the correct collection table in `docs/schema.md`?
- [ ] Have I added a Changelog entry at the bottom of `docs/schema.md`?
- [ ] Do I bump the version number in the `docs/schema.md` header?

---

### Violation Examples

❌ **Wrong:** Add a migration that adds `questions.is_zscore` but only update the JS file.

✅ **Correct:** Add migration + add `is_zscore` row to the `questions` table in `docs/schema.md` + add Changelog entry.

❌ **Wrong:** Discover `questions.is_zscore` is undocumented and ignore it to finish the current task faster.

✅ **Correct:** Immediately document `questions.is_zscore` in `docs/schema.md` before continuing.

---

### Related Principles
- Database Design Principles @/Users/panjiadhiemusthofa/project/kesprimkom/kes-ssi/data-reference-module/kes-ref-screening-directus/.agents/rules/database-design-principles.md
- Directus Extension Patterns @/Users/panjiadhiemusthofa/project/kesprimkom/kes-ssi/data-reference-module/kes-ref-screening-directus/.agents/rules/directus-extension-patterns.md
- Code Completion Mandate @/Users/panjiadhiemusthofa/project/kesprimkom/kes-ssi/data-reference-module/kes-ref-screening-directus/.agents/rules/code-completion-mandate.md
- Project Structure @/Users/panjiadhiemusthofa/project/kesprimkom/kes-ssi/data-reference-module/kes-ref-screening-directus/.agents/rules/project-structure.md
