package adapters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func TestMCPlannedManagedConnectionGuardsAndCancellation(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeInt16, ""))
	tr := &mcPlannedTransport{response: mcPlannedWords(7)}
	managed := &connector.ManagedConnection{Protocol: mcPlannedConnector(tr), ProtocolType: schema.ProtocolMC3E, Config: "{}", DeviceID: "device"}
	got, err := managed.ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, int16(7), got[0].Value)
	for _, bad := range []string{"config", "device", "protocol"} {
		t.Run(bad, func(t *testing.T) {
			copy := &connector.ManagedConnection{Protocol: managed.Protocol, ProtocolType: schema.ProtocolMC3E, Config: "{}", DeviceID: "device"}
			switch bad {
			case "config":
				copy.Config = `{"host":"foreign"}`
			case "device":
				copy.DeviceID = "foreign"
			case "protocol":
				copy.ProtocolType = schema.ProtocolFatekFBs
			}
			_, err := copy.ReadPlanned(t.Context(), req)
			require.Error(t, err)
			require.Len(t, tr.frames, 1)
		})
	}
	managed.Lock()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = managed.ReadPlanned(ctx, req)
	require.ErrorIs(t, err, context.Canceled)
	managed.Unlock()
	require.Len(t, tr.frames, 1)
}

func TestMCPlannedManagedConcurrentReadersStaySerialized(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "D0", schema.DataTypeInt16, ""))
	entered, release := make(chan struct{}), make(chan struct{})
	tr := &mcPlannedTransport{response: mcPlannedWords(7), hook: func() { close(entered); <-release }}
	managed := &connector.ManagedConnection{Protocol: mcPlannedConnector(tr), ProtocolType: schema.ProtocolMC3E, Config: "{}", DeviceID: "device"}
	done := make(chan error, 1)
	go func() { _, err := managed.ReadPlanned(t.Context(), req); done <- err }()
	<-entered
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := managed.ReadPlanned(ctx, req)
	require.ErrorIs(t, err, context.Canceled)
	close(release)
	require.NoError(t, <-done)
	require.Len(t, tr.frames, 1)
}
