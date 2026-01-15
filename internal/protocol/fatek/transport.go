package fatek

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"go.bug.st/serial"
)

type Transport interface {
	Connect() error
	Close() error
	SendReceive(data []byte) ([]byte, error)
}

// TCPTransport implementation
type TCPTransport struct {
	Host    string
	Port    int
	Timeout time.Duration
	conn    net.Conn
	reader  *bufio.Reader
	mu      sync.Mutex // 保護連線狀態的並發存取
}

func NewTCPTransport(host string, port int) *TCPTransport {
	return &TCPTransport{
		Host:    host,
		Port:    port,
		Timeout: 2 * time.Second,
	}
}

func (t *TCPTransport) Connect() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 若已有連線則先關閉
	if t.conn != nil {
		t.conn.Close()
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr := net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	t.reader = bufio.NewReader(conn)
	return nil
}

func (t *TCPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.internalClose()
}

// internalClose 內部關閉（不加鎖，供已持鎖的方法呼叫）
func (t *TCPTransport) internalClose() error {
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		t.reader = nil
		return err
	}
	return nil
}

func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return nil, ErrConnectionClosed
	}

	// 設定讀寫截止時間
	if err := t.conn.SetDeadline(time.Now().Add(t.Timeout)); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("set deadline failed: %w", err)
	}

	// 寫入請求
	// 注意：若上次操作逾時，緩衝區可能殘留資料
	// 建議：錯誤時關閉連線，由呼叫方偵測 ErrConnectionClosed 並重連
	if _, err := t.conn.Write(data); err != nil {
		t.internalClose() // 錯誤時強制關閉以重置狀態
		return nil, err
	}

	// 讀取直到 ETX
	response, err := t.reader.ReadBytes(ETX)
	if err != nil {
		t.internalClose() // 錯誤時（包括逾時）強制關閉以重置狀態
		return nil, err
	}

	return response, nil
}

// SerialTransport 實作串列埠通訊
type SerialTransport struct {
	Port     string
	BaudRate int
	DataBits int
	Parity   serial.Parity
	StopBits serial.StopBits
	Timeout  time.Duration
	port     serial.Port
	reader   *bufio.Reader
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
		timeout = 1 * time.Second
	}

	// 轉換 Parity
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

	// 轉換 StopBits
	var sb serial.StopBits
	switch stopBits {
	case 1:
		sb = serial.OneStopBit
	case 2:
		sb = serial.TwoStopBits
	default:
		sb = serial.OneStopBit
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
	mode := &serial.Mode{
		BaudRate: s.BaudRate,
		DataBits: s.DataBits,
		Parity:   s.Parity,
		StopBits: s.StopBits,
	}

	port, err := serial.Open(s.Port, mode)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %w", s.Port, err)
	}

	s.port = port
	s.reader = bufio.NewReader(port)
	return nil
}

func (s *SerialTransport) Close() error {
	if s.port != nil {
		err := s.port.Close()
		s.port = nil
		s.reader = nil
		return err
	}
	return nil
}

func (s *SerialTransport) SendReceive(data []byte) ([]byte, error) {
	if s.port == nil {
		return nil, ErrConnectionClosed
	}

	// 設定讀取逾時
	s.port.SetReadTimeout(s.Timeout)

	// 寫入資料
	_, err := s.port.Write(data)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("serial write error: %w", err)
	}

	// 讀取直到 ETX
	response, err := s.reader.ReadBytes(ETX)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("serial read error: %w", err)
	}

	return response, nil
}
