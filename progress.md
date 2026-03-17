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

### Phase 2 後續實作 - 建立統一表單元件 [完成]
- 建立 `frontend/src/components/forms/` 目錄與統一表單元件：
  - `FormLabel.tsx`：統一的表單標籤元件，支援必填標記（`required` prop）
  - `FormInput.tsx`：整合的表單輸入元件，包含 Label、Input、Error、Hint、Success 訊息
  - `FormError.tsx`：統一的錯誤訊息元件，支援圖示（AlertCircle）
  - `FormHint.tsx`：統一的提示訊息元件
  - `FormSuccess.tsx`：統一的成功訊息元件，支援圖示（CheckCircle2）
  - `index.ts`：統一匯出所有表單元件
- 所有元件特性：
  - 使用 `designSystem.forms.*` 樣式系統
  - 完整的 TypeScript 型別定義與 JSDoc 註解
  - 支援 accessibility（`aria-invalid`、`aria-describedby`、`htmlFor`）
  - 支援自訂 className 與 id
- 驗證所有元件通過 ESLint

### Phase 2 後續實作 - 元件遷移到 designSystem.components [進行中]
- 遷移 `SmartDashboardWorkspaceContent` icon 按鈕：
  - Grid 切換按鈕（ChevronLeft/Right）使用 `designSystem.components.button.icon`
- 遷移 `SmartDashboardSidebar` icon 按鈕：
  - Workbench、Import、Export、Undo、Redo、Shortcuts 按鈕使用 `designSystem.components.button.icon`
- 遷移 `SmartDashboardCommitPanel` 卡片與 badge：
  - Segment feedback 卡片使用 `designSystem.components.card.panel`
  - Commit action message 卡片使用 `designSystem.components.card.panel`
  - Chunk results 卡片使用 `designSystem.components.card.panel`
  - Badge（pending、linked、conflict、committed、failed）使用 `designSystem.components.badge.*`
- 驗證所有變更通過 ESLint

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

## Session: 2026-03-15

### Datalink UI 深度分析與重整方向確認
- **Status:** complete
- Actions taken:
  - 建立 SQL todos，拆分為前端主流程盤查、後端能力盤查、可重用資產/技術債盤查、整合改版方案四個任務
  - 以 `claude-opus-4.6` 平行啟動三個分析代理，並確認 todos 狀態全數完成
  - 快速 spot-check `SmartDashboardPage.tsx`、`LocalModbusWorkbenchPage.tsx`、`legacyRoutes.ts`
  - 使用 `ask_user` 確認使用者偏好：
    - 同一畫面先完成 `Tag -> Local Modbus`，資料庫作為同流程下一步
    - 可視化偏好為格狀視覺化為主、表格為輔
  - 產出整合結論：建議採 **混合式過渡**，新增 `/datalink/workbench` 逐步取代現有 SmartDashboard 主線
  - 完成設計 spec：`docs/superpowers/specs/2026-03-15-datalink-workbench-design.md`
  - 完成 spec review 並取得使用者核准進入 implementation plan
  - 以繁中 commit 設計 spec：`5afef10 新增 datalink workbench 設計規格`
  - 建立 execution backlog：foundation / shell UI / source canvas / tag binding / local modbus / integration / quality / phase2 contracts
  - 依使用者指示先做「2、3 再做 1」，補齊兩份 detail spec：
    - `docs/superpowers/specs/2026-03-15-datalink-workbench-ui-detail.md`
    - `docs/superpowers/specs/phase2-runtime-dbtarget-detail.md`
  - 補查 implementation 風險：確認 `MemoryGrid` props 已足以承接 workbench 主視覺；確認 `POST /datalink/points/:id/poll` 與 `POST /datalink/points/poll` 皆已掛上 router
- Files created/modified:
  - `findings.md`
  - `task_plan.md`
  - `progress.md`
  - `/Users/yishow/.copilot/session-state/63ec8c96-8e62-43c4-969c-5c845467d5be/plan.md`

## Test Results
| Test | Input | Expected | Actual | Status |
|------|-------|----------|--------|--------|
| SQL todo workflow | 建立/更新分析 todos | todo 狀態依序變更 | 四個 todos 皆完成 | ✓ |
| 平行分析代理 | `frontend-flow-audit` / `backend-capability-audit` / `reuse-debt-audit` | 各自完成盤查並回報 | 三個代理皆完成且 SQL 狀態為 `done` | ✓ |
| Spec review loop | reviewer 檢查 spec 完整性 | spec 可進 implementation planning | reviewer 核准通過 | ✓ |
| Spec commit | `git commit` spec 文件 | 設計 spec 獨立提交 | commit `5afef10` 完成 | ✓ |
| Detail planning | UI 細節與 Phase 2 契約補齊 | 形成可開工前的完整規劃 | 兩份 detail spec 已完成 | ✓ |

