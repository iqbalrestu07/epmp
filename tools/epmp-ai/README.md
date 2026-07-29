---
id: EPMP-AI-README
title: EPMP AI Operating System
version: 1.0.0
status: Stable
owner: EPMP Architecture Team
audience:
  - Software Engineers
  - AI Coding Agents
  - Solution Architects
purpose: >
  Defines the AI Engineering Operating System used across the EPMP repository.
---

# EPMP AI Operating System

## Overview

The **EPMP AI Operating System (AI-OS)** is the official engineering framework that governs how Artificial Intelligence participates in software development within the Enterprise Property Management Platform (EPMP).

This directory is not a prompt collection.

It is not documentation for humans only.

It is the operational specification that allows human engineers and AI agents to collaborate consistently while preserving architecture quality.

---

# Goals

The AI Operating System exists to achieve five objectives.

1. Maintain architectural consistency.

2. Reduce implementation time.

3. Minimize AI hallucination.

4. Standardize AI behavior.

5. Enable multi-agent collaboration.

---

# Design Principles

The AI Operating System follows these principles.

## Human-Led

Humans own business decisions, architecture, and production responsibility.

---

## AI-Accelerated

AI accelerates implementation but never replaces architectural judgment.

---

## Specification Driven

Implementation follows specification.

AI never invents requirements.

---

## Context Driven

AI loads only the context required for the current task.

Smaller context produces higher-quality output.

---

## Deterministic

The same input should produce consistent output regardless of the AI model being used.

---

# Repository Structure

```
tools/

└── epmp-ai/

    README.md

    AI.md

    RULES.md

    SKILLS.md

    CONTEXT.md

    WORKFLOW.md

    CHECKLIST.md

    CONVENTIONS.md

    PROMPTS.md

    templates/

        module.md.template

        work-package.md.template

    roles/

        backend.md

        frontend.md

        database.md

        api.md

        reviewer.md

        documenter.md
```

---

# Reading Order

Every AI agent should read documents in the following order.

```
README.md

↓

AI.md

↓

RULES.md

↓

SKILLS.md

↓

CONTEXT.md

↓

Role File

↓

MODULE.md

↓

WORK_PACKAGE.md
```

Do not load unnecessary documentation.

---

# Human Workflow

```
Business Requirement

↓

Architecture

↓

Specification

↓

Work Package

↓

AI Implementation

↓

Human Review

↓

Merge
```

---

# AI Workflow

```
Load Context

↓

Understand Scope

↓

Implement

↓

Self Review

↓

Return Result
```

---

# Core Documents

| Document       | Purpose                         |
| -------------- | ------------------------------- |
| AI.md          | Defines AI identity and mission |
| RULES.md       | Immutable engineering rules     |
| SKILLS.md      | Engineering responsibilities    |
| CONTEXT.md     | Context loading strategy        |
| WORKFLOW.md    | Engineering workflow            |
| CHECKLIST.md   | Definition of Done              |
| CONVENTIONS.md | Naming and coding conventions   |
| PROMPTS.md     | Standard prompt templates       |

---

# Templates

Templates define reusable engineering artifacts.

Examples include:

- Module Specification
- Work Package
- ADR
- Feature Specification

Templates should be treated as the source of truth for newly created documents.

---

# Roles

Each AI agent should assume one primary role.

Examples:

- Backend Engineer
- Frontend Engineer
- Database Engineer
- API Engineer
- Reviewer
- Documentation Engineer

Role files contain role-specific instructions.

---

# Compatibility

The AI Operating System is designed to be model-agnostic.

It should work consistently with:

- ChatGPT
- Codex CLI
- Claude Code
- Gemini CLI
- Cursor
- Cline
- Roo Code
- Continue
- OpenHands
- Future AI coding agents

---

# Guiding Principle

Architecture is protected by standards.

Implementation is accelerated by AI.

Quality is enforced by review.

Documentation is the source of truth.
