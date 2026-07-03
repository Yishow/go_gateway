# workspace-database-row-groups Specification

## Purpose

TBD - created by archiving change 'support-step4-database-row-groups'. Update Purpose after archive.

## Requirements

### Requirement: Workspace database row groups capture shared-column intent

The workspace SHALL persist database row groups as first-class Step 4 planning entities whenever an operator wants multiple points to reuse the same business columns in one table.

#### Scenario: Persist one shared-column row group

- **WHEN** an operator creates a row group for one connector and one table, assigns multiple points to it, and saves Step 4
- **THEN** the workspace persists a stable row-group record with its id, connector/table scope, member points, and group-key metadata
- **AND** the member database targets persist their referenced row-group id instead of being treated as unrelated single-row targets

#### Scenario: Reload row-group planning on next workspace load

- **WHEN** `/studio/v2` reloads a workspace that already contains persisted row groups
- **THEN** Step 4 restores the same row-group structure and shared-column memberships
- **AND** the workspace SHALL NOT silently flatten those members back into global unique-column bindings


<!-- @trace
source: support-step4-database-row-groups
updated: 2026-06-09
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/RowGroupPlanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/state/rowGroupValidation.ts
  - internal/datalink/workspace/service_readiness_row_groups.go
  - internal/api/handlers/studio_v2_workspace_database_row_groups.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/workspace/service_database_row_groups.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/workspace/service_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - Makefile
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step4-database-row-groups.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/row-group-validation.test.ts
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/datalink/workspace/service_readiness_row_groups_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/database-row-groups-autosave.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
-->

---
### Requirement: Row groups define the legality of shared business columns

The system SHALL treat repeated business columns as legal only when every repeated binding belongs to the same row group and that row group defines a valid row identity contract.

#### Scenario: Shared column stays legal inside one row group

- **WHEN** two points inside the same row group both bind to `temperature_c`
- **THEN** validation marks the plan as legal if the row group contains a valid row identity contract
- **AND** the plan SHALL NOT surface a generic duplicate-column blocker

#### Scenario: Shared column remains blocking across row groups

- **WHEN** two points in different row groups bind to `temperature_c` for the same connector/table scope
- **THEN** validation reports a blocking conflict
- **AND** the conflict SHALL identify the cross-group reuse rather than a generic single-point column error


<!-- @trace
source: support-step4-database-row-groups
updated: 2026-06-09
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/RowGroupPlanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/state/rowGroupValidation.ts
  - internal/datalink/workspace/service_readiness_row_groups.go
  - internal/api/handlers/studio_v2_workspace_database_row_groups.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/workspace/service_database_row_groups.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/workspace/service_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - Makefile
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step4-database-row-groups.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/row-group-validation.test.ts
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/datalink/workspace/service_readiness_row_groups_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/database-row-groups-autosave.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
-->

---
### Requirement: Row-group write mode guardrails protect unsafe upsert reuse

The system SHALL block row-group configurations that reuse business columns under `upsert` unless the row group defines a stable uniqueness key that can distinguish row instances.

#### Scenario: Insert-mode shared-column row group is allowed

- **WHEN** a row group reuses business columns and the connector write mode is `insert`
- **THEN** readiness allows the plan to proceed if all other validation passes

#### Scenario: Upsert without uniqueness key is rejected

- **WHEN** a row group reuses business columns and the connector write mode is `upsert` without a stable uniqueness key
- **THEN** readiness reports a blocking issue
- **AND** apply SHALL NOT continue until the operator supplies a valid uniqueness contract or changes the write mode

<!-- @trace
source: support-step4-database-row-groups
updated: 2026-06-09
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/RowGroupPlanner.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/state/rowGroupValidation.ts
  - internal/datalink/workspace/service_readiness_row_groups.go
  - internal/api/handlers/studio_v2_workspace_database_row_groups.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/workspace/service_database_row_groups.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/workspace/service_readiness.go
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/datalink/workspace/service_readiness_connector_issue.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - Makefile
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
tests:
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step4-database-row-groups.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/row-group-validation.test.ts
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/datalink/workspace/service_readiness_row_groups_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_row_groups_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/dbtarget/writer_row_groups_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/database-row-groups-autosave.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database-controls.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
-->