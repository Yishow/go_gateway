# datalink-api Specification

## Purpose
TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.

## Requirements

### Requirement: Management APIs

The system SHALL expose APIs for devices, points, tags, and mappings with persistent storage.

All data SHALL be persisted to SQLite (for development/testing) or PostgreSQL (for production) using a Repository pattern with factory-based selection.

#### Scenario: Create mapping via API

- WHEN a client posts a mapping
- THEN the mapping is saved to database and returned with an id

#### Scenario: Data persistence across restarts

- WHEN the server restarts
- THEN all previously saved data is available

---
### Requirement: Connection and preview APIs
The system SHALL expose APIs to test device connections and preview mapping results.

#### Scenario: Mapping preview API
- WHEN a client requests a preview for a mapping
- THEN the API returns raw and transformed values

---
### Requirement: Write settings API

The system SHALL expose APIs to read and update write settings, including timestamp precision and partition interval.

Settings SHALL be persisted to database (not hardcoded).

#### Scenario: Update write precision

- WHEN a client updates the write precision to milliseconds
- THEN subsequent writes use millisecond precision

#### Scenario: Get settings

- WHEN a client requests settings
- THEN the API returns current settings from database

---
### Requirement: Polling Groups API

The system SHALL expose APIs for managing polling groups.

Endpoints:

- `GET /datalink/polling-groups` - List all polling groups
- `POST /datalink/polling-groups` - Create polling group
- `GET /datalink/polling-groups/:id` - Get polling group by ID
- `PUT /datalink/polling-groups/:id` - Update polling group
- `DELETE /datalink/polling-groups/:id` - Delete polling group

#### Scenario: Create polling group

- WHEN a client posts a polling group with name and interval
- THEN the polling group is saved and returned with an id

#### Scenario: List polling groups

- WHEN a client requests polling groups
- THEN all polling groups are returned

---
### Requirement: Protocols API

The system SHALL expose an API to retrieve supported protocols and their configuration schemas.

Endpoint:

- `GET /datalink/protocols` - List supported protocols

#### Scenario: List protocols

- WHEN a client requests protocols
- THEN the API returns protocol list with type, name, description, and config_schema

---
### Requirement: Health API

The system SHALL expose a health check API for the Datalink service.

Endpoint:

- `GET /datalink/health` - Check service health

#### Scenario: Health check

- WHEN a client requests health status
- THEN the API returns status and service name

---
### Requirement: Tag lifecycle APIs

The system SHALL expose APIs for tag lifecycle management.

Endpoints:

- `POST /datalink/tags/:id/activate` - Activate a draft tag
- `POST /datalink/tags/:id/retire` - Retire an active tag
- `POST /datalink/tags/batch` - Batch create tags
- `POST /datalink/tags/validate-key` - Validate tag key format

#### Scenario: Activate tag

- WHEN a client activates a draft tag
- THEN the tag status changes to active

#### Scenario: Retire tag

- WHEN a client retires an active tag
- THEN the tag status changes to retired

#### Scenario: Batch create tags

- WHEN a client posts multiple tags
- THEN all valid tags are created and errors are reported

#### Scenario: Validate tag key

- WHEN a client validates a tag key
- THEN the API returns validity, normalized key, and existence status

---
### Requirement: Device batch test API

The system SHALL expose an API to batch test device connections.

Endpoint:

- `POST /datalink/devices/test-batch` - Test multiple device connections

#### Scenario: Batch test connections

- WHEN a client posts device IDs for testing
- THEN the API tests all devices and returns individual results

---
### Requirement: Point polling APIs

The system SHALL expose APIs for on-demand point polling.

Endpoints:

- `POST /datalink/points/:id/poll` - Poll single point
- `POST /datalink/points/poll` - Batch poll multiple points

#### Scenario: Poll single point

- WHEN a client requests a point poll
- THEN the API returns the current value and timestamp

#### Scenario: Batch poll points

- WHEN a client posts point IDs for polling
- THEN the API polls all points and returns results

---
### Requirement: Pipeline validation API

The system SHALL expose an API to validate transform pipelines.

Endpoint:

- `POST /datalink/mappings/validate-pipeline` - Validate pipeline

#### Scenario: Validate valid pipeline

- WHEN a client submits a valid pipeline
- THEN the API returns valid=true

#### Scenario: Validate invalid pipeline

- WHEN a client submits an invalid pipeline
- THEN the API returns valid=false with error details

---
### Requirement: Preview stream API

The system SHALL expose an SSE endpoint for live mapping preview.

Endpoint:

- `GET /datalink/preview/stream` - Server-Sent Events stream

#### Scenario: Subscribe to preview stream

- WHEN a client connects to the preview stream with a mapping ID
- THEN the server pushes raw and transformed values in real-time

---
### Requirement: Device readiness API

The system SHALL expose an API endpoint to check device readiness status using the explicit planning, activation, and apply eligibility contract.

Endpoint:

- `GET /api/v1/datalink/devices/:id/readiness` - Check device readiness

