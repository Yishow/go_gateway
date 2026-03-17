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

**後續實作**（Phase 2 後續 - 進行中）：
- [complete] 補齊表單元件的 `autocomplete`、`inputmode`、`name` 屬性
  - 已完成 `SmartDashboardWorkspaceContent`、`SmartDashboardWorkflowModal`、`LocalModbusWorkbenchPage` 的表單屬性補齊
  - 使用 `designSystem.forms.name.*`、`designSystem.forms.autocomplete.*`、`designSystem.forms.inputmode.*`
- [complete] 統一 microcopy 使用 `designSystem.microcopy.*`
  - 已完成 `LocalModbusWorkbenchPage` loading 狀態統一（`designSystem.microcopy.loading.*`）
  - 已完成 `useSmartDashboardModbusActions` 與 `LocalModbusWorkbenchPage` 的 feedback 訊息統一（`designSystem.microcopy.feedback.*`）
- [complete] 逐步將現有元件遷移到使用 `designSystem.components.*`
  - 已完成 SmartDashboard 中的 icon 按鈕遷移（`SmartDashboardWorkspaceContent`、`SmartDashboardSidebar`）
  - 已完成 SmartDashboard 中的卡片與 badge 遷移（`SmartDashboardCommitPanel`）
  - 其他頁面的按鈕與卡片樣式可在後續逐步遷移
- [complete] 建立統一的 FormInput、FormLabel、FormError 元件
  - 已建立 `frontend/src/components/forms/` 目錄下的統一表單元件
  - `FormLabel`：使用 `designSystem.forms.label.*`，支援必填標記
  - `FormInput`：整合 Label、Input、Error、Hint、Success，使用 `designSystem.forms.*`
  - `FormError`：使用 `designSystem.forms.error.*`，支援圖示
  - `FormHint`：使用 `designSystem.forms.hint.*`
  - `FormSuccess`：使用 `designSystem.forms.success.*`，支援圖示
  - 所有元件已通過 ESLint 驗證

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

## 2026-03-15：Datalink Workbench 重整決策

### Phase A: 深度盤查目前 datalink UI [complete]
- [x] 以平行子任務盤查前端主流程、後端能力、可重用資產與技術債
- [x] 確認目前核心問題集中在頁面層資訊架構與 SmartDashboard God Component
- [x] 確認 `MemoryGrid`、`hooks/services/types/designSystem` 可作為新畫面的重用底座

### Phase B: 收斂改版方向 [complete]
- [x] 與使用者確認第一版偏好：同畫面先完成 `Tag -> Local Modbus`，資料庫作為同流程下一步
- [x] 與使用者確認可視化偏好：格狀視覺化為主，表格為輔
- [x] 比較三條路徑：漸進重整、全新單頁工作台、混合式過渡
- [x] 建議採 **混合式過渡**：新增 `/datalink/workbench`，保留舊頁作 fallback

### Phase C: 設計確認與 spec 文件 [complete]
- [x] 取得使用者對混合式過渡方案與新畫面資訊架構的確認
- [x] 將確認後的設計寫成 spec 文件
- [x] 完成 spec review 並取得使用者同意進入 implementation plan

### Phase D: Implementation planning [complete]
- [x] 產出 `docs/superpowers/specs/2026-03-15-datalink-workbench-design.md`
- [x] 將 spec 以繁中 commit：`5afef10 新增 datalink workbench 設計規格`
- [x] 建立 execution backlog 與依賴關係

### Phase E: 先做 2、3 再做 1 的細化規劃 [complete]
- [x] 補齊 Phase 1 UI 細節：`docs/superpowers/specs/2026-03-15-datalink-workbench-ui-detail.md`
- [x] 補齊 Phase 2 runtime/live value/database target 細化規劃：`docs/superpowers/specs/phase2-runtime-dbtarget-detail.md`
- [x] 再確認 `MemoryGrid` 可沿用、`points poll` 路由已存在，避免 implementation planning 建立在錯誤假設上

### Phase G: Master-detail polish 設計定稿 [complete]
- [x] 延續 quiet desktop v2 後的使用者回饋，重新聚焦 Step 1/2/3 的桌面操作模型
- [x] 以 brainstorming + frontend-design 收斂 master-detail 方向
- [x] 完成設計 spec：`docs/superpowers/specs/2026-03-17-datalink-workbench-master-detail-polish-design.md`
- [x] 完成 spec review loop，修正 data type、Step 2 summary 分工、Step 3 unlock 條件等 drift
- [x] 已取得使用者 review 與核准，進入 implementation planning

