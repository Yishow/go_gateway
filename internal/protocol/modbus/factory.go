package modbus

import "time"

// CreateTCPClient 建立 Modbus TCP 客戶端
//
// Args:
//   - host: PLC 的 IP 位址
//   - port: TCP 埠號 (預設 502)
//   - unitID: Modbus 單元 ID (站號，預設 1)
//   - timeout: Socket 逾時時間 (預設 2 秒)
//
// Returns:
//   - 配置好的 Modbus 客戶端實例
func CreateTCPClient(host string, port int, unitID byte, timeout time.Duration) *ModbusClient {
	transport := NewTCPTransport(host, port)
	if timeout > 0 {
		transport.Timeout = timeout
	}
	return NewClient(transport, unitID)
}

// CreateUDPClient 建立 Modbus UDP 客戶端
//
// Args:
//   - host: PLC 的 IP 位址
//   - port: UDP 埠號 (預設 502)
//   - unitID: Modbus 單元 ID (站號，預設 1)
//   - timeout: Socket 逾時時間 (預設 2 秒)
//
// Returns:
//   - 配置好的 Modbus 客戶端實例
func CreateUDPClient(host string, port int, unitID byte, timeout time.Duration) *ModbusClient {
	transport := NewUDPTransport(host, port)
	if timeout > 0 {
		transport.Timeout = timeout
	}
	return NewClient(transport, unitID)
}

// CreateRTUClient 建立 Modbus RTU 客戶端 (串列埠)
//
// Args:
//   - port: 串列埠名稱 (例如 "COM3" 或 "/dev/ttyUSB0")
//   - baudRate: 波特率 (預設 9600)
//   - dataBits: 資料位元數 (預設 8)
//   - stopBits: 停止位元數 (預設 1)
//   - parity: 同位檢查 ('N'=None, 'E'=Even, 'O'=Odd，預設 'N')
//   - timeout: 讀取逾時時間 (預設 2 秒)
//   - unitID: Modbus 單元 ID (站號，預設 1)
//
// Returns:
//   - 配置好的 Modbus 客戶端實例
func CreateRTUClient(port string, baudRate, dataBits, stopBits int, parity string, timeout time.Duration, unitID byte) *ModbusClient {
	transport := NewRTUTransport(port, baudRate, dataBits, stopBits, parity, timeout)
	return NewClient(transport, unitID)
}
