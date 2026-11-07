## Why

補跑先前跳過的真資料庫測試後，發現 fixture 沒完整保存資料庫證據或會受環境時間影響，造成不實失敗、清理失敗與後續跳過。修正測試自身，保留所有產品安全驗證與精確斷言。

## What Changes

- Timescale extension 清理改用新連線池重建，避免重用已載入 extension 的 backend 而失敗。
- 同一批 recovery fixture 的群組 Apply 使用同一個 clock，離開 fixture 後還原，避免跨分鐘邊界造成不同生效時間。
- SQL Server 測試 converter 傳遞 NativeMonthly partition 證據，與既有完整 converter 一致。
- SQL Server runtime timestamp 以 RFC3339Nano 解析後 time.Equal 比較同一瞬間，不因合法字串表示差異失敗。

## Non-Goals (optional)

不修改產品行為、不放寬精度、不移除 cleanup/assert、不將 skip 當 pass、不修理或清理正式資料庫。

## Impact

- Affected specs: none
- Affected code:
  - New: none
  - Modified: internal/api/router_studio_v2_write_group_recovery_denied_pg_test.go, internal/api/router_studio_v2_write_group_recovery_rig_test.go, internal/api/router_studio_v2_sqlserver_circuit_schema_fault_test.go, internal/api/router_studio_v2_successor_sqlserver_runtime_test.go
  - Removed: none
- Compatibility: no capability-level observable behavior changes
