## 1. State 層�- [x] 1.1 在 `state/types.ts` 補完 `Mapping` 型別（含 10 種 target_type union），落地設計決策「拆檔策略：6 個元件 + 2 個 state module」型別契約。**行為**：`import type { Mapping } from '.../state/types'` 可在所有 step3 子元件正確解析。**驗證**：`tsc --noEmit` 通過；新增 type-level assertion `types-step3.test-d.ts`。
- [x] 1.2 建立 `state/transformPipeline.ts`，匯出 `runScale`、`castValue`、`formatFinal`、`buildPayload` 4 個純函式，落地設計決策「轉換管線數值流：使用 raw seed 序列」並實作需求 **Transform pipeline preview** 與 **API payload preview** 的計算部分。**行為**：純函式輸入相同輸出相同；`runScale(243, 0.1, 0) === 24.3`；`castValue(0, 'bool') === false`；`buildPayload` 結構含 `point_id`、`tag_id`、`transform_pipeline`、`enabled`。**驗證**：新增 `frontend/tests/workbench-v2/transformPipeline.test.ts`，4 個函式各 2–3 個 case，覆蓋設計 Decisions 的 format table。
- [x] 1.3 建立 `state/mappingDefaults.ts`，匯出 `POINT_SEMANTIC`（8 個 entry）、`RAW_VALUE_SEEDS`、`buildDefaultMapping(point, idx)`，落地需求 **Automatic mapping initialization** 的 default 來源。**行為**：`POINT_SEMANTIC.length === 8`；`buildDefaultMapping` 用 `idx % 8` 對應 semantic。**驗證**：新增 `frontend/tests/workbench-v2/mappingDefaults.test.ts`。
- [x] 1.4 擴充 `useWorkbenchV2State` reducer 加 4 個 mapping actions（`initMappingsForPoints`、`updateMapping`、`toggleMappingEnabled`、`bulkApplyTransform`），落地設計決策「自動 mapping 初始化：reducer action 而非 component setState」與「批次套用：reducer action」並實作需求 **Automatic mapping initialization** 與 **Bulk apply transform parameters**。**行為**：`initMappingsForPoints` 保留既有、新增缺漏、移除孤兒；`bulkApplyTransform` 把 fromPointId 的 fields 同步到所有 mapping；`updateMapping` 為 immutable patch。**驗證**：新增 `frontend/tests/workbench-v2/reducer-step3.test.ts`，4 個 action 各 2 case；覆蓋孤兒清理、保留既有編輯、批次套用範圍。
- [x] 1.5 在 `state/selectors.ts` 加 `useMappingValidation()` 與 `useSelectedMapping(selectedIdx)`，落地需求 **Continue gate with empty tag_key detection** 與設計決策「繼續閘：tag_key 不能為空」。**行為**：`useMappingValidation()` 回傳 `{ canContinue, emptyTagCount, enabledCount, totalCount }`；`useSelectedMapping(idx)` 回傳 `{ point, mapping, rawSeed }` 或 null。**驗證**：擴充 `selectors.test.tsx` 覆蓋 2 個新 hook。

## 2. UI 子元件

- [x] 2.1 建立 `steps/step3/PipelineSteps.tsx`，渲染 4 步驟 ol（decode / scale / cast / final）+ 圓點 + 1px 連線，落地需求 **Transform pipeline preview** 與 **Pipeline color semantics**。**行為**：每步依 selected mapping 計算顯示；顏色語意 raw=slate-500、scale/offset=blue-300、scaled=emerald-300、final=emerald-200 2xl；final 數字依 target_type 格式化（`formatFinal`）。**驗證**：新增 `frontend/tests/workbench-v2/step3-mapping.test.tsx` 中 `describe('PipelineSteps')` 覆蓋 float / int / bool 三種 target_type 的格式化、4 步驟順序與顏色 class。
- [x] 2.2 建立 `steps/step3/PayloadPreview.tsx`，渲染 `<pre>` mono JSON `POST /mappings` 結構，落地需求 **API payload preview**。**行為**：payload 即時反映 selected mapping；不發 any network。**驗證**：`step3-mapping.test.tsx` 中 `describe('PayloadPreview')` 測 payload 結構正確、spy fetch 斷言未被呼叫。
- [x] 2.3 建立 `steps/step3/TransformPreview.tsx`，組合 SectionCard 副標 + PipelineSteps + PayloadPreview。**行為**：副標顯示 `{point.name} @ {point.address}`。**驗證**：`step3-mapping.test.tsx` 中測切換 selectedIdx 後副標與內容同步。
- [x] 2.4 建立 `steps/step3/MappingRow.tsx`，渲染 1 列含 7 個 inline 控件 + 列點擊 setSelectedIdx，落地需求 **Inline-editable mapping table** 的 row 部分與設計決策「行內編輯 + 列選取共存：stopPropagation 策略」。**行為**：所有 input/select/toggle 加 `onClick`/`onMouseDown stopPropagation`；列點擊 setSelectedIdx；selected 列 outline+bg；rule color dot 從 `useDeviceColor` 取（搭配 rule color cycle helper）。**驗證**：`step3-mapping.test.tsx` 中 `describe('MappingRow')` 覆蓋三種輸入元素點擊不觸發列選取、列空白處點擊觸發選取、控件 onChange 即時 dispatch `updateMapping`。
- [x] 2.5 建立 `steps/step3/MappingTable.tsx`，渲染 thead + tbody（map MappingRow）+ 表尾批次套用連結，落地需求 **Bulk apply transform parameters** 的 UI 部分。**行為**：表尾連結文字依 selected mapping 動態顯示「使用 {tagShort} 的 Scale / Offset / 型態」；selected===null 時連結 disabled。**驗證**：`step3-mapping.test.tsx` 中 `describe('MappingTable')` 覆蓋連結文字動態、bulk apply dispatch、disabled 條件。

