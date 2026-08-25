# datalink-ui Specification

## Purpose

TBD - created by archiving change add-device-data-pipeline. Update Purpose after archive.

## Requirements

### Requirement: Guided workflow

The UI SHALL provide a guided workflow that keeps the complete datalink product flow inside /studio/v2 as the canonical user-facing setup workspace: Device -> SourceRule -> Tag review -> Output review/apply. The workflow SHALL remain independent of the removed `/studio` workspace and SHALL NOT make `/studio` a product or compatibility workspace.

The workflow SHALL ensure:

1. Operators establish device capability context before source-rule planning.
2. Tag and output review remain bound to the active source-rule revision.
3. Manual point, tag, or mapping construction tools are secondary and SHALL NOT be required to complete the normal product path.
4. Database and Local Modbus output readiness remain visible in the same workspace.
5. Failures remain localized to the owning step, selection, or output target and use safe localized copy.
6. Step 2 SHALL surface the current active source rule as the primary planning unit, so the operator can immediately understand what is currently planned before moving into Tag review.
7. When the active output target is database, Step 2 SHALL reflect database-aware planning guidance without turning Step 2 into connector or schema setup.
8. A legacy `/studio` bookmark after deletion SHALL follow the existing generic unknown-route policy and SHALL NOT mount a legacy workspace or require query migration.

#### Scenario: End-to-end guided configuration uses v2

- **WHEN** an operator selects a device, saves a source rule, reviews tag candidates, and applies one or both output targets
- **THEN** the primary workflow stays inside /studio/v2
- **AND** the operator does not need a separate manual point-first or mapping-first route to complete the normal path

#### Scenario: Source step exposes the active rule before downstream review

- **WHEN** an operator opens Step 2 with an active source rule
- **THEN** the UI shows the current rule planning summary before the operator moves to Tag review
- **AND** the operator can identify the rule address coverage, type intent, and planning state without reading the full canvas first

#### Scenario: Deleted legacy bookmark does not re-enter legacy workspace

- **WHEN** an operator follows a `/studio` bookmark after deletion
- **THEN** the browser follows the same generic unknown-route policy as an arbitrary unknown path
- **AND** no legacy workspace, query migration, or legacy-only chunk is mounted


<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->

---
### Requirement: Drag-drop mapping canvas
The UI MAY provide a drag-drop mapping canvas as a secondary engineering tool, but the primary `/studio` workflow SHALL NOT depend on manual drag-drop point-to-tag mapping.

#### Scenario: Primary workflow completes without drag-drop mapping
- **WHEN** an operator completes the normal `/studio` workflow
- **THEN** the system does not require opening a drag-drop mapping canvas to create the effective rule-driven tag or output state

---
### Requirement: Transform builder

The UI SHALL provide a transform builder with ordered steps, parameters, and validation.

The builder SHALL support the following transform types:

- `decode`: Decode raw bytes to typed value
- `cast`: Type conversion
- `scale`: Linear scaling with parameters (multiplier, offset)
- `lookup`: Table lookup mapping
- `conditional`: Conditional branching
- `formula`: Expression-based calculation

#### Scenario: Configure scaling step

- WHEN a user adds a scaling step with parameters
- THEN the UI validates and saves the step

#### Scenario: Reorder transform steps

- WHEN a user drags a step to a new position
- THEN the pipeline order is updated

#### Scenario: Pipeline validation

- WHEN a user saves the transform pipeline
- THEN the UI calls `/mappings/validate-pipeline` to verify

---
### Requirement: Live preview

The UI SHALL provide a live preview of raw and transformed values for a selected mapping using the server preview response and Server-Sent Events (SSE). The preview panel SHALL display raw value, intermediate step results, final transformed value, quality indicator, connection state, last successful timestamp, and actionable error or retry copy. The UI SHALL use idle, connecting, live, reconnecting, degraded, error, and stale states.

#### Scenario: View preview from a valid server response

- **WHEN** a user opens preview and the server returns a valid snapshot
- **THEN** the UI displays raw, step results, final value, quality, and live/last-success context
- **AND** the panel does not claim live before the valid server response is accepted

