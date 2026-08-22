## Why

目前有 9 個 frontend 單元測試檔案超過檔案行數硬上限，最大檔案達 1902 行，造成 line-limit gate 阻擋、測試 setup 與 fixture 維護困難，也提高局部修改的回歸風險。現在拆分可在不改變產品行為、測試契約、斷言或測試探索路徑的前提下，讓每個測試單位可獨立理解與驗證。

## What Changes

- 將以下 9 個超過 500 行的測試 monolith 依行為與責任拆成可獨立探索的測試檔與 test-only typed harness：
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.test.tsx（1084 行）
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.test.tsx（1353 行）
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.test.tsx（592 行）
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.test.tsx（1902 行）
  - frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.test.tsx（985 行）
  - frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.test.tsx（633 行）
  - frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx（700 行）
  - frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx（618 行）
  - frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx（729 行）
- 移除上述原始 monolith，保留所有測試名稱、斷言、mock、fixture、路由與測試探索能力；共建立必要的 test-only harness 與行為分割檔，每個產出 TypeScript/TSX 檔案目標不超過 300 行。
- 以每一組原始測試先做 characterization，再執行 focused moved tests、line-limit gate、frontend lint、完整測試與 build，建立可重跑的等價性證據。
- 將前端測試檔可維護性、每檔行數上限、共享 harness 邊界與測試契約保留方式記錄為新的 capability 規格。
- 單獨量測並記錄 Workbench full-suite 的既有 timeout 或失敗基線；不以提高 timeout 或修改測試語意取代拆分。

## Non-Goals

- 不修改 frontend/src、Go backend、API、runtime、資料庫、協議或任何 production code。
- 不修改測試的語意、斷言、測試名稱、fixture 資料或 mock 契約；只允許拆檔與必要的 import/harness 調整。
- 不新增或修改 .line-limit-ignore，不以提高 timeout、跳過測試或改變 discovery 設定繞過問題。
- 不保證拆分本身消除既有 full-suite timeout；該問題須以基線證據另行評估。
- 不把測試移出 frontend/tests/unit，也不建立平行的 production 測試入口。

## Capabilities

### New Capabilities

- frontend-test-file-maintainability: 定義前端單元測試檔案的行數、拆分、test-only harness、測試契約保留與品質 gate。

### Modified Capabilities

- 無。這是測試組織與可維護性契約的新增，不改變既有產品或功能需求。

## Impact

- 主要影響 frontend/tests/unit/pages/datalink/workbench/ 下 6 個 Workbench 測試檔，以及 frontend/tests/unit/workbench-v2/ 下 3 個 autosave 測試檔。
- 會新增同目錄的 53 個 test-only TS/TSX harness 或行為分割檔，並移除 9 個原始 monolith；完整 ownership、destination 與測試群組會在 design.md 固定。
- 會新增 openspec/specs/frontend-test-file-maintainability/spec.md；不修改現有 capability spec。
- 不增加 runtime dependency，不改變 npm scripts、Vite/Vitest 設定、production bundle 或對外介面。
