package collector

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// Exercise the real MC transport/decoder together with each single-point caller.
func TestMCScalarPointCountAcrossEntrypoints(t *testing.T) {
	cases := []struct {
		typ   schema.DataType
		words []uint16
		want  any
	}{
		{schema.DataTypeInt32, []uint16{65535, 65535}, int32(-1)},
		{schema.DataTypeUint32, []uint16{34464, 1}, uint32(100000)},
		{schema.DataTypeFloat32, []uint16{0, 16840}, float32(25)},
		{schema.DataTypeInt64, []uint16{65535, 65535, 65535, 32767}, int64(9223372036854775807)},
		{schema.DataTypeUint64, []uint16{1, 0, 0, 32}, uint64(9007199254740993)},
		{schema.DataTypeFloat64, []uint16{0, 0, 0, 16441}, float64(25)},
	}
	for _, tc := range cases {
		for _, entry := range []string{"periodic", "manual", "read-now"} {
			t.Run(string(tc.typ)+"/"+entry, func(t *testing.T) {
				config, requests := mcWordFixture(t, tc.words)
				connector.Register(schema.ProtocolMC3E, adapters.NewMC3EConnector)
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				t.Cleanup(func() { manager.CloseAll() })
				dev := &schema.Device{ID: "mc-count", Protocol: schema.ProtocolMC3E, ConnectionConfig: config}
				point := &schema.Point{ID: "point", DeviceID: dev.ID, Address: "D3", DataType: tc.typ}
				scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
				scheduler.AddDevice(dev)
				scheduler.AddPoint(point)
				var value any
				switch entry {
				case "periodic":
					conn, err := manager.GetOrCreate(t.Context(), dev.ID, dev.Protocol, config)
					require.NoError(t, err)
					require.NoError(t, scheduler.pollPointWithAcquisition(t.Context(), conn, scheduler.pointInfos[point.ID], "acquisition"))
					value = (<-scheduler.ValueChannel()).Value
				case "manual":
					values := scheduler.PollNowContext(t.Context(), []string{point.ID})
					require.Len(t, values, 1)
					value = values[0].Value
				case "read-now":
					repo := device.NewMemoryRepository()
					require.NoError(t, repo.Create(t.Context(), dev))
					result, err := device.NewService(repo, manager).ReadPointNow(t.Context(), point)
					require.NoError(t, err)
					value = result.Value
				}
				require.Equal(t, tc.want, value)
				require.Equal(t, len(tc.words), <-requests)
			})
		}
	}
}

func TestMCDirectTypedBatchCountIsPreserved(t *testing.T) {
	config, requests := mcWordFixture(t, []uint16{34464, 1, 34465, 1})
	protocol := adapters.NewMC3EConnector()
	require.NoError(t, protocol.Connect(t.Context(), config))
	t.Cleanup(func() { _ = protocol.Close() })
	result, err := protocol.Read(t.Context(), connector.ReadRequest{Address: "D3", DataType: schema.DataTypeUint32, Count: 2})
	require.NoError(t, err)
	require.Equal(t, []any{uint32(100000), uint32(100001)}, result.Value)
	require.Equal(t, 4, <-requests)
}

func mcWordFixture(t *testing.T, words []uint16) (configJSON string, counts <-chan int) {
	t.Helper()
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	requests := make(chan int, 16)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var connections []net.Conn
	wg.Go(func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			connections = append(connections, conn)
			mu.Unlock()
			wg.Go(func() {
				defer conn.Close()
				for {
					header := make([]byte, 9)
					if _, err := io.ReadFull(conn, header); err != nil {
						return
					}
					body := make([]byte, binary.LittleEndian.Uint16(header[7:]))
					if _, err := io.ReadFull(conn, body); err != nil {
						return
					}
					if len(body) != 12 {
						return
					}
					count := int(binary.LittleEndian.Uint16(body[10:]))
					requests <- count
					response := make([]byte, 2+2*count)
					for i := range count {
						if i < len(words) {
							binary.LittleEndian.PutUint16(response[2+2*i:], words[i])
						}
					}
					header[0], header[1] = 0xd0, 0
					binary.LittleEndian.PutUint16(header[7:], uint16(len(response))) //nolint:gosec // test response is bounded by small fixture requests
					if _, err := conn.Write(append(header, response...)); err != nil {
						return
					}
				}
			})
		}
	})
	t.Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		for _, conn := range connections {
			_ = conn.Close()
		}
		mu.Unlock()
		wg.Wait()
	})
	return fmt.Sprintf(`{"host":"127.0.0.1","port":%d}`, listener.Addr().(*net.TCPAddr).Port), requests
}

type pointCountProtocol struct {
	acquisitionFactsProtocol
	count int
}

func (p *pointCountProtocol) Read(_ context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	p.count = req.Count
	return connector.ReadResult{Value: uint64(9007199254740993), Quality: schema.QualityGood}, nil
}

func TestOtherProtocolPointWordCountsRemainUnchanged(t *testing.T) {
	for _, typ := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP, schema.ProtocolFatekFBs} {
		t.Run(string(typ), func(t *testing.T) {
			old, exists := connector.GetFactory(typ)
			protocol := &pointCountProtocol{}
			connector.Register(typ, func() connector.Protocol { return protocol })
			t.Cleanup(func() {
				if exists {
					connector.Register(typ, old)
				} else {
					connector.Unregister(typ)
				}
			})
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			t.Cleanup(func() { manager.CloseAll() })
			config := "{}"
			if typ == schema.ProtocolModbusRTU {
				config = `{"serial_port":"synthetic-point-count-port"}`
			}
			dev := &schema.Device{ID: "device", Protocol: typ, ConnectionConfig: config}
			address := "40001"
			if typ == schema.ProtocolFatekFBs {
				address = "D0"
			}
			point := &schema.Point{ID: "point", DeviceID: dev.ID, Address: address, DataType: schema.DataTypeUint64}
			s := NewScheduler(DefaultSchedulerConfig(), manager)
			s.AddDevice(dev)
			s.AddPoint(point)
			require.Len(t, s.PollNowContext(t.Context(), []string{point.ID}), 1)
			require.Equal(t, 4, protocol.count)
			conn, err := manager.GetOrCreate(t.Context(), dev.ID, typ, config)
			require.NoError(t, err)
			require.NoError(t, s.pollPointWithAcquisition(t.Context(), conn, s.pointInfos[point.ID], "periodic"))
			require.Equal(t, 4, protocol.count)
			repo := device.NewMemoryRepository()
			require.NoError(t, repo.Create(t.Context(), dev))
			_, err = device.NewService(repo, manager).ReadPointNow(t.Context(), point)
			require.NoError(t, err)
			require.Equal(t, 4, protocol.count)
		})
	}
}
