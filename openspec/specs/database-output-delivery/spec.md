# database-output-delivery Specification

## Purpose

TBD - created by archiving change 'guarantee-database-output-delivery'. Update Purpose after archive.

## Requirements

### Requirement: Database delivery operates only on the live enabled mapping set
The system SHALL build new schema plans and accept new recording inputs only from the current live enabled target set. Legacy immediate writes SHALL continue to use that live set. Records already accepted by the durable recording path SHALL retain their validated destination and definition snapshots and remain deliverable after a later disable or edit unless explicitly cancelled with an auditable affected-record decision. Hidden stale mappings MUST NOT create new delivery records or block unrelated live targets.

#### Scenario: Hidden stale target does not block live delivery
- **WHEN** a connector contains an old hidden mapping no longer tied to a live rule-owned point or tag
- **THEN** new schema planning and intake ignore that mapping
- **AND** unrelated live targets remain eligible for activation and delivery.

#### Scenario: Accepted records survive later disable
- **WHEN** valid durable records are pending and their source rule is subsequently disabled
- **THEN** no new samples are accepted for that disabled scope and previously accepted records remain pending for their original destination.

#### Scenario: Explicit cancellation
- **WHEN** an operator explicitly cancels pending records after reviewing scope and count
- **THEN** cancellation is audited and affected reports expose incomplete coverage rather than claiming successful delivery.

---
### Requirement: Database delivery records last schema and write outcomes

The system SHALL record the last schema ensure outcome and the last write delivery outcome for each relevant connector or target scope.

#### Scenario: Failed write updates delivery outcome

- **WHEN** runtime attempts a database write and the write fails
- **THEN** the system records the failed delivery outcome with timestamp and failure reason for the affected connector or target scope
- **AND** the next operator-facing surface can show that failure without scraping raw logs

##### Example: permission failure updates connector delivery truth

- **GIVEN** connector db-main is selected and a write fails with a permission error at 10:03:00Z
- **WHEN** the failed write outcome is recorded
- **THEN** db-main reports the failed delivery outcome with 10:03:00Z and the latest failure reason


<!-- @trace
source: guarantee-database-output-delivery
updated: 2026-06-09
code:
  - internal/datalink/workspace/service_readiness.go
  - internal/api/handlers/runtime_stream_handler.go
  - Makefile
  - internal/api/router.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - tests/shell/start-frontend-install-failure.sh
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/runtime/service_workspace_projection.go
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - internal/datalink/workspace/service_runtime_projection.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - tests/shell/start-frontend-deps.sh
  - internal/api/handlers/studio_v2_workspace_audit.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/schema/schema_dbtarget_models.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/studio-v2-diagnostic.cjs
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/migrator.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/audit/service.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/studio-v2-real-check.cjs
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - internal/datalink/runtime/truth_state.go
  - internal/datalink/dbtarget/live_projection.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/datalink/migrator_database_delivery_outcomes.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/source_rule_handler.go
  - internal/datalink/audit/types.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/service.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/types/studioV2RuntimeContext.ts
  - internal/datalink/runtime/ingestor.go
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/dbtarget/tooling_service.go
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - internal/datalink/audit/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/api/handlers/runtime_handler.go
  - frontend/src/types/databaseDelivery.ts
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - tests/shell/start-frontend-readiness.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/runtime/service_device_sync.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
tests:
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/datalink/runtime/ingestor_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/api/router_runtime_database_delivery_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/datalink/audit/service_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - internal/datalink/runtime/service_test.go
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/workspace/service_activation_projection_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_readiness_empty_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
-->

---
### Requirement: Runtime write delivery is diagnosable end to end

The system SHALL make the path from collected value to database write result diagnosable for the affected connector and target scope.

#### Scenario: Collected value reaches a failed write path

- **WHEN** runtime collects a value for a point whose database target write later fails
- **THEN** the system can report that the value was collected, mapped, and rejected by the database delivery path
- **AND** the operator SHALL NOT have to infer the failure only from missing rows in the external database

##### Example: collected point value fails only at DB write stage

