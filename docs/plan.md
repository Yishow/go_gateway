# Go 工業數據收集器 (go-gateway) 實作計劃

## 專案概述

使用 Go 語言開發一個工業數據收集器，支援 **Windows/Linux/嵌入式** 跨平台部署，提供多種工控協議支援，並具備 **可自訂資料庫表/欄位映射** 功能。

### 支援平台
- **Windows** (主要開發環境)
- Linux (伺服器/嵌入式)
- ARM (樹莓派、工業閘道器)

### 支援協議
- Modbus RTU (RS-485 串列)
- Modbus TCP/UDP (網路)
- FATEK FBs (RS-485 + Ethernet)
- MC Protocol 3E (三菱 PLC)
- MQTT (雲端上傳)

### 核心功能
- **數據映射**：採集點 → 資料庫表/欄位可自訂對應
- **多資料庫支援**：SQLite (本地) / MySQL / PostgreSQL / SQL Server
- **即時快取**：環形緩衝區 + 記憶體暫存器

---

## 專案結構

```
go-gateway/
├── cmd/gateway/main.go           # 主程式入口
├── internal/
│   ├── app/app.go               # 應用程式生命週期
│   ├── config/                  # 配置管理
│   ├── protocol/                # 協議抽象層
│   │   ├── protocol.go          # 介面定義
│   │   ├── modbus/              # Modbus 實作
│   │   ├── fatek/               # FATEK 實作 (自行開發)
│   │   └── mcprotocol/          # MC Protocol 實作 (自行開發)
│   ├── connection/              # 連線管理
│   │   ├── manager.go           # 連線管理器 (Singleton)
│   │   └── pool.go              # 連線池
│   ├── task/                    # 任務調度
│   │   ├── executor.go          # 執行引擎
│   │   └── scheduler.go         # 排程器
│   ├── store/                   # 數據存儲
│   │   ├── memory.go            # 記憶體快取 (環形緩衝區)
│   │   └── database.go          # 資料庫抽象介面
│   ├── datalink/                # 數據映射模組 (新增)
│   │   ├── mapper.go            # 映射引擎
│   │   ├── config.go            # 映射配置定義
│   │   └── writer.go            # 資料庫寫入器
│   ├── database/                # 多資料庫支援 (新增)
│   │   ├── driver.go            # 資料庫驅動介面
│   │   ├── sqlite.go            # SQLite 實作
│   │   ├── mysql.go             # MySQL 實作
│   │   ├── postgres.go          # PostgreSQL 實作
│   │   └── sqlserver.go         # SQL Server 實作
│   ├── mqtt/                    # MQTT 客戶端
│   ├── api/                     # REST API
│   └── web/embed.go             # 前端靜態資源嵌入
├── web/                         # React 前端原始碼
├── configs/                     # YAML 配置檔
├── Makefile
└── go.mod
```

---

## 核心依賴

| 功能 | 套件 | 說明 |
|------|------|------|
| Modbus | `github.com/grid-x/modbus` | TCP/RTU 支援 |
| 串列埠 | `go.bug.st/serial` | **跨平台** (Windows COM / Linux tty) |
| MQTT | `github.com/eclipse/paho.mqtt.golang` | Eclipse 官方 |
| SQLite | `modernc.org/sqlite` | 純 Go，無 CGO |
| MySQL | `github.com/go-sql-driver/mysql` | MySQL 驅動 |
| PostgreSQL | `github.com/lib/pq` | PostgreSQL 驅動 |
| SQL Server | `github.com/denisenkom/go-mssqldb` | MSSQL 驅動 |
| HTTP | `github.com/gin-gonic/gin` | 高效能框架 |
| 配置 | `github.com/spf13/viper` | YAML 支援 |
| 日誌 | `go.uber.org/zap` | 高效能日誌 |

### Windows 特別注意
- 串列埠路徑：Windows 使用 `COM1`, `COM3` 等格式
- 服務註冊：可使用 `golang.org/x/sys/windows/svc` 註冊為 Windows 服務

---

## 實作階段

### 階段 1：專案骨架
- [ ] 初始化 Go module
- [ ] 建立目錄結構
- [ ] 配置管理 (Viper + YAML)
- [ ] 日誌系統 (Zap)

### 階段 2：協議抽象層
- [ ] 定義 `Protocol` 介面 (`internal/protocol/protocol.go`)
- [ ] 實作 Modbus TCP (`internal/protocol/modbus/tcp.go`)
- [ ] 實作 Modbus RTU (`internal/protocol/modbus/rtu.go`)
- [ ] 實作 Modbus UDP (擴展 grid-x/modbus)

### 階段 3：連線管理
- [ ] 連線管理器 Singleton (`internal/connection/manager.go`)
- [ ] TCP 連線池
- [ ] 串列埠管理
- [ ] 心跳檢測與自動重連

### 階段 4：任務執行引擎
- [ ] 任務定義與排程 (`internal/task/executor.go`)
- [ ] Worker Pool 並發控制
- [ ] 週期性採集調度

### 階段 5：數據存儲與多資料庫支援
- [ ] 記憶體快取 - 環形緩衝區 (`internal/store/memory.go`)
- [ ] 資料庫驅動抽象介面 (`internal/database/driver.go`)
- [ ] SQLite 實作 (`internal/database/sqlite.go`)
- [ ] MySQL 實作 (`internal/database/mysql.go`)
- [ ] PostgreSQL 實作 (`internal/database/postgres.go`)
- [ ] SQL Server 實作 (`internal/database/sqlserver.go`)
- [ ] 數據清理策略

### 階段 5.5：數據映射引擎 (Datalink)
- [ ] 映射配置定義 (`internal/datalink/config.go`)
- [ ] 映射引擎實作 (`internal/datalink/mapper.go`)
- [ ] 資料庫批次寫入器 (`internal/datalink/writer.go`)
- [ ] UI 映射設定介面

### 階段 6：MQTT 整合
- [ ] MQTT 客戶端管理 (`internal/mqtt/client.go`)
- [ ] 數據發布器
- [ ] 自動重連機制

