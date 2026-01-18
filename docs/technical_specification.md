# go-gateway 技術規格與維運手冊 (Draft)

## 1. 系統概述 (System Overview)

`go-gateway` 是一個專為工業物聯網 (IIoT) 設計的高效能數據收集器 (Data Collector)。它作為 OT (Operation Technology) 與 IT (Information Technology) 之間的橋樑，能夠從各種工業設備 (PLC) 採集數據，經過清洗與轉換後，即時寫入資料庫或發送至雲端 MQTT Broker。

### 1.1 核心價值

- **跨平台支援**：單一執行檔 (Single Binary) 即可運作於 Windows Server、Linux 伺服器及 ARM 嵌入式裝置 (如 Raspberry Pi)。
- **多協議整合**：原生支援 Modbus (TCP/RTU/UDP)、FATEK FBs、Mitsubishi MC Protocol，解決設備孤島問題。
- **數據映射引擎 (Datalink)**：內建強大的 ETL 邏輯，可直接將設備暫存器 (Registers) 映射至資料庫欄位，無需額外開發轉接程式。
- **高可靠性**：具備自動重連、心跳檢測、資料緩存 (Buffering) 與斷線補傳機制。

### 1.2 支援協議與資料庫

| 類別         | 支援項目                                                                                         |
| ------------ | ------------------------------------------------------------------------------------------------ |
| **採集協議** | Modbus TCP, Modbus RTU/ASCII (RS-485), FATEK FBs (Serial/Net), Mitsubishi MC Protocol (3E Frame) |
| **上傳協議** | MQTT v3.1.1/v5 (支援 TLS SSL), RESTful API, WebSocket                                            |
| **資料庫**   | SQLite (本地儲存), MySQL, PostgreSQL, SQL Server (MSSQL)                                         |

---

## 2. 部署與安裝 (Deployment & Installation)

本系統採用 Go 語言編譯，部署時僅需一個執行檔 (`gateway.exe` 或 `gateway`) 與配置檔目錄。

### 2.1 系統需求

- **Windows**: Windows 10/11, Windows Server 2016+ (x64)
- **Linux**: Ubuntu 20.04+, Debian 10+, CentOS 7+ (x64/ARM64)
- **硬體**: 至少 1 CPU Core, 256MB RAM (視採集點數而定)

### 2.2 安裝步驟 (Windows)

1.  **下載程式**：將 `gateway.exe` 放置於 `C:\go-gateway\`。
2.  **建立配置檔**：在同目錄下建立 `configs\` 資料夾，並放入 `gateway.yaml` 等配置檔。
3.  **測試執行**：
    ```powershell
    cd C:\go-gateway
    .\gateway.exe -c configs\gateway.yaml
    ```
4.  **註冊為服務** (開機自動啟動)：
    ```powershell
    # 使用系統管理員權限執行 PowerShell
    New-Service -Name "GoGateway" -BinaryPathName "C:\go-gateway\gateway.exe -c C:\go-gateway\configs\gateway.yaml" -Description "Industrial Data Collector Service" -StartupType Automatic
    Start-Service "GoGateway"
    ```

### 2.3 安裝步驟 (Linux/Raspberry Pi)

1.  **部署檔案**：將 `gateway` 執行檔複製到 `/opt/go-gateway/` 並賦予執行權限。
    ```bash
    mkdir -p /opt/go-gateway/configs
    chmod +x /opt/go-gateway/gateway
    ```
2.  **設定 Systemd 服務**：建立 `/etc/systemd/system/go-gateway.service`

    ```ini
    [Unit]
    Description=Go Gateway Service
    After=network.target

    [Service]
    ExecStart=/opt/go-gateway/gateway -c /opt/go-gateway/configs/gateway.yaml
    WorkingDirectory=/opt/go-gateway
    Restart=always
    User=root

    [Install]
    WantedBy=multi-user.target
    ```

3.  **啟動服務**：
    ```bash
    systemctl enable go-gateway
    systemctl start go-gateway
    systemctl status go-gateway
    ```

---

## 3. 配置指南 (Configuration Guide)

系統配置分為三個主要檔案，通常位於 `configs/` 目錄下：

| 檔案            | 用途                                                                        |
| --------------- | --------------------------------------------------------------------------- |
| `gateway.yaml`  | **系統主配置**：定義 HTTP 服務埠口、系統資料庫、目標資料庫連線、MQTT 設定。 |
| `devices.yaml`  | **設備定義**：定義 PLC 連線參數 (IP, Port, Station ID)。                    |
| `datalink.yaml` | **數據映射**：定義如何將採集到的 Tag 寫入資料庫表。                         |

### 3.1 系統主配置 (`gateway.yaml`)

```yaml
server:
  port: 8080 # Web UI 與 API 埠口

system_database:
  driver: sqlite
  path: "./data/gateway.db" # 存放任務狀態與日誌

target_databases: # 定義數據寫入的目標資料庫
  - id: "db-production"
    driver: mysql # 支援: mysql, postgres, sqlserver
    host: "192.168.1.200"
    port: 3306
    database: "factory_data"
    username: "gateway_user"
    password: "${DB_PASSWORD}" # 支援環境變數

mqtt:
  broker: "tcp://mqtt.example.com:1883"
  client_id: "gateway-001"
  enabled: true
