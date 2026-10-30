package fatek

import (
	"bufio"
	"context"
	"go-gateway/internal/protocol/wirebudget"
	"net"
	"strings"
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
	mu      wirebudget.Mutex // 保護連線狀態的並發存取
}

func NewTCPTransport(host string, port int) *TCPTransport {
	return &TCPTransport{
		Host:    host,
		Port:    port,
		Timeout: 2 * time.Second,
	}
}

func (t *TCPTransport) Connect() error { return t.ConnectContext(context.Background()) }

func (t *TCPTransport) Close() error { return t.CloseContext(context.Background()) }

// internalClose 內部關閉（不加鎖，供已持鎖的方法呼叫）
func (t *TCPTransport) internalClose() error {
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return err
		}
		t.conn = nil
		t.reader = nil
		return err
	}
	return nil
}

func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	return t.SendReceiveContext(context.Background(), data)
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
	mu       wirebudget.Mutex
	closeErr error
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

	// 轉換 Parity（支援單字 N/E/O 與 none/even/odd，與前端連線設定一致）
	parityKey := strings.ToLower(strings.TrimSpace(parity))
	var p serial.Parity
	switch parityKey {
	case "n", "none":
		p = serial.NoParity
	case "e", "even":
		p = serial.EvenParity
	case "o", "odd":
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

func (s *SerialTransport) Connect() error { return s.ConnectContext(context.Background()) }

func (s *SerialTransport) Close() error { return s.CloseContext(context.Background()) }

func (s *SerialTransport) SendReceive(data []byte) ([]byte, error) {
	return s.SendReceiveContext(context.Background(), data)
}
