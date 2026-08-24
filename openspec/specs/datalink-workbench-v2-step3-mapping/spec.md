# datalink-workbench-v2-step3-mapping Specification

## Purpose

TBD - created by archiving change 'datalink-workbench-v2-step3-mapping'. Update Purpose after archive.

## Requirements

### Requirement: Inline-editable mapping table

The Step 3 mapping table SHALL render one row per enabled point in `state.points`. Each row MUST display the rule color dot (sourced from `useDeviceColor` / rule color cycle), point name (mono), point address (mono blue-200), and seven inline-editable controls: Tag Key text input (28% column width), display name text input, unit text input, target type select (10 options from DATA_TYPES), × Scale number input (step 0.01, width 16), + Offset number input (step 0.01, width 16), and an enabled toggle (size sm). Clicking anywhere outside these inline controls MUST set the row as the current `selectedIdx`; the selected row MUST receive a blue outline and 6% blue background. Clicking, mousedown, or focusing any inline control MUST NOT change `selectedIdx`.

#### Scenario: Default row rendering

- **GIVEN** state has 8 enabled points
- **WHEN** Step 3 mounts
- **THEN** the mapping table renders exactly 8 rows
- **AND** each row's Tag Key input value matches the per-index default from `POINT_SEMANTIC` (e.g. row 0: `line01.temp.inlet`)
- **AND** each row's scale input value equals `state.points[i]._rule_scale`
- **AND** each row's offset input value equals `state.points[i]._rule_offset`
- **AND** each row's target type select value is `float64`

#### Scenario: Row selection by clicking blank area

- **WHEN** the operator clicks the row body (not on any input/select/toggle)
- **THEN** `selectedIdx` becomes that row's index
- **AND** the row gains an outline `outline-1 outline-blue-500/30` and background `bg-blue-500/[0.06]`

#### Scenario: Clicking inline control does not select row

- **GIVEN** `selectedIdx === 0`
- **WHEN** the operator clicks (or focuses) the Tag Key input in row 3
- **THEN** `selectedIdx` remains `0`
- **AND** the row 3 input gains focus
- **AND** typing into the input updates `state.mappings[points[3].id].tag_key` on every change

#### Scenario: Toggle enabled does not select row

- **GIVEN** `selectedIdx === 0`
- **WHEN** the operator clicks the enabled toggle in row 2
- **THEN** `selectedIdx` remains `0`
- **AND** `state.mappings[points[2].id].enabled` is toggled


<!-- @trace
source: datalink-workbench-v2-step3-mapping
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/package.json
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
tests:
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
-->

---
### Requirement: Automatic mapping initialization

The system SHALL automatically initialize a mapping for every enabled point that does not yet have one. Initialization MUST be triggered by Step 3 mounting or by the set of enabled point IDs changing. For each new mapping, defaults MUST be: `tag_key` from `POINT_SEMANTIC[i % 8].tag_key`, `display_name` from `POINT_SEMANTIC[i % 8].display`, `unit` from `POINT_SEMANTIC[i % 8].unit`, `target_type` = `float64`, `scale` = `point._rule_scale`, `offset` = `point._rule_offset`, `enabled` = `true`. Existing mappings MUST NOT be overwritten. Mappings whose point no longer exists in `state.points` MUST be removed.

#### Scenario: Initialize new mappings on first mount

- **GIVEN** `state.mappings === {}` and `state.points` contains 8 enabled points
- **WHEN** Step 3 mounts and the `initMappingsForPoints` action runs
- **THEN** `state.mappings` contains 8 entries keyed by each point's `id`
- **AND** each mapping's defaults match the table above

#### Scenario: Preserve existing mappings

- **GIVEN** the operator has edited `state.mappings[p3].tag_key` to `custom.tag.x`
- **WHEN** Step 3 remounts (or any other trigger of `initMappingsForPoints`)
- **THEN** `state.mappings[p3].tag_key` remains `custom.tag.x`
- **AND** no other mapping fields for `p3` are reset

#### Scenario: Remove orphan mappings

- **GIVEN** `state.mappings` contains entry for `pt-X` but `state.points` no longer has any point with id `pt-X` (e.g. Step 2 deleted the source rule)
- **WHEN** `initMappingsForPoints` runs
- **THEN** `state.mappings.pt-X` is removed


<!-- @trace
source: datalink-workbench-v2-step3-mapping
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/package.json
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
tests:
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
-->

---
### Requirement: Transform pipeline preview

