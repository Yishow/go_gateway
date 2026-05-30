## Why

使用者已明確決定 Step 4 的最後按鈕不再負責資料寫入，而是只負責第一次啟動尚未啟動、目前合法且可用的設備。現在的 Step 4 仍以 commit/成功卡為主語意，如果不把這件事拆成單獨 change，前面 Step 1~4 的 autosave 會再次被最後一顆按鈕綁回「整批 commit」心智。

## What Changes

- 將 Step 4 最後按鈕改成第一次啟動 flow
- 啟動 singleton workspace 內所有合法、可用、尚未啟動的設備
- 回傳逐台成功/失敗結果，不要求整批成功
- 即使部分失敗，仍允許前往 runtime 觀察成功設備

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step4-database`: 將 Step 4 從 commit flow 改成 first activation flow
- `datalink-api`: 新增 workspace first activation API 與逐台結果契約

## Impact

- Affected specs: `datalink-workbench-v2-step4-database`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step4/`, `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`
  - Modified: `internal/api/router.go`, `internal/api/handlers/runtime_handler.go`, `internal/datalink/runtime/`, `internal/datalink/device/`
  - New: `frontend/src/types/studioV2Activation.ts`
  - Removed: none