## 2026-03-16：Source templates / Tag batch 收尾

### Completed
- `redesign-source-templates`
  - 在 `SourceCanvasSection.tsx` 接上本地 template save/load UI，不再停留在 disabled buttons
  - `sourceTemplateStorage.ts` schema 升到 v3，新增 `preferredViewMode` 與 `capabilitySnapshot`
  - `sourcePlannerContract.ts` 補齊 template contract / create helper 對新欄位的支援，同時保留 SmartDashboard 舊資料的向下相容
  - Source step 測試新增儲存模板、重新套用 planner inputs / view mode 的覆蓋
- `redesign-tag-batch`
  - 整合外部變更後，完成 `TagBindingStudio.tsx` 的 batch diff preview 與結果總表
  - `tagBindingModel.ts` 新增 `buildBatchDiffPreview()` 與相關 diff contracts
  - Tag step / model tests 補齊 create 與 existing flow 的 diff preview、skip reasons、結果 summary 覆蓋

### Validation
- `cd frontend && npm run test -- --run tests/unit/features/datalink/sourcePlannerContract.test.ts tests/unit/features/datalink/sourceTemplateStorage.test.ts tests/unit/features/datalink/tag-binding-model.test.ts tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`
- `cd frontend && npm run build`

## 2026-03-16：Workbench desktop polish（scan-first）spec 與 implementation planning

### Completed
- 重新做 workbench 桌面 audit，並用 brainstorming 與使用者逐段確認 shell 規則、step 改法、互動規則與測試策略
- 新增 spec：`docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-polish-design.md`
- 完成 spec review loop
  - reviewer 首輪指出 3 個問題：
    - 未明確說明與前一份 redesign spec 的 supersede 邊界
    - `handoff summary` 有 scope creep
    - 驗收與測試條件太主觀
  - 已修正後取得 `APPROVED`
- 完成兩個繁中 commit：
  - `ce2eae0` 補齊 workbench 掃描優先 polish 設計規格
  - `eb4a2c2` 修正 workbench polish 規格邊界與驗收條件
- 取得使用者同意進入 implementation planning
- 已將 planning 檔同步到新的 `scan-first` phase breakdown

### Completed
- `scan-first` implementation planning
  - `scan-shell-density`
  - `scan-device-browser`
  - `scan-source-toolbar-canvas`
  - `scan-tag-row-board`
  - `scan-output-active-target`
  - `scan-cross-step-regressions`

### Notes
- 這一輪尚未開始新的前端程式碼實作，先完成 spec 與 implementation plan，避免又在 UI 方向上走偏。
- 使用者新增偏好：前端設計 / 提案若能選模型品質，優先採用 Opus 級別輸出。

## 2026-03-16：scan-shell-density

### Completed
- 以 TDD 先補 `WorkbenchContextBar` 的紅燈：
  - compact step summary 取代 capability chips 常駐顯示
  - single primary action 取代 quick-action cluster
- reviewer 補抓 `source -> output` 跳步 regression 後，再補一輪 TDD：
  - `source` step CTA 改成 `gotoTag`
  - `sourceReady` 未滿足時 disabled
  - locale 補上 `workbench.contextBar.actions.gotoTag`
- reviewer 第二輪再抓到 Step 4 CTA 語義不對：
  - output step 不再退回 `switchDevice`
  - 改為 output-focused CTA 文案（`configureOutput` / `completeTagBinding`）
  - output-ready 時，CTA 會 focus `output-primary-anchor`，不再是 enabled no-op
- 以 TDD 補 `WorkbenchBottomSummaryBar` 的紅燈：
  - 只高亮 active step
  - 其餘 step 改成 compact readiness marker
- 將 related foundation tests 從舊的 context-bar assumptions 對齊到 scan-first spec

### Files modified
- `frontend/src/pages/datalink/workbench/WorkbenchContextBar.tsx`
- `frontend/src/pages/datalink/workbench/WorkbenchBottomSummaryBar.tsx`
- `frontend/src/pages/datalink/workbench/LocalModbusBoard.tsx`
- `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx`
- `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchFoundation.test.tsx`

## 2026-03-17：scan-device-browser

### Completed
- 先在 `DatalinkWorkbenchFoundation.test.tsx` 補紅燈：
  - Step 1 要有 `device-primary-toolbar`
  - compact row 要只保留 endpoint / health / 兩個 capability hints
  - MQTT / MC3E compact row endpoint 要顯示真正 connection target，而不是 `protocol-traits`
- `WorkbenchDeviceStep.tsx` 已改為 scan-first Step 1：
  - toolbar 收成單列（search / protocol / status / refresh / create）
  - main surface 移除 clone / continue cluster
  - device list 改成 compact rows，詳細能力與 recent tests 留在 inspector
