## Why

`datalink-workbench-v2-shell`、`datalink-workbench-v2-step1-device`、`datalink-workbench-v2-step2-rule` 三個 change 已落地 shell + Step 1 + Step 2。本 change 把 Step 3「點位映射」完整靜態 UI 移植進 v2 容器，取代 `Step3MappingPlaceholder`。

Step 3 把 Step 2 衍生出來的點位（Point）逐一對應到 Tag（時序資料庫的命名鍵），同時：

- **行內可編輯表格**：每列直接在表內 inline 改 Tag Key / 顯示名稱 / 單位 / 目標型態 / × Scale / + Offset / 啟用 toggle；不再像舊 `/studio` 那樣彈出 modal。
- **轉換管線即時預覽**：點選任一列，右側面板即時呈現該點位的「raw → scale → cast → final → Tag」管線，並以顏色語意（灰/藍/綠/emerald）標示資料流向。
- **個別轉換覆寫**：scale/offset 預設值來自 Step 2 rule 層（`point._rule_scale` / `point._rule_offset`），但每個 mapping 可獨立調整。
- **批次套用**：表尾「將選中列的 Scale / Offset / 型態 套用到全部」一鍵同步所有 mapping，省去逐列調整成本。
- **API payload 預覽**：右側面板底部顯示 `POST /mappings` 將送出的 JSON 結構（含 `transform_pipeline`），方便工程師檢視。

本 change 不接後端；後續 backend-wiring 會把 mapping 寫入接到 `POST /tags × N` 與 `POST /mappings × N`。

## What Changes

- 新增 Step 3 完整 UI 於 `frontend/src/features/datalink/workbench-v2/steps/step3/`，2 欄版面：
  - **左 col 8** 點位 → Tag 映射表（`SectionCard` content p-0）：
    - 第一欄 1.5px 規則色點（沿用 Step 2 規則色票）。
    - 點位 / 位址 mono 純文字。
    - Tag Key（28% 寬）/ 顯示名稱 / 單位：透明 `<input>`，hover 出現邊框，focus 變藍。
    - 目標型態：透明 `<select>` 列 10 種 DATA_TYPES。
    - × Scale / + Offset：`<input type="number" step=0.01 w-16 text-right>`，初值來自 `point._rule_scale` / `point._rule_offset`。
    - 啟用：size sm Toggle。
    - 點任一列（除 input/select/toggle 外）→ setSelectedIdx，row 套 `outline outline-1 outline-blue-500/30 bg-blue-500/[0.06]`。
  - **右 col 4** 轉換管線預覽（`SectionCard`）：
    - 副標：選中點位名稱 + 位址。
    - 4 步驟 `<ol>`（連續直線連接）：
      1. `decode (raw)`：模擬 raw 值（種子 `[243, 251, 1024, 985, 67, 542, 18, 1450]` 循環取）。
      2. `scale (linear)`：算式 `{raw} × {scale} + {offset} = {scaled}`，顏色語意 raw=灰 / scale=藍 / 結果=綠粗體。
      3. `cast to {target_type}`：`cast({scaled} → {target_type})`。
      4. `final → Tag`：emerald 強調卡，巨大 mono 數字依 target_type 格式化（bool → true/false，floatN → toFixed(2)，其他 → Math.round），單位顯示在右下，底部 `→ tag.{tag_key}`。
    - API payload 預覽（`<pre>` mono）：完整 `POST /mappings` JSON 含 `transform_pipeline`。
- **mappings 自動初始化**：useEffect 依 `points.length` 觸發；對每個 enabled point 若 `state.mappings[p.id]` 不存在，依 `POINT_SEMANTIC[i % 8]` 配出 default mapping（tag_key 走 `line01.temp.inlet` 等預設語意；scale/offset 從 `_rule_scale` / `_rule_offset` 帶入；target_type 預設 `float64`）。
- **批次套用**：表尾連結點擊後把 selected 列的 `scale` / `offset` / `target_type` 同步寫入所有 mapping。
- **繼續閘**：底部「設定資料庫寫入」按鈕在「任一 enabled mapping `tag_key` 留空」時 disabled；aside chip 顯示 `{enabled}/{total} 已啟用`（warning tone 當有空 tag_key）。
- 擴充 `useWorkbenchV2State` reducer 加入 mapping actions（`upsertMapping`、`updateMapping`、`bulkApplyTransform`、`toggleMappingEnabled`）。
- 在 `state/mappingDefaults.ts` 落地 `POINT_SEMANTIC`（8 個預設語意：TEMP_IN / TEMP_OUT / PRES_01 / PRES_02 / FLOW_01 / HUM_01 / VIB_01 / RPM_01）與 `RAW_VALUE_SEEDS`（8 個模擬 raw seed）。
- 在 `state/transformPipeline.ts` 落地純函式：`runScale(raw, scale, offset)`、`castValue(scaled, targetType)`、`formatFinal(value, targetType)`、`buildPayload(mapping, point)`。
- 修改 `WorkbenchV2Shell` 在 `current === 3` 時改 render `Step3Mapping`；shell spec MODIFIED placeholder 從名單移除 Step 3。
- 擴充 i18n `workbench-v2` namespace 加入 Step 3 microcopy（column 名、轉換步驟標題、按鈕、警告 chip）。

## Non-Goals (optional)

- 不接 `POST /tags` / `POST /mappings` 真實後端；payload 預覽是顯示文字，不會發送。
- 不在本 change 處理 transform_pipeline 進階步驟（lookup table、conditional rules、formula expressions）；只實作原型涵蓋的 decode → scale → cast → final 四步驟，與 `mapping-pipeline` spec 中其他 transform 步驟的整合留給後續 change。
- 不變更 `mapping-pipeline` spec（既有 spec 保持不動）。
- 不處理 mapping 拖拉排序、批次匯入 / 匯出。
- 不變更 Step 1/2/4/Settings 行為。

## Capabilities

### New Capabilities

- `datalink-workbench-v2-step3-mapping`：定義 v2 Step 3 點位 → Tag 映射工作區的需求、行內表格編輯契約、轉換管線預覽（decode → scale → cast → final）規則、自動 mapping 初始化、批次套用行為、API payload 預覽結構、繼續閘條件。

### Modified Capabilities

- `datalink-workbench-v2-shell`：縮窄 `Placeholder step and settings surfaces` 需求，把 Step 3 從 placeholder 名單移除（剩餘 Step 4 與 Settings 仍走 placeholder）。

## Impact

- Affected specs:
  - 新增 `openspec/specs/datalink-workbench-v2-step3-mapping/spec.md`
  - 修改 `openspec/specs/datalink-workbench-v2-shell/spec.md`
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step3/MappingTable.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step3/PayloadPreview.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step3/index.ts`
    - `frontend/src/features/datalink/workbench-v2/state/mappingDefaults.ts`
    - `frontend/src/features/datalink/workbench-v2/state/transformPipeline.ts`
    - `frontend/tests/workbench-v2/step3-mapping.test.tsx`
    - `frontend/tests/workbench-v2/transformPipeline.test.ts`
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`（新增 mapping actions + auto-init effect 對應 reducer hook）
    - `frontend/src/features/datalink/workbench-v2/state/types.ts`（補完 `Mapping` 型別）
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`（current=3 切到 Step3Mapping）
    - `frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx`（標記為 rollback only）
    - `frontend/src/i18n/locales/zh-TW/workbench-v2.json`、`frontend/src/i18n/locales/en/workbench-v2.json`
- 不變更後端任何檔案。
