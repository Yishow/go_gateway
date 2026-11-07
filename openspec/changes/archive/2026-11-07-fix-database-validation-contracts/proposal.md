## Problem

先前因缺少資料庫而跳過的整合測試，實際連接可丟棄 PostgreSQL/MariaDB 後揭露兩個產品缺陷與一個 fixture 缺漏：PostgreSQL 初始 migration 無法呼叫月份分區函式；已成功建立的記錄表因 NUMERIC 欄位被當成不相容而回 unknown；MySQL 系列手工 runtime snapshot 缺少不可變 counter binding facts。

## Root Cause

DATE + INTERVAL 會變成 timestamp，但函式參數是 DATE。managed table 驗證器將 val_dec 分類為文字，與 PostgreSQL DDL 的無界 NUMERIC 不一致。手工 fixture 未跟上 Apply 的封存 facts 契約。

## Proposed Solution

- 對下一月分區函式參數明確轉 DATE，測試改用新 schema 與 transaction rollback，避免刪除 public 表。
- val_dec 使用獨立 decimal 分類，只接受既有文字表示與無界 NUMERIC，不放寬浮點或有界小數。
- MySQL/MariaDB fixture 經既有計算與 seal 流程加入 counter bindings，不改產品 fail-closed 行為。
- Modified capability: recording-plan-storage，補齊 PostgreSQL 產生 DDL 與實測驗證一致性。

## Non-Goals (optional)

不改儲存版型、不擴充DB支援、不移除 counter proof、不用 skip 或降低斷言消除失敗、不操作正式資料庫或真 PLC。

## Success Criteria

兩個真 PostgreSQL 失敗測試修復後通過；無界 NUMERIC 接受且有界/浮點拒絕回歸通過；MySQL/MariaDB fixture 實際整合通過，未跑引擎明列。相關 package、vet/lint與行數檢查通過。

## Impact

- Affected code:
  - Modified: internal/datalink/schema/migrations/001_initial_schema.up.sql, internal/datalink/storage/postgres_test.go, internal/datalink/recordingplan/schema_tables.go, internal/api/router_studio_v2_mysql_circuit_plain_fixture_test.go
  - New: internal/datalink/recordingplan/schema_decimal_compatibility_test.go, internal/api/router_studio_v2_mysql_counter_proof_fixture_test.go
  - Removed: none
- Affected specs: recording-plan-storage
