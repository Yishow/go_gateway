# Findings

## 2026-03-15

### datalink UI 深度分析：快速 spot-check
- `frontend/src/pages/datalink/SmartDashboardPage.tsx` 仍是主入口，單檔很大，並同時承擔流程控制、URL intent、面板切換、overlay、匯入匯出、commit 等多重責任。
- `SmartDashboard` 仍透過 query 參數與 `legacyRoutes.ts` 承接多種 modal intent，代表畫面入口與使用者真正想做的主線任務還有歷史包袱混在一起。
- `SmartDashboardSidebar` 目前包含 `plan` / `tag` / `modbus` / `commit` 四分頁，主流程被拆成多個區域與狀態來源，對單純工作流不夠聚焦。
- `SmartDashboardWorkflowModal`、`SmartDashboardOverlays`、`SmartDashboardPanels`、`SmartDashboardGridOverlaysSection` 顯示目前交互大量依賴 modal / slide panel / popover 疊加。
- `LocalModbusWorkbenchPage.tsx` 仍是獨立頁面，代表 Tag -> Local Modbus 的後段操作沒有真正整合回單一工作台。
- Repo 中仍存在 `SmartDashboardSidebarTools.tsx`、`SmartDashboardTagAndModbusPanel.tsx`、`DatalinkLayout.tsx` 等 legacy/未引用檔案，說明 datalink UI 歷經多輪演變，維護認知成本偏高。

### 本輪分析方向
- 優先確認是否應改成單頁主線工作台：來源設定 -> 記憶體/位址視覺化 -> Tag 綁定 -> 輸出目標（Local Modbus / DB）。
- 優先保留既有資料模型、hooks、測試與 design system；除非結構耦合太深，再考慮另開新頁取代 `SmartDashboardPage`。

### SmartDashboard / LocalModbus 原始碼補讀
- `SmartDashboardPage.tsx` 一開頭就引入大量 hook、service、feature helper 與多種 UI 容器，顯示主畫面仍是 orchestration god component，維護成本高。
- 頁面 state 同時管理：device 切換、grid 選取、planner、batch naming、template、guide stage、active tab、sidebar tab、workflow modal、create/edit/delete device、import/export 等，說明目前單一主線被太多旁支任務共用同一個畫面狀態。
- `legacyRoutes.ts` 仍定義 `points` / `mappings` / `wizard` 等 legacy intent，且透過 `?legacy=`、`?modal=`、`?section=` 注入 SmartDashboard，表示新舊入口尚未真正切開。
- `LocalModbusWorkbenchPage.tsx` 雖然本身流程相對單純，但仍是獨立頁，且以「映射管理工具」思維設計，不是從來源規劃一路往下完成輸出綁定的單頁體驗。
- `LocalModbusWorkbenchPage` 已具備可重用的後段能力：讀 status、列 mappings、upsert/delete mapping、sync、write test、匯入匯出；若新 UI 要重做，這些 API 與互動可以保留，但應改成嵌入同一條主流程，而不是跳另一頁。

### 使用者需求與偏好確認
- 第一版新流程：同一畫面先完成 `Tag -> Local Modbus`，資料庫保留在同流程中的下一步。
- 來源位置可視化：採「格狀視覺化為主，表格為輔助編輯」。

### 三份盤查整合結論
- **最嚴重問題** 是主線被切碎：SmartDashboard 以 Dashboard + Sidebar Tabs + Modal/Popover/SlidePanel 組成，無法讓使用者連續完成「來源 -> 可視化 -> Tag -> 輸出」。
- **頁面層技術債** 集中在 `SmartDashboardPage.tsx` 與其高度耦合的 SmartDashboard* 子元件；底層 `hooks/services/types/MemoryGrid/designSystem` 大多可以沿用。
- **LocalModbusWorkbenchPage** 的 domain 能力完整，但它以獨立工具頁存在，造成主線後段斷裂。
- **後端大致夠用**：Device / Point / Tag / Mapping / PollingGroup / ModbusShare CRUD 與前端 service/hook 可直接承接；大 blocker 是 runtime HTTP API 缺失、point poll stub 與 database target 契約尚未成形。

### 改版 approach 對照
- **Approach A：漸進重整既有 SmartDashboard**
  - 優點：短期改動小、較少新路由與導流成本。
  - 缺點：仍被既有 Dashboard/Tab/Modal 架構限制，主線體感改善有限。
- **Approach B：全新單頁工作台**
  - 優點：最能對齊使用者主線，資訊架構最乾淨。
  - 缺點：開發量最大，過渡期要維護新舊兩套入口。
- **Approach C：混合式過渡（推薦）**
  - 先新增 `/datalink/workbench` 作為新入口，不碰舊 `SmartDashboardPage`；
  - Phase 1 先完成純前端主線：Device -> Grid/Point -> Tag -> Local Modbus；
  - Phase 2 再接 runtime 即時值與 database target。

### 推薦方案
- **推薦採用混合式過渡**：新做單頁工作台，但以漸進方式導入。
- 理由：
  - 直接對齊使用者偏好與主線任務。
  - 可大量重用既有底層資產，避免重寫 domain 層。
  - 不必在 1498 行的 God Component 上做高風險大手術。
  - 第一階段不受後端 runtime API 缺口阻擋，能先交付可操作的新流程骨架。

### 補充風險盤查（implementation planning 前）
- `MemoryGrid` 目前 props 已具備新 workbench 需要的大部分核心能力：`selectedAddresses`、`onSelect`、`onCellClick`、`plannedAllocations`、`linkedAddresses`、`showConflictsOnly`。這代表格狀主視覺可以沿用，但 **table mode 仍需要獨立的新元件**，不在 `MemoryGrid` 本身內處理。
- 後端 router 已實際掛上：
  - `POST /datalink/points/:id/poll`
  - `POST /datalink/points/poll`
  - `GET /datalink/preview/stream`
