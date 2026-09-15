---
description: Requirement Gathering and Specification update workflow for out-of-scope requests
---

# /create-spec Workflow

When the user runs the `/create-spec` slash command, follow these exact steps to formally draft and integrate new feature requests into the official specification documents (PRD and TSD). Do not proceed to write application code during this workflow.

## Step 1: Gather Requirements
1. Ask the user to comprehensively describe the new feature, user stories, and business logic they want to add to the module.
2. Ask if there are any specific UI/UX dependencies or API constraints.
3. Stop and **wait for the user's response**.

## Step 2: Analyze Impact
Once the user provides the requirements:
1. Analyze how this change affects the existing PostgreSQL Database Schema / Directus Collections (`docs/schema.md`).
2. Analyze how this affects the Directus JSON Extension API adapter (`extensions/survey-forms/index.js`).
3. Propose a list of Functional Requirements (FR) and Non-Functional Requirements (NFR). 
4. Assing strict `Requirement IDs` for each new feature (e.g., `REQ-002`, `FR-11`).

## Step 3: Propose Implementation Plan
1. Create an `implementation_plan.md` artifact detailing how the PRD and TSD will be modified to accommodate the new specs. 
2. Set `request_feedback = true` and yield control. **Wait for user approval**.

## Step 4: Update Documentation (PRD & TSD)
Once the user approves the plan:
1. Modify the official PRD (`/docs/PRD/PRD_v1.0.0.md` or the latest iteration) and TSD (`/docs/TSD/TSD_v1.0.0.md`).
2. Append the new `Requirement ID`s, update the System Scope, and detail the newly defined API attributes or collections.
3. Update the "Versi Dokumen" / Version History table at the top of both PRD and TSD to reflect the version bump (e.g., bumping to `v1.1.0`).
4. Congratulate the user and remind them that the specifications are now officially "In Scope". Encourage them to use standard workflows (like `/2-implement` or `/orchestrator`) to begin the actual code execution.