### Phase H: Master-detail polish implementation planning [complete]
- [x] 依最新 spec 將實作拆為「共享契約 -> Step 1 -> Step 2 -> Step 3/4 -> 品質驗證」
- [x] 為每個 phase 定義主要檔案、依賴與驗證焦點
- [x] 將 execution backlog 回寫到 planning files 與 SQL todos

### Phase I: Master-detail polish implementation [complete]

#### Phase P0：Shared readiness contract [complete]
- [x] `master-detail-readiness-contract`
  - 對齊 Step 2 右上摘要與 `BottomSummaryBar` 的責任分工
  - 對齊 Step 2/3 的 point lifecycle 語意（`planned/selected` vs `used/linked`）
  - 先穩定 selector / helper / microcopy 契約，再進 UI 改造
  - 已完成：
    - `useWorkbenchSummary.ts` 將 downstream prerequisite 缺口改成 `blocked`
    - `WorkbenchBottomSummaryBar.tsx` 補 `data-reason`
    - shell test 補齊 `no-points -> tag blocked`、`no-tags-linked -> output blocked`
    - 已通過 `workbench-shell-ui.test.tsx`、`npm run lint`、`npx tsc --noEmit`

#### Phase P1：Step 1 device studio [complete]
- [x] `master-detail-step1-studio`
  - `WorkbenchDeviceStep.tsx` 主工作區改成裝置列 + selected-device detail panel 的 master-detail 兩欄
  - 搜尋 / 篩選 / 建立維持單列，保留 scan-first 的 device row browser
  - detail panel 補齊 selected device 的 endpoint、capability summary、last test 與 continue CTA
  - create/edit/clone dialog 改成 desktop 大型置中姿態，保留全螢幕 overlay
  - regression seam 已補：shell inspector query scope 與 foundation master-detail assertions
  - 已通過 `workbench-shell-ui.test.tsx`、`workbench-foundation.test.tsx`、`npm run lint`、`npx tsc --noEmit`、`npm run build`

#### Phase P2：Step 2 source studio [complete]
- [x] `master-detail-step2-creation-flow`
  - `SourceCanvasSection.tsx` 補右上健康摘要與雙主 CTA
  - 補齊 selection toolbar
  - 補齊 rule inline edit 與 point create 雙路徑
  - `source-step-summary`：`ready to create / in conflict / protected`
  - `Create selected points` / `Create rule points` 雙主 CTA
  - point create 主動作自 secondary controls 退位，`More tools` 回到純工具區
  - selection-first 支援 logical span create，並能從 merge continuation 回推 root address
- [x] `master-detail-step2-conflict-types`
  - actionable conflict queue 已落地
  - data type grouped selector 對齊現有 `DataType`
  - `Protect plan` 文案 / 說明已落地
  - 32/64-bit logical cell selection guard 已補齊
  - fixed 16-bit lattice、coverage summary 與 Step 2 inspector seam 已穩定

#### Phase P3：Step 3/4 continuity [complete]
- [x] `master-detail-step3-guidance`
  - Step 3 只在 persisted/usable point 存在時解鎖
  - empty state 會顯示可建立 span 數量與返回 Step 2 的快捷動作
  - `countEligibleSpans` 已對齊 Step 2 `readyToCreateCount`
  - row-board guidance 與 inspector handoff 已對齊
- [x] `master-detail-step4-continuity`
  - Step 3 bind success 與 bound row focus 都會交接真實 `focusedTagIds`
  - Step 4 shared tag selection surface 會優先接續 Step 3 焦點，再 fallback 到既有候選
  - active-target contract 持續成立，未重開另一套 selection UI

#### Phase P4：Quality and rollout [complete]
- [x] `master-detail-polish-regressions`
  - 補/改 unit/page regressions，最終 workbench suite `134/134` 通過
  - `lint` / `tsc --noEmit` / `build` 全數通過
  - 1920×1080 desktop check 已完成，截圖保存在 session artifacts
  - 補齊 deep-link route-sync regression，修正 `?step=...` stale query 會鎖住步驟導航的問題
  - requirement review 與 final code review 皆為 `APPROVED`
  - 文件同步完成，route-fix 已繁中 commit：`efe0e0b 修正 workbench 深連結後步驟導航被鎖定的回歸問題`

