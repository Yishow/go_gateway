## Why

`/studio/v2/settings` 仍保留明確 mock / noop 行為，包括 connector test 的隨機成功失敗與 `儲存所有設定` 的 warning-only noop。若目標是 `/studio/v2` 全站對接，settings 不能繼續停在 demo-grade local state。

## What Changes

- 讓 settings page 啟動時從 backend 讀取真實資料
- 將 connector pool CRUD / test 接到既有 backend connector APIs
- 將 save bar 從 noop 改成真實 persistence flow
- 移除 operator-visible 的 settings mock / noop 行為

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-settings`: 由 demo state 改成 backend-backed settings surface
- `datalink-api`: 明確收斂 V2 settings 所需的 read / write / connector test 契約

## Impact

- Affected specs: `datalink-workbench-v2-settings`, `datalink-api`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/settings/`, `frontend/src/features/datalink/workbench-v2/state/settingsReducer.ts`
  - Modified: `frontend/src/hooks/datalink/`, `frontend/src/services/`
  - Modified: `internal/api/router.go`, `internal/api/handlers/settings_handler.go`, `internal/api/handlers/dbtarget_handler.go`
  - Modified: `internal/datalink/settings/`, `internal/datalink/dbtarget/`
  - Removed: settings mock test / noop save behavior
