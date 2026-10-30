package adapters

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestMCAdapterCanceledWaiter(t *testing.T) {
	for _, operation := range []string{"read", "reconnect", "managed-read", "managed-write"} {
		t.Run(operation, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer listener.Close()
			entered := make(chan struct{})
			serverDone := make(chan struct{})
			var extraWireBytes atomic.Int64
			go func() {
				defer close(serverDone)
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				defer conn.Close()
				request := make([]byte, 21)
				if _, err = io.ReadFull(conn, request); err != nil {
					return
				}
				close(entered)
				n, _ := io.Copy(io.Discard, conn)
				extraWireBytes.Store(n)
			}()
			address := listener.Addr().(*net.TCPAddr)
			configJSON, err := json.Marshal(schema.ConnectionConfigMC3E{Host: address.IP.String(), Port: address.Port, Timeout: 5})
			require.NoError(t, err)
			adapter := &MC3EConnector{}
			require.NoError(t, adapter.Connect(t.Context(), string(configJSON)))
			defer adapter.Close()
			managed := &connector.ManagedConnection{Protocol: adapter, ProtocolType: schema.ProtocolMC3E}
			first, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			request := connector.ReadRequest{Address: "D0", DataType: schema.DataTypeInt16, Count: 1}
			go func() {
				var err error
				if operation == "managed-read" || operation == "managed-write" {
					_, err = managed.Read(first, request)
				} else {
					_, err = adapter.Read(first, request)
				}
				done <- err
			}()
			<-entered
			second, stop := context.WithTimeout(t.Context(), 10*time.Millisecond)
			defer stop()
			waiter := make(chan error, 1)
			go func() {
				if operation == "managed-read" {
					_, err := managed.Read(second, request)
					waiter <- err
					return
				}
				if operation == "managed-write" {
					waiter <- managed.Write(second, connector.WriteRequest{})
					return
				}
				if operation == "reconnect" {
					waiter <- adapter.Reconnect(second)
					return
				}
				_, err := adapter.Read(second, request)
				waiter <- err
			}()
			select {
			case err := <-waiter:
				require.ErrorIs(t, err, context.DeadlineExceeded)
			case <-time.After(150 * time.Millisecond):
				t.Error("canceled adapter waiter blocked behind connection operation lock")
				cancel()
				<-waiter
			}
			select {
			case <-serverDone:
				t.Error("canceled waiter closed another request's live socket")
			default:
			}
			cancel()
			require.ErrorIs(t, <-done, context.Canceled)
			<-serverDone
			require.Zero(t, extraWireBytes.Load(), "canceled waiter sent additional bytes")
		})
	}
}
