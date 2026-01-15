package modbus_test

import (
	"fmt"
	"time"

	"go-gateway/internal/protocol/modbus"
)

// ExampleCreateTCPClient 示範如何建立 Modbus TCP 客戶端
func ExampleCreateTCPClient() {
	// 建立 TCP 客戶端
	client := modbus.CreateTCPClient("192.168.1.10", 502, 1, 2*time.Second)

	// 連線
	if err := client.Connect(); err != nil {
		fmt.Printf("連線失敗: %v\n", err)
		return
	}
	defer client.Close()

	// 讀取保持暫存器 (地址 0, 讀取 10 個)
	registers, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		fmt.Printf("讀取暫存器失敗: %v\n", err)
		return
	}
	fmt.Printf("保持暫存器值: %v\n", registers)

	// Output:
	// 保持暫存器值: [0 0 0 0 0 0 0 0 0 0]
}

// ExampleCreateUDPClient 示範如何建立 Modbus UDP 客戶端
func ExampleCreateUDPClient() {
	// 建立 UDP 客戶端
	client := modbus.CreateUDPClient("192.168.1.10", 502, 1, 2*time.Second)

	if err := client.Connect(); err != nil {
		fmt.Printf("連線失敗: %v\n", err)
		return
	}
	defer client.Close()

	// 使用方式與 TCP 相同
	registers, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		fmt.Printf("讀取暫存器失敗: %v\n", err)
		return
	}
	fmt.Printf("保持暫存器值: %v\n", registers)

	// Output:
	// 保持暫存器值: [0 0 0 0 0 0 0 0 0 0]
}

// ExampleCreateRTUClient 示範如何建立 Modbus RTU 客戶端
func ExampleCreateRTUClient() {
	// 建立 RTU 客戶端 (串列埠)
	// Windows: "COM3", Linux: "/dev/ttyUSB0"
	client := modbus.CreateRTUClient("COM3", 9600, 8, 1, "N", 2*time.Second, 1)

	if err := client.Connect(); err != nil {
		fmt.Printf("連線失敗: %v\n", err)
		return
	}
	defer client.Close()

	// 使用方式與 TCP 相同
	registers, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		fmt.Printf("讀取暫存器失敗: %v\n", err)
		return
	}
	fmt.Printf("保持暫存器值: %v\n", registers)

	// Output:
	// 保持暫存器值: [0 0 0 0 0 0 0 0 0 0]
}

// ExampleModbusClient_ReadCoils 示範讀取線圈
func ExampleModbusClient_ReadCoils() {
	client := modbus.CreateTCPClient("192.168.1.10", 502, 1, 2*time.Second)
	if err := client.Connect(); err != nil {
		return
	}
	defer client.Close()

	// 讀取線圈 (地址 0, 讀取 8 個)
	coils, err := client.ReadCoils(0, 8)
	if err != nil {
		fmt.Printf("讀取線圈失敗: %v\n", err)
		return
	}
	fmt.Printf("線圈狀態: %v\n", coils)

	// Output:
	// 線圈狀態: [false false false false false false false false]
}

// ExampleModbusClient_WriteSingleRegister 示範寫入單個暫存器
func ExampleModbusClient_WriteSingleRegister() {
	client := modbus.CreateTCPClient("192.168.1.10", 502, 1, 2*time.Second)
	if err := client.Connect(); err != nil {
		return
	}
	defer client.Close()

	// 寫入單個暫存器
	if err := client.WriteSingleRegister(0, 12345); err != nil {
		fmt.Printf("寫入暫存器失敗: %v\n", err)
		return
	}
	fmt.Println("寫入暫存器成功")

	// Output:
	// 寫入暫存器成功
}