### Phase F: Phase 2 契約與輸出目標落地 [complete]
- [x] `runtime-live-value-phase2`
  - 已補齊 runtime status / stream、source runtime summary 與 point poll contract
- [x] `database-target-phase2`
  - 同一個 output step 已可同時承接 Local Modbus 與 Database Target
  - 已完成 dbtarget schema / migration / repo / service / writer / API routes
  - 已完成 reviewer 收尾：secret redaction、`clear_password`、upsert unique 驗證、insert 清空 timestamp、sqlite introspection deadlock、point handler large-list baseline
  - 已完成前後端驗證與 `git diff --check`

## Implementation backlog（已建立 SQL todos）
- `Phase 1 / foundation`
  - `workbench-foundation`
  - 建立 `/datalink/workbench`、`DatalinkWorkbenchPage`、WorkbenchProvider/store、fallback 導流與第一批測試骨架
- `Phase 2 / 可平行 UI 任務`
  - `workbench-shell-ui`
  - `source-canvas-step`
  - `tag-binding-step`
  - `local-modbus-step`
- `Phase 3 / integration`
  - `workbench-integration`
  - 整合來源格位、tag 綁定、local modbus 主線與 validation/readiness
- `Phase 4 / quality`
  - `workbench-quality-pass`
  - 單元/整合/UI regression、a11y、responsive、lint/build/code-review、繁中 commit
- `Phase 5 / phase 2 contracts`
  - `runtime-live-value-phase2`
  - `database-target-phase2`
  - 補 runtime/live value 與 database target

## 2026-03-16：Workbench implementation 進度更新

### 已完成
- [x] `workbench-foundation`
  - `/datalink/workbench` 路由、`DatalinkWorkbenchPage`、`WorkbenchProvider`、`WorkbenchStepNavigator`、`WorkbenchHeaderBar`、`WorkbenchActionDock`
  - foundation route/provider/locale tests
- [x] `frontend-build-blockers`
  - 補上 `designSystem.forms.autocomplete.search`
  - 透過 `npm install --no-audit --no-fund` 補齊前端依賴，解除 `zod` 缺失
- [x] `source-canvas-step`
  - `SourceCanvasSection`、`AddressCanvas`、`AddressLedger`
  - `sourceCanvasModel` 補齊 wide type span、occupied cells、conflict blocking 與 partial failure batch create
- [x] `tag-binding-step`
  - `TagBindingStudio`
  - `tagBindingModel`
  - batch preview key、已連結/既有 key 衝突阻擋、partial failure summary
- [x] `local-modbus-step`
  - `LocalModbusBoard`
  - 將 Local Modbus 輸出整合回 workbench 主線，完成 server 啟停、register 綁定、conflict block sync、push current value
- [x] `workbench-shell-ui`
  - `WorkbenchHeaderBar` / `WorkbenchActionDock` 改為實際摘要卡
  - 新增 `useWorkbenchSummary`，顯示 point/tag/output counts 與下一步建議
- [x] `workbench-integration`
  - 手動 review `useWorkbenchSummary`、`WorkbenchHeaderBar`、`WorkbenchActionDock`、`TagBindingStudio`、`LocalModbusBoard`
  - 確認 Step 2/3/4、shell summary、route wiring、i18n 與 state handoff 已連成完整主線
- [x] `workbench-quality-pass`
  - 以分批測試 + lint/build/diff check 完成最終驗證
  - 背景 code-review / explore / general-purpose agents 持續出現 zero-turn 卡住或 `429`，本輪改由主代理手動收尾

### 下一個 ready phase
- [ ] 收斂 `/datalink/workbench` 與舊 `SmartDashboard` / legacy redirect 的導流策略
- [ ] 視使用者下一步決定是否把 Local Modbus / Database Target 再收斂成單一輸出策略面板