### 階段 7：FATEK 協議 (自行開發)
- [ ] 指令構建器 (`internal/protocol/fatek/command.go`)
- [ ] 回應解析器 (`internal/protocol/fatek/parser.go`)
- [ ] RS-485 實作 (`internal/protocol/fatek/serial.go`)
- [ ] Ethernet 實作 (`internal/protocol/fatek/ethernet.go`)

### 階段 8：MC Protocol 3E (自行開發)
- [ ] 3E Frame 編解碼 (`internal/protocol/mcprotocol/frame3e.go`)
- [ ] 設備代碼映射
- [ ] TCP 客戶端實作

### 階段 9：REST API
- [ ] Gin 路由設定 (`internal/api/router.go`)
- [ ] 設備管理 API
- [ ] 任務管理 API
- [ ] 數據查詢 API
- [ ] WebSocket 即時推播

### 階段 10：前端 UI
- [ ] React + TypeScript 專案初始化
- [ ] 設備配置頁面
- [ ] 任務管理頁面
- [ ] 數據監控儀表板
- [ ] 使用 `embed` 打包進執行檔

---

## 關鍵檔案

1. `internal/protocol/protocol.go` - 協議抽象介面
2. `internal/connection/manager.go` - 連線管理器
3. `internal/task/executor.go` - 任務執行引擎
4. `internal/store/memory.go` - 記憶體快取
5. `internal/datalink/mapper.go` - 數據映射引擎 (新增)
6. `internal/database/driver.go` - 多資料庫驅動介面 (新增)
7. `internal/protocol/fatek/command.go` - FATEK 協議實作

---

## 配置檔範例

```yaml
# configs/gateway.yaml
server:
  port: 8080

# 系統資料庫 (存放配置、任務等)
system_database:
  driver: sqlite
  path: "./data/gateway.db"

# 目標資料庫 (存放採集數據)
target_databases:
  - id: "db-local"
    driver: sqlite
    path: "./data/collected.db"

  - id: "db-remote"
    driver: mysql
    host: "192.168.1.200"
    port: 3306
    database: "factory_data"
    username: "gateway"
    password: "${DB_PASSWORD}"  # 環境變數

  - id: "db-mssql"
    driver: sqlserver
    host: "192.168.1.201"
    port: 1433
    database: "SCADA"
    username: "sa"
    password: "${MSSQL_PASSWORD}"

mqtt:
  broker: "tcp://mqtt.example.com:1883"
  client_id: "gateway-001"

task:
  worker_count: 4
  default_timeout: 10s
```

```yaml
# configs/devices.yaml
devices:
  - id: "plc-001"
    name: "主控 PLC"
    protocol: modbus_tcp
    host: "192.168.1.100"
    port: 502

  - id: "fatek-001"
    name: "永宏 FBs"
    protocol: fatek_serial
    port: "COM3"           # Windows 格式
    # port: "/dev/ttyUSB0" # Linux 格式
    baud_rate: 9600
    station: 1
```

```yaml
# configs/datalink.yaml - 數據映射配置
mappings:
  # 範例 1: 溫度感測器 → MySQL 表
  - id: "mapping-001"
    name: "溫度數據映射"
    source:
      device_id: "plc-001"
      task_id: "task-temp"
      tag: "temperature"
    target:
      database_id: "db-remote"
      table: "sensor_readings"
      columns:
        - source_field: "value"
          target_column: "temperature"
          data_type: float
        - source_field: "timestamp"
          target_column: "recorded_at"
          data_type: datetime
        - static_value: "plc-001"
          target_column: "device_id"
          data_type: string
    strategy: "insert"      # insert | upsert | update
    batch_size: 100
    flush_interval: 5s

  # 範例 2: 多個採集點 → 單一表
  - id: "mapping-002"
    name: "設備狀態彙整"
    source:
      device_id: "fatek-001"
      task_id: "task-status"
      tags: ["M0", "M1", "M2", "M3"]  # 多個暫存器
    target:
      database_id: "db-mssql"
      table: "equipment_status"
      columns:
        - source_field: "M0"
          target_column: "running"
          data_type: bool
        - source_field: "M1"
          target_column: "alarm"
          data_type: bool
        - source_field: "M2"
          target_column: "mode"
          data_type: int
        - source_field: "M3"
          target_column: "error_code"
          data_type: int
    strategy: "upsert"
    upsert_key: "device_id"
```

---

## 驗證方式

### 單元測試
```bash
go test -v -race ./...
```

### 整合測試
1. 啟動 gateway：`./bin/gateway -c configs/gateway.yaml`
2. 訪問 Web UI：`http://localhost:8080/ui`
3. 新增測試設備（Modbus 模擬器）
4. 建立採集任務
5. 驗證數據收集與 MQTT 發布

### 數據映射測試
1. 建立測試 MySQL 資料庫和表
2. 在 UI 配置數據映射（採集點 → 表欄位）
3. 啟動採集任務
4. 驗證數據正確寫入目標資料庫
5. 測試批次寫入效能

### Windows 部署測試
```powershell
# 編譯 Windows 版本
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o gateway.exe ./cmd/gateway

# 執行測試
.\gateway.exe -c configs\gateway.yaml

# 註冊為 Windows 服務 (可選)
sc create GoGateway binPath="C:\gateway\gateway.exe -c C:\gateway\config.yaml"
```

### 嵌入式部署測試
```bash
# 交叉編譯 ARM 版本
make build-arm

# 複製到樹莓派
scp bin/gateway-linux-arm pi@raspberry:/opt/gateway/

# 執行測試
ssh pi@raspberry "/opt/gateway/gateway-linux-arm -c /opt/gateway/config.yaml"
```

---

## 技術選型理由

