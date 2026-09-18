# 完成真實欄位查詢與受確認保護的建表

## Why
[workflow evidence](../fix-studio-v2-database-workflow/evidence.md) 的 E03／E06／E08 及本次工作樹核對顯示：欄位來自 SAMPLE_DB_TABLES、舊預覽可能回寫新畫面，而現有 ApplyManagedSchema 沒有完整核對方案與目標版本。一般 activation 的 EnsureWorkspaceSchema 會發出非預覽建表要求；這與既有「先預覽再明確建立」規格需要收斂。不能只把現成函式接到按鈕就視為完成。

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
定位：`internal/api/handlers/dbtarget_handler_tooling.go` 的 generic GenerateSchema 與對應 router，`internal/api/router_studio_v2_recording_routes.go`、operation status handler 與 `GET /studio-v2/workspace/database-operations/:operation_id`，`internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`、`studio_v2_workspace_database_handler.go`、`internal/datalink/recordingplan/service.go`、`internal/datalink/recordingplan/` 的 SQLRepository／MemoryRepository／types、`internal/datalink/schema/migrations/`、`internal/datalink/dbtarget/`、`frontend/src/services/datalink.ts`、`studioV2WorkspaceDatabase.ts`、Step 4 的 `SchemaSetupSection.tsx` 與欄位資料來源。需要持久化的 preview/operation 記錄與 API 文件、真資料庫測試。

## Non-goals
不執行任意使用者 SQL，不修改不相容既有表，不支援無法證明安全的 adapter，不重新設計資料採集或歷史保留政策。

## Dependencies and Delivery
前置為 `fix-studio-v2-database-result-truthfulness` 1.1–1.8，以及 `fix-studio-v2-database-workflow` 的 1.2、1.3、2.1、2.2、2.4；需核對已存 connector/member ownership、版本與原子儲存的測試證據，不要求 workflow 整案完成。完成本案 3.1–3.10 後只開放已驗證 schema operation，並解鎖 workflow 2.3 與試寫階段；試寫直到 workflow 3.3／3.4 通過才開放。唯一分工與順序見 [workflow design](../fix-studio-v2-database-workflow/design.md)。取代原來對未發布 C2/C4 的依賴，不新增 change。本次不建表、不接正式資料庫。
