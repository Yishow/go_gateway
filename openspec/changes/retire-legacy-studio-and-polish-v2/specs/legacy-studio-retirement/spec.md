## ADDED Requirements

### Requirement: Canonical studio route and non-target boundaries

The product SHALL use `/studio/v2` as the only user-facing setup entry and `/studio/runtime` as the focused post-setup monitor. The product SHALL preserve `/test` as an independent engineering tool and `/gateway/*` as an experimental surface. This change MUST NOT remove, redirect, or reinterpret `/studio/v2`, `/studio/runtime`, `/test`, or `/gateway/*`.

#### Scenario: Canonical product entry

- **WHEN** an operator opens the root route, a generic datalink landing route, or `/studio/v2`
- **THEN** the application resolves to the v2 guided setup shell
- **AND** no default route renders the legacy workspace

#### Scenario: Non-target routes remain distinct

- **WHEN** an operator opens `/studio/runtime` with supported context, `/test`, or a `/gateway/*` route
- **THEN** `/studio/runtime` remains the focused monitor, `/test` remains the engineering tool, and `/gateway/*` remains experimental
- **AND** none of these route families redirects to a different product surface

### Requirement: Immediate legacy studio deletion is bounded and complete

The implementation SHALL immediately remove the dedicated `/studio` route registration, any dedicated `/studio` handler or tombstone, all proven legacy-only frontend imports/chunks/assets, obsolete legacy tests, and obsolete legacy docs references. The deletion SHALL be driven by a reviewed pre-delete inventory. After deletion, `/studio` SHALL be indistinguishable from an arbitrary unknown route under the existing generic unknown-route policy, including its standard redirect or fallback when configured. The implementation SHALL NOT add a `/studio` special redirect, fallback, compatibility flag, or legacy mount.

#### Scenario: Pre-delete inventory controls the deletion scope

- **WHEN** the deletion batch is prepared
- **THEN** the inventory lists every `/studio` route registration, direct/lazy import, chunk, asset, test, docs reference, and preserved-surface dependency
- **AND** a file with an unproven legacy-only boundary is retained and recorded as out of deletion scope

#### Scenario: Deleted path follows generic unknown-route behavior

- **WHEN** a browser or HTTP client requests `/studio` after deletion alongside an arbitrary unknown path such as `/unknown-route-for-retirement`
- **THEN** both paths produce the same repository-policy response/navigation, status/fallback, error boundary, and network asset pattern
- **AND** neither path matches a dedicated `/studio` route, handler, tombstone, special redirect, or legacy workspace

#### Scenario: Preserved surfaces remain available

- **WHEN** an operator opens `/studio/v2`, `/studio/runtime`, `/test`, or `/gateway/quick-setup` after deletion
- **THEN** each route retains its established identity and required assets
- **AND** no preserved surface loads a legacy-only chunk or is redirected to a new product surface

### Requirement: Git rollback and recovery boundaries are explicit

The deletion SHALL be delivered as a dedicated Git commit that references the pre-delete inventory, deletion list, verification results, and rollback command. Recovery SHALL use a Git revert of that commit followed by rebuild, redeploy, and the same route/static/browser smoke. The recovery contract SHALL state that Git restores tracked source and docs but does not restore deployment/runtime data, external database contents, device state, or browser bookmarks.

#### Scenario: Deletion commit has a runnable rollback

- **WHEN** the release owner evaluates the deletion commit
- **THEN** the evidence names the commit, revert command, rebuild command, redeploy target, smoke results, and source/data/bookmark limitations
- **AND** a clean checkout can reproduce the deletion and revert verification without manual file copying

#### Scenario: Rollback scope is reported honestly

- **WHEN** the dedicated deletion commit is reverted and the binary is rebuilt
- **THEN** tracked legacy route/source/docs definitions are restored and route/static/browser smoke is rerun
- **AND** deployment/runtime data and browser bookmarks are reported as outside Git recovery scope

### Requirement: Legacy deletion is independent of unrelated lifecycle gates

The immediate legacy deletion batch SHALL NOT require telemetry, usage windows, deprecation phases, query migration, redirect windows, database changes, or A/B acceptance evidence. This independence SHALL NOT change, skip, or weaken the separate Modbus Share lifecycle, CAS, concurrency, ownership, or acceptance requirements.

#### Scenario: Deletion can start without unrelated evidence