| 決策 | 理由 |
|------|------|
| Go 語言 | 單一執行檔、低資源、高並發、**跨平台** (Windows/Linux/ARM) |
| modernc.org/sqlite | 純 Go 無 CGO，支援 ARM 交叉編譯 |
| go.bug.st/serial | **跨平台串列埠**，Windows COM 和 Linux tty 統一 API |
| grid-x/modbus | 社群維護良好，錯誤處理完善 |
| 內嵌前端 | 單一執行檔部署，無需額外 web server |
| database/sql | Go 標準庫，統一介面支援多種資料庫 |

---

## 風險與應對

| 風險 | 應對策略 |
|------|----------|
| FATEK 協議文件不完整 | 使用示波器/串口監控工具抓取實際封包 |
| MC Protocol 版本差異 | 優先支援 3E Frame，後續擴展 4E |
| 嵌入式記憶體不足 | 使用 sync.Pool、限制緩衝區、定期 GC |
| Windows 服務管理 | 使用 `kardianos/service` 簡化服務註冊 |
| 資料庫連線池管理 | 使用 `database/sql` 內建連線池 + 重連機制 |

---

## 數據映射架構說明

```
採集任務 (Task)
    │
    ├── 讀取設備數據 (Protocol.Read)
    │
    ├── 存入記憶體快取 (MemoryStore)
    │
    └── 數據映射引擎 (Datalink)
         │
         ├── 映射規則匹配 (Mapper)
         │   └── source → target 欄位轉換
         │
         ├── 數據型別轉換
         │   └── int16 → float, bool → int, etc.
         │
         └── 批次寫入資料庫 (Writer)
             ├── INSERT (時序數據)
             ├── UPSERT (狀態更新)
             └── UPDATE (單一記錄)
```

### UI 映射設定功能
- 選擇來源設備/任務/Tag
- 選擇目標資料庫/表
- 拖拉式欄位對應
- 數據型別自動偵測/手動設定
- 即時預覽寫入 SQL

---

## 詳細介面定義

### 1. 協議抽象介面 (`internal/protocol/protocol.go`)

```go
package protocol

import (
    "context"
    "time"
)

// DataType 數據類型枚舉
type DataType int

const (
    DataTypeCoil           DataType = iota // 線圈 (1 bit)
    DataTypeDiscreteInput                   // 離散輸入 (1 bit)
    DataTypeHoldingRegister                 // 保持暫存器 (16 bit)
    DataTypeInputRegister                   // 輸入暫存器 (16 bit)
)

// ValueType 值類型枚舉 (用於數據解析)
type ValueType int

const (
    ValueTypeBool    ValueType = iota // 布林值
    ValueTypeInt16                    // 有號 16 位元
    ValueTypeUint16                   // 無號 16 位元
    ValueTypeInt32                    // 有號 32 位元 (2 暫存器)
    ValueTypeUint32                   // 無號 32 位元
    ValueTypeFloat32                  // 32 位元浮點數
    ValueTypeFloat64                  // 64 位元浮點數 (4 暫存器)
    ValueTypeString                   // ASCII 字串
)

// ReadRequest 讀取請求
type ReadRequest struct {
    DeviceID    string    // 設備 ID
    SlaveID     uint8     // 從站地址 (Modbus)
    DataType    DataType  // 數據類型
    Address     uint16    // 起始位址
    Quantity    uint16    // 讀取數量
    ValueType   ValueType // 值解析類型
}

// ReadResult 讀取結果
type ReadResult struct {
    DeviceID    string        // 設備 ID
    Address     uint16        // 位址
    RawData     []byte        // 原始資料
    ParsedValue interface{}   // 解析後的值
    Timestamp   time.Time     // 時間戳記
    Latency     time.Duration // 讀取延遲
    Error       error         // 錯誤 (如有)
}

// WriteRequest 寫入請求
type WriteRequest struct {
    DeviceID  string      // 設備 ID
    SlaveID   uint8       // 從站地址
    DataType  DataType    // 數據類型
    Address   uint16      // 起始位址
    Values    interface{} // 寫入值
}

// Protocol 協議抽象介面
type Protocol interface {
    // 基本資訊
    Name() string              // 協議名稱
    Version() string           // 協議版本

    // 連線管理
    Connect(ctx context.Context) error
    Disconnect() error
    IsConnected() bool
    Reconnect(ctx context.Context) error

    // 數據操作
    Read(ctx context.Context, req *ReadRequest) (*ReadResult, error)
    Write(ctx context.Context, req *WriteRequest) error
    BatchRead(ctx context.Context, reqs []*ReadRequest) ([]*ReadResult, error)

    // 健康檢查
    Ping(ctx context.Context) error
    GetLatency() time.Duration
}

// ProtocolFactory 協議工廠介面
type ProtocolFactory interface {
    Create(config *DeviceConfig) (Protocol, error)
    SupportedTypes() []string
}
```

### 2. 設備配置結構 (`internal/config/device.go`)

```go
package config

import "time"

// DeviceConfig 設備配置
type DeviceConfig struct {
    ID          string            `yaml:"id" json:"id"`
    Name        string            `yaml:"name" json:"name"`
    Protocol    string            `yaml:"protocol" json:"protocol"`
    Enabled     bool              `yaml:"enabled" json:"enabled"`
    Connection  ConnectionConfig  `yaml:"connection" json:"connection"`
    Options     map[string]any    `yaml:"options,omitempty" json:"options,omitempty"`
}

// ConnectionConfig 連線配置
type ConnectionConfig struct {
    // 網路連線 (TCP/UDP)
    Host    string `yaml:"host,omitempty" json:"host,omitempty"`
    Port    int    `yaml:"port,omitempty" json:"port,omitempty"`

    // 串列埠連線 (RTU)
    SerialPort string `yaml:"serial_port,omitempty" json:"serial_port,omitempty"`
    BaudRate   int    `yaml:"baud_rate,omitempty" json:"baud_rate,omitempty"`
    DataBits   int    `yaml:"data_bits,omitempty" json:"data_bits,omitempty"`
    Parity     string `yaml:"parity,omitempty" json:"parity,omitempty"`      // N, E, O
    StopBits   int    `yaml:"stop_bits,omitempty" json:"stop_bits,omitempty"`

    // 通用設定
    SlaveID        uint8         `yaml:"slave_id,omitempty" json:"slave_id,omitempty"`
    Timeout        time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
    RetryCount     int           `yaml:"retry_count,omitempty" json:"retry_count,omitempty"`
    RetryDelay     time.Duration `yaml:"retry_delay,omitempty" json:"retry_delay,omitempty"`
    HeartbeatInterval time.Duration `yaml:"heartbeat_interval,omitempty" json:"heartbeat_interval,omitempty"`
}

// 協議類型常數
const (
    ProtocolModbusTCP      = "modbus_tcp"
    ProtocolModbusUDP      = "modbus_udp"
    ProtocolModbusRTU      = "modbus_rtu"
    ProtocolModbusRTUOverTCP = "modbus_rtu_over_tcp"
    ProtocolFATEKSerial    = "fatek_serial"
    ProtocolFATEKEthernet  = "fatek_ethernet"
    ProtocolMCProtocol3E   = "mc_protocol_3e"
)
```

