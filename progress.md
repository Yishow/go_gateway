# Progress

## 2026-03-08

### Phase 2 後續實作 - 補齊表單屬性 [進行中]
- 補齊 `SmartDashboardWorkspaceContent` 表單屬性：
  - 起始位址 input：`name={designSystem.forms.name.pointAddress}`、`autoComplete={designSystem.forms.autocomplete.text}`、`inputMode={designSystem.forms.inputmode.text}`
  - 數量 input：`name="plan-count"`、`autoComplete={designSystem.forms.autocomplete.number}`、`inputMode={designSystem.forms.inputmode.numeric}`
  - Modbus 區域 select：`name="modbus-area"`
  - 資料型別 select：`name={designSystem.forms.name.pointDataType}`
- 補齊 `SmartDashboardWorkflowModal` 表單屬性：
  - 搜尋 input：`name={designSystem.forms.name.searchQuery}`、`autoComplete={designSystem.forms.autocomplete.search}`、`inputMode={designSystem.forms.inputmode.search}`
- 補齊 `LocalModbusWorkbenchPage` 表單屬性：
  - Server Port input：`name={designSystem.forms.name.modbusPort}`、`autoComplete={designSystem.forms.autocomplete.port}`、`inputMode={designSystem.forms.inputmode.numeric}`
  - Register input：`name={designSystem.forms.name.modbusRegister}`、`autoComplete={designSystem.forms.autocomplete.number}`、`inputMode={designSystem.forms.inputmode.numeric}`
  - 測試數值 input：`name="test-value"`、`autoComplete={designSystem.forms.autocomplete.number}`、`inputMode={designSystem.forms.inputmode.decimal}`
- 統一 `LocalModbusWorkbenchPage` loading 狀態 microcopy：
  - 狀態訊息：`{isBusy ? designSystem.microcopy.loading.default : message}`
  - 新增/更新映射按鈕：`{isBusy ? designSystem.microcopy.loading.saving : '新增 / 更新映射'}`
  - 執行測試寫入按鈕：`{isBusy ? designSystem.microcopy.loading.saving : '執行測試寫入'}`
- 統一 `useSmartDashboardModbusActions` feedback 訊息：
  - 同步成功：`${designSystem.microcopy.feedback.success.validated}：updated=...`
  - 同步失敗：`designSystem.microcopy.feedback.error.validationFailed`
- 統一 `LocalModbusWorkbenchPage` feedback 訊息：
  - 同步成功：`${designSystem.microcopy.feedback.success.validated}：updated=...`
  - 同步失敗：`designSystem.microcopy.feedback.error.validationFailed`
  - 匯入成功：`${designSystem.microcopy.feedback.success.created}：${imported.length} 筆`
  - 匯入失敗：`designSystem.microcopy.feedback.error.createFailed`
- 修正 ESLint 錯誤：
  - 移除 `smart-dashboard-regression.test.tsx` 中未使用的 `screen` import
  - 補齊 `useSmartDashboardModbusActions.ts` 中缺少的 `designSystem` import
- 驗證所有變更通過 ESLint 與測試（4 tests passed）

### Session Start
- 執行 `planning-with-files` session catchup。
- 確認專案根目錄存在，可建立持久化工作檔。

