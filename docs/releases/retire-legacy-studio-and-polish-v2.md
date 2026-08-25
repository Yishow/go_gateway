# Legacy `/studio` Immediate Deletion Record

## Decision and status

- Decision owner: project owner
- Decision date: `2026-08-24`
- Decision: delete the dedicated legacy `/studio` product surface immediately.
- Current batch: the bounded application-code deletion was anchored by Git commit `efa355f8c7943941b246a750cf7650b5ba211f88` (`刪除 legacy Studio 產品面`). The commit records this change and the deletion batch; a clean revert/rebuild/redeploy execution remains an operator acceptance step, not an unrecorded source change.
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

The deletion is delivered as dedicated Git commit `efa355f8c7943941b246a750cf7650b5ba211f88`, which refers to this record. Rollback is:

```text
git revert efa355f8c7943941b246a750cf7650b5ba211f88
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
- [x] The deletion commit, revert command, and rollback limits are recorded; a clean revert/rebuild/redeploy run remains pending external deployment-owner execution.

## C release-polish contract evidence (2026-08-25)

- Preview, runtime snapshot/stream, readiness, and activation failures use the generated `handlers.APIErrorResponse` envelope: `error.code`, `error.message`, `error.retryable`, and opaque `error.request_id`.
- Generated Swagger also includes the complete Modbus Share bootstrap/status, start/stop, mappings, sync, write, and reconcile routes plus durable settings update and activation request/error paths (`StartModbusShareRequest`, `WriteTagValueRequest`, `UpsertMirrorMappingRequest`, `ReconcileRequest`, `UpdateSettingRequest`, `ActivateWorkspaceRequest`) with the same typed envelope refs.
- `/studio/runtime` consumes snapshot-first data, layers valid SSE values, and polls snapshots at 5 seconds after stream error/unavailability while preserving the last successful snapshot. The page shows loading, empty, degraded, or error states and removes the live tone after stream loss.
- Runtime known codes are `runtime_device_not_found`, `runtime_snapshot_unavailable`, and `runtime_stream_unavailable`. Unknown codes use generic localized fallback with request ID; raw backend exception text is excluded from the normal dashboard error view.
- Generated Swagger artifacts are `docs/swagger/docs.go`, `docs/swagger/swagger.yaml`, and `docs/swagger/swagger.json`. Registry ownership and runtime inventory state are updated in the same batch; inventory edits are recorded in `changelog.sqlite`.

## Historical acceptance run (2026-08-25; superseded by the refresh below)

- Historical snapshot only: focused frontend runtime/typed-error/locale suite (40 tests), full frontend Vitest suite (122 files / 651 tests), frontend lint, frontend build, focused Go handler/runtime tests (17), and full `go test ./...` (1035 tests) passed at that point, along with `go vet ./...`, `make check-lines`, Swagger generation via `go run github.com/swaggo/swag/cmd/swag@v1.16.6`, `spectra analyze`, `spectra validate --strict`, and scoped `git diff --check`. The later refresh below is the current repository evidence.
- Passed: embedded Playwright delivery smoke **8/8 cases in 13.6s** after the harness built and synchronized `frontend/dist` into the embedded static tree. The run covered `/studio/v2`, generic `/studio` and unknown-route equivalence, `/studio/runtime`, `/test`, `/gateway/quick-setup`, and deliberate same-origin asset 404 detection.
- Passed: route-level runtime dashboard tests and preserved-surface route assertions in the frontend suite. GUI, PLC, external service, and field acceptance were not exercised.
- Not a change failure: repository-wide `golangci-lint run ./...` reports 1120 existing issues across 297 files, dominated by unrelated `noctx`, `gosec`, and `gocritic` findings; no broad lint cleanup was attempted.
- Debug evidence: Phase 1 reproduced the timeout with only `embed-placeholder.txt` present (`GET /` returned 500); Phase 2 isolated the product/static boundary and separately reproduced a `react-vendor` circular-chunk initialization error after static sync; Phase 3 confirmed no port/lock or product startup regression (the server reached its configured port); Phase 4 applied three bounded fixes: self-contained static synchronization in the smoke harness, removal of the circular manual vendor partitions, and no timeout increase. The final elevated run passed all 8 cases; the initial Go cache permission denial remains an environment-only limitation.

## Acceptance refresh (2026-08-26; review25 final B10 evidence)

- Passed: review17 annotation-drift check found the existing runtime status schema's `projection_code` field missing from the three canonical artifacts; only that drift was regenerated into `docs/swagger/docs.go`, `docs/swagger/swagger.json`, and `docs/swagger/swagger.yaml`. The synchronized contract includes direct Modbus Share mapping delete `422`, preview SSE `workspace_id` scope, and `DesiredMapping` full identity/geometry plus reconcile/settings CAS fields.
- Passed: current B10 pair `b10-exe-acceptance-20260825T155805Z.json` and `.../b10-exe-acceptance-20260825T155901Z.json` are both **13/13 PASS**, with focused Go **22/22**, Python **16/16**, `py_compile` PASS, and binary SHA-256 `8398ffdf4b8977848c7dffd22a7fd3df367643b2867df795d3a5f61c3d0b2de7`. The 155314Z failure and earlier pairs remain historical; its witness was removed by retention and is intentionally not linked.
- Passed: startup hydration semantics are fail-closed: persisted settings load without binding; source-rule candidates restore the canonical projection under non-empty workspace/settings revisions and readiness token; only ready hydration starts the configured listener. The process-global projection transaction lock protects memory-bank snapshot/swap/rollback across workspaces, trading some throughput for rollback isolation while prevalidation/CAS run before the global lock. Uncertain rollback fail-closes first; a `MarkDirty` failure does not reopen listener/readiness.
- Passed: readiness, projection, `runtime_apply_*`, and connector messages are safe operator copy with endpoint/credential/DSN/transport details redacted. Preview stream requires persisted `mapping_id` plus `workspace_id` ownership before subscription; runtime/preview typed events carry `retryable`, opaque `request_id`, and actionable `action` metadata. Runtime reconnect rebuilds EventSource and snapshot fallback polls every 5 seconds; preview retries are bounded and explicit disconnect cancels them.
- Passed: current ownership/runtime evidence is fail-closed: Share list/delete/upsert/activation/restore require ready hydration, durable workspace relationships, matching revisions/readiness token, and the global listener gate; cross-workspace preview/reconcile attempts are rejected before subscription/projection mutation. Typed runtime/preview SSE envelopes and en/zh-TW i18n safe-copy assertions preserve request IDs without raw diagnostics; stale span/capacity and bind-conflict failures retain the prior projection.
- Passed: current target-writer evidence is documented: collector target fan-out writes the database target and authoritative Modbus Share projection only after hydration, while per point/tag delivery diagnostics are exposed for runtime observation; diagnostic sync endpoints are not the runtime delivery path.
- Passed: `go test ./... -count=1` (all packages), race tests for `internal/datalink/modbusshare`, `internal/datalink/settings`, `internal/api`, and `internal/api/handlers`, and `go vet ./...`.
- Recorded local frontend evidence: the prior working-tree Vitest run reported **134 files / 722 tests**; after the contract-integration changes this refresh makes no current frontend full-suite count claim. Earlier test-count results remain historical. The prior local Playwright delivery smoke **8/8 cases in 15.4s** covered `/studio/v2`, generic `/studio` and unknown-route equivalence, `/studio/runtime`, `/test`, `/gateway/quick-setup`, and deliberate same-origin asset 404 detection. These are static/delivery evidence only, not an embedded Chromium, Windows, LAN, PLC, SCADA, or field acceptance gate.
- Passed: `make check-lines`, scoped unstaged/staged/HEAD `git diff --check`, `spectra analyze` and `spectra validate --strict` for both changes. Swagger was regenerated after the `DesiredMapping.ownership_proof` source contract became reachable; generated JSON/YAML compare byte-for-byte with an independent temporary regeneration, and the handler response-matrix contract test passes. `retire-legacy-studio-and-polish-v2` analyze reports 19 existing Suggestion-level ambiguity findings; strict validation remains valid. Registry/runtime owner surfaces, typed SSE/i18n contracts, direct-mutation diagnostics, stale/capacity invariance, and inventory changelog readback match the current contract.
- No current source unit failure is recorded in the docs-only verification scope; B10 current acceptance is represented by the parent-reviewed 155805Z/155901Z PASS pair. No product-code change was made during this documentation/contract refresh.
- The missing normal `proc.wait()` success `return True` explains the earlier historical shutdown false negative only; it is not used as evidence for the current pair.
- GUI, external service, embedded Chromium, external Chromium, PLC, and field acceptance remain outside this local run. The dedicated Git rollback/rebuild/redeploy witness also remains an operator acceptance step.

## Contract refresh (2026-08-25; review22)

- Modbus Share status is now documented as workspace-scoped: `workspace_id`, `workspace_revision`, `settings_revision`, `listener_state`, and workspace-only `mapping_count` share the hydrated status boundary. Canonical desired mappings include full register geometry and persisted `ownership_proof`; dirty/uncertain recovery keeps listener/readiness fail-closed.
- The backend lifecycle contract now serializes settings/reconcile with a shared epoch (a staged reconcile after disable is rejected without durable/runtime mutation); status keeps only safe typed bind failure `code`/`retryable`/`action` and clears it on success/disable; candidate snapshot replacement uses repository revision CAS; device attach/detach advances scope through the durable workspace revision store and restores the workspace record if that CAS fails.
- Candidate review, recompute, tag-review decision, and output-apply requests share `workspace_id`, `expected_workspace_revision`, and `revision_id` gates. Step 4 reads persisted candidate snapshots for status, blocking reason, collision, and Local Modbus geometry; browser-only candidate values are not authoritative.
- Backend rejects stride values smaller than datatype span, foreign workspace/device scope, stale candidate revisions, and unproven ownership before mutation. Reconcile returns outcome, removed/invalidated spans, and safe diagnostics; `dirty_unknown` and `invalidated_unknown` are not ready/applied states.
- The generated Swagger response matrix remains the reachable handler contract (Start `400/409/422/503`; Stop `400/409/422/500/503`; mapping list `403/422/503`; mapping upsert `400/403/422/503`; mapping delete `403/422/503`; reconcile `400/403/409/422/503`). This closeout records the current B10 pair above and does not modify the other agent's setup/database change.

## Contract/docs closeout (2026-08-25; review20)

- Activation now has one public request contract: the JSON body is mandatory and `workspace_revision`, `settings_revision`, and `readiness_token` are required fields in Go binding validation and all three generated Swagger artifacts. Missing body remains a typed `422` save-barrier failure; missing required fields are a typed `400` invalid-request failure. Valid pending/error saves remain typed `422`, and stale revisions remain typed `409`.
- `projection_code` is present exactly once in each generated runtime status schema (`docs/swagger/docs.go`, `docs/swagger/swagger.json`, and `docs/swagger/swagger.yaml`); the activation request required metadata is likewise present exactly once per artifact. Swagger was regenerated from the current annotations with `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/test_ui/main.go -o docs/swagger`.
- Router-owned Modbus Share activation diagnostics use lower-case operator messages/actions. No frontend, B10 harness, witness, commit, or archive work was included in this closeout.
- Modbus Share handler annotations now document only the reachable typed failure matrix: Start `400/409/422/503`; Stop `400/409/422/500/503`; mapping list `403/422/503`; mapping upsert `400/403/422/503`; mapping delete `403/422/503`; reconcile `400/403/409/422/503`. The contract test asserts the generated Swagger JSON matches these exact operation response sets.
- Swagger was regenerated into `docs/swagger/docs.go`, `docs/swagger/swagger.json`, and `docs/swagger/swagger.yaml` from the current annotations; an independent regeneration into a temporary `swagger` output directory compares byte-for-byte (`cmp`) with all three artifacts (drift 0). The accidental zero-byte root `gateway_metrics` artifact was moved to `/private/tmp/go_gateway-review20-artifacts/gateway_metrics.zero-byte` for recovery and is not part of this change.
- Current `make check-lines` passes. The two modified frontend locale JSON files are each `500` lines (warning threshold only, not a blocker); `scripts/lib/b10_acceptance_helpers.py` is now exactly `300` lines, so it needs no line-plan entry. Unstaged, staged, and `HEAD`-relative `git diff --check` all pass.

### File-line warning handoff

The line-limit checker reports the following files in this change above the `300`-line warning threshold. These are warnings, not evidence that the files were split; no file was split in review20. Each row records the precise reason for retaining the current boundary and the follow-up split plan.

| File | Lines | Current-boundary reason | Follow-up split plan (not completed) |
| --- | ---: | --- | --- |
| `frontend/src/components/datalink/wizard/MappingWizard.tsx` | 345 | Existing V2 wizard orchestration remains the owner of step sequencing. | Extract step navigation/state adapters, then retain only shell composition. |
| `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx` | 308 | Runtime page still owns snapshot/stream state composition and route-level rendering. | Move stream/error state composition into a dashboard controller hook. |
| `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx` | 473 | V2 shell owns preserved route layout, settings status, and step handoff. | Split route shell, status rail, and step content boundaries. |
| `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts` | 454 | Shared reducer hook still owns the cross-step V2 state contract. | Extract settings, device, and step-local reducers behind a composed hook. |
| `frontend/src/i18n/locales/en/workbench-v2.json` | 500 | Locale file is a generated/curated namespace boundary for the V2 surface. | Split locale namespaces by step while preserving one locale loader. |
| `frontend/src/i18n/locales/zh-TW/workbench-v2.json` | 500 | Locale file mirrors the English V2 namespace for parity. | Split locale namespaces by step together with the English file. |
| `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts` | 488 | Autosave state coordinates all durable V2 save barriers. | Extract per-domain save state and compose a narrow page-level barrier. |
| `frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts` | 371 | Database autosave still owns target/schema ordering and failure recovery. | Split target persistence from schema-generation orchestration. |
| `frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx` | 444 | Page test covers the retained runtime route's integrated state/render contract. | Extract focused fixtures for stream, polling, and render-state assertions. |
| `frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx` | 393 | Rule autosave tests cover the page-level durable barrier path. | Split persistence, error recovery, and navigation assertions. |
| `frontend/tests/unit/workbench-v2/settings.test.tsx` | 414 | Settings tests cover the integrated settings page and operation ownership. | Separate query/mutation ownership tests from view interaction tests. |
| `frontend/tests/unit/workbench-v2/step2-rule.test.tsx` | 439 | Step 2 test retains the rule editor's end-to-end local contract. | Split editor rendering, candidate loading, and autosave behavior. |
| `frontend/tests/unit/workbench-v2/step3-components.test.tsx` | 331 | Step 3 component tests still share mapping fixtures and preview seams. | Group table, cell, and preview component tests into focused files. |
| `frontend/tests/unit/workbench-v2/step3-live-values.test.tsx` | 454 | Live-value tests cover the mapping preview stream contract as one scenario set. | Split stream lifecycle, malformed event, and value rendering cases. |
| `frontend/tests/unit/workbench-v2/step4-commit.test.tsx` | 359 | Step 4 commit tests cover activation barrier and handoff in one page flow. | Separate request-contract, error-state, and success-handoff cases. |
| `internal/api/handlers/polling_group_handler_test.go` | 415 | Existing polling-group handler suite predates this closeout and remains its package boundary. | Split CRUD and validation cases into focused handler test files. |
| `internal/api/handlers/runtime_handler.go` | 425 | Runtime handler still composes the preserved status/setup-context response contract. | Extract setup-context and readiness response builders. |
| `internal/api/handlers/runtime_handler_test.go` | 365 | Runtime handler tests cover the integrated preserved-surface response contract. | Split status, setup context, and error envelope tests. |
| `internal/api/handlers/runtime_workspace_setup_context_regression_test.go` | 315 | Regression suite preserves a historical cross-context behavior boundary. | Move fixtures/helpers into a shared test helper and split scenario groups. |
| `internal/api/handlers/source_rule_handler_database_candidates_test.go` | 391 | Candidate isolation tests cover a broad persisted-output contract. | Split connector discovery, ownership, and failure cases. |
| `internal/api/handlers/studio_v2_workspace_devices_handler.go` | 339 | Workspace device handler owns the preserved V2 device CRUD and readiness wiring. | Extract request parsing and readiness projection helpers. |
| `internal/api/router.go` | 452 | Router remains the single registration boundary for all retained API surfaces. | Extract route groups into domain registration files without changing identity. |
| `internal/datalink/migrator.go` | 485 | Migration orchestration must remain ordered in one transaction boundary. | Split migration discovery, execution, and compatibility checks. |
| `internal/datalink/modbusshare/service_test.go` | 325 | Share service tests retain lifecycle and persistence behavior in one package suite. | Split settings, listener, and projection service cases. |
| `internal/datalink/point/service_point_crud.go` | 316 | Point CRUD service keeps repository and runtime synchronization atomic. | Extract repository CRUD from runtime synchronization. |
| `internal/datalink/runtime/delivery_diagnostic.go` | 308 | Delivery diagnostic types and safe projection builders share one contract. | Split diagnostic collection from response projection. |
| `internal/datalink/runtime/service.go` | 470 | Runtime service owns lifecycle, projection, and delivery coordination. | Extract lifecycle coordination and target-delivery orchestration. |
| `internal/datalink/runtime/status.go` | 394 | Runtime status projection must preserve one coherent response schema. | Split status assembly from per-target diagnostic projection. |
| `internal/datalink/schema/migrations/001_initial_schema.up.sql` | 388 | Initial schema is intentionally one ordered bootstrap migration. | Future migrations should carry new domain tables; do not split history retroactively. |
| `internal/datalink/workspace/service_readiness.go` | 487 | Readiness aggregation preserves one ordered device/connector/workspace gate. | Extract device, connector, and summary aggregation into separate services. |
| `internal/virtual/server/modbus/server_test.go` | 431 | Virtual Modbus server tests cover protocol lifecycle and register behavior together. | Split lifecycle, register handlers, and concurrency fixtures. |
| `frontend/src/services/datalink.ts` | 751 | Legacy file is tracked above the hard limit and is already shrinking; this closeout did not enlarge it. | Continue extracting typed API clients by domain without changing legacy exports. |
| `frontend/src/types/datalink.ts` | 691 | Legacy type barrel is tracked above the hard limit and remains a compatibility export surface. | Move domain types into per-feature modules and retain a compatibility barrel. |
| `internal/api/handlers/device_handler_extended_test.go` | 725 | Legacy test file is above the hard limit and was not enlarged by this closeout. | Extract activation, CRUD, and runtime-sync scenarios incrementally. |
| `internal/api/handlers/mapping_handler_extended_test.go` | 544 | Legacy test file is above the hard limit and was not enlarged by this closeout. | Split preview, CRUD, and validation scenarios. |
| `internal/datalink/sourcerule/service.go` | 1204 | Legacy service is above the hard limit; this closeout preserves its established source-rule seam. | Extract candidate, ownership, and runtime-reconcile services in a dedicated change. |
