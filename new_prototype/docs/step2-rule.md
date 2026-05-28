# Step 2 · 接入規則 (Source Rules)

## 目的

為每一個設備定義一條以上的「位址範圍規則」(SourceRule)，每條規則展開後就是一組點位 (Point)。同時設定每條規則在 Local Modbus Share 對外發布的位址，並於合併點位表中偵測跨規則衝突。

對應後端：`POST /source-rules` → 系統依此 bulk-create `POST /points` 並加入 polling_group。

---

## 完整 Prompt (可直接生成此頁)

> 為 go_gateway 系統做一個「接入規則」頁，支援多條規則同時編輯。版面三層：上方規則 tab 列、中層編輯器+網格、下方合併點位表。
>
> ---
>
> **第一層 · 規則 tab 列 (SectionCard)**
>
> - 標題「接入規則」，副標「可同時定義多條範圍規則，每條展開為一組點位」
> - aside：`{N} 條規則 · 共 {M} 點位` chip；若有位址衝突再加紅色 chip「{X} 位址衝突」
> - tab 列水平捲動，每個 tab 是 `min-w-[200px]` 卡片：
>   - 第一行：顏色點 + inline 可編輯規則名稱 + hover 才顯示的啟用 toggle + 刪除 X
>   - 第二行：`{start_address} · {data_type}` (mono) + `{enabled}/{count} pts`
>   - 第三行 (僅當有 2+ 裝置時顯示)：裝置顏色點 + 裝置名稱
>   - 選中時用該規則顏色 (blue/emerald/amber/fuchsia/cyan/rose 循環)
> - 最右邊「+ 新增規則」虛線邊框按鈕
>
> ---
>
> **第二層 · 編輯選中規則 (12-col grid)**
>
> 左 (col 5) `SectionCard "編輯：{rule_name}"`：
>
> 1. **所屬裝置** — Select dropdown，只有 1 個裝置時 disabled，>=2 顯示所有裝置「{name} ({protocol})」
> 2. **2x2 grid 基本參數**：
>    - 「起始位址」(required, mono, hint: "Modbus 4xxxx = Holding Register")
>    - 「點位數量」(required, number, min 1, max 64)
>    - 「資料型態」(required, hint: "寬度 N register") - 從 DATA_TYPES (10 種)
>    - 「命名前綴」(hint: "生成 SENSOR_01, SENSOR_02 …")
>    - 修改 start_address/count/data_type 時清空 skipped_addresses
> 3. **範圍摘要區塊** (rounded-xl bg-slate-950/40)：
>    - 「{startAddr} → {endAddr}」(mono)，右側「共 N register」
>    - 顯示 Modbus function code chip：FC 03/06/16 (Holding) / FC 04 (Input) / FC 02 (Discrete Input) / FC 01/05 (Coil) — 依起始位址第一位數字推斷
> 4. **線性轉換 `<details open>`**：副標「線性轉換 / 字節序 (規則預設值)」
>    - 提示：「這裡設的值會作為下一步點位映射的初始值，仍可在映射表中個別調整」
>    - Multiplier (number step 0.01) + Offset (number step 0.01)
>    - Byte order：ABCD/BADC/CDAB/DCBA 四選項 + 空值預設 ABCD
> 5. **Local Modbus Share 摺疊區** `<details>`：副標「Local Modbus Share 位址」
>    - 若全域 Share 未啟用：顯示「全域未啟用」chip
>    - 若啟用：顯示該規則占用的 register 範圍 chip (emerald)
>    - 表單：Share 起始 Register (留空 = 自動)、Stride (留空 = 用資料型態寬度)
>    - Toggle：「透過 Modbus Share 對外發布此規則的點位」
>    - 若有手動設定，顯示「改回自動分配」reset 連結
> 6. **底部 toggle**：「啟用此規則 (建立後立即輪詢)」
>
> 右 (col 7) `SectionCard "當前規則點位網格"`：副標「點擊切換略過 · Shift+點擊：範圍切換 · Ctrl/⌘+點擊：加入選取」
>
> - aside chips：啟用 N、略過 M、已選 X (有選取時才顯示)
> - **批次工具列** (rounded-lg, 內嵌)：
>   - 全部啟用 / 全部略過 / 反轉啟用 (ghost size sm 按鈕含 icon)
>   - 略過選取 (N) / 啟用選取 / 清除選取 — disabled if N===0
>   - 右側 mono text：`total {count} · stride {width}`
> - **網格** (rounded-xl bg-grid 格線底)：
>   - 上方 mono summary：`{startAddr} … {endAddr}` ⇆ `{enabled} active · {skipped} skipped · share@{shareStart}`
>   - 格子：4-col on mobile, 6-col on sm, 8-col on md+，每格高 48px (mem-cell class)
>     - 正常：規則顏色背景；略過：amber
>     - 衝突：紅色 ring 圍邊 + 右上角紅色 ! 徽章
>     - 多選中：blue ring 圍邊
>     - 啟用且有 share：右下角 mono 小字「→40001」浮動標籤
>     - 內容：位址 (addr class) + 點位名稱 (meta class)
>     - title 屬性：點位名 + share 位址 + 衝突警告
>   - 點擊行為：
>     - 一般點擊 → toggle skip
>     - Shift+click → 對 lastClickedIdx → 當前 idx 範圍批次套用 (若當前已 skip 則全部啟用，反之全部略過)
>     - Ctrl/⌘+click → 加入或移出多選 set
>   - 底部 mono legend：width / start / 衝突警示
>
> ---
>
> **第三層 · 合併點位表 (SectionCard, content p-0)**
>
> - 標題「將要建立的點位 (N)」，副標「所有啟用規則合併後的結果」
> - aside：若 Share 啟用則顯示 info chip「Modbus Share 開啟」；mono：`POST /source-rules × {rules} · POST /points × {points}`
> - 表格 (max-h 280px overflow-y-auto, sticky thead)：
>   | 所屬設備 | 所屬規則 | 點位名稱 | 位址 | 型態 | 功能碼 | Share 位址 | 狀態 |
>   - 所屬設備/規則：顏色點 + 名稱
>   - 位址：有衝突時紅字 + ⚠
>   - Share 位址：值是 mono emerald 文字；off 或 skipped 顯示「—」
>   - 狀態 chip：建立 (success) / 跳過 (warning) / 衝突 (error)
> - 表尾 (border-t)：info icon + 「指派到輪詢群組 快速輪詢 (1s)」 + 若 Share 未啟用提示前往設定；右側「繼續到映射」Button (disabled 條件：totalEnabled === 0 或 conflictAddrs.size > 0)

