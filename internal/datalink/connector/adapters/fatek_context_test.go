package adapters

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"net"
	"testing"
	"time"
)

type fatekCanceledTransport struct{ sends int }

func (*fatekCanceledTransport) Connect() error                       { return nil }
func (*fatekCanceledTransport) Close() error                         { return nil }
func (t *fatekCanceledTransport) SendReceive([]byte) ([]byte, error) { t.sends++; return nil, io.EOF }
func TestFatekCanceledReadDoesNotSend(t *testing.T) {
	transport := &fatekCanceledTransport{}
	c := &FatekConnector{client: fatek.NewClient(transport, 1), transport: transport, connected: true, persistentMode: true}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := c.Read(ctx, connector.ReadRequest{Address: "D0", Count: 1})
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, transport.sends)
}

func TestFatekProbeSharesCallerWireBudget(t *testing.T) {
	for _, limit := range []int{1, 2} {
		transport := &fatekCanceledTransport{}
		c := &FatekConnector{client: fatek.NewClient(transport, 1), transport: transport, connected: true, persistentMode: true}
		ctx := wirebudget.With(t.Context(), limit, 0)
		require.Error(t, c.TestConnection(ctx))
		require.Equal(t, limit, transport.sends)
		require.Equal(t, limit, wirebudget.From(ctx).Used())
	}
}

func TestFatekAdapterWaiterHonorsCallerDeadline(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	sent := make(chan struct{})
	finished := make(chan struct{})
	defer close(finished)
	go func() {
		conn, e := listener.Accept()
		if e != nil {
			return
		}
		defer conn.Close()
		_, _ = bufio.NewReader(conn).ReadBytes(fatek.ETX)
		close(sent)
		<-finished
	}()
	address := listener.Addr().(*net.TCPAddr)
	config, _ := json.Marshal(schema.ConnectionConfigFatekFBs{Mode: "tcp", Host: "127.0.0.1", Port: address.Port, Timeout: 1})
	c := &FatekConnector{}
	require.NoError(t, c.Connect(t.Context(), string(config)))
	defer c.Close()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, e := c.Read(ctx, connector.ReadRequest{Address: "D0", Count: 1}); done <- e }()
	<-sent
	waiter, stop := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer stop()
	start := time.Now()
	_, err = c.Read(waiter, connector.ReadRequest{Address: "D0", Count: 1})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), 150*time.Millisecond)
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
}