#### Scenario: Live update via SSE

- **WHEN** the server publishes a valid device value event
- **THEN** the preview panel updates the raw, intermediate, final, quality, and last-success values in real time
- **AND** the panel remains visibly live

#### Scenario: SSE reconnection is visible

- **WHEN** the SSE connection is lost after a valid value
- **THEN** the UI preserves the last value with stale context, enters reconnecting or degraded state, and retries with bounded backoff
- **AND** the UI returns to live only after a valid server event

#### Scenario: Server preview failure is fail-closed

- **WHEN** the server preview returns a typed failure before a valid value
- **THEN** the panel renders an actionable localized error and retry action
- **AND** it does not use local or mock values to present a normal or live preview


<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->

---
### Requirement: Write precision settings

The UI SHALL allow operators to configure write timestamp precision (seconds or milliseconds).

#### Scenario: Select millisecond precision

- WHEN a user selects millisecond precision
- THEN the UI saves the setting for subsequent writes

#### Scenario: Configure partition interval

- WHEN a user selects a partition interval (daily/weekly/monthly)
- THEN the UI saves the setting for time-series storage

---
### Requirement: Query-based state management

The UI SHALL use TanStack Query for server state management with caching, automatic refetching, and optimistic updates.

#### Scenario: Cached data loading

- WHEN a user navigates to a previously visited page
- THEN the UI displays cached data immediately while revalidating in background

#### Scenario: Optimistic status toggle

- WHEN a user toggles device status
- THEN the UI reflects the change immediately
- AND reverts if the backend operation fails

---
### Requirement: Device onboarding wizard

The system SHALL provide a guided onboarding wizard component to help users complete device configuration.

#### Scenario: Wizard opens after device creation

- **WHEN** a new device is created successfully
- **THEN** the onboarding wizard automatically opens

#### Scenario: Wizard guides through steps

- **WHEN** a user follows the onboarding wizard
- **THEN** the wizard guides through 7 steps: device creation, connection test, activation, point creation, polling group assignment, tag creation, and mapping creation

#### Scenario: Wizard shows progress

- **WHEN** a user progresses through the wizard
- **THEN** each step shows completion status and validation results

#### Scenario: Wizard validates each step

- **WHEN** a user attempts to proceed to the next step
- **THEN** the wizard validates the current step before allowing progression

---
### Requirement: Device status dashboard

The system SHALL provide a dashboard page showing device collection status and configuration completeness.

#### Scenario: Display device statistics

- **WHEN** a user views the status dashboard
- **THEN** the dashboard displays summary statistics: total devices, active devices, devices collecting data

#### Scenario: Display device list with status

- **WHEN** a user views the status dashboard
- **THEN** the dashboard displays a list of devices with their collection status, last collection time, and error count

#### Scenario: Real-time status updates

- **WHEN** device collection status changes
- **THEN** the dashboard updates automatically using SSE

#### Scenario: Filter devices by status

- **WHEN** a user filters devices by status
- **THEN** the dashboard shows only devices matching the filter criteria

---
### Requirement: Device readiness indicator

The system SHALL display readiness status indicators in the device list and device detail pages.

#### Scenario: Show readiness badge

- **WHEN** a user views the device list
- **THEN** each device displays a readiness badge (ready/not ready)

#### Scenario: Show readiness details

- **WHEN** a user clicks on a device readiness badge
- **THEN** the system shows detailed readiness check results and missing configurations

#### Scenario: Show completion progress

- **WHEN** a user views a device detail page
- **THEN** the page displays a configuration completion progress bar

---
### Requirement: Sidebar navigation improvements

The system SHALL migrate legacy sidebar-aligned feature routes into dashboard modal workflows, except the independent `/test` page.

#### Scenario: Sidebar feature routes open modal workflows
- **WHEN** a user accesses legacy feature paths (`/datalink/devices`, `/datalink/settings`, `/datalink/points`, `/datalink/mappings`, `/datalink/wizard`)
- **THEN** the system redirects to `/datalink`
- **AND** opens the corresponding dashboard modal context

