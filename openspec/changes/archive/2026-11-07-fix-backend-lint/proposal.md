## Why

完整後端 lint 實測有 39 項既有問題，讓發布檢查無法清楚判斷新增缺陷。此變更清除已確認的問題，維持設備採集與資料寫入契約。

## What Changes

- 修正未檢查的 rollback、型別斷言、整數窄化，補相關邊界測試。
- 整理等價 switch、命名回傳、常數、context、預配置與無用程式。
- 執行完整後端測試、vet 與 lint，記錄真實結果及跳過條件。

## Non-Goals

不變更 lint 規則，不新增忽略，不改 UI 或資料庫功能，不接觸正式設備。既有無條件 skip 的修復由獨立變更負責。

## Impact

- Affected specs: none
- Affected code:
  - Modified: cmd/test_ui、internal/api/handlers、internal/datalink、internal/protocol/mcprotocol 的 lint 回報檔案及相關回歸測試。
  - Removed: 無呼叫的 rulePointPlan.keeps；flushDue 未使用的 all 選項。
- Compatibility: no capability-level observable behavior changes
