# Go 工業數據收集器 (go-gateway) 詳細工作分解結構 (WBS)

**專案代碼**: GO-GATEWAY
**目標**: 構建跨平台、多協議、支援自訂數據映射的工業數據採集閘道器。

---

## 1.0 階段一：基礎架構與配置 (Infrastructure & Core)
此階段目標為建立專案骨架，確保開發環境與基礎建設一致。

*   **1.1 專案初始化**
    *   1.1.1 建立 Git Repository 與 `.gitignore` 設定。
    *   1.1.2 執行 `go mod init` 初始化模組。
    *   1.1.3 建立標準目錄結構 (`/cmd`, `/internal`, `/pkg`, `/configs`, `/scripts`)。
    *   1.1.4 設定 `Makefile` (包含 build, test, clean, run 指令)。
*   **1.2 配置管理系統**
    *   1.2.1 引入 `spf13/viper` 套件。
    *   1.2.2 定義 `configs/gateway.yaml` 結構與預設值。
    *   1.2.3 實作配置載入邏輯 (支援環境變數覆寫)。
    *   1.2.4 實作配置熱重載 (Hot Reload) 機制 (選用)。
*   **1.3 日誌系統**
    *   1.3.1 引入 `go.uber.org/zap` 套件。
    *   1.3.2 封裝 Logger 工具 (支援 Console/File 輸出、Log Rotation)。
    *   1.3.3 設定日誌等級動態調整功能。
*   **1.4 系統資料庫初始化 (SQLite)**
    *   1.4.1 引入 `modernc.org/sqlite`。
    *   1.4.2 設計並建立系統 Schema (`devices`, `tasks`, `settings`, `logs` 表)。
    *   1.4.3 實作系統資料庫的 Migration 機制。

## 2.0 階段二：協議抽象層與 Modbus (Protocol Layer)
此階段確立核心介面，並完成最通用的 Modbus 協議支援。

*   **2.1 協議介面定義**
    *   2.1.1 定義 `internal/protocol/protocol.go` 中的 `Protocol` 介面 (Connect, Read, Write)。
    *   2.1.2 定義通用數據結構 (`ReadRequest`, `ReadResult`, `DataType`, `ValueType`)。
    *   2.1.3 定義 `ProtocolFactory` 介面。
*   **2.2 連線管理 (Connection Manager)**
    *   2.2.1 實作 `ConnectionManager` (Singleton Pattern)。
    *   2.2.2 實作 TCP 連線池 (Connection Pool)。
    *   2.2.3 實作串列埠互斥鎖 (Serial Port Mutex) 避免多工衝突。
    *   2.2.4 實作自動重連 (Reconnect) 與心跳 (Heartbeat) 機制。
*   **2.3 Modbus 協議實作**
    *   2.3.1 整合 `github.com/grid-x/modbus`。
    *   2.3.2 實作 Modbus TCP Driver (`Read`, `Write` 封裝)。
    *   2.3.3 實作 Modbus RTU Driver (整合 `go.bug.st/serial`)。
    *   2.3.4 實作 Modbus UDP Driver。
    *   2.3.5 撰寫 Modbus 模擬器單元測試。

## 3.0 階段三：任務執行引擎 (Task Engine)
此階段負責調度採集任務，確保數據源源不斷地被採集。

*   **3.1 任務調度器 (Scheduler)**
    *   3.1.1 定義 `TaskConfig` 與 `PointConfig` 結構。
    *   3.1.2 實作基於 `time.Ticker` 或 Cron 的排程邏輯。
    *   3.1.3 實作任務優先級 (Priority) 排序邏輯。
*   **3.2 執行引擎 (Executor)**
    *   3.2.1 建立 Worker Pool (並發控制)。
    *   3.2.2 實作 `ExecuteTask` 流程 (讀取配置 -> 呼叫 Protocol -> 寫入 Cache)。
    *   3.2.3 實作任務超時 (Timeout) 與錯誤重試 (Retry) 機制。
*   **3.3 記憶體快取 (Memory Store)**
    *   3.3.1 實作環形緩衝區 (Ring Buffer) 或 Map 結構儲存即時值。
    *   3.3.2 實作 `sync.RWMutex` 確保讀寫安全。
    *   3.3.3 提供快取查詢介面 (`Get`, `Set`, `GetSnapshot`)。

## 4.0 階段四：多資料庫支援 (Multi-Database Support)
此階段建立數據落地的能力，支援多種資料庫後端。

*   **4.1 資料庫驅動抽象**
    *   4.1.1 定義 `internal/database/driver.go` 介面 (`Connect`, `BatchInsert`, `Upsert`)。
    *   4.1.2 實作連線測試 (`Ping`) 與 Schema 查詢 (`GetTables`, `GetColumns`) 方法。
*   **4.2 資料庫實作**
    *   4.2.1 **SQLite**: 實作本地數據存儲驅動。
    *   4.2.2 **MySQL**: 整合 `go-sql-driver/mysql`，實作 Driver 介面。
    *   4.2.3 **PostgreSQL**: 整合 `lib/pq`，實作 Driver 介面 (注意 `$1, $2` 參數語法差異)。
    *   4.2.4 **SQL Server**: 整合 `go-mssqldb`，實作 Driver 介面。
*   **4.3 數據清理 (Retention Policy)**
    *   4.3.1 實作定期刪除過期數據的 Background Job。

## 5.0 階段五：數據映射引擎 (Datalink Engine) - **核心亮點**
此階段實作「採集點」到「資料庫欄位」的 ETL 邏輯。

