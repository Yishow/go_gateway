package fatek

import (
	"fmt"
	"strings"
	"sync"
)

type FatekClient struct {
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

// execute performs the Frame Build -> Send -> Receive -> Parse cycle
func (c *FatekClient) execute(cmd, body string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := BuildFrame(c.station, cmd, body)
	
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return "", err
	}

	return ParseResponse(resp, cmd)
}

// ReadStatus (Cmd 44): Continuous Read of Discrete Status
func (c *FatekClient) ReadStatus(symbol string, startAddr int, count int) ([]bool, error) {
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
	addrStr := FormatAddress(comp, startAddr)
	
	body := countHex + addrStr
	
	dataStr, err := c.execute("44", body)
	if err != nil {
		return nil, err
	}

	result := make([]bool, len(dataStr))
	for i, char := range dataStr {
		result[i] = (char == '1')
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
	if !comp.IsDiscrete {
		return fmt.Errorf("component %s is not discrete", symbol)
	}

	countHex := IntToHex(count, 2)
	addrStr := FormatAddress(comp, startAddr)
	
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
func (c *FatekClient) ReadRegisters(symbol string, startAddr int, count int) ([]int, error) {
	comp, err := GetComponentType(symbol)
	if err != nil {
		return nil, err
	}

	// Limit check: 16-bit max 64, 32-bit max 32 (roughly)
	maxCount := 64
	if comp.Width == 32 {
		maxCount = 32
	}
	if count > maxCount {
		return nil, fmt.Errorf("max count for %d-bit register is %d", comp.Width, maxCount)
	}

	countHex := IntToHex(count, 2)
	addrStr := FormatAddress(comp, startAddr)
	
	body := countHex + addrStr
	dataStr, err := c.execute("46", body)
	if err != nil {
		return nil, err
	}

	// 16-bit = 4 chars, 32-bit = 8 chars
	charsPerVal := comp.Width / 4
	expectedLen := count * charsPerVal
	if len(dataStr) != expectedLen {
		return nil, fmt.Errorf("response length mismatch: expected %d chars, got %d", expectedLen, len(dataStr))
	}

	result := make([]int, 0, count)
	for i := 0; i < len(dataStr); i += charsPerVal {
		if i+charsPerVal > len(dataStr) {
			break
		}
		valHex := dataStr[i : i+charsPerVal]
		val, err := HexToInt(valHex)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	
	// 驗證結果數量是否與請求一致
	if len(result) != count {
		return nil, fmt.Errorf("incomplete response: expected %d values, got %d", count, len(result))
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
	addrStr := FormatAddress(comp, startAddr)
	
	var sb strings.Builder
	charsPerVal := comp.Width / 4
	
	for _, val := range data {
		// Masking to correct bit width
		mask := (1 << comp.Width) - 1
		maskedVal := val & mask
		sb.WriteString(IntToHex(maskedVal, charsPerVal))
	}
	
	body := countHex + addrStr + sb.String()
	_, err = c.execute("47", body)
	return err
}

type RandomReadItem struct {
	Symbol string
	Addr   int
}

// ReadRandom (Cmd 48): Mixed Read
func (c *FatekClient) ReadRandom(items []RandomReadItem) (map[string]interface{}, error) {
	if len(items) > 64 {
		return nil, fmt.Errorf("max items is 64")
	}

	countHex := IntToHex(len(items), 2)
	var sb strings.Builder
	sb.WriteString(countHex)

	comps := make([]ComponentType, len(items))

	for i, item := range items {
		comp, err := GetComponentType(item.Symbol)
		if err != nil {
			return nil, err
		}
		comps[i] = comp
		sb.WriteString(FormatAddress(comp, item.Addr))
	}

	dataStr, err := c.execute("48", sb.String())
	if err != nil {
		return nil, err
	}

	// Parse mixed response
	results := make(map[string]interface{})
	ptr := 0
	
	for i, item := range items {
		comp := comps[i]
		key := fmt.Sprintf("%s%d", item.Symbol, item.Addr)
		
		if comp.IsDiscrete {
			// 1 char
			if ptr+1 > len(dataStr) {
				return nil, fmt.Errorf("malformed response for item %d", i)
			}
			valChar := dataStr[ptr]
			results[key] = (valChar == '1')
			ptr += 1
		} else {
			// Register (4 or 8 chars)
			chars := comp.Width / 4
			if ptr+chars > len(dataStr) {
				return nil, fmt.Errorf("malformed response for item %d", i)
			}
			valHex := dataStr[ptr : ptr+chars]
			val, err := HexToInt(valHex)
			if err != nil {
				return nil, err
			}
			results[key] = val
			ptr += chars
		}
	}

	return results, nil
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

	addrStr := FormatAddress(comp, addr)
	body := code + addrStr
	_, err = c.execute("42", body)
	return err
}