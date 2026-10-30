package collector

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"

	"github.com/stretchr/testify/require"
)

type fatekScalarWireRequest struct {
	Command string
	Count   int
	Address string
}

type fatekScalarCase struct {
	name      string
	address   string
	dataType  schema.DataType
	payload   string
	want      any
	wireCount int
}

func TestFatekScalarPointCountAcrossEntrypoints(t *testing.T) {
	cases := []fatekScalarCase{
		{name: "D float32", address: "D0", dataType: schema.DataTypeFloat32, payload: "3F800000", want: float32(1), wireCount: 2},
		{name: "R uint64", address: "R0", dataType: schema.DataTypeUint64, payload: "0102030405060708", want: uint64(0x0102030405060708), wireCount: 4},
		{name: "DR float32", address: "DR0", dataType: schema.DataTypeFloat32, payload: "3F800000", want: float32(1), wireCount: 1},
		{name: "DD uint64", address: "DD0", dataType: schema.DataTypeUint64, payload: "0102030405060708", want: uint64(0x0102030405060708), wireCount: 2},
		{name: "D int16", address: "D0", dataType: schema.DataTypeInt16, payload: "0001", want: int16(1), wireCount: 1},
		{name: "D uint16", address: "D0", dataType: schema.DataTypeUint16, payload: "0001", want: uint16(1), wireCount: 1},
		{name: "X bool", address: "X0", dataType: schema.DataTypeBool, payload: "0", want: false, wireCount: 1},
	}
	for _, tc := range cases {
		for _, entry := range []string{"periodic", "manual", "read-now"} {
			t.Run(tc.name+"/"+entry, func(t *testing.T) {
				value, request := readFatekScalarThroughEntryPoint(t, tc, entry)
				require.Equal(t, tc.want, value)
				require.Equal(t, tc.wireCount, request.Count)
				expectedCommand := "46"
				if tc.address == "X0" {
					expectedCommand = "44"
				}
				require.Equal(t, expectedCommand, request.Command)
			})
		}
	}
}

func readFatekScalarThroughEntryPoint(t *testing.T, tc fatekScalarCase, entry string) (any, fatekScalarWireRequest) {
	t.Helper()
	config, requests := newFatekScalarLoopback(t, tc.payload)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	deviceRecord := &schema.Device{ID: "fatek-count", Protocol: schema.ProtocolFatekFBs, ConnectionConfig: config}
	point := &schema.Point{ID: "fatek-point", DeviceID: deviceRecord.ID, Address: tc.address, DataType: tc.dataType}
	scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
	scheduler.AddDevice(deviceRecord)
	scheduler.AddPoint(point)

	var value any
	switch entry {
	case "periodic":
		conn, err := manager.GetOrCreate(t.Context(), deviceRecord.ID, deviceRecord.Protocol, config)
		require.NoError(t, err)
		err = scheduler.pollPointWithAcquisition(t.Context(), conn, scheduler.pointInfos[point.ID], "fatek-acquisition")
		require.NoError(t, err)
		value = (<-scheduler.ValueChannel()).Value
	case "manual":
		values := scheduler.PollNowContext(t.Context(), []string{point.ID})
		require.Len(t, values, 1)
		value = values[0].Value
	case "read-now":
		repo := device.NewMemoryRepository()
		require.NoError(t, repo.Create(t.Context(), deviceRecord))
		result, err := device.NewService(repo, manager).ReadPointNow(t.Context(), point)
		require.NoError(t, err)
		value = result.Value
	default:
		t.Fatalf("unknown entrypoint %q", entry)
	}

	select {
	case request := <-requests:
		return value, request
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for FATEK wire request")
		return nil, fatekScalarWireRequest{}
	}
}

