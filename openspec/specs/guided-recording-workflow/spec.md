# guided-recording-workflow Specification

## Purpose

TBD - created by archiving change 'redesign-studio-recording-flow'. Update Purpose after archive.

## Requirements

### Requirement: Four-step intent-led setup
Operators MUST be able to save and navigate offline drafts with unresolved probes; draft navigation MUST NOT imply verified readiness, and only dependent reads or activation SHALL be blocked. The existing Studio V2 route SHALL guide device connection, acquisition points, Point-to-Tag mapping and grouped database output in four steps labelled 連接設備 → 選要讀的資料 → 命名與換算 → 儲存並開始記錄. Basic/Advanced presentation changes SHALL preserve configuration values and protocol/type capabilities. Basic typed Tag writing MUST remain usable without required reporting, retention setup, aggregation or physical measurement semantics. Explicitly selected derived operations SHALL require their own confirmed semantics. The system MUST retain the Go backend and existing React framework and MUST NOT add a parallel V3 or restore a dedicated legacy /studio route. The basic path SHALL derive and persist routine mappings, show only necessary operator choices, and offer whole-device managed wide snapshots or selected-source managed long recording in the selected device scope without converting existing configurations. Table layout SHALL be independent of instantaneous or cumulative semantics. Existing canonical same-kind ownership SHALL remain unique and existing latest snapshots SHALL not become means through navigation.

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
source: unify-device-recording-workflow
updated: 2026-10-25
code:
  - internal/datalink/workspace/write_group_long_recording_shape.go
  - internal/datalink/workspace/write_group_member_metadata.go
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/features/datalink/workbench-v2/state/basicLegacyDraft.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/technical/studio-v2-circuit-recording.md
  - internal/datalink/grouppipeline/delivery_sources.go
  - internal/datalink/groupdelivery/source_evidence_wide.go
  - internal/datalink/workspace/write_group_cutover_repository.go
  - internal/datalink/measurement/mean.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - internal/datalink/dbtarget/group_long_row.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/workspace/write_group_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitPurposeFields.tsx
  - internal/datalink/workspace/write_group_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_wide_profile_shape.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitRequiredPolicy.tsx
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/swagger/docs.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/workspace/write_group_repository_validation.go
  - frontend/src/features/datalink/workbench-v2/state/deviceScopeSelection.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/runtime/group_boundary_wide_profile.go
  - internal/datalink/groupdelivery/source_evidence_aggregation.go
  - internal/datalink/snapshot/config.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - internal/datalink/runtime/group_boundary_build.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/groupdelivery/source_evidence.go
  - docs/swagger/swagger.json
  - internal/datalink/grouppipeline/recovery_long.go
  - internal/datalink/workspace/write_group_readiness.go
  - internal/datalink/snapshot/assembler_profile.go
  - internal/datalink/snapshot/outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SourceReview.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/types/studioV2WriteGroup.ts
  - internal/datalink/dbtarget/group_row_layout.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/runtime/group_boundary.go
  - internal/datalink/schema/migrations/038_write_group_recording_reduction_sqlite.up.sql
  - internal/datalink/snapshot/assembler_counter.go
  - frontend/src/utils/studioV2RecordingProfile.ts
  - internal/datalink/snapshot/assembler_mean.go
  - internal/datalink/snapshot/mean_evidence.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - internal/datalink/workspace/write_group_counter_readiness.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/workspace/write_group_profile_interval.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - internal/datalink/grouppipeline/recovery.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - internal/datalink/snapshot/assembler_wide_profile.go
  - internal/datalink/runtime/group_boundary_long.go
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/assembler.go
  - internal/datalink/workspace/write_group_cutover_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeviceScopeTiles.tsx
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
tests:
  - internal/datalink/grouppipeline/mean_result_type_proof_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - internal/api/router_studio_v2_reviewed_wide_mean_transition_test.go
  - internal/datalink/snapshot/mean_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/api/router_studio_v2_wide_profile_guard_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - internal/api/router_studio_v2_recording_start_scope_test.go
  - internal/api/router_studio_v2_wide_profile_runtime_test.go
  - frontend/tests/unit/workbench-v2/write-group-aggregation-evidence.test.ts
  - internal/datalink/dbtarget/group_long_row_mean_test.go
  - internal/datalink/dbtarget/group_row_layout_wide_profile_test.go
  - internal/datalink/grouppipeline/average_recovery_test.go
  - internal/datalink/grouppipeline/recovery_reduction_test.go
  - internal/datalink/runtime/group_boundary_wide_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - frontend/tests/unit/workbench-v2/step4-evidence-boundaries.test.tsx
  - internal/datalink/grouppipeline/wide_profile_result_type_proof_test.go
  - internal/datalink/workspace/write_group_member_recording_shapes_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-tiles.test.tsx
  - frontend/tests/unit/workbench-v2/wide-profile-member-metadata.test.ts
  - frontend/tests/e2e/fixtures/device-recording-workflow.ts
  - frontend/tests/unit/services/studioV2WriteGroupDeliveryEnvelope.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose.test.tsx
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/runtime/group_boundary_long_reduction_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - internal/datalink/runtime/group_boundary_long_counter_raw_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - frontend/tests/e2e/device-recording-workflow.spec.ts
  - internal/api/router_studio_v2_recording_workflow_coexistence_test.go
  - internal/datalink/snapshot/assembler_wide_cumulative_test.go
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_legacy_provenance_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/basic-recording-panel-regressions.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/datalink/grouppipeline/average_fixture_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/datalink/grouppipeline/average_delivery_evidence_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/snapshot/assembler_wide_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/datalink/workspace/write_group_member_metadata_test.go
  - frontend/tests/unit/workbench-v2/profile-reduction-equivalence.test.ts
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/recording-workflow-unified.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recording_reduction_test.go
  - frontend/tests/unit/workbench-v2/device-recording-summary.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/router_studio_v2_wide_counter_parity_test.go
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/datalink/grouppipeline/average_unrelated_test.go
  - internal/datalink/groupdelivery/source_evidence_wide_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_recording_start_fixture_test.go
  - frontend/tests/fixtures/cutoverFixtures.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - internal/datalink/dbtarget/group_row_layout_wide_counter_test.go
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - internal/datalink/measurement/mean_test.go
  - internal/api/router_studio_v2_mean_long_runtime_test.go
  - internal/datalink/workspace/write_group_wide_cumulative_shape_test.go
  - frontend/tests/unit/workbench-v2/writeGroupProposal.test.tsx
  - internal/datalink/runtime/group_boundary_wide_profile_test.go
  - internal/datalink/snapshot/mean_evidence_test.go
  - frontend/tests/unit/workbench-v2/device-recording-purpose-flow.test.tsx
  - internal/api/device_recording_integration_test.go
  - frontend/tests/fixtures/sourceReview.ts
  - frontend/tests/unit/workbench-v2/schema-operation.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
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
Basic setup SHALL complete managed schema preparation without entering Advanced. Connection controls SHALL configure only connection identity; the canonical group SHALL own recording members, layout, storage, interval and reduction. New numeric gauge recording SHALL default to the explicitly labelled arithmetic mean; existing latest and counter roles SHALL retain their persisted interpretation. Operator primary surfaces SHALL show names, values, types, units and truthful recording stages; technical IDs, API payloads and revisions SHALL remain available in diagnostics.

#### Scenario: Fresh basic recording
- **WHEN** an operator configures eight measurement points with type-appropriate register spans, maps tags and selects an empty managed SQLite destination
- **THEN** Basic provides schema preview and explicit confirmation and proceeds through the selected wide or long recording flow without Advanced and without an additional workspace activation action.

#### Scenario: Inspect effective settings
- **WHEN** an operator prepares a canonical group
- **THEN** no competing connector table, INSERT/UPSERT or recording interval control is presented as effective, and technical identifiers remain inspectable through diagnostics.

#### Scenario: Destination before recording choices
- **WHEN** no database destination is saved and Share is not enabled
- **THEN** Basic first presents destination setup before recording members, frequency and start; a saved destination is shown as a summary with an edit entry preserving advanced connection and password semantics.

#### Scenario: Prepare action and focused repair
- **WHEN** Basic preparation is requested or a prepare_schema error needs repair
- **THEN** one preparation action loads the read-only preview in the Basic schema card, with explicit DDL confirmation still required and no forced switch into all Advanced settings.


