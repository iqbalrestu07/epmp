# Document 001

```
Document ID : EPMP-README-001
Title       : Enterprise Property Management Platform (EPMP)
Version     : 0.1.0
Status      : Draft
Owner       : Product & Architecture
Last Update : 2026-07-29
```

---

# Enterprise Property Management Platform (EPMP)

## Welcome

Welcome to the **Enterprise Property Management Platform (EPMP)** documentation.

EPMP is a modular, scalable, AI-friendly platform designed to manage various kinds of rental properties through a single unified system.

Unlike traditional property management software that is built for one specific business model (boarding house, apartment, hotel, warehouse, etc.), EPMP is designed as a **general-purpose platform** whose behavior is driven by configuration rather than hardcoded business logic.

The long-term vision is to build a platform that can support thousands of properties, multiple organizations, and multiple business models without requiring changes to the core application.

---

# Vision

To become a flexible, extensible, and enterprise-grade Property Management Platform that enables any property rental business to digitize, automate, and scale its operations through configurable business rules and modular architecture.

---

# Mission

EPMP aims to provide a modern software platform that allows property owners and operators to manage:

- Properties
- Buildings
- Floors
- Units / Rooms
- Tenants
- Reservations
- Contracts
- Payments
- Assets
- Maintenance
- Reports

through one integrated system.

The platform should minimize hardcoded business logic and maximize flexibility through configuration.

---

# Product Philosophy

EPMP is built upon one fundamental philosophy.

> **Everything should be configurable. Nothing should be hardcoded.**

Every module must follow this philosophy.

Business requirements should be represented as configuration whenever possible instead of application logic.

---

# What EPMP is NOT

EPMP is **not**:

- a Boarding House Application
- an Apartment Application
- a Hotel Application
- a Dormitory Application

Instead,

those are merely different implementations of the same platform.

---

# Supported Business Types

The platform is designed to support (but is not limited to):

- Boarding House (Kost)
- Apartment
- Dormitory
- Student Housing
- Co-Living
- Guest House
- Villa
- Office Rental
- Commercial Building
- Warehouse
- Storage Unit
- Serviced Apartment
- Mixed Property

without changing the core architecture.

---

# Core Concepts

EPMP is organized around several core business domains.

```
Organization
    │
    ├── Property
    │
    ├── Building
    │
    ├── Floor
    │
    ├── Zone
    │
    ├── Room / Unit
    │
    ├── Bed (optional)
    │
    ├── Tenant
    │
    ├── Reservation
    │
    ├── Contract
    │
    ├── Invoice
    │
    ├── Payment
    │
    ├── Asset
    │
    ├── Maintenance
    │
    └── Reports
```

Each domain is independent and connected through clearly defined relationships.

---

# Design Principles

EPMP follows the following principles.

## 1. Configuration First

Every configurable business rule should live inside configuration.

Never inside application code.

---

## 2. Modular Architecture

Every module has a single responsibility.

Modules should communicate through clearly defined interfaces.

---

## 3. Domain Driven

Business domains define the software architecture.

Not UI.

Not Database.

---

## 4. API First

Every feature available in the UI must also be available through API.

---

## 5. AI First Documentation

Documentation is written not only for humans.

Documentation is also optimized for AI agents.

---

## 6. Event Driven

Important business activities generate business events.

Examples:

- TenantCheckedIn
- ContractActivated
- InvoiceIssued
- PaymentReceived
- DepositReturned
- AssetAssigned

---

## 7. Audit Everything

Every important business operation must be traceable.

No data should silently disappear.

---

## 8. Scalable by Design

Every architecture decision should assume future growth.

Current implementation may manage one property.

Future implementation should support thousands.

---

# Documentation Structure

The documentation repository is organized as follows.

```
00-project/
    Product documentation

01-ai-context/
    AI optimized documentation

02-business/
    Business rules

03-architecture/
    System architecture

04-domain/
    Business entities

05-modules/
    Functional specifications

06-ui/
    UI specifications

07-api/
    API specifications

08-database/
    Database specifications

09-development/
    Engineering guides

10-roadmap/
    Product roadmap

11-adr/
    Architecture Decision Records

12-rfc/
    Request For Change
```

---

# Reading Order

New contributors should read the documentation in the following order.

```
README

↓

Project Overview

↓

Product Principles

↓

Architecture Overview

↓

Domain Model

↓

Business Rules

↓

Module Specifications

↓

Database

↓

API

↓

Development Guide
```

---

# Documentation Rules

Each document should:

- Have a unique document ID.
- Have a version number.
- Have a clear owner.
- Define its dependencies.
- Focus on a single responsibility.
- Avoid duplicated information.
- Reference other documents instead of repeating them.

---

# Project Goals

The first version of EPMP aims to provide:

- Property Management
- Building Management
- Floor Management
- Room / Unit Management
- Tenant Management
- Reservation Management
- Contract Management
- Billing
- Payment Recording
- Deposit Management
- Asset Management
- Maintenance Management
- Dashboard
- Reports

through one integrated platform.

---

# Long-Term Goals

Future versions of EPMP may include:

- Mobile Applications
- Tenant Portal
- Owner Portal
- Payment Gateway
- Smart Lock Integration
- IoT Integration
- Visitor Management
- AI Analytics
- Dynamic Pricing Engine
- Public API
- Plugin Marketplace

without redesigning the core architecture.

---

# Guiding Principle

Every architecture decision should answer the following question.

> **Can this feature be configured instead of being hardcoded?**

If the answer is yes,

configuration should always be preferred.

---

# Closing Statement

EPMP is intended to become more than a property management application.

It is designed to become a configurable enterprise platform capable of supporting diverse property rental businesses while remaining maintainable, extensible, and understandable by both software engineers and AI-assisted development tools.

---
