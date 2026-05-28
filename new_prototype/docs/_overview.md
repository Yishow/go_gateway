# Datalink Workbench — 架構總覽

## 一、產品定位

Datalink Workbench 是 go_gateway 系統的設定流程主介面，把原本散落在不同頁面的 Device、SourceRule、Point、Tag、Mapping、DatabaseTarget 操作整併為**單一可推進的線性流程**，配上一個獨立的「設定」頁管理系統級設定。

完整流程：

```
新增裝置 ──→ 接入規則 ──→ 點位映射 ──→ 儲存資料庫
(多裝置)   (多規則·多範圍)  (Tag+轉換)   (Connector+欄位)
```

---

## 二、資料模型 (前端 state shape)

```ts
type State = {
  // 多裝置 - Step 1 管理
  devices: Array<{
    id: string;                 // dev-xxx
    name: string;
    description: string;
    protocol: 'modbus_tcp' | 'modbus_rtu' | 'modbus_udp' | 'fatek_fbs' | 'mc_3e' | 'mqtt';
    config: Record<string, any>; // host/port/slave_id/timeout/... (依 protocol 而異)
    status: 'draft' | 'tested' | 'active';
    test: null | {
      status: 'running' | 'success' | 'failed';
      latency_ms?: number;
      stages: Record<string, { status: string; latency_ms: number }>;
      tested_at?: Date;
    };
  }>;

  // 多規則 - Step 2 管理；每條規則綁定一個 device
  rules: Array<{
    id: string;                 // rule-xxx
    device_id: string;          // 所屬裝置
    name: string;               // 規則名稱 (e.g. "Holding Registers")
    start_address: string;      // "40001"
    count: number;              // 1-64
    data_type: 'bool'|'int16'|'int32'|'int64'|'uint16'|'uint32'|'uint64'|'float32'|'float64'|'string';
    naming_prefix: string;      // "SENSOR_" → SENSOR_01, SENSOR_02 …
    enabled: boolean;
    scale_multiplier: number;   // 規則層級的線性轉換預設 (Step 3 可個別覆寫)
    scale_offset: number;
    data_format: '' | 'ABCD' | 'BADC' | 'CDAB' | 'DCBA';  // byte order
    skipped_addresses: string[];          // 被略過的位址
    // Local Modbus Share 對外發布
    share_enabled: boolean;
    share_start_register: number | null;  // null = 自動接續配置
    share_stride: number | null;          // null = 用資料型態寬度
  }>;

  // 由 Step 2 的 useEffect 依 rules 衍生；Step 3/4 消費
  points: Array<{
    id: string;
    device_id: string;
    rule_id: string;
    rule_name: string;
    name: string;               // SENSOR_01
    address: string;            // "40001"
    data_type: string;
    function: 'coil'|'discrete_input'|'input_register'|'holding_register';
    width: number;              // register 寬度
    enabled: boolean;
    skipped: boolean;
    _rule_scale: number;        // 來自規則的預設
    _rule_offset: number;
  }>;

  // 點位 → Tag 映射 (Step 3)
  mappings: Record<pointId, {
    point_id: string;
    tag_key: string;            // "line01.temp.inlet"
    display_name: string;
    unit: string;
    target_type: string;
    scale: number;              // 個別覆寫
    offset: number;
    enabled: boolean;
  }>;

  // 資料庫 - Step 4
  db: {
    connector: {
      kind: 'sqlite'|'postgres'|'mysql'|'sqlserver';
      name: string;
      host: string;
      port: number;
      database: string;
      username: string;
      schema: string;
      table: string;
      write_mode: 'insert' | 'upsert';
      write_interval_seconds: number;
      timestamp_column: string;
      status: string;
    };
    targets: Record<pointId, {
      tag_id: string;
      column_name: string;
      enabled: boolean;
    }>;
  };

  // 設定頁
  settings: {
    connectors: Connector[];                                    // 連接器池
    timeseries: { write_precision, partition_interval, batch_size, retention_days };
    scheduler: { default_interval_ms, default_retry_count, default_retry_delay_ms, breaker_threshold, auto_start };
    modbus_share: { enabled, bind_address, port, slave_id, base_register };
    general: { theme, locale, addr_format, api_base, api_version, timeout_seconds, log_level, sse_heartbeat_seconds, enable_debug_panel, enable_audit_log };
  };

  // Commit
  commit?: { status, logs: Array<{ label, detail, status }>, started_at, finished_at };
  committed: boolean;
};
```

---

## 三、共用版面 (App shell)