<!-- @trace
source: unify-device-recording-workflow
updated: 2026-10-25
code:
  - internal/datalink/workspace/write_group_long_recording_shape.go
  - internal/datalink/workspace/write_group_member_metadata.go
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/features/datalink/workbench-v2/state/basicLegacyDraft.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/technical/studio-v2-circuit-recording.md
  - internal/datalink/grouppipeline/delivery_sources.go
  - internal/datalink/groupdelivery/source_evidence_wide.go
  - internal/datalink/workspace/write_group_cutover_repository.go
  - internal/datalink/measurement/mean.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - internal/datalink/dbtarget/group_long_row.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/workspace/write_group_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitPurposeFields.tsx
  - internal/datalink/workspace/write_group_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_wide_profile_shape.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitRequiredPolicy.tsx
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/swagger/docs.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/workspace/write_group_repository_validation.go
  - frontend/src/features/datalink/workbench-v2/state/deviceScopeSelection.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/runtime/group_boundary_wide_profile.go
  - internal/datalink/groupdelivery/source_evidence_aggregation.go
  - internal/datalink/snapshot/config.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - internal/datalink/runtime/group_boundary_build.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/groupdelivery/source_evidence.go
  - docs/swagger/swagger.json
  - internal/datalink/grouppipeline/recovery_long.go
  - internal/datalink/workspace/write_group_readiness.go
  - internal/datalink/snapshot/assembler_profile.go
  - internal/datalink/snapshot/outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SourceReview.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/types/studioV2WriteGroup.ts
  - internal/datalink/dbtarget/group_row_layout.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/runtime/group_boundary.go
  - internal/datalink/schema/migrations/038_write_group_recording_reduction_sqlite.up.sql
  - internal/datalink/snapshot/assembler_counter.go
  - frontend/src/utils/studioV2RecordingProfile.ts
  - internal/datalink/snapshot/assembler_mean.go
  - internal/datalink/snapshot/mean_evidence.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - internal/datalink/workspace/write_group_counter_readiness.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/workspace/write_group_profile_interval.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - internal/datalink/grouppipeline/recovery.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - internal/datalink/snapshot/assembler_wide_profile.go
  - internal/datalink/runtime/group_boundary_long.go
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/assembler.go
  - internal/datalink/workspace/write_group_cutover_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeviceScopeTiles.tsx
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
tests:
  - internal/datalink/grouppipeline/mean_result_type_proof_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - internal/api/router_studio_v2_reviewed_wide_mean_transition_test.go
  - internal/datalink/snapshot/mean_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/api/router_studio_v2_wide_profile_guard_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - internal/api/router_studio_v2_recording_start_scope_test.go
  - internal/api/router_studio_v2_wide_profile_runtime_test.go
  - frontend/tests/unit/workbench-v2/write-group-aggregation-evidence.test.ts
  - internal/datalink/dbtarget/group_long_row_mean_test.go
  - internal/datalink/dbtarget/group_row_layout_wide_profile_test.go
  - internal/datalink/grouppipeline/average_recovery_test.go
  - internal/datalink/grouppipeline/recovery_reduction_test.go
  - internal/datalink/runtime/group_boundary_wide_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - frontend/tests/unit/workbench-v2/step4-evidence-boundaries.test.tsx
  - internal/datalink/grouppipeline/wide_profile_result_type_proof_test.go
  - internal/datalink/workspace/write_group_member_recording_shapes_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-tiles.test.tsx
  - frontend/tests/unit/workbench-v2/wide-profile-member-metadata.test.ts
  - frontend/tests/e2e/fixtures/device-recording-workflow.ts
  - frontend/tests/unit/services/studioV2WriteGroupDeliveryEnvelope.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose.test.tsx
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/runtime/group_boundary_long_reduction_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - internal/datalink/runtime/group_boundary_long_counter_raw_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - frontend/tests/e2e/device-recording-workflow.spec.ts
  - internal/api/router_studio_v2_recording_workflow_coexistence_test.go
  - internal/datalink/snapshot/assembler_wide_cumulative_test.go
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_legacy_provenance_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/basic-recording-panel-regressions.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/datalink/grouppipeline/average_fixture_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/datalink/grouppipeline/average_delivery_evidence_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/snapshot/assembler_wide_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/datalink/workspace/write_group_member_metadata_test.go
  - frontend/tests/unit/workbench-v2/profile-reduction-equivalence.test.ts
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/recording-workflow-unified.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recording_reduction_test.go
  - frontend/tests/unit/workbench-v2/device-recording-summary.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/router_studio_v2_wide_counter_parity_test.go
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/datalink/grouppipeline/average_unrelated_test.go
  - internal/datalink/groupdelivery/source_evidence_wide_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_recording_start_fixture_test.go
  - frontend/tests/fixtures/cutoverFixtures.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - internal/datalink/dbtarget/group_row_layout_wide_counter_test.go
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - internal/datalink/measurement/mean_test.go
  - internal/api/router_studio_v2_mean_long_runtime_test.go
  - internal/datalink/workspace/write_group_wide_cumulative_shape_test.go
  - frontend/tests/unit/workbench-v2/writeGroupProposal.test.tsx
  - internal/datalink/runtime/group_boundary_wide_profile_test.go
  - internal/datalink/snapshot/mean_evidence_test.go
  - frontend/tests/unit/workbench-v2/device-recording-purpose-flow.test.tsx
  - internal/api/device_recording_integration_test.go
  - frontend/tests/fixtures/sourceReview.ts
  - frontend/tests/unit/workbench-v2/schema-operation.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
-->

---
### Requirement: Selection is separate from acquisition

Step 2 point click SHALL select or inspect points without changing acquisition. Multi-selection, modifier selection and contiguous range selection SHALL remain available. Acquisition changes SHALL require an explicit collection checkbox or a clearly scoped selected-point enable/skip action. Read range descriptions SHALL distinguish item count from register/bit span, expose required types and offer editable name suggestions.

#### Scenario: Inspect and select a range
- **WHEN** the operator clicks points and uses Shift or Ctrl to select multiple points
- **THEN** only selection/detail state changes; collection remains unchanged until an explicit collection action.

#### Scenario: Multiword types are explained
- **WHEN** a range contains four float32 items
- **THEN** the read summary identifies four data items occupying eight words and retains the required type control.


<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Narrow basic managed membership updates

Basic SHALL present a reviewed member-addition/removal difference for eligible canonical managed groups and save it through existing workspace/group/connector version checks. Existing columns, required flags, deduplication, entity/provenance rules, frozen historical data and custom groups SHALL NOT be overwritten. Customized or unsaved-ineligible mappings SHALL be blocked or directed to existing Advanced controls. Required storage changes SHALL use the existing explicit DDL preview rather than an implicit mutation.

#### Scenario: Add a newly mapped point
- **WHEN** an eligible basic managed group has a newly saved point mapping
- **THEN** Basic offers the narrow member difference and explicit update, preserves existing group policies/columns, and links to confirmed schema preparation when a column is missing.

#### Scenario: Existing customization or stale version
- **WHEN** the group has custom columns/entity/deduplication semantics or the expected version is stale
- **THEN** Basic does not overwrite the configuration; existing advanced editing or version repair remains available.


<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Recording claims follow committed evidence

Runtime diagnostics with no failure_stage SHALL state only that no error was detected. Acquisition, waiting for the first record and confirmed SQL commit SHALL remain distinguishable. Basic recording success SHALL require the existing workspace/group/applied revision/connector-scoped committed-effect evidence. Share-only success or acquisition alone SHALL NOT claim all points reached a database.

#### Scenario: Collection without SQL evidence
- **WHEN** reads or Share delivery succeed but no matching committed SQL effect is observed
- **THEN** diagnostics do not claim database delivery and Basic continues to show the applicable waiting or evidence state.

#### Scenario: Scoped database commit
- **WHEN** matching committed-effect evidence exists for the effective group and connector
- **THEN** Basic can show confirmed database writing using that evidence without treating missing errors as proof.

<!-- @trace
source: close-studio-v2-simplification
updated: 2026-10-23
code:
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-logs-20261023.json.gz
  - docs/releases/studio-v2-change-closeout-20261023.md
  - docs/plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png
  - docs/plans/studio-v2-flow-completion/README.md
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/recordingplan/schema_group_extend_test.go
  - frontend/tests/unit/workbench-v2/schemaOperationRecovery.test.ts
  - frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx
-->

---
### Requirement: Basic missing-value policy reveals required gating
Wide and custom recording SHALL state that required missing or invalid members still block a partial row, show the effective required state, and provide a narrow entry to existing per-member required controls without silently changing any existing member. New mean wide groups SHALL default to optional partial recording only with verified nullable/provenance support. Long recording SHALL show independent per-source missing/coverage results and SHALL NOT offer a whole-row skip policy. Counter baseline or unknown delta SHALL not be described as a zero value or a generic save failure.

#### Scenario: Required gap remains blocked
- **WHEN** a basic group selects partial and a required member is missing or invalid
- **THEN** the UI SHALL explain that the row is skipped and SHALL expose the existing member settings without changing required automatically.

#### Scenario: Optional gap and recovery
- **WHEN** an explicitly optional member is missing in a partial-capable destination and subsequently recovers
- **THEN** the system SHALL preserve NULL and its reason for the partial row and write its actual value after recovery; skip_row and required protection SHALL remain unchanged.

#### Scenario: Edit existing exact managed members
- **WHEN** an operator opens required settings for saved SQLite managed uint16 or uint32 members in canonical generated TEXT columns
- **THEN** the editor SHALL retain their exact storage compatibility and require fresh schema confirmation after a group revision change without weakening custom-table checks.

#### Scenario: Confirm unchanged schema after member policy edit
- **WHEN** a saved managed group changes required or completeness policy without changing its physical columns
- **THEN** a compatible preview SHALL return an empty statements array, retain its revision-bound token and operation identity, and require explicit confirmation before applying the new group revision.


<!-- @trace
source: unify-device-recording-workflow
updated: 2026-10-25
code:
  - internal/datalink/workspace/write_group_long_recording_shape.go
  - internal/datalink/workspace/write_group_member_metadata.go
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/features/datalink/workbench-v2/state/basicLegacyDraft.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/technical/studio-v2-circuit-recording.md
  - internal/datalink/grouppipeline/delivery_sources.go
  - internal/datalink/groupdelivery/source_evidence_wide.go
  - internal/datalink/workspace/write_group_cutover_repository.go
  - internal/datalink/measurement/mean.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - internal/datalink/dbtarget/group_long_row.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/workspace/write_group_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitPurposeFields.tsx
  - internal/datalink/workspace/write_group_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_wide_profile_shape.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitRequiredPolicy.tsx
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/swagger/docs.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/workspace/write_group_repository_validation.go
  - frontend/src/features/datalink/workbench-v2/state/deviceScopeSelection.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/runtime/group_boundary_wide_profile.go
  - internal/datalink/groupdelivery/source_evidence_aggregation.go
  - internal/datalink/snapshot/config.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - internal/datalink/runtime/group_boundary_build.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/groupdelivery/source_evidence.go
  - docs/swagger/swagger.json
  - internal/datalink/grouppipeline/recovery_long.go
  - internal/datalink/workspace/write_group_readiness.go
  - internal/datalink/snapshot/assembler_profile.go
  - internal/datalink/snapshot/outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SourceReview.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/types/studioV2WriteGroup.ts
  - internal/datalink/dbtarget/group_row_layout.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/runtime/group_boundary.go
  - internal/datalink/schema/migrations/038_write_group_recording_reduction_sqlite.up.sql
  - internal/datalink/snapshot/assembler_counter.go
  - frontend/src/utils/studioV2RecordingProfile.ts
  - internal/datalink/snapshot/assembler_mean.go
  - internal/datalink/snapshot/mean_evidence.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - internal/datalink/workspace/write_group_counter_readiness.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/workspace/write_group_profile_interval.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - internal/datalink/grouppipeline/recovery.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - internal/datalink/snapshot/assembler_wide_profile.go
  - internal/datalink/runtime/group_boundary_long.go
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/assembler.go
  - internal/datalink/workspace/write_group_cutover_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeviceScopeTiles.tsx
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
tests:
  - internal/datalink/grouppipeline/mean_result_type_proof_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - internal/api/router_studio_v2_reviewed_wide_mean_transition_test.go
  - internal/datalink/snapshot/mean_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/api/router_studio_v2_wide_profile_guard_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - internal/api/router_studio_v2_recording_start_scope_test.go
  - internal/api/router_studio_v2_wide_profile_runtime_test.go
  - frontend/tests/unit/workbench-v2/write-group-aggregation-evidence.test.ts
  - internal/datalink/dbtarget/group_long_row_mean_test.go
  - internal/datalink/dbtarget/group_row_layout_wide_profile_test.go
  - internal/datalink/grouppipeline/average_recovery_test.go
  - internal/datalink/grouppipeline/recovery_reduction_test.go
  - internal/datalink/runtime/group_boundary_wide_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - frontend/tests/unit/workbench-v2/step4-evidence-boundaries.test.tsx
  - internal/datalink/grouppipeline/wide_profile_result_type_proof_test.go
  - internal/datalink/workspace/write_group_member_recording_shapes_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-tiles.test.tsx
  - frontend/tests/unit/workbench-v2/wide-profile-member-metadata.test.ts
  - frontend/tests/e2e/fixtures/device-recording-workflow.ts
  - frontend/tests/unit/services/studioV2WriteGroupDeliveryEnvelope.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose.test.tsx
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/runtime/group_boundary_long_reduction_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - internal/datalink/runtime/group_boundary_long_counter_raw_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - frontend/tests/e2e/device-recording-workflow.spec.ts
  - internal/api/router_studio_v2_recording_workflow_coexistence_test.go
  - internal/datalink/snapshot/assembler_wide_cumulative_test.go
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_legacy_provenance_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/basic-recording-panel-regressions.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/datalink/grouppipeline/average_fixture_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/datalink/grouppipeline/average_delivery_evidence_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/snapshot/assembler_wide_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/datalink/workspace/write_group_member_metadata_test.go
  - frontend/tests/unit/workbench-v2/profile-reduction-equivalence.test.ts
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/recording-workflow-unified.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recording_reduction_test.go
  - frontend/tests/unit/workbench-v2/device-recording-summary.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/router_studio_v2_wide_counter_parity_test.go
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/datalink/grouppipeline/average_unrelated_test.go
  - internal/datalink/groupdelivery/source_evidence_wide_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_recording_start_fixture_test.go
  - frontend/tests/fixtures/cutoverFixtures.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - internal/datalink/dbtarget/group_row_layout_wide_counter_test.go
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - internal/datalink/measurement/mean_test.go
  - internal/api/router_studio_v2_mean_long_runtime_test.go
  - internal/datalink/workspace/write_group_wide_cumulative_shape_test.go
  - frontend/tests/unit/workbench-v2/writeGroupProposal.test.tsx
  - internal/datalink/runtime/group_boundary_wide_profile_test.go
  - internal/datalink/snapshot/mean_evidence_test.go
  - frontend/tests/unit/workbench-v2/device-recording-purpose-flow.test.tsx
  - internal/api/device_recording_integration_test.go
  - frontend/tests/fixtures/sourceReview.ts
  - frontend/tests/unit/workbench-v2/schema-operation.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
