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

---
### Requirement: Database target workflow uses the live visible target set

The database target workflow SHALL use the live visible target set for planning, schema ensure, and apply behavior.

#### Scenario: Workflow excludes stale hidden mappings

- **WHEN** a database connector still contains stale hidden target mappings that are no longer visible from the current live rule-owned state
- **THEN** the workflow excludes those stale mappings from planning and delivery operations
- **AND** it SHALL NOT let a stale hidden mapping silently change or block the live target set

##### Example: planner excludes one orphaned hidden mapping

- **GIVEN** the current live planner shows targets for tags tag-A and tag-B while connector metadata still retains a hidden orphaned mapping for tag-old
- **WHEN** the workflow computes planning and delivery operations
- **THEN** tag-old is excluded from the live target set and cannot block tag-A or tag-B

<!-- @trace
source: guarantee-database-output-delivery
updated: 2026-06-09
code:
  - internal/datalink/workspace/service_readiness.go
  - internal/api/handlers/runtime_stream_handler.go
  - Makefile
  - internal/api/router.go
  - frontend/src/hooks/datalink/useStudioV2Rules.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.up.sql
  - frontend/src/features/datalink/workbench-v2/steps/step4/CollapsibleSupportCard.tsx
  - internal/datalink/runtime/delivery_diagnostic.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DestinationOverviewCard.tsx
  - frontend/e2e-studio-v2-live.cjs
  - frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts
  - tests/shell/start-frontend-install-failure.sh
  - internal/datalink/workspace/service_activation.go
  - internal/datalink/runtime/service_workspace_projection.go
  - frontend/src/services/studioV2WorkspaceAudit.ts
  - frontend/src/types/runtimeDiagnostics.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - frontend/src/features/datalink/workbench-v2/shell/WorkspaceAuditHistoryPanel.tsx
  - frontend/src/services/studioV2RuntimeContext.ts
  - internal/datalink/dbtarget/sql_repository.go
  - frontend/src/hooks/datalink/keys.ts
  - frontend/src/types/studioV2Workspace.ts
  - internal/datalink/dbtarget/sql_repository_connector_scan.go
  - internal/datalink/workspace/service_runtime_projection.go
  - internal/datalink/dbtarget/delivery_outcome.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/index.ts
  - tests/shell/start-frontend-deps.sh
  - internal/api/handlers/studio_v2_workspace_audit.go
  - internal/api/handlers/studio_v2_workspace_devices_handler.go
  - internal/datalink/schema/schema_dbtarget_models.go
  - internal/datalink/workspace/service.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx
  - frontend/src/i18n/locales/en/workbench-v2.json
  - internal/api/handlers/studio_v2_workspace_mappings_delete.go
  - frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - internal/api/handlers/studio_v2_workspace_mapping_types.go
  - frontend/src/features/datalink/workbench-v2/state/selectors.ts
  - frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - frontend/studio-v2-diagnostic.cjs
  - frontend/src/App.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPayloadDialog.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx
  - frontend/src/i18n/locales/zh-TW/runtime-dashboard.json
  - frontend/src/types/runtimeTruth.ts
  - frontend/src/services/studioV2Mappings.ts
  - internal/api/handlers/studio_v2_workspace_audit_history_handler.go
  - internal/datalink/migrator.go
  - internal/datalink/runtime/service_source_rule_reconcile.go
  - frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx
  - frontend/src/hooks/datalink/useStudioV2Mappings.ts
  - internal/datalink/audit/service.go
  - frontend/src/i18n/locales/en/runtime-dashboard.json
  - internal/api/handlers/studio_v2_workspace_database_handler_helpers.go
  - internal/datalink/dbtarget/writer_statements.go
  - frontend/studio-v2-real-check.cjs
  - frontend/src/features/datalink/workbench-v2/state/studioV2RuleAutosave.ts
  - internal/datalink/workspace/service_runtime_projection_scope.go
  - internal/api/handlers/studio_v2_workspace_mappings_recovery.go
  - frontend/src/features/datalink/workbench-v2/steps/step4/DeliveryTruthStrip.tsx
  - internal/datalink/runtime/truth_state.go
  - internal/datalink/dbtarget/live_projection.go
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx
  - start.sh
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - internal/api/handlers/studio_v2_runtime_apply.go
  - frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts
  - internal/api/handlers/studio_v2_workspace_activation_handler.go
  - internal/api/handlers/runtime_workspace_setup_context.go
  - internal/datalink/migrator_database_delivery_outcomes.go
  - frontend/src/types/studioV2RuntimeApply.ts
  - frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - frontend/src/types/studioV2WorkspaceAudit.ts
  - frontend/src/features/datalink/workbench-v2/state/types.ts
  - frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingPreviewCells.tsx
  - internal/datalink/workspace/service_runtime_projection_link_scope.go
  - frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - internal/api/handlers/source_rule_handler.go
  - internal/datalink/audit/types.go
  - internal/api/handlers/studio_v2_workspace_mappings_handler.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceAuditHistory.ts
  - frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - frontend/src/types/studioV2WorkspaceReadiness.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes_sqlite.up.sql
  - internal/datalink/runtime/service_projection_state.go
  - internal/api/handlers/studio_v2_workspace_handler.go
  - internal/datalink/sourcerule/runtime_reconcile.go
  - frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts
  - frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx
  - frontend/src/features/datalink/workbench-v2/state/dbReducer.ts
  - frontend/src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave.ts
  - frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx
  - frontend/src/features/datalink/workbench-v2/components/WorkspaceReadinessPanel.tsx
  - frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx
  - internal/datalink/dbtarget/service.go
  - internal/api/handlers/studio_v2_workspace_database_handler.go
  - frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - internal/datalink/runtime/service.go
  - internal/datalink/dbtarget/service_validate.go
  - frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts
  - frontend/src/types/studioV2RuntimeContext.ts
  - internal/datalink/runtime/ingestor.go
  - internal/datalink/audit/memory_repository.go
  - internal/datalink/dbtarget/tooling_service.go
  - internal/api/handlers/runtime_workspace_setup_types.go
  - frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts
  - internal/datalink/schema/migrations/015_database_delivery_outcomes.down.sql
  - internal/datalink/audit/sql_repository.go
  - frontend/src/features/datalink/workbench-v2/steps/step2/RuleTabRail.tsx
  - internal/api/handlers/runtime_handler.go
  - frontend/src/types/databaseDelivery.ts
  - frontend/src/types/datalink.ts
  - internal/datalink/dbtarget/writer.go
  - internal/datalink/schema/migrations/016_workspace_audit_history_sqlite.up.sql
  - tests/shell/start-frontend-readiness.sh
  - frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx
  - frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx
  - internal/datalink/runtime/status.go
  - frontend/src/features/datalink/workbench-v2/steps/step3/index.ts
  - frontend/src/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel.tsx
  - internal/datalink/runtime/service_device_sync.go
  - internal/api/handlers/studio_v2_workspace_source_rules_handler.go
  - frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts
