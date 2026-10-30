package mcprotocol

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestMCContextCancelInterruptsOnlyCapturedSocket(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()
	transport := NewTCPTransport("unused", 0)
	transport.conn = client
	defer transport.Close()
	entered := make(chan struct{})
	go func() { var b [1]byte; _, _ = io.ReadFull(server, b[:]); close(entered) }()
	ctx, cancel := context.WithCancel(wirebudget.With(t.Context(), 2, 0))
	done := make(chan error, 1)
	go func() { _, err := transport.SendReceiveContext(ctx, []byte{1}); done <- err }()
	<-entered
	waiting, stopWaiting := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer stopWaiting()
	_, err := transport.SendReceiveContext(waiting, []byte{2})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	cancel()
	select {
	case err = <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("cancel did not release read")
	}
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	// A completed cancellation callback cannot poison a replacement socket.
	replacement, peer := net.Pipe()
	defer peer.Close()
	transport.mu.Lock()
	transport.conn = replacement
	transport.mu.Unlock()
	go func() {
		var b [1]byte
		_, _ = io.ReadFull(peer, b[:])
		_, _ = peer.Write([]byte{0xD0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0})
	}()
	response, err := transport.SendReceiveContext(wirebudget.With(t.Context(), 1, 0), []byte{3})
	require.NoError(t, err)
	require.Equal(t, []byte{0, 0}, response)
}

func TestMCTCPBudgetPreventsHiddenReconnectResend(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	var sends atomic.Int32
	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer c.Close()
				var b [1]byte
				if _, err := io.ReadFull(c, b[:]); err == nil {
					sends.Add(1)
				}
			}()
		}
	}()
	addr := listener.Addr().(*net.TCPAddr)
	transport := NewTCPTransport(addr.IP.String(), addr.Port)
	defer transport.Close()
	ctx := wirebudget.With(t.Context(), 1, 0)
	_, err = transport.SendReceiveContext(ctx, []byte{1})
	require.Error(t, err)
	_, err = transport.SendReceiveContext(ctx, []byte{1})
	require.ErrorIs(t, err, wirebudget.ErrExhausted)
	require.EqualValues(t, 1, sends.Load())
	require.Equal(t, 1, wirebudget.From(ctx).Used())
}

func TestMCCommandRejectionIsTypedAndNarrow(t *testing.T) {
	require.True(t, IsCommandRejected(&MCError{Code: 0xC059}))
	for _, err := range []error{io.EOF, context.DeadlineExceeded, errors.New("unsupported command"), &MCError{Code: 0xC05E}, &MCError{Code: 0xC05B}} {
		require.False(t, IsCommandRejected(err))
	}
}

type commandTransport struct {
	calls int
	code  uint16
}

func (t *commandTransport) Connect() error { return nil }
func (t *commandTransport) Close() error   { return nil }
func (t *commandTransport) SendReceive([]byte) ([]byte, error) {
	t.calls++
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, t.code)
	return b, nil
}
func TestMCClientContextViewDoesNotMutateOriginal(t *testing.T) {
	transport := &commandTransport{code: 0xC059}
	client := NewClientWithTransport(transport)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := client.WithContext(ctx).BatchReadWord("D", 0, 1)
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, transport.calls)
	_, err = client.BatchReadWord("D", 0, 1)
	require.True(t, IsCommandRejected(err))
	require.Equal(t, 1, transport.calls)
}
