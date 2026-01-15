/*
Package modbus 實作 Modbus 通訊協議，支援 TCP、UDP 和 RTU 三種傳輸方式。

Modbus 是工業自動化領域最廣泛使用的通訊協議之一，支援讀寫線圈（Coils）、
離散輸入（Discrete Inputs）、保持暫存器（Holding Registers）和輸入暫存器（Input Registers）。

## 支援的功能碼

### 讀取功能
- 0x01: Read Coils - 讀取線圈狀態 (1-2000 個)
- 0x02: Read Discrete Inputs - 讀取離散輸入狀態 (1-2000 個)
- 0x03: Read Holding Registers - 讀取保持暫存器 (1-125 個)
- 0x04: Read Input Registers - 讀取輸入暫存器 (1-125 個)

### 寫入功能
- 0x05: Write Single Coil - 寫入單個線圈
- 0x06: Write Single Register - 寫入單個暫存器
- 0x0F: Write Multiple Coils - 寫入多個線圈 (1-1968 個)
- 0x10: Write Multiple Registers - 寫入多個暫存器 (1-123 個)

## 使用範例

### Modbus TCP

	import "go-gateway/internal/protocol/modbus"

	// 建立 TCP 傳輸
	transport := modbus.NewTCPTransport("192.168.1.10", 502)
	client := modbus.NewClient(transport, 1) // Unit ID = 1

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

	// 讀取線圈
	coils, err := client.ReadCoils(0, 8)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Coils: %v\n", coils)

	// 寫入多個線圈
	err = client.WriteMultipleCoils(0, []bool{true, false, true, false})
	if err != nil {
		log.Fatal(err)
	}

### Modbus RTU

	// 建立 RTU 傳輸 (串列埠)
	transport := modbus.NewRTUTransport("COM3", 9600, 8, 1, "N", 2*time.Second)
	client := modbus.NewClient(transport, 1)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 使用方式與 TCP 相同
	registers, err := client.ReadHoldingRegisters(0, 10)

### Modbus UDP

	// 建立 UDP 傳輸
	transport := modbus.NewUDPTransport("192.168.1.10", 502)
	client := modbus.NewClient(transport, 1)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 使用方式與 TCP 相同

## 錯誤處理

所有 Modbus 操作可能返回以下錯誤類型：

- ProtocolError: Modbus 協議層級錯誤（異常碼）
- CommunicationError: 通訊層級錯誤（連線、逾時等）

	registers, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		if modbusErr, ok := err.(*modbus.ProtocolError); ok {
			fmt.Printf("Modbus 異常碼: 0x%02X\n", modbusErr.Code())
		} else {
			fmt.Printf("通訊錯誤: %v\n", err)
		}
	}

## 地址範圍

- 線圈 (Coils): 0x0000 - 0xFFFF (可讀寫)
- 離散輸入 (Discrete Inputs): 0x10000 - 0x1FFFF (只讀)
- 輸入暫存器 (Input Registers): 0x30000 - 0x3FFFF (只讀)
- 保持暫存器 (Holding Registers): 0x40000 - 0x4FFFF (可讀寫)

注意：實際使用時，地址參數應使用相對地址（0-65535），而非絕對地址。
*/
package modbus
