# Project Governance & Documentation Rules

This rule dictates how AI agents and developers must operate regarding project specifications, technical documentation, and functional traceability.

## 1. Single Source of Truth
All code implementations, refactoring, or additions **MUST** refer to the approved specifications located in the following folders:
- `/docs/PRD/` (Product Requirements Document)
- `/docs/TSD/` (Technical Specification Document)

It is strictly prohibited to make any architectural or functional assumptions that contradict the contents of these documents without updating them first.

## 2. Scope Enforcement
If the user requests features or code changes that are NOT present in the PRD/TSD documents, you must:
1. Politely reject the request.
2. Inform the user that the feature is "Out of Scope".
3. Suggest the user run the `/create-spec` workflow first.

## 3. Directus Schema Management (Schema-as-Code)
Every implemented change regarding the database structure, field configurations, Directus Admin UI Panels, or table relations **MUST** be recorded in detail, including its Changelog, inside:
- `docs/schema.md`

Do not execute schema modification scripts (like `setup-survey-schema.js`) if the documentation does not reflect the latest logic to be released.

## 4. Requirement ID Traceability
Every new feature, bug fix, or module implementation must refer to a specific `Requirement ID` (e.g., `REQ-001`, `FR1`, or `NFR2` from the PRD).

Always mention this ID in your reasoning, pull requests, commit messages, or release notes when explaining task resolutions to the User. This guarantees that the written code directly answers the functional/non-functional needs of the project.
