# datalink-workbench-desktop Specification

## Purpose
TBD - created by archiving change redesign-datalink-workbench-desktop-flow. Update Purpose after archive.

## Requirements

### Requirement: Desktop workbench shell
The system SHALL provide a desktop-first workbench shell at `/studio` as the primary product route with persistent workbench regions for step navigation, context, main review surfaces, inspection, and summary feedback.

Any legacy `/datalink/workbench` entry SHALL redirect to `/studio` instead of remaining a separate primary workflow.

#### Scenario: Primary workbench opens through `/studio`
- **WHEN** the operator opens the primary datalink workbench
- **THEN** the system loads the desktop-first shell at `/studio`
- **AND** does not require a separate primary `/datalink/workbench` route for the same workflow

---
### Requirement: Device workspace establishes capability context
The system SHALL provide a `DeviceWorkspace` that lets the operator select or create a source device and understand the capability context required for source planning.

#### Scenario: Selected device exposes planning-critical capabilities
- **WHEN** the operator selects a device in Step 1
- **THEN** the UI shows `address base`, `word order`, `unit id`, and protocol traits in the workbench context
- **AND** those capability hints remain visible in the device inspector and context bar

#### Scenario: Clone flow starts from selected device defaults
- **WHEN** the operator chooses to clone the selected device
- **THEN** the UI opens a workbench-specific create flow prefilled with reusable connection defaults
- **AND** still requires a distinct device identity before save

---
### Requirement: Address canvas supports rule-driven continuous planning
The system SHALL render source planning as a continuous 16-bit lattice that can expand multiple rules into merged spans on one shared canvas.

#### Scenario: Typed rules expand into merged spans
- **WHEN** the operator applies a rule with `32-bit`, `float32`, or `64-bit` data width
- **THEN** the canvas merges the correct number of adjacent 16-bit cells per item
- **AND** the operator can visually distinguish merged spans from empty cells

#### Scenario: Multiple rules share one continuous address range
- **WHEN** the operator defines multiple source rules with different start addresses
- **THEN** the canvas spans from the minimum rule start address to the maximum rule end address
- **AND** gaps between rule-covered ranges remain visible as unplanned address space

---
### Requirement: Address canvas uses stable multi-mode visualization
The system SHALL let the operator switch between `Plan`, `Live`, and `Link` modes without replacing the underlying address geometry.

#### Scenario: View-mode switching preserves geometry
- **WHEN** the operator switches between `Plan`, `Live`, and `Link`
- **THEN** the address positions, merged spans, and current selection remain geometrically stable
- **AND** only the information layer changes between planning, value validation, and downstream linkage states

#### Scenario: Toolbar exposes live validation controls
- **WHEN** the operator uses Step 2 toolbar controls
- **THEN** the toolbar provides value-format switching, live freeze, snapshot compare, jump-to-address, and template load/save actions
- **AND** those controls apply to the current canvas state without leaving the workbench route

---
### Requirement: Tag binding board preserves source-to-tag review context
The system SHALL provide a `TagBindingBoard` that exposes source context, proposed tag identity, and binding status densely enough for batch review.

#### Scenario: Source-to-tag information is visible before inspector open
- **WHEN** the operator reviews items in Step 3
- **THEN** each row or card shows tag key preview, source address/span, raw vs transformed value, data width context, and binding status
- **AND** the operator can compare multiple items without opening them one by one

#### Scenario: Batch review includes diff and result feedback
- **WHEN** the operator performs a batch tag action
- **THEN** the UI shows a pre-apply diff preview before confirmation
- **AND** shows a post-apply result summary with created, linked, skipped, and failed items

---
### Requirement: Unified output workspace supports both targets
The system SHALL provide one Output workspace inside `/studio` that keeps Local Modbus and Database outputs in the same workbench step.

#### Scenario: Operator changes output target without leaving `/studio`
- **WHEN** the operator switches between `Local Modbus` and `Database` in the Output step
- **THEN** the output review surface remains in the same `/studio` route
- **AND** the target-specific studio changes without requiring navigation to a separate primary page

#### Scenario: Output target switching preserves one primary workflow
- **WHEN** legacy navigation or bookmarks would previously open a separate workbench route for output work
- **THEN** the system keeps the operator inside `/studio`
- **AND** treats any retained legacy entry as a compatibility redirect rather than a second primary workflow

