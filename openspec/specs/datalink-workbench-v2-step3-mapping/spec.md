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

The right-side preview panel SHALL render four ordered pipeline steps for the currently selected mapping: `decode (raw)`, `scale (linear)`, `cast to {target_type}`, `final → Tag`. The preview MUST update reactively whenever the selected row's mapping (tag_key, scale, offset, target_type) changes. Raw values MUST come from the deterministic seed array `[243, 251, 1024, 985, 67, 542, 18, 1450]` indexed by selected row position modulo 8. The final value MUST be formatted per target type: `bool` → `"true"` / `"false"`, float types → `value.toFixed(2)`, integer types → `Math.round(value).toString()`, `string` → `String(value)`.

#### Scenario: Default preview with row 0 selected

- **GIVEN** `selectedIdx === 0`, the row 0 mapping has `scale=0.1`, `offset=0`, `target_type='float64'`, `tag_key='line01.temp.inlet'`
- **WHEN** the preview renders
- **THEN** decode step shows raw `243` with caption `register @ {address}, type={data_type}`
- **AND** scale step shows expression `243 × 0.1 + 0 = 24.300`
- **AND** cast step shows `cast(24.300 → float64)`
- **AND** final step shows `24.30` with unit and footer `→ tag.line01.temp.inlet`

#### Scenario: Float to int cast

- **GIVEN** the same row 0 but target_type is `int16`
- **WHEN** the preview renders
- **THEN** cast step shows `cast(24.300 → int16)`
- **AND** final step shows `24`

#### Scenario: Bool cast and format

- **GIVEN** the same row 0 but target_type is `bool` and scaled value is 24.3
- **WHEN** the preview renders
- **THEN** cast step shows `cast(24.300 → bool)`
- **AND** final step shows `true`

##### Example: format table

| Scaled value | Target type | Final display |
| ------------ | ----------- | ------------- |
| 24.3         | float64     | 24.30         |
| 24.3         | int16       | 24            |
| 24.3         | bool        | true          |
| 0            | bool        | false         |
| -1.7         | int16       | -2            |


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

The preview panel SHALL render a `POST /mappings` JSON payload preview below the pipeline steps. The payload MUST include `point_id`, `tag_id` (`tag.${tag_key}`), `transform_pipeline` (three steps: decode, scale with params `{scale, offset}`, cast with param `{to: target_type}`), and `enabled`. The preview MUST update on every change to the selected mapping. The preview MUST be display-only and MUST NOT initiate any real network request.

#### Scenario: Payload structure

- **GIVEN** the selected mapping has `tag_key='line01.temp.inlet'`, `scale=0.1`, `offset=0`, `target_type='float64'`, `enabled=true`, and the underlying point id is `pt-abc`
- **WHEN** the preview renders
- **THEN** the payload preview displays JSON with `point_id: "pt-abc"`, `tag_id: "tag.line01.temp.inlet"`, `transform_pipeline: [ {type:"decode"}, {type:"scale", params:{scale:0.1, offset:0}}, {type:"cast", params:{to:"float64"}} ]`, `enabled: true`

#### Scenario: Payload preview is non-issuing

- **WHEN** any value in the selected mapping changes
- **THEN** the payload preview updates
- **AND** no `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by the preview


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