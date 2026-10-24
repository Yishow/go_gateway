package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func longGroupFixture(kind WriteGroupProfileKind, mode WriteGroupStorageMode, measurements ...bool) *WriteGroup {
	group := &WriteGroup{
		ID: "group-L", LayoutVersion: WriteGroupLayoutManagedLongV1, StorageMode: mode,
		RecordingProfile: &WriteGroupRecordingProfile{Kind: kind, Interval: 15, Unit: WriteGroupIntervalUnitSecond},
	}
	for i, hasMeasurement := range measurements {
		member := WriteGroupMember{DeviceID: "device", PointID: "point", TagID: string(rune('a' + i))}
		if hasMeasurement {
			id := "meas-" + member.TagID
			member.MeasurementID = &id
		}
		group.Members = append(group.Members, member)
	}
	return group
}

func finalizedLong(group *WriteGroup) *WriteGroupReadiness {
	readiness := newWriteGroupReadiness(&Record{ID: "workspace"}, group)
	readiness.ConfigReady, readiness.SchemaReady = true, true
	return finalizeWriteGroupReadiness(readiness)
}

func issueCodesOf(readiness *WriteGroupReadiness) []string {
	codes := make([]string, 0, len(readiness.Issues))
	for _, issue := range readiness.Issues {
		codes = append(codes, issue.Code)
	}
	return codes
}

// TestLongRuntimeGateOpensOnlyForTheSupportedProfiles pins what the runtime
// really writes: plain or explicitly selected timescale storage with an
// instantaneous or cumulative profile. Anything else stays blocked, and a cumulative group also needs a saved
// measurement for every source before it can be called ready.
func TestLongRuntimeGateOpensOnlyForTheSupportedProfiles(t *testing.T) {
	instantaneous := finalizedLong(longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false, false))
	require.True(t, instantaneous.Ready, "%v", issueCodesOf(instantaneous))
	require.Empty(t, instantaneous.Issues)

	cumulative := finalizedLong(longGroupFixture(WriteGroupProfileCumulative, WriteGroupStorageModePlain, true, true))
	require.True(t, cumulative.Ready, "%v", issueCodesOf(cumulative))

	unconfirmed := finalizedLong(longGroupFixture(WriteGroupProfileCumulative, WriteGroupStorageModePlain, true, false))
	require.False(t, unconfirmed.Ready)
	require.Equal(t, []string{"counter-semantics-required"}, issueCodesOf(unconfirmed),
		"a cumulative source without a saved measurement cannot be counted")

	timescale := finalizedLong(longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModeTimescale, false))
	require.True(t, timescale.Ready, "an explicitly selected timescale group is written by the same runtime: %v", issueCodesOf(timescale))
	timescaleCumulative := finalizedLong(longGroupFixture(WriteGroupProfileCumulative, WriteGroupStorageModeTimescale, true))
	require.True(t, timescaleCumulative.Ready, "%v", issueCodesOf(timescaleCumulative))

	for name, group := range map[string]*WriteGroup{
		"an unknown storage mode": longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageMode("native"), false),
		"no storage mode":         longGroupFixture(WriteGroupProfileInstantaneous, "", false),
		"an unknown profile":      longGroupFixture(WriteGroupProfileKind("average"), WriteGroupStorageModePlain, false),
		"no profile": func() *WriteGroup {
			g := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
			g.RecordingProfile = nil
			return g
		}(),
	} {
		readiness := finalizedLong(group)
		require.False(t, readiness.Ready, name)
		require.Contains(t, issueCodesOf(readiness), "long-runtime-unavailable", name)
	}

	wide := finalizedLong(&WriteGroup{ID: "wide"})
	require.True(t, wide.Ready, "a wide group never carried the long gate")
}
