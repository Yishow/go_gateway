# 進度日誌

## 會話：2026-05-29

### 階段 1：建立分支與規劃
- **狀態：** complete
- **開始時間：** 03:00
- 執行的操作：
  - 檢查 git status 與 spectra list 狀態
  - 建立 `task_plan.md`、`findings.md`、`progress.md`
- 建立/修改的檔案：
  - `task_plan.md`
  - `findings.md`
  - `progress.md`

### 階段 2：實作 datalink-workbench-v2-shell
- **狀態：** complete
- **開始時間：** 03:01
- 執行的操作：
  - 安裝 `@fontsource/inter` 與 `@fontsource/jetbrains-mono` 字型套件
  - 建立 `tokens.ts` 與 Scoped `workbench-v2.css`
  - 建立共用 UI 元件集 (Icon, Button, inputs, Field, Toggle, StatusChip, SectionCard)
  - 建立狀態管理 Hook `useWorkbenchV2State` (載入預設裝置與規則)
  - 建立佈局 Shell 各區塊 (TopBar, StepRail, SummaryRail, TweaksPanel, WorkbenchV2Shell)
  - 註冊 `/studio/v2` 路由，並確保舊版路由完全不受影響
  - 撰寫單元測試及路由測試，確保全部通過且 bundle 建置正常
