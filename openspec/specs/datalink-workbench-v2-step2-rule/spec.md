# datalink-workbench-v2-step2-rule Specification

## Purpose

TBD - created by archiving change 'datalink-workbench-v2-step2-rule'. Update Purpose after archive.

## Requirements

### Requirement: Multi-rule tab management

The Step 2 source-rule workspace SHALL render a horizontally scrollable tab rail listing all rules in `state.rules`. The tab rail MUST allow the operator to add a new rule, select an existing rule for editing, rename a rule inline, toggle a rule's `enabled` flag, and remove a non-only rule. Each rule tab MUST display the rule color dot, name, `{start_address} · {data_type}` summary, enabled point count, and (when the workspace has 2+ devices) the owning device's color dot and name. When creating a new rule, the system SHALL determine the default start address based on the target device's communication protocol (`40001` for Modbus TCP/RTU/UDP, and `D0` for FATEK FBs / Mitsubishi MC 3E).

#### Scenario: Default rule on first render
- **WHEN** the operator opens Step 2 with default state
- **THEN** the tab rail shows exactly one rule with start address corresponding to the default device protocol
- **AND** for a Modbus device, start address is `40001` with data type `int16` and enabled point count `8/8`
- **AND** for an MC 3E or FATEK device, start address defaults to `D0`

#### Scenario: Add rule
- **WHEN** the operator clicks `+ 新增規則`
- **THEN** a new rule is appended to `state.rules` with name `規則 N` (N is next ordinal), default count 8, naming_prefix `BLOCK{N}_`, and device_id defaulting to the currently selected rule's device_id or `state.devices[0].id`
- **AND** `start_address` is populated using `getDefaultPlannerStartAddress(device.protocol)` (`40001` for Modbus, `D0` for FATEK/MC 3E)
- **AND** the new rule is selected

#### Scenario: Remove rule
- **WHEN** the operator hovers a rule tab and clicks the close icon while `state.rules.length >= 2`
- **THEN** the rule is removed from `state.rules`
- **AND** if the removed rule was selected, the first remaining rule becomes selected
- **AND** the close icon is hidden when `state.rules.length === 1`

#### Scenario: Multi-device device-row visibility
- **WHEN** `state.devices.length === 1`
- **THEN** the rule tab does NOT render the device color dot or device name row

#### Scenario: Multi-device device-row appears with 2+ devices
- **WHEN** `state.devices.length >= 2`
- **THEN** every rule tab renders a third row showing the owning device's color dot and name


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

---
### Requirement: Rule editor with linked reset

The Step 2 rule editor SHALL render input controls for rule properties (name, owning device, start address, point count, data type, naming prefix, scale multiplier, scale offset, data format, skipped addresses, and Modbus Share settings). All placeholder hints and range calculation summaries MUST dynamically adapt to the communication protocol of the owning device.

#### Scenario: Start address placeholder adapts to protocol
- **WHEN** editing a rule owned by an MC 3E or FATEK device
- **THEN** the start address input displays a hint indicating alphanumeric register notation (e.g. `D0`, `M0`, `W0`)
- **AND** the address range preview shows valid alphanumeric bounds (e.g. `D0 ~ D7`)

#### Scenario: Start address placeholder for Modbus
- **WHEN** editing a rule owned by a Modbus device
- **THEN** the start address input displays a hint for 5-digit Modbus registers (e.g. `40001`)
- **AND** the address range preview shows numeric bounds (e.g. `40001 ~ 40008`)

#### Scenario: Reset skipped on count change
- **GIVEN** a rule with `count=8`, `skipped_addresses=['40003','40005']`
- **WHEN** the operator changes count to 6
- **THEN** the rule's `count` becomes 6
- **AND** `skipped_addresses` becomes `[]`
- **AND** the point grid below renders 6 cells, all enabled

#### Scenario: Reset skipped on data_type change
- **GIVEN** a rule with `data_type='int16'`, `skipped_addresses=['40002']`
- **WHEN** the operator changes data_type to `int32`
- **THEN** the rule's `data_type` becomes `int32`
- **AND** `skipped_addresses` becomes `[]`
- **AND** the grid renders cells using stride 2 (int32 width)

#### Scenario: Count clamped to 1-64
- **WHEN** the operator enters count 0 or count 100
- **THEN** the reducer clamps `count` to `max(1, min(64, input))`
- **AND** the grid reflects the clamped value

