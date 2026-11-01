package connector

import (
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionConnectionPreviewNeverCreatesEntry(t *testing.T) {
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	view := cm.PeekAcquisitionConnection(schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`})
	require.Equal(t, "absent", view.State)
	require.Empty(t, cm.entries)
	cm.entries["d"] = &connectionEntry{done: make(chan struct{})}
	require.Equal(t, "unavailable", cm.PeekAcquisitionConnection(schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}).State)
}

func TestAcquisitionConnectionPreviewExactTokenAndConcurrentGeneration(t *testing.T) {
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	device := schema.Device{ID: "d", Protocol: schema.ProtocolMC3E, ConnectionConfig: `{}`}
	identity, _, err := connectionIdentity(device.Protocol, device.ConnectionConfig)
	require.NoError(t, err)
	protocol := &tokenProtocol{}
	protocol.generation.Store(1)
	done := make(chan struct{})
	close(done)
	conn := &ManagedConnection{managerGeneration: 7, Protocol: protocol, DeviceID: "d", ProtocolType: device.Protocol, Config: device.ConnectionConfig}
	cm.entries["d"] = &connectionEntry{identity: identity, generation: 7, conn: conn, done: done}
	cm.connections["d"] = conn
	first := cm.PeekAcquisitionConnection(device)
	require.Equal(t, "current", first.State)
	require.EqualValues(t, 7, first.Token.ManagerGeneration)
	different := device
	different.ConnectionConfig = `{"port":1000}`
	require.Equal(t, "unavailable", cm.PeekAcquisitionConnection(different).State)
	completed := make(chan struct{})
	go func() {
		defer close(completed)
		for range 100 {
			protocol.generation.Add(1)
		}
	}()
	for range 100 {
		require.Equal(t, "current", cm.PeekAcquisitionConnection(device).State)
	}
	<-completed
	require.NotEqual(t, first.Token, cm.PeekAcquisitionConnection(device).Token)
	require.Len(t, cm.entries, 1)
}
