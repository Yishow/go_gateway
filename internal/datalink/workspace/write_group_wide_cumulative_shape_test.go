package workspace

import (
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func cumulativeWideGroup(mutate func(*WriteGroup)) (group *WriteGroup, tagTypes map[string]schema.DataType) {
	water, gas := "meas-water", "meas-gas"
	group = &WriteGroup{
		ID: "group-W", LayoutVersion: WriteGroupLayoutManagedWideProfileV1,
		RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 10, Unit: WriteGroupIntervalUnitMinute},
		Members: []WriteGroupMember{
			{
				DeviceID: "device", PointID: "p-water", TagID: "tag-water", MeasurementID: &water,
				SourceValueType: schema.DataTypeFloat64, ResultValueType: measurement.ExactFloat64,
			},
			{
				DeviceID: "device", PointID: "p-gas", TagID: "tag-gas", MeasurementID: &gas,
				SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactInt64,
			},
		},
	}
	if mutate != nil {
		mutate(group)
	}
	return group, map[string]schema.DataType{"tag-water": schema.DataTypeFloat64, "tag-gas": schema.DataTypeInt64}
}

// TestWideCumulativeShapeNeedsConfirmedRawMembers pins what a cumulative wide
// profile may save and prepare: raw numeric members, each with the saved
// measurement that is to confirm it, no role and no reduction. A gauge mean can
// never be a counter baseline, and an id is only structure: readiness and the
// runtime provider still prove the counter semantics before anything writes.
func TestWideCumulativeShapeNeedsConfirmedRawMembers(t *testing.T) {
	group, types := cumulativeWideGroup(nil)
	require.True(t, WideProfileShapeSupported(group, types))

	blank := "  "
	for name, mutate := range map[string]func(*WriteGroup){
		"a mean as a counter baseline": func(g *WriteGroup) { g.RecordingProfile.SampleReduction = WriteGroupSampleReductionMean },
		"an instantaneous role":        func(g *WriteGroup) { g.Members[0].ValueRole = WriteGroupValueRoleGauge },
		"a member without a measurement": func(g *WriteGroup) {
			g.Members[1].MeasurementID = nil
		},
		"a blank measurement id": func(g *WriteGroup) { g.Members[1].MeasurementID = &blank },
		"a result type that is not the raw type": func(g *WriteGroup) {
			g.Members[1].ResultValueType = measurement.ExactDecimal
		},
		"an omitted source type": func(g *WriteGroup) { g.Members[0].SourceValueType = "" },
		"a source type that is not the tag's": func(g *WriteGroup) {
			g.Members[0].SourceValueType = schema.DataTypeFloat32
		},
		"a storage mode":       func(g *WriteGroup) { g.StorageMode = WriteGroupStorageModeTimescale },
		"no members":           func(g *WriteGroup) { g.Members = nil },
		"a legacy wide layout": func(g *WriteGroup) { g.LayoutVersion = WriteGroupLayoutManagedWideV1 },
	} {
		forged, forgedTypes := cumulativeWideGroup(mutate)
		require.False(t, WideProfileShapeSupported(forged, forgedTypes), name)
	}

	text, textTypes := cumulativeWideGroup(func(g *WriteGroup) {
		g.Members[0].SourceValueType, g.Members[0].ResultValueType = schema.DataTypeString, measurement.ExactText
	})
	textTypes["tag-water"] = schema.DataTypeString
	require.False(t, WideProfileShapeSupported(text, textTypes), "a text tag has no counter")

	missing, missingTypes := cumulativeWideGroup(nil)
	delete(missingTypes, "tag-gas")
	require.False(t, WideProfileShapeSupported(missing, missingTypes), "a tag whose type is unknown is refused")
}

// TestWideCumulativeReadinessWaitsForCounterSemantics pins that a saved
// cumulative wide group is not ready on its structure alone: a member without
// its saved measurement is blocked for counter semantics, as in the long layout.
func TestWideCumulativeReadinessWaitsForCounterSemantics(t *testing.T) {
	confirmed, _ := cumulativeWideGroup(nil)
	require.NotContains(t, issueCodesOf(finalizedLong(confirmed)), "counter-semantics-required")

	unconfirmed, _ := cumulativeWideGroup(func(g *WriteGroup) { g.Members[1].MeasurementID = nil })
	readiness := finalizedLong(unconfirmed)
	require.False(t, readiness.Ready)
	require.Contains(t, issueCodesOf(readiness), "counter-semantics-required")
}
