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

// UDPTransport Modbus UDP 傳輸實作
type UDPTransport struct {
	Host          string
	Port          int
	Timeout       time.Duration
	conn          *net.UDPConn
	mu            wirebudget.Mutex
	transactionID uint16
}

// NewUDPTransport 建立新的 UDP 傳輸實例
func NewUDPTransport(host string, port int) *UDPTransport {
	if port == 0 {
		port = UDPDefaultPort
	}
	return &UDPTransport{
		Host:          host,
		Port:          port,
		Timeout:       2 * time.Second,
		transactionID: 0,
	}
}

func (u *UDPTransport) Connect() error { return u.ConnectContext(context.Background()) }
func (u *UDPTransport) ConnectContext(ctx context.Context) error {
	ctx, cancel := transportContext(ctx, u.Timeout)
	defer cancel()
	if err := u.mu.LockContext(ctx); err != nil {
		return err
	}
	defer u.mu.Unlock()

	if u.conn != nil {
		u.conn.Close()
		u.conn = nil
	}

	conn, err := (&net.Dialer{}).DialContext(ctx, "udp", net.JoinHostPort(u.Host, strconv.Itoa(u.Port)))
	if err != nil {
		return fmt.Errorf("UDP 連線失敗: %w", err)
	}
	if err := ctx.Err(); err != nil {
		_ = conn.Close()
		return err
	}
	u.conn = conn.(*net.UDPConn)
	return nil
}

func (u *UDPTransport) Close() error { return u.CloseContext(context.Background()) }
func (u *UDPTransport) CloseContext(ctx context.Context) error {
	if err := u.mu.LockContext(ctx); err != nil {
		return err
	}
	defer u.mu.Unlock()

	if u.conn != nil {
		err := u.conn.Close()
		u.conn = nil
		return err
	}
	return nil
}

// IsConnected 檢查 UDP 連線狀態
func (u *UDPTransport) IsConnected() bool {
	live, _ := u.IsConnectedContext(context.Background())
	return live
}
func (u *UDPTransport) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := u.mu.LockContext(ctx); err != nil {
		return false, err
	}
	defer u.mu.Unlock()
	return u.conn != nil, nil
}

func (u *UDPTransport) SendReceive(data []byte) ([]byte, error) {
	return u.SendReceiveContext(context.Background(), data)
}
func (u *UDPTransport) SendReceiveContext(ctx context.Context, data []byte) (response []byte, resultErr error) {
	ctx, cancel := transportContext(ctx, u.Timeout)
	defer cancel()
	if err := u.mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer u.mu.Unlock()

	if u.conn == nil {
		return nil, ErrConnectionClosed
	}

	stop, err := socketContext(ctx, u.conn)
	if err != nil {
		u.internalClose()
		return nil, fmt.Errorf("UDP 設定逾時失敗: %w", err)
	}

	defer stop()
	defer func() {
		resultErr = transportContextError(ctx, resultErr)
	}()
	if err := wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	// 寫入請求
	if n, err := u.conn.Write(data); err != nil || n != len(data) {
		if err == nil {
			err = io.ErrShortWrite
		}
		return nil, fmt.Errorf("UDP 寫入失敗: %w", err)
	}

	// The existing MBAP-over-UDP contract allows the full 260-byte ADU.
	// One sentinel byte distinguishes an oversized datagram from a valid prefix.
	buffer := make([]byte, 261)
	n, err := u.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("UDP 讀取失敗: %w", err)
	}

	if n > 260 {
		return nil, ErrInvalidFrame
	}
	return buffer[:n], nil
}

func (u *UDPTransport) GetNextTransactionID() uint16 {
	id, _ := u.GetNextTransactionIDContext(context.Background())
	return id
}
func (u *UDPTransport) GetNextTransactionIDContext(ctx context.Context) (uint16, error) {
	if err := u.mu.LockContext(ctx); err != nil {
		return 0, err
	}
	defer u.mu.Unlock()
	u.transactionID++
	if u.transactionID == 0 {
		u.transactionID = 1
	}
	return u.transactionID, nil
}

func (u *UDPTransport) internalClose() {
	if u.conn != nil {
		if err := u.conn.Close(); err != nil {
			u.conn = nil
			return
		}
		u.conn = nil
	}
}
