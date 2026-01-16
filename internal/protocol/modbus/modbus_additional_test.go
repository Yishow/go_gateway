package modbus

import (
	"testing"
)

// 測試更多 frame 解析邊界情況
func TestParseReadResponse_EdgeCases(t *testing.T) {
	// 測試空數據
	_, err := ParseReadResponse([]byte{})
	if err != ErrResponseTooShort {
		t.Errorf("Expected ErrResponseTooShort, got %v", err)
	}

	// 測試 ByteCount 為 0
	data := []byte{0}
	_, err = ParseReadResponse(data)
	if err != nil {
		t.Logf("ParseReadResponse with byteCount 0: %v", err)
	}

	// 測試數據長度不足
	data = []byte{10, 0x01, 0x02} // ByteCount=10, 但只有2字節數據
	_, err = ParseReadResponse(data)
	if err != ErrResponseTooShort {
		t.Errorf("Expected ErrResponseTooShort, got %v", err)
	}
}

// 測試 RTU Frame 邊界情況
func TestParseRTUFrame_EdgeCases(t *testing.T) {
	// 測試過短的數據
	_, _, _, err := ParseRTUFrame([]byte{0x01})
	if err != ErrResponseTooShort {
		t.Errorf("Expected ErrResponseTooShort, got %v", err)
	}

	// 測試 CRC 錯誤
	frame := BuildRTUFrame(1, 0x03, []byte{0x00, 0x0A})
	frame[len(frame)-1] ^= 0xFF // 破壞 CRC
	_, _, _, err = ParseRTUFrame(frame)
	if err != ErrCRCError {
		t.Errorf("Expected ErrCRCError, got %v", err)
	}
}

// 測試 TCP Frame 邊界情況
func TestParseTCPFrame_EdgeCases(t *testing.T) {
	// 測試過短的數據
	_, _, _, _, err := ParseTCPFrame([]byte{0x01})
	if err == nil {
		t.Error("Expected error for short frame")
	}

	// 測試無效的協議 ID
	header := BuildMBAPHeader(0x1234, 1, 10)
	header[2] = 0xFF // 破壞協議 ID
	header[3] = 0xFF
	_, _, _, _, err = ParseTCPFrame(header)
	if err == nil {
		t.Error("Expected error for invalid protocol ID")
	}
}

// 測試 PDU 解析邊界情況
func TestParsePDU_EdgeCases(t *testing.T) {
	// 測試空數據
	_, _, err := ParsePDU([]byte{})
	if err != ErrResponseTooShort {
		t.Errorf("Expected ErrResponseTooShort, got %v", err)
	}

	// 測試異常回應但數據不足
	_, _, err = ParsePDU([]byte{0x83}) // 異常功能碼但缺少異常碼
	if err == nil {
		t.Error("Expected error for incomplete exception response")
	}
}

// 測試 CRC 計算
func TestCalculateCRC16_Various(t *testing.T) {
	// 測試空數據
	crc := CalculateCRC16([]byte{})
	if crc != 0xFFFF {
		t.Errorf("Expected CRC 0xFFFF for empty data, got 0x%04X", crc)
	}

	// 測試單字節
	crc = CalculateCRC16([]byte{0x01})
	if crc == 0 {
		t.Error("CRC should not be zero")
	}

	// 測試多字節
	data := []byte{0x01, 0x03, 0x00, 0x00, 0x00, 0x0A}
	crc1 := CalculateCRC16(data)
	crc2 := CalculateCRC16(data)
	if crc1 != crc2 {
		t.Error("CRC should be deterministic")
	}
}

// 測試 UnpackBits 邊界情況
func TestUnpackBits_EdgeCases(t *testing.T) {
	// 測試空數據
	bits := UnpackBits([]byte{}, 0)
	if len(bits) != 0 {
		t.Errorf("Expected 0 bits, got %d", len(bits))
	}

	// 測試數據不足
	bits = UnpackBits([]byte{0x01}, 20)
	if len(bits) != 20 {
		t.Errorf("Expected 20 bits, got %d", len(bits))
	}
	// 超出數據範圍的位元應該是 false
	if bits[10] {
		t.Error("Bit beyond data should be false")
	}

	// 測試奇數個位元
	bits = UnpackBits([]byte{0xFF}, 7)
	if len(bits) != 7 {
		t.Errorf("Expected 7 bits, got %d", len(bits))
	}
}

// 測試 PackBits 邊界情況
func TestPackBits_EdgeCases(t *testing.T) {
	// 測試空數組
	data := PackBits([]bool{})
	if len(data) != 0 {
		t.Errorf("Expected 0 bytes, got %d", len(data))
	}

	// 測試單個位元
	data = PackBits([]bool{true})
	if len(data) != 1 {
		t.Errorf("Expected 1 byte, got %d", len(data))
	}
	if data[0]&0x01 == 0 {
		t.Error("First bit should be set")
	}

	// 測試奇數個位元
	bits := []bool{true, false, true, false, true, false, true}
	data = PackBits(bits)
	if len(data) != 1 {
		t.Errorf("Expected 1 byte for 7 bits, got %d", len(data))
	}
}

