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

---
### Requirement: Runtime snapshot and stream expose projection reconciliation state

The runtime snapshot and stream contract SHALL expose whether each selected device is running the latest persisted workspace projection.

#### Scenario: Snapshot and stream carry projection status

- **WHEN** a client requests runtime snapshot or subscribes to runtime stream for a selected device
- **THEN** the response includes whether the device projection is aligned, deferred, or stale against the latest persisted workspace state
- **AND** the client SHALL NOT need to invent a projection-truth guess from unrelated counters

##### Example: stale snapshot reports its own projection status

- **GIVEN** runtime snapshot for dev-A still runs projection v12 while persisted workspace is at projection v13
- **WHEN** the client fetches snapshot or stream status for dev-A
- **THEN** the response marks dev-A as stale or equivalent instead of implying aligned truth

<!-- @trace
source: reconcile-runtime-with-workspace
updated: 2026-06-09
code:
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - internal/datalink/workspace/service_activation.go
  - tests/shell/start-frontend-readiness.sh
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - start.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/runtime/service_device_sync.go
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/dbtarget/live_projection.go
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/datalink/audit/sql_repository.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/datalink/audit/service.go
  - frontend/src/types/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/workspace/service_runtime_projection.go
  - frontend/studio-v2-diagnostic.cjs
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - internal/datalink/migrator_database_delivery_outcomes.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - Makefile
  - internal/datalink/migrator.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/App.tsx
  - internal/api/handlers/studio_v2_workspace_audit.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/runtime/status.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/schema/schema_dbtarget_models.go
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - tests/shell/start-frontend-install-failure.sh
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/dbtarget/tooling_service.go
  - frontend/src/types/studioV2Workspace.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/runtime_stream_handler.go
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - internal/datalink/workspace/service_readiness.go
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/types/databaseDelivery.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/studio-v2-real-check.cjs
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/dbtarget/writer.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - tests/shell/start-frontend-deps.sh
  - internal/datalink/runtime/service_workspace_projection.go
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - internal/datalink/audit/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/dbtarget/writer_statements.go
tests:
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - internal/api/router_runtime_database_delivery_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/datalink/runtime/service_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/runtime/ingestor_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/datalink/workspace/service_activation_projection_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/workspace/service_readiness_empty_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/audit/service_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/datalink/runtime/status_test.go
-->

---
### Requirement: Runtime backend expresses degraded and empty states explicitly

The runtime backend SHALL express empty, degraded, unavailable, and stale states explicitly in snapshot and stream responses.

#### Scenario: Backend marks stream unavailable explicitly

- **WHEN** runtime stream is unavailable for the selected device
- **THEN** the backend response or event model marks that unavailability explicitly
- **AND** the client SHALL NOT need to infer it from missing events alone

##### Example: selected stream outage is explicit

- **GIVEN** dev-A snapshot is available but the live stream channel for dev-A cannot attach
- **WHEN** the backend responds to the client
- **THEN** the contract marks stream unavailable for dev-A instead of silently returning no events

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
### Requirement: Runtime dashboard backend exposes diagnostics summary

The runtime dashboard backend SHALL expose diagnostics summary and latest failure context for the selected runtime scope.

#### Scenario: Backend returns diagnostics summary for selected device

- **WHEN** a client requests runtime monitoring data for a selected device
- **THEN** the backend response can include diagnostics summary with latest success and failure context for that device scope
- **AND** the client SHALL NOT need to parse raw server logs to present that summary

##### Example: diagnostics summary returns latest failure for dev-A

- **GIVEN** dev-A last failed at 10:12:00Z during database delivery and last succeeded at 10:10:00Z
- **WHEN** the client requests runtime monitoring data for dev-A
- **THEN** the response includes both timestamps and the failure context for dev-A

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