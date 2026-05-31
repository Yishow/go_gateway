# datalink-workbench-v2-shell Specification

## Purpose

TBD - created by archiving change 'datalink-workbench-v2-step1-device'. Update Purpose after archive.

## Requirements

### Requirement: Scheduler indicator reflects committed state

The shell top bar SHALL render a scheduler status indicator whose tone reflects `state.committed`. When `state.committed === false` the indicator MUST render with an amber dot and the text `scheduler idle` (or i18n equivalent). When `state.committed === true` the indicator MUST render with an emerald dot, a `pulse-dot` animation, and the text `scheduler running` (or i18n equivalent). The transition MUST occur on the same render tick as the reducer applies the `completeCommit` action.

#### Scenario: Idle indicator before commit

- **WHEN** the v2 page mounts with `state.committed === false`
- **THEN** the top bar shows `scheduler idle` text
- **AND** the indicator dot uses an amber tone with no `pulse-dot` class

#### Scenario: Running indicator after commit

- **WHEN** the Step 4 commit animation completes and `state.committed` becomes `true`
- **THEN** the top bar shows `scheduler running` text
- **AND** the indicator dot uses an emerald tone with the `pulse-dot` animation class

#### Scenario: Indicator persists across step navigation

- **GIVEN** `state.committed === true`
- **WHEN** the operator navigates from Step 4 to Step 1 via the StepRail
- **THEN** the top bar continues to show `scheduler running`

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
### Requirement: Coexisting v2 workbench route

The system SHALL expose a parallel workbench route at `/studio/v2` that loads the Datalink Workbench v2 shell, while `/studio` continues to load the legacy `DatalinkWorkbenchPage` without modification. The two routes MUST coexist; navigating to either route MUST NOT redirect to the other.

#### Scenario: Operator opens /studio/v2

- **WHEN** the operator navigates to `/studio/v2`
- **THEN** the system loads the Workbench v2 shell with top bar, collapsible step rail, central step content, summary rail, and tweaks panel
- **AND** does not redirect to `/studio`

#### Scenario: Legacy /studio remains unchanged

- **WHEN** the operator navigates to `/studio`
- **THEN** the system loads the existing `DatalinkWorkbenchPage` exactly as before this change
- **AND** does not redirect to `/studio/v2`

#### Scenario: Direct deep link to v2

- **WHEN** the operator opens a fresh browser tab with URL ending in `/studio/v2`
- **THEN** the v2 shell mounts with the default state: Step 1 highlighted, view mode `flow`, no completed steps, summary rail visible on viewports ≥ 1280px


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Shell layout regions

The Workbench v2 shell SHALL render five persistent regions in a single page: a sticky top bar, a collapsible left step rail, a central step content surface, an optional right summary rail, and a floating tweaks panel. The shell MUST remain functional when the right summary rail is hidden, when the step rail is collapsed, or when the tweaks panel is closed.

#### Scenario: Default render of all regions

- **WHEN** the v2 page mounts at viewport width 1440px and the user has not changed any layout preference
- **THEN** the top bar is sticky at the top of the viewport
- **AND** the left step rail is expanded with width 232px showing four steps and the settings entry
- **AND** the central step content shows the Step 1 header (`STEP 01 / 04`, step title, progress bar) and the Step 1 placeholder
- **AND** the right summary rail is visible with width 240px showing device / rule / mapping / database summary cards
- **AND** the tweaks panel is hidden by default

#### Scenario: Step rail collapsed state

- **WHEN** the user clicks the collapse button in the top bar
- **THEN** the left step rail width transitions to 64px showing only step icons
- **AND** each icon shows a tooltip with the step title on hover
- **AND** the central step content surface expands to fill the released horizontal space

##### Example: collapsed widths

| State     | Step rail width | Has expanded labels | Tooltip on hover |
| --------- | --------------- | ------------------- | ---------------- |
| Expanded  | 232px           | Yes                 | No               |
| Collapsed | 64px            | No                  | Yes              |

#### Scenario: Summary rail hidden at narrow viewports

- **WHEN** the viewport width is less than 1280px
- **THEN** the right summary rail is hidden
- **AND** the central step content uses the released horizontal space
- **AND** all other shell regions remain functional


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Step rail navigation

