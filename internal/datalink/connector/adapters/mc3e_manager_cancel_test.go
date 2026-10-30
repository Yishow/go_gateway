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

func TestMCManagerCanceledReuseAndRetirement(t *testing.T) {
	for _, replace := range []bool{false, true} {
		t.Run(map[bool]string{false: "reuse", true: "retire"}[replace], func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer listener.Close()
			entered, serverDone := make(chan struct{}), make(chan struct{})
			var extra atomic.Int64
			go func() {
				defer close(serverDone)
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				defer conn.Close()
				var request [21]byte
				if _, err = io.ReadFull(conn, request[:]); err != nil {
					return
				}
				close(entered)
				n, _ := io.Copy(io.Discard, conn)
				extra.Store(n)
			}()
			address := listener.Addr().(*net.TCPAddr)
			config := schema.ConnectionConfigMC3E{Host: address.IP.String(), Port: address.Port, Timeout: 5}
			data, err := json.Marshal(config)
			require.NoError(t, err)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			managed, err := manager.GetOrCreate(t.Context(), "device", schema.ProtocolMC3E, string(data))
			require.NoError(t, err)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			firstDone := make(chan error, 1)
			go func() {
				_, err := managed.Read(ctx, connector.ReadRequest{Address: "D0", Count: 1, DataType: schema.DataTypeInt16})
				firstDone <- err
			}()
			<-entered
			if replace {
				config.StationNo = 1
				data, err = json.Marshal(config)
				require.NoError(t, err)
			}
			waiter, stop := context.WithTimeout(t.Context(), 10*time.Millisecond)
			defer stop()
			done := make(chan error, 1)
			go func() {
				_, err := manager.GetOrCreate(waiter, "device", schema.ProtocolMC3E, string(data))
				done <- err
			}()
			select {
			case err := <-done:
				require.ErrorIs(t, err, context.DeadlineExceeded)
			case <-time.After(150 * time.Millisecond):
				t.Error("manager wait exceeded caller deadline")
				cancel()
				<-done
			}
			select {
			case <-serverDone:
				t.Error("canceled manager waiter closed the running socket")
			default:
			}
			cancel()
			require.ErrorIs(t, <-firstDone, context.Canceled)
			<-serverDone
			require.Zero(t, extra.Load())
			if replace {
				status, ok := manager.GetStatus("device")
				require.True(t, ok)
				require.True(t, status.CleanupPending)
			}
			require.NoError(t, manager.CloseAll())
			_, ok := manager.GetStatus("device")
			require.False(t, ok)
		})
	}
}
