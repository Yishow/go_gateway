# datalink-workbench-v2-settings Specification

## Purpose

TBD - created by archiving change 'datalink-workbench-v2-settings'. Update Purpose after archive.

## Requirements

### Requirement: Settings page layout

The Settings page SHALL render a header banner card, five SectionCards in vertical order (`Connector Pool`, `Timeseries Strategy` + `Scheduler Defaults` side by side, `Local Modbus Share`, `UI` + `API` + `Diagnostics` side by side), and a sticky bottom save bar. The right summary rail of the shell MUST be hidden while the view is `settings`. The save bar MUST remain visible during page scroll.

#### Scenario: Default render

- **WHEN** the operator switches the shell view to `settings`
- **THEN** the header banner card is rendered first with the slider icon and `系統設定` heading
- **AND** the five SectionCards appear in the specified order
- **AND** the right summary rail is not present in the DOM
- **AND** the save bar at the bottom uses `sticky bottom-4 z-10` positioning

#### Scenario: Save bar visible during scroll

- **WHEN** the operator scrolls the settings page downward
- **THEN** the save bar remains visible at the bottom of the viewport (sticky positioning intact)
- **AND** the save bar contains the `重設為預設` ghost button and `儲存所有設定` success button


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Connector pool CRUD

The Connector Pool section SHALL list every entry in `state.settings.connectors`, allow the operator to add a new connector via the `+ 新增連接器` aside button, edit per-row fields (name, kind, host, port, database, username, default_write_interval_seconds, enabled), test the connection, and remove a connector. Adding a connector MUST append a new entry with default name `新連線 N`, kind `postgres`, status `unknown`. Editing `kind`, `host`, or `port` MUST set the connector's status back to `unknown`. Removing a connector MUST be immediate (no confirmation in this change).

#### Scenario: Default connector

- **WHEN** the settings page mounts with the application's default settings
- **THEN** `state.settings.connectors` contains exactly one entry with id `conn-prod`, name `TimeSeries Prod`, kind `postgres`, host `tsdb.internal`, status `ready`

#### Scenario: Add connector

- **GIVEN** the pool has 1 connector
- **WHEN** the operator clicks `+ 新增連接器`
- **THEN** `state.settings.connectors.length === 2`
- **AND** the new entry has name `新連線 2`, kind `postgres`, status `unknown`, port `5432`, database `metrics`, username `gw_writer`, default_write_interval_seconds `5`

#### Scenario: Edit kind resets status

- **GIVEN** a connector has status `ready` and last_check_at set
- **WHEN** the operator changes the connector's `kind` from `postgres` to `mysql`
- **THEN** the connector's status becomes `unknown`
- **AND** the StatusChip displays `未測試` (or i18n equivalent)

#### Scenario: Remove connector

- **GIVEN** the pool has 2 connectors
- **WHEN** the operator clicks the remove button on connector 2
- **THEN** `state.settings.connectors.length === 1`
- **AND** no confirmation dialog is shown


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Connector mock test

The system SHALL provide a mock connector test triggered by the test button on each connector row. When triggered, the connector's status MUST transition to `testing` and the test button MUST become disabled with a spinning refresh icon. After exactly 900ms, the connector's status MUST become either `ready` (85% probability) with `last_check_at` set to the current ISO timestamp, or `unreachable` (15% probability) with `last_check_at` set and `last_check_error` set to `connection refused`. The test MUST NOT issue any real network request.

#### Scenario: Success path

- **GIVEN** `Math.random` is mocked to return `0.5` (i.e. > 0.15)
- **WHEN** the operator clicks the test button on a connector and `vi.useFakeTimers()` advances 900ms
- **THEN** the connector's status becomes `ready`
- **AND** `last_check_at` is set to a valid ISO timestamp
- **AND** the StatusChip displays `已就緒` (or i18n equivalent) in success tone

#### Scenario: Failure path

- **GIVEN** `Math.random` is mocked to return `0.05` (i.e. ≤ 0.15)
- **WHEN** the operator clicks the test button and 900ms elapses
- **THEN** the connector's status becomes `unreachable`
- **AND** `last_check_error` becomes `connection refused`
- **AND** the StatusChip displays `無法連線` in error tone
- **AND** the error text `connection refused` is rendered next to the chip

#### Scenario: No network call

- **WHEN** the operator clicks the test button
- **THEN** no `fetch`, `axios`, `useQuery`, or `useMutation` call is initiated by the settings page


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Timeseries strategy fields

