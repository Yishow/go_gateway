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
	reader  *bufio.Reader
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
	t.reader = bufio.NewReader(conn)
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

	// Flush Buffer logic:
	// We want to discard any pending bytes from previous (timeout/error) operations.
	// Since bufio.Reader can buffer bytes, simply creating a new one isn't enough if data is in OS stack.
	// The best way is to ensure we read everything before writing, but "everything" is undefined if silent.
	// 
	// Optimization: If we trust the request-response lock-step, buffer should be empty.
	// If previous op timed out, we might be out of sync.
	// Recommendation: On timeout, Close() the connection. The caller (Client) should detect ErrConnectionClosed and Reconnect.
	// For this transport, if write/read fails, we close.
	
	// Write
	_, err := t.conn.Write(data)
	if err != nil {
		t.Close() // Force close on error to reset state
		return nil, err
	}

	// Read until ETX
	response, err := t.reader.ReadBytes(ETX)
	if err != nil {
		t.Close() // Force close on error (including timeout) to reset state
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