-->

---
### Requirement: Simple circuit choices and scoped batch edits

Basic SHALL provide independently labelled instantaneous and cumulative controls, their interval units and the 15-second and 10-minute defaults. Operators SHALL be able to select several circuits and apply a shared profile or interval while retaining unselected circuits. Existing individual overrides SHALL be listed before the operator explicitly confirms their replacement. Cumulative selection SHALL explain that both raw readings and differences are stored, with baseline and uncertain differences displayed as unknown rather than zero. Controls and errors SHALL be keyboard accessible and localized in en and zh-TW.

#### Scenario: Batch edit with override
- **WHEN** circuits A and B are selected for a 30-second instantaneous interval, B already has an individual override, and C is unselected
- **THEN** Basic previews B's override replacement, saves A and B only after confirmation, preserves C, and retains these choices after reload.

#### Scenario: Cumulative first record
- **WHEN** the first cumulative raw reading has reached SQL but no comparable prior baseline exists
- **THEN** the interface distinguishes successful raw recording from an unknown first interval difference and does not display zero consumption or a write failure.


<!-- @trace
source: add-circuit-recording-profiles
updated: 2026-10-25
code:
  - frontend/src/utils/recordingPlanJson.ts
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/measurement/types.go
  - frontend/src/features/datalink/workbench-v2/state/cutoverFlow.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingSaveInfo.tsx
  - internal/datalink/snapshot/identity.go
  - tests/shell/start-backend-logfile.sh
  - tests/shell/start-backend-cleanup-order.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CounterPolicyEditor.tsx
  - internal/datalink/groupdelivery/sample_codec.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - tests/powershell/start-build-single-failure.ps1
  - internal/datalink/schema/migrations/030_write_group_profile_keys_sqlite.up.sql
  - internal/datalink/workspace/write_group_runtime_layout.go
  - tests/shell/start-dev-bootstrap.sh
  - internal/datalink/recordingplan/schema_preview.go
  - tests/shell/start-backend-before-frontend.sh
  - internal/datalink/recordingplan/storage_restore.go
  - README.md
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileSchemaPreparation.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaPanelHelpers.ts
  - internal/api/handlers/studio_v2_workspace_write_group_profile_handler.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - docs/releases/studio-mapping-recovery-20261024.md
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - tests/powershell/start-embedded-bootstrap.ps1
  - internal/datalink/groupdelivery/status_quality.go
  - internal/datalink/recordingplan/storage_lineage.go
  - frontend/tsconfig.json
  - internal/datalink/dbtarget/managed_destination.go
  - tests/powershell/start-dev-mode-order.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicHeader.tsx
  - internal/datalink/grouppipeline/storage_recovery.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
  - internal/datalink/dbtarget/service.go
  - tests/shell/start-port-management.sh
  - internal/datalink/measurement/counter_confirmation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchEditor.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/recordingplan/schema_group_long.go
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/technical/evidence/circuit-recording-browser-20261024.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/schema/migrations/032_write_group_storage_scopes_sqlite.up.sql
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/index.css
  - docs/releases/frontend-toolchain-upgrade-20261024.md
  - start.ps1
  - internal/datalink/workspace/recording_start_steps.go
  - internal/datalink/workspace/write_group_circuit_profile.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/config.go
  - tests/powershell/start-environment-restore.ps1
  - internal/datalink/snapshot/assembler_counter.go
  - internal/datalink/workspace/write_group_schema_confirm.go
  - internal/datalink/groupdelivery/storage_recovery.go
  - docs/releases/evidence/mapping-layout-20261025/basic-1024.png
  - frontend/src/utils/recordingGroupLayoutJson.ts
  - internal/datalink/dbtarget/managed_timescale_inspection.go
  - docs/technical/evidence/circuit-recording-20261024/final-reconstruction-warning.png
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusList.tsx
  - internal/datalink/migrator_circuit_recording_profiles.go
  - internal/datalink/runtime/group_boundary_counter_source.go
  - internal/datalink/workspace/write_group_schema.go
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/schema/migrations/031_write_group_storage_mode_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4FocusTargets.ts
  - frontend/src/utils/backendErrorCodes.ts
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/recordingplan/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useCircuitProfileStart.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/utils/preciseInstant.ts
  - tests/shell/embedded-frontend-delivery.sh
  - internal/datalink/schema/migrations/033_write_group_counter_baselines_sqlite.up.sql
  - docs/technical/evidence/circuit-recording-20261024/final-batch-preview.png
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/group_boundary.go
  - scripts/start-port-utils.ps1
  - docs/technical/evidence/circuit-recording-20261024/final-prepare-only.png
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/candidates.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - docs/technical/studio-v2-circuit-recording.md
  - frontend/postcss.config.cjs
  - internal/datalink/recordingplan/storage_restore_exec.go
  - internal/datalink/dbtarget/postgres_column_inspection.go
  - internal/datalink/groupdelivery/store_closure.go
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - frontend/src/components/datalink/MemoryGrid.tsx
  - internal/api/router_studio_v2_write_group_routes.go
  - internal/datalink/grouppipeline/delivery_view.go
  - internal/datalink/runtime/group_boundary_counter.go
  - internal/datalink/schema/migrations/037_write_group_cutover_intent_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - internal/datalink/workspace/write_group_lifecycle.go
  - docs/swagger/swagger.json
  - internal/datalink/dbtarget/timescale.go
  - frontend/package.json
  - internal/datalink/measurement/counter_policy.go
  - internal/datalink/measurement/epoch.go
  - internal/datalink/recordingplan/schema_group_inspect.go
  - internal/datalink/dbtarget/group_row_layout.go
  - tests/powershell/start-port-ownership.ps1
  - internal/datalink/workspace/write_group_apply.go
  - tests/shell/start-owned-readiness.sh
  - internal/datalink/workspace/write_group_readiness.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - internal/datalink/schema/migrations/034_write_group_storage_lineage_sqlite.up.sql
  - docs/releases/remaining-wip-review-20261025.md
  - frontend/src/features/gateway/gatewayAdapter.ts
  - frontend/src/utils/studioV2RecordingProfile.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - tests/shell/start-frontend-deps.sh
  - internal/datalink/workspace/write_group_profile_interval.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingBulkApply.tsx
  - docs/technical/evidence/circuit-recording-scenario-coverage.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - internal/datalink/grouppipeline/recovery.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/setupTests.ts
  - internal/datalink/workspace/write_group_cutover_rollback.go
  - internal/datalink/snapshot/assembler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/releases/launcher-workflows-20261023.md
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - internal/datalink/dbtarget/group_long_row_columns.go
  - frontend/src/components/MonitorControl.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBlockers.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WorkspaceScopeNote.tsx
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchOutcome.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverPanel.tsx
  - internal/datalink/aggregation/profile_counter_types.go
  - internal/datalink/groupdelivery/ledger.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - internal/datalink/grouppipeline/storage_recovery_target.go
  - internal/datalink/migrator_write_group_delivery.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - internal/datalink/aggregation/profile_counter.go
  - scripts/start-process-utils.ps1
  - docs/swagger/docs.go
  - internal/api/router_logger.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/basicStartRequest.ts
  - tests/shell/start-dev-lifecycle.sh
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - docs/technical/studio-v2-circuit-recording-browser-verification.md
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - internal/datalink/recordingplan/repository_preview_token.go
  - internal/datalink/groupdelivery/store.go
  - internal/api/handlers/studio_v2_workspace_write_group_cutover_handler.go
  - docs/technical/evidence/circuit-recording-db-matrix-20261024.json
  - internal/datalink/groupdelivery/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileBatchExecutor.ts
  - internal/datalink/groupdelivery/store_baseline.go
  - docs/technical/studio-v2-write-groups.md
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/groupdelivery/sender.go
  - frontend/src/types/measurement.ts
  - internal/datalink/schema/migrations/036_write_group_cutovers_sqlite.up.sql
  - internal/datalink/aggregation/profile_counter_checkpoint.go
  - frontend/src/types/recordingPlan.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/recordingplan/repository_storage_repair.go
  - internal/datalink/recordingplan/schema_group.go
  - tests/shell/start-frontend-install-failure.sh
  - tests/powershell/start-browser-prompt.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - internal/datalink/workspace/write_group_long_storage.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverRun.tsx
  - docs/technical/studio-v2-circuit-recording-db-verification.md
  - docs/technical/evidence/circuit-recording-20261024/final-cutover-history.png
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/datalink/grouppipeline/delivery_sources.go
  - frontend/src/features/datalink/workbench-v2/state/circuitCounterConfirmation.ts
  - internal/datalink/dbtarget/writer_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicActions.tsx
  - tests/powershell/start-owned-readiness.ps1
  - AGENTS.md
  - internal/datalink/dbtarget/group_long_row.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useCircuitProfileBatchRun.ts
  - internal/datalink/dbtarget/storage_restore_target.go
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - tests/shell/start-frontend-readiness.sh
  - docs/technical/studio-v2-circuit-recording-final-review.md
  - internal/datalink/snapshot/outcome.go
  - scripts/start-process-utils.sh
  - internal/datalink/runtime/group_boundary_long.go
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useCircuitProfileBatchHost.ts
  - frontend/src/types/studioV2WriteGroupCutover.ts
  - internal/datalink/workspace/write_group_readiness_evidence.go
  - tests/shell/start-browser-prompt.sh
  - internal/datalink/workspace/write_group_cutover_plan.go
  - internal/datalink/dbtarget/types.go
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/groupdelivery/source_evidence.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - frontend/src/utils/studioV2WriteGroupDeliveryJson.ts
  - internal/datalink/dbtarget/managed_index_inspection.go
  - frontend/src/utils/studioV2WriteGroupStorageJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - scripts/start-log-utils.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/aggregation/profile_counter_events.go
  - internal/datalink/aggregation/profile_counter_walk.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/schema/migrations/035_write_group_storage_events_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileRoles.tsx
  - internal/datalink/workspace/write_group_readiness_eval.go
  - frontend/src/utils/studioV2WriteGroupCutoverJson.ts
  - tests/powershell/start-process-log-stream.ps1
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_repository.go
  - internal/datalink/workspace/write_group_create_once.go
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/snapshot/assembler_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/runtime/group_boundary_sample.go
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/hooks/datalink/useStudioV2WriteGroupCutover.ts
  - internal/datalink/workspace/write_group_repository_validation.go
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverOutcomeView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/LegacyMigrationEntry.tsx
  - internal/datalink/recordingplan/schema_group_extend.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchPreview.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_schema_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverReviewView.tsx
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/measurement/exact_arith.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBlockers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useReportOperationActive.ts
  - internal/datalink/grouppipeline/pipeline.go
  - internal/datalink/grouppipeline/recovery_long.go
  - docs/releases/evidence/mapping-layout-20261025/advanced-1024.png
  - internal/datalink/workspace/write_group_counter_readiness.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/state/cutoverMessages.ts
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitCounterConfirmation.tsx
  - internal/datalink/dbtarget/table_inspection.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - internal/datalink/workspace/write_group_cutover_repository.go
  - tests/shell/start-menu-args-root.sh
  - frontend/src/features/datalink/workbench-v2/state/basicRecordingMessages.ts
  - internal/datalink/workspace/write_group_cutover_apply.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - docs/swagger/swagger.yaml
  - frontend/src/utils/studioV2WriteGroupJsonPrimitives.ts
  - internal/datalink/schema/migrations/029_write_group_circuit_profiles_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatch.ts
  - frontend/src/features/datalink/workbench-v2/state/counterPolicyForm.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/useCutoverFlow.ts
  - frontend/src/features/gateway/gatewayDraftStore.ts
  - scripts/build-single.mjs
  - .air.toml
  - docs/technical/evidence/circuit-recording-test-results-20261024.json
  - internal/datalink/workspace/write_group_cutover.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2WriteGroup.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/runtime/group_boundary_build.go
