# datalink-workbench-v2-step4-database Specification

## Purpose

TBD - created by archiving change 'datalink-workbench-v2-step4-database'. Update Purpose after archive.

## Requirements

### Requirement: Connector configuration form

The Step 4 connector section SHALL render a card-style selector for four database kinds (`sqlite`, `postgres`, `mysql`, `sqlserver`) and a configuration form with fields: connection name, host (mono), port (number), database (mono), schema (mono), table (mono), username (mono), write_mode (radio: `insert` / `upsert`), write_interval_seconds (number with `秒` unit suffix), and timestamp_column (mono, default `ts`). Switching the kind MUST update `state.db.connector.kind` and trigger re-evaluation of the auto-assign algorithm using the new kind's column set.

#### Scenario: Default connector on first render

- **WHEN** Step 4 mounts with no prior connector state
- **THEN** the connector is initialized with name `TimeSeries Prod`, kind `postgres`, host `tsdb.internal`, port `5432`, database `gateway_metrics`, username `gw_writer`, schema `public`, table `sensor_readings`, write_mode `insert`, write_interval_seconds `5`, timestamp_column `ts`, status `ready`

#### Scenario: Kind switch triggers auto-assign

- **GIVEN** the workspace has 8 enabled points with auto-assigned targets and kind is `postgres`
- **WHEN** the operator clicks the `mysql` kind card
- **THEN** `state.db.connector.kind` becomes `mysql`
- **AND** the auto-assign algorithm re-runs using the mysql column set
- **AND** existing target entries that match a still-valid column are preserved

#### Scenario: Write strategy persistence

- **WHEN** the operator selects the `upsert` radio
- **THEN** `state.db.connector.write_mode` becomes `upsert`
- **AND** the radio is rendered as selected


<!-- @trace
source: datalink-workbench-v2-step4-database
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - task_plan.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
tests:
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->

---
### Requirement: Tag-to-column auto-assignment

The system SHALL provide an auto-assign function `autoAssignTargets(enabledPoints, mappings, columnNames, existingTargets)` that returns a `Record<pointId, DbTarget>`. The function MUST follow this priority order per point: (1) if `existingTargets[pointId]` exists, preserve it; (2) compute `tagShort = mapping.tag_key.split('.').pop()` and find an unused column that equals `tagShort`, ends with `_${tagShort}`, or starts with `${tagShort}_`; (3) fall back to `columnNames[i % length]` if not used; (4) otherwise find the first unused column. The function MUST mark used columns to avoid auto-generated conflicts when enough columns exist.

#### Scenario: Eight points to nine columns

- **GIVEN** 8 enabled points with mappings whose tag_keys end in `temp.inlet`, `temp.outlet`, `pressure.main`, `pressure.sub`, `flow.q1`, `humidity.amb`, `vibration.motor`, `motor.rpm`, and the postgres sample table has 8 non-primary-key columns (temp_in_c, temp_out_c, pressure_main_kpa, pressure_sub_kpa, flow_lpm, humidity_pct, vibration_mms, motor_rpm)
- **WHEN** `autoAssignTargets` runs with empty existingTargets
- **THEN** each point is mapped to a column via the endsWith / startsWith / index-fallback chain
- **AND** no two points share the same column

#### Scenario: Existing target preserved

- **GIVEN** existingTargets contains `{ pt-X: { column_name: 'flow_lpm', enabled: true, tag_id: 'tag.foo' } }`
- **WHEN** `autoAssignTargets` runs
- **THEN** the output for `pt-X` equals the existing entry (column_name unchanged)
- **AND** `flow_lpm` is treated as used for subsequent points

#### Scenario: Fewer columns than points causes wrap and conflict

- **GIVEN** 8 enabled points but only 4 columns
- **WHEN** `autoAssignTargets` runs with empty existingTargets and no exact matches
- **THEN** each point still receives a target (via `i % length`)
- **AND** at least two points share the same column (downstream conflict detection flags this)


<!-- @trace
source: datalink-workbench-v2-step4-database
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - task_plan.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
tests:
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->

---
### Requirement: Column conflict detection

The system SHALL detect column conflicts in `state.db.targets`: when two or more enabled targets share the same `column_name`, the column MUST be flagged as conflicting. Conflicting target rows MUST render the column select with a red border (e.g. `border-red-500/50`), an `alert` icon adjacent to the select, and the merged table footer MUST display a red banner stating that multiple Tags write to the same column. The commit button MUST be disabled when any conflict exists.

#### Scenario: Two enabled targets collide on temp_in_c

- **GIVEN** target for point A has `column_name='temp_in_c'`, enabled=true; target for point B has `column_name='temp_in_c'`, enabled=true
- **WHEN** the merged table renders
- **THEN** both rows show a red border on the column select and an alert icon
- **AND** the table footer shows a red banner with text equivalent to "偵測到多個 Tag 寫入同一資料表欄位。請調整以避免覆寫。"
- **AND** the commit button (`提交並啟動排程器`) is disabled

#### Scenario: Disabling one side resolves conflict

