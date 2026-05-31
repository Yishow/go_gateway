## Why

使用者已經決定 Step 4 的按鈕不再負責真正寫資料，而 Step 4 裡的 connector / database target 編輯也要像其它 step 一樣合法即自動存。若不把 Step 4 資料保存和第一次啟動拆開，之後很容易又回到「按最後一顆按鈕才真正落地」的舊模式。

## What Changes

- 讓 Step 4 的 connector 與 database target 編輯改成 workspace-scoped 真實 autosave
- 只有合法 database 設定才寫入後端；不合法內容留在畫面上但不覆蓋最後成功版本
- 多個 target 逐項保存、逐項錯誤隔離
- 明確規定 Step 4 autosave 不會直接啟動 runtime

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step4-database`: 增加 connector / target autosave 與 pre-activation boundary
- `datalink-api`: 新增 workspace-scoped Step 4 autosave API 契約

## Impact

- Affected specs: `datalink-workbench-v2-step4-database`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step4/`, `frontend/src/features/datalink/workbench-v2/state/dbReducer.ts`
  - New: `frontend/src/services/studioV2DatabaseTargets.ts`, `frontend/src/hooks/datalink/useStudioV2DatabaseTargets.ts`
  - Modified: `internal/api/router.go`, `internal/api/handlers/db_target_handler.go`, `internal/datalink/dbtarget/`, `internal/datalink/workspace/`
  - Removed: none
