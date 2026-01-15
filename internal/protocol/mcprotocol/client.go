package mcprotocol

import (
	"encoding/binary"
	"fmt"
)

type MCClient struct {
	transport Transport
	frame     RequestFrame
}

// NewClient creates a new MCClient with default TCP transport
func NewClient(host string, port int) *MCClient {
	return &MCClient{
		transport: NewTCPTransport(host, port),
		frame:     NewRequestFrame(0, 0xFF, 0),
	}
}

// NewClientWithTransport 使用自訂傳輸層建立新的 MC 客戶端
//
// 此函數允許使用自訂的傳輸實作（例如測試用的 Mock Transport）。
// 對於一般使用，建議使用 NewClient 函數。
//
// Args:
//   - transport: 實作 Transport 介面的傳輸層實例
//
// Returns:
//   - 配置好的 MC 客戶端實例
//
// Example:
//   transport := NewTCPTransport("192.168.1.10", 5000)
//   client := NewClientWithTransport(transport)
func NewClientWithTransport(transport Transport) *MCClient {
	return &MCClient{
		transport: transport,
		frame:     NewRequestFrame(0, 0xFF, 0),
	}
}

func (c *MCClient) Connect() error {
	return c.transport.Connect()
}

func (c *MCClient) Close() error {
	return c.transport.Close()
}

// BatchReadWord reads 16-bit words
func (c *MCClient) BatchReadWord(device string, addr int, count int) ([]int, error) {
	// 輸入驗證
	if addr < 0 || addr > 0xFFFFFF {
		return nil, fmt.Errorf("address out of range: %d (valid: 0-16777215)", addr)
	}
	if count < 1 || count > 125 {
		return nil, fmt.Errorf("count out of range: %d (valid: 1-125 for word devices)", count)
	}

	devType, err := GetDeviceType(device)
	if err != nil {
		return nil, err
	}

	// 驗證設備類型為字組設備
	if devType.IsBit {
		return nil, fmt.Errorf("device type %s is a bit device, use BatchReadBit instead", device)
	}

	// Prepare Data: Head(3) + Code(1) + Count(2)
	data := make([]byte, 6)
	data[0] = byte(addr & 0xFF)
	data[1] = byte((addr >> 8) & 0xFF)
	data[2] = byte((addr >> 16) & 0xFF)
	data[3] = devType.Code
	binary.LittleEndian.PutUint16(data[4:], uint16(count))

	// Build & Send
	req := c.frame.BuildPacket(CmdBatchRead, SubCmdWord, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return nil, err
	}

	// Check EndCode (First 2 bytes)
	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return nil, &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for device %s at address %d (count: %d)", endCode, device, addr, count),
		}
	}

	// Data follows EndCode
	raw := resp[2:]
	if len(raw) != count*2 {
		return nil, fmt.Errorf("response length mismatch: expected %d, got %d", count*2, len(raw))
	}

	res := make([]int, count)
	for i := 0; i < count; i++ {
		res[i] = int(binary.LittleEndian.Uint16(raw[i*2:]))
	}
	return res, nil
}

// BatchWriteWord writes 16-bit words
func (c *MCClient) BatchWriteWord(device string, addr int, values []int) error {
	// 輸入驗證
	if addr < 0 || addr > 0xFFFFFF {
		return fmt.Errorf("address out of range: %d (valid: 0-16777215)", addr)
	}
	count := len(values)
	if count < 1 || count > 123 {
		return fmt.Errorf("count out of range: %d (valid: 1-123 for word write)", count)
	}

	devType, err := GetDeviceType(device)
	if err != nil {
		return err
	}

	// 驗證設備類型為字組設備
	if devType.IsBit {
		return fmt.Errorf("device type %s is a bit device, use BatchWriteBit instead", device)
	}

	// Data: Head(3) + Code(1) + Count(2) + Values(2*N)
	data := make([]byte, 6+count*2)
	data[0] = byte(addr & 0xFF)
	data[1] = byte((addr >> 8) & 0xFF)
	data[2] = byte((addr >> 16) & 0xFF)
	data[3] = devType.Code
	binary.LittleEndian.PutUint16(data[4:], uint16(count))

	for i, val := range values {
		binary.LittleEndian.PutUint16(data[6+i*2:], uint16(val))
	}

	req := c.frame.BuildPacket(CmdBatchWrite, SubCmdWord, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return err
	}

	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for device %s at address %d (count: %d)", endCode, device, addr, count),
		}
	}
	return nil
}

