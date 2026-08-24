# Legacy `/studio` Immediate Deletion Record

## Decision and status

- Decision owner: project owner
- Decision date: `2026-08-24`
- Decision: delete the dedicated legacy `/studio` product surface immediately.
- Current batch: the bounded application-code deletion is present in the worktree and its focused/source/embedded checks are recorded below. The dedicated deletion commit and clean revert/rebuild/redeploy smoke remain pending (task 1.5).
- After deletion: `/studio` must have no dedicated route, handler, tombstone, compatibility flag, or special redirect. It must follow the same generic unknown-route policy as an arbitrary unknown path.
- This batch has no telemetry, database migration, usage window, 14-day wait, query migration, or 307/308 compatibility stage.

## Preserved surfaces

The following surfaces are explicitly retained and must keep their route identity, required assets, and existing product/tool semantics:

- `/studio/v2` — the only user-facing datalink setup entry.
- `/studio/runtime` — focused post-setup runtime observer.
- `/test` — engineering test tooling.
- `/gateway/*` — experimental/prototype surfaces.
- Backend datalink domain APIs and shared modules remain unless a source/import graph proves they are legacy-only. Removing the frontend route does not by itself authorize deleting backend capabilities or runtime data.

## Pre-delete source evidence

The following evidence was re-checked against the current worktree before writing this record:

| Surface | Evidence | Decision |
| --- | --- | --- |
| Dedicated route and lazy import | `frontend/src/App.tsx` pre-delete lines 17-19 import `DatalinkWorkbenchPage`; line 122 registers `path="/studio"`; lines 54-57 define `LegacyStudioRedirect`; lines 126-127 route legacy workbench/dashboard paths through it. | Delete the dedicated route, legacy page lazy import, and redirect branch only after the compatibility paths are explicitly re-routed or retired. |
| `/studio` redirect builders | `frontend/src/features/datalink/legacyRoutes.ts` pre-delete `buildWorkbenchRedirect` and `buildLocalModbusCompatRedirect` return `/studio` paths. | The functions/branches returning `/studio` are delete or rewrite candidates; the file is shared by other legacy landing paths and must not be deleted wholesale without a new reference scan. |
| Legacy page graph | `git ls-files frontend/src/pages/datalink/workbench` returned 79 tracked files. Production-source search found the directory reached from `App.tsx` only through the dedicated lazy import; no `/studio/v2`, runtime, test, or gateway production import was found. | Candidate for deletion as a bounded legacy-only graph. Final build/import graph must confirm no preserved surface reaches it. |
| Legacy unit-test graph | `git ls-files frontend/tests/unit/pages/datalink/workbench` returned 66 tracked tests. They import the legacy page/components or exercise its workbench contracts. | Candidate for deletion or narrow replacement; do not remove preserved V2/runtime assertions merely because their fixtures mock the old page. |
| Cross-route fixtures | `frontend/tests/unit/workbench-v2/routing.test.tsx`, `frontend/tests/unit/app-routing-lazy-load.test.tsx`, `frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx`, and `frontend/tests/unit/utils/appAgentationRemoval.test.tsx` reference the legacy page/mock. | Retain the tests that cover `/studio/v2` or `/studio/runtime`; remove only obsolete legacy mocks/assertions and add generic unknown-route equivalence coverage. |
| Backend route scan | Current Go/API/static source scan found no exact dedicated backend `/studio` route or handler registration. The existing backend datalink APIs are domain capabilities, not proof of a frontend owner. | Preserve backend APIs unless a separate source review proves them unused by retained surfaces. |
| Embedded delivery | Tracked `cmd/test_ui/static` contains only `embed-placeholder.txt`; build output is copied from `frontend/dist` by `scripts/build.ps1`. | Do not hand-delete static files. Rebuild, inspect the generated asset graph, and remove only orphaned legacy chunks/assets proven absent from retained routes. |

## Delete candidates

These are candidates supported by the source evidence above, not permission to delete unreviewed files:

1. The `/studio` `Route`, `DatalinkWorkbenchPage` lazy import, `LegacyStudioRedirect`, and route-level comments in `frontend/src/App.tsx`.
2. The `/studio`-returning branches in `frontend/src/features/datalink/legacyRoutes.ts`, after each caller is mapped. Keep shared V2 and preserved compatibility helpers.
3. The tracked production graph under `frontend/src/pages/datalink/workbench/**` (79 files), subject to final import/asset graph verification.
4. The legacy-only tests under `frontend/tests/unit/pages/datalink/workbench/**` (66 files), subject to retaining any assertions that are actually preserved-surface contracts.
5. Obsolete legacy route mocks/assertions in the four cross-route test files listed above.
6. Generated legacy chunks/assets found only after a clean frontend build; never infer this list from filenames alone.

## Explicitly retained or uncertain