- 因此 spec review 提到的 `points/:id/poll` 路由疑慮已排除；Phase 2 真正要補的是更高階的 runtime/status/stream 契約，而不是 point poll 路由不存在。

### 2026-03-16：Database Target Phase 2 收尾發現
- `LocalModbusBoard` 下方直接嵌入 `DatabaseTargetBoard`，能維持「來源 -> Tag -> 輸出」單一路徑，不需要再跳另一個工具頁。
- `timestamp_column` 必須明確收斂成 **只有 upsert 模式有效**；若 insert 模式仍殘留舊值，UI 顯示與 runtime writer 實際行為會分離。
- Secret redaction 之後，空字串不能同時代表「保留舊密碼」與「清除密碼」；若要支援編輯既有 connector，必須提供額外明確訊號（本輪採 `clear_password`）。
- SQLite schema introspection 在 `MaxOpenConns=1` 情境下，不能一邊遍歷 `PRAGMA index_list` rows 一邊再查 `PRAGMA index_info`；必須先收完 index names 再逐一查 detail，否則會自我阻塞。
- `point_handler_extended_test.go` 的 large-list baseline 問題是測試資料生成錯誤（第 10 筆位址被拼成 `4000:`），不是 handler 主邏輯。

### 2026-03-16：Workbench Device Step 重做發現
- `/datalink/workbench` 進頁只有骨架的直接原因，不是資料沒載入，而是 `device` step 根本還停留在 placeholder render path。
- 這個問題不能用「把舊 `DeviceForm` / `DeviceOnboardingWizard` 塞回來」快修；使用者已明確要求依 spec 重做，且 spec 也要求 Step 1 應是 `cards/list + inspector + embedded panel` 的工作台式體驗。
- `ContextBar` 若只做摘要數字卡，不足以承接 Step 1 的「立即顯示選中設備上下文」需求；至少要把 selected device / protocol / status / last test 與 quick actions 放進同一個頂部脈絡列。
- 設備編輯流程有一個容易漏掉的契約：`description` 清空時前端不能送 `undefined`，否則後端 pointer update 會把它解讀成「不更新」，舊描述會殘留。
- Vite proxy error 不代表 workbench UI regression；在本輪 browser smoke 中，實際 root cause 是 dev server 還活著，但 `http://localhost:8080` 後端未啟動。
- browser smoke 額外揭露兩個值得收斂的 UI polish：
  - Step 1 已選設備時，`ActionDock` 不能再顯示「先選設備」；next action 要改成前往 Source。
  - `last_test_at` 若來自後端 zero timestamp（`0001-01-01T00:00:00Z`），前端要顯示 fallback microcopy，而不是直接把零值時間戳丟給使用者。

### 2026-03-16：Source templates / Tag batch 收斂發現
- Source template 不能只存 planner 欄位；若不一併保存 `preferredViewMode`，重新套用模板後會回不到使用者儲存時的工作上下文。
- Source template 也不能把能力摘要直接存成翻譯後文字；應改存 `capabilitySnapshot` 這種穩定結構，才能跨語系比較並在套用時產生 warning。
- `SourceTemplateRecord` 升 schema 時要保留向下相容，否則會直接打壞仍在使用舊 template shape 的 SmartDashboard flow。
- Tag batch 新增 diff preview 後，原本用純文字搜尋 row 的測試會被 preview 面板污染；測試應改鎖定 candidate row `data-testid`，精準驗證 filter 行為而不是整頁文字存在性。

## 2026-03-08

### 規範來源盤查
- `AGENTS.md` 已作為專案共通規範入口，並已納入 UI/UX 主流程目標、TDD 與 `planning-with-files`。
- `CLAUDE.md` 已更新為：先讀 `AGENTS.md`，再依檔案類型讀 `.github/instructions/`，最後讀 `CLAUDE.md`。
- `GEMINI.md` 已改寫成與目前 repo 現況一致的 Gemini Agent 補充文件，並納入同樣的規範分工。
- `.github/instructions/` 目前共有三份正式規範：
  - `go.instructions.md`
  - `reactjs.instructions.md`
  - `typescript-5-es2022.instructions.md`

### 建議的規範分工
- `AGENTS.md` 管專案怎麼做。
- `CLAUDE.md` / `GEMINI.md` 管 Agent 怎麼工作。
- `.github/instructions/` 管程式怎麼寫。

### 已完成的文件對齊
- `AGENTS.md`
  新增「規範來源分層」，明確宣告 `.github/instructions/` 為正式實作規範來源，並定義優先順序。
- `CLAUDE.md`
  補上文件分工、閱讀順序與規範優先順序。
- `GEMINI.md`
  從舊版原型描述整理為目前專案可直接採用的 Agent 補充文件。

### 與目前 UI/UX 計畫的關聯
- `SmartDashboard` 應收斂成單人主工作流。
- `LocalModbusWorkbenchPage` 應作為主流程後段工作台。
- `TestPage` 保持工程工具定位，只做風格一致化。
- 後續前端實作需要同時遵守 React 與 TypeScript instructions。

### TDD 基線盤查
- `SmartDashboard` 已有多組互動與 hook 測試，包含：
  - `frontend/src/pages/datalink/__tests__/SmartDashboard.interaction.test.tsx`
  - `frontend/src/pages/datalink/__tests__/SmartDashboardGridOverlaysSection.test.tsx`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardCommitFlow.test.ts`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardWorkspaceContentState.test.tsx`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardWorkspaceState.test.ts`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardPanelsState.test.ts`
- `LocalModbusWorkbenchPage` 已有基礎頁面測試：
  - `frontend/src/pages/datalink/__tests__/LocalModbusWorkbenchPage.test.tsx`
- 目前未找到 `TestPage` 對應測試檔，代表後續若要做樣式一致化，必須先建立 `TestPage` 測試基線。

