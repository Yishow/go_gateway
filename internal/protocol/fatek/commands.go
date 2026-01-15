package fatek

import (
	"fmt"
	"strconv"
)

// ReadStatus 讀取連續單點狀態 (Cmd 44)
//
// 用於 X, Y, M, S, T(status), C(status)
//
// Args:
//   - symbol: 元件類型 (例如 'X', 'Y')
//   - startAddr: 起始位址
//   - count: 要讀取的點數 (最大 255)
//
// Returns:
//   - 狀態列表 (True 為 ON, False 為 OFF)
func (c *Client) ReadStatus(symbol string, startAddr, count int) ([]bool, error) {
	if count > MaxDiscreteCount {
		return nil, fmt.Errorf("max count is %d", MaxDiscreteCount)
	}

	addrStr, err := NormalizeAddress(symbol, startAddr)
	if err != nil {
		return nil, err
	}

	countHex := fmt.Sprintf("%02X", count)
	body := countHex + addrStr

	dataStr, err := c.sendCommand("44", body)
	if err != nil {
		return nil, err
	}

	// 資料是 '0' 和 '1' 的字串
	results := make([]bool, len(dataStr))
	for i, char := range dataStr {
		results[i] = (char == '1')
	}

	return results, nil
}

// WriteStatus 寫入連續單點狀態 (Cmd 45)
//
// Args:
//   - symbol: 元件類型
//   - startAddr: 起始位址
//   - data: 要寫入的布林值列表
func (c *Client) WriteStatus(symbol string, startAddr int, data []bool) error {
	count := len(data)
	if count > MaxDiscreteCount {
		return fmt.Errorf("max count is %d", MaxDiscreteCount)
	}

	addrStr, err := NormalizeAddress(symbol, startAddr)
	if err != nil {
		return err
	}

	countHex := fmt.Sprintf("%02X", count)
	dataStr := ""
	for _, b := range data {
		if b {
			dataStr += "1"
		} else {
			dataStr += "0"
		}
	}

	body := countHex + addrStr + dataStr
	_, err = c.sendCommand("45", body)
	return err
}

// ReadRegisters 讀取連續 16-bit 暫存器 (Cmd 46)
//
// 用於 R, D, RT, RC
//
// Args:
//   - symbol: 元件類型 (例如 'D', 'R')
//   - startAddr: 起始位址
//   - count: 要讀取的暫存器數量 (最大 64)
//
// Returns:
//   - 暫存器值列表 (無符號 16-bit 整數)
func (c *Client) ReadRegisters(symbol string, startAddr, count int) ([]int, error) {
	if count > MaxRegisterCount {
		return nil, fmt.Errorf("max count is %d", MaxRegisterCount)
	}

	addrStr, err := NormalizeAddress(symbol, startAddr)
	if err != nil {
		return nil, err
	}

	countHex := fmt.Sprintf("%02X", count)
	body := countHex + addrStr

	dataStr, err := c.sendCommand("46", body)
	if err != nil {
		return nil, err
	}

	// 資料是每個 word 4 個 hex 字元
	results := make([]int, 0, count)
	for i := 0; i < len(dataStr); i += 4 {
		if i+4 > len(dataStr) {
			break
		}
		valHex := dataStr[i : i+4]
		val, err := strconv.ParseInt(valHex, 16, 16)
		if err != nil {
			return nil, fmt.Errorf("failed to parse hex value %s: %v", valHex, err)
		}
		results = append(results, int(val))
	}

	return results, nil
}

// WriteRegisters 寫入連續 16-bit 暫存器 (Cmd 47)
//
// Args:
//   - symbol: 元件類型
//   - startAddr: 起始位址
//   - data: 要寫入的整數值列表
func (c *Client) WriteRegisters(symbol string, startAddr int, data []int) error {
	count := len(data)
	if count > MaxRegisterCount {
		return fmt.Errorf("max count is %d", MaxRegisterCount)
	}

	addrStr, err := NormalizeAddress(symbol, startAddr)
	if err != nil {
		return err
	}

	countHex := fmt.Sprintf("%02X", count)
	dataBody := ""
	for _, val := range data {
		// 確保 16-bit 限制
		val = val & 0xFFFF
		dataBody += fmt.Sprintf("%04X", val)
	}

	body := countHex + addrStr + dataBody
	_, err = c.sendCommand("47", body)
	return err
}