### 2026-03-16：Slice 1 補完（完成）
- [x] `workbench-device-step-redesign`
  - 問題：目前 `device` step 只有 placeholder，導致進頁面後無法設定或選取設備，使用者只看到骨架
  - 依 spec 補齊範圍：
    - `WorkbenchHeaderBar / ContextBar` 的設備上下文與 quick actions
    - Device Step 的 `list/cards + 搜尋/篩選`
    - Step 1 Inspector（設備詳情 / 連線摘要 / 測試連線 / 編輯）
    - embedded create/edit panel flow
    - 連線測試互動與回饋
  - 原則：**不重用舊 `DeviceForm` / `DeviceOnboardingWizard` UI**，僅參考既有 domain 欄位與 hooks / API
  - 已完成內容：
    - `WorkbenchDeviceStep.tsx` 與 `workbenchDeviceFormModel.ts` 新增 workbench 專屬設備設定體驗
    - `/datalink/workbench` 不再落回 `device` placeholder
    - `WorkbenchHeaderBar` 收斂成較接近 spec 的 ContextBar
    - 補齊 zh-TW / en i18n、foundation regression 與其餘 workbench mocks

### 最新驗證摘要
- 已通過：
  - `cd frontend && npm run test -- tests/unit/utils/designSystemForms.test.ts tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-locale.test.ts tests/unit/features/datalink/workbench-source-canvas-model.test.ts tests/unit/features/datalink/tag-binding-model.test.ts tests/unit/features/datalink/legacyRoutes.test.ts --run`
  - `cd frontend && npm run test -- tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx --run`
  - `cd frontend && npm run test -- tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/pages/datalink/local-modbus-workbench.test.ts --run`
  - `cd frontend && npm run lint`
  - `cd frontend && npm run build`
  - `git --no-pager diff --check`
  - `go test ./internal/api -run 'TestNewRouter_RuntimeStatusEndpoint|TestNewRouter_RuntimeStreamEndpointRequiresDeviceID|TestNewRouter_RuntimeStatusUsesRuntimeServiceState' -count=1`
  - `go test ./internal/api/handlers -run 'TestRuntimeStreamHandler_StreamWritesValueEvent|TestPointHandler_Poll$|TestPointHandler_(CreateSyncsRuntimePoint|DeleteRemovesRuntimePoint)|TestMappingHandler_(CreateRefreshesRuntimeMappings|DeleteRefreshesRuntimeMappings)' -count=1`
  - `go test ./internal/datalink/runtime -run 'TestService_SubscribeValueEvents_BroadcastsMatchingPoint|TestHandleCollectedValue_RunPipelineAndWrite|TestHandleCollectedValue_NoMapping_NoWrite' -count=1`
  - `go build ./cmd/test_ui`
  - `PORT=18080 AUTO_OPEN_BROWSER=0 go run ./cmd/test_ui` + `curl http://127.0.0.1:18080/api/v1/datalink/runtime/status`
- 注意：
  - 一次串太多 workbench Vitest 檔案時，曾出現 worker 在測試通過後仍不正常結束的情況；改為分批驗證可穩定完成
  - `go test ./internal/api/...` 全包驗證仍會撞到既有 `TestPointHandler_PollBatch_LargeList`（測試本身用 `4000:` 當位址建立第 10 筆 point），本輪未順手修改該既有 baseline 問題

## 2026-03-16：Workbench desktop redesign round 2 [design approved]
- [x] 以 brainstorming 重新審視 1920×1080 的 workbench 桌面版資訊架構
- [x] 確認回到原始 spec 骨架：`StepRail + ContextBar + PrimaryWorkArea + InspectorPanel + BottomSummaryBar`
- [x] 重新定義 Step 1：`DeviceWorkspace`（能力摘要、Clone、測試歷史）
- [x] 重新定義 Step 2：`AddressCanvasWorkspace`（Plan/Live/Link、多組 rule、templates、coverage/live controls）
- [x] 重新定義 Step 3：`TagBindingBoard`（source→tag 高資訊密度、diff preview、結果總表）
- [x] 重新定義 Step 4：`OutputWorkspace`（Local Modbus + Database 同工作台）
- [x] 核准產出新 spec：`docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-redesign.md`
- [x] scaffold OpenSpec change：`openspec/changes/redesign-datalink-workbench-desktop-flow/`
- [x] 完成 OpenSpec proposal/design/specs/tasks artifacts
- [x] 完成 implementation execution plan（SQL backlog + phase breakdown）

## 2026-03-16：Workbench desktop polish（scan-first）implementation planning

