package fatek

import (
	"sync"
)

// Client 高階 FATEK PLC 客戶端
//
// 此類別處理 FATEK ASCII 協定命令的邏輯、
// 位址格式化，以及回應解析
type Client struct {
	transport Transport
	station   int
	mu        sync.Mutex // 用於確保「一問一答」順序執行
}

// NewClient 建立新的 FATEK 客戶端
//
// Args:
//   - transport: 傳輸實例 (Serial 或 TCP)
//   - station: PLC 站號 (0-255)，預設 1
func NewClient(transport Transport, station int) *Client {
	if station == 0 {
		station = DefaultStation
	}

	return &Client{
		transport: transport,
		station:   station,
	}
}

// Connect 連線底層傳輸
func (c *Client) Connect() error {
	return c.transport.Connect()
}

// Close 關閉底層傳輸
func (c *Client) Close() error {
	return c.transport.Close()
}

// buildFrame 構建完整的 FATEK ASCII 訊框（內部方法）
func (c *Client) buildFrame(command, body string) []byte {
	return BuildFrame(c.station, command, body)
}

// sendCommand 發送命令並解析回應（內部方法，執行緒安全）
func (c *Client) sendCommand(command, body string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req := c.buildFrame(command, body)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return "", err
	}

	return ParseResponse(resp, command)
}