The Timeseries section SHALL provide four fields with the following defaults and constraints:

- `write_precision`: select with options `second` and `millisecond`, default `millisecond`.
- `partition_interval`: select with options `daily`, `weekly`, `monthly`, default `daily`.
- `batch_size`: number, default `500`, hint `一次最多寫入的筆數`.
- `retention_days`: number, default `90`, hint `超過天數的資料將被歸檔`.

Changes MUST write through to `state.settings.timeseries` via `updateSettingsSection`.

#### Scenario: Default values

- **WHEN** the settings page mounts with default state
- **THEN** the timeseries section selects show `millisecond` and `daily`
- **AND** the number inputs show `500` and `90`

#### Scenario: Update write_precision

- **WHEN** the operator selects `second` from the write_precision select
- **THEN** `state.settings.timeseries.write_precision` becomes `second`


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Scheduler defaults fields

The Scheduler section SHALL provide five fields:

- `default_interval_ms`: number, default `1000`, hint `毫秒 (ms)`.
- `default_retry_count`: number, default `3`.
- `default_retry_delay_ms`: number, default `500`, hint `毫秒 (ms)`.
- `breaker_threshold`: number, default `10`, hint `連續錯誤次數`.
- `auto_start`: toggle, default `true`, label `開機時自動啟動 collector`.

Changes MUST write through to `state.settings.scheduler`.

#### Scenario: Toggle auto_start

- **WHEN** the operator toggles the `auto_start` switch from on to off
- **THEN** `state.settings.scheduler.auto_start` becomes `false`


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Local Modbus Share configuration

The Local Modbus Share section SHALL provide a master toggle in the aside slot and a 4-column form below: `bind_address` (text mono, default `0.0.0.0`, hint `0.0.0.0 = 全部介面`), `port` (number, default `5020`), `slave_id` (number, default `1`), `base_register` (number, default `40001`). The form MUST be hidden when the master toggle is off and a help text MUST be shown in its place. Toggling the master toggle MUST update `state.settings.modbus_share.enabled` and MUST propagate to Step 2's share layout computation via the existing `useShareLayout` selector.

#### Scenario: Toggle hides the form

- **GIVEN** the master toggle is on and the form is visible
- **WHEN** the operator toggles it off
- **THEN** the 4-column form is removed from the DOM
- **AND** the help text `啟用後可將所有 Tag 再次經由 Modbus TCP 提供給其他系統。` (or i18n equivalent) is shown

#### Scenario: base_register change affects Step 2

- **GIVEN** `state.settings.modbus_share.enabled === true` and `base_register === 40001`
- **WHEN** the operator changes `base_register` to `50001`
- **AND** the operator navigates to Step 2 with rules that have `share_start_register === null`
- **THEN** the Step 2 share layout for those rules starts at `50001` (the new base)


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: UI / API / Diagnostics options

The General settings SHALL be split across three side-by-side SectionCards:

- **UI**: `theme` select (`dark` / `light` / `auto`, default `dark`), `locale` select (`zh-TW` / `en`, default `zh-TW`), `addr_format` select (`modbus` / `hex` / `raw`, default `modbus`).
- **API**: `api_base` text (mono, default `http://localhost:8080`, hint `不含尾斜線`), `api_version` select (`v1` / `v2`, default `v1`), `timeout_seconds` number (default `30`, hint `秒`).
- **Diagnostics**: `log_level` select (`trace`/`debug`/`info`/`warn`/`error`, default `info`), `sse_heartbeat_seconds` number (default `15`, hint `秒`), `enable_debug_panel` toggle (default `false`), `enable_audit_log` toggle (default `true`).

Changes MUST write through to `state.settings.general`. The theme select MUST NOT actually toggle the UI theme during this change (dark-only); writing the value is sufficient.

#### Scenario: Change theme stored but UI stays dark

- **WHEN** the operator selects `light` from the theme select
- **THEN** `state.settings.general.theme` becomes `light`
- **AND** the v2 page subtree continues to use dark tokens (the visual theme does not change)

#### Scenario: Toggle audit log

- **WHEN** the operator toggles `enable_audit_log` from on to off
- **THEN** `state.settings.general.enable_audit_log` becomes `false`


<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Save bar actions

The sticky save bar SHALL contain a left-side info text equivalent to `設定會立即套用，並於下次重啟後生效。`, a `重設為預設` ghost button, and a `儲存所有設定` success button. The `重設為預設` button MUST trigger a native `window.confirm`; on confirmation, `state.settings` MUST be replaced with `DEFAULT_SETTINGS`. The `儲存所有設定` button MUST be a no-op in this change and MUST log a warning (`console.warn('pending backend-wiring')`) when clicked; it MUST NOT issue any network request.