- **GIVEN** pt-A is collected at 10:05:00Z, mapping succeeds, and dbtarget writer later fails for connector db-main
- **WHEN** the operator inspects delivery truth
- **THEN** the system reports pt-A as collected and mapped before failing in database delivery

<!-- @trace
source: guarantee-database-output-delivery
updated: 2026-06-09
code:
  - internal/datalink/workspace/service_readiness.go
  - internal/api/handlers/runtime_stream_handler.go
  - Makefile
  - internal/api/router.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - tests/shell/start-frontend-install-failure.sh
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/runtime/service_workspace_projection.go
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - internal/datalink/workspace/service_runtime_projection.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - tests/shell/start-frontend-deps.sh
  - internal/api/handlers/studio_v2_workspace_audit.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/schema/schema_dbtarget_models.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/studio-v2-diagnostic.cjs
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/migrator.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/audit/service.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/studio-v2-real-check.cjs
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - internal/datalink/runtime/truth_state.go
  - internal/datalink/dbtarget/live_projection.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/datalink/migrator_database_delivery_outcomes.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/source_rule_handler.go
  - internal/datalink/audit/types.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/service.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/types/studioV2RuntimeContext.ts
  - internal/datalink/runtime/ingestor.go
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/dbtarget/tooling_service.go
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - internal/datalink/audit/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/api/handlers/runtime_handler.go
  - frontend/src/types/databaseDelivery.ts
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - tests/shell/start-frontend-readiness.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/runtime/service_device_sync.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
tests:
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/datalink/runtime/ingestor_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/api/router_runtime_database_delivery_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/datalink/audit/service_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - internal/datalink/runtime/service_test.go
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/workspace/service_activation_projection_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_readiness_empty_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
-->

---
### Requirement: Automatic database creation on probe and initialization

When probing a database connector or establishing a target connection, if the target database does not exist (such as MySQL Error 1049 `ER_BAD_DB_ERROR` or PostgreSQL error `3D000`), the backend probe service SHALL attempt to automatically create the database on the server:
- For MySQL, connect without specifying the target database and execute `CREATE DATABASE IF NOT EXISTS <dbname>`.
- For PostgreSQL, connect to the administrative database (`postgres` or `template1`) and execute `CREATE DATABASE <dbname>`.
- Upon successful database creation, the backend SHALL re-establish the connection to the newly created database and complete verification.

#### Scenario: MySQL connection test with non-existent database
- **GIVEN** a MySQL server where database `gateway_metrics` does not yet exist
- **WHEN** the backend probes the connector configuration
- **THEN** the backend SHALL automatically execute database creation and return status `ready`

#### Scenario: Database creation error due to missing privileges
- **GIVEN** a database user without `CREATE DATABASE` privileges on the target server
- **WHEN** the backend attempts to auto-create the database and receives an access denied error
- **THEN** the probe SHALL fail with a clear authorization error message and SHALL NOT panic

