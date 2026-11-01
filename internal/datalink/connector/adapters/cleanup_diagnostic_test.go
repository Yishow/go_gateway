package adapters

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/require"
)

type fatekCleanupPeer struct {
	*fatekPlannedPeer
	closeErr error
	closes   int
}

func (p *fatekCleanupPeer) Close() error { p.closes++; return p.closeErr }

type modbusCleanupPeer struct {
	*modbusPlannedPeer
	closeErr error
	closes   int
}

func (p *modbusCleanupPeer) Close() error { p.closes++; return p.closeErr }

func cleanupReadFixture(kind schema.ProtocolType, primary, cleanup error) (adapterResult connector.Protocol, requestResult connector.ReadRequest, closeCount func() int, connectedState func() bool) {
	req := connector.ReadRequest{Address: "D0", DataType: schema.DataTypeUint16}
	if kind == schema.ProtocolFatekFBs {
		peer := &fatekCleanupPeer{fatekPlannedPeer: &fatekPlannedPeer{memory: map[string]any{"D0": 7}, err: primary}, closeErr: cleanup}
		c := &FatekConnector{client: fatek.NewClient(peer, 1), transport: peer, connected: true}
		return c, req, func() int { return peer.closes }, func() bool { return c.connected }
	}
	peer := &modbusCleanupPeer{modbusPlannedPeer: &modbusPlannedPeer{kind: kind, memory: map[modbusPlannedAddress]uint16{{3, 0}: 7}, failure: primary}, closeErr: cleanup}
	client := modbus.NewClient(peer, 1)
	req.Address = "HR0"
	switch kind {
	case schema.ProtocolModbusRTU:
		c := &ModbusRTUConnector{client: client, transport: peer, connected: true}
		return c, req, func() int { return peer.closes }, c.isConnected
	case schema.ProtocolModbusUDP:
		c := &ModbusUDPConnector{client: client, transport: peer, connected: true}
		return c, req, func() int { return peer.closes }, c.isConnected
	default:
		c := &ModbusTCPConnector{client: client, transport: peer, connected: true}
		return c, req, func() int { return peer.closes }, c.isConnected
	}
}

// Deliberately not parallel: restore the existing process logger after capture.
func TestAdapterCleanupDiagnosticPreservesReadOutcome(t *testing.T) {
	original := slog.Default()
	t.Cleanup(func() { slog.SetDefault(original) })
	primary := errors.New("primary read rejected")
	cleanup := errors.New("private-cleanup tcp://user:secret@hidden-endpoint")
	for _, kind := range []schema.ProtocolType{schema.ProtocolFatekFBs, schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, tc := range []struct {
			name             string
			primary, cleanup error
		}{
			{name: "successful-read-cleanup-failed", cleanup: cleanup},
			{name: "successful-read-cleanup-ok"},
			{name: "primary-error-cleanup-failed", primary: primary, cleanup: cleanup},
		} {
			t.Run(string(kind)+"/"+tc.name, func(t *testing.T) {
				var captured bytes.Buffer
				slog.SetDefault(slog.New(slog.NewJSONHandler(&captured, nil)))
				adapter, req, closes, connected := cleanupReadFixture(kind, tc.primary, tc.cleanup)
				result, err := adapter.Read(t.Context(), req)
				if tc.primary == nil {
					require.NoError(t, err, "cleanup cannot turn a known good read into a failed operation")
					require.Equal(t, schema.QualityGood, result.Quality)
					require.EqualValues(t, 7, result.Value)
					expectedRaw := []byte{0, 7}
					if kind == schema.ProtocolFatekFBs {
						expectedRaw = []byte{0, 0, 0, 7}
					}
					require.Equal(t, expectedRaw, result.RawBytes)
					require.Empty(t, result.Error)
				} else {
					require.ErrorIs(t, err, primary)
					require.NotErrorIs(t, err, cleanup, "cleanup must not change the operation error chain")
				}
				require.Equal(t, 1, closes())
				require.Equal(t, tc.cleanup != nil, connected(), "failed CloseContext retains live state")
				if tc.cleanup == nil {
					require.Empty(t, captured.String())
				} else {
					require.Contains(t, captured.String(), "adapter connection cleanup failed; operation outcome unchanged")
					require.Contains(t, captured.String(), `"protocol":"`+string(kind)+`"`)
					require.NotContains(t, captured.String(), "private-cleanup")
					require.NotContains(t, captured.String(), "hidden-endpoint")
					require.NotContains(t, captured.String(), "secret")
				}
			})
		}
	}
}
