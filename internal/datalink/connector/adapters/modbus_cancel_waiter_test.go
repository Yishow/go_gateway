package adapters

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"io"
	"net"
	"testing"
	"time"
)

func TestModbusReadWaiterCancellationAcrossManagerAndAdapter(t *testing.T) {
	for _, managed := range []bool{false, true} {
		t.Run(map[bool]string{false: "adapter", true: "manager"}[managed], func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer listener.Close()
			sent := make(chan struct{})
			wireClosed := make(chan struct{})
			serverDone := make(chan struct{})
			go func() {
				defer close(serverDone)
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				defer conn.Close()
				req := make([]byte, 12)
				if _, err = io.ReadFull(conn, req); err != nil {
					return
				}
				close(sent)
				var b [1]byte
				_, _ = conn.Read(b[:])
				close(wireClosed)
			}()
			defer func() { listener.Close(); <-serverDone }()
			address := listener.Addr().(*net.TCPAddr)
			raw, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: address.Port, Timeout: 1})
			require.NoError(t, err)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			var read func(context.Context, connector.ReadRequest) (connector.ReadResult, error)
			if managed {
				conn, err := manager.GetOrCreate(t.Context(), "d", schema.ProtocolModbusTCP, string(raw))
				require.NoError(t, err)
				read = conn.Read
			} else {
				adapter := &ModbusTCPConnector{}
				require.NoError(t, adapter.Connect(t.Context(), string(raw)))
				defer adapter.Close()
				read = adapter.Read
			}
			firstCtx, cancelFirst := context.WithCancel(t.Context())
			defer cancelFirst()
			firstDone := make(chan struct{})
			request := connector.ReadRequest{Address: "0", DataType: schema.DataTypeUint16}
			go func() { defer close(firstDone); _, _ = read(firstCtx, request) }()
			select {
			case <-sent:
			case <-time.After(time.Second):
				t.Fatal("first request not sent")
			}
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
			defer cancel()
			done := make(chan error, 1)
			go func() { _, err := read(ctx, request); done <- err }()
			select {
			case err := <-done:
				require.ErrorIs(t, err, context.DeadlineExceeded)
			case <-time.After(150 * time.Millisecond):
				t.Error("waiter ignored deadline")
				cancelFirst()
				<-done
			}
			select {
			case <-wireClosed:
				t.Error("waiter touched the active socket")
			default:
			}
			cancelFirst()
			select {
			case <-firstDone:
			case <-time.After(time.Second):
				t.Fatal("active read did not cancel")
			}
		})
	}
}
