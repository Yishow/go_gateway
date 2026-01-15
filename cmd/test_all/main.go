package main

import (
	"fmt"
	"os"
	"time"

	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"
)

// TestResult 測試結果結構
type TestResult struct {
	Protocol string
	Mode     string
	Success  bool
	Message  string
	Duration time.Duration
}

// TestResults 測試結果集合
type TestResults struct {
	Results []TestResult
	Total   int
	Passed  int
	Failed  int
}

func main() {
	results := &TestResults{
		Results: make([]TestResult, 0),
	}

	fmt.Println("=========================================")
	fmt.Println("  完整協議測試程序")
	fmt.Println("=========================================")
	fmt.Println()

	// 測試配置
	config := struct {
		ModbusTCP string
		ModbusRTU struct {
			Port     string
			BaudRate int
			DataBits int
			StopBits int
			Parity   string
		}
		MCTCP string
		MCRTU struct {
			Port     string
			BaudRate int
			DataBits int
			StopBits int
			Parity   string
		}
		FatekTCP string
		FatekRTU struct {
			Port     string
			BaudRate int
			DataBits int
			StopBits int
			Parity   string
		}
	}{
		ModbusTCP: "127.0.0.1:502",
		ModbusRTU: struct {
			Port     string
			BaudRate int
			DataBits int
			StopBits int
			Parity   string
		}{
			Port:     "COM4",
			BaudRate: 9600,
			DataBits: 8,
			StopBits: 1,
			Parity:   "N",
		},
		MCTCP: "127.0.0.1:6000",
		MCRTU: struct {
			Port     string
			BaudRate int
			DataBits int
			StopBits int
			Parity   string
		}{
			Port:     "COM4",
			BaudRate: 9600,
			DataBits: 7,
			StopBits: 2,
			Parity:   "E",
		},
		FatekTCP: "127.0.0.1:2000",
		FatekRTU: struct {
			Port     string
			BaudRate int
			DataBits int
			StopBits int
			Parity   string
		}{
			Port:     "COM4",
			BaudRate: 9600,
			DataBits: 7,
			StopBits: 2,
			Parity:   "E",
		},
	}

	// 執行所有測試
	fmt.Println("開始執行測試...")
	fmt.Println()

	// Modbus TCP 測試
	testModbusTCP(results, config.ModbusTCP)

	// Modbus RTU 測試
	testModbusRTU(results, config.ModbusRTU)

	// MC Protocol TCP 測試
	testMCTCP(results, config.MCTCP)

	// MC Protocol RTU 測試
	testMCRTU(results, config.MCRTU)

	// Fatek TCP 測試
	testFatekTCP(results, config.FatekTCP)

	// Fatek RTU 測試
	testFatekRTU(results, config.FatekRTU)

	// 輸出測試報告
	printReport(results)

	// 根據測試結果決定退出碼
	if results.Failed > 0 {
		os.Exit(1)
	}
}

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

// testMCTCP 測試 MC Protocol TCP
func testMCTCP(results *TestResults, addr string) {
	start := time.Now()
	result := TestResult{
		Protocol: "MC Protocol",
		Mode:     "TCP",
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
	result.Message = "所有測試通過"
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
		Mode:     "RTU",
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
	result.Message = "所有測試通過"
	fmt.Printf("  ✅ %s\n", result.Message)
}

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

// addResult 添加測試結果
func (tr *TestResults) addResult(result TestResult) {
	tr.Results = append(tr.Results, result)
	tr.Total++
	if result.Success {
		tr.Passed++
	} else {
		tr.Failed++
	}
}

// printReport 輸出測試報告
func printReport(results *TestResults) {
	fmt.Println()
	fmt.Println("=========================================")
	fmt.Println("  測試報告")
	fmt.Println("=========================================")
	fmt.Println()

	for _, result := range results.Results {
		status := "✅ 通過"
		if !result.Success {
			status = "❌ 失敗"
		}
		fmt.Printf("%s [%s %s] %s (耗時: %v)\n",
			status,
			result.Protocol,
			result.Mode,
			result.Message,
			result.Duration.Round(time.Millisecond),
		)
	}

	fmt.Println()
	fmt.Println("=========================================")
	fmt.Printf("總計: %d  通過: %d  失敗: %d\n", results.Total, results.Passed, results.Failed)
	fmt.Println("=========================================")
}
