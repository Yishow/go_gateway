# database-target-workbench Specification

## Purpose
Define the canonical Studio-side database output workflow so connector selection, schema context, and tag-to-column mapping remain scoped, coherent, and limited to the supported SQLite and PostgreSQL connectors.

## Requirements

### Requirement: Database output workflow is layered by connector, schema context, and mapping
The workbench SHALL present database output as four linked layers: source-rule revision context, connector selection, schema/table context, and grouped row planner review/apply.

#### Scenario: Connector selection scopes schema context and candidate board
- **WHEN** an operator selects a database connector in the Output step for a source rule
- **THEN** the available schema and table context SHALL scope to that connector only
- **AND** the rule-scoped database candidate board SHALL update using only that connector context

#### Scenario: Apply uses current layered context
- **WHEN** an operator applies a database output candidate
- **THEN** the mapping SHALL be applied against the current source-rule revision and the selected connector, schema, table, column, group key, and effective interval context
- **AND** the UI SHALL NOT apply the action to stale prior selections

---
### Requirement: Database connector scope for this workflow is SQLite and PostgreSQL
The system SHALL support SQLite and PostgreSQL as the database connector kinds for this workbench flow in this change.

#### Scenario: Supported connectors are first-class options
- **WHEN** an operator creates or edits a database target connector in the workbench
- **THEN** the workbench SHALL present SQLite and PostgreSQL as supported connector kinds
- **AND** SHALL NOT advertise unsupported connector kinds as first-class options in this workflow

---
### Requirement: Database output is rule-scoped and connector-aware
The workbench SHALL generate database output candidates from rule-owned tag state using persisted connector and table context.

The system SHALL revalidate blocked or `out_of_sync` database candidates whenever the referenced connector context becomes valid again.

#### Scenario: Connector invalidation blocks database candidate apply
- **WHEN** the selected connector, schema, table, or column becomes invalid for an existing database output candidate
- **THEN** the candidate or applied mapping is marked as `blocked` or `out_of_sync`
- **AND** database apply SHALL NOT silently rebind to a different connector or table

#### Scenario: Restored connector context re-enables database apply
- **WHEN** a previously invalid connector, schema, table, and column context becomes valid again for an existing database output candidate
- **THEN** the system revalidates the candidate against the restored context
- **AND** returns the candidate to an apply-eligible state without requiring the operator to recreate it from scratch

---
### Requirement: Database planner organizes compatible tags into grouped rows
The database workbench SHALL organize compatible database candidates into grouped row plans keyed by connector context, table context, group key, and interval semantics.

#### Scenario: Prefix-inferred tags appear as one row plan
- **WHEN** the candidate set includes `meter/A1`, `meter/A2`, `meter/A3`, and `meter/kw`
- **THEN** the database planner shows one `meter` row group with member columns `a1`, `a2`, `a3`, and `kw`
- **AND** the operator is not forced to bind each tag one by one as the primary first-pass workflow

---
### Requirement: Database grouped rows expose editable grouping and interval controls
The database workbench SHALL let the operator review and edit grouped row metadata before apply, including `group_key`, `column_name`, and effective interval.

#### Scenario: Operator overrides grouped row metadata before apply
- **WHEN** an inferred grouped row needs a different group key, column name, or interval
- **THEN** the planner lets the operator edit those values before apply
- **AND** the chosen overrides are preserved in the apply request

#### Scenario: Interval edits validate grouped compatibility
- **WHEN** the operator changes the effective interval for a grouped row
- **THEN** the planner validates whether every member still shares the same effective interval
- **AND** the row stays blocked until the interval mismatch is resolved

#### Scenario: Incompatible grouped members stay blocked
- **WHEN** a proposed grouped row mixes members with different connector, schema, table, write mode, timestamp column, or effective interval
- **THEN** the planner marks that row as blocked with an explicit reason
- **AND** dry-run/apply do not silently split or merge the incompatible members

---
### Requirement: Legacy single-member mappings remain first-class
The database workbench SHALL preserve legacy mappings and slashless tag keys as single-member rows when no group key is set.

#### Scenario: Slashless tag stays single-member
- **WHEN** a database candidate has no slash in its tag key or an existing mapping persists with `group_key = null`
- **THEN** the planner keeps it as a single-member row
- **AND** grouped row logic does not merge it with unrelated mappings

---
### Requirement: Database target context informs Source-step planning hints
The workbench SHALL let the active `database` output target influence Source-step planning guidance before the operator enters the Output step.

This guidance SHALL remain limited to planning defaults and grouped-row hints. It SHALL NOT require the operator to configure connector, schema, table, or column setup in Step 2.

#### Scenario: Database context changes Source-step grouping guidance
- **WHEN** the active output target is `database` while the operator is planning a source rule in Step 2
- **THEN** the workbench shows grouping-oriented hints that help the operator understand how the rule is likely to map into future table, row, or column groupings
- **AND** the workbench does not expose connector or schema setup controls in Step 2

#### Scenario: Database-aware preview remains rule-scoped
- **WHEN** Step 2 previews downstream Tag or Database implications for the focused source rule
- **THEN** the preview remains scoped to the active rule only
- **AND** the operator does not have to infer grouped-row implications from unrelated rules or device-wide state

<!-- @trace
source: source-step-rule-first-database-preview
updated: 2026-05-09
code:
  - .github/prompts/spectra-ask.prompt.md
  - .github/prompts/spectra-archive.prompt.md
  - tests/shell/start-port-management.sh
  - .github/prompts/spectra-debug.prompt.md
  - .github/skills/spectra-ask/SKILL.md
  - .github/skills/spectra-debug/SKILL.md
  - frontend/src/i18n/locales/en/common.json
  - .github/skills/spectra-audit/SKILL.md
  - .github/prompts/spectra-propose.prompt.md
  - .github/skills/spectra-apply/SKILL.md
  - frontend/src/i18n/locales/zh-TW/common.json
  - tests/shell/start-backend-before-frontend.sh
  - .github/prompts/spectra-apply.prompt.md
  - CLAUDE.md
  - .github/prompts/spectra-audit.prompt.md
  - .github/prompts/spectra-discuss.prompt.md
  - frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx
  - .github/skills/spectra-discuss/SKILL.md
  - .github/prompts/opsx-propose.prompt.md
  - start.sh
  - .github/skills/spectra-ingest/SKILL.md
  - tests/shell/start-backend-cleanup-order.sh
  - frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx
  - .spectra.yaml
  - .github/skills/spectra-commit/SKILL.md
  - frontend/src/main.tsx
  - .github/prompts/opsx-apply.prompt.md
  - frontend/src/pages/datalink/workbench/sourceStepRuleSummaryModel.ts
  - frontend/package.json
  - frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx
  - .github/prompts/opsx-archive.prompt.md
  - frontend/src/App.tsx
  - .github/skills/spectra-propose/SKILL.md
  - frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx
  - .github/prompts/spectra-commit.prompt.md
  - start.ps1
  - frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - frontend/src/components/DevAgentation.tsx
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx
  - tests/shell/start-backend-logfile.sh
  - .github/skills/spectra-archive/SKILL.md
  - .github/prompts/spectra-ingest.prompt.md
tests:
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx
  - cmd/test_ui/static/assets/index-CcV2SQjv.css
  - cmd/test_ui/static/assets/index-Bl1MOChG.js
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourcePreview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
  - frontend/tests/unit/components/DevAgentation.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
  - frontend/tests/unit/pages/datalink/workbench-source-preview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - cmd/test_ui/static/index.html
  - frontend/tests/unit/utils/appAgentationRemoval.test.tsx
-->