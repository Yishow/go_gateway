## Context

Settings 是獨立於 4 步驟主流程之外的系統級頁面，由 5 個彼此獨立但相互關聯的區塊組成。原型 `settings.jsx`（374 行）幾乎沒有複雜互動邏輯，主要是長表單與 connector 列表的 CRUD + mock 測試。

關鍵約束：

- Settings state 已部分由前面 change 使用（Step 2 讀 `modbus_share`、Step 4 讀 `connectors`）；本 change 必須確保寫入路徑與既有讀取契約相容。
- Connector mock 測試是 setTimeout 動畫；多個 connector 可同時測試（每筆獨立的 status）；UI 上 status chip 即時反映。
- 底部 sticky 儲存列在頁面捲動時始終可見，但本 change 的「儲存」按鈕為 noop。

## Goals / Non-Goals

**Goals:**

- 1:1 視覺與互動復刻原型 Settings。
- 把 settings 預設值集中於 `state/settingsDefaults.ts`。
- 每個區塊一個元件，方便後續 phase 2 單獨接 backend endpoint。
- Connector mock 測試與 Step 1 device 測試共用「兩階段 status 切換 + setTimeout」模式，但 stage 結構不同（connector 沒有分階段，只有最後成功/失敗）。

**Non-Goals:**

- 不接後端。
- 不處理設定匯入 / 匯出。
- 不實作 light mode 主題切換。
- 不處理語言切換即時 reload。
- 不變更其他 step 行為。

## Decisions

### 拆檔策略：12 個元件 + 1 個 state module

原型 374 行拆為：

- `SettingsPage.tsx`（~150 行）：容器，組合所有區塊。
- `SettingsHeader.tsx`（~50 行）：標題列漸層卡。
- `ConnectorPoolSection.tsx`（~120 行）：第 1 個 SectionCard 容器 + 空狀態 + 列表。
- `ConnectorRow.tsx`（~200 行）：單一 connector 列。
- `TimeseriesSection.tsx`（~80 行）：第 2 區塊。
- `SchedulerSection.tsx`（~80 行）：第 3 區塊。
- `ModbusShareSection.tsx`（~120 行）：第 4 區塊（總開關 + 4-col grid）。
- `UiSection.tsx`（~80 行）：第 5 區塊左（介面）。
- `ApiSection.tsx`（~80 行）：第 5 區塊中（API 端點）。
- `DiagnosticsSection.tsx`（~80 行）：第 5 區塊右（診斷）。
- `SaveBar.tsx`（~60 行）：sticky 底部儲存列。
- `index.ts` barrel。

state module：
- `state/settingsDefaults.ts`（~100 行）：`DEFAULT_SETTINGS`、`makeDefaultConnector(idx)`、`makeDefaultSection(key)`。

### Connector 池：以陣列管理 + per-row status

每個 connector 結構（與 Step 4 共用 type 的擴充）：

```ts
type SettingsConnector = {
  id: string;                     // conn-xxx
  name: string;
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  host: string;
  port: number;
  database: string;
  username: string;
  enabled: boolean;
  status: 'unknown' | 'testing' | 'ready' | 'unreachable' | 'auth_failed';
  last_check_at?: string;        // ISO timestamp
  last_check_error?: string;
  default_write_interval_seconds: number;
};
```

新增、測試、刪除等動作走 reducer：

- `addConnector()`：append 一筆 default connector，name `新連線 N`，kind 預設 `postgres`，status `unknown`。
- `updateConnector(id, patch)`：merge 欄位變更；改 kind/host/port 時 status 自動回 `unknown`。
- `removeConnector(id)`：filter 移除；不連動 Step 4（Step 4 的 connector 是本 change 完成前的內建，後續 backend-wiring 才會切到從 pool 選）。
- `startConnectorTest(id)`：把該 connector status 設為 `testing`。
- `completeConnectorTest(id, result)`：把 status 設為 `ready` 或 `unreachable`，更新 `last_check_at` 與 `last_check_error`。

Component 內 `testConnector(id)`：

```ts
function testConnector(id: string) {
  dispatch({ type: 'startConnectorTest', id });
  setTimeout(() => {
    const ok = Math.random() > 0.15;
    dispatch({
      type: 'completeConnectorTest',
      id,
      result: ok
        ? { status: 'ready', last_check_at: new Date().toISOString() }
        : { status: 'unreachable', last_check_at: new Date().toISOString(), last_check_error: 'connection refused' },
    });
  }, 900);
}
```