tests:
  - internal/datalink/recordingplan/schema_preview_recreate_test.go
  - internal/api/router_studio_v2_circuit_profile_timescale_start_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-dual.test.tsx
  - internal/datalink/recordingplan/storage_restore_failure_test.go
  - internal/datalink/dbtarget/group_long_row_counter_events_postgres_test.go
  - frontend/tests/e2e/fixtures/mapping-live-overflow.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - frontend/tests/fixtures/cutoverFixtures.ts
  - internal/datalink/grouppipeline/recovery_long_timescale_test.go
  - frontend/tests/unit/workbench-v2/mapping-save-info.test.tsx
  - internal/datalink/recordingplan/schema_group_timescale_test.go
  - internal/api/router_studio_v2_circuit_profile_start_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-counter-confirmation.test.tsx
  - internal/datalink/runtime/group_boundary_counter_source_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/api/router_studio_v2_write_group_recovery_frozen_test.go
  - internal/datalink/runtime/group_boundary_long_durable_test.go
  - frontend/tests/unit/utils/viteConfig.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recreate_policy_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/api/router_studio_v2_write_group_cutover_guard_test.go
  - internal/api/router_studio_v2_write_group_schema_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/runtime/group_boundary_counter_durable_test.go
  - internal/datalink/runtime/group_boundary_long_gates_test.go
  - frontend/tests/unit/workbench-v2/mapping-save-error-ui.test.tsx
  - internal/datalink/dbtarget/group_row_layout_truncated_test.go
  - internal/datalink/groupdelivery/source_evidence_test.go
  - scripts/tests/frontend-styles.test.mjs
  - frontend/tests/unit/workbench-v2/circuit-profile-identity.test.ts
  - internal/datalink/dbtarget/group_long_row_counter_postgres_test.go
  - internal/api/router_studio_v2_write_group_recovery_cumulative_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - internal/datalink/groupdelivery/store_baseline_test.go
  - internal/datalink/groupdelivery/sender_storage_recovery_test.go
  - internal/api/router_studio_v2_write_group_recovery_denied_test.go
  - internal/api/router_studio_v2_write_group_recovery_observe_test.go
  - frontend/tests/unit/workbench-v2/recreate-missing-contract.test.ts
  - frontend/tests/unit/types/measurement-counter-policy.test.ts
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_owner_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-profiles.test.ts
  - internal/datalink/aggregation/profile_counter_unproven_test.go
  - internal/api/circuit_recording_integration_test.go
  - internal/api/router_studio_v2_write_group_cutover_observe_test.go
  - internal/datalink/groupdelivery/sender_storage_recovery_settlement_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-contract.test.ts
  - frontend/tests/e2e/fixtures/mapping-live-overflow.tsx
  - internal/api/router_studio_v2_write_group_long_rollover_test.go
  - internal/datalink/dbtarget/group_long_row_roundtrip_test.go
  - internal/datalink/groupdelivery/sample_codec_received_test.go
  - frontend/tests/unit/workbench-v2/step3-source-display-names.test.tsx
  - internal/api/router_studio_v2_circuit_profile_types_test.go
  - internal/datalink/dbtarget/postgres_numeric_inspection_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-blockers.test.tsx
  - internal/api/router_studio_v2_write_group_recovery_rig_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - internal/datalink/dbtarget/group_long_row_test.go
  - cmd/test_ui/service_wiring.go
  - scripts/tests/launcher-live-smoke.py
  - frontend/tests/fixtures/circuitBatchPanel.ts
  - internal/api/router_studio_v2_write_groups_test.go
  - internal/datalink/dbtarget/circuit_recording_integration_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover-locks.test.tsx
  - internal/datalink/dbtarget/managed_index_inspection_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-browser-wiring.test.tsx
  - internal/datalink/workspace/write_group_legacy_decode_test.go
  - internal/api/router_studio_v2_write_group_timescale_role_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/api/router_studio_v2_write_group_long_test.go
  - internal/datalink/dbtarget/timescale_test.go
  - internal/api/router_studio_v2_write_group_timescale_drift_test.go
  - internal/datalink/recordingplan/storage_restore_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-queue.missing.test.tsx
  - internal/api/router_studio_v2_write_group_timescale_guard_test.go
  - frontend/tests/unit/workbench-v2/long-layout-contract.test.ts
  - frontend/tests/e2e/fixtures/mapping-scroll.tsx
  - internal/datalink/groupdelivery/source_evidence_review_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/e2e/mapping-scroll.spec.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - internal/datalink/measurement/exact_arith_test.go
  - internal/api/router_studio_v2_write_group_cutover_draft_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-policy-binding.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-recovery.test.tsx
  - frontend/tests/e2e/studio-mapping-scroll.spec.ts
  - internal/api/router_studio_v2_write_group_source_evidence_cutover_test.go
  - internal/api/router_studio_v2_circuit_profile_counter_http_test.go
  - internal/datalink/recordingplan/schema_group_long_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-outcome.test.ts
  - internal/datalink/runtime/group_boundary_counter_discontinuity_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - internal/datalink/recordingplan/schema_group_long_catalog_test.go
  - frontend/tests/e2e/fixtures/studio-mapping-scroll.tsx
  - internal/api/router_studio_v2_write_group_timescale_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/api/router_studio_v2_write_group_recovery_denied_pg_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-review.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/aggregation/profile_counter_gap_test.go
  - internal/datalink/grouppipeline/circuit_recording_gap_integration_test.go
  - internal/datalink/grouppipeline/storage_recovery_wiring_test.go
  - frontend/tests/e2e/fixtures/mapping-scroll.html
  - internal/api/router_studio_v2_write_group_cutover_review_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-storage-notice.test.tsx
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-batch.test.tsx
  - frontend/tests/e2e/mapping-live-overflow.spec.ts
  - scripts/tests/build-single.test.mjs
  - internal/datalink/grouppipeline/long_timescale_test.go
  - internal/api/router_studio_v2_profile_legacy_coexistence_rig_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-guards.test.tsx
  - internal/datalink/grouppipeline/long_build_test.go
  - internal/datalink/runtime/group_boundary_long_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/measurement/counter_policy_test.go
  - internal/datalink/runtime/group_boundary_long_fixture_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - internal/datalink/runtime/group_boundary_counter_test.go
  - frontend/tests/unit/workbench-v2/database-connector-save-race.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-controller.test.tsx
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - frontend/tests/e2e/fixtures/studio-mapping-scroll.html
  - frontend/tests/unit/utils/preciseInstant.test.ts
  - internal/api/router_studio_v2_write_group_profile_acquisition_test.go
  - internal/api/router_studio_v2_write_group_long_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/aggregation/profile_counter_test.go
  - internal/datalink/recordingplan/schema_group_long_postgres_test.go
  - internal/api/router_studio_v2_write_group_long_recovery_test.go
  - internal/datalink/workspace/write_group_long_storage_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/dbtarget/timescale_interval_test.go
  - internal/datalink/recordingplan/storage_lineage_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover-storage.test.ts
  - internal/datalink/aggregation/profile_counter_checkpoint_test.go
  - internal/api/router_studio_v2_write_group_recovery_test.go
  - internal/datalink/workspace/write_group_long_runtime_gate_test.go
  - frontend/tests/e2e/mapping-save-info.spec.ts
  - frontend/tests/unit/workbench-v2/cutover-contract.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
  - internal/datalink/grouppipeline/circuit_recording_integration_test.go
  - internal/datalink/dbtarget/group_long_row_counter_events_test.go
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/workspace/write_group_circuit_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-policy-review.test.tsx
  - internal/datalink/snapshot/profile_counter_test.go
  - internal/api/router_studio_v2_circuit_profile_partial_test.go
  - internal/datalink/dbtarget/group_long_row_counter_events_view_test.go
  - internal/datalink/aggregation/profile_counter_discontinuity_test.go
  - internal/datalink/dbtarget/group_long_row_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-policy.test.tsx
  - internal/api/router_studio_v2_write_group_cutover_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - internal/api/router_studio_v2_write_group_source_evidence_test.go
  - internal/datalink/snapshot/profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-proof.test.ts
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - internal/api/router_studio_v2_profile_legacy_coexistence_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - internal/datalink/runtime/group_boundary_cutover_test.go
  - internal/datalink/grouppipeline/long_recovery_test.go
  - internal/datalink/grouppipeline/long_fixture_test.go
  - internal/api/router_studio_v2_workspace_measurements_policy_test.go
  - internal/datalink/workspace/write_group_timescale_test.go
  - internal/datalink/dbtarget/storage_restore_target_test.go
  - internal/datalink/groupdelivery/source_evidence_consistency_test.go
