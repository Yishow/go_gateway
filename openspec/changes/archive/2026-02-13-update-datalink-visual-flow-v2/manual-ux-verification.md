# Manual UX Verification

- Date: 2026-02-13
- Scope: Flow-first workspace (`SmartDashboard`) for source -> memory grid -> tag -> database flow

## Environment

- Frontend build: `npm run build` successful
- Test baseline: vitest flow and parser tests passing
- Browser profile: desktop layout simulation (>= 1024px)

## Scenarios

1. Device switch state reset
- Step: Select device A, select memory addresses, open point detail, switch to device B.
- Expected: Address selections clear, detail panel closes, flow returns to current device context.
- Result: Pass

2. Flow validation and activation
- Step: Select source context with point/tag linkage, run validate, then activate.
- Expected: Status transitions `draft -> validated -> active`.
- Result: Pass

3. Segment error and recovery
- Step: Trigger flow error state, run recovery action.
- Expected: Status transitions `error -> draft`, diagnostics reset to unknown baseline.
- Result: Pass

4. Keyboard accessibility for core actions
- Step: Trigger shortcuts (`Ctrl+K`, `Ctrl+Enter`, `Ctrl+Shift+Enter`, `Alt+R`, `Ctrl+I`, `Ctrl+E`).
- Expected: Search focus and flow actions are reachable without mouse.
- Result: Pass

5. Visible focus and contrast baseline
- Step: Tab through core action buttons and flow controls.
- Expected: Focus ring remains clearly visible on dark background; critical action labels remain readable.
- Result: Pass

## Notes

- Integration flow behavior is additionally covered by `flowLifecycle.integration.test.ts`.
- Existing project-wide warnings are not introduced by this change set.
