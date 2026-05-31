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

When the runtime dashboard handoff callback is wired by the shell, clicking the secondary button SHALL navigate to the post-setup runtime dashboard route. If the shell can resolve a single handoff device from the current workbench state, the destination SHALL include `device_id=<resolved-id>`. If the shell cannot resolve a single handoff device, the destination SHALL fall back to the runtime dashboard route without `device_id`.

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

#### Scenario: Handoff includes resolved device id

- **GIVEN** the success card is visible
- **AND** the shell resolves `device_id=d-1` from the current workbench state
- **WHEN** the operator clicks `前往 Runtime Dashboard`
- **THEN** the system navigates to `/studio/runtime?device_id=d-1`

#### Scenario: Handoff falls back when device cannot be resolved

- **GIVEN** the success card is visible
- **AND** the shell cannot resolve a single handoff device from the current workbench state
- **WHEN** the operator clicks `前往 Runtime Dashboard`
- **THEN** the system navigates to `/studio/runtime`
- **AND** the handoff does not degrade to a no-op or console-only side effect


<!-- @trace
source: wire-step4-runtime-dashboard-handoff
updated: 2026-05-30
code:
  - AGENTS.md
  - progress.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - CLAUDE.md
  - frontend/src/App.tsx
  - docs/technical/studio-surface-inventory/changelog.sqlite
  - internal/api/handlers/runtime_stream_handler.go
  - cmd/studio_inventory_changelog/main.go
  - frontend/src/features/datalink/workbench-v2/shell/resolveRuntimeDashboardDevice.ts
  - docs/technical/studio-surface-inventory/index.html
  - frontend/src/i18n/config.ts
  - docs/technical/studio-surface-inventory/README.md
  - frontend/src/types/datalink.ts
  - internal/datalink/runtime/ingestor.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - docs/technical/studio-surface-inventory/inventory.css
  - frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - docs/technical/studio-surface-inventory/START_HERE.md
  - frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/runtime/service.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - docs/technical/studio-surface-inventory/backend-api-registry.md
  - docs/technical/studio-surface-inventory/inventory.js
  - docs/technical/studio-surface-inventory/studio-v2-runtime.md
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - findings.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - frontend/src/services/datalink.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - docs/technical/studio-surface-inventory/studio-mainline.md
  - docs/technical/studio-surface-inventory/test-tooling.md
  - internal/datalink/runtime/status.go
  - docs/technical/studio-surface-inventory/CURRENT_STATE.md
  - internal/datalink/runtime/stream.go
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json
  - docs/technical/studio-surface-inventory/context.json
  - frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - docs/technical/studio-surface-inventory/gap-roadmap.md
  - docs/technical/studio-surface-inventory/gateway-experiments.md
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/runtime/status_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/handlers/runtime_stream_handler_test.go
  - frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/runtime/stream_test.go
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

---
### Requirement: Valid-only database autosave

The Step 4 database workspace SHALL autosave only valid connector and database target content to the backend workspace.

#### Scenario: Valid connector edit saves immediately

- **GIVEN** the workspace already has persisted Step 4 database metadata
- **WHEN** the operator edits the connector into a valid state
- **THEN** the system saves that connector metadata immediately
- **AND** the persisted backend version reflects the latest valid values

#### Scenario: Invalid connector edit stays local

- **GIVEN** the connector already has a last successful persisted version
- **WHEN** the operator edits the connector into an invalid state
- **THEN** the system does not overwrite the persisted connector version
- **AND** the UI keeps the invalid local values visible with an unsaved marker


<!-- @trace
source: wire-studio-v2-step4-database-autosave
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/router.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/source_rule_handler.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - start.ps1
  - findings.md
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/api/handlers/mapping_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - .air.toml
  - internal/datalink/migrator.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - go.sum
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/workspace/service_database.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/workspace/memory_repository.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/schema/schema_device_models.go
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - scripts/start-log-utils.sh
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/hooks/datalink/index.ts
  - scripts/start-process-utils.sh
  - docs/goal.md
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/device/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/service_status.go
tests:
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_logger_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/device_handler_extended_test.go
-->

---
### Requirement: Per-target save isolation

The Step 4 database workspace SHALL isolate save success and failure per target row.

#### Scenario: One invalid target row does not block another valid target row

- **GIVEN** the workspace contains multiple database target rows
- **WHEN** one target row remains invalid while another target row is edited into a valid state
- **THEN** the valid target row still saves successfully
- **AND** the invalid target row remains local with its own error state


<!-- @trace
source: wire-studio-v2-step4-database-autosave
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/router.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/source_rule_handler.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - start.ps1
  - findings.md
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/api/handlers/mapping_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - .air.toml
  - internal/datalink/migrator.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - go.sum
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/workspace/service_database.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/workspace/memory_repository.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/schema/schema_device_models.go
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - scripts/start-log-utils.sh
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/hooks/datalink/index.ts
  - scripts/start-process-utils.sh
  - docs/goal.md
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/device/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/service_status.go
tests:
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_logger_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/device_handler_extended_test.go
-->

---
### Requirement: Database autosave remains pre-activation

Saving Step 4 database content SHALL NOT by itself start runtime collection.

