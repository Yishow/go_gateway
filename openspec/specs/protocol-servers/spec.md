# protocol-servers Specification

## Purpose
TBD - created by archiving change p3-enable-virtual-device-mode. Update Purpose after archive.

## Requirements

### Requirement: Virtual server hosting

The system SHALL support hosting one or more protocol servers (listeners) to accept incoming connections from external clients (SCADA/HMI).

#### Scenario: Host Modbus TCP server

- GIVEN the gateway is configured to start a Modbus TCP server on port 10502
- WHEN the gateway starts
- THEN it listens on port 10502 and accepts Modbus queries

---
### Requirement: Virtual memory mapping

The system SHALL provide a strictly addressed virtual memory space that serves as the data source for hosted servers.

#### Scenario: Read from virtual memory

- GIVEN a value 12345 is stored at virtual address D100
- WHEN an external Modbus client reads Holding Register 100
- THEN the server responds with value 12345

---
### Requirement: Data bridging

The system SHALL support mapping collected data from real devices into the virtual memory space.

#### Scenario: Bridge data to virtual server

- GIVEN a real device tag `PLC1.D200` is collected
- AND a mapping exists from `PLC1.D200` to Virtual Address `V.D100`
- WHEN `PLC1.D200` changes value
- THEN `V.D100` is automatically updated

---
### Requirement: Local Modbus server sink on port 5020

The system SHALL support a local Modbus TCP server sink that mirrors committed tag values into virtual Modbus memory and exposes them to external clients using the persisted Share bind address, port, slave id, and register capacity. Port 5020 SHALL be treated as a persisted or explicitly selected value, not as an unconditional handler fallback.

#### Scenario: Start local Modbus sink from durable settings

- **GIVEN** persisted Share settings have enabled=true, bind address 127.0.0.1, port 15020, valid slave id, and register capacity
- **WHEN** the gateway starts or the operator enables local Modbus sharing
- **THEN** the gateway starts a Modbus TCP server bound to 127.0.0.1:15020
- **AND** the startup status exposes the configured address, port, slave id, capacity, and running state

#### Scenario: Disabled Share has no listener

- **GIVEN** persisted settings.modbus_share.enabled is false
- **WHEN** the gateway starts, hydrates workspace state, or receives a rule-level share_enabled=true
- **THEN** no local Modbus TCP listener is bound
- **AND** the rule-level flag does not create a listener or mutate virtual memory
- **AND** the status reports disabled rather than ready

#### Scenario: Mirror tag value across its complete register span

- **GIVEN** a committed workspace-owned tag mapping has a validated register start, datatype span, and stride
- **WHEN** the mapped tag value is updated in pipeline runtime
- **THEN** the corresponding virtual Modbus registers are updated using two bytes per register and the declared encoding
- **AND** external Modbus clients reading that span receive the mirrored value

#### Scenario: Port or address bind conflict

- **GIVEN** the configured bind address or port is already occupied or unavailable
- **WHEN** local Modbus sharing starts
- **THEN** no partially started listener remains
- **AND** the system reports a typed bind failure with the configured endpoint and an operator action
- **AND** local Modbus sharing remains failed or disabled until the durable setting is corrected and retried

#### Scenario: Device scope isolation

- **GIVEN** pipeline data is processed for multiple devices and workspaces
- **WHEN** each workspace reconciles its persisted Share desired set
- **THEN** each source device mapping remains isolated by workspace and source-rule ownership
- **AND** no cross-device or cross-workspace deletion occurs outside the authorized scope

#### Scenario: Local sink target conflict detection

- **GIVEN** multiple mappings target overlapping register ranges in the same local Share capacity
- **WHEN** the complete desired set is validated
- **THEN** the system flags the range collision before mutation
- **AND** prevents enabling or activating the conflicting set
- **AND** retains the last valid projection until the collision is resolved

