---
name: frontend-redesign-integration-review
description: Use when redesigning a frontend, comparing screen directions, or planning how current and Stitch-generated pages will reconnect to an existing backend with durable page-by-page baselines.
---

# Frontend Redesign Integration Review

Use this workflow when the task is to redesign a frontend and reconnect it to an existing backend.

## Workflow

1. Read repo rules first.
2. Audit current routes, pages, services, hooks, backend APIs, and current product specs.
3. Build the current integration map: `route -> page -> hook -> service -> API`.
4. Create or refresh page-by-page baseline documents before redesigning.
5. Redefine workflow, IA, screen map, workspace skeletons, and readiness rules.
6. Generate AI-friendly outputs for Stitch and/or Pencil.
7. Plan reconnection using adapters before proposing backend changes.
8. End only when generated frontend output has a concrete reconnection and verification path.

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
4. Page-by-page baseline and comparison documents for every page in scope
5. Stitch prompt package when requested
6. Pencil planning package when requested
7. Backend reconnection plan

## Core rules

- Preserve backend contracts first.
- Reuse existing frontend integration layers before inventing new ones.
- Add adapters or view-models before changing backend APIs.
- Do not stop at mockups or prompts.
- For independent tool consoles, distinguish shell problems from missing capability by documenting the current baseline first.
- For pages that do not yet exist in code, use the accepted Stitch page as a `Stitch-planned baseline` instead of skipping documentation.

## Audit sequence

Always inspect these categories:

- repo rules and agent-specific rules
- route entrypoints
- page entrypoints
- frontend services
- data hooks or query hooks
- backend router
- current specs or design docs

When tool or utility routes are in scope, also inspect:

- page-level state ownership
- component-level state ownership
- shell structure
- diagnostics capabilities
- overlaps between tool modes

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

When redefining an engineer console or utility route, decide in this order:

1. Session or context shell
2. Setup area
3. Active workspace modes
4. Diagnostics area
5. State ownership and adapters

## Stitch guidance

When the user wants Google Stitch output:

- Start with plain language.
- Generate one screen at a time.
- Change one major thing at a time.
- Do not put the entire product into one prompt.
- Always produce Stitch prompts in two parts:
  - a short preamble that says this is an existing product redesign
  - the screen-specific prompt body
- For follow-up prompts, tell Stitch to keep the accepted direction and focus only on one screen or one refinement request.

When Stitch MCP is available:

- Prefer operating on the actual Stitch project instead of stopping at prompt files.
- Find the correct project with `list_projects`, then confirm it with `get_project`.
- Use `list_screens` to identify the current screen set.
- Use `get_screen` before editing so the edit request is anchored to the current page.
- Use one focused `edit_screens` request per screen and per design goal.
- Use `generate_screen_from_text` for first-pass generation and `edit_screens` for iterative refinement.
- After interrupted or slow Stitch operations, re-check with `get_screen` and `get_project` before retrying.
- Do not assume an aborted Stitch operation failed; remote generation may continue after the client stops waiting.
- If no fresh screenshot can be inspected, report the result as unverified rather than claiming visual compliance.

When refining tool-console screens:

- Use the current architecture baseline to decide whether the blocker is shell design, state ownership, or missing backend capability.
- Do not mark the design implementation-ready until the existing component and API mapping has been compared.

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

## Baseline protocol

Produce or refresh durable baseline documents for every page in scope.

For code-backed pages, record:

- current routes
- page entrypoints
- shell structure
- component composition
- page-level state
- component-level state
- service and API usage
- accepted-screen-to-current-code mapping when Stitch outputs exist

For Stitch-only pages, record:

- accepted Stitch screen id and title
- planned route or mode
- intended domain object
- missing code owner
- missing API or adapter gap

When tool or utility routes are in scope, also split redesign blockers into:

- shell
- adapter
- backend

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