#### Scenario: Device select disabled with single device
- **WHEN** `state.devices.length === 1`
- **THEN** the owning-device select is disabled
- **AND** the rule is auto-assigned to that device

#### Scenario: Invalid start address blocks continuation
- **WHEN** the operator enters an address that is invalid for the owning device protocol
- **THEN** the editor displays a localized validation error and marks the input invalid
- **AND** no derived points are produced from the invalid address
- **AND** Step 2 does not allow continuation to Step 3

#### Scenario: One invalid enabled rule blocks a multi-rule workspace
- **GIVEN** a workspace has at least one enabled rule with valid derived points
- **AND** another enabled rule has an address that is invalid for its owning device protocol
- **WHEN** the operator attempts to continue from Step 2
- **THEN** Step 2 remains active and identifies the invalid rule
- **AND** valid points from other rules do not mask the invalid rule


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

---
### Requirement: Modbus function code inference

The range summary block SHALL display a chip indicating the Modbus function code derived from the rule's `start_address` first digit. The inference table MUST be: starting digit `0` → `Coil (R/W)` / `FC 01/05`, `1` → `Discrete Input (RO)` / `FC 02`, `3` → `Input Register (RO)` / `FC 04`, `4` (default) → `Holding Register` / `FC 03/06/16`.

#### Scenario: Holding register from 4xxxx

- **WHEN** the rule's `start_address` is `40001`
- **THEN** the range summary chip shows `Holding Register` with code `FC 03/06/16` in a blue tone

#### Scenario: Coil from 0xxxx

- **WHEN** the operator changes `start_address` to `00100`
- **THEN** the range summary chip shows `Coil (R/W)` with code `FC 01/05` in an amber tone

##### Example: function code mapping table

| start_address | Inferred function | Code         | Chip tone |
| ------------- | ----------------- | ------------ | --------- |
| 00100         | Coil (R/W)        | FC 01/05     | amber     |
| 10100         | Discrete Input    | FC 02        | slate     |
| 30100         | Input Register    | FC 04        | cyan      |
| 40001         | Holding Register  | FC 03/06/16  | blue      |


<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Point grid with modifier-key interactions

The current rule's point grid SHALL render the rule's expanded addresses as a 4-column (mobile), 6-column (sm), or 8-column (md+) grid. Each cell MUST display the address (mono) and the point name (smaller mono). Cell interactions MUST follow three modifier-key modes: plain click toggles the cell's skip state; Shift+click toggles a range from `lastClickedIdx` to the current index using the current cell's skip state as the inversion base; Cmd/Ctrl+click toggles the cell's address membership in a per-rule selection set. Switching the selected rule MUST reset the selection set and `lastClickedIdx`.

#### Scenario: Plain click toggles skip

- **GIVEN** a rule with cells `40001-40008` all enabled
- **WHEN** the operator clicks the `40003` cell
- **THEN** `state.rules[i].skipped_addresses` becomes `['40003']`
- **AND** the cell tone switches to amber

#### Scenario: Shift+click range

- **GIVEN** a rule with cells all enabled, and `lastClickedIdx` is 0 (40001)
- **WHEN** the operator Shift+clicks the cell at idx 3 (40004)
- **THEN** `skipped_addresses` includes `40001`, `40002`, `40003`, `40004`
- **AND** cells idx 0-3 show amber tone

#### Scenario: Cmd+click multi-select

- **GIVEN** a rule with cells all enabled
- **WHEN** the operator Cmd+clicks cells `40002`, `40004`, `40006`
- **THEN** those three addresses appear in the selection set
- **AND** those three cells show a blue ring border
- **AND** the chip `已選 3` is rendered in the section header aside

#### Scenario: Switching rule resets selection

- **GIVEN** the operator has Cmd+clicked three cells in rule A
- **WHEN** the operator selects rule B in the tab rail
- **THEN** the selection set becomes empty
- **AND** `lastClickedIdx` becomes null
- **AND** no cell shows a blue ring border in rule B's grid


<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Point grid batch toolbar

The point grid toolbar SHALL provide six batch actions: `全部啟用`, `全部略過`, `反轉啟用`, `略過選取 (N)`, `啟用選取 (N)`, `清除選取`. The `略過選取` and `啟用選取` actions MUST be disabled when the selection set is empty. The `清除選取` action MUST only render when the selection set is non-empty.

#### Scenario: 全部啟用 clears skipped

