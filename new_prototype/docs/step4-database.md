# Step 4 · 儲存資料庫 (DB Target + Commit)

## 目的

選擇/設定資料庫 Connector，把每個 Tag 對應到資料表的欄位，最後**一鍵 commit**整套設定（10 個 API 呼叫依序送出）並啟動 collector。

對應後端：`POST /db-connectors/:id/test` → `POST /db-targets × N` → `POST /scheduler/start`。

---

## 完整 Prompt (可直接生成此頁)

> 為 go_gateway 系統做一個「儲存資料庫」頁。頂部是 Connector 設定，下方左邊是 Tag → 欄位映射表，右邊是提交摘要與 Commit 動畫。
>
> ---
>
> **第一層 · Connector (SectionCard, col-span-12)**
>
> - 標題「資料庫連接器 (Connector)」，副標「寫入時序值的目標資料庫」
> - icon：db；aside：success chip 顯示已連線 host
> - 表單橫向排列 (12-col grid)：
>   - **資料庫類型** (col 3)：4 個方形選擇器，emoji icon + 名稱 (SQLite 🗄️ / PostgreSQL 🐘 / MySQL 🐬 / SQL Server 🪟)
>   - 連線名稱 / Host / Port / Database / Schema / Table / 使用者 ... (混合排佈)
> - 寫入策略 (border-t 區塊)：radio 選 INSERT (時序追加) / UPSERT (依時間戳)；右側「寫入間隔」number + 「秒」
>
> ---
>
> **第二層左 (col 7) · Tag → 資料表欄位 (SectionCard, content p-0)**
>
> - 標題`Tag → 資料表欄位`，副標 `寫入 {schema}.{table}`
> - icon：table
> - aside：若有衝突→ error chip「欄位衝突」；否則 success chip「{N} 對應」
> - 表格：
>
>   | Tag | 點位 | → | 資料表欄位 | 欄位型態 | 啟用 |
>
>   - Tag：tag_key (mono blue-200) + display_name (小字)
>   - 點位：name + @address (兩行 mono)
>   - 資料表欄位：select dropdown，列出所有 `SAMPLE_DB_TABLES[kind][0].columns` 中非 PK 的欄位
>     - 衝突時用紅框 + 旁邊 alert icon
>   - 欄位型態：mono 小字顯示資料表欄位的 SQL 型別
>   - 啟用：size sm Toggle
>
> - 若有衝突，表尾紅色提示 banner：「偵測到多個 Tag 寫入同一資料表欄位。請調整以避免覆寫。」
>
> **自動欄位匹配演算法**
>
> ```js
> const used = new Set();
> enabledPoints.forEach((p, i) => {
>   const tagShort = mapping.tag_key.split('.').pop();
>   // 1. 嘗試精確匹配：欄位名 === tagShort，或前綴/後綴 _
>   const exact = columnNames.find(c =>
>     !used.has(c) && (c === tagShort || c.endsWith('_' + tagShort) || c.startsWith(tagShort + '_'))
>   );
>   // 2. 沒匹配到，依 index 後備；若 index 那欄已用，找第一個未用
>   const fallback = columnNames[i % columnNames.length];
>   const auto = exact || (used.has(fallback) ? columnNames.find(c => !used.has(c)) : fallback);
>   used.add(auto);
>   target[p.id] = { column_name: auto, ... };
> });
> ```
>
> 此邏輯確保 8 個點位剛好對應 8 個欄位，不會自動產生衝突。
>
> ---
>
> **第二層右 (col 5) · 提交摘要 (SectionCard)**
>
> - 標題「提交摘要 (Commit)」，副標「一次原子寫入所有設定」
> - icon：save
>
> **未開始 commit 時 - 摘要列表 (5 項)**
>
> ```
> ┌ 新增設備 ── {N} 個                    {names · protocols}
> ├ 接入規則 ── {N} 條                    {rules summary}
> ├ 點位 ──── {N} 個                      加入快速輪詢 (1s)
> ├ 標籤映射 ── {N} 個 Tag                含 scale/offset/cast 管線
> └ 資料庫寫入 ─ {N} 欄位                  {DBKind} → {schema}.{table}
> ```
>
> 每列是 rounded-lg card，左邊 label (uppercase 微小字) + value (mono 大字)，右邊 sub 文字 (10px slate-500 右對齊)。
>
> 下方一個寬版 success Button「提交並啟動排程器」(disabled if 有衝突或無啟用 target)
>
> 底部 info：「將執行 10 個 API 呼叫並啟動 collector」
>
> **執行 commit 時 - 動畫式 log**
>
> ```
> ✓ POST /devices × {N}                  PLC-生產線-01 (modbus_tcp) · PLC-02 (mc_3e)
>   200
> ✓ POST /devices/:id/activate × {N}     draft → active
>   200
> ✓ POST /source-rules × {N}             Rule1: 40001×8 · Rule2: 00100×4
>   200
> ✓ POST /points × {N}                   bulk create
> ✓ POST /polling-groups                 快速輪詢 1s, enabled
> ✓ POST /tags × {N}                     register tag keys
> ✓ POST /mappings × {N}                 point ↔ tag, scale pipeline
> ✓ POST /db-connectors/:id/test         postgres tsdb.internal:5432
> ✓ POST /db-targets × {N}               → public.sensor_readings
> ✓ POST /scheduler/start                collectors started
> ● 正在執行下一個指令…  (running 中才顯示)
> ```
>
> 每 280ms 推進一個 log，sweep-in 動畫；右側 200 status code 用 emerald 小字。
>
> **完成後**
>
> emerald 卡片：圓形勾 icon + 「設定已套用 · 開始收集資料」+ 「Scheduler 已啟動 · 第一筆資料預計在 ~Ns 後寫入」+ 「前往 Runtime Dashboard」secondary button (`onCommit` callback)