---

## 關鍵衍生計算

```js
// 全部點位
const allPoints = deriveAllPoints(rules, devices[0]?.id);

// 衝突偵測
const addrUsage = {};
allPoints.forEach(p => {
  if (!p.skipped) (addrUsage[p.address] ??= []).push(p.rule_id);
});
const conflictAddrs = new Set(
  Object.entries(addrUsage).filter(([_, arr]) => arr.length > 1).map(([k]) => k)
);

// Share 位址布局 (computeShareLayout in shared.jsx)
// 每條規則自動接續配置，或用 share_start_register 手動指定
const shareLayout = computeShareLayout(rules, settings.modbus_share.base_register);
```

---

## 互動規範

| 行為 | 觸發 | 結果 |
|---|---|---|
| 新增規則 | 「+ 新增規則」 | 用 makeDefaultRule 建一條，預設 device_id 用當前選中規則的或第一個 device |
| 刪除規則 | tab 上 X 鈕 | 至少保留一條 |
| 啟用/停用規則 | tab hover 顯示的 toggle | 不刪除規則，只是不展開點位 |
| 改名 | tab 上的 inline input | onChange 直接 setState |
| 切換 skip | 格子單擊 | 加入/移除 skipped_addresses |
| 範圍切換 skip | Shift+click | 以 lastClickedIdx 為基準 |
| 加入多選 | Ctrl/⌘+click | 加入/移除 gridSelection set；不影響 skip |
| 略過選取 | 工具列按鈕 | gridSelection 全部加入 skipped_addresses |
| 變更裝置 | 「所屬裝置」select | 該規則所有 points 重新衍生並改 device_id |
| 自動 share 配置 | share_start_register 留空 | computeShareLayout 依規則順序接續分配 |

---

## 規則顏色與裝置顏色

兩者使用同一個 6 色循環 `['blue','emerald','amber','fuchsia','cyan','rose']`，但 index 來源不同：

- 裝置色 = `devices.findIndex(d => d.id === deviceId) % 6`
- 規則色 = `rules.findIndex(r => r.id === ruleId) % 6`

在合併點位表上會同時看到兩個顏色點：第一欄是裝置色，第二欄是規則色。

---

## 驗收標準

- [ ] 初始顯示 1 條規則 (40001 起 8 個 int16)，網格顯示 8 格
- [ ] 點「+ 新增規則」後可以為第 2 條規則選擇所屬裝置
- [ ] Shift+Click 範圍 skip 正確運作
- [ ] 兩條規則位址重疊時，重疊格子顯示紅色 ring + ! 徽章，且「繼續」按鈕 disabled
- [ ] 預設啟用 Modbus Share 時，網格格子右下角顯示對應 share 位址
- [ ] 修改 share_start_register 為手動值後，「改回自動分配」按鈕出現並可重設
- [ ] 修改 data_type 時，網格 stride 與 share 位址間距同步更新
