## Why

`/studio/runtime` 目前主要仍是單台 `device_id` 視角，但使用者已決定它要承接「V2 工作區裡目前所有已成功存進後端的設備」，預設先開 V2 排列最前面的成功設備，並能切換其它設備。若沒有 workspace 視角，前面多台設備 autosave 與 partial activation 的結果無法被正確觀察。

## What Changes

- 讓 `/studio/runtime` 以 singleton V2 workspace 為主載入設備集合
- 預設選擇 V2 排列最前面的可用設備
- 讓使用者在 runtime 頁內切換工作區中的其它設備
- 沒有可用設備時保留在 runtime 顯示空狀態，不自動跳回 `/studio/v2`
- 保留 unavailable device 在清單中，讓使用者知道它還存在但目前有錯

## Capabilities

### New Capabilities

- `studio-runtime-workspace-view`: 定義 runtime 以 V2 workspace 為主的設備集合、預設選擇、切換與空狀態

### Modified Capabilities

- `datalink-api`: 新增 workspace runtime context API 與 ordered device payload

## Impact

- Affected specs: `studio-runtime-workspace-view`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/runtime-dashboard/`, `frontend/src/App.tsx`
  - New: `frontend/src/services/studioV2RuntimeContext.ts`, `frontend/src/hooks/datalink/useStudioV2RuntimeContext.ts`
  - Modified: `internal/api/router.go`, `internal/api/handlers/runtime_handler.go`, `internal/datalink/runtime/`, `internal/datalink/workspace/`
  - Removed: none