```
┌───────────────────────────────────────────────────────────────┐
│ [☰] [⚡] Datalink Workbench  Datalink › Workbench › <step>   │
│         go_gateway · v1                  [儲存草稿][取消][scheduler ●idle/running] │
├──────┬──────────────────────────────────────┬─────────────────┤
│      │  STEP 0X / 04  |  <subtitle>          │  即時設定摘要     │
│ 側邊 │  <Step Title>            ▢▢▢▢       │  ┌─設備──┐       │
│ 步驟 │                                       │  ┌─規則──┐       │
│ 列   │  ─ Step content ─                     │  ┌─映射──┐       │
│      │                                       │  ┌─資料庫─┐      │
│ 設定 │                                       │  ┌─已部署─┐      │
│      │  [← 上一步]    ⌘B 收合側欄             │                 │
└──────┴──────────────────────────────────────┴─────────────────┘
                                                         [Tweaks]
```

### App shell 元件
- **頂列 (sticky)**：collapse 鈕、品牌標識、breadcrumb、儲存草稿/取消按鈕、scheduler 狀態指示
- **左側 StepRail**：4 個流程步驟 + 分隔線 + 「設定」入口；可收合（⌘B 或 collapse 鈕）；收合時只剩 icon，hover 顯示 tooltip
- **中央 main**：步驟內容或設定頁
- **右側 SummaryRail (≥xl 顯示)**：即時設定摘要，可在 Tweaks 開關

### 設計系統 tokens
- 主背景：`#0b1220` (略深於 slate-950) + 漸層光暈 + 24px 格線底
- Surface：slate-900/60 (含 backdrop-blur)，border slate-700/60
- Primary：blue-500/600
- Success：emerald-500
- Warning：amber-500
- Error：red-500
- 文字：slate-100 (主) / slate-400 (次) / slate-500 (弱)
- mono font：JetBrains Mono (數值、位址、API)
- chip：圓邊 `rounded-full`，4 種 tone (info/success/warning/error)
- card：`rounded-2xl`，標題列含 icon + title + subtitle + aside (chip/按鈕)

---

## 四、互動規範

| 行為 | 觸發 | 備註 |
|---|---|---|
| 收合側邊欄 | 頂列按鈕 / ⌘B | flex item 用 `flex: 0 0 64px/232px` 控制寬度 |
| 切換 step | StepRail 點擊 | 只有 reachable (已完成或前一步已完成) 才可跳轉 |
| 進入下一步 | Step 內 onContinue | 驗證通過才啟用；同時把 step id 加入 completed set |
| 重置/示範 | Tweaks 面板 | 包含「重置流程」「加入第二台設備」「加入第二條規則」「跳到 Step 4」|
| 切換摘要欄顯示 | Tweaks 面板 | 持久化到 localStorage |

---

## 五、檔案結構

```
index.html          - 入口 + Tailwind/Babel CDN + script 引入順序
shared.jsx          - 共用元件 (Icon/Button/Field/Toggle/StatusChip/SectionCard) + 工具 (derivePoints/deriveAllPoints/makeDefaultRule/makeDefaultDevice/computeShareLayout) + 常數 (PROTOCOLS/DATA_TYPES/DB_KINDS/POINT_SEMANTIC)
tweaks-panel.jsx    - Tweaks 浮動面板與控制元件
step1-device.jsx    - Step 1 (多裝置)
step2-rule.jsx      - Step 2 (多規則 + 記憶體網格 + Share 位址)
step3-mapping.jsx   - Step 3 (Tag 映射 + 個別 scale/offset + 轉換管線預覽)
step4-database.jsx  - Step 4 (Connector + 欄位映射 + Commit 動畫)
settings.jsx        - 設定頁
app.jsx             - 主 App shell + StepRail + SummaryRail + state
```

> Script 載入順序：tweaks-panel → shared → step1 → step2 → step3 → step4 → settings → app
> 各元件透過 `window.XXX = XXX` 暴露成全域，避免 babel script 跨檔案 scope 問題

---

## 六、繁中 microcopy 規範

- 省略號統一用 `…`（非 `...`）
- loading 標籤：「載入中…」「儲存中…」「測試中…」「提交中…」
- 按鈕：「儲存」「取消」「建立」「下一步」「繼續」「提交並啟動排程器」
- 狀態：draft=草稿 / active=啟用 / disabled=停用 / tested=已測試 / committed=已部署

---

## 七、各頁文件索引

- [Step 1 · 新增裝置](./step1-device.md)
- [Step 2 · 接入規則](./step2-rule.md)
- [Step 3 · 點位映射](./step3-mapping.md)
- [Step 4 · 儲存資料庫](./step4-database.md)
- [設定頁](./settings.md)
