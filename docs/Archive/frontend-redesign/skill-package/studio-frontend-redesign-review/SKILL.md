---
name: studio-frontend-redesign-review
description: Use when redesigning this repo's `/studio` or `/test`, comparing accepted Stitch pages against current code, or planning how page-by-page baselines will reconnect to the existing backend.
---

# Studio Frontend Redesign Review

Use this skill when the task is to rethink or regenerate the frontend for this repo, especially when the user wants:

- a new `/studio` workflow
- a redesigned `/test` engineer console
- AI-first UI generation via Stitch or Pencil
- API inventory and frontend-backend integration mapping
- a plan to reconnect generated frontend output to the current backend

## Required workflow

### 1. Read project rules first

Always read:

- `AGENTS.md`
- agent-specific file
- relevant `.github/instructions/*`

### 2. Audit current frontend and backend

Before proposing redesign, inspect:

- routes and page entrypoints
- `frontend/src/services/*.ts`
- `frontend/src/hooks/datalink/*.ts`
- backend router and major API groups
- current OpenSpec

Use the reference templates in `references/` for inventory output.

### 3. Build the current integration map

Document:

- `route -> page -> hook -> service -> API`
- `/studio` dependencies
- `/test` dependencies
- stable contracts vs likely adapter boundaries

Create or refresh reusable baseline documents for every page in scope.

For code-backed pages, record:

- current shell structure
- page-level state
- component-level state
- service and API usage
- accepted-screen-to-current-code mapping

For Stitch-only pages, record:

- accepted Stitch screen id and title
- planned route or mode
- intended domain object
- missing code owner
- missing API or adapter gap

### 4. Redesign IA and object model

Output:

- new workflow
- page map
- object model
- workspace skeletons
- readiness and handoff rules

### 5. Produce AI generation outputs

If the user wants Stitch:

- write plain-language, screen-by-screen prompts
- avoid one huge prompt
- keep one screen or one change per prompt
- if Stitch MCP is available, prefer operating on the actual Stitch project instead of stopping at prompt text
- for this repo's Stitch generation and edit work, default to `modelId: GEMINI_3_1_PRO` unless the user explicitly requests another model
- identify the target project first, then list screens, then edit one screen at a time
- treat `generate_screen_from_text` as first-pass generation and `edit_screens` as refinement
- after any interrupted or long-running Stitch edit, verify with `get_screen` and `get_project` before retrying
- do not assume an aborted edit failed; Stitch edits may finish asynchronously

If the user wants Pencil:

- write screen structure plans
- define frame regions and planning units

### 6. Plan backend reconnection

Do not stop at mockup generation.

The skill is only complete when it also specifies:

- which current APIs are reused unchanged
- which screens need adapters/view-models
- which backend gaps really require changes
- the order for reconnecting generated frontend output to current backend

### 7. Define verification

The terminal state is not "frontend generated".

The terminal state is:

`frontend generated + backend reconnected + verification completed`

Verification should include, as applicable:

- typecheck
- build
- targeted tests
- route-by-route integration sanity checks

## Output package

At minimum, produce:

1. workflow and IA spec
2. API inventory and integration map
3. `/test` current architecture baseline when `/test` is in scope
4. page-by-page baseline / comparison documents for all other pages in scope
5. Stitch prompt package when requested
6. Pencil planning package when requested
7. backend reconnection plan

## Design rules for this repo

Treat `/studio` and `/test` as separate first-class products inside one repo:

- `/studio` = workflow-first product UI
- `/test` = field engineer debug console

Keep the same design language, but allow `/test` to be denser and more tool-like.

For `/test`, do not jump straight from accepted screens to implementation sequencing.
First compare accepted screens against the current baseline and explicitly separate:

- shell problems
- adapter problems
- backend gaps

For pages that do not yet exist in code, use the accepted Stitch page as a `Stitch-planned baseline` instead of skipping the page.

## References

- For API inventory format, read `references/api-inventory-template.md`
- For frontend-backend integration mapping, read `references/integration-map-template.md`
