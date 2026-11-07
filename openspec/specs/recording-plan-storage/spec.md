# recording-plan-storage Specification

## Purpose

TBD - created by archiving change 'add-recording-plans'. Update Purpose after archive.

## Requirements

### Requirement: Versioned recording plans with independent streams
The system SHALL persist RecordingPlan revisions and separate source acquisition, raw retention, summary intervals and destination flush policies.

#### Scenario: One source serves curve and usage
- **WHEN** a meter plan selects raw history, one-minute summaries and hourly usage
- **THEN** collection is shared while outputs keep independent stream identities and policies.

---
### Requirement: Record mode semantics
The system MUST support raw history, window summary, usage interval, state changes, events, batch snapshots and latest-only with truthful capabilities.

#### Scenario: Latest-only has no history
- **WHEN** an operator chooses only latest_only
- **THEN** the UI and query capability state that past curves cannot be reconstructed.

#### Scenario: Flush delay does not reduce readings
- **WHEN** twelve samples are acquired and buffered before one transport flush
- **THEN** every_sample retains twelve samples rather than only the last value.

---
### Requirement: Managed typed storage and identity
The system SHALL store source identity, observation time, definition revision, quality and exact typed value in managed tables with stable record keys.

#### Scenario: Many devices share the same field name
- **WHEN** two meters both expose active_power
- **THEN** their samples remain distinct by identity without creating new per-device tables.

---
### Requirement: Explicit reduction without corrupting statistics
The system MUST calculate summaries and usage from pre-decimation inputs and preserve necessary calculation evidence for the configured correction horizon.

#### Scenario: Short spike under on-change reduction
- **WHEN** a valid sampled spike reaches the calculation input
- **THEN** its sampled maximum remains in the summary even when the raw display stream is reduced.

---
### Requirement: Batch snapshot freshness and trigger identity
The system SHALL bind a batch snapshot to a stable trigger identity and record member-level observation times and completeness.

#### Scenario: Held trigger with one missing member
- **WHEN** the same trigger is polled repeatedly and pressure is missing
- **THEN** one partial batch is recorded with missing pressure, not repeated batches filled with old pressure.

---
### Requirement: Safe compatibility with existing tables
The system MUST preserve existing single-point/grouped mappings and require explicit review before creating managed tables or changing output mode.

#### Scenario: Existing table name collision
- **WHEN** the requested managed table name already belongs to an incompatible table
- **THEN** preparation is blocked and no existing table or data is altered.

---
### Requirement: Managed exact decimal schema compatibility

PostgreSQL managed schema verification SHALL accept the unrestricted NUMERIC representation generated for val_dec, while preserving the existing text representation used by SQLite. Bounded numeric precision/scale, floating-point and integer columns SHALL NOT be accepted as an equivalent exact decimal representation by this compatibility check.

#### Scenario: 驗證剛產生的 PostgreSQL 小數欄位
- **WHEN** managed samples 表由產生的 DDL 建立 val_dec NUMERIC
- **THEN** schema apply 成功並回傳 verified digest，再次 preview 是相容的 no-op

##### Example: 無界小數與文字表示
- **GIVEN** val_dec 的型別為 numeric 或 text
- **WHEN** 執行 managed schema 相容性檢查
- **THEN** 此欄位不會被當成不相容，其他欄位仍照既有規則檢查

#### Scenario: 拒絕可能截斷或近似的小數欄位
- **WHEN** val_dec 是浮點、整數或限定 precision/scale 的型別
- **THEN** preview 拒絕不相容表，且不儲存可執行 token

##### Example: 有界或近似型別
- **GIVEN** val_dec 為 real、double precision、integer、boolean、numeric(38,12) 或 decimal(20,6)
- **WHEN** 執行 preview
- **THEN** 每種型別均回傳不相容錯誤，token 儲存次數為 0


<!-- @trace
source: fix-database-validation-contracts
updated: 2026-11-07
code:
  - internal/datalink/recordingplan/schema_tables.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
tests:
  - internal/datalink/storage/postgres_test.go
  - internal/datalink/recordingplan/schema_decimal_compatibility_test.go
  - internal/api/router_studio_v2_mysql_circuit_plain_fixture_test.go
  - internal/api/router_studio_v2_mysql_counter_proof_fixture_test.go
-->

---
### Requirement: PostgreSQL 初始月份分區可建立

PostgreSQL 初始時序儲存 migration SHALL 使用 DATE 型別參數建立當月及次月分區，不能因 timestamp 與 DATE 的函式簽名差異而失敗。

#### Scenario: 新資料庫建立月份分區
- **WHEN** 在空白 schema 執行 PostgreSQL 初始 migration
- **THEN** migration 完成，當月分區存在，未來月份分區可建立且資料可寫入讀回

##### Example: 讀回時序樣本
- **GIVEN** 新 schema 已完成 migration
- **WHEN** 將有效 UUID tag 的數值 123.45 寫入當月 timeseries
- **THEN** 讀回同一 tag 得到 123.45

<!-- @trace
source: fix-database-validation-contracts
updated: 2026-11-07
code:
  - internal/datalink/recordingplan/schema_tables.go
  - internal/datalink/schema/migrations/001_initial_schema.up.sql
tests:
  - internal/datalink/storage/postgres_test.go
  - internal/datalink/recordingplan/schema_decimal_compatibility_test.go
  - internal/api/router_studio_v2_mysql_circuit_plain_fixture_test.go
  - internal/api/router_studio_v2_mysql_counter_proof_fixture_test.go
-->