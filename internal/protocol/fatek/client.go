package fatek

import (
	"fmt"
	"strings"
	"sync"
)

type FatekClient struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	transport Transport
	station   int
	mu        sync.Mutex
}

func NewClient(transport Transport, station int) *FatekClient {
	if station == 0 {
		station = DefaultStation
	}
	return &FatekClient{
		transport: transport,
		station:   station,
	}
}

func (c *FatekClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transport.Connect()
}

func (c *FatekClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transport.Close()
}

// execute performs the Frame Build -> Send -> Receive -> Parse cycle using Buffer Pool
func (c *FatekClient) execute(cmd, body string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get a buffer from the pool
	buf := GetBuffer()
	defer PutBuffer(buf)

	// Build frame directly into the buffer (Zero Allocation)
	BuildFrameToBuffer(buf, c.station, cmd, body)

	// Send raw bytes from buffer
	resp, err := c.transport.SendReceive(buf.Bytes())
	if err != nil {
		return "", err
	}

	// Reject foreign-station status before it can be interpreted as a capability
	// rejection. Short frames still use the existing framing error contract.
	if len(resp) >= 3 && string(resp[1:3]) != fmt.Sprintf("%02X", c.station) {
		return "", fmt.Errorf("response station mismatch: expected %02X, got %s", c.station, resp[1:3])
	}
	return ParseResponse(resp, cmd)
}

// ReadStatus (Cmd 44): Continuous Read of Discrete Status
func (c *FatekClient) ReadStatus(symbol string, startAddr, count int) ([]bool, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be greater than 0")
	}
	if count > 255 {
		return nil, fmt.Errorf("max count is 255")
	}

	comp, err := GetComponentType(symbol)
	if err != nil {
		return nil, err
	}
	if !comp.IsDiscrete {
		return nil, fmt.Errorf("component %s is not discrete", symbol)
	}

	countHex := IntToHex(count, 2)
	addrStr, err := FormatAddress(comp, startAddr)
	if err != nil {
		return nil, err
	}

	body := countHex + addrStr

	dataStr, err := c.execute("44", body)
	if err != nil {
		return nil, err
	}
	if len(dataStr) != count {
		return nil, fmt.Errorf("response length mismatch for %s: expected %d status bytes, got %d", symbol, count, len(dataStr))
	}

	result := make([]bool, count)
	for i := range len(dataStr) {
		switch dataStr[i] {
		case '0':
			result[i] = false
		case '1':
			result[i] = true
		default:
			return nil, fmt.Errorf("invalid status payload byte at index %d", i)
		}
	}
	return result, nil
}

// WriteStatus (Cmd 45): Continuous Write of Discrete Status
func (c *FatekClient) WriteStatus(symbol string, startAddr int, data []bool) error {
	count := len(data)
	if count > 255 {
		return fmt.Errorf("max count is 255")
	}

	comp, err := GetComponentType(symbol)
	if err != nil {
		return err
	}

	countHex := IntToHex(count, 2)
	addrStr, err := FormatAddress(comp, startAddr)
	if err != nil {
		return err
	}

	var sb strings.Builder
	for _, b := range data {
		if b {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}

	body := countHex + addrStr + sb.String()
	_, err = c.execute("45", body)
	return err
}

// ReadRegisters (Cmd 46): Continuous Read of Registers
func (c *FatekClient) ReadRegisters(symbol string, startAddr, count int) ([]int, error) {
	comp, err := GetComponentType(symbol)
	if err != nil {
		return nil, err
	}

	// Limit check: 16-bit max 64, 32-bit max 32 (roughly)
	maxCount := 64
	if comp.Width == 32 {
		maxCount = 32
	}
	if count < 1 {
		return nil, fmt.Errorf("count must be greater than 0")
	}
	if count > maxCount {
		return nil, fmt.Errorf("max count for %d-bit register is %d", comp.Width, maxCount)
	}

	countHex := IntToHex(count, 2)
	addrStr, err := FormatAddress(comp, startAddr)
	if err != nil {
		return nil, err
	}

	if _, err = FormatAddress(comp, startAddr+(count-1)*comp.AddressStep()); err != nil {
		return nil, err
	}

	body := countHex + addrStr
	dataStr, err := c.execute("46", body)
	if err != nil {
		return nil, err
	}

	// 16-bit = 4 chars, 32-bit = 8 chars
	charsPerVal := comp.Width / 4
	expectedLen := count * charsPerVal

	// 嚴格驗證回應長度
	if len(dataStr) != expectedLen {
		return nil, fmt.Errorf(
			"response length mismatch for %s: expected %d chars, got %d",
			symbol, expectedLen, len(dataStr),
		)
	}

	result := make([]int, count)
	for i := range count {
		valHex := dataStr[i*charsPerVal : (i+1)*charsPerVal]
		val, err := HexToInt(valHex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value at index %d: %w", i, err)
		}
		result[i] = val
	}
	return result, nil
}

// WriteRegisters (Cmd 47): Continuous Write of Registers
func (c *FatekClient) WriteRegisters(symbol string, startAddr int, data []int) error {
	count := len(data)
	comp, err := GetComponentType(symbol)
	if err != nil {
		return err
	}

	maxCount := 64
	if comp.Width == 32 {
		maxCount = 32
	}
	if count > maxCount {
		return fmt.Errorf("max count for %d-bit register is %d", comp.Width, maxCount)
	}

	countHex := IntToHex(count, 2)
	addrStr, err := FormatAddress(comp, startAddr)
	if err != nil {
		return err
	}

	var sb strings.Builder
	charsPerVal := comp.Width / 4

	for _, val := range data {
		// Use uint64 to safely handle masking without overflow on 32-bit int
		mask := (uint64(1) << uint64(comp.Width)) - 1         // #nosec G115 -- component widths are validated fixed-width FATEK register sizes.
		maskedVal := uint64(val) & mask                       // #nosec G115 -- preserve the low register bits for signed two's-complement wire values.
		sb.WriteString(IntToHex(int(maskedVal), charsPerVal)) // #nosec G115 -- the masked value is rendered as the configured low-width hexadecimal word.
	}

	body := countHex + addrStr + sb.String()
	_, err = c.execute("47", body)
	return err
}

// LoopbackTest (Cmd 4E)
func (c *FatekClient) LoopbackTest(data string) (bool, error) {
	// Execute already handles the special parsing for 4E
	respBody, err := c.execute("4E", data)
	if err != nil {
		return false, err
	}
	return respBody == data, nil
}

// Run (Cmd 41)
func (c *FatekClient) Run() error {
	_, err := c.execute("41", "1")
	return err
}

// Stop (Cmd 41)
func (c *FatekClient) Stop() error {
	_, err := c.execute("41", "0")
	return err
}

// SingleAction (Cmd 42)
func (c *FatekClient) SingleAction(symbol string, addr int, action string) error {
	comp, err := GetComponentType(symbol)
	if err != nil {
		return err
	}
	if !comp.IsDiscrete {
		return fmt.Errorf("component must be discrete")
	}

	var code string
	switch strings.ToUpper(action) {
	case "DISABLE":
		code = "1"
	case "ENABLE":
		code = "2"
	case "SET":
		code = "3"
	case "RESET":
		code = "4"
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	addrStr, err := FormatAddress(comp, addr)
	if err != nil {
		return err
	}

	body := code + addrStr
	_, err = c.execute("42", body)
	return err
}
