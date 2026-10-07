# guided-recording-workflow Specification

## Purpose

TBD - created by archiving change 'redesign-studio-recording-flow'. Update Purpose after archive.

## Requirements

### Requirement: Four-step intent-led setup
Operators MUST be able to save and navigate offline drafts with unresolved probes; draft navigation MUST NOT imply verified readiness, and only dependent reads or activation SHALL be blocked. The existing Studio V2 route SHALL guide device connection, acquisition points, Point-to-Tag mapping and grouped database output in four steps. Basic typed Tag writing MUST remain usable without required reporting, retention setup, aggregation or physical measurement semantics. Explicitly selected derived operations SHALL require their own confirmed semantics. The system MUST retain the Go backend and existing React framework and MUST NOT add a parallel V3 or restore a dedicated legacy /studio route. The basic path SHALL derive and persist routine mappings, show only necessary operator choices and default new recording to one managed group per device without converting existing advanced configurations.

#### Scenario: Mixed sensor configured without SQL
- **WHEN** an operator maps temperature, pressure, flow and counter source values as basic typed Tags
- **THEN** the operator can select those Tags for a managed write group without hand-written SQL or mandatory derived-value policies; choosing derived usage then requests the required confirmed roles.

#### Scenario: Reload and edited connection
- **WHEN** an operator reloads setup or changes a saved connection identity
- **THEN** the same persisted validity and revision checks determine progress, and the old probe does not keep the changed device ready.

#### Scenario: Unavailable protocol
- **WHEN** MQTT or another protocol lacks the complete V2 setup/parser/collector capability
- **THEN** the UI clearly marks it unavailable with a reason rather than letting the user enter an unsupported range flow.

#### Scenario: Same address on different devices
- **WHEN** two devices both expose address 40001
- **THEN** the address conflict checker keeps them separate while still detecting true overlapping ranges within the same device/address space.

#### Scenario: Missing advanced measurement
- **WHEN** a saved Point-to-Tag mapping has no persisted MeasurementDefinition
- **THEN** basic raw group setup proceeds without fake measurement IDs; only explicitly selected advanced semantics remain blocked.

#### Scenario: Offline draft navigation
- **WHEN** a device is offline or its saved configuration has no successful probe
- **THEN** the operator can save a draft and use Next/Back to configure later steps, including after reload, while the device remains visibly unverified and its dependent activation is blocked with a repair action.

#### Scenario: Routine mapping does not require duplicate entry
- **WHEN** the operator confirms known points and supported types in the basic path
- **THEN** the existing mapping services persist real scoped Tag identities once, display actual reads and allow correction without a second manual mapping exercise.

#### Scenario: Existing advanced configuration is reopened
- **WHEN** an existing custom-table, multi-entity or cross-device group is loaded
- **THEN** its persisted identity and layout remain intact with accessible advanced controls; the basic default neither splits nor overwrites it.

#### Scenario: New basic setup is resumed
- **WHEN** a new per-device managed draft is reloaded or its save response is retried
- **THEN** the persisted (workspace_id, device_id, canonical managed role) key reuses the same group and mappings without duplicate writers or selecting another group. A changed destination requires an explicit CAS-checked draft and new preparation/apply evidence; accepted historical data keeps its original destination. Existing advanced groups remain outside automatic conversion.


