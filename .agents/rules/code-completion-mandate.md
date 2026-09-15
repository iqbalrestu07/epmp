---
trigger: always_on
---

## Code Completion Mandate

### Universal Requirement

**Before marking any code task as complete, you MUST run automated quality checks and remediate all issues.**

This is NOT OPTIONAL. Delivering code without validation violates the Rugged Software Constitution @rugged-software-constitution.md.

### The Completion Checklist

Every code generation task follows this workflow:

1. **Generate** - Write the code based on requirements
2. **Validate** - Run language-appropriate quality checks (see below)
3. **Remediate** - Fix all detected issues
4. **Verify** - Re-run checks to confirm fixes
5. **Deliver** - Mark task complete only after all checks pass

**Never skip validation "to save time." Validation IS the work.**

### Quality Commands for This Project

This project uses **plain ESM JavaScript** (no TypeScript, no build step).

| Check                | Command                                                        | Notes                                         |
| -------------------- | -------------------------------------------------------------- | --------------------------------------------- |
| **Syntax check**     | `node --check extensions/survey-forms/index.js`               | Ekstensi API utama |
| **Syntax check**     | `node --check scripts/seed-samples.js`                        | Seeder data |
| **Syntax check**     | `node --check scripts/run-migrations.js`                      | Migration runner |
| **Dependency audit** | `pnpm audit` or `npm audit`                                    | Check for known CVEs |
| **Docker build**     | `docker compose build`                                         | Verifies Dockerfile and dependency resolution |
| **Manual smoke test**| Start Docker, call endpoints, verify JSON                      | Primary validation for endpoint extensions |

> Since this project has no ESLint/Prettier configured, syntax checking via `node --check` is the minimum automated validation.
> For more thorough validation, manually test endpoints against a running Directus instance.

### Failure Protocol

**If any quality check fails:**

1. Read the error output completely
2. Fix the identified issues in the code
3. Re-run the failing command
4. Do not proceed until all checks pass

> Never disable a lint rule or suppress a warning to make checks pass. Fix the root cause.

### Related Principles
- Rugged Software Constitution @rugged-software-constitution.md
- Code Idioms and Conventions @code-idioms-and-conventions.md
- Directus Extension Patterns @directus-extension-patterns.md