# runtime-dashboard-backend Specification

## Purpose

TBD - created by archiving change 'add-runtime-dashboard-backend-contract'. Update Purpose after archive.

## Requirements

### Requirement: Device-scoped runtime monitoring contract

The system SHALL provide a device-scoped runtime monitoring contract so a post-setup client can monitor one committed device without subscribing to fleet-wide point traffic.

#### Scenario: Request snapshot for one device

- **WHEN** a client requests the runtime snapshot with `device_id=device-A`
- **THEN** the system returns runtime metrics plus exactly one collector summary for `device-A`
- **AND** the response does not require the client to filter unrelated devices locally

#### Scenario: Subscribe with point filter inside one device

- **WHEN** a client opens the runtime stream with `device_id=device-A&point_ids=pt-1,pt-2`
- **THEN** the system emits live `value` events only for `pt-1` and `pt-2`
- **AND** the stream remains scoped to `device-A`


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
### Requirement: Snapshot and stream use one derived device status model

The system SHALL derive runtime device status consistently across the snapshot response and stream `status` events using the same freshness, point error, and breaker-state rules.

#### Scenario: Stale points produce warning status consistently

- **WHEN** the latest read time for one or more enabled points exceeds the effective freshness window for a device
- **THEN** the runtime snapshot reports the device status as `warning`
- **AND** the next emitted `status` event for that device reports the same `warning` status with matching stale-point counts

##### Example: stale counts stay aligned

| Source | points_total | points_healthy | points_stale | points_error | breaker_state | derived status |
| ----- | ------------ | -------------- | ------------ | ------------ | ------------- | -------------- |
| snapshot | 4 | 3 | 1 | 0 | closed | warning |
| status event | 4 | 3 | 1 | 0 | closed | warning |

#### Scenario: Open breaker produces error status consistently

- **WHEN** a device breaker state becomes `open`
- **THEN** the runtime snapshot reports the device status as `error`
- **AND** the stream emits a `status` event with `breaker_state=open` and `status=error`


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
### Requirement: Live monitoring stream remains loss-tolerant

The system SHALL protect runtime collection progress when a monitoring subscriber is slow, even if some intermediate live events cannot be delivered to that subscriber.

#### Scenario: Slow subscriber does not block collection

- **WHEN** a monitoring subscriber stops consuming events and its stream buffer is saturated
- **THEN** runtime collection continues without waiting for that subscriber
- **AND** the system preserves connection availability for other subscribers and future snapshot requests

##### Example: saturated subscriber buffer

| Subscriber state | Incoming value events | Expected collector behavior | Expected slow-subscriber behavior |
| ----- | --------------------- | --------------------------- | --------------------------------- |
| buffer full | 10 new values for `device-A` | collection loop continues processing all 10 values | some intermediate SSE deliveries are dropped |

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