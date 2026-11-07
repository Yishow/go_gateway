# simple-recording-lifecycle Specification

## Purpose

Make existing recording lifecycle and recovery operations understandable from the simple Studio interface while retaining server authority, forward-only storage changes and truthful outcomes.

## Requirements

### Requirement: Recording actions follow matching server availability
The system SHALL derive saved-recording actions from matching workspace, group and revision evidence. Stop and delete SHALL require explicit scope and impact confirmation. Deleting a recording SHALL permanently delete its internal configuration and delivery records and SHALL NOT delete historical database rows or tables in the destination database.

#### Scenario: The action facts are stale
- **WHEN** lifecycle evidence does not match the displayed recording revision
- **THEN** the interface disables consequential actions, explains the need to reload evidence and does not guess availability from status.

#### Scenario: A recording is deleted
- **WHEN** an operator acknowledges that deletion cannot be undone and confirms an allowed deletion
- **THEN** only the recording deletion contract is invoked, the interface states that the recording configuration was permanently deleted while the destination table and its rows were retained, and the recording no longer appears in the home or editor lists

#### Scenario: Deletion is refused by the server
- **WHEN** the deletion contract returns `WRITE_GROUP_DELETE_ACTIVE`, `WRITE_GROUP_DELETE_BACKLOG`, `WRITE_GROUP_DELETE_OPERATION_PENDING` or `WRITE_GROUP_DELETE_HAS_SUCCESSOR`
- **THEN** the interface keeps the recording, shows an actionable explanation with the matching stop, open-recording or refresh action, and does not display raw backend text

#### Scenario: A referenced source is changed
- **WHEN** an operator attempts to delete or redefine a source used by existing recordings
- **THEN** the interface preserves existing dependency protections and identifies affected recordings or opens the scoped advanced repair flow without creating dangling references.


<!-- @trace
source: improve-recording-home-status-and-delete
updated: 2026-11-08
code:
  - internal/datalink/dbtarget/writer_snapshot.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar.tsx
  - frontend/studio-v2-diagnostic.cjs
  - internal/datalink/workspace/write_group_service.go
  - internal/datalink/runtimeview/contract.go
  - docs/validation/memory-views-20261106/skill/SKILL.md
  - frontend/src/services/runtimeLogLegacyProjection.ts
  - internal/datalink/dbtarget/writer_workers.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeResultFacts.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-status.json
  - internal/datalink/workspace/write_group_permanent_delete_audit.go
  - docs/testing.md
  - frontend/src/features/datalink/workbench-v2/recording-home/DeviceDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingLifecycleHistory.tsx
  - internal/datalink/workspace/device_recording_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary.tsx
  - internal/datalink/workspace/write_group_permanent_delete.go
  - internal/diagnostics/projection.go
  - internal/api/router.go
  - frontend/src/services/runtimeLogDomainProjection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/recording-delete-dialog.css
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/modbusshare/service_projection.go
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - frontend/src/i18n/locales/zh-TW/runtime-logs.json
  - internal/api/router_studio_v2_write_group_routes.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatusCounts.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx
  - internal/datalink/workspace/write_group_permanent_delete_purge.go
  - internal/diagnostics/domain_event.go
  - docs/swagger/swagger.json
  - docs/swagger/docs.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeLifecycleMatches.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeWriteStages.tsx
  - internal/api/handlers/typed_errors.go
  - internal/datalink/device/workspace_delete.go
  - internal/datalink/groupdelivery/quota.go
  - frontend/src/services/runtimeRuleEvidenceParser.ts
  - AGENTS.md
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - internal/datalink/workspace/write_group_lifecycle_actions.go
  - frontend/src/types/runtimeRuleEvidence.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeStatusDistribution.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start.go
  - frontend/src/i18n/locales/en/recording-home.json
  - frontend/src/services/studioV2WorkspaceWriteGroups.ts
  - internal/datalink/workspace/write_group_permanent_delete_shared.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_delete.go
  - internal/datalink/audit/types.go
  - internal/datalink/dbtarget/writer_failure_log.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDelete.ts
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_mutate.go
  - internal/datalink/dbtarget/writer_flush.go
  - frontend/src/hooks/datalink/useStudioV2WriteGroups.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx
  - internal/testutil/legacyrecording/tombstone.go
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete.go
  - frontend/src/features/datalink/workbench-v2/recording-home/deviceConnectionState.ts
  - frontend/src/i18n/locales/en/recording-home-status.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRowActions.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHome.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDestinationCell.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeWriteStageMapping.ts
  - internal/datalink/dbtarget/writer_group_flush.go
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeRows.css
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/dbtarget/writer_retry.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeUnconfirmed.tsx
  - frontend/src/features/runtime-logs/RuntimeLogLegacyEvent.tsx
  - internal/datalink/workspace/write_group_permanent_delete_operations.go
  - internal/datalink/workspace/write_group_permanent_delete_storage.go
  - frontend/src/i18n/config.ts
  - internal/datalink/workspace/write_group_permanent_delete_scope.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeStatus.css
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceActions.tsx
  - frontend/src/features/datalink/workbench-v2/recording-home/useRecordingHomeDeviceStatus.ts
  - docs/swagger/swagger.yaml
  - frontend/src/i18n/locales/en/recording-home-rows.json
  - frontend/src/types/studioV2Delete.ts
  - frontend/src/types/deviceConnection.ts
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRecordRow.tsx
  - internal/datalink/workspace/write_group_permanent_delete_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeDeviceStrip.tsx
  - frontend/src/i18n/locales/zh-TW/recording-home-rows.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/dev/previewI18n.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - docs/validation/memory-views-20261106/skill/references/lessons.md
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController.tsx
  - internal/datalink/grouppipeline/backlog_guard.go
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeCollection.tsx
  - frontend/src/i18n/locales/en/recording-home-delete.json
  - frontend/src/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogError.tsx
  - frontend/src/i18n/locales/en/runtime-logs.json
