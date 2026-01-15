package mcprotocol

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"go.bug.st/serial"
)

// Transport defines the interface for MC Protocol communication
type Transport interface {
	Connect() error
	Close() error
	SendReceive(req []byte) ([]byte, error)
}

type TCPTransport struct {
	Host        string
	Port        int
	Timeout     time.Duration
	MaxBodySize int // Safety limit for response body allocation
	conn        net.Conn
	mu          sync.Mutex
}

func NewTCPTransport(host string, port int) *TCPTransport {
	return &TCPTransport{
		Host:        host,
		Port:        port,
		Timeout:     2 * time.Second,
		MaxBodySize: 32 * 1024, // 32KB Limit
	}
}

func (t *TCPTransport) Connect() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if t.conn != nil {
		t.conn.Close()
	}
	
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *TCPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}

// SendReceive sends a packet and receives the response
func (t *TCPTransport) SendReceive(req []byte) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		// Auto-reconnect attempt
		if err := t.internalConnect(); err != nil {
			return nil, err
		}
	}

	t.conn.SetDeadline(time.Now().Add(t.Timeout))

	// Write
	if _, err := t.conn.Write(req); err != nil {
		t.internalClose()
		return nil, err
	}

	// Read Header (9 bytes)
	// Sub(2)+Net(1)+PC(1)+IO(2)+Station(1)+Len(2)
	// Use a small fixed buffer for header to avoid allocation
	var header [9]byte
	if _, err := io.ReadFull(t.conn, header[:]); err != nil {
		t.internalClose()
		return nil, err
	}

	// Parse Length
	dataLen, err := ParseResponseHeader(header[:])
	if err != nil {
		// Invalid header, maybe out of sync
		t.internalClose()
		return nil, err
	}

	// Safety check for body size
	if dataLen > t.MaxBodySize {
		t.internalClose()
		return nil, fmt.Errorf("response body too large: %d bytes (limit: %d)", dataLen, t.MaxBodySize)
	}

	// Read Body (EndCode + Data)
	body := make([]byte, dataLen)
	if _, err := io.ReadFull(t.conn, body); err != nil {
		t.internalClose()
		return nil, err
	}

	return body, nil
}

// internalConnect (no lock)
func (t *TCPTransport) internalConnect() error {
	if t.conn != nil {
		t.conn.Close()
	}
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

// internalClose (no lock)
func (t *TCPTransport) internalClose() {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
}

// SerialTransport MC Protocol 串列埠傳輸實作
type SerialTransport struct {
	Port     string
	BaudRate int
	DataBits int
	Parity   serial.Parity
	StopBits serial.StopBits
	Timeout  time.Duration
	port     serial.Port
	reader   *bufio.Reader
	mu       sync.Mutex
}

// NewSerialTransport 建立新的串列埠傳輸實例
func NewSerialTransport(port string, baudRate, dataBits, stopBits int, parity string, timeout time.Duration) *SerialTransport {
	if baudRate == 0 {
		baudRate = 9600
	}
	if dataBits == 0 {
		dataBits = 7
	}
	if timeout == 0 {
		timeout = 2 * time.Second
	}

	var p serial.Parity
	switch parity {
	case "N", "n":
		p = serial.NoParity
	case "E", "e":
		p = serial.EvenParity
	case "O", "o":
		p = serial.OddParity
	default:
		p = serial.EvenParity
	}

	var sb serial.StopBits
	switch stopBits {
	case 1:
		sb = serial.OneStopBit
	case 2:
		sb = serial.TwoStopBits
	default:
		sb = serial.TwoStopBits
	}

	return &SerialTransport{
		Port:     port,
		BaudRate: baudRate,
		DataBits: dataBits,
		Parity:   p,
		StopBits: sb,
		Timeout:  timeout,
	}
}

func (s *SerialTransport) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.port != nil {
		s.port.Close()
	}

	mode := &serial.Mode{
		BaudRate: s.BaudRate,
		DataBits: s.DataBits,
		Parity:   s.Parity,
		StopBits: s.StopBits,
	}

	port, err := serial.Open(s.Port, mode)
	if err != nil {
		return fmt.Errorf("串列埠開啟失敗: %w", err)
	}

	s.port = port
	s.reader = bufio.NewReader(port)
	return nil
}

func (s *SerialTransport) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.port != nil {
		err := s.port.Close()
		s.port = nil
		s.reader = nil
		return err
	}
	return nil
}

func (s *SerialTransport) SendReceive(req []byte) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.port == nil {
		return nil, fmt.Errorf("連線已關閉")
	}

	s.port.SetReadTimeout(s.Timeout)

	// 寫入請求
	if _, err := s.port.Write(req); err != nil {
		s.internalClose()
		return nil, fmt.Errorf("串列埠寫入失敗: %w", err)
	}

	// 等待一段時間（根據波特率計算）
	charTime := time.Duration(10000000/s.BaudRate) * time.Microsecond
	time.Sleep(charTime * 35 / 10) // 3.5 字符時間

	// 讀取回應標頭 (9 bytes)
	header := make([]byte, 9)
	if _, err := io.ReadFull(s.reader, header); err != nil {
		s.internalClose()
		return nil, fmt.Errorf("串列埠讀取標頭失敗: %w", err)
	}

	// 解析長度
	dataLen, err := ParseResponseHeader(header)
	if err != nil {
		s.internalClose()
		return nil, err
	}

	// 安全檢查
	if dataLen > 32*1024 {
		s.internalClose()
		return nil, fmt.Errorf("回應資料過大: %d bytes (限制: %d)", dataLen, 32*1024)
	}

	// 讀取資料本體 (EndCode + Data)
	body := make([]byte, dataLen)
	if _, err := io.ReadFull(s.reader, body); err != nil {
		s.internalClose()
		return nil, fmt.Errorf("串列埠讀取資料失敗: %w", err)
	}

	return body, nil
}

func (s *SerialTransport) internalClose() {
	if s.port != nil {
		s.port.Close()
		s.port = nil
		s.reader = nil
	}
}