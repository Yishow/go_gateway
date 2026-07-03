# source-rule-runtime Specification

## Purpose
Define the canonical runtime contract for persisted source rules so rule state survives restarts, activation is gated by probe success, derived relationships remain intact across enable/disable transitions, and live values surface back into the source-planning grid.

## Requirements

### Requirement: Persisted source-rule lifecycle
The system SHALL persist each source rule as a first-class record with device association, planning definition, revision metadata, and enabled/disabled state.

#### Scenario: Save rule definition
- **WHEN** an operator creates or updates a source rule
- **THEN** the system stores the rule in the database with its device association, planning inputs, current enabled/disabled state, and a new revision id
- **AND** the saved rule can be reloaded independently of the current browser session

#### Scenario: Restore rule state after restart
- **WHEN** the service restarts
- **THEN** the system reloads persisted source rules from the database together with the latest persisted revision metadata
- **AND** restores each rule's enabled/disabled state before resuming runtime behavior

---
### Requirement: Rule revisions produce downstream candidate snapshots
The system SHALL persist a dedicated downstream candidate snapshot for each saved source-rule revision before any later review/apply workflow can consume that revision.

Source-rule runtime SHALL own candidate generation and persistence for downstream review/apply workflows.

Each persisted candidate snapshot MUST include at least:

- `source_rule_id`
- `revision_id`
- `candidate_type`
- `payload`
- `status`
- `generated_at`

#### Scenario: Save rule writes candidate snapshot
- **WHEN** a source rule is created or updated
- **THEN** the system persists the new revision and synchronizes runtime rule-owned point state
- **AND** persists the candidate snapshot for that revision before later review flows can query it

#### Scenario: Candidate snapshot survives restart
- **WHEN** the service restarts after a source-rule revision has generated candidate snapshots
- **THEN** the system restores the persisted candidate snapshots for that revision
- **AND** downstream review/apply workflows can query them without regenerating an incompatible contract

---
### Requirement: Rule-owned identity and signature are canonical across revisions
The system SHALL compare rule-derived tags, mappings, and output candidates using one canonical rule-owned identity and one deterministic proposed-signature model.

A rule-owned identity MUST include at least:

- `source_rule_id`
- `candidate_type`
- `derived_from_rule_address`
- target binding scope fields needed to distinguish one downstream object from another

The proposed signature MUST be derived from the effective candidate payload for that identity.

#### Scenario: Same identity with changed payload is detectable
- **WHEN** a later rule revision preserves a downstream object's rule-owned identity but changes its effective payload
- **THEN** the system preserves the same rule-owned identity for comparison
- **AND** computes a different proposed signature for the changed candidate

---
### Requirement: Snapshot completeness records ready, blocked, and deferred target state explicitly
The system SHALL treat a candidate snapshot as complete only when it contains the current tag candidate set and an explicit state entry for each supported downstream target.

Each downstream target entry MUST be present as one of `ready`, `blocked`, or `deferred`, and blocked or deferred entries MUST include the reason.

Snapshot completeness SHALL NOT require every downstream target to be apply-ready.

#### Scenario: Missing database context still yields a complete snapshot
- **WHEN** a source-rule revision generates tag candidates but the database target lacks connector context
- **THEN** the snapshot is still marked complete
- **AND** the database target entry is persisted as `blocked` or `deferred` with its reason instead of being omitted

---
### Requirement: Rule changes preserve applied downstream state until explicit reapply
The system SHALL mark changed downstream assets as `out_of_sync` rather than silently replacing applied tag or output state.

#### Scenario: Applied downstream asset becomes out of sync
- **WHEN** a saved source-rule revision changes the derived downstream result
- **THEN** the system preserves the previously applied asset
- **AND** marks it as `out_of_sync` until the operator explicitly reapplies the new candidate

---
### Requirement: Rule revision rollback restores the last applied downstream state
The system SHALL treat rollback as a state-based restore of the last successfully applied downstream state per target.

