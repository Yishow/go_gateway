package fatek_test

import (
	"fmt"

	"go-gateway/internal/protocol/fatek"
)

// ExampleFatekClient 展示如何使用 TCP 客戶端
func ExampleFatekClient() {
	// 1. 建立傳輸層 (TCP)

	transport := fatek.NewTCPTransport("192.168.1.5", 500)

	// 2. 建立客戶端 (站號 1)
	client := fatek.NewClient(transport, 1)

	// 3. 連線
	if err := client.Connect(); err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer client.Close()

	// 4. 讀取 X0-X5 狀態
	status, err := client.ReadStatus("X", 0, 5)
	if err != nil {
		fmt.Printf("ReadStatus error: %v\n", err)
		return
	}
	fmt.Printf("X0-X5: %v\n", status)

	// 5. 讀取 D0-D9 暫存器
	regs, err := client.ReadRegisters("D", 0, 10)
	if err != nil {
		fmt.Printf("ReadRegisters error: %v\n", err)
		return
	}
	fmt.Printf("D0-D9: %v\n", regs)
}
