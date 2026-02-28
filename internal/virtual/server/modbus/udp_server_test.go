package modbus

import (
	"testing"
	"time"

	protocolmodbus "go-gateway/internal/protocol/modbus"
	"go-gateway/internal/virtual/memory"
)

func TestModbusUDPServer_StartStop(t *testing.T) {
	bank := memory.NewMemoryBank(2048)
	server := NewUDPServer(bank)

	if err := server.Start(0); err != nil {
		t.Fatalf("啟動 UDP 伺服器失敗: %v", err)
	}

	if server.Port() == 0 {
		t.Fatal("預期啟動後有可用埠號")
	}

	if err := server.Stop(); err != nil {
		t.Fatalf("停止 UDP 伺服器失敗: %v", err)
	}
}

func TestModbusUDPServer_ReadHoldingRegisters(t *testing.T) {
	bank := memory.NewMemoryBank(2048)
	_ = bank.WriteWord(0, 4321)

	server := NewUDPServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("啟動 UDP 伺服器失敗: %v", err)
	}
	defer server.Stop()

	time.Sleep(20 * time.Millisecond)

	transport := protocolmodbus.NewUDPTransport("127.0.0.1", server.Port())
	transport.Timeout = 2 * time.Second
	client := protocolmodbus.NewClient(transport, 1)

	if err := client.Connect(); err != nil {
		t.Fatalf("UDP client connect 失敗: %v", err)
	}
	defer client.Close()

	regs, err := client.ReadHoldingRegisters(0, 1)
	if err != nil {
		t.Fatalf("讀取保持暫存器失敗: %v", err)
	}
	if len(regs) != 1 {
		t.Fatalf("預期 1 筆暫存器，實際 %d", len(regs))
	}
	if regs[0] != 4321 {
		t.Fatalf("預期暫存器值 4321，實際 %d", regs[0])
	}
}

func TestModbusUDPServer_WriteSingleRegister(t *testing.T) {
	bank := memory.NewMemoryBank(2048)

	server := NewUDPServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("啟動 UDP 伺服器失敗: %v", err)
	}
	defer server.Stop()

	time.Sleep(20 * time.Millisecond)

	transport := protocolmodbus.NewUDPTransport("127.0.0.1", server.Port())
	transport.Timeout = 2 * time.Second
	client := protocolmodbus.NewClient(transport, 1)

	if err := client.Connect(); err != nil {
		t.Fatalf("UDP client connect 失敗: %v", err)
	}
	defer client.Close()

	if err := client.WriteSingleRegister(0, 9999); err != nil {
		t.Fatalf("寫入單一暫存器失敗: %v", err)
	}

	value, err := bank.ReadWord(0)
	if err != nil {
		t.Fatalf("讀取 bank 失敗: %v", err)
	}
	if value != 9999 {
		t.Fatalf("預期值 9999，實際 %d", value)
	}
}