tests:
  - internal/api/router_studio_v2_write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_permanent_delete_delivery_fixture_test.go
  - internal/api/router_studio_v2_successor_first_start_fixture_test.go
  - frontend/tests/unit/runtime-logs/legacy-event-fixtures.ts
  - frontend/tests/unit/workbench-v2/saved-recording-actions.test.tsx
  - internal/api/handlers/studio_v2_workspace_write_groups_handler_delete_test.go
  - frontend/tests/unit/workbench-v2/groupDeliveryStrip.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-accessibility.test.tsx
  - internal/api/router_studio_v2_locator_entry_test.go
  - frontend/tests/unit/dev/recording-style-bootstrap.test.tsx
  - internal/datalink/workspace/recording_start_disabled_test.go
  - internal/datalink/dbtarget/writer_failure_revision_test.go
  - cmd/test_ui/group_pipeline_outage_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-editor-range-mc-prepare.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-memory-scope.test.tsx
  - frontend/tests/unit/workbench-v2/recording-guided-saved-card.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_test.go
  - internal/datalink/workspace/write_group_mixed_safety_test.go
  - frontend/tests/unit/workbench-v2/recording-home-route.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_storage_test.go
  - frontend/tests/unit/workbench-v2/device-connection-state.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-status-counts.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-shell.test.tsx
  - frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx
  - cmd/test_ui/gateway_logging.go
  - internal/api/router_studio_v2_write_group_boundaries_test.go
  - frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_lifecycle_test.go
  - internal/datalink/workspace/write_group_circuit_profile_rules_test.go
  - frontend/tests/unit/workbench-v2/recording-delete-dialog.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_test.go
  - internal/api/router_studio_v2_write_group_action_errors_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogLegacyProjection.test.ts
  - cmd/test_ui/delivery_events_unavailable_test.go
  - frontend/tests/e2e/fixtures/recording-lifecycle-api.ts
  - frontend/tests/unit/runtime-logs/RuntimeLogCompatibility.test.tsx
  - internal/datalink/dbtarget/writer_failure_log_test.go
  - internal/api/router_runtime_rule_evidence_test.go
  - internal/api/router_studio_v2_write_group_stop_remove_safety_test.go
  - internal/datalink/groupdelivery/quota_reclaim_receipts_test.go
  - internal/datalink/workspace/write_group_basic_managed_test.go
  - frontend/tests/unit/workbench-v2/recording-modal-return.test.tsx
  - frontend/tests/unit/services/studioV2WorkspaceWriteGroups.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-source-switching.test.tsx
  - cmd/test_ui/group_pipeline_outage_test.go
  - frontend/tests/unit/workbench-v2/recording-home-row-actions.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-legacy-delivery-safety.test.tsx
  - frontend/tests/unit/workbench-v2/device-delete-cache.test.tsx
  - internal/datalink/workspace/write_group_lifecycle_error_preservation_test.go
  - internal/datalink/workspace/write_group_reclaim_delete_test.go
  - cmd/test_ui/group_pipeline_lifecycle_crash_test.go
  - frontend/tests/unit/workbench-v2/recording-home-visual-status.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-status.test.tsx
  - internal/datalink/workspace/write_group_permanent_delete_legacy_test.go
  - frontend/tests/unit/services/studioV2DeleteResponses.test.ts
  - internal/api/router_studio_v2_recording_start_delete_test.go
  - frontend/tests/unit/workbench-v2/recording-home-memory-view.test.tsx
  - frontend/tests/unit/dev/recording-style-preview.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx
  - internal/datalink/workspace/write_group_successor_replay_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_fixture_test.go
  - internal/datalink/workspace/write_group_lifecycle_actions_test.go
  - frontend/tests/e2e/recording-home-status-layout.spec.ts
  - frontend/tests/e2e/recording-home-tombstone-delete.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-inline-evidence.test.tsx
  - frontend/tests/unit/workbench-v2/recordingLifecycleBrowserFixture.test.ts
  - internal/datalink/workspace/write_group_service_test.go
  - internal/datalink/workspace/write_group_tombstone_remove_actions_test.go
  - cmd/test_ui/service_wiring.go
  - cmd/test_ui/group_pipeline_legacy_takeover_test.go
  - frontend/tests/unit/runtime-logs/runtimeLogDomainProjection.test.ts
  - frontend/tests/unit/workbench-v2/recording-home-delete-integration.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-device-actions.ts
  - frontend/tests/unit/workbench-v2/recording-home-mc-resave.test.tsx
  - frontend/tests/unit/workbench-v2/writeGroupLifecycle.test.tsx
  - internal/datalink/grouppipeline/backlog_guard_delete_count_test.go
  - cmd/test_ui/group_pipeline_lifecycle_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/datalink/grouppipeline/applied_configuration_read_test.go
  - internal/datalink/workspace/write_group_permanent_delete_recording_start_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go
  - internal/api/router_studio_v2_recording_delete_contract_test.go
  - frontend/tests/unit/workbench-v2/device-delete-dialog.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-management.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx
  - frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx
  - frontend/tests/e2e/recording-database-unavailable.spec.ts
  - internal/datalink/workspace/write_group_successor_guard_test.go
  - internal/datalink/dbtarget/writer_failure_log_outcome_test.go
  - internal/api/router_studio_v2_write_group_lifecycle_history_test.go
  - internal/datalink/grouppipeline/stop_remove_intake_fence_test.go
  - cmd/test_ui/gateway_logging_database_test.go
  - frontend/tests/e2e/recording-home-narrow-records.spec.ts
  - frontend/tests/unit/workbench-v2/recording-home-tombstone-delete.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_recordings_test.go
  - internal/datalink/workspace/write_group_lifecycle_error_test.go
  - internal/datalink/workspace/write_group_permanent_delete_schema_test.go
  - internal/datalink/workspace/write_group_permanent_delete_shared_test.go
  - frontend/tests/unit/workbench-v2/recording-home-new-range-strict-ack.test.tsx
