package modbus

import (
	"context"
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

	return t.internalConnect()
}

// IsConnected 檢查底層連線是否已建立且尚未關閉
func (t *TCPTransport) IsConnected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.conn != nil
}

func (t *TCPTransport) internalConnect() error {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr := net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
	dialCtx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(dialCtx, "tcp", addr)
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

	// 若連線已被關閉，嘗試自動重新連線
	if t.conn == nil {
		if err := t.internalConnect(); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrConnectionClosed, err)
		}
	}

	resp, err := t.sendReceiveLocked(data)
	if err != nil {
		// 若因對端逾時斷開導致寫入/讀取失敗，嘗試自癒重連並重試一次
		if t.conn == nil {
			if connErr := t.internalConnect(); connErr == nil {
				retryResp, retryErr := t.sendReceiveLocked(data)
				if retryErr == nil {
					return retryResp, nil
				}
			}
		}
		return nil, err
	}

	return resp, nil
}

func (t *TCPTransport) sendReceiveLocked(data []byte) ([]byte, error) {
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

	// Return the complete response (MBAP header plus PDU).
	response := make([]byte, 0, MBAPHeaderLength+len(pdu))
	response = append(response, header...)
	response = append(response, pdu...)
	return response, nil
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