- `workbenchDeviceFormModel.ts` 新增 `buildDeviceEndpointSummary()`，將 host/port、serial、broker URL 等 endpoint derivation 集中處理。
- foundation 內「前進到 source」測試改為走 ContextBar CTA，對齊 shell-first 流程。

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`

## 2026-03-17：scan-source-toolbar-canvas

### Completed
- 先在 `DatalinkWorkbenchSourceStep.test.tsx` 補紅燈：
  - Step 2 需拆成 `source-primary-toolbar` 與 `source-secondary-controls`
  - `source-canvas-workspace` 必須是 primary，`source-rule-layer` 必須是 supporting
- `SourceCanvasSection.tsx` 已改為 scan-first Step 2：
  - primary toolbar 承接 view mode / value format / planner inputs / apply
  - secondary controls 承接 jump / freeze / snapshot / showAudit / batch create / save-load
  - coverage overview 留在 primary workspace 內，但改為 lighter strip
  - rule layer 改為 supporting side rail，減少與 canvas 競爭

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`

## 2026-03-17：scan-tag-row-board

### Completed
- 先在 `DatalinkWorkbenchTagStep.test.tsx` 補紅燈：
  - row board 需要明確的 `tag-candidate-board` / `data-layout="row-board"` 結構
  - batch diff preview 與 bind CTA 只有在至少選取一筆 row 後才顯示
  - conflict detail 不再出現在預設 row，而是改由 inspector 承接
- `TagBindingStudio.tsx` 已改為 scan-first Step 3：
  - `selectedPointIds` 進入 Step 3 時預設為空，不再自動全選所有 point
  - candidate list 改成 row-board / table-like 密度，將 source meta、raw/transformed value、preview/status 壓回單列
  - `alreadyLinked` / `existingKey` / `duplicatePreview` badge 從 row 主表面移除，交由 inspector 顯示
  - batch diff preview / bind CTA 僅在 selection 存在時顯示

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-tag-step.test.tsx`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`

## 2026-03-17：scan-output-active-target

### Completed
- 先在 `DatalinkWorkbenchOutputStep.test.tsx` 補紅燈：
  - shared candidate rows 只顯示 active target 的 mapping field
  - Modbus operational panels 要明確標成 supporting sections
  - Database 的 schema snapshot / write preview 要被 grouped 成 supporting secondary panels
  - shared output candidate selection 要與 Database mapping form 維持同步，不能 split-brain
- `LocalModbusBoard.tsx` 已改為 scan-first Step 4：
  - shared candidate row 依 `activeOutputTarget` 僅顯示 Modbus 或 Database 其中一種狀態 badge
  - Modbus metrics / server controls 區塊加上 `modbus-secondary-panels` supporting 語意
  - active target 切換後，shared candidate selection 與 database form 共用同一個 `selectedTagId`
- `DatabaseTargetBoard.tsx` 已改為 scan-first Step 4：
  - selection 改由 parent 注入，消除 shared board / database board 的雙重選取來源
  - schema snapshot 與 write-row-preview 改收進 `database-secondary-panels`

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-output-step.test.tsx`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`

## 2026-03-17：scan-cross-step-regressions

### Completed
- 針對 shell / foundation / source / tag / output / locale 跑一輪 workbench regression，確認 Step 3 與 Step 4 的 progressive disclosure 沒有打壞既有路徑。
- 驗證 shell CTA、shared output target switcher、source canvas hierarchy、tag row-board disclosure、database form sync 均維持可用。

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-shell-ui.test.tsx`
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx`
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`

## 2026-03-17：quiet desktop v2 收尾

### Completed
- Step 1 收斂完成：
  - `WorkbenchDeviceStep.tsx` 移除裝飾 hero block
  - device row 改為更緊湊的 desktop grid
  - capability hints 改成 inline pills
  - create/edit panel overlay 改成 viewport-level `fixed`
  - `WorkbenchInspectorPanel.tsx` 補上 scrollable aside
- Step 2 收斂完成：
  - `SourceCanvasSection.tsx` 將 planner 移入 rule layer
  - `apply` 正式改為 `addRule`
  - 加入 delete rule（workspace + inspector 兩條路徑）
  - utility tools 收進 `moreTools`
  - `AddressCanvas` 與 coverage strip 的主次層級重新排序
- Step 3 收斂完成：
  - `TagBindingStudio.tsx` 補齊 `flowModeHint` / `selectionHint`
  - raw/transformed value 扁平化為 row-board 欄位
  - existing mode 隱藏 preview key，只保留 existing tag select
  - sidebar metrics 壓縮，讓候選列掃描更穩
- Step 4 收斂完成：
  - `LocalModbusBoard.tsx` / `DatabaseTargetBoard.tsx` 改成 shared candidate board 為唯一 tag selection surface
  - Modbus / Database studio 以 read-only selected-tag summary 取代各自的 tag select
  - Database connector form 預設收合，並新增 `configureConnector` / `hideConnector`
- Code review 已執行，並依 reviewer 建議補強：
  - 點選既有 connector 時會直接展開 connector editor
  - inspector delete-rule 路徑與 source workspace delete-rule 路徑維持同樣的 focus/selection 清理條件

### Validation
- `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
- `cd frontend && npm run lint`
- `cd frontend && npx tsc --noEmit`
- `cd frontend && npm run build`
- code review：reviewer 先抓出 connector editor 展開互動缺口與 delete-rule contract 不一致，兩項已修正後再完成驗證