#### Scenario: Get readiness status

- **WHEN** a client requests readiness status for a device
- **THEN** the API returns `connect_status`, `probe_status`, `planning_allowed`, `activation_allowed`, `apply_allowed`, and blocking metadata

#### Scenario: Readiness response format

- **WHEN** a readiness check is performed
- **THEN** the API returns a response with explicit planning/activation/apply eligibility fields
- **AND** does not require clients to infer workflow state from one final `ready` boolean alone

---
### Requirement: Rule-scoped candidate query and recompute APIs
The API SHALL expose rule-scoped endpoints to query and recompute candidate snapshots for a source-rule revision.

Endpoints:

- `GET /api/v1/datalink/source-rules/:id/candidates` - Get the current candidate snapshot for a source rule
- `POST /api/v1/datalink/source-rules/:id/candidates/recompute` - Recompute the candidate snapshot for the latest source-rule revision

The candidate query response MUST include at least:

- `source_rule_id`
- `revision_id`
- `tags`
- `database_outputs`
- `local_modbus_outputs`

Each returned candidate object MUST include at least:

- `id`
- `status`
- `proposed_signature`
- `last_applied_signature`
- `blocking_reason`

#### Scenario: Query rule candidates
- **WHEN** a client requests candidates for a source rule
- **THEN** the API returns the current tag, database, and Local Modbus candidate sets with revision and status metadata

#### Scenario: Recompute candidate snapshot
- **WHEN** a client requests candidate recompute for a source rule
- **THEN** the API regenerates the current revision's candidate snapshot using the shared identity/signature contract
- **AND** returns the recomputed snapshot without introducing target-specific apply side effects

---
### Requirement: Device collection statistics API

The system SHALL expose an API endpoint to retrieve device collection statistics.

Endpoint:

- `GET /api/v1/datalink/devices/:id/stats` - Get device collection statistics

#### Scenario: Get device statistics

- **WHEN** a client requests device statistics
- **THEN** the API returns collection statistics including last collection time, collection count, and error count

---
### Requirement: Device status stream API

The system SHALL expose an SSE endpoint for real-time device status updates.

Endpoint:

- `GET /api/v1/datalink/devices/stream` - Server-Sent Events stream for device status

#### Scenario: Subscribe to device status stream

- **WHEN** a client connects to the device status stream
- **THEN** the server pushes device status updates in real-time

---
### Requirement: Rule-scoped tag apply APIs
The API SHALL expose explicit rule-scoped endpoints to apply approved tag candidates for a source-rule revision.

Endpoint:

- `POST /api/v1/datalink/source-rules/:id/tags/apply` - Apply approved tag candidates for the current source-rule revision

The apply request MUST identify the source-rule revision and the tag candidate ids being applied.

The apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply approved tag candidates
- **WHEN** a client applies tag candidates for a source rule
- **THEN** the API applies only the approved tag candidates for that revision
- **AND** leaves database and Local Modbus outputs as downstream review targets instead of silently applying them

#### Scenario: Partial tag apply returns per-item results
- **WHEN** a client applies multiple tag candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result

---
### Requirement: Rule-scoped database output apply APIs
The API SHALL expose explicit rule-scoped endpoints to apply approved database output candidates for a source-rule revision.

Endpoint:

- `POST /api/v1/datalink/source-rules/:id/database-outputs/apply` - Apply approved database output candidates for the current source-rule revision

The apply request MUST identify the source-rule revision and the database candidate ids being applied.

The apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply database outputs without affecting Local Modbus
- **WHEN** a client applies database output candidates for a source rule
- **THEN** the API applies only the database target changes
- **AND** preserves Local Modbus candidate and apply state for the same rule

#### Scenario: Partial database apply returns per-item results
- **WHEN** a client applies multiple database candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result

#### Scenario: Stale database apply revision is rejected
- **WHEN** a client applies database output candidates with a stale `revision_id`
- **THEN** the API rejects the request with a `revision-conflict` error
- **AND** returns enough revision context for the UI to refresh candidates instead of silently retrying

---
### Requirement: Rule-scoped Local Modbus apply APIs
The API SHALL expose explicit rule-scoped endpoints to apply approved Local Modbus candidates for a source-rule revision.

Endpoint:

- `POST /api/v1/datalink/source-rules/:id/local-modbus/apply` - Apply approved Local Modbus candidates for the current source-rule revision

The apply request MUST identify the source-rule revision and the Local Modbus candidate ids being applied.

The apply response MUST report per-item success or failure so partial success is explicit.

#### Scenario: Apply Local Modbus outputs without affecting Database
- **WHEN** a client applies Local Modbus candidates for a source rule
- **THEN** the API applies only the Local Modbus target changes
- **AND** preserves database candidate and apply state for the same rule

#### Scenario: Partial Local Modbus apply returns per-item results
- **WHEN** a client applies multiple Local Modbus candidates and only some succeed
- **THEN** the API returns per-candidate success and failure results
- **AND** does not collapse the response into one ambiguous pass/fail result