### 已補上的 TestPage 測試基線
- 已新增 `frontend/src/pages/__tests__/TestPage.test.tsx`
- 目前覆蓋的穩定頁面行為：
  - 初始化時會批量最小化非 config 卡片
  - 連線成功後會自動收折 config 卡片，且可手動展開
  - 協議切換時至少能穩定驗證模式同步到 `serial`
- 測試刻意聚焦頁面殼層與使用者可觀察行為，避免綁定過多內部實作細節，方便後續做 UI 風格一致化。

### TestPage 第一輪 UI 收斂
- `frontend/src/pages/TestPage.tsx` 已先完成低風險的外殼收斂：
  - 頁面背景改為與 datalink 主介面相近的深色漸層
  - 頂部加入工程工作台定位說明
  - 主要卡片表面統一為深色面板語言
  - 將本次碰到的 `transition-all` 改為更明確的 transition 屬性
- 此輪不改變 `TestPage` 的工具定位與核心流程，只收斂視覺語言與資訊層級。

### 前端測試入口重整
- 目前前端測試原本散落在：
  - `frontend/src/**/__tests__/*`
  - `frontend/src/**/*.test.ts(x)`
  - `frontend/tests/e2e/*`
- 已建立新的正式測試入口與分類：
  - `frontend/tests/unit/`
  - `frontend/tests/integration/`
  - `frontend/tests/e2e/`
- 本輪採用「root 測試入口收斂 + 舊測試內容逐步遷移」策略：
  - 先讓 Vitest 只收斂執行 `frontend/tests/unit` 與 `frontend/tests/integration`
  - 既有 source-adjacent 測試先由 root wrapper 匯入
  - `TestPage` 測試則已實際移到 `frontend/tests/unit/pages/test-page.test.tsx`
- `Playwright` 原本就使用 `frontend/tests/e2e/`，此次維持不變。
- `Go` 測試不適合集中到 root `tests/`，需維持 `*_test.go` 與實作檔相鄰，這一點已同步寫進 `AGENTS.md`。

### Root 測試入口的實務發現
- 若把多個有 `vi.mock` 隔離需求的 page test 直接聚合到同一個 wrapper，會因模組先後載入順序讓 mock 失效。
- 因此 root `tests/` 的分類不應只看資料夾層級，還要考慮測試隔離邊界：
  - `Gateway` page 測試拆成 `expert-workbench.test.ts` 與 `quick-setup.test.ts`
  - `SmartDashboard` page 測試拆成 interaction / grid-overlays / state 三組 wrapper
- 這個粒度可同時滿足：
  - 測試入口集中在 root `tests/`
  - 類型與領域有清楚分類
  - `vi.mock` 與模組快取不互相污染

### TDD 基線規劃（SmartDashboard / LocalModbusWorkbench / TestPage）

| 頁面 | 測試入口 | 覆蓋範圍 | 基線狀態 |
|------|----------|----------|----------|
| **TestPage** | `tests/unit/pages/test-page.test.tsx` | 批次最小化、config 收折/展開、協議切換模式同步 | 已完成實體遷移，3 tests |
| **SmartDashboard** | `tests/unit/pages/datalink/smart-dashboard-*.test.ts` | interaction（section/modal intent、設備切換、unsaved guard、draft 啟用）、grid-overlays（popover、delete）、state hooks（commit flow、panels、workspace） | wrapper 匯入 src，21 tests |
| **LocalModbusWorkbenchPage** | `tests/unit/pages/datalink/local-modbus-workbench.test.ts` | 啟動/停止 server、register 衝突阻擋、dashboard 返回連結 | wrapper 匯入 src，4 tests |

**後續 UI 變更前**：先跑對應頁面的測試確保基線通過，再改樣式或 IA。SmartDashboard / LocalModbusWorkbench 若要收斂視覺，需維持或補強現有測試。

### .github/instructions 在 UI/UX 實作中的套用點

| 規範檔 | 套用於 | 具體 touchpoint |
|--------|--------|-----------------|
| **reactjs.instructions.md** | `*.tsx`、`*.jsx`、`*.css`、`*.scss` | 元件設計（單一職責、composition）、hooks 依賴陣列、state 結構、styling（CSS 變數、responsive、ARIA）、form 受控元件、routing、測試用 RTL、accessibility |
| **typescript-5-es2022.instructions.md** | `*.ts` | props/state 型別、 discriminated unions（例如 flow 狀態機）、`unknown` + narrowing、避免 `any`、utility types、async 錯誤處理、命名與格式 |
| **go.instructions.md** | 後端、API | 非前端 UI/UX 直接套用，但 API 契約會影響前端型別與 data fetching |

**實務對齊**：新增或修改 React 元件時，同時遵守 reactjs（元件設計、hooks、樣式）與 typescript（型別、async）兩份 instructions；既有的 `react-hooks/exhaustive-deps`、`no-unused-vars` 等 lint 規則與 instructions 一致。

### SmartDashboard Legacy 盤查結果

#### 未使用檔案（可安全刪除）

| 檔案路徑 | 狀態 | 說明 |
|---------|------|------|
| `frontend/src/layouts/DatalinkLayout.tsx` | **未引用** | 舊版 layout，目前 SmartDashboard 直接渲染，不使用 Outlet |
| `frontend/src/pages/datalink/smart-dashboard/SmartDashboardSidebarTools.tsx` | **未引用** | 功能已整合到 `SmartDashboardSidebar.tsx` 的 Mini Toolbar 區塊 |
| `frontend/src/pages/datalink/smart-dashboard/SmartDashboardTagAndModbusPanel.tsx` | **未引用** | `SmartDashboardSidebar` 直接使用 `SmartDashboardTagPanel` 與 `SmartDashboardModbusPanel`，不再需要組合層 |