## 2026-03-17：master-detail polish spec 定稿

### Completed
- 已完成新一輪 workbench polish 設計定稿：
  - `docs/superpowers/specs/2026-03-17-datalink-workbench-master-detail-polish-design.md`
- 本輪設計重點：
  - Step 1 回到真正的 master-detail device studio
  - Step 2 明確定義 rule-first + selection-first 的雙 point 建立路徑
  - Step 2 右上健康摘要與底部 `BottomSummaryBar` 完成分工切線
  - `32-bit` / `64-bit` 明確定義為單一邏輯格，不允許半格選取
  - data type 分組對齊現有 `DataType` union，避免 spec drift
  - Step 3 改成至少有一筆 persisted/usable point 才解鎖
- 已依 frontend-design / UX psychology 原則將設計收斂到：
  - Hick's Law：減少同時選項
  - Fitts' Law：主 CTA 固定且可快速命中
  - Miller's Law：摘要與衝突資訊分塊呈現

### Validation
- spec review loop：
  - 第一輪抓出 3 個實質問題：data type drift、Step 2/BottomSummary 重複計數風險、Step 3 解鎖門檻模糊
  - 修正後重新審稿，結果：`Approved`
- 目前狀態：
  - spec 已寫成 repo 文件
  - 等待使用者 review spec 後，再決定是否進入 implementation planning

## 5-Question Reboot Check
| Question | Answer |
|----------|--------|
| Where am I? | implementation plan 已建立，等待是否啟動 Phase 1 |
| Where am I going? | 若使用者同意，即用 fleet mode 派發 foundation 與第一批 UI 子任務；Phase 2 契約也已有提前規劃 |
| What's the goal? | 讓 datalink UI 收斂成「來源 -> 可視化 -> Tag -> 輸出」的單純主線 |
| What have I learned? | 問題集中在頁面層，底層 domain 資產大多可沿用；runtime API 與 DB target 已完成細化規劃但尚未實作 |
| What have I done? | 已完成深度盤查、設計 spec、spec review、spec commit、implementation backlog、Phase 1/2 detail planning |

## Session: 2026-03-16

### Workbench 主線實作補齊
- **Status:** complete
- Actions taken:
  - 完成 Step 3 `TagBindingStudio`
    - 新增 `tagBindingModel.ts`
    - 支援 prefix / strategy preview key
    - 支援 existing key / duplicate preview / already linked 阻擋
    - 支援 partial failure summary
  - 完成 Step 4 `LocalModbusBoard`
    - 整合 Local Modbus server 啟停、register 綁定、conflict block sync、push current value
    - 以 selected device 的 points + mappings + tags 推導 output candidates
  - 強化 shell UI
    - 新增 `useWorkbenchSummary.ts`
    - `WorkbenchHeaderBar` 顯示 point/tag/output counts
    - `WorkbenchActionDock` 顯示 readiness 與 next action
  - 補齊 workbench i18n keys（en / zh-TW）
  - 補齊測試：
    - `workbench-tag-step.test.tsx`
    - `workbench-output-step.test.tsx`
    - `workbench-shell-ui.test.tsx`
    - `tag-binding-model.test.ts`
  - 完成 integration / quality 收尾
    - 手動 review `useWorkbenchSummary`、`WorkbenchHeaderBar`、`WorkbenchActionDock`、`TagBindingStudio`、`LocalModbusBoard`
    - 將 `workbench-integration`、`workbench-quality-pass` 回寫為 `done`
    - 確認 `runtime-live-value-phase2`、`database-target-phase2` 轉為 ready todo
  - 手動完成最終驗證
    - 分批執行 workbench 測試
    - 執行 `npm run lint`
    - 執行 `npm run build`
    - 執行 `git --no-pager diff --check`
