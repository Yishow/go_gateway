## Why

使用者已將目標收斂為 `/studio/v2` 全站對接，但 Step 1 的 `執行測試` 目前仍是前端 mock animation。這讓畫面在沒有真實設備、沒有真實連線、甚至沒有任何後端診斷結果時也可能顯示通過，與「前端完整接後端」目標直接衝突。

## What Changes

- 將 Step 1 `執行測試` 從 mock animation 改成真實後端 diagnostics flow
- 讓 connect / probe 階段狀態來自 backend response，而不是本地 timer 與隨機值
- 讓 Step 1 continue gate 只接受真實成功的 diagnostics 結果
- 保留 Step 1 autosave 與 live diagnostics 的責任分離，不把測試流程再塞回 Step 4

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step1-device`: 將測試面板從 mock 成功動畫改成真實後端診斷
- `datalink-api`: 提供或收斂 V2 可直接使用的 Step 1 diagnostics 契約

## Impact

- Affected specs: `datalink-workbench-v2-step1-device`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step1/`, `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`
  - Modified: `frontend/src/hooks/datalink/`, `frontend/src/services/`
  - Modified: `internal/api/router.go`, `internal/api/handlers/device_handler.go`, `internal/datalink/device/`
  - Removed: mock-only Step 1 diagnostics behavior
