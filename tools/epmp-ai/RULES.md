---
id: EPMP-AI-RULES
title: AI Engineering Rules
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the immutable engineering rules that every AI agent MUST follow
  when participating in EPMP software development.

audience:
  - AI Coding Agents
  - Software Engineers

depends_on:
  - README.md
  - AI.md
---

# Overview

This document defines the non-negotiable engineering rules of the EPMP AI Operating System.

These rules exist to protect:

- Architecture
- Business Rules
- Repository Consistency
- Code Quality
- Long-term Maintainability

Unless explicitly instructed by the Software Architect, these rules MUST NOT be violated.

---

# RFC 2119 Keywords

The key words:

- MUST
- MUST NOT
- REQUIRED
- SHALL
- SHALL NOT
- SHOULD
- SHOULD NOT
- MAY

are to be interpreted as described in RFC 2119.

---

# 1. Architecture Rules

### AR-001

AI MUST preserve the approved architecture.

### AR-002

AI MUST NOT redesign system architecture.

### AR-003

AI MUST NOT introduce new architectural patterns.

### AR-004

AI MUST respect Clean Architecture dependency direction.

Allowed

```
Presentation
↓

Application
↓

Domain

↑

Infrastructure
```

Forbidden

```
Domain

↓

Infrastructure
```

---

# 2. Business Rules

### BR-001

Business rules MUST exist only inside the Domain Layer.

### BR-002

AI MUST NOT duplicate business rules.

### BR-003

Business terminology MUST follow the Ubiquitous Language.

### BR-004

AI MUST NOT invent business behavior.

If behavior is missing:

STOP.

Request clarification.

---

# 3. Module Rules

### MR-001

Each module owns its business capability.

### MR-002

Modules MUST remain loosely coupled.

### MR-003

Cross-module communication SHOULD use:

- Application Services
- Domain Events

### MR-004

AI MUST NOT access another module's persistence directly.

---

# 4. Repository Rules

### RR-001

Repository interfaces belong to the Domain Layer.

### RR-002

Repository implementations belong to Infrastructure.

### RR-003

Repositories persist data only.

Repositories MUST NOT contain business logic.

---

# 5. Application Rules

Application Services MAY:

- orchestrate use cases
- manage transactions
- coordinate repositories
- publish events

Application Services MUST NOT contain business rules.

---

# 6. Domain Rules

Entities MUST protect invariants.

Aggregates MUST guarantee consistency.

Value Objects MUST be immutable.

Domain Events MUST describe business facts.

---

# 7. API Rules

REST APIs MUST be:

- Stateless
- Versioned
- Validated

Request validation MUST happen before entering the Domain Layer.

---

# 8. Database Rules

Schema changes MUST use migrations.

Production data MUST NOT be modified manually by AI.

AI MUST NEVER generate destructive SQL unless explicitly requested.

---

# 9. Documentation Rules

Every meaningful behavior change MUST update documentation.

Every new module MUST include:

- MODULE.md

Every completed task SHOULD update:

- WORK_PACKAGE.md

---

# 10. Testing Rules

New business logic MUST include tests.

Bug fixes SHOULD include regression tests.

AI MUST NOT remove failing tests to make builds pass.

---

# 11. Refactoring Rules

Refactoring MUST preserve behavior.

Public API changes REQUIRE approval.

Large refactors SHOULD be divided into smaller commits.

---

# 12. Security Rules

Secrets MUST NEVER be hardcoded.

Credentials MUST come from configuration.

Sensitive data MUST NOT appear in logs.

---

# 13. AI Behavior Rules

AI MUST:

- Ask when uncertain
- Preserve architecture
- Follow specifications
- Keep changes minimal
- Explain assumptions

AI MUST NOT:

- Guess requirements
- Invent APIs
- Rename modules
- Change folder structure
- Ignore standards

---

# 14. Conflict Resolution

Priority order:

1. Human Instructions

2. Approved Specification

3. Architecture Documents

4. RULES.md

5. AI Preferences

If conflicts occur:

Stop.

Explain the conflict.

Request clarification.

---

# Compliance Checklist

Before completing any task, verify:

✓ Architecture preserved

✓ Business rules unchanged

✓ Naming consistent

✓ Tests updated

✓ Documentation updated

✓ No duplicated logic

✓ No unnecessary complexity

If any item fails, revise the implementation before returning the result.

---

# Final Principle

Good engineering favors consistency over cleverness.

AI MUST optimize for maintainability, not novelty.
