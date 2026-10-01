## Why

recording-plans/test-write 目前以501安全封鎖，既有schema operation ledger已有scope、revision與confirm保護。需沿用此機制，讓所選群組真實試寫、讀回與只清理自有測試資料各自有證據，避免把連線成功或buffered當作試寫成功。

## What Changes

- 沿用 schema operation repository，增加test_write種類的preview、claim、status與結果；canonical group為scope，舊plan入口以相容解析器進入相同authority。
- 預覽零外部mutation；確認綁定operation、content digest、workspace/group/connector revisions與expiry，不能拿建表token試寫。
- 用production row encoder/sender向可驗證目標寫入，同operation重送／重啟不重複；uncertain outcome只查同operation。
- 分開written_verified、written_unverified、failed、unknown與cleanup_status；讀回比較型別與實際值，清理只刪本operation擁有的測試列。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `recording-database-setup`: 接手舊active change的Truthful explicit test writes完整需求，加入write-group authority與真實交付證據

## Impact

internal/api/handlers/studio_v2_workspace_recording_plan_operations.go、現有schema operation handlers/repository、internal/datalink/dbtarget/、frontend services/hooks/types。前置：runtime-write-group-delivery及其前置驗證；schema setup已歸檔能力保留。只在後續實作授權下使用可丟棄DB，不碰正式資料。

本次僅起草文件，沒有產品實作。來源、依賴與移交見 [總覽](../../../docs/plans/studio-v2-write-groups/README.md) 及 [現況證據](../../../docs/plans/studio-v2-write-groups/evidence.md)。
