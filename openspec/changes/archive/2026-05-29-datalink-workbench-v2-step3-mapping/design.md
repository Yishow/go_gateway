## Context

Step 3 把 Step 2 衍生點位逐一綁定到 Tag，並設定個別轉換參數。原型 `step3-mapping.jsx`（286 行）整體結構不複雜，但有兩個關鍵互動：行內表格編輯（避免 modal）與右側轉換管線即時預覽。

行內表格的挑戰：每列含 7 個 input/select/toggle，點任一控件不應觸發列選取（用 `stopPropagation`）；同時要保持「點空白處選取列」直覺。

轉換管線預覽的挑戰：每次 raw / scale / offset / target_type / tag_key 任一更新，預覽要即時反算，避免延遲；同時數值格式化要依 target_type 切換。

## Goals / Non-Goals

**Goals:**

- 1:1 視覺與互動復刻原型 Step 3。
- 把轉換管線運算（runScale / castValue / formatFinal）抽到純函式，方便單元測試。
- 自動 mapping 初始化邏輯放在 reducer 層，由 component 透過 effect 觸發 `initMappings` action。
- 跨規則情境下 rule color dot 沿用 Step 2 規則色票，不重複定義。

**Non-Goals:**

- 不接後端。
- 不處理 transform_pipeline 進階步驟（lookup / conditional / formula）。
- 不變更 `mapping-pipeline` spec。
- 不處理 mapping 排序、匯入 / 匯出。
- 不支援多列同時選取（只支援單列 selectedIdx）。

## Decisions

### 拆檔策略：6 個元件 + 2 個 state module

原型 286 行拆為：

- `Step3Mapping.tsx`（~180 行）：容器，state derivation + auto-init effect + onContinue gate。
- `MappingTable.tsx`（~80 行）：左 col 8 表格殼層（thead + tbody）。
- `MappingRow.tsx`（~180 行）：單列；含 7 個 inline 編輯控件與 stopPropagation 處理。
- `TransformPreview.tsx`（~120 行）：右 col 4 容器，組合 `PipelineSteps` + `PayloadPreview`。
- `PipelineSteps.tsx`（~150 行）：4 步驟 ol + 圓點連線 + 顏色語意。
- `PayloadPreview.tsx`（~60 行）：`<pre>` mono JSON。

state module：
- `state/mappingDefaults.ts`（~80 行）：`POINT_SEMANTIC`、`RAW_VALUE_SEEDS`、`buildDefaultMapping(point, idx)`。
- `state/transformPipeline.ts`（~80 行）：`runScale` / `castValue` / `formatFinal` / `buildPayload` 純函式。

### 自動 mapping 初始化：reducer action 而非 component setState

原型在 component 內用 useEffect + setState 自動產生 mappings；移植時改為 reducer action `initMappingsForPoints(points)`，由 Step3Mapping 內 useEffect 監聽 `points` 觸發：

```ts
useEffect(() => {
  dispatch({ type: 'initMappingsForPoints', points: enabledPoints });
}, [enabledPoints.map(p => p.id).join(',')]);
```

reducer 行為：對每個 point，若 `state.mappings[p.id]` 已存在則保留（避免覆寫使用者編輯）；否則用 `buildDefaultMapping(p, i)` 產生。同時刪除 `state.mappings` 中對應 point 已不存在的 entry（清理孤兒）。

**Alternatives considered**：(A) 在 component 內 setState — 原型作法但難測；(B) 在 Step 2 衍生 points 時直接產 mapping — 違反「mapping 屬於 Step 3 概念」分離。

### 行內編輯 + 列選取共存：stopPropagation 策略

`MappingRow` 整列綁 `onClick={() => setSelectedIdx(i)}`，所有 input / select / toggle 內部用 `onClick={(e) => e.stopPropagation()}`。額外加上 `onMouseDown stopPropagation` 避免某些瀏覽器先觸發 mousedown 才 click。

key navigation：Tab 順序跟著 DOM 順序，不會誤觸列選取。

**Alternatives considered**：(A) 用 row 內非控件區域加 explicit click handler — 邊界判斷複雜；(B) 把選取改成 hover-only — 失去「明確選定一列來預覽」的直覺。

### 轉換管線數值流：使用 raw seed 序列

`buildRawSeed(idx)` 從 `RAW_VALUE_SEEDS = [243, 251, 1024, 985, 67, 542, 18, 1450]` 循環取（`seeds[idx % 8]`）。確定性 raw 讓 UI 預覽穩定、易於測試（不隨機）。

`runScale(raw, scale, offset)` 回傳 `raw * scale + offset`，浮點以 JS 標準算術。

`castValue(scaled, targetType)`：
- bool → `scaled !== 0`
- intN → `Math.round(scaled)` 並夾在型態範圍（不溢出）
- floatN → `scaled`
- string → `String(scaled)`

