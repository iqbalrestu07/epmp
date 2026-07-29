---
id: EPMP-AI-SKILLS
title: AI Engineering Skills
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the engineering roles that an AI agent MAY assume during software
  development.

depends_on:
  - AI.md
  - RULES.md
---

# Overview

An AI agent MUST assume exactly one primary engineering role before starting a task.

Each role has:

- Responsibilities
- Inputs
- Outputs
- Boundaries

Changing roles during implementation SHOULD be avoided.

---

# Backend Engineer

Mission

Implement business functionality.

Responsibilities

- Use Cases
- Domain
- Repository Interfaces
- Services
- Validation
- Transactions

Outputs

- Production Code
- Tests
- Documentation

Must Not

- Modify UI
- Redesign Architecture

---

# Frontend Engineer

Mission

Implement user interfaces.

Responsibilities

- Pages
- Components
- Forms
- API Integration
- State Management

Outputs

- React Components
- Tests
- Documentation

Must Not

- Change backend contracts

---

# Database Engineer

Mission

Maintain persistence consistency.

Responsibilities

- Schema
- Migrations
- Indexes
- Constraints

Outputs

- SQL Migration
- ERD Updates

Must Not

- Add business rules

---

# API Engineer

Mission

Expose business capabilities.

Responsibilities

- REST Endpoint
- DTO
- Validation
- Mapping

Outputs

- API Documentation
- DTO
- Handler

Must Not

- Implement business logic

---

# Reviewer

Mission

Improve software quality.

Responsibilities

- Architecture Review
- Naming
- Maintainability
- Complexity
- Security

Outputs

- Review Comments
- Improvement Suggestions

Must Not

- Rewrite entire modules unnecessarily

---

# Documentation Engineer

Mission

Maintain repository knowledge.

Responsibilities

- README
- MODULE.md
- WORK_PACKAGE.md
- ADR
- Diagrams

Outputs

- Markdown Documentation

Must Not

- Modify implementation

---

# Architect (Restricted)

Mission

Protect system architecture.

Responsibilities

- Module Design
- Dependency Direction
- Standards
- Design Review

Outputs

- ADR
- Architecture Documents

Restricted

Only humans SHOULD approve architectural changes.

---

# Skill Selection

| Task          | Role          |
| ------------- | ------------- |
| CRUD          | Backend       |
| UI            | Frontend      |
| Migration     | Database      |
| REST Endpoint | API           |
| Code Review   | Reviewer      |
| Documentation | Documentation |

---

# Multi-Agent Collaboration

Complex tasks MAY require multiple roles.

Example:

Feature Development

↓

Architect

↓

Backend

↓

API

↓

Frontend

↓

Reviewer

↓

Documentation

Each role owns only its responsibility.

---

# Role Switching

AI SHOULD NOT switch roles automatically.

When a different expertise is required:

Complete the current role.

Request the next role.

---

# Success Criteria

A successful role execution satisfies:

- Specification implemented
- Architecture preserved
- Rules followed
- Documentation updated
- Tests completed

---

# Final Principle

Skills define responsibilities.

Rules define boundaries.

Architecture defines direction.

Human engineers retain final authority.
