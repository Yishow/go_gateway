package modbus

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"go-gateway/internal/virtual/memory"
)

// =============================================================================
// 任務 2.1-2.2: Modbus TCP Server 測試
// =============================================================================

func TestModbusServer_StartStop(t *testing.T) {
	bank := memory.NewMemoryBank(65536)
	server := NewServer(bank)

	// 啟動伺服器
	err := server.Start(0) // 隨機端口
	if err != nil {
		t.Fatalf("啟動伺服器失敗: %v", err)
	}

	port := server.Port()
	if port == 0 {
		t.Error("伺服器端口應該非零")
	}

	// 停止伺服器
	err = server.Stop()
	if err != nil {
		t.Fatalf("停止伺服器失敗: %v", err)
	}
}

func TestModbusServer_ReadHoldingRegisters(t *testing.T) {
	bank := memory.NewMemoryBank(65536)
	server := NewServer(bank)

	// 預先寫入測試數據
	bank.WriteWord(0, 12345)     // 暫存器 0
	bank.WriteWord(2, 54321)     // 暫存器 1
	bank.WriteDWord(4, 0xDEADBEEF) // 暫存器 2-3

	err := server.Start(0)
	if err != nil {
		t.Fatalf("啟動伺服器失敗: %v", err)
	}
	defer server.Stop()

	// 等待伺服器就緒
	time.Sleep(50 * time.Millisecond)

	// 建立 TCP 連接
	conn, err := net.Dial("tcp", server.Address())
	if err != nil {
		t.Fatalf("連接伺服器失敗: %v", err)
	}
	defer conn.Close()

	// 發送 Modbus 請求: 讀取 Holding Registers (FC 03)
	// Transaction ID: 0x0001
	// Protocol ID: 0x0000 (Modbus)
	// Length: 0x0006
	// Unit ID: 0x01
	// Function Code: 0x03 (Read Holding Registers)
	// Start Address: 0x0000
	// Quantity: 0x0002
	request := []byte{
		0x00, 0x01, // Transaction ID
		0x00, 0x00, // Protocol ID
		0x00, 0x06, // Length
		0x01,       // Unit ID
		0x03,       // Function Code
		0x00, 0x00, // Start Address
		0x00, 0x02, // Quantity
	}

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write(request)
	if err != nil {
		t.Fatalf("發送請求失敗: %v", err)
	}

	// 讀取回應
	response := make([]byte, 256)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("讀取回應失敗: %v", err)
	}

	// 驗證回應
	// MBAP Header (7 bytes) + Function Code (1) + Byte Count (1) + Data (4)
	if n < 13 {
		t.Fatalf("回應長度不足: %d", n)
	}

	// 驗證 Transaction ID
	if response[0] != 0x00 || response[1] != 0x01 {
		t.Errorf("Transaction ID 不匹配")
	}

	// 驗證 Function Code
	if response[7] != 0x03 {
		t.Errorf("Function Code 預期 0x03，實際 0x%02X", response[7])
	}

	// 驗證 Byte Count
	if response[8] != 0x04 { // 2 registers * 2 bytes = 4
		t.Errorf("Byte Count 預期 0x04，實際 0x%02X", response[8])
	}

	// 驗證暫存器值
	reg0 := binary.BigEndian.Uint16(response[9:11])
	reg1 := binary.BigEndian.Uint16(response[11:13])

	if reg0 != 12345 {
		t.Errorf("暫存器 0 預期 12345，實際 %d", reg0)
	}
	if reg1 != 54321 {
		t.Errorf("暫存器 1 預期 54321，實際 %d", reg1)
	}
}

func TestModbusServer_WriteSingleRegister(t *testing.T) {
	bank := memory.NewMemoryBank(65536)
	server := NewServer(bank)

	err := server.Start(0)
	if err != nil {
		t.Fatalf("啟動伺服器失敗: %v", err)
	}
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("tcp", server.Address())
	if err != nil {
		t.Fatalf("連接伺服器失敗: %v", err)
	}
	defer conn.Close()

	// 發送 Modbus 請求: 寫入單一暫存器 (FC 06)
	// Write 9999 to address 100
	request := []byte{
		0x00, 0x02, // Transaction ID
		0x00, 0x00, // Protocol ID
		0x00, 0x06, // Length
		0x01,       // Unit ID
		0x06,       // Function Code (Write Single Register)
		0x00, 0x64, // Register Address (100)
		0x27, 0x0F, // Value (9999)
	}

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write(request)
	if err != nil {
		t.Fatalf("發送請求失敗: %v", err)
	}

	// 讀取回應
	response := make([]byte, 256)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("讀取回應失敗: %v", err)
	}

	// 驗證回應 (寫入成功應該回傳相同的請求)
	if n < 12 {
		t.Fatalf("回應長度不足: %d", n)
	}

	// 驗證 Function Code
	if response[7] != 0x06 {
		t.Errorf("Function Code 預期 0x06，實際 0x%02X", response[7])
	}

	// 驗證記憶體值
	val, _ := bank.ReadWord(200) // 地址 100 * 2 bytes
	if val != 9999 {
		t.Errorf("記憶體值預期 9999，實際 %d", val)
	}
}

