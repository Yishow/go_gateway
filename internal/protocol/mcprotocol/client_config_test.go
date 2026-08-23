package mcprotocol

import (
	"encoding/binary"
	"testing"
	"time"
)

// TestClient_SetTimeout 驗證 SetTimeout 對支援的傳輸層生效，且不影響不支援的傳輸層
func TestClient_SetTimeout(t *testing.T) {
	tcp := NewTCPTransport("127.0.0.1", 5000)
	client := NewClientWithTransport(tcp)
	client.SetTimeout(3 * time.Second)
	if tcp.Timeout != 3*time.Second {
		t.Errorf("Expected TCP timeout 3s, got %v", tcp.Timeout)
	}

	// d <= 0 不應變更逾時
	client.SetTimeout(0)
	if tcp.Timeout != 3*time.Second {
		t.Errorf("Expected timeout unchanged for d=0, got %v", tcp.Timeout)
	}

	// 串列埠傳輸層同樣支援（不開啟實體埠，僅驗證欄位）
	serial := NewSerialTransport("COM1", 9600, 7, 2, "E", time.Second)
	serialClient := NewClientWithTransport(serial)
	serialClient.SetTimeout(7 * time.Second)
	if serial.Timeout != 7*time.Second {
		t.Errorf("Expected serial timeout 7s, got %v", serial.Timeout)
	}

	// 未支援 timeoutSetter 的傳輸層（Mock）不應 panic
	mockClient := NewClientWithTransport(NewMockMCTransport())
	mockClient.SetTimeout(time.Second)
}

// TestClient_SetFrame 驗證 SetFrame 完整覆寫請求標頭欄位並反映到封包
func TestClient_SetFrame(t *testing.T) {
	client := NewClient("127.0.0.1", 5000)

	frame := NewRequestFrame(0x01, 0x02, 0x03)
	frame.IONo = 0x0300
	client.SetFrame(frame)

	packet := client.frame.BuildPacket(CmdBatchRead, SubCmdWord, []byte{0x00})

	if got := binary.BigEndian.Uint16(packet[0:]); got != ReqSubHeader {
		t.Errorf("Expected subheader 0x%04X, got 0x%04X", ReqSubHeader, got)
	}
	if packet[2] != 0x01 {
		t.Errorf("Expected NetworkNo 0x01, got 0x%02X", packet[2])
	}
	if packet[3] != 0x02 {
		t.Errorf("Expected PCNo 0x02, got 0x%02X", packet[3])
	}
	if got := binary.LittleEndian.Uint16(packet[4:]); got != 0x0300 {
		t.Errorf("Expected IONo 0x0300, got 0x%04X", got)
	}
	if packet[6] != 0x03 {
		t.Errorf("Expected StationNo 0x03, got 0x%02X", packet[6])
	}
	// Timer 應保留 NewRequestFrame 的預設值（0x0010）
	if got := binary.LittleEndian.Uint16(packet[9:]); got != 0x0010 {
		t.Errorf("Expected default Timer 0x0010, got 0x%04X", got)
	}
}
