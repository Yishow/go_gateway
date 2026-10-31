package adapters

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
)

func TestModbusPlannedTCPActiveCancelDoesNotPoisonReplacement(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	peer := &modbusPlannedPeer{kind: schema.ProtocolModbusTCP, memory: map[modbusPlannedAddress]uint16{{3, 0}: 7, {3, 1}: 8}}
	received := make(chan error, 1)
	serverDone := make(chan error, 1)
	go func() {
		first, err := listener.Accept()
		if err != nil {
			received <- err
			return
		}
		frame := make([]byte, 12)
		_, err = io.ReadFull(first, frame)
		if err == nil {
			_, err = peer.respond(frame)
		}
		received <- err
		var b [1]byte
		_, _ = first.Read(b[:])
		_ = first.Close()
		second, err := listener.Accept()
		if err != nil {
			serverDone <- err
			return
		}
		defer second.Close()
		for range 2 {
			frame := make([]byte, 12)
			if _, err = io.ReadFull(second, frame); err != nil {
				serverDone <- err
				return
			}
			response, readErr := peer.respond(frame)
			if readErr != nil {
				serverDone <- readErr
				return
			}
			if _, err = second.Write(response); err != nil {
				serverDone <- err
				return
			}
		}
		serverDone <- nil
	}()
	raw, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: listener.Addr().(*net.TCPAddr).Port, Timeout: 1})
	require.NoError(t, err)
	c := &ModbusTCPConnector{}
	require.NoError(t, c.Connect(t.Context(), string(raw)))
	defer c.Close()
	req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""), modbusPlannedPoint("b", "HR1", schema.DataTypeUint16, ""))
	base, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx := wirebudget.With(base, 3, 0)
	type outcome struct {
		got []connector.PlannedPointResult
		err error
	}
	done := make(chan outcome, 1)
	go func() { got, err := c.ReadPlanned(ctx, req); done <- outcome{got, err} }()
	select {
	case err := <-received:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("planned request not sent")
	}
	cancel()
	select {
	case result := <-done:
		require.ErrorIs(t, result.err, context.Canceled)
		require.Len(t, result.got, 2)
		for _, p := range result.got {
			require.Equal(t, schema.QualityBad, p.Quality)
			require.Nil(t, p.Value)
		}
	case <-time.After(time.Second):
		t.Fatal("planned read did not cancel")
	}
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	// The old operation has joined its cancellation callback. Later reads create
	// a fresh socket and retain their own identities and complete response values.
	for i := range 2 {
		req.RequestID = "replacement"
		got, err := c.ReadPlanned(t.Context(), req)
		require.NoError(t, err)
		assertModbusPlannedResults(t, req, got, map[string]any{"a": uint16(7), "b": uint16(8)})
		if i == 0 {
			cancel()
		}
	}
	require.NoError(t, <-serverDone)
	require.Len(t, peer.frames, 3)
	require.Equal(t, 1, wirebudget.From(ctx).Used())
}