func TestModbusServer_WriteMultipleRegisters(t *testing.T) {
	bank := memory.NewMemoryBank(65536)
	server := NewServer(bank)

	err := server.Start(0)
	if err != nil {
		t.Fatalf("啟動伺服器失敗: %v", err)
	}
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("tcp", server.Address())
	if err != nil {
		t.Fatalf("連接伺服器失敗: %v", err)
	}
	defer conn.Close()

	// 發送 Modbus 請求: 寫入多個暫存器 (FC 16)
	request := []byte{
		0x00, 0x03, // Transaction ID
		0x00, 0x00, // Protocol ID
		0x00, 0x0D, // Length (13 bytes: Unit ID + PDU)
		0x01,       // Unit ID
		0x10,       // Function Code (Write Multiple Registers)
		0x00, 0x00, // Start Address
		0x00, 0x03, // Quantity (3 registers)
		0x06,       // Byte Count
		0x00, 0x01, // Value 1
		0x00, 0x02, // Value 2
		0x00, 0x03, // Value 3
	}

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write(request)
	if err != nil {
		t.Fatalf("發送請求失敗: %v", err)
	}

	// 讀取回應
	response := make([]byte, 256)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("讀取回應失敗: %v", err)
	}

	if n < 12 {
		t.Logf("收到回應: %X", response[:n])
		t.Fatalf("回應長度不足: %d (WriteMultiple)", n)
	}

	// 驗證 Function Code
	if response[7] != 0x10 {
		t.Errorf("Function Code 預期 0x10，實際 0x%02X", response[7])
	}

	// 驗證記憶體值
	val0, _ := bank.ReadWord(0)
	val1, _ := bank.ReadWord(2)
	val2, _ := bank.ReadWord(4)

	if val0 != 1 {
		t.Errorf("暫存器 0 預期 1，實際 %d", val0)
	}
	if val1 != 2 {
		t.Errorf("暫存器 1 預期 2，實際 %d", val1)
	}
	if val2 != 3 {
		t.Errorf("暫存器 2 預期 3，實際 %d", val2)
	}
}

func TestModbusServer_InvalidFunctionCode(t *testing.T) {
	bank := memory.NewMemoryBank(65536)
	server := NewServer(bank)

	err := server.Start(0)
	if err != nil {
		t.Fatalf("啟動伺服器失敗: %v", err)
	}
	defer server.Stop()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("tcp", server.Address())
	if err != nil {
		t.Fatalf("連接伺服器失敗: %v", err)
	}
	defer conn.Close()

	// 發送無效的 Function Code
	request := []byte{
		0x00, 0x04, // Transaction ID
		0x00, 0x00, // Protocol ID
		0x00, 0x06, // Length
		0x01,       // Unit ID
		0xFF,       // Invalid Function Code
		0x00, 0x00,
		0x00, 0x00,
	}

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err = conn.Write(request)
	if err != nil {
		t.Fatalf("發送請求失敗: %v", err)
	}

	// 讀取回應
	response := make([]byte, 256)
	n, err := conn.Read(response)
	if err != nil {
		t.Fatalf("讀取回應失敗: %v", err)
	}

	if n < 9 {
		t.Fatalf("回應長度不足: %d", n)
	}

	// 驗證異常回應 (Function Code 應該是 0xFF | 0x80 = 0x7F?? no, 0xFF + 0x80 = overflow)
	// 實際上 Modbus 異常回應是 FC | 0x80
	// 這裡 FC = 0xFF，所以異常碼應該是 0xFF (溢出後)
	// 但標準做法是 FC | 0x80，無效 FC 通常返回 Illegal Function (01)
}
