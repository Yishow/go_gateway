## ADDED Requirements

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

### Requirement: PostgreSQL 初始月份分區可建立

PostgreSQL 初始時序儲存 migration SHALL 使用 DATE 型別參數建立當月及次月分區，不能因 timestamp 與 DATE 的函式簽名差異而失敗。

#### Scenario: 新資料庫建立月份分區
- **WHEN** 在空白 schema 執行 PostgreSQL 初始 migration
- **THEN** migration 完成，當月分區存在，未來月份分區可建立且資料可寫入讀回

##### Example: 讀回時序樣本
- **GIVEN** 新 schema 已完成 migration
- **WHEN** 將有效 UUID tag 的數值 123.45 寫入當月 timeseries
- **THEN** 讀回同一 tag 得到 123.45
