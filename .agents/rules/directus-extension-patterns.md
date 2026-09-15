---
trigger: always_on
---

## Directus Extension Patterns

> This project is a Directus v11 headless CMS instance with custom endpoint extensions.
> All code is plain **ESM JavaScript** (not TypeScript). There is no build step.

### Architecture Context

- **Runtime:** Directus v11 on Node.js (Docker)
- **Database:** PostgreSQL (via Knex.js, managed by Directus)
- **Language:** JavaScript (ESM — `"type": "module"`)
- **Extensions:** Endpoint extensions under `extensions/`
- **Schema management:** `setup-survey-schema.js` via Directus REST API
- **SurveyJS target:** v1 (major version 1) — semua output JSON dari extension harus valid dan kompatibel dengan SurveyJS v1

### Extension Authoring Rules

**1. Default export signature**
```javascript
// Endpoint extension
export default (router, context) => {
  const { services, getSchema } = context;
  const { ItemsService } = services;

  router.get('/', async (req, res) => {
    // handler
  });
};
```

**2. Use Directus services, not raw SQL**
```javascript
// ✅ Correct — uses Directus abstraction
const service = new ItemsService('collection_name', {
  schema: await getSchema(),
  accountability: { admin: true },
});
const items = await service.readByQuery({ filter: { ... } });

// ❌ Wrong — bypasses Directus permissions & hooks
const items = await knex('collection_name').where({ ... });
```

**3. Error handling**
```javascript
router.get('/:id', async (req, res) => {
  try {
    const schema = await getSchema();
    const service = new ItemsService('forms', { schema, accountability: { admin: true } });

    let item;
    try {
      item = await service.readOne(req.params.id, { fields: ['*'] });
    } catch {
      return res.status(404).json({ error: 'Not found' });
    }

    res.json({ data: item });
  } catch (err) {
    console.error('[extension-name] Error:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});
```

**4. Logging**
- Use `console.log` / `console.error` in extensions (Directus captures stdout/stderr)
- Prefix with extension name: `[survey-forms]`
- Do not import external logging libraries — keep extensions zero-dependency

**5. Response format**
- Always wrap responses in `{ data: ... }` to match Directus API convention
- Use standard HTTP status codes (200, 404, 500)

### Schema Management

The schema is defined declaratively in `setup-survey-schema.js` and applied via Directus REST API:

```bash
# Apply schema (destructive — resets all data)
node setup-survey-schema.js
```

**Rules:**
- Schema changes go in `scripts/migrations/` — **NEVER via GUI Directus Admin**
- Document ALL changes in `docs/schema.md` with a Changelog entry — see Schema Documentation Mandate @schema-documentation-mandate.md
- Bump the version in `docs/schema.md` header on each change
- **Strict 1-to-1 Field Synchronization**: Every single field created or touched via migration MUST be explicitly mapped in `docs/schema.md` without any missing fields, even for standard or system fields.
- **Unique Identifier (`internal_name`)**: A unique string slug (`internal_name`) is mandatory for standard payload deduplication checks inside Seeder Scripts (e.g. `seed-samples.js`). Never ignore or skip this core element.

### Testing Approach

Directus extensions are thin API adapters. Testing strategy:

1. **Manual testing** — primary method: run Docker, hit endpoints
2. **Integration testing** — spin up Directus in Docker, seed data, call extension endpoints
3. **Unit testing** — only for pure transformation functions extracted from handlers

> Because extensions tightly depend on Directus `ItemsService`, mocking is complex.
> Prefer integration tests over isolated unit tests for endpoint logic.

### Docker Workflow

```bash
# Start Directus
docker compose up -d

# Rebuild after Dockerfile changes
docker compose up -d --build

# Apply schema
node setup-survey-schema.js

# View logs
docker compose logs -f directus

# Stop
docker compose down
```

### Extension Versioning

**Every extension has a `version` field in its `package.json`. When making changes to an extension, the version MUST be bumped before committing.**

**Semver rules:**
- **patch** (`1.0.1`): bug fixes, typos, cosmetic tweaks
- **minor** (`1.1.0`): new features, UI additions, behavioral changes
- **major** (`2.0.0`): breaking changes, major rewrites

**How to bump:**
1. Edit `extensions/<name>/package.json` → update `"version"` field
2. Commit with: `chore(<extension-name>): bump version to X.Y.Z`

**When to bump:**
- After all functional changes to the extension are done (not per-commit during a session)
- Before the final commit/push of a work session that touched the extension

> **Note:** Some extensions (e.g. `playground`) display the version in the UI footer.
> The `index.js` reads the version from `package.json` at startup — no need to edit HTML templates.

### Related Principles
- Project Structure @project-structure.md
- API Design Principles @api-design-principles.md
- Error Handling Principles @error-handling-principles.md
- Schema Documentation Mandate @schema-documentation-mandate.md
- Deployment and Data Safety Mandate @deployment-and-data-safety-mandate.md