<!-- @trace
source: fix-studio-v2-setup-and-db-flow
updated: 2026-09-06
code:
  - scripts/lib/b10_focused_suite.py
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/migrator_database_target_mysql_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/dbtarget/tooling_service.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/modbusshare/settings_lifecycle.go
  - internal/datalink/runtime/status.go
  - docs/swagger/swagger.yaml
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - internal/datalink/sourcerule/candidate_api.go
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/api/router.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/settings/errors.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - internal/datalink/modbusshare/service_projection.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - docs/swagger/docs.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/api/handlers/device_health_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/sourcerule/service_state.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/sourcerule/share_desired_mappings.go
  - internal/virtual/server/modbus/server.go
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/geometry.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/api/handlers/mapping_handler.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/datalink/modbusshare/types.go
  - internal/datalink/workspace/service_readiness_device.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/datalink/settings/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/device/repository_memory.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/typed_errors.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/collector/scheduler_dispatch.go
  - docs/swagger/swagger.json
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/hooks/previewStreamEvents.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - frontend/src/services/datalink.ts
  - internal/api/handlers/response_keys.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/datalink/modbusshare/errors.go
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/modbusshare/canonical_plan.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/template.go
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/connection.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/types/modbusShare.test-d.ts
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/types/sourceRuleCandidates.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/utils/typedErrors.ts
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/sourcerule/interfaces.go
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/modbusshare/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/vite.config.ts
  - scripts/lib/b10_acceptance_helpers.py
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/datalink/modbusshare/service_persistence.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/types/datalink.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/services/sourceRuleCandidates.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/dashboard_handler.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/tag/service_crud.go
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - internal/api/handlers/test.go
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - internal/datalink/modbusshare/sql_revision_store.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/api/handlers/settings_handler.go
  - internal/api/handlers/source_rule_handler_output_apply.go
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/datalink/device/service_readiness_safe.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/virtual/server/modbus/server_connection.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - internal/datalink/sourcerule/validation.go
  - internal/api/handlers/transport_wrapper.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - internal/api/router_modbus_share.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/sourcerule/tag_apply_service.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - internal/api/handlers/modbus_share_handler_status.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/services/datalinkClient.ts
  - internal/datalink/collector/scheduler.go
  - scripts/lib/b10_projection_assertions.py
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/virtual/server/modbus/server_register_handlers.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - scripts/check_file_lines.sh
  - internal/api/handlers/datalink_sse_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/types/modbusShare.ts
  - internal/datalink/dbtarget/service_validate.go
  - internal/api/handlers/modbus_share_handler.go
  - scripts/lib/b10_negative_steps.py
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/runtime/service.go
  - frontend/src/services/modbusShare.ts
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - go.mod
tests:
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/modbusshare/service_test.go
  - internal/datalink/modbusshare/status_contract_test.go
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - internal/datalink/device/service_readiness_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - internal/api/handlers/polling_group_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - scripts/lib/test_check_file_lines.py
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
  - internal/datalink/collector/scheduler_refactor_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/api/handlers/test_client_factory.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - internal/datalink/sourcerule/candidate_scope_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - cmd/test_ui/harness_config_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/datalink_sse_handler_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - cmd/test_ui/harness_config.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/modbusshare/reconciler_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - cmd/test_ui/main.go
  - internal/virtual/server/modbus/server_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/utils/safeJson.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/device/service_refactor_test.go
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/api/handlers/test_connection_handler.go
  - cmd/test_ui/share_startup.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - internal/api/handlers/test_monitor_handler.go
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - cmd/test_ui/server_runtime.go
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/api/device_handler_test.go
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - internal/api/modbus_share_swagger_contract_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - internal/api/handlers/test_request_helpers_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/migrator_test.go
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/sourcerule/share_gates_test.go
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
-->

---
### Requirement: Generated schema statements use dialect-valid syntax

Schema generation SHALL emit statements accepted by the target database dialect. For MySQL the unique-index statement SHALL NOT use the IF NOT EXISTS clause, which MySQL does not accept for index creation. For SQLite and PostgreSQL the existing IF NOT EXISTS form SHALL be retained. Duplicate index creation SHALL continue to be prevented by the inspected column metadata — a column already reported as a primary key or unique SHALL produce no index statement — and by deduplication of identical index names within one generation batch.

#### Scenario: MySQL unique index statement omits IF NOT EXISTS

- **GIVEN** a MySQL connector with an upsert mapping whose timestamp column exists but is neither a primary key nor unique
- **WHEN** schema generation runs
- **THEN** the generated statement creates a unique index without the IF NOT EXISTS clause
- **AND** executing the generated statements against MySQL succeeds

#### Scenario: SQLite and PostgreSQL keep the guarded form

- **GIVEN** a SQLite or PostgreSQL connector with an upsert mapping whose timestamp column is neither a primary key nor unique
- **WHEN** schema generation runs
- **THEN** the generated statement creates a unique index using the IF NOT EXISTS clause

#### Scenario: An already-unique timestamp column produces no index statement

- **GIVEN** a connector whose inspected timestamp column is reported as a primary key or unique
- **WHEN** schema generation runs
- **THEN** no unique-index statement is generated for that column

##### Example: unique index statement by dialect and column state

