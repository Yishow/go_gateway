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

// NewClientWithTransport creates a new MCClient with custom transport
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
	devType, err := GetDeviceType(device)
	if err != nil {
		return nil, err
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
		return nil, &MCError{Code: int(endCode), Msg: "PLC returned error code"}
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
	count := len(values)
	devType, err := GetDeviceType(device)
	if err != nil {
		return err
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
		return &MCError{Code: int(endCode), Msg: "PLC returned error code"}
	}
	return nil
}

// BatchReadBit reads bits
func (c *MCClient) BatchReadBit(device string, addr int, count int) ([]bool, error) {
	devType, err := GetDeviceType(device)
	if err != nil {
		return nil, err
	}
	if !devType.IsBit {
		// Can read words as bits? Usually strictly separate in MC Protocol 3E
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
		return nil, &MCError{Code: int(endCode), Msg: "PLC returned error code"}
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
	count := len(values)
	devType, err := GetDeviceType(device)
	if err != nil {
		return err
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
		return &MCError{Code: int(endCode), Msg: "PLC returned error code"}
	}
	return nil
}

// RandomReadItem
type RandomReadItem struct {
	Device string
	Addr   int
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
		return nil, &MCError{Code: int(endCode), Msg: "PLC returned error code"}
	}

	raw := resp[2:]
	if len(raw) != count*2 {
		return nil, fmt.Errorf("response length mismatch")
	}

	res := make([]int, count)
	for i := 0; i < count; i++ {
		res[i] = int(binary.LittleEndian.Uint16(raw[i*2:]))
	}
	return res, nil
}

// RandomWrite (Placeholder) - usually involves complex command structure
func (c *MCClient) RandomWrite(wordItems []struct{Device string; Addr int; Value int}, bitItems []struct{Device string; Addr int; Value bool}) error {
	// TODO: Implement Command 1402 (Random Write) if needed
	return fmt.Errorf("random write not implemented")
}