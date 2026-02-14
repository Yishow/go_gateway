package modbus

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// UDPTransport Modbus UDP 傳輸實作
type UDPTransport struct {
	Host          string
	Port          int
	Timeout       time.Duration
	conn          *net.UDPConn
	mu            sync.Mutex
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

func (u *UDPTransport) Connect() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.conn != nil {
		u.conn.Close()
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(u.Host, strconv.Itoa(u.Port)))
	if err != nil {
		return fmt.Errorf("UDP 地址解析失敗: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("UDP 連線失敗: %w", err)
	}

	u.conn = conn
	return nil
}

func (u *UDPTransport) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.conn != nil {
		err := u.conn.Close()
		u.conn = nil
		return err
	}
	return nil
}

func (u *UDPTransport) SendReceive(data []byte) ([]byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.conn == nil {
		return nil, ErrConnectionClosed
	}

	if err := u.conn.SetDeadline(time.Now().Add(u.Timeout)); err != nil {
		return nil, fmt.Errorf("UDP 設定逾時失敗: %w", err)
	}

	// 寫入請求
	if _, err := u.conn.Write(data); err != nil {
		return nil, fmt.Errorf("UDP 寫入失敗: %w", err)
	}

	// 讀取回應
	buffer := make([]byte, 256)
	n, err := u.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("UDP 讀取失敗: %w", err)
	}

	return buffer[:n], nil
}

func (u *UDPTransport) GetNextTransactionID() uint16 {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.transactionID++
	if u.transactionID == 0 {
		u.transactionID = 1
	}
	return u.transactionID
}
