## 1. 資料庫預設值與密碼輸入穩定化

- [x] [P] 1.1 實作資料庫預設值與智慧連動計算邏輯，落實「決策 1：統一資料庫連線預設值與切換策略（智慧連動）」與 `Database kind switch auto-populates defaults`，切換類型時自動帶入 Port/Schema/User 並保留自訂 Host。驗證目標：`tests/unit/workbench-v2/dbSchemas.test.ts` 新增測試驗證各 DB 類型切換產出之 patch。
- [x] [P] 1.2 在系統設定頁 `ConnectorRow` 中整合智慧連動與密碼隔離機制，落實「決策 3：隔離密碼輸入受控狀態」與 `Protected connector password input`，確保使用者輸入密碼時不被清空。驗證目標：`tests/unit/workbench-v2/settings-connectors.test.tsx` 驗證切換 kind 帶入預設值及密碼輸入連續性。

## 2. Step 4 資料庫連線複用與智慧切換

- [x] [P] 2.1 在 Step 4 資料庫設定提供連接器複用選擇器，落實「決策 2：Step 4 支援 Connector Pool 連線選取與快速套用」與 `Select existing connector from pool in Step 4`，允許直接載入 Settings Connector Pool 連線。驗證目標：`tests/unit/workbench-v2/step4-database.test.tsx` 驗證選取 connector 後即時套用各欄位。
- [x] [P] 2.2 在 Step 4 連接器表單落實 `Step 4 kind switch auto-populates defaults`，切換資料庫種類時同步套用智慧預設值。驗證目標：`tests/unit/workbench-v2/reducer-step4.test.ts` 驗證 kind 切換後的 connector 屬性更新。

## 3. Step 1 全通訊協議預設值與前後端欄位對齊

- [x] [P] 3.1 完善全通訊協議預設值矩陣（涵蓋 Modbus TCP/RTU/UDP, FATEK, MC Protocol, MQTT），落實「決策 4：全通訊協議預設值矩陣與前後端對齊」與 `Protocol switch auto-populates dedicated communication presets`，MC Protocol 預設 Port 6000/站號 0、Fatek 預設 Port 500/站號 1、Modbus RTU 預設 Baud 9600 等。驗證目標：`tests/unit/workbench-v2/protocols.test.ts` 測試所有通訊協議回傳之 `getDefaultConfig`。
- [x] [P] 3.2 檢查並對齊 Step 1 設備設定表單（`ConnectionConfigForm`）與後端各 Driver 欄位映射及標籤，落實 `Protocol form frontend-backend field alignment`。驗證目標：`tests/unit/workbench-v2/step1.test.tsx` 驗證各通訊協議之欄位渲染與資料提交格式。

## 4. 後端自動建立資料庫機制

- [x] 4.1 在後端 probe 模組實作 MySQL 自動建庫，落實「決策 5：後端 MySQL 資料庫缺失自動建立機制」與 `Automatic database creation on probe and initialization`，當捕獲 Error 1049 (Unknown database) 時自動執行 CREATE DATABASE IF NOT EXISTS 並重新驗證連線。驗證目標：`internal/datalink/dbtarget/service_probe_test.go` (或 `service_mysql_test.go`) 單元測試模擬 Error 1049 建庫流程。
- [x] 4.2 整合全端驗證流程，執行前端與後端測試套件確保全數綠燈。驗證目標：執行 `npx vitest run` 與 `go test ./...` 均以 0 error 退出。