### 3. 任務配置結構 (`internal/config/task.go`)

```go
package config

import "time"

// TaskConfig 採集任務配置
type TaskConfig struct {
    ID          string        `yaml:"id" json:"id"`
    Name        string        `yaml:"name" json:"name"`
    DeviceID    string        `yaml:"device_id" json:"device_id"`
    Enabled     bool          `yaml:"enabled" json:"enabled"`
    Interval    time.Duration `yaml:"interval" json:"interval"`
    Priority    int           `yaml:"priority,omitempty" json:"priority,omitempty"`
    Points      []PointConfig `yaml:"points" json:"points"`
}

// PointConfig 採集點配置
type PointConfig struct {
    Tag         string    `yaml:"tag" json:"tag"`                 // 標籤名稱
    Description string    `yaml:"description,omitempty" json:"description,omitempty"`
    DataType    string    `yaml:"data_type" json:"data_type"`     // coil, holding, input, etc.
    Address     uint16    `yaml:"address" json:"address"`         // 起始位址
    Quantity    uint16    `yaml:"quantity,omitempty" json:"quantity,omitempty"` // 數量，預設 1
    ValueType   string    `yaml:"value_type,omitempty" json:"value_type,omitempty"` // int16, float32, etc.
    Scale       float64   `yaml:"scale,omitempty" json:"scale,omitempty"`       // 縮放係數
    Offset      float64   `yaml:"offset,omitempty" json:"offset,omitempty"`     // 偏移量
    Unit        string    `yaml:"unit,omitempty" json:"unit,omitempty"`         // 單位
    ByteOrder   string    `yaml:"byte_order,omitempty" json:"byte_order,omitempty"` // big, little, mid_big, mid_little
}
```

### 4. 數據映射配置 (`internal/datalink/config.go`)

```go
package datalink

import "time"

// MappingConfig 數據映射配置
type MappingConfig struct {
    ID             string          `yaml:"id" json:"id"`
    Name           string          `yaml:"name" json:"name"`
    Enabled        bool            `yaml:"enabled" json:"enabled"`
    Source         SourceConfig    `yaml:"source" json:"source"`
    Target         TargetConfig    `yaml:"target" json:"target"`
    Strategy       WriteStrategy   `yaml:"strategy" json:"strategy"`
    BatchSize      int             `yaml:"batch_size,omitempty" json:"batch_size,omitempty"`
    FlushInterval  time.Duration   `yaml:"flush_interval,omitempty" json:"flush_interval,omitempty"`
    RetryOnError   bool            `yaml:"retry_on_error,omitempty" json:"retry_on_error,omitempty"`
}

// SourceConfig 來源配置
type SourceConfig struct {
    DeviceID string   `yaml:"device_id" json:"device_id"`
    TaskID   string   `yaml:"task_id,omitempty" json:"task_id,omitempty"`
    Tags     []string `yaml:"tags,omitempty" json:"tags,omitempty"`       // 多個 Tag
    Tag      string   `yaml:"tag,omitempty" json:"tag,omitempty"`         // 單一 Tag
}

// TargetConfig 目標配置
type TargetConfig struct {
    DatabaseID string         `yaml:"database_id" json:"database_id"`
    Table      string         `yaml:"table" json:"table"`
    Columns    []ColumnMapping `yaml:"columns" json:"columns"`
}

// ColumnMapping 欄位映射
type ColumnMapping struct {
    SourceField   string `yaml:"source_field,omitempty" json:"source_field,omitempty"`     // 來源欄位
    StaticValue   any    `yaml:"static_value,omitempty" json:"static_value,omitempty"`     // 靜態值
    TargetColumn  string `yaml:"target_column" json:"target_column"`                        // 目標欄位
    DataType      string `yaml:"data_type" json:"data_type"`                               // 目標數據類型
    Transform     string `yaml:"transform,omitempty" json:"transform,omitempty"`           // 轉換表達式
}

// WriteStrategy 寫入策略
type WriteStrategy string

const (
    StrategyInsert WriteStrategy = "insert"  // 直接插入 (時序數據)
    StrategyUpsert WriteStrategy = "upsert"  // 更新或插入
    StrategyUpdate WriteStrategy = "update"  // 僅更新
)

// UpsertConfig Upsert 配置
type UpsertConfig struct {
    KeyColumns []string `yaml:"key_columns" json:"key_columns"` // 主鍵欄位
}
```

### 5. 資料庫驅動介面 (`internal/database/driver.go`)

