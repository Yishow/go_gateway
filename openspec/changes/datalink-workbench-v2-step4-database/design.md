## Context

Step 4 是整個 4 步驟流程的終點，使用者在這頁完成最後設定並按下「提交」啟動排程器。原型 `step4-database.jsx`（339 行）有兩個關鍵互動：(1) 自動欄位匹配演算法須在 connector kind 切換或 enabled point 數量變化時自動執行；(2) 10 步 commit 動畫須與 shell 頂列 scheduler 指示燈連動。

由於 Step 4 是流程終點，UI 還需要在 commit 完成後顯示綠色 emerald 大卡與「前往 Runtime Dashboard」按鈕，並把 form 變只讀（commit 後不允許再改）。

## Goals / Non-Goals

**Goals:**

- 1:1 視覺與互動復刻原型 Step 4。
- 自動欄位匹配演算法抽到純函式，便於單元測試。
- 10 步 commit 動畫由 reducer 管理 log state、component 管理 timer；timer cleanup 完整。
- 完成 commit 後 form 變只讀，避免使用者誤改。
- shell 頂列 scheduler 指示燈依 `state.committed` 即時切換 tone（emerald running / amber idle）。

**Non-Goals:**

- 不接後端。
- 不支援 commit 失敗 / retry。
- 不處理 schema 探測；用 `SAMPLE_DB_TABLES` 靜態資料。
- 不變更 `database-target-workbench` spec。
- 不處理 SQLite 檔案路徑欄位 UX。

## Decisions

### 拆檔策略：8 個元件 + 3 個 state module

原型 339 行拆為：

- `Step4Database.tsx`（~180 行）：容器，state derivation + auto-assign effect + commit timer + onCommit callback。
- `KindSelector.tsx`（~80 行）：4 個 db kind 卡片。
- `WriteStrategy.tsx`（~80 行）：INSERT/UPSERT radio + 寫入間隔 number。
- `ConnectorSection.tsx`（~150 行）：第一層 SectionCard 容器，組合 KindSelector + 連線欄位 + WriteStrategy。
- `TargetMappingTable.tsx`（~200 行）：第二層左，含表格 + 衝突 banner。
- `CommitSummary.tsx`（~120 行）：第二層右未開始狀態（5 列摘要 + 按鈕）。
- `CommitProgress.tsx`（~120 行）：第二層右進行中狀態（動畫 log 列表）。
- `CommitSuccessCard.tsx`（~80 行）：第二層右完成狀態（emerald 大卡 + 前往 Dashboard 按鈕）。

state module：
- `state/dbSchemas.ts`（~100 行）：`SAMPLE_DB_TABLES`（PostgreSQL `sensor_readings` 9 欄）、`getColumnsFor(kind)`、`getDefaultConnector(kind)`。
- `state/autoAssignTargets.ts`（~100 行）：純函式 `autoAssignTargets(enabledPoints, mappings, columnNames, existingTargets)` 回傳新 targets dict。
- `state/commitLog.ts`（~100 行）：`buildCommitLogSequence(state)` 回傳 10 個 log entry。

### 自動欄位匹配演算法

`autoAssignTargets(enabledPoints, mappings, columnNames, existingTargets)`：

```ts
function autoAssignTargets(
  enabledPoints: Point[],
  mappings: Record<string, Mapping>,
  columnNames: string[],
  existingTargets: Record<string, DbTarget>,
): Record<string, DbTarget> {
  const used = new Set<string>();
  const out: Record<string, DbTarget> = {};
  for (let i = 0; i < enabledPoints.length; i++) {
    const p = enabledPoints[i];
    const m = mappings[p.id];
    if (!m) continue;

    if (existingTargets[p.id]) {
      out[p.id] = existingTargets[p.id];
      used.add(existingTargets[p.id].column_name);
      continue;
    }

    const tagShort = m.tag_key.split('.').pop() ?? `col_${i + 1}`;
    const exact = columnNames.find(c =>
      !used.has(c) && (c === tagShort || c.endsWith('_' + tagShort) || c.startsWith(tagShort + '_'))
    );
    const fallback = columnNames[i % columnNames.length];
    const auto = exact ?? (used.has(fallback) ? columnNames.find(c => !used.has(c)) : fallback) ?? fallback;
    used.add(auto);
    out[p.id] = {
      tag_id: `tag.${m.tag_key}`,
      column_name: auto,
      enabled: true,
    };
  }
  return out;
}
```

