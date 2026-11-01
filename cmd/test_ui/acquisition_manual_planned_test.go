package main

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// The peer parses actual production adapter frames; these are loopback reads,
// not a physical PLC or SQL-delivery benchmark.
func TestManualPlannedExactSubsetUsesProductionAdapters(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, raw := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 11, 1: 22, 10: 33, 20: 44})
			var config map[string]any
			require.NoError(t, json.Unmarshal([]byte(raw), &config))
			profile := "safe-v1"
			wantRequests := 2
			if protocol == schema.ProtocolMC3E {
				profile, wantRequests = "mitsubishi-fx5-cpu-3e-binary-v1", 1
			}
			if protocol == schema.ProtocolFatekFBs {
				profile, wantRequests = "fatek-fbs-ascii-v1", 1
			}
			config["acquisition"] = map[string]any{"version": 1, "profile": profile}
			encoded, err := json.Marshal(config)
			require.NoError(t, err)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
			settings := collector.DefaultSchedulerConfig()
			settings.DefaultRetryCount = 0
			observations := []collector.AcquisitionObservation{}
			settings.AcquisitionObserver = func(v collector.AcquisitionObservation) { observations = append(observations, v) }
			scheduler := collector.NewScheduler(settings, manager)
			scheduler.AddDevice(&schema.Device{ID: "device", Protocol: protocol, ConnectionConfig: string(encoded)})
			addresses := []string{"D0", "D1", "D10", "D20"}
			if protocol == schema.ProtocolModbusTCP {
				addresses = []string{"40001", "40002", "40011", "40021"}
			}
			for i, id := range []string{"a", "b", "c", "unrequested"} {
				scheduler.AddPoint(&schema.Point{ID: id, DeviceID: "device", Address: addresses[i], Function: "03", DataType: schema.DataTypeInt16})
			}
			values := scheduler.PollNowContext(t.Context(), []string{"c", "a", "b", "a", "missing"})
			for _, value := range values {
				require.Equal(t, schema.QualityGood, value.Quality, value.Error)
			}
			require.Len(t, values, 3, "duplicate request ID must have only one final result")
			require.Len(t, peer.log(), wantRequests, "manual must use the existing immutable plan, not scalar loop")
			got := map[string]any{}
			for _, value := range values {
				require.Equal(t, schema.QualityGood, value.Quality, value.Error)
				require.NotEmpty(t, value.RequestID)
				require.NotEmpty(t, value.PlanFingerprint)
				got[value.PointID] = value.Value
			}
			require.Equal(t, map[string]any{"a": int16(11), "b": int16(22), "c": int16(33)}, got)
			require.Len(t, observations, 3)
			for _, request := range peer.log() {
				require.NotContains(t, request.Addresses, 20)
			}
			select {
			case <-scheduler.ValueChannel():
				t.Fatal("manual return must not create duplicate scheduled delivery")
			default:
			}
			require.Empty(t, scheduler.PollNowContext(t.Context(), nil))
			require.Empty(t, scheduler.PollNowContext(t.Context(), []string{"missing"}))
		})
	}
}

func TestManualPlannedOverlappingPointsKeepIndependentFormats(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolMC3E, schema.ProtocolFatekFBs, schema.ProtocolModbusTCP} {
		t.Run(string(protocol), func(t *testing.T) {
			peer, config := newAcquisitionWirePeer(t, protocol, map[int]uint16{0: 0x1122, 1: 0x3344})
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
			settings := collector.DefaultSchedulerConfig()
			settings.DefaultRetryCount = 0
			scheduler := collector.NewScheduler(settings, manager)
			scheduler.AddDevice(&schema.Device{ID: "device", Protocol: protocol, ConnectionConfig: config})
			address := "D0"
			if protocol == schema.ProtocolModbusTCP {
				address = "40001"
			}
			for _, format := range []string{"ABCD", "CDAB"} {
				point := &schema.Point{ID: format, DeviceID: "device", Address: address, Function: "03", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatVersion: 1, FormatOverride: format, RevisionID: "point-v1"}
				profile, err := connector.ResolveReadProfile(*point, decoding.FormatContext{Protocol: protocol, ConnectionFormat: "ABCD", ConnectionRevision: connector.ReadConnectionRevision(config)})
				require.NoError(t, err)
				require.NotZero(t, scheduler.AddPointWithObservationLease(point, profile))
			}
			values := scheduler.PollNowContext(t.Context(), []string{"ABCD", "CDAB"})
			for _, value := range values {
				require.Equal(t, schema.QualityGood, value.Quality, value.Error)
			}
			require.Len(t, values, 2)
			require.Len(t, peer.log(), 1, "one raw span must fan out to two canonical typed Points")
			got := map[string]any{}
			for _, value := range values {
				require.Equal(t, schema.QualityGood, value.Quality, value.Error)
				got[value.PointID] = value.Value
				require.Equal(t, values[0].RequestID, value.RequestID)
				require.Equal(t, values[0].ReceivedAt, value.ReceivedAt)
			}
			require.Equal(t, map[string]any{"ABCD": uint32(0x11223344), "CDAB": uint32(0x33441122)}, got)
			select {
			case <-scheduler.ValueChannel():
				t.Fatal("manual must not emit into scheduled channel")
			default:
			}
		})
	}
}

func TestManualPlannedMultipleGroupsKeepOneDeviceAcquisition(t *testing.T) {
	peer, config := newAcquisitionWirePeer(t, schema.ProtocolMC3E, map[int]uint16{0: 11, 1: 22})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	settings := collector.DefaultSchedulerConfig()
	settings.DefaultRetryCount = 0
	nextID := 0
	settings.AcquisitionIDFactory = func() string { nextID++; return "device-intent" }
	scheduler := collector.NewScheduler(settings, manager)
	scheduler.AddDevice(&schema.Device{ID: "device", Protocol: schema.ProtocolMC3E, ConnectionConfig: config})
	for _, item := range []struct{ id, address, group string }{{"a", "D0", "group-a"}, {"b", "D1", "group-b"}} {
		scheduler.AddPoint(&schema.Point{ID: item.id, DeviceID: "device", Address: item.address, DataType: schema.DataTypeInt16, PollingGroupID: &item.group})
	}
	values := scheduler.PollNowContext(t.Context(), []string{"a", "b"})
	require.Len(t, values, 2)
	require.Equal(t, 1, nextID, "manual caller still creates one device acquisition, not one per original group")
	require.Len(t, peer.log(), 2, "original groups must remain separate acquisition requests")
	require.Equal(t, "device-intent", values[0].AcquisitionID)
	require.Equal(t, values[0].AcquisitionID, values[1].AcquisitionID)
	require.NotEmpty(t, values[0].RequestID)
	require.NotEmpty(t, values[1].RequestID)
	require.NotEqual(t, values[0].RequestID, values[1].RequestID, "request index cannot reset inside one device acquisition")
}
