# Step 3 · 點位映射 (Point → Tag)

## 目的

把 Step 2 衍生出來的所有點位 (Point) 一一對應到 Tag (時序資料庫的命名鍵)，並設定每個點位「個別」的線性轉換與目標型態。右側即時呈現轉換管線的數值流。

對應後端：`POST /tags × N` → `POST /mappings × N` (內含 transform_pipeline)。

---

## 完整 Prompt (可直接生成此頁)

> 為 go_gateway 系統做一個「點位映射」頁。把所有點位列在左側可編輯的表格，每列獨立調整 Tag Key / 顯示名稱 / 單位 / 目標型態 / scale / offset / enabled；右側選中一列時即時呈現該點位的轉換管線。
>
> ---
>
> **左 (col 8) `SectionCard "點位 → Tag 映射"` (content p-0)**
>
> - 副標「每個點位對應一個 Tag (時序資料庫的命名鍵)」
> - aside chip：「{enabled} / {total} 已啟用」(success 或 warning tone)
> - 表格 (overflow-x-auto, text-xs)，欄位：
>
>   | (color dot) | 點位 | 位址 | Tag Key (28%) | 顯示名稱 | 單位 | 目標型態 | × Scale | + Offset | 啟用 |
>
>   - 第一欄是 1.5px 的規則顏色點 (從 step 2 沿用)，可一眼分辨點位來自哪條規則
>   - 點位 / 位址：mono 純文字 (位址 blue-200)
>   - Tag Key / 顯示名稱 / 單位：透明背景 `<input>`，hover 顯示邊框，focus 變藍框
>   - 目標型態：透明 `<select>` 列出所有 DATA_TYPES
>   - × Scale / + Offset：`<input type="number" step=0.01 className="w-16 text-right">`，初值來自規則層 (_rule_scale / _rule_offset)
>   - 啟用：size sm 的 Toggle
>   - 點任一列 (除了輸入元素) → setSelectedIdx，row 套上 `outline outline-1 outline-blue-500/30 bg-blue-500/[0.06]`
>   - 所有 input/select/toggle 都要 `onClick={e=>e.stopPropagation()}` 避免冒泡到列選取
>
> - 表尾 (border-t)：
>   - 左：info icon「點選列以在右側看轉換管線預覽」
>   - 右：bolt icon「將選中列的 Scale / Offset / 型態 套用到全部」連結，按下後把選中列的這三個值同步到所有 mapping
>
> ---
>
> **右 (col 4) `SectionCard "轉換管線預覽"`**
>
> - 副標：當前選中的點位名稱與位址
> - icon：flow
>
> **轉換步驟條列 `<ol>`** (連續直線連接)
>
> 1. **decode (raw)** - 顯示模擬 raw 值 (大字體, mono)：「243」，下方 mono 小字「register @ 40001, type=int16」
> 2. **scale (linear)** - 顯示算式：`{raw} × {scale} + {offset} = {scaled}`，其中
>    - raw 用 slate-500
>    - scale/offset 用 blue-300
>    - 結果 (toFixed 3) 用 emerald-300 粗體
> 3. **cast to {target_type}** - 顯示「cast({scaled} → {target_type})」
> 4. **final → Tag** (emerald 強調卡)
>    - 巨大數字 (2xl, mono semibold emerald-200)：依 target_type 格式化
>      - bool → 'true'/'false'
>      - floatN → toFixed(2)
>      - 其它 → Math.round
>    - 右下：單位 (emerald-300)
>    - 底部 mono：`→ tag.{tag_key}`
>
> 步驟編號圓點 (7x7 grid place-items-center)，連在一條 1px 直線上 (`before:absolute before:left-3.5 before:top-2 before:bottom-2 before:w-px before:bg-slate-700`)
>
> **API payload 預覽** (rounded-lg bg-slate-950/60)
>
> ```
> POST /mappings
> {
>   "point_id": "pt-xxx…",
>   "tag_id": "tag.line01.temp.inlet",
>   "transform_pipeline": [
>     { "type": "decode" },
>     { "type": "scale", "params": { "scale": 0.1, "offset": 0 } },
>     { "type": "cast", "params": { "to": "float64" } }
>   ],
>   "enabled": true
> }
> ```
>
> **底部繼續按鈕**：「設定資料庫寫入」(disabled if 任一 mapping tag_key 為空)
>
> ---
>
> **狀態初始化邏輯**
>
> ```js
> React.useEffect(() => {
>   enabledPoints.forEach((p, i) => {
>     const semantic = POINT_SEMANTIC[i % POINT_SEMANTIC.length];  // 預設語意 (溫度/壓力/流量…)
>     mappings[p.id] = existing || {
>       point_id: p.id,
>       tag_key: semantic.tag_key,        // "line01.temp.inlet" …
>       display_name: semantic.display,
>       unit: semantic.unit,
>       target_type: 'float64',
>       scale: p._rule_scale ?? 1,         // 從規則層帶入
>       offset: p._rule_offset ?? 0,
>       enabled: true,
>     };
>   });
> }, [points.length]);
> ```
>
> **模擬 raw 值**：固定 seeds `[243, 251, 1024, 985, 67, 542, 18, 1450]`，按 index 循環取。

---

## State 介面

```ts
Props: {
  state: { points: Point[], rules: Rule[], mappings: Record<id, Mapping>, ... };
  setState: (updater) => void;
  onContinue: () => void;
}

Local state:
- selectedIdx: number       // 哪一列被選中以顯示在右側預覽
```

---

## 互動規範

| 行為 | 觸發 | 結果 |
|---|---|---|
| 編輯任一欄 | input/select onChange | updateMapping(p.id, {...}) 即時寫入 state.mappings |
| 切換啟用 | toggle | mapping.enabled = bool |
| 選中列以預覽 | 點擊列 (非控件區域) | setSelectedIdx |
| 批次套用轉換 | 表尾「套用到全部」連結 | 將選中列的 scale/offset/target_type 同步到所有 mapping |
| 繼續 | 底部按鈕 | onContinue (前提：所有 mapping 都有 tag_key) |

---

## 視覺重點

- 規則顏色點放在每列最左，跨規則的點位混排時可一眼分辨來源
- 選中列：藍色 outline + 淡藍背景；hover 列：較淡的 slate 變色
- 表格欄位寬度：Tag Key 給 28%、scale/offset 各 16px 寬，其餘自動
- 右側預覽用大字體 + 顏色語意：raw(灰) → scale(藍) → cast → final(綠)，讓使用者理解資料流向

---

## 驗收標準

- [ ] 進入時所有點位都自動產生 mapping，scale/offset 預設從規則帶入 (0.1, 0)
- [ ] 每列獨立調整 scale/offset，右側預覽即時反映新算式與結果
- [ ] 修改 Tag Key 之後預覽底部的 `→ tag.xxx` 同步更新
- [ ] 「套用到全部」可一鍵把選中列的轉換參數複製到所有列
- [ ] tag_key 留空時「設定資料庫寫入」disabled，並 aside chip 轉 warning
- [ ] 規則顏色點與 Step 2 對應的規則色一致
