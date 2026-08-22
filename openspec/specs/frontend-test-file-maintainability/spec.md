# frontend-test-file-maintainability Specification

## Purpose

本規格定義前端大型單元測試的拆分與驗收契約，確保行為檔可被 Vitest 探索、共享 test-only harness 不改變測試語意、所有產出遵守 300 行上限，並將完整測試的既有 timeout 與 focused correctness 分開量測與回報。

## Requirements

### Requirement: Oversized test suites are split into discoverable behavior files

The implementation SHALL replace each of the nine named frontend unit-test monoliths with behavior-oriented .test.ts or .test.tsx files under frontend/tests/unit. The implementation SHALL remove the original monolith after its complete test-title set has been verified in the destination files. Every resulting TypeScript or TSX test or test-only harness file SHALL be at most 300 physical lines.

#### Scenario: all nine monoliths have bounded discoverable replacements

- **WHEN** the split is complete for DatalinkWorkbenchFoundation, DatalinkWorkbenchOutputStep, DatalinkWorkbenchShellUi, DatalinkWorkbenchSourceStep, DatalinkWorkbenchTagStep, MuiOutputIncidentDesk.reopen, database-autosave-page, device-autosave-page, and mapping-autosave-page
- **THEN** each original monolith is absent, each behavior destination remains under frontend/tests/unit, Vitest discovers every .test.ts or .test.tsx destination, and every resulting TS/TSX file is at most 300 physical lines

##### Example: source and destination boundary

- **GIVEN** frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.test.tsx is replaced by its planning, inspector, rule-actions, canvas-overlays, batch-selection, inline-editing, conflict-queue, datatype-span, and visual-span destinations
- **WHEN** the destination test glob is collected
- **THEN** every original test title occurs exactly once across those destinations and the source monolith does not remain

#### Scenario: behavior ownership is explicit for each destination

- **WHEN** an implementer assigns tests from any named monolith to destination files
- **THEN** each destination has one behavior ownership boundary, no test title is duplicated, and all nine source title sets have an equal destination title set before source removal


<!-- @trace
source: split-oversized-frontend-tests
updated: 2026-08-23
tests:
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.legacy-routes.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.operator-flow.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-list.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-editor.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-surface.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-bindings.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-canvas.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-targets.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-schema.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-schema-columns.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.output-inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.cross-step-handoff.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.frame.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.context.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.summary.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.switching.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.planning.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.rule-actions.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.canvas-overlays.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.batch-selection.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inline-editing.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.conflict-queue.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.datatype-span.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.visual-span.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.review-overview.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.review-selection.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.tag-binding.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.tag-results.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.command.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.recovery.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.validation.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.bulk-toggle.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.bulk-toggle.test.tsx
-->

### Requirement: Shared test-only harness preserves the test contract

The implementation SHALL place reusable typed render helpers, router wrappers, fixture factories, mock registration, and cleanup in test-only .testHarness.tsx files where sharing reduces duplication. Harness files SHALL NOT register standalone tests, and moved tests SHALL preserve their test titles, assertions, fixture values, mock request shapes, route transitions, and cleanup semantics. Harness extraction SHALL NOT add production imports or change application behavior.

#### Scenario: moved tests retain observable assertions and setup

- **WHEN** a behavior test is moved from one of the nine monoliths into a destination file
- **THEN** the test keeps its original title and assertions, receives equivalent typed fixtures and mocks from local test-only setup, and passes the focused destination test command without a semantic assertion rewrite

#### Scenario: harness files are not independent test entries

- **WHEN** Vitest scans frontend/tests/unit after the split
- **THEN** files named .testHarness.tsx provide imported helpers only, contain no top-level test registration, and do not create duplicate test execution

#### Scenario: mock and fixture state remains isolated

- **WHEN** two destination files import the same local harness
- **THEN** each test resets mutable mock and fixture state according to the original cleanup contract, and a focused failure in one destination does not depend on execution order in another destination


<!-- @trace
source: split-oversized-frontend-tests
updated: 2026-08-23
tests:
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.testHarness.tsx
-->

### Requirement: Line-limit and frontend quality gates remain enforceable

The implementation SHALL keep all resulting TypeScript and TSX test or harness files at most 300 physical lines without changing .line-limit-ignore, test discovery configuration, or timeout settings. The implementation SHALL run the focused moved tests for each group and the repository-supported line gate, frontend lint, frontend test, and frontend build commands.

#### Scenario: line gate rejects oversized destinations

- **WHEN** any resulting test or harness file exceeds 300 physical lines
- **THEN** the change remains incomplete and the owning behavior group is subdivided before source removal; no ignore entry or line-limit exception is added

#### Scenario: quality commands cover the moved test contract

- **WHEN** all nine groups have passed their focused destination tests
- **THEN** the implementer runs cd frontend && npm run lint, cd frontend && npm run test, cd frontend && npm run build, and make check-lines or the repository-supported equivalent, recording each command result

#### Scenario: production and test configuration remain unchanged

- **WHEN** the final diff is reviewed
- **THEN** all changed source-code paths are under frontend/tests/unit, no frontend/src or Go production path is changed, and npm/Vitest/Vite configuration, dependency manifests, timeout settings, and .line-limit-ignore are unchanged


<!-- @trace
source: split-oversized-frontend-tests
updated: 2026-08-23
tests:
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.legacy-routes.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.operator-flow.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-list.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-editor.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-surface.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-bindings.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-canvas.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-targets.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-schema.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-schema-columns.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.output-inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.cross-step-handoff.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.frame.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.context.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.summary.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.switching.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.planning.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inspector.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.rule-actions.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.canvas-overlays.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.batch-selection.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inline-editing.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.conflict-queue.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.datatype-span.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.visual-span.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.review-overview.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.review-selection.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.tag-binding.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.tag-results.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.command.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.recovery.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.validation.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.bulk-toggle.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.bulk-toggle.test.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.testHarness.tsx
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.testHarness.tsx
  - frontend/tests/unit/workbench-v2/database-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.testHarness.tsx
-->

### Requirement: Full-suite timeout is measured separately from the split

The implementation SHALL record the pre-split and post-split full frontend test duration and exact timeout or failure output when a Workbench full-suite timeout exists. The implementation SHALL report that evidence separately from focused split correctness and SHALL NOT increase a blanket timeout, skip tests, or claim that the split resolves an unrelated baseline timeout.

#### Scenario: existing timeout remains an explicit baseline

- **WHEN** the full frontend test command times out before or after the split
- **THEN** the report includes the exact command, elapsed or configured timeout, affected suite, and failure output, while focused destination results remain separately classified

#### Scenario: no timeout workaround is introduced

- **WHEN** the split is implemented and validated
- **THEN** no timeout setting, test skip, discovery exclusion, or unrelated performance workaround is changed, and any follow-up timeout remediation is identified as a separate change

<!-- @trace
source: split-oversized-frontend-tests
updated: 2026-08-23
-->