*   **5.1 映射配置模組**
    *   5.1.1 定義 `MappingConfig` (Source, Target, Strategy)。
    *   5.1.2 實作配置驗證邏輯 (檢查型別相容性)。
*   **5.2 映射轉換器 (Mapper)**
    *   5.2.1 實作欄位對應邏輯 (Source Field -> Target Column)。
    *   5.2.2 實作數據型別轉換 (Type Casting: Int -> Float, Bool -> Int 等)。
    *   5.2.3 實作靜態值注入 (Static Value Injection)。
    *   5.2.4 (選用) 簡單表達式計算 (Transform Expression)。
*   **5.3 批次寫入器 (Batch Writer)**
    *   5.3.1 實作緩衝區機制 (達到 `BatchSize` 或 `FlushInterval` 時觸發寫入)。
    *   5.3.2 實作 `Strategy` 邏輯 (區分 Insert, Upsert, Update SQL 生成)。
    *   5.3.3 錯誤處理：寫入失敗時的暫存或 Log 記錄。

## 6.0 階段六：進階協議擴充 (Advanced Protocols)
實作非標準或廠商特定的通訊協議。

*   **6.1 FATEK FBs 協議**
    *   6.1.1 實作指令構建器 (`BuildReadCommand`, `BuildWriteCommand`)。
    *   6.1.2 實作 LRC 校驗演算法。
    *   6.1.3 實作回應解析器 (`ParseResponse`) 處理 ASCII 格式。
    *   6.1.4 實作 FATEK over TCP 與 FATEK Serial。
*   **6.2 Mitsubishi MC Protocol (3E Frame)**
    *   6.2.1 實作 3E Binary Frame 編碼器。
    *   6.2.2 實作 3E Binary Frame 解碼器。
    *   6.2.3 建立設備代碼映射表 (D, M, X, Y -> Hex Code)。
    *   6.2.4 實作 MC Protocol TCP Client。

## 7.0 階段七：MQTT 雲端整合 (Cloud Integration)
*   **7.1 MQTT Client**
    *   7.1.1 整合 `paho.mqtt.golang`。
    *   7.1.2 實作連線管理 (含 TLS 支援)。
*   **7.2 數據發布**
    *   7.2.1 定義 Payload JSON 格式。
    *   7.2.2 實作數據變更推送 (Report on Change) 或週期推送。
    *   7.2.3 斷線緩存機制 (QoS 1/2 支援)。

## 8.0 階段八：API 服務層 (Backend API)
提供前端與外部系統操作的 RESTful 介面。

*   **8.1 API 框架設定**
    *   8.1.1 初始化 Gin Engine 與 Middleware (CORS, Logger, Recovery)。
    *   8.1.2 規劃 API Versioning (`/api/v1`).
*   **8.2 資源管理 API**
    *   8.2.1 **Devices API**: CRUD, Test Connection.
    *   8.2.2 **Tasks API**: CRUD, Start/Stop, Status.
    *   8.2.3 **Datalinks API**: CRUD, Preview SQL.
    *   8.2.4 **Databases API**: CRUD, Fetch Tables/Columns (Schema Browsing).
*   **8.3 即時數據 API**
    *   8.3.1 實作 WebSocket Endpoint (`/api/v1/ws`) 推送即時數據。
    *   8.3.2 實作歷史數據查詢 API (查詢 SQLite/Logger)。

## 9.0 階段九：前端監控介面 (Frontend UI)
*   **9.1 前端專案建置**
    *   9.1.1 初始化 React + TypeScript + Vite 專案。
    *   9.1.2 設定 UI Component Library (e.g., MUI, Ant Design)。
    *   9.1.3 設定 API Client (Axios/TanStack Query)。
*   **9.2 頁面開發**
    *   9.2.1 **Dashboard**: 儀表板與狀態概覽。
    *   9.2.2 **Device/Task Manager**: 表單與列表操作。
    *   9.2.3 **Datalink Configurator**: 拖拉式或選單式映射設定介面 (UX 重點)。
    *   9.2.4 **Data Viewer**: 即時數據表格與簡單趨勢圖。
*   **9.3 整合**
    *   9.3.1 建置前端靜態檔案 (`npm run build`)。
    *   9.3.2 使用 Go `embed` 將靜態檔打包進 Binary。
    *   9.3.3 設定 Gin 路由服務靜態檔案。

## 10.0 階段十：部署與測試 (Deploy & Test)
*   **10.1 測試**
    *   10.1.1 執行全系統整合測試。
    *   10.1.2 長時間穩定性測試 (Soak Testing)。
*   **10.2 跨平台編譯**
    *   10.2.1 Windows Build (`GOOS=windows`).
    *   10.2.2 Linux Build (`GOOS=linux`).
    *   10.2.3 ARM Build (Raspberry Pi).
*   **10.3 系統整合**
    *   10.3.1 Windows Service 註冊功能實作。
    *   10.3.2 Linux Systemd Unit File 撰寫。
    *   10.3.3 Dockerfile 撰寫 (容器化部署)。

---

### 建議優先執行項目 (Critical Path)

1.  **Phase 2.1 & 4.1 (Interfaces)**: 先將 `Protocol` 和 `Database` 的介面定義清楚，這是整個系統解耦的關鍵。
2.  **Phase 3 (Task Engine)**: 這是系統的心臟，確保排程和並發執行是穩定的。
3.  **Phase 5 (Datalink)**: 這是此產品與普通 Modbus Poll 工具最大的差異化功能，應儘早驗證 "Config -> Map -> SQL Write" 的流程原型。
