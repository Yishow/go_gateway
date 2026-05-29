## Why

`datalink-workbench-v2-shell` 與 step1/step2/step3 四個 change 已落地 shell + 前三步驟。本 change 把 Step 4「資料庫寫入 + Commit」的完整靜態 UI 移植進 v2 容器，取代 `Step4DatabasePlaceholder`，完成 phase 1 主流程的最後一步。

Step 4 同時扮演兩個角色：

- **DB Connector + Tag 欄位映射**：選擇 connector kind（SQLite / PostgreSQL / MySQL / SQL Server）、設定 host/port/database/schema/table、寫入策略（INSERT 時序追加 / UPSERT 依時間戳）、寫入間隔；然後把每個 enabled Tag 對應到資料表欄位（含自動匹配演算法）。
- **整套設定 Commit**：底部「提交並啟動排程器」按鈕觸發一個 10 步驟原子寫入動畫（每 280ms 推進一個 log），完成後顯示綠色大勾並提示「前往 Runtime Dashboard」。Commit 後切換 `state.committed === true`，shell 頂列 scheduler 指示燈改為 `running`。

本 change 仍不接後端；10 個 commit log 是純動畫示意，phase 2 會接到真實 `POST /devices` × N、`POST /source-rules`、`POST /scheduler/start` 等鏈式呼叫。

## What Changes

- 新增 Step 4 完整 UI 於 `frontend/src/features/datalink/workbench-v2/steps/step4/`，三層版面：
  - **第一層** Connector（`SectionCard` col-span-12）：
    - 標題「資料庫連接器 (Connector)」+ icon `db` + aside success chip「已連線 · {host}」。
    - 12-col grid：4 個 kind 卡片選擇器（emoji `🗄️/🐘/🐬/🪟` + 名稱）+ 連線名稱 / Host / Port / Database / Schema / Table 一字排開。
    - 寫入策略 border-t 區塊：radio `INSERT (時序追加)` / `UPSERT (依時間戳)` + 「寫入間隔」number + 「秒」。
  - **第二層左 col 7** Tag → 欄位映射（`SectionCard` content p-0）：
    - 標題 `Tag → 資料表欄位` + 副標 `寫入 {schema}.{table}` + icon `table`。
    - aside：衝突時 error chip「欄位衝突」否則 success chip「{N} 對應」。
    - 表格欄：Tag（tag_key mono + display_name）+ 點位（name + @address）+ → + 資料表欄位（select dropdown，衝突時紅框）+ 欄位型態（SQL data_type）+ 啟用（Toggle）。
    - 表尾紅 banner（衝突時）。
  - **第二層右 col 5** 提交摘要（`SectionCard`）：
    - 未開始：5 列摘要卡（新增設備 / 接入規則 / 點位 / 標籤映射 / 資料庫寫入），每列左 label + 大字 value mono + 右 sub 文字 10px slate-500。
    - 寬版 success 按鈕「提交並啟動排程器」（disabled if 衝突或無啟用 target）。
    - 底部 info：「將執行 10 個 API 呼叫並啟動 collector」。
  - **提交動畫**：每 280ms 推進一個 log，sweep-in 動畫；每列含 emerald 勾 + label（含 method + path × N）+ detail + 右側 `200` emerald 小字。
  - **完成後**：emerald 大卡含圓勾 icon + 「設定已套用 · 開始收集資料」+ 「Scheduler 已啟動 · 第一筆資料預計在 ~{interval}s 後寫入」+ 「前往 Runtime Dashboard」secondary button（`onCommit` callback）。
