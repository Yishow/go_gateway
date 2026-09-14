# 完成真實欄位查詢與受確認保護的建表

## Why
F06/F07/F11 顯示：欄位來自 SAMPLE_DB_TABLES、舊預覽可能回寫新畫面，而現有 ApplyManagedSchema 沒有完整核對方案與目標版本。一般 activation 的 EnsureWorkspaceSchema 會發出非預覽建表要求；這與既有「先預覽再明確建立」規格需要收斂。不能只把現成函式接到按鈕就視為完成。

## What Changes
- 真正讀取所選連線、資料庫及資料表，分清楚表不存在、無權查詢與查詢失敗。
- 以版本綁定且會過期的預覽確認資料保護建表；預覽本身零副作用。
- 真正操作前再次檢查歸屬、版本、目標與允許功能，支援重複要求及不確定結果查詢。
- 現成表不破壞、不自動改名刪欄；啟用設備只檢查已確認的結構，不暗中建表。

## Capabilities
### New Capabilities
None.
### Modified Capabilities
- `recording-database-setup`: 更新 Capability-backed database preparation、Revision-bound schema preview and explicit creation；新增失敗與重試的建表結果規範。

## Impact
定位：`internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`、`studio_v2_workspace_database_handler.go`、`internal/datalink/recordingplan/service.go`、`internal/datalink/dbtarget/`、`frontend/src/services/datalink.ts`、`studioV2WorkspaceDatabase.ts`、Step 4 的 `SchemaSetupSection.tsx` 與欄位資料來源。需要持久化的 preview/operation 記錄與 API 文件、真資料庫測試。

## Non-goals
不執行任意使用者 SQL，不修改不相容既有表，不支援無法證明安全的 adapter，不重新設計資料採集或歷史保留政策。

## Dependencies and Delivery
前置 C1、C2；完成本案驗證後只開放已通過的 schema operation，試寫保持 C1 封鎖直到 C4 完成。C3 與 C4 都改 recordingplan service，必須串行。本次不建表、不接正式資料庫。