- **GIVEN** a rule with `skipped_addresses=['40003','40005']`
- **WHEN** the operator clicks `全部啟用`
- **THEN** `skipped_addresses` becomes `[]`
- **AND** the selection set becomes empty

#### Scenario: 反轉啟用 swaps skipped vs enabled

- **GIVEN** a rule with `count=8`, `skipped_addresses=['40001','40003','40005']`
- **WHEN** the operator clicks `反轉啟用`
- **THEN** `skipped_addresses` becomes `['40002','40004','40006','40007','40008']` (the previously enabled addresses)

#### Scenario: 略過選取 applies selection and clears it

- **GIVEN** the operator has Cmd+clicked `40002`, `40004` (selection set size 2)
- **WHEN** the operator clicks `略過選取 (2)`
- **THEN** `skipped_addresses` adds `40002` and `40004`
- **AND** the selection set becomes empty


<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Cross-rule address conflict detection

The system SHALL detect address conflicts across enabled rules: when two or more enabled rules have an enabled (non-skipped) point at the same address, that address MUST be flagged as a conflict. Conflict cells MUST be rendered with a red ring border and a red `!` badge in the upper-right corner. The merged point table MUST mark conflicting rows with a red address value and a `衝突` chip. The continue button MUST be disabled when any conflict exists.

#### Scenario: Two rules collide at 40001

- **GIVEN** rule A starts at 40001 with int16 count 4 (addresses 40001-40004), and rule B starts at 40001 with int16 count 2 (addresses 40001-40002)
- **WHEN** the workspace renders
- **THEN** addresses `40001` and `40002` are in the conflict set
- **AND** in rule A's grid, cells `40001` and `40002` show red ring and `!` badge
- **AND** the merged point table shows conflict chip on rows for both rules at those addresses
- **AND** the continue button is disabled

#### Scenario: Skipped point does not conflict

- **GIVEN** rule A and rule B both target `40001` but rule A has `40001` in its `skipped_addresses`
- **WHEN** the workspace renders
- **THEN** `40001` is NOT in the conflict set
- **AND** the continue button is enabled (assuming no other conflicts)


<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Modbus Share layout

When `state.settings.modbus_share.enabled === true`, the system SHALL compute a share layout assigning each share-enabled rule a starting register, stride, and end register. The computation MUST follow: walk rules in array order; for each rule with `share_enabled === true`, derive `stride = rule.share_stride ?? dataTypeWidth(rule.data_type)`; if `rule.share_start_register` is non-null use it as `start`, otherwise use the current cursor; `end = start + enabledCount * stride`; advance cursor to `max(cursor, end)`. Share-disabled rules MUST have layout `null`. Each non-skipped point in a share-enabled rule MUST have its share address rendered as a floating mono label `→{shareAddr}` in the grid cell's bottom-right corner. When the global `state.settings.modbus_share.enabled === false`, share addresses MUST NOT be rendered in the grid and the rule's share details summary MUST show a `全域未啟用` chip.

#### Scenario: Auto layout for two share-enabled rules

- **GIVEN** `state.settings.modbus_share.enabled=true`, `base_register=40001`, two rules with share_enabled=true: rule A (count=4, int16, no manual start), rule B (count=2, int32, no manual start)
- **WHEN** computeShareLayout runs
- **THEN** rule A's layout is `{ start: 40001, stride: 1, end: 40005, auto: true }`
- **AND** rule B's layout is `{ start: 40005, stride: 2, end: 40009, auto: true }`

#### Scenario: Manual share_start_register

- **GIVEN** the same rules but rule B has `share_start_register=50001`
- **WHEN** computeShareLayout runs
- **THEN** rule B's layout is `{ start: 50001, stride: 2, end: 50005, auto: false }`
- **AND** rule A's layout remains `{ start: 40001, stride: 1, end: 40005, auto: true }`

#### Scenario: Share globally disabled hides cell labels

- **WHEN** `state.settings.modbus_share.enabled === false` and a rule has `share_enabled=true`
- **THEN** the grid cells for that rule do NOT render any `→{shareAddr}` label
- **AND** the rule editor's Modbus Share details summary shows a `全域未啟用` chip

#### Scenario: Share-disabled rule has null layout

- **GIVEN** `state.settings.modbus_share.enabled=true` and a rule has `share_enabled=false`
- **WHEN** computeShareLayout runs
- **THEN** that rule's layout entry is `null`
- **AND** the grid cells for that rule do NOT render any share label

#### Scenario: Reset to auto

