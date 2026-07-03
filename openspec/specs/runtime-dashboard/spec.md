# runtime-dashboard Specification

## Purpose

TBD - created by archiving change 'runtime-dashboard'. Update Purpose after archive.

## Requirements

### Requirement: Focused runtime dashboard surface
The system SHALL render a focused post-setup runtime dashboard surface for one selected device, without replacing the `/studio/v2` setup shell or introducing a fleet-wide card wall.

#### Scenario: Render focused device monitoring page
- **WHEN** the runtime dashboard receives a selected device context and a successful runtime snapshot
- **THEN** the page renders a focused device header, runtime summary, collector health panel, and live point values table
- **AND** the page does not require queue backlog, diagnostic logs, or a global navigation card wall to become usable

##### Example: Focused page for one committed device
- **GIVEN** the selected device is `device-A` and the runtime snapshot for `device-A` loaded successfully
- **WHEN** the runtime dashboard page renders
- **THEN** the operator can read one page centered on `device-A` instead of a multi-device overview


<!-- @trace
source: runtime-dashboard
updated: 2026-05-30
code:
  - frontend/src/types/datalink.ts
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/test-tooling.md
  - task_plan.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/i18n/config.ts
  - docs/technical/studio-surface-inventory/inventory.js
  - internal/datalink/runtime/service.go
  - CLAUDE.md
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - progress.md
  - internal/api/handlers/runtime_handler.go
  - docs/technical/studio-surface-inventory/index.html
  - internal/datalink/runtime/stream.go
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/README.md
  - internal/datalink/runtime/status.go
  - cmd/studio_inventory_changelog/main.go
  - findings.md
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/App.tsx
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - AGENTS.md
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - docs/technical/studio-surface-inventory/gap-roadmap.md
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
-->

---
### Requirement: Runtime summary presents backend-supported metrics
The runtime dashboard SHALL present only metrics and summary fields guaranteed by the runtime backend contract.

#### Scenario: Render summary from runtime snapshot
- **WHEN** the runtime dashboard receives a runtime snapshot for the selected device
- **THEN** the summary panel shows the backend-supported runtime metrics and collector summary values
- **AND** the page does not infer unsupported backlog or log counters from unrelated endpoints

##### Example: Summary uses only formal contract fields
- **GIVEN** the runtime snapshot includes `collected_total=128`, `write_success_total=128`, `write_error_total=0`, and collector `status=healthy`
- **WHEN** the summary panel renders
- **THEN** it shows those runtime metrics and collector summary values without inventing queue backlog or log counters


<!-- @trace
source: runtime-dashboard
updated: 2026-05-30
code:
  - frontend/src/types/datalink.ts
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/test-tooling.md
  - task_plan.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/i18n/config.ts
  - docs/technical/studio-surface-inventory/inventory.js
  - internal/datalink/runtime/service.go
  - CLAUDE.md
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - progress.md
  - internal/api/handlers/runtime_handler.go
  - docs/technical/studio-surface-inventory/index.html
  - internal/datalink/runtime/stream.go
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/README.md
  - internal/datalink/runtime/status.go
  - cmd/studio_inventory_changelog/main.go
  - findings.md
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/App.tsx
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - AGENTS.md
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - docs/technical/studio-surface-inventory/gap-roadmap.md
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
-->

---
### Requirement: Collector health panel
The runtime dashboard SHALL show collector health for the selected device using the runtime status fields defined by the backend contract.

#### Scenario: Show collector health counters and status
- **WHEN** the runtime snapshot or aligned status update reports `status`, `points_total`, `points_healthy`, `points_stale`, `points_error`, `last_read_at`, and `breaker_state`
- **THEN** the collector health panel renders those fields for the selected device
- **AND** the panel uses the same status vocabulary as the backend contract instead of introducing page-only health states


<!-- @trace
source: runtime-dashboard
updated: 2026-05-30
code:
  - frontend/src/types/datalink.ts
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/test-tooling.md
  - task_plan.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/i18n/config.ts
  - docs/technical/studio-surface-inventory/inventory.js
  - internal/datalink/runtime/service.go
  - CLAUDE.md
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - progress.md
  - internal/api/handlers/runtime_handler.go
  - docs/technical/studio-surface-inventory/index.html
  - internal/datalink/runtime/stream.go
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/README.md
  - internal/datalink/runtime/status.go
  - cmd/studio_inventory_changelog/main.go
  - findings.md
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/App.tsx
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - AGENTS.md
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - docs/technical/studio-surface-inventory/gap-roadmap.md
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
-->