The right-side preview panel SHALL request backend preview results for the currently selected mapping instead of computing the operator-visible transform steps entirely on the client. The request SHALL use the deterministic raw seed array `[243, 251, 1024, 985, 67, 542, 18, 1450]` indexed by selected row position modulo 8 as `raw_value`, together with the current draft transform pipeline. The panel MUST update reactively whenever the selected row's mapping `scale`, `offset`, or `target_type` changes. The panel MUST render explicit loading, success, and error states.

#### Scenario: Selected row triggers backend preview request

- **GIVEN** `selectedIdx === 0`, the row 0 mapping has `scale=0.1`, `offset=0`, `target_type='float64'`, and `tag_key='line01.temp.inlet'`
- **WHEN** the preview updates
- **THEN** the system sends `POST /api/v1/datalink/mappings/preview`
- **AND** the request includes raw seed `243`
- **AND** the request includes a transform pipeline derived from the current draft mapping values

#### Scenario: Preview loading state

- **GIVEN** a selected mapping exists
- **WHEN** a new preview request is in flight
- **THEN** the preview panel shows a loading state instead of stale success copy

#### Scenario: Preview renders backend result

- **GIVEN** backend preview returns successful `step_results` and `final_value`
- **WHEN** the response is received
- **THEN** the preview panel renders ordered decode / scale / cast / final steps from that backend result
- **AND** the final value shown to the operator comes from backend preview output

#### Scenario: Preview failure remains visible

- **GIVEN** backend preview returns an execution error
- **WHEN** the response is received
- **THEN** the preview panel shows an actionable error state
- **AND** the UI does not synthesize a local success preview

#### Scenario: Older preview response is ignored

- **GIVEN** the operator changes the selected mapping twice in quick succession
- **WHEN** the earlier preview response arrives after the later one
- **THEN** the UI keeps the latest preview result
- **AND** the earlier response does not overwrite the newer preview state


<!-- @trace
source: wire-studio-v2-step3-live-preview-and-target-type-shortcuts
updated: 2026-05-31
code:
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/errors.go
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/runtime/status.go
  - scripts/start-log-utils.sh
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_direct_reader.go
  - docs/goal.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/device/service_status.go
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - scripts/start-process-utils.sh
  - frontend/src/types/datalink.ts
  - internal/api/router.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - go.sum
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - gateway.db
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/device/availability.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/hooks/datalink/keys.ts
  - .air.toml
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/sourcerule/service_links.go
  - findings.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/api/handlers/mapping_handler.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/source_rule_handler.go
  - progress.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
tests:
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/router_runtime_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
-->

---
### Requirement: Pipeline color semantics

The pipeline step display SHALL use color semantics to communicate data flow: raw values in `slate-500` (muted), scale and offset values in `blue-300`, intermediate scaled result in `emerald-300` bold, final cast result in `emerald-200` 2xl mono. The four step circles MUST be connected by a single 1px vertical line.

#### Scenario: Color tokens applied

- **WHEN** the preview renders the scale step expression
- **THEN** the `raw` literal uses a slate-500 color class
- **AND** the `scale` and `offset` literals use a blue-300 color class
- **AND** the `= scaled` result uses an emerald-300 bold color class


<!-- @trace
source: datalink-workbench-v2-step3-mapping
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/package.json
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
tests:
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
-->

---
### Requirement: API payload preview

The preview panel SHALL render a `POST /mappings` JSON payload preview below the transform preview. The payload MUST include `point_id`, `tag_id` (`tag.${tag_key}`), `transform_pipeline` (three steps: decode, scale with params `{scale, offset}`, cast with param `{to: target_type}`), and `enabled`. The payload preview MUST update on every change to the selected mapping. The payload preview MUST remain display-only and MUST NOT initiate its own network request.

#### Scenario: Payload preview remains display-only

