package grouppipeline

import (
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// wideProfileProof is the frozen descriptor an Apply of a managed wide profile
// mean group leaves behind: a float and an integer gauge, the integer stored as
// an exact decimal column, every member naming its role and both types.
func wideProfileProof(t *testing.T, mutate func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout)) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	group, layout := canonicalProof(t)
	group.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
	injectMeanReduction(t, group)
	group.Members[0].ValueRole, group.Members[0].SourceValueType, group.Members[0].ResultValueType =
		workspace.WriteGroupValueRoleGauge, schema.DataTypeFloat64, measurement.ExactFloat64
	group.Members[0].Required = false
	group.Members[1].ValueRole, group.Members[1].SourceValueType, group.Members[1].ResultValueType =
		workspace.WriteGroupValueRoleGauge, schema.DataTypeInt64, measurement.ExactDecimal
	group.Members[1].Required = false
	layout.Columns[5] = dbtarget.ColumnInfo{Name: "pressure", DataType: "TEXT", Nullable: true}
	if mutate != nil {
		mutate(group, layout)
	}
	return group, layout
}

// TestWideProfileResultTypeProof pins the frozen proof of the wide profile
// layout: the source type is the Tag's own schema.DataType, the result type is
// what the column holds, and a partial, contradictory or unservable descriptor
// never proves a writer. The legacy wide proof keeps refusing a descriptor with
// the new fields.
func TestWideProfileResultTypeProof(t *testing.T) {
	t.Run("a verified mean descriptor proves a writer", func(t *testing.T) {
		group, layout := wideProfileProof(t, nil)
		require.True(t, validCanonicalManagedRuntimeProof(group, layout))
		require.Equal(t, schema.DataTypeInt64, layout.TagTypes["tag-p"], "the frozen source Tag type is never rewritten")
	})

	t.Run("state and counter_raw members keep their latest typed contract", func(t *testing.T) {
		id := "meas-1"
		group, layout := wideProfileProof(t, func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole, g.Members[0].ResultValueType = workspace.WriteGroupValueRoleState, measurement.ExactFloat64
			g.Members[1].ValueRole, g.Members[1].ResultValueType, g.Members[1].MeasurementID =
				workspace.WriteGroupValueRoleCounterRaw, measurement.ExactInt64, &id
			l.Columns[5] = dbtarget.ColumnInfo{Name: "pressure", DataType: "INTEGER", Nullable: true}
		})
		require.True(t, validCanonicalManagedRuntimeProof(group, layout))
	})

	negatives := map[string]func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout){
		"an omitted role": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[0].ValueRole = "" },
		"an unknown role": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole = "average"
		},
		"an omitted source type": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[1].SourceValueType = "" },
		"a source type that is not the tag's": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[1].SourceValueType = schema.DataTypeInt32
		},
		"an omitted result type": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[1].ResultValueType = "" },
		"a result type the mean cannot store": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[1].ResultValueType = measurement.ExactFloat64
		},
		"a member column that cannot hold NULL": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Columns[4].Nullable = false
		},
		"a missing member column": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Columns = append(l.Columns[:4], l.Columns[5:]...)
		},
		"a mean on a cumulative profile": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
		},
		"a storage mode": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.StorageMode = workspace.WriteGroupStorageModePlain
		},
		"an unknown reduction": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReduction("average")
		},
		"a frozen tag type that is missing": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			delete(l.TagTypes, "tag-p")
		},
		"counter_raw without a saved measurement id": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[1].ValueRole, g.Members[1].ResultValueType, g.Members[1].MeasurementID =
				workspace.WriteGroupValueRoleCounterRaw, measurement.ExactInt64, nil
		},
		"a legacy member without any metadata": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole, g.Members[0].SourceValueType, g.Members[0].ResultValueType = "", "", ""
		},
	}
	for name, mutate := range negatives {
		t.Run("a forged wide profile proof is refused: "+name, func(t *testing.T) {
			group, layout := wideProfileProof(t, mutate)
			require.False(t, validCanonicalManagedRuntimeProof(group, layout))
		})
	}

	t.Run("a legacy wide descriptor with the new fields is still refused", func(t *testing.T) {
		group, layout := canonicalProof(t)
		require.True(t, validCanonicalManagedRuntimeProof(group, layout))
		group.Members[0].ValueRole = workspace.WriteGroupValueRoleGauge
		require.False(t, validCanonicalManagedRuntimeProof(group, layout))
	})
}

// wideCumulativeProof is the frozen descriptor an Apply of a cumulative wide
// profile group leaves behind: a float and an integer counter keeping their raw
// types, each naming the saved measurement that confirms it, with no role.
func wideCumulativeProof(t *testing.T, mutate func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout)) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	group, layout := wideProfileProof(t, nil)
	group.RecordingProfile.Kind, group.RecordingProfile.SampleReduction = workspace.WriteGroupProfileCumulative, ""
	water, gas := "meas-water", "meas-gas"
	group.Members[0].ValueRole, group.Members[0].SourceValueType, group.Members[0].ResultValueType, group.Members[0].MeasurementID =
		"", schema.DataTypeFloat64, measurement.ExactFloat64, &water
	group.Members[1].ValueRole, group.Members[1].SourceValueType, group.Members[1].ResultValueType, group.Members[1].MeasurementID =
		"", schema.DataTypeInt64, measurement.ExactInt64, &gas
	layout.Columns[5] = dbtarget.ColumnInfo{Name: "pressure", DataType: "INTEGER", Nullable: true}
	if mutate != nil {
		mutate(group, layout)
	}
	return group, layout
}

// TestWideCumulativeProof pins the frozen proof of a cumulative wide group: raw
// numeric types, a saved measurement per member, no role and no reduction. A mean
// is never a counter baseline, and an id alone proves nothing the provider has not
// confirmed; every forged or partial descriptor is refused.
func TestWideCumulativeProof(t *testing.T) {
	group, layout := wideCumulativeProof(t, nil)
	require.True(t, validCanonicalManagedRuntimeProof(group, layout))
	require.Equal(t, schema.DataTypeInt64, layout.TagTypes["tag-p"], "the frozen source Tag type is never rewritten")

	blank := "  "
	negatives := map[string]func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout){
		"a mean as a counter baseline": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
		},
		"an instantaneous role": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole = workspace.WriteGroupValueRoleGauge
		},
		"a member without a saved measurement": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[1].MeasurementID = nil },
		"a blank measurement id": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[1].MeasurementID = &blank
		},
		"a result type that is not the raw type": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[1].ResultValueType = measurement.ExactDecimal
		},
		"an omitted source type": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[0].SourceValueType = "" },
		"a source type that is not the tag's": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].SourceValueType = schema.DataTypeFloat32
		},
		"a text tag": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.TagTypes["tag-t"] = schema.DataTypeString
			g.Members[0].SourceValueType, g.Members[0].ResultValueType = schema.DataTypeString, measurement.ExactText
		},
		"a frozen tag type that is missing": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) { delete(l.TagTypes, "tag-p") },
		"a member column that cannot hold NULL": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			l.Columns[4].Nullable = false
		},
		"a storage mode": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.StorageMode = workspace.WriteGroupStorageModePlain
		},
	}
	for name, mutate := range negatives {
		t.Run("a forged cumulative wide proof is refused: "+name, func(t *testing.T) {
			forged, forgedLayout := wideCumulativeProof(t, mutate)
			require.False(t, validCanonicalManagedRuntimeProof(forged, forgedLayout))
		})
	}
}
