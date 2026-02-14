package main

import (
	"fmt"
	"time"

	"go-gateway/internal/protocol/modbus"
)

// testModbusTCP 測試 Modbus TCP
func testModbusTCP(results *TestResults, addr string) {
	start := time.Now()
	result := TestResult{
		Protocol: "Modbus",
		Mode:     "TCP",
		Success:  false,
	}

	defer func() {
		result.Duration = time.Since(start)
		results.addResult(result)
	}()

	fmt.Printf("[測試] Modbus TCP (%s)\n", addr)

	// 建立客戶端
	client := modbus.CreateTCPClient("127.0.0.1", 502, 1, 2*time.Second)

	// 連線
	if err := client.Connect(); err != nil {
		result.Message = fmt.Sprintf("連線失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	defer client.Close()

	// 測試讀取保持暫存器
	fmt.Println("  → 測試讀取保持暫存器...")
	registers, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取保持暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", registers)

	// 測試寫入單個暫存器
	fmt.Println("  → 測試寫入單個暫存器...")
	err = client.WriteSingleRegister(0, 12345)
	if err != nil {
		result.Message = fmt.Sprintf("寫入暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 寫入成功\n")

	// 驗證寫入
	registers, err = client.ReadHoldingRegisters(0, 1)
	if err != nil {
		result.Message = fmt.Sprintf("驗證讀取失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	if len(registers) > 0 && registers[0] == 12345 {
		fmt.Printf("  ✓ 驗證成功: 值為 %d\n", registers[0])
	}

	// 測試讀取線圈
	fmt.Println("  → 測試讀取線圈...")
	coils, err := client.ReadCoils(0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取線圈失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", coils)

	result.Success = true
	result.Message = "所有測試通過"
	fmt.Printf("  ✅ %s\n", result.Message)
}

// testModbusRTU 測試 Modbus RTU
func testModbusRTU(results *TestResults, config struct {
	Port     string
	BaudRate int
	DataBits int
	StopBits int
	Parity   string
}) {
	start := time.Now()
	result := TestResult{
		Protocol: "Modbus",
		Mode:     "RTU",
		Success:  false,
	}

	defer func() {
		result.Duration = time.Since(start)
		results.addResult(result)
	}()

	fmt.Printf("[測試] Modbus RTU (%s, %d/%d/%d/%s)\n",
		config.Port, config.BaudRate, config.DataBits, config.StopBits, config.Parity)

	// 建立客戶端
	client := modbus.CreateRTUClient(
		config.Port,
		config.BaudRate,
		config.DataBits,
		config.StopBits,
		config.Parity,
		2*time.Second,
		1,
	)

	// 連線
	if err := client.Connect(); err != nil {
		result.Message = fmt.Sprintf("連線失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	defer client.Close()

	// 測試讀取保持暫存器
	fmt.Println("  → 測試讀取保持暫存器...")
	registers, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取保持暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", registers)

	// 測試寫入單個暫存器
	fmt.Println("  → 測試寫入單個暫存器...")
	err = client.WriteSingleRegister(0, 12345)
	if err != nil {
		result.Message = fmt.Sprintf("寫入暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 寫入成功\n")

	// 測試讀取線圈
	fmt.Println("  → 測試讀取線圈...")
	coils, err := client.ReadCoils(0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取線圈失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", coils)

	result.Success = true
	result.Message = "所有測試通過"
	fmt.Printf("  ✅ %s\n", result.Message)
}
