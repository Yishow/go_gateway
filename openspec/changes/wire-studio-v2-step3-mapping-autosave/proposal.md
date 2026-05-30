## Why

Step 3 的 mapping 目前仍是本地表格資料，使用者看不到哪些 mapping 已經進到後端、哪些只是暫時輸入值。若不把 mapping 層變成合法即自動存，Step 4 之後的 database output 與 runtime 行為都還是建立在假資料上。

## What Changes

- 讓 Step 3 mapping 編輯改成 workspace-scoped 真實後端 autosave
- 只有合法 mapping 列會寫入後端；不合法列留在畫面上但不覆蓋最後成功版本
- 多列 mapping 逐列保存、逐列錯誤隔離
- 補齊前後端對 mapping save state 的契約

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step3-mapping`: 增加逐列 autosave 與逐列錯誤隔離
- `datalink-api`: 新增 workspace-scoped mapping autosave API 契約

## Impact

- Affected specs: `datalink-workbench-v2-step3-mapping`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step3/`, `frontend/src/features/datalink/workbench-v2/state/mappingReducer.ts`
  - New: `frontend/src/services/studioV2Mappings.ts`, `frontend/src/hooks/datalink/useStudioV2Mappings.ts`
  - Modified: `internal/api/router.go`, `internal/api/handlers/mapping_handler.go`, `internal/datalink/mapping/`, `internal/datalink/workspace/`
  - Removed: none
