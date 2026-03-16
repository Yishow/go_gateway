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
- 本輪驗證：
  - `cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-readiness.test.ts tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx`（59 tests passed）
  - `cd frontend && npm run lint`
  - `cd frontend && npx tsc --noEmit`
  - `cd frontend && npm run build`
- Reviewer sub-agent 多次因 429 未能返回有效 review，本輪改以 TDD + targeted validation + controller manual spec spot-check 收斂 Step 1。
- 下一步：dispatch Phase 1 下一波 ready todos（source-rule-model / tag-board / output-board）