優先級：existing > exact match > index fallback (skip if used) > index fallback (fall through)。確保 8 個點位剛好對應 8 個欄位、無自動衝突。

**Alternatives considered**：(A) 沿用原型 useEffect inline 邏輯 — 難以單元測試；(B) 用 Hungarian algorithm 求最佳分配 — 過度工程，當前 8/8 場景手動 fallback 足夠。

### Commit log 序列：純函式 + reducer 串聯

`buildCommitLogSequence(state)` 回傳：

```ts
type CommitLog = {
  label: string;          // 含 method + path + count，例如 'POST /devices × 2'
  detail: string;         // 例如 'PLC-生產線-01 (modbus_tcp) · PLC-02 (mc_3e)'
  status_code: 200;
};

[
  { label: `POST /devices × ${devices.length}`, detail: devices.map(...).join(' · ') },
  { label: `POST /devices/:id/activate × ${devices.length}`, detail: 'draft → active' },
  { label: `POST /source-rules × ${enabledRules.length}`, detail: rules.map(...).join(' · ') },
  { label: `POST /points × ${enabledPoints.length}`, detail: 'bulk create' },
  { label: 'POST /polling-groups', detail: '快速輪詢 1s, enabled' },
  { label: `POST /tags × ${enabledPoints.length}`, detail: 'register tag keys' },
  { label: `POST /mappings × ${enabledPoints.length}`, detail: 'point ↔ tag, scale pipeline' },
  { label: 'POST /db-connectors/:id/test', detail: `${connector.kind} ${connector.host}:${connector.port}` },
  { label: `POST /db-targets × ${enabledTargetCount}`, detail: `→ ${connector.schema}.${connector.table}` },
  { label: 'POST /scheduler/start', detail: 'collectors started' },
]
```

Component 內 `useEffect` 監聽 `state.commit.status === 'running'`：

```ts
useEffect(() => {
  if (state.commit?.status !== 'running') return;
  const seq = buildCommitLogSequence(state);
  let i = state.commit.logs?.length ?? 0;
  if (i >= seq.length) {
    dispatch({ type: 'completeCommit' });
    return;
  }
  const handle = setTimeout(() => {
    dispatch({ type: 'appendCommitLog', log: seq[i] });
  }, 280);
  return () => clearTimeout(handle);
}, [state.commit?.status, state.commit?.logs?.length]);
```

優點：每個 log 是獨立 reducer dispatch，timer 自動配合 React 重渲染節奏；unmount 時 cleanup 取消 pending tick。

**Alternatives considered**：(A) component 內 setTimeout 鏈式遞迴（原型作法）— 難 cleanup、難測試；(B) 用 RxJS observable — 引入新依賴。

### Commit 後 form 只讀

`state.committed === true` 時：

- `KindSelector` 四個卡片 disabled、不允許切換。
- `ConnectorSection` 所有 input/select disabled。
- `WriteStrategy` radio 與 number disabled。
- `TargetMappingTable` 每列 select 與 toggle disabled。
- `CommitSummary` 不顯示；`CommitSuccessCard` 顯示。

把這個邏輯抽到 hook `useStep4Readonly()` 回傳 boolean，每個子元件用 prop drilling 或 context 取。

**Alternatives considered**：(A) 不變唯讀讓使用者可重新提交 — 邏輯複雜（已 committed 後修改要怎麼處理？）；(B) commit 後直接跳到 Runtime Dashboard 不留在 Step 4 — 失去「commit 完成可看 log 回顧」價值。

### Scheduler 指示燈與 committed 連動