#### Legacy 路由與 Redirect（App.tsx）

| 路由 | 目標 | 類型 | 建議 |
|------|------|------|------|
| `/datalink/dashboard-legacy` | `/datalink` | 直接重定向 | 可移除，無外部引用 |
| `/datalink/devices-legacy` | `/datalink?modal=devices` | Modal intent | 可移除，`/datalink/devices` 已提供相同功能 |
| `/datalink/points` | `/datalink?legacy=points&modal=points` | Legacy migration | 保留，可能仍有外部書籤 |
| `/datalink/mappings` | `/datalink?legacy=mappings&modal=mappings` | Legacy migration | 保留，可能仍有外部書籤 |
| `/datalink/wizard` | `/datalink?legacy=wizard&modal=wizard` | Legacy migration | 保留，可能仍有外部書籤 |
| `/datalink/devices` | `/datalink?modal=devices` | Modal intent | 保留，現役路由 |
| `/datalink/polling-groups` | `/datalink?modal=polling-groups` | Modal intent | 保留，現役路由 |
| `/datalink/tags` | `/datalink?modal=tags` | Modal intent | 保留，現役路由 |
| `/datalink/settings` | `/datalink?modal=settings&section=settings` | Section + Modal | 保留，現役路由 |

#### SmartDashboard 現役結構

**主入口**：
- `SmartDashboard.tsx`：僅 re-export `SmartDashboardPage.tsx`
- `SmartDashboardPage.tsx`：主元件（1498 行），包含完整狀態管理與業務邏輯

**子元件層級**：
- `SmartDashboardHeader`：頂部標題與狀態顯示
- `SmartDashboardIntentNotices`：section/modal intent 通知
- `SmartDashboardControlBar`：設備選擇與建立按鈕
- `SmartDashboardWorkspaceSection`：工作區（記憶體格 + 側欄）
- `SmartDashboardWorkflowModal`：設備管理中心 modal
- `SmartDashboardOverlays`：切換確認、刪除確認、建立設備 modal
- `SmartDashboardPanels`：BatchCreate、Shortcuts、PointDetail 滑出面板
- `SmartDashboardGridOverlaysSection`：格位 popover（Tag 連結、點位詳情）
- `ImportDialog` / `ExportDialog`：匯入/匯出對話框

**側欄結構**（`SmartDashboardSidebar`）：
- Mini Toolbar：Workbench、Import/Export、Undo/Redo、Shortcuts
- Tab Bar：`plan` / `tag` / `modbus` / `commit` 四分頁
- Tab Content：對應 `SmartDashboardPlanningTab`、`SmartDashboardTagPanel`、`SmartDashboardModbusPanel`、`SmartDashboardCommitPanel`

**狀態管理 hooks**：
- `useSmartDashboardWorkspaceState`：工作區狀態
- `useSmartDashboardTagLinking`：Tag 連結邏輯
- `useSmartDashboardModbusActions`：Modbus 操作
- `useSmartDashboardCommitFlow`：提交流程
- `useSmartDashboardGridOverlays`：格位 overlay 狀態
- `useSmartDashboardSidebarPanelMotion`：側欄動畫

#### Legacy 相容層評估

**`legacyRoutes.ts`**：
- 定義 `LEGACY_DECOMMISSION_ROUTES`：`['points', 'mappings', 'wizard']`
- 定義 `DASHBOARD_MODAL_INTENTS`：`['devices', 'settings', 'points', 'mappings', 'wizard', 'polling-groups', 'tags']`
- 提供 `buildLegacyMigrationRedirect`、`buildDashboardModalRedirect` 等 helper
- **狀態**：現役使用中，`SmartDashboardPage` 透過 URL query 參數解析 intent

**建議**：
- 三個未使用檔案可立即刪除
- Legacy 路由可先觀察 1-2 週使用量，再決定是否移除
- `legacyRoutes.ts` 需保留，但可簡化為只支援現役 modal intents

### TDD 基線狀態（Phase 1 驗證）

#### SmartDashboard 主流程測試覆蓋

| 主流程 | 測試入口 | 覆蓋狀態 | 測試數量 |
|--------|----------|----------|----------|
| **選設備** | `smart-dashboard-interaction.test.ts` | ✅ 完整 | 6 tests（設備切換、unsaved guard、draft 啟用） |
| **規劃來源** | `smart-dashboard-interaction.test.ts` + `useSmartDashboardWorkspaceContentState.test.tsx` | ✅ 完整 | 4 tests（typed occupancy、planner bindings、grid view 優先） |
| **格子顯示/選取** | `smart-dashboard-grid-overlays.test.ts` | ✅ 完整 | 2 tests（popover tag panel、context menu delete） |
| **Tag 設定** | `smart-dashboard-grid-overlays.test.ts`（透過 TagPanel mock） | ⚠️ 部分 | TagPanel 本身未獨立測試，但透過 grid overlay 有覆蓋 |
| **前往 Local Modbus Workbench** | `local-modbus-workbench.test.ts` | ✅ 完整 | 4 tests（啟動/停止 server、衝突阻擋、返回連結） |

**測試執行結果**（2026-03-08）：
- `tests/unit/pages/datalink/`：**25 tests passed**（4 test files）
- 包含：interaction（10 tests）、grid-overlays（2 tests）、state hooks（13 tests）、local-modbus（4 tests）

#### 測試缺口分析

**已覆蓋**：
- ✅ 設備選擇與切換流程
- ✅ 來源規劃（modbus area、data type、contiguous rules）
- ✅ 記憶體格選取與 context menu
- ✅ Local Modbus Workbench 基本操作
- ✅ Commit flow（validation、activation、rollback）
- ✅ Workspace state management

