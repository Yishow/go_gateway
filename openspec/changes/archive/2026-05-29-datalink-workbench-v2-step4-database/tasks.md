## 1. State 層擴充

- [x] 1.1 在 `state/types.ts` 補完 `DbConnector`、`DbTarget`、`CommitLog`、`CommitState` 型別，落地設計決策「拆檔策略：8 個元件 + 3 個 state module」型別契約。**行為**：所有 step4 子元件 import 型別解析正確。**驗證**：`tsc --noEmit` 通過；新增 `types-step4.test-d.ts` type-level assertion。
- [x] 1.2 建立 `state/dbSchemas.ts`，匯出 `SAMPLE_DB_TABLES`（PostgreSQL `sensor_readings` 9 欄）、`getColumnsFor(kind)`、`getDefaultConnector(kind)`，落地設計決策「Sample table schema」。**行為**：四種 kind 都回傳同份 sample（暫定）；`getColumnsFor('postgres')[0].name === 'ts'` 且 `primary_key === true`。**驗證**：新增 `frontend/tests/workbench-v2/dbSchemas.test.ts`，覆蓋 4 種 kind 的 column 順序與 default connector 預設值。
- [x] 1.3 建立 `state/autoAssignTargets.ts`，匯出純函式 `autoAssignTargets(enabledPoints, mappings, columnNames, existingTargets)`，落地設計決策「自動欄位匹配演算法」並實作需求 **Tag-to-column auto-assignment**。**行為**：優先 existing > exact > index fallback (skip if used) > first unused；不主動產生衝突當 column ≥ point。**驗證**：新增 `frontend/tests/workbench-v2/autoAssignTargets.test.ts`，覆蓋 8 points × 9 columns 無衝突、existingTargets 保留、8 points × 4 columns 產生 wrap conflict、tag_key 後綴 match。
- [x] 1.4 建立 `state/commitLog.ts`，匯出 `buildCommitLogSequence(state)` 回傳 10 個 CommitLog，落地設計決策「Commit log 序列：純函式 + reducer 串聯」並實作需求 **Commit sequence and animation** 的 log 順序部分。**行為**：10 個 log label 與 detail 與設計決策的順序表一致；count 變數依 state 動態填入。**驗證**：新增 `frontend/tests/workbench-v2/commitLog.test.ts`，斷言 10 個 log 順序、label 字串包含「POST /devices」等 substring、detail 反映 state.devices/rules/points 數量。
- [x] 1.5 擴充 `useWorkbenchV2State` reducer 加 8 個 db actions（`updateDbConnector`、`upsertDbTarget`、`updateDbTarget`、`autoAssignDbTargets`、`startCommit`、`appendCommitLog`、`completeCommit`、`resetCommit`），落地設計決策「Commit log 序列：純函式 + reducer 串聯」並實作需求 **Connector configuration form**、**Commit sequence and animation**。**行為**：`startCommit` 把 `commit.status='running'`、`logs=[]`、`started_at=now`；`appendCommitLog` 累加；`completeCommit` 設 `committed=true`、`status='success'`、`finished_at=now`；`resetCommit` 清空。**驗證**：新增 `frontend/tests/workbench-v2/reducer-step4.test.ts`，8 個 action 各 1–2 case；覆蓋 commit lifecycle 完整鏈、immutable 更新。

## 2. UI 子元件