```go
package database

import (
    "context"
    "database/sql"
)

// Driver 資料庫驅動介面
type Driver interface {
    // 連線管理
    Connect(ctx context.Context) error
    Close() error
    Ping(ctx context.Context) error
    IsConnected() bool

    // 基本操作
    Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
    Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRow(ctx context.Context, query string, args ...any) *sql.Row

    // 交易
    Begin(ctx context.Context) (*sql.Tx, error)

    // 批次操作
    BatchInsert(ctx context.Context, table string, columns []string, rows [][]any) error
    Upsert(ctx context.Context, table string, columns []string, values []any, keyColumns []string) error

    // 結構查詢
    GetTables(ctx context.Context) ([]string, error)
    GetColumns(ctx context.Context, table string) ([]ColumnInfo, error)

    // 驅動資訊
    DriverName() string
    DSN() string
}

// ColumnInfo 欄位資訊
type ColumnInfo struct {
    Name       string `json:"name"`
    DataType   string `json:"data_type"`
    Nullable   bool   `json:"nullable"`
    PrimaryKey bool   `json:"primary_key"`
    Default    string `json:"default,omitempty"`
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
    ID       string `yaml:"id" json:"id"`
    Driver   string `yaml:"driver" json:"driver"`           // sqlite, mysql, postgres, sqlserver
    Host     string `yaml:"host,omitempty" json:"host,omitempty"`
    Port     int    `yaml:"port,omitempty" json:"port,omitempty"`
    Database string `yaml:"database,omitempty" json:"database,omitempty"`
    Username string `yaml:"username,omitempty" json:"username,omitempty"`
    Password string `yaml:"password,omitempty" json:"password,omitempty"`
    Path     string `yaml:"path,omitempty" json:"path,omitempty"`           // SQLite 檔案路徑
    Options  map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
}
```

---

## 資料庫 Schema 設計

### 系統資料庫 (SQLite - gateway.db)

```sql
-- 設備表
CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    protocol TEXT NOT NULL,
    connection_config TEXT NOT NULL,  -- JSON
    options TEXT,                      -- JSON
    enabled INTEGER DEFAULT 1,
    status TEXT DEFAULT 'disconnected', -- connected, disconnected, error
    last_seen_at INTEGER,
    error_message TEXT,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- 採集任務表
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    device_id TEXT NOT NULL,
    interval_ms INTEGER NOT NULL,
    priority INTEGER DEFAULT 0,
    points_config TEXT NOT NULL,      -- JSON
    enabled INTEGER DEFAULT 1,
    status TEXT DEFAULT 'stopped',    -- running, stopped, paused, error
    last_run_at INTEGER,
    next_run_at INTEGER,
    error_count INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- 目標資料庫配置表
CREATE TABLE target_databases (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    driver TEXT NOT NULL,             -- sqlite, mysql, postgres, sqlserver
    connection_config TEXT NOT NULL,  -- JSON (加密儲存敏感資訊)
    enabled INTEGER DEFAULT 1,
    status TEXT DEFAULT 'disconnected',
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- 數據映射表
CREATE TABLE datalink_mappings (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    source_device_id TEXT NOT NULL,
    source_task_id TEXT,
    source_tags TEXT NOT NULL,        -- JSON array
    target_database_id TEXT NOT NULL,
    target_table TEXT NOT NULL,
    column_mappings TEXT NOT NULL,    -- JSON
    strategy TEXT DEFAULT 'insert',   -- insert, upsert, update
    batch_size INTEGER DEFAULT 100,
    flush_interval_ms INTEGER DEFAULT 5000,
    enabled INTEGER DEFAULT 1,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (source_device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (target_database_id) REFERENCES target_databases(id) ON DELETE CASCADE
);

-- MQTT Broker 配置表
CREATE TABLE mqtt_brokers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    broker_url TEXT NOT NULL,
    client_id TEXT NOT NULL,
    username TEXT,
    password TEXT,                    -- 加密儲存
    qos INTEGER DEFAULT 1,
    keep_alive_seconds INTEGER DEFAULT 60,
    clean_session INTEGER DEFAULT 0,
    enabled INTEGER DEFAULT 1,
    status TEXT DEFAULT 'disconnected',
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- 系統日誌表
CREATE TABLE system_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT NOT NULL,              -- debug, info, warn, error
    module TEXT NOT NULL,             -- protocol, task, datalink, mqtt, etc.
    message TEXT NOT NULL,
    details TEXT,                     -- JSON
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX idx_logs_level_time ON system_logs(level, created_at);
CREATE INDEX idx_logs_module ON system_logs(module, created_at);

-- 設備錯誤日誌表
CREATE TABLE device_error_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL,
    error_type TEXT NOT NULL,         -- connection, read, write, timeout
    error_code TEXT,
    error_message TEXT NOT NULL,
    consecutive_count INTEGER DEFAULT 1,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

CREATE INDEX idx_device_errors ON device_error_logs(device_id, created_at);

-- 系統設定表
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    value_type TEXT DEFAULT 'string', -- string, int, bool, json
    description TEXT,
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- 預設設定
INSERT INTO settings (key, value, value_type, description) VALUES
('task.worker_count', '4', 'int', '任務執行緒數'),
('task.default_timeout_ms', '10000', 'int', '預設任務超時時間'),
('heartbeat.interval_ms', '30000', 'int', '心跳檢測間隔'),
('datalink.default_batch_size', '100', 'int', '預設批次寫入大小'),
('log.retention_days', '30', 'int', '日誌保留天數');
```

---

## REST API 端點清單

### 設備管理 API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/devices` | 取得所有設備列表 |
| POST | `/api/v1/devices` | 新增設備 |
| GET | `/api/v1/devices/:id` | 取得單一設備詳情 |
| PUT | `/api/v1/devices/:id` | 更新設備配置 |
| DELETE | `/api/v1/devices/:id` | 刪除設備 |
| POST | `/api/v1/devices/:id/connect` | 連線設備 |
| POST | `/api/v1/devices/:id/disconnect` | 斷開設備 |
| POST | `/api/v1/devices/:id/test` | 測試設備連線 |
| GET | `/api/v1/devices/:id/status` | 取得設備狀態 |

### 任務管理 API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/tasks` | 取得所有任務列表 |
| POST | `/api/v1/tasks` | 新增採集任務 |
| GET | `/api/v1/tasks/:id` | 取得單一任務詳情 |
| PUT | `/api/v1/tasks/:id` | 更新任務配置 |
| DELETE | `/api/v1/tasks/:id` | 刪除任務 |
| POST | `/api/v1/tasks/:id/start` | 啟動任務 |
| POST | `/api/v1/tasks/:id/stop` | 停止任務 |
| GET | `/api/v1/tasks/:id/status` | 取得任務執行狀態 |

