# V2 Calm Summary Harvest Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Soften the noisiest `v2` Shell and Output summaries by harvesting selected `v1` calm-framing cues without changing `v2`'s incident-desk structure, action order, or domain behavior.

**Architecture:** Implement only in `.worktrees/woe-v2-mui`, touching the two `v2` summary wrapper components and the exact locale keys they render. Use test-first changes in the existing reopen test suites to lock ownership/action invariants, add translation-aware copy assertions, and then validate the result on the fixed Phase 6 browser fixture.

**Tech Stack:** React 19, TypeScript 5, MUI, React Query, Vitest, Testing Library, Vite

---

## File Map

- `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx`
  - `v2` Shell summary wrapper; owns readiness/blocker/refresh/return presentation only.
- `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx`
  - `v2` Output incident wrapper; owns the priority card, summary strip, command dock wrapper, and handoff panel.
- `.worktrees/woe-v2-mui/frontend/src/i18n/locales/en/common.json`
  - English copy for `workbench.shell.*` and `workbench.output.incident.*` keys touched by the calming pass.
- `.worktrees/woe-v2-mui/frontend/src/i18n/locales/zh-TW/common.json`
  - zh-TW copy for the same keys.
- `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx`
  - Shell reopen harness; currently mocks `t(key) => key`, so it must gain a tiny translation map before it can verify calmer copy.
- `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx`
  - Output reopen harness; same translation-map issue, plus this is where priority/handoff/surface invariants are already locked.
- `.worktrees/woe-v2-mui/frontend/tests/unit/pages/datalink/workbench-output-step.test.tsx`
  - Existing shared Output contract test; use it to confirm Local Modbus / Database behavior stays intact.
- `docs/superpowers/specs/2026-04-09-v2-calm-summary-harvest-design.md`
  - Approved spec; read-only during implementation.

## Implementation Notes

- Do **not** touch `LocalModbusBoard.tsx`, `DatabaseTargetBoard.tsx`, routing, or `WorkbenchProvider`.
- Keep `v2` obviously `v2`:
  - calm the copy and supporting hierarchy
  - do **not** import `v1` rail structure
  - do **not** remove blocker nouns or reorder command-dock actions
- Prefer updating the existing `vi.mock('react-i18next', ...)` blocks inside each reopen test file instead of creating new helper files.
- Keep commit scopes small: one commit for Shell, one commit for Output, one final validation/polish commit if browser-led tweaks are needed.

### Task 0: Bootstrap the v2 follow-up worktree

**Files:**
- Create if missing: `.worktrees/woe-v2-mui/`
- Verify branch: `woe-v2-mui`

- [x] **Step 1: Verify or create the worktree before touching files**

Run:

```bash
cd /Users/yishow/prj/go_gateway
if test -d .worktrees/woe-v2-mui; then
  echo "worktree exists"
else
  if git show-ref --verify --quiet refs/heads/woe-v2-mui; then
    git worktree add .worktrees/woe-v2-mui woe-v2-mui
  else
    git worktree add -b woe-v2-mui .worktrees/woe-v2-mui HEAD
  fi
fi
```

Expected: `.worktrees/woe-v2-mui` exists and is attached to branch `woe-v2-mui`.

- [x] **Step 2: Verify the worktree is clean enough to start**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui
git status --short
git branch --show-current
```

Expected: branch name is `woe-v2-mui`; if unrelated changes exist, stop and decide whether they are safe to work around before editing.

- [x] **Step 3: Confirm the target frontend files are reachable**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui
test -f frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx
test -f frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx
test -f frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
test -f frontend/src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx
```

Expected: all four checks succeed before Task 1 starts.

### Task 1: Calm the v2 Shell incident strip

**Files:**
- Modify: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx`
- Modify: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx`
- Modify: `.worktrees/woe-v2-mui/frontend/src/i18n/locales/en/common.json`
- Modify: `.worktrees/woe-v2-mui/frontend/src/i18n/locales/zh-TW/common.json`

- [x] **Step 1: Write the failing Shell test with a tiny translation map**

Replace the current `t: (key) => key` mock in `MuiWorkbenchShellIncidentDesk.reopen.test.tsx` with a local translation map for the touched Shell keys and add assertions for the calmer copy while keeping the same action contract.

