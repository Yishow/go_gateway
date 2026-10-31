package adapters

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"sync"
	"testing"
)

func TestModbusPlannedManagedIdentityAndCanceledWaiter(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""))
			peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 7}}
			managed := &connector.ManagedConnection{DeviceID: "device", ProtocolType: kind, Config: "{}", Protocol: modbusPlannedProtocol(kind, peer)}
			got, err := managed.ReadPlanned(t.Context(), req)
			require.NoError(t, err)
			assertModbusPlannedResults(t, req, got, map[string]any{"a": uint16(7)})
			for _, fault := range []string{"device", "config", "protocol"} {
				copy := &connector.ManagedConnection{DeviceID: "device", ProtocolType: kind, Config: "{}", Protocol: managed.Protocol}
				switch fault {
				case "device":
					copy.DeviceID = "foreign"
				case "config":
					copy.Config = `{"slaveId":2}`
				case "protocol":
					copy.ProtocolType = schema.ProtocolFatekFBs
				}
				_, err = copy.ReadPlanned(t.Context(), req)
				require.Error(t, err)
				require.Len(t, peer.frames, 1)
			}
			entered, release := make(chan struct{}), make(chan struct{})
			peer.hook = func() { close(entered); <-release }
			done := make(chan error, 1)
			go func() { _, err := managed.ReadPlanned(t.Context(), req); done <- err }()
			<-entered
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			_, err = managed.ReadPlanned(ctx, req)
			require.ErrorIs(t, err, context.Canceled)
			close(release)
			require.NoError(t, <-done)
			require.Len(t, peer.frames, 2)
		})
	}
}
func TestModbusPlannedConcurrentManagedReadsUseFrozenProfiles(t *testing.T) {
	req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint32, "ABCD"), modbusPlannedPoint("b", "HR0", schema.DataTypeUint32, "CDAB"))
	peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 0x0102, {3, 1}: 0x0304}}
	managed := &connector.ManagedConnection{DeviceID: "device", ProtocolType: schema.ProtocolModbusTCP, Config: "{}", Protocol: modbusPlannedProtocol(schema.ProtocolModbusTCP, peer)}
	var wg sync.WaitGroup
	type outcome struct {
		got []connector.PlannedPointResult
		err error
	}
	done := make(chan outcome, 20)
	for range 20 {
		wg.Add(1)
		go func() { defer wg.Done(); got, err := managed.ReadPlanned(t.Context(), req); done <- outcome{got, err} }()
	}
	wg.Wait()
	close(done)
	for r := range done {
		require.NoError(t, r.err)
		assertModbusPlannedResults(t, req, r.got, map[string]any{"a": uint32(0x01020304), "b": uint32(0x03040102)})
	}
	require.Len(t, peer.frames, 20)
}
func TestModbusPlannedSnapshotCannotBeReboundDuringSend(t *testing.T) {
	req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint32, "ABCD"), modbusPlannedPoint("b", "HR0", schema.DataTypeUint32, "CDAB"))
	original := append([]connector.PlannedPointBinding(nil), req.Bindings...)
	peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 0x0102, {3, 1}: 0x0304}, hook: func() {
		req.Bindings[0].InstallationLease = 20
		req.Bindings[0].PointRevision = "new"
		req.Bindings[0].Point.SourceRevision = "new"
	}}
	got, err := requireModbusPlannedReader(t, modbusPlannedProtocol(schema.ProtocolModbusTCP, peer)).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	req.Bindings = original
	assertModbusPlannedResults(t, req, got, map[string]any{"a": uint32(0x01020304), "b": uint32(0x03040102)})
}

func TestModbusPlannedWrongStationAndFrameDoNotSend(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, fault := range []string{"station", "frame"} {
			t.Run(string(kind)+fault, func(t *testing.T) {
				req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""))
				r := req.Plan.Requests()[0]
				source := optimizer.AcquisitionSource{Identity: r.Group.Identity, Capability: optimizer.AcquisitionCapability{Revision: "cap-v1", ProtocolLimit: 125}, Points: []optimizer.EffectivePoint{req.Bindings[0].Point}}
				if fault == "station" {
					source.Identity.Station = "2"
				} else {
					source.Identity.Frame = "unknown"
				}
				req.Plan = optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{ConfigurationRevision: "config-v1", Sources: []optimizer.AcquisitionSource{source}})
				require.Empty(t, req.Plan.Errors())
				peer := &modbusPlannedPeer{}
				_, err := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer)).ReadPlanned(t.Context(), req)
				require.ErrorIs(t, err, connector.ErrPlannedRead)
				require.Empty(t, peer.frames)
			})
		}
	}
}
