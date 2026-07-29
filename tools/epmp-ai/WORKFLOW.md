---
id: EPMP-AI-WORKFLOW
title: AI Engineering Workflow
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the standard engineering workflow followed by AI agents from
  receiving a task until completing implementation.

depends_on:
  - CONTEXT.md
  - RULES.md
---

# Overview

Every engineering task MUST follow the same workflow.

The workflow ensures consistency,
predictability,
and reviewability.

Skipping workflow steps is NOT allowed.

---

# Workflow Stages

The AI Engineering Workflow consists of eight stages.

```

Task

↓

Understand

↓

Load Context

↓

Plan

↓

Implement

↓

Self Review

↓

Document

↓

Complete

```

---

# Stage 1 — Receive Task

Input:

- Work Package
- Human Request

Output:

- Task Understanding

AI MUST identify:

- objective
- scope
- module
- engineering role

---

# Stage 2 — Understand

Before writing code,
AI MUST understand:

- business objective
- acceptance criteria
- architecture impact

Questions MUST be asked if requirements are incomplete.

---

# Stage 3 — Load Context

Load context according to CONTEXT.md.

Context MUST remain minimal.

Do not load unrelated modules.

---

# Stage 4 — Plan

Before implementation,
AI SHOULD create a short execution plan.

Example:

1. Update DTO

2. Create Use Case

3. Update Handler

4. Add Tests

5. Update Documentation

The plan SHOULD be small enough to review easily.

---

# Stage 5 — Implement

Implementation MUST:

- preserve architecture
- respect module boundaries
- minimize changes
- avoid duplicated logic

Generated code SHOULD be production-ready.

---

# Stage 6 — Self Review

AI MUST review:

Architecture

Naming

Error Handling

Tests

Documentation

Security

Performance

If issues are found,

fix them before continuing.

---

# Stage 7 — Documentation

Implementation is incomplete until documentation is updated.

Possible updates:

- MODULE.md

- WORK_PACKAGE.md

- API.md

- README.md

Documentation SHOULD explain why,
not only what.

---

# Stage 8 — Complete

Before completing the task,
AI MUST verify:

✓ Acceptance Criteria

✓ Tests

✓ Documentation

✓ No Architecture Violations

✓ No TODO Left Behind

Only then may the task be considered complete.

---

# Failure Handling

If implementation cannot continue:

Stop.

Explain:

- what is missing
- why it blocks progress
- what information is required

Never guess.

---

# Definition of Success

A task is successful when:

- Business requirement satisfied
- Architecture preserved
- Tests pass
- Documentation updated
- Human review ready

---

# Workflow Diagram

```

Receive Task

↓

Understand

↓

Load Context

↓

Plan

↓

Implement

↓

Self Review

↓

Documentation

↓

Complete

↓

Human Review

↓

Merge

```

---

# Final Principle

Fast implementation is valuable.

Correct implementation is mandatory.

The goal is not to finish quickly.

The goal is to finish correctly.
