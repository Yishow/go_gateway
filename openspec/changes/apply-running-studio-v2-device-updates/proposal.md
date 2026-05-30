## Why

使用者已經明確決定：設備第一次啟動後，只要之後的合法自動存成功，就應直接套用到正在執行的設備，不需要再回 Step 4 按一次。若沒有這個契約，V2 的「每一步直接接後端」會停在保存層，無法真正反映到 runtime。

## What Changes

- 定義已啟動設備在合法 autosave 成功後的直接套用行為
- 補齊 autosave 成功回應中的 runtime apply 狀態
- 固定「未啟動設備 autosave 不會偷啟動」的邊界
- 讓前端可以看見哪次 autosave 已成功套用到正在執行的設備

## Capabilities

### New Capabilities

- `studio-v2-live-config-apply`: 定義已啟動設備在合法 autosave 成功後如何直接套用新設定

### Modified Capabilities

- `datalink-api`: 擴充 V2 autosave 回應中的 runtime apply 狀態

## Impact

- Affected specs: `studio-v2-live-config-apply`, `datalink-api`
- Affected code:
  - Modified: `internal/datalink/runtime/`, `internal/datalink/device/`, `internal/datalink/sourcerule/`, `internal/datalink/mapping/`, `internal/datalink/dbtarget/`
  - Modified: `internal/api/handlers/`, `frontend/src/hooks/datalink/`, `frontend/src/features/datalink/workbench-v2/`
  - New: `frontend/src/types/studioV2RuntimeApply.ts`
  - Removed: none
