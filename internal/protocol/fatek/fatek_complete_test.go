package fatek

import (
	"errors"
	"testing"
)

// MockTransport 模擬傳輸層用於測試
type MockFatekTransport struct {
	responses map[string][]byte
	errors    map[string]error
	sentData  [][]byte
}

func NewMockFatekTransport() *MockFatekTransport {
	return &MockFatekTransport{
		responses: make(map[string][]byte),
		errors:    make(map[string]error),
		sentData:  make([][]byte, 0),
	}
}

func (m *MockFatekTransport) Connect() error {
	return nil
}

func (m *MockFatekTransport) Close() error {
	return nil
}

func (m *MockFatekTransport) SendReceive(data []byte) ([]byte, error) {
	m.sentData = append(m.sentData, data)

	// 檢查是否有預設錯誤
	if err, ok := m.errors[string(data)]; ok {
		return nil, err
	}

	// 檢查是否有預設回應
	if resp, ok := m.responses[string(data)]; ok {
		return resp, nil
	}

	// 預設回應：根據命令生成
	if len(data) < 9 {
		return nil, ErrResponseTooShort
	}

	// 解析命令
	cmd := string(data[3:5])

	var response []byte

	switch cmd {
	case "44": // ReadStatus
		// STX + Station + Cmd + Status(1) + Data(N) + LRC(2) + ETX
		response = make([]byte, 9+10) // 最小長度 + 10 個狀態
		response[0] = STX
		copy(response[1:3], data[1:3]) // Station
		copy(response[3:5], data[3:5]) // Cmd
		response[5] = '0'              // Status
		// 10 個 '0' (false)
		for i := 0; i < 10; i++ {
			response[6+i] = '0'
		}
		lrc := CalculateLRC(response[:16])
		copy(response[16:18], lrc)
		response[18] = ETX

	case "45": // WriteStatus
		// STX + Station + Cmd + Status(1) + LRC(2) + ETX
		response = make([]byte, 9)
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		response[5] = '0'
		lrc := CalculateLRC(response[:6])
		copy(response[6:8], lrc)
		response[8] = ETX

	case "46": // ReadRegisters
		// STX + Station + Cmd + Status(1) + Data(4*N) + LRC(2) + ETX
		count := 10
		response = make([]byte, 9+count*4)
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		response[5] = '0'
		// 10 個暫存器值 (每個 4 個 hex 字符)
		for i := 0; i < count; i++ {
			val := IntToHex(i*10, 4)
			copy(response[6+i*4:6+(i+1)*4], val)
		}
		lrc := CalculateLRC(response[:46])
		copy(response[46:48], lrc)
		response[48] = ETX

	case "47": // WriteRegisters
		// STX + Station + Cmd + Status(1) + LRC(2) + ETX
		response = make([]byte, 9)
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		response[5] = '0'
		lrc := CalculateLRC(response[:6])
		copy(response[6:8], lrc)
		response[8] = ETX

	case "48": // ReadRandom
		// STX + Station + Cmd + Status(1) + Data(Variable) + LRC(2) + ETX
		response = make([]byte, 9+8) // 假設 2 個項目，每個 4 字符
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		response[5] = '0'
		copy(response[6:14], "00010002") // 兩個值
		lrc := CalculateLRC(response[:14])
		copy(response[14:16], lrc)
		response[16] = ETX

	case "4E": // LoopbackTest
		// STX + Station + Cmd + Data + LRC(2) + ETX
		bodyLen := len(data) - 8 // 減去 STX + Station + Cmd + LRC + ETX
		response = make([]byte, 8+bodyLen)
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		copy(response[5:5+bodyLen], data[5:5+bodyLen])
		lrc := CalculateLRC(response[:5+bodyLen])
		copy(response[5+bodyLen:5+bodyLen+2], lrc)
		response[5+bodyLen+2] = ETX

	case "41": // Run/Stop
		// STX + Station + Cmd + Status(1) + LRC(2) + ETX
		response = make([]byte, 9)
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		response[5] = '0'
		lrc := CalculateLRC(response[:6])
		copy(response[6:8], lrc)
		response[8] = ETX

	case "42": // SingleAction
		// STX + Station + Cmd + Status(1) + LRC(2) + ETX
		response = make([]byte, 9)
		response[0] = STX
		copy(response[1:3], data[1:3])
		copy(response[3:5], data[3:5])
		response[5] = '0'
		lrc := CalculateLRC(response[:6])
		copy(response[6:8], lrc)
		response[8] = ETX

	default:
		return nil, ErrInvalidCommand
	}

	return response, nil
}

