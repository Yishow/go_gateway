---
name: studio-frontend-redesign-review
description: Review and redesign this repo's frontend end-to-end: audit current routes, APIs, and frontend-backend integration; redefine product IA for /studio and /test; generate Google Stitch and Pencil MCP outputs; and carry the work through until the AI-generated frontend is reconnected to the current backend contract and verified.
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
3. Stitch prompt package when requested
4. Pencil planning package when requested
5. backend reconnection plan

## Design rules for this repo

Treat `/studio` and `/test` as separate first-class products inside one repo:

- `/studio` = workflow-first product UI
- `/test` = field engineer debug console

Keep the same design language, but allow `/test` to be denser and more tool-like.

## References

- For API inventory format, read `references/api-inventory-template.md`
- For frontend-backend integration mapping, read `references/integration-map-template.md`