unmount 時不需 cleanup（不像 Step 1 / Step 4 推進多階段，這裡只是一次 setTimeout；若使用者離開頁面後測試 callback 仍會 fire，但 reducer dispatch 是 idempotent）。

**Alternatives considered**：(A) 真實 backend test — 違反不接後端邊界；(B) 永遠成功 — 失去「失敗 status chip 視覺驗證」。

### Settings 預設值集中管理

`state/settingsDefaults.ts` 匯出：

```ts
export const DEFAULT_SETTINGS: Settings = {
  connectors: [{
    id: 'conn-prod',
    name: 'TimeSeries Prod',
    kind: 'postgres',
    host: 'tsdb.internal',
    port: 5432,
    database: 'gateway_metrics',
    username: 'gw_writer',
    enabled: true,
    status: 'ready',
    last_check_at: '<initial ISO timestamp>',
    default_write_interval_seconds: 5,
  }],
  timeseries: { write_precision: 'millisecond', partition_interval: 'daily', batch_size: 500, retention_days: 90 },
  scheduler: { default_interval_ms: 1000, default_retry_count: 3, default_retry_delay_ms: 500, breaker_threshold: 10, auto_start: true },
  modbus_share: { enabled: true, bind_address: '0.0.0.0', port: 5020, slave_id: 1, base_register: 40001 },
  general: { theme: 'dark', locale: 'zh-TW', addr_format: 'modbus', api_base: 'http://localhost:8080', api_version: 'v1', timeout_seconds: 30, log_level: 'info', sse_heartbeat_seconds: 15, enable_debug_panel: false, enable_audit_log: true },
};
```

注意 `modbus_share.enabled` 預設為 `true`，與 step2 設計時的假設一致；否則 Step 2 share 預覽會看不到任何 share 位址。

### 「重設為預設」與「儲存所有設定」按鈕

- 「重設為預設」ghost button：點擊 dispatch `resetSettingsToDefaults`，把整個 `state.settings` 換回 `DEFAULT_SETTINGS`（包括 connector 池）。
- 「儲存所有設定」success button：本 change 為 noop（console.warn `pending backend-wiring`）；後續 backend-wiring change 才實作真實 PATCH。

UI 上「儲存所有設定」按鈕仍維持 success tone，不視覺上標示「未實作」，避免使用者誤判已壞掉；底部 info 文字「設定會立即套用，並於下次重啟後生效」誠實表達 mock 行為。

**Alternatives considered**：(A) 把按鈕在本 change 移除 — 失去視覺完整性；(B) 按鈕標 disabled — 違反「按鈕應該對應實際 action」原則。

### 連動效應

settings 變更後，Step 2 的 share 計算（透過 `useShareLayout` selector）會即時讀新值；Step 4 的 connector pool 本 change 落地後仍未取代 Step 4 內建 connector（那是 backend-wiring 的事），所以 Step 4 行為不變。

Step 2 已在前面 change 接好 selector，本 change 只負責提供寫入路徑；測試應覆蓋「在 Settings 改 base_register 後切到 Step 2 看 share 重算」。

### Sticky 儲存列定位

底部儲存列用 `sticky bottom-4 z-10` 配 `backdrop-blur`；確保頁面捲動時始終可見。z-10 不會蓋過 Tweaks panel（後者用 z-2147483646）。

## Implementation Contract

#### Behavior