- **GIVEN** the same conflict as above
- **WHEN** the operator toggles the enabled switch off for one of the conflicting rows
- **THEN** the conflict is cleared
- **AND** the red border and banner are removed
- **AND** the commit button becomes enabled (assuming no other conflicts and at least one enabled target)


<!-- @trace
source: datalink-workbench-v2-step4-database
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - task_plan.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
tests:
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->

---
### Requirement: Commit sequence and animation

The commit button click SHALL initiate a deterministic 10-step animation. Each step MUST be appended to `state.commit.logs` exactly 280ms after the previous step. The 10 steps MUST be (in order):

1. `POST /devices × {deviceCount}` — device names and protocols
2. `POST /devices/:id/activate × {deviceCount}` — `draft → active`
3. `POST /source-rules × {enabledRuleCount}` — rule names and counts
4. `POST /points × {enabledPointCount}` — `bulk create`
5. `POST /polling-groups` — `快速輪詢 1s, enabled` (or i18n equivalent)
6. `POST /tags × {enabledPointCount}` — `register tag keys`
7. `POST /mappings × {enabledPointCount}` — `point ↔ tag, scale pipeline`
8. `POST /db-connectors/:id/test` — `{kind} {host}:{port}`
9. `POST /db-targets × {enabledTargetCount}` — `→ {schema}.{table}`
10. `POST /scheduler/start` — `collectors started`

During the animation, the commit button MUST be disabled. After the 10th log, `state.committed` MUST become `true`, `state.commit.status` MUST become `'success'`, `state.commit.finished_at` MUST be set to the current ISO timestamp, and the CommitSuccessCard MUST be rendered.

#### Scenario: Animation timing

- **GIVEN** the operator clicks the commit button at t=0
- **WHEN** `vi.useFakeTimers()` advances 10 × 280ms
- **THEN** `state.commit.logs.length === 10`
- **AND** the logs appear in the exact order specified
- **AND** `state.committed === true`
- **AND** `state.commit.status === 'success'`

#### Scenario: No backend call issued

- **WHEN** the commit animation runs end-to-end
- **THEN** no `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by Step 4 or its children
- **AND** the log labels reflect the API names that the future backend-wiring change will eventually call, but the strings are display-only in this change

#### Scenario: Resume on remount mid-commit

- **GIVEN** the operator clicks commit and three logs have been appended (status `running`)
- **WHEN** the operator navigates to Step 3 and then back to Step 4 within the next 1s
- **THEN** `state.commit.logs.length` remains 3 at the moment of remount
- **AND** the effect continues ticking from log 4 until all 10 logs are present


<!-- @trace
source: datalink-workbench-v2-step4-database
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - task_plan.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
tests:
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->

---
### Requirement: Commit completion card

After `state.committed === true`, the right column SHALL render an emerald-tone success card containing: a circular check icon, a heading equivalent to `設定已套用 · 開始收集資料`, a subline equivalent to `Scheduler 已啟動 · 第一筆資料預計在 ~{write_interval_seconds}s 後寫入`, and a secondary button `前往 Runtime Dashboard` that invokes the `onCommit` callback. The CommitSummary and CommitProgress views MUST be hidden when the success card is shown.

#### Scenario: Success card visible

- **GIVEN** `state.committed === true` and `state.commit.status === 'success'`
- **WHEN** the right column renders
- **THEN** the emerald success card is present
- **AND** CommitSummary (the 5-row summary) is NOT in the DOM
- **AND** the secondary button text equals `前往 Runtime Dashboard` (or i18n equivalent)
- **AND** clicking the button invokes `onCommit` exactly once

#### Scenario: Subline reflects write_interval_seconds

- **GIVEN** the success card is visible and `state.db.connector.write_interval_seconds === 10`
- **WHEN** the subline renders
- **THEN** the subline text contains `~10s`


<!-- @trace
source: datalink-workbench-v2-step4-database
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - task_plan.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
tests:
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->

---
### Requirement: Read-only form after commit

When `state.committed === true`, all interactive controls in Step 4 (kind selector cards, connector form inputs, write strategy radios, target mapping column selects, and target enabled toggles) MUST be disabled. The commit button MUST be replaced by the CommitSuccessCard rather than being rendered as a disabled button.

#### Scenario: Inputs disabled after commit

- **GIVEN** `state.committed === true`
- **WHEN** the operator attempts to type into the host input
- **THEN** the input is rendered with `disabled` attribute
- **AND** typing does not produce any change to `state.db.connector.host`

#### Scenario: Column select disabled after commit

- **GIVEN** `state.committed === true`
- **WHEN** the operator inspects a target mapping row's column select
- **THEN** the select is disabled
- **AND** clicking it does not open any options menu (browser-native disabled behavior)

#### Scenario: Kind selector locked after commit

- **GIVEN** `state.committed === true` and kind is `postgres`
- **WHEN** the operator clicks the `mysql` kind card
- **THEN** the kind remains `postgres`
- **AND** the card does not show focus / active styling

<!-- @trace
source: datalink-workbench-v2-step4-database
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - task_plan.md
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
tests:
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
-->