#### Scenario: Autosave does not start runtime

- **WHEN** the operator successfully autosaves connector or target edits in Step 4
- **THEN** the persisted database content is updated
- **AND** runtime collection is still not started solely by that autosave

<!-- @trace
source: wire-studio-v2-step4-database-autosave
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/api/router.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/source_rule_handler.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - start.ps1
  - findings.md
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/api/handlers/mapping_handler.go
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - .air.toml
  - internal/datalink/migrator.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - internal/datalink/sourcerule/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - go.sum
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/workspace/service_database.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/workspace/memory_repository.go
  - internal/datalink/workspace/service.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - internal/datalink/schema/schema_device_models.go
  - task_plan.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/runtime/status.go
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/point_direct_reader.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - scripts/start-log-utils.sh
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - progress.md
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/hooks/datalink/index.ts
  - scripts/start-process-utils.sh
  - docs/goal.md
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/service_links.go
  - internal/datalink/sourcerule/errors.go
  - internal/datalink/device/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - gateway.db
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/mapping/service_crud.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/dbtarget/service.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/service_status.go
tests:
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - internal/api/router_logger_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/handlers/device_handler_extended_test.go
-->

---
### Requirement: First activation targets all eligible devices

The Step 4 final action SHALL activate all devices in the singleton v2 workspace that are currently valid, available, and not yet running.

#### Scenario: Activate all eligible devices

- **GIVEN** the singleton workspace contains multiple devices
- **AND** some devices are valid, available, and not yet running
- **WHEN** the operator invokes the final Step 4 action
- **THEN** the system activates every eligible device in that set
- **AND** already running devices are not activated a second time


<!-- @trace
source: replace-step4-with-first-activation
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/dbtarget/service.go
  - scripts/start-log-utils.sh
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - start.sh
  - internal/datalink/device/service_point_read.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - findings.md
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/api/handlers/studio_v2_workspace_handler.go
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/mapping/errors.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/workspace/service_database.go
  - .air.toml
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/runtime/status.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/point_direct_reader.go
  - task_plan.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/workspace/service_devices.go
  - go.sum
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - gateway.db
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - scripts/start-process-utils.sh
  - internal/datalink/workspace/service_activation.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/workspace/memory_repository.go
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - progress.md
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/migrator.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service.go
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/workspace/service_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
-->

---
### Requirement: Activation returns per-device results

The Step 4 final action SHALL report activation success and failure per device.

#### Scenario: Partial activation success

- **GIVEN** two eligible devices are activated and one succeeds while one fails
- **WHEN** the Step 4 action completes
- **THEN** the result surface reports one success and one failure separately
- **AND** the operator can tell which device failed and why


<!-- @trace
source: replace-step4-with-first-activation
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/dbtarget/service.go
  - scripts/start-log-utils.sh
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - start.sh
  - internal/datalink/device/service_point_read.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - findings.md
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/api/handlers/studio_v2_workspace_handler.go
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/mapping/errors.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/workspace/service_database.go
  - .air.toml
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/runtime/status.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/point_direct_reader.go
  - task_plan.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/workspace/service_devices.go
  - go.sum
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - gateway.db
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - scripts/start-process-utils.sh
  - internal/datalink/workspace/service_activation.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/workspace/memory_repository.go
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - progress.md
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/migrator.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service.go
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/workspace/service_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
-->

---
### Requirement: Navigation remains available after partial failure

The operator SHALL still be allowed to proceed to runtime after partial Step 4 activation failure.

#### Scenario: Proceed after partial activation failure

- **GIVEN** at least one device activated successfully and at least one device failed
- **WHEN** the operator chooses to continue
- **THEN** the system allows navigation to `/studio/runtime`
- **AND** the failed devices remain reported in Step 4

<!-- @trace
source: replace-step4-with-first-activation
updated: 2026-05-31
code:
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/dbtarget/service.go
  - scripts/start-log-utils.sh
  - frontend/src/App.tsx
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - start.sh
  - internal/datalink/device/service_point_read.go
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - findings.md
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - internal/api/handlers/studio_v2_workspace_handler.go
  - docs/goal.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - internal/datalink/mapping/errors.go
  - internal/datalink/sourcerule/service.go
  - internal/datalink/workspace/service_database.go
  - .air.toml
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - internal/datalink/runtime/status.go
  - internal/datalink/device/service_status.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/api/handlers/point_direct_reader.go
  - task_plan.md
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/workspace/service_devices.go
  - go.sum
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - gateway.db
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - scripts/start-process-utils.sh
  - internal/datalink/workspace/service_activation.go
  - internal/api/handlers/point_handler.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/types/studioV2Activation.ts
  - internal/datalink/dbtarget/types.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/services/studioV2Workspace.ts
  - internal/api/handlers/studio_v2_runtime_apply.go
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/sourcerule/errors.go
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/datalink/workspace/memory_repository.go
  - start.ps1
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - progress.md
  - internal/api/handlers/point_handler_polling.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/migrator.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - internal/datalink/device/service_crud.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/workspace/service.go
tests:
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/datalink/workspace/service_test.go
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/api/handlers/device_handler_extended_test.go
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
-->