-->

---
### Requirement: Resume revalidates and copying creates only a new draft
The system SHALL resume through existing scoped review and start contracts. Copy SHALL prefill a new inactive draft without old identity, operation, readiness, receipt or ownership proof. Unsupported representations or managed-profile deduplication SHALL be explained rather than reported as successful new recordings.

#### Scenario: A stopped recording is resumed
- **WHEN** an operator requests resume
- **THEN** current source, destination, revision and readiness evidence is checked and explicit start confirmation is required before waiting for a new write result.

#### Scenario: A supported recording is copied
- **WHEN** an operator copies a representable recording
- **THEN** an inactive draft is opened with editable descriptive settings, the original is unchanged, and no command authority or success evidence is copied.

#### Scenario: Managed ownership prevents a second copy
- **WHEN** submitting a copied managed profile would resolve to the existing canonical role
- **THEN** the interface explains that a distinct duplicate is not available and does not claim a new group or writer was created.




---
### Requirement: Table changes use reviewed forward-only successors
The system SHALL use existing successor review, reservation, schema planning, confirmation and scoped start or cutover for supported layout or destination changes. The interface SHALL show affected sources, old and new destinations, retained history and the boundary. It MUST NOT perform destructive in-place layout conversion or historical data migration.

#### Scenario: A supported wide recording is changed to long layout
- **WHEN** the operator selects an allowed transition
- **THEN** the interface reviews a new target table and cutover impact through existing guards, retains the old table and history, and distinguishes preparing, scheduled and effective states.

