package adapters

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/wirebudget"
)

func TestFatekPlannedMalformedWholeResponseFailsEveryPoint(t *testing.T) {
	for _, fault := range []string{"short", "extra", "nonbinary", "partialhex", "checksum", "status", "status-extra-9", "status-extra-4", "command", "station"} {
		t.Run(fault, func(t *testing.T) {
			req := fatekPlannedRequest(t, fatekPlannedPoint("word", "D0", schema.DataTypeUint16, ""), fatekPlannedPoint("bit", "X20", schema.DataTypeBool, ""))
			peer := &fatekPlannedPeer{memory: map[string]any{"D0": 0x1234, "X20": true}, corrupt: func(frame []byte) []byte {
				payload := string(frame[6 : len(frame)-3])
				status, command, station := "0", "48", "01"
				switch fault {
				case "short":
					payload = payload[:len(payload)-1]
				case "extra":
					payload += "0"
				case "nonbinary":
					payload = "12342"
				case "partialhex":
					payload = "123Z1"
				case "checksum":
					frame[len(frame)-2] = '!'
					return frame
				case "status":
					status = "9"
					payload = ""
				case "status-extra-9":
					status, payload = "9", "garbage"
				case "status-extra-4":
					status, payload = "4", "0000"
				case "command":
					command = "46"
				case "station":
					station = "02"
				}
				return fatekPeerFrame(station, command, status+payload)
			}}
			ctx := wirebudget.With(t.Context(), 3, 0)
			got, err := readFatekPlanned(t, ctx, fatekPlannedConnector(peer), req)
			require.Error(t, err)
			require.Len(t, got, 2)
			require.Len(t, peer.frames, 1)
			require.Equal(t, 1, wirebudget.From(ctx).Used(), "no hidden scalar fallback/retry")
			var pe *fatek.FatekProtocolError
			require.Equal(t, fault == "status", errors.As(err, &pe), "malformed frames cannot authorize capability fallback")
			for _, result := range got {
				require.Equal(t, schema.QualityBad, result.Quality)
				require.Nil(t, result.Value)
				require.Nil(t, result.RawBytes)
				require.NotEmpty(t, result.Error)
				require.Equal(t, req.AcquisitionID, result.AcquisitionID)
				require.Equal(t, "opaque-source-proof", result.Binding.Point.SourceRevision)
				require.Equal(t, uint64(19), result.Binding.InstallationLease)
				require.Equal(t, result.Timestamp, result.ObservedAt)
				require.Equal(t, result.Timestamp, result.ReceivedAt)
			}
		})
	}
}
func TestFatekPlannedMissingMapEntryDoesNotBecomeGoodZero(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("wide", "D0", schema.DataTypeUint32, "ABCD"), fatekPlannedPoint("far", "D20", schema.DataTypeUint16, ""), fatekPlannedPoint("bit", "X20", schema.DataTypeBool, ""))
	request, bindings, err := req.Snapshot()
	require.NoError(t, err)
	for _, tc := range []struct {
		name   string
		values map[string]any
		badIDs map[string]bool
	}{
		{"missing word", map[string]any{"D0": 0x0102, "D20": 0, "X20": false}, map[string]bool{"wide": true}},
		{"wrong word type", map[string]any{"D0": true, "D1": 0x0304, "D20": 0, "X20": false}, map[string]bool{"wide": true}},
		{"word out of range", map[string]any{"D0": 65536, "D1": 0x0304, "D20": 0, "X20": false}, map[string]bool{"wide": true}},
		{"missing bit", map[string]any{"D0": 0x0102, "D1": 0x0304, "D20": 0}, map[string]bool{"bit": true}},
		{"wrong bit type", map[string]any{"D0": 0x0102, "D1": 0x0304, "D20": 0, "X20": 0}, map[string]bool{"bit": true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := fatekPlannedResults(req, request, bindings, fatekPlannedRaw{random: tc.values}, time.Now(), nil)
			for _, result := range got {
				id := result.Binding.Point.ID
				if tc.badIDs[id] {
					require.Equal(t, schema.QualityBad, result.Quality, id)
					require.Nil(t, result.Value, id)
					require.Nil(t, result.RawBytes, id)
					require.NotEmpty(t, result.Error, id)
				} else {
					require.Equal(t, schema.QualityGood, result.Quality, id)
					want := map[string]any{"wide": uint32(0x01020304), "far": uint16(0), "bit": false}
					require.Equal(t, want[id], result.Value, id)
				}
			}
		})
	}
}
func TestFatekPlannedPeerRejectsWrongOrUnmappedWireAddresses(t *testing.T) {
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 7}}
	for _, body := range []string{"01D00001", "01D000000", "01WX00012", "00", "02D00000"} {
		_, err := peer.SendReceive(fatekPeerFrame("01", "48", body))
		require.Error(t, err, body)
	}
	// A legitimate zero is available only if explicitly present in device memory.
	peer.memory["D0"] = 0
	got, err := peer.SendReceive(fatekPeerFrame("01", "48", "01D00000"))
	require.NoError(t, err)
	require.Equal(t, "0000", string(got[6:len(got)-3]))
}
func TestFatekPlannedReadKeepsFrozenMetadataAtCompletion(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "D0", schema.DataTypeUint32, "ABCD"), fatekPlannedPoint("b", "D10", schema.DataTypeUint16, ""))
	original := req.Bindings[0]
	entered := time.Time{}
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 0x0102, "D1": 0x0304, "D10": 7}, hook: func() {
		entered = time.Now()
		req.Bindings[0].PointRevision = "new-point-revision"
		req.Bindings[0].InstallationLease = 99
		req.Bindings[0].Point.SourceRevision = "new-source"
	}}
	got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
	require.NoError(t, err)
	found := false
	for _, r := range got {
		require.False(t, r.Timestamp.Before(entered))
		if r.Binding.Point.ID == original.Point.ID {
			require.Equal(t, original, r.Binding)
			found = true
		}
	}
	require.True(t, found)
}
func TestFatekPlannedDecodeFailureIsOnlyAffectedPoint(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("word", "D0", schema.DataTypeUint16, ""), fatekPlannedPoint("other", "D10", schema.DataTypeUint16, ""))
	r, bindings, err := req.Snapshot()
	require.NoError(t, err)
	r.Points[0].Offset = 99
	got := fatekPlannedResults(req, r, bindings, fatekPlannedRaw{random: map[string]any{"D0": 7, "D10": 8}}, time.Now(), nil)
	require.Equal(t, schema.QualityBad, got[0].Quality)
	require.Nil(t, got[0].Value)
	require.Equal(t, schema.QualityGood, got[1].Quality)
	require.Equal(t, uint16(8), got[1].Value)
}