-->

---
### Requirement: Explain advanced groups and workspace responsibilities

Basic SHALL identify Advanced as the editor for canonical source mapping, custom layouts and confirmed signed/rollover policies, and Workspace as the scope for readiness and activation. It MUST NOT describe them as alternative data-table formats or competing recording settings. A blocked recording SHALL identify the affected circuit, source or destination and offer its repair action. Identical blocker causes SHALL not be repeated as independent failures, and collecting or local acceptance SHALL not be shown as SQL commit evidence.

#### Scenario: Advanced and workspace review
- **WHEN** an operator opens Advanced or workspace confirmation after selecting circuit profiles
- **THEN** Advanced shows the same saved profile settings, workspace shows scoped readiness and activation, and neither silently changes intervals or storage.

#### Scenario: Database preparation failed
- **WHEN** a selected destination lacks schema confirmation or Timescale capability
- **THEN** the operator sees one actionable explanation for that cause with a destination repair entry rather than an unexplained disabled start button or raw driver error.

<!-- @trace
source: add-circuit-recording-profiles
updated: 2026-10-25
code:
  - frontend/src/utils/recordingPlanJson.ts
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/measurement/types.go
  - frontend/src/features/datalink/workbench-v2/state/cutoverFlow.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingSaveInfo.tsx
  - internal/datalink/snapshot/identity.go
  - tests/shell/start-backend-logfile.sh
  - tests/shell/start-backend-cleanup-order.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CounterPolicyEditor.tsx
  - internal/datalink/groupdelivery/sample_codec.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - tests/powershell/start-build-single-failure.ps1
  - internal/datalink/schema/migrations/030_write_group_profile_keys_sqlite.up.sql
  - internal/datalink/workspace/write_group_runtime_layout.go
  - tests/shell/start-dev-bootstrap.sh
  - internal/datalink/recordingplan/schema_preview.go
  - tests/shell/start-backend-before-frontend.sh
  - internal/datalink/recordingplan/storage_restore.go
  - README.md
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileSchemaPreparation.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaPanelHelpers.ts
  - internal/api/handlers/studio_v2_workspace_write_group_profile_handler.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - docs/releases/studio-mapping-recovery-20261024.md
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - tests/powershell/start-embedded-bootstrap.ps1
  - internal/datalink/groupdelivery/status_quality.go
  - internal/datalink/recordingplan/storage_lineage.go
  - frontend/tsconfig.json
  - internal/datalink/dbtarget/managed_destination.go
  - tests/powershell/start-dev-mode-order.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicHeader.tsx
  - internal/datalink/grouppipeline/storage_recovery.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
  - internal/datalink/dbtarget/service.go
  - tests/shell/start-port-management.sh
  - internal/datalink/measurement/counter_confirmation.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchEditor.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/recordingplan/schema_group_long.go
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - docs/technical/evidence/circuit-recording-browser-20261024.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/schema/migrations/032_write_group_storage_scopes_sqlite.up.sql
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - frontend/src/index.css
  - docs/releases/frontend-toolchain-upgrade-20261024.md
  - start.ps1
  - internal/datalink/workspace/recording_start_steps.go
  - internal/datalink/workspace/write_group_circuit_profile.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/config.go
  - tests/powershell/start-environment-restore.ps1
  - internal/datalink/snapshot/assembler_counter.go
  - internal/datalink/workspace/write_group_schema_confirm.go
  - internal/datalink/groupdelivery/storage_recovery.go
  - docs/releases/evidence/mapping-layout-20261025/basic-1024.png
  - frontend/src/utils/recordingGroupLayoutJson.ts
  - internal/datalink/dbtarget/managed_timescale_inspection.go
  - docs/technical/evidence/circuit-recording-20261024/final-reconstruction-warning.png
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusList.tsx
  - internal/datalink/migrator_circuit_recording_profiles.go
  - internal/datalink/runtime/group_boundary_counter_source.go
  - internal/datalink/workspace/write_group_schema.go
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/schema/migrations/031_write_group_storage_mode_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/step4FocusTargets.ts
  - frontend/src/utils/backendErrorCodes.ts
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/recordingplan/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useCircuitProfileStart.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/utils/preciseInstant.ts
  - tests/shell/embedded-frontend-delivery.sh
  - internal/datalink/schema/migrations/033_write_group_counter_baselines_sqlite.up.sql
  - docs/technical/evidence/circuit-recording-20261024/final-batch-preview.png
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/runtime/group_boundary.go
  - scripts/start-port-utils.ps1
  - docs/technical/evidence/circuit-recording-20261024/final-prepare-only.png
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/candidates.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - docs/technical/studio-v2-circuit-recording.md
  - frontend/postcss.config.cjs
  - internal/datalink/recordingplan/storage_restore_exec.go
  - internal/datalink/dbtarget/postgres_column_inspection.go
  - internal/datalink/groupdelivery/store_closure.go
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - frontend/src/components/datalink/MemoryGrid.tsx
  - internal/api/router_studio_v2_write_group_routes.go
  - internal/datalink/grouppipeline/delivery_view.go
  - internal/datalink/runtime/group_boundary_counter.go
  - internal/datalink/schema/migrations/037_write_group_cutover_intent_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - internal/datalink/workspace/write_group_lifecycle.go
  - docs/swagger/swagger.json
  - internal/datalink/dbtarget/timescale.go
  - frontend/package.json
  - internal/datalink/measurement/counter_policy.go
  - internal/datalink/measurement/epoch.go
  - internal/datalink/recordingplan/schema_group_inspect.go
  - internal/datalink/dbtarget/group_row_layout.go
  - tests/powershell/start-port-ownership.ps1
  - internal/datalink/workspace/write_group_apply.go
  - tests/shell/start-owned-readiness.sh
  - internal/datalink/workspace/write_group_readiness.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - internal/datalink/schema/migrations/034_write_group_storage_lineage_sqlite.up.sql
  - docs/releases/remaining-wip-review-20261025.md
  - frontend/src/features/gateway/gatewayAdapter.ts
  - frontend/src/utils/studioV2RecordingProfile.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/LocalModbusReviewSurface.tsx
  - tests/shell/start-frontend-deps.sh
  - internal/datalink/workspace/write_group_profile_interval.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingBulkApply.tsx
  - docs/technical/evidence/circuit-recording-scenario-coverage.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - internal/datalink/grouppipeline/recovery.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/setupTests.ts
  - internal/datalink/workspace/write_group_cutover_rollback.go
  - internal/datalink/snapshot/assembler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/releases/launcher-workflows-20261023.md
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - internal/datalink/dbtarget/group_long_row_columns.go
  - frontend/src/components/MonitorControl.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBlockers.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WorkspaceScopeNote.tsx
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchOutcome.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverPanel.tsx
  - internal/datalink/aggregation/profile_counter_types.go
  - internal/datalink/groupdelivery/ledger.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - internal/datalink/grouppipeline/storage_recovery_target.go
  - internal/datalink/migrator_write_group_delivery.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - internal/datalink/aggregation/profile_counter.go
  - scripts/start-process-utils.ps1
  - docs/swagger/docs.go
  - internal/api/router_logger.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/basicStartRequest.ts
  - tests/shell/start-dev-lifecycle.sh
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - docs/technical/studio-v2-circuit-recording-browser-verification.md
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - internal/datalink/recordingplan/repository_preview_token.go
  - internal/datalink/groupdelivery/store.go
  - internal/api/handlers/studio_v2_workspace_write_group_cutover_handler.go
  - docs/technical/evidence/circuit-recording-db-matrix-20261024.json
  - internal/datalink/groupdelivery/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileBatchExecutor.ts
  - internal/datalink/groupdelivery/store_baseline.go
  - docs/technical/studio-v2-write-groups.md
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/groupdelivery/sender.go
  - frontend/src/types/measurement.ts
  - internal/datalink/schema/migrations/036_write_group_cutovers_sqlite.up.sql
  - internal/datalink/aggregation/profile_counter_checkpoint.go
  - frontend/src/types/recordingPlan.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/recordingplan/repository_storage_repair.go
  - internal/datalink/recordingplan/schema_group.go
  - tests/shell/start-frontend-install-failure.sh
  - tests/powershell/start-browser-prompt.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields.tsx
  - internal/datalink/workspace/write_group_long_storage.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverRun.tsx
  - docs/technical/studio-v2-circuit-recording-db-verification.md
  - docs/technical/evidence/circuit-recording-20261024/final-cutover-history.png
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/datalink/grouppipeline/delivery_sources.go
  - frontend/src/features/datalink/workbench-v2/state/circuitCounterConfirmation.ts
  - internal/datalink/dbtarget/writer_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicActions.tsx
  - tests/powershell/start-owned-readiness.ps1
  - AGENTS.md
  - internal/datalink/dbtarget/group_long_row.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/useCircuitProfileBatchRun.ts
  - internal/datalink/dbtarget/storage_restore_target.go
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - tests/shell/start-frontend-readiness.sh
  - docs/technical/studio-v2-circuit-recording-final-review.md
  - internal/datalink/snapshot/outcome.go
  - scripts/start-process-utils.sh
  - internal/datalink/runtime/group_boundary_long.go
  - frontend/vite.config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useCircuitProfileBatchHost.ts
  - frontend/src/types/studioV2WriteGroupCutover.ts
  - internal/datalink/workspace/write_group_readiness_evidence.go
  - tests/shell/start-browser-prompt.sh
  - internal/datalink/workspace/write_group_cutover_plan.go
  - internal/datalink/dbtarget/types.go
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/groupdelivery/source_evidence.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - frontend/src/utils/studioV2WriteGroupDeliveryJson.ts
  - internal/datalink/dbtarget/managed_index_inspection.go
  - frontend/src/utils/studioV2WriteGroupStorageJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel.tsx
  - scripts/start-log-utils.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/aggregation/profile_counter_events.go
  - internal/datalink/aggregation/profile_counter_walk.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/schema/migrations/035_write_group_storage_events_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileRoles.tsx
  - internal/datalink/workspace/write_group_readiness_eval.go
  - frontend/src/utils/studioV2WriteGroupCutoverJson.ts
  - tests/powershell/start-process-log-stream.ps1
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_repository.go
  - internal/datalink/workspace/write_group_create_once.go
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/snapshot/assembler_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - internal/datalink/runtime/group_boundary_sample.go
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/hooks/datalink/useStudioV2WriteGroupCutover.ts
  - internal/datalink/workspace/write_group_repository_validation.go
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverOutcomeView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/LegacyMigrationEntry.tsx
  - internal/datalink/recordingplan/schema_group_extend.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchPreview.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_schema_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverReviewView.tsx
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/datalink/measurement/exact_arith.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBlockers.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/useReportOperationActive.ts
  - internal/datalink/grouppipeline/pipeline.go
  - internal/datalink/grouppipeline/recovery_long.go
  - docs/releases/evidence/mapping-layout-20261025/advanced-1024.png
  - internal/datalink/workspace/write_group_counter_readiness.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/state/cutoverMessages.ts
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitCounterConfirmation.tsx
  - internal/datalink/dbtarget/table_inspection.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - internal/datalink/workspace/write_group_cutover_repository.go
  - tests/shell/start-menu-args-root.sh
  - frontend/src/features/datalink/workbench-v2/state/basicRecordingMessages.ts
  - internal/datalink/workspace/write_group_cutover_apply.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - docs/swagger/swagger.yaml
  - frontend/src/utils/studioV2WriteGroupJsonPrimitives.ts
  - internal/datalink/schema/migrations/029_write_group_circuit_profiles_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatch.ts
  - frontend/src/features/datalink/workbench-v2/state/counterPolicyForm.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/useCutoverFlow.ts
  - frontend/src/features/gateway/gatewayDraftStore.ts
  - scripts/build-single.mjs
  - .air.toml
  - docs/technical/evidence/circuit-recording-test-results-20261024.json
  - internal/datalink/workspace/write_group_cutover.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2WriteGroup.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/runtime/group_boundary_build.go