```

### 3.2 設備定義 (`devices.yaml`)

```yaml
devices:
  - id: "plc-assembly-01"
    name: "組裝線主控 PLC"
    protocol: modbus_tcp # 支援: modbus_tcp, modbus_rtu, fatek_serial
    connection:
      host: "192.168.1.10"
      port: 502
      slave_id: 1
      timeout: 2s

  - id: "meter-power-01"
    name: "電力計 (RS-485)"
    protocol: modbus_rtu
    connection:
      serial_port: "COM3" # Windows 格式。Linux 使用 "/dev/ttyUSB0"
      baud_rate: 9600
      data_bits: 8
      parity: "N"
      stop_bits: 1
```

### 3.3 數據映射 (`datalink.yaml`)

此為本系統核心功能，定義如何將 PLC 的暫存器 (Source) 轉換並寫入資料庫 (Target)。

```yaml
mappings:
  - id: "map-temp-monitoring"
    name: "溫度監控數據"
    source:
      device_id: "plc-assembly-01"
      tags: ["D100", "D101"] # PLC 暫存器位置
    target:
      database_id: "db-production"
      table: "sensor_logs"
      columns:
        - source_field: "D100"
          target_column: "temperature_c"
          data_type: float
        - source_field: "D101"
          target_column: "humidity_percent"
          data_type: float
        - static_value: "LINE-01"
          target_column: "line_id"
    strategy: "insert" # insert (記錄歷程) 或 upsert (即時狀態)
    flush_interval: 5s # 每 5 秒批次寫入一次
```

---

## 4. 系統架構 (System Architecture)

### 4.1 數據流向 (Data Flow)

系統運作遵循「採集 (Collect) -> 轉換 (Transform) -> 派送 (Dispatch)」的流程：

1.  **Protocol Adapter (協議適配器)**：根據 `devices.yaml` 設定，透過 Modbus/FATEK 等驅動程式主動詢問 (Polling) 設備數據。
2.  **Memory Store (記憶體快取)**：採集到的原始數據 (Raw Bytes) 會先暫存於記憶體中的環形緩衝區，供快速查詢。
3.  **Datalink Engine (映射引擎)**：
    - 監聽快取更新。
    - 根據 `datalink.yaml` 規則，將 Raw Data 轉換為有意義的數值 (如 Int16 -> Float)。
    - 組裝 SQL 語句 (Insert/Upsert)。
4.  **Database Writer (資料庫寫入器)**：利用批次處理 (Batch Processing) 機制，將 SQL 寫入目標資料庫 (`target_databases`)。
5.  **MQTT Publisher**：若有設定，同步將變更數據發布至 MQTT Broker。

### 4.2 模組互動

- **Task Scheduler**: 負責管理所有採集任務的生命週期，確保依設定頻率執行。
- **Connection Manager**: 維護設備連線池 (Connection Pool)，處理斷線重連與資源釋放。

---

## 5. API 參考 (API Reference)

閘道器內建 RESTful API，預設埠口為 `8080`。

### 5.1 設備管理

- **GET /api/v1/devices**
  - 功能：取得所有設備列表與連線狀態。
  - 回應：

  ```json
  [
    {
      "id": "plc-01",
      "name": "Main PLC",
      "status": "connected",
      "last_seen": "2024-01-01T12:00:00Z"
    }
  ]
  ```

- **GET /api/v1/devices/{id}/test**
  - 功能：測試指定設備連線。

### 5.2 即時數據

- **GET /api/v1/data/live?device={id}**
  - 功能：取得指定設備的最新快取數據。
- **WS /api/v1/ws**
  - 功能：WebSocket 連線，訂閱即時數據變更。

---

## 6. 維運與故障排除 (Maintenance & Troubleshooting)

### 6.1 日誌系統 (Logging)

- **路徑**：預設位於 `./logs/system.log`。
- **格式**：JSON 格式，方便 ELK Stack 或其他日誌工具收集。
- **等級**：可於配置中調整 (DEBUG, INFO, WARN, ERROR)。
  - `INFO`: 一般運作記錄 (啟動、停止、連線成功)。
  - `WARN`: 暫時性錯誤 (如單次採集超時、資料庫忙碌)。
  - `ERROR`: 持續性故障 (如已達最大重試次數仍無法連線、配置檔錯誤)。

### 6.2 常見問題

**Q1: 設備狀態顯示 Disconnected**

- 檢查實體線路與網路 Ping 值。
- 檢查 `devices.yaml` 中的 IP、Port、Slave ID 是否正確。
- 查看日誌中是否有 `Connection refused` 或 `Timeout` 錯誤。

**Q2: 資料庫無數據寫入**

- 檢查 `datalink.yaml` 中的 `target.database_id` 是否對應 `gateway.yaml` 中的定義。
- 檢查資料庫連線帳密權限 (是否允許 INSERT/UPDATE)。
- 確認 `flush_interval` 設定，數據可能仍在緩衝區待寫入。

**Q3: 記憶體使用量過高**

- 檢查 MQTT 是否斷線導致訊息堆積。
- 確認是否開啟了 `DEBUG` 模式日誌。

### 6.3 資料庫維護

- **SQLite**: 建議每月執行一次 `VACUUM` 指令以釋放空間。
- **Retention**: 可設定排程定期清理 `system_database` 中的歷史日誌與舊任務記錄。
