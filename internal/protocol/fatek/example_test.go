package fatek_test

import (
	"fmt"
	"time"

	"go-gateway/internal/protocol/fatek"
)

// ExampleClient_TCP 展示如何使用 TCP 客戶端
func ExampleClient_TCP() {
	// 建立 TCP 客戶端
	client := fatek.CreateTCPClient("192.168.1.5", 500, 1, 2*time.Second)

	// 連線
	err := client.Connect()
	if err != nil {
		fmt.Printf("連線失敗: %v\n", err)
		return
	}
	defer client.Close()

	// 讀取 X0-X5
	status, err := client.ReadStatus("X", 0, 5)
	if err != nil {
		fmt.Printf("讀取失敗: %v\n", err)
		return
	}
	fmt.Printf("X0-X5 狀態: %v\n", status)

	// 讀取 D0-D9
	registers, err := client.ReadRegisters("D", 0, 10)
	if err != nil {
		fmt.Printf("讀取失敗: %v\n", err)
		return
	}
	fmt.Printf("D0-D9 值: %v\n", registers)
}

// ExampleClient_Serial 展示如何使用串列埠客戶端
func ExampleClient_Serial() {
	// 建立串列埠客戶端
	client := fatek.CreateSerialClient(
		"COM3",        // 埠號
		1,             // 站號
		9600,          // 波特率
		7,             // 資料位元
		1,             // 停止位元
		"E",           // 同位檢查
		1*time.Second, // 逾時
	)

	// 連線
	err := client.Connect()
	if err != nil {
		fmt.Printf("連線失敗: %v\n", err)
		return
	}
	defer client.Close()

	// 寫入 Y0-Y2
	data := []bool{true, false, true}
	err = client.WriteStatus("Y", 0, data)
	if err != nil {
		fmt.Printf("寫入失敗: %v\n", err)
		return
	}
	fmt.Println("寫入成功")
}
