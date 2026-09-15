---
trigger: model_decision
description: When creating branches, committing code, managing PRs, or working with version control
---

## Git Workflow Principles

### Commit Messages — Conventional Commits

**Format:**
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
| Type       | Purpose                          |
| ---------- | -------------------------------- |
| `feat`     | New feature                      |
| `fix`      | Bug fix                          |
| `docs`     | Documentation only               |
| `style`    | Formatting, semicolons, etc.     |
| `refactor` | Code change (no new feature/fix) |
| `test`     | Adding or updating tests         |
| `chore`    | Maintenance, dependencies        |
| `perf`     | Performance improvement          |
| `ci`       | CI/CD configuration changes      |

**Rules:**
- Description is imperative mood ("add" not "added", "fix" not "fixes")
- Scope matches the feature area (e.g., `task`, `auth`, `ui`)
- Description is concise (<72 characters)
- Body explains **why**, not what (the diff shows what)

### Branch Naming

**Format:** `<type>/<ticket-or-short-description>`

**Examples:**
```
feat/task-crud-api
fix/auth-token-expiry
refactor/storage-layer
chore/update-deps
```

**Rules:**
- Use lowercase with hyphens (kebab-case)
- Prefix matches commit type
- Keep branch names short but descriptive

### Commit Hygiene

- **One logical change per commit** — don't mix unrelated changes
- **Never commit broken tests** — all tests must pass before committing
- **Don't commit debug code** — remove console.log, print statements, TODO hacks
- **Don't commit secrets** — use `.gitignore` and environment variables

### Push Policy

> **CRITICAL: Jangan pernah `git push` otomatis.**

- Agent hanya boleh `git add` + `git commit`. **TIDAK BOLEH** langsung `git push`.
- Setelah commit, biarkan user yang memutuskan kapan push ke remote.
- Jika user meminta push secara eksplisit, barulah jalankan `git push`.

### PR Size Guidelines

- **Ideal:** <400 lines changed
- **Acceptable:** 400-800 lines
- **Too large:** >800 lines — split into smaller PRs

**Why:** Large PRs get rubber-stamped. Small PRs get thoughtful reviews.

### Merge Strategy

- **Feature branches → main:** Squash merge (clean history)
- **Release branches:** Merge commit (preserve history)
- **Hotfixes:** Cherry-pick to affected branches

### Branch Strategy (Tim)

Repositori ini menggunakan **3-tier GitFlow**:

| Branch | Environment | Fungsi | Siapa yang Push |
|---|---|---|---|
| `master` | Production | Release final yang stabil | Tidak langsung — hanya via MR dari `staging` |
| `staging` | Staging/UAT | Validasi sebelum ke production | Merge dari `development` atau `hotfix/` |
| `development` | Development | Integrasi semua fitur harian | Merge dari `feature/` |
| `feature/<deskripsi>` | Local | Branch kerja per fitur/task | Developer individual |
| `hotfix/<deskripsi>` | — | Perbaikan kritis production | Developer → `master` + cherry-pick ke `staging` & `development` |

**Alur normal:**
```
feature/* → development → staging → master
```

**Aturan:**
- Selalu buat branch dari `development`: `git checkout -b feature/nama-fitur development`
- PR `feature/*` → `development`
- PR `development` → `staging` (setelah fitur siap QA)
- PR `staging` → `master` (setelah UAT lulus)
- PR `hotfix/*` → `master`, lalu cherry-pick ke `staging` dan `development`
- Hapus branch `feature/` dan `hotfix/` setelah di-merge

### Git Workflow Checklist

- [ ] Branch named with correct type prefix?
- [ ] All commits follow conventional format?
- [ ] No debug code or secrets committed?
- [ ] All tests pass before committing?
- [ ] PR is <400 lines (or justified if larger)?
- [ ] Commit messages explain why, not just what?

### Related Principles
- Code Completion Mandate @code-completion-mandate.md
- Testing Strategy @testing-strategy.md
- Security Mandate @security-mandate.md
