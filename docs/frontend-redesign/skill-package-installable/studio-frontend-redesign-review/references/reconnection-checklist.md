# Backend Reconnection Checklist

Use this after AI-generated frontend screens exist.

## 1. Freeze screen ownership first

Confirm:

- which routes belong to `/studio`
- which routes belong to `/test`
- which screen owns which domain object
- which screens are entry hubs versus workspaces

Do not wire backend calls before this is stable.

## 2. Build the screen contract map

For each screen, write:

- route
- page or workspace name
- main operator object
- required read APIs
- required write APIs
- dry-run or preview APIs
- stream or live-status APIs

Prefer this chain:

`route -> page -> adapter -> hook -> service -> API`

## 3. Add adapters before backend changes

Create adapters when the generated UI speaks a different model than the current API.

Common cases:

- source lattice view-model
- semantic refinement board model
- delivery group planner model
- database row planner projection
- local modbus register block projection
- test console session summary

## 4. Reuse existing integration layers

Prefer existing files first:

- `frontend/src/services/datalink.ts`
- `frontend/src/services/api.ts`
- `frontend/src/hooks/datalink/*`

Do not put raw fetch logic inside screen components.

## 5. Wire interaction states explicitly

Verify every important action has:

- loading state
- success feedback
- error handling
- retry behavior where needed
- dry-run or preview state where applicable
- empty state

## 6. Verify contract coverage

Check whether each screen can already be served by current backend contracts.

If not, classify the gap:

- adapter-only gap
- frontend state-management gap
- real backend contract gap

Only the third case should create a new backend change request.

## 7. Run final verification

Minimum verification should include the affected frontend checks.

Typical set:

- `git diff --check`
- `cd frontend && npm run lint`
- `cd frontend && npm run test`
- `cd frontend && npm run build`

If backend contracts changed, also run:

- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...`

## Success condition

The redesign is only complete when all of these are true:

- the generated UI direction is accepted
- route and screen ownership are stable
- the current backend contract is reconnected or explicitly amended
- verification has been run and recorded
