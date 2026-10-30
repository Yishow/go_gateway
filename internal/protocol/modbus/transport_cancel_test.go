package modbus

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"go.bug.st/serial"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTCPContextBudgetCountsActualWriteWithoutHiddenRetry(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	var sends atomic.Int32
	done := make(chan struct{})
	defer func() { listener.Close(); <-done }()
	go func() {
		defer close(done)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			buf := make([]byte, 12)
			if _, err = io.ReadFull(conn, buf); err == nil {
				sends.Add(1)
			}
			conn.Close()
		}
	}()
	addr := listener.Addr().(*net.TCPAddr)
	tr := NewTCPTransport("127.0.0.1", addr.Port)
	defer tr.Close()
	ctx := wirebudget.With(t.Context(), 1, 0)
	_, err = tr.SendReceiveContext(ctx, BuildTCPFrame(1, 1, 3, []byte{0, 0, 0, 1}))
	require.Error(t, err)
	require.EqualValues(t, 1, sends.Load())
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	_, err = tr.SendReceiveContext(ctx, []byte{1})
	require.ErrorIs(t, err, wirebudget.ErrExhausted)
	require.EqualValues(t, 1, sends.Load())
}

func TestTCPContextCanceledWaiterAndCapturedSocket(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()
	tr := NewTCPTransport("", 0)
	tr.conn = client
	first, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := tr.SendReceiveContext(wirebudget.With(first, 1, 0), []byte{1}); done <- err }()
	buf := make([]byte, 1)
	_, err := io.ReadFull(server, buf)
	require.NoError(t, err)
	ctx, stop := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer stop()
	waiter := make(chan error, 1)
	go func() { _, err := tr.GetNextTransactionIDContext(ctx); waiter <- err }()
	select {
	case err := <-waiter:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(150 * time.Millisecond):
		t.Fatal("transaction gate ignored cancellation")
	}
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	next, peer := net.Pipe()
	defer peer.Close()
	tr.mu.Lock()
	tr.conn = next
	tr.mu.Unlock()
	defer tr.Close()
	go func() {
		b := make([]byte, 1)
		_, _ = peer.Read(b)
		_, _ = peer.Write(BuildTCPFrame(1, 1, 3, []byte{2, 0, 1}))
	}()
	_, err = tr.SendReceiveContext(t.Context(), []byte{1})
	require.NoError(t, err)
}

type cancelSerialPort struct {
	serial.Port
	closed   chan struct{}
	once     sync.Once
	writes   atomic.Int32
	closes   atomic.Int32
	closeErr error
	zero     bool
	data     io.Reader
}

func (p *cancelSerialPort) Write(b []byte) (int, error) { p.writes.Add(1); return len(b), nil }
func (p *cancelSerialPort) Read(b []byte) (int, error) {
	if p.data != nil {
		return p.data.Read(b)
	}
	if p.zero {
		return 0, nil
	}
	<-p.closed
	return 0, io.EOF
}
func (p *cancelSerialPort) SetReadTimeout(time.Duration) error { return nil }
func (p *cancelSerialPort) Close() error {
	p.closes.Add(1)
	p.once.Do(func() { close(p.closed) })
	return p.closeErr
}
func TestRTUContextCancellationClosesCapturedPortOnly(t *testing.T) {
	port := &cancelSerialPort{closed: make(chan struct{})}
	tr := NewRTUTransport("", 9600, 8, 1, "N", time.Second)
	tr.port = port
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := tr.SendReceiveContext(wirebudget.With(ctx, 1, 0), []byte{1}); done <- err }()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("serial read not interrupted")
	}
	require.EqualValues(t, 1, port.writes.Load())
	require.EqualValues(t, 1, port.closes.Load())
	require.Nil(t, tr.port)
	replacement := &cancelSerialPort{closed: make(chan struct{})}
	tr.port = replacement
	time.Sleep(10 * time.Millisecond)
	require.Zero(t, replacement.closes.Load())
	require.NoError(t, tr.Close())
}
func TestRTUContextRetainsFailedCloseAndStopsZeroRead(t *testing.T) {
	failure := errors.New("close failed")
	port := &cancelSerialPort{closed: make(chan struct{}), closeErr: failure, zero: true}
	tr := NewRTUTransport("", 9600, 8, 1, "N", time.Second)
	tr.port = port
	_, err := tr.SendReceiveContext(wirebudget.With(t.Context(), 1, 0), []byte{1})
	require.ErrorIs(t, err, io.ErrNoProgress)
	require.Same(t, port, tr.port)
	require.False(t, tr.IsConnected())
	_, err = tr.SendReceiveContext(t.Context(), []byte{1})
	require.ErrorIs(t, err, failure)
	require.EqualValues(t, 1, port.writes.Load())
	require.ErrorIs(t, tr.Close(), failure)
	require.Same(t, port, tr.port)
	port.closeErr = nil
	require.NoError(t, tr.Close())
	require.Nil(t, tr.port)
	require.Nil(t, tr.closeErr)
}

func TestRTUContextBudgetSuccessThenExhausted(t *testing.T) {
	port := &cancelSerialPort{closed: make(chan struct{}), data: bytes.NewReader(BuildRTUFrame(1, 3, []byte{2, 0, 1}))}
	tr := NewRTUTransport("", 9600, 8, 1, "N", time.Second)
	tr.port = port
	defer tr.Close()
	ctx := wirebudget.With(t.Context(), 1, 0)
	response, err := tr.SendReceiveContext(ctx, BuildRTUFrame(1, 3, []byte{0, 0, 0, 1}))
	require.NoError(t, err)
	require.NotEmpty(t, response)
	_, err = tr.SendReceiveContext(ctx, []byte{1})
	require.ErrorIs(t, err, wirebudget.ErrExhausted)
	require.EqualValues(t, 1, port.writes.Load())
	require.Zero(t, port.closes.Load())
}