#### Scenario: Test page remains independent
- **WHEN** a user accesses `/test`
- **THEN** the system keeps `/test` as an independent page
- **AND** does not convert it into dashboard modal flow

---
### Requirement: Flow-first workspace visualization

The UI SHALL provide a flow-first workspace with a desktop shell that combines `StepRail`, `ContextBar`, `PrimaryWorkArea`, `InspectorPanel`, and `BottomSummaryBar`.

The workspace SHALL preserve the following flow-reading order:
- Device context
- Source planning and value visualization
- Tag linkage
- Output readiness and target mapping

#### Scenario: Persistent desktop workbench shell
- **WHEN** the operator changes selected devices, source rules, tag bindings, or output targets
- **THEN** the workbench updates the relevant shell regions cohesively
- **AND** preserves a stable reading order across the full flow

#### Scenario: Source-to-output linkage visibility
- **WHEN** a source span becomes linked to a tag and prepared for output
- **THEN** the UI shows that linkage in source, tag, and output contexts
- **AND** displays readiness as `draft`, `ready`, `partial`, `blocked`, or `applied`

---
### Requirement: Accessible and responsive operator workspace

The UI SHALL remain fully operable by keyboard, preserve i18n compatibility, avoid horizontal overflow at supported desktop breakpoints, and expose truthful accessible state for controls that change server-backed configuration.

#### Scenario: Keyboard-only operation in v2 workbench

- **WHEN** an operator uses keyboard-only navigation in /studio/v2
- **THEN** the operator can switch steps, operate the current toolbar, review selections, activate Share/global controls, and trigger inspector actions
- **AND** all interactive controls expose visible focus and screen-reader-readable state changes

#### Scenario: Share state is announced accessibly

- **WHEN** an operator focuses the Share/global toggle
- **THEN** the control exposes a localized accessible name and aria-pressed state
- **AND** Enter and Space activation works while pending state disables repeat activation

#### Scenario: Desktop responsive stability

- **WHEN** the workbench is rendered at 1280px or above
- **THEN** the desktop shell shows all core workflow regions without horizontal scrolling
- **AND** the main working surface remains usable at 1920x1080 without overflow traps or collapsed critical controls


<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->

---
### Requirement: Source template library

The UI SHALL allow operators to save, load, update, and delete source planning templates for Step 2 rule groups.

Each template SHALL store at least:
- source rules
- rule ordering and lock state
- preferred Step 2 view mode
- naming defaults when applicable

Template persistence for this workflow SHALL be browser-local in this redesign round.

#### Scenario: Save and reuse source template
- **WHEN** an operator saves the current Step 2 rule set as a template
- **THEN** the template is persisted locally
- **AND** the operator can later load it to prefill planning controls

#### Scenario: Warn on capability mismatch
- **WHEN** an operator applies a saved template to a device with different capability assumptions
- **THEN** the UI warns about mismatches such as `address base` or `word order`
- **AND** requires the operator to confirm before applying the template

---
### Requirement: Motion-guided operator flow

The UI SHALL use motion cues to guide stage transitions between source planning, grid allocation, tag linkage, and DB commit.

Animations MUST use short transitions (150-300ms) and MUST support reduced-motion preference.

#### Scenario: Guided transition after source planning
- **WHEN** an operator confirms source plan
- **THEN** the grid allocation region receives a transition cue indicating next action

#### Scenario: Reduced-motion mode
- **WHEN** user preference is `prefers-reduced-motion`
- **THEN** motion cues are replaced with static visual emphasis without animation

---
### Requirement: Planning intelligence and safety checks

The UI SHALL provide planning assistance and safety checks for one-screen operator execution.

#### Scenario: Naming preview and duplicate detection
- **WHEN** an operator defines batch naming rules
- **THEN** the UI previews generated names
- **AND** flags duplicates or naming conflicts before commit

#### Scenario: Two-stage validation execution
- **WHEN** an operator validates pending changes
- **THEN** the UI runs structural validation first
- **AND** runs executable validation only after structural validation succeeds

