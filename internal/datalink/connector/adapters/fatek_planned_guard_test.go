package adapters

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/wirebudget"
)

func TestFatekPlannedManagedGuardsAndFrozenProfiles(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "D0", schema.DataTypeUint16, ""))
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 7}}
	managed := &connector.ManagedConnection{DeviceID: "device", ProtocolType: schema.ProtocolFatekFBs, Config: "{}", Protocol: fatekPlannedConnector(peer)}
	got, err := managed.ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertFatekPlannedResults(t, req, got, map[string]any{"a": uint16(7)})
	for _, fault := range []string{"device", "config", "protocol", "binding", "empty", "request index"} {
		t.Run(fault, func(t *testing.T) {
			m := &connector.ManagedConnection{DeviceID: "device", ProtocolType: schema.ProtocolFatekFBs, Config: "{}", Protocol: managed.Protocol}
			bad := req
			bad.Bindings = append([]connector.PlannedPointBinding(nil), req.Bindings...)
			switch fault {
			case "device":
				m.DeviceID = "foreign"
			case "config":
				m.Config = `{"stationNo":2}`
			case "protocol":
				m.ProtocolType = schema.ProtocolMC3E
			case "binding":
				bad.Bindings[0].PointRevision = "new"
			case "empty":
				bad.Bindings = nil
			case "request index":
				bad.RequestIndex = 1
			}
			_, err := m.ReadPlanned(t.Context(), bad)
			require.Error(t, err)
			require.Len(t, peer.frames, 1)
		})
	}
	managed.Lock()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = managed.ReadPlanned(ctx, req)
	managed.Unlock()
	require.ErrorIs(t, err, context.Canceled)
	require.Len(t, peer.frames, 1)
}
func TestFatekPlannedManagedConcurrentReadersAreSerialized(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "D0", schema.DataTypeUint16, ""))
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 7}}
	managed := &connector.ManagedConnection{DeviceID: "device", ProtocolType: schema.ProtocolFatekFBs, Config: "{}", Protocol: fatekPlannedConnector(peer)}
	var wg sync.WaitGroup
	failures := make(chan error, 20)
	for range 20 {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := managed.ReadPlanned(t.Context(), req); failures <- err }()
	}
	wg.Wait()
	close(failures)
	for err := range failures {
		require.NoError(t, err)
	}
	require.Len(t, peer.frames, 20)
}
func TestFatekPlannedPreservesCallerAttemptBudget(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "D0", schema.DataTypeUint16, ""))
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 7}, err: io.EOF}
	c := fatekPlannedConnector(peer)
	ctx := wirebudget.With(t.Context(), 2, 0)
	for range 2 {
		got, err := readFatekPlanned(t, ctx, c, req)
		require.ErrorIs(t, err, io.EOF)
		require.Equal(t, schema.QualityBad, got[0].Quality)
	}
	got, err := readFatekPlanned(t, ctx, c, req)
	require.ErrorIs(t, err, wirebudget.ErrExhausted)
	require.Equal(t, schema.QualityBad, got[0].Quality)
	require.Equal(t, 2, wirebudget.From(ctx).Used())
	require.Len(t, peer.frames, 2)
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = readFatekPlanned(t, canceled, c, req)
	require.ErrorIs(t, err, context.Canceled)
	require.Len(t, peer.frames, 2)
}
func TestFatekPlannedTCPReadCancellationIsFinalBad(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	received := make(chan error, 1)
	release := make(chan struct{})
	defer close(release)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			received <- err
			return
		}
		defer connection.Close()
		frame, err := bufio.NewReader(connection).ReadBytes(3)
		if err == nil {
			_, _, err = parseFatekPlannedFrame(frame)
		}
		received <- err
		<-release
	}()
	address := listener.Addr().(*net.TCPAddr)
	transport := fatek.NewTCPTransport("127.0.0.1", address.Port)
	require.NoError(t, transport.Connect())
	defer transport.Close()
	c := &FatekConnector{client: fatek.NewClient(transport, 1), connected: true, persistentMode: true}
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "D0", schema.DataTypeUint16, ""), fatekPlannedPoint("b", "D10", schema.DataTypeUint16, ""))
	base, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx := wirebudget.With(base, 3, 0)
	type outcome struct {
		values []connector.PlannedPointResult
		err    error
	}
	done := make(chan outcome, 1)
	go func() { values, err := c.ReadPlanned(ctx, req); done <- outcome{values, err} }()
	require.NoError(t, <-received)
	cancel()
	select {
	case result := <-done:
		require.ErrorIs(t, result.err, context.Canceled)
		require.Len(t, result.values, 2)
		for _, value := range result.values {
			require.Equal(t, schema.QualityBad, value.Quality)
			require.Nil(t, value.Value)
		}
	case <-time.After(time.Second):
		t.Fatal("planned read did not release canceled transport")
	}
	require.Equal(t, 1, wirebudget.From(ctx).Used())
}