```ts
const shellTranslations: Record<string, string> = {
  'workbench.shell.eyebrow': 'Incident context',
  'workbench.shell.labels.blocker': 'Current blocker',
  'workbench.shell.labels.refresh': 'Diagnostics status',
  'workbench.shell.blockers.outputPending':
    'Output still needs review. Check Local Modbus or Database, then apply the mapping.',
  'workbench.shell.refresh.success': 'Diagnostics refreshed',
  'workbench.shell.actions.returnOutput': 'Open output workspace',
};

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, vars?: Record<string, unknown>) => {
      if (key === 'workbench.shell.readinessSummary') {
        return `${vars?.readyCount}/${vars?.stepCount} steps ready · ${vars?.attentionCount} need attention`;
      }
      return shellTranslations[key] ?? key;
    },
  }),
}));
```

Add/update expectations so the test proves:

- the Shell still renders the same strip/action test IDs
- the blocker copy uses the calmer string above
- the return action still routes to Output
- refresh success still surfaces explicit success text

- [x] **Step 2: Run the Shell test to verify it fails**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui/frontend
npm test -- --run src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
```

Expected: FAIL because the current component/locale strings do not match the new calmer copy assertions.

- [x] **Step 3: Implement the minimal Shell change**

Update `MuiWorkbenchIncidentStrip.tsx` and the two locale files.

Concrete implementation targets:

- keep both buttons in the same order
- keep the four Shell-owned surfaces
- lower non-primary emphasis in the strip:
  - `shell-readiness-summary` should use a quieter text color than the blocker line
  - `shell-refresh-status` should no longer compete visually with the blocker except for active/error states
- update only the Shell keys listed in the spec scope

Representative component change:

```tsx
<Typography
  data-testid="shell-readiness-summary"
  sx={{ ...SENTRY_SX.monoData, fontSize: '11px', color: tokens.text.secondary }}
>
  {t('workbench.shell.readinessSummary', readiness)}
</Typography>
```

Representative locale edits:

```json
"eyebrow": "Incident context",
"blocker": "Current blocker",
"refresh": "Diagnostics status",
"outputPending": "Output still needs review. Check Local Modbus or Database, then apply the mapping.",
"success": "Diagnostics refreshed"
```

- [x] **Step 4: Run the Shell test to verify it passes**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui/frontend
npm test -- --run src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx
```

Expected: PASS with the calmer copy rendered and the same return/refresh behavior intact.

- [x] **Step 5: Commit the Shell pass**

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui
git add \
  frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx \
  frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx \
  frontend/src/i18n/locales/en/common.json \
  frontend/src/i18n/locales/zh-TW/common.json
git commit -m "[v2-followup][Task 1] 收斂 Shell incident strip 文案與層級" \
  -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 2: Calm the v2 Output summary strip and handoff panel

**Files:**
- Modify: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx`
- Modify: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx`
- Modify: `.worktrees/woe-v2-mui/frontend/src/i18n/locales/en/common.json`
- Modify: `.worktrees/woe-v2-mui/frontend/src/i18n/locales/zh-TW/common.json`

- [x] **Step 1: Write the failing Output test with translation-aware copy assertions**

Update the Output reopen harness so it no longer returns raw keys for the touched incident strings. Add a minimal translation map and failing assertions for the calmer summary/handoff layer while keeping command-dock behavior locked.

Use these concrete copy targets:

```ts
const outputTranslations: Record<string, string> = {
  'workbench.output.incident.summaryEyebrow': 'Incident context',
  'workbench.output.incident.metrics.linked': 'linked',
  'workbench.output.incident.metrics.scoped': 'scoped',
  'workbench.output.incident.metrics.attention': 'needs review',
  'workbench.output.incident.handoff.label': 'Next step',
  'workbench.output.incident.handoff.noRule':
    'Return to Tag to restore the active source rule, then continue output review.',
  'workbench.output.incident.handoff.ready':
    'The desk context is aligned with the shared workboard for this target.',
};
```

Also add structural assertions for the small layout adjustment by introducing two explicit containers:

- `output-incident-summary-primary`
- `output-incident-summary-secondary`

The failing test should prove:

- the priority card still exists and remains the first review surface
- linked/scoped context stays in the primary summary row
- attention/revision context moves into the secondary summary row
- `Go to Tag` still exists and keeps the repair path intact

- [x] **Step 2: Run the Output test to verify it fails**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui/frontend
npm test -- --run src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx
```

Expected: FAIL because the current Output incident desk still uses the old copy and one-row summary strip.

- [x] **Step 3: Implement the minimal Output change**

Update `MuiOutputIncidentDesk.tsx` and the two locale files.

Concrete implementation targets:

- keep the same four-part wrapper:
  - priority card
  - summary strip
  - command dock
  - handoff panel
- keep command-dock buttons and target-switch behavior unchanged
- split the summary strip into two grouped rows:
  - primary row: linked + scoped
  - secondary row: attention + revision
- soften only the wrapper copy:
  - calmer summary eyebrow
  - calmer handoff label and repair copy
  - calmer attention chip wording

Representative structure:

```tsx
<Box data-testid="output-incident-summary-strip" ...>
  <Typography ...>{t('workbench.output.incident.summaryEyebrow')}</Typography>
  <Stack data-testid="output-incident-summary-primary" direction="row" ...>
    {/* linked + scoped chips */}
  </Stack>
  <Stack data-testid="output-incident-summary-secondary" direction="row" ...>
    {/* attention + revision chips */}
  </Stack>
