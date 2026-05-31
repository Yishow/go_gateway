# 任務計畫：實作 Datalink Workbench V2 變更

## 目標
依序完成 6 個 `spectra-apply` 變更的實作，確保前端與後端功能整合完畢，並通過所有測試與 lint 檢查。

## 目前階段
階段 1：建立分支與規劃

## 各階段

### 階段 1：建立分支與規劃
- [x] 另開 git branch `feature/datalink-workbench-v2`
- [x] 建立並初始化規劃檔案（`task_plan.md`, `findings.md`, `progress.md`）
- **狀態：** complete

### 階段 2：實作 datalink-workbench-v2-shell
- [x] 執行 `spectra status --change "datalink-workbench-v2-shell" --json` 取得任務
- [x] 依序完成 shell 的任務
- [x] 執行單元測試與驗證
- [x] 將變更標記為 done 並封存或進行下一階段
- **狀態：** complete

### 階段 3：實作 datalink-workbench-v2-step1-device
- [x] 執行 `spectra status --change "datalink-workbench-v2-step1-device" --json` 取得任務
- [x] 依序完成 device 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 4：實作 datalink-workbench-v2-step2-rule
- [x] 執行 `spectra status --change "datalink-workbench-v2-step2-rule" --json` 取得任務
- [x] 依序完成 rule 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 5：實作 datalink-workbench-v2-step3-mapping
- [x] 執行 `spectra status --change "datalink-workbench-v2-step3-mapping" --json` 取得任務
- [x] 依序完成 mapping 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 6：實作 datalink-workbench-v2-step4-database
- [x] 執行 `spectra status --change "datalink-workbench-v2-step4-database" --json` 取得任務
- [x] 依序完成 database 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 7：實作 datalink-workbench-v2-settings
- [x] 執行 `spectra status --change "datalink-workbench-v2-settings" --json` 取得任務
- [x] 依序完成 settings 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 8：最終驗證與整合
- [x] 執行完整建置：前端 build、複製 embed static、後端 build
- [x] 執行 Go 和 React 的所有單元與整合測試
- [x] 檢查單檔行數限制
- **狀態：** complete

### 階段 9：優化項目補齊與日誌系統升級
- [x] 實作診斷日誌面板 (RealtimeLogsPanel.tsx) 與 `stale` / `recovered` 狀態去重日誌。
- [x] 將日誌改為 append 末尾以匹配 auto-scroll 效果。
- [x] 撰寫整合單元測試並通過 `npm run test` 與 `go test`。
- **狀態：** complete

### 階段 10：建立 studio surface 與 API 維護文件
- [x] 盤點 `/studio`、`/studio/v2`、`/studio/runtime`、`/gateway/*`、`/test` 的定位與成熟度
- [x] 盤點前端實際使用的 query / mutation / SSE 與對應後端 API
- [x] 建立多份 Markdown 文件，分開記錄主線頁面、V2/runtime、gateway/test、API registry 與 gap roadmap
- [x] 建立一份可直接開啟的 HTML 總覽頁
- **狀態：** complete

### 階段 11：依產品優先順序重構 surface inventory 文件
- [x] 將 `/studio` 標示為完整版主線但先暫停
- [x] 將 `/studio/v2` 標示為預設入口與重點施作
- [x] 將 `/test` 拆成獨立文件並標示先暫停
- [x] 將 `/gateway/*` 拆成獨立 experimental 記錄
- [x] 更新 HTML 總覽以反映上述優先順序
- **狀態：** complete

### 階段 12：修正 HTML 文件台無資料並拆分為 html/js/css
- [x] 找出 HTML 無資料的 root cause
- [x] 將文件台拆為 `index.html`、`inventory.css`、`inventory.js`
- [x] 驗證 JS 語法與基本檔案關聯
- **狀態：** complete

### 階段 13：為 studio surface inventory 建立 changelog 機制
- [x] 在 `AGENTS.md` / `CLAUDE.md` 寫明 inventory 更新必須記錄
- [x] 新增 SQLite changelog CLI
- [x] 初始化 changelog database 並寫入本次變更
- [x] 更新 inventory README 使用說明
- [x] 驗證 CLI 與資料庫內容
- **狀態：** complete