### 數據映射 API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/datalinks` | 取得所有映射配置 |
| POST | `/api/v1/datalinks` | 新增映射配置 |
| GET | `/api/v1/datalinks/:id` | 取得單一映射詳情 |
| PUT | `/api/v1/datalinks/:id` | 更新映射配置 |
| DELETE | `/api/v1/datalinks/:id` | 刪除映射配置 |
| POST | `/api/v1/datalinks/:id/test` | 測試映射 (預覽 SQL) |

### 目標資料庫 API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/databases` | 取得所有目標資料庫 |
| POST | `/api/v1/databases` | 新增目標資料庫 |
| GET | `/api/v1/databases/:id` | 取得單一資料庫詳情 |
| PUT | `/api/v1/databases/:id` | 更新資料庫配置 |
| DELETE | `/api/v1/databases/:id` | 刪除資料庫配置 |
| POST | `/api/v1/databases/:id/test` | 測試資料庫連線 |
| GET | `/api/v1/databases/:id/tables` | 取得資料庫表列表 |
| GET | `/api/v1/databases/:id/tables/:table/columns` | 取得表欄位資訊 |

### 數據查詢 API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/data/memory/:device_id` | 取得記憶體快取數據 |
| GET | `/api/v1/data/memory/:device_id/:tag` | 取得特定 Tag 數據 |
| GET | `/api/v1/data/history/:device_id` | 取得歷史數據 |
| WS | `/api/v1/data/realtime` | WebSocket 即時數據推播 |

### MQTT API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/mqtt/brokers` | 取得所有 MQTT Broker |
| POST | `/api/v1/mqtt/brokers` | 新增 MQTT Broker |
| PUT | `/api/v1/mqtt/brokers/:id` | 更新 Broker 配置 |
| DELETE | `/api/v1/mqtt/brokers/:id` | 刪除 Broker |
| POST | `/api/v1/mqtt/brokers/:id/connect` | 連線 Broker |
| POST | `/api/v1/mqtt/brokers/:id/disconnect` | 斷開 Broker |
| POST | `/api/v1/mqtt/publish` | 手動發布訊息 |

### 系統管理 API

| 方法 | 端點 | 說明 |
|------|------|------|
| GET | `/api/v1/system/status` | 系統狀態總覽 |
| GET | `/api/v1/system/stats` | 統計資訊 |
| GET | `/api/v1/system/logs` | 系統日誌 |
| GET | `/api/v1/system/settings` | 系統設定 |
| PUT | `/api/v1/system/settings/:key` | 更新設定 |
| POST | `/api/v1/system/restart` | 重啟服務 |

---

## FATEK FBs 協議規格

### 通訊格式

```
請求訊框:
┌─────┬──────────┬─────────┬───────────────┬───────┬─────┬─────┐
│ STX │ Station  │ Command │     Data      │  ETX  │ LRC │ LRC │
│ 02h │  2 bytes │ 1 byte  │   Variable    │  03h  │ Hi  │ Lo  │
└─────┴──────────┴─────────┴───────────────┴───────┴─────┴─────┘

回應訊框:
┌─────┬──────────┬─────────┬───────┬───────────────┬─────┬─────┬─────┐
│ STX │ Station  │ Command │ Error │     Data      │ ETX │ LRC │ LRC │
│ 02h │  2 bytes │ 1 byte  │1 byte │   Variable    │ 03h │ Hi  │ Lo  │
└─────┴──────────┴─────────┴───────┴───────────────┴─────┴─────┴─────┘
```

### 指令碼

| 指令碼 | ASCII | 功能 |
|--------|-------|------|
| 40h | @ | 讀取離散暫存器狀態 (X, Y, M, S, T, C) |
| 41h | A | 寫入離散暫存器狀態 |
| 42h | B | 讀取暫存器資料 (R, D, T, C) |
| 43h | C | 寫入暫存器資料 |
| 44h | D | 混合讀取 |
| 45h | E | 寫入離散暫存器 (單點) |
| 46h | F | 讀取連續暫存器 |
| 47h | G | 寫入連續暫存器 |
| 48h | H | 混合讀取 (進階) |

### 暫存器類型

| 類型 | 代碼 | 範圍 | 說明 |
|------|------|------|------|
| X | X | 0-255 | 輸入點 |
| Y | Y | 0-255 | 輸出點 |
| M | M | 0-4095 | 內部繼電器 |
| S | S | 0-999 | 步進繼電器 |
| T | T | 0-255 | 定時器接點 |
| C | C | 0-255 | 計數器接點 |
| R | R | 0-4167 | 資料暫存器 (16-bit) |
| D | D | 0-4999 | 資料暫存器 (16-bit) |
| RT | RT | 0-255 | 定時器當前值 (32-bit) |
| RC | RC | 0-255 | 計數器當前值 (32-bit) |

### Go 實作範例