#### Scenario: Stale Local Modbus apply revision is rejected
- **WHEN** a client applies Local Modbus candidates with a stale `revision_id`
- **THEN** the API rejects the request with a `revision-conflict` error
- **AND** returns enough revision context for the UI to refresh candidates instead of silently retrying

---
### Requirement: Database connector and mapping APIs expose grouped-row metadata

The system SHALL expose grouped database-row planning metadata through connector and mapping APIs so `/studio` can plan interval-driven row writes without inventing a version-specific backend contract.

#### Scenario: Connector returns default write interval
- **WHEN** a client creates, updates, or reads a database connector for this workflow
- **THEN** the payload includes `default_write_interval_seconds`
- **AND** the field defaults to `15` unless explicitly changed

#### Scenario: Mapping preserves grouped compatibility fields
- **WHEN** a client creates, updates, or reads a database target mapping
- **THEN** the payload includes nullable `group_key` and optional `write_interval_seconds`
- **AND** a `null` `group_key` remains legacy single-member behavior rather than being coerced into a shared group

#### Scenario: Existing records keep compatible migration behavior
- **WHEN** existing database connectors and mappings are read after the grouped-row migration
- **THEN** connectors expose `default_write_interval_seconds = 15` unless a different value was explicitly persisted
- **AND** existing mappings keep `group_key = null` so they remain single-member rows

#### Scenario: Effective interval resolves predictably
- **WHEN** a client reads or validates a database mapping
- **THEN** the effective interval equals `write_interval_seconds` when that override is present
- **AND** otherwise the effective interval equals the connector `default_write_interval_seconds`

---
### Requirement: Database schema generation and mapping dry-run APIs

The system SHALL expose dedicated database tooling endpoints so `/studio` can preview schema mutations and validate mapping candidates against the current connector state.

Endpoints:

- `POST /api/v1/datalink/db-targets/connectors/:id/schema/generate` - Generate or preview connector schema changes
- `POST /api/v1/datalink/db-targets/connectors/:id/mappings/dry-run` - Dry-run selected mapping candidates before apply

The schema generate request MUST support `dry_run=true|false`.

The dry-run response MUST return candidate-level validation outcomes and machine-readable blocking categories for blocked candidates.

#### Scenario: Schema dry-run previews SQL without mutating connector schema
- **WHEN** a client calls schema generate with `dry_run=true`
- **THEN** the API returns SQL preview and validation feedback
- **AND** does not execute schema mutations

#### Scenario: Mapping dry-run identifies blocked candidates before apply
- **WHEN** a client requests mapping dry-run for selected candidates
- **THEN** the API returns which candidates are eligible to apply or blocked
- **AND** includes machine-readable blocking categories such as `schema_missing`, `connector_unavailable`, and `type_conflict`

---
### Requirement: Rule-scoped database candidate and apply APIs carry grouped planning metadata

The system SHALL return grouped database planning metadata in rule-scoped candidate, dry-run, and apply payloads so the frontend can review row groups before commit.

#### Scenario: Candidate payload exposes inferred group planning
- **WHEN** a client requests rule-scoped database output candidates for tags such as `meter/A1` and `meter/kw`
- **THEN** each candidate includes the inferred `group_key`, normalized `column_name`, and effective `write_interval_seconds`
- **AND** slashless tags continue to return `group_key = null`

#### Scenario: Column-name normalization is deterministic
- **WHEN** a client requests rule-scoped database output candidates for slash-based tag keys
- **THEN** the remaining suffix after the first slash is normalized to lowercase for `column_name`
- **AND** any additional `/` separators are converted to `_`

#### Scenario: Apply rejects incompatible grouped members
- **WHEN** a client dry-runs or applies grouped database candidates whose connector, schema, table, write mode, timestamp column, or effective interval do not match
- **THEN** the API returns an explicit blocking reason
- **AND** the system does not silently merge or remap incompatible members

---
### Requirement: Database write history reflects grouped row flushes

The system SHALL expose grouped database write history through `GET /api/v1/datalink/db-targets/connectors/:id/write-history` and SHALL report row outcomes per grouped flush rather than per raw tag event.

#### Scenario: Write history query returns grouped flush records
- **WHEN** a client requests grouped database write history for a connector
- **THEN** the response returns recent grouped flush records in reverse chronological order
- **AND** each record includes at least `observed_at`, `status`, `row_count`, `group_key`, `table_name`, `effective_interval_seconds`, and `error_summary` when failed

#### Scenario: Grouped flush records one row outcome
- **WHEN** four compatible `meter/*` tags flush into one database row for the same interval bucket
- **THEN** write history records one row outcome for that grouped flush
- **AND** the history remains directly renderable in Studio diagnostics

#### Scenario: Interval buckets use wall-clock alignment
- **WHEN** grouped database writes are bucketed for an effective interval
- **THEN** the bucket start is aligned by `observed_at.UTC().Truncate(effective_interval)`
- **AND** all writes within the same aligned bucket contribute to the same grouped row outcome