---
### Requirement: Source desk switching preserves guided-flow continuity
The UI SHALL preserve guided-flow continuity when the operator switches between Step 2 Source desk modes.

Step 2 desk switching SHALL keep the same active source-rule context, grouped Tag handoff cues, and blocker diagnostics so the operator can continue toward Tag review and Output review without reconstructing planning state.

#### Scenario: Desk switching keeps Source-to-Tag continuity visible
- **WHEN** an operator changes Step 2 from `Inspect` to `Build` or `Triage`
- **THEN** the UI keeps the active rule summary and grouped Tag handoff cues visible for the same rule
- **AND** the operator can still understand what Step 3 will review next without reopening or reloading the workspace

#### Scenario: Desk switching does not behave like a separate workflow branch
- **WHEN** an operator uses the Step 2 desk selector repeatedly during one planning session
- **THEN** the workflow remains a single `Device -> SourceRule -> Tag review -> Output review/apply` path inside `/studio`
- **AND** changing desks SHALL NOT clear progress, hide the primary next-step call to action, or create a second primary workflow inside Step 2

<!-- @trace
source: rework-source-step-desk-modes
updated: 2026-05-09
code:
  - tests/shell/start-backend-before-frontend.sh
  - frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx
  - frontend/src/App.tsx
  - .github/skills/spectra-ask/SKILL.md
  - frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - .github/prompts/spectra-debug.prompt.md
  - .github/prompts/spectra-audit.prompt.md
  - .github/prompts/spectra-propose.prompt.md
  - .github/prompts/opsx-archive.prompt.md
  - .github/skills/spectra-apply/SKILL.md
  - .github/skills/spectra-audit/SKILL.md
  - .github/skills/spectra-commit/SKILL.md
  - .github/skills/spectra-discuss/SKILL.md
  - .spectra.yaml
  - .github/prompts/opsx-propose.prompt.md
  - CLAUDE.md
  - frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx
  - .github/skills/spectra-debug/SKILL.md
  - tests/shell/start-backend-cleanup-order.sh
  - start.sh
  - frontend/src/i18n/locales/zh-TW/common.json
  - frontend/src/pages/datalink/workbench/sourceStepRuleSummaryModel.ts
  - .github/skills/spectra-ingest/SKILL.md
  - .github/skills/spectra-archive/SKILL.md
  - frontend/src/i18n/locales/en/common.json
  - tests/shell/start-backend-logfile.sh
  - frontend/src/main.tsx
  - tests/shell/start-port-management.sh
  - .github/skills/spectra-propose/SKILL.md
  - .github/prompts/spectra-discuss.prompt.md
  - .github/prompts/spectra-ingest.prompt.md
  - .github/prompts/spectra-ask.prompt.md
  - frontend/src/components/DevAgentation.tsx
  - frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx
  - .github/prompts/spectra-commit.prompt.md
  - frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx
  - frontend/package.json
  - .github/prompts/spectra-archive.prompt.md
  - start.ps1
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - .github/prompts/opsx-apply.prompt.md
  - .github/prompts/spectra-apply.prompt.md
  - frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx
tests:
  - cmd/test_ui/static/assets/index-CcV2SQjv.css
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/components/DevAgentation.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
  - cmd/test_ui/static/assets/index-Bl1MOChG.js
  - frontend/tests/unit/pages/datalink/workbench-source-preview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - frontend/tests/unit/utils/appAgentationRemoval.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourcePreview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Commit progress reports actual result states

The v2 commit progress surface SHALL render pending, success, failed, and skipped states from the backend result model. A failed item SHALL display a safe localized error code/message and an actionable retry or recovery action. The progress surface SHALL NOT replace a failed result with a fixed HTTP 200 or success tone.

#### Scenario: Failed commit item is actionable

- **WHEN** one commit item returns failed with a typed error code
- **THEN** the item displays failed status, localized safe copy, and the retry or recovery action
- **AND** the item does not display 200 as a success result


