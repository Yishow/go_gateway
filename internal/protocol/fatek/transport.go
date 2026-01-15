package fatek

import (
	"fmt"
	"net"
	"time"

	"go.bug.st/serial"
)

// Transport 定義傳輸層實作的介面
//
// 用於發送和接收原始位元組
type Transport interface {
	// Connect 建立與 PLC 的連線
	Connect() error

	// Close 關閉連線
	Close() error

	// SendReceive 發送資料到 PLC 並等待回應
	//
	// Args:
	//   - data: 原始請求訊框，包含 STX、內容、LRC 和 ETX
	//
	// Returns:
	//   - 原始回應訊框
	//
	// Raises:
	//   - FatekCommunicationError: 如果發送/接收失敗或逾時
	SendReceive(data []byte) ([]byte, error)
}

// SerialTransport 使用 go.bug.st/serial 的串列埠傳輸實作
type SerialTransport struct {
	port      string
	baudrate  int
	dataBits  int
	parity    string
	stopBits  int
	timeout   time.Duration
	serialPort serial.Port
}

// NewSerialTransport 建立新的串列埠傳輸實例
//
// Args:
//   - port: COM 埠名稱 (例如 'COM1', '/dev/ttyUSB0')
//   - baudrate: 波特率，預設 9600
//   - dataBits: 資料位元數，預設 7 (FATEK ASCII 標準)
//   - parity: 同位檢查 ('N', 'E', 'O')，預設 'E' (Even)
//   - stopBits: 停止位元數，預設 1
//   - timeout: 讀取逾時時間（秒），預設 1.0
func NewSerialTransport(port string, baudrate, dataBits, stopBits int, parity string, timeout time.Duration) *SerialTransport {
	if baudrate == 0 {
		baudrate = DefaultBaudRate
	}
	if dataBits == 0 {
		dataBits = DefaultDataBits
	}
	if parity == "" {
		parity = DefaultParity
	}
	if stopBits == 0 {
		stopBits = DefaultStopBits
	}
	if timeout == 0 {
		timeout = time.Duration(DefaultTimeout) * time.Second
	}

	return &SerialTransport{
		port:     port,
		baudrate: baudrate,
		dataBits: dataBits,
		parity:   parity,
		stopBits: stopBits,
		timeout:  timeout,
	}
}

// Connect 開啟串列埠
func (s *SerialTransport) Connect() error {
	mode := &serial.Mode{
		BaudRate: s.baudrate,
		DataBits: s.dataBits,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// 設定同位檢查
	switch s.parity {
	case "E":
		mode.Parity = serial.EvenParity
	case "O":
		mode.Parity = serial.OddParity
	case "N":
		mode.Parity = serial.NoParity
	}

	// 設定停止位元
	if s.stopBits == 2 {
		mode.StopBits = serial.TwoStopBits
	}

	port, err := serial.Open(s.port, mode)
	if err != nil {
		return NewFatekCommunicationError("failed to open serial port %s: %v", s.port, err)
	}

	port.SetReadTimeout(s.timeout)
	s.serialPort = port
	return nil
}

// Close 關閉串列埠
func (s *SerialTransport) Close() error {
	if s.serialPort != nil {
		err := s.serialPort.Close()
		s.serialPort = nil
		return err
	}
	return nil
}

// SendReceive 透過串列埠發送資料並讀取直到找到 ETX
func (s *SerialTransport) SendReceive(data []byte) ([]byte, error) {
	if s.serialPort == nil {
		return nil, NewFatekCommunicationError("serial port not open")
	}

	// 清除緩衝區以確保沒有舊資料干擾
	err := s.serialPort.ResetInputBuffer()
	if err != nil {
		return nil, NewFatekCommunicationError("failed to reset input buffer: %v", err)
	}

	// 發送資料
	_, err = s.serialPort.Write(data)
	if err != nil {
		return nil, NewFatekCommunicationError("failed to write to serial port: %v", err)
	}

	// 逐位元組讀取直到 ETX (0x03)
	response := make([]byte, 0, 256)
	buf := make([]byte, 1)
	for {
		n, err := s.serialPort.Read(buf)
		if err != nil {
			return nil, NewFatekCommunicationError("timeout waiting for response: %v", err)
		}
		if n == 0 {
			return nil, NewFatekCommunicationError("timeout waiting for response")
		}

		response = append(response, buf[0])
		if buf[0] == ETX {
			break
		}
	}

	return response, nil
}

// TCPTransport TCP/IP 傳輸實作
//
// 用於具有乙太網模組或 FATEK-Ethernet 轉換器的 FATEK FBs PLC
type TCPTransport struct {
	host    string
	port    int
	timeout time.Duration
	conn    net.Conn
}

// NewTCPTransport 建立新的 TCP 傳輸實例
//
// Args:
//   - host: PLC IP 位址
//   - port: TCP 埠號，預設 500
//   - timeout: Socket 逾時時間（秒），預設 2.0
func NewTCPTransport(host string, port int, timeout time.Duration) *TCPTransport {
	if port == 0 {
		port = DefaultTCPPort
	}
	if timeout == 0 {
		timeout = time.Duration(DefaultTimeout) * time.Second
	}

	return &TCPTransport{
		host:    host,
		port:    port,
		timeout: timeout,
	}
}

// Connect 建立 TCP socket 連線
func (t *TCPTransport) Connect() error {
	address := net.JoinHostPort(t.host, fmt.Sprintf("%d", t.port))
	conn, err := net.DialTimeout("tcp", address, t.timeout)
	if err != nil {
		return NewFatekCommunicationError("failed to connect to %s: %v", address, err)
	}

	// 設定讀寫逾時
	err = conn.SetDeadline(time.Now().Add(t.timeout))
	if err != nil {
		conn.Close()
		return NewFatekCommunicationError("failed to set deadline: %v", err)
	}

	t.conn = conn
	return nil
}

// Close 關閉 TCP socket
func (t *TCPTransport) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}

// SendReceive 透過 TCP 發送資料並讀取直到找到 ETX
func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	if t.conn == nil {
		return nil, NewFatekCommunicationError("socket not connected")
	}

	// 設定寫入逾時
	err := t.conn.SetWriteDeadline(time.Now().Add(t.timeout))
	if err != nil {
		return nil, NewFatekCommunicationError("failed to set write deadline: %v", err)
	}

	// 發送資料
	_, err = t.conn.Write(data)
	if err != nil {
		return nil, NewFatekCommunicationError("failed to write to socket: %v", err)
	}

	// 設定讀取逾時
	err = t.conn.SetReadDeadline(time.Now().Add(t.timeout))
	if err != nil {
		return nil, NewFatekCommunicationError("failed to set read deadline: %v", err)
	}

	// 讀取直到 ETX
	response := make([]byte, 0, 256)
	buf := make([]byte, 1)
	for {
		n, err := t.conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, NewFatekCommunicationError("socket timeout")
			}
			if err.Error() == "EOF" {
				return nil, NewFatekCommunicationError("connection closed by peer")
			}
			return nil, NewFatekCommunicationError("socket error: %v", err)
		}
		if n == 0 {
			return nil, NewFatekCommunicationError("connection closed by peer")
		}

		response = append(response, buf[0])
		if buf[0] == ETX {
			break
		}
	}

	return response, nil
}