- Files created/modified:
  - `frontend/src/pages/datalink/workbench/TagBindingStudio.tsx`
  - `frontend/src/pages/datalink/workbench/tagBindingModel.ts`
  - `frontend/src/pages/datalink/workbench/LocalModbusBoard.tsx`
  - `frontend/src/pages/datalink/workbench/useWorkbenchSummary.ts`
  - `frontend/src/pages/datalink/workbench/WorkbenchHeaderBar.tsx`
  - `frontend/src/pages/datalink/workbench/WorkbenchActionDock.tsx`
  - `frontend/src/pages/datalink/workbench/DatalinkWorkbenchPage.tsx`
  - `frontend/src/i18n/locales/en/common.json`
  - `frontend/src/i18n/locales/zh-TW/common.json`
  - `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchTagStep.test.tsx`
  - `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx`
  - `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchShellUi.test.tsx`
  - `frontend/src/pages/datalink/workbench/__tests__/tagBindingModel.test.ts`
  - `frontend/tests/unit/features/datalink/tag-binding-model.test.ts`
  - `frontend/tests/unit/pages/datalink/workbench-tag-step.test.tsx`
  - `frontend/tests/unit/pages/datalink/workbench-output-step.test.tsx`
  - `frontend/tests/unit/pages/datalink/workbench-shell-ui.test.tsx`

### Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
| `claude-opus-4.6` background agents returned 429 / rate limit | 1 | 改由主代理手動接手 Step 3 / Step 4 實作 |
| `gpt-5.4` general-purpose / explore background agents 長時間 running 但沒有第一輪 turn | 1 | 停止空等，直接由主代理手動做 quality pass 與 phase2 接手準備 |
| `LocalModbusBoard` 在測試中因 `loadData` 依賴 `t` 而重複觸發 effect | 1 | 以 `ref` 穩定 fallback translation，讓 callback 不受 `t` identity 影響 |
| `LocalModbusBoard` register input 被 effect 重設回 `0` | 1 | 將 effect 依賴收斂到 `selectedTagId` 與 `selectedCandidateRegister` |
| `TagBindingStudio` selection 因 fresh array dependency 反覆重設 | 1 | 改為 `pointIdsKey -> split` 產生穩定 id list |
| 多個 workbench Vitest 檔一次串跑時 worker 在測試通過後不正常結束 | 1 | 改採分批驗證，保留定位能力並避免白等 |

## Test Results
| Test | Input | Expected | Actual | Status |
|------|-------|----------|--------|--------|
| Tag binding model | `tests/unit/features/datalink/tag-binding-model.test.ts` | preview/build request contract 正常 | 3 tests passed | ✓ |
| Tag step UI | `tests/unit/pages/datalink/workbench-tag-step.test.tsx` | empty state / preview / conflict / partial failure 正常 | 4 tests passed | ✓ |
| Output + Local Modbus | `tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/pages/datalink/local-modbus-workbench.test.ts` | server start / register bind / conflict block / legacy page 正常 | 9 tests passed | ✓ |
| Shell UI summary | `tests/unit/pages/datalink/workbench-shell-ui.test.tsx` | header/action dock 反映 counts 與 next action | 2 tests passed | ✓ |
| Unit/meta batch | `tests/unit/utils/designSystemForms.test.ts ... tests/unit/features/datalink/legacyRoutes.test.ts` | foundation/meta contracts 正常 | 19 tests passed | ✓ |
| Shell/foundation batch | `tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx` | route + summary UI 正常 | 4 tests passed | ✓ |
| Manual quality pass | 三批 workbench 測試 + `npm run lint` + `npm run build` + `git --no-pager diff --check` | workbench 可 commit，且無新增格式問題 | pass（僅既有 Vite chunk size warning） | ✓ |
| Frontend lint | `cd frontend && npm run lint` | 無 lint error/warning | pass | ✓ |
| Frontend build | `cd frontend && npm run build` | TypeScript + Vite build 通過 | pass（僅既有 chunk size warning） | ✓ |

## Next Focus
- runtime/live value vertical slice 已完成大部分核心：
  - backend：`runtime/status`、`runtime/stream`、runtime service pubsub、point/mapping refresh path
  - app wiring：`cmd/test_ui/main.go` 已接上 scheduler/runtime，smoke 成功
  - frontend：source step runtime summary + EventSource live raw value
- `runtime-live-value-phase2` 尚餘 `runtime-poll-contract`
- 下一個主目標：`database-target-phase2`

## Session: 2026-03-16（續）

### Database Target Phase 2 收尾
- **Status:** complete
- Actions taken:
  - 在同一個 workbench output step 新增 `DatabaseTargetBoard`，讓 Database Target 與 Local Modbus 同時存在於同一條主線
  - 落地 dbtarget schema / migration / repo / service / writer / API handler / router / runtime writer fan-out
  - 補齊前端型別、API client、i18n 與 `workbench-output-step` 測試
  - 修正 reviewer 指出的契約問題：
    - redaction 後的密碼更新路徑新增 `clear_password`
    - upsert timestamp 欄位必須為 PK / single-column unique
    - insert / upsert 切換時清空並忽略 `timestamp_column`
    - sqlite introspection self-deadlock
  - 修正 `point_handler_extended_test.go` 的 large-list baseline 測試資料
