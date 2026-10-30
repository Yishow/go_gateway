package collector

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mcBudgetServer counts actual protocol requests, not Protocol.Read invocations.
func mcBudgetServer(t *testing.T, mode string) (string, *atomic.Int32) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	var count atomic.Int32
	var workers sync.WaitGroup
	stop := make(chan struct{})
	workers.Go(func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			workers.Go(func() {
				defer conn.Close()
				for {
					header := make([]byte, 9)
					if _, err := io.ReadFull(conn, header); err != nil {
						return
					}
					size := int(binary.LittleEndian.Uint16(header[7:]))
					if _, err := io.CopyN(io.Discard, conn, int64(size)); err != nil {
						return
					}
					count.Add(1)
					switch mode {
					case "eof":
						return
					case "stall":
						<-stop
						return
					case "unsupported":
						_, _ = conn.Write([]byte{0xD0, 0, 0, 0, 0, 0, 0, 2, 0, 0x59, 0xC0})
					case "short":
						_, _ = conn.Write([]byte{0xD0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0})
					}
				}
			})
		}
	})
	t.Cleanup(func() { close(stop); _ = listener.Close(); workers.Wait() })
	return listener.Addr().String(), &count
}

func mcBudgetFixture(t *testing.T, address string, retry int, delay time.Duration) (*Scheduler, *connector.ManagedConnection, pointInfo) {
	t.Helper()
	host, portText, err := net.SplitHostPort(address)
	require.NoError(t, err)
	port, err := net.LookupPort("tcp", portText)
	require.NoError(t, err)
	configJSON, err := json.Marshal(schema.ConnectionConfigMC3E{Host: host, Port: port, Timeout: 1})
	require.NoError(t, err)
	adapter := &adapters.MC3EConnector{}
	require.NoError(t, adapter.Connect(t.Context(), string(configJSON)))
	t.Cleanup(func() { _ = adapter.Close() })
	conn := &connector.ManagedConnection{DeviceID: "device", ProtocolType: schema.ProtocolMC3E, Config: string(configJSON), Protocol: adapter}
	config := DefaultSchedulerConfig()
	config.DefaultRetryCount = retry
	config.DefaultRetryDelay = delay
	scheduler := NewScheduler(config, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	scheduler.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: string(configJSON)})
	scheduler.AddPoint(&schema.Point{ID: "point", DeviceID: "device", Address: "D0", DataType: schema.DataTypeUint16})
	return scheduler, conn, scheduler.pointInfos["point"]
}

func TestSchedulerMCActualWireBudgetAcrossAllLayers(t *testing.T) {
	for _, retries := range []int{0, 1, 3} {
		t.Run(string(rune('0'+retries)), func(t *testing.T) {
			address, sends := mcBudgetServer(t, "eof")
			scheduler, conn, point := mcBudgetFixture(t, address, retries, 0)
			err := scheduler.pollPointWithAcquisition(t.Context(), conn, point, "scan")
			require.Error(t, err)
			require.EqualValues(t, retries+1, sends.Load())
			value := <-scheduler.ValueChannel()
			require.Equal(t, schema.QualityBad, value.Quality)
			select {
			case <-scheduler.ValueChannel():
				t.Fatal("duplicate final result")
			default:
			}
		})
	}
}
func TestSchedulerMCUnsupportedDoesNotMultiplyRequests(t *testing.T) {
	address, sends := mcBudgetServer(t, "unsupported")
	scheduler, conn, point := mcBudgetFixture(t, address, 3, 0)
	require.Error(t, scheduler.pollPointWithAcquisition(t.Context(), conn, point, "scan"))
	require.EqualValues(t, 1, sends.Load())
}
func TestSchedulerMCCancelAndBackoffDoNotSendAgain(t *testing.T) {
	for _, mode := range []string{"eof", "stall"} {
		t.Run(mode, func(t *testing.T) {
			address, sends := mcBudgetServer(t, mode)
			scheduler, conn, point := mcBudgetFixture(t, address, 3, time.Second)
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Millisecond)
			defer cancel()
			begin := time.Now()
			err := scheduler.pollPointWithAcquisition(ctx, conn, point, "scan")
			require.ErrorIs(t, err, context.DeadlineExceeded)
			require.Less(t, time.Since(begin), 500*time.Millisecond)
			require.EqualValues(t, 1, sends.Load())
		})
	}
}
func TestSchedulerMCShortResponseRetriesAreBoundedWithoutFallback(t *testing.T) {
	address, sends := mcBudgetServer(t, "short")
	scheduler, conn, point := mcBudgetFixture(t, address, 1, 0)
	require.Error(t, scheduler.pollPointWithAcquisition(t.Context(), conn, point, "scan"))
	require.EqualValues(t, 2, sends.Load())
	value := <-scheduler.ValueChannel()
	require.Equal(t, schema.QualityBad, value.Quality)
}

func TestSchedulerMCConfiguredDeadlineIncludesAllRetries(t *testing.T) {
	address, sends := mcBudgetServer(t, "eof")
	scheduler, conn, point := mcBudgetFixture(t, address, 3, 3*time.Second)
	begin := time.Now()
	require.ErrorIs(t, scheduler.pollPointWithAcquisition(t.Context(), conn, point, "scan"), context.DeadlineExceeded)
	require.Less(t, time.Since(begin), 1500*time.Millisecond)
	require.EqualValues(t, 1, sends.Load())
}

func TestManualMCUsesConfiguredRequestBudget(t *testing.T) {
	for _, mode := range []string{"retry-zero", "pre-canceled"} {
		t.Run(mode, func(t *testing.T) {
			address, sends := mcBudgetServer(t, "eof")
			scheduler, _, _ := mcBudgetFixture(t, address, 0, 0)
			t.Cleanup(func() { _ = scheduler.connMgr.CloseAll() })
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if mode == "pre-canceled" {
				cancel()
			}
			result := scheduler.PollNowContext(ctx, []string{"point"})
			require.Len(t, result, 1)
			require.Equal(t, schema.QualityBad, result[0].Quality)
			expected := int32(1)
			if mode == "pre-canceled" {
				expected = 0
			}
			require.Equal(t, expected, sends.Load())
		})
	}
}