**部分覆蓋**：
- ⚠️ Tag 設定：目前透過 `SmartDashboardGridOverlaysSection` 測試 TagPanel 的渲染與關閉，但缺少：
  - Tag 連結到格位的完整流程
  - 建立新 Tag 並連結
  - Tag 編輯（display name、unit、description）
  - Tag 與 Mapping 的影響範圍計算

**建議補強**（Phase 2 或後續）：
- 新增 `SmartDashboardTagPanel` 獨立測試，覆蓋 tag linking、creation、editing 流程
- 或擴充 `smart-dashboard-interaction.test.ts` 加入端對端 tag 設定場景

#### TDD 基線結論

**現有測試基線足以支撐 Phase 2 的 UI 規範建立與 Phase 3 的 SmartDashboard 核心流程重整**：
- 主要互動流程（設備選擇、來源規劃、格位選取）已有穩定測試保護
- State hooks 與 commit flow 有完整覆蓋
- 後續 UI 變更可先跑現有測試確保基線通過，再進行樣式或資訊架構調整

**風險**：
- Tag 設定流程的測試覆蓋較弱，若 Phase 3 要大幅調整 TagPanel UI，需先補強對應測試

### Phase 2：設計系統基礎盤查

#### 現有設計 tokens 與樣式系統

| 檔案 | 用途 | 狀態 |
|------|------|------|
| `frontend/src/styles/tokens.ts` | 基礎設計 tokens（色彩、間距、圓角、陰影、字體） | ✅ 現役 |
| `frontend/src/pages/gateway/styleTokens.ts` | Gateway 頁面專用 tokens（light mode 導向） | ⚠️ 僅 Gateway 使用 |
| `frontend/src/components/ui/button.tsx` | shadcn/ui Button 元件（使用 CVA） | ✅ 現役 |
| `frontend/src/components/ui/alert.tsx` | shadcn/ui Alert 元件 | ✅ 現役 |
| `frontend/src/components/ui/dialog.tsx` | shadcn/ui Dialog 元件 | ✅ 現役 |
| `frontend/src/index.css` | Tailwind 基礎設定、scrollbar、動畫 | ✅ 現役 |

#### 設計系統問題

1. **多套 tokens 並存**：
   - `tokens.ts` 為 dark mode 導向（slate-900/800/700）
   - `styleTokens.ts` 為 light mode 導向（slate-50/100/200）
   - 兩者未統一，導致不同頁面風格不一致

2. **元件樣式未統一**：
   - 按鈕：`button.tsx`（shadcn/ui）、`styleTokens.ts`（Gateway）、直接 Tailwind class（SmartDashboard）
   - 卡片：各頁面手寫 Tailwind class，未統一
   - Badge：未統一元件，各處手寫

3. **表單規範缺失**：
   - `placeholder` 不一致（有些有，有些無）
   - `autocomplete` 幾乎未使用
   - `inputmode` 未使用
   - `name` 屬性命名不一致

4. **Microcopy 未規範化**：
   - 載入狀態混用 `...` 與 `…`
   - 按鈕文字不一致（「儲存」vs「儲存中」）
   - 錯誤訊息語氣不一致

#### Phase 2 產出

**已建立**：
- `frontend/src/styles/designSystem.ts`：統一設計系統規範
  - 整合 `tokens.ts` 作為單一來源
  - 提供元件樣式類別（button、card、badge、sectionHeader）
  - 定義表單規範（input、label、error、autocomplete、inputmode、name）
  - 定義 microcopy 規範（loading、button、feedback、ellipsis）

**後續實作方向**：
- 逐步將現有元件遷移到使用 `designSystem.components.*`
- 補齊表單元件的 `autocomplete`、`inputmode`、`name` 屬性
- 統一 microcopy 使用 `designSystem.microcopy.*`
- 考慮建立統一的 FormInput、FormLabel、FormError 元件

### Phase 3：SmartDashboard 核心流程重整實作

#### 已完成的改進

1. **資料流向說明**（`SmartDashboardCommitPanel.tsx`）：
   - 在提交按鈕上方加入資料流向說明區塊
   - 明確說明資料會流向：
     - **資料庫**：點位與 Tag 資料寫入資料庫，供其他 UI 專案透過 API 取用
     - **本地 Modbus**：映射設定寫入本地 Modbus Server，供外部設備讀取
   - 使用藍色邊框與背景，與提交按鈕視覺一致

2. **流程指引**（`SmartDashboardSidebar.tsx`）：
   - 在 Tab Bar 上方加入流程指引區塊（僅在已選擇設備時顯示）
   - 顯示主流程：規劃 → Tag → Modbus → 提交
   - 當前步驟以藍色高亮顯示
   - 使用藍色邊框與背景，與 Tab 樣式一致

3. **WorkflowModal 收斂**（`SmartDashboardWorkflowModal.tsx`）：
   - 收斂為設備入口，標題改為「設備管理中心」
   - 非設備 intent 顯示簡化訊息：
     - 說明該功能已整合至主工作流程
     - 提供具體指引（例如：點位管理請使用「規劃」分頁）
     - 提供「前往設備管理中心」按鈕
   - 保留設備管理完整功能（搜尋、篩選、建立、編輯、測試、啟用/停用、刪除）

#### 改進效果

- ✅ 使用者清楚知道資料會流向資料庫與本地 Modbus
- ✅ 主流程（規劃 → Tag → Modbus → 提交）更加明確
- ✅ WorkflowModal 責任更清晰，不再承擔過多 legacy intent
- ✅ 非設備 intent 有明確的引導，不會讓使用者困惑

#### Phase 3 完成項目

- [x] 將 `WorkspaceContent` 聚焦在「資料來源設定 + 格子可視化」：
  - Flow Status Section 改為可收折，預設收合
  - 保留 Source Planner 緊湊列與 Memory Grid 作為核心功能
  - 使用者可依需要展開 Flow Status 查看詳細診斷資訊

#### Phase 4 完成項目