- **GIVEN** a rule has `share_start_register=50001` and `share_stride=3`
- **WHEN** the operator clicks `改回自動分配`
- **THEN** `share_start_register` becomes null
- **AND** `share_stride` becomes null
- **AND** the rule's layout reverts to auto computation


<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Merged point table summary

The merged point table SHALL render one row per derived point across all enabled rules. Columns MUST include: device color dot + name, rule color dot + name, point name, address (red if in conflict set), data_type, function code, share address (mono emerald if share enabled and not skipped, otherwise `—`), and status chip (`建立` success / `跳過` warning / `衝突` error). The table header MUST stick to the top during vertical scroll. The table footer MUST display the polling group assignment (default `快速輪詢 (1s)`), a hint about share when globally disabled, and the `繼續到映射` button.

#### Scenario: Skipped row is dim

- **GIVEN** a rule has a skipped address
- **WHEN** the merged table renders
- **THEN** the row for the skipped point has reduced opacity (~40%) and a `跳過` warning chip

#### Scenario: Conflict row is highlighted

- **GIVEN** two rules conflict at `40001`
- **WHEN** the merged table renders
- **THEN** both rows for address `40001` display the address text in red with a `⚠` icon and a `衝突` error chip

#### Scenario: Continue button gate

- **WHEN** the total enabled point count is 0 OR the conflict set is non-empty
- **THEN** the `繼續到映射` button is disabled


<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Rule-level scale inheritance to points

Each derived point produced by `derivePoints(rule, ...)` MUST carry `_rule_scale` and `_rule_offset` fields copied from the rule's `scale_multiplier` and `scale_offset`. These fields MUST be available to Step 3 mapping initialization so per-mapping scale defaults inherit from the rule layer, while still allowing per-mapping override in Step 3.

#### Scenario: Point inherits rule scale

- **GIVEN** a rule with `scale_multiplier=0.1`, `scale_offset=0`
- **WHEN** `derivePoints` runs
- **THEN** every produced point has `_rule_scale === 0.1` and `_rule_offset === 0`

#### Scenario: Updating rule scale updates derived points on next render

- **GIVEN** a rule with `scale_multiplier=0.1`
- **WHEN** the operator changes `scale_multiplier` to 0.5
- **THEN** the next call to `useAllPoints()` returns points with `_rule_scale === 0.5`
- **AND** Step 3 (when applied) will use 0.5 as the initial scale for newly created mappings

<!-- @trace
source: datalink-workbench-v2-step2-rule
updated: 2026-05-29
code:
  - frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/components/Button.tsx
  - frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGridToolbar.tsx
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/features/datalink/workbench-v2/components/inputs.tsx
  - frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx
  - frontend/src/main.tsx
  - frontend/src/features/datalink/workbench-v2/tokens.ts
  - frontend/src/features/datalink/workbench-v2/components/Icon.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ProtocolSelector.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/ShareSection.tsx
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx
  - findings.md
  - frontend/src/features/datalink/workbench-v2/state/ruleReducer.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/state/protocols.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/PointGrid.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.test-d.ts
  - frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx
  - frontend/src/features/datalink/workbench-v2/components/Toggle.tsx
  - .antigravitycli/4252526d-bebd-463d-84e9-9145d2a0eb40.json
  - frontend/src/features/datalink/workbench-v2/state/sourceRule.ts
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx
  - frontend/package.json
  - frontend/src/features/datalink/workbench-v2/steps/step2/ScaleSection.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/state/deviceColors.ts
  - frontend/src/i18n/config.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/MergedPointTable.tsx
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx
  - progress.md
  - frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - frontend/src/features/datalink/workbench-v2/components/index.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/Field.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step2/index.ts
tests:
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/protocols.test.ts
  - frontend/tests/unit/workbench-v2/step2-grid.test.tsx
  - frontend/tests/unit/workbench-v2/types-step2.test-d.ts
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/deviceColors.test.tsx
  - frontend/tests/unit/workbench-v2/step1-readiness.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - frontend/tests/unit/workbench-v2/components.test.tsx
  - frontend/tests/unit/workbench-v2/sourceRule.test.ts
  - frontend/tests/unit/workbench-v2/step2-share.test.tsx
  - frontend/tests/unit/workbench-v2/tokens.test.ts
  - frontend/tests/unit/workbench-v2/state.test.ts
  - frontend/tests/unit/workbench-v2/reducer-step2.test.ts
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/workbench-v2/types.test-d.ts
-->