// 測試 BuildRTUFrame 和 ParseRTUFrame 往返
func TestRTUFrame_RoundTrip(t *testing.T) {
	originalData := []byte{0x00, 0x0A, 0x00, 0x0A}
	frame := BuildRTUFrame(1, 0x03, originalData)

	addr, funcCode, pduData, err := ParseRTUFrame(frame)
	if err != nil {
		t.Fatalf("ParseRTUFrame failed: %v", err)
	}
	if addr != 1 {
		t.Errorf("Expected address 1, got %d", addr)
	}
	if funcCode != 0x03 {
		t.Errorf("Expected function code 0x03, got 0x%02X", funcCode)
	}
	if len(pduData) != len(originalData) {
		t.Errorf("Expected PDU data length %d, got %d", len(originalData), len(pduData))
	}
}

// 測試 BuildTCPFrame 和 ParseTCPFrame 往返
func TestTCPFrame_RoundTrip(t *testing.T) {
	originalData := []byte{0x00, 0x0A, 0x00, 0x0A}
	frame := BuildTCPFrame(0x1234, 1, 0x03, originalData)

	transID, unitID, funcCode, pduData, err := ParseTCPFrame(frame)
	if err != nil {
		t.Fatalf("ParseTCPFrame failed: %v", err)
	}
	if transID != 0x1234 {
		t.Errorf("Expected transaction ID 0x1234, got 0x%04X", transID)
	}
	if unitID != 1 {
		t.Errorf("Expected unit ID 1, got %d", unitID)
	}
	if funcCode != 0x03 {
		t.Errorf("Expected function code 0x03, got 0x%02X", funcCode)
	}
	if len(pduData) != len(originalData) {
		t.Errorf("Expected PDU data length %d, got %d", len(originalData), len(pduData))
	}
}

// 測試異常碼處理
func TestExceptionCodes(t *testing.T) {
	exceptionCodes := []byte{
		ExceptionIllegalFunction,
		ExceptionIllegalDataAddress,
		ExceptionIllegalDataValue,
		ExceptionServerDeviceFailure,
		ExceptionAcknowledge,
		ExceptionServerDeviceBusy,
		ExceptionMemoryParityError,
		ExceptionGatewayPathUnavailable,
		ExceptionGatewayTargetNoResponse,
	}

	for _, code := range exceptionCodes {
		err := NewProtocolError(code, "Test error")
		if err == nil {
			t.Errorf("NewProtocolError returned nil for code 0x%02X", code)
			continue
		}
		if err.Code() != code {
			t.Errorf("Expected error code 0x%02X, got 0x%02X", code, err.Code())
		}
		if err.Error() == "" {
			t.Errorf("Error message should not be empty for code 0x%02X", code)
		}
	}
}

// 測試錯誤訊息
func TestExceptionMessages(t *testing.T) {
	for code, msg := range ExceptionMessages {
		if msg == "" {
			t.Errorf("Exception message should not be empty for code 0x%02X", code)
		}
	}
}

// 測試 Transport GetNextTransactionID
func TestTCPTransport_GetNextTransactionID(t *testing.T) {
	transport := NewTCPTransport("127.0.0.1", 502)

	id1 := transport.GetNextTransactionID()
	id2 := transport.GetNextTransactionID()

	if id1 == 0 {
		t.Error("Transaction ID should not be zero")
	}
	if id2 != id1+1 {
		t.Errorf("Expected transaction ID %d, got %d", id1+1, id2)
	}

	// 測試循環（如果達到最大值）
	for i := 0; i < 65535; i++ {
		transport.GetNextTransactionID()
	}
	id := transport.GetNextTransactionID()
	if id == 0 {
		t.Error("Transaction ID should wrap to 1, not 0")
	}
}

func TestUDPTransport_GetNextTransactionID(t *testing.T) {
	transport := NewUDPTransport("127.0.0.1", 502)

	id1 := transport.GetNextTransactionID()
	id2 := transport.GetNextTransactionID()

	if id1 == 0 {
		t.Error("Transaction ID should not be zero")
	}
	if id2 != id1+1 {
		t.Errorf("Expected transaction ID %d, got %d", id1+1, id2)
	}
}

// 測試工廠函數的預設值
func TestFactoryDefaults(t *testing.T) {
	// TCP 客戶端預設端口
	client := CreateTCPClient("127.0.0.1", 0, 1, 0)
	if client == nil {
		t.Fatal("CreateTCPClient returned nil")
	}

	// UDP 客戶端預設端口
	client = CreateUDPClient("127.0.0.1", 0, 1, 0)
	if client == nil {
		t.Fatal("CreateUDPClient returned nil")
	}

	// RTU 客戶端預設值
	client = CreateRTUClient("COM1", 0, 0, 0, "", 0, 0)
	if client == nil {
		t.Fatal("CreateRTUClient returned nil")
	}
}
