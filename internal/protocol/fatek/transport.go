package fatek

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
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
	reader  *bufio.Reader // 持久化的讀取器，避免緩衝區數據丟失
}

func NewTCPTransport(host string, port int) *TCPTransport {
	return &TCPTransport{
		Host:    host,
		Port:    port,
		Timeout: 2 * time.Second,
	}
}

func (t *TCPTransport) Connect() error {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	t.reader = bufio.NewReader(conn) // 建立持久化的讀取器
	return nil
}

func (t *TCPTransport) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		t.reader = nil
		return err
	}
	return nil
}

func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	if t.conn == nil {
		return nil, ErrConnectionClosed
	}

	// Set Deadline
	t.conn.SetDeadline(time.Now().Add(t.Timeout))

	// Write
	_, err := t.conn.Write(data)
	if err != nil {
		return nil, err
	}

	// Read until ETX
	// 使用持久化的 bufio.Reader 以避免每次創建新讀取器導致的緩衝區數據丟失
	if t.reader == nil {
		// 如果讀取器不存在（不應該發生，但為了安全）
		t.reader = bufio.NewReader(t.conn)
	}
	
	response, err := t.reader.ReadBytes(ETX)
	if err != nil {
		return nil, err
	}
	
	return response, nil
}

// SerialTransport (Stub for io.ReadWriteCloser compatible serial libs)
type SerialTransport struct {
	Port string
	// Underlying serial port interface usually matches io.ReadWriteCloser
	dev io.ReadWriteCloser 
}
// Implementation would depend on "go.bug.st/serial" or similar