### 已完成
- [x] 重新盤查 workbench 桌面殼層與四個步驟的資訊密度問題
- [x] 以 brainstorming 收斂 3 個方向：`靜音桌面` / `緊湊儀表板` / `嚮導接力`
- [x] 與使用者逐段確認 shell 規則、step-by-step 改法、互動與測試策略
- [x] 產出 spec：`docs/superpowers/specs/2026-03-16-datalink-workbench-desktop-polish-design.md`
- [x] 完成 spec review loop
- [x] 取得使用者核准進入 implementation planning

### 設計決策
- 本輪只做 **scan-first polish**，不再重寫 workflow。
- 優先 bundle：
  - `01` ContextBar 只留步驟摘要與主 CTA
  - `02` BottomSummary 改成細進度條，只高亮目前步驟
  - `05` 所有步驟主 CTA 固定在同一位置
  - `07` Device 搜尋 / 篩選 / 建立收成同一列
  - `11` Source 只保留一條緊湊 toolbar
  - `13` Address canvas 放大成主角，rule list 降權
  - `14` Tag candidate 改成 row board / table-like
  - `15` Tag batch action 只有勾選後才出現
  - `17` Output candidate 只顯示當前 target 必要資訊
  - `20` 各步驟未選取時隱藏不必要操作
- 前端設計 / 提案 / review 若可選模型，優先採用 Opus 品質等級。

### 實作 phase

#### Phase P0：Shell density polish [pending]
- [x] `scan-shell-density`
  - 調整 `WorkbenchContextBar`
  - 調整 `WorkbenchBottomSummaryBar`
  - 固定主 CTA 區域
  - 更新 shell / foundation regression tests
  - 已完成：
    - ContextBar 改成 `step summary + single primary action`
    - 移除預設 capability chip / latest test 狀態常駐顯示
    - BottomSummary 新增 `active / compact` emphasis
    - primary CTA 改回逐步流程：`device -> source -> tag -> output`
    - source / tag 未就緒時，CTA 直接 disabled，避免跳步
    - output step CTA 改成 output-focused 文案，不再退回 `switchDevice`
    - output-ready 狀態下，CTA 會把焦點拉回輸出主控區，不再是 enabled no-op
    - shell / foundation tests 與 scan-first spec 對齊

#### Phase P1：Step surfaces（可平行） [complete]
- [x] `scan-device-browser`
  - Step 1 改成 single-row primary toolbar（search / protocol / status / refresh / create）
  - 移除 main surface 的 clone / continue cluster，前進改由 shell CTA 承接
  - device list 改成 compact row layout，只保留 identity / endpoint / health / 兩個 capability hints
  - inspector 持續承接 edit / clone / test 細節
- [ ] `scan-source-toolbar-canvas`
- [x] `scan-source-toolbar-canvas`
  - Step 2 收成 `source-primary-toolbar`
  - jump / freeze / snapshot / save-load / batch create 改進 `source-secondary-controls`
  - `AddressCanvas` workspace 標成 primary，rule layer 降為 supporting
  - coverage overview 留在 primary workspace，但改成 lighter supporting strip
- [x] `scan-tag-row-board`
  - Step 3 厚卡改 row board
  - batch action 僅在 selection 後出現
  - conflict / already-linked 細節進 inspector
- [x] `scan-output-active-target`
  - Step 4 row 只顯示 active target 必要資訊
  - Modbus / DB 次層資訊分組降權
  - non-active-target detail 移入 inspector

#### Phase P2：Cross-step regression + rollout [complete]
- [x] `scan-cross-step-regressions`
  - progressive disclosure regression
  - `step/target` deep-link regression
  - `1920×1080` toolbar / CTA / shell 穩定性驗證
  - 文件、OpenSpec tasks、繁中 commit、code review

## 2026-03-17：quiet desktop v2 收尾 [complete]
- [x] `ui-compact-step1`
  - 移除 Step 1 裝飾 hero
  - 收斂 device row density 與 capability pills
  - device create/edit panel 改成 viewport overlay
  - inspector aside 支援捲動
- [x] `ui-guide-step2`
  - planner 搬進 rule layer
  - `apply` 改成 `addRule`
  - 加入 delete rule
  - utility tools 收進 `moreTools`
  - canvas / coverage / rule layer 主次層級重排
- [x] `ui-align-step3`
  - 補 flow mode / selection guidance
  - row-board 扁平化 raw/transformed value
  - existing mode 移除 preview，保留 existing tag select
  - sidebar metrics 壓縮
