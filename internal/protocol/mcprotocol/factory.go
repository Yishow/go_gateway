package mcprotocol

import "time"

// CreateTCPClient 建立 MC Protocol TCP 客戶端
//
// Args:
//   - host: PLC 的 IP 位址
//   - port: TCP 埠號
//   - timeout: Socket 逾時時間 (預設 2 秒)
//
// Returns:
//   - 配置好的 MC 客戶端實例
func CreateTCPClient(host string, port int, timeout time.Duration) *MCClient {
	transport := NewTCPTransport(host, port)
	if timeout > 0 {
		transport.Timeout = timeout
	}
	return NewClientWithTransport(transport)
}

// CreateSerialClient 建立 MC Protocol 串列埠客戶端
//
// Args:
//   - port: 串列埠名稱 (例如 "COM4" 或 "/dev/ttyUSB0")
//   - baudRate: 波特率 (預設 9600)
//   - dataBits: 資料位元數 (預設 7)
//   - stopBits: 停止位元數 (預設 2)
//   - parity: 同位檢查 ('N'=None, 'E'=Even, 'O'=Odd，預設 'E')
//   - timeout: 讀取逾時時間 (預設 2 秒)
//
// Returns:
//   - 配置好的 MC 客戶端實例
func CreateSerialClient(port string, baudRate, dataBits, stopBits int, parity string, timeout time.Duration) *MCClient {
	transport := NewSerialTransport(port, baudRate, dataBits, stopBits, parity, timeout)
	return NewClientWithTransport(transport)
}
