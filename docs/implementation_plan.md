# go_gateway 完整強化與 hsllogic 共通包實作計劃 (v2)

## 1. 核心目標

建立一個仿照 **HslCommunication** 核心邏輯的 Go 套件 (`hsllogic`)，並將其整合至 **所有** 現有與新增的適配器 (`mc3e`, `modbus`, `fatek`, `mqtt`, `siemens`, `websocket`)。

---

## 2. hsllogic 套件架構

### 套件路徑

```text
go_gateway/
└── lib/
    └── hsllogic/
        ├── bytetransform.go   # 字節序轉換核心 (ABCD/BADC/CDAB/DCBA)
        ├── address_parser.go  # 工業地址解析標準化
        ├── operate_result.go  # 統一結果封裝 (仿 OperateResult<T>)
        ├── data_types.go      # 數據型別與轉換工具
        └── packet_logger.go   # 通訊報文軌跡日誌
```

---

## 3. 現有協議分析摘要

| 協議           | 檔案                                                                                                 | 現狀            | 升級方向                            |
| :------------- | :--------------------------------------------------------------------------------------------------- | :-------------- | :---------------------------------- |
| **MC3E**       | [adapters/mc3e.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/mc3e.go)     | 基礎可用        | 整合 hsllogic 地址解析與數據轉換    |
| **Modbus**     | [adapters/modbus.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/modbus.go) | 基礎可用        | 整合 hsllogic ByteTransform         |
| **FATEK**      | [adapters/fatek.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go)   | 基礎可用        | 整合 hsllogic 結果封裝              |
| **MQTT**       | [adapters/mqtt.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/mqtt.go)     | 功能完整        | 僅需整合結果封裝 (OperateResult)    |
| **WebSocket**  | [handlers/websocket.go](file:///c:/AIProject/go_gateway/internal/api/handlers/websocket.go)          | 僅用於 API 廣播 | 規劃升級為設備連接器 (可選)         |
| **Siemens S7** | **尚未實作**                                                                                         | -               | **新增**：參考 HSL 文件建立完整驅動 |

---

## 4. 詳細模組設計

### A. `bytetransform.go` (字節序轉換)

```go
// DataFormat 定義字節序格式
type DataFormat string

const (
    DataFormatABCD DataFormat = "ABCD" // Big-Endian
    DataFormatBADC DataFormat = "BADC" // Big-Endian with Byte Swap
    DataFormatCDAB DataFormat = "CDAB" // Little-Endian with Word Swap
    DataFormatDCBA DataFormat = "DCBA" // Little-Endian
)

// ByteTransform 提供字節轉換介面
type ByteTransform interface {
    TransformUint16(buffer []byte, index int) uint16
    TransformInt32(buffer []byte, index int) int32
    TransformFloat32(buffer []byte, index int) float32
    // ...
}
```

### B. `address_parser.go` (統一地址解析)

```go
// ParsedAddress 解析後的地址結構
type ParsedAddress struct {
    DeviceType string // "D", "M", "HR", etc.
    Offset     int
    BitIndex   int    // 用於位元設備 (如 M10.1)
    Raw        string // 原始地址字串
}

// ParseAddress 統一解析各協議地址
func ParseAddress(protocol string, address string) (*ParsedAddress, error)
```

### C. `operate_result.go` (結果封裝)

```go
// OperateResult 統一結果封裝
type OperateResult[T any] struct {
    IsSuccess bool
    ErrorCode int
    Message   string
    Content   T
}

// Success 建立成功結果
func Success[T any](content T) OperateResult[T]

// Fail 建立失敗結果
func Fail[T any](code int, message string) OperateResult[T]
```

---

## 5. 新增：Siemens S7 適配器

### 路徑

`go_gateway/internal/datalink/connector/adapters/siemens_s7.go`

### 實作要點

1. 使用開源 Go 庫：`github.com/robinson/gos7`
2. 地址格式：參考 HSL 文件 (`DB1.DBX0.0`, `M0.0`, `I0.0`, `Q0.0`)
3. 支持 PLC 類型：S7-200, S7-300, S7-400, S7-1200, S7-1500

### 核心接口

```go
type SiemensS7Connector struct {
    client    *gos7.S7Client
    config    schema.ConnectionConfigSiemensS7
    connected bool
}

func (c *SiemensS7Connector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error)
func (c *SiemensS7Connector) Write(ctx context.Context, req connector.WriteRequest) error
```

---

## 6. 升級：WebSocket 設備連接器 (可選)

### 路徑

`go_gateway/internal/datalink/connector/adapters/websocket.go`

### 用途

某些設備（如 IoT 閘道）透過 WebSocket 提供數據訂閱，可統一作為 `connector.Protocol` 管理。

### 與現有 WebSocket Handler 的區別

- **現有**: [handlers/websocket.go](file:///c:/AIProject/go_gateway/internal/api/handlers/websocket.go) 用於前端 UI 推播。
- **新增**: `adapters/websocket.go` 用於從遠端設備接收數據。

---

## 7. 所有可升級功能清單

### 優先級 P0 (必做)

- [x] 建立 `hsllogic/bytetransform.go`
- [x] 建立 `hsllogic/address_parser.go`
- [x] 建立 `hsllogic/operate_result.go`
- [x] 建立 `hsllogic/data_types.go`

### 優先級 P1 (重要)

- [ ] 更新 [mc3e.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/mc3e.go) 整合 hsllogic
- [ ] 更新 [modbus.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/modbus.go) 整合 hsllogic (DataFormat 參數)
- [ ] 更新 [fatek.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go) 整合 hsllogic
- [ ] 更新 [mqtt.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/mqtt.go) 使用 OperateResult 封裝

### 優先級 P2 (加分/新功能)

- [ ] **新增** `siemens_s7.go` Siemens S7 適配器
- [ ] **新增** `adapters/websocket.go` WebSocket 設備連接器
- [ ] 建立 `hsllogic/packet_logger.go` 報文軌跡日誌
- [ ] 在 `connector.ReadRequest` 中增加 `DataFormat` 欄位
- [ ] 在 `schema` 套件中增加 `ProtocolSiemensS7` 常數

---

## 8. 驗證計劃

1. **單元測試**: 為 `hsllogic` 每個模組撰寫覆蓋率 80%+ 的測試。
2. **整合測試**: 確認各適配器在使用新邏輯後，對現有功能無破壞性影響。
3. **模擬測試**: 使用 mockgen 模擬 PLC 響應進行 E2E 驗證。