`WorkbenchV2Shell` 內 `TopBar` 已渲染 scheduler 狀態文字。本 change 擴充 `TopBar` 接 `scheduler` prop（`'idle' | 'running'`），依 `state.committed` 切：

- `idle`：amber dot + `scheduler idle`。
- `running`：emerald dot + pulse-dot animation + `scheduler running`。

切換時機是 `state.committed` reducer 變 true 即同步生效。

### Sample table schema

`state/dbSchemas.ts` 含 `SAMPLE_DB_TABLES.postgres[0]`：

| Column            | data_type           | nullable | primary_key |
| ----------------- | ------------------- | -------- | ----------- |
| ts                | timestamptz         | no       | yes         |
| temp_in_c         | double precision    | yes      | no          |
| temp_out_c        | double precision    | yes      | no          |
| pressure_main_kpa | double precision    | yes      | no          |
| pressure_sub_kpa  | double precision    | yes      | no          |
| flow_lpm          | double precision    | yes      | no          |
| humidity_pct      | double precision    | yes      | no          |
| vibration_mms     | double precision    | yes      | no          |
| motor_rpm         | integer             | yes      | no          |

`getColumnsFor(kind)` 在 sqlite/mysql/sqlserver 也回傳上述 schema（共用）；phase 2 會替換成真實 schema 探測。`getDefaultConnector(kind)` 回傳對應 kind 預設連線參數。

## Implementation Contract

#### Behavior

- 進入 `current === 4`：connector 預設 PostgreSQL `tsdb.internal:5432` / database `gateway_metrics` / schema `public` / table `sensor_readings` / write_mode `insert` / interval 5 秒；右側 5 列摘要顯示 device/rule/point/mapping/target 數字。
- 切 connector kind 卡片：連線欄位保留 host/port（不重置），但 sample schema 透過 `getColumnsFor(kind)` 換；觸發 auto-assign 重算 targets。
- 改 schema / table：摘要文字 `→ {schema}.{table}` 同步更新。
- 改 write_mode radio：state 更新。
- 改 write_interval：summary card 「資料庫寫入」sub 文字 `{kind} → {schema}.{table}` 不變但 commit 後 emerald 卡訊息 `~{N}s` 反映新值。
- TargetMappingTable：每列 select 顯示 columnNames（非 primary_key）；初始 auto-assign 後每筆 mapping 對應一個 column；改某筆 column 變成已用值 → 該列紅框 + alert icon + 表尾紅 banner + Commit 按鈕 disabled。
- 點 enabled toggle：state.db.targets[p.id].enabled 切換；衝突 / 啟用數 / commit 按鈕 enabled 條件即時重算。
- 點「提交並啟動排程器」：reducer dispatch `startCommit`、`state.commit.status='running'`；component effect 啟動 10 步 tick 每 280ms 推 1 個 log；按鈕 disabled、icon 換 refresh 旋轉；CommitSummary 隱藏、CommitProgress 顯示。
- log 推進中：每筆 sweep-in 動畫；右側 200 emerald 小字；底部 `正在執行下一個指令…` pulse-dot。
- 10 個 log 全部完成：dispatch `completeCommit`、`state.committed=true`、`state.commit.status='success'`、`state.commit.finished_at=now`；CommitProgress 變只列出最終 log、CommitSuccessCard 顯示 emerald 大卡；form 全部變 disabled；頂列 scheduler 指示燈切到 `running`（emerald + pulse-dot）。
- 點「前往 Runtime Dashboard」secondary button：呼叫 `onCommit` callback；本 change 留為 noop（show alert 或 noop），由後續 change 接 Runtime Dashboard 路由。
- 切到其他 step（StepRail）後再切回 Step 4：`state.commit.logs` 保留、CommitSuccessCard 仍顯示；不重跑 commit。

#### Interface / Data Shape

- `DbConnector`：
  ```ts
  type DbConnector = {
    name: string;
    kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
    host: string;
    port: number;
    database: string;
    username: string;
    schema: string;
    table: string;
    write_mode: 'insert' | 'upsert';
    write_interval_seconds: number;
    timestamp_column: string;
    status: 'unknown' | 'testing' | 'ready' | 'unreachable';
  };
  ```
