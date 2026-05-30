## Why

後端監看契約就緒後，產品仍缺少一個真正可用的 post-setup runtime 入口。若沒有獨立 route，操作者在 Step 4 commit 後仍只能停在 setup flow 或回退舊 dashboard，無法用同一套 v2 語彙直接確認「這台剛上線的 device 現在到底有沒有在跑」。

## What Changes

- 新增一個專供 post-setup 使用的 runtime dashboard route，承接 v2 setup 完成後的監看需求。
- 讓 runtime dashboard 以單一 `device_id` 作為預設上下文，首屏聚焦剛完成設定的 device，而不是 fleet-wide 首頁。
- 以 backend runtime contract 為唯一資料來源，採 snapshot-first + SSE live update 的載入方式，並定義 stream 不可用時的 polling fallback。
- 提供在 dashboard 內切換其他 device 的能力，但不把頁面改造成 generic KPI dashboard 或新的 `/studio/v2` 首頁。
- 這個 change 只擁有 route、query-state 與 route-level dashboard state contract；真正的 page surface 與 panels 由 `runtime-dashboard` change 擁有。

## Capabilities

### New Capabilities

- `post-setup-runtime-dashboard`: 定義 post-setup runtime dashboard 的 route、device context、資料載入與降級行為。

### Modified Capabilities

(none)

## Impact

- Affected specs: `post-setup-runtime-dashboard`
- Affected code:
  - New:
    - `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx`
    - `frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts`
    - `frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts`
    - `frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts`
    - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx`
    - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`
  - Modified:
    - `frontend/src/App.tsx`
    - `frontend/src/services/datalink.ts`