- **自動欄位匹配演算法**：useEffect 依 `connector.kind` 與 `enabledPoints.length` 觸發；對每個 enabled point 嘗試三段比對：(1) `tagShort === columnName`、(2) `columnName.endsWith('_'+tagShort)` 或 `startsWith(tagShort+'_')`、(3) 退到 `columnNames[i % length]` 並避開已用；確保不會自動產生衝突。
- **衝突偵測**：對 `state.db.targets`（enabled + column_name 非空）做 column 用次統計；衝突列紅框 + alert icon + 表尾 banner + Commit 按鈕 disabled。
- **Sample table schema**：在 `state/dbSchemas.ts` 落地原型 `SAMPLE_DB_TABLES`（PostgreSQL `public.sensor_readings` 9 欄）；未來 phase 2 換成從後端拉真實 schema。
- **10 步 Commit 動畫**：reducer action `startCommit()` 啟動；component 內 timer 每 280ms dispatch `appendCommitLog`；完成後 dispatch `completeCommit` 同時把 `committed=true`、`commit.status='success'`、頂列 scheduler 切到 running。
- 擴充 `useWorkbenchV2State` reducer 加 db actions（`updateDbConnector`、`upsertDbTarget`、`updateDbTarget`、`autoAssignDbTargets`、`startCommit`、`appendCommitLog`、`completeCommit`、`resetCommit`）。
- 修改 `WorkbenchV2Shell` 在 `current === 4` 時改 render `Step4Database`；頂列 scheduler 指示燈依 `state.committed` 切換 `running` / `idle` 顏色；shell spec MODIFIED placeholder 從名單移除 Step 4。
- 擴充 i18n 加入 Step 4 microcopy（kind 名稱、column 名、commit log labels、emerald 完成卡訊息、寫入策略 radio）。

## Non-Goals (optional)

- 不接 `POST /db-connectors`、`POST /db-targets`、`POST /scheduler/start` 真實後端；10 個 commit log 是純動畫示意。
- 不變更 `database-target-workbench` spec（既有 spec 保持不動）。
- 不處理 schema 探測（從後端拉 connector 真實 schema），仍用 `SAMPLE_DB_TABLES`。
- 不支援 SQLite 檔案路徑欄位（原型用 emoji icon 顯示但 form 仍是 host/port/db，留給 backend-wiring 處理真實 SQLite UX）。
- 不處理 commit 失敗 retry（mock 不會失敗）。
- 不變更 Step 1/2/3/Settings 行為。

## Capabilities

### New Capabilities

- `datalink-workbench-v2-step4-database`：定義 v2 Step 4 資料庫寫入與 Commit 工作區的需求、connector 欄位契約、Tag → column 自動匹配演算法、衝突偵測規則、10 步驟 commit 動畫序列、完成卡與 scheduler 狀態切換。

### Modified Capabilities

- `datalink-workbench-v2-shell`：縮窄 `Placeholder step and settings surfaces` 需求，把 Step 4 從 placeholder 名單移除（剩餘 Settings 仍走 placeholder）；同時新增 `committed === true` 時頂列 scheduler 指示燈切到 `running` 行為。

## Impact

- Affected specs:
  - 新增 `openspec/specs/datalink-workbench-v2-step4-database/spec.md`
  - 修改 `openspec/specs/datalink-workbench-v2-shell/spec.md`
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/KindSelector.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/WriteStrategy.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/CommitProgress.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/index.ts`
    - `frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts`
    - `frontend/src/features/datalink/workbench-v2/state/commitLog.ts`
    - `frontend/src/features/datalink/workbench-v2/state/autoAssignTargets.ts`
    - `frontend/tests/workbench-v2/step4-database.test.tsx`
    - `frontend/tests/workbench-v2/step4-commit.test.tsx`
    - `frontend/tests/workbench-v2/autoAssignTargets.test.ts`
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`（新增 db actions）
    - `frontend/src/features/datalink/workbench-v2/state/types.ts`（補完 `DbConnector`、`DbTarget`、`CommitLog`、`CommitState` 型別）
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`（current=4 切到 Step4Database；scheduler 指示燈依 committed 切換）
    - `frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx`（scheduler 指示燈 running tone）
    - `frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx`（標 rollback only）
    - `frontend/src/i18n/locales/zh-TW/workbench-v2.json`、`frontend/src/i18n/locales/en/workbench-v2.json`
- 不變更後端任何檔案。
