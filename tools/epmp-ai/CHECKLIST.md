---
id: EPMP-AI-CHECKLIST
title: Engineering Checklist
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the mandatory engineering checklist that MUST be completed before
  any implementation is considered finished.

depends_on:
  - RULES.md
  - CONVENTIONS.md
---

# Overview

Every engineering task MUST pass this checklist.

Completion is defined by quality, not by writing code.

---

# Architecture Checklist

□ Clean Architecture preserved

□ Module boundaries respected

□ No forbidden dependencies introduced

□ Domain layer unchanged unless required

---

# Business Checklist

□ Business rules correctly implemented

□ No duplicated business logic

□ Ubiquitous Language preserved

□ Acceptance criteria satisfied

---

# Code Quality Checklist

□ Naming follows conventions

□ Functions remain small

□ Complexity acceptable

□ Dead code removed

□ TODO comments removed

---

# API Checklist

□ Request validation completed

□ Response contract unchanged

□ Error handling implemented

□ Status codes correct

---

# Database Checklist

□ Migration added if required

□ Existing schema preserved

□ Indexes reviewed

□ No destructive SQL

---

# Testing Checklist

□ Unit tests added

□ Existing tests pass

□ Regression tests included for bug fixes

□ Edge cases verified

---

# Documentation Checklist

□ MODULE.md updated if behavior changed

□ WORK_PACKAGE.md updated

□ API documentation updated

□ README updated if necessary

---

# Security Checklist

□ Secrets not hardcoded

□ Sensitive data not logged

□ Validation completed

□ Authorization preserved

---

# AI Checklist

□ Context loaded correctly

□ Scope respected

□ No unrelated files modified

□ No architectural assumptions made

□ Human review ready

---

# Definition of Done

A task is complete only if:

✓ Acceptance criteria satisfied

✓ Architecture preserved

✓ Tests pass

✓ Documentation updated

✓ Checklist completed

---

# Human Review

The reviewer SHOULD verify:

- Correctness
- Simplicity
- Maintainability
- Readability
- Architectural compliance

Human approval is required before merge.

---

# Final Principle

Shipping code is easy.

Shipping maintainable software is the goal.