- Validation:
  - `go test ./internal/datalink/dbtarget ./internal/datalink/runtime ./internal/api/handlers ./internal/api -count=1`
  - `go build ./cmd/test_ui`
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-output-step.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npm run build`
  - `git --no-pager diff --check`

### Workbench Device Step 補完（完成）
- 使用者實機回報：進入 `/datalink/workbench` 後沒有設備可設定與選取，頁面只剩骨架。
- root cause 已確認：
  - `WorkbenchProvider` 預設 active step 為 `device`
  - `DatalinkWorkbenchPage` 對 `device` step 仍走 placeholder render path，沒有真正的設備設定 UI
- 本輪採用的修正策略：
  - 按 spec 的 Slice 1 一次補齊 `HeaderBar/ContextBar + Device Step + Inspector + embedded create/edit panel + connection test`
  - 不重用舊 `DeviceForm` / `DeviceOnboardingWizard` UI，只參考其欄位模型與 hooks / API
- 本輪完成項目：
  - 新增 `WorkbenchDeviceStep.tsx`
    - Device cards/list
    - 搜尋 / protocol / status 篩選
    - Step 1 inspector
    - embedded create/edit panel
    - 連線測試互動與 inline status
  - 新增 `workbenchDeviceFormModel.ts`
    - 協定欄位預設值
    - connection config parse / sanitize
    - connection summary builder
  - `WorkbenchHeaderBar` 改為較接近 spec 的 context bar：
    - selected device / protocol / status / last test
    - 快速跳回 device / source / output
  - `DatalinkWorkbenchPage` 將 device step placeholder 換成真正的 `WorkbenchDeviceStep`
  - 補齊 zh-TW / en 的 workbench device i18n
- 本輪手動 review 額外修掉：
  - 編輯設備時若把 description 清空，原本會送 `undefined` 導致後端保留舊值；已改成送空字串，並補 regression test
- 本輪驗證：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/pages/datalink/workbench-runtime-phase.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npm run build`
- 雜項：
  - 清除先前驗證殘留的 `.tmp_dbtarget_verify.go` 與 `test_ui`

### Workbench browser smoke + polish（完成）
- 使用者回報 Vite 開發環境出現 `/api/v1/datalink/*` proxy errors。
- 實際查證結果：
  - Vite dev server 仍在 `localhost:5173`
  - proxy target 預設為 `http://localhost:8080`
  - 當時真正的 root cause 是後端 `test_ui` 沒有在 `8080` listen，不是剛重做的 Device Step UI 壞掉
- smoke 流程驗證：
  - 啟動 `go run ./cmd/test_ui` 後，`/api/v1/datalink/devices` direct / proxy 都恢復 `200`
  - 使用 `agent-browser` 驗證 `Create device -> Select device -> Continue to source`
  - 驗證 `Test connection` 失敗時會顯示明確錯誤，而不是靜默失敗
- 依 smoke 補做 polish：
  - `WorkbenchActionDock` 在 Step 1 已選設備時，改為顯示前往 Source 的 next action
  - Step 1 inspector 對 `0001-01-01T00:00:00Z` 改顯示 fallback microcopy，不直接暴露 zero timestamp
- smoke 後驗證：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/pages/datalink/workbench-runtime-phase.test.tsx`
  - `cd frontend && npm run lint`
  - `cd frontend && npm run build`

## 2026-03-16：Workbench desktop redesign round 2（brainstorming -> spec drafting -> execution plan）
- 已依使用者要求停止把 layout 問題當成小修，改以第二輪 desktop redesign 處理。
- 已確認 root cause：workbench shell 偏離原 spec 骨架，導致 1920×1080 下主區過窄，Step 1 內嵌 inspector 壓縮主欄，Step 2 缺 detail surface。
- 已與使用者逐段確認：
  - 回到原始 desktop 骨架
  - Step 1 全新 `DeviceWorkspace`
  - Step 2 全新 `AddressCanvasWorkspace`
  - Step 3 全新 `TagBindingBoard`
  - Step 4 全新 `OutputWorkspace`
- 已補入使用者挑選的實用設計：
  - Step 1：能力摘要、clone device
  - Step 2：rule layer、coverage/gap overview、value format、freeze/snapshot
  - Step 3：命名規則預覽、source/value/merge 詳情、existing/new tag 分流、diff preview、結果總表
  - Step 4：register map、base/offset、schema snapshot、required/missing、output readiness、filter/search、auto-map、dry-run、health summary、preview、sync result
- 已建立新 spec：`docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-redesign.md`
- 已 scaffold OpenSpec change：`openspec/changes/redesign-datalink-workbench-desktop-flow/`
- 已完成 OpenSpec proposal/design/specs/tasks artifacts

### Execution plan 完成（本 session）
- 建立 20 個 SQL todos（`redesign-*`）與 31 條 dependency edges
- Phase 0（Shell）→ Phase 1（四條平行 track）→ Phase 2（Quality + Rollout）
- 最大平行度：Phase 0 完成後可同時啟動 4 個 step workspace track
- 每個 track 內部有序列依賴（如 rule-model → canvas → viewmodes）
- 更新 `task_plan.md` 加入完整 phase breakdown 與 parallelism map
- 更新 `findings.md` 加入 redesign follow-up 記錄
- 更新 session `plan.md` 加入 execution plan 指引
- 已完成 Phase 0 shell：
  - `WorkbenchFrame` / `WorkbenchStepRail` / `WorkbenchContextBar` / `WorkbenchInspectorPanel` / `WorkbenchBottomSummaryBar`
  - `WorkbenchProvider` 新增 `inspectorSelection`、`activeOutputTarget`、`crossStepContext`
  - `useWorkbenchSummary` 新增 device/source/tag/output readiness 匯總
- Phase 0 驗證：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-readiness.test.ts tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx`（54 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
- 已完成 Phase 1 / Step 1 slice：
  - `WorkbenchDeviceStep` 改為新 DeviceBrowser 主區，不再內嵌第二個 desktop inspector
  - `WorkbenchInspectorPanel` 在 Step 1 會顯示 device identity / capability summary / connection summary / recent test timeline / clone action
  - `WorkbenchContextBar` 顯示 selected device capability chips
  - clone flow 以 workbench drawer 形式預填 connection defaults，但保留空白 name 強制新識別
