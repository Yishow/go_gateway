# Modbus 協議實作

完整的 Modbus TCP、UDP 和 RTU 協議實作，支援所有標準功能碼。

## 功能特性

- ✅ **Modbus TCP** - 完整的 TCP/IP 實作，支援 MBAP 標頭
- ✅ **Modbus UDP** - UDP 傳輸支援
- ✅ **Modbus RTU** - 串列埠通訊，支援 CRC-16 校驗
- ✅ **所有標準功能碼** - 讀寫線圈、暫存器
- ✅ **錯誤處理** - 完整的異常碼和通訊錯誤處理
- ✅ **線程安全** - 使用 mutex 保護並發操作

## 支援的功能碼

### 讀取功能
- `0x01` - Read Coils (讀取線圈狀態)
- `0x02` - Read Discrete Inputs (讀取離散輸入)
- `0x03` - Read Holding Registers (讀取保持暫存器)
- `0x04` - Read Input Registers (讀取輸入暫存器)

### 寫入功能
- `0x05` - Write Single Coil (寫入單個線圈)
- `0x06` - Write Single Register (寫入單個暫存器)
- `0x0F` - Write Multiple Coils (寫入多個線圈)
- `0x10` - Write Multiple Registers (寫入多個暫存器)

## 快速開始

### Modbus TCP

```go
import "go-gateway/internal/protocol/modbus"
import "time"

// 建立客戶端
client := modbus.CreateTCPClient("192.168.1.10", 502, 1, 2*time.Second)

// 連線
if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Close()

// 讀取保持暫存器
registers, err := client.ReadHoldingRegisters(0, 10)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Registers: %v\n", registers)

// 寫入單個暫存器
err = client.WriteSingleRegister(0, 12345)
if err != nil {
    log.Fatal(err)
}
```

### Modbus RTU

```go
// 建立 RTU 客戶端 (串列埠)
client := modbus.CreateRTUClient("COM3", 9600, 8, 1, "N", 2*time.Second, 1)

if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用方式與 TCP 相同
registers, err := client.ReadHoldingRegisters(0, 10)
```

### Modbus UDP

```go
// 建立 UDP 客戶端
client := modbus.CreateUDPClient("192.168.1.10", 502, 1, 2*time.Second)

if err := client.Connect(); err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## API 參考

### 讀取操作

```go
// 讀取線圈 (1-2000 個)
coils, err := client.ReadCoils(address uint16, quantity uint16) ([]bool, error)

// 讀取離散輸入 (1-2000 個)
inputs, err := client.ReadDiscreteInputs(address uint16, quantity uint16) ([]bool, error)

// 讀取保持暫存器 (1-125 個)
registers, err := client.ReadHoldingRegisters(address uint16, quantity uint16) ([]uint16, error)

// 讀取輸入暫存器 (1-125 個)
registers, err := client.ReadInputRegisters(address uint16, quantity uint16) ([]uint16, error)
```

### 寫入操作

```go
// 寫入單個線圈
err := client.WriteSingleCoil(address uint16, value bool) error

// 寫入單個暫存器
err := client.WriteSingleRegister(address uint16, value uint16) error

// 寫入多個線圈 (1-1968 個)
err := client.WriteMultipleCoils(address uint16, values []bool) error

// 寫入多個暫存器 (1-123 個)
err := client.WriteMultipleRegisters(address uint16, values []uint16) error
```

## 錯誤處理

所有操作可能返回以下錯誤類型：

- `ProtocolError` - Modbus 協議層級錯誤（異常碼）
- `CommunicationError` - 通訊層級錯誤（連線、逾時等）

```go
registers, err := client.ReadHoldingRegisters(0, 10)
if err != nil {
    if protocolErr, ok := err.(*modbus.ProtocolError); ok {
        fmt.Printf("Modbus 異常碼: 0x%02X\n", protocolErr.Code())
    } else {
        fmt.Printf("通訊錯誤: %v\n", err)
    }
}
```

## 地址範圍

- **線圈 (Coils)**: 0-65535 (可讀寫)
- **離散輸入 (Discrete Inputs)**: 0-65535 (只讀)
- **輸入暫存器 (Input Registers)**: 0-65535 (只讀)
- **保持暫存器 (Holding Registers)**: 0-65535 (可讀寫)

注意：地址參數使用相對地址（0-65535），而非絕對地址。

## 限制

- 讀取線圈/離散輸入：最多 2000 個
- 讀取暫存器：最多 125 個
- 寫入線圈：最多 1968 個
- 寫入暫存器：最多 123 個

## 技術細節

### Modbus TCP/UDP
- 使用 MBAP (Modbus Application Protocol) 標頭
- Transaction ID 自動管理
- 支援完整的錯誤檢測

### Modbus RTU
- CRC-16 校驗 (多項式 0x8005 反向)
- 自動計算 3.5 字符時間延遲
- 完整的封包解析和驗證

## 範例程式碼

詳細範例請參考 `example_test.go`。