- 進入 `view === 'settings'`：渲染標題列 + 5 個 SectionCard + 底部 sticky 儲存列；摘要欄（SummaryRail）隱藏（shell 既有規則）。
- Connector 池：預設 1 筆 `TimeSeries Prod`（postgres / tsdb.internal / 5432）status ready。
- 點「+ 新增連接器」：append 一筆 `新連線 N`，kind `postgres`，status `unknown`，預設 port 5432。
- 改 kind / host / port：status 自動回 `unknown`；status chip 顯示「未測試」。
- 點測試 button：button icon 換 refresh 旋轉、disabled；status chip 變 testing（藍）；900ms 後 85% 機率成功（emerald chip + `上次檢查 {time}`）/ 15% 失敗（red chip + `connection refused` 紅字）。
- 點刪除 button：直接 dispatch removeConnector（不 confirm，與原型一致）。
- 時序儲存策略：4 個欄位互不影響；onChange 即時寫入。
- 排程預設：5 個欄位含 auto_start toggle 跨欄。
- Modbus Share 總開關 toggle：切 enabled 後下方 4-col grid 顯示/隱藏；變更立即生效（Step 2 share 計算重跑）。
- Modbus Share base_register 改變：Step 2 share layout 透過 selector 即時重算。
- 介面 / API / 診斷：所有 Select 與 Toggle 即時寫入；本 change 不實作 theme 與 locale 即時生效。
- 「重設為預設」：彈出原生 confirm 後 dispatch resetSettingsToDefaults（confirm 取消則 noop）。
- 「儲存所有設定」：本 change 點擊只 console.warn；無 UI 變化。

#### Interface / Data Shape

- `Settings` 型別（state/types.ts 既有）：
  ```ts
  type Settings = {
    connectors: SettingsConnector[];
    timeseries: TimeseriesSettings;
    scheduler: SchedulerSettings;
    modbus_share: ModbusShareSettings;
    general: GeneralSettings;
  };

  type TimeseriesSettings = {
    write_precision: 'second' | 'millisecond';
    partition_interval: 'daily' | 'weekly' | 'monthly';
    batch_size: number;
    retention_days: number;
  };

  type SchedulerSettings = {
    default_interval_ms: number;
    default_retry_count: number;
    default_retry_delay_ms: number;
    breaker_threshold: number;
    auto_start: boolean;
  };

  type ModbusShareSettings = {
    enabled: boolean;
    bind_address: string;
    port: number;
    slave_id: number;
    base_register: number;
  };

  type GeneralSettings = {
    theme: 'dark' | 'light' | 'auto';
    locale: 'zh-TW' | 'en';
    addr_format: 'modbus' | 'hex' | 'raw';
    api_base: string;
    api_version: 'v1' | 'v2';
    timeout_seconds: number;
    log_level: 'trace' | 'debug' | 'info' | 'warn' | 'error';
    sse_heartbeat_seconds: number;
    enable_debug_panel: boolean;
    enable_audit_log: boolean;
  };
  ```
- Reducer settings actions：
  - `{ type: 'updateSettings'; patch: Partial<Settings> }`
  - `{ type: 'updateSettingsSection'; section: keyof Settings; patch: object }`
  - `{ type: 'addConnector' }`（不帶 payload；reducer 自己生 default）
  - `{ type: 'updateConnector'; id: string; patch: Partial<SettingsConnector> }`
  - `{ type: 'removeConnector'; id: string }`
  - `{ type: 'startConnectorTest'; id: string }`
  - `{ type: 'completeConnectorTest'; id: string; result: { status: 'ready' | 'unreachable' | 'auth_failed'; last_check_at: string; last_check_error?: string } }`
  - `{ type: 'resetSettingsToDefaults' }`

#### Failure Modes

- mock connector test 隨機失敗：85/15 比例；失敗的 status chip 顯示 `無法連線`（red tone）+ `connection refused` 紅字。
- 「儲存所有設定」按鈕本 change 為 noop：console.warn 一次「pending backend-wiring」；不顯示 toast / alert（避免使用者誤判）。
- 「重設為預設」沒按 confirm：reducer 不執行、state 不變。
- enabled toggle 對所有 connector 都 disabled：UI 仍允許 enable=true 但 Step 4 仍用內建 connector（本 change 不接 pool ↔ Step 4）。
- 第一次測試成功後切換 kind：status 回 `unknown`，但 `last_check_at` 保留（與 last_check_error 一起）；UI 上 `last_check_at` 只在 status === 'ready' 時顯示。

#### Acceptance Criteria

- `frontend/tests/workbench-v2/settings.test.tsx` 覆蓋：
  - 進入 settings：5 個 SectionCard + 標題列 + sticky save bar 全部 render。
  - Timeseries section：改 write_precision 從 ms 到 s、state 更新。
  - Scheduler section：改 default_interval_ms、state 更新；切 auto_start toggle。
  - Modbus Share section：切 enabled toggle → 4-col grid 顯隱；改 base_register。
  - UI section：切 theme select → state.general.theme 更新（不實際換主題）。
  - API section：改 api_base → state 更新。
  - Diagnostics section：切 enable_audit_log → state 更新。
  - 「重設為預設」點擊 → mock `window.confirm` 回 true → state 全部 reset；confirm 回 false → state 不變。
  - 「儲存所有設定」點擊 → console.warn 出現「pending」字串、state 不變。