- [x] `LocalModbusWorkbenchPage` 重整區塊層級：
  - 系統狀態區塊：Server 狀態、映射數量、衝突狀態，含 Server 控制按鈕
  - Mapping 編輯區塊：Tag 選擇、Register 輸入、映射列表
  - 衝突治理區塊：顯示衝突列表與解決指引
  - 寫入測試區塊：Tag 選擇、測試數值輸入、執行測試寫入
- [x] 改善標題與說明文字，使其更符合單人工作流程
- [x] 更新測試以匹配新 UI（4 tests passed）

#### Phase 5：可近用性與響應式驗證

**Accessibility 補強**：
- [x] 在 `SmartDashboardCommitPanel` 的 `commitActionMessage` 加入 `aria-live="polite"` 與 `role="status"`
- [x] 在 `LocalModbusWorkbenchPage` 的狀態訊息加入 `aria-live="polite"` 與 `role="status"`
- [x] 在 `SmartDashboardSidebar` 的流程指引加入 `role="region"`、`aria-label` 與 `aria-current="step"`
- [x] 在 `SmartDashboardCommitPanel` 的資料流向說明加入 `role="region"` 與 `aria-label`
- [x] 裝飾性圖示加入 `aria-hidden="true"`

**響應式設計檢查結果**：

| 元件 | 窄螢幕處理 | 狀態 |
|------|-----------|------|
| **SmartDashboardWorkspaceSection** | `grid-cols-1 xl:grid-cols-[1fr_380px]` | ✅ 側欄在 xl 以下堆疊 |
| **SmartDashboardSidebar** | `flex-wrap` 用於 Mini Toolbar | ✅ Toolbar 按鈕可換行 |
| **SmartDashboardWorkspaceContent** | Source Planner 使用 `flex-wrap`，`px-2 sm:px-4` | ✅ 窄螢幕減少 padding |
| **SmartDashboardWorkflowModal** | `grid-cols-1 xl:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]` | ✅ 設備列表與設定在 xl 以下堆疊 |
| **LocalModbusWorkbenchPage** | `grid-cols-1 lg:grid-cols-2 xl:grid-cols-[1fr_400px]` | ✅ Mapping 編輯與衝突/測試在 lg 以下堆疊 |
| **Flow Status Section** | `grid-cols-1 sm:grid-cols-2 xl:grid-cols-4` | ✅ 流程區段響應式排列 |

**建議**：
- 窄螢幕下側欄可能需要改為 drawer/modal 模式（目前為堆疊）
- Source Planner 緊湊列在極窄螢幕下可能需要進一步優化（目前使用 flex-wrap）

**UI Regression Test**：
- [x] 建立 `tests/integration/ui/smart-dashboard-regression.test.tsx`
- [x] 驗證主要 UI 結構渲染正常

#### 後續優化建議（Phase 4+）

- [ ] 考慮加入更多視覺提示，引導使用者完成主流程（例如：完成某步驟後顯示下一步提示）
- [ ] 優化 Sidebar Tab 切換的動畫與過渡效果
- [ ] 考慮在 WorkspaceContent 加入「快速開始」引導（首次使用時）
- [ ] 窄螢幕下側欄可考慮改為 drawer/modal 模式

## 2026-03-16：Workbench implementation 關鍵發現

### Step 3 / TagBindingStudio
- `TagBindingStudio` 應沿用既有 `tagAPI.create` / `mappingAPI.create` 能力，但 UI 必須整合進 workbench，不能再回到舊 SmartDashboard/獨立 panel 心智模型。
- `tagBindingModel.ts` 必須把三種狀態拆清楚：
  - preview key 是否與既有 Tag 重複
  - preview key 是否在本批次內重複
  - point 是否已經有 mapping（already linked）
- batch bind 不能用 happy path 假設；必須逐筆建立 Tag 與 mapping，並把 partial failure 顯示回 UI，否則使用者會誤判整批已成功。
- `useEffect` 若直接依賴每次 render 都新建的 point id array，會造成 Step 3 selection 被反覆重設；改用 `pointIdsKey -> split` 後才能穩定。

### Step 4 / LocalModbusBoard
- `LocalModbusWorkbenchPage` 的 domain 能力可重用，但呈現層必須重做為 workbench 內嵌 board。
- output candidates 不能只看 Tag 或只看 Modbus mapping；必須以「selected device 的 points + mappings + linked tags」交集推導，否則會把別台設備的輸出候選混進來。
- Local Modbus sync 需明確以 duplicate register conflict 作為 block 條件；若 register 衝突未阻擋 sync，使用者會把錯誤映射推到 server。
- `loadData` 若直接依賴 `t`，在測試 mock `useTranslation()` 時可能造成 callback identity 改變，進而重複觸發 effect；以 `ref` 穩定 fallback translation 後可避免這個問題。
- register input 若在 `selectedCandidate` object identity 改變時每次都重設，會蓋掉使用者剛輸入的值；依賴應收斂到 `selectedTagId` 與 `selectedCandidateRegister`。
- `RegisterMapCanvas`、auto-map、dry-run helper 若只存在於檔案上半但沒有接進主 render，對使用者來說等同功能不存在；這類「已實作 helper / 未接線 UI」是 Step 4 的主要落差來源。
- Local Modbus register 衝突不能只比對 start register；multi-word data type（如 `int32` / `float64`）必須以實際占用 slot 範圍判斷 overlap，否則會漏掉 `HR10~11` 與 `HR11` 這類交疊衝突。
- Database schema snapshot 不應在 metadata 尚未回來時先 render 空殼 testid；否則測試與使用者都會在「面板已出現但內容尚未掛上」的中間態讀到錯誤結論。
- Step 4 inspector 的 readiness 不能只看 active target；正確來源是 `outputCandidate` 對應 tag 在 Local Modbus 與 Database 兩邊的 mapping 完整度。
- Step 2 inspector 不需要再從 `SourceCanvasSection` 手動傳 detail props；既然 `rules / selectedRuleId / selectedAddress` 已在 provider，直接在 shared inspector 以 provider state + queries 重建 detail model，比跨元件穿線更穩定。
- Step 3 若要在 shared inspector 看單筆 detail，candidate row 本身就必須成為 selection source；只靠 checkbox 無法表達「我想看哪一筆」。