---
### Requirement: Runtime status snapshot API

The system SHALL expose `GET /api/v1/datalink/runtime/status` as the polling-compatible snapshot endpoint for runtime monitoring clients.

The endpoint SHALL accept an optional `device_id` query parameter.

Each successful response SHALL include top-level fields `running`, `uptime_seconds`, `metrics`, and `collectors`.

The `metrics` object SHALL include at least:

- `collected_total`
- `write_success_total`
- `write_error_total`
- `mapping_error_total`
- `point_state_error_total`

Each collector object SHALL include at least:

- `device_id`
- `device_name`
- `protocol`
- `status`
- `points_total`
- `points_healthy`
- `points_stale`
- `points_error`
- `last_read_at`
- `last_error`
- `breaker_state`

#### Scenario: Request one device snapshot

- **WHEN** a client requests `GET /api/v1/datalink/runtime/status?device_id=device-A`
- **THEN** the response includes runtime top-level metrics
- **AND** the `collectors` array contains exactly the collector summary for `device-A`

#### Scenario: Request fleet snapshot without device filter

- **WHEN** a client requests `GET /api/v1/datalink/runtime/status` without `device_id`
- **THEN** the response includes runtime top-level metrics
- **AND** the `collectors` array SHALL contain all device collector summaries visible to the runtime service


<!-- @trace
source: add-runtime-dashboard-backend-contract
updated: 2026-05-30
code:
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_handler.go
  - task_plan.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - cmd/studio_inventory_changelog/main.go
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - CLAUDE.md
  - AGENTS.md
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - progress.md
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/gap-roadmap.md
  - docs/technical/studio-surface-inventory/index.html
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/inventory.js
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - internal/datalink/runtime/status.go
  - frontend/src/i18n/config.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - docs/technical/studio-surface-inventory/test-tooling.md
  - docs/technical/studio-surface-inventory/inventory.css
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - findings.md
  - frontend/src/types/datalink.ts
  - docs/technical/studio-surface-inventory/README.md
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/runtime_stream_handler.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - docs/technical/studio-surface-inventory/context.json
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/datalink/runtime/stream.go
tests:
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - cmd/test_ui/static/index.html
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/runtime/stream_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
-->

---
### Requirement: Runtime value stream API

The system SHALL expose `GET /api/v1/datalink/runtime/stream` as the live runtime monitoring stream for one device.

The endpoint SHALL require `device_id` and SHALL accept an optional comma-separated `point_ids` query parameter.

The stream SHALL emit the following SSE event types:

- `value`
- `status`
- `heartbeat`

A `value` event payload SHALL include at least `device_id`, `point_id`, `address`, `raw_value`, `transformed_value`, `quality`, `stale`, and `timestamp`.

A `status` event payload SHALL include at least `device_id`, `status`, `points_total`, `points_healthy`, `points_stale`, `points_error`, `last_read_at`, `last_error`, and `breaker_state`.

A `heartbeat` event payload SHALL include `ts`.

#### Scenario: Reject stream request without device id

- **WHEN** a client requests `GET /api/v1/datalink/runtime/stream` without `device_id`
- **THEN** the server returns HTTP 400
- **AND** the request does not create a runtime subscription

#### Scenario: Stream only selected points for one device

- **WHEN** a client requests `GET /api/v1/datalink/runtime/stream?device_id=device-A&point_ids=pt-1,pt-2`
- **THEN** the stream emits `value` events only for `pt-1` and `pt-2`
- **AND** the stream continues to emit `status` and `heartbeat` events for `device-A`

#### Scenario: Emit device status change on runtime degradation

- **WHEN** runtime monitoring detects that `device-A` transitions from healthy collection to stale, point-error, or breaker-open state
- **THEN** the stream emits a `status` event for `device-A`
- **AND** the payload reflects the same derived status model used by the runtime snapshot endpoint

<!-- @trace
source: add-runtime-dashboard-backend-contract
updated: 2026-05-30
code:
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_handler.go
  - task_plan.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - cmd/studio_inventory_changelog/main.go
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - CLAUDE.md
  - AGENTS.md
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - progress.md
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/gap-roadmap.md
  - docs/technical/studio-surface-inventory/index.html
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/inventory.js
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - internal/datalink/runtime/status.go
  - frontend/src/i18n/config.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - docs/technical/studio-surface-inventory/test-tooling.md
  - docs/technical/studio-surface-inventory/inventory.css
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - findings.md
  - frontend/src/types/datalink.ts
  - docs/technical/studio-surface-inventory/README.md
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/runtime_stream_handler.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - docs/technical/studio-surface-inventory/context.json
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/datalink/runtime/stream.go
tests:
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - cmd/test_ui/static/index.html
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/runtime/stream_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
-->

---
### Requirement: Studio v2 workspace device autosave APIs

The system SHALL expose workspace-scoped device APIs for Step 1 autosave.

