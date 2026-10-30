package modbus

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"net"
	"testing"
	"time"
)

func TestUDPContextDeadlineBudgetAndReplacement(t *testing.T) {
	server, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	defer server.Close()
	tr := NewUDPTransport("127.0.0.1", server.LocalAddr().(*net.UDPAddr).Port)
	require.NoError(t, tr.Connect())
	defer tr.Close()
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	ctx = wirebudget.With(ctx, 1, 0)
	done := make(chan error, 1)
	go func() { _, err := tr.SendReceiveContext(ctx, []byte{1}); done <- err }()
	require.NoError(t, server.SetReadDeadline(time.Now().Add(time.Second)))
	buf := make([]byte, 256)
	_, _, err = server.ReadFromUDP(buf)
	require.NoError(t, err)
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(150 * time.Millisecond):
		t.Fatal("UDP deadline ignored")
	}
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	require.NoError(t, tr.ConnectContext(t.Context()))
	go func() {
		n, addr, err := server.ReadFromUDP(buf)
		if err == nil {
			_, _ = server.WriteToUDP(buf[:n], addr)
		}
	}()
	response, err := tr.SendReceiveContext(t.Context(), []byte{2})
	require.NoError(t, err)
	require.Equal(t, []byte{2}, response)
}

func TestTransportContextAllLifecycleWaitersCancel(t *testing.T) {
	tcp := NewTCPTransport("127.0.0.1", 1)
	udp := NewUDPTransport("127.0.0.1", 1)
	rtu := NewRTUTransport("missing", 0, 0, 0, "", 0)
	for _, entry := range []struct {
		name      string
		gate      *wirebudget.Mutex
		transport interface {
			ConnectContext(context.Context) error
			CloseContext(context.Context) error
			IsConnectedContext(context.Context) (bool, error)
			SendReceiveContext(context.Context, []byte) ([]byte, error)
		}
	}{{"tcp", &tcp.mu, tcp}, {"udp", &udp.mu, udp}, {"rtu", &rtu.mu, rtu}} {
		t.Run(entry.name, func(t *testing.T) {
			entry.gate.Lock()
			defer entry.gate.Unlock()
			methods := []func(context.Context) error{entry.transport.ConnectContext, entry.transport.CloseContext, func(ctx context.Context) error { _, err := entry.transport.IsConnectedContext(ctx); return err }, func(ctx context.Context) error {
				_, err := entry.transport.SendReceiveContext(ctx, []byte{1})
				return err
			}}
			for _, method := range methods {
				ctx, cancel := context.WithTimeout(t.Context(), 5*time.Millisecond)
				done := make(chan error, 1)
				go func() { done <- method(ctx) }()
				select {
				case err := <-done:
					require.ErrorIs(t, err, context.DeadlineExceeded)
				case <-time.After(150 * time.Millisecond):
					t.Fatal("gate ignored deadline")
				}
				cancel()
			}
		})
	}
}
