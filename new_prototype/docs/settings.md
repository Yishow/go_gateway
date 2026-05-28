# 設定頁 (Settings)

## 目的

集中管理系統級設定，獨立於 4 步驟流程之外。包含資料庫連接器池、時序儲存策略、排程器預設、Local Modbus Share、介面/API/診斷選項。

側邊欄「系統 → 設定」入口進入；不在 step 流程的 next/back 中。

---

## 完整 Prompt (可直接生成此頁)

> 為 go_gateway 系統做一個獨立的「設定」頁。整頁是一連串的 SectionCard，由上而下依序：標題列、資料庫連接器池、時序儲存策略 + 排程預設 (並排)、Local Modbus Share、介面/API/診斷 (3 欄並排)、底部 sticky 儲存列。
>
> ---
>
> **1. 標題列**
>
> rounded-2xl 漸層卡 (slate-900/80 → slate-900/40)：
> - 左邊 11x11 圓角藍色 icon (sliders)
> - 中間：「系統設定」(lg semibold) + 一行說明「資料庫連接器池、時序儲存策略、排程器、Local Modbus Share、介面偏好」
>
> ---
>
> **2. 資料庫連接器池 (SectionCard)**
>
> - 標題「資料庫連接器池」，副標「可註冊多個目標資料庫，於映射時挑選」
> - icon：db
> - aside：「+ 新增連接器」secondary button
> - 內容 (content p-0)：
>   - 空狀態：「尚未註冊任何連接器。點上方『新增連接器』開始。」
>   - 列表 (divide-y)，每行：
>     - 左邊 10x10 emoji icon 大格 (依 kind 顯示 🗄️/🐘/🐬/🪟)
>     - 中間 12-col grid 表單：連線名稱(3) / 類型(2) / Host(3) / Port(1) / Database(3) / 使用者(3) / 預設寫入間隔(2) / 啟用 toggle(2)
>     - 底部跨欄：status chip + last_check_at 時間戳 (若 ready) 或 last_check_error (紅) + 右側對齊「測試」/「✕ 刪除」ghost button
>   - 測試按鈕：點下後該 connector 進入 testing status，900ms 後 85% 機率 ready / 15% unreachable
>
> ---
>
> **3. 時序儲存策略 + 排程預設 (grid col 6 + 6)**
>
> 左 `SectionCard "時序儲存策略"`：icon spark
> - 時間精度：秒 / 毫秒
> - 分區間隔：每日 / 每週 / 每月
> - 批次寫入大小：number, hint「一次最多寫入的筆數」
> - 保留天數：number, hint「超過天數的資料將被歸檔」
>
> 右 `SectionCard "排程預設"`：icon refresh
> - 預設輪詢間隔 (ms)
> - 重試次數
> - 重試延遲 (ms)
> - 斷路器門檻 (連續錯誤次數)
> - 跨欄 toggle：「開機時自動啟動 collector」
>
> ---
>
> **4. Local Modbus Share (SectionCard)**
>
> - 標題「Local Modbus Share」，副標「把採集到的 Tag 重新發佈成 Modbus TCP，供下游 SCADA 訂閱」
> - icon：cable
> - aside：總開關 Toggle (size md)
> - 若啟用，顯示 4-col grid：
>   - 綁定位址 (hint: 0.0.0.0 = 全部介面)
>   - Port (default 5020)
>   - Slave ID
>   - 起始 Register (default 40001)
> - 若停用，僅顯示「啟用後可將所有 Tag 再次經由 Modbus TCP 提供給其他系統。」
>
> ---
>
> **5. 介面 / API / 診斷 (grid col 4 + 4 + 4)**
>
> 「介面」icon eye：
> - 主題 (深色 / 淺色 / 跟隨系統)
> - 語言 (繁中 / English)
> - 位址顯示格式 (Modbus 5 位 / 十六進位 / 原始)
>
> 「API 端點」icon flow：
> - Base URL (不含尾斜線)
> - API 版本 (v1 / v2 beta)
> - 請求逾時 (秒)
>
> 「診斷」icon alert：
> - 日誌等級 (trace/debug/info/warn/error)
> - SSE 心跳 (秒)
> - Toggle「顯示 Debug 面板」
> - Toggle「操作審計記錄」
>
> ---
>
> **6. 底部 sticky 儲存列**
>
> `<div className="sticky bottom-4 rounded-xl border bg-slate-900/95 backdrop-blur px-4 py-3 flex items-center justify-between shadow-2xl">`
>
> - 左：info icon「設定會立即套用，並於下次重啟後生效。」
> - 右：「重設為預設」ghost button + 「儲存所有設定」success button (size md, icon save)

---

## State 介面

```ts
Props: {
  state: { settings: SettingsObject, ... };
  setState: (updater) => void;
}

settings = {
  connectors: Connector[],
  timeseries: { write_precision, partition_interval, batch_size, retention_days },
  scheduler: { default_interval_ms, default_retry_count, default_retry_delay_ms, breaker_threshold, auto_start },
  modbus_share: { enabled, bind_address, port, slave_id, base_register },
  general: { theme, locale, addr_format, api_base, api_version, timeout_seconds, log_level, sse_heartbeat_seconds, enable_debug_panel, enable_audit_log },
};

Local state:
- testingId: string | null   // 哪個 connector 正在測試
```

---

## 主要函式

```js
const update = (patch) =>
  setState(s => ({ ...s, settings: { ...s.settings, ...patch } }));

const updateNested = (key, patch) =>
  setState(s => ({ ...s, settings: { ...s.settings, [key]: { ...s.settings[key], ...patch } } }));

const addConnector = () => {
  update({ connectors: [...connectors, { id: nextId(), name: '新連線 N', kind: 'postgres', ... }] });
};

const testConnector = (id) => {
  updateConnector(id, { status: 'testing' });
  setTimeout(() => {
    const ok = Math.random() > 0.15;
    updateConnector(id, {
      status: ok ? 'ready' : 'unreachable',
      last_check_at: new Date().toISOString(),
      last_check_error: ok ? '' : 'connection refused',
    });
  }, 900);
};
```

---

## 與 Step 2 / Step 4 的連動

- `settings.modbus_share.enabled` → Step 2 規則的「Share 位址」欄位是否生效
- `settings.modbus_share.base_register` → Step 2 自動配置 share register 的起點
- `settings.scheduler.default_interval_ms` → 新建立 polling group 的預設值
- `settings.timeseries.*` → 影響後端寫入時序庫的策略 (commit 時送出)
- `settings.connectors` → Step 4 的 connector 來源 (未來會把 Step 4 內建的單一 connector 改為從這裡選)

---

## 視覺重點

- 整頁長卷軸，靠 SectionCard 分區
- 底部 sticky 儲存列確保使用者任何時候都能存檔，不用滾回頂部
- Connector 列表用 divide-y 而不是個別卡片，緊湊呈現
- status chip 多種 tone (ready=success / testing=info / unreachable=error / unknown=draft)

---

## 驗收標準

- [ ] 預設載入 1 個 ready 的 connector「TimeSeries Prod」
- [ ] 點「新增連接器」立即新增一筆並進入編輯狀態
- [ ] 「測試」按鈕轉圈 900ms 後出現結果，狀態 chip 更新
- [ ] Local Modbus Share toggle 切換時，下方表單區會顯示/收起
- [ ] 設定變更時立即反映到 Step 2 的 Share 位址計算
- [ ] 底部儲存列在頁面捲動時始終可見