- `frontend/tests/workbench-v2/settings-connectors.test.tsx` 覆蓋：
  - 預設 1 個 `TimeSeries Prod` connector ready。
  - 點「+ 新增連接器」→ connectors 變 2 個、第 2 個 status unknown。
  - 改 host → status 回 unknown。
  - 點測試 button → status 變 testing；用 `vi.useFakeTimers()` 推 900ms + mock `Math.random` 回 0.5（>0.15） → status 變 ready + last_check_at 存在；mock 回 0.1 → status 變 unreachable + last_check_error。
  - 點刪除 → connector 數量 -1（不 confirm）。
- `cd frontend && npm run lint && npm run typecheck && npm run test -- --run workbench-v2 && npm run build` 全綠。
- 手動 smoke：完整跑「進入 settings → 新增 connector → 測試 connector（看 mock 動畫）→ 改 Modbus Share base_register → 切到 Step 2 看 share 位址重算 → 切回 settings 點重設為預設 confirm → 全部回到初值」。

#### Scope Boundaries

**In scope:**

- Settings 完整 UI（12 個元件）。
- `state/settingsDefaults.ts` 預設值集中。
- reducer 8 個 settings actions。
- `state/types.ts` 補完所有 Settings 子型別。
- shell 換掉 SettingsPlaceholder；shell spec 清空 placeholder 殘留。
- i18n `workbench-v2` 加入 Settings microcopy（5 個區塊完整）。
- unit test。

**Out of scope:**

- 後端 API（PATCH /settings / POST /db-connectors/test）。
- 主題 light/dark 即時切換。
- 語言即時切換 reload。
- 設定匯入 / 匯出。
- E2E。
- pool → Step 4 connector 整合（屬於 backend-wiring）。

## Risks / Trade-offs

- [「儲存所有設定」按鈕為 noop 可能讓使用者誤判已存] → 底部 info 文字「設定會立即套用，並於下次重啟後生效」實際上前半截是真的（state 即時更新）、後半截是 mock；明確在 phase 1 完成後加註「Settings 變更立即影響 Step 2/4 預覽，PATCH 後端 endpoint 由 backend-wiring change 落地」。
- [connector 測試 mock 隨機失敗讓 e2e 不確定] → 在 e2e 中 mock `Math.random` 為固定值；reducer 測試也用相同手法。
- [`window.confirm` 在 e2e 中需 hook 處理] → Playwright `page.on('dialog')` 可處理；本 change 不寫 e2e，留給後續 e2e change。
- [重設為預設可能誤砍使用者已建好的多個 connector] → confirm dialog 提示「將清除所有 connector 與設定變更」。

## Migration Plan

1. 在 `state/types.ts` 補完 Settings 全部子型別。
2. 建立 `state/settingsDefaults.ts`，TDD：settingsDefaults.test.ts 覆蓋 DEFAULT_SETTINGS 與 makeDefaultConnector。
3. 擴充 reducer 加 8 個 settings actions + reducer-settings.test.ts。
4. 從子元件實作：`SettingsHeader` → `TimeseriesSection` → `SchedulerSection` → `ModbusShareSection` → `UiSection` / `ApiSection` / `DiagnosticsSection`。
5. `ConnectorRow` → `ConnectorPoolSection`。
6. `SaveBar`。
7. 組裝 `SettingsPage.tsx`，覆寫 shell `view === 'settings'`。
8. 補 i18n。
9. 跑 lint/typecheck/test/build/check-lines。
10. 手動 smoke。

**Rollback strategy**：不刪除 placeholder；rollback 把 shell `view === 'settings'` 改回 placeholder、移除 settings 目錄、reducer actions 變 dead code。

## Open Questions

- 是否要在 connector row 加「複製」按鈕？傾向 **不**，原型沒有；使用者通常只有少數 connector。
- 「重設為預設」要不要分區重設（只重設某一塊）？傾向 **不**，原型只有整體 reset；分區增加 UI 複雜度。