- [x] `ui-compact-step4`
  - shared candidate board 成為唯一 tag selection surface
  - Modbus / Database 改成 selected-tag summary
  - Database connector form 預設收合
- [x] `ui-regressions-quiet-desktop-v2`
  - targeted workbench regression
  - lint / tsc / build
  - code review feedback fix-up

### Implementation Phase Breakdown（Redesign Round 2）

**Phase 0: Shell Infrastructure（序列，已完成）**
- [x] `redesign-shell-frame`：以五區 shell 取代現有 HeaderBar/ActionDock/StepNavigator
- [x] `redesign-shell-state`：擴充 WorkbenchProvider，加入 readiness/inspector/output-target/cross-step 狀態
  - 已補齊 `WorkbenchProvider` 的 `inspectorSelection`、`activeOutputTarget`、`crossStepContext`
  - 已補齊 `workbenchTypes.ts` readiness / inspector / cross-step contracts
  - 已補齊 `useWorkbenchSummary.ts` per-step readiness 匯總
  - 驗證：`cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-readiness.test.ts tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-shell-ui.test.tsx tests/unit/pages/datalink/workbench-foundation.test.tsx`
  - 驗證：`cd frontend && npm run lint`
  - 驗證：`cd frontend && npx tsc --noEmit`

**Phase 1: 四個 Step Workspace（可平行，依賴 Phase 0）**

Step 1 — DeviceWorkspace：
- [x] `redesign-device-workspace`：DeviceBrowser（搜尋/篩選/建立/重新整理/cards）
- [x] `redesign-device-inspector`：Device inspector + ContextBar 能力摘要
- [x] `redesign-device-clone`：Clone flow drawer（依賴 device-workspace）
  - 已完成卡片 capability summary、右側 inspector 詳情、context bar capability chips、clone drawer
  - 已以 session-local recent test history 補足 inspector 最近三次測試時間線
  - 驗證：`cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-readiness.test.ts tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-foundation.test.tsx tests/unit/pages/datalink/workbench-shell-ui.test.tsx`
  - 驗證：`cd frontend && npm run lint`
  - 驗證：`cd frontend && npx tsc --noEmit`
  - 驗證：`cd frontend && npm run build`

Step 2 — AddressCanvasWorkspace：
- [x] `redesign-source-rule-model`：Rule 資料模型 + RuleLayerBar
- [x] `redesign-address-canvas`：連續 16-bit lattice + 多 rule merged spans（依賴 rule-model）
- [x] `redesign-source-viewmodes`：Plan/Live/Link + PlannerToolbar（依賴 canvas）
- [x] `redesign-source-templates`：本地 template 持久化（依賴 rule-model）
- [x] `redesign-source-inspector`：Step 2 inspector + audit surface（依賴 canvas）
  - 已完成 RuleLayerBar、連續 gap item model、coverage overview、Plan/Live/Link、value format、freeze/snapshot、jump to address、audit drawer
  - 已補齊 source rules 跨步驟 persistence；從 Source 切到 Tag/Output 再回來，不會把已套用規則清空
  - 已補齊 source template 本地儲存 / 載入、view mode 回填、capability snapshot warning 與 template-origin provenance
  - 已補齊 shared inspector 的 rule / span detail，會顯示 coverage、origin、bit width、raw value、link state
  - 2026-03-16 fleet audit 補記：`AddressCanvas.tsx` 目前仍是卡片式 responsive grid，尚未完全符合 OpenSpec `3.2` 要求的固定 16-bit lattice 幾何；此項需在 Phase 2 收尾前補做對齊。
  - 驗證：`cd frontend && npm run test -- --run tests/unit/features/datalink/workbench-source-canvas-model.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
  - 驗證：`cd frontend && npm run test -- --run tests/unit/features/datalink/sourcePlannerContract.test.ts tests/unit/features/datalink/sourceTemplateStorage.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
  - 驗證：`cd frontend && npm run lint`
  - 驗證：`cd frontend && npx tsc --noEmit`
  - 驗證：`cd frontend && npm run build`