#### Scenario: Unapplied revision is rolled back
- **WHEN** a newer source-rule revision is abandoned or rolled back before downstream apply is completed
- **THEN** the system removes unapplied candidate snapshots for the rolled-back revision
- **AND** preserves the last successfully applied downstream state for tags, database output, and Local Modbus output independently

---
### Requirement: Rule enablement controls collection without deleting derived relationships
The system SHALL treat rule enablement as collection control only and SHALL preserve Point, Tag, Mapping, and Output relationships when a rule is disabled.

#### Scenario: Disable rule preserves derived records
- **WHEN** an operator disables a rule
- **THEN** the system stops collection for that rule
- **AND** preserves the derived Point, Tag, Mapping, and Output relationships in storage

#### Scenario: Re-enable rule resumes existing definition
- **WHEN** an operator re-enables a previously disabled rule
- **THEN** the system resumes collection using the stored rule definition and derived relationships
- **AND** does not require the operator to recreate the rule or downstream bindings

---
### Requirement: Rule activation is gated by protocol probe success
The system SHALL use the device readiness contract to distinguish planning from runtime activation and collection.

#### Scenario: Connect succeeds but probe fails
- **WHEN** a device test reports transport connect success and protocol probe failure
- **THEN** the operator MAY save the device and save source-rule planning state
- **AND** the system SHALL reject rule activation and data collection until probe succeeds

---
### Requirement: Rule-driven live values appear in the source grid
The system SHALL display parsed live values for enabled rules in the Step 2 source grid.

#### Scenario: Enabled rule updates grid
- **WHEN** runtime collection succeeds for an enabled source rule
- **THEN** the Step 2 grid shows the parsed value, timestamp, and error state for the rule-derived span
- **AND** the displayed value follows the rule's configured data-type and merge semantics

---
### Requirement: Persisted rule MAY store intended target data type
The system SHALL persist an optional **target data type** field on each source rule record. The field SHALL be reloadable after service restart and SHALL participate in rule enable/disable lifecycle without being silently discarded.

#### Scenario: Reload after restart preserves target intent
- **WHEN** the service restarts after a rule with a non-null target data type was saved
- **THEN** the system SHALL reload the target data type with the rest of the rule definition
- **AND** SHALL apply the same synchronization behavior as before restart

---
### Requirement: Rule synchronization derives tag type and default mapping pipeline from target intent
When synchronizing a rule to Points, Tags, and Mappings, the system SHALL set **Tag.data_type** to the rule’s target data type (or protocol read type when target is unset). When target differs from the derived Point read type, the system SHALL ensure the mapping includes a validated default **`cast`** transform pipeline bridging point semantics to the tag type, without requiring manual pipeline entry for the cast alone.

#### Scenario: New tag uses target type when provided
- **WHEN** a rule declares a target data type and synchronization creates a new tag for a derived point
- **THEN** the created tag SHALL use the target data type
- **AND** the point SHALL retain the protocol read data type

#### Scenario: Mapping receives default cast when types differ
- **WHEN** target data type differs from the point read data type for a rule-derived mapping
- **THEN** the mapping SHALL include a `cast` step appropriate to the declared target type
- **AND** the pipeline SHALL pass mapping validation

---
### Requirement: Persisted rule MAY store optional scale parameters
The system SHALL persist optional **linear scale** fields on a source rule (multiplier and/or offset, or an equivalent single structure) when the operator configures engineering-unit conversion. When unset, the system SHALL treat scaling as identity and SHALL NOT insert a redundant `scale` step.

#### Scenario: Scale fields persist across restarts
- **WHEN** a rule saved with non-default scale parameters is reloaded after restart
- **THEN** the system SHALL restore those parameters
- **AND** SHALL regenerate or validate the derived mapping `scale` step consistently with the saved values

---
### Requirement: Downstream integrity gaps surface as readiness issues before activation

The system SHALL surface missing or inconsistent downstream rule-derived relationships as readiness issues before activation.

#### Scenario: Missing downstream relationships become readiness blockers

