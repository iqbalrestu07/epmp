---
id: EPMP-AI-CONTEXT
title: AI Context Management
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines how AI loads, manages, and limits context in order to produce
  deterministic, high-quality implementations while minimizing token usage.

audience:
  - AI Coding Agents
  - Software Engineers

depends_on:
  - README.md
  - AI.md
  - RULES.md
  - SKILLS.md
---

# Overview

Large repositories cannot be understood by loading every file.

Doing so reduces AI performance,
increases token usage,
and introduces unrelated information into the reasoning process.

Therefore, AI MUST load only the minimum context required to complete the current task.

This document defines the official Context Management Strategy of EPMP.

---

# Context Engineering Principles

## CE-001 — Minimum Required Context

AI MUST load only documents that are directly related to the current task.

More context does not necessarily produce better output.

The preferred context is the smallest context that enables a correct implementation.

---

## CE-002 — Specification Before Implementation

Implementation MUST NEVER begin before loading:

- AI.md
- RULES.md
- Relevant Role
- MODULE.md
- WORK_PACKAGE.md

---

## CE-003 — No Repository Scanning

AI SHOULD NOT scan the entire repository.

Repository-wide searches SHOULD only be performed when explicitly requested.

---

## CE-004 — Context Isolation

Every task is independent.

Context loaded for one task MUST NOT automatically influence another task.

---

# Context Hierarchy

The EPMP AI Operating System organizes context into layers.

```

Global Context

↓

Role Context

↓

Module Context

↓

Task Context

↓

Implementation

```

---

# Layer 1 — Global Context

Loaded for every task.

Required documents:

- README.md
- AI.md
- RULES.md

Purpose:

- Identity
- Engineering Rules
- Architecture Protection

---

# Layer 2 — Role Context

Loaded according to the assigned engineering role.

Examples:

Backend Engineer

```

roles/backend.md

```

Frontend Engineer

```

roles/frontend.md

```

Reviewer

```

roles/reviewer.md

```

Purpose:

Provide role-specific responsibilities.

---

# Layer 3 — Module Context

Loaded only for the module being modified.

Examples:

```

modules/

reservation/

MODULE.md

```

Purpose:

Understand:

- module boundaries
- dependencies
- events
- APIs
- ownership

---

# Layer 4 — Task Context

Loaded only for the current work item.

Example:

```

work-packages/

reservation/

create-reservation.md

```

Purpose:

Describe:

- scope
- acceptance criteria
- implementation notes

---

# Context Loading Order

AI MUST load context in the following order.

```

README

↓

AI

↓

RULES

↓

Role

↓

MODULE

↓

WORK_PACKAGE

↓

Implementation

```

Loading additional documents SHOULD require justification.

---

# Context Budget

To reduce unnecessary token consumption, AI SHOULD follow these limits.

| Context Type         | Recommendation        |
| -------------------- | --------------------- |
| Global               | Always                |
| Role                 | Always                |
| Module               | One module only       |
| Work Package         | One work package only |
| Additional Documents | Only if required      |

Avoid loading unrelated modules.

---

# Context Expansion

AI MAY expand context only when:

- required interface is unknown
- dependency cannot be determined
- architecture document is referenced
- human explicitly requests repository analysis

Otherwise,
keep context minimal.

---

# Context Switching

When switching tasks:

Unload:

- previous MODULE
- previous WORK_PACKAGE

Load:

- new MODULE
- new WORK_PACKAGE

Global context remains loaded.

---

# Context Package

A Context Package is the complete set of documents loaded for one task.

Example:

```

Context Package

README.md

AI.md

RULES.md

roles/backend.md

reservation/MODULE.md

work-packages/create-reservation.md

```

This package SHOULD be sufficient for most backend tasks.

---

# CLI Integration

Future versions of the EPMP CLI SHOULD support automatic context loading.

Example:

```

epmp context backend reservation create-reservation

```

Expected output:

```

Loading...

✓ README.md

✓ AI.md

✓ RULES.md

✓ roles/backend.md

✓ reservation/MODULE.md

✓ create-reservation.md

Done.

```

---

# AI Context Lifecycle

```

Receive Task

↓

Select Role

↓

Load Global Context

↓

Load Role Context

↓

Load Module Context

↓

Load Work Package

↓

Implement

↓

Self Review

↓

Unload Module Context

↓

Task Completed

```

---

# Context Validation

Before implementation begins, AI MUST verify:

✓ Required documents loaded

✓ Module identified

✓ Work package identified

✓ Role assigned

✓ Scope understood

If any item fails,

implementation MUST NOT begin.

---

# Final Principle

Good AI is not the AI that reads everything.

Good AI is the AI that reads the right things.