- 已完成 Step 1 review follow-up：
  - create/clone 後以 pending selection guard 避免 list refresh race 清掉新選取設備
  - clone flow 明確阻擋與來源設備同名
  - context bar latest test tone 改以前端 session-local latest test result 為準
  - inspector `Test connection` pending 時 disable，避免重複點擊
- 已完成 Phase 1 / Step 2 第一輪核心 slice：
  - `SourceCanvasSection` 改為真正的 AddressCanvasWorkspace：RuleLayerBar + PlannerToolbar + CoverageOverview + audit drawer
  - `sourceCanvasModel` 改為連續 16-bit lattice，會把 `planned / used / conflict / gap` 放進同一條連續位址帶
  - 已支援 `Plan / Live / Link` 三種資訊層、value format、freeze live / snapshot compare、jump to address
  - Link view 先用 point→mapping→tag 狀態導出 `needsPoint / unbound / draft / ready / blocked`
- 已完成 Phase 1 / Step 3 主板第一輪 slice：
  - `TagBindingStudio` 改為 dense board，單列直接顯示 source address、span、raw/transformed value、preview key、status badge
  - 已支援 `create new` / `bind existing` 明確切換、search、status filter、existing tag select
  - `tagBindingModel` 已補齊 `bitWidth / cellSpan / rawValue / bindingStatus / existingTagOptions`
- 本輪驗證：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-readiness.test.ts tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx`（59 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- Step 2 補充驗證：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-source-canvas-model.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`（12 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- Step 3 補充驗證：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/tag-binding-model.test.ts tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`（14 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- Reviewer sub-agent 多次因 429 未能返回有效 review，本輪改以 TDD + targeted validation + controller manual spec spot-check 收斂 Step 1。
- `claude-opus-4.6` 的 `redesign-source-rule-model` / `redesign-tag-board` 子代理都因 429 失敗，已改為 controller 手動接手。
- 已由 controller 手動接手並完成 Step 4 `OutputWorkspace` shared board：
  - `LocalModbusBoard` 現在先顯示 shared candidate board，再以 target switcher 切換 Local Modbus / Database studio
  - 同一列會同時顯示 Modbus badge 與 database path/status，避免切頁才知道另一個 target 是否已綁定
  - 為避免 shared board 測試留下 React `act(...)` 警告，test 會在切換到 database 後等待 mapping table 掛載完成
- 已補修 Step 2 source rule persistence：
  - `SourceCanvasSection` 的 `rules / selectedRuleId / selectedAddress` 已提升到 `WorkbenchProvider`
  - `Source → Tag → Source` 往返不再遺失已套用規則
  - rule id 會依目前 persisted rules 推導，避免 component remount 後再度產生重複 `rule-1`
- Step 4 補充驗證：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`（11 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- broader workbench regression：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`（69 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- 已完成 Step 4 細化收尾：
  - `LocalModbusBoard` 已把既有 `RegisterMapCanvas` / auto-map / dry-run helper 接回主 UI
  - `DatabaseTargetBoard` 已補 `schema-snapshot`、required highlight、`write-row-preview`
  - `WorkbenchInspectorPanel` 已補 `outputCandidate` trace panel，顯示 source→tag→Local Modbus / Database readiness
  - controller 手動 review 後，再補上 multi-word register overlap 判斷與 inspector mapping 載入失敗提示
- Step 4 收尾驗證：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`（25 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- 已完成 shared inspector 第二波收尾：
  - Step 2 rule / span selection 會在右側 inspector 顯示 coverage、origin、bit width、raw value、downstream link state
  - Step 3 candidate row 已接上 `tag` inspector selection，右側可顯示單筆 tag key、source address、raw/transformed value、span/bit width 與 conflict/already-linked 提示
