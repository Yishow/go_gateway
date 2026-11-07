## Problem

完整 Go 驗證中的 524 個 skip 混合了外部 fixture 缺失、刻意的環境限制，以及三個無條件停用的 handler 測試。設備啟用測試沒有可重現的本機 peer；重複輪詢群組名稱測試只有 TODO。

## Root Cause

兩個設備測試依賴固定 502 埠；輪詢群組的 SQL 儲存庫已有名稱 UNIQUE 約束，記憶體儲存庫卻沒有同等檢查。實體設備測試還會在一般完整測試時自動探測固定目標。

## Proposed Solution

- 以可丟棄的 loopback Modbus peer 啟用設備 handler 測試，驗證回應與持久狀態。
- 對齊記憶體儲存庫的名稱唯一性，補上新增、修改與失敗不變更資料的回歸測試。
- 修正資料庫 successor proof 測試將 stale revision 誤當 lineage corruption 的 fixture；保留真 corruption 的 500 斷言，另測 stale CAS 的 409。
- 實體 PLC 測試改為明確 opt-in，保留外部 fixture、平台限制與子程序 helper 的真實 skip 原因。
- New capability document: polling-group-name-integrity（補齊 datalink-api 既有 Polling Groups API 的名稱約束）；只補既有輪詢群組名稱唯一契約，不修改 HTTP error 格式。

## Non-Goals (optional)

不移除合法 skip、不將未執行視為通過、不連真 PLC、不新增資料庫能力或改寫 UI；資料庫 fixture 補測由獨立驗證工作處理。

## Success Criteria

三個停用測試實際執行通過；重複名稱回歸先失敗後修復；受影響 package、lint、vet 與行數檢查通過。524 原始 skip 依測試事件精確分類，可追到原因。

## Impact

- Affected code:
  - Modified: internal/api/handlers/device_handler_extended_test.go, internal/api/handlers/polling_group_handler_test.go, internal/datalink/pollinggroup/memory_repo.go, internal/datalink/real_device_integration_test.go
  - Modified: internal/api/router_studio_v2_successor_activation_postgres_proof_test.go, internal/api/router_studio_v2_successor_mysql_proof_test.go, internal/api/router_studio_v2_successor_sqlserver_proof_test.go
  - New: internal/api/handlers/device_activation_fixture_test.go, internal/api/handlers/polling_group_duplicate_test.go, internal/datalink/pollinggroup/memory_repo_name_test.go
  - New: docs/testing/skip-audit-20261107.md
  - Removed: none
- Affected specs: polling-group-name-integrity（新增契約文件；不改變既有 SQL 行為）