- `frontend/src/pages/datalink/workbench-v2/**` and `frontend/src/features/datalink/workbench-v2/**` are retained for `/studio/v2`.
- `frontend/src/features/datalink/runtime-dashboard/**` is retained for `/studio/runtime`.
- `frontend/src/pages/TestPage.tsx`, `frontend/src/pages/TestPageShell.tsx`, and `/test` support are retained.
- `frontend/src/router/gateway/**` and gateway feature-flag wiring are retained.
- `frontend/src/features/datalink/legacyRoutes.ts` is shared/uncertain until all callers are re-mapped; only its `/studio`-specific branches are candidates.
- Shared services, i18n keys, CSS, API clients, backend handlers, schema/migrations, and runtime collectors are retained unless a separate reachability review proves a component legacy-only.
- `docs/technical/studio-surface-inventory/studio-mainline.md` remains as historical pre-delete evidence. It is not a current product contract and must not be used to reintroduce `/studio`.
- Existing runtime/deployment data and external DB/device state are retained and are not modified by this deletion batch.

## Rollback boundary

The deletion must be delivered as a dedicated Git commit referring to this record. Rollback is:

```text
git revert <legacy-deletion-commit>
pwsh -NoProfile -File scripts/build.ps1
redeploy the rebuilt bin/test-ui.exe using the deployment owner procedure
repeat route, static-asset, and browser smoke checks
```

Git rollback can restore tracked source/docs and assets that are regenerated and redeployed from tracked source. It cannot restore runtime/deployment data, external database contents, already-issued device state, browser bookmarks, localStorage, or an already-running process without a rebuild/redeploy/restart. Those require the relevant deployment or data owner procedure.

## Immediate deletion evidence (2026-08-24)

The implementation worker's bounded diff contains **160 proven legacy-only file deletions**. The deletion inventory was re-counted from `git diff --name-only --diff-filter=D` and grouped as follows:

| Deletion group | Files |
| --- | ---: |
| `frontend/src/pages/datalink/workbench/**` | 79 |
| `frontend/tests/unit/pages/datalink/workbench/**` | 66 |
| `frontend/src/features/datalink/` legacy-only contracts/storage | 2 |
| `frontend/tests/unit/features/datalink/` legacy-only contracts | 4 |
| Legacy-only workbench styles and adjacent model/guard tests | 9 |
| **Total** | **160** |

The source graph now removes the dedicated `/studio` route, the legacy page lazy import, and the `/studio`-specific redirect branch. Generic `*` fallback remains the single route for `/studio` and arbitrary unknown paths; generic datalink landing routes continue to converge on `/studio/v2`. No backend `/studio` route was introduced or changed by this deletion batch.

Recorded validation evidence:

- Worker full frontend run: **119 files / 639 tests passed**.
- Primary focused route/deletion run: **4 files / 17 tests passed**.
- Frontend lint: passed.
- Frontend build script: passed.
- Embedded Playwright smoke: **8/8 passed**.
- File line gate: passed.
- `git diff --check`: passed.
- Static graph check: no legacy chunk token remained; `/studio/v2`, `/studio/runtime`, `/test`, and `/gateway/*` remained present and preserved.
- An initial sandbox `spawn EPERM` was rerun outside the sandbox and succeeded; it was an execution-environment issue, not a product failure.

The evidence above proves tasks 1.1, 1.2, 1.3, 1.4, and 1.6 only. It does not claim the dedicated Git commit, revert, rebuild, redeploy, or full C-polish acceptance.

## Dependency reachability evidence (2026-08-24)

The same legacy-only cleanup also removes four unused direct frontend dependencies from `frontend/package.json` and both tracked lockfiles (`frontend/package-lock.json` and `frontend/pnpm-lock.yaml`):

- `@mui/material`
- `@mui/icons-material`
- `@emotion/react`
- `@emotion/styled`

Reachability checks found no remaining consumer or string for these four packages in `frontend/src`, `frontend/tests`, `frontend/package.json`, `frontend/package-lock.json`, or `frontend/pnpm-lock.yaml`. The package-lock `dev` flag insertions are dependency-reachability reclassification after the production dependency roots were removed; retained packages did not undergo version, resolved URL, or integrity-hash churn.

Recorded dependency and regression validation:

- npm offline dry run: `npm install --package-lock-only --dry-run --offline --ignore-scripts --no-audit --no-fund` passed.
- pnpm offline frozen-lockfile check: `pnpm install --offline --frozen-lockfile --ignore-scripts --lockfile-only` passed.
- Default frontend run after cleanup: `npm --prefix frontend run test -- --run` — **119 files / 639 tests passed in 42.07s**, with no timeout or watchdog result.
- Focused route/deletion run: **4 files / 17 tests passed**.
- Frontend lint: passed.
- Frontend build: passed after rerunning outside the restricted sandbox; the initial sandbox `spawn EPERM` was an execution-environment limitation.
- File line gate: passed.
- `git diff --check`: passed after dependency cleanup.

## Acceptance checklist for the implementation worker

- [x] No dedicated frontend/backend `/studio` route, handler, tombstone, special redirect, or hidden lazy import remains.
- [x] `/studio` and an arbitrary unknown path have equivalent navigation/status/fallback/error-boundary and network-asset behavior.
- [x] No legacy-only chunk/asset/request is loaded by `/studio/v2`, `/studio/runtime`, `/test`, or `/gateway/*`.
- [x] Preserved route identity and focused smoke checks pass for `/studio/v2`, `/studio/runtime`, `/test`, and `/gateway/*`.
- [x] Frontend tests are updated with exact retained-surface and unknown-route evidence.
- [ ] The deletion commit, revert command, rebuild/redeploy evidence, and rollback limits are recorded.