#### Scenario: List workspace devices

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/devices`
- **THEN** the API returns the devices currently attached to the singleton v2 workspace in persisted order

#### Scenario: Save one valid workspace device

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/devices` or `PUT /api/v1/datalink/studio-v2/workspace/devices/:id` with a valid device payload
- **THEN** the API persists that device and returns the saved record
- **AND** the response is scoped to that one device rather than a bulk transaction

#### Scenario: Reject one invalid workspace device without touching others

- **WHEN** a client sends one invalid Step 1 device payload
- **THEN** the API returns a validation error for that device
- **AND** previously saved devices in the same workspace remain unchanged


<!-- @trace
source: wire-studio-v2-step1-device-autosave
updated: 2026-05-31
code:
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_status.go
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/sourcerule/service.go
  - task_plan.md
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/migrator.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2Availability.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/router.go
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/device/service_point_read.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/services/studioV2Mappings.ts
  - .air.toml
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - go.sum
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - scripts/start-process-utils.sh
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/availability.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - findings.md
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/errors.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/index.ts
  - internal/datalink/workspace/service_activation.go
  - progress.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/workspace/service_devices.go
  - scripts/start-log-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/db.go
  - internal/datalink/workspace/memory_repository.go
tests:
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_device_order_test.go
-->

---
### Requirement: Studio v2 workspace device order API

The system SHALL expose `PUT /api/v1/datalink/studio-v2/workspace/device-order` to persist the device tab order used by V2.

#### Scenario: Persist device order

- **WHEN** a client submits an ordered list of workspace device ids
- **THEN** the API stores that order for the singleton workspace
- **AND** a later workspace bootstrap returns the same order

##### Example: reordered tabs survive reload

- **GIVEN** the current workspace device order is `[dev-A, dev-B, dev-C]`
- **WHEN** a client submits `[dev-C, dev-A, dev-B]` to `PUT /api/v1/datalink/studio-v2/workspace/device-order`
- **THEN** the next workspace bootstrap returns `[dev-C, dev-A, dev-B]`

<!-- @trace
source: wire-studio-v2-step1-device-autosave
updated: 2026-05-31
code:
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_status.go
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/sourcerule/service.go
  - task_plan.md
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/migrator.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2Availability.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/router.go
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/device/service_point_read.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/services/studioV2Mappings.ts
  - .air.toml
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - go.sum
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - scripts/start-process-utils.sh
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/availability.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - findings.md
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/errors.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/index.ts
  - internal/datalink/workspace/service_activation.go
  - progress.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/workspace/service_devices.go
  - scripts/start-log-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/db.go
  - internal/datalink/workspace/memory_repository.go
tests:
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_device_order_test.go
-->

---
### Requirement: Studio v2 workspace rule autosave APIs

The system SHALL expose workspace-scoped source-rule APIs for Step 2 autosave.

#### Scenario: List workspace rules

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/source-rules`
- **THEN** the API returns the persisted rules currently attached to the singleton v2 workspace

#### Scenario: Save one valid workspace rule

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/source-rules` or `PUT /api/v1/datalink/studio-v2/workspace/source-rules/:id` with a valid rule payload
- **THEN** the API persists that rule and returns the saved record including `revision_id`
- **AND** the response is scoped to that one rule rather than a bulk transaction

#### Scenario: Reject one invalid workspace rule without touching others

- **WHEN** a client sends one invalid Step 2 rule payload
- **THEN** the API returns a validation error for that rule
- **AND** previously saved rules in the same workspace remain unchanged

<!-- @trace
source: wire-studio-v2-step2-rule-autosave
updated: 2026-05-31
code:
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/sourcerule/errors.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - docs/goal.md
  - internal/datalink/dbtarget/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/point_handler.go
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/db.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/migrator.go
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - scripts/start-process-utils.sh
  - start.ps1
  - findings.md
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service.go
  - go.sum
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - scripts/start-log-utils.sh
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/workspace/service_database.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - progress.md
  - start.sh
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/schema/schema_device_models.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/hooks/datalink/index.ts
tests:
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
-->

---
### Requirement: Studio v2 workspace mapping autosave APIs

The system SHALL expose workspace-scoped mapping APIs for Step 3 autosave.

#### Scenario: List workspace mappings

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/mappings`
- **THEN** the API returns the persisted mappings currently attached to the singleton v2 workspace

#### Scenario: Save one valid workspace mapping

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/mappings` or `PUT /api/v1/datalink/studio-v2/workspace/mappings/:id` with a valid mapping payload
- **THEN** the API persists that mapping and returns the saved record
- **AND** the response is scoped to that one row rather than a bulk transaction

#### Scenario: Reject one invalid workspace mapping without touching others

- **WHEN** a client sends one invalid Step 3 mapping payload
- **THEN** the API returns a validation error for that mapping row
- **AND** previously saved mappings in the same workspace remain unchanged

