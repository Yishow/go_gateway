package main

import (
	"context"
	"testing"

	"go-gateway/internal/protocol/modbus"
	"go-gateway/internal/protocol/wirebudget"

	"github.com/stretchr/testify/require"
)

func TestWordRTUFixtureExactContextBudgetAndCancel(t *testing.T) {
	plc := &wordWirePLC{words: map[int]uint16{0: 7}}
	peer := &wordRTUPeer{plc: plc}
	exact, ok := any(peer).(interface {
		SendReceiveExactContext(context.Context, []byte) ([]byte, error)
	})
	require.True(t, ok, "a whole-frame RTU peer must explicitly support exact context reads")
	request := modbus.BuildRTUFrame(1, 3, modbus.BuildReadRequest(0, 1))
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	response, err := exact.SendReceiveExactContext(canceled, request)
	require.ErrorIs(t, err, context.Canceled)
	require.Nil(t, response)
	require.Zero(t, plc.calls)
	ctx := wirebudget.With(t.Context(), 1, 0)
	response, err = exact.SendReceiveExactContext(ctx, request)
	require.NoError(t, err)
	unit, function, payload, err := modbus.ParseRTUFrame(response)
	require.NoError(t, err)
	require.Equal(t, byte(1), unit)
	require.Equal(t, byte(3), function)
	require.Equal(t, []byte{2, 0, 7}, payload)
	require.Equal(t, 1, plc.calls)
	response, err = exact.SendReceiveExactContext(ctx, request)
	require.ErrorIs(t, err, wirebudget.ErrExhausted)
	require.Nil(t, response)
	require.Equal(t, 1, plc.calls, "budget rejection must occur before simulated I/O")
}
