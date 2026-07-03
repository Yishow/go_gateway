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

---
### Requirement: Post-setup runtime route shows truthful empty and degraded states

The post-setup runtime route SHALL show truthful empty and degraded states when runtime data is missing or unavailable.

#### Scenario: Route opens without real runtime data

- **WHEN** an operator opens the runtime route for a selected device that has no current runtime data yet
- **THEN** the page shows a truthful waiting, empty, or degraded state
- **AND** the page SHALL NOT present a synthetic ready dashboard just because the route resolved successfully

##### Example: direct link resolves but selected device is still cold

- **GIVEN** the operator opens /studio/runtime for dev-A immediately after activation but runtime has not produced real data yet
- **WHEN** the route resolves successfully
- **THEN** the page shows waiting or empty state for dev-A instead of a synthetic ready dashboard

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