### 階段 14：補齊 inventory onboarding 入口
- [x] 新增 `START_HERE.md` 作為第一入口
- [x] 新增 `context.json` 作為 machine-readable 摘要
- [x] 更新 `AGENTS.md` / `CLAUDE.md` 指向 onboarding 入口
- [x] 更新 `README.md` 閱讀順序與 onboarding 規則
- [x] 將本次文件更新寫入 SQLite changelog
- [x] 驗證文件與 changelog
- **狀態：** complete

### 階段 15：強化新對話自動接手入口
- [x] 新增 `CURRENT_STATE.md` 作為最新狀態快照
- [x] 更新 onboarding 文件與 `context.json`
- [x] 強化 `AGENTS.md` / `CLAUDE.md`，把 `studio` surface 任務必讀順序寫死
- [x] 補一筆 changelog 記錄這次接手機制升級
- [x] 補一筆 Codex memory note，讓新對話更容易直接接上
- [x] 驗證文件、JSON 與 changelog
- **狀態：** complete

### 階段 16：依 `docs/goal.md` 實作預設 `/studio/v2` 入口
- [x] 讀取 `docs/goal.md`、`AGENTS.md`、`CLAUDE.md` 與 studio inventory onboarding 入口
- [x] 確認第一個 change `make-studio-v2-default-entry` 的 task / SHALL / scope / 驗證目標
- [x] 建立分支 `feature/docs-goal-studio-v2-default-entry`
- [x] 先補 routing failing tests，再只修改 route-entry surface
- [x] 執行指定 Vitest 與 `git diff --check`
- [x] 更新 `openspec/changes/make-studio-v2-default-entry/tasks.md` 完成狀態
- **狀態：** complete

### 階段 17：依 `docs/goal.md` 建立 Studio V2 singleton workspace foundation
- [x] 讀取 `add-studio-v2-single-workspace-foundation` proposal / design / tasks / specs
- [x] 先補 backend/frontend failing tests，再只實作 workspace metadata/bootstrap surface
- [x] 建立 singleton workspace service / repo、bootstrap handler / router wiring、前端 service / hook / page boot
- [x] 執行指定 Go / Vitest 驗證與 `git diff --check`
- [x] 更新 `openspec/changes/add-studio-v2-single-workspace-foundation/tasks.md` 完成狀態
- **狀態：** complete

### 階段 18：依 `docs/goal.md` 串上 Studio V2 Step 1 device autosave
- [x] 讀取 `wire-studio-v2-step1-device-autosave` proposal / design / tasks / specs
- [x] 先補 backend/frontend failing tests，再只實作 workspace-scoped Step 1 device CRUD / order / delete contract
- [x] 建立 workspace device APIs、前端 autosave hook/service 與 save-state markers
- [x] 執行指定 Go / Vitest / build 驗證與 `git diff --check`
- [x] 更新 `openspec/changes/wire-studio-v2-step1-device-autosave/tasks.md` 完成狀態
- **狀態：** complete

### 階段 19：依 `docs/goal.md` 串上 Studio V2 Step 2 rule autosave
- [x] 讀取 `wire-studio-v2-step2-rule-autosave` proposal / design / tasks / specs
- [x] 先補 backend/frontend failing tests，再只實作 workspace-scoped Step 2 rule CRUD / ownership contract
- [x] 建立 workspace rule APIs、前端 autosave hook/service 與 save-state markers
- [x] 執行指定 Go / Vitest / build 驗證與 `git diff --check`
- [x] 更新 `openspec/changes/wire-studio-v2-step2-rule-autosave/tasks.md` 完成狀態
- **狀態：** complete

### 階段 20：依 `docs/goal.md` 串上 Studio V2 Step 3 mapping autosave
- [x] 讀取 `wire-studio-v2-step3-mapping-autosave` proposal / design / tasks / specs
- [x] 先補 backend/frontend failing tests，再只實作 workspace-scoped Step 3 mapping CRUD / ownership contract
- [x] 建立 workspace mapping APIs、前端 autosave hook/service 與 row `local_value` / `persisted_value` / `save_state`
- [x] 執行指定 Go / Vitest / build 驗證與 `git diff --check`
- [x] 更新 `openspec/changes/wire-studio-v2-step3-mapping-autosave/tasks.md` 完成狀態
- **狀態：** complete