### Shell UI / Summary
- `WorkbenchHeaderBar` 與 `WorkbenchActionDock` 若只顯示 active step / selected device，不足以支撐單頁主流程；需要補 point/tag/output counts 與 next action，使用者才知道目前流程停在哪裡。
- `useWorkbenchSummary` 適合作為 shell 層的單一摘要來源，讓 header/action dock 不必各自重複查詢與計算。

### 驗證策略
- 某些 workbench Vitest 組合在同一個 command 下，會出現「測試邏輯已通過，但 worker 不正常結束」的情況；分批跑：
  - unit/meta tests
  - shell/foundation
  - source/tag
  - output/local-modbus
  可以穩定完成驗證並保留定位能力。

## 2026-03-16：Workbench 收尾與 phase2 接手

### 收尾結論
- 手動 integration review 後，`useWorkbenchSummary`、`WorkbenchHeaderBar`、`WorkbenchActionDock`、`TagBindingStudio`、`LocalModbusBoard` 沒有再發現 correctness blocker；目前 `/datalink/workbench` 已可完整承接 `來源可視化 -> Tag -> Local Modbus` 主線。
- workbench 這段的可靠 quality pass 做法，是把測試拆成三批再跑 lint/build/diff check；直接把多個 workbench 測試檔一次串成超大批次，容易碰到 Vitest worker 不正常結束。

### 代理執行策略修正
- 本 session 多個 background agents（code-review、general-purpose、explore）都出現兩種問題：
  - 長時間 `running` 但沒有第一輪 turn
  - `claude-opus-4.6` code review 回覆 `429 rate limit`
- 因此後續 phase2 若再遇到相同狀況，不應空等；應優先改為主代理手動執行，或改派更小範圍的 sync 子任務。

### phase2 接手判斷
- `runtime-live-value-phase2` 與 `database-target-phase2` 在 workbench quality pass 完成後已成為 ready todo。
- 依目前 spec 與先前掃描結果，runtime / db target 很可能不是「把現有 UI 接上」即可，而是要從 backend contract、frontend service/hook、頁面區塊三個面向補一個新的 vertical slice。

### runtime/live value vertical slice（本輪完成）
- `runtime.Service` 現在除了寫 timeseries，也會廣播 device-scoped `value` 事件；`/api/v1/datalink/runtime/stream` 已可接收 `device_id` + `point_ids` 過濾並輸出 `value` / `heartbeat` SSE。
- `cmd/test_ui/main.go` 已接上 `collector.Scheduler`、`runtime.Service` 與 batch writer；實際 smoke 測到 `/api/v1/datalink/runtime/status` 回 `{"running":true,...}`，代表 app wiring 生效。
- workbench source step 現在會用 `runtime/status` 顯示 collector summary，並透過 `useRuntimeStream` 將 source cell / ledger 更新為 point raw live value。
- review 補強後，point create/update/delete 會同步 runtime scheduler/meta，mapping create/update/delete 會 refresh runtime mappings，避免 runtime 只吃啟動時快照。
- source 視圖必須顯示 point raw value，不應吃 mapping `transformed_value`；tag/output 相關面板才適合顯示 transform 後的值。

## 2026-03-16：Workbench desktop redesign round 2 關鍵發現
- 1920×1080 下的主要問題不是單一 CSS 細修，而是 shell IA drift：主工作區被 `max-w` 內容框與常駐 summary dock 共同壓縮。
- Step 1 不能在 desktop main column 內再自帶第二個 inline inspector，否則即使 viewport 很寬，實際可用內容仍過窄。
- Step 1 的 capability summary 目前沒有後端專用欄位可直接讀取 `address base / word order`；安全做法是用 `connection_config` + protocol defaults/traits 派生，缺值時顯示 `protocol default / not specified / n/a`，而不是虛構 schema。
- Step 1 的「最近三次連線測試」可先以前端 session-local history 補齊：沿用後端 `last_test_*` 當 fallback latest snapshot，並在 workbench 內每次 `Test connection` 後將結果推入最近三筆時間線。
- Step 2 的價值不只是地址規畫，還包含即時值監看與規則模板；因此 `AddressCanvas` 需要穩定幾何 + 可切換資訊層，而不是切成多套獨立頁面。
- Step 2 的 lattice 不能只畫「有點位或有規劃的格子」；必須從最小位址畫到最大位址，把 gap 保留下來，操作者才看得出空洞與衝突位置。
- Step 2 的 `Plan / Live / Link` 不應切成三套不同 layout；正確做法是固定同一套幾何，只替換每格要顯示的資訊層。
- `sourceCanvasModel` 需要同時兼容新 `rules[]` 模型與舊 `plannedPointAddresses` 測試輸入，才能在重構期間保住既有單元測試 contract。
- Step 2 的 source rules 不能只放在 `SourceCanvasSection` local state；desktop flow 會頻繁做 `Source → Tag → Output → Source` 往返，若切步驟就清空，整個規畫體驗會斷掉。至少 `rules / selectedRuleId / selectedAddress` 要提升到 provider 或其他持久層。
- 當 source rules 被提升為跨步驟 state 後，rule id 不能再依賴 component-local ref；重新掛載後若 ref 歸零，下一次 apply 會重複產生 `rule-1`，造成 key/selection 混亂。安全做法是由現存 rule id 推導下一個序號。
- Step 3 的 candidate row 若只顯示 preview key，使用者仍不知道「這筆綁定代表哪個來源值」；至少要把 address/span/raw value/bit width/status 與 preview key 放在同一列。
- Step 3 的 `create new tag` 與 `bind existing tag` 必須明確分流；若只留一個 batch bind 按鈕，使用者無法理解自己是在建立新 Tag 還是沿用舊 Tag。
- Step 3 使用者更在意 tag 管理資訊密度，而不只是批次操作按鈕；主區必須直接露出 naming/source/value/merge/status 等資訊。
- Step 4 若只保留表格式設定，無法達到「能理解、能驗證、能維運」；Local Modbus 與 Database 都需要各自的可視化與驗證面。
- Step 4 的第一層不應直接掉進單一 target 的詳細表單；正確做法是先有 shared candidate board，讓使用者一眼看到每個 tag 在 Modbus / Database 兩邊的綁定狀態，再切換 active studio 做細部設定。
- shared board 上的 Modbus 狀態不適合再用長句 microcopy；桌面情境下直接顯示 `HR12` 這種短 badge，比「Mapped to HR12」更容易掃描與比對。
- output step 測試若會切到 `DatabaseTargetBoard`，必須等待 database mapping table 掛載完成，否則 async load 會在測試結束後補跑 setState，留下 React `act(...)` 警告。
- shell state 需要明確拆成兩種脈絡：
  - `inspectorSelection`：只代表目前步驟的局部選取，切步驟時必須清掉
  - `crossStepContext`：保存 Step 2 → Step 3 → Step 4 的焦點 rule/tag 線索，只在切換 device 時重置