---
### Requirement: Valid-only rule autosave

The Step 2 source-rule workspace SHALL autosave only valid rule content to the backend workspace.

#### Scenario: Valid rule change saves immediately

- **GIVEN** a rule already belongs to the singleton v2 workspace
- **WHEN** the operator edits that rule into a valid state
- **THEN** the system saves the rule immediately
- **AND** the persisted rule revision reflects the latest valid values

#### Scenario: Invalid rule change stays local

- **GIVEN** a rule already has a last successful persisted revision
- **WHEN** the operator edits that rule into an invalid state
- **THEN** the system does not overwrite the persisted revision
- **AND** the UI keeps the invalid local values visible with an unsaved marker


<!-- @trace
source: wire-studio-v2-step2-rule-autosave
updated: 2026-05-31
code:
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/sourcerule/errors.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - docs/goal.md
  - internal/datalink/dbtarget/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/point_handler.go
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/db.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/migrator.go
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - scripts/start-process-utils.sh
  - start.ps1
  - findings.md
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service.go
  - go.sum
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - scripts/start-log-utils.sh
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/workspace/service_database.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - progress.md
  - start.sh
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/schema/schema_device_models.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/hooks/datalink/index.ts
tests:
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
-->

---
### Requirement: Per-rule save isolation

The Step 2 source-rule workspace SHALL isolate save success and failure per rule.

#### Scenario: One invalid rule does not block another valid rule

- **GIVEN** the workspace contains multiple rules
- **WHEN** one rule remains invalid while another rule is edited into a valid state
- **THEN** the valid rule still saves successfully
- **AND** the invalid rule remains local with its own error state


<!-- @trace
source: wire-studio-v2-step2-rule-autosave
updated: 2026-05-31
code:
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/sourcerule/errors.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - docs/goal.md
  - internal/datalink/dbtarget/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/point_handler.go
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/db.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/migrator.go
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - scripts/start-process-utils.sh
  - start.ps1
  - findings.md
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service.go
  - go.sum
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - scripts/start-log-utils.sh
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/workspace/service_database.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - progress.md
  - start.sh
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/schema/schema_device_models.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/hooks/datalink/index.ts
tests:
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
-->

---
### Requirement: Workspace-scoped rule ownership

Every Step 2 source rule SHALL preserve its owning workspace and device relationship across reloads.

#### Scenario: Reload restores device ownership

- **WHEN** the operator reloads `/studio/v2`
- **THEN** each persisted rule reappears under the same workspace and owning device as before reload
- **AND** the system does not silently reassign that rule to another device

<!-- @trace
source: wire-studio-v2-step2-rule-autosave
updated: 2026-05-31
code:
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts
  - internal/datalink/dbtarget/types.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - internal/datalink/sourcerule/errors.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - internal/api/router.go
  - frontend/src/features/datalink/legacyRoutes.ts
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/datalink/mapping/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx
  - docs/goal.md
  - internal/datalink/dbtarget/service.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/device/service_status.go
  - frontend/src/hooks/datalink/useStudioV2Workspace.ts
  - internal/api/handlers/runtime_handler.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - internal/api/handlers/point_handler.go
  - go.mod
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - internal/datalink/workspace/service_devices.go
  - frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx
  - docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md
  - frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/state/settingsDefaults.ts
  - frontend/src/services/studioV2Workspace.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx
  - internal/datalink/db.go
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/datalink/migrator.go
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - internal/api/handlers/mapping_handler.go
  - internal/datalink/device/sql_repo.go
  - internal/datalink/device/availability.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
  - frontend/src/types/studioV2Activation.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useSettings.ts
  - internal/datalink/workspace/sql_repository.go
  - internal/datalink/workspace/memory_repository.go
  - internal/api/handlers/point_direct_reader.go
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_handler.go
  - frontend/src/services/studioV2Rules.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - scripts/start-process-utils.sh
  - start.ps1
  - findings.md
  - internal/datalink/sourcerule/activation_readiness.go
  - frontend/src/services/studioV2WorkspaceActivation.ts
  - frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - internal/datalink/sourcerule/service_links.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx
  - task_plan.md
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - internal/datalink/sourcerule/service.go
  - go.sum
  - internal/api/handlers/source_rule_handler.go
  - frontend/src/i18n/locales/en/workbench-v2.json
  - docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md
  - frontend/src/types/datalink.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/services/studioV2Mappings.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts
  - frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - scripts/start-log-utils.sh
  - gateway.db
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - frontend/src/types/studioV2Availability.ts
  - internal/datalink/device/service_crud.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - internal/api/handlers/point_handler_polling.go
  - internal/datalink/workspace/service_database.go
  - internal/datalink/mapping/service_crud.go
  - internal/datalink/device/errors.go
  - frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - progress.md
  - start.sh
  - .air.toml
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/workspace/service_activation.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - frontend/src/hooks/datalink/keys.ts
  - internal/datalink/device/service_point_read.go
  - frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts
  - internal/datalink/schema/schema_device_models.go
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/services/studioV2WorkspaceDevices.ts
  - frontend/src/services/studioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/hooks/datalink/index.ts