`formatFinal(value, targetType)`：
- bool → `"true"` / `"false"`
- floatN → `value.toFixed(2)`
- intN → `Math.round(value).toString()`
- string → `value`

**Alternatives considered**：(A) 隨機 raw — 測試易失敗；(B) 用 API 拿真實 raw — 違反「不接後端」邊界。

### 批次套用：reducer action

`bulkApplyTransform({ from: point_id, fields: ['scale','offset','target_type'] })` reducer action：讀 `state.mappings[from]`，把指定欄位複製到所有 mapping。

UI 觸發點：表尾連結「將選中列的 Scale / Offset / 型態 套用到全部」；連結文字依當前 `selectedIdx` 對應的 mapping 顯示「使用 line01.temp.inlet 的 Scale / Offset / 型態」。

**Alternatives considered**：(A) Modal 確認「將套用到 N 筆 mapping」— 增加摩擦；(B) 在表頭加「✓ Apply」按鈕 — 表頭已擁擠。

### 繼續閘：tag_key 不能為空

「設定資料庫寫入」按鈕在「任一 enabled mapping `tag_key` 為空字串或只含空白」時 disabled。aside chip 從 success（`{enabled}/{total} 已啟用`）換成 warning（`{count} 個 Tag Key 留空`）。

把驗證邏輯抽成 selector `useMappingValidation()` 回傳 `{ canContinue, emptyTagCount, enabledCount, totalCount }`。

## Implementation Contract

#### Behavior

- 進入 `current === 3`：表格列出所有 enabled point；每列依 `POINT_SEMANTIC` 預設 tag_key（line01.temp.inlet 等 8 個循環）、display_name、unit；target_type 預設 `float64`；scale/offset 從 `point._rule_scale` / `point._rule_offset` 帶入。
- 點任一列（非控件區域）：`selectedIdx` 設為該列；該列加 outline blue + bg blue 6%；右側 TransformPreview 切換到該 mapping。
- 在 Tag Key input 改值：`state.mappings[p.id].tag_key` 即時更新；右側預覽底部 `→ tag.{tag_key}` 同步刷新；aside chip 若新值為空則 warning。
- 改 × Scale 或 + Offset：右側預覽算式即時重算。
- 改目標型態：右側預覽 cast 步驟更新，final 數字格式化跟著切換。
- 切換 enabled toggle：row 視覺變淡 + 不算 enabled；繼續按鈕重新評估 gate。
- 點表尾「使用 {currentTagShort} 的 Scale / Offset / 型態」連結：所有 mapping 的 scale/offset/target_type 對齊 selected 列。
- 任一 enabled mapping `tag_key` 為空時：aside chip warning `{N} 個 Tag Key 留空`；底部「設定資料庫寫入」disabled。
- 全部 tag_key 填妥、至少 1 個 enabled mapping：底部按鈕啟用；點擊呼叫 `onContinue`。

#### Interface / Data Shape

- `Mapping` 型別（state/types.ts 既有）：
  ```ts
  type Mapping = {
    point_id: string;
    tag_key: string;
    display_name: string;
    unit: string;
    target_type: 'bool' | 'int16' | 'int32' | 'int64' | 'uint16' | 'uint32' | 'uint64' | 'float32' | 'float64' | 'string';
    scale: number;
    offset: number;
    enabled: boolean;
  };
  ```
- Reducer mapping actions：
  - `{ type: 'initMappingsForPoints'; points: Point[] }`
  - `{ type: 'updateMapping'; pointId: string; patch: Partial<Mapping> }`
  - `{ type: 'toggleMappingEnabled'; pointId: string }`
  - `{ type: 'bulkApplyTransform'; fromPointId: string; fields: ('scale' | 'offset' | 'target_type')[] }`
- 純函式：
  - `runScale(raw: number, scale: number, offset: number): number`
  - `castValue(scaled: number, targetType: Mapping['target_type']): number | boolean | string`
  - `formatFinal(value: number | boolean | string, targetType: Mapping['target_type']): string`
  - `buildPayload(mapping: Mapping, point: Point): object`
  - `buildDefaultMapping(point: Point, idx: number): Mapping`
- Selectors：
  - `useMappingValidation(): { canContinue: boolean; emptyTagCount: number; enabledCount: number; totalCount: number }`
  - `useSelectedMapping(): { point: Point; mapping: Mapping; rawSeed: number } | null`

#### Failure Modes

- enabled point 與 mapping 不同步（例如 Step 2 改了 count 把舊 mapping 留下）：reducer `initMappingsForPoints` 同時做孤兒清理，移除 `state.mappings` 中對應已不存在 point 的 entry。
- tag_key 含空白前後：trim 後判斷是否為空；不主動 trim 寫入 state（保留使用者輸入），只用 `.trim() === ''` 判斷 emptiness。
- scale=0：runScale 結果為 offset；UI 不阻擋（這是合法值，但會讓所有點位輸出同值）。
- target_type 改為 bool 後 scaled 為小數：castValue 回傳 `scaled !== 0`，預覽顯示 `true`/`false`。
- 批次套用時 selected 列被刪：connector 變更後 selected 失效；reducer 在 selected 列不存在時靜默 noop。