## 3. 容器與整合

- [x] 3.1 建立 `steps/step3/Step3Mapping.tsx`，組裝 12-col grid（MappingTable col 8 + TransformPreview col 4）+ 底部 footer（aside chip + 繼續按鈕），含 useEffect 觸發 `initMappingsForPoints`，落地設計決策「自動 mapping 初始化：reducer action 而非 component setState」並實作需求 **Continue gate with empty tag_key detection**。**行為**：mount 觸發 initMappings；繼續按鈕 disabled 條件 `!canContinue`；點繼續呼叫 `onContinue`；aside chip 3 種狀態切換（success / 留空 warning / 無啟用 warning）。**驗證**：`step3-mapping.test.tsx` 主測試 `renders full Step 3 with default state` + 繼續按鈕條件測試 3 個 case。
- [x] 3.2 建立 `steps/step3/index.ts` barrel re-export。**行為**：`from 'steps/step3'` 可取得 Step3Mapping。**驗證**：`tsc --noEmit` 通過。
- [x] 3.3 修改 `shell/WorkbenchV2Shell.tsx` 在 `current === 3 && view === 'flow'` 時改 render `<Step3Mapping onContinue={...} />`，落地 modified 需求 **Placeholder step and settings surfaces**。**行為**：current=3 顯示 Step 3；4 與 settings 仍 placeholder。**驗證**：`shell.test.tsx` 既有 placeholder 測試更新。
- [x] 3.4 在 `steps/Step3MappingPlaceholder.tsx` 加 rollback only 註解。**行為**：保留供 rollback；未被 import。**驗證**：`grep -rn "Step3MappingPlaceholder" frontend/src` 只在自身或註解出現。

## 4. i18n

- [x] 4.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `step3.*` keys：column 名（Tag Key / 顯示名稱 / 單位 / 目標型態 / × Scale / + Offset / 啟用）、轉換步驟標題（decode (raw) / scale (linear) / cast to {target} / final → Tag）、批次套用連結、繼續按鈕、aside chip 三種狀態文字。**行為**：覆蓋。**驗證**：`step3-mapping.test.tsx` 在 zh-TW 斷言 `點位 → Tag 映射` 出現。
- [x] 4.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應翻譯（Point → Tag Mapping、Apply selected row's Scale / Offset / Type to all、Configure Database Write 等）。**行為**：對等翻譯。**驗證**：`step3-mapping.test.tsx` 在 en 斷言 `Point → Tag Mapping` 出現。

## 5. 驗證

- [x] 5.1 跑 `cd frontend && npm run lint`，零錯誤。
- [x] 5.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [x] 5.3 跑 `cd frontend && npm run test -- --run workbench-v2`，全綠。
- [x] 5.4 跑 `cd frontend && npm run build`，產出乾淨。
- [x] 5.5 跑 `bash scripts/check_file_lines.sh`，確認 step3 所有新檔 ≤ 300 行。
- [x] 5.6 手動 smoke：完整跑「進入 Step 3 → 改 Tag Key → 改 Scale → 點批次套用 → 切到不同列看右側預覽 → 清空一筆 Tag Key 看繼續鎖住 → 補回 → 繼續到 Step 4 placeholder」。p3Mapping。**驗證**：`tsc --noEmit` 通過。
- [ ] 3.3 修改 `shell/WorkbenchV2Shell.tsx` 在 `current === 3 && view === 'flow'` 時改 render `<Step3Mapping onContinue={...} />`，落地 modified 需求 **Placeholder step and settings surfaces**。**行為**：current=3 顯示 Step 3；4 與 settings 仍 placeholder。**驗證**：`shell.test.tsx` 既有 placeholder 測試更新。
- [ ] 3.4 在 `steps/Step3MappingPlaceholder.tsx` 加 rollback only 註解。**行為**：保留供 rollback；未被 import。**驗證**：`grep -rn "Step3MappingPlaceholder" frontend/src` 只在自身或註解出現。

## 4. i18n

- [ ] 4.1 擴充 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，新增 `step3.*` keys：column 名（Tag Key / 顯示名稱 / 單位 / 目標型態 / × Scale / + Offset / 啟用）、轉換步驟標題（decode (raw) / scale (linear) / cast to {target} / final → Tag）、批次套用連結、繼續按鈕、aside chip 三種狀態文字。**行為**：覆蓋。**驗證**：`step3-mapping.test.tsx` 在 zh-TW 斷言 `點位 → Tag 映射` 出現。
- [ ] 4.2 擴充 `frontend/src/i18n/locales/en/workbench-v2.json` 對應翻譯（Point → Tag Mapping、Apply selected row's Scale / Offset / Type to all、Configure Database Write 等）。**行為**：對等翻譯。**驗證**：`step3-mapping.test.tsx` 在 en 斷言 `Point → Tag Mapping` 出現。

## 5. 驗證

- [ ] 5.1 跑 `cd frontend && npm run lint`，零錯誤。
- [ ] 5.2 跑 `cd frontend && npm run typecheck`，零錯誤。
- [ ] 5.3 跑 `cd frontend && npm run test -- --run workbench-v2`，全綠。
- [ ] 5.4 跑 `cd frontend && npm run build`，產出乾淨。
- [ ] 5.5 跑 `bash scripts/check_file_lines.sh`，確認 step3 所有新檔 ≤ 300 行。
- [ ] 5.6 手動 smoke：完整跑「進入 Step 3 → 改 Tag Key → 改 Scale → 點批次套用 → 切到不同列看右側預覽 → 清空一筆 Tag Key 看繼續鎖住 → 補回 → 繼續到 Step 4 placeholder」。