- `DbTarget`：
  ```ts
  type DbTarget = {
    tag_id: string;
    column_name: string;
    enabled: boolean;
  };
  ```
- `CommitLog`：`{ label: string; detail: string; status_code: 200 }`
- `CommitState`：`{ status: 'idle' | 'running' | 'success'; logs: CommitLog[]; started_at?: string; finished_at?: string }`
- Reducer db actions：
  - `{ type: 'updateDbConnector'; patch: Partial<DbConnector> }`
  - `{ type: 'upsertDbTarget'; pointId: string; target: DbTarget }`
  - `{ type: 'updateDbTarget'; pointId: string; patch: Partial<DbTarget> }`
  - `{ type: 'autoAssignDbTargets'; targets: Record<string, DbTarget> }`
  - `{ type: 'startCommit' }`
  - `{ type: 'appendCommitLog'; log: CommitLog }`
  - `{ type: 'completeCommit' }`
  - `{ type: 'resetCommit' }`（dev-only，給 Tweaks panel「重置流程」用）

#### Failure Modes

- mock commit 不會失敗；若使用者在 commit 進行中切到其他 step，timer 因 useEffect cleanup 取消，回到 Step 4 時 commit 仍處於 running 但 log 數量未滿 10，effect 重新啟動 tick 繼續推。
- 自動 auto-assign 在 enabled point 與 column 數量不相等時：若 point > column，多出來的 point 依 `i % length` 環繞造成衝突，使用者需手動調整；若 point < column，多餘 column 留白；衝突 banner 與 disabled 按鈕引導使用者解決。
- mapping 不存在的 point（極端：Step 3 還沒 init 完成就跳到 Step 4）：autoAssign 跳過該 point，console.warn 一次；UI 列數 = mappings 中有 entry 的 point 數。
- 重新點擊「前往 Runtime Dashboard」：本 change 為 noop，不會多次觸發 onCommit。

#### Acceptance Criteria

- `frontend/tests/workbench-v2/autoAssignTargets.test.ts` 覆蓋：
  - 8 enabledPoints + 9 columnNames（postgres sample）→ 每個 point 取得一個 column，無衝突。
  - 改 row 0 tag_key 為 `flow_lpm` → exact match `flow_lpm` 欄。
  - existingTargets 內已有 entry → 保留不覆寫。
  - 8 points + 4 columns → 前 4 個 exact 或 fallback、後 4 個 index wrap 造成衝突（測試 expect 衝突存在）。
- `frontend/tests/workbench-v2/step4-database.test.tsx` 覆蓋：
  - 預設 render：connector 預設值正確、表格 8 列無衝突、Commit 按鈕啟用。
  - 切到 mysql → kind state 更新、按鈕仍啟用。
  - 改某列 column 變成已用值 → 紅框 + alert + 表尾 banner + Commit 按鈕 disabled。
  - 點 toggle 把某列 enabled=false → 衝突重算、按鈕 enabled 條件更新。
  - 寫入策略改 UPSERT → state.write_mode 更新。
- `frontend/tests/workbench-v2/step4-commit.test.tsx` 覆蓋：
  - 點 Commit 後 reducer state.commit.status === 'running'、按鈕 disabled。
  - 用 `vi.useFakeTimers()` 推 10 × 280ms → 10 個 log 出現、state.committed === true、state.commit.status === 'success'、form 變 disabled。
  - shell 頂列 scheduler 指示燈在 commit 完成後顯示 `running`。
  - 在 commit 進行中（推 3 × 280ms 後）unmount → 重新 mount → state.commit.logs 仍是 3 筆、effect 繼續從第 4 步推。
  - 點 emerald 卡內「前往 Runtime Dashboard」→ `onCommit` 被呼叫一次。
