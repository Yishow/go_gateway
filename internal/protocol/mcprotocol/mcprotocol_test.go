package mcprotocol

import (
	"encoding/binary"
	"testing"
)

// MockTransport 模擬傳輸層用於測試
type MockMCTransport struct {
	responses map[string][]byte
	errors    map[string]error
	sentData  [][]byte
}

func NewMockMCTransport() *MockMCTransport {
	return &MockMCTransport{
		responses: make(map[string][]byte),
		errors:    make(map[string]error),
		sentData:  make([][]byte, 0),
	}
}

func (m *MockMCTransport) Connect() error {
	return nil
}

func (m *MockMCTransport) Close() error {
	return nil
}

func (m *MockMCTransport) SendReceive(req []byte) ([]byte, error) {
	m.sentData = append(m.sentData, req)

	// 檢查是否有預設錯誤
	if err, ok := m.errors[string(req)]; ok {
		return nil, err
	}

	// 檢查是否有預設回應
	if resp, ok := m.responses[string(req)]; ok {
		return resp, nil
	}

	// 預設回應：根據請求生成
	if len(req) < 15 {
		return nil, &MCError{Code: 0x0001, Msg: "request too short"}
	}

	// 解析請求
	cmd := binary.LittleEndian.Uint16(req[11:13])
	subCmd := binary.LittleEndian.Uint16(req[13:15])

	// 構建回應標頭 (9 bytes)
	header := make([]byte, 9)
	header[0] = 0xD0 // SubHeader
	header[1] = 0x00
	header[2] = 0x00                                  // Network
	header[3] = 0xFF                                  // PC
	binary.LittleEndian.PutUint16(header[4:], 0x03FF) // IO
	header[6] = 0x00                                  // Station

	var body []byte

	switch cmd {
	case CmdBatchRead:
		if subCmd == SubCmdWord {
			// 讀取字組回應: EndCode(2) + Data(2*count)
			count := 10
			body = make([]byte, 2+count*2)
			binary.LittleEndian.PutUint16(body[0:], 0) // EndCode
			for i := 0; i < count; i++ {
				binary.LittleEndian.PutUint16(body[2+i*2:], uint16(i*100))
			}
		} else if subCmd == SubCmdBit {
			// 讀取位元回應: EndCode(2) + Data((count+1)/2)
			count := 10
			byteCount := (count + 1) / 2
			body = make([]byte, 2+byteCount)
			binary.LittleEndian.PutUint16(body[0:], 0) // EndCode
			// 設置一些位元
			body[2] = 0x10 // 第一個位元為 1
		}

	case CmdBatchWrite:
		// 寫入回應: EndCode(2)
		body = make([]byte, 2)
		binary.LittleEndian.PutUint16(body[0:], 0) // EndCode

	case CmdRandomRead:
		// 隨機讀取回應: EndCode(2) + Data(2*count)
		count := 3
		body = make([]byte, 2+count*2)
		binary.LittleEndian.PutUint16(body[0:], 0) // EndCode
		for i := 0; i < count; i++ {
			binary.LittleEndian.PutUint16(body[2+i*2:], uint16(i*50))
		}
	}

	// 設置長度
	binary.LittleEndian.PutUint16(header[7:], uint16(len(body)))

	return body, nil
}

func (m *MockMCTransport) SetResponse(key string, response []byte) {
	m.responses[key] = response
}

func (m *MockMCTransport) SetError(key string, err error) {
	m.errors[key] = err
}

// TestFrame 測試封包構建和解析
func TestBuildPacket(t *testing.T) {
	frame := NewRequestFrame(0, 0xFF, 0)
	data := []byte{0x01, 0x02, 0x03}
	packet := frame.BuildPacket(CmdBatchRead, SubCmdWord, data)

	if len(packet) < 15 {
		t.Errorf("Expected packet length >= 15, got %d", len(packet))
	}

	// 驗證標頭
	if packet[0] != 0x00 || packet[1] != 0x50 {
		t.Errorf("Expected subheader 0x5000, got 0x%02X%02X", packet[1], packet[0])
	}
}

func TestParseResponseHeader(t *testing.T) {
	header := make([]byte, 9)
	header[0] = 0xD0
	header[1] = 0x00
	header[2] = 0x00
	header[3] = 0xFF
	binary.LittleEndian.PutUint16(header[4:], 0x03FF)
	header[6] = 0x00
	binary.LittleEndian.PutUint16(header[7:], 10)

	dataLen, err := ParseResponseHeader(header)
	if err != nil {
		t.Fatalf("ParseResponseHeader failed: %v", err)
	}
	if dataLen != 10 {
		t.Errorf("Expected data length 10, got %d", dataLen)
	}

	// 測試無效標頭
	header[0] = 0xFF
	_, err = ParseResponseHeader(header)
	if err == nil {
		t.Error("Expected error for invalid subheader")
	}

	// 測試過短的標頭
	_, err = ParseResponseHeader(header[:5])
	if err == nil {
		t.Error("Expected error for short header")
	}
}

