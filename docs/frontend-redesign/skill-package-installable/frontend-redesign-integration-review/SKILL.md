---
name: frontend-redesign-integration-review
description: Redesign a frontend for any repo end-to-end: audit routes, screens, backend APIs, and frontend-backend integration; redefine IA and workflow; generate Google Stitch and Pencil MCP outputs; and continue until the AI-generated frontend is reconnected to the current backend contract and verified. Use when rethinking product flow, comparing frontend directions, preparing AI-generated UI, or planning how redesigned screens will reconnect to an existing backend.
---

# Frontend Redesign Integration Review

Use this workflow when the task is to redesign a frontend and reconnect it to an existing backend.

## Workflow

1. Read repo rules first.
2. Audit current routes, pages, services, hooks, backend APIs, and current product specs.
3. Build the current integration map: `route -> page -> hook -> service -> API`.
4. Redefine workflow, IA, screen map, workspace skeletons, and readiness rules.
5. Generate AI-friendly outputs for Stitch and/or Pencil.
6. Plan reconnection using adapters before proposing backend changes.
7. End only when generated frontend output has a concrete reconnection and verification path.

## Required inputs

- repo root
- current product routes
- current tool or utility routes if any
- redesign goal
- whether to include Stitch, Pencil, or both

## Required outputs

Produce all of the following unless the user narrows scope:

1. Workflow and IA spec
2. API inventory
3. Frontend-backend integration map
4. Stitch prompt package when requested
5. Pencil planning package when requested
6. Backend reconnection plan

## Core rules

- Preserve backend contracts first.
- Reuse existing frontend integration layers before inventing new ones.
- Add adapters or view-models before changing backend APIs.
- Do not stop at mockups or prompts.

## Audit sequence

Always inspect these categories:

- repo rules and agent-specific rules
- route entrypoints
- page entrypoints
- frontend services
- data hooks or query hooks
- backend router
- current specs or design docs

Use the templates in `references/` when writing inventory and integration outputs.

## Redesign sequence

When redefining product flow, decide in this order:

1. Product areas and route ownership
2. Main workflow steps
3. Core object model
4. Screen map
5. Workspace skeletons
6. Readiness and handoff rules
7. Destination or tool-specific planners

## Stitch guidance

When the user wants Google Stitch output:

- Start with plain language.
- Generate one screen at a time.
- Change one major thing at a time.
- Do not put the entire product into one prompt.

## Pencil guidance

When the user wants Pencil MCP output:

- Plan screens first.
- Define region structure before visual polish.
- Keep frame hierarchy explicit.
- Reserve clear adapter or view-model attachment points.

## Backend reconnection protocol

Do not treat generated frontend output as integrated by default.

Reconnect in this order:

1. Fix route and screen ownership.
2. Map each screen to domain objects.
3. Add adapters or view-models.
4. Reuse existing services and hooks.
5. Wire mutations, dry-run states, retries, and errors.
6. Run verification.

The success condition is:

`frontend generated + backend reconnected + verification completed`

Read `references/reconnection-checklist.md` when moving from generated layouts into actual wiring.

## If the generated UI is not good enough

Use this escalation order.

### 1. Change prompts for structural or stylistic problems

Change prompts when the issue is:

- wrong IA
- wrong hierarchy
- wrong page composition
- wrong visual direction
- wrong operating tone

Do not try to patch these only in the generated output.

### 2. Refine generated output for local interaction problems

Refine the generated output when the issue is:

- spacing
- labels
- panel ordering
- local controls
- specific card or table content

Use focused follow-up prompts or local design edits.

### 3. Patch code only after the direction is stable

Only move to direct frontend implementation after:

- page map is accepted
- screen structure is accepted
- backend integration path is clear

## References

- Read `references/api-inventory-template.md` to document backend endpoints and current screen consumers.
- Read `references/integration-map-template.md` to document how generated UI reconnects to the current backend.
- Read `references/iteration-rules.md` when deciding whether to revise prompts, generated layouts, or implementation code.
- Read `references/reconnection-checklist.md` when turning an accepted generated direction into a real integrated frontend.