- 第二波 inspector 驗證：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/sourcePlannerContract.test.ts tests/unit/features/datalink/sourceTemplateStorage.test.ts tests/unit/features/datalink/tag-binding-model.test.ts tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx`（80 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- 下一步：進入 `redesign-i18n-a11y` / `redesign-regression-tests` / `redesign-legacy-compat`。

## Session: 2026-03-16（fleet follow-up）

### Phase 2 / spec alignment audit（進行中）
- 重新查 SQL 後確認尚未完成的 todo 只剩：
  - `redesign-i18n-a11y`（in_progress）
  - `redesign-regression-tests`（in_progress）
  - `spec-alignment-audit`（in_progress）
  - `redesign-legacy-compat`（pending）
- `redesign-legacy-compat` 目前在 SQL 依賴 `redesign-regression-tests`，因此尚未進入 ready queue。
- 手動補查 OpenSpec `openspec/changes/redesign-datalink-workbench-desktop-flow/tasks.md`，發現未勾項目為：
  - `3.2 Rebuild AddressCanvas ... continuous 16-bit lattice`
  - `3.3 Plan / Live / Link ...`
  - `6.1 i18n / accessibility / keyboard`
  - `6.2 desktop regression`
  - `6.3 compatibility strategy`
- 進一步 spot-check `SourceCanvasSection.tsx` / `AddressCanvas.tsx`：
  - `3.3` 所需的大部分控制列已在（view mode、value format、freeze/snapshot、jump-to-address、coverage overview）
  - 但 `AddressCanvas.tsx` 仍是卡片式 responsive grid，而非批准 spec 的固定 16-bit lattice，表示 `3.2` 仍有 drift
- 進一步 spot-check rollout 路由：
  - `App.tsx` 仍將 `/datalink/local-modbus` 直接導向 `LocalModbusWorkbenchPage`
  - `SmartDashboardPage.tsx` 仍把後段 CTA 導向舊 `/datalink/local-modbus`
  - 代表 `6.3` 的新舊入口導流策略尚未落地
- 背景代理狀態：
  - `agent-72` / `agent-73` / `agent-74` 仍在 running 且尚未產生第一輪 turn
  - `agent-77` spec audit 直接以 `429` 結束，未提供可用摘要

### Errors Encountered
- `agent-77`（spec alignment explore）再次遭遇 `429 rate limit`，代表本輪 fleet 仍不能只依賴 Opus 子代理產出。
- 重新跑 `planning-with-files` catchup 時，誤用了 `~/.codex/.../session-catchup.py`；本機 skill 實際位於 `~/.copilot/skills/...`。

## Session: 2026-03-17（Phase 2 收尾完成）

- 已手動完成 `redesign-address-canvas-revisit`：
  - `AddressCanvas.tsx` 改為固定 16 欄的連續 lattice row render
  - row 補上 `data-row-start-address` / `data-row-end-address` / `data-lattice-columns`
  - cell 補上 `data-merge-span` / `data-merge-offset` 與 `aria-pressed`
- 已完成 `redesign-i18n-a11y`：
  - `WorkbenchDeviceStep` notice banner 補 `aria-live="polite"`
  - `WorkbenchInspectorPanel` output inspector load error 補 `aria-live="polite"`
  - `WorkbenchDeviceStep` 的 parity / data format option 改為 `t()`，移除 `None / Even / Odd / Binary / ASCII` 硬編碼 label
  - `DatalinkWorkbenchFoundation` / `WorkbenchLocaleContract` 新增對應回歸測試
- 已完成 `redesign-regression-tests`：
  - 補 `WorkbenchFrame` 桌面殼層 sizing + `overflow-hidden` regression test
  - 補 Step 2 多 rule / merged span / gap cell regression test
  - runtime tests 改成顯式切到 `live` overlay，避免依賴過時預設值
- 已完成 `redesign-legacy-compat`：
  - `App.tsx` 將 `/datalink/local-modbus` 改為 compat redirect，導向 `/datalink/workbench?step=output&target=modbus`
  - 保留舊頁 fallback 路徑 `/datalink/local-modbus/legacy`
  - `SmartDashboardPage.tsx` 後段 CTA 直接導向新 workbench，而不再先進舊 Local Modbus 頁
  - `DatalinkWorkbenchPage` 新增 query bootstrap，支援 `step` / `target` deep link，且在無 Router 的單元測試環境下安全退化
- 本輪驗證已通過：
  - `cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/pages/datalink/workbench-runtime-phase.test.tsx tests/unit/features/datalink/workbench-locale.test.ts tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-readiness.test.ts tests/unit/features/datalink/workbench-source-canvas-model.test.ts`
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
