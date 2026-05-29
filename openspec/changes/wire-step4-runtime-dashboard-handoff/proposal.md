## Why

Step 4 的成功卡目前已經顯示 `前往 Runtime Dashboard`，但實際 callback 仍只是 `console.log`，沒有真正把操作者從 setup flow 帶到 post-setup runtime dashboard。若這個 handoff 不補上，前兩個 change 開完之後，產品仍會停在「有頁面、有 API、但沒有入口」的半完成狀態。

## What Changes

- 將 Step 4 成功卡的 `前往 Runtime Dashboard` 按鈕從 placeholder callback 接成真正的 route handoff。
- 定義從現有 workbench state 解析 dashboard `device_id` 的 deterministic 規則，讓 commit 後預設打開剛完成設定的 device。
- 若無法安全解析單一 device context，仍導向 runtime dashboard route，但保留 route change 定義的 `missing-device-context` 行為，而不是靜默無動作。
- 補齊 Step 4 commit flow 與 shell routing 測試，確保成功卡入口能被後續實作穩定依賴。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step4-database`: 補強 commit completion card 的 handoff 語意，讓 success card 的 Runtime Dashboard 按鈕導向新的 post-setup runtime dashboard route。

## Impact

- Affected specs: `datalink-workbench-v2-step4-database`
- Affected code:
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx`
    - `frontend/tests/unit/workbench-v2/step4-commit.test.tsx`
    - `frontend/tests/unit/workbench-v2/shell.test.tsx`
  - New:
    - `frontend/src/features/datalink/workbench-v2/state/resolveRuntimeDashboardDevice.ts`
    - `frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts`
