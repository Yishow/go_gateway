## Why

Step 4 的 recording plan、measurement membership、row group 與 target mapping 同時決定能否開始，造成 managed 路線仍被另一套欄位計數阻擋。先建立唯一寫入群組 authority 與相容移轉，才能在既有 Go／React 與 /studio/v2 內收斂，而不重寫整個產品。

## What Changes

- 新增版本化 WriteGroup，連接已儲存 workspace/device/point/tag 與 connector/table/columns，保存寫入週期、資料列識別及適用狀態。
- managed／custom 是同一群組的 storage strategy；measurement 語意只在使用者選擇衍生運算時要求，基本 Tag 寫入不要求報表、retention、聚合或電量差分。
- 逐批匯入既有 single-point、row-group、target mappings 與可無損映射的 recording plans；模糊或衍生模式保留原資料並要求明確審閱，不猜來源。
- 以交易／CAS 保存版本、一次性切換 applied projection；同一輸出不可同時由舊 writer 與新群組寫兩次。舊 API 經相容 adapter 或安全衝突回覆，不另成 authority。

- 補齊rename／member edit／disable／delete lifecycle；delete以tombstone保留accepted backlog與收據歸屬，不cascade刪除、不重用ID。

## Capabilities

### New Capabilities

- `studio-v2-write-groups`: 群組契約、來源身分、原子儲存、相容移轉與單一輸出 authority

### Modified Capabilities

- `workspace-database-row-groups`: 保留共享欄位與重載結果，將新群組的 row-group UI 轉為 canonical group 的相容 projection

## Impact

internal/datalink/workspace/、internal/datalink/schema/、internal/datalink/recordingplan/、internal/datalink/dbtarget/、internal/api/handlers/studio_v2_workspace_database*、frontend/src/types/ 與 services。保留既有 scope、schema-operation、readiness 與 Share 保護。無前置新 change；移轉／部署只在後續獲授權實作時執行。

本次僅起草文件，沒有產品實作。來源、依賴與移交見 [總覽](../../../docs/plans/studio-v2-write-groups/README.md) 及 [現況證據](../../../docs/plans/studio-v2-write-groups/evidence.md)。
