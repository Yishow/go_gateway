package modbus

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

// Transport 定義 Modbus 傳輸介面
type Transport interface {
	Connect() error
	Close() error
	SendReceive(data []byte) ([]byte, error)
}

// TCPTransport Modbus TCP 傳輸實作
type TCPTransport struct {
	Host          string
	Port          int
	Timeout       time.Duration
	conn          net.Conn
	mu            sync.Mutex
	transactionID uint16
}

// NewTCPTransport 建立新的 TCP 傳輸實例
func NewTCPTransport(host string, port int) *TCPTransport {
	if port == 0 {
		port = TCPDefaultPort
	}
	return &TCPTransport{
		Host:          host,
		Port:          port,
		Timeout:       2 * time.Second,
		transactionID: 0,
	}
}

func (t *TCPTransport) Connect() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		t.conn.Close()
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr := net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return fmt.Errorf("TCP 連線失敗: %w", err)
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

func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return nil, ErrConnectionClosed
	}

	if err := t.conn.SetDeadline(time.Now().Add(t.Timeout)); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 設定逾時失敗: %w", err)
	}

	// 寫入請求
	if _, err := t.conn.Write(data); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 寫入失敗: %w", err)
	}

	// 讀取 MBAP 標頭 (7 bytes)
	header := make([]byte, MBAPHeaderLength)
	if _, err := io.ReadFull(t.conn, header); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 讀取標頭失敗: %w", err)
	}

	mbap, err := ParseMBAPHeader(header)
	if err != nil {
		t.internalClose()
		return nil, err
	}

	// 讀取 PDU (Length - 1, 因為 Length 包含 Unit ID)
	pduLength := int(mbap.Length) - 1
	if pduLength < 0 || pduLength > 256 {
		t.internalClose()
		return nil, ErrInvalidFrame
	}

	pdu := make([]byte, pduLength)
	if _, err := io.ReadFull(t.conn, pdu); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 讀取 PDU 失敗: %w", err)
	}

	// 返回完整回應 (MBAP + PDU)
	return append(header, pdu...), nil
}

func (t *TCPTransport) internalClose() {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
}

// GetNextTransactionID 獲取下一個交易 ID
func (t *TCPTransport) GetNextTransactionID() uint16 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.transactionID++
	if t.transactionID == 0 {
		t.transactionID = 1
	}
	return t.transactionID
}
