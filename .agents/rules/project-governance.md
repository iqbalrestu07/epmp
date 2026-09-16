---
trigger: always_on
---

# Project Governance & Documentation Rules

This rule dictates how AI agents and developers must operate regarding project specifications, technical documentation, and functional traceability.

## 1. Single Source of Truth

All code implementations, refactoring, or additions **MUST** refer to the approved specifications:

- `epmp-docs/epmp-001.md` … `epmp-docs/epmp-013-notification.md` — product, architecture, domain, module, and standard documents (see `README.md` § Documentation Structure)
- `epmp-docs/epmp-010.md` — Engineering Platform Standard (stack, branches, commits, quality gates)
- `epmp-docs/epmp-011.md` — Backend Architecture Standard
- `tools/epmp-ai/RULES.md`, `CONVENTIONS.md`, `WORKFLOW.md` — AI Operating System
- `backend/internal/modules/*/MODULE.md` — per-module field/behaviour contracts

It is strictly prohibited to make architectural or functional assumptions that contradict these documents without updating them first.

## 2. Scope Enforcement

If the user requests features or behaviour that are NOT described in `epmp-docs/` or an existing `MODULE.md`, you must:

1. Say so explicitly and label the request "Out of Scope" relative to current specs.
2. Offer to run the `/create-spec` workflow to document it first (spec + ADR if architectural).
3. Proceed with implementation only after the user confirms (a human owns requirements — `tools/epmp-ai/RULES.md` BR-004).

## 3. Schema-as-Code

Every database structure change **MUST** be a migration pair in `backend/migrations/` and be reflected in `MODULE.md`, `backend/migrations/README.md`, and the generator schema — see @schema-documentation-mandate.md. Never change the schema by hand or via GUI tools.

## 4. Traceability

Every feature, bug fix, or module change must reference where it is specified:

- the `epmp-docs/epmp-0xx.md` document / section, or
- the `MODULE.md` of the module, or
- the ADR in `epmp-docs/adr/` that introduced it.

Mention this reference in your reasoning, commit body, or PR description so that code is traceable to an approved specification. When the spec has requirement IDs (e.g. `EP-01`, `AR-004`, `MR-004`), cite them.

## 5. Documentation Currency

Every meaningful behaviour change MUST update documentation in the same change (EPMP-010 §8, `tools/epmp-ai/RULES.md` §9): `MODULE.md`, the relevant `epmp-docs` document, `PROGRESS.md` when a feature milestone is reached, and `E2E_TESTING.md`/`PAGES_TO_TEST` when routes change.