- [x] 2.1 建立 `steps/step4/KindSelector.tsx`，渲染 4 個 db kind 卡片（emoji + 名稱），落地需求 **Connector configuration form** 的 kind 選擇部分。**行為**：active kind 藍底；點卡片觸發父元件 `onChange(kind)`，父元件 dispatch `updateDbConnector({ kind })` 與 `autoAssignDbTargets(...)`；`state.committed === true` 時 disabled。**驗證**：新增 `frontend/tests/workbench-v2/step4-database.test.tsx` 中 `describe('KindSelector')` 覆蓋 4 卡片渲染、active highlight、disabled 條件、kind 切換觸發 auto-assign。
- [x] 2.2 建立 `steps/step4/WriteStrategy.tsx`，渲染 INSERT/UPSERT radio + 寫入間隔 number + 「秒」label。**行為**：radio 與 number 兩向綁定；committed 時 disabled。**驗證**：`step4-database.test.tsx` 中 `describe('WriteStrategy')` 測 radio 切換、間隔 onChange、disabled 條件。
- [x] 2.3 建立 `steps/step4/ConnectorSection.tsx`，組合 KindSelector + 連線欄位（連線名稱 / Host / Port / Database / Schema / Table / 使用者）+ WriteStrategy，落地需求 **Connector configuration form**。**行為**：每欄輸入 dispatch `updateDbConnector`；committed 時所有 input/select disabled。**驗證**：`step4-database.test.tsx` 中 `describe('ConnectorSection')` 測預設值、輸入更新 state、committed 後 readonly。
- [x] 2.4 建立 `steps/step4/TargetMappingTable.tsx`，渲染表格（Tag / 點位 / → / 資料表欄位 / 欄位型態 / 啟用）+ 衝突 banner，落地需求 **Column conflict detection**。**行為**：每列 select 選 columnNames（非 PK）；衝突列紅框 + alert icon；表尾 banner；committed 時 select 與 toggle disabled。**驗證**：`step4-database.test.tsx` 中 `describe('TargetMappingTable')` 測列數、衝突偵測（兩列同 column）、紅框 + banner 出現、停用一列衝突消失、committed 後 readonly。
- [x] 2.5 建立 `steps/step4/CommitSummary.tsx`，渲染 5 列摘要卡 + 寬版「提交並啟動排程器」按鈕 + 底部 info，落地需求 **Commit sequence and animation** 的 trigger UI 部分。**行為**：按鈕在衝突或無啟用 target 時 disabled；點擊 dispatch `startCommit`。**驗證**：`step4-database.test.tsx` 中 `describe('CommitSummary')` 測 5 列摘要值、按鈕 disabled 條件、點擊 dispatch。
- [x] 2.6 建立 `steps/step4/CommitProgress.tsx`，渲染動畫 log 列表（每筆 emerald 勾 + label + detail + 200）+ 底部 pulse-dot「正在執行下一個指令…」。**行為**：sweep-in 動畫；列表來自 `state.commit.logs`；commit.status === 'running' 時顯示 pulse-dot。**驗證**：新增 `frontend/tests/workbench-v2/step4-commit.test.tsx` 中 `describe('CommitProgress')` 測 log 列表渲染、pulse-dot 條件。
- [x] 2.7 建立 `steps/step4/CommitSuccessCard.tsx`，渲染 emerald 大卡（圓勾 icon + 設定已套用標題 + scheduler 啟動 subline + 前往 Runtime Dashboard secondary button），落地需求 **Commit completion card**。**行為**：subline 文字含 `~{write_interval_seconds}s`；按鈕點擊呼叫 `onCommit` prop。**驗證**：`step4-commit.test.tsx` 中 `describe('CommitSuccessCard')` 測 subline 含 interval、按鈕點擊 onCommit 一次。

## 3. 容器與整合

- [x] 3.1 建立 `steps/step4/Step4Database.tsx`，組裝 ConnectorSection + 12-col grid（TargetMappingTable col 7 + 右 col 5 三態切換 CommitSummary/CommitProgress/CommitSuccessCard）+ commit timer effect，落地設計決策「Commit log 序列：純函式 + reducer 串聯」與「Commit 後 form 只讀」並實作需求 **Commit sequence and animation** 的 `Animation timing` + `Resume on remount mid-commit` Scenario 與 **Read-only form after commit**。**行為**：useEffect 在 `commit.status === 'running' && logs.length < 10` 時啟動 280ms timer 推下一個 log；unmount cleanup clearTimeout；committed 時 useStep4Readonly() 回傳 true 並向下傳遞 disabled prop。**驗證**：`step4-commit.test.tsx` 主測試用 `vi.useFakeTimers()` 推 10 × 280ms、斷言完整 lifecycle；另一個 it 模擬 unmount → remount 後 effect 繼續從第 N 步推。
- [x] 3.2 建立 `steps/step4/index.ts` barrel re-export。**行為**：`from 'steps/step4'` 可取得 Step4Database。**驗證**：`tsc --noEmit` 通過。
- [x] 3.3 修改 `shell/WorkbenchV2Shell.tsx`：current=4 改 render `<Step4Database onCommit={...} />`；同時把 `state.committed` 透過 prop 傳給 TopBar 控制 scheduler 指示燈，落地 modified 需求 **Placeholder step and settings surfaces** 與新增需求 **Scheduler indicator reflects committed state**。**行為**：current=4 顯示 Step 4；TopBar scheduler 指示燈依 committed 切 idle/running。**驗證**：`shell.test.tsx` 更新並新增 it 測 indicator切換。
- [x] 3.4 修改 `shell/TopBar.tsx` 接 `scheduler: 'idle' | 'running'` prop；running 時 emerald + pulse-dot、idle 時 amber、文字依 i18n，落地設計決策「Scheduler 指示燈與 committed 連動」。**行為**：兩種狀態的 tone 與 class 對應。**驗證**：`shell.test.tsx` 中 `describe('TopBar scheduler')` 兩個 it 對應兩種狀態。
- [x] 3.5 在 `steps/Step4DatabasePlaceholder.tsx` 加 rollback only 註解。**驗證**：`grep -rn "Step4DatabasePlaceholder" frontend/src` 只在自身或註解出現。