- **WHEN** any value in the selected mapping changes
- **THEN** the payload preview updates to reflect the latest draft values
- **AND** no additional `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by the payload preview itself


<!-- @trace
source: wire-studio-v2-step3-live-preview-and-target-type-shortcuts
updated: 2026-05-31
code:
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/errors.go
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/runtime/status.go
  - scripts/start-log-utils.sh
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_direct_reader.go
  - docs/goal.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/device/service_status.go
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - scripts/start-process-utils.sh
  - frontend/src/types/datalink.ts
  - internal/api/router.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - go.sum
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - gateway.db
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/device/availability.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/hooks/datalink/keys.ts
  - .air.toml
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/sourcerule/service_links.go
  - findings.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/api/handlers/mapping_handler.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/source_rule_handler.go
  - progress.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
tests:
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/router_runtime_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
-->

---
### Requirement: Bulk apply transform parameters

The mapping table footer SHALL provide a `將選中列的 Scale / Offset / 型態 套用到全部` (or i18n equivalent) link. Clicking the link MUST copy the selected mapping's `scale`, `offset`, and `target_type` to every mapping in `state.mappings`. The link MUST be a no-op when there is no selected row.

#### Scenario: Bulk apply propagates fields

- **GIVEN** `selectedIdx === 0` with row 0 mapping `scale=0.5`, `offset=2`, `target_type='int32'`, and there are 8 total mappings
- **WHEN** the operator clicks the bulk apply link
- **THEN** every mapping in `state.mappings` has `scale === 0.5`, `offset === 2`, `target_type === 'int32'`
- **AND** other fields (tag_key, display_name, unit, enabled) are unchanged

#### Scenario: Bulk apply with no selection

- **GIVEN** `selectedIdx === null`
- **WHEN** the operator clicks the bulk apply link
- **THEN** no mapping is modified


<!-- @trace
source: datalink-workbench-v2-step3-mapping
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/package.json
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
tests:
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
-->

---
### Requirement: Continue gate with empty tag_key detection

The bottom continue button (`設定資料庫寫入` or i18n equivalent) MUST be disabled when any enabled mapping has an empty `tag_key` (string with `.trim() === ''`) or when there are zero enabled mappings. The mapping table aside chip MUST switch between three states: success `{enabled}/{total} 已啟用` when validation passes, warning `{count} 個 Tag Key 留空` when one or more enabled mappings have empty tag_key, and warning `無啟用 mapping` when `enabledCount === 0`.

#### Scenario: Empty tag_key disables continue

- **GIVEN** 8 enabled mappings and the row 2 mapping has `tag_key=''`
- **WHEN** the operator inspects the bottom action bar
- **THEN** the continue button is disabled
- **AND** the table aside chip shows `1 個 Tag Key 留空` in warning tone

#### Scenario: All valid enables continue

- **GIVEN** all 8 enabled mappings have non-empty `tag_key`
- **WHEN** the operator inspects the bottom action bar
- **THEN** the continue button is enabled
- **AND** the aside chip shows `8/8 已啟用` in success tone
- **AND** clicking the continue button invokes the shell `onContinue` callback for Step 3

<!-- @trace
source: datalink-workbench-v2-step3-mapping
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/package.json
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
tests:
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
-->

---
### Requirement: Valid-only mapping autosave

The Step 3 mapping workspace SHALL autosave only valid mapping rows to the backend workspace.

#### Scenario: Valid mapping row saves immediately

- **GIVEN** a mapping row already belongs to the singleton v2 workspace
- **WHEN** the operator edits that row into a valid state
- **THEN** the system saves that row immediately
- **AND** the persisted mapping reflects the latest valid values

#### Scenario: Invalid mapping row stays local

- **GIVEN** a mapping row already has a last successful persisted version
- **WHEN** the operator edits that row into an invalid state
- **THEN** the system does not overwrite the persisted mapping
- **AND** the UI keeps the invalid local values visible with an unsaved marker


<!-- @trace
source: wire-studio-v2-step3-mapping-autosave
updated: 2026-05-31
code:
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/mapping/errors.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - start.sh
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - task_plan.md
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/mapping/service_crud.go
  - start.ps1
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - findings.md
  - go.sum
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/runtime/status.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/datalink/device/availability.go
  - scripts/start-process-utils.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - scripts/start-log-utils.sh
  - frontend/src/hooks/datalink/useSettings.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/runtime_handler.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/db.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - go.mod
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/device/service_status.go
  - docs/goal.md
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_logger_test.go
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
-->

---
### Requirement: Per-point mapping save isolation

The Step 3 mapping workspace SHALL isolate save success and failure per point row.

#### Scenario: One invalid row does not block another valid row

- **GIVEN** the workspace contains multiple mapping rows
- **WHEN** one row remains invalid while another row is edited into a valid state
- **THEN** the valid row still saves successfully
- **AND** the invalid row remains local with its own error state

<!-- @trace
source: wire-studio-v2-step3-mapping-autosave
updated: 2026-05-31
code:
  - internal/datalink/sourcerule/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/mapping/errors.go
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - start.sh
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - internal/api/handlers/mapping_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/types/datalink.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - task_plan.md
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/mapping/service_crud.go
  - start.ps1
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - findings.md
  - go.sum
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/runtime/status.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/datalink/device/availability.go
  - scripts/start-process-utils.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - scripts/start-log-utils.sh
  - frontend/src/hooks/datalink/useSettings.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/runtime_handler.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/db.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - go.mod
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/device/service_status.go
  - docs/goal.md
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_logger_test.go
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
-->

---
### Requirement: Target type quick actions

The Step 3 preview panel SHALL provide quick target type actions for the currently selected mapping and a one-click apply-all action for `target_type` only. The quick actions MUST include at least `bool`, `int16`, `float64`, and `string`.

#### Scenario: Quick action updates the current row

- **GIVEN** a mapping row is selected
- **WHEN** the operator clicks one of the target type quick actions
- **THEN** the selected row's `target_type` is updated to that value
- **AND** the change flows through the same Step 3 reducer/autosave path as the inline select

#### Scenario: Apply current target type to all rows

- **GIVEN** `selectedIdx === 0` and the selected row has `target_type='int16'`
- **WHEN** the operator clicks `套用到全部列`
- **THEN** every mapping in `state.mappings` has `target_type === 'int16'`
- **AND** other mapping fields remain unchanged

#### Scenario: No selected row means no quick target type action

- **GIVEN** `selectedIdx === null`
- **WHEN** the operator views the preview area
- **THEN** no target type quick action is offered as an active control

<!-- @trace
source: wire-studio-v2-step3-live-preview-and-target-type-shortcuts
updated: 2026-05-31
code:
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/errors.go
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/runtime/status.go
  - scripts/start-log-utils.sh
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/api/handlers/point_direct_reader.go
  - docs/goal.md
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/datalink/device/service_status.go
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - scripts/start-process-utils.sh
  - frontend/src/types/datalink.ts
  - internal/api/router.go
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - go.sum
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - gateway.db
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/datalink/device/availability.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/hooks/datalink/keys.ts
  - .air.toml
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/sourcerule/service_links.go
  - findings.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/api/handlers/mapping_handler.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/source_rule_handler.go
  - progress.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
tests:
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - cmd/test_ui/static/index.html
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/router_runtime_test.go
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/router_studio_v2_runtime_context_test.go
-->

---
### Requirement: Step 3 reload restores persisted mapping rows against current point identities

The Step 3 mapping workspace SHALL restore persisted mappings against the current derived point identities after reload.

#### Scenario: Reload reattaches persisted mappings to current points

- **WHEN** the operator reloads Studio V2 after mappings have already been persisted for derived points
- **THEN** Step 3 rehydrates those persisted mappings onto the correct current point rows
- **AND** Step 3 SHALL NOT silently drop persisted mappings only because point rows were reconstructed during bootstrap

<!-- @trace
source: preserve-studio-v2-setup-state
updated: 2026-06-09
code:
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - frontend/src/types/runtimeDiagnostics.ts
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - internal/datalink/workspace/service_readiness.go
  - internal/datalink/dbtarget/live_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/runtime/service_projection_state.go
  - internal/datalink/dbtarget/writer.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/studio-v2-real-check.cjs
  - internal/datalink/workspace/service_runtime_projection.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - Makefile
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/api/router.go
  - frontend/src/types/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/runtime/service_device_sync.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/dbtarget/writer_statements.go
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - internal/datalink/audit/types.go
  - frontend/src/types/studioV2Workspace.ts
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - internal/datalink/dbtarget/service.go
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/migrator_database_delivery_outcomes.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/workspace/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - internal/api/handlers/runtime_stream_handler.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/datalink/dbtarget/tooling_service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - tests/shell/start-frontend-install-failure.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/api/handlers/studio_v2_workspace_audit.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/audit/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/types/databaseDelivery.ts
  - internal/datalink/runtime/service_workspace_projection.go
  - internal/datalink/runtime/status.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/schema/schema_dbtarget_models.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - tests/shell/start-frontend-readiness.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/datalink/audit/sql_repository.go
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/App.tsx
  - internal/datalink/dbtarget/delivery_outcome.go
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - frontend/studio-v2-diagnostic.cjs
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/datalink/dbtarget/sql_repository.go
  - internal/datalink/runtime/truth_state.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - tests/shell/start-frontend-deps.sh
tests:
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - internal/datalink/workspace/service_activation_projection_test.go
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - internal/datalink/runtime/service_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/datalink/audit/service_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/api/router_runtime_database_delivery_test.go
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/workspace/service_readiness_empty_test.go
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - cmd/test_ui/main.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/runtime/ingestor_test.go
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
-->

---
### Requirement: Protocol-aware default tag key generation

The system SHALL generate default tag keys and display names in Step 3 based on the point's protocol-adapted address. The generated tag key MUST retain the register prefix identifier for non-Modbus devices (e.g. `dev.sensor.rd0` for MC/FATEK `D0` instead of stripped numeric index).

#### Scenario: Default tag key for MC 3E point
- **WHEN** entering Step 3 with an MC 3E point at address `D100`
- **THEN** the suggested tag key contains `rd100` and accurately reflects the alphanumeric address

#### Scenario: Default tag key for Modbus point
- **WHEN** entering Step 3 with a Modbus point at address `40001`
- **THEN** the suggested tag key contains `r40001`

<!-- @trace
source: fix-protocol-address-adaptation-v2
updated: 2026-08-24
code:
  - .github/workflows/backend-ci.yml
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - .github/prompts/spectra-commit.prompt.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - .github/prompts/spectra-discuss.prompt.md
  - .github/instructions/go.instructions.md
  - .github/skills/spectra-audit/SKILL.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - internal/datalink/sourcerule/repository_memory.go
  - .github/skills/openspec-apply-change/SKILL.md
  - internal/datalink/migrator.go
  - .github/prompts/spectra-audit.prompt.md
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - .github/workflows/frontend-ci.yml
  - frontend/src/features/datalink/workbench-v2/state/defaults.ts
  - .github/skills/openspec-propose/SKILL.md
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - .github/skills/spectra-drift/SKILL.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2ShareActivation.ts
  - frontend/src/types/datalink.ts
  - .github/prompts/spectra-apply.prompt.md
  - .github/skills/spectra-archive/SKILL.md
  - internal/api/handlers/source_rule_handler.go
  - internal/datalink/schema/schema_source_rule_models.go
  - internal/datalink/sourcerule/sql_rule_scan.go
  - .github/prompts/spectra-archive.prompt.md
  - .github/skills/spectra-propose/SKILL.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/datalink/schema/migrations/017_source_rule_modbus_share.down.sql
  - .github/skills/openspec-archive-change/SKILL.md
  - .github/prompts/spectra-ingest.prompt.md
  - internal/datalink/schema/migrations/017_source_rule_modbus_share_sqlite.up.sql
  - internal/datalink/sourcerule/protocol_address_planner.go
  - .github/workflows/file-line-limit.yml
  - .github/pull_request_template.md
  - .github/skills/openspec-explore/SKILL.md
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - .github/skills/spectra-ask/SKILL.md
  - .github/prompts/spectra-propose.prompt.md
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - .github/skills/spectra-commit/SKILL.md
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/sourcerule/sql_repo.go
  - .github/instructions/typescript-5-es2022.instructions.md
  - frontend/src/features/datalink/workbench-v2/state/deviceState.ts
  - .github/prompts/opsx-explore.prompt.md
  - .github/skills/spectra-discuss/SKILL.md
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - .github/prompts/spectra-debug.prompt.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .github/skills/spectra-debug/SKILL.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - .github/skills/spectra-ingest/SKILL.md
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - .github/skills/spectra-apply/SKILL.md
  - frontend/src/utils/addressParser.ts
  - .github/prompts/spectra-ask.prompt.md
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - .github/instructions/reactjs.instructions.md
  - .github/prompts/spectra-drift.prompt.md
  - internal/datalink/schema/migrations/017_source_rule_modbus_share.up.sql
  - internal/datalink/sourcerule/validation.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
tests:
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/database-autosave-page.validation.test.tsx
  - frontend/tests/unit/utils/addressParser.test.ts
  - frontend/tests/unit/workbench-v2/database-autosave-page.bulk-toggle.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx
  - frontend/tests/unit/workbench-v2/step4-share.test.tsx
  - internal/api/handlers/source_rule_handler_share_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_share_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - internal/datalink/sourcerule/protocol_address_test.go
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/studioV2RuleAutosave.test.ts
  - internal/datalink/migrator_test.go
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.planning.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/rule-tab-rail.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - internal/datalink/sourcerule/share_persistence_test.go
-->