- `useWorkbenchSummary` 已成為 shell 層的單一 readiness 匯總來源；StepRail 與 BottomSummaryBar 應共用這份導出，避免各自重算。
- frontend Vitest 的正式入口仍以 `frontend/tests/unit/**` wrapper 為主；若要做 targeted validation，應優先跑 wrapper 檔而不是直接把 `src/pages/.../__tests__` 路徑丟給 `npm run test`.


### Execution plan 關鍵決策
- 拆為 20 個 todos，以最大化平行度：Phase 0 完成後可同時展開 4 條 track（Device/Source/Tag/Output）。
- 每條 track 內部有序列依賴，但 track 之間完全獨立，適合 fleet dispatch。
- 現有 workbench 檔案（WorkbenchDeviceStep、SourceCanvasSection、TagBindingStudio、LocalModbusBoard 等）視為「可參考但需重做」，不直接 patch。
- workbenchDeviceFormModel.ts、sourceCanvasModel.ts、tagBindingModel.ts 的計算邏輯可重用，但 UI composition 需重寫。
- i18n 鍵在各 step 實作時同步擴充，Phase 2 做最終稽核。
- 不在此輪刪除 legacy routes（SmartDashboard、LocalModbusWorkbenchPage），僅決策導流策略。

### 2026-03-16：Phase 2 / spec alignment 再盤查
- OpenSpec `tasks.md` 目前仍有 `3.2`、`3.3`、`6.1`、`6.2`、`6.3` 未勾；這表示目前剩餘工作不只是 quality/rollout，Step 2 視覺幾何也需要再對齊 spec。
- `SourceCanvasSection.tsx` 的 toolbar 與 state 基本已具備 `Plan / Live / Link`、value format、freeze live / snapshot compare、jump-to-address、coverage overview；因此 `3.3` 比較像是「已大致實作，但 OpenSpec 勾選尚未同步」。
- 真正的 spec drift 在 `AddressCanvas.tsx`：目前主畫面仍是 `sm:grid-cols-2 xl:grid-cols-5 2xl:grid-cols-6` 的卡片格，而不是批准規格要求的固定 16-bit 連續 lattice。換句話說，item model 雖已會補 gap，但呈現幾何尚未到位。
- 目前 `DatalinkWorkbenchSourceStep.test.tsx` 只驗證 `address-cell-*` 的狀態與 view mode 切換，並沒有驗證 row geometry / 固定欄數 / 連續 lattice 版面，因此這個 drift 先前能通過測試。
- `/datalink/local-modbus` 目前在 `App.tsx` 仍直接掛到 `LocalModbusWorkbenchPage`，`SmartDashboardPage.tsx` 也仍用 `navigate('/datalink/local-modbus?...')` 導向舊頁；這代表 `6.3` 的 rollout / compatibility strategy 尚未真正落地。
- rollout spec 的核心原則仍應維持：legacy route 保留作 fallback，但要把「Try new workbench」與新舊入口關係說清楚，避免把 legacy deep link 直接混入尚未完成的新 state model。
- 本輪 `agent-77` 的 spec audit 再次因 `429` 失敗，說明 fleet 在 Phase 2 收尾仍不能只靠 Opus background agents；controller 必須持續手動 spot-check 並修正 planning artifacts。

### 2026-03-17：Phase 2 已收斂
- `AddressCanvas.tsx` 已改為固定 16-bit lattice；先前 `3.2` drift 已由 row geometry / lattice columns / merged span regression tests 補齊。
- `redesign-i18n-a11y` 的高風險缺口集中在：
  - source cell 缺 `aria-pressed`
  - device notice / output inspector error 缺 `aria-live`
  - Step 1 protocol-specific option label 仍有硬編碼英文
  以上現已補齊並有對應測試。
- `WorkbenchFrame` 加上 `overflow-hidden` 與 row-2 `min-h-0` 後，桌面殼層更符合 1920×1080 的固定區域滾動模型。
- `/datalink/local-modbus` 現在改為 compat redirect 到新 workbench output；舊 `LocalModbusWorkbenchPage` 保留在 `/datalink/local-modbus/legacy` 作 fallback。
- `SmartDashboardPage` 後段 CTA 已直接導向 `/datalink/workbench?step=output&target=modbus`，新舊入口關係已收斂到單一路徑。