| Dialect | Timestamp column state | Index statement generated |
| ------- | ---------------------- | ------------------------- |
| MySQL | plain column | CREATE UNIQUE INDEX without IF NOT EXISTS |
| MySQL | primary key or unique | none |
| PostgreSQL | plain column | CREATE UNIQUE INDEX with IF NOT EXISTS |
| SQLite | plain column | CREATE UNIQUE INDEX with IF NOT EXISTS |


<!-- @trace
source: fix-activation-barrier-and-mysql-target-defects
updated: 2026-09-06
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/migrator_database_target_mysql_schema.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/device/service_crud.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service.go
  - internal/datalink/dbtarget/tooling_service.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/datalink/device/repository_memory.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
tests:
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/device/service_refactor_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/api/device_handler_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
-->

---
### Requirement: Upsert mappings require a guaranteed unique key on MySQL

MySQL upsert delivery relies on a unique or primary key covering the timestamp column; without it the upsert clause silently degrades to append-only inserts. Schema generation SHALL guarantee that key for every enabled MySQL upsert mapping, either by reporting the column as already unique or by emitting a statement that creates the unique index. When the index cannot be created, schema generation SHALL fail and SHALL surface the underlying database error rather than completing successfully.

#### Scenario: MySQL upsert mapping ends with a unique key present

- **GIVEN** an enabled MySQL upsert mapping against an existing table whose timestamp column has no unique key
- **WHEN** schema generation completes successfully
- **THEN** the timestamp column is covered by a unique index

#### Scenario: Unique index creation failure fails schema generation

- **GIVEN** an existing MySQL table containing duplicate values in the timestamp column
- **WHEN** schema generation attempts to create the unique index
- **THEN** schema generation is recorded as failed
- **AND** the underlying database error is preserved in the reported failure


<!-- @trace
source: fix-activation-barrier-and-mysql-target-defects
updated: 2026-09-06
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/migrator_database_target_mysql_schema.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/device/service_crud.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service.go
  - internal/datalink/dbtarget/tooling_service.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/datalink/device/repository_memory.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
tests:
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/device/service_refactor_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/api/device_handler_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
-->

---
### Requirement: Schema generation resolves the default schema from the connector

When a target mapping record carries an empty table schema, schema generation SHALL resolve the default schema from the connector, using the same resolution the mapping create and update paths use. For MySQL that resolution SHALL yield the database name from the connector connection configuration. Schema generation SHALL NOT fall back to a kind-only default that cannot match the inspected table set.

#### Scenario: MySQL mapping with an empty table schema targets the connection database

- **GIVEN** a MySQL connector whose connection configuration names database gateway_metrics
- **AND** an enabled mapping whose table schema is empty
- **WHEN** schema generation runs
- **THEN** the generated statements qualify the table with gateway_metrics
- **AND** no generated statement qualifies a table with the kind-only default schema name

#### Scenario: Mapping with an empty schema matches existing inspected tables

- **GIVEN** a MySQL connector whose target table already exists in the connection database
- **AND** an enabled mapping for that table whose table schema is empty
- **WHEN** schema generation runs
- **THEN** the existing table is matched and no create-table statement is generated for it

##### Example: resolved schema by connector kind

| Connector kind | Connection configuration | Mapping table schema | Resolved schema |
| -------------- | ------------------------ | -------------------- | --------------- |
| MySQL | database=gateway_metrics | empty | gateway_metrics |
| MySQL | database=gateway_metrics | reporting | reporting |
| PostgreSQL | not applicable | empty | public |


<!-- @trace
source: fix-activation-barrier-and-mysql-target-defects
updated: 2026-09-06
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/migrator_database_target_mysql_schema.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/device/service_crud.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service.go
  - internal/datalink/dbtarget/tooling_service.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/datalink/device/repository_memory.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
tests:
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/device/service_refactor_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/api/device_handler_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
-->

---
### Requirement: MySQL connections fail closed rather than sending cleartext passwords

MySQL connection assembly SHALL NOT enable cleartext password authentication by default, and SHALL NOT silently fall back to an unencrypted connection when cleartext authentication is enabled. The same rules SHALL apply to both the target connection and the administrative connection used by connector probing, from a single shared resolution.