tests:
  - internal/datalink/recordingplan/schema_preview_recreate_test.go
  - internal/api/router_studio_v2_circuit_profile_timescale_start_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-dual.test.tsx
  - internal/datalink/recordingplan/storage_restore_failure_test.go
  - internal/datalink/dbtarget/group_long_row_counter_events_postgres_test.go
  - frontend/tests/e2e/fixtures/mapping-live-overflow.html
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - frontend/tests/fixtures/cutoverFixtures.ts
  - internal/datalink/grouppipeline/recovery_long_timescale_test.go
  - frontend/tests/unit/workbench-v2/mapping-save-info.test.tsx
  - internal/datalink/recordingplan/schema_group_timescale_test.go
  - internal/api/router_studio_v2_circuit_profile_start_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-counter-confirmation.test.tsx
  - internal/datalink/runtime/group_boundary_counter_source_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/api/router_studio_v2_write_group_recovery_frozen_test.go
  - internal/datalink/runtime/group_boundary_long_durable_test.go
  - frontend/tests/unit/utils/viteConfig.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recreate_policy_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/api/router_studio_v2_write_group_cutover_guard_test.go
  - internal/api/router_studio_v2_write_group_schema_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/runtime/group_boundary_counter_durable_test.go
  - internal/datalink/runtime/group_boundary_long_gates_test.go
  - frontend/tests/unit/workbench-v2/mapping-save-error-ui.test.tsx
  - internal/datalink/dbtarget/group_row_layout_truncated_test.go
  - internal/datalink/groupdelivery/source_evidence_test.go
  - scripts/tests/frontend-styles.test.mjs
  - frontend/tests/unit/workbench-v2/circuit-profile-identity.test.ts
  - internal/datalink/dbtarget/group_long_row_counter_postgres_test.go
  - internal/api/router_studio_v2_write_group_recovery_cumulative_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - internal/datalink/groupdelivery/store_baseline_test.go
  - internal/datalink/groupdelivery/sender_storage_recovery_test.go
  - internal/api/router_studio_v2_write_group_recovery_denied_test.go
  - internal/api/router_studio_v2_write_group_recovery_observe_test.go
  - frontend/tests/unit/workbench-v2/recreate-missing-contract.test.ts
  - frontend/tests/unit/types/measurement-counter-policy.test.ts
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_owner_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-profiles.test.ts
  - internal/datalink/aggregation/profile_counter_unproven_test.go
  - internal/api/circuit_recording_integration_test.go
  - internal/api/router_studio_v2_write_group_cutover_observe_test.go
  - internal/datalink/groupdelivery/sender_storage_recovery_settlement_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-contract.test.ts
  - frontend/tests/e2e/fixtures/mapping-live-overflow.tsx
  - internal/api/router_studio_v2_write_group_long_rollover_test.go
  - internal/datalink/dbtarget/group_long_row_roundtrip_test.go
  - internal/datalink/groupdelivery/sample_codec_received_test.go
  - frontend/tests/unit/workbench-v2/step3-source-display-names.test.tsx
  - internal/api/router_studio_v2_circuit_profile_types_test.go
  - internal/datalink/dbtarget/postgres_numeric_inspection_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-blockers.test.tsx
  - internal/api/router_studio_v2_write_group_recovery_rig_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - internal/datalink/dbtarget/group_long_row_test.go
  - cmd/test_ui/service_wiring.go
  - scripts/tests/launcher-live-smoke.py
  - frontend/tests/fixtures/circuitBatchPanel.ts
  - internal/api/router_studio_v2_write_groups_test.go
  - internal/datalink/dbtarget/circuit_recording_integration_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover-locks.test.tsx
  - internal/datalink/dbtarget/managed_index_inspection_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-browser-wiring.test.tsx
  - internal/datalink/workspace/write_group_legacy_decode_test.go
  - internal/api/router_studio_v2_write_group_timescale_role_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/api/router_studio_v2_write_group_long_test.go
  - internal/datalink/dbtarget/timescale_test.go
  - internal/api/router_studio_v2_write_group_timescale_drift_test.go
  - internal/datalink/recordingplan/storage_restore_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-queue.missing.test.tsx
  - internal/api/router_studio_v2_write_group_timescale_guard_test.go
  - frontend/tests/unit/workbench-v2/long-layout-contract.test.ts
  - frontend/tests/e2e/fixtures/mapping-scroll.tsx
  - internal/datalink/groupdelivery/source_evidence_review_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/e2e/mapping-scroll.spec.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - internal/datalink/measurement/exact_arith_test.go
  - internal/api/router_studio_v2_write_group_cutover_draft_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-policy-binding.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-recovery.test.tsx
  - frontend/tests/e2e/studio-mapping-scroll.spec.ts
  - internal/api/router_studio_v2_write_group_source_evidence_cutover_test.go
  - internal/api/router_studio_v2_circuit_profile_counter_http_test.go
  - internal/datalink/recordingplan/schema_group_long_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-outcome.test.ts
  - internal/datalink/runtime/group_boundary_counter_discontinuity_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - internal/datalink/recordingplan/schema_group_long_catalog_test.go
  - frontend/tests/e2e/fixtures/studio-mapping-scroll.tsx
  - internal/api/router_studio_v2_write_group_timescale_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/api/router_studio_v2_write_group_recovery_denied_pg_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-review.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - internal/datalink/aggregation/profile_counter_gap_test.go
  - internal/datalink/grouppipeline/circuit_recording_gap_integration_test.go
  - internal/datalink/grouppipeline/storage_recovery_wiring_test.go
  - frontend/tests/e2e/fixtures/mapping-scroll.html
  - internal/api/router_studio_v2_write_group_cutover_review_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-storage-notice.test.tsx
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-batch.test.tsx
  - frontend/tests/e2e/mapping-live-overflow.spec.ts
  - scripts/tests/build-single.test.mjs
  - internal/datalink/grouppipeline/long_timescale_test.go
  - internal/api/router_studio_v2_profile_legacy_coexistence_rig_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-guards.test.tsx
  - internal/datalink/grouppipeline/long_build_test.go
  - internal/datalink/runtime/group_boundary_long_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/measurement/counter_policy_test.go
  - internal/datalink/runtime/group_boundary_long_fixture_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - internal/datalink/runtime/group_boundary_counter_test.go
  - frontend/tests/unit/workbench-v2/database-connector-save-race.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-controller.test.tsx
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - frontend/tests/e2e/fixtures/studio-mapping-scroll.html
  - frontend/tests/unit/utils/preciseInstant.test.ts
  - internal/api/router_studio_v2_write_group_profile_acquisition_test.go
  - internal/api/router_studio_v2_write_group_long_runtime_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - internal/datalink/aggregation/profile_counter_test.go
  - internal/datalink/recordingplan/schema_group_long_postgres_test.go
  - internal/api/router_studio_v2_write_group_long_recovery_test.go
  - internal/datalink/workspace/write_group_long_storage_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/dbtarget/timescale_interval_test.go
  - internal/datalink/recordingplan/storage_lineage_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover-storage.test.ts
  - internal/datalink/aggregation/profile_counter_checkpoint_test.go
  - internal/api/router_studio_v2_write_group_recovery_test.go
  - internal/datalink/workspace/write_group_long_runtime_gate_test.go
  - frontend/tests/e2e/mapping-save-info.spec.ts
  - frontend/tests/unit/workbench-v2/cutover-contract.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
  - internal/datalink/grouppipeline/circuit_recording_integration_test.go
  - internal/datalink/dbtarget/group_long_row_counter_events_test.go
  - frontend/tests/unit/workbench-v2/workbench-local-modbus-review-surface.test.tsx
  - internal/datalink/workspace/write_group_circuit_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-policy-review.test.tsx
  - internal/datalink/snapshot/profile_counter_test.go
  - internal/api/router_studio_v2_circuit_profile_partial_test.go
  - internal/datalink/dbtarget/group_long_row_counter_events_view_test.go
  - internal/datalink/aggregation/profile_counter_discontinuity_test.go
  - internal/datalink/dbtarget/group_long_row_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-policy.test.tsx
  - internal/api/router_studio_v2_write_group_cutover_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - internal/api/router_studio_v2_write_group_source_evidence_test.go
  - internal/datalink/snapshot/profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-proof.test.ts
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - internal/api/router_studio_v2_profile_legacy_coexistence_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - internal/datalink/runtime/group_boundary_cutover_test.go
  - internal/datalink/grouppipeline/long_recovery_test.go
  - internal/datalink/grouppipeline/long_fixture_test.go
  - internal/api/router_studio_v2_workspace_measurements_policy_test.go
  - internal/datalink/workspace/write_group_timescale_test.go
  - internal/datalink/dbtarget/storage_restore_target_test.go
  - internal/datalink/groupdelivery/source_evidence_consistency_test.go