- **WHEN** the owner has approved the deletion and the pre-delete inventory is complete while unrelated A/B evidence is absent
- **THEN** the legacy deletion batch can proceed through its own route/import/asset/rollback gates
- **AND** the batch does not add or modify Modbus Share behavior or claim that Modbus Share is accepted

### Requirement: Typed safe errors are mapped to v2 i18n

Preview, runtime snapshot/stream, workspace readiness, and activation failures SHALL use an error object containing `code`, `message`, `retryable`, and `request_id`. The supported code set SHALL include `preview_invalid_request`, `preview_unavailable`, `runtime_device_not_found`, `runtime_snapshot_unavailable`, `runtime_stream_unavailable`, `workspace_not_ready`, and `activation_failed`. The v2 UI SHALL map each code to en and zh-TW copy and SHALL keep raw exception details in a separate diagnostics channel.

#### Scenario: Known typed error renders safe action

- **WHEN** the preview endpoint returns code `preview_unavailable` with `retryable` true
- **THEN** the v2 page renders the localized preview-unavailable message and retry action
- **AND** it does not render the raw backend message or stack trace

#### Scenario: Unknown typed error fails safely

- **WHEN** a supported endpoint returns an unrecognized error code
- **THEN** the UI renders a localized generic error with the opaque `request_id` and an actionable retry path
- **AND** the raw response string is not inserted into the page

### Requirement: Step 3 preview and SSE states are truthful

Step 3 SHALL expose `idle`, `connecting`, `live`, `reconnecting`, `degraded`, `error`, and `stale` states for server preview and SSE values. An initial server preview failure SHALL render error without ready/live data; an SSE disconnect after a successful value SHALL preserve the last value with stale context and render reconnecting or degraded; a recovered stream SHALL restore live only after a valid server event. The UI SHALL use bounded retry delays of 1, 2, 4, 8, 16, and 30 seconds and SHALL NOT silently substitute local or mock values for a server failure.

#### Scenario: Preview failure is not a silent success

- **WHEN** the server preview returns a typed error before any value is received
- **THEN** the Step 3 panel renders an actionable error state
- **AND** no local/mock value is presented as live or ready

#### Scenario: Stream recovery is visible

- **WHEN** the SSE connection disconnects after a valid value and later reconnects
- **THEN** the panel preserves the last value with stale/reconnecting context, retries with bounded backoff, and renders live only after a valid event
- **AND** the state transition is available to the operator

### Requirement: Share and CommitProgress presentation is accessible and single-flight

The Share/global toggle SHALL expose an i18n accessible name, `aria-pressed` state, visible focus, and Enter/Space keyboard behavior. While its UI mutation is pending it SHALL be disabled and SHALL accept at most one request for the triggering gesture. CommitProgress SHALL render pending, success, failed, and skipped rows from actual results; a failed row SHALL show safe code/message and retry action and SHALL NOT show a fixed HTTP 200 success value. This requirement SHALL NOT define backend CAS or concurrency semantics.

#### Scenario: Keyboard toggle prevents double submit

- **WHEN** an operator activates the Share/global toggle with keyboard or pointer while no request is pending
- **THEN** the control exposes the localized name and current `aria-pressed` state, sends one mutation, and becomes disabled while pending
- **AND** a repeated activation during pending sends no second UI request

#### Scenario: Failed progress is not reported as 200

- **WHEN** the backend returns a failed result for one CommitProgress item
- **THEN** that row is rendered with failed tone, safe error copy, and retry action
- **AND** the row does not display a fixed 200 or success status

### Requirement: Documentation and inventory delivery is traceable

The release SHALL regenerate the affected Swagger artifacts, update the backend API registry and runtime/inventory docs with immediate deletion, generic unknown-route equivalence, preserved route identities, pre-delete inventory, and Git rollback limitations, add a release note, and write an inventory changelog entry for every modified file under the studio-surface-inventory directory. The release SHALL NOT claim unrelated data migration, telemetry, or historical Swagger work.

#### Scenario: Formal docs match deletion contract

- **WHEN** the implementation regenerates API and release documentation
- **THEN** the docs describe the removed dedicated `/studio` route and generic unknown-route-equivalent behavior
- **AND** the docs identify `/studio/v2`, `/studio/runtime`, `/test`, and `/gateway/*` as preserved surfaces and describe Git rollback limitations

#### Scenario: Inventory edit has a changelog record

- **WHEN** an implementation task modifies any inventory document
- **THEN** `go run ./cmd/studio_inventory_changelog add` records summary, surface, files, and reason in `changelog.sqlite`
- **AND** the readback entry names every modified inventory file
