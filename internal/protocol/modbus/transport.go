package modbus

import (
	"context"
	"fmt"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"strconv"
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
	mu            wirebudget.Mutex
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

func (t *TCPTransport) Connect() error { return t.ConnectContext(context.Background()) }
func (t *TCPTransport) ConnectContext(ctx context.Context) error {
	ctx, cancel := transportContext(ctx, t.Timeout)
	defer cancel()
	if err := t.mu.LockContext(ctx); err != nil {
		return err
	}
	defer t.mu.Unlock()

	return t.internalConnectContext(ctx)
}

// IsConnected 檢查底層連線是否已建立且尚未關閉
func (t *TCPTransport) IsConnected() bool {
	live, _ := t.IsConnectedContext(context.Background())
	return live
}
func (t *TCPTransport) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := t.mu.LockContext(ctx); err != nil {
		return false, err
	}
	defer t.mu.Unlock()
	return t.conn != nil, nil
}

func (t *TCPTransport) internalConnectContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr := net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
	dialCtx, cancel := transportContext(ctx, t.Timeout)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(dialCtx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("TCP 連線失敗: %w", err)
	}

	if err := ctx.Err(); err != nil {
		_ = conn.Close()
		return err
	}
	t.conn = conn
	return nil
}

func (t *TCPTransport) Close() error { return t.CloseContext(context.Background()) }
func (t *TCPTransport) CloseContext(ctx context.Context) error {
	if err := t.mu.LockContext(ctx); err != nil {
		return err
	}
	defer t.mu.Unlock()

	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}

func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	return t.SendReceiveContext(context.Background(), data)
}
func (t *TCPTransport) SendReceiveContext(ctx context.Context, data []byte) ([]byte, error) {
	ctx, cancel := transportContext(ctx, t.Timeout)
	defer cancel()
	if err := t.mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer t.mu.Unlock()
	if b := wirebudget.From(ctx); b != nil && b.Remaining() == 0 {
		return nil, wirebudget.ErrExhausted
	}
	if t.conn == nil {
		if err := t.internalConnectContext(ctx); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrConnectionClosed, err)
		}
	}
	response, err := t.sendReceiveLocked(ctx, data)
	if ctxErr := transportContextError(ctx, nil); ctxErr != nil {
		return nil, ctxErr
	}
	if err != nil && wirebudget.From(ctx) == nil && t.conn == nil {
		if connectErr := t.internalConnectContext(ctx); connectErr == nil {
			return t.sendReceiveLocked(ctx, data)
		}
	}
	return response, err
}

func (t *TCPTransport) sendReceiveLocked(ctx context.Context, data []byte) ([]byte, error) {
	if t.conn == nil {
		return nil, ErrConnectionClosed
	}

	conn := t.conn
	stop, err := socketContext(ctx, conn)
	if err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 設定逾時失敗: %w", err)
	}

	defer stop()
	if err := wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	// 寫入請求
	if n, err := conn.Write(data); err != nil || n != len(data) {
		if err == nil {
			err = io.ErrShortWrite
		}
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
	id, _ := t.GetNextTransactionIDContext(context.Background())
	return id
}
func (t *TCPTransport) GetNextTransactionIDContext(ctx context.Context) (uint16, error) {
	if err := t.mu.LockContext(ctx); err != nil {
		return 0, err
	}
	defer t.mu.Unlock()
	t.transactionID++
	if t.transactionID == 0 {
		t.transactionID = 1
	}
	return t.transactionID, nil
}