---
### Requirement: Device testing distinguishes transport and protocol phases
The workbench SHALL show and judge transport connectivity and protocol probe results as separate outcomes in Step 1.

#### Scenario: Connect succeeds but probe fails
- **WHEN** an operator tests a device and transport connect succeeds while protocol probe fails
- **THEN** Step 1 SHALL show connect as successful and probe as failed
- **AND** SHALL allow the device to be saved
- **AND** SHALL mark rule activation and data collection as blocked until probe succeeds

#### Scenario: Route failure is classified as connect-stage failure
- **WHEN** a device test fails with a network error such as `no route to host`
- **THEN** Step 1 SHALL classify the failure as a connect-stage failure
- **AND** SHALL present diagnostics without mislabeling the protocol probe as the failing stage

---
### Requirement: Source planning reflects device capability and persisted rule state
The workbench SHALL align Step 2 planning behavior with the selected device capability context and persisted rule state.

#### Scenario: Device capability updates planner behavior
- **WHEN** an operator selects a device with different address model or protocol traits
- **THEN** Step 2 SHALL update available rule options, address semantics, merge behavior, and warning/block behavior to match that device

#### Scenario: Existing unmanaged points are explained in the grid
- **WHEN** Step 2 loads points that exist without an active persisted rule association
- **THEN** the grid SHALL distinguish them from rule-planned spans
- **AND** SHALL explain whether each address is unmanaged used state, conflict state, or rule-derived state

---
### Requirement: Tag step is review-first for rule-derived mappings
The workbench SHALL treat Step 3 as a review and exception-handling surface for rule-derived tag candidates and pending mapping intent, not as an immediate active tag-and-mapping persistence step.

#### Scenario: Rule creation pre-populates review surface
- **WHEN** an operator creates or restores a rule in the primary flow
- **THEN** Step 3 SHALL load the generated tag candidates and pending mapping intent without requiring a manual first-pass bind
- **AND** SHALL show the generated status for operator review

#### Scenario: Review step supports exception handling
- **WHEN** automatic tag generation needs correction or approval
- **THEN** Step 3 SHALL surface rename, skip, override, and apply actions inline
- **AND** SHALL allow corrective review actions without reverting the entire workflow to manual binding

---
### Requirement: Tag review consumes source-rule revision candidates
The workbench SHALL bind Step 3 to tag candidates derived from the active source-rule revision.

#### Scenario: Tag review uses the active revision
- **WHEN** the operator opens Step 3 for a source rule
- **THEN** the workbench loads the tag candidate set for the current active revision
- **AND** does not mix candidates from older revisions into the same review session

---
### Requirement: Tag review detects stale candidate views
The workbench SHALL detect when Step 3 is showing candidates from an older source-rule revision.

#### Scenario: Stale tag review blocks apply
- **WHEN** a newer source-rule revision exists for the rule currently open in Step 3
- **THEN** the workbench marks the review state as stale
- **AND** blocks tag apply until the operator refreshes to the latest candidate revision

---
### Requirement: Source rule builder exposes target data type selection
The system SHALL provide, within Step 2 source planning, a control to set the **target / intended tag data type** alongside the protocol read data type when creating or editing a source rule. The control SHALL default to matching the protocol read data type and SHALL clearly distinguish protocol read semantics from target output semantics.

#### Scenario: Operator sets distinct target type before applying rule
- **WHEN** an operator configures a new rule and selects a target data type different from the protocol read type
- **THEN** the workbench SHALL persist that intent with the rule according to backend contract
- **AND** SHALL show validation feedback when the pair is unsupported

#### Scenario: Defaults avoid extra friction
- **WHEN** an operator does not change the default target data type
- **THEN** the workbench SHALL behave as today with a single data type selection driving protocol read planning

---
### Requirement: Device connection form exposes data format when the protocol uses multi-register decoding
The workbench SHALL render **`data_format`** (or the documented equivalent) in the device connection editor for each protocol whose connector uses multi-register / floating-point byte ordering (for example **Modbus TCP, Modbus RTU, Modbus UDP**, and **Mitsubishi MC 3E**). The control SHALL list the supported orderings (for example **ABCD, BADC, CDAB, DCBA**) and SHALL persist the value in `connection_config` for the backend to consume during reads.

