# source-step-desk-modes Specification

## Purpose

TBD - created by archiving change 'rework-source-step-desk-modes'. Update Purpose after archive.

## Requirements

### Requirement: Source desk modes share one planning skeleton
The system SHALL provide `Inspect`, `Build`, and `Triage` desks inside Step 2 while preserving one shared planning skeleton for the active source rule.

The shared skeleton SHALL keep the active rule summary, grouped Tag handoff context, diagnostic cues, and primary rule action visible without requiring route changes or mode-specific re-entry.

#### Scenario: Switching desks preserves the active rule context
- **WHEN** an operator switches between `Inspect`, `Build`, and `Triage` while the same source rule is active
- **THEN** Step 2 keeps the same active rule summary, grouped handoff copy, and blocker diagnostics visible
- **AND** the operator does not need to reselect the rule to continue planning

#### Scenario: Desk switching stays inside one Source workspace
- **WHEN** an operator changes the Source desk mode
- **THEN** the system updates the Step 2 workspace in place
- **AND** the workflow remains inside the same `/studio` Source step without opening a separate route or modal-only flow


<!-- @trace
source: rework-source-step-desk-modes
updated: 2026-05-09
code:
  - tests/shell/start-backend-before-frontend.sh
  - frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx
  - frontend/src/App.tsx
  - .github/skills/spectra-ask/SKILL.md
  - frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - .github/prompts/spectra-debug.prompt.md
  - .github/prompts/spectra-audit.prompt.md
  - .github/prompts/spectra-propose.prompt.md
  - .github/prompts/opsx-archive.prompt.md
  - .github/skills/spectra-apply/SKILL.md
  - .github/skills/spectra-audit/SKILL.md
  - .github/skills/spectra-commit/SKILL.md
  - .github/skills/spectra-discuss/SKILL.md
  - .spectra.yaml
  - .github/prompts/opsx-propose.prompt.md
  - CLAUDE.md
  - frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx
  - .github/skills/spectra-debug/SKILL.md
  - tests/shell/start-backend-cleanup-order.sh
  - start.sh
  - frontend/src/i18n/locales/zh-TW/common.json
  - frontend/src/pages/datalink/workbench/sourceStepRuleSummaryModel.ts
  - .github/skills/spectra-ingest/SKILL.md
  - .github/skills/spectra-archive/SKILL.md
  - frontend/src/i18n/locales/en/common.json
  - tests/shell/start-backend-logfile.sh
  - frontend/src/main.tsx
  - tests/shell/start-port-management.sh
  - .github/skills/spectra-propose/SKILL.md
  - .github/prompts/spectra-discuss.prompt.md
  - .github/prompts/spectra-ingest.prompt.md
  - .github/prompts/spectra-ask.prompt.md
  - frontend/src/components/DevAgentation.tsx
  - frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx
  - .github/prompts/spectra-commit.prompt.md
  - frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx
  - frontend/package.json
  - .github/prompts/spectra-archive.prompt.md
  - start.ps1
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - .github/prompts/opsx-apply.prompt.md
  - .github/prompts/spectra-apply.prompt.md
  - frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx
tests:
  - cmd/test_ui/static/assets/index-CcV2SQjv.css
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/components/DevAgentation.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
  - cmd/test_ui/static/assets/index-Bl1MOChG.js
  - frontend/tests/unit/pages/datalink/workbench-source-preview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - frontend/tests/unit/utils/appAgentationRemoval.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourcePreview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Inspect and Build desks expose distinct working emphasis
The system SHALL make `Inspect` the canvas-first verification desk and SHALL make `Build` the rule-editor-first planning desk.

`Inspect` SHALL prioritize reading planned spans, live address context, and downstream readiness on the shared canvas. `Build` SHALL prioritize editing the active rule while still keeping a readable canvas region visible for immediate planning feedback.

#### Scenario: Inspect desk emphasizes canvas verification
- **WHEN** the operator opens `Inspect`
- **THEN** the main Step 2 workspace prioritizes the shared canvas and planning verification surfaces
- **AND** the rule-editing controls remain available as secondary support instead of dominating the viewport

#### Scenario: Build desk keeps the editor primary without hiding the canvas
- **WHEN** the operator opens `Build` for an active source rule
- **THEN** the rule editor becomes the primary workspace for editing
- **AND** the canvas remains visible with enough size to review planned spans, unmanaged gaps, and current address coverage without leaving the desk

##### Example: desk emphasis matrix
| Desk | Primary surface | Secondary surface |
| ---- | --------------- | ----------------- |
| Inspect | Shared source canvas | Rule details and diagnostics |
| Build | Active rule editor | Readable source canvas |


