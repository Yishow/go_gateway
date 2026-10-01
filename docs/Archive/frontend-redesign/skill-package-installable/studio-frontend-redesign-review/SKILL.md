---
name: studio-frontend-redesign-review
description: Use when redesigning this repo's `/studio` or `/test`, comparing frontend directions, or planning how accepted Stitch or Pencil pages will reconnect to the existing backend with durable page-by-page baselines.
---

# Studio Frontend Redesign Review

Follow this workflow when the task is to redesign the frontend for this repo.

## Workflow

1. Read repo rules first.
2. Audit current routes, pages, services, hooks, backend APIs, and current OpenSpec.
3. Build the current integration map: `route -> page -> hook -> service -> API`.
4. Create or refresh page-by-page baseline documents before redesigning.
5. Redefine workflow, IA, object model, workspace skeletons, and readiness rules.
6. Generate AI-friendly outputs for Stitch and/or Pencil.
7. Plan backend reconnection using adapters before proposing backend changes.
8. End only when generated frontend output has a concrete reconnection and verification path.

## Required inputs

- Repo root
- Current product routes
- Current engineering-tool routes
- Redesign goal
- Whether to include Stitch, Pencil, or both

## Required outputs

Produce all of the following unless the user narrows scope:

1. Workflow and IA spec
2. API inventory
3. Frontend-backend integration map
4. Page-by-page baseline and comparison documents for every page in scope
5. Stitch prompt package when requested
6. Pencil planning package when requested
7. Backend reconnection plan

## Repo-specific rules

- Treat `/studio` as the product workflow.
- Treat `/test` as an independent field engineer debug console.
- Do not rely only on generated screens. First document current or planned page ownership as reusable baselines.
- Preserve backend contracts first.
- Reuse `frontend/src/services/*.ts` and `frontend/src/hooks/datalink/*.ts`.
- Add adapters or view-models before changing backend APIs.
- Do not stop at mockups or prompts.

## Audit sequence

Always inspect these categories:

- `AGENTS.md`, agent-specific file, `.github/instructions/*`
- route entrypoints
- page entrypoints
- frontend services
- TanStack Query hooks
- backend router
- current OpenSpec

When `/test` is in scope, also inspect:

- page-level state ownership
- component-level state ownership
- diagnostics capabilities
- tool-mode overlap such as command, monitor, scan, and polling

When destination-family pages are in scope, also inspect:

- whether the page already exists in code or only in Stitch
- whether current ownership is a dedicated route or still collapsed under `output`
- which parts are already backed by `dbTargetAPI` or `modbusShareAPI`
- whether MQTT is still design-only

Use the templates in `references/` when writing inventory and integration outputs.

## Redesign sequence

When redefining the product flow, decide in this order:

1. Product split between `/studio` and `/test`
2. Main workflow steps
3. Core object model
4. Screen map
5. Workspace skeletons
6. Readiness and handoff rules
7. Destination-specific planners

When redefining an independent tool console such as `/test`, decide in this order:

1. Session shell
2. Setup rail
3. Active workspace modes
4. Diagnostics rail
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

- Prefer direct Stitch project operations over prompt-only output.
- For this repo's Stitch generation and edit work, default to `modelId: GEMINI_3_1_PRO` unless the user explicitly asks for a different model.
- Resolve the correct project first with `list_projects`.
- Confirm the target project with `get_project`.
- Inspect the current screen set with `list_screens`.
- Read the target screen with `get_screen` before editing.
- Edit only one screen per `edit_screens` call unless the user explicitly wants a coordinated multi-screen change.
- Keep each edit request focused on one structural goal, such as shell, device workbench, or diagnostics rail.
- Treat `generate_screen_from_text` as first-pass generation and `edit_screens` as follow-up refinement.
- After any long-running or interrupted Stitch operation, re-check `get_screen` and `get_project` before retrying.
- Do not assume an aborted Stitch edit failed. Stitch may still complete the change asynchronously.
- When visual inspection is unavailable, do not claim the result is correct based only on project update time. Mark it as pending screenshot review.

When refining `/test` screens:

- Use the current architecture baseline to distinguish shell problems from missing backend capability.
- Avoid claiming a `/test` screen is implementation-ready until its current component and state ownership gaps are explicitly mapped.

## Pencil guidance

When the user wants Pencil MCP output:

- Plan screens first.
- Define region structure before visual polish.
- Keep frame hierarchy explicit.
- Reserve clear adapter/view-model attachment points.

## Backend reconnection protocol

Do not treat generated frontend output as integrated by default.

Reconnect in this order:

1. Fix route and screen ownership.
2. Map each screen to domain objects.
3. Add adapters/view-models.
4. Reuse existing services and hooks.
5. Wire mutations, dry-run states, and retries.
6. Run verification.

The success condition is:

`frontend generated + backend reconnected + verification completed`

Read `references/reconnection-checklist.md` when you are moving from generated layouts into actual frontend-backend wiring.

## Baseline protocol

Produce or refresh durable baseline documents for every page in scope.

For code-backed pages, record:

- current routes
- page entrypoints
- current shell structure
- current component composition
- page-level state
- component-level state
- service and API usage
- accepted-screen-to-current-code mapping

For Stitch-only pages, record:

- accepted Stitch screen id and title
- planned route or mode
- intended domain object
- current missing code owner
- current missing API or adapter gap

If `/test` is included in scope, explicitly split redesign blockers into:

- shell problems
- adapter problems
- backend gaps

Do this before proposing implementation sequencing.

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
- specific card content

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