#### Acceptance Criteria

- `frontend/tests/workbench-v2/transformPipeline.test.ts` 覆蓋：
  - `runScale(243, 0.1, 0) === 24.3`。
  - `castValue(24.3, 'int16') === 24`。
  - `castValue(0, 'bool') === false`、`castValue(1, 'bool') === true`。
  - `formatFinal(24.3456, 'float64') === '24.35'`。
  - `formatFinal(true, 'bool') === 'true'`。
  - `buildPayload(mapping, point)` 結構含 `point_id`、`tag_id`、`transform_pipeline`（3 step）、`enabled`。
- `frontend/tests/workbench-v2/step3-mapping.test.tsx` 覆蓋：
  - render 預設 8 個 enabled point → 表格 8 列、每列依 `POINT_SEMANTIC` 帶預設 tag_key。
  - 改 Tag Key input → state.mappings 更新、右側預覽 `→ tag.xxx` 更新。
  - 改 scale 從 0.1 到 0.5 → 右側預覽算式從 `243 × 0.1` 變 `243 × 0.5`。
  - 點列選取 → selectedIdx 變該列、右側預覽切過去；點列內 input → selectedIdx 不變（stopPropagation）。
  - 點批次套用連結 → 所有 mapping 的 scale/offset/target_type 對齊 selected。
  - tag_key 留空 → 繼續按鈕 disabled、aside chip warning。
  - 全部填妥 → 繼續按鈕啟用、onContinue 觸發。
- `cd frontend && npm run lint && npm run typecheck && npm run test -- --run workbench-v2 && npm run build` 全綠。
- 手動 smoke：完整跑「進入 Step 3 → 改 Tag Key → 改 Scale → 點批次套用 → 切到不同點位看預覽 → 把一筆 Tag Key 清空看繼續鎖住 → 補回 → 繼續到 Step 4 placeholder」。

#### Scope Boundaries

**In scope:**

- Step 3 完整 UI（6 個元件）。
- `state/mappingDefaults.ts` + `state/transformPipeline.ts` 純函式。
- reducer 4 個 mapping actions + 1 個 selector hook。
- `state/types.ts` 補完 Mapping 型別。
- shell 換掉 Step3MappingPlaceholder；shell spec 縮窄 placeholder。
- i18n `workbench-v2` 加入 Step 3 microcopy。
- unit test。

**Out of scope:**

- 後端 API。
- transform_pipeline 進階步驟（lookup / conditional / formula）。
- mapping 拖拉、匯入 / 匯出。
- 多列同時選取。
- E2E。
- Step 1/2/4/Settings 行為。

## Risks / Trade-offs

- [inline 編輯與列選取共存的 click 行為可能跨瀏覽器差異] → 同時對 `onClick` 與 `onMouseDown` 都加 stopPropagation；測試覆蓋三種輸入元素（input/select/toggle）。
- [自動 init 用 useEffect 依賴 point id 序列字串可能在 1000+ point 時造成 string 比對開銷] → 預期 enabled point ≤ 256（單裝置 64 reg × 4 條規則），可接受；若未來放寬可改用 ref 比對。
- [POINT_SEMANTIC 8 個語意循環在 9+ 個點位時會重複] → 預期行為；使用者通常會改名；UI hint 提示「重複 Tag Key 會在後端被視為相同 Tag」。
- [castValue 對 string target_type 直接 String() 失去精度] → 預期，string target 用於描述性欄位（如 status text），不負責保留浮點精度。

## Migration Plan

1. 在 `state/types.ts` 補完 Mapping 型別。
2. 建立 `state/transformPipeline.ts`，TDD：先寫 transformPipeline.test.ts 對 4 個函式各 case。
3. 建立 `state/mappingDefaults.ts` 落地 POINT_SEMANTIC + RAW_VALUE_SEEDS + buildDefaultMapping。
4. 擴充 reducer 加入 4 個 mapping actions + selectors，補 reducer-step3.test。
5. 從子元件實作：`PipelineSteps` → `PayloadPreview` → `TransformPreview`。
6. `MappingRow` → `MappingTable`。
7. 組裝 `Step3Mapping.tsx`，覆寫 shell `current === 3`。
8. 補 i18n。
9. 跑 lint/typecheck/test/build/check-lines。
10. 手動 smoke。

**Rollback strategy**：不刪除 placeholder；rollback 把 shell `current === 3` 改回 placeholder、移除 step3 目錄。

## Open Questions

- 是否要顯示「scale 為 0 提示」？傾向 **不**，原型沒提示，使用者明知 0 的後果。
- 批次套用是否要支援「只套用啟用列」？傾向 **不**，套用後 disabled 列若再啟用會繼承新值，符合直覺。