---
### Requirement: Live point values table
The runtime dashboard SHALL present a live point values table for the selected device, including raw values, transformed values, quality or stale state, and timestamps.

#### Scenario: Show live values after snapshot-first load
- **WHEN** the selected device snapshot has loaded but live point events have not arrived yet
- **THEN** the page keeps the runtime summary visible
- **AND** the live point values table shows a waiting or placeholder state instead of disappearing

##### Example: Snapshot visible before first live point event
- **GIVEN** `device-A` snapshot loaded at `10:00:00Z` and no live point events have arrived yet
- **WHEN** the page enters its initial post-snapshot state
- **THEN** the runtime summary remains visible and the live values area shows a waiting placeholder

#### Scenario: Update table with live values
- **WHEN** live point value updates arrive for the selected device
- **THEN** the table updates the corresponding rows with raw value, transformed value, quality or stale state, and timestamp

##### Example: Raw and transformed values update one row
- **GIVEN** the table already contains point `pt-1` for `device-A`
- **WHEN** a live update arrives with `raw_value=150`, `transformed_value=20.0`, `quality=good`, and timestamp `2026-05-29T10:00:05Z`
- **THEN** the `pt-1` row shows `150`, `20.0`, `good`, and `2026-05-29T10:00:05Z`


<!-- @trace
source: runtime-dashboard
updated: 2026-05-30
code:
  - frontend/src/types/datalink.ts
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/test-tooling.md
  - task_plan.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/i18n/config.ts
  - docs/technical/studio-surface-inventory/inventory.js
  - internal/datalink/runtime/service.go
  - CLAUDE.md
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - progress.md
  - internal/api/handlers/runtime_handler.go
  - docs/technical/studio-surface-inventory/index.html
  - internal/datalink/runtime/stream.go
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/README.md
  - internal/datalink/runtime/status.go
  - cmd/studio_inventory_changelog/main.go
  - findings.md
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/App.tsx
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - AGENTS.md
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - docs/technical/studio-surface-inventory/gap-roadmap.md
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
-->

---
### Requirement: Degraded runtime monitoring banner
The runtime dashboard SHALL clearly indicate when live runtime monitoring is degraded while preserving the last successful runtime data on screen.

#### Scenario: Live stream disconnect preserves last data
- **WHEN** the live stream disconnects after the page has already rendered runtime summary and point values
- **THEN** the page shows a degraded monitoring banner
- **AND** the last successful runtime summary and point values remain visible until fresher data arrives

##### Example: Disconnect after live data keeps last snapshot on screen
- **GIVEN** the page is already showing runtime summary and the latest values for `device-A`
- **WHEN** the live stream disconnects
- **THEN** the page switches to degraded monitoring while keeping the last successful summary and point values visible

<!-- @trace
source: runtime-dashboard
updated: 2026-05-30
code:
  - frontend/src/types/datalink.ts
  - frontend/src/services/datalink.ts
  - docs/technical/studio-surface-inventory/test-tooling.md
  - task_plan.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/i18n/config.ts
  - docs/technical/studio-surface-inventory/inventory.js
  - internal/datalink/runtime/service.go
  - CLAUDE.md
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - progress.md
  - internal/api/handlers/runtime_handler.go
  - docs/technical/studio-surface-inventory/index.html
  - internal/datalink/runtime/stream.go
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/README.md
  - internal/datalink/runtime/status.go
  - cmd/studio_inventory_changelog/main.go
  - findings.md
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/App.tsx
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - AGENTS.md
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - docs/technical/studio-surface-inventory/gap-roadmap.md
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
-->

---
### Requirement: Runtime dashboard does not invent unsupported runtime values

The runtime dashboard SHALL NOT invent runtime values, statuses, or summaries that are not supplied by the backend contract.

#### Scenario: Unsupported field stays absent instead of guessed

- **WHEN** a dashboard panel does not receive a required backend-supported field
- **THEN** that panel renders an empty or degraded state for that field
- **AND** it SHALL NOT guess a replacement from unrelated client state

##### Example: missing write error metric stays absent

