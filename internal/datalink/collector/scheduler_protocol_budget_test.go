package collector

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func acquisitionBudgetTCPServer(t *testing.T, fatek bool) (int, *atomic.Int32) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	var sends atomic.Int32
	var workers sync.WaitGroup
	workers.Go(func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			workers.Go(func() {
				defer conn.Close()
				var err error
				if fatek {
					_, err = bufio.NewReader(conn).ReadBytes(3)
				} else {
					header := make([]byte, 7)
					if _, err = io.ReadFull(conn, header); err == nil {
						_, err = io.CopyN(io.Discard, conn, int64(binary.BigEndian.Uint16(header[4:]))-1)
					}
				}
				if err == nil {
					sends.Add(1)
				}
			})
		}
	})
	t.Cleanup(func() { listener.Close(); workers.Wait() })
	return listener.Addr().(*net.TCPAddr).Port, &sends
}
func TestSchedulerProtocolActualWireBudget(t *testing.T) {
	for _, typ := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolFatekFBs} {
		for _, retries := range []int{0, 1, 3} {
			t.Run(string(typ)+string(rune('0'+retries)), func(t *testing.T) {
				fatek := typ == schema.ProtocolFatekFBs
				port, sends := acquisitionBudgetTCPServer(t, fatek)
				var config any = schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: port, Timeout: 1}
				address := "0"
				if fatek {
					config = map[string]any{"host": "127.0.0.1", "port": port, "timeout": 1, "mode": "tcp"}
					address = "D0"
				}
				raw, err := json.Marshal(config)
				require.NoError(t, err)
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				defer manager.CloseAll()
				settings := DefaultSchedulerConfig()
				settings.DefaultRetryCount = retries
				settings.DefaultRetryDelay = time.Millisecond
				scheduler := NewScheduler(settings, manager)
				scheduler.AddDevice(&schema.Device{ID: "d", Protocol: typ, ConnectionConfig: string(raw)})
				scheduler.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: address, DataType: schema.DataTypeUint16})
				conn, err := manager.GetOrCreate(t.Context(), "d", typ, string(raw))
				require.NoError(t, err)
				err = scheduler.pollPointWithAcquisition(t.Context(), conn, scheduler.pointInfos["p"], "scan")
				require.Error(t, err)
				require.EqualValues(t, retries+1, sends.Load())
			})
		}
	}
}

func TestSchedulerProtocolBackoffCancellation(t *testing.T) {
	for _, typ := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolFatekFBs} {
		t.Run(string(typ), func(t *testing.T) {
			isFatek := typ == schema.ProtocolFatekFBs
			port, sends := acquisitionBudgetTCPServer(t, isFatek)
			config := map[string]any{"host": "127.0.0.1", "port": port, "timeout": 1, "mode": "tcp"}
			raw, err := json.Marshal(config)
			require.NoError(t, err)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			settings := DefaultSchedulerConfig()
			settings.DefaultRetryCount = 3
			settings.DefaultRetryDelay = time.Second
			scheduler := NewScheduler(settings, manager)
			address := "0"
			if isFatek {
				address = "D0"
			}
			scheduler.AddDevice(&schema.Device{ID: "d", Protocol: typ, ConnectionConfig: string(raw)})
			scheduler.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: address, DataType: schema.DataTypeUint16})
			conn, err := manager.GetOrCreate(t.Context(), "d", typ, string(raw))
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
			defer cancel()
			started := time.Now()
			err = scheduler.pollPointWithAcquisition(ctx, conn, scheduler.pointInfos["p"], "scan")
			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.Less(t, time.Since(started), 150*time.Millisecond)
			require.EqualValues(t, 1, sends.Load())
		})
	}
}