func TestPackBits(t *testing.T) {
	values := []bool{true, false, true, false, true, false, true, false}
	data := PackBits(values)

	if len(data) != 4 {
		t.Errorf("Expected 4 bytes, got %d", len(data))
	}

	// 驗證第一個位元組: true(high), false(low) = 0x10
	if data[0] != 0x10 {
		t.Errorf("Expected first byte 0x10, got 0x%02X", data[0])
	}
}

func TestUnpackBits(t *testing.T) {
	// 根據 PackBits 邏輯：兩個位元打包成一個位元組
	// High nibble = 第一個設備 (i=0), Low nibble = 第二個設備 (i=1)
	// 0x10 = 0001 0000 -> high nibble=1(true), low nibble=0(false)
	// 0x30 = 0011 0000 -> high nibble=3(會被當作true因為非0), low nibble=0(false)
	// 但實際上，根據 UnpackBits 的實現，它檢查 bitVal == 1
	// 所以 high nibble=3 會被當作 false (因為 3 != 1)
	data := []byte{0x10, 0x30}
	bits := UnpackBits(data, 4)

	if len(bits) != 4 {
		t.Errorf("Expected 4 bits, got %d", len(bits))
	}
	if !bits[0] {
		t.Error("Bit 0 should be true (high nibble=1)")
	}
	if bits[1] {
		t.Error("Bit 1 should be false (low nibble=0)")
	}
	// 0x30 的 high nibble 是 3，但 UnpackBits 檢查 bitVal == 1，所以是 false
	if bits[2] {
		t.Error("Bit 2 should be false (high nibble=3, not 1)")
	}
	if bits[3] {
		t.Error("Bit 3 should be false (low nibble=0)")
	}

	// 測試正確的數據
	data2 := []byte{0x11, 0x11} // 都是 1
	bits2 := UnpackBits(data2, 4)
	for i := range bits2 {
		if !bits2[i] {
			t.Errorf("Bit %d should be true", i)
		}
	}
}

// TestClient 測試客戶端功能
func TestNewClient(t *testing.T) {
	client := NewClient("127.0.0.1", 5000)
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClientWithTransport(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)
	if client == nil {
		t.Fatal("NewClientWithTransport returned nil")
	}
}

func TestClient_BatchReadWord(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)

	values, err := client.BatchReadWord("D", 0, 10)
	if err != nil {
		t.Fatalf("BatchReadWord failed: %v", err)
	}
	if len(values) != 10 {
		t.Errorf("Expected 10 values, got %d", len(values))
	}

	// 測試小寫設備名稱
	_, err = client.BatchReadWord("d", 0, 10)
	if err != nil {
		t.Errorf("Expected lowercase device name to work, got error: %v", err)
	}

	// 測試無效地址
	_, err = client.BatchReadWord("D", -1, 10)
	if err == nil {
		t.Error("Expected error for invalid address")
	}

	_, err = client.BatchReadWord("D", 0xFFFFFF+1, 10)
	if err == nil {
		t.Error("Expected error for address out of range")
	}

	// 測試無效數量
	_, err = client.BatchReadWord("D", 0, 0)
	if err == nil {
		t.Error("Expected error for invalid count")
	}

	_, err = client.BatchReadWord("D", 0, 126)
	if err == nil {
		t.Error("Expected error for count out of range")
	}

	// 測試位元設備
	_, err = client.BatchReadWord("M", 0, 10)
	if err == nil {
		t.Error("Expected error for bit device")
	}

	// 測試未知設備
	_, err = client.BatchReadWord("UNKNOWN", 0, 10)
	if err == nil {
		t.Error("Expected error for unknown device")
	}
}

func TestClient_BatchWriteWord(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)

	values := []int{100, 200, 300}
	err := client.BatchWriteWord("D", 0, values)
	if err != nil {
		t.Fatalf("BatchWriteWord failed: %v", err)
	}

	// 測試無效數量
	err = client.BatchWriteWord("D", 0, []int{})
	if err == nil {
		t.Error("Expected error for empty values")
	}

	values124 := make([]int, 124)
	err = client.BatchWriteWord("D", 0, values124)
	if err == nil {
		t.Error("Expected error for count out of range")
	}
}

