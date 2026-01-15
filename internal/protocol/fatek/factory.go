package fatek

import "time"

// CreateSerialClient 建立配置為串列埠通訊的 FatekClient
//
// 注意：目前 SerialTransport 尚未完整實作，此函數為預留介面
//
// Args:
//   - port: 串列埠名稱
//   - station: PLC 站號 ID
//   - baudrate: 波特率，0 使用預設值 9600
//   - dataBits: 資料位元數，0 使用預設值 7
//   - parity: 同位檢查 ('N', 'E', 'O')，空字串使用預設值 'E'
//   - stopBits: 停止位元數，0 使用預設值 1
//   - timeout: 逾時時間，0 使用預設值 1.0 秒
//
// Returns:
//   - 配置好的客戶端實例
func CreateSerialClient(port string, station, baudrate, dataBits, stopBits int, parity string, timeout time.Duration) *FatekClient {
	// TODO: 實作完整的 SerialTransport
	// 目前返回 nil，需要整合 "go.bug.st/serial" 或類似套件
	// transport := NewSerialTransport(port, baudrate, dataBits, stopBits, parity, timeout)
	// return NewClient(transport, station)
	return nil
}

// CreateTCPClient 建立配置為 TCP 通訊的 FatekClient
//
// Args:
//   - host: PLC 的 IP 位址
//   - port: TCP 埠號 (通常為 500)，0 使用預設值 500
//   - station: PLC 站號 ID
//   - timeout: Socket 逾時時間，0 使用預設值 2.0 秒
//
// Returns:
//   - 配置好的客戶端實例
func CreateTCPClient(host string, port, station int, timeout time.Duration) *FatekClient {
	if port == 0 {
		port = DefaultTCPPort
	}
	transport := NewTCPTransport(host, port)
	if timeout > 0 {
		transport.Timeout = timeout
	}
	return NewClient(transport, station)
}