```go
// internal/protocol/fatek/command.go

package fatek

import (
    "fmt"
)

const (
    STX = 0x02
    ETX = 0x03
)

// BuildReadCommand 建立讀取指令
func BuildReadCommand(station int, regType string, address int, count int) []byte {
    stationStr := fmt.Sprintf("%02d", station)

    // 根據暫存器類型選擇指令
    var cmd byte
    var addrFormat string

    switch regType {
    case "X", "Y", "M", "S", "T", "C":
        cmd = 0x40 // 讀取離散暫存器
        addrFormat = fmt.Sprintf("%s%04d", regType, address)
    case "R", "D":
        cmd = 0x46 // 讀取連續暫存器
        addrFormat = fmt.Sprintf("%s%05d", regType, address)
    case "RT", "RC":
        cmd = 0x46
        addrFormat = fmt.Sprintf("%s%03d", regType, address)
    }

    countStr := fmt.Sprintf("%02d", count)
    data := stationStr + string(cmd) + addrFormat + countStr

    lrcHi, lrcLo := calculateLRC(data)

    frame := make([]byte, 0, len(data)+4)
    frame = append(frame, STX)
    frame = append(frame, []byte(data)...)
    frame = append(frame, ETX, lrcHi, lrcLo)

    return frame
}

// calculateLRC 計算 LRC 校驗碼
func calculateLRC(data string) (byte, byte) {
    var sum byte
    for i := 0; i < len(data); i++ {
        sum += data[i]
    }

    // 轉換為 ASCII 十六進位
    hi := "0123456789ABCDEF"[(sum>>4)&0x0F]
    lo := "0123456789ABCDEF"[sum&0x0F]

    return byte(hi), byte(lo)
}

// ParseResponse 解析回應
func ParseResponse(response []byte) ([]byte, error) {
    if len(response) < 7 {
        return nil, fmt.Errorf("response too short: %d bytes", len(response))
    }

    // 檢查 STX
    if response[0] != STX {
        return nil, fmt.Errorf("invalid STX: 0x%02X", response[0])
    }

    // 找到 ETX 位置
    etxPos := -1
    for i := len(response) - 3; i >= 0; i-- {
        if response[i] == ETX {
            etxPos = i
            break
        }
    }

    if etxPos < 0 {
        return nil, fmt.Errorf("ETX not found")
    }

    // 檢查錯誤碼 (位於資料開頭)
    errorCode := response[4]
    if errorCode != '0' {
        return nil, fmt.Errorf("FATEK error code: %c", errorCode)
    }

    // 提取資料 (錯誤碼之後到 ETX 之前)
    data := response[5:etxPos]

    return data, nil
}
```

---

## MC Protocol 3E Frame 規格

### 訊框結構

```
請求訊框 (Binary):
┌───────────┬─────────┬──────┬───────────┬───────────┬─────────────┬─────────┬───────────┬─────────┐
│ SubHeader │ Network │  PC  │ Unit I/O  │  Unit     │ Data Length │   CPU   │  Command  │  Data   │
│  2 bytes  │   No    │  No  │    No     │ Station   │   2 bytes   │ Timer   │ + SubCmd  │Variable │
│  5000h    │ 1 byte  │1byte │  2 bytes  │  1 byte   │             │ 2 bytes │  4 bytes  │         │
└───────────┴─────────┴──────┴───────────┴───────────┴─────────────┴─────────┴───────────┴─────────┘

回應訊框:
┌───────────┬─────────┬──────┬───────────┬───────────┬─────────────┬──────────────┬─────────┐
│ SubHeader │ Network │  PC  │ Unit I/O  │  Unit     │ Data Length │ Complete     │  Data   │
│  D000h    │   No    │  No  │    No     │ Station   │   2 bytes   │ Code 2bytes  │Variable │
└───────────┴─────────┴──────┴───────────┴───────────┴─────────────┴──────────────┴─────────┘
```

### 指令碼

| 指令 | SubCommand | 功能 |
|------|------------|------|
| 0401h | 0000h | 批次讀取 (字元單位) |
| 0401h | 0001h | 批次讀取 (位元單位) |
| 1401h | 0000h | 批次寫入 (字元單位) |
| 1401h | 0001h | 批次寫入 (位元單位) |
| 0403h | 0000h | 隨機讀取 |
| 1402h | 0000h | 隨機寫入 |

### 設備代碼

| 設備 | ASCII | 二進位代碼 | 範圍 |
|------|-------|------------|------|
| D | D* | A8h | 0-65535 |
| M | M* | 90h | 0-65535 |
| X | X* | 9Ch | 0-2047 |
| Y | Y* | 9Dh | 0-2047 |
| B | B* | A0h | 0-32767 |
| W | W* | B4h | 0-32767 |
| R | R* | AFh | 0-65535 |
| ZR | ZR | B0h | 0-65535 |

### Go 實作範例

```go
// internal/protocol/mcprotocol/frame3e.go

package mcprotocol

import (
    "encoding/binary"
    "fmt"
)

// DeviceCode 設備代碼映射
var DeviceCode = map[string]byte{
    "D":  0xA8,
    "M":  0x90,
    "X":  0x9C,
    "Y":  0x9D,
    "B":  0xA0,
    "W":  0xB4,
    "R":  0xAF,
    "ZR": 0xB0,
}

// Frame3E 3E 訊框
type Frame3E struct {
    SubHeader     uint16
    NetworkNo     byte
    PCNo          byte
    UnitIONo      uint16
    UnitStationNo byte
    DataLength    uint16
    CPUTimer      uint16
    Command       uint16
    SubCommand    uint16
    Data          []byte
}

// NewReadRequest 建立讀取請求
func NewReadRequest(device string, address uint32, count uint16) *Frame3E {
    // 建構資料部分
    data := make([]byte, 8)

    // 起始位址 (3 bytes, Little Endian)
    data[0] = byte(address & 0xFF)
    data[1] = byte((address >> 8) & 0xFF)
    data[2] = byte((address >> 16) & 0xFF)

    // 設備代碼
    code, ok := DeviceCode[device]
    if !ok {
        code = 0xA8 // 預設 D
    }
    data[3] = code

    // 讀取點數
    binary.LittleEndian.PutUint16(data[4:], count)

    return &Frame3E{
        SubHeader:     0x5000,
        NetworkNo:     0x00,
        PCNo:          0xFF,
        UnitIONo:      0x03FF,
        UnitStationNo: 0x00,
        CPUTimer:      0x0010, // 4 秒超時
        Command:       0x0401, // 批次讀取
        SubCommand:    0x0000, // 字元單位
        Data:          data,
    }
}

// Encode 編碼訊框
func (f *Frame3E) Encode() []byte {
    // 計算資料長度 (監視定時器 + 指令 + 子指令 + 資料)
    f.DataLength = uint16(2 + 2 + 2 + len(f.Data))

    buf := make([]byte, 0, 21+len(f.Data))

    // SubHeader
    buf = binary.LittleEndian.AppendUint16(buf, f.SubHeader)

    // Network No, PC No
    buf = append(buf, f.NetworkNo, f.PCNo)

    // Unit I/O No
    buf = binary.LittleEndian.AppendUint16(buf, f.UnitIONo)

    // Unit Station No
    buf = append(buf, f.UnitStationNo)

    // Data Length
    buf = binary.LittleEndian.AppendUint16(buf, f.DataLength)

    // CPU Timer
    buf = binary.LittleEndian.AppendUint16(buf, f.CPUTimer)

    // Command, SubCommand
    buf = binary.LittleEndian.AppendUint16(buf, f.Command)
    buf = binary.LittleEndian.AppendUint16(buf, f.SubCommand)

    // Data
    buf = append(buf, f.Data...)

    return buf
}

// ParseResponse 解析回應
func ParseResponse(data []byte) ([]byte, error) {
    if len(data) < 11 {
        return nil, fmt.Errorf("response too short: %d bytes", len(data))
    }

    // 檢查 SubHeader
    subHeader := binary.LittleEndian.Uint16(data[0:2])
    if subHeader != 0xD000 {
        return nil, fmt.Errorf("invalid response sub-header: 0x%04X", subHeader)
    }

    // 檢查完成代碼
    completeCode := binary.LittleEndian.Uint16(data[9:11])
    if completeCode != 0x0000 {
        return nil, fmt.Errorf("MC protocol error: 0x%04X", completeCode)
    }

    // 返回資料部分
    return data[11:], nil
}
```

