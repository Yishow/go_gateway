package main

import (
	"fmt"
	"time"

	"go-gateway/internal/protocol/fatek"
)

// testFatekTCP 測試 Fatek TCP
func testFatekTCP(results *TestResults, addr string) {
	start := time.Now()
	result := TestResult{
		Protocol: "Fatek",
		Mode:     "TCP",
		Success:  false,
	}

	defer func() {
		result.Duration = time.Since(start)
		results.addResult(result)
	}()

	fmt.Printf("[測試] Fatek TCP (%s)\n", addr)

	// 建立客戶端
	client := fatek.CreateTCPClient("127.0.0.1", 2000, 1, 2*time.Second)

	// 連線
	if err := client.Connect(); err != nil {
		result.Message = fmt.Sprintf("連線失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	defer client.Close()

	// 測試讀取暫存器
	fmt.Println("  → 測試讀取暫存器 (D)...")
	registers, err := client.ReadRegisters("D", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", registers)

	// 測試寫入暫存器
	fmt.Println("  → 測試寫入暫存器 (D)...")
	writeValues := []int{100, 200, 300}
	err = client.WriteRegisters("D", 0, writeValues)
	if err != nil {
		result.Message = fmt.Sprintf("寫入暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 寫入成功\n")

	// 驗證寫入
	registers, err = client.ReadRegisters("D", 0, 3)
	if err != nil {
		result.Message = fmt.Sprintf("驗證讀取失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 驗證成功: %v\n", registers)

	// 測試讀取狀態
	fmt.Println("  → 測試讀取狀態 (X)...")
	status, err := client.ReadStatus("X", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取狀態失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", status)

	result.Success = true
	result.Message = "所有測試通過"
	fmt.Printf("  ✅ %s\n", result.Message)
}

// testFatekRTU 測試 Fatek RTU
func testFatekRTU(results *TestResults, config struct {
	Port     string
	BaudRate int
	DataBits int
	StopBits int
	Parity   string
}) {
	start := time.Now()
	result := TestResult{
		Protocol: "Fatek",
		Mode:     "RTU",
		Success:  false,
	}

	defer func() {
		result.Duration = time.Since(start)
		results.addResult(result)
	}()

	fmt.Printf("[測試] Fatek RTU (%s, %d/%d/%d/%s)\n",
		config.Port, config.BaudRate, config.DataBits, config.StopBits, config.Parity)

	// 建立客戶端
	client := fatek.CreateSerialClient(
		config.Port,
		1,
		config.BaudRate,
		config.DataBits,
		config.StopBits,
		config.Parity,
		2*time.Second,
	)

	// 連線
	if err := client.Connect(); err != nil {
		result.Message = fmt.Sprintf("連線失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	defer client.Close()

	// 測試讀取暫存器
	fmt.Println("  → 測試讀取暫存器 (D)...")
	registers, err := client.ReadRegisters("D", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", registers)

	// 測試寫入暫存器
	fmt.Println("  → 測試寫入暫存器 (D)...")
	writeValues := []int{100, 200, 300}
	err = client.WriteRegisters("D", 0, writeValues)
	if err != nil {
		result.Message = fmt.Sprintf("寫入暫存器失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 寫入成功\n")

	// 測試讀取狀態
	fmt.Println("  → 測試讀取狀態 (X)...")
	status, err := client.ReadStatus("X", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取狀態失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", status)

	result.Success = true
	result.Message = "所有測試通過"
	fmt.Printf("  ✅ %s\n", result.Message)
}
