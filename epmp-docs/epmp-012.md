# EPMP-012 — AI Engineering Standard

| Metadata      | Value                                                               |
| ------------- | ------------------------------------------------------------------- |
| Document ID   | EPMP-012                                                            |
| Version       | 2.0.0                                                               |
| Status        | Approved                                                            |
| Owner         | Software Architecture Team                                          |
| Depends On    | EPMP-001 ~ EPMP-011                                                 |
| Referenced By | AI.md, RULES.md, SKILLS.md, MODULE.md, Work Package, Prompt Library |

---

# 1. Purpose

Dokumen ini mendefinisikan standar resmi penggunaan Artificial Intelligence pada seluruh siklus pengembangan Enterprise Property Management Platform (EPMP).

AI diperlakukan sebagai **Engineering Participant** yang bekerja bersama engineer manusia.

Standar ini bertujuan untuk:

- menjaga kualitas implementasi,
- mengurangi hallucination,
- meningkatkan konsistensi,
- menghemat waktu implementasi,
- mengoptimalkan penggunaan context/token,
- menjaga Architecture Integrity.

---

# 2. Engineering Philosophy

EPMP menggunakan filosofi berikut.

```

Human Leads

↓

Architecture

↓

Specification

↓

AI Implements

↓

Human Reviews

↓

AI Refines

↓

Human Approves

↓

Merge

```

AI bukan pengambil keputusan.

AI adalah implementation partner.

---

# 3. Human + AI Collaboration Model

## Human Responsibilities

Human memiliki tanggung jawab penuh terhadap:

- Business Requirement
- Domain Model
- Architecture
- Technical Decision
- Code Review
- Merge
- Production Deployment

---

## AI Responsibilities

AI membantu:

- CRUD
- Use Case
- Repository
- DTO
- API
- Unit Test
- Integration Test
- Refactoring
- Documentation

---

## Shared Responsibilities

Dilakukan bersama:

- Debugging
- Optimization
- Security Review
- Database Evolution
- Refactoring Planning

---

# 4. AI Engineering Principles

Semua AI Agent wajib mengikuti prinsip berikut.

## AI-01

Specification is Source of Truth.

---

## AI-02

Architecture Before Code.

---

## AI-03

One Task One Context.

---

## AI-04

Documentation Before Memory.

---

## AI-05

Verification Over Generation.

---

## AI-06

Human Owns The Outcome.

---

## AI-07

Generator Before Boilerplate.

---

# 5. AI Capability Declaration

AI memiliki capability berikut.

| Capability             | Allowed |
| ---------------------- | ------- |
| Generate CRUD          | YES     |
| Generate DTO           | YES     |
| Generate Test          | YES     |
| Generate API           | YES     |
| Generate Documentation | YES     |
| Refactor               | YES     |
| Modify Business Rule   | NO      |
| Modify Architecture    | NO      |
| Create New Pattern     | NO      |
| Change Domain Model    | NO      |

Semua perubahan arsitektur harus dilakukan oleh manusia.

---

# 6. AI Context Strategy

AI hanya diberikan context yang diperlukan.

Prioritas:

1. Work Package
2. Module Specification
3. MODULE.md
4. AI.md
5. EPMP Standard
6. Source Code

AI tidak perlu membaca seluruh repository.

---

# 7. Context Package

Minimal context:

```

Task

Acceptance Criteria

Referenced Documents

Files Allowed

Files Forbidden

Review Checklist

```

Semakin kecil context,
semakin tinggi kualitas AI.

---

# 8. AI Workflow

```

Requirement

↓

Specification

↓

Work Package

↓

AI Context

↓

Implementation

↓

Testing

↓

Human Review

↓

Merge

```

Tidak diperbolehkan melewati Specification.

---

# 9. AI Task Granularity

Ukuran task AI:

Target:

1–4 jam pekerjaan engineer.

Contoh baik:

```

Implement Create Reservation

```

Contoh buruk:

```

Build Reservation Module

```

---

# 10. Prompt Standard

Prompt harus deterministic.

Selalu menyebutkan:

- Task
- Scope
- References
- Acceptance Criteria
- Files Allowed
- Files Forbidden

---

# 11. Generator Policy

Apabila helper-package menyediakan generator,

AI **WAJIB** menggunakan generator tersebut.

AI tidak boleh membuat CRUD baru.

Contoh:

```

helper-package/

be/

crud/

repository/

migration/

```

dan

```

helper-package/

fe/

crud/

table/

form/

```

Generator merupakan implementasi resmi EPMP.

---

# 12. AI Review Policy

AI wajib melakukan self review.

Checklist:

- Clean Architecture
- DDD
- Naming
- Test
- Documentation

Setelah itu dilakukan Human Review.

---

# 13. Human Review Policy

Semua hasil AI wajib direview.

Minimal review:

- Architecture
- Domain
- Security
- Maintainability

---

# 14. AI Memory Policy

AI tidak boleh diasumsikan mengingat percakapan sebelumnya.

Semua informasi berasal dari:

- Documentation
- Source Code
- Context Package

---

# 15. AI Governance

AI tidak memiliki hak:

- Merge
- Release
- Architecture Decision
- Business Decision

Hak tersebut dimiliki manusia.

---

# 16. AI Knowledge Sources

Urutan sumber informasi.

```

Specification

↓

MODULE.md

↓

AI.md

↓

RULES.md

↓

Source Code

↓

Conversation

```

Conversation memiliki prioritas paling rendah.

---

# 17. AI Performance Metrics

Diukur berdasarkan:

- Acceptance Rate
- Rework Rate
- Review Time
- Bug Escape Rate
- Test Coverage
- Documentation Coverage

Bukan berdasarkan Lines of Code.

---

# 18. AI Failure Recovery

Jika AI gagal:

1.  Periksa Specification.

2.  Periksa Context.

3.  Periksa Work Package.

4.  Perkecil Scope.

5.  Generate ulang.

---

# 19. AI Evolution

EPMP dirancang untuk mendukung banyak AI.

Contoh:

- Backend Agent
- Frontend Agent
- API Agent
- Database Agent
- Reviewer Agent
- Testing Agent
- Documentation Agent

Semua bekerja menggunakan standar yang sama.

---

# 20. Required AI Documents

Seluruh repository wajib memiliki:

```

AI.md

RULES.md

SKILLS.md

PROMPTS.md

MODULE.md

```

Semua AI Agent menggunakan dokumen tersebut sebagai konteks utama.

---

# Closing Statement

AI Engineering Standard menjadikan AI sebagai anggota resmi tim engineering EPMP.

Dengan spesifikasi yang jelas, context yang kecil, generator resmi, dan review manusia, AI dapat menghasilkan implementasi yang cepat tanpa mengorbankan kualitas maupun integritas arsitektur.