---

## State 介面

```ts
Props: {
  state: { devices, rules, points, mappings, db: {connector, targets}, ... };
  setState: (updater) => void;
  onCommit: () => void;   // 點「前往 Runtime Dashboard」時觸發
}

Local state:
- committing: boolean        // 提交動畫進行中
- committed: boolean         // 已完成
```

state.commit:
```ts
{
  status: 'running' | 'success';
  logs: Array<{ label, detail, status }>;
  started_at: Date;
  finished_at?: Date;
}
```

---

## 互動規範

| 行為 | 觸發 | 結果 |
|---|---|---|
| 切換資料庫類型 | 4 個 kind 卡片 | 重新依該 kind 找 sample table，重置 status 為 unknown |
| 改欄位映射 | dropdown | updateTarget(p.id, { column_name }) |
| 啟用/停用單筆 | Toggle | target.enabled = bool |
| 改寫入策略 | radio | INSERT or UPSERT |
| 提交 | success button | 觸發 startCommit() — 10 個 log 依序加入 state.commit.logs |
| 前往 Dashboard | success card 內按鈕 | onCommit callback |

---

## 衝突偵測

```js
const colUsage = {};
Object.entries(targets).forEach(([pid, t]) => {
  if (t.enabled && t.column_name) {
    (colUsage[t.column_name] ??= []).push(pid);
  }
});
const hasConflict = Object.values(colUsage).some(arr => arr.length > 1);
```

衝突會：
- aside chip 變紅「欄位衝突」
- 衝突 row 的 select 邊框變紅、旁邊出現 alert icon
- 表尾紅色 banner 提示
- 「提交並啟動排程器」按鈕 disabled

---

## 視覺重點

- Connector 區塊用標準的卡片式類型選擇器 (上方 col 3 + 文字欄位)
- 提交摘要列表用 ` ┌├└` 視覺對齊 (內部結構)，每列 label 在左、 value 大字 mono、sub 右側 10px
- Commit log 一條一條動畫式進入；完成後的 emerald 大勾 + 訊息要明顯
- 已完成 commit 後表單變只讀；摘要區只顯示 logs

---

## 驗收標準

- [ ] 預設 8 個點位剛好對應 8 個欄位 (Postgres sample table)，無衝突
- [ ] 改某筆 column_name 變成已用值，立即顯示紅框 + alert + banner
- [ ] 「提交並啟動排程器」按下後，10 條 log 依序冒出 (約 3 秒)
- [ ] 完成後出現 success 大勾卡，scheduler 狀態指示燈在頂列也變 running
- [ ] 摘要中「新增設備」會列出所有設備 (測試多裝置情境)
- [ ] commit 動畫期間 button 變 disabled，無法重複觸發