- `cd frontend && npm run lint && npm run typecheck && npm run test -- --run workbench-v2 && npm run build` 全綠。
- 手動 smoke：完整跑「進入 Step 4 → 切 kind 看 schema 換 → 改某列 column 製造衝突 → 改回 → 點 Commit 看 10 個 log 動畫 → 完成後看 emerald 卡 + 頂列 scheduler running → 切到 Step 3 看 form 已 commit」。

#### Scope Boundaries

**In scope:**

- Step 4 完整 UI（8 個元件）。
- 3 個 state module（dbSchemas / autoAssignTargets / commitLog）。
- reducer 8 個 db actions。
- `state/types.ts` 補完 DbConnector / DbTarget / CommitLog / CommitState 型別。
- shell `current === 4` 改 render Step4Database；TopBar scheduler 指示燈依 committed 切換；shell spec 縮窄 placeholder 與新增 scheduler indicator 行為。
- i18n `workbench-v2` 加入 Step 4 microcopy。
- unit test。

**Out of scope:**

- 後端 API、SSE。
- commit 失敗 retry。
- schema 探測。
- SQLite 檔案路徑 UX。
- Runtime Dashboard 路由（onCommit 為 noop）。
- E2E。
- Step 1/2/3/Settings。

## Risks / Trade-offs

- [auto-assign 在 connector kind 切換時觸發，但 user 已手動調整過某列 → existingTargets 保留邏輯確保不會被覆寫] → 同時用 useEffect 依賴 `[connector.kind, enabledPoints.length]` 而非 `enabledPoints`，避免 mapping 內部變動觸發重算。
- [commit 動畫 effect 依賴 `state.commit?.logs?.length` 可能在 React 18 strict mode 開發環境 double-fire effect] → 用 `dispatch` 是 idempotent；append 重複 log 會被 reducer 用「if logs.length === seq.length 不再 append」防護。
- [完成後 form 變 disabled 但使用者切回 Step 3 改 mapping 後切回 Step 4，UI 還是 disabled] → 預期行為，commit 已完成；若要重做需 Tweaks panel「重置流程」（dispatch `resetCommit` + 其他 reset）。
- [10 步驟總共 ~2.8 秒動畫，使用者可能覺得慢] → 每 280ms 對應原型節奏，保留視覺節奏；可在 Tweaks panel 加速度倍率 slider（留給後續 change）。

## Migration Plan

1. 在 `state/types.ts` 補完 DbConnector / DbTarget / CommitLog / CommitState 型別。
2. 建立 `state/dbSchemas.ts`，TDD：dbSchemas.test.ts 對 PostgreSQL sample 與 getColumnsFor 各 case。
3. 建立 `state/autoAssignTargets.ts`，TDD：autoAssignTargets.test.ts 對演算法各分支。
4. 建立 `state/commitLog.ts`，TDD：commitLog.test.ts 對 buildCommitLogSequence 與 state 變化的對應。
5. 擴充 reducer 加 8 個 db actions + reducer-step4.test.ts。
6. 從子元件實作：`KindSelector` → `WriteStrategy` → `ConnectorSection`。
7. `TargetMappingTable`。
8. `CommitSummary` → `CommitProgress` → `CommitSuccessCard`。
9. 組裝 `Step4Database.tsx`，覆寫 shell `current === 4`；TopBar 接 scheduler prop。
10. 補 i18n。
11. 跑 lint/typecheck/test/build/check-lines。
12. 手動 smoke。

**Rollback strategy**：不刪除 placeholder；rollback 把 shell `current === 4` 改回 placeholder、移除 step4 目錄、TopBar scheduler indicator 改回固定 idle。

## Open Questions

- emerald 完成卡的「前往 Runtime Dashboard」按鈕在本 change 暫為 noop（show alert 或 console.log），未來 Runtime Dashboard 何時上線？傾向後續獨立 change，phase 2 backend-wiring 完成後評估。
- mock commit 是否要在某機率隨機 fail 來測 fail UI？傾向 **不**，與 step1 一致；fail 路徑留給後續 change（backend-wiring 完成後真實 fail 才實作 UI）。
