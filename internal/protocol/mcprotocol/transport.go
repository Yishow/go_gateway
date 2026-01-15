package mcprotocol

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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