<!-- @trace
source: wire-studio-v2-step3-mapping-autosave
updated: 2026-05-31
code:
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/mapping/errors.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - start.sh
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - task_plan.md
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/mapping/service_crud.go
  - start.ps1
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - findings.md
  - go.sum
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/runtime/status.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/datalink/device/availability.go
  - scripts/start-process-utils.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - scripts/start-log-utils.sh
  - frontend/src/hooks/datalink/useSettings.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/runtime_handler.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/db.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - go.mod
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/device/service_status.go
  - docs/goal.md
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_logger_test.go
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
-->

---
### Requirement: Studio v2 workspace database autosave APIs

The system SHALL expose workspace-scoped Step 4 database autosave APIs.

#### Scenario: Save workspace database config

- **WHEN** a client requests `PUT /api/v1/datalink/studio-v2/workspace/database-config` with a valid connector payload
- **THEN** the API persists the workspace database config and returns the saved metadata

#### Scenario: Save one valid target row

- **WHEN** a client requests `PUT /api/v1/datalink/studio-v2/workspace/database-targets/:point_id` with a valid target payload
- **THEN** the API persists that target row and returns the saved row
- **AND** the response is scoped to that one target row rather than a bulk transaction

#### Scenario: Reject one invalid target row without touching others

- **WHEN** a client sends one invalid Step 4 target payload
- **THEN** the API returns a validation error for that target row
- **AND** previously saved workspace database rows remain unchanged

<!-- @trace
source: wire-studio-v2-step4-database-autosave
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/router.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/source_rule_handler.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - start.ps1
  - findings.md
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/api/handlers/mapping_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - .air.toml
  - internal/datalink/migrator.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - go.sum
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/workspace/service_database.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/workspace/memory_repository.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/schema/schema_device_models.go
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - scripts/start-log-utils.sh
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/hooks/datalink/index.ts
  - scripts/start-process-utils.sh
  - docs/goal.md
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/device/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/service_status.go
tests:
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_logger_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/device_handler_extended_test.go
-->

---
### Requirement: Studio v2 device availability status

The system SHALL include availability status in workspace and runtime-facing device payloads used by V2.

#### Scenario: Available device payload

- **WHEN** a workspace device is currently valid and usable
- **THEN** the payload includes `availability_status = "available"`

#### Scenario: Unavailable device payload

- **WHEN** a workspace device becomes invalid after editing
- **THEN** the payload includes `availability_status = "unavailable"`
- **AND** the payload includes `availability_reason`

<!-- @trace
source: stop-invalid-studio-v2-devices
updated: 2026-05-31
code:
  - internal/datalink/workspace/service_devices.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/dbtarget/types.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/router.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - go.sum
  - .air.toml
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/device/availability.go
  - findings.md
  - scripts/start-process-utils.sh
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/types/datalink.ts
  - docs/goal.md
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/hooks/datalink/index.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/device/errors.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - start.ps1
  - internal/api/handlers/runtime_handler.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/runtime/status.go
  - internal/datalink/device/service_crud.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/mapping/errors.go
  - gateway.db
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - scripts/start-log-utils.sh
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/sourcerule/service.go
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/workspace/service.go
  - task_plan.md
  - progress.md
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/schema/schema_device_models.go
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/datalink/workspace/service_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - cmd/test_ui/static/index.html
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/device/service_crud_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/device/availability_test.go
  - cmd/test_ui/main.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
-->

---
### Requirement: Studio v2 autosave responses expose runtime apply state

The system SHALL include runtime apply state in successful V2 autosave responses for device-affecting updates.

#### Scenario: Running device response reports applied

- **WHEN** a valid autosave change succeeds for a running workspace device
- **THEN** the success payload includes `runtime_apply_status = "applied"`

#### Scenario: Not-running device response reports not running

- **WHEN** a valid autosave change succeeds for a workspace device that has never been activated
- **THEN** the success payload includes `runtime_apply_status = "not_running"`

#### Scenario: Apply failure is distinct from save success

- **WHEN** the persisted save succeeds but runtime apply fails for a running workspace device
- **THEN** the success payload includes `runtime_apply_status = "apply_failed"`
- **AND** the payload includes a failure message explaining why runtime apply did not complete

<!-- @trace
source: apply-running-studio-v2-device-updates
updated: 2026-05-31
code:
  - frontend/src/hooks/datalink/keys.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/device/errors.go
  - go.sum
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/datalink/sourcerule/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/device/availability.go
  - task_plan.md
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - scripts/start-process-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/hooks/datalink/index.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/db.go
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/device/service_crud.go
  - progress.md
  - start.sh
  - internal/api/handlers/point_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/device/service_status.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/services/studioV2Rules.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/migrator.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - start.ps1
  - internal/api/router.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - docs/goal.md
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/runtime/status.go
  - internal/datalink/dbtarget/service.go
  - gateway.db
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/mapping/errors.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/runtime_handler.go
tests:
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_logger_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - cmd/test_ui/static/index.html
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->

---
### Requirement: Studio v2 workspace first activation API

The system SHALL expose `POST /api/v1/datalink/studio-v2/workspace/activate` as the first activation endpoint for Step 4.