func (m *MockFatekTransport) SetResponse(key string, response []byte) {
	m.responses[key] = response
}

func (m *MockFatekTransport) SetError(key string, err error) {
	m.errors[key] = err
}

// TestFrame 測試封包構建和解析
func TestCalculateLRC(t *testing.T) {
	data := []byte{STX, '0', '1', '4', '4', 'T', 'E', 'S', 'T'}
	lrc := CalculateLRC(data)

	if len(lrc) != 2 {
		t.Errorf("Expected LRC length 2, got %d", len(lrc))
	}
}

func TestBuildFrame(t *testing.T) {
	frame := BuildFrame(1, "44", "TEST")

	if len(frame) < 9 {
		t.Error("Frame too short")
	}
	if frame[0] != STX {
		t.Error("Frame should start with STX")
	}
	if frame[len(frame)-1] != ETX {
		t.Error("Frame should end with ETX")
	}
}

func TestBuildFrameToBuffer_Complete(t *testing.T) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	BuildFrameToBuffer(buf, 1, "44", "TEST")

	if buf.Len() < 9 {
		t.Error("Frame too short")
	}

	data := buf.Bytes()
	if data[0] != STX {
		t.Error("Frame should start with STX")
	}
	if data[len(data)-1] != ETX {
		t.Error("Frame should end with ETX")
	}
}