Step 3 — TagBindingBoard：
- [x] `redesign-tag-board`：Dense board 主面
- [x] `redesign-tag-batch`：Batch diff preview + 結果總表（依賴 tag-board）
- [x] `redesign-tag-inspector`：Step 3 inspector（依賴 tag-board）
  - 已完成 dense candidate rows、create/existing flow switch、search/status filter、raw/transformed value 區、status badge、existing tag 選擇
  - 已補齊 batch diff preview、skip reasons、create/existing result summary counters 與對應 model helper
  - 已補齊 shared inspector 的單筆 tag detail，會顯示 preview/bound tag key、source address、raw/transformed value、span/bit width、衝突提示
  - 驗證：`cd frontend && npm run test -- --run tests/unit/features/datalink/tag-binding-model.test.ts tests/unit/pages/datalink/workbench-tag-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`
  - 驗證：`cd frontend && npm run test -- --run tests/unit/features/datalink/sourcePlannerContract.test.ts tests/unit/features/datalink/sourceTemplateStorage.test.ts tests/unit/features/datalink/tag-binding-model.test.ts tests/unit/features/datalink/workbench-provider.test.tsx tests/unit/features/datalink/workbench-locale.test.ts tests/unit/pages/datalink/workbench-source-step.test.tsx tests/unit/pages/datalink/workbench-tag-step.test.tsx`
  - 驗證：`cd frontend && npm run lint`
  - 驗證：`cd frontend && npx tsc --noEmit`
  - 驗證：`cd frontend && npm run build`

Step 4 — OutputWorkspace：
- [x] `redesign-output-board`：OutputCandidateBoard + TargetSwitcher
- [x] `redesign-modbus-studio`：RegisterMapCanvas + auto-map + dry-run（依賴 output-board）
- [x] `redesign-database-studio`：Schema snapshot + mapping UI（依賴 output-board）
- [x] `redesign-output-inspector`：Source→tag→output trace（依賴 output-board）
  - 已完成 shared candidate board，單列同時顯示 Local Modbus / Database 狀態
  - 已完成 target switcher，能在同一個 Step 4 內切換 active Modbus / Database studio
  - shared board 的 Modbus badge 改為直接顯示 `HRxx`，較適合桌面掃描
  - 已完成 `RegisterMapCanvas` 主畫面接線，補齊 auto-map 三種策略與 dry-run result surface
  - 已完成 Database schema snapshot / required highlight / write-row preview
  - 已完成 Step 4 inspector trace，候選項可顯示 `source -> tag -> modbus/database` readiness
  - 驗證：`cd frontend && npm run test -- --run tests/unit/pages/datalink/workbench-output-step.test.tsx tests/unit/features/datalink/workbench-locale.test.ts`（25 tests passed）
  - 驗證：`cd frontend && npm run lint`
  - 驗證：`cd frontend && npx tsc --noEmit`
  - 驗證：`cd frontend && npm run build`

**Phase 2: Quality + Rollout（依賴所有 Step 完成）**
- [x] `redesign-address-canvas`（spec drift revisit）：補齊固定 16-bit lattice 幾何，校正 OpenSpec `3.2`
- [x] `redesign-i18n-a11y`：i18n + 鍵盤導覽 + aria
- [x] `redesign-regression-tests`：1920×1080 regression + 功能覆蓋
- [x] `redesign-legacy-compat`：Legacy route 相容策略

### Parallelism Map

```
Phase 0 (serial):
  shell-frame → shell-state

Phase 1 (4 parallel tracks after Phase 0):
  Track A: device-workspace ──┬─→ device-clone
                              └─→ device-inspector
  Track B: source-rule-model ─┬─→ address-canvas ──┬─→ source-viewmodes
                              │                     └─→ source-inspector
                              └─→ source-templates
  Track C: tag-board ─────────┬─→ tag-batch
                              └─→ tag-inspector
  Track D: output-board ──────┬─→ modbus-studio
                              ├─→ database-studio
                              └─→ output-inspector

Phase 2 (serial after all tracks):
  i18n-a11y → regression-tests → legacy-compat
```

### 實作原則
- 不重用舊 presentational components，僅重用 hooks/services/types/helper
- 所有新 UI 元件以 Opus-first 品質為目標
- 每個 todo 含對應測試，測試先行
- 每個 step 完成後獨立驗證 lint + build
- i18n 字典隨各 step 同步更新，Phase 2 做最終稽核

### 當前 ready queue
- Phase 2 已完成，下一步進入 code review / commit / OpenSpec tasks 同步