func newFatekScalarLoopback(t *testing.T, payload string) (configJSON string, observed <-chan fatekScalarWireRequest) {
	t.Helper()
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	requests := make(chan fatekScalarWireRequest, 8)
	var connections []net.Conn
	var connectionsMu sync.Mutex
	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			connectionsMu.Lock()
			connections = append(connections, conn)
			connectionsMu.Unlock()
			wg.Go(func() {
				defer conn.Close()
				reader := bufio.NewReader(conn)
				for {
					frame, readErr := reader.ReadBytes(fatek.ETX)
					if readErr != nil {
						return
					}
					if len(frame) < 9 || frame[0] != fatek.STX {
						t.Errorf("invalid FATEK request frame: %q", frame)
						return
					}
					command := string(frame[3:5])
					body := string(frame[5 : len(frame)-3])
					if len(body) < 3 {
						t.Errorf("invalid FATEK request body: %q", body)
						return
					}
					count, parseErr := strconv.ParseInt(body[:2], 16, 0)
					if parseErr != nil {
						t.Errorf("invalid FATEK request count: %q: %v", body[:2], parseErr)
						return
					}
					request := fatekScalarWireRequest{Command: command, Count: int(count), Address: body[2:]}
					requests <- request

					responsePayload := payload
					if command == "44" {
						responsePayload = strings.Repeat(payload, int(count))
					} else {
						componentWidth := 4
						if strings.HasPrefix(request.Address, "DR") || strings.HasPrefix(request.Address, "DD") {
							componentWidth = 8
						}
						responsePayload += strings.Repeat("0", max(int(count)*componentWidth-len(responsePayload), 0))
					}
					if _, writeErr := conn.Write(fatek.BuildFrame(1, command, "0"+responsePayload)); writeErr != nil {
						return
					}
				}
			})
		}
	})
	t.Cleanup(func() {
		_ = listener.Close()
		connectionsMu.Lock()
		for _, conn := range connections {
			_ = conn.Close()
		}
		connectionsMu.Unlock()
		wg.Wait()
	})
	return fmt.Sprintf(`{"mode":"tcp","host":"127.0.0.1","port":%d,"station_no":1,"timeout":1}`, listener.Addr().(*net.TCPAddr).Port), requests
}

type fatekCountProbeProtocol struct {
	readCalls int
}

func (p *fatekCountProbeProtocol) Connect(context.Context, string) error { return nil }
func (p *fatekCountProbeProtocol) Close() error                          { return nil }
func (p *fatekCountProbeProtocol) IsConnected() bool                     { return true }
func (p *fatekCountProbeProtocol) ProtocolType() schema.ProtocolType     { return schema.ProtocolFatekFBs }
func (p *fatekCountProbeProtocol) TestConnection(context.Context) error  { return nil }
func (p *fatekCountProbeProtocol) Write(context.Context, connector.WriteRequest) error {
	return nil
}
func (p *fatekCountProbeProtocol) Read(context.Context, connector.ReadRequest) (connector.ReadResult, error) {
	p.readCalls++
	return connector.ReadResult{Value: uint16(1), Quality: schema.QualityGood}, nil
}

func TestFatekInvalidScalarAddressFailsBeforeReadAcrossEntrypoints(t *testing.T) {
	for _, entry := range []string{"periodic", "manual", "read-now"} {
		t.Run(entry, func(t *testing.T) {
			probe := &fatekCountProbeProtocol{}
			old, exists := connector.GetFactory(schema.ProtocolFatekFBs)
			connector.Register(schema.ProtocolFatekFBs, func() connector.Protocol { return probe })
			t.Cleanup(func() {
				if exists {
					connector.Register(schema.ProtocolFatekFBs, old)
				} else {
					connector.Unregister(schema.ProtocolFatekFBs)
				}
			})

			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
			deviceRecord := &schema.Device{ID: "fatek-invalid", Protocol: schema.ProtocolFatekFBs, ConnectionConfig: "{}"}
			point := &schema.Point{ID: "invalid-point", DeviceID: deviceRecord.ID, Address: "RT0", DataType: schema.DataTypeFloat32}
			scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
			scheduler.AddDevice(deviceRecord)
			scheduler.AddPoint(point)

			switch entry {
			case "periodic":
				conn, err := manager.GetOrCreate(t.Context(), deviceRecord.ID, deviceRecord.Protocol, "{}")
				require.NoError(t, err)
				require.Error(t, scheduler.pollPointWithAcquisition(t.Context(), conn, scheduler.pointInfos[point.ID], "invalid"))
				result := <-scheduler.ValueChannel()
				require.Equal(t, schema.QualityBad, result.Quality)
				require.Nil(t, result.Value)
				require.NotEmpty(t, result.Error)
			case "manual":
				results := scheduler.PollNowContext(t.Context(), []string{point.ID})
				require.Len(t, results, 1)
				require.Equal(t, schema.QualityBad, results[0].Quality)
				require.Nil(t, results[0].Value)
				require.NotEmpty(t, results[0].Error)
			case "read-now":
				repo := device.NewMemoryRepository()
				require.NoError(t, repo.Create(t.Context(), deviceRecord))
				result, err := device.NewService(repo, manager).ReadPointNow(t.Context(), point)
				require.Error(t, err)
				require.Equal(t, schema.QualityBad, result.Quality)
				require.Nil(t, result.Value)
				require.NotEmpty(t, result.Error)
			}
			require.Equal(t, 0, probe.readCalls)
		})
	}
}