<!-- @trace
source: productionize-modbus-share-lifecycle
updated: 2026-08-26
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/modbus_share_handler_mapping_query.go
  - internal/datalink/sourcerule/local_modbus_mapping_reader.go
  - internal/datalink/collector/scheduler_lifecycle.go
  - internal/api/handlers/modbus_share_handler.go
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3RuntimeStreams.ts
  - internal/datalink/sourcerule/service_state.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/typed_errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - internal/datalink/point/service_point_crud.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - internal/datalink/runtime/target_delivery.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/workspace/service_readiness_device.go
  - internal/datalink/modbusshare/sql_revision_store.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/sourcerule/candidate_api.go
  - internal/datalink/sourcerule/local_modbus_restore.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - scripts/b10_exe_acceptance.py
  - internal/datalink/workspace/service.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/hooks/datalink/useModbusShareCandidateReview.ts
  - internal/datalink/modbusshare/geometry.go
  - frontend/src/components/datalink/wizard/LivePreviewPanel.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/device/service_readiness_safe.go
  - internal/datalink/sourcerule/repository_memory.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - docs/swagger/swagger.json
  - go.mod
  - frontend/src/types/modbusShare.test-d.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/modbusshare/sql_revision_store_helpers.go
  - internal/datalink/runtime/modbus_share_delivery.go
  - frontend/src/features/datalink/workbench-v2/settings/useSettingsOperations.ts
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/modbusshare/errors.go
  - internal/datalink/schema/migrations/001_initial_schema_sqlite.sql
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/modbusshare/service_projection.go
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/tag_apply_service.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_conflicts.go
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - internal/datalink/sourcerule/tag_review_decision_service.go
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/runtime/ingestor.go
  - internal/virtual/server/modbus/server_register_handlers.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/utils/safeJson.ts
  - internal/api/handlers/source_rule_handler_tag_review_decisions.go
  - internal/datalink/modbusshare/canonical_plan.go
  - internal/api/handlers/mapping_handler.go
  - frontend/src/hooks/datalink/useModbusShareStatus.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/types/studioV2Workspace.ts
  - internal/datalink/settings/errors.go
  - frontend/src/components/datalink/wizard/MappingWizard.tsx
  - internal/api/handlers/modbus_share_handler_reconcile.go
  - scripts/lib/b10_projection_assertions.py
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/sourcerule/interfaces.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/services/datalinkClient.ts
  - internal/datalink/workspace/service_readiness.go
  - frontend/src/hooks/previewStreamEvents.ts
  - internal/api/handlers/modbus_share_handler_mapping_mutations.go
  - frontend/src/services/sourceRuleCandidates.ts
  - frontend/src/pages/datalink/workbench-v2/studioV2AutosaveBarrier.ts
  - internal/api/handlers/modbus_share_handler_gates.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/status.go
  - internal/datalink/sourcerule/validation.go
  - docs/swagger/docs.go
  - internal/api/handlers/response_keys.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/datalink/sourcerule/share_desired_mappings.go
  - frontend/vite.config.ts
  - internal/api/handlers/connection.go
  - frontend/src/services/sourceRuleTagReviewDecisions.ts
  - frontend/src/services/modbusShare.test-d.ts
  - internal/datalink/modbusshare/reconciler.go
  - internal/datalink/modbusshare/settings_lifecycle.go
  - internal/api/handlers/source_rule_handler_candidates.go
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - internal/api/handlers/dashboard_handler.go
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/test.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/services/datalink.ts
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4Activation.ts
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - internal/datalink/device/service_status.go
  - internal/api/handlers/template.go
  - internal/datalink/sourcerule/candidate_snapshot_revision.go
  - internal/datalink/modbusshare/service_network_helpers.go
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/api/router_modbus_share.go
  - internal/datalink/modbusshare/service_persistence.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/components/datalink/wizard/steps/PreviewStep.tsx
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - internal/datalink/modbusshare/reconciler_state.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/datalink/modbusshare/reconciler_transaction.go
  - internal/datalink/sourcerule/share_restore_projection.go
  - internal/virtual/server/modbus/server_connection.go
  - frontend/src/hooks/usePreviewStream.ts
  - frontend/src/features/datalink/workbench-v2/settings/settingsOperationOwnership.ts
  - scripts/lib/b10_acceptance_helpers.py
  - internal/api/handlers/source_rule_handler_output_apply.go
  - frontend/src/hooks/datalink/useRuntimeStream.ts
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
  - frontend/src/features/datalink/workbench-v2/settings/SettingsStatus.tsx
  - frontend/src/types/runtimeDiagnostics.ts
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - internal/api/handlers/settings_handler.go
  - internal/api/handlers/modbus_share_handler_status.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - internal/datalink/migrator.go
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/settings/sql_repo.go
  - internal/virtual/server/modbus/server.go
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/datalink/settings/service.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.ts
  - internal/api/handlers/datalink_sse_handler.go
  - frontend/src/types/sourceRuleCandidates.ts
  - frontend/src/services/modbusShare.ts
  - internal/datalink/sourcerule/candidate_snapshot_local_modbus_outputs.go
  - internal/datalink/sourcerule/mutation_rollback_created.go
  - internal/datalink/modbusshare/service.go
  - internal/datalink/tag/service_crud.go
  - internal/api/handlers/transport_wrapper.go
  - internal/api/router.go
  - internal/datalink/sourcerule/output_apply_service.go
  - internal/datalink/sourcerule/share_ownership.go
  - internal/datalink/modbusshare/types.go
  - scripts/lib/b10_focused_suite.py
  - internal/api/handlers/source_rule_handler_tag_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/sourcerule/candidate_scope.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - internal/datalink/sourcerule/mutation_rollback.go
  - internal/datalink/modbusshare/service_lifecycle.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/modbusshare/reconciler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/step3LiveSubscription.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/dbtarget/service_probe.go
  - internal/api/handlers/source_rule_handler.go
  - internal/api/handlers/debug.go
  - docs/swagger/swagger.yaml
  - internal/datalink/collector/scheduler.go
  - frontend/src/utils/typedErrors.ts
  - docs/releases/retire-legacy-studio-and-polish-v2.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - scripts/lib/b10_negative_steps.py
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/types/modbusShare.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsSections.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/api/handlers/config.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/collector/scheduler_dispatch.go
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - internal/datalink/sourcerule/candidate_snapshot.go
  - scripts/check_file_lines.sh
  - internal/api/handlers/device_health_handler.go
  - internal/datalink/modbusshare/service_write.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/runtime/delivery_diagnostic.go