<!-- @trace
source: streamline-studio-v2-recording-setup
updated: 2026-10-04
code:
  - internal/datalink/runtime/group_boundary.go
  - internal/datalink/recordingplan/schema_preview.go
  - docs/plans/studio-v2-flow-completion/evidence/review-repairs-uint64-postgres.png
  - internal/datalink/runtime/service_workspace_projection.go
  - internal/datalink/runtime/service_workspace_projection_scope.go
  - docs/plans/studio-v2-flow-completion/evidence-c/error-1440.png
  - internal/datalink/workspace/write_group.go
  - internal/datalink/workspace/write_group_managed.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - frontend/src/utils/safeJson.ts
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/grouppipeline/pipeline.go
  - docs/plans/studio-v2-flow-completion/evidence-e/partial-768.png
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - internal/datalink/migrator_recording_plans.go
  - docs/plans/studio-v2-flow-completion/evidence-c/error-390.png
  - internal/datalink/workspace/write_group_repository.go
  - frontend/src/types/recordingPlan.ts
  - internal/datalink/recordingplan/repository_memory.go
  - internal/datalink/grouppipeline/recovery.go
  - frontend/src/hooks/datalink/useStudioV2RecordingStart.ts
  - internal/datalink/dbtarget/group_row_layout.go
  - internal/datalink/dbtarget/service.go
  - docs/plans/studio-v2-flow-completion/E-setup-verification.md
  - internal/datalink/dbtarget/managed_destination.go
  - internal/datalink/workspace/service_readiness_scope.go
  - docs/plans/studio-v2-flow-completion/evidence-d/partial-390.png
  - internal/datalink/modbusshare/service_recording_scope.go
  - docs/plans/studio-v2-flow-completion/evidence-c/error-768.png
  - internal/datalink/recordingplan/schema_group.go
  - internal/datalink/workspace/write_group_readiness.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - internal/datalink/groupdelivery/runtime_versions.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - internal/datalink/workspace/write_group_readiness_table.go
  - internal/datalink/workspace/recording_start_steps.go
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/columns.ts
  - frontend/src/types/studioV2RecordingStart.ts
  - docs/plans/studio-v2-flow-completion/review-repairs.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupColumnProposal.tsx
  - internal/datalink/workspace/write_group_schema_confirm.go
  - internal/datalink/workspace/write_group_service.go
  - docs/plans/studio-v2-flow-completion/evidence/review-repairs-uint64-reject-bigint.png
  - docs/plans/studio-v2-flow-completion/evidence-e/error-1440.png
  - frontend/src/utils/studioV2RecordingStartJson.ts
  - internal/datalink/recordingplan/repository_recording_start.go
  - docs/technical/studio-v2-write-groups.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/plans/studio-v2-flow-completion/B-deadline-verification.md
  - docs/plans/studio-v2-flow-completion/evidence-d/partial-768.png
  - internal/datalink/workspace/write_group_runtime_layout.go
  - docs/plans/studio-v2-flow-completion/evidence-e/runtime-handoff-390.png
  - docs/plans/studio-v2-flow-completion/evidence-d/error-390.png
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecordingMessages.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/dbtarget/table_inspection_readonly.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - internal/datalink/schema/migrations/028_write_group_basic_keys_sqlite.up.sql
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/workspace/service_activation_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingEvidence.tsx
  - internal/datalink/workspace/write_group_schema.go
  - docs/plans/studio-v2-flow-completion/evidence-c/main-390.png
  - internal/datalink/workspace/service_readiness_downstream.go
  - docs/plans/studio-v2-flow-completion/evidence-d/main-768.png
  - internal/datalink/grouppipeline/delivery_view.go
  - internal/datalink/grouppipeline/reconcile.go
  - docs/plans/studio-v2-flow-completion/evidence-d/error-768.png
  - internal/datalink/dbtarget/postgres_column_inspection.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/runtime/group_boundary_encoding.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/utils/studioV2WriteGroupJson.ts
  - internal/api/handlers/studio_v2_workspace_recording_plan_apply.go
  - internal/api/router_studio_v2_write_group_routes.go
  - docs/plans/studio-v2-flow-completion/evidence-e/main-1440.png
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/workspace/recording_start_types.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingProgress.tsx
  - docs/plans/studio-v2-flow-completion/evidence-e/main-768.png
  - internal/datalink/runtime/group_boundary_layout.go
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/recordingplan/repository.go
  - internal/datalink/recordingplan/schema_apply.go
  - docs/plans/studio-v2-flow-completion/evidence-e/partial-1440.png
  - docs/plans/studio-v2-flow-completion/evidence-c/main-768.png
  - internal/datalink/workspace/write_group_readiness_partial.go
  - internal/datalink/groupdelivery/sender.go
  - docs/plans/studio-v2-flow-completion/evidence-c/partial-1440.png
  - docs/plans/studio-v2-flow-completion/evidence-c/main-1440.png
  - docs/plans/studio-v2-flow-completion/D-managed-verification.md
  - docs/plans/studio-v2-flow-completion/evidence-e/share-only-1440.png
  - docs/plans/studio-v2-flow-completion/evidence/review-repairs-uint64-sqlite.png
  - internal/api/handlers/studio_v2_recording_start_handler.go
  - docs/plans/studio-v2-flow-completion/README.md
  - internal/datalink/recordingplan/repository_preview_token.go
  - internal/datalink/dbtarget/managed_schema_execution.go
  - docs/swagger/swagger.json
  - frontend/src/utils/recordingPlanJson.ts
  - internal/datalink/schema/migrations/027_write_group_runtime_sqlite.up.sql
  - docs/plans/studio-v2-flow-completion/evidence-e/share-only-390.png
  - docs/plans/studio-v2-flow-completion/evidence-d/main-390.png
  - internal/datalink/workspace/write_group_readiness_eval.go
  - internal/datalink/workspace/service_readiness_write_groups.go
  - internal/datalink/workspace/write_group_lifecycle.go
  - frontend/src/types/studioV2WriteGroup.ts
  - internal/api/handlers/studio_v2_workspace_database_metadata.go
  - internal/datalink/collector/scheduler_config_methods.go
  - docs/plans/studio-v2-flow-completion/evidence-e/offline-neighbor-1440.png
  - docs/plans/studio-v2-flow-completion/evidence-e/offline-neighbor-step4-1440.png
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - docs/plans/studio-v2-flow-completion/evidence-e/main-390.png
  - docs/plans/studio-v2-flow-completion/C-entity-verification.md
  - internal/api/router_studio_v2_recording_start.go
  - docs/plans/studio-v2-flow-completion/evidence-d/error-1440.png
  - internal/datalink/recordingplan/types.go
  - frontend/src/utils/studioV2WriteGroupDeliveryJson.ts
  - frontend/src/services/studioV2RecordingStart.ts
  - frontend/src/utils/typedErrors.ts
  - docs/plans/studio-v2-flow-completion/evidence-e/partial-390.png
  - frontend/src/utils/backendErrorCodes.ts
  - docs/plans/studio-v2-flow-completion/evidence-e/waiting-1440.png
  - internal/api/router.go
  - docs/plans/studio-v2-flow-completion/evidence-c/partial-390.png
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - internal/datalink/groupdelivery/committed_effect.go
  - docs/plans/studio-v2-flow-completion/evidence-e/error-768.png
  - docs/plans/studio-v2-flow-completion/evidence-d/main-1440.png
  - internal/datalink/groupdelivery/status.go
  - internal/datalink/workspace/recording_start.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns.ts
  - internal/api/handlers/studio_v2_workspace_write_group_schema_handler.go
  - docs/plans/studio-v2-flow-completion/evidence-e/error-390.png
  - docs/plans/studio-v2-flow-completion/evidence-e/share-only-768.png
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/proposal.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/datalink/dbtarget/schema_execution.go
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupMemberTable.tsx
  - docs/plans/studio-v2-flow-completion/A-lifecycle-verification.md
  - internal/datalink/recordingplan/recording_start_operation.go
  - internal/datalink/workspace/recording_start_scope.go
  - docs/plans/studio-v2-flow-completion/evidence-d/partial-1440.png
  - docs/plans/studio-v2-flow-completion/evidence-e/mapping-1440.png
  - internal/datalink/workspace/service_database_setup.go
  - docs/plans/studio-v2-flow-completion/evidence-c/partial-768.png
  - internal/datalink/workspace/write_group_basic_managed.go
  - docs/plans/studio-v2-flow-completion/evidence-e/offline-draft-1440.png
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
tests:
  - cmd/test_ui/group_pipeline_managed_sql_test.go
  - internal/datalink/modbusshare/service_recording_scope_test.go
  - internal/datalink/groupdelivery/sender_postgres_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_write_group_schema_conflict_test.go
  - frontend/tests/unit/hooks/useStudioV2GroupMetadata.test.tsx
  - frontend/tests/unit/utils/recordingSchemaPreviewToken.test.ts
  - frontend/tests/unit/workbench-v2/writeGroupManagedStorage.test.ts
  - internal/datalink/grouppipeline/reconcile_managed_recovery_test.go
  - internal/datalink/workspace/write_group_managed_schema_guard_test.go
  - frontend/tests/unit/services/studioV2RecordingStart.test.ts
  - cmd/test_ui/group_pipeline_managed_test.go
  - internal/datalink/grouppipeline/reconcile_recovery_test.go
  - internal/datalink/groupdelivery/sender_deadline_recovery_test.go
  - internal/datalink/workspace/write_group_runtime_layout_test.go
  - internal/datalink/runtime/group_boundary_cutoff_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/writeGroupSection.test.tsx
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/router_studio_v2_recording_start_failures_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx
  - internal/datalink/dbtarget/group_row_layout_test.go
  - internal/api/router_studio_v2_write_group_schema_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - internal/api/router_studio_v2_write_group_schema_drift_test.go
  - internal/datalink/groupdelivery/sender_deadline_test.go
  - internal/api/router_studio_v2_recording_start_barrier_test.go
  - cmd/test_ui/group_pipeline_entity_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/writeGroupDraft.test.ts
  - internal/datalink/workspace/write_group_reenable_test.go
  - cmd/test_ui/group_pipeline_managed_quality_test.go
  - internal/api/router_studio_v2_write_group_metadata_readonly_test.go
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroupSchema.test.ts
  - internal/datalink/recordingplan/schema_group_numeric_test.go
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
  - frontend/tests/unit/utils/backendErrorCodes.test.ts
  - internal/api/router_studio_v2_delivery_effect_test.go
  - cmd/test_ui/recording_start_share_test.go
  - frontend/tests/unit/workbench-v2/writeGroupProposal.test.tsx
  - internal/datalink/groupdelivery/committed_effect_test.go
  - internal/api/router_studio_v2_write_group_schema_partial_test.go
  - frontend/tests/unit/workbench-v2/writeGroupMembers.test.tsx
  - internal/datalink/recordingplan/recording_start_operation_test.go
  - cmd/test_ui/group_pipeline_managed_multidevice_test.go
  - internal/datalink/runtime/group_boundary_entity_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel-regressions.test.tsx
  - frontend/tests/fixtures/managedSchemaPreview.ts
  - internal/api/handlers/studio_v2_workspace_database_metadata_readonly_test.go
  - internal/datalink/recordingplan/schema_group_test.go
  - frontend/tests/unit/services/studioV2WriteGroupDelivery.test.ts
  - internal/datalink/workspace/service_readiness_scope_test.go
  - internal/api/router_studio_v2_write_group_schema_replay_test.go
  - cmd/test_ui/service_wiring.go
  - internal/datalink/workspace/service_activation_scope_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-state.test.ts
  - internal/datalink/workspace/write_group_readiness_table_test.go
  - internal/api/handlers/studio_v2_workspace_metadata_scope_test.go
  - cmd/test_ui/recording_start_share.go
  - internal/datalink/workspace/write_group_readiness_managed_inspector_test.go
  - internal/api/handlers/studio_v2_workspace_group_metadata_test.go
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/runtime/group_boundary_durable_test.go
  - internal/datalink/dbtarget/managed_destination_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_share_create_test.go
  - internal/datalink/dbtarget/postgres_numeric_inspection_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/writeGroupColumns.test.ts
  - internal/api/router_studio_v2_recording_start_recovery_test.go
  - internal/api/router_studio_v2_recording_start_initial_test.go
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/dbtarget/group_row_entity_test.go
  - internal/datalink/dbtarget/managed_inspection_scope_test.go
  - internal/datalink/runtime/service_workspace_projection_scope_test.go
  - internal/api/router_studio_v2_write_group_schema_scope_test.go
  - internal/datalink/grouppipeline/reconcile_revision_test.go
  - frontend/tests/unit/services/studioV2WorkspaceDatabaseMetadata.test.ts
  - cmd/test_ui/group_pipeline_entity_sql_test.go
  - internal/datalink/dbtarget/managed_schema_execution_test.go
  - internal/api/router_studio_v2_recording_start_test.go
  - internal/datalink/runtime/service_workspace_projection_shared_scope_test.go
  - internal/datalink/groupdelivery/sender_settlement_test.go
  - internal/api/router_studio_v2_recording_start_fixture_test.go
  - internal/api/router_studio_v2_write_group_schema_recovery_test.go
