package main

import (
	"fmt"
	"os"
	"time"
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
			Port:     defaultSerialPort,
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
			Port:     defaultSerialPort,
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
			Port:     defaultSerialPort,
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
