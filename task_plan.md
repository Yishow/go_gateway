# UI/UX 收斂任務計畫

## 任務目標
- 將專案規範來源整理為一致的層級結構。
- 明確要求 Agent 同時遵守 `AGENTS.md`、agent 專屬文件與 `.github/instructions/`。
- 以單人操作情境重整前端 UI/UX 改造方向。
- 所有後續 UI 變更採 TDD 先行，並保留完整工作記錄。

## 規範來源層級
- `AGENTS.md`
  專案共通目標、目錄分工、測試、TDD、文件化工作流。
- `CLAUDE.md` / `GEMINI.md`
  Agent 工作流程、架構脈絡、執行邊界與補充說明。
- `.github/instructions/*.md`
  依檔案類型套用的語言與框架實作規範。

## 階段
- [complete] Phase 1: 更新專案規範文件
- [complete] Phase 2: 建立工作記錄檔
- [complete] Phase 3: 後續 UI/UX 實作前置準備
- [complete] Phase 1 (UI/UX Plan): Legacy 盤查與 TDD 基線建立
- [complete] Phase 2 (UI/UX Plan): 建立最小 UI 規範底座
- [complete] Phase 3 (UI/UX Plan): 重整 SmartDashboard 核心流程
- [complete] Phase 4 (UI/UX Plan): 整理 LocalModbusWorkbench 與 TestPage
- [complete] Phase 5 (UI/UX Plan): 可近用性、響應式與品質驗證

## 待辦
- [x] 在 `AGENTS.md` 宣告 `.github/instructions/` 為正式規範來源
- [x] 在 `CLAUDE.md` 補充文件分工與閱讀順序
- [x] 在 `GEMINI.md` 對齊相同規範來源層級
- [x] 建立 `findings.md`
- [x] 建立 `progress.md`
- [x] 規劃 SmartDashboard / LocalModbusWorkbench / TestPage 的 TDD 基線
- [x] 補上 `TestPage` 第一批基線測試
- [x] 完成 `TestPage` 第一輪風格一致化收斂
- [x] 盤點前端測試分布，將正式測試入口統一到 `frontend/tests/` 並分類
- [x] 更新 `AGENTS.md`，明確規定前端測試入口與 Go 測試例外
- [x] 盤點 `.github/instructions/` 在後續 UI/UX 實作中的具體套用點

## UI/UX 收斂階段（依附加 plan 執行）

### Phase 1: Legacy 盤查與 TDD 基線建立 [complete]
- [x] 盤查 SmartDashboard 現役 render 鏈、舊路由、未引用元件
- [x] 識別可安全刪除檔案：`DatalinkLayout.tsx`、`SmartDashboardSidebarTools.tsx`、`SmartDashboardTagAndModbusPanel.tsx`
- [x] 盤查 Legacy 路由與 redirect 策略（`App.tsx`、`legacyRoutes.ts`）
- [x] 驗證 SmartDashboard 主流程測試覆蓋（25 tests passed）
- [x] 記錄 TDD 基線狀態與測試缺口分析

### Phase 2: 建立最小 UI 規範底座 [complete]
- [x] 盤查現有設計 tokens 與樣式系統
- [x] 識別設計系統問題（多套 tokens、元件未統一、表單規範缺失、microcopy 未規範化）
- [x] 建立 `designSystem.ts` 統一設計系統規範
  - 整合 tokens 作為單一來源
  - 提供元件樣式類別（button、card、badge、sectionHeader）
  - 定義表單規範（input、label、error、autocomplete、inputmode、name）
  - 定義 microcopy 規範（loading、button、feedback、ellipsis）
- [x] 驗證 `designSystem.ts` 無 lint 錯誤

**後續實作**（Phase 3+ 時逐步應用）：
- 逐步將現有元件遷移到使用 `designSystem.components.*`
- 補齊表單元件的 `autocomplete`、`inputmode`、`name` 屬性
- 統一 microcopy 使用 `designSystem.microcopy.*`
- 考慮建立統一的 FormInput、FormLabel、FormError 元件

### Phase 3: 重整 SmartDashboard 核心流程 [complete]
- [x] 重新定義唯一主線：建立資料來源 -> 格子檢視 -> Tag 設定 -> Local Modbus -> 資料庫（已在 Sidebar 加入流程指引）
- [x] 將 `WorkflowModal` 收斂為設備入口（非設備 intent 顯示簡化訊息與引導）
- [x] 將 `Sidebar` 改為配合主流程的輕量 stage panel（加入流程指引，顯示當前步驟）
- [x] 將 `WorkspaceContent` 聚焦在「資料來源設定 + 格子可視化」（Flow Status 改為可收折，預設收合）
- [x] 加入資料流向說明（資料庫與本地 Modbus 去向，已加入 CommitPanel）
- [x] 驗證所有測試通過（25 tests passed）

### Phase 4: 整理 LocalModbusWorkbench 與 TestPage [complete]
- [x] `TestPage` 樣式收斂（已完成第一輪，Phase 1 時完成）
- [x] `LocalModbusWorkbenchPage` 依「系統狀態 / Mapping 編輯 / 衝突治理 / 寫入測試」重整區塊層級
  - 重新組織為清晰的區塊結構：系統狀態（含 Server 控制）、Mapping 編輯、衝突治理、寫入測試
  - 改善標題與說明文字，使其更符合單人工作流程
  - 更新測試以匹配新 UI（4 tests passed）

### Phase 5: 可近用性、響應式與品質驗證 [complete]
- [x] 補做鍵盤導覽、focus order、screen reader label、aria-live
  - 在 `SmartDashboardCommitPanel` 與 `LocalModbusWorkbenchPage` 加入 `aria-live="polite"` 與 `role="status"`
  - 在 `SmartDashboardSidebar` 流程指引加入 `role="region"`、`aria-label` 與 `aria-current="step"`
  - 在 `SmartDashboardCommitPanel` 資料流向說明加入 `role="region"` 與 `aria-label`
- [x] 檢視窄螢幕下的 toolbar、側欄、sticky panel、固定寬欄位
  - 記錄響應式設計檢查結果（WorkspaceSection、Sidebar、WorkflowModal、LocalModbusWorkbenchPage）
  - 確認主要元件在窄螢幕下使用 `flex-wrap` 與響應式 grid
- [x] 建立 UI regression / interaction test
  - 建立 `tests/integration/ui/smart-dashboard-regression.test.tsx`
  - 驗證主要 UI 結構渲染正常

## 下一階段（Phase 4 候選 - 舊）
- LocalModbusWorkbenchPage 視覺收斂（TDD 已有，可依 findings TDD 基線先跑測再改）
- SmartDashboard 資訊架構收斂（modal/sidebar 複雜度、主流程認知負荷）
- 頁面測試實體遷移（Gateway、SmartDashboard 自 `src/pages/.../__tests__/` 搬至 `tests/unit/pages/`）

## 錯誤紀錄
- `tests/unit/pages/datalink.test.ts`
  過度聚合 `SmartDashboard` 相關測試後，`vi.mock` 隔離被模組快取污染，導致 `SmartDashboardGridOverlaysSection` 的 mock 未生效；改為更細粒度的 root wrapper 後解決。
- `test:gateway:unit`
  `GatewayQuickSetupPage` 舊測試在 root wrapper 驗證時暴露出既有不穩定斷言；改成重新查詢啟用後的提交按鈕並以獨立 wrapper 執行 page suites 後恢復穩定。