func TestParseResponse(t *testing.T) {
	// 正常回應 - 需要正確的 LRC
	response := make([]byte, 19)
	response[0] = STX
	copy(response[1:3], []byte{'0', '1'})
	copy(response[3:5], []byte{'4', '4'})
	response[5] = '0'
	copy(response[6:16], []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'})
	lrc := CalculateLRC(response[:16])
	copy(response[16:18], lrc)
	response[18] = ETX

	body, err := ParseResponse(response, "44")
	if err != nil {
		t.Fatalf("ParseResponse failed: %v", err)
	}
	if len(body) != 10 {
		t.Errorf("Expected body length 10, got %d", len(body))
	}

	// 測試無效 STX
	response[0] = 0xFF
	_, err = ParseResponse(response, "44")
	if !errors.Is(err, ErrInvalidSTX) {
		t.Errorf("Expected ErrInvalidSTX, got %v", err)
	}

	// 測試無效 ETX
	response[0] = STX
	response[len(response)-1] = 0xFF
	_, err = ParseResponse(response, "44")
	if !errors.Is(err, ErrInvalidETX) {
		t.Errorf("Expected ErrInvalidETX, got %v", err)
	}

	// 測試過短的回應
	_, err = ParseResponse([]byte{STX, '0', '1'}, "44")
	if !errors.Is(err, ErrResponseTooShort) {
		t.Errorf("Expected ErrResponseTooShort, got %v", err)
	}

	// 測試錯誤狀態
	errorResponse := []byte{STX, '0', '1', '4', '4', '1', 'F', 'F', 0x03} // Status '1' = Error
	_, err = ParseResponse(errorResponse, "44")
	if err == nil {
		t.Error("Expected error for error status")
	}
}

// TestAddress 測試地址處理
func TestGetComponentType(t *testing.T) {
	comp, err := GetComponentType("X")
	if err != nil {
		t.Fatalf("GetComponentType failed: %v", err)
	}
	if comp.Name != "X" {
		t.Errorf("Expected component name X, got %s", comp.Name)
	}
	if !comp.IsDiscrete {
		t.Error("Component X should be discrete")
	}

	comp, err = GetComponentType("D")
	if err != nil {
		t.Fatalf("GetComponentType failed: %v", err)
	}
	if comp.IsDiscrete {
		t.Error("Component D should not be discrete")
	}

	comp, err = GetComponentType("DD")
	if err != nil {
		t.Fatalf("GetComponentType failed: %v", err)
	}
	if comp.Name != "DD" {
		t.Errorf("Expected component name DD, got %s", comp.Name)
	}
	if comp.Width != 32 {
		t.Errorf("Expected DD width 32, got %d", comp.Width)
	}

	comp, err = GetComponentType("DF")
	if err != nil {
		t.Fatalf("GetComponentType failed: %v", err)
	}
	if comp.Name != "DF" {
		t.Errorf("Expected component name DF, got %s", comp.Name)
	}

	comp, err = GetComponentType("WX")
	if err != nil {
		t.Fatalf("GetComponentType failed: %v", err)
	}
	if comp.Name != "WX" {
		t.Errorf("Expected component name WX, got %s", comp.Name)
	}
	if comp.IsDiscrete {
		t.Error("Component WX should be word access (not discrete)")
	}

	_, err = GetComponentType("UNKNOWN")
	if err == nil {
		t.Error("Expected error for unknown component")
	}
}

func TestFormatAddress(t *testing.T) {
	comp, _ := GetComponentType("D")
	addr, err := FormatAddress(comp, 123)
	if err != nil {
		t.Fatalf("FormatAddress failed: %v", err)
	}
	if addr != "D00123" {
		t.Errorf("Expected address D00123, got %s", addr)
	}

	comp, _ = GetComponentType("DD")
	addr, err = FormatAddress(comp, 123)
	if err != nil {
		t.Fatalf("FormatAddress failed: %v", err)
	}
	if addr != "DD00123" {
		t.Errorf("Expected address DD00123, got %s", addr)
	}

	comp, _ = GetComponentType("DF")
	addr, err = FormatAddress(comp, 9)
	if err != nil {
		t.Fatalf("FormatAddress failed: %v", err)
	}
	if addr != "DF00009" {
		t.Errorf("Expected address DF00009, got %s", addr)
	}

	comp, _ = GetComponentType("WX")
	addr, err = FormatAddress(comp, 12)
	if err != nil {
		t.Fatalf("FormatAddress failed: %v", err)
	}
	if addr != "WX00012" {
		t.Errorf("Expected address WX00012, got %s", addr)
	}

	// 測試超出範圍
	_, err = FormatAddress(comp, 100000)
	if err == nil {
		t.Error("Expected error for address out of range")
	}

	// 測試負數
	_, err = FormatAddress(comp, -1)
	if err == nil {
		t.Error("Expected error for negative address")
	}
}

// TestClient 測試客戶端功能
func TestNewClient(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.station != 1 {
		t.Errorf("Expected station 1, got %d", client.station)
	}

	// 測試預設站號
	client2 := NewClient(transport, 0)
	if client2.station != DefaultStation {
		t.Errorf("Expected default station %d, got %d", DefaultStation, client2.station)
	}
}

func TestClient_ReadStatus(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	status, err := client.ReadStatus("X", 0, 10)
	if err != nil {
		t.Fatalf("ReadStatus failed: %v", err)
	}
	if len(status) != 10 {
		t.Errorf("Expected 10 status values, got %d", len(status))
	}

	// 測試無效組件
	_, err = client.ReadStatus("D", 0, 10)
	if err == nil {
		t.Error("Expected error for non-discrete component")
	}

	// 測試數量過大
	_, err = client.ReadStatus("X", 0, 256)
	if err == nil {
		t.Error("Expected error for count > 255")
	}
}

func TestClient_WriteStatus(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	values := []bool{true, false, true, false}
	err := client.WriteStatus("X", 0, values)
	if err != nil {
		t.Fatalf("WriteStatus failed: %v", err)
	}
}

func TestClient_ReadRegisters(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	registers, err := client.ReadRegisters("D", 0, 10)
	if err != nil {
		t.Fatalf("ReadRegisters failed: %v", err)
	}
	if len(registers) != 10 {
		t.Errorf("Expected 10 registers, got %d", len(registers))
	}
}

func TestClient_WriteRegisters(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	values := []int{100, 200, 300}
	err := client.WriteRegisters("D", 0, values)
	if err != nil {
		t.Fatalf("WriteRegisters failed: %v", err)
	}
}

func TestClient_ReadRandom(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	items := []RandomReadItem{
		{Symbol: "D", Addr: 0},
		{Symbol: "D", Addr: 10},
	}

	results, err := client.ReadRandom(items)
	if err != nil {
		t.Fatalf("ReadRandom failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// 測試數量過大
	items65 := make([]RandomReadItem, 65)
	for i := range items65 {
		items65[i] = RandomReadItem{Symbol: "D", Addr: i}
	}
	_, err = client.ReadRandom(items65)
	if err == nil {
		t.Error("Expected error for count > 64")
	}
}

func TestClient_LoopbackTest(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	ok, err := client.LoopbackTest("TEST")
	if err != nil {
		t.Fatalf("LoopbackTest failed: %v", err)
	}
	if !ok {
		t.Error("LoopbackTest should return true")
	}
}

func TestClient_Run(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	err := client.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
}

func TestClient_Stop(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	err := client.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
}

func TestClient_SingleAction(t *testing.T) {
	transport := NewMockFatekTransport()
	client := NewClient(transport, 1)

	err := client.SingleAction("X", 0, "SET")
	if err != nil {
		t.Fatalf("SingleAction failed: %v", err)
	}

	// 測試無效動作
	err = client.SingleAction("X", 0, "INVALID")
	if err == nil {
		t.Error("Expected error for invalid action")
	}

	// 測試非離散組件
	err = client.SingleAction("D", 0, "SET")
	if err == nil {
		t.Error("Expected error for non-discrete component")
	}
}

// TestFactory 測試工廠函數
func TestCreateTCPClient(t *testing.T) {
	client := CreateTCPClient("127.0.0.1", 500, 1, 0)
	if client == nil {
		t.Fatal("CreateTCPClient returned nil")
	}
}

func TestCreateSerialClient(t *testing.T) {
	client := CreateSerialClient("COM1", 1, 9600, 7, 2, "E", 0)
	if client == nil {
		t.Fatal("CreateSerialClient returned nil")
	}
}

// TestUtils 測試工具函數
func TestIntToHex(t *testing.T) {
	hex := IntToHex(255, 4)
	if hex != "00FF" {
		t.Errorf("Expected 00FF, got %s", hex)
	}

	hex = IntToHex(10, 2)
	if hex != "0A" {
		t.Errorf("Expected 0A, got %s", hex)
	}
}

func TestHexToInt(t *testing.T) {
	val, err := HexToInt("00FF")
	if err != nil {
		t.Fatalf("HexToInt failed: %v", err)
	}
	if val != 255 {
		t.Errorf("Expected 255, got %d", val)
	}

	_, err = HexToInt("INVALID")
	if err == nil {
		t.Error("Expected error for invalid hex")
	}
}

// TestPool 測試緩衝池
func TestBufferPool(t *testing.T) {
	buf1 := GetBuffer()
	PutBuffer(buf1)

	buf2 := GetBuffer()
	PutBuffer(buf2)

	// 驗證池可以重用緩衝區
	if buf1 != buf2 {
		// 這不是錯誤，只是說明池可能沒有重用（取決於實現）
		t.Log("Buffer pool may not be reusing buffers")
	}
}

// TestErrors 測試錯誤處理
func TestErrors(t *testing.T) {
	if ErrConnectionClosed.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if ErrResponseTooShort.Error() == "" {
		t.Error("Error message should not be empty")
	}
}