tests:
  - internal/datalink/workspace/service_readiness_stale_relationships_test.go
  - internal/datalink/workspace/service_readiness_test.go
  - frontend/tests/unit/workbench-v2/shell-readiness.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx
  - internal/api/handlers/studio_v2_workspace_mappings_handler_test.go
  - cmd/test_ui/main.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-diagnostics-state.test.tsx
  - internal/api/handlers/runtime_handler_projection_test.go
  - internal/datalink/runtime/service_workspace_projection_test.go
  - frontend/tests/unit/workbench-v2/reducer-step3.test.ts
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx
  - frontend/tests/unit/workbench-v2/workspace-boot.test.tsx
  - internal/api/handlers/studio_v2_workspace_database_readiness_test.go
  - internal/api/handlers/studio_v2_workspace_database_handler_test.go
  - internal/api/handlers/runtime_workspace_setup_context_recovery_test.go
  - internal/datalink/workspace/service_activation_test.go
  - internal/api/handlers/studio_v2_workspace_devices_handler_test.go
  - internal/api/handlers/studio_v2_workspace_audit_handler_test.go
  - internal/api/handlers/studio_v2_workspace_activation_handler_test.go
  - frontend/tests/unit/workbench-v2/workspace-readiness-panel.test.tsx
  - frontend/tests/unit/workbench-v2/step1.test.tsx
  - frontend/tests/unit/hooks/useStudioV2WorkspaceAuditHistory.test.ts
  - internal/datalink/dbtarget/service_postgres_test.go
  - internal/datalink/workspace/service_runtime_projection_helpers_test.go
  - frontend/tests/unit/workbench-v2/shell-redesign.test.tsx
  - internal/datalink/runtime/ingestor_test.go
  - internal/datalink/sourcerule/service_runtime_reconcile_test.go
  - internal/datalink/workspace/service_readiness_recovery_test.go
  - internal/api/router_runtime_database_delivery_test.go
  - internal/api/router_studio_v2_workspace_test.go
  - internal/datalink/runtime/status_test.go
  - frontend/tests/unit/workbench-v2/step4-commit.test.tsx
  - frontend/tests/unit/workbench-v2/mappingDefaults.test.ts
  - frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts
  - frontend/tests/unit/workbench-v2/step4-database.test.tsx
  - internal/api/handlers/runtime_workspace_setup_context_regression_test.go
  - internal/datalink/audit/service_test.go
  - frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx
  - internal/datalink/dbtarget/live_projection_test.go
  - internal/datalink/runtime/service_test.go
  - frontend/tests/unit/workbench-v2/step4-delivery-truth.test.tsx
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - frontend/tests/unit/hooks/useStudioV2Rules.test.ts
  - internal/api/handlers/runtime_stream_handler_test.go
  - internal/datalink/workspace/service_readiness_connector_missing_test.go
  - internal/api/handlers/studio_v2_workspace_devices_reconcile_test.go
  - frontend/tests/unit/runtime-dashboard/runtime-dashboard-truth-state.test.tsx
  - frontend/tests/unit/workbench-v2/routing.test.tsx
  - internal/api/handlers/runtime_handler_test.go
  - internal/api/handlers/studio_v2_workspace_handler_test.go
  - frontend/tests/unit/hooks/useStudioV2Mappings.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_readiness_test.go
  - cmd/test_ui/static/index.html
  - internal/datalink/workspace/service_activation_projection_test.go
  - frontend/tests/unit/workbench-v2/step3-mapping.test.tsx
  - frontend/tests/unit/workbench-v2/reducer-step4.test.ts
  - internal/api/handlers/studio_v2_workspace_mappings_recovery_regression_test.go
  - internal/datalink/runtime/service_source_rule_reconcile_test.go
  - frontend/tests/unit/workbench-v2/shell.test.tsx
  - internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go
  - frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx
  - internal/datalink/workspace/service_runtime_projection_test.go
  - frontend/tests/unit/runtime-dashboard/runtimeSetupFixture.ts
  - internal/api/handlers/studio_v2_workspace_database_delivery_truth_test.go
  - internal/api/router_runtime_test.go
  - frontend/tests/unit/workbench-v2/step2-rule.test.tsx
  - frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx
  - internal/datalink/workspace/service_readiness_empty_test.go
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx
  - internal/api/router_studio_v2_workspace_activation_test.go
  - frontend/tests/unit/workbench-v2/step3-components.test.tsx
-->