package collector

import (
	"encoding/binary"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"net"
	"sync/atomic"
	"testing"
)

func TestSchedulerModbusUDPCountsRealDatagrams(t *testing.T) {
	for _, unsupported := range []bool{false, true} {
		t.Run(map[bool]string{false: "short", true: "unsupported"}[unsupported], func(t *testing.T) {
			listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
			require.NoError(t, err)
			var sends atomic.Int32
			done := make(chan struct{})
			go func() {
				defer close(done)
				var buf [256]byte
				for {
					n, peer, err := listener.ReadFromUDP(buf[:])
					if err != nil {
						return
					}
					sends.Add(1)
					response := []byte{0}
					if unsupported && n >= 8 {
						response = append([]byte(nil), buf[:7]...)
						binary.BigEndian.PutUint16(response[4:6], 3)
						response = append(response, buf[7]|0x80, 1)
					}
					_, _ = listener.WriteToUDP(response, peer)
				}
			}()
			defer func() { listener.Close(); <-done }()
			raw, err := json.Marshal(schema.ConnectionConfigModbusUDP{Host: "127.0.0.1", Port: listener.LocalAddr().(*net.UDPAddr).Port, Timeout: 1})
			require.NoError(t, err)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			settings := DefaultSchedulerConfig()
			settings.DefaultRetryCount = 3
			settings.DefaultRetryDelay = 0
			scheduler := NewScheduler(settings, manager)
			scheduler.AddDevice(&schema.Device{ID: "d", Protocol: schema.ProtocolModbusUDP, ConnectionConfig: string(raw)})
			scheduler.AddPoint(&schema.Point{ID: "p", DeviceID: "d", Address: "0", DataType: schema.DataTypeUint16})
			conn, err := manager.GetOrCreate(t.Context(), "d", schema.ProtocolModbusUDP, string(raw))
			require.NoError(t, err)
			err = scheduler.pollPointWithAcquisition(t.Context(), conn, scheduler.pointInfos["p"], "scan")
			require.Error(t, err)
			expected := 4
			if unsupported {
				expected = 1
			}
			require.EqualValues(t, expected, sends.Load())
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
