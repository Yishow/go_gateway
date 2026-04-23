## 1. Contract-first desk coverage

- [x] 1.1 Implement **Lock desk contracts with targeted source-step tests** by updating `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx` and `frontend/tests/unit/pages/datalink/workbench-source-step.test.tsx` to cover **Source desk modes share one planning skeleton**, **Source desk switching preserves guided-flow continuity**, and **Source workspace keeps a persistent desk skeleton**.
- [x] 1.2 Update `frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx` so **Triage desk is issue-first and explicit when clear** is verified for both recovery and no-issue states.

## 2. Shared Source-step shell refactor

- [x] 2.1 Implement **Preserve one shared source-step planning skeleton** in `frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx`, `frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`, and `frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx` so the shared summary, handoff strip, diagnostics, and primary action stay anchored across desk switches.
- [x] 2.2 Implement **Split mode-specific panels out of SourceCanvasSection** by introducing `frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx` and `frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx`, wiring them through the shared Step 2 orchestrator without duplicating Source-rule state.

## 3. Distinct desk-mode layouts and behavior

- [x] 3.1 Implement **Inspect and Build desks expose distinct working emphasis**, **Keep build mode editor-primary without collapsing the canvas**, and **Source desk layouts allocate distinct primary regions on desktop** in `frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx` and the Step 2 desk panels.
- [x] 3.2 Implement **Make triage mode issue-first with filtered recovery cues** and **Triage desk is issue-first and explicit when clear** by surfacing recovery actions, filtered canvas context, and a dedicated clear-state message in the new triage panel.

## 4. Workflow continuity, copy, and verification

- [x] 4.1 Implement **Source desk switching preserves guided-flow continuity** by updating the Step 2 handoff, blocker, and CTA copy in `frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx`, `frontend/src/i18n/locales/zh-TW/common.json`, and `frontend/src/i18n/locales/en/common.json`.
- [x] 4.2 Run `cd frontend && npm run lint && npm run test -- --run && npm run build` and verify **Source desk modes share one planning skeleton** after switching between `Inspect`, `Build`, and `Triage` in `/studio`.