#### Scenario: Modbus connection shows byte order
- **WHEN** an operator edits a Modbus-class device connection
- **THEN** the UI SHALL expose `data_format` selection
- **AND** saving the device SHALL persist the chosen value

---
### Requirement: Source rule builder MAY expose scale or offset for unit conversion
The system SHALL provide optional inputs for **linear scaling** (multiplier and/or offset) on the source rule form when the product enables engineering-unit conversion from rule planning, with clear labels that distinguish scaling from protocol read type and from target tag type.

#### Scenario: Operator sets scale without changing target type
- **WHEN** an operator enters scale parameters only
- **THEN** the workbench SHALL persist scale intent with the rule
- **AND** SHALL not require a target type different from the read type

---
### Requirement: Source planning surfaces grouped-tag handoff context

The workbench SHALL show, in Step 2, enough Source planning context for the operator to understand how planned points will move into grouped Tag review and downstream Database planning.

Step 2 SHALL treat the active source rule as the primary planning unit and SHALL present a rule-first summary before the larger canvas surface. The canvas MAY remain available for detailed address inspection, but it SHALL NOT be the only place where the operator can understand the current plan.

#### Scenario: Source summary explains grouped handoff
- **WHEN** an operator plans source points for a rule
- **THEN** Step 2 shows the planned point count and current naming-prefix context
- **AND** the handoff copy explains that the next step will review grouped tags before database planning

#### Scenario: Active rule summary precedes canvas reading
- **WHEN** an operator focuses a source rule in Step 2
- **THEN** the workbench shows that rule’s address coverage, rule-scoped planning status, and data-type intent before the operator reads the full lattice canvas
- **AND** the operator does not need to infer the current rule only from canvas cells or secondary tabs


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
### Requirement: Tag review exposes grouped database suggestions before Output

The workbench SHALL show, in Step 3, the inferred database grouping context for rule-derived tags before the operator moves into Database planning.

#### Scenario: Tag review shows inferred grouping metadata
- **WHEN** the operator opens Step 3 for slash-based tags such as `meter/A1` and `meter/kw`
- **THEN** the Tag review surface shows inferred `group_key`, suggested `column_name`, and group membership for each candidate
- **AND** the operator can understand which tags are expected to land in the same database row

#### Scenario: Tag overrides carry forward to Output planning
- **WHEN** the operator edits grouped row suggestions in Step 3
- **THEN** the chosen group and column overrides are preserved for the Database planner in Output
- **AND** the operator does not have to re-enter the same override state from scratch

---
### Requirement: Output mainline completes when either target is configured

The workbench SHALL treat the Output step as mainline-complete when either Local Modbus or Database is configured, while still keeping both target states visible in the same workspace.

#### Scenario: Database apply satisfies the mainline
- **WHEN** at least one database output mapping has been applied for the active source rule and Local Modbus is still incomplete
- **THEN** the Output step is reported as complete for mainline progression
- **AND** the UI still shows Local Modbus as incomplete instead of hiding it

#### Scenario: Local Modbus apply satisfies the mainline
- **WHEN** at least one Local Modbus binding has been applied for the active source rule and Database is still incomplete
- **THEN** the Output step is reported as complete for mainline progression
- **AND** the UI still shows Database planning status and required follow-up work

---
### Requirement: Recovery cues point to the true blocked surface

The workbench SHALL align the step rail, context bar, and shell recovery action to the same blocker truth after grouped Output planning changes.

#### Scenario: Locked steps explain blockers without acting as a second router
- **WHEN** a downstream step is blocked by missing grouped Tag review or Output configuration
- **THEN** the step rail explains the blocking reason without behaving like an alternate forward-action router
- **AND** the shell return action points to the actual recovery surface

---
### Requirement: Source workspace keeps a persistent desk skeleton
The workbench SHALL reserve a persistent Step 2 desk skeleton for the active source rule before rendering desk-specific primary surfaces.

The persistent skeleton SHALL include the active rule summary, grouped handoff strip, desk selector, blocker diagnostics, and current primary action so the desktop workspace keeps a stable reading order across desk changes.