## 4. i18n

- [x] 4.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `step4.*` keys：kind 名稱（SQLite / PostgreSQL / MySQL / SQL Server）、column 名（Tag / 點位 / 資料表欄位 / 欄位型態 / 啟用）、寫入策略 radio（INSERT 時序追加 / UPSERT 依時間戳）、commit log labels（10 個 API 名與 detail）、emerald 完成卡文案（設定已套用 · 開始收集資料 / Scheduler 已啟動 · 第一筆資料預計在 ~{N}s 後寫入 / 前往 Runtime Dashboard）、衝突 banner 文字、scheduler idle/running 文字。**行為**：覆蓋。**驗證**：`step4-database.test.tsx` 在 zh-TW 斷言 `資料庫連接器` 出現。
- [x] 4.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應翻譯（Database Connector / Configure Database Write / Commit and Start Scheduler / Setting Applied · Collecting Data / Scheduler Running · First record expected in ~{N}s / Go to Runtime Dashboard 等）。**行為**：對等。**驗證**：`step4-database.test.tsx` 在 en 斷言 `Database Connector` 出現。

## 5. 驗證

- [x] 5.1 跑 `cd frontend && npm run lint`，零錯誤。
- [x] 5.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [x] 5.3 跑 `cd frontend && npm run test -- --run workbench-v2`，全綠。
- [x] 5.4 跑 `cd frontend && npm run build`，產出乾淨。
- [x] 5.5 跑 `bash scripts/check_file_lines.sh`，確認 step4 所有新檔 ≤ 300 行。
- [x] 5.6 手動 smoke：完整跑「進入 Step 4 → 切 kind 看 schema 換 → 改某列 column 製造衝突 → 改回 → 點 Commit 看 10 個 log 動畫 → 完成後看 emerald 卡 + 頂列 scheduler 切 running → 切到 Step 3 看 form 已 commit → 切回 Step 4 看 success card 仍在」。mmitted state**。**行為**：current=4 顯示 Step 4；TopBar scheduler 指示燈依 committed 切 idle/running。**驗證**：`shell.test.tsx` 更新並新增 it 測 indicator 切換。
- [ ] 3.4 修改 `shell/TopBar.tsx` 接 `scheduler: 'idle' | 'running'` prop；running 時 emerald + pulse-dot、idle 時 amber、文字依 i18n，落地設計決策「Scheduler 指示燈與 committed 連動」。**行為**：兩種狀態的 tone 與 class 對應。**驗證**：`shell.test.tsx` 中 `describe('TopBar scheduler')` 兩個 it 對應兩種狀態。
- [ ] 3.5 在 `steps/Step4DatabasePlaceholder.tsx` 加 rollback only 註解。**驗證**：`grep -rn "Step4DatabasePlaceholder" frontend/src` 只在自身或註解出現。

## 4. i18n

- [ ] 4.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `step4.*` keys：kind 名稱（SQLite / PostgreSQL / MySQL / SQL Server）、column 名（Tag / 點位 / 資料表欄位 / 欄位型態 / 啟用）、寫入策略 radio（INSERT 時序追加 / UPSERT 依時間戳）、commit log labels（10 個 API 名與 detail）、emerald 完成卡文案（設定已套用 · 開始收集資料 / Scheduler 已啟動 · 第一筆資料預計在 ~{N}s 後寫入 / 前往 Runtime Dashboard）、衝突 banner 文字、scheduler idle/running 文字。**行為**：覆蓋。**驗證**：`step4-database.test.tsx` 在 zh-TW 斷言 `資料庫連接器` 出現。
- [ ] 4.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應翻譯（Database Connector / Configure Database Write / Commit and Start Scheduler / Setting Applied · Collecting Data / Scheduler Running · First record expected in ~{N}s / Go to Runtime Dashboard 等）。**行為**：對等。**驗證**：`step4-database.test.tsx` 在 en 斷言 `Database Connector` 出現。

## 5. 驗證

- [ ] 5.1 跑 `cd frontend && npm run lint`，零錯誤。
- [ ] 5.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [ ] 5.3 跑 `cd frontend && npm run test -- --run workbench-v2`，全綠。
- [ ] 5.4 跑 `cd frontend && npm run build`，產出乾淨。
- [ ] 5.5 跑 `bash scripts/check_file_lines.sh`，確認 step4 所有新檔 ≤ 300 行。
- [ ] 5.6 手動 smoke：完整跑「進入 Step 4 → 切 kind 看 schema 換 → 改某列 column 製造衝突 → 改回 → 點 Commit 看 10 個 log 動畫 → 完成後看 emerald 卡 + 頂列 scheduler 切 running → 切到 Step 3 看 form 已 commit → 切回 Step 4 看 success card 仍在」。