### 階段 21：依 `docs/goal.md` 串上 Studio runtime workspace device switching
- [x] 讀取 `add-studio-runtime-workspace-device-switching` proposal / design / tasks / specs
- [x] 先補 backend/frontend failing tests，再只實作 workspace runtime context 與 runtime page switching surface
- [x] 建立 `GET /studio-v2/workspace/runtime-context`、前端 runtime context service/hook、workspace-first runtime boot 與 empty state
- [x] 執行指定 Go / Vitest / build 驗證與 `git diff --check`
- [x] 更新 `openspec/changes/add-studio-runtime-workspace-device-switching/tasks.md` 完成狀態
- **狀態：** complete

### 階段 22：修正 Studio V2 空 workspace seeded fallback
- [x] 確認空 workspace boot 仍保留 `INITIAL_STATE` seeded device/rule 的 root cause
- [x] 讓 workspace hydration 在 `devices=[] && rules=[]` 時明確清空本地 seeded state
- [x] 讓 `/studio/v2` 等待 workspace hydration 完成後再渲染 shell，避免假資料閃現
- [x] 補齊 workspace boot 測試，驗證空 workspace 進入時 `devices=0`、`rules=0`
- **狀態：** complete

### 階段 23：將 Studio V2 Step 1 執行測試改為真實 diagnostics
- [x] 補 frontend RED 測試，要求 `run test` 發送 `test-draft` request 並以 backend success/failure 決定 continue gate
- [x] 補 backend handler RED/GREEN 測試，覆蓋 draft diagnostics 的 success、connect failure、probe failure
- [x] 移除 `Step1Device` timer/mock animation，改接真實 `test-draft` mutation
- [x] 將 Step 1 測試狀態收斂為 backend `connect/probe` 兩段結果，不再依賴前端假步驟推進
- [x] 補手動 demo，驗證沒有真實 diagnostics success 時 Step 1 不能宣告通過
- **狀態：** complete

### 階段 24：盤點 Studio V2 settings backend wiring 缺口
- [x] 對照 `wire-studio-v2-settings-backend` change 契約與現有 settings surface
- [x] 確認前端已存在 `settingsAPI` / `dbTargetAPI`，缺口主要是 boot/save/connector wiring
- [x] 確認 `/settings` handler 雖只有少數預設 key，但可作為 generic key/value store 持久化更多欄位
- [x] 確認 connector pool 可直接重用既有 `db-targets/connectors/*` 契約
- [x] 完成 settings boot/save 與 connector CRUD/test backend wiring，移除 operator-visible mock / noop
- [x] 執行前端 settings tests、frontend build、backend handler tests
- [x] 補手動 demo，確認 `/studio/v2/settings` 實際操作不再暴露 mock / noop
- **狀態：** complete

### 階段 25：定義 Studio V2 Step 3 live preview 與 target type 快捷設定
- [x] 確認 Step 3 `轉換管線預覽` 目前仍是本地 mock / deterministic seed flow
- [x] 確認 backend 已存在可重用的 `POST /api/v1/datalink/mappings/preview`
- [x] 完成 approved design doc，固定 preview / payload / target type 快捷設定邊界
- [x] 建立 Spectra change `wire-studio-v2-step3-live-preview-and-target-type-shortcuts`
- [x] 補 proposal / design / spec / tasks，讓後續 implementation 可直接接上
- **狀態：** complete

### 階段 26：為 Studio V2 Step 3 live preview 撰寫 implementation plan
- [x] 讀取 approved design doc 與新建 Spectra change artifacts
- [x] 盤點 Step 3 現有檔案責任與測試入口，決定不把新 async coverage 繼續塞進既有 375 行測試檔
- [x] 建立 `docs/superpowers/plans/2026-05-30-v2-step3-live-preview-and-target-type-shortcuts.md`
- [x] 完成手動 plan review，確認 task 邊界、驗證命令與修改檔案集合可直接執行
- **狀態：** complete

### 階段 27：實作 Studio V2 Step 3 live preview 與 target type 快捷設定
- [x] 依 `wire-studio-v2-step3-live-preview-and-target-type-shortcuts` change 補 Step 3 live preview RED tests
- [x] 將 preview 改接既有 backend preview API，補 loading / error / stale response guard
- [x] 在 preview 卡新增 `target_type` 快捷設定與套用到全部列
- [x] 執行 Step 3 targeted tests、frontend build 與 `git diff --check`
- [x] 更新 Spectra tasks 與 repo 內 progress 文件
- **狀態：** complete