#### Scenario: Activate eligible workspace devices

- **WHEN** a client requests `POST /api/v1/datalink/studio-v2/workspace/activate`
- **THEN** the API activates all valid, available, not-yet-running devices in the singleton workspace
- **AND** the response returns one result entry per attempted device

#### Scenario: Partial failure does not roll back successes

- **WHEN** one device activation succeeds and another fails in the same request
- **THEN** the response reports both outcomes separately
- **AND** the already successful activation remains active

##### Example: one success and one timeout

- **GIVEN** `device-A` is eligible and starts successfully while `device-B` is eligible but times out during activation
- **WHEN** the client calls `POST /api/v1/datalink/studio-v2/workspace/activate`
- **THEN** the response includes one `success` result for `device-A`
- **AND** one `failed` result for `device-B`
- **AND** `device-A` remains running after the request returns

<!-- @trace
source: replace-step4-with-first-activation
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/dbtarget/service.go
  - scripts/start-log-utils.sh
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - start.sh
  - internal/datalink/device/service_point_read.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - findings.md
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/api/handlers/studio_v2_workspace_handler.go
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/mapping/errors.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/workspace/service_database.go
  - .air.toml
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/runtime/status.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/point_direct_reader.go
  - task_plan.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/workspace/service_devices.go
  - go.sum
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - gateway.db
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - scripts/start-process-utils.sh
  - internal/datalink/workspace/service_activation.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/workspace/memory_repository.go
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - progress.md
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/migrator.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service.go
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/workspace/service_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
-->

---
### Requirement: Studio v2 workspace runtime context API

The system SHALL expose `GET /api/v1/datalink/studio-v2/workspace/runtime-context` as the workspace-scoped runtime context endpoint.

#### Scenario: Ordered runtime context payload

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace/runtime-context`
- **THEN** the API returns the workspace id, ordered device list, per-device availability fields, and `default_device_id`

#### Scenario: Empty runtime context payload

- **WHEN** the singleton workspace currently has no available devices
- **THEN** the API still returns HTTP 200
- **AND** the device list is empty or contains only unavailable devices
- **AND** the payload still includes enough metadata for the runtime page to render an empty state

<!-- @trace
source: add-studio-runtime-workspace-device-switching
updated: 2026-05-31
code:
  - progress.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/runtime/status.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/hooks/datalink/index.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - scripts/start-log-utils.sh
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/router.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - findings.md
  - frontend/src/types/studioV2Availability.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/datalink/sourcerule/service_links.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - docs/goal.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - go.sum
  - frontend/src/services/studioV2Rules.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/types/datalink.ts
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/sourcerule/activation_readiness.go
  - task_plan.md
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/datalink/mapping/errors.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/service_status.go
  - frontend/src/App.tsx
  - scripts/start-process-utils.sh
  - internal/datalink/device/service_point_read.go
  - internal/datalink/workspace/memory_repository.go
  - internal/datalink/device/availability.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/device/errors.go
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - start.sh
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/dbtarget/service.go
  - gateway.db
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - go.mod
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/db.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/workspace/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/datalink/sourcerule/errors.go
tests:
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/router_logger_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_runtime_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
-->

---
### Requirement: Studio v2 workspace bootstrap API

The system SHALL expose `GET /api/v1/datalink/studio-v2/workspace` as the bootstrap endpoint for the singleton v2 workspace.

#### Scenario: Bootstrap returns existing workspace

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace` and the singleton workspace already exists
- **THEN** the API returns HTTP 200 with the persisted workspace metadata
- **AND** the payload includes at least `id`, `kind`, `status`, `ordered_device_ids`, `created_at`, and `updated_at`

#### Scenario: Bootstrap auto-creates missing workspace

- **WHEN** a client requests `GET /api/v1/datalink/studio-v2/workspace` and the singleton workspace does not exist yet
- **THEN** the API creates the singleton workspace before responding
- **AND** the response still returns HTTP 200 with the created workspace metadata

<!-- @trace
source: add-studio-v2-single-workspace-foundation
updated: 2026-05-31
code:
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/device/availability.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - findings.md
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/types/datalink.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/router.go
  - internal/api/handlers/point_handler_polling.go
  - docs/goal.md
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - progress.md
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/services/studioV2Workspace.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/source_rule_handler.go
  - internal/datalink/migrator.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/device/sql_repo.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/runtime/status.go
  - scripts/start-process-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/datalink/device/service_status.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - .air.toml
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/workspace/sql_repository.go
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - start.ps1
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service_devices.go
  - go.sum
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - scripts/start-log-utils.sh
  - internal/datalink/sourcerule/service_links.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/api/handlers/point_handler.go
  - internal/datalink/db.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/datalink/mapping/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/index.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/mapping/errors.go
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/App.tsx
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/api/handlers/runtime_handler.go
tests:
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/router_runtime_test.go
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/datalink/workspace/service_device_order_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
-->

---
### Requirement: Studio v2 Step 1 diagnostics API

