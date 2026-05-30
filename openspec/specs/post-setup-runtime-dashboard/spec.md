# post-setup-runtime-dashboard Specification

## Purpose

TBD - created by archiving change 'add-post-setup-runtime-dashboard-route'. Update Purpose after archive.

## Requirements

### Requirement: Dedicated post-setup runtime dashboard route

The system SHALL expose a dedicated runtime dashboard route for post-setup monitoring without replacing the `/studio/v2` setup entry flow.

#### Scenario: Open runtime dashboard by URL

- **WHEN** an operator navigates to `/studio/runtime?device_id=device-A`
- **THEN** the system renders the runtime dashboard page
- **AND** the system does not redirect to `/studio` or `/studio/v2`

#### Scenario: Setup flow remains separate

- **WHEN** an operator navigates to `/studio/v2`
- **THEN** the system continues to render the setup workbench flow
- **AND** the runtime dashboard does not replace the v2 landing page


<!-- @trace
source: add-post-setup-runtime-dashboard-route
updated: 2026-05-30
code:
  - AGENTS.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - docs/technical/studio-surface-inventory/gap-roadmap.md
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - task_plan.md
  - CLAUDE.md
  - progress.md
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - cmd/studio_inventory_changelog/main.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/runtime/stream.go
  - frontend/src/i18n/config.ts
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - docs/technical/studio-surface-inventory/test-tooling.md
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - docs/technical/studio-surface-inventory/START_HERE.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - docs/technical/studio-surface-inventory/index.html
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - docs/technical/studio-surface-inventory/inventory.js
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - findings.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/README.md
  - frontend/src/services/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
-->

---
### Requirement: Runtime dashboard uses device-focused context

The runtime dashboard SHALL treat `device_id` as the canonical selected-device context.

#### Scenario: Missing device context shows focused empty state

- **WHEN** an operator opens `/studio/runtime` without `device_id`
- **THEN** the system shows a missing-device-context state
- **AND** the page does not silently fall back to a fleet-wide dashboard

#### Scenario: Switching device updates route context

- **WHEN** an operator switches the selected device from `device-A` to `device-B` inside the runtime dashboard
- **THEN** the system updates the route query to `device_id=device-B`
- **AND** the dashboard reloads its runtime data using `device-B` as the only selected device context


<!-- @trace
source: add-post-setup-runtime-dashboard-route
updated: 2026-05-30
code:
  - AGENTS.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - docs/technical/studio-surface-inventory/gap-roadmap.md
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - task_plan.md
  - CLAUDE.md
  - progress.md
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - cmd/studio_inventory_changelog/main.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/runtime/stream.go
  - frontend/src/i18n/config.ts
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - docs/technical/studio-surface-inventory/test-tooling.md
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - docs/technical/studio-surface-inventory/START_HERE.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - docs/technical/studio-surface-inventory/index.html
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - docs/technical/studio-surface-inventory/inventory.js
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - findings.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/README.md
  - frontend/src/services/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
-->

---
### Requirement: Runtime dashboard layers snapshot, live stream, and degraded fallback

The runtime dashboard SHALL load runtime data by fetching a snapshot first, layering live stream updates on top, and degrading to snapshot polling when the live stream is unavailable.

#### Scenario: Snapshot renders before live stream attaches

- **WHEN** an operator opens `/studio/runtime?device_id=device-A`
- **THEN** the system first fetches the runtime snapshot for `device-A`
- **AND** the page renders the latest available runtime summary before live stream events arrive

#### Scenario: Stream disconnect degrades to snapshot polling

- **WHEN** the runtime live stream for `device-A` disconnects after the initial snapshot succeeds
- **THEN** the system preserves the last successful runtime summary on screen
- **AND** the page enters a degraded state that continues refreshing the snapshot until live streaming recovers

<!-- @trace
source: add-post-setup-runtime-dashboard-route
updated: 2026-05-30
code:
  - AGENTS.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - docs/technical/studio-surface-inventory/gap-roadmap.md
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - task_plan.md
  - CLAUDE.md
  - progress.md
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - cmd/studio_inventory_changelog/main.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/runtime/stream.go
  - frontend/src/i18n/config.ts
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - docs/technical/studio-surface-inventory/test-tooling.md
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - docs/technical/studio-surface-inventory/START_HERE.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - docs/technical/studio-surface-inventory/index.html
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - docs/technical/studio-surface-inventory/inventory.js
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - findings.md
  - docs/technical/studio-surface-inventory/context.json
  - docs/technical/studio-surface-inventory/README.md
  - frontend/src/services/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/router_runtime_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/runtime/stream_test.go
-->