</Box>
```

- [x] **Step 4: Run Output tests to verify the wrapper still works**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui/frontend
npm test -- --run \
  src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx \
  tests/unit/pages/datalink/workbench-output-step.test.tsx
```

Expected: PASS; the incident wrapper copy/layout is calmer while Local Modbus / Database behavior remains unchanged.

- [x] **Step 5: Commit the Output pass**

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui
git add \
  frontend/src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx \
  frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx \
  frontend/src/i18n/locales/en/common.json \
  frontend/src/i18n/locales/zh-TW/common.json
git commit -m "[v2-followup][Task 2] 收斂 Output incident summary 與 handoff 文案" \
  -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 3: Run the full v2 validation and browser gate

**Files:**
- Modify: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx` (only if browser validation exposes a small Shell polish gap)
- Modify: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx` (only if browser validation exposes a small Output polish gap)
- Test: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx`
- Test: `.worktrees/woe-v2-mui/frontend/src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx`
- Test: `.worktrees/woe-v2-mui/frontend/tests/unit/pages/datalink/workbench-output-step.test.tsx`

- [x] **Step 1: Run the combined automated validation**

Run:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui/frontend
npm test -- --run \
  src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx \
  src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx \
  tests/unit/pages/datalink/workbench-output-step.test.tsx
npx tsc --noEmit
npx eslint \
  src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx \
  src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx
npm run build
```

Expected: all commands pass.

- [x] **Step 2: Reproduce the fixed browser fixture**

Use `agent-browser` against:

- URL: `http://127.0.0.1:4175/studio`
- device: `UI 4.3 Modbus TCP`
- path: `Device -> Source -> Tag -> Output -> Shell`
- source-rule context: `40001 · MBT · int16`
- revision context: `82b8eea1-0876-4284-a9d7-e15db8779ca3`

Expected: the same fixture used by the approved spec is available; if not, stop and resolve the fixture instead of substituting a different path.

- [x] **Step 3: Capture the required browser evidence**

Capture at minimum:

- `v2-followup-shell-summary.png`
- `v2-followup-output-local-modbus-summary.png`
- `v2-followup-output-database-summary.png`

The screenshots must show:

- Shell still reads blocker-first
- Output still leads with the priority card
- the calmer summary/handoff copy does not make the screen read like `v1`

- [x] **Step 4: Apply only minimal browser-led polish if needed**

If browser validation exposes a small wording/spacing issue, fix it only in the two wrapper components already in scope, then rerun:

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui/frontend
npm test -- --run \
  src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx \
  src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx \
  tests/unit/pages/datalink/workbench-output-step.test.tsx
npm run build
```

Expected: PASS again after the small polish fix.

- [x] **Step 5: Commit the validated follow-up**

```bash
cd /Users/yishow/prj/go_gateway/.worktrees/woe-v2-mui
git add \
  frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx \
  frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx \
  frontend/src/i18n/locales/en/common.json \
  frontend/src/i18n/locales/zh-TW/common.json \
  frontend/src/pages/datalink/workbench/__tests__/MuiWorkbenchShellIncidentDesk.reopen.test.tsx \
  frontend/src/pages/datalink/workbench/__tests__/MuiOutputIncidentDesk.reopen.test.tsx
git commit -m "[v2-followup][Task 3] 完成 calm summary harvest 驗證與收尾" \
  -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

Completed outcome: browser gate surfaced a shared backend blocker in `internal/datalink/mapping/sql_repo.go`, so the final validating commit was `ba75251` (`[v2][Task 3] 修正 mappings 可空欄位掃描並完成瀏覽器驗收`). No extra wrapper polish commit was needed after the real fixture recovered.

## Done Checklist

- [x] Shell / Output wrapper surfaces stayed unchanged during Task 3 browser gate; the only extra code change was the shared backend mappings nullable-scan fix needed to restore the real fixture
- [x] `v2` still reads as incident-desk first
- [x] Shell actions are unchanged
- [x] Output command-dock actions are unchanged
- [x] Local Modbus and Database flows still pass shared tests
- [x] Browser evidence captured on the fixed `4175` Phase 6 fixture
