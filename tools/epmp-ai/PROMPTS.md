---
id: EPMP-AI-PROMPTS
title: Prompt Specification
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the standard prompt specification used to communicate with AI
  agents. Prompts describe intent, context, constraints, and expected
  deliverables rather than raw instructions.

depends_on:
  - AI.md
  - RULES.md
  - CONTEXT.md
  - CONVENTIONS.md
---

# Overview

Prompts are contracts.

A prompt MUST describe:

- Who the AI should become
- What the AI should accomplish
- What context should be loaded
- What constraints apply
- What output is expected

A good prompt minimizes ambiguity.

---

# Prompt Structure

Every prompt consists of five sections.

```
ROLE

↓

TASK

↓

CONTEXT

↓

CONSTRAINTS

↓

OUTPUT
```

---

# ROLE

Defines the engineering role.

Examples

Backend Engineer

Frontend Engineer

Database Engineer

Reviewer

Documentation Engineer

Architect

---

# TASK

Describes the objective.

Examples

Create Reservation

Update Payment Flow

Review Pull Request

Generate API Documentation

Refactor Tenant Module

---

# CONTEXT

Defines which documents must be loaded.

Example

README.md

AI.md

RULES.md

roles/backend.md

modules/reservation/MODULE.md

work-packages/create-reservation.md

---

# CONSTRAINTS

Examples

Follow Clean Architecture

Do not modify other modules

Follow naming conventions

Do not change API contracts

Do not introduce new dependencies

---

# OUTPUT

Examples

Go Source Code

Unit Tests

Migration

Updated Documentation

Review Notes

---

# Prompt Quality Rules

A prompt SHOULD:

- Define one objective
- Load minimal context
- Define expected output
- Define engineering role

A prompt SHOULD NOT:

- Combine unrelated tasks
- Leave scope ambiguous
- Assume repository knowledge
- Require repository-wide scanning

---

# Standard Prompt Template

ROLE

...

TASK

...

CONTEXT

...

CONSTRAINTS

...

OUTPUT

...

---

# Final Principle

A clear prompt produces predictable software.