### Completed
- 建立 `task_plan.md`
- 建立 `findings.md`
- 建立 `progress.md`
- 更新 `AGENTS.md`，納入 UI/UX 主流程、TDD 與文件化工作流
- 更新 `AGENTS.md`，加入 `.github/instructions/` 規範來源分層
- 更新 `CLAUDE.md`，加入文件分工、閱讀順序與規範優先順序
- 更新 `GEMINI.md`，對齊目前專案脈絡與規範分工
- 盤查前端既有測試覆蓋，確認 `SmartDashboard` 與 `LocalModbusWorkbenchPage` 已有基礎測試
- 確認 `TestPage` 目前缺少對應測試檔，列為下一個 TDD 缺口
- 新增 `frontend/src/pages/__tests__/TestPage.test.tsx`
- 完成 `TestPage` 第一批頁面基線測試，覆蓋初始化最小化、連線後 config 收折、協議切換模式同步
- 驗證 `TestPage` 測試檔可通過 Vitest
- 驗證新增測試檔無 IDE lint 問題
- 完成 `TestPage` 第一輪 UI 外殼收斂，包含頁面背景、頂部定位說明、卡片表面與縮小配置卡樣式
- 驗證 `TestPage` 樣式調整後，基線測試仍通過
- 驗證 `TestPage.tsx` 與 `TestPage.test.tsx` 無 linter 錯誤
- 盤點前端測試分布，確認目前單元/互動測試散落於 `frontend/src/**/__tests__/*` 與 `frontend/src/**/*.test.ts(x)`，E2E 已位於 `frontend/tests/e2e/`
- 建立 `frontend/tests/README.md`，定義 `unit`、`integration`、`e2e` 三層分類
- 建立 root test wrappers，將 Vitest 正式入口統一收斂到 `frontend/tests/unit/` 與 `frontend/tests/integration/`
- 將 `TestPage` 測試實際移動到 `frontend/tests/unit/pages/test-page.test.tsx`
- 更新 `frontend/vite.config.ts`，限制 Vitest 只收斂 root 測試入口
- 更新 `frontend/package.json` 的 `test:gateway:unit`，改跑 root `tests/` 分類入口
- 更新 `frontend/tsconfig.json`，納入 root Vitest 測試入口
- 更新 `AGENTS.md`，加入前端測試正式入口規則與 Go `*_test.go` 例外
- 修正 `GatewayQuickSetupPage` 既有測試不穩定斷言，避免 root wrapper 驗證時誤判
- 修正 `gatewayAdapter.ts` 未使用例外變數與 `GatewayQuickSetupPage.tsx` 的 hook 依賴 lint 問題

### Completed（本 session）
- 規劃 SmartDashboard / LocalModbusWorkbench / TestPage 的 TDD 基線（寫入 findings.md）
- 盤點 `.github/instructions/` 在 UI/UX 實作中的具體套用點（寫入 findings.md）
- 將 Phase 3 標記為 complete

### In Progress
- Phase 1：Legacy 盤查與 TDD 基線建立（進行中）

### Completed（本 session - Phase 1）
- 盤查 SmartDashboard 現役結構與未使用檔案
- 識別三個可安全刪除的未使用檔案：`DatalinkLayout.tsx`、`SmartDashboardSidebarTools.tsx`、`SmartDashboardTagAndModbusPanel.tsx`
- 盤查 Legacy 路由與 redirect 策略
- 驗證 SmartDashboard 主流程測試覆蓋（25 tests passed）
- 記錄 TDD 基線狀態與測試缺口分析

### Completed（本 session - Phase 1, Phase 2 & Phase 3 開始）

**Phase 3: 重整 SmartDashboard 核心流程（完成）**
- 在 `SmartDashboardCommitPanel` 加入資料流向說明區塊，明確說明資料會流向資料庫與本地 Modbus
- 在 `SmartDashboardSidebar` 加入流程指引（規劃 → Tag → Modbus → 提交），顯示當前步驟
- 簡化 `SmartDashboardWorkflowModal`，收斂為設備入口：
  - 非設備 intent（points、mappings、tags、polling-groups、settings、wizard）顯示簡化訊息
  - 引導使用者前往設備管理中心或使用主工作流程
  - 保留設備管理完整功能（搜尋、篩選、建立、編輯、測試、啟用/停用、刪除）
- 將 `SmartDashboardWorkspaceContent` 聚焦在「資料來源設定 + 格子可視化」：
  - Flow Status Section 改為可收折，預設收合以聚焦主流程
  - 保留 Source Planner 緊湊列與 Memory Grid 作為核心功能
- 驗證所有測試通過（25 tests passed，4 test files）
- 驗證 lint 檢查通過

**Phase 4: 整理 LocalModbusWorkbench 與 TestPage（完成）**
- 重整 `LocalModbusWorkbenchPage` 區塊層級：
  - 系統狀態區塊：Server 狀態、映射數量、衝突狀態，含 Server 控制按鈕
  - Mapping 編輯區塊：Tag 選擇、Register 輸入、映射列表
  - 衝突治理區塊：顯示衝突列表與解決指引
  - 寫入測試區塊：Tag 選擇、測試數值輸入、執行測試寫入