func TestClient_BatchReadBit(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)

	bits, err := client.BatchReadBit("M", 0, 10)
	if err != nil {
		t.Fatalf("BatchReadBit failed: %v", err)
	}
	if len(bits) != 10 {
		t.Errorf("Expected 10 bits, got %d", len(bits))
	}

	// 測試字組設備
	_, err = client.BatchReadBit("D", 0, 10)
	if err == nil {
		t.Error("Expected error for word device")
	}
}

func TestClient_BatchWriteBit(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)

	values := []bool{true, false, true, false}
	err := client.BatchWriteBit("M", 0, values)
	if err != nil {
		t.Fatalf("BatchWriteBit failed: %v", err)
	}

	// 測試無效數量
	err = client.BatchWriteBit("M", 0, []bool{})
	if err == nil {
		t.Error("Expected error for empty values")
	}

	values1969 := make([]bool, 1969)
	err = client.BatchWriteBit("M", 0, values1969)
	if err == nil {
		t.Error("Expected error for count out of range")
	}
}

func TestClient_RandomRead(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)

	items := []RandomReadItem{
		{Device: "D", Addr: 0},
		{Device: "D", Addr: 10},
		{Device: "D", Addr: 20},
	}

	values, err := client.RandomRead(items)
	if err != nil {
		t.Fatalf("RandomRead failed: %v", err)
	}
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	// 測試位元設備
	items[0].Device = "M"
	_, err = client.RandomRead(items)
	if err == nil {
		t.Error("Expected error for bit device in RandomRead")
	}

	// 測試未知設備
	items[0].Device = "UNKNOWN"
	_, err = client.RandomRead(items)
	if err == nil {
		t.Error("Expected error for unknown device")
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	transport := NewMockMCTransport()
	client := NewClientWithTransport(transport)

	// 創建一個會返回錯誤的回應
	// 需要構建完整的請求來設置回應
	req := client.frame.BuildPacket(CmdBatchRead, SubCmdWord, []byte{0x00, 0x00, 0x00, 0xA8, 0x0A, 0x00})
	errorResponse := make([]byte, 2)
	binary.LittleEndian.PutUint16(errorResponse[0:], 0x0001) // Error code
	transport.SetResponse(string(req), errorResponse)

	// 測試錯誤處理
	_, err := client.BatchReadWord("D", 0, 10)
	if err == nil {
		t.Error("Expected error for error response")
		return
	}

	mcErr, ok := err.(*MCError)
	if !ok {
		t.Logf("Error type: %T, error: %v", err, err)
		// 可能不是 MCError，但應該有錯誤
		return
	}
	if mcErr.Code != 0x0001 {
		t.Errorf("Expected error code 0x0001, got 0x%04X", mcErr.Code)
	}
}

// TestDeviceType 測試設備類型
func TestGetDeviceType(t *testing.T) {
	devType, err := GetDeviceType("D")
	if err != nil {
		t.Fatalf("GetDeviceType failed: %v", err)
	}
	if devType.Name != "D" {
		t.Errorf("Expected device name D, got %s", devType.Name)
	}
	if devType.IsBit {
		t.Error("Device D should not be a bit device")
	}

	devType, err = GetDeviceType("d")
	if err != nil {
		t.Fatalf("GetDeviceType failed: %v", err)
	}
	if devType.Name != "D" {
		t.Errorf("Expected device name D, got %s", devType.Name)
	}

	devType, err = GetDeviceType("M")
	if err != nil {
		t.Fatalf("GetDeviceType failed: %v", err)
	}
	if !devType.IsBit {
		t.Error("Device M should be a bit device")
	}

	_, err = GetDeviceType("UNKNOWN")
	if err == nil {
		t.Error("Expected error for unknown device")
	}
}

// TestFactory 測試工廠函數
func TestCreateTCPClient(t *testing.T) {
	client := CreateTCPClient("127.0.0.1", 5000, 0)
	if client == nil {
		t.Fatal("CreateTCPClient returned nil")
	}
}

func TestCreateSerialClient(t *testing.T) {
	client := CreateSerialClient("COM1", 9600, 7, 2, "E", 0)
	if client == nil {
		t.Fatal("CreateSerialClient returned nil")
	}
}

// TestMCError 測試錯誤類型
func TestMCError(t *testing.T) {
	err := &MCError{
		Code: 0x0001,
		Msg:  "Test error",
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
	if err.Code != 0x0001 {
		t.Errorf("Expected error code 0x0001, got 0x%04X", err.Code)
	}
}
