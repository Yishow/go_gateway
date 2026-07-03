## Why

目前 runtime dashboard 與 runtime-related UI 仍有機會顯示推測值、fallback 值、或與 workspace context 不一致的資料，這會讓操作員看到「看起來在跑」但其實沒有對應到真實 runtime/back-end 契約的畫面。對你要求的「可運行、可觀察」來說，這是高風險，因為假資料比明確的空白或失敗更危險。

## What Changes

- 移除 runtime dashboard 與 related UI 中任何不是由後端正式 contract 支撐的 synthetic/fallback runtime 值。
- 將 runtime dashboard 的 selected device、snapshot、stream、degraded state 全部對齊 workspace-scoped runtime context。
- 當後端沒有真實資料時，顯示 explicit empty/degraded/error state，而不是虛構正常數據。

## Capabilities

### New Capabilities

- `runtime-dashboard-truthfulness`: 定義 runtime dashboard 僅能顯示 backend-backed runtime data，並在缺資料時顯式降級，而不是回退到 synthetic values。

### Modified Capabilities

- `runtime-dashboard`: runtime dashboard 畫面必須移除假資料與未定義 fallback。
- `runtime-dashboard-backend`: backend contract 必須明確表示無資料、degraded、stream unavailable、projection stale 等狀態。
- `post-setup-runtime-dashboard`: runtime route 進入後的空狀態與錯誤狀態必須為真實 contract 驅動。
- `studio-runtime-workspace-view`: workspace-scoped runtime view 必須與 selected device context 對齊，不能再退回 synthetic summary。

## Impact

- Affected specs:
  - New: runtime-dashboard-truthfulness
  - Modified: runtime-dashboard, runtime-dashboard-backend, post-setup-runtime-dashboard, studio-runtime-workspace-view
- Affected code:
  - Modified: frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - Modified: frontend/src/features/datalink/runtime-dashboard/useRuntimeStatus.ts
  - Modified: frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts
  - Modified: frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/CollectorHealthPanel.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/LivePointsTable.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/LiveStateBanner.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx
  - Modified: internal/api/handlers/runtime_handler.go
  - Modified: internal/api/handlers/runtime_stream_handler.go
  - Modified: internal/datalink/runtime/status.go
  - Modified: internal/datalink/runtime/stream.go