- 建立/修改的檔案：
  - `frontend/src/main.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/tokens.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css` (建立)
  - `frontend/src/features/datalink/workbench-v2/components/*` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/types.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/shell/*` (建立)
  - `frontend/src/features/datalink/workbench-v2/steps/*` (建立)
  - `frontend/src/features/datalink/workbench-v2/settings/*` (建立)
  - `frontend/src/i18n/locales/zh-TW/workbench-v2.json` (建立)
  - `frontend/src/i18n/locales/en/workbench-v2.json` (建立)
  - `frontend/src/i18n/config.ts` (修改)
  - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx` (建立)
  - `frontend/src/App.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/tokens.test.ts` (建立)
  - `frontend/tests/unit/workbench-v2/components.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/shell.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/routing.test.tsx` (建立)

### 階段 3：實作 datalink-workbench-v2-step1-device
- **狀態：** complete
- **開始時間：** 08:15
- 執行的操作：
  - 補完 `state/types.ts` 中的 `ProtocolId`, `ReadinessStage`, `DeviceTest` 型別與 tests/d-level。
  - 建立 `state/protocols.ts` 和 `state/deviceColors.ts` 及其對應之 helper 函數與 hooks。
  - 在 `state/useWorkbenchV2State.ts` 擴充 10 個 device-related reducer actions 與級聯刪除函式 `cascadeRemoveDevice`。
  - 建立 `ProtocolSelector`, `ConnectionConfigForm`, `ReadinessStages`, `ConnectionTestPanel`, `DeviceTabRail`, `DeviceEditor` 元件。
  - 組裝並嵌入 `Step1Device` 到 `WorkbenchV2Shell.tsx` 以替換佔位元件，擴充 zh-TW / en i18n key-value。
- 建立/修改的檔案：
  - `frontend/src/features/datalink/workbench-v2/state/types.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/protocols.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/deviceColors.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step1/` (建立元件目錄)
  - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx` (修改)
  - `frontend/src/i18n/locales/zh-TW/workbench-v2.json` (修改)
  - `frontend/src/i18n/locales/en/workbench-v2.json` (修改)
  - `frontend/tests/unit/workbench-v2/step1.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/step1-readiness.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/deviceColors.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/protocols.test.ts` (建立)
  - `frontend/tests/unit/workbench-v2/types.test-d.ts` (建立)

### 階段 4：實作 datalink-workbench-v2-step2-rule
- **狀態：** complete
- **開始時間：** 08:30
- 執行的操作：
  - 建立 `ruleReducer.ts` 將 `useWorkbenchV2State.ts` 內的 11 個 rule cases 移出委派，以滿足單檔行數限制目標。
  - 將 action `WorkbenchV2Action` 強型別化，清除 type assertion 的 any warning。
  - 完成格狀畫布中 active source rule 顯示與 rules 清單。
  - 新增並通過 `reducer-step2.test.ts` 與 `sourceRule.test.ts`。

### 階段 5：實作 datalink-workbench-v2-step3-mapping
- **狀態：** complete
- **開始時間：** 08:45
- 執行的操作：
  - 新增 `Step3Mapping` 組件，整合 Tag/Mapping 列表與 review-first 主線。
  - 修復 `MappingTable.tsx` 元件中 `selectedPoint` 為 null 的 type error。
  - 移除 ESLint 未使用的 variables 警告。
  - 新增並通過 `reducer-step3.test.ts` 與 `step3-mapping.test.tsx`。

### 階段 6：實作 datalink-workbench-v2-step4-database
- **狀態：** complete
- **開始時間：** 09:00
- 執行的操作：
  - 實作資料庫與 Local Modbus 的目標輸出綁定組件。
  - 實作 autoAssignTargets 自動分配功能。
  - 建立 CommitSummary、CommitProgress、CommitSuccessCard 提交介面，並實作 10 步驟動畫。
  - 修正測試 `reducer-step1.test.ts` 的 target_type 'double' 改為 'float64' 以對齊強型別限制。
  - 新增並通過 `reducer-step4.test.ts`、`step4-database.test.tsx`、`step4-commit.test.tsx`、`dbSchemas.test.ts`。

### 階段 7：實作 datalink-workbench-v2-settings
- **狀態：** complete
- **開始時間：** 09:15
- 執行的操作：
  - 實作 `SettingsPlaceholder` 及完整的預設設定檔。
  - 補齊 settingsReducer 的型別與 logic。
  - 新增並通過 `reducer-settings.test.ts`。

### 階段 8：最終驗證與整合
- **狀態：** complete
- **開始時間：** 09:30
- 執行的操作：
  - 使用 `spectra archive` 完成所有 6 個 workbench-v2 變更的 Spectra 歸檔。
  - 執行全量 `tsc && vite build`，建置成功且無 error。
  - 執行 `npm run test -- --run workbench-v2` 通過 30 個測試檔案、201 個測試案例。

### 階段 9：優化項目補齊與日誌系統升級
- **狀態：** complete
- **開始時間：** 12:15
- 執行的操作：
  - 補齊並優化即時診斷日誌面板 (Real-time Log Panel) 與 `stale` / `recovered` 點位狀態轉換日誌。
  - 實作 pointStaleStatesRef 狀態過濾器，解決高頻 stale 點位推送造成的 `setLogs` 狀態狂刷與 DOM 渲染效能瓶頸。
  - 調整日誌輸出順序，變更為 append 方式插入日誌陣列末尾，以實現時間遞增 (Chronological) 顯示，並與日誌置底滾動 UX 完美契合。
  - 在 `runtime-dashboard-state.test.tsx` 新增針對 logs 轉移與去重的整合單元測試。
  - 執行 `npm run test -- --run runtime-dashboard` 通過 3 個測試檔案、13 個測試案例（全量綠燈）。
  - 執行後端 `go test ./...` 整合測試與前端全量 production build (tsc && vite build)，均完全通過。

### 階段 10：建立 studio surface 與 API 維護文件
- **狀態：** complete
- **開始時間：** 14:40
- 執行的操作：
  - 盤點 router 中所有與 `/studio`、`/studio/v2`、`/studio/runtime`、`/gateway/*`、`/test` 相關的產品 surface。
  - 盤點 `/studio` 主線各步驟實際使用的 query / mutation / runtime SSE，以及候選 review / output / runtime 相關 API。
  - 確認 `/studio/v2` 現況主要仍是 local reducer draft，Step 4 commit 為模擬流程，runtime handoff 只完成 route 層。
  - 確認 `/gateway/*` 當前主要依賴 feature flag 與 `/test/connect` connect-only APIs。
  - 新增 `docs/technical/studio-surface-inventory/` 下的多份 Markdown 文件與 HTML 總覽頁。
- 建立/修改的檔案：
  - `docs/technical/studio-surface-inventory/README.md` (建立)
  - `docs/technical/studio-surface-inventory/studio-mainline.md` (建立)
  - `docs/technical/studio-surface-inventory/studio-v2-runtime.md` (建立)
  - `docs/technical/studio-surface-inventory/gateway-and-test.md` (建立)
  - `docs/technical/studio-surface-inventory/backend-api-registry.md` (建立)
  - `docs/technical/studio-surface-inventory/gap-roadmap.md` (建立)
  - `docs/technical/studio-surface-inventory/index.html` (建立)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 11：依產品優先順序重構 surface inventory 文件
- **狀態：** complete
- **開始時間：** 15:05
- 執行的操作：
  - 依使用者新指示將 `/studio` 標示為完整版但先暫停。
  - 將 `/studio/v2` 與 `/studio/runtime` 文件改成預設入口與重點施作語境。
  - 將原本混寫的 gateway/test 文件拆成 `/test` 與 `/gateway/*` 兩份獨立文件。
  - 更新 HTML 總覽頁，把重點施作 / 暫停 / 僅記錄狀態做成明確視覺標示。
- 建立/修改的檔案：
  - `docs/technical/studio-surface-inventory/README.md` (修改)
  - `docs/technical/studio-surface-inventory/studio-mainline.md` (修改)
  - `docs/technical/studio-surface-inventory/studio-v2-runtime.md` (修改)
  - `docs/technical/studio-surface-inventory/test-tooling.md` (建立)
  - `docs/technical/studio-surface-inventory/gateway-experiments.md` (建立)
  - `docs/technical/studio-surface-inventory/backend-api-registry.md` (修改)
  - `docs/technical/studio-surface-inventory/gap-roadmap.md` (修改)
  - `docs/technical/studio-surface-inventory/index.html` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 12：修正 HTML 文件台無資料並拆分為 html/js/css
- **狀態：** complete
- **開始時間：** 15:20
- 執行的操作：
  - 檢查第一版 `index.html` 無資料的根因，確認為內嵌 JS 語法錯誤導致 render 未執行。
  - 將文件台拆為 `index.html`、`inventory.css`、`inventory.js` 三檔，避免大段 inline script/style 難維護。
  - 移除會破壞 template literal 的未跳脫反引號內容，改用 `<code>` 標記輸出。
  - 以 `node --check` 驗證 `inventory.js` 語法正確。
- 建立/修改的檔案：
  - `docs/technical/studio-surface-inventory/index.html` (重建)
  - `docs/technical/studio-surface-inventory/inventory.css` (建立)
  - `docs/technical/studio-surface-inventory/inventory.js` (建立)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 13：為 studio surface inventory 建立 changelog 機制
- **狀態：** complete
- **開始時間：** 15:45
- 執行的操作：
  - 檢查 `AGENTS.md`、`CLAUDE.md` 與現有 inventory 文件，確認目前沒有明文要求在更新 inventory 時同步留下 changelog。
  - 確認 `go.mod` 已含 `modernc.org/sqlite`，適合用最小 Go CLI 落地 SQLite 記錄，而不是只寫文字規範。
  - 規劃新增 `cmd/studio_inventory_changelog`，提供 `init`、`add`、`list` 三個最小命令。
  - 在 `AGENTS.md`、`CLAUDE.md` 與 inventory `README.md` 補上 changelog 維護規則與使用方式。
  - 以 `go run ./cmd/studio_inventory_changelog init` 初始化 `docs/technical/studio-surface-inventory/changelog.sqlite`。
  - 以 `go run ./cmd/studio_inventory_changelog add ...` 寫入首筆 changelog，記錄這次機制建立本身。
  - 以 `go run ./cmd/studio_inventory_changelog list -limit 5` 驗證資料可讀。
  - 執行 `go test ./cmd/studio_inventory_changelog` 與 `git diff --check` 完成收尾驗證。
- 建立/修改的檔案：
  - `AGENTS.md` (修改)
  - `CLAUDE.md` (修改)
  - `cmd/studio_inventory_changelog/main.go` (建立)
  - `docs/technical/studio-surface-inventory/README.md` (修改)
  - `docs/technical/studio-surface-inventory/changelog.sqlite` (建立)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 14：補齊 inventory onboarding 入口
- **狀態：** complete
- **開始時間：** 19:35
- 執行的操作：
  - 檢查現有 inventory 文件結構，確認目前雖有 README 與 HTML，但仍缺接手 AI 的固定第一入口與 machine-readable 摘要。
  - 規劃新增 `START_HERE.md` 與 `context.json`，並要求 `AGENTS.md` / `CLAUDE.md` 明確導向這兩個入口。
  - 建立 `START_HERE.md`，把產品決策、閱讀順序、常見任務入口與 changelog 用法濃縮成單一入口。
  - 建立 `context.json`，將重點 surface、策略、canonical docs、known gaps 與 changelog 路徑轉成 machine-readable 摘要。
  - 更新 `README.md`、`AGENTS.md`、`CLAUDE.md`，要求接手 inventory 任務時先讀 `START_HERE.md` 與 `context.json`。
  - 以 `go run ./cmd/studio_inventory_changelog add ...` 寫入這次 onboarding 更新的 changelog。
  - 執行 `go run ./cmd/studio_inventory_changelog list -limit 5`、`python3 -m json.tool docs/technical/studio-surface-inventory/context.json`、`git diff --check` 驗證。
- 建立/修改的檔案：
  - `docs/technical/studio-surface-inventory/START_HERE.md` (建立)
  - `docs/technical/studio-surface-inventory/context.json` (建立)
  - `docs/technical/studio-surface-inventory/README.md` (修改)
  - `AGENTS.md` (修改)

### 階段 29：Step 3 顯示真實裝置值並以真值計算 preview
- **狀態：** complete
- **開始時間：** 00:45
- 執行的操作：
  - 補 `selectors.test.tsx`、`step3-live-preview.test.tsx`、`step3-mapping.test.tsx`，並新增 `step3-live-values.test.tsx`，先鎖住「不再吃 raw seed」與「Point → Tag 表格顯示裝置即時值」。
  - 新增 `useStep3LiveValues.ts`，以每台 device 一條 runtime SSE stream 聚合 `liveValues` 與 `connectionByDevice`。
  - 移除 `useSelectedMapping()` 的 deterministic `rawSeed` 回傳，讓 `TransformPreview` 與 `PipelineSteps` 改吃真實 `rawValue`。
  - 在 `MappingTable` / `MappingRow` 新增「裝置值」欄，沒值時顯示 `尚未收到 / 連線中 / 串流錯誤`。
  - 在 `TransformPreview` 新增 waiting state，沒收到 live value 時不送 preview request，也不顯示假 preview。
  - 補用 `persisted_point_id` 訂閱 runtime stream 的 RED/GREEN 測試，修正 `useStep3LiveValues` 不該拿本地 row id (`p-*`) 訂閱 SSE 的根因。
  - 執行 Step 3 targeted Vitest、`npm run build` 與 `git diff --check`，全部通過。
- 建立/修改的檔案：
  - `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/selectors.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx` (修改)
  - `frontend/src/i18n/locales/zh-TW/workbench-v2.json` (修改)
  - `frontend/src/i18n/locales/en/workbench-v2.json` (修改)
  - `frontend/tests/unit/workbench-v2/selectors.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/step3-live-values.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/step3-mapping.test.tsx` (修改)

### 階段 30：修正 Step 3 runtime stream 重連風暴
- **狀態：** complete
- **開始時間：** 01:20
- 執行的操作：
  - 補 `step3-live-values.test.tsx` RED 測試，證明 mapping 僅變更 `save_state` 時不應重新建立 runtime stream。
  - 將 `useStep3LiveValues.ts` 改為先建立 `Step3LiveSubscription`，並以 `local point id / device_id / address / persisted_point_id` 組成穩定 `subscriptionKey`。
  - 讓 SSE effect 不再依賴整個 `mappings` 物件或每次重建的 lookup object，避免 autosave 期間持續 close/reopen `EventSource`。
  - 執行 targeted Vitest、`npm run build` 與 `git diff --check`，全部通過。
- 建立/修改的檔案：
  - `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts` (修改)
  - `frontend/tests/unit/workbench-v2/step3-live-values.test.tsx` (修改)

### 階段 31：修正 Step 3 數值顯示只吃 SSE、未退回 last_value
- **狀態：** complete
- **開始時間：** 01:32
- 執行的操作：
  - 補 `step3-live-values.test.tsx` RED 測試，要求沒有 SSE 時也能顯示 persisted point `last_value`。
  - 在 `useStep3LiveValues.ts` 加入 points snapshot 查詢，依 device 讀 `/points`，把 persisted point `last_value` 映回本地 row。
  - 讓 Step 3 顯示值改成 `snapshotRawValues + SSE override`，因此左表與右側 preview 都能先顯示最後一次設備值。
  - 更新 `step3-mapping.test.tsx` 的 QueryClientProvider / pointAPI mock，以支撐新的 query 路徑。
  - 執行 targeted Vitest、`npm run build` 與 `git diff --check`，全部通過。
- 建立/修改的檔案：
  - `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/step3-live-values.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/step3-mapping.test.tsx` (修改)

### 階段 25：定義 Studio V2 Step 3 live preview 與 target type 快捷設定
- **狀態：** complete
- **開始時間：** 22:10
- 執行的操作：
  - 檢查 Step 3 spec、前端元件與既有 mapping preview API，確認 `轉換管線預覽` 仍是本地 mock，而 backend 已有 `POST /api/v1/datalink/mappings/preview` 可重用。
  - 與使用者逐段確認設計：重用既有 preview API、payload preview 保持 display-only、`目標型態` 同時支援當前列與全部列快捷操作。
  - 新增 `docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md`，正式記錄 approved design。
  - 建立 Spectra change `wire-studio-v2-step3-live-preview-and-target-type-shortcuts`，補齊 `proposal.md`、`design.md`、`specs/` 與 `tasks.md`。
- 建立/修改的檔案：
  - `docs/superpowers/specs/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts-design.md` (建立)
  - `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/proposal.md` (建立)
  - `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/design.md` (建立)
  - `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/specs/datalink-workbench-v2-step3-mapping/spec.md` (建立)
  - `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/specs/datalink-api/spec.md` (建立)
  - `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/tasks.md` (建立)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 26：為 Studio V2 Step 3 live preview 撰寫 implementation plan
- **狀態：** complete
- **開始時間：** 22:25
- 執行的操作：
  - 讀取 approved design doc 與新建 Spectra change artifacts，整理成可直接執行的 implementation plan。
  - 確認 Step 3 async preview coverage 不再擴寫既有大型測試檔，而是拆成獨立 targeted tests。
  - 建立 `docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md`，明確列出 task、驗證命令與涉及檔案。
- 建立/修改的檔案：
  - `docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md` (建立)
  - `task_plan.md` (修改)
  - `progress.md` (修改)

### 階段 27：實作 Studio V2 Step 3 live preview 與 target type 快捷設定
- **狀態：** complete
- **開始時間：** 22:35
- 執行的操作：
  - 新增 `step3-live-preview.test.tsx`，覆蓋 loading、error、stale response guard 與無選取列不發 request。
  - 新增 `step3-target-type-shortcuts.test.tsx`，覆蓋當前列快捷設定與 `套用到全部列`。
  - 建立 `useStep3LivePreview.ts`，以 debounce + request id guard 重用 `mappingAPI.preview(...)`。
  - 將 `TransformPreview`、`PipelineSteps`、`Step3Mapping` 接到真實 preview 與 `target_type` 快捷操作。
  - 補齊中英文 i18n key，並修正 build 時暴露的 preview scalar / callback 型別問題。
  - 執行 `npm run test -- --run tests/unit/workbench-v2/step3-live-preview.test.tsx tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx tests/unit/workbench-v2/step3-mapping.test.tsx`、`npm run build`、`git diff --check`，全數通過。
- 建立/修改的檔案：
  - `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LivePreview.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/TargetTypeQuickActions.tsx` (建立)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/index.ts` (修改)
  - `frontend/src/i18n/locales/zh-TW/workbench-v2.json` (修改)
  - `frontend/src/i18n/locales/en/workbench-v2.json` (修改)
  - `frontend/tests/unit/workbench-v2/step3-live-preview.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/step3-target-type-shortcuts.test.tsx` (建立)
  - `openspec/changes/wire-studio-v2-step3-live-preview-and-target-type-shortcuts/tasks.md` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 28：修正 Studio V2 rule autosave 500 與 mapping stale-row 連鎖錯誤
- **狀態：** complete
- **開始時間：** 00:10
- 執行的操作：
  - 依照 user 提供的 HTTP log 與 UI 錯誤訊息，先針對 `/studio/v2` source-rule autosave 與 Step 3 mapping autosave 做 root-cause tracing。
  - 補 `mapping-autosave-page.test.tsx` RED test，重現 point row address 改變後仍沿用舊 `mapping_id`，導致後續 save-error。
  - 補 `router_studio_v2_workspace_source_rules_test.go` RED test，重現 device 尚未 readiness 時 `PUT /studio-v2/workspace/source-rules/:id` 被當成 500。
  - 將 `ruleActivationBlockedError` 改為 `source rule validation failed` 類型，讓 handler 正確回 400。
  - 在 `mappingReducer.ts` 加入 point row identity guard，當 `rule_id/device_id/address` 變更時重建預設 mapping，避免 stale `mapping_id`、`persisted_value` 留在新 row。
  - 在 `useStudioV2MappingAutosave.ts` 加入 owning rule gate：當來源規則尚未 persisted/saved 時，不再對 backend 發 mapping autosave。
  - 執行 `go test ./internal/api -run 'TestNewRouter_StudioV2WorkspaceSourceRule(UpdateReturnsValidationWhenDeviceIsNotReady|AutosaveEndpoints)$'`、`npm run test -- --run tests/unit/workbench-v2/mapping-autosave-page.test.tsx tests/unit/workbench-v2/reducer-step3.test.ts tests/unit/workbench-v2/step3-mapping.test.tsx`、`npm run build`、`git diff --check`，全部通過。
- 建立/修改的檔案：
  - `internal/datalink/sourcerule/activation_readiness.go` (修改)
  - `internal/api/router_studio_v2_workspace_source_rules_test.go` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts` (修改)
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts` (修改)
  - `frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 26：為 Studio V2 Step 3 live preview 撰寫 implementation plan
- **狀態：** complete
- **開始時間：** 22:35
- 執行的操作：
  - 讀取 `writing-plans` skill 與既有 plan 範例，依 repo 慣例建立新的 implementation plan。
  - 盤點 Step 3 檔案大小與責任邊界，確認 `step3-mapping.test.tsx` 已 375 行，因此新增 coverage 應拆成獨立測試檔，而不是繼續膨脹既有檔案。
  - 規劃新 hook `useStep3LivePreview.ts` 與新元件 `TargetTypeQuickActions.tsx`，讓 async preview 邏輯與快捷操作各自聚焦。
  - 完成 `docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md`，將實作切成 4 個可執行 task，附具體檔案、測試、指令與 commit 邊界。
- 建立/修改的檔案：
  - `docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md` (建立)
  - `task_plan.md` (修改)
  - `progress.md` (修改)

### 階段 22：修正 Studio V2 空 workspace seeded fallback
- **狀態：** complete
- **開始時間：** 20:45
- 執行的操作：
  - 檢查 `useStudioV2AutosaveState` 與 `workspace-boot` 測試，確認空 workspace 時 hydration 直接 return，導致 `INITIAL_STATE` 的預設 device/rule 被保留。
  - 將空 workspace hydration 改為明確 dispatch 空 `devices` / `rules`，並新增 `workspaceHydrated` gate 讓 page 在 hydration 完成前維持 bootstrap loading。
  - 補齊 `workspace-boot.test.tsx`，驗證空 workspace 進入時 shell 收到 `device-count=0`、`rule-count=0`。
  - 執行 `cd frontend && npm run test -- --run tests/unit/workbench-v2/workspace-boot.test.tsx`，3 個測試全數通過。
- 建立/修改的檔案：
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts` (修改)
  - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 19：依 `docs/goal.md` 串上 Studio V2 Step 2 rule autosave
- **狀態：** complete
- **開始時間：** 14:10
- 執行的操作：
  - 先前已補齊 backend red tests 與 workspace-scoped source rule CRUD / ownership handler；本輪完成前端 Step 2 rule autosave orchestration。
  - 建立 `useStudioV2RuleAutosave.ts`，把 per-rule save queue、valid-only save、delete interception 與 per-rule save-state 隔離從 page 入口統一管理。
  - 在 `useStudioV2AutosaveState.ts` 追加 workspace rule query / hydration，讓 `/studio/v2` 重新載入時可還原 persisted rule 的 `workspace_id`、`device_id`、`revision_id`，並在只有 persisted devices、尚無 persisted rules 時把 draft rule rebind 到第一台 workspace device。
  - 保持 scope 最小：只 autosave backend 已有 source-rule 契約欄位；`name` 與 `share_*` 仍維持 local-only，不假裝已持久化。
  - 更新 `DatalinkWorkbenchV2Page.tsx` 的 bootstrap loading / error gate，讓 device + rule bootstrap 都完成後才進 shell。
  - 補齊 page-level tests 的 `studioV2RulesAPI.list` mocks，並讓新增的 Step 2 autosave tests 與既有 device/workspace/shell tests 全部轉綠。
- 建立/修改的檔案：
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts` (修改)
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts` (建立)
  - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/rule-autosave-page.test.tsx` (建立)
  - `frontend/tests/unit/hooks/useStudioV2Rules.test.ts` (建立)
  - `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/shell.test.tsx` (修改)
  - `openspec/changes/wire-studio-v2-step2-rule-autosave/tasks.md` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 20：依 `docs/goal.md` 串上 Studio V2 Step 3 mapping autosave
- **狀態：** complete
- **開始時間：** 14:35
- 執行的操作：
  - 先補齊 backend red tests，確認 workspace mappings handler/router 與 mapping service rebind/validation sentinel 都是缺口，而不是測試本身編譯壞掉。
  - 在 backend 新增 workspace-scoped mapping list/create/update/delete handler，使用 `rule_id + address` 解析真正的 persisted point / source-rule link，並把 row payload 轉成 tag + mapping + transform pipeline。
  - 在 `mapping.Service` 補上 `ErrValidation`、`TagID` rebind update 與 validation error chain，讓 workspace autosave 與既有 generic handler 都能正確回傳 400。
  - 在 frontend 新增 `studioV2Mappings` service / React Query hook，並把 Step 3 row state 擴充為 `local_value`、`persisted_value`、`save_state`、`save_error` 與 persisted metadata。
  - 建立 `useStudioV2MappingAutosave.ts`，把 Step 3 hydration、valid-only save、per-row save isolation 與 orphaned persisted row delete 都集中在 page-level orchestration。
  - 在 `MappingRow.tsx` 補上最小 save-state 可視 badge，讓使用者能直接看出 `saving` / `saved` / `save-error` / `draft-invalid`。
  - 修正 mapping reconcile 兩個實作陷阱：避免 hydration effect 自我 dispatch loop，並避免 dirty local 值在 save queue 啟動前被 persisted hydration 蓋回去。
- 建立/修改的檔案：
  - `internal/datalink/mapping/errors.go` (修改)
  - `internal/datalink/mapping/service_crud.go` (修改)
  - `internal/datalink/sourcerule/service_links.go` (建立)
  - `internal/api/handlers/mapping_handler.go` (修改)
  - `internal/api/handlers/studio_v2_workspace_mappings_handler.go` (建立)
  - `internal/api/router.go` (修改)
  - `internal/datalink/mapping/service_workspace_scope_test.go` (建立)
  - `internal/api/handlers/studio_v2_workspace_mappings_handler_test.go` (建立)
  - `internal/api/router_studio_v2_workspace_mappings_test.go` (建立)
  - `frontend/src/types/datalink.ts` (修改)
  - `frontend/src/services/studioV2Mappings.ts` (建立)
  - `frontend/src/hooks/datalink/keys.ts` (修改)
  - `frontend/src/hooks/datalink/useStudioV2Mappings.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/types.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/studioV2MappingAutosave.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx` (修改)
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts` (建立)
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts` (修改)
  - `frontend/tests/unit/hooks/useStudioV2Mappings.test.ts` (建立)
  - `frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx` (建立)
  - `openspec/changes/wire-studio-v2-step3-mapping-autosave/tasks.md` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)
- 驗證：
  - `go test ./internal/datalink/mapping ./internal/api ./internal/api/handlers`
  - `cd frontend && npm run test -- --run tests/unit/hooks/useStudioV2Mappings.test.ts tests/unit/workbench-v2/mapping-autosave-page.test.tsx tests/unit/workbench-v2/step3-mapping.test.tsx tests/unit/workbench-v2/reducer-step3.test.ts tests/unit/workbench-v2/shell.test.tsx tests/unit/workbench-v2/workspace-boot.test.tsx`
  - `cd frontend && npm run build`
  - `git diff --check`
- 驗證：
  - `cd frontend && npm run test -- --run tests/unit/hooks/useStudioV2Rules.test.ts tests/unit/workbench-v2/rule-autosave-page.test.tsx`
  - `cd frontend && npm run test -- --run tests/unit/workbench-v2/device-autosave-page.test.tsx tests/unit/workbench-v2/workspace-boot.test.tsx tests/unit/workbench-v2/shell.test.tsx tests/unit/workbench-v2/step2-rule.test.tsx tests/unit/workbench-v2/reducer-step2.test.ts tests/unit/hooks/useStudioV2Workspace.test.tsx tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts tests/unit/workbench-v2/routing.test.tsx`
  - `cd frontend && npm run build`
  - `git diff --check`
  - `CLAUDE.md` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

### 階段 15：強化新對話自動接手入口
- **狀態：** complete
- **開始時間：** 19:50
- 執行的操作：
  - 規劃新增 `CURRENT_STATE.md`，把「目前做到哪、不要重做什麼、下一步最可能做什麼」濃縮成最新快照。
  - 擴大 `AGENTS.md` / `CLAUDE.md` 的觸發條件，不只 inventory 任務，連 `/studio`、`/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 相關任務都要先讀 onboarding 入口。
  - 建立 `CURRENT_STATE.md`，固定記錄產品決策、最近新增內容、下次必讀順序與不要重做的事情。
  - 更新 `START_HERE.md`、`context.json`、`README.md`，把 onboarding 順序擴成 `START_HERE -> context.json -> CURRENT_STATE`。
  - 更新 `AGENTS.md` 與 `CLAUDE.md`，要求凡是碰 `studio` surfaces 或 inventory 任務都必須先走這個順序。
  - 新增 Codex memory note 到 `~/.codex/memories/extensions/ad_hoc/notes/2026-05-29T19-55-studio-surface-inventory-onboarding.md`，讓新對話更容易直接接上。
  - 以 `go run ./cmd/studio_inventory_changelog add ...` 寫入接手機制升級紀錄。
  - 執行 `go run ./cmd/studio_inventory_changelog list -limit 5`、`python3 -m json.tool docs/technical/studio-surface-inventory/context.json`、`git diff --check` 驗證。
- 建立/修改的檔案：
  - `docs/technical/studio-surface-inventory/CURRENT_STATE.md` (建立)
  - `docs/technical/studio-surface-inventory/START_HERE.md` (修改)

### 階段 18：依 `docs/goal.md` 串上 Studio V2 Step 1 device autosave
- **狀態：** complete
- **開始時間：** 13:40
- 執行的操作：
  - 讀取 `wire-studio-v2-step1-device-autosave` proposal / design / tasks / specs，確認 scope 只落在 Step 1 device autosave / order / delete cascade。
  - 先以 TDD 補 backend service/router tests，鎖定 workspace-scoped device list/create/update/delete/order 契約與 invalid isolation。
  - 後端新增 workspace attach/detach/order service、workspace device handler 與路由；device service 補上 name validation 與 client-supplied id passthrough。
  - 前端新增 `studioV2WorkspaceDevices` service / hooks、device hydrate/validate helpers、page-level autosave hook，讓合法裝置逐台 create/update、不合法保留本地值。
  - 在 Step 1 tab/editor 補上 per-device save markers，顯示 `draft-invalid`、`saving`、`saved`、`save-error` 與錯誤訊息。
  - 更新既有 shell/workspace boot tests 的 QueryClient 與 devices bootstrap mocks，避免新 bootstrap query 讓舊 tests 停在 loading/error。
  - 執行 targeted Go tests、Vitest autosave/shell/workspace/routing suites、`npm run build` 與 `git diff --check`。
- 建立/修改的檔案：
  - `internal/datalink/device/errors.go` (建立)
  - `internal/datalink/device/service_crud.go` (修改)
  - `internal/datalink/workspace/service_devices.go` (建立)
  - `internal/datalink/workspace/service_device_order_test.go` (建立)
  - `internal/api/handlers/studio_v2_workspace_devices_handler.go` (建立)
  - `internal/api/handlers/studio_v2_workspace_devices_handler_test.go` (建立)
  - `internal/api/router.go` (修改)
  - `internal/api/router_studio_v2_workspace_devices_test.go` (建立)
  - `frontend/src/services/studioV2Workspace.ts` (修改)
  - `frontend/src/services/studioV2WorkspaceDevices.ts` (建立)
  - `frontend/src/hooks/datalink/keys.ts` (修改)
  - `frontend/src/hooks/datalink/index.ts` (修改)
  - `frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/state/types.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/studioV2DeviceAutosave.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/steps/step1/DeviceTabRail.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx` (修改)
  - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx` (修改)
  - `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts` (建立)
  - `frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts` (建立)
  - `frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/step1.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/shell.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx` (修改)
  - `openspec/changes/wire-studio-v2-step1-device-autosave/tasks.md` (修改)
  - `docs/technical/studio-surface-inventory/context.json` (修改)
  - `docs/technical/studio-surface-inventory/README.md` (修改)
  - `AGENTS.md` (修改)
  - `CLAUDE.md` (修改)
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)

## 測試結果
| 測試 | 輸入 | 預期結果 | 實際結果 | 狀態 |
|------|------|---------|---------|------|
| Vitest 測試集 (V2) | `npm run test -- --run workbench-v2` | 31 檔案 207 測試全部通過 | 31 檔案 207 測試全數通過 | success |
| Vitest 測試集 (Dashboard) | `npm run test -- --run runtime-dashboard` | 3 檔案 13 測試全部通過 | 3 檔案 13 測試全數通過 | success |
| Go 測試集 | `go test ./...` | 後端測試全部通過 | 後端測試全部通過 | success |
| 前端 Lint 檢查 | `npm run lint` | 0 errors | 0 errors, 22 warnings | success |
| 前端打包建置 | `npm run build` | 建置成功無錯誤 | 建置成功無錯誤 | success |
| 行數規範檢查 | `make check-lines` | 通過 | 通過 (僅 state 測試檔 warning，實作檔全數符合上限) | success |

## 錯誤日誌
| 時間戳記 | 錯誤 | 嘗試次數 | 解決方案 |

## 會話：2026-05-30

### 階段 16：依 `docs/goal.md` 實作預設 `/studio/v2` 入口
- **狀態：** complete
- **開始時間：** 10:15
- 執行的操作：
  - 讀取 `docs/goal.md`、`AGENTS.md`、`CLAUDE.md`、`docs/technical/studio-surface-inventory/START_HERE.md`、`context.json`、`CURRENT_STATE.md`。
  - 確認 `spectra list` 中 10 個 Studio V2 change 都存在，並依 `docs/goal.md` 決定先從 `make-studio-v2-default-entry` 開始。
  - 執行 `spectra status --change "make-studio-v2-default-entry" --json`、`spectra instructions apply --change "make-studio-v2-default-entry" --json`、`spectra analyze make-studio-v2-default-entry --json`，確認 change ready、preflight clean、artifact analysis clean。
  - 讀取 proposal / design / tasks / spec，確認 `tasks=4`、`SHALL=2`、受影響路徑為 `/`、unknown fallback、`/datalink`、`/datalink/workbench`、`/studio`。
  - 讀取 `App.tsx`、`legacyRoutes.ts`、`routing.test.tsx` 與 `.github/instructions/*`，確認目前 route tree 仍把首頁與 fallback 導向 `/studio`。
  - 取得 Spectra `tdd` / `audit` 指令，準備先寫 failing routing tests。
  - 嘗試建立分支 `feature/docs-goal-studio-v2-default-entry` 時遇到 `.git/refs/heads/...lock: Operation not permitted`；改以提升權限重跑後成功切換到新分支。
  - 先修改 `frontend/tests/unit/workbench-v2/routing.test.tsx`，把 `/`、unknown fallback、`/datalink`、`/datalink/workbench` 的期望改為 V2，並保留 `/studio` 為 legacy。
  - 依 `docs/goal.md` 文字先執行 `cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/routing.test.tsx`，確認因 Vitest include 規則只收 `tests/**` 而找不到測試檔；之後改用有效命令 `cd frontend && npm run test -- --run tests/unit/workbench-v2/routing.test.tsx` 取得真正 RED。
  - RED 結果顯示 6 個測試中有 4 個失敗，且 DOM 都落在 `legacy-workbench`，證明 `/`、unknown fallback、`/datalink`、`/datalink/workbench` 仍走 `/studio`。
  - 在 `frontend/src/features/datalink/legacyRoutes.ts` 新增 `buildWorkbenchV2EntryRedirect()`，並於 `frontend/src/App.tsx` 新增 `GuidedWorkbenchEntryRedirect`，讓 `/`、`*`、`/datalink`、`/datalink/workbench` 改導向 `/studio/v2`；保留 `/studio` 與 `/datalink/workbench/*` 不變。
  - 重跑 `cd frontend && npm run test -- --run tests/unit/workbench-v2/routing.test.tsx`，6/6 通過。
  - 執行 `git diff --check` 通過，並檢查 route-entry diff 只涉及 `App.tsx`、`legacyRoutes.ts`、`routing.test.tsx` 的產品面變更。
  - 將 `openspec/changes/make-studio-v2-default-entry/tasks.md` 四個 task 全部勾選完成。
- 建立/修改的檔案：
  - `task_plan.md` (修改)
  - `findings.md` (修改)
  - `progress.md` (修改)
  - `frontend/tests/unit/workbench-v2/routing.test.tsx` (修改)
  - `frontend/src/features/datalink/legacyRoutes.ts` (修改)
  - `frontend/src/App.tsx` (修改)
  - `openspec/changes/make-studio-v2-default-entry/tasks.md` (修改)

### 階段 17：依 `docs/goal.md` 建立 Studio V2 singleton workspace foundation
- **狀態：** complete
- **開始時間：** 13:20
- 執行的操作：
  - 讀取 `add-studio-v2-single-workspace-foundation` 的 proposal / design / tasks / specs，確認 `tasks=6`，目標是 singleton workspace、bootstrap API、前端 page boot。
  - 補出 backend RED 測試：
    - `internal/datalink/workspace/service_test.go`
    - `internal/api/router_studio_v2_workspace_test.go`
    - `internal/api/handlers/studio_v2_workspace_handler_test.go`
  - 補出 frontend RED 測試：
    - `frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx`
    - `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx`
  - RED 證據：
    - `go test ./internal/datalink/workspace ./internal/api ./internal/api/handlers` 一開始因 package 尚不存在而編譯失敗。
    - `npm run test -- --run tests/unit/hooks/useStudioV2Workspace.test.tsx tests/unit/workbench-v2/workspace-boot.test.tsx` 一開始因 hook/service 檔案不存在而失敗。
  - 實作 `internal/datalink/workspace/` service + SQL/memory repository，使用既有 `system_settings` 表持久化 singleton workspace metadata。
  - 實作 `internal/api/handlers/studio_v2_workspace_handler.go`，在 `internal/api/router.go` 註冊 `GET /api/v1/datalink/studio-v2/workspace`，並在 `cmd/test_ui/main.go` 注入 workspace service。
  - 實作前端 `frontend/src/services/studioV2Workspace.ts`、`frontend/src/hooks/datalink/useStudioV2Workspace.ts`、`StudioV2Workspace` 型別與 query key，讓 `DatalinkWorkbenchV2Page.tsx` 先 bootstrap workspace API，再顯示 shell；失敗時顯示可行動錯誤訊息。
  - 更新 `frontend/tests/unit/workbench-v2/shell.test.tsx`，以 mock workspace hook 保持既有 shell 測試聚焦在 UI 本身。
  - 執行 `gofmt -w ...` 格式化所有新增/修改的 Go 檔。
  - 第一次後端 GREEN 驗證時，`TestService_GetOrCreateDoesNotImportLegacyDevices` 因走 `device.Service` 需要 protocol adapter 而失敗；改成直接用 `device.NewSQLRepository(db).Create(...)` 建立 legacy device 後通過，避免把測試焦點帶離 workspace isolation contract。
  - 後端測試在 sandbox 內因 Go build cache 權限失敗，需用提升權限重跑。
  - GREEN 驗證通過：
    - `go test ./internal/datalink/workspace ./internal/api ./internal/api/handlers`
    - `go test ./cmd/test_ui ./internal/datalink/workspace ./internal/api ./internal/api/handlers`
    - `cd frontend && npm run test -- --run tests/unit/hooks/useStudioV2Workspace.test.tsx tests/unit/workbench-v2/workspace-boot.test.tsx tests/unit/workbench-v2/shell.test.tsx tests/unit/workbench-v2/routing.test.tsx`
    - `git diff --check`
  - 將 `openspec/changes/add-studio-v2-single-workspace-foundation/tasks.md` 六個 task 全部勾選完成。
- 建立/修改的檔案：
  - `internal/datalink/workspace/service.go` (建立)
  - `internal/datalink/workspace/sql_repository.go` (建立)
  - `internal/datalink/workspace/memory_repository.go` (建立)
  - `internal/datalink/workspace/service_test.go` (建立)
  - `internal/api/handlers/studio_v2_workspace_handler.go` (建立)
  - `internal/api/handlers/studio_v2_workspace_handler_test.go` (建立)
  - `internal/api/router_studio_v2_workspace_test.go` (建立)
  - `internal/api/router.go` (修改)
  - `cmd/test_ui/main.go` (修改)
  - `frontend/src/services/studioV2Workspace.ts` (建立)
  - `frontend/src/hooks/datalink/useStudioV2Workspace.ts` (建立)
  - `frontend/src/hooks/datalink/index.ts` (修改)
  - `frontend/src/hooks/datalink/keys.ts` (修改)
  - `frontend/src/types/datalink.ts` (修改)
  - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx` (修改)
  - `frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/shell.test.tsx` (修改)
  - `openspec/changes/add-studio-v2-single-workspace-foundation/tasks.md` (修改)

### 階段 21：依 `docs/goal.md` 串上 Studio runtime workspace device switching
- **狀態：** complete
- **開始時間：** 16:09
- 執行的操作：
  - 讀取 `add-studio-runtime-workspace-device-switching` proposal / design / tasks / specs，確認 `tasks=6`、`SHALL=6`，最小 payload 是 `workspace_id`、ordered `devices[]`、`default_device_id`。
  - 補出 backend RED 測試：
    - `internal/api/handlers/runtime_handler_test.go`
    - `internal/api/router_studio_v2_runtime_context_test.go`
  - 補出 frontend RED 測試：
    - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`
    - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`
    - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx`
  - RED 證據：
    - Go 測試先失敗於 `WorkspaceContext` method / route 尚不存在。
    - Vitest 先失敗於 runtime 仍是 `missing-device-context` boot、無 empty workspace state、無 unavailable visibility。
  - 實作 backend `GET /api/v1/datalink/studio-v2/workspace/runtime-context`，在 `RuntimeHandler` 內以 workspace order 組出 runtime context，並保留 `default_device_id` 為第一台 available device。
  - 將 router wiring 接到既有 `RuntimeHandler`，不重做 Step 4 handoff；`device_id` query 只保留 optional override。
  - 實作前端 `studioV2RuntimeContext` service / hook / query key，改寫 `useRuntimeDashboardState` 先吃 workspace runtime context，再決定選中設備、snapshot 與 stream。
  - 補 `empty-workspace` route state，讓沒有任何 available device 時留在 `/studio/runtime` 顯示引導，而不是 redirect。
  - 改寫 `FocusedDeviceHeader`，讓 unavailable device 保留在切換清單中並顯示原因；route tests 另外 mock `/studio/v2` page，避免把 V2 bootstrap loading 混進純 routing 斷言。
  - GREEN 驗證通過：
    - `go test ./internal/api ./internal/api/handlers`
    - `cd frontend && npm run test -- --run tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx tests/unit/workbench-v2/shell.test.tsx`
    - `cd frontend && npm run build`
    - `git diff --check`
  - 將 `openspec/changes/add-studio-runtime-workspace-device-switching/tasks.md` 六個 task 全部勾選完成。
- 建立/修改的檔案：
  - `internal/api/handlers/runtime_handler.go` (修改)
  - `internal/api/handlers/runtime_handler_test.go` (建立)
  - `internal/api/router.go` (修改)
  - `internal/api/router_studio_v2_runtime_context_test.go` (建立)
  - `frontend/src/services/studioV2RuntimeContext.ts` (建立)
  - `frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts` (建立)
  - `frontend/src/hooks/datalink/keys.ts` (修改)
  - `frontend/src/hooks/datalink/index.ts` (修改)
  - `frontend/src/types/datalink.ts` (修改)
  - `frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts` (修改)
  - `frontend/src/features/datalink/runtime-dashboard/FocusedDeviceHeader.tsx` (修改)
  - `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx` (修改)
  - `frontend/src/i18n/locales/zh-TW/runtime-dashboard.json` (修改)
  - `frontend/src/i18n/locales/en/runtime-dashboard.json` (修改)
  - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx` (修改)
  - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx` (修改)
  - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx` (修改)
  - `openspec/changes/add-studio-runtime-workspace-device-switching/tasks.md` (修改)
|----------|------|---------|---------|
| 03:01 | targetFile 寫入錯誤 (IsArtifact 與 Path 不匹配) | 1 | 移除非 artifact 標識直接寫入專案目錄 |
| 03:02 | vitest 無法識別 tests/workbench-v2 目錄中的測試 | 1 | 將測試檔案移至符合 include 規則的 tests/unit/workbench-v2 目錄 |
| 03:04 | React is not defined 錯誤 | 1 | 在 useWorkbenchV2State.ts 導入 React |
| 03:04 | TestingLibraryElementError (重疊 text 斷言衝突) | 1 | 改用 data-testid 或 getAllByText 長度斷言 |
| 03:06 | tsc 類別型別導入與 db.targets 解構 TypeScript 錯誤 | 1 | 加上 import type 與調整 db.targets 為強型別安全判斷 |
| 08:16 | DeviceTabRail Hook 違反與 ConnectionConfigForm tsc compile 錯誤 | 1 | 移除 map 內的 hook 改呼叫純函數 getColorTheme；將 config 屬性加 type assertions |
| 09:02 | 舊版 SourceStep 相關測試失敗 | 1 | 驗證發現為基礎 commit 歷史中重構 Source Step 後未同步更新舊測試所致，非本分支變更引起。 |
| 12:16 | 日誌滾動與 prepending logs 順序不咬合 | 1 | 將新 log 插入行為改為 append 至陣列末尾，並更新 logs index render，使其符合 terminal 排版與 auto-scroll。 |

### 階段 23：將 Studio V2 Step 1 執行測試改為真實 diagnostics
- **狀態：** complete
- **執行的操作：**
  - 補出 frontend RED 測試：
    - `frontend/tests/unit/workbench-v2/step1.test.tsx`
    - `frontend/tests/unit/workbench-v2/reducer-step1.test.ts`
  - RED 證據：
    - `cd frontend && npm run test -- --run tests/unit/workbench-v2/step1.test.tsx`
    - 初始失敗點證明 `Step1Device` 並未發送 `testDraftConnection` request，failure path 也不會顯示 backend 訊息。
  - 實作 `Step1Device` 改接 `useTestDraftConnectionMutation`，移除 `setTimeout` / `Math.random` mock animation。
  - 調整 reducer / panel，讓 Step 1 測試狀態直接以 backend `connect/probe` 兩段結果為準，continue gate 僅接受真實 diagnostics success。
  - 補出 backend draft diagnostics tests：
    - `internal/api/handlers/device_handler_extended_test.go`
    - 覆蓋 success、connect failure、probe failure。
  - 第一次執行 `go test ./internal/api/handlers -run 'TestDeviceHandler_TestDraftConnection'` 因 sandbox Go build cache 權限失敗，改用提升權限重跑。
  - GREEN 驗證通過：
    - `cd frontend && npm run test -- --run tests/unit/workbench-v2/step1.test.tsx tests/unit/workbench-v2/reducer-step1.test.ts`
    - `go test ./internal/api/handlers -run 'TestDeviceHandler_TestDraftConnection'`
    - `cd frontend && npm run build`
    - `git diff --check`
  - 手動 demo：
    - 在本地 `/studio/v2` 將 host 改為 `bad-host.invalid` 後執行測試，UI 直接顯示 backend DNS failure `lookup bad-host.invalid: no such host`。
    - 同一畫面 `全部建立並繼續` 維持 disabled，證明沒有真實 diagnostics success 時不能宣告通過。
  - 更新 `openspec/changes/replace-step1-mock-test-with-live-diagnostics/tasks.md`，勾選全部 task。
- **建立/修改的檔案：**
  - `frontend/src/features/datalink/workbench-v2/steps/step1/Step1Device.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/steps/step1/ReadinessStages.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/types.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts` (修改)
  - `frontend/tests/unit/workbench-v2/step1.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/reducer-step1.test.ts` (修改)
  - `internal/api/handlers/device_handler_extended_test.go` (修改)
  - `openspec/changes/replace-step1-mock-test-with-live-diagnostics/tasks.md` (修改)

### 階段 24：盤點 Studio V2 settings backend wiring 缺口
- **狀態：** complete
- **執行的操作：**
  - 讀取 `wire-studio-v2-settings-backend` proposal / design / tasks / specs。
  - 確認前端已存在 `settingsAPI` 與 `dbTargetAPI` service，settings surface 目前缺的是 boot/save/connector CRUD/test wiring，不是缺 endpoint。
  - 確認 `/settings` handler/service 雖只有少數預設 key，但 `PUT /settings/:key` 實際可持久化 arbitrary key/value，足夠承接多數 V2 settings 欄位。
  - 確認 connector pool 可直接重用 `db-targets/connectors/*` 契約，現行 mock 只是在 `ConnectorPoolSection.tsx` 層用 timer/Math.random 假裝測試。
  - 新增 `frontend/src/hooks/datalink/useSettings.ts` 與 `frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts`，把 `/settings` 與 `db-targets/connectors/*` 封裝成 V2 settings boot/save/connector CRUD/test 所需的 query/mutation 與 mapping layer。
  - 改寫 `SettingsPage.tsx`，加入 backend hydration gate、loading/error UI、save failure banner，並把 connector add/update/delete/test 全部改為真實 backend mutations。
  - 改寫 `ConnectorPoolSection.tsx`，移除 timer/`Math.random` mock test，讓 row 操作完全由 page-level callbacks 驅動。
  - 補齊前端測試：
    - `frontend/tests/unit/workbench-v2/settings-backend.test.tsx`
    - `frontend/tests/unit/workbench-v2/settings.test.tsx`
    - `frontend/tests/unit/workbench-v2/settings-connectors.test.tsx`
    - `frontend/tests/unit/workbench-v2/reducer-settings.test.ts`
  - 補齊後端 handler round-trip 測試：
    - `internal/api/handlers/dbtarget_handler_connectors_test.go`
    - 既有 `internal/api/handlers/settings_handler_test.go`
  - 第一次 `npm run build` 失敗於 `SettingsConnector.status` / `completeConnectorTest` union 未同步；補齊 `state/useWorkbenchV2State.ts` 與 `state/types-settings.test-d.ts` 後通過。
  - 第一次 `go test ./internal/api/handlers -run 'TestDatabaseTargetHandler_ConnectorCrudAndTestRoundTrip|TestSettingsHandler_List|TestSettingsHandler_Update'` 因 sandbox Go build cache 權限失敗，改用提升權限重跑。
  - GREEN 驗證通過：
    - `cd frontend && npm run test -- --run tests/unit/workbench-v2/settings-backend.test.tsx tests/unit/workbench-v2/settings.test.tsx tests/unit/workbench-v2/settings-connectors.test.tsx tests/unit/workbench-v2/reducer-settings.test.ts`
    - `cd frontend && npm run build`
    - `go test ./internal/api/handlers -run 'TestDatabaseTargetHandler_ConnectorCrudAndTestRoundTrip|TestSettingsHandler_List|TestSettingsHandler_Update'`
    - `git diff --check`
  - 手動 demo：
    - 本地 `/studio/v2/settings` 新增 connector 後，畫面立即顯示真實 backend 診斷 `dial tcp 127.0.0.1:5432: connect: connection refused`，不是 mock timer。
    - 將 `批次寫入限制條數` 從 `1000` 改成 `1001`、按 `儲存所有設定`、reload 後再進入 settings，值仍保留為 `1001`，證明 save bar 不再是 noop；之後已恢復回 `1000` 並刪除測試 connector。
  - 更新 `openspec/changes/wire-studio-v2-settings-backend/tasks.md`，勾選全部 task。
 - **建立/修改的檔案：**
  - `frontend/src/hooks/datalink/useSettings.ts` (建立)
  - `frontend/src/hooks/datalink/index.ts` (修改)
  - `frontend/src/hooks/datalink/keys.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/settings/backendMappings.ts` (建立)
  - `frontend/src/features/datalink/workbench-v2/settings/SettingsPage.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/settings/ConnectorPoolSection.tsx` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/types.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts` (修改)
  - `frontend/src/features/datalink/workbench-v2/state/types-settings.test-d.ts` (修改)
  - `frontend/tests/unit/workbench-v2/settings-backend.test.tsx` (建立)
  - `frontend/tests/unit/workbench-v2/settings.test.tsx` (修改)
  - `frontend/tests/unit/workbench-v2/settings-connectors.test.tsx` (修改)
  - `internal/api/handlers/dbtarget_handler_connectors_test.go` (建立)
  - `openspec/changes/wire-studio-v2-settings-backend/tasks.md` (修改)

## 五問重啟檢查
| 問題 | 答案 |
|------|------|
| 我在哪裡？ | 階段 9：優化項目補齊與日誌系統升級 |
| 我要去哪裡？ | 任務完成，提交代碼並回報使用者 |
| 目標是什麼？ | 完美交付 datalink-workbench-v2 與 runtime-dashboard 變更與效能優化 |
| 我學到了什麼？ | SSE 高頻推送時除了數值要 Throttling 之外，狀態轉移日誌若去重機制設計不周，會造成嚴重的 log 爆量與 rendering 效能問題；應以 Ref 緩衝前一次點位狀態來進行 transition logs 處理。 |
| 我做了什麼？ | 補齊並深度優化了診斷日誌系統的 transition/deduplication 機制、修正了滾動 UX，並通過全部 220+ 測試與 production 打包。 |

---
*每個階段完成後或遇到錯誤時更新此檔案*
