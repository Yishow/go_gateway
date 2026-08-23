## Problem

前端完整 Vitest suite 的 timeout 與 failure identity 會隨執行漂移，現有證據不足以判定是測試本身、worker/resource saturation、生命週期洩漏或環境因素。拆分後基線為 181/938、928 pass、10 個 5 秒 timeout、wall 78.9 秒；拆分前為 913/938、25 failures、wall 149.5 秒，而 focused 201/201、standalone 22/22、targeted Modbus surface 5/5、canvas overlay 8/8、SourcePreview 2/2、SourceRuleTargetDatatype 3/3，以及 default/single-worker 四檔 18/18 均通過。現在需要可重跑且不預設固定 suite identity 的診斷契約，才能安全決定是否存在可修復的 test-local root cause。

## Root Cause

Root cause 尚未證實。最高優先假說是多 worker 下的資源或程序飽和，但 failure identity 漂移表示不能把任何一次出現的十個 suite 當成固定故障集合。現有 focused 與 group evidence 只能排除部分局部 assertion 問題，尚不足以排除 jsdom lifecycle、mock cleanup、worker/process/heap 壓力或執行環境差異。

## Proposed Solution

- 新增 diagnosis-only matrix harness，動態收集每次完整 suite 的 failure set，並以 normal settings、每檔 fresh isolated repeat、default group、single-worker group 及一次只改一個變因的 isolate/no-isolate、order、pool 矩陣建立可比較 evidence。
- 對每次觀測記錄 command、settings、failure identity、phase timing、process、heap、worker 與 timeout output，將結果分類為 assertion、test-timeout、worker-hang、resource leak 或 environment。
- 只有在單一變因矩陣證實 test-local lifecycle 或 cleanup root cause 後，才允許提出 bounded test-only repair；修復必須維持既有 title、assertion、mock、fixture、route 與 discovery contract。
- 若無法重現，或唯一有效手段是調高 timeout、改 worker/config、skip/retry/exclude 或降低 coverage，則保留完整 evidence 並將 diagnosis 標記為 blocked，不引入 workaround。

## Non-Goals

- 不提高 testTimeout 或 hookTimeout，不加入 skip、only、retry、exclude，不降低 coverage。
- 不修改 Vite/Vitest config、package manifest、lockfile、dependency、production frontend/src 或 Go production code。
- 不把目前觀察到的十個 timeout suite 固定成永久 identity，也不宣稱降低 worker 數量是產品修復。
- 不執行或代替現場 acceptance；本 change 只建立可重跑診斷與有條件的 test-only 修復門檻。

## Capabilities

### New Capabilities

- frontend-test-suite-stability: 以動態 failure-set、單變因矩陣與分類門檻診斷前端完整測試 timeout，並約束有證據的 bounded test-only repair。

### Modified Capabilities

- none

## Impact

- Affected specs: new frontend-test-suite-stability.
- Affected code:
  - New: scripts/diagnose-frontend-vitest-matrix.ps1
  - Modified (conditional only after root-cause proof):
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchRuntimePhase.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourcePreview.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-surface.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.batch-selection.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.canvas-overlays.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.conflict-queue.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inline-editing.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.planning.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.rule-actions.test.tsx
    - frontend/tests/unit/pages/datalink/workbench/MuiSourceIncidentDesk.reopen.test.tsx
    - frontend/src/setupTests.ts only if evidence proves shared test-setup cleanup is the root cause.
  - Removed: none

The eleven listed test paths are an observed cross-run union only, not a fixed suite identity. If dynamic evidence selects a root-relative path outside this list, stop implementation and run spectra-ingest to update the change before modifying that path; do not expand scope locally.