### 階段 28：修正 Studio V2 rule autosave 500 與 mapping stale-row 連鎖錯誤
- [x] 重現 `PUT /studio-v2/workspace/source-rules/:id` 在 device 未 ready 時回 500 的行為
- [x] 重現 Step 3 在 point row 身分變更後沿用舊 `mapping_id` 的 stale-row 錯誤
- [x] 修正 source-rule readiness block 改回 validation 4xx
- [x] 修正 mapping row identity reset，並阻止 rule 未保存成功時繼續打 mapping autosave
- [x] 執行 targeted Go/Vitest/build 驗證與 `git diff --check`
- **狀態：** complete

### 階段 29：讓 Studio V2 Step 3 顯示真實裝置值並以真值計算 preview
- [x] 補 Step 3 測試，要求 Point → Tag 表格顯示裝置即時值，preview 僅能吃真實 `raw_value`
- [x] 新增多裝置 `useStep3LiveValues` hook，聚合 runtime SSE live values
- [x] 移除 `useSelectedMapping` / `TransformPreview` 對 deterministic `rawSeed` 的依賴
- [x] 在 MappingTable / MappingRow 新增「裝置值」欄並顯示等待狀態
- [x] 執行 targeted Vitest、frontend build 與 `git diff --check`
- **狀態：** complete

### 階段 30：修正 Step 3 runtime stream 重連風暴
- [x] 補 RED 測試，證明 mapping 無關欄位變更時不應重建 runtime streams
- [x] 將 `useStep3LiveValues` 依賴收斂為 `device_id/address/persisted_point_id/localPointId`
- [x] 確保 `save_state/tag_key` 等 UI 狀態變更不再觸發 SSE 重連
- [x] 執行 targeted Vitest、frontend build 與 `git diff --check`
- **狀態：** complete

### 階段 31：修正 Step 3 數值顯示只吃 SSE、未退回 last_value
- [x] 補 RED 測試，要求沒有 SSE 時也能顯示 `points.last_value`
- [x] 在 `useStep3LiveValues` 補 persisted point snapshot 查詢
- [x] 讓 Step 3 顯示邏輯改為 `SSE raw_value` 優先，否則退回 `last_value`
- [x] 執行 targeted Vitest、frontend build 與 `git diff --check`
- **狀態：** complete

## 關鍵問題
1. 每一步驟的變更有無互相依賴？（需依序實作，後面的變更可能基於前面已實作的程式碼）
2. 測試與 Lint 是否在每個階段實作完都要通過？（是，確保每個階段提交時皆為 clean 狀態）

## 已做決策
| 決策 | 理由 |
|------|------|
| 新建分支 `feature/datalink-workbench-v2` | 使用者要求另開分支實作，不直接污染 main |
| 日誌改為時間遞增排序 (Append to array end) | 與 auto-scroll-to-bottom 對齊，提供流暢的控制台滾動 UX。 |
| 將 surface inventory 拆成多份 md + 一份 html 總覽 | 避免單一文件過大，讓後續維護與 change 規畫更容易定位 |
| 這一輪先只實作 `make-studio-v2-default-entry` | `docs/goal.md` 明確要求一次只跑一個 `/goal`，且建議順序先做預設入口 |
| 第二個 `/goal` 只落實 workspace metadata/bootstrap，不偷做 step autosave | spec 明確把 multi-workspace、legacy import、step autosave 排除在本 change 之外 |
| 第三個 `/goal` 維持 client `dev-*` ids 直接存後端 | 避免 Step 1 create 後需要立即 remap Step 2/3/4 本地關聯資料，先用後端 passthrough id 保持最小改動 |
| 第四個 `/goal` 只 autosave backend 既有 source-rule 契約欄位 | 後端目前沒有 `name` / `share_*` persisted shape，若硬接會製造 reload 漂移與假保存 |
| 空 workspace boot 先修 hydration gate，不先改全域 `INITIAL_STATE` | 可保留既有 reducer/單元測試預設資料，同時阻止 `/studio/v2` 真空工作區誤顯示 seeded draft |

## 遇到的錯誤
| 錯誤 | 嘗試次數 | 解決方案 |
|------|---------|---------|

## 備註
- 隨著進度更新階段狀態：pending → in_progress → complete
- 做重大決策前重新讀取此計畫（注意力操縱）
- 記錄所有錯誤，避免重複
