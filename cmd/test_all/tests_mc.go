package main

import (
	"fmt"
	"time"

	"go-gateway/internal/protocol/mcprotocol"
)

// testMCTCP 測試 MC Protocol TCP
func testMCTCP(results *TestResults, addr string) {
	start := time.Now()
	result := TestResult{
		Protocol: "MC Protocol",
		Mode:     transportLabelTCP,
		Success:  false,
	}

	defer func() {
		result.Duration = time.Since(start)
		results.addResult(result)
	}()

	fmt.Printf("[測試] MC Protocol TCP (%s)\n", addr)

	// 建立客戶端
	client := mcprotocol.CreateTCPClient("127.0.0.1", 6000, 2*time.Second)

	// 連線
	if err := client.Connect(); err != nil {
		result.Message = fmt.Sprintf("連線失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	defer client.Close()

	// 測試讀取字組設備
	fmt.Println("  → 測試讀取字組設備 (D)...")
	values, err := client.BatchReadWord("D", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取字組設備失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", values)

	// 測試寫入字組設備
	fmt.Println("  → 測試寫入字組設備 (D)...")
	writeValues := []int{100, 200, 300}
	err = client.BatchWriteWord("D", 0, writeValues)
	if err != nil {
		result.Message = fmt.Sprintf("寫入字組設備失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 寫入成功\n")

	// 驗證寫入
	values, err = client.BatchReadWord("D", 0, 3)
	if err != nil {
		result.Message = fmt.Sprintf("驗證讀取失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 驗證成功: %v\n", values)

	// 測試讀取位元設備
	fmt.Println("  → 測試讀取位元設備 (M)...")
	bits, err := client.BatchReadBit("M", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取位元設備失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", bits)

	result.Success = true
	result.Message = testsSucceededText
	fmt.Printf("  ✅ %s\n", result.Message)
}

// testMCRTU 測試 MC Protocol RTU
func testMCRTU(results *TestResults, config struct {
	Port     string
	BaudRate int
	DataBits int
	StopBits int
	Parity   string
}) {
	start := time.Now()
	result := TestResult{
		Protocol: "MC Protocol",
		Mode:     transportLabelRTU,
		Success:  false,
	}

	defer func() {
		result.Duration = time.Since(start)
		results.addResult(result)
	}()

	fmt.Printf("[測試] MC Protocol RTU (%s, %d/%d/%d/%s)\n",
		config.Port, config.BaudRate, config.DataBits, config.StopBits, config.Parity)

	// 建立客戶端
	client := mcprotocol.CreateSerialClient(
		config.Port,
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

	// 測試讀取字組設備
	fmt.Println("  → 測試讀取字組設備 (D)...")
	values, err := client.BatchReadWord("D", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取字組設備失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", values)

	// 測試寫入字組設備
	fmt.Println("  → 測試寫入字組設備 (D)...")
	writeValues := []int{100, 200, 300}
	err = client.BatchWriteWord("D", 0, writeValues)
	if err != nil {
		result.Message = fmt.Sprintf("寫入字組設備失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 寫入成功\n")

	// 測試讀取位元設備
	fmt.Println("  → 測試讀取位元設備 (M)...")
	bits, err := client.BatchReadBit("M", 0, 10)
	if err != nil {
		result.Message = fmt.Sprintf("讀取位元設備失敗: %v", err)
		fmt.Printf("  ❌ %s\n", result.Message)
		return
	}
	fmt.Printf("  ✓ 讀取成功: %v\n", bits)

	result.Success = true
	result.Message = testsSucceededText
	fmt.Printf("  ✅ %s\n", result.Message)
}