- **GIVEN** the backend response for dev-A omits a write-error detail field that one panel expects
- **WHEN** the panel renders
- **THEN** it shows that field as unavailable instead of deriving a substitute from another counter

<!-- @trace
source: make-runtime-dashboard-truthful
updated: 2026-06-09
code:
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - internal/datalink/schema/schema_dbtarget_models.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/runtime/ingestor.go
  - internal/datalink/workspace/service_runtime_projection.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/audit/service.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - internal/datalink/runtime/service_projection_state.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - start.sh
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - tests/shell/start-frontend-install-failure.sh
  - internal/api/handlers/studio_v2_workspace_audit.go
  - internal/datalink/audit/types.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/runtime/service_workspace_projection.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - internal/datalink/runtime/service_device_sync.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/dbtarget/service_validate.go
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - frontend/e2e-studio-v2-live.cjs
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/workspace/service_readiness.go
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - tests/shell/start-frontend-deps.sh
  - Makefile
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/datalink/runtime/truth_state.go
  - internal/datalink/dbtarget/tooling_service.go
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/runtime/status.go
  - frontend/src/types/datalink.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/runtime/service.go
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - internal/datalink/dbtarget/writer.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/App.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/dbtarget/live_projection.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/types/runtimeTruth.ts
  - internal/api/handlers/runtime_stream_handler.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - frontend/src/types/runtimeDiagnostics.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/runtime_workspace_setup_types.go
  - internal/api/router.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/dbtarget/writer_statements.go
  - tests/shell/start-frontend-readiness.sh
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/audit/sql_repository.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/types/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/datalink/workspace/service.go
  - internal/api/handlers/runtime_workspace_setup_context.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/types/databaseDelivery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/migrator_database_delivery_outcomes.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/src/types/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/studio-v2-real-check.cjs
  - internal/datalink/runtime/service_source_rule_reconcile.go
tests:
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/datalink/audit/service_test.go
  - internal/api/router_runtime_database_delivery_test.go
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/datalink/workspace/service_readiness_empty_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/api/router_runtime_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - internal/api/handlers/runtime_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/runtime/ingestor_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/datalink/workspace/service_activation_projection_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/runtime/service_test.go
-->

---
### Requirement: Runtime dashboard shows operator-facing diagnostics

The runtime dashboard SHALL show operator-facing diagnostics for recent runtime and delivery failures.

#### Scenario: Dashboard shows latest failure context

- **WHEN** runtime or database delivery has a recent failure for the selected device scope
- **THEN** the dashboard shows the latest failure context and timestamp for that selected scope
- **AND** the operator can see more than aggregate counters alone

##### Example: dashboard shows latest DB delivery failure

- **GIVEN** dev-A has a recent failed database delivery at 10:12:00Z
- **WHEN** the runtime dashboard opens for dev-A
- **THEN** the dashboard shows that latest delivery failure context and timestamp for dev-A

<!-- @trace
source: add-runtime-observability-and-audit
updated: 2026-06-09
code:
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/audit/sql_repository.go
  - internal/datalink/dbtarget/tooling_service.go
  - internal/datalink/workspace/service_runtime_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - Makefile
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/dbtarget/writer_statements.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/audit/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_audit.go
  - tests/shell/start-frontend-readiness.sh
  - start.sh
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - internal/datalink/dbtarget/delivery_outcome.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/datalink/audit/service.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/dbtarget/live_projection.go
  - internal/datalink/runtime/service_workspace_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - frontend/e2e-studio-v2-live.cjs
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - internal/datalink/runtime/service.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - tests/shell/start-frontend-install-failure.sh
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/types/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/databaseDelivery.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/studio-v2-real-check.cjs
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - internal/datalink/runtime/status.go
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/api/handlers/runtime_workspace_setup_types.go
  - internal/datalink/workspace/service_readiness.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/migrator_database_delivery_outcomes.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/runtime/service_device_sync.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/schema/schema_dbtarget_models.go
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/App.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - tests/shell/start-frontend-deps.sh
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/datalink/dbtarget/writer.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/types/datalink.ts
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
tests:
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - cmd/test_ui/main.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/datalink/audit/service_test.go
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_activation_projection_test.go
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/datalink/runtime/service_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/router_runtime_database_delivery_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/workspace/service_readiness_empty_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/datalink/runtime/ingestor_test.go
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/runtime_handler_test.go
-->