#### Scenario: Shared desk skeleton remains anchored across desk changes
- **WHEN** the operator switches the Source workspace between `Inspect`, `Build`, and `Triage`
- **THEN** the same summary and handoff band remains anchored in the desktop shell above the mode-specific workspace
- **AND** the operator does not lose the current rule, planning status, or blocker explanation

#### Scenario: Desk skeleton keeps the primary action in a stable location
- **WHEN** the operator edits, reviews, or triages the active source rule
- **THEN** the desktop workspace keeps the current primary action in the same Step 2 shell region
- **AND** the operator does not need to scan different panels to find the next planning action after each desk switch


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
### Requirement: Source desk layouts allocate distinct primary regions on desktop
The workbench SHALL allocate distinct desktop-first primary regions for `Inspect`, `Build`, and `Triage` instead of reusing one nearly identical layout.

`Inspect` SHALL prioritize the source canvas, `Build` SHALL prioritize the rule editor while keeping a readable secondary canvas region, and `Triage` SHALL prioritize the incident queue and affected recovery context.

#### Scenario: Build keeps a readable secondary canvas region
- **WHEN** the operator opens `Build` on a desktop breakpoint
- **THEN** the rule editor becomes the primary panel
- **AND** the remaining canvas region stays large enough to inspect current address coverage without collapsing into a token-height preview

#### Scenario: Triage elevates incident review above general canvas browsing
- **WHEN** the operator opens `Triage` with pending planning issues
- **THEN** the incident queue and recovery controls occupy the primary desktop region
- **AND** the canvas is presented only as supporting context for the affected ranges instead of as the dominant browsing surface

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
### Requirement: Source step previews rule-scoped Tag review before navigation
The workbench SHALL provide, within Step 2, a rule-scoped preview of the next Tag review result for the currently focused source rule.

The preview SHALL reuse the same active-rule continuity that Step 3 consumes, so the operator can trust that the Source-step preview and Tag-step review are describing the same candidate scope.

#### Scenario: Clicked rule previews next Tag review scope
- **WHEN** an operator selects a source rule in Step 2
- **THEN** the workbench shows a preview of the Tag candidates that the next step will review for that rule
- **AND** the preview remains bound to the same focused rule when the operator navigates to Step 3

#### Scenario: Hovered rule temporarily previews downstream scope
- **WHEN** an operator hovers a different source rule in Step 2
- **THEN** the workbench MAY temporarily preview that rule’s Tag-review scope
- **AND** leaving the hover state restores the previously selected rule preview unless the operator explicitly changes focus


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
### Requirement: Source step prioritizes planning actions over mode switching
The workbench SHALL make applying the current rule to the planning surface the primary action in Step 2, while keeping mode switching as a secondary control.

#### Scenario: Apply planning is the primary action
- **WHEN** an operator is editing or reviewing a source rule in Step 2
- **THEN** the most prominent action communicates applying the rule to the planning surface
- **AND** the workbench gives immediate confirmation that the rule has taken shape for the next review step

#### Scenario: Mode tabs remain secondary
- **WHEN** Step 2 renders `plan`, `live`, or `link` modes
- **THEN** those controls remain available without overpowering the active-rule summary or primary apply action
- **AND** the operator can still identify the current plan without switching modes first

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
### Requirement: Coexisting v2 workbench at `/studio/v2`

The system SHALL accept `/studio/v2` as a parallel workbench route alongside the existing `/studio` primary route. The existing `/studio` route MUST continue to load the legacy `DatalinkWorkbenchPage` without modification. The `/studio/v2` route MUST load the Workbench v2 shell defined by capability `datalink-workbench-v2-shell`. Neither route MUST redirect to the other automatically.

#### Scenario: /studio remains the primary route

- **WHEN** the operator navigates to `/studio`
- **THEN** the system loads the pre-existing `DatalinkWorkbenchPage` shell
- **AND** does not redirect to `/studio/v2`

#### Scenario: /studio/v2 loads the v2 shell

- **WHEN** the operator navigates to `/studio/v2`
- **THEN** the system loads the Workbench v2 shell
- **AND** does not redirect to `/studio`

#### Scenario: Legacy datalink routes still redirect to /studio

- **WHEN** the operator navigates to `/datalink/workbench` or any `/datalink/workbench/*` path
- **THEN** the system redirects to `/studio` exactly as before this change
- **AND** does not redirect to `/studio/v2`

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