- 改善標題與說明文字，使其更符合單人工作流程
- 更新測試以匹配新 UI（4 tests passed）

### 修改檔案清單（Phase 3 & Phase 4）
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardCommitPanel.tsx`：加入資料流向說明區塊
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardSidebar.tsx`：加入流程指引（規劃 → Tag → Modbus → 提交）
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardWorkflowModal.tsx`：收斂為設備入口，非設備 intent 顯示簡化訊息
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardWorkspaceContent.tsx`：Flow Status 改為可收折，預設收合
- `frontend/src/pages/datalink/LocalModbusWorkbenchPage.tsx`：重整區塊層級（系統狀態 / Mapping 編輯 / 衝突治理 / 寫入測試）
- `frontend/src/pages/datalink/__tests__/LocalModbusWorkbenchPage.test.tsx`：更新測試以匹配新 UI
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardSidebar.tsx`：補強 accessibility（role、aria-label、aria-current）
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardCommitPanel.tsx`：補強 accessibility（aria-live、role="status"）
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardWorkspaceContent.tsx`：改善響應式設計（px-2 sm:px-4）
- `frontend/src/pages/datalink/LocalModbusWorkbenchPage.tsx`：補強 accessibility（aria-live、role="status"）
- `frontend/tests/integration/ui/smart-dashboard-regression.test.tsx`：建立 UI regression test

### Completed（本 session - Phase 1 & Phase 2）
**Phase 1: Legacy 盤查與 TDD 基線建立**
- 盤查 SmartDashboard 現役結構與未使用檔案
- 識別三個可安全刪除的未使用檔案：`DatalinkLayout.tsx`、`SmartDashboardSidebarTools.tsx`、`SmartDashboardTagAndModbusPanel.tsx`
- 盤查 Legacy 路由與 redirect 策略
- 驗證 SmartDashboard 主流程測試覆蓋（25 tests passed）
- 記錄 TDD 基線狀態與測試缺口分析

**Phase 2: 建立最小 UI 規範底座**
- 盤查現有設計 tokens 與樣式系統
- 識別設計系統問題（多套 tokens、元件未統一、表單規範缺失、microcopy 未規範化）
- 建立 `frontend/src/styles/designSystem.ts` 統一設計系統規範
  - 整合 tokens 作為單一來源
  - 提供元件樣式類別（button、card、badge、sectionHeader）
  - 定義表單規範（input、label、error、autocomplete、inputmode、name）
  - 定義 microcopy 規範（loading、button、feedback、ellipsis）
- 驗證 `designSystem.ts` 無 lint 錯誤

### Validation
- 已執行：`cd frontend && npm run test -- src/pages/__tests__/TestPage.test.tsx --run`
- 結果：通過（3 tests，於新增測試後與樣式調整後各驗證一次）
- 已檢查：`TestPage.test.tsx` 無 linter 錯誤
- 已檢查：`TestPage.tsx` 無 linter 錯誤
- 已執行：`cd frontend && npm run test -- tests/unit/pages/test-page.test.tsx --run`
- 結果：通過（3 tests）
- 已執行：`cd frontend && npm run test:gateway:unit`
- 結果：通過（32 tests）
- 已執行：`cd frontend && npm run test -- tests/unit tests/integration --run`
- 結果：通過（14 files, 207 tests）
- 已執行：`cd frontend && npm run lint`
- 結果：通過
- 已執行：`cd frontend && npm run build`
- 結果：通過

### Notes
- 後續若開始 UI/UX 實作，需先補主流程測試基線
- `.github/instructions/` 已被提升為專案正式實作規範來源之一
- `TestPage` 現在已有最小可用的 TDD 保護網，可安全進入樣式一致化階段
- 下一個自然階段可選：`LocalModbusWorkbenchPage` 的視覺收斂，或 `SmartDashboard` 主流程的 TDD 缺口補強
- 前端測試現已正式收斂到 root `frontend/tests/` 作為執行入口
- 第二階段實體遷移：utils、hooks、features、components 已搬至 `frontend/tests/`，來源檔已刪除；頁面測試（Gateway、SmartDashboard）仍以 wrapper 匯入 `src/pages/.../__tests__/`，可於後續 session 遷移