- **WHEN** a persisted source rule is missing a required derived point, tag, mapping, or required database target relationship for the chosen activation scope
- **THEN** the workspace readiness result records a normalized issue for that gap before activation starts
- **AND** the operator SHALL NOT first discover that gap from a late runtime or write-path failure

##### Example: missing tag and missing database target both surface early

- **GIVEN** rule rule-A still exists but point pt-A has no persisted tag and pt-B has no required database target
- **WHEN** readiness evaluates the workspace before activation
- **THEN** readiness returns normalized issues for rule-A or its affected points instead of waiting for runtime to fail

<!-- @trace
source: enforce-workspace-readiness-contract
updated: 2026-06-09
code:
  - frontend/src/types/studioV2Workspace.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - frontend/src/App.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/dbtarget/live_projection.go
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/types/runtimeTruth.ts
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - tests/shell/start-frontend-readiness.sh
  - internal/api/handlers/runtime_stream_handler.go
  - internal/datalink/dbtarget/tooling_service.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/runtime/service_device_sync.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/runtime/status.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - start.sh
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/datalink/schema/schema_dbtarget_models.go
  - internal/api/handlers/studio_v2_workspace_audit.go
  - frontend/src/types/datalink.ts
  - frontend/studio-v2-diagnostic.cjs
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - internal/datalink/workspace/service_readiness.go
  - tests/shell/start-frontend-deps.sh
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/migrator.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - internal/datalink/audit/types.go
  - internal/datalink/workspace/service.go
  - internal/datalink/runtime/service.go
  - internal/datalink/runtime/truth_state.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - frontend/src/types/databaseDelivery.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/datalink/audit/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/types/runtimeDiagnostics.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/datalink/runtime/service_projection_state.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/runtime/service_workspace_projection.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/studio-v2-real-check.cjs
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - internal/datalink/migrator_database_delivery_outcomes.go
  - internal/api/handlers/source_rule_handler.go
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - frontend/src/types/studioV2RuntimeContext.ts
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - internal/datalink/workspace/service_runtime_projection.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - internal/datalink/audit/sql_repository.go
  - tests/shell/start-frontend-install-failure.sh
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - Makefile
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/audit/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
tests:
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/workspace/service_activation_projection_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/runtime/ingestor_test.go
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/audit/service_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/api/router_runtime_database_delivery_test.go
  - internal/datalink/runtime/service_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/workspace/service_readiness_empty_test.go
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
-->

---
### Requirement: Rule lifecycle changes trigger runtime reconciliation results

The system SHALL turn source-rule lifecycle changes into explicit runtime reconciliation results.

#### Scenario: Rule change reports reconcile result

- **WHEN** a persisted source rule is enabled, disabled, updated, or removed while the workspace is running
- **THEN** the system records and returns the runtime reconciliation result for that rule scope
- **AND** the operator SHALL NOT have to infer from stale runtime behavior whether the rule change actually took effect

##### Example: disabling one running rule returns immediate reconcile

- **GIVEN** rule-A is currently active in runtime for dev-A
- **WHEN** the operator persists a disable action for rule-A
- **THEN** the system returns the reconcile result for rule-A instead of leaving the operator to guess whether collection stopped

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
### Requirement: Rule-derived database target drift does not block unrelated live runtime

The system SHALL keep stale rule-derived database target drift from blocking unrelated live runtime and database delivery scope.

#### Scenario: Orphaned database target is isolated from live runtime scope

- **WHEN** a rule-derived database target becomes orphaned from its live rule-owned relationship
- **THEN** the system isolates that orphaned target from live schema ensure and write delivery
- **AND** unrelated live rule-owned targets continue to activate and deliver normally

##### Example: removed rule leaves one orphaned target behind

- **GIVEN** rule-A was removed, target row-A is now orphaned, and rule-B still owns live target row-B
- **WHEN** activation and delivery evaluate current runtime scope
- **THEN** row-A is isolated and row-B continues to activate and deliver normally

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