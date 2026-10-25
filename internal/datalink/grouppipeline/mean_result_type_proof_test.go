package grouppipeline

import (
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// meanLongShape turns the common long proof into a fully verified mean gauge
// source: an int16 tag whose exact mean is a decimal.
func meanLongShape(mutate func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout)) func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	return func(group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) {
		group.RecordingProfile.Kind = workspace.WriteGroupProfileInstantaneous
		group.RecordingProfile.Interval = 15
		group.RecordingProfile.Unit = workspace.WriteGroupIntervalUnitSecond
		group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
		group.RowPolicy.IntervalSeconds = 15
		group.Members[0].ValueRole = workspace.WriteGroupValueRoleGauge
		group.Members[0].SourceValueType = schema.DataTypeInt16
		group.Members[0].ResultValueType = measurement.ExactDecimal
		layout.TagTypes["tag-1"] = schema.DataTypeInt16
		if mutate != nil {
			mutate(group, layout)
		}
	}
}

// TestMeanResultTypeProof pins the frozen proof of the mean long shape: the
// source type is the tag's own schema.DataType, the result type is the exact
// type the mean writer stores, and a partial or contradictory descriptor never
// proves a writer. New wide, unknown reduction and unverified metadata stay
// rejected.
func TestMeanResultTypeProof(t *testing.T) {
	t.Run("a fully verified mean gauge proves a writer", func(t *testing.T) {
		group, layout := longProof(t, workspace.WriteGroupStorageModePlain, meanLongShape(nil))
		require.True(t, validLongRuntimeProof(group, layout),
			"int16 source, decimal result, gauge role and mean reduction are complete frozen facts")
		require.Equal(t, schema.DataTypeInt16, layout.TagTypes["tag-1"], "the frozen source Tag type is never rewritten")
	})

	t.Run("a counter_raw member is structurally frozen beside its confirmed id", func(t *testing.T) {
		id := "meas-1"
		group, layout := longProof(t, workspace.WriteGroupStorageModePlain, meanLongShape(func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole = workspace.WriteGroupValueRoleCounterRaw
			g.Members[0].ResultValueType = measurement.ExactInt64
			g.Members[0].MeasurementID = &id
		}))
		require.True(t, validLongRuntimeProof(group, layout),
			"structure proves the id; the confirmed meaning is checked by readiness and the counter provider")
	})

	t.Run("a state member keeps its latest typed contract", func(t *testing.T) {
		group, layout := longProof(t, workspace.WriteGroupStorageModePlain, meanLongShape(func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			g.Members = append(g.Members, workspace.WriteGroupMember{
				DeviceID: "device-1", PointID: "point-2", TagID: "tag-state", TargetColumn: "s",
				ValueRole: workspace.WriteGroupValueRoleState, SourceValueType: schema.DataTypeString, ResultValueType: measurement.ExactText,
			})
			l.TagTypes["tag-state"] = schema.DataTypeString
		}))
		require.True(t, validLongRuntimeProof(group, layout), "a state member is latest text, not part of the mean")
	})

	negatives := map[string]func(*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout){
		"an omitted role": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[0].ValueRole = "" },
		"an unknown role": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole = workspace.WriteGroupValueRole("average")
		},
		"an omitted source type": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[0].SourceValueType = "" },
		"a source type that is not the tag's": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].SourceValueType = schema.DataTypeInt32
		},
		"an omitted result type": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) { g.Members[0].ResultValueType = "" },
		"a result type the mean cannot store": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ResultValueType = measurement.ExactFloat64
		},
		"a mean on a cumulative profile": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
		},
		"a gauge mean on a text tag": func(g *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].SourceValueType = schema.DataTypeString
			g.Members[0].ResultValueType = measurement.ExactText
			l.TagTypes["tag-1"] = schema.DataTypeString
		},
		"an unknown reduction": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReduction("average")
		},
		"a mean gauge whose frozen tag type is missing": func(_ *workspace.WriteGroup, l *workspace.WriteGroupRuntimeLayout) {
			delete(l.TagTypes, "tag-1")
		},
		"counter_raw without a saved measurement id": func(g *workspace.WriteGroup, _ *workspace.WriteGroupRuntimeLayout) {
			g.Members[0].ValueRole = workspace.WriteGroupValueRoleCounterRaw
			g.Members[0].ResultValueType = measurement.ExactInt64
			g.Members[0].MeasurementID = nil
		},
	}
	for name, mutate := range negatives {
		t.Run("a forged mean proof is refused: "+name, func(t *testing.T) {
			group, layout := longProof(t, workspace.WriteGroupStorageModePlain, meanLongShape(mutate))
			require.False(t, validLongRuntimeProof(group, layout))
		})
	}
}
