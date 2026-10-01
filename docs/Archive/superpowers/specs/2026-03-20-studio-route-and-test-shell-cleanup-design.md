# Studio route and test shell cleanup

## Goal

Collapse the old multi-entry datalink and test-tool routing into two clear surfaces:

- `/studio` for the operator-facing source-to-output workflow
- `/test` for the engineering test page

## Approved decisions

1. Replace the public datalink entry with `/studio` instead of keeping `datalink` or `workbench` in the URL.
2. Keep legacy datalink URLs as redirects where practical, so old links do not hard-break.
3. Remove the old test-tool sidebar shell and keep `/test` as a single page without extra navigation.
4. Remove clearly obsolete legacy test-tool pages (`/templates`, `/history`, `/compare`, `/analyzer`) and redirect their old URLs to `/test`.

## Route design

- `/` redirects to `/studio`
- `/studio` renders the existing workbench flow
- `/datalink`, `/datalink/workbench`, and local-modbus compat routes redirect into `/studio`
- `/test` renders the standalone test shell
- old test utility routes redirect into `/test`

## UI shell design

`/test` should no longer inherit the old application sidebar. Instead, it gets a page-owned shell that:

- keeps the current TestPage workflow intact
- provides full-width, no-sidebar presentation
- widens the content container for 1920 desktop usage while preserving responsive padding

## Cleanup scope

- remove obsolete legacy test-tool page files and route wiring
- remove the legacy local modbus standalone page and its page-level test
- remove the deeper SmartDashboard implementation tree, related tests, and now-unused helper files once they are confirmed to be detached from the new `/studio` flow

## Verification

- route foundation tests cover `/studio`, `/datalink -> /studio`, `/test`, and legacy test-route redirects
- frontend validation runs targeted tests first, then lint, type-check, and build
