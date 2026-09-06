## Why

在 Studio v2 Setup 與 Database 設定流程中，使用者在操作與連線時面臨以下關鍵問題：
1. **設定預設值缺失**：系統設定頁切換資料庫類型時未自動填入對應的標準 Host 與 Port。
2. **重工問題**：在系統設定頁已建置完成 Connector Pool 後，Step 4 無法直接選取已建立的連接器，導致使用者必須重新手動設定。
3. **密碼輸入體驗不穩定**：資料庫設定在輸入密碼時，容易受 autosave 或父層狀態重新渲染影響而被意外清空。
4. **所有設備通訊協議預設值缺失與欄位不對齊**：Step 1 新增設備時，切換任何通訊協議（Modbus TCP, Modbus RTU, Modbus UDP, FATEK, MC Protocol, MQTT）未依所選協議自動帶入專屬標準預設值（例如 MC Protocol 應預設 Port 6000、站號 0；FATEK 應預設 Port 500；Modbus RTU 應預設 Baud 9600 等），且部分協議通訊參數與標籤未與後端完整對齊。
5. **缺少自動建立 Database 機制**：連線測試若目標資料庫尚不存在（例如 MySQL Error 1049: Unknown database），未提供自動建立資料庫機制導致連線失敗。

本變更旨在一次性修復並完善上述 5 個設定與連線流程體驗。

## What Changes

- **資料庫預設值智慧連動**：
  - PostgreSQL：預設 Host `127.0.0.1`、Port `5432`、Schema `public`、User `postgres`。
  - MySQL：預設 Host `127.0.0.1`、Port `3306`、User `root`。
  - SQL Server：預設 Host `127.0.0.1`、Port `1433`、Schema `dbo`、User `sa`。
  - SQLite：清空 Host/Port/User，預設資料庫路徑為 `gateway.db`。
  - 保留使用者自訂的遠端 Host，但切換類型時自動更新 Port 與 Schema。
- **Step 4 連接器快速複用**：在 Step 4 資料庫設定提供下拉清單或選取器，可直接選取系統設定頁的 Connector Pool 既有連線，無需重複手動填寫。
- **密碼輸入穩定化**：修復密碼輸入欄位之受控狀態與更新機制，避免在使用者輸入過程中被非預期清空。
- **所有通訊協議預設值與前後端對齊**：
  - 切換至 MC Protocol (`mc_3e`) 時自動帶入 Port `6000`、站號 `0`、網路編號 `0` 等標準預設值。
  - 切換至 FATEK (`fatek_fbs`) 時自動帶入 Port `500`、站號 `1`。
  - 切換至 Modbus TCP (`modbus_tcp`) / UDP (`modbus_udp`) 時自動帶入 Port `502`、站號 `1`。
  - 切換至 Modbus RTU (`modbus_rtu`) 時自動帶入 Port `/dev/ttyUSB0`、Baud `9600`、Parity `N`、站號 `1`。
  - 切換至 MQTT (`mqtt`) 時自動帶入 Broker `mqtt://127.0.0.1:1883`、Client ID `gw-01`。
  - 全面檢查並補齊所有協議在 UI 表單與後端 Driver 之間的欄位名稱、型別與標籤對齊。
- **後端自動建立資料庫**：後端在執行資料庫連線測試與初始化時，若遇到資料庫不存在錯誤（例如 MySQL Error 1049），自動執行建立資料庫語法或完成初始化。

## Capabilities

### Modified Capabilities

- `datalink-workbench-v2-settings`: 增強系統設定頁資料庫連接器管理，包含類型切換自動帶入預設 Host/Port/Schema/User、密碼輸入保護。
- `datalink-workbench-v2-step4-database`: Step 4 新增直接選取 Connector Pool 連接器能力，並對齊類型切換預設值連動。
- `datalink-workbench-v2-step1-device`: Step 1 設備設定支援所有通訊協議切換自動帶入專屬標準參數預設值，並完整對齊前後端協議欄位。
- `database-output-delivery`: 後端資料庫連線測試與初始化支援目標資料庫不存在時自動建立資料庫。

## Impact

- Affected specs:
  - `datalink-workbench-v2-settings`
  - `datalink-workbench-v2-step4-database`
  - `datalink-workbench-v2-step1-device`
  - `database-output-delivery`
- Affected code:
  - Modified:
    - `frontend/src/features/datalink/workbench-v2/settings/ConnectorRow.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/step1/DeviceEditor.tsx`
    - `frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts`
    - `frontend/src/features/datalink/workbench-v2/state/protocols.ts`
    - `internal/datalink/connector/service.go`
    - `internal/datalink/dbtarget/postgres.go`
    - `internal/datalink/dbtarget/sqlite.go`