// ReadRandom 混合/隨機讀取 (Cmd 48)
//
// 允許一次讀取多個任意位址 (單點或暫存器)
//
// Args:
//   - items: (symbol, address) 元組列表
//            例如 [('X', 0), ('R', 100), ('M', 20)]
//
// Returns:
//   - 映射 'SymbolAddress' 到值的字典
//     例如 {'X0': true, 'R100': 1234}
func (c *Client) ReadRandom(items []struct {
	Symbol  string
	Address int
}) (map[string]interface{}, error) {
	if len(items) > 64 { // 協定限制約為 64
		return nil, fmt.Errorf("too many items (max 64)")
	}

	countHex := fmt.Sprintf("%02X", len(items))
	body := countHex

	for _, item := range items {
		addrStr, err := NormalizeAddress(item.Symbol, item.Address)
		if err != nil {
			return nil, err
		}
		body += addrStr
	}

	rawData, err := c.sendCommand("48", body)
	if err != nil {
		return nil, err
	}

	// 解析混合回應
	results := make(map[string]interface{})
	ptr := 0
	for _, item := range items {
		key := fmt.Sprintf("%s%d", item.Symbol, item.Address)
		symbol := item.Symbol

		// 根據 symbol 判斷是單點還是暫存器
		if symbol == "X" || symbol == "Y" || symbol == "M" || symbol == "S" || symbol == "T" || symbol == "C" {
			// 單點：1 字元
			if ptr >= len(rawData) {
				return nil, fmt.Errorf("unexpected end of response data")
			}
			valStr := rawData[ptr : ptr+1]
			results[key] = (valStr == "1")
			ptr++
		} else {
			// 暫存器：4 字元
			if ptr+4 > len(rawData) {
				return nil, fmt.Errorf("unexpected end of response data")
			}
			valStr := rawData[ptr : ptr+4]
			val, err := strconv.ParseInt(valStr, 16, 16)
			if err != nil {
				return nil, fmt.Errorf("failed to parse hex value %s: %v", valStr, err)
			}
			results[key] = int(val)
			ptr += 4
		}
	}

	return results, nil
}

// LoopbackTest 執行 Loopback 測試 (Cmd 4E)
//
// 發送資料到 PLC 並檢查是否完全回傳
//
// Args:
//   - data: 要發送的字串資料 (通常是 Hex 字串)
//
// Returns:
//   - True 如果回傳匹配，否則 False
func (c *Client) LoopbackTest(data string) (bool, error) {
	req := c.buildFrame("4E", data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return false, err
	}

	// Cmd 4E 回應結構特殊：它只是回傳
	// 沒有狀態碼 '0' 欄位
	// 回應：STX + Station + 4E + Data + LRC + ETX
	if len(resp) < 6 {
		return false, nil
	}
	if resp[0] != STX || resp[len(resp)-1] != ETX {
		return false, nil
	}

	// 提取內容
	receivedCmd := string(resp[3:5])
	if receivedCmd != "4E" {
		return false, nil
	}

	receivedData := string(resp[5 : len(resp)-3])
	return receivedData == data, nil
}

// Run 啟動 PLC (RUN 模式) 使用 Cmd 41
func (c *Client) Run() error {
	return c.controlRunStop(true)
}

// Stop 停止 PLC (STOP 模式) 使用 Cmd 41
func (c *Client) Stop() error {
	return c.controlRunStop(false)
}

// controlRunStop RUN/STOP 命令的輔助方法
func (c *Client) controlRunStop(run bool) error {
	ctrlCode := "1"
	if !run {
		ctrlCode = "0"
	}
	_, err := c.sendCommand("41", ctrlCode)
	return err
}

// SingleAction 執行單一離散控制動作 (Cmd 42)
//
// Args:
//   - symbol: 元件類型 (例如 'Y', 'M')
//   - addr: 元件位址
//   - action: 要執行的操作
//             選項：'DISABLE', 'ENABLE', 'SET', 'RESET'
//
// Raises:
//   - error: 如果 action 無效
func (c *Client) SingleAction(symbol string, addr int, action string) error {
	actionMap := map[string]string{
		"DISABLE": "1",
		"ENABLE":  "2",
		"SET":     "3",
		"RESET":   "4",
	}

	code, ok := actionMap[action]
	if !ok {
		return fmt.Errorf("invalid action. use: DISABLE, ENABLE, SET, RESET")
	}

	addrStr, err := NormalizeAddress(symbol, addr)
	if err != nil {
		return err
	}

	body := code + addrStr
	_, err = c.sendCommand("42", body)
	return err
}