tests:
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/handlers/point_handler_poll_contract_test.go
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/router_logger_test.go
  - internal/datalink/mapping/service_workspace_scope_test.go
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - frontend/tests/unit/workbench-v2/settings-backend.test.tsx
  - internal/datalink/workspace/service_device_order_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/api/router_studio_v2_workspace_devices_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/settings-connectors.test.tsx
  - internal/datalink/sourcerule/service_workspace_scope_test.go
  - internal/api/router_studio_v2_runtime_context_test.go
  - internal/api/router_studio_v2_workspace_source_rules_test.go
  - frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - internal/api/handlers/device_handler_extended_test.go
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/datalink/device/availability_test.go
  - frontend/tests/unit/workbench-v2/step3-live-values.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx
  - internal/datalink/workspace/service_test.go
  - frontend/tests/unit/workbench-v2/settings.test.tsx
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - internal/datalink/dbtarget/service_mysql_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - internal/datalink/workspace/service_activation_test.go
  - frontend/tests/unit/workbench-v2/selectors.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step1.test.ts
  - frontend/tests/unit/workbench-v2/settings-backend-mappings.test.ts
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/device/service_crud_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx
  - internal/api/router_studio_v2_workspace_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/api/handlers/dbtarget_handler_connectors_test.go
  - internal/api/router_studio_v2_workspace_mappings_test.go
-->

---
### Requirement: Step 2 reload restores persisted source-rule definitions and save truth

The Step 2 source-rule workspace SHALL restore persisted source-rule definitions and their persisted save truth after reload.

#### Scenario: Reload shows persisted source rules

- **WHEN** the operator reloads Studio V2 after source rules have been persisted for the workspace
- **THEN** Step 2 renders those persisted source rules with the correct owning device, enabled state, and persisted save state indicators
- **AND** Step 2 SHALL NOT regenerate only local placeholder rules when persisted rules already exist

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
### Requirement: Protocol-aware point address derivation

The system SHALL derive point addresses by taking the rule's `start_address`, `count`, `data_type` (width stride), and owning device's `protocol`. The address calculation MUST preserve alphanumeric register area prefixes (such as `D`, `M`, `W`, `X`, `Y` for Mitsubishi/FATEK) and correctly advance address offsets without stripping characters or falling back to Modbus defaults.

#### Scenario: Derive points for MC 3E device
- **WHEN** a rule for an MC 3E device is defined with start address `D100`, count `4`, and data type `int16` (width `1`)
- **THEN** the system derives 4 points with addresses `D100`, `D101`, `D102`, and `D103`
- **AND** no Modbus `40001` fallback is applied

##### Example: Multi-protocol point derivation
| Protocol | Start Address | Count | Data Type | Expected Addresses |
| --- | --- | --- | --- | --- |
| `modbus_tcp` | `40001` | `3` | `int16` | `40001, 40002, 40003` |
| `mc_3e` | `D0` | `4` | `int16` | `D0, D1, D2, D3` |
| `mc_3e` | `D100` | `2` | `int32` | `D100, D102` |
| `mc_3e` | `X0` | `3` | `bool` | `X0, X1, X2` |
| `fatek_fbs` | `R0` | `3` | `int16` | `R0, R1, R2` |

#### Scenario: Backend address offset matches protocol radix
- **WHEN** the source-rule service offsets MC 3E address `X0` by `16`
- **THEN** it returns `X10`
- **AND** MC 3E areas `X`, `Y`, and `B` use hexadecimal suffixes while supported decimal areas retain decimal suffixes

#### Scenario: Backend rejects invalid protocol address
- **WHEN** source-rule validation or address offset receives an address invalid for the selected protocol
- **THEN** the operation returns an explicit validation error
- **AND** it does not return the original input as a successful derived address

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