package fatek

import (
	"bufio"
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"go.bug.st/serial"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeContextPort struct {
	serial.Port
	written  chan struct{}
	closed   chan struct{}
	once     sync.Once
	sends    atomic.Int32
	closes   atomic.Int32
	zero     bool
	closeErr error
}

func newContextPort() *fakeContextPort {
	return &fakeContextPort{written: make(chan struct{}, 1), closed: make(chan struct{})}
}
func (p *fakeContextPort) SetReadTimeout(time.Duration) error { return nil }
func (p *fakeContextPort) Write(b []byte) (int, error) {
	p.sends.Add(1)
	select {
	case p.written <- struct{}{}:
	default:
	}
	return len(b), nil
}
func (p *fakeContextPort) Read([]byte) (int, error) {
	if p.zero {
		return 0, nil
	}
	<-p.closed
	return 0, net.ErrClosed
}
func (p *fakeContextPort) Close() error {
	p.closes.Add(1)
	p.once.Do(func() { close(p.closed) })
	return p.closeErr
}
func TestFatekSerialCancellationOwnsCapturedPort(t *testing.T) {
	p := newContextPort()
	s := &SerialTransport{port: p, reader: bufio.NewReader(p), Timeout: time.Second}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { _, err := s.SendReceiveContext(wirebudget.With(ctx, 1, 0), []byte("request")); done <- err }()
	<-p.written
	waiter, stop := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer stop()
	_, err := s.SendReceiveContext(waiter, []byte("waiter"))
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.EqualValues(t, 1, p.sends.Load())
	require.Zero(t, p.closes.Load())
	cancel()
	select {
	case err = <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(150 * time.Millisecond):
		t.Fatal("serial read did not unblock")
	}
	require.Nil(t, s.port)
	require.EqualValues(t, 1, p.closes.Load())
	replacement := newContextPort()
	s.mu.Lock()
	s.port = replacement
	s.mu.Unlock()
	require.Zero(t, replacement.closes.Load())
	require.NoError(t, s.Close())
}
func TestFatekSerialFailedCloseRetainsHandle(t *testing.T) {
	p := newContextPort()
	p.closeErr = errors.New("close failure")
	s := &SerialTransport{port: p, Timeout: time.Second}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { _, err := s.SendReceiveContext(ctx, []byte("read")); done <- err }()
	<-p.written
	cancel()
	err := <-done
	require.ErrorIs(t, err, context.Canceled)
	require.ErrorIs(t, err, p.closeErr)
	require.Same(t, p, s.port)
	_, err = s.SendReceiveContext(t.Context(), []byte("no replay"))
	require.ErrorIs(t, err, p.closeErr)
	require.EqualValues(t, 1, p.sends.Load())
	require.Error(t, s.Close())
	require.Same(t, p, s.port)
}
func TestFatekSerialZeroReadDoesNotLoop(t *testing.T) {
	p := newContextPort()
	p.zero = true
	s := &SerialTransport{port: p, Timeout: time.Second}
	start := time.Now()
	_, err := s.SendReceiveContext(t.Context(), []byte("request"))
	require.Error(t, err)
	require.Less(t, time.Since(start), 150*time.Millisecond)
}
func TestFatekTCPBudgetAndCancellation(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()
	tr := &TCPTransport{conn: client, reader: bufio.NewReader(client), Timeout: time.Second}
	read := make(chan struct{})
	go func() { buf := make([]byte, 1); _, _ = server.Read(buf); close(read) }()
	ctx, cancel := context.WithCancel(t.Context())
	budget := wirebudget.With(ctx, 1, 0)
	done := make(chan error, 1)
	go func() { _, err := tr.SendReceiveContext(budget, []byte{1}); done <- err }()
	<-read
	waiter, stop := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer stop()
	_, err := tr.IsConnectedContext(waiter)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	require.Equal(t, 1, wirebudget.From(budget).Used())
	require.Nil(t, tr.conn)
}
func TestFatekTransportCanceledContextConsumesNoBudget(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	ctx = wirebudget.With(ctx, 1, 0)
	for _, send := range []func(context.Context, []byte) ([]byte, error){(&TCPTransport{}).SendReceiveContext, (&SerialTransport{}).SendReceiveContext} {
		_, err := send(ctx, nil)
		require.ErrorIs(t, err, context.Canceled)
	}
	require.Zero(t, wirebudget.From(ctx).Used())
}