The step rail SHALL expose four ordered steps (`新增裝置`, `接入規則`, `點位映射`, `儲存資料庫`) plus a settings entry below a divider. Each step entry MUST surface its step number, title, subtitle, and an icon. Reachability MUST follow the rule: step 1 is always reachable, step N (N>1) is reachable when step N-1 has been completed or when step N is the current step.

#### Scenario: Initial reachability

- **WHEN** the v2 page mounts with no completed steps
- **THEN** step 1 is reachable and current
- **AND** steps 2, 3, and 4 are not reachable and rendered with `cursor: not-allowed` and reduced opacity

#### Scenario: Step completion unlocks the next step

- **WHEN** the user has completed step 1 (step 1 ID exists in `completed` state)
- **THEN** step 2 becomes reachable
- **AND** clicking step 2 in the rail switches the current step to 2 and updates the breadcrumb in the top bar to `Datalink › Workbench › 接入規則`

#### Scenario: Settings entry switches view mode

- **WHEN** the user clicks the settings entry in the step rail
- **THEN** the view state changes from `flow` to `settings`
- **AND** the central step content surface renders the Settings placeholder instead of any step placeholder
- **AND** the right summary rail is hidden
- **AND** the top bar breadcrumb shows `Datalink › Workbench › 設定`

#### Scenario: Returning from settings to flow

- **WHEN** the view state is `settings` and the user clicks any reachable step in the step rail
- **THEN** the view state changes back to `flow`
- **AND** the current step updates to the clicked step ID
- **AND** the right summary rail returns to visible (on viewports ≥ 1280px)


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Keyboard shortcut for step rail collapse

The shell SHALL bind a keyboard shortcut (`Cmd+B` on macOS, `Ctrl+B` on Windows/Linux) to toggle the step rail collapsed state. The shortcut MUST only fire while the v2 page is mounted, MUST call `preventDefault` on the original keyboard event, and MUST NOT fire when focus is inside an editable form control.

#### Scenario: Shortcut toggles step rail

- **WHEN** the user presses `Cmd+B` while the step rail is expanded and focus is on the document body
- **THEN** the step rail collapses to width 64px
- **AND** the keydown event default action is prevented

#### Scenario: Shortcut suppressed inside text input

- **WHEN** the user focuses an `<input>` element in the central step content and presses `Cmd+B`
- **THEN** the step rail collapsed state does NOT change
- **AND** the keydown event default action proceeds normally for the input


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Design tokens scope

The v2 shell SHALL apply Dark Industrial Telemetry design tokens (background `#0b1220`, surface `slate-900` with 60% opacity and backdrop blur, primary `blue-500/600`, success `emerald-500`, warning `amber-500`, error `red-500`, 24px grid background, Inter and JetBrains Mono fonts) only within the v2 subtree. Tokens MUST NOT leak to the `/studio` route or any other page in the application.

#### Scenario: Tokens scoped to v2 subtree

- **WHEN** the user opens `/studio` after visiting `/studio/v2` in the same browser session
- **THEN** the `/studio` page renders with its pre-existing styling
- **AND** `#0b1220` background and 24px grid background do NOT appear on `/studio`

#### Scenario: Tokens applied on v2

- **WHEN** the v2 page is mounted
- **THEN** the page root element has a data attribute identifying the v2 subtree (for example `data-workbench-v2="true"`)
- **AND** the background uses the canvas gradient + 24px grid pattern
- **AND** the default UI font-family resolves to Inter (with system fallback) and mono regions resolve to JetBrains Mono


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Shared component library for v2

The v2 feature module SHALL provide a TypeScript component library consisting of `Icon`, `Button`, `Field`, `Input`, `Select`, `Textarea`, `Toggle`, `StatusChip`, and `SectionCard`. Each component MUST export a named TypeScript props interface, MUST NOT use the `any` type, and MUST support keyboard interaction and accessible labels.

#### Scenario: Icon renders without emoji

- **WHEN** a component renders `<Icon name="device" />`
- **THEN** the rendered output is an inline SVG element with `aria-hidden="true"`
- **AND** no emoji character is used as the visual glyph

#### Scenario: Toggle exposes ARIA switch role

- **WHEN** a `Toggle` component is rendered with `checked={false}`
- **THEN** the rendered button has `role="switch"` and `aria-checked="false"`
- **AND** clicking the button calls `onChange(true)`

#### Scenario: SectionCard renders header and content

