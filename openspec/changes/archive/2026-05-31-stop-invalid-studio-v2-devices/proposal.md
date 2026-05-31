## Why

使用者已經明確決定：如果一台原本正在跑的設備被改成不合法，它要立刻視為不可用，runtime 不再把它當可用設備，但仍要讓使用者看得到這台設備現在出了問題。這需要明確的「停止規則」與可觀察狀態，不然 autosave 只會留下模糊的半套行為。

## What Changes

- 定義 running device 因 invalid edit 轉成不可用時的立即停止規則
- 定義 workspace/device status 中的 availability state 與 reason
- 固定 invalid edit 不得繼續沿用舊的 running 狀態
- 補齊前後端對 unavailable device 的狀態表達

## Capabilities

### New Capabilities

- `studio-v2-device-validity`: 定義 V2 設備在 invalid edit 後的不可用與停止行為

### Modified Capabilities

- `datalink-api`: 擴充 workspace device status / runtime status 的 availability 欄位

## Impact

- Affected specs: `studio-v2-device-validity`, `datalink-api`
- Affected code:
  - Modified: `internal/datalink/runtime/`, `internal/datalink/device/`, `internal/api/handlers/`
  - Modified: `frontend/src/features/datalink/workbench-v2/`, `frontend/src/features/datalink/runtime-dashboard/`
  - New: `frontend/src/types/studioV2Availability.ts`
  - Removed: none