<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->

---
### Requirement: Preview error diagnostics are separated

The live preview UI SHALL display safe localized operator messages and opaque request identifiers only. Raw backend exception text, stack traces, connection strings, and internal diagnostics SHALL remain outside the normal UI and SHALL be available only through the controlled diagnostics channel.

#### Scenario: Raw preview detail is not rendered

- **WHEN** a preview request fails with a raw backend detail and a typed code
- **THEN** the UI renders the typed code mapping and safe action
- **AND** the raw backend detail is absent from rendered text and DOM content

<!-- @trace
source: retire-legacy-studio-and-polish-v2
updated: 2026-08-26
code:
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/sourcerule/candidate_api.go
  - internal/api/router_modbus_share.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/config.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/modbusshare/service_persistence.go
  - internal/datalink/modbusshare/types.go
  - frontend/src/services/sourceRuleCandidates.ts
  - internal/datalink/sourcerule/validation.go
  - internal/datalink/dbtarget/service.go
  - internal/datalink/workspace/service_readiness_device.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/virtual/server/modbus/server.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/collector/scheduler_lifecycle.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/device/service_readiness_safe.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/target_delivery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - frontend/src/services/datalinkClient.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/settings/sql_repo.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - docs/swagger/swagger.json
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/runtime/delivery_diagnostic.go
  - internal/api/handlers/debug.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - internal/api/handlers/connection.go
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/tag/service_crud.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/collector/scheduler.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/dashboard_handler.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/datalink.ts
  - internal/datalink/modbusshare/errors.go
  - internal/api/handlers/test.go
  - internal/api/handlers/typed_errors.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - internal/datalink/modbusshare/service.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/services/modbusShare.ts
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/b10_exe_acceptance.py
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/modbusShare.test-d.ts
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/point/service_point_crud.go
  - internal/datalink/runtime/status.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - internal/datalink/sourcerule/output_apply_service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/types/datalink.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - internal/api/handlers/runtime_handler.go
  - docs/swagger/docs.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/candidate_snapshot.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/api/handlers/source_rule_handler.go
  - scripts/lib/b10_projection_assertions.py
  - internal/api/handlers/modbus_share_handler.go
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/sourcerule/share_ownership.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/modbusshare/reconciler_transaction.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - go.mod
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/typedErrors.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/api/handlers/response_keys.go
  - internal/datalink/modbusshare/service_write.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/types/sourceRuleCandidates.ts
  - internal/api/handlers/settings_handler.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/api/handlers/datalink_sse_handler.go
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/settings/service.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/modbusShare.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - internal/api/handlers/source_rule_handler_candidates.go
  - internal/datalink/settings/errors.go
  - internal/datalink/sourcerule/candidate_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/source_rule_handler_output_apply.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - scripts/check_file_lines.sh
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/usePreviewStream.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/hooks/previewStreamEvents.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - internal/api/router.go
  - internal/api/handlers/transport_wrapper.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/api/handlers/template.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/modbusshare/service_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - scripts/lib/b10_focused_suite.py
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - scripts/lib/b10_negative_steps.py
  - internal/virtual/server/modbus/server_register_handlers.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
tests:
  - frontend/tests/unit/utils/safeJson.test.ts
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - cmd/test_ui/main.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/share_startup_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/source_rule_handler_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - cmd/test_ui/server_runtime.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - internal/virtual/server/modbus/server_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - frontend/tests/unit/services/modbusShare.test.ts
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - internal/datalink/collector/scheduler_refactor_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/api/router_modbus_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - internal/datalink/modbusshare/status_contract_test.go
  - internal/api/handlers/test_client_operations.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/api/handlers/test_request_helpers_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/test_connection_handler.go
  - internal/api/handlers/test_script_handler.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - internal/datalink/sourcerule/share_gates_test.go
  - cmd/test_ui/harness_config_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/api/handlers/test_client_factory.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/api/handlers/test_monitor_handler.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/modbus_share_swagger_contract_test.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - scripts/lib/test_check_file_lines.py
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
-->