#### Scenario: Reset with confirm

- **GIVEN** the operator has modified several settings fields
- **WHEN** the operator clicks `重設為預設` and confirms the native dialog
- **THEN** `state.settings` equals `DEFAULT_SETTINGS`

#### Scenario: Reset cancelled

- **GIVEN** the operator has modified settings fields
- **WHEN** the operator clicks `重設為預設` and cancels the dialog
- **THEN** `state.settings` remains unchanged

#### Scenario: Save is no-op

- **WHEN** the operator clicks `儲存所有設定`
- **THEN** `console.warn` is called once with a message indicating backend-wiring is pending
- **AND** no `fetch` or other network call is initiated
- **AND** `state.settings` remains unchanged

<!-- @trace
source: datalink-workbench-v2-settings
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
tests:
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
-->

---
### Requirement: Backend-backed V2 settings surface

The `/studio/v2/settings` surface SHALL boot from real backend state instead of demo-only local defaults.

#### Scenario: Settings boot from backend

- **WHEN** the operator opens `/studio/v2/settings`
- **THEN** the page loads settings and connector data from backend APIs
- **AND** the page does not silently treat local defaults as a completed backend load


<!-- @trace
source: wire-studio-v2-settings-backend
updated: 2026-05-31
code:
  - scripts/start-process-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/services/studioV2Rules.ts
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/sourcerule/service.go
  - progress.md
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - go.sum
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/device/errors.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/runtime/status.go
  - .air.toml
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - gateway.db
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/mapping_handler.go
  - internal/api/router.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - start.ps1
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/mapping/service_crud.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/types/datalink.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - docs/goal.md
  - internal/datalink/db.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
tests:
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - cmd/test_ui/main.go
-->

---
### Requirement: Real connector pool operations

The connector pool SHALL use real backend CRUD and test APIs.

#### Scenario: Real connector test

- **GIVEN** a connector row exists in the pool
- **WHEN** the operator clicks the test button
- **THEN** the system sends a real backend connector test request
- **AND** the result shown in the UI comes from backend response, not Math.random or a mock timer


<!-- @trace
source: wire-studio-v2-settings-backend
updated: 2026-05-31
code:
  - scripts/start-process-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/services/studioV2Rules.ts
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/sourcerule/service.go
  - progress.md
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - go.sum
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/device/errors.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/runtime/status.go
  - .air.toml
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - gateway.db
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/mapping_handler.go
  - internal/api/router.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - start.ps1
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/mapping/service_crud.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/types/datalink.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - docs/goal.md
  - internal/datalink/db.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
tests:
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - cmd/test_ui/main.go
-->

---
### Requirement: Save bar persists settings

The `儲存所有設定` action SHALL persist settings through backend APIs.

#### Scenario: Save settings

- **GIVEN** the operator has changed one or more settings fields
- **WHEN** the operator clicks `儲存所有設定`
- **THEN** the system writes the changed settings to backend
- **AND** save failure is reported as an actionable error instead of a warning-only noop

<!-- @trace
source: wire-studio-v2-settings-backend
updated: 2026-05-31
code:
  - scripts/start-process-utils.sh
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - frontend/src/services/studioV2Rules.ts
  - scripts/start-log-utils.sh
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - task_plan.md
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - internal/datalink/device/service_crud.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - internal/datalink/sourcerule/service.go
  - progress.md
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - go.sum
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/migrator.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/datalink/workspace/service_devices.go
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - internal/api/handlers/point_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/schema/schema_device_models.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/runtime_handler.go
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/workspace/service.go
  - internal/datalink/device/errors.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - internal/datalink/device/service_status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - internal/datalink/runtime/status.go
  - .air.toml
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - gateway.db
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/mapping_handler.go
  - internal/api/router.go
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - start.ps1
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_database.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/sourcerule/errors.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/mapping/service_crud.go
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - internal/datalink/device/availability.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/types/studioV2Availability.ts
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - start.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/types/datalink.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - internal/datalink/device/sql_repo.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - docs/goal.md
  - internal/datalink/db.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
tests:
  - internal/api/router_logger_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - internal/api/router_runtime_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/device/availability_test.go
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/workspace/service_test.go
  - internal/datalink/device/service_crud_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - cmd/test_ui/main.go
-->