package main

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func exerciseWordTransportBoundaries(t *testing.T, f wordProductionFixture) {
	t.Helper()
	_, zero, _ := f.createRule(t, schema.ProtocolModbusTCP, 0, schema.DataTypeUint32, "ABCD")
	f.plc.set(0, []uint16{0, 0})
	_, legacy, _ := f.createRule(t, schema.ProtocolModbusTCP, 1, schema.DataTypeUint32, "ABCD")
	legacy.FormatMode = schema.PointFormatLegacy
	legacy.FormatVersion = 0
	legacy.FormatOverride = ""
	legacy.DataFormat = "CDAB"
	require.NoError(t, f.services.point.Restore(t.Context(), legacy))
	f.plc.set(8, []uint16{0x3344, 0x1122})
	_, single, _ := f.createRule(t, schema.ProtocolModbusTCP, 2, schema.DataTypeUint16, "")
	f.plc.set(16, []uint16{0x1234})
	require.NoError(t, f.services.runtime.Start(t.Context()))
	for _, tc := range []struct {
		p    *schema.Point
		want any
	}{{zero, uint32(0)}, {legacy, uint32(0x11223344)}, {single, uint16(0x1234)}} {
		values := f.services.scheduler.PollNowContext(t.Context(), []string{tc.p.ID})
		require.Len(t, values, 1)
		require.Equal(t, schema.QualityGood, values[0].Quality)
		require.Equal(t, tc.want, values[0].Value)
	}
	f.plc.mu.Lock()
	f.plc.bad = true
	f.plc.mu.Unlock()
	bad := f.services.scheduler.PollNowContext(t.Context(), []string{zero.ID})
	require.Len(t, bad, 1)
	require.NotEqual(t, schema.QualityGood, bad[0].Quality)
	require.Nil(t, bad[0].Value)
	require.NotEmpty(t, bad[0].Error)
}

func TestProductionWordFormatUDPZeroShortAndLegacy(t *testing.T) {
	plc, config := wordDatagramPeer(t)
	exerciseWordTransportBoundaries(t, newWordTransportFixture(t, schema.ProtocolModbusUDP, plc, config))
}