-->

---
### Requirement: Device tiles change only the viewed scope

Step 4 SHALL provide named device tiles with one visibly selected persisted device ID, keyboard focus and native activation. It SHALL NOT use a device dropdown. Tile selection SHALL filter the current device's recording configuration, sources, preparation and status without creating, applying, disabling, activating or deleting resources. Drafts SHALL remain scoped by workspace/device/purpose/circuit; late replies and unmounted requests SHALL not update another scope. A legacy wide group on device A SHALL not force device B onto A's legacy workflow. Cross-device custom groups SHALL retain their existing advanced identity and SHALL not be automatically split.

#### Scenario: Same address and late reply
- **WHEN** A and B share the same address/display name, the operator edits A, switches to B and A's delayed save/read reply arrives
- **THEN** B retains its own draft and status, A retains its pending work under its own IDs, and switching back restores A without creating a writer or changing acquisition.

#### Scenario: Keyboard and reload
- **WHEN** the operator selects a device tile using keyboard and reloads or navigates Back/Next
- **THEN** the selected valid device and its scoped drafts are restored without a dropdown, false saved state or automatic Apply.


<!-- @trace
source: unify-device-recording-workflow
updated: 2026-10-25
code:
  - internal/datalink/workspace/write_group_long_recording_shape.go
  - internal/datalink/workspace/write_group_member_metadata.go
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/features/datalink/workbench-v2/state/basicLegacyDraft.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/technical/studio-v2-circuit-recording.md
  - internal/datalink/grouppipeline/delivery_sources.go
  - internal/datalink/groupdelivery/source_evidence_wide.go
  - internal/datalink/workspace/write_group_cutover_repository.go
  - internal/datalink/measurement/mean.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - internal/datalink/dbtarget/group_long_row.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/workspace/write_group_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitPurposeFields.tsx
  - internal/datalink/workspace/write_group_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_wide_profile_shape.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitRequiredPolicy.tsx
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/swagger/docs.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/workspace/write_group_repository_validation.go
  - frontend/src/features/datalink/workbench-v2/state/deviceScopeSelection.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/runtime/group_boundary_wide_profile.go
  - internal/datalink/groupdelivery/source_evidence_aggregation.go
  - internal/datalink/snapshot/config.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - internal/datalink/runtime/group_boundary_build.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/groupdelivery/source_evidence.go
  - docs/swagger/swagger.json
  - internal/datalink/grouppipeline/recovery_long.go
  - internal/datalink/workspace/write_group_readiness.go
  - internal/datalink/snapshot/assembler_profile.go
  - internal/datalink/snapshot/outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SourceReview.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/types/studioV2WriteGroup.ts
  - internal/datalink/dbtarget/group_row_layout.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/runtime/group_boundary.go
  - internal/datalink/schema/migrations/038_write_group_recording_reduction_sqlite.up.sql
  - internal/datalink/snapshot/assembler_counter.go
  - frontend/src/utils/studioV2RecordingProfile.ts
  - internal/datalink/snapshot/assembler_mean.go
  - internal/datalink/snapshot/mean_evidence.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - internal/datalink/workspace/write_group_counter_readiness.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/workspace/write_group_profile_interval.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - internal/datalink/grouppipeline/recovery.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - internal/datalink/snapshot/assembler_wide_profile.go
  - internal/datalink/runtime/group_boundary_long.go
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/assembler.go
  - internal/datalink/workspace/write_group_cutover_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeviceScopeTiles.tsx
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
tests:
  - internal/datalink/grouppipeline/mean_result_type_proof_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - internal/api/router_studio_v2_reviewed_wide_mean_transition_test.go
  - internal/datalink/snapshot/mean_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/api/router_studio_v2_wide_profile_guard_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - internal/api/router_studio_v2_recording_start_scope_test.go
  - internal/api/router_studio_v2_wide_profile_runtime_test.go
  - frontend/tests/unit/workbench-v2/write-group-aggregation-evidence.test.ts
  - internal/datalink/dbtarget/group_long_row_mean_test.go
  - internal/datalink/dbtarget/group_row_layout_wide_profile_test.go
  - internal/datalink/grouppipeline/average_recovery_test.go
  - internal/datalink/grouppipeline/recovery_reduction_test.go
  - internal/datalink/runtime/group_boundary_wide_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - frontend/tests/unit/workbench-v2/step4-evidence-boundaries.test.tsx
  - internal/datalink/grouppipeline/wide_profile_result_type_proof_test.go
  - internal/datalink/workspace/write_group_member_recording_shapes_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-tiles.test.tsx
  - frontend/tests/unit/workbench-v2/wide-profile-member-metadata.test.ts
  - frontend/tests/e2e/fixtures/device-recording-workflow.ts
  - frontend/tests/unit/services/studioV2WriteGroupDeliveryEnvelope.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose.test.tsx
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/runtime/group_boundary_long_reduction_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - internal/datalink/runtime/group_boundary_long_counter_raw_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - frontend/tests/e2e/device-recording-workflow.spec.ts
  - internal/api/router_studio_v2_recording_workflow_coexistence_test.go
  - internal/datalink/snapshot/assembler_wide_cumulative_test.go
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_legacy_provenance_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/basic-recording-panel-regressions.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/datalink/grouppipeline/average_fixture_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/datalink/grouppipeline/average_delivery_evidence_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/snapshot/assembler_wide_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/datalink/workspace/write_group_member_metadata_test.go
  - frontend/tests/unit/workbench-v2/profile-reduction-equivalence.test.ts
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/recording-workflow-unified.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recording_reduction_test.go
  - frontend/tests/unit/workbench-v2/device-recording-summary.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/router_studio_v2_wide_counter_parity_test.go
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/datalink/grouppipeline/average_unrelated_test.go
  - internal/datalink/groupdelivery/source_evidence_wide_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_recording_start_fixture_test.go
  - frontend/tests/fixtures/cutoverFixtures.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - internal/datalink/dbtarget/group_row_layout_wide_counter_test.go
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - internal/datalink/measurement/mean_test.go
  - internal/api/router_studio_v2_mean_long_runtime_test.go
  - internal/datalink/workspace/write_group_wide_cumulative_shape_test.go
  - frontend/tests/unit/workbench-v2/writeGroupProposal.test.tsx
  - internal/datalink/runtime/group_boundary_wide_profile_test.go
  - internal/datalink/snapshot/mean_evidence_test.go
  - frontend/tests/unit/workbench-v2/device-recording-purpose-flow.test.tsx
  - internal/api/device_recording_integration_test.go
  - frontend/tests/fixtures/sourceReview.ts
  - frontend/tests/unit/workbench-v2/schema-operation.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
-->

---
### Requirement: Purpose-led layout with one advanced authority

The system SHALL label the two new managed paths as whole-device snapshot (wide) and selected measurements (long), explain their row granularity, and allow instantaneous and cumulative roles in either layout. Advanced SHALL open the same saved group with the same profile/reduction/interval and explicit saved-versus-applied state. It SHALL not expose a competing interval or start authority. Every database device SHALL enter this shared recording setup by default, including a device with a saved legacy managed wide group. The presence of that group SHALL NOT automatically replace the shared setup with its legacy editor or block configuration of a new recording profile. Legacy wide/custom settings SHALL remain editable under their own label without silent conversion; the legacy managed editor SHALL require an explicit operator action and SHALL retain its saved values and scoped drafts independently of new profile drafts. A legacy applied revision SHALL be described as saved applied settings rather than proof of current collection. New starts SHALL use the existing scoped request and server ownership checks; layout changes and takeover of still-owned sources SHALL use reviewed version/ownership transitions instead of changing IDs, deleting groups or disabling an existing writer automatically. The existing explicit preparation-only context SHALL remain unable to restore or send a recording start.

#### Scenario: Full meter and mixed selected sources
- **WHEN** one device selects full voltage/current/power/counter snapshots and another selects temperature plus water/gas/electric counters
- **THEN** the first can use managed wide and the second managed long, each preserves explicitly selected roles and periods, and neither layout implies that counters are averaged.

#### Scenario: Advanced round trip
- **WHEN** an operator opens a profile from Basic in Advanced and returns without saving
- **THEN** group ID, reduction, interval, members and destination remain identical and no mutation or additional start occurs.

#### Scenario: Legacy and new devices share the default setup
- **WHEN** A has a saved legacy managed wide group and B has none, and the operator switches between them
- **THEN** both devices show the shared purpose and instantaneous/cumulative setup, their new drafts remain distinct, and viewing either device changes no saved group or writer.

#### Scenario: Explicit legacy editing preserves independent drafts
- **WHEN** the operator opens legacy editing on A, changes its 60-second draft to 73 seconds, returns to the new setup with a 47-second draft, or reloads
- **THEN** the corresponding draft is restored only in its own scope and editor, and neither value is automatically saved or copied to the other draft or device.

#### Scenario: Never-applied legacy draft does not join a new start
- **WHEN** A retains a never-applied legacy draft and the operator starts an already prepared new profile
- **THEN** the existing scoped request includes only A and its reviewed profile groups; the legacy draft is not updated, applied, disabled or deleted.