---

## 前端 UI 頁面設計

### 1. 儀表板 (Dashboard)
- 系統狀態總覽
- 設備連線狀態卡片
- 任務執行統計
- 最近錯誤日誌
- 即時數據圖表

### 2. 設備管理頁面
- 設備列表表格 (篩選、排序、搜尋)
- 新增/編輯設備對話框
- 設備連線狀態即時更新
- 設備測試連線功能
- 批次操作 (啟用/停用/刪除)

### 3. 任務管理頁面
- 任務列表表格
- 新增/編輯任務對話框
- 採集點配置編輯器
- 任務狀態控制 (啟動/停止)
- 任務執行歷史

### 4. 數據映射頁面 (Datalink)
```
┌─────────────────────────────────────────────────────────────────┐
│  數據映射配置                                          [+ 新增]  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐        ┌─────────────────┐                 │
│  │    來  源       │   →    │    目  標       │                 │
│  ├─────────────────┤        ├─────────────────┤                 │
│  │ 設備: [下拉選單] │        │ 資料庫: [下拉選單]│                 │
│  │ 任務: [下拉選單] │        │ 表: [下拉選單]   │                 │
│  │ Tag:  [多選]    │        │                 │                 │
│  └─────────────────┘        └─────────────────┘                 │
│                                                                 │
│  欄位映射:                                                       │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ 來源欄位          目標欄位          資料類型     動作     │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │ temperature   →   temp_value        FLOAT      [刪除]    │   │
│  │ timestamp     →   recorded_at       DATETIME   [刪除]    │   │
│  │ [靜態值: plc-1] → device_id         VARCHAR    [刪除]    │   │
│  │ [+ 新增欄位映射]                                          │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                 │
│  寫入策略: [● INSERT ○ UPSERT ○ UPDATE]                         │
│  批次大小: [100]  刷新間隔: [5] 秒                               │
│                                                                 │
│  預覽 SQL:                                                       │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ INSERT INTO sensor_readings (temp_value, recorded_at,    │   │
│  │   device_id) VALUES (?, ?, ?)                            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                 │
│                              [測試] [取消] [儲存]               │
└─────────────────────────────────────────────────────────────────┘
```

### 5. 目標資料庫管理頁面
- 資料庫連線列表
- 新增/編輯資料庫對話框
- 連線測試功能
- 表結構瀏覽器
- 欄位資訊查看

### 6. 即時監控頁面
- 記憶體快取數據即時顯示
- 按設備/Tag 篩選
- 數據趨勢圖表
- 數據匯出功能

### 7. 系統設定頁面
- 一般設定 (任務執行緒數、超時時間等)
- 日誌設定 (等級、保留天數)
- 備份/還原配置
- 關於/版本資訊

---

## 錯誤處理策略

### 1. 連線錯誤
```go
type ConnectionError struct {
    DeviceID  string
    ErrorType string // timeout, refused, reset, unreachable
    Message   string
    Timestamp time.Time
    Retry     int
}

// 重試策略
type RetryPolicy struct {
    MaxRetries    int           // 最大重試次數
    InitialDelay  time.Duration // 初始延遲
    MaxDelay      time.Duration // 最大延遲
    Multiplier    float64       // 延遲倍增因子
}
```

### 2. 協議錯誤
- Modbus 異常碼處理
- FATEK 錯誤碼對應
- MC Protocol 完成碼處理

### 3. 資料庫錯誤
- 連線斷開自動重連
- 寫入失敗本地緩存
- 批次寫入部分失敗處理

### 4. 日誌記錄
```go
// 日誌等級
const (
    LogDebug = "debug"
    LogInfo  = "info"
    LogWarn  = "warn"
    LogError = "error"
)

// 日誌模組
const (
    ModuleProtocol   = "protocol"
    ModuleConnection = "connection"
    ModuleTask       = "task"
    ModuleDatalink   = "datalink"
    ModuleMQTT       = "mqtt"
    ModuleAPI        = "api"
    ModuleSystem     = "system"
)
```

---

## 效能優化策略

### 1. 記憶體管理
- 使用 `sync.Pool` 重用緩衝區
- 限制環形緩衝區大小
- 定期執行 GC

### 2. 並發控制
- Worker Pool 限制並發任務數
- 設備級別的讀寫鎖
- 資料庫連線池

### 3. 批次處理
- 資料庫批次寫入
- MQTT 訊息批次發布
- 日誌批次寫入

### 4. 快取策略
- 記憶體快取熱點數據
- 配置快取避免頻繁讀取
- 資料庫查詢結果快取
