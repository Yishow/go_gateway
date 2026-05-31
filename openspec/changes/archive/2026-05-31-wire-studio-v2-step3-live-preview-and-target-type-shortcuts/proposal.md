## Why

`/studio/v2` 已經把 Step 1、workspace autosave 與 settings 接上真後端，但 Step 3 的 `轉換管線預覽` 仍是前端本地 mock。這讓 Point → Tag 頁面的 preview 不是 backend 真相，也和「`/studio/v2` 全站對接」目標不一致。

此外，Step 3 操作大量 mapping 列時，`目標型態` 只有逐列 select，缺少快速設定當前列與一鍵套用到全部列的高頻操作。

## What Changes

- 將 Step 3 `轉換管線預覽` 改接既有 `POST /api/v1/datalink/mappings/preview`
- 保留 deterministic sample raw value，但改由 backend 執行 transform pipeline
- 為 preview 補上 loading、error、stale response guard
- 保持 `API payload preview` 為 display-only，不新增額外 request
- 在 Step 3 preview 卡新增 `目標型態` 快捷設定，支援：
  - 當前列快速設定
  - 一鍵套用到全部列

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step3-mapping`: 將右側 transform preview 從本地 mock 改成 backend preview，並新增目標型態快捷設定
- `datalink-api`: 正式納入 Step 3 使用的 mapping preview API 契約

## Impact

- Affected specs: `datalink-workbench-v2-step3-mapping`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step3/Step3Mapping.tsx`
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step3/TransformPreview.tsx`
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step3/PipelineSteps.tsx`
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx`
  - Modified: `frontend/src/features/datalink/workbench-v2/state/selectors.ts`
  - Modified: `frontend/src/i18n/locales/zh-TW/workbench-v2.json`
  - Modified: `frontend/src/i18n/locales/en/workbench-v2.json`
  - Modified: `frontend/tests/unit/workbench-v2/step3-mapping.test.tsx`
  - Modified: `openspec/specs/datalink-workbench-v2-step3-mapping/spec.md`
  - Modified: `openspec/specs/datalink-api/spec.md`
