## Context

使用者在 Studio v2 Setup（Step 1 ~ Step 4）以及系統設定頁（Settings）的操作與資料庫交付流程中，遇到了五個阻礙流暢性與正確性的問題：
1. 系統設定頁切換資料庫類型時未自動帶入標準 Host/Port/Schema/User。
2. Step 4 資料庫設定需重複手動配置，未能直接選取 Settings 中已建置的 Connector Pool 連線。
3. 資料庫密碼在輸入時容易被狀態重新渲染或非預期重設機制清空。
4. Step 1 設備切換任何通訊協議（Modbus TCP/RTU/UDP, FATEK, MC Protocol, MQTT）時，預設值未能按各協議標準正確初始化，且協議設定欄位與後端有缺漏不對齊處。
5. 後端在進行資料庫連線測試或交付寫入時，若目標資料庫尚不存在（例如 MySQL Error 1049），未像 PostgreSQL 一樣具備自動建立資料庫的機制。

## Goals / Non-Goals

**Goals:**
- **資料庫預設值智慧連動**：在 Settings ConnectorRow 與 Step 4 ConnectorSection 切換 kind 時，採用策略 A（智慧連動）：自動更新 Port、Schema、User；保留非本機自訂 Host；SQLite 自動清空連線資訊。
- **Step 4 複用 Connector Pool**：在 Step 4 增加 Connector 選擇器或下拉快速套用機制，讓使用者可直接挑選 Settings Connector Pool 既有連線，無需重複輸入。
- **密碼輸入穩定受控**：確保密碼欄位在輸入期間維持隔離受控狀態，不受非預期 state 覆蓋。
- **所有協議預設值矩陣與前後端對齊**：完善 `getDefaultConfig` 與表單結構，針對六種通訊協議提供專屬標準預設值，並全面對齊前後端欄位命名與型別。
- **後端多資料庫 Auto Create**：在後端 probe/connect 流程中，為 MySQL 等資料庫補齊像 PostgreSQL `ensurePostgresDatabaseIfMissing` 般的自動建立資料庫機制。

**Non-Goals:**
- 變更底層 Datalink 映射核心演算法。
- 引入新的未支援資料庫引擎（維持既有 PostgreSQL、MySQL、SQL Server、SQLite）。

## Decisions

### 決策 1：統一資料庫連線預設值與切換策略（智慧連動）
- **選定方案**：採用策略 A（智慧連動）。
  - **Port**：一律切換為目標 DB 之標準預設 Port（Postgres: 5432, MySQL: 3306, SQL Server: 1433, SQLite: 0）。
  - **Host**：若為空或為本機預設值（`127.0.0.1`、`localhost`、`tsdb.internal`），更新為 `127.0.0.1`；若為自訂遠端 IP 則予以保留。
  - **Schema**：切換為目標 DB 之標準預設 Schema（Postgres: `public`、SQL Server: `dbo`、其餘為空）。
  - **Username**：若為各 DB 預設帳號或空值，切換至新 DB 的預設帳號（Postgres: `postgres`、MySQL: `root`、SQL Server: `sa`）。
  - **SQLite**：清空 Host/Port/Username/Password，Database 設為 `gateway.db`。

### 決策 2：Step 4 支援 Connector Pool 連線選取與快速套用
- **選定方案**：在 Step 4 連接器配置上方提供「從設定集選擇既有連線」下拉或按鈕，當使用者在 Settings 已設定多組 Connector 時，可一鍵載入設定至 Step 4 當前連線。

### 決策 3：隔離密碼輸入受控狀態
- **選定方案**：在 ConnectorRow / ConnectorSection 的 Password 輸入控制上，確保受控 input 綁定與 debounced/local state 機制不會在父層其他欄位或 autosave 觸發時將輸入中的未提交密碼覆蓋或重置為 empty。

### 決策 4：全通訊協議預設值矩陣與前後端對齊
- **選定方案**：重構 `protocols.ts` 中的 `getDefaultConfig` 與 `ConnectionConfigForm.tsx`：
  - `mc_3e`：`host: '192.168.1.100'`, `port: 6000`, `station: 0`, `network_no: 0`, `pc_no: 255`, `module_io_no: 1023`, `module_station_no: 0`
  - `fatek_fbs`：`host: '192.168.1.100'`, `port: 500`, `station: 1`
  - `modbus_tcp`：`host: '192.168.1.100'`, `port: 502`, `slave_id: 1`, `timeout: 5`
  - `modbus_udp`：`host: '192.168.1.100'`, `port: 502`, `slave_id: 1`, `timeout: 5`
  - `modbus_rtu`：`port: '/dev/ttyUSB0'`, `baud: 9600`, `parity: 'N'`, `data_bits: 8`, `stop_bits: 1`, `slave_id: 1`
  - `mqtt`：`broker: 'mqtt://127.0.0.1:1883'`, `client_id: 'gw-01'`, `username: ''`, `password: ''`
  - 前後端欄位解析器與驗證器完整對齊上述欄位名稱與標籤。

### 決策 5：後端 MySQL 資料庫缺失自動建立機制
- **選定方案**：在 `internal/datalink/dbtarget/service_probe.go` 中，實作 `ensureMySQLDatabaseIfMissing`，當捕捉到 MySQL Error 1049 (`ER_BAD_DB_ERROR`) 時，透過不指定 database 的 root/admin 連線執行 `CREATE DATABASE IF NOT EXISTS <dbname>`，建立成功後再重新建立連線。

## Implementation Contract

- **Behavior**:
  - 在 Settings 切換 DB 類型時，Port/Host/Schema/User 自動更新為對應預設值。
  - 在 Step 4 可直接從下拉選單選取 Settings Connector Pool 連線。
  - 密碼欄位輸入時不再發生清空現象。
  - Step 1 新增設備時，切換至任何通訊協議皆即時帶入該協議專屬的標準預設值與欄位。
  - 測試 MySQL 連線時若 database 不存在，後端自動建立並回傳測試成功。
- **Interface / Data Shape**:
  - `getDbKindPatch(nextKind, currentConnector)`: 回傳依智慧策略計算出的 patch 物件。
  - `getDefaultConfig(protocol)`: 回傳涵蓋 6 種通訊協議的標準預設設定物件。
  - `ensureMySQLDatabaseIfMissing(ctx, config, connectErr)`: 後端 Go 函式處理 Error 1049。
- **Failure Modes**:
  - 若 DB 權限不足無法自動建庫，清楚拋出權限不足錯誤（如 `Access denied for user to database`），不進行靜默吞沒。
- **Acceptance Criteria**:
  - 前端單元測試與後端 Go 測試全數通過（`npx vitest run` & `go test ./...`）。
  - 各切換與建庫情境均有相應單元測試覆蓋。
- **Scope Boundaries**:
  - In Scope: Settings 連接器設定、Step 4 連接器選取與切換、密碼欄位受控更新、Step 1 所有通訊協議預設值與前後端對齊、後端 MySQL/PostgreSQL 自動建立資料庫。
  - Out of Scope: 擴充額外非 SQL 儲存介質或大幅改動 ETL 映射排程。

## Risks / Trade-offs

- **[Risk]** MySQL 自動建庫若目前使用者無 `CREATE DATABASE` 權限會失敗。
  - **Mitigation**：捕獲錯誤並回傳可閱讀的清晰錯誤提示，提示使用者手動建立或更換權限帳號。
- **[Risk]** Step 4 選取 Connector Pool 連線時覆蓋未儲存的暫存草稿。
  - **Mitigation**：選取切換時自動同步至 Step 4 connector state，並標記為已選取連線。
