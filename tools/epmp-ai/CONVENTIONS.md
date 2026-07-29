---
id: EPMP-AI-CONVENTIONS
title: Engineering Conventions
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the engineering conventions used across the EPMP repository.
  These conventions ensure that all human engineers and AI agents produce
  consistent, maintainable, and predictable artifacts.

depends_on:
  - AI.md
  - RULES.md
---

# Overview

Conventions reduce cognitive load.

A consistent codebase is easier to:

- Understand
- Maintain
- Review
- Extend
- Generate using AI

Every implementation SHOULD follow these conventions unless explicitly documented otherwise.

---

# Naming Principles

Names MUST be:

- Clear
- Predictable
- Consistent
- Business-oriented

Avoid abbreviations unless universally understood.

Good

Reservation

Property

Tenant

Payment

Invoice

Bad

Res

Prop

Tbl

Obj

Manager2

HelperFinal

---

# File Naming

Markdown

```

kebab-case.md

```

Go

```

reservation_service.go

payment_repository.go

```

React

```

ReservationPage.tsx

ReservationForm.tsx

ReservationTable.tsx

```

---

# Folder Naming

Folders MUST use:

```

kebab-case

```

Example

```

reservation/

payment/

tenant/

```

---

# Package Naming (Go)

Packages MUST

- be singular
- lowercase
- concise

Good

```

reservation

payment

invoice

```

Bad

```

ReservationPackage

ReservationService

ReservationModule

```

---

# Function Naming

Functions MUST describe behavior.

Good

```

CreateReservation()

CancelReservation()

CalculateOutstandingBalance()

```

Avoid

```

Handle()

Run()

Do()

Execute()

```

unless context already provides meaning.

---

# Variable Naming

Variables SHOULD reveal intent.

Good

```

reservation

tenant

invoice

remainingBalance

paymentAmount

```

Avoid

```

x

data

obj

value

item

```

---

# Interface Naming

Interfaces describe capability.

Good

```

ReservationRepository

PaymentGateway

NotificationSender

```

Avoid

```

IReservation

BaseRepository

CommonInterface

```

---

# Error Naming

Errors SHOULD describe business failures.

Good

```

ErrReservationNotFound

ErrRoomAlreadyOccupied

ErrPaymentExpired

```

Avoid

```

ErrSomethingWrong

ErrUnknown

```

---

# DTO Naming

Request DTO

```

CreateReservationRequest

```

Response DTO

```

ReservationResponse

```

Update DTO

```

UpdateReservationRequest

```

---

# Event Naming

Domain Events MUST describe completed facts.

Good

```

ReservationCreated

ReservationCancelled

PaymentCompleted

```

Avoid

```

CreateReservation

DoPayment

```

---

# Test Naming

Test names MUST describe behavior.

```

TestCreateReservation_ShouldCreateReservation

TestCancelReservation_ShouldRejectPaidReservation

```

---

# Constant Naming

Exported

```

DefaultPageSize

MaximumReservationDays

```

Internal

```

defaultTimeout

maxRetry

```

---

# API Convention

Endpoint

```

POST /reservations

GET /reservations/{id}

PUT /reservations/{id}

DELETE /reservations/{id}

```

Plural resource names SHOULD be used.

---

# Commit Convention

Format

```

type(scope): summary

```

Example

```

feat(reservation): add create reservation

fix(payment): handle expired invoice

docs(ai): update workflow

refactor(tenant): simplify validation

```

---

# Documentation Convention

Every markdown SHOULD contain:

- Metadata
- Purpose
- Overview
- Rules
- Examples
- References

---

# AI Convention

AI responses SHOULD:

- Explain assumptions
- Keep changes minimal
- Preserve architecture
- Avoid unrelated modifications

---

# Final Principle

Consistency beats cleverness.

Readable code outlives smart code.