The resolution SHALL be:

- An explicit TLS value in the connection configuration SHALL be used as given.
- Otherwise, when TLS is requested through the use-TLS option, the TLS mode SHALL be an encrypting mode that does not permit falling back to plaintext.
- Otherwise, when cleartext authentication is explicitly enabled, the TLS mode SHALL be an encrypting mode that does not permit falling back to plaintext.
- Otherwise, no TLS mode SHALL be forced.

#### Scenario: Default configuration does not enable cleartext authentication

- **GIVEN** a MySQL connector configuration that does not set the cleartext password option
- **WHEN** the connection descriptor is assembled
- **THEN** cleartext password authentication is disabled

#### Scenario: Explicitly enabled cleartext requires an encrypted transport

- **GIVEN** a MySQL connector configuration that explicitly enables cleartext password authentication and sets no explicit TLS value
- **WHEN** the connection descriptor is assembled
- **THEN** the TLS mode is an encrypting mode
- **AND** falling back to an unencrypted connection is not permitted
- **AND** a server without TLS causes the connection to fail with the existing connection error classification

#### Scenario: An explicit TLS value is honoured

- **GIVEN** a MySQL connector configuration that sets an explicit TLS value
- **WHEN** the connection descriptor is assembled
- **THEN** that TLS value is used unchanged

#### Scenario: Probe and target connections share the resolution

- **GIVEN** any MySQL connector configuration
- **WHEN** both the target connection descriptor and the probe administrative connection descriptor are assembled
- **THEN** both carry the same cleartext and TLS resolution

##### Example: TLS and cleartext resolution

| Explicit TLS value | use-TLS requested | cleartext explicitly enabled | Cleartext | TLS mode | Plaintext fallback |
| ------------------ | ----------------- | ---------------------------- | --------- | -------- | ------------------ |
| none | no | no | disabled | none forced | not applicable |
| none | yes | no | disabled | encrypting | not permitted |
| none | no | yes | enabled | encrypting | not permitted |
| set | any | any | as configured | the explicit value | as implied by that value |


<!-- @trace
source: fix-activation-barrier-and-mysql-target-defects
updated: 2026-09-06
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/migrator_database_target_mysql_schema.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/device/service_crud.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service.go
  - internal/datalink/dbtarget/tooling_service.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/datalink/device/repository_memory.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
tests:
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/device/service_refactor_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/api/device_handler_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
-->

---
### Requirement: Connector probing attempts database creation once per kind

Connector probing SHALL attempt automatic database creation at most once per probe for a given connector kind. For MySQL, where the connection helper already performs create-then-retry internally, the probe SHALL NOT contain a second create-and-retry branch.

#### Scenario: MySQL probe does not retry database creation twice

- **GIVEN** a MySQL connector whose target database does not exist
- **WHEN** the connector is probed
- **THEN** database creation is attempted once
- **AND** the connection is retried once after creation

#### Scenario: A non-missing-database MySQL error is not treated as a creation candidate

- **GIVEN** a MySQL connector that fails to connect for a reason other than a missing database
- **WHEN** the connector is probed
- **THEN** no database creation is attempted
- **AND** the original failure is classified and reported

<!-- @trace
source: fix-activation-barrier-and-mysql-target-defects
updated: 2026-09-06
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/migrator_database_target_mysql_schema.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - internal/datalink/dbtarget/writer.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/dbtarget/service_probe.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/device/service_crud.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - internal/datalink/dbtarget/tooling_service_helpers.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service.go
  - internal/datalink/dbtarget/tooling_service.go
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - internal/datalink/device/repository_memory.go
  - internal/datalink/dbtarget/service_mysql_inspection.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
tests:
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/settings-save-state-convergence.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/autosave-settlement-timeout.test.ts
  - internal/datalink/dbtarget/service_mysql_inspection_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/datalink/dbtarget/tooling_service_mysql_schema_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx
  - internal/datalink/migrator_database_target_mysql_schema_test.go
  - internal/datalink/dbtarget/writer_statements_mysql_test.go
  - internal/datalink/device/service_refactor_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/api/device_handler_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
-->