<!-- @trace
source: rework-source-step-desk-modes
updated: 2026-05-09
code:
  - tests/shell/start-backend-before-frontend.sh
  - frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx
  - frontend/src/App.tsx
  - .github/skills/spectra-ask/SKILL.md
  - frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - .github/prompts/spectra-debug.prompt.md
  - .github/prompts/spectra-audit.prompt.md
  - .github/prompts/spectra-propose.prompt.md
  - .github/prompts/opsx-archive.prompt.md
  - .github/skills/spectra-apply/SKILL.md
  - .github/skills/spectra-audit/SKILL.md
  - .github/skills/spectra-commit/SKILL.md
  - .github/skills/spectra-discuss/SKILL.md
  - .spectra.yaml
  - .github/prompts/opsx-propose.prompt.md
  - CLAUDE.md
  - frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx
  - .github/skills/spectra-debug/SKILL.md
  - tests/shell/start-backend-cleanup-order.sh
  - start.sh
  - frontend/src/i18n/locales/zh-TW/common.json
  - frontend/src/pages/datalink/workbench/sourceStepRuleSummaryModel.ts
  - .github/skills/spectra-ingest/SKILL.md
  - .github/skills/spectra-archive/SKILL.md
  - frontend/src/i18n/locales/en/common.json
  - tests/shell/start-backend-logfile.sh
  - frontend/src/main.tsx
  - tests/shell/start-port-management.sh
  - .github/skills/spectra-propose/SKILL.md
  - .github/prompts/spectra-discuss.prompt.md
  - .github/prompts/spectra-ingest.prompt.md
  - .github/prompts/spectra-ask.prompt.md
  - frontend/src/components/DevAgentation.tsx
  - frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx
  - .github/prompts/spectra-commit.prompt.md
  - frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx
  - frontend/package.json
  - .github/prompts/spectra-archive.prompt.md
  - start.ps1
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - .github/prompts/opsx-apply.prompt.md
  - .github/prompts/spectra-apply.prompt.md
  - frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx
tests:
  - cmd/test_ui/static/assets/index-CcV2SQjv.css
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/components/DevAgentation.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
  - cmd/test_ui/static/assets/index-Bl1MOChG.js
  - frontend/tests/unit/pages/datalink/workbench-source-preview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - frontend/tests/unit/utils/appAgentationRemoval.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourcePreview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx
  - cmd/test_ui/static/index.html
-->

---
### Requirement: Triage desk is issue-first and explicit when clear
The system SHALL make `Triage` the issue-first recovery desk for conflict, blocked, and unmanaged planning states.

`Triage` SHALL surface recovery actions and the affected planning context before general inspection content. When no triage-worthy issues exist, the desk SHALL render an explicit clear-state message instead of falling back to an inspect-like layout.

#### Scenario: Triage desk elevates incident recovery
- **WHEN** conflict, blocked, or unmanaged items exist for the active source rule
- **THEN** the `Triage` desk shows those items and their recovery actions as the primary review surface
- **AND** the desk limits canvas context to the affected ranges needed to understand and repair the issue

#### Scenario: Triage desk shows a clear-state message when no issues exist
- **WHEN** the operator opens `Triage` and the active source rule has no conflict, blocked, or unmanaged items requiring intervention
- **THEN** the desk shows a clear-state explanation that no triage work is pending
- **AND** the desk guides the operator back to `Inspect` or `Build` instead of rendering a generic canvas-first workspace

<!-- @trace
source: rework-source-step-desk-modes
updated: 2026-05-09
code:
  - tests/shell/start-backend-before-frontend.sh
  - frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx
  - frontend/src/App.tsx
  - .github/skills/spectra-ask/SKILL.md
  - frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx
  - tests/shell/start-backend-cleanup-no-log-wait.sh
  - .github/prompts/spectra-debug.prompt.md
  - .github/prompts/spectra-audit.prompt.md
  - .github/prompts/spectra-propose.prompt.md
  - .github/prompts/opsx-archive.prompt.md
  - .github/skills/spectra-apply/SKILL.md
  - .github/skills/spectra-audit/SKILL.md
  - .github/skills/spectra-commit/SKILL.md
  - .github/skills/spectra-discuss/SKILL.md
  - .spectra.yaml
  - .github/prompts/opsx-propose.prompt.md
  - CLAUDE.md
  - frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx
  - .github/skills/spectra-debug/SKILL.md
  - tests/shell/start-backend-cleanup-order.sh
  - start.sh
  - frontend/src/i18n/locales/zh-TW/common.json
  - frontend/src/pages/datalink/workbench/sourceStepRuleSummaryModel.ts
  - .github/skills/spectra-ingest/SKILL.md
  - .github/skills/spectra-archive/SKILL.md
  - frontend/src/i18n/locales/en/common.json
  - tests/shell/start-backend-logfile.sh
  - frontend/src/main.tsx
  - tests/shell/start-port-management.sh
  - .github/skills/spectra-propose/SKILL.md
  - .github/prompts/spectra-discuss.prompt.md
  - .github/prompts/spectra-ingest.prompt.md
  - .github/prompts/spectra-ask.prompt.md
  - frontend/src/components/DevAgentation.tsx
  - frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx
  - .github/prompts/spectra-commit.prompt.md
  - frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx
  - frontend/package.json
  - .github/prompts/spectra-archive.prompt.md
  - start.ps1
  - node_modules/.vite/vitest/da39a3ee5e6b4b0d3255bfef95601890afd80709/results.json
  - .github/prompts/opsx-apply.prompt.md
  - .github/prompts/spectra-apply.prompt.md
  - frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx
tests:
  - cmd/test_ui/static/assets/index-CcV2SQjv.css
  - cmd/test_ui/static/assets/index-ykctqrgN.css
  - frontend/tests/unit/components/DevAgentation.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
  - cmd/test_ui/static/assets/index-Bl1MOChG.js
  - frontend/tests/unit/pages/datalink/workbench-source-preview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx
  - cmd/test_ui/static/assets/index-Bw9ae6Dz.js
  - frontend/tests/unit/utils/appAgentationRemoval.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourcePreview.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx
  - frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx
  - cmd/test_ui/static/index.html
-->