tests:
  - internal/datalink/modbusshare/settings_hydration_failure_test.go
  - internal/api/handlers/test_monitor_handler.go
  - internal/datalink/modbusshare/reconciler_nil_store_test.go
  - frontend/tests/unit/workbench-v2/settings-operation-ownership.test.ts
  - internal/api/handlers/modbus_share_handler_diagnostics_test.go
  - internal/datalink/modbusshare/reconciler_concurrency_test.go
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - internal/datalink/modbusshare/canonical_plan_test.go
  - internal/api/handlers/modbus_share_handler_production_gate_test.go
  - internal/api/handlers/typed_errors_contract_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/datalink/sourcerule/share_desired_mappings_test.go
  - internal/api/handlers/modbus_share_handler_gates_test.go
  - internal/datalink/migrator_test.go
  - internal/datalink/sourcerule/share_gates_test.go
  - internal/datalink/modbusshare/reconciler_lifecycle_test.go
  - internal/datalink/sourcerule/candidate_scope_test.go
  - internal/datalink/modbusshare/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - internal/api/handlers/modbus_share_handler_status_contract_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/utils/typedErrors.test.ts
  - internal/datalink/modbusshare/reconciler_idempotency_test.go
  - internal/api/handlers/source_rule_handler_test.go
  - cmd/test_ui/target_writer_test.go
  - internal/api/handlers/test_client_operations.go
  - internal/api/handlers/mapping_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - cmd/test_ui/harness_config_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/studioV2WorkspaceActivation.test.ts
  - internal/datalink/modbusshare/reconciler_transaction_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/datalink/settings/sql_repo_storage_errors_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview-changes.test.tsx
  - internal/api/handlers/test_request_helpers_test.go
  - frontend/tests/unit/services/modbusShare.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - internal/datalink/settings/sql_repo_cas_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/modbusshare/geometry_validation_test.go
  - frontend/tests/unit/hooks/usePreviewStream.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/modbusshare/stale_span_invalidation_test.go
  - frontend/tests/unit/workbench-v2/step4-share-helpers.ts
  - cmd/test_ui/server_runtime.go
  - internal/api/handlers/source_rule_handler_database_candidates_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/utils/safeJson.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/test_script_handler.go
  - internal/datalink/modbusshare/settings_lifecycle_cas_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/sourcerule/share_restore_projection_test.go
  - frontend/tests/unit/workbench-v2/settings-operations-race.test.tsx
  - internal/api/handlers/source_rule_handler_local_modbus_candidates_test.go
  - internal/datalink/collector/scheduler_refactor_test.go
  - internal/api/handlers/studio_v2_workspace_activation_contract_test.go
  - internal/datalink/runtime/ingestor_target_outcomes_test.go
  - internal/datalink/modbusshare/service_write_hydration_test.go
  - cmd/test_ui/service_wiring.go
  - frontend/tests/unit/runtime-dashboard/runtime-stream-contract.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler_barrier_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - scripts/lib/test_check_file_lines.py
  - internal/api/handlers/test_connection_handler.go
  - frontend/tests/unit/workbench-v2/autosave-barrier.test.ts
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - cmd/test_ui/target_writer.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - internal/datalink/modbusshare/reconciler_empty_revision_b10_test.go
  - frontend/tests/unit/workbench-v2/commit-progress-accessibility.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-stream-state.test.tsx
  - internal/datalink/workspace/service_readiness_connector_issue_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/sourcerule/mutation_rollback_read_failure_test.go
  - internal/api/handlers/test_client_factory.go
  - internal/api/handlers/studio_v2_runtime_apply_connector_safety_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.fixture.tsx
  - internal/api/handlers/studio_v2_runtime_apply_test.go
  - internal/datalink/modbusshare/reconciler_test.go
  - internal/datalink/modbusshare/status_contract_test.go
  - scripts/lib/test_b10_acceptance_helpers.py
  - frontend/tests/unit/workbench-v2/step4-database-components.test.tsx
  - cmd/test_ui/harness_config.go
  - internal/api/handlers/settings_handler_test.go
  - internal/datalink/modbusshare/sql_revision_store_test.go
  - cmd/test_ui/share_startup.go
  - internal/datalink/modbusshare/reconciler_validation_test.go
  - internal/virtual/server/modbus/server_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/handlers/datalink_sse_handler_test.go
  - internal/api/handlers/modbus_share_handler_lifecycle_test.go
  - internal/api/modbus_share_swagger_contract_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/modbusshare/settings_lifecycle_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/handlers/modbus_share_handler_parse_test.go
  - internal/datalink/device/service_readiness_contract_test.go
  - internal/datalink/sourcerule/repository_memory_test.go
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/locale-parity.test.ts
  - frontend/tests/unit/workbench-v2/step4-share-summary.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/router_modbus_share_test.go
  - cmd/test_ui/share_runtime_reconcile.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/datalink/modbusshare/settings_hydration_lifecycle_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route-state.test.tsx
  - cmd/test_ui/share_startup_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - frontend/tests/unit/workbench-v2/step3-live-stream-malformed.test.tsx
-->