#### Scenario: Applied legacy ownership remains authoritative
- **WHEN** a device has an applied legacy revision and the server refuses a new start that would take over its owned sources
- **THEN** the UI exposes the refusal with the existing settings and a review-cutover entry, preserves the old group, and does not bypass ownership through a second Apply or workspace activation; explicit preparation-only mode still performs no recording start.


<!-- @trace
source: fix-device-recording-shared-entry
updated: 2026-10-25
code:
  - frontend/src/features/datalink/workbench-v2/steps/step4/ExistingDeviceRecording.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - docs/technical/studio-v2-circuit-recording.md
tests:
  - frontend/tests/e2e/fixtures/device-recording-entry.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-workflow-unified.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-summary.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-tiles.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/e2e/fixtures/device-recording-entry.html
-->

---
### Requirement: Compact source review preserves actionable exceptions

The source review SHALL default to a collapsed item count with an explicit expand action and an internally scrollable expanded list. Names SHALL use saved display names then source addresses, with technical keys restricted to diagnostics. Missing/foreign-source counts and repair actions SHALL remain visible while collapsed. Acquisition, recording intervals and later report windows SHALL be explained separately without adding a report engine or new acquisition cadence controls.

#### Scenario: Large source list and an unavailable member
- **WHEN** a saved group contains 32 members including one unavailable source
- **THEN** the default view shows the count and actionable unavailable-source summary without 32 permanently expanded rows, and expansion reveals the scoped source details without changing configuration.

#### Scenario: Recording versus report interval
- **WHEN** the selected instantaneous interval is 15 seconds and the cumulative interval is 10 minutes
- **THEN** the summary explains that accepted acquisition feeds those independent recording windows and does not claim to restore 15-second history from only 10-minute records or create reports.

<!-- @trace
source: unify-device-recording-workflow
updated: 2026-10-25
code:
  - internal/datalink/workspace/write_group_long_recording_shape.go
  - internal/datalink/workspace/write_group_member_metadata.go
  - internal/datalink/workspace/write_group_readiness_partial.go
  - frontend/src/utils/studioV2WriteGroupSourceEvidenceJson.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingMembers.tsx
  - frontend/src/features/datalink/workbench-v2/state/basicLegacyDraft.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/technical/studio-v2-circuit-recording.md
  - internal/datalink/grouppipeline/delivery_sources.go
  - internal/datalink/groupdelivery/source_evidence_wide.go
  - internal/datalink/workspace/write_group_cutover_repository.go
  - internal/datalink/measurement/mean.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel.tsx
  - internal/datalink/dbtarget/group_long_row.go
  - internal/datalink/grouppipeline/counter_sources.go
  - internal/datalink/workspace/write_group_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusView.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitPurposeFields.tsx
  - internal/datalink/workspace/write_group_schema.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/workspace/write_group_cutover_forward.go
  - internal/datalink/workspace/write_group_wide_profile_shape.go
  - internal/datalink/dbtarget/group_long_row_counter.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitRequiredPolicy.tsx
  - internal/datalink/workspace/write_group_repository_members.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/circuitProfileActions.ts
  - docs/swagger/docs.go
  - internal/datalink/workspace/write_group_profile.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/workspace/write_group_repository_validation.go
  - frontend/src/features/datalink/workbench-v2/state/deviceScopeSelection.ts
  - frontend/src/features/datalink/workbench-v2/state/basicRecording.ts
  - internal/datalink/runtime/group_boundary_wide_profile.go
  - internal/datalink/groupdelivery/source_evidence_aggregation.go
  - internal/datalink/snapshot/config.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection.tsx
  - internal/datalink/runtime/group_boundary_build.go
  - internal/datalink/snapshot/selector.go
  - internal/datalink/groupdelivery/source_evidence.go
  - docs/swagger/swagger.json
  - internal/datalink/grouppipeline/recovery_long.go
  - internal/datalink/workspace/write_group_readiness.go
  - internal/datalink/snapshot/assembler_profile.go
  - internal/datalink/snapshot/outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupProfileSummary.tsx
  - internal/api/handlers/studio_v2_workspace_write_group_errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SourceReview.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/groupdelivery/source_evidence_consistency.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBatchIntent.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus.tsx
  - internal/datalink/grouppipeline/reconcile.go
  - frontend/src/types/studioV2WriteGroup.ts
  - internal/datalink/dbtarget/group_row_layout.go
  - internal/datalink/workspace/write_group.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel.tsx
  - internal/datalink/workspace/write_group_repository_update.go
  - internal/datalink/runtime/group_boundary.go
  - internal/datalink/schema/migrations/038_write_group_recording_reduction_sqlite.up.sql
  - internal/datalink/snapshot/assembler_counter.go
  - frontend/src/utils/studioV2RecordingProfile.ts
  - internal/datalink/snapshot/assembler_mean.go
  - internal/datalink/snapshot/mean_evidence.go
  - internal/datalink/workspace/write_group_basic_managed.go
  - internal/datalink/workspace/write_group_counter_readiness.go
  - internal/datalink/workspace/write_group_long_readiness.go
  - internal/datalink/workspace/write_group_profile_interval.go
  - docs/swagger/swagger.yaml
  - frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts
  - internal/datalink/grouppipeline/recovery.go
  - frontend/src/features/datalink/workbench-v2/state/circuitProfileBasic.ts
  - internal/datalink/snapshot/assembler_wide_profile.go
  - internal/datalink/runtime/group_boundary_long.go
  - internal/datalink/dbtarget/group_row_layout_provenance.go
  - frontend/src/utils/safeJson.ts
  - internal/datalink/migrator_write_groups.go
  - internal/datalink/snapshot/assembler.go
  - internal/datalink/workspace/write_group_cutover_apply.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeviceScopeTiles.tsx
  - frontend/src/types/studioV2WriteGroupDelivery.ts
  - internal/datalink/groupdelivery/source_evidence_verify.go
  - internal/datalink/workspace/write_group_profile_repository.go
  - internal/datalink/workspace/write_group_readiness_table.go
  - frontend/src/utils/studioV2WriteGroupJson.ts
tests:
  - internal/datalink/grouppipeline/mean_result_type_proof_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-cutover.test.tsx
  - internal/api/router_studio_v2_reviewed_wide_mean_transition_test.go
  - internal/datalink/snapshot/mean_test.go
  - internal/api/router_studio_v2_write_group_profile_test.go
  - internal/api/router_studio_v2_wide_profile_guard_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-status-review.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-status-scope.test.tsx
  - internal/api/router_studio_v2_recording_start_scope_test.go
  - internal/api/router_studio_v2_wide_profile_runtime_test.go
  - frontend/tests/unit/workbench-v2/write-group-aggregation-evidence.test.ts
  - internal/datalink/dbtarget/group_long_row_mean_test.go
  - internal/datalink/dbtarget/group_row_layout_wide_profile_test.go
  - internal/datalink/grouppipeline/average_recovery_test.go
  - internal/datalink/grouppipeline/recovery_reduction_test.go
  - internal/datalink/runtime/group_boundary_wide_counter_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx
  - frontend/tests/unit/workbench-v2/step4-evidence-boundaries.test.tsx
  - internal/datalink/grouppipeline/wide_profile_result_type_proof_test.go
  - internal/datalink/workspace/write_group_member_recording_shapes_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose-advanced.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-tiles.test.tsx
  - frontend/tests/unit/workbench-v2/wide-profile-member-metadata.test.ts
  - frontend/tests/e2e/fixtures/device-recording-workflow.ts
  - frontend/tests/unit/services/studioV2WriteGroupDeliveryEnvelope.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx
  - frontend/tests/unit/workbench-v2/device-recording-purpose.test.tsx
  - internal/datalink/grouppipeline/delivery_sources_test.go
  - internal/datalink/runtime/group_boundary_long_reduction_test.go
  - frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-intent.test.ts
  - internal/datalink/runtime/group_boundary_long_counter_raw_test.go
  - internal/datalink/workspace/write_group_recording_profile_test.go
  - frontend/tests/e2e/device-recording-workflow.spec.ts
  - internal/api/router_studio_v2_recording_workflow_coexistence_test.go
  - internal/datalink/snapshot/assembler_wide_cumulative_test.go
  - internal/datalink/workspace/write_group_circuit_profile_fixture_test.go
  - internal/datalink/dbtarget/group_long_row_legacy_provenance_test.go
  - internal/datalink/migrator_circuit_recording_profiles_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/basic-recording-panel-regressions.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-actions-reduction.test.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-status-panel.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-stale.test.tsx
  - internal/datalink/grouppipeline/average_fixture_test.go
  - frontend/tests/unit/workbench-v2/step4-source-display-names.test.tsx
  - internal/datalink/grouppipeline/average_delivery_evidence_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-prepare.test.tsx
  - internal/datalink/snapshot/assembler_wide_profile_test.go
  - frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx
  - internal/datalink/workspace/write_group_member_metadata_test.go
  - frontend/tests/unit/workbench-v2/profile-reduction-equivalence.test.ts
  - frontend/tests/unit/workbench-v2/profile-boundaries.test.ts
  - frontend/tests/unit/workbench-v2/recording-workflow-unified.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-display-names.test.tsx
  - internal/datalink/workspace/write_group_recording_reduction_test.go
  - frontend/tests/unit/workbench-v2/device-recording-summary.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/api/router_studio_v2_wide_counter_parity_test.go
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-migration-timescale.test.tsx
  - internal/datalink/grouppipeline/average_unrelated_test.go
  - internal/datalink/groupdelivery/source_evidence_wide_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_recording_start_fixture_test.go
  - frontend/tests/fixtures/cutoverFixtures.ts
  - frontend/tests/unit/workbench-v2/circuit-profile-preparation.test.tsx
  - internal/datalink/dbtarget/group_row_layout_wide_counter_test.go
  - internal/api/router_studio_v2_write_group_cutover_rig_test.go
  - internal/datalink/measurement/mean_test.go
  - internal/api/router_studio_v2_mean_long_runtime_test.go
  - internal/datalink/workspace/write_group_wide_cumulative_shape_test.go
  - frontend/tests/unit/workbench-v2/writeGroupProposal.test.tsx
  - internal/datalink/runtime/group_boundary_wide_profile_test.go
  - internal/datalink/snapshot/mean_evidence_test.go
  - frontend/tests/unit/workbench-v2/device-recording-purpose-flow.test.tsx
  - internal/api/device_recording_integration_test.go
  - frontend/tests/fixtures/sourceReview.ts
  - frontend/tests/unit/workbench-v2/schema-operation.test.tsx
  - frontend/tests/unit/workbench-v2/circuit-profile-sources-a11y.test.tsx
-->