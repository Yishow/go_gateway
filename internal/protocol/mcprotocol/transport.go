package mcprotocol

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Transport struct {
	Host    string
	Port    int
	Timeout time.Duration
	conn    net.Conn
	mu      sync.Mutex
}

func NewTransport(host string, port int) *Transport {
	return &Transport{
		Host:    host,
		Port:    port,
		Timeout: 2 * time.Second,
	}
}

func (t *Transport) Connect() error {
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

func (t *Transport) Close() error {
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
func (t *Transport) SendReceive(req []byte) ([]byte, error) {
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
	header := make([]byte, 9)
	if _, err := io.ReadFull(t.conn, header); err != nil {
		t.internalClose()
		return nil, err
	}

	// Parse Length
	dataLen, err := ParseResponseHeader(header)
	if err != nil {
		// Invalid header, maybe out of sync
		t.internalClose()
		return nil, err
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
func (t *Transport) internalConnect() error {
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
func (t *Transport) internalClose() {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
}