-->

---
### Requirement: Safe reusable templates and overrides
The system SHALL preview template changes per device and preserve confirmed user overrides and stable identities.

#### Scenario: Reapply a template
- **WHEN** a template is applied again after a user edits a pressure multiplier
- **THEN** the difference is shown and the multiplier is not silently overwritten or duplicated.

---
### Requirement: Natural-language recording policies
The system SHALL distinguish acquisition, raw retention, summaries, transport batching and history retention with examples and a final readable summary.

#### Scenario: Five-second polling and one-minute summaries
- **WHEN** those settings are selected
- **THEN** the summary explains that sampled readings feed the minute statistics rather than recording only one instantaneous value per minute.

---
### Requirement: Truthful save and activation states
The system MUST separate saved draft, applied revision, collecting, durable, delivered and verified states, including bounded waiting and partial failure.

#### Scenario: Save during an active plan
- **WHEN** the operator edits a running plan
- **THEN** the current applied plan remains active until explicit validated apply and both versions are visible.

---
### Requirement: Actionable accessible failures
The system SHALL provide field-level reasons, focused repair actions and keyboard-accessible controls instead of unexplained disabled buttons or raw technical errors.

#### Scenario: Write permission denied
- **WHEN** a connected database refuses write access
- **THEN** the operator sees which destination needs permission repair and may retry it without redoing successful device setup.

---
### Requirement: Optional outputs remain independent
The workflow MUST allow a selected Local Modbus-only setup without requiring database configuration while retaining the existing Share activation gates.

#### Scenario: No database output selected
- **WHEN** the user selects only validated Local Modbus forwarding
- **THEN** database fields and tests are skipped and no database writer starts.

---
### Requirement: Complete basic preparation with one effective configuration
Basic setup SHALL complete managed schema preparation without entering Advanced. Connection controls SHALL configure only connection identity; the canonical group SHALL own recording members, storage and interval. Operator primary surfaces SHALL show names, values, types, units and truthful recording stages; technical IDs, API payloads and revisions SHALL remain available in diagnostics.

#### Scenario: Fresh basic recording
- **WHEN** an operator configures eight measurement points with type-appropriate register spans, maps tags and selects an empty managed SQLite destination
- **THEN** Basic provides schema preview and explicit confirmation and proceeds to timed one-row recording without Advanced.

#### Scenario: Inspect effective settings
- **WHEN** an operator prepares a canonical group
- **THEN** no competing connector table, INSERT/UPSERT or recording interval control is presented as effective, and technical identifiers remain inspectable through diagnostics.