// BatchReadBit reads bits
func (c *MCClient) BatchReadBit(device string, addr int, count int) ([]bool, error) {
	// 輸入驗證
	if addr < 0 || addr > 0xFFFFFF {
		return nil, fmt.Errorf("address out of range: %d (valid: 0-16777215)", addr)
	}
	if count < 1 || count > 2000 {
		return nil, fmt.Errorf("count out of range: %d (valid: 1-2000 for bit devices)", count)
	}

	devType, err := GetDeviceType(device)
	if err != nil {
		return nil, err
	}

	// 驗證設備類型為位元設備（MC Protocol 3E 嚴格區分位元和字組設備）
	if !devType.IsBit {
		return nil, fmt.Errorf("device type %s is not a bit device, use BatchReadWord instead", device)
	}

	data := make([]byte, 6)
	data[0] = byte(addr & 0xFF)
	data[1] = byte((addr >> 8) & 0xFF)
	data[2] = byte((addr >> 16) & 0xFF)
	data[3] = devType.Code
	binary.LittleEndian.PutUint16(data[4:], uint16(count))

	req := c.frame.BuildPacket(CmdBatchRead, SubCmdBit, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return nil, err
	}

	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return nil, &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for device %s at address %d (count: %d)", endCode, device, addr, count),
		}
	}

	raw := resp[2:]
	expectedBytes := (count + 1) / 2
	if len(raw) != expectedBytes {
		return nil, fmt.Errorf("response length mismatch: expected %d, got %d", expectedBytes, len(raw))
	}

	return UnpackBits(raw, count), nil
}

// BatchWriteBit writes bits
func (c *MCClient) BatchWriteBit(device string, addr int, values []bool) error {
	// 輸入驗證
	if addr < 0 || addr > 0xFFFFFF {
		return fmt.Errorf("address out of range: %d (valid: 0-16777215)", addr)
	}
	count := len(values)
	if count < 1 || count > 1968 {
		return fmt.Errorf("count out of range: %d (valid: 1-1968 for bit write)", count)
	}

	devType, err := GetDeviceType(device)
	if err != nil {
		return err
	}

	// 驗證設備類型為位元設備
	if !devType.IsBit {
		return fmt.Errorf("device type %s is not a bit device, use BatchWriteWord instead", device)
	}

	bitData := PackBits(values)
	
	// Data: Head(3) + Code(1) + Count(2) + BitData...
	data := make([]byte, 6+len(bitData))
	data[0] = byte(addr & 0xFF)
	data[1] = byte((addr >> 8) & 0xFF)
	data[2] = byte((addr >> 16) & 0xFF)
	data[3] = devType.Code
	binary.LittleEndian.PutUint16(data[4:], uint16(count))
	copy(data[6:], bitData)

	req := c.frame.BuildPacket(CmdBatchWrite, SubCmdBit, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return err
	}

	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for device %s at address %d (count: %d)", endCode, device, addr, count),
		}
	}
	return nil
}

// RandomReadItem 定義隨機讀取項目
type RandomReadItem struct {
	Device string
	Addr   int
}

// RandomWriteItem 定義隨機寫入字組項目
type RandomWriteItem struct {
	Device string
	Addr   int
	Value  int // 16-bit word value
}

// RandomWriteBitItem 定義隨機寫入位元項目
type RandomWriteBitItem struct {
	Device string
	Addr   int
	Value  bool // bit value
}

// RandomRead reads mixed words
func (c *MCClient) RandomRead(items []RandomReadItem) ([]int, error) {
	count := len(items)
	// Request: Count(1) + DoubleCount(1) + [Code(1)+Head(3)]...
	// Note: DoubleCount is usually 0 for 3E frame word access command 0403
	
	data := make([]byte, 2+count*4)
	data[0] = byte(count)
	data[1] = 0 // Double word access points
	
	for i, item := range items {
		devType, err := GetDeviceType(item.Device)
		if err != nil {
			return nil, err
		}
		
		offset := 2 + i*4
		data[offset] = devType.Code
		data[offset+1] = byte(item.Addr & 0xFF)
		data[offset+2] = byte((item.Addr >> 8) & 0xFF)
		data[offset+3] = byte((item.Addr >> 16) & 0xFF)
	}

	req := c.frame.BuildPacket(CmdRandomRead, SubCmdWord, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return nil, err
	}

	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return nil, &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for random read (count: %d)", endCode, count),
		}
	}

	raw := resp[2:]
	if len(raw) != count*2 {
		return nil, fmt.Errorf("response length mismatch: expected %d bytes, got %d bytes for %d items", count*2, len(raw), count)
	}

	res := make([]int, count)
	for i := 0; i < count; i++ {
		res[i] = int(binary.LittleEndian.Uint16(raw[i*2:]))
	}
	return res, nil
}

// RandomWrite 隨機寫入多個字組和位元設備
//
// 注意：此功能需要 Command 0x1402，目前未實作。
// 參考: MC Protocol 3E Binary Frame Specification
//
// Args:
//   - wordItems: 要寫入的字組設備列表
//   - bitItems: 要寫入的位元設備列表
//
// Returns:
//   - 錯誤（目前總是返回未實作錯誤）
func (c *MCClient) RandomWrite(wordItems []RandomWriteItem, bitItems []RandomWriteBitItem) error {
	// TODO: Implement Command 0x1402 (Random Write)
	// 參考: MC Protocol 3E Binary Frame Specification
	// 需要構建複雜的請求結構，包含字組和位元設備的混合寫入
	return fmt.Errorf("random write not implemented: requires Command 0x1402")
}