- **WHEN** a `SectionCard` is rendered with `title="設備列表"`, `subtitle="可同時設定多個設備"`, an icon prop, an aside prop, and children
- **THEN** the rendered output contains a header region with title, subtitle, icon, and aside
- **AND** a content region containing the children


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Placeholder step and settings surfaces

The v2 shell SHALL render placeholder content for each of the four steps and the settings page during the shell-only delivery phase. Each placeholder MUST display a message identifying which step it represents and that full content is coming in a subsequent change. Placeholders MUST NOT make backend API calls, MUST NOT consume `useQuery` or `useMutation` hooks, and MUST NOT block the `onContinue` action for development purposes.

#### Scenario: Step 1 placeholder renders identifying content

- **WHEN** the v2 page mounts with current step 1
- **THEN** the central step content displays text identifying it as Step 1 (`新增裝置`)
- **AND** the placeholder text indicates that full Step 1 functionality is delivered by a subsequent change
- **AND** no network request is initiated by the placeholder

#### Scenario: Settings placeholder renders identifying content

- **WHEN** the view state is `settings`
- **THEN** the central content displays text identifying it as the Settings page
- **AND** the placeholder text indicates that full Settings functionality is delivered by a subsequent change


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Tweaks panel availability

The shell SHALL provide a floating tweaks panel that is hidden by default. The panel MUST be visible only when the application is running in development mode (`import.meta.env.DEV === true`) or when `localStorage.WBV2_TWEAKS === '1'`. The panel MUST control the `sidebarCollapsed` and `showSummaryRail` layout preferences and persist them to `localStorage`.

#### Scenario: Panel hidden in production by default

- **WHEN** the v2 page mounts with `import.meta.env.DEV === false` and `localStorage.WBV2_TWEAKS` is unset
- **THEN** the tweaks panel is not visible in the DOM
- **AND** the shell still renders correctly

#### Scenario: Panel visible in development

- **WHEN** the v2 page mounts with `import.meta.env.DEV === true`
- **THEN** the tweaks panel is rendered as a floating element in the viewport
- **AND** toggling `收合側邊欄` inside the panel updates the `sidebarCollapsed` state and persists `"sidebarCollapsed":true` to `localStorage`

#### Scenario: localStorage fallback when unavailable

- **WHEN** `localStorage.setItem` throws an exception (for example in a private browsing context) and the user toggles a layout preference in the tweaks panel
- **THEN** the UI state still updates in memory for the current session
- **AND** the shell continues to function without crashing


<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: i18n namespace for v2

The v2 shell SHALL load all user-facing strings from a dedicated i18n namespace `workbench-v2`. The namespace MUST be available in at least `zh-TW` and `en`. All button labels, breadcrumb fragments, step titles, step subtitles, placeholder copy, and summary card labels MUST resolve through `useTranslation('workbench-v2')`.

#### Scenario: zh-TW locale shows traditional Chinese microcopy

- **WHEN** the locale is `zh-TW` and the v2 page mounts
- **THEN** the step rail shows step titles `新增裝置`, `接入規則`, `點位映射`, `儲存資料庫` and settings entry `設定`
- **AND** the save draft button shows `儲存草稿`

#### Scenario: en locale shows English microcopy

- **WHEN** the locale is `en` and the v2 page mounts
- **THEN** the step rail shows step titles in English (for example `Add Device`, `Source Rule`, `Point Mapping`, `Save to Database`) and the settings entry shows `Settings`
- **AND** the save draft button shows `Save Draft`

<!-- @trace
source: datalink-workbench-v2-shell
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - .line-limit-ignore
  - frontend/src/features/datalink/workbench-v2/settings/UiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsHeader.tsx
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ModbusShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SchedulerSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/index.ts
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/settings/DiagnosticsSection.tsx
  - frontend/src/features/datalink/workbench-v2/settings/TimeseriesSection.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/main.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/state/types-step4.test-d.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ApiSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/commitLog.ts
  - frontend/src/features/datalink/workbench-v2/settings/SaveBar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
tests:
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - frontend/tests/unit/workbench-v2/types-step3.test-d.ts
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/types.test-d.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/transformPipeline.test.ts
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/settingsDefaults.test.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/commitLog.test.ts
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/workbench-v2/reducer-settings.test.ts
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/dbSchemas.test.ts
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/autoAssignTargets.test.ts
-->

---
### Requirement: Default v2 entry route