#### Scenario: The reverse layout is not supported
- **WHEN** an operator requests a layout or destination combination that the existing backend rejects
- **THEN** the interface displays the matching blocker and repair or advanced path without claiming conversion or expanding backend capabilities.

#### Scenario: A first start differs from a cutover
- **WHEN** the successor has no prior effective recording
- **THEN** the interface uses the existing first-start review instead of presenting a completed historical cutover.




---
### Requirement: Unknown operations recover without replacement commands
The system SHALL retain submitted operation identity and truthful pending or unknown states when dialogs close or navigation changes. Closing a dialog SHALL NOT imply cancellation. Recovery SHALL query the original operation before a replacement mutation is permitted by the existing contract. Existing exact-retry or repreview actions SHALL remain available only after fresh matching status evidence satisfies their original guards and the operator explicitly requests the action.

#### Scenario: A response is lost and the dialog closes
- **WHEN** a schema or activation request has an unknown outcome and the operator closes then reopens the dialog
- **THEN** the same operation is shown for status lookup without automatically resending it or creating a new operation.

#### Scenario: A late result belongs to another scope
- **WHEN** an old operation response arrives after switching the selected recording
- **THEN** it does not update the new recording or enable its actions.

#### Scenario: Fresh evidence permits exact retry or repreview
- **WHEN** recovery returns fresh matching evidence that satisfies the existing not-submitted or preview-not-submitted retry guards
- **THEN** the interface offers explicit exact retry with the original request identity, or the existing repreview flow with its permitted new preview identity and retained root-operation lineage and prior previews; neither action weakens its original proof requirements.

#### Scenario: Recovery does not prove retry is safe
- **WHEN** the operation status remains unknown, stale or scope-mismatched
- **THEN** exact retry and repreview remain unavailable and no replacement command is submitted.




---
### Requirement: Results and repair paths are actionable without invented values
The system SHALL display source and database writing rules, current scoped state, last committed evidence and scoped repair actions before optional technical detail. Partial, stale, failed and unknown outcomes SHALL remain distinct. General delivery evidence SHALL NOT be represented as raw database rows. Unsupported existing configurations SHALL retain read-only summaries and advanced access.

#### Scenario: One source is stale while the database is available
- **WHEN** current evidence identifies partial or stale source data
- **THEN** the interface identifies the affected source and opens its point or device repair context without labelling the database as the failing component.

#### Scenario: SQL writing fails after earlier success
- **WHEN** a recording has a previous receipt and a current database-write failure
- **THEN** the interface retains the historical timestamp, explains the current database issue and offers the scoped destination or logs action without displaying sample rows as actual results.

#### Scenario: Advanced-only configuration is inspected
- **WHEN** the simple view cannot represent an existing configuration
- **THEN** its rules are summarized without changing them and its advanced locator remains reachable.