The system SHALL provide a real backend diagnostics contract that Studio V2 Step 1 can use for the current device draft.

#### Scenario: Successful diagnostics response

- **WHEN** a client requests Step 1 diagnostics for the current device draft
- **THEN** the API returns connect / probe outcomes produced by backend logic
- **AND** the response is sufficient for the UI to render success without local synthesis

#### Scenario: Failed diagnostics response

- **WHEN** backend diagnostics fails at connect or probe stage
- **THEN** the API returns the failing stage and actionable message
- **AND** the client does not need to guess failure state from a generic exception alone

<!-- @trace
source: replace-step1-mock-test-with-live-diagnostics
updated: 2026-05-31
code:
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/migrator.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/workspace/service_devices.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - findings.md
  - internal/api/handlers/point_handler.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/api/handlers/point_direct_reader.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - docs/goal.md
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/datalink/mapping/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/device/sql_repo.go
  - internal/datalink/sourcerule/service.go
  - internal/api/router.go
  - scripts/start-process-utils.sh
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/services/studioV2Workspace.ts
  - start.ps1
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/schema/schema_device_models.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/runtime_handler.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/device/errors.go
  - frontend/src/services/studioV2Mappings.ts
  - go.mod
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/datalink/dbtarget/types.go
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/mapping_handler.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/sourcerule/service_links.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/datalink/runtime/status.go
  - internal/datalink/db.go
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - scripts/start-log-utils.sh
  - go.sum
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - start.sh
  - .air.toml
  - gateway.db
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/workspace/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
tests:
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - cmd/test_ui/static/index.html
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/router_logger_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/device/availability_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_runtime_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
-->

---
### Requirement: Studio v2 settings API reuse

The system SHALL provide a backend contract that Studio V2 settings can use for boot and persistence without mock-only fallbacks.

#### Scenario: Read settings for V2 boot

- **WHEN** a client requests settings data for `/studio/v2/settings`
- **THEN** the API returns the persisted backend values needed for the settings surface

#### Scenario: Persist changed settings

- **WHEN** a client submits changed settings from `/studio/v2/settings`
- **THEN** the backend persists those values
- **AND** the client receives success or actionable failure rather than relying on console warnings


<!-- @trace
source: wire-studio-v2-settings-backend
updated: 2026-05-31
code:
  - scripts/start-process-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/services/studioV2Rules.ts
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/sourcerule/service.go
  - progress.md
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - go.sum
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/device/errors.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/runtime/status.go
  - .air.toml
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - gateway.db
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/mapping_handler.go
  - internal/api/router.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - start.ps1
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/mapping/service_crud.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/types/datalink.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - docs/goal.md
  - internal/datalink/db.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
tests:
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - cmd/test_ui/main.go
-->

---
### Requirement: Studio v2 connector pool API reuse

The system SHALL let Studio V2 settings use real connector CRUD and connector test APIs.

#### Scenario: Test one connector

- **WHEN** a client requests connector test for a settings connector
- **THEN** the backend performs a real connector test
- **AND** the response is sufficient for the UI to render success or failure without randomness

<!-- @trace
source: wire-studio-v2-settings-backend
updated: 2026-05-31
code:
  - scripts/start-process-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/services/studioV2Rules.ts
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/sourcerule/service.go
  - progress.md
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - go.sum
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/device/errors.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/runtime/status.go
  - .air.toml
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - gateway.db
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/mapping_handler.go
  - internal/api/router.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - start.ps1
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/mapping/service_crud.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/types/datalink.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - docs/goal.md
  - internal/datalink/db.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
tests:
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - cmd/test_ui/main.go
-->

---
### Requirement: Mapping preview request API

The system SHALL expose `POST /api/v1/datalink/mappings/preview` for draft transform preview requests used by `/studio/v2` Step 3 and similar clients.

#### Scenario: Preview a draft transform pipeline

- **WHEN** a client posts `raw_value` and `transform_pipeline` to `POST /api/v1/datalink/mappings/preview`
- **THEN** the API returns preview data containing `raw_value`, `final_value`, and ordered `step_results`
- **AND** the response can be rendered directly as transform preview output

#### Scenario: Preview failure returns actionable error

- **WHEN** preview execution fails for the submitted draft pipeline
- **THEN** the API response includes an actionable error message for the client
- **AND** the client can surface that failure without synthesizing a local success result

<!-- @trace
source: wire-studio-v2-step3-live-preview-and-target-type-shortcuts
updated: 2026-05-31
code:
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/errors.go
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/runtime/status.go
  - scripts/start-log-utils.sh
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_direct_reader.go
  - docs/goal.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/device/service_status.go
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - scripts/start-process-utils.sh
  - frontend/src/types/datalink.ts
  - internal/api/router.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - go.sum
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - gateway.db
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/device/availability.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/hooks/datalink/keys.ts
  - .air.toml
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/sourcerule/service_links.go
  - findings.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/api/handlers/mapping_handler.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/source_rule_handler.go
  - progress.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
tests:
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/router_runtime_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
-->