The system SHALL treat `/studio/v2` as the default guided entry for the datalink setup flow.

#### Scenario: Root path enters v2

- **WHEN** an operator opens `/`
- **THEN** the system redirects to `/studio/v2`
- **AND** the first rendered setup surface is the v2 shell rather than the legacy `/studio` page

#### Scenario: Unknown path falls back to v2

- **WHEN** an operator opens an unknown route that does not match any explicit surface
- **THEN** the system redirects to `/studio/v2`
- **AND** the fallback does not land on `/studio`

#### Scenario: Direct legacy mainline remains available

- **WHEN** an operator opens `/studio`
- **THEN** the system still renders the legacy workbench page
- **AND** the route does not auto-redirect to `/studio/v2`


<!-- @trace
source: make-studio-v2-default-entry
updated: 2026-05-31
code:
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - start.ps1
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/device/service_status.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/datalink/workspace/service_database.go
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/migrator.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/dbtarget/service.go
  - internal/datalink/runtime/status.go
  - frontend/src/types/studioV2Availability.ts
  - go.sum
  - progress.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - .air.toml
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - gateway.db
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - task_plan.md
  - internal/datalink/sourcerule/errors.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - start.sh
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - scripts/start-log-utils.sh
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - scripts/start-process-utils.sh
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - docs/goal.md
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/device/availability.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/App.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/workspace/service.go
  - internal/api/handlers/point_handler.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/router.go
  - internal/datalink/device/service_point_read.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/device/service_crud.go
tests:
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - cmd/test_ui/main.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/router_logger_test.go
-->

---
### Requirement: Legacy landing paths converge on /studio/v2

The system SHALL route generic legacy datalink landing paths to `/studio/v2` while preserving task-specific deep-link behavior outside this change.

#### Scenario: Generic legacy datalink root lands on v2

- **WHEN** an operator opens `/datalink`
- **THEN** the system redirects to `/studio/v2`
- **AND** the operator does not first land on `/studio`

#### Scenario: Generic legacy workbench path lands on v2

- **WHEN** an operator opens `/datalink/workbench`
- **THEN** the system redirects to `/studio/v2`
- **AND** the operator does not first land on `/studio`

<!-- @trace
source: make-studio-v2-default-entry
updated: 2026-05-31
code:
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - start.ps1
  - go.mod
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/device/service_status.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - internal/datalink/workspace/service_database.go
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/migrator.go
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/datalink/dbtarget/service.go
  - internal/datalink/runtime/status.go
  - frontend/src/types/studioV2Availability.ts
  - go.sum
  - progress.md
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/api/handlers/mapping_handler.go
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - .air.toml
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - gateway.db
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/hooks/datalink/index.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/hooks/datalink/useSettings.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - internal/api/handlers/studio_v2_workspace_handler.go
  - task_plan.md
  - internal/datalink/sourcerule/errors.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - start.sh
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - scripts/start-log-utils.sh
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - scripts/start-process-utils.sh
  - frontend/src/services/studioV2Workspace.ts
  - internal/datalink/workspace/service_devices.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - docs/goal.md
  - frontend/src/types/studioV2RuntimeApply.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/datalink/device/availability.go
  - internal/datalink/workspace/sql_repository.go
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - internal/datalink/db.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/types/datalink.ts
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - frontend/src/App.tsx
  - internal/datalink/dbtarget/types.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/types/studioV2Activation.ts
  - internal/api/handlers/runtime_handler.go
  - internal/datalink/mapping/service_crud.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - internal/datalink/workspace/service.go
  - internal/api/handlers/point_handler.go
  - internal/datalink/sourcerule/activation_readiness.go
  - internal/api/handlers/point_handler_polling.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/router.go
  - internal/datalink/device/service_point_read.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - internal/datalink/schema/schema_device_models.go
  - internal/datalink/workspace/service_activation.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service.go
  - internal/datalink/device/sql_repo.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/services/studioV2Rules.ts
  - internal/datalink/workspace/memory_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/device/service_crud.go
tests:
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - cmd/test_ui/main.go
  - internal/api/handlers/runtime_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/device/availability_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - internal/datalink/workspace/service_device_order_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/router_runtime_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - internal/api/router_studio_v2_workspace_mappings_test.go
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - internal/api/router_studio_v2_workspace_devices_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - internal/api/router_logger_test.go
-->