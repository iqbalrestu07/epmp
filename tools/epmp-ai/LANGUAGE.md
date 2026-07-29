---
id: EPMP-AI-LANGUAGE
title: Ubiquitous Language
version: 1.0.0
status: Stable
owner: EPMP Architecture Team

purpose: >
  Defines the official business vocabulary used across the EPMP repository.
  All documentation, code, APIs, database objects, and AI responses MUST
  use these terms consistently.

depends_on:
  - RULES.md
  - CONVENTIONS.md
---

# Overview

Software quality depends on shared language.

When different people use different words for the same concept,
confusion follows.

The purpose of this document is to establish a single vocabulary shared by:

- Business
- Architects
- Developers
- QA
- AI Agents

This vocabulary becomes the source of truth.

---

# Core Principle

One Business Concept

↓

One Official Name

↓

Everywhere

---

# Naming Rules

Each business concept MUST have exactly one official term.

Example

GOOD

Reservation

Reservation API

ReservationRepository

ReservationCreated

reservation_id

BAD

Booking

Reservation

Reserve

ReserveRequest

BookingRepository

ReservationAPI

These all refer to the same concept and MUST NOT coexist.

---

# Domain Terms

## Reservation

Official Meaning

A customer's commitment to occupy a property during a defined period.

Allowed

Reservation

Forbidden

Booking

Order

Schedule

---

## Property

Official Meaning

A rentable physical asset.

Allowed

Property

Forbidden

House

Building

Room

unless explicitly modeled as different domain concepts.

---

## Tenant

Official Meaning

A customer renting a property.

Forbidden

User

Customer

Member

unless those represent different business entities.

---

## Invoice

Official Meaning

A financial document requesting payment.

Forbidden

Bill

Receipt

Payment Sheet

---

## Payment

Official Meaning

A completed financial transaction.

Payment is NOT an Invoice.

Payment is NOT Billing.

---

# Code Convention

Business names MUST propagate consistently.

Example

Entity

```
Reservation
```

Repository

```
ReservationRepository
```

DTO

```
CreateReservationRequest
```

API

```
POST /reservations
```

Database

```
reservation
```

Event

```
ReservationCreated
```

---

# Documentation Convention

Documentation MUST use official terminology.

Never mix synonyms.

---

# API Convention

Resource names MUST match business language.

GOOD

```
/reservations
```

BAD

```
/booking
/reserve
/orders
```

---

# Database Convention

Table names SHOULD match domain language.

GOOD

```
reservation
tenant
property
invoice
```

BAD

```
booking_tbl
customer_master
```

---

# AI Behavior

When encountering synonyms:

AI MUST translate them into the official business language before implementation.

Example

Human Request

```
Create booking feature
```

Internal Interpretation

```
Reservation
```

Generated Code

```
ReservationService
ReservationRepository
ReservationCreated
```

---

# Introducing New Terms

A new business term MAY only be introduced when:

- It represents a genuinely new business concept.
- It is approved by the Architecture Team.
- It is documented here.

---

# Final Principle

Shared language creates shared understanding.

Shared understanding creates consistent software.
