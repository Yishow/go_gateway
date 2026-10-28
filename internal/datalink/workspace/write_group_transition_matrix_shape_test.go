package workspace

import (
	"maps"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupTransitionMatrixEquivalentInstantaneousShapes(t *testing.T) {
	for _, tc := range []struct {
		name       string
		role       WriteGroupValueRole
		reduction  WriteGroupSampleReduction
		sourceType schema.DataType
		resultType measurement.ExactType
	}{
		{"text state latest", WriteGroupValueRoleState, "latest", schema.DataTypeString, measurement.ExactText},
		{"bool state mean profile", WriteGroupValueRoleState, "mean", schema.DataTypeBool, measurement.ExactBool},
		{"gauge latest", WriteGroupValueRoleGauge, "latest", schema.DataTypeUint64, measurement.ExactUint64},
		{"integer mean", WriteGroupValueRoleGauge, "mean", schema.DataTypeInt64, measurement.ExactDecimal},
		{"float mean", WriteGroupValueRoleGauge, "mean", schema.DataTypeFloat64, measurement.ExactFloat64},
		{"raw counter latest", WriteGroupValueRoleCounterRaw, "latest", schema.DataTypeUint64, measurement.ExactUint64},
		{"raw counter mean profile", WriteGroupValueRoleCounterRaw, "mean", schema.DataTypeInt64, measurement.ExactInt64},
	} {
		for _, engine := range []string{"sqlite", "postgres"} {
			t.Run(tc.name+"/"+engine, func(t *testing.T) {
				source, target, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1)
				for _, g := range []*WriteGroup{source, target} {
					g.RecordingProfile.SampleReduction = tc.reduction
					g.Members[0].ValueRole, g.Members[0].SourceValueType, g.Members[0].ResultValueType = tc.role, tc.sourceType, tc.resultType
					if tc.role == WriteGroupValueRoleCounterRaw {
						id := "counter-1"
						g.Members[0].MeasurementID = &id
					}
				}
				types["tag-1"] = tc.sourceType
				require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, engine, engine, types))
				require.NoError(t, ValidateWriteGroupTransitionStructure(target, source, engine, engine, types))
			})
		}
	}
}

func transitionMatrixCumulativePair() (*WriteGroup, *WriteGroup, map[string]schema.DataType) {
	long, wide, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1)
	for _, g := range []*WriteGroup{long, wide} {
		g.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileCumulative, Interval: 10, Unit: WriteGroupIntervalUnitMinute}
		id := "counter-1"
		g.Members[0].ValueRole, g.Members[0].MeasurementID = "", &id
	}
	// Canonical unscaled long has no frozen raw metadata; wide requires it.
	long.Members[0].SourceValueType, long.Members[0].ResultValueType = "", ""
	wide.Members[0].SourceValueType, wide.Members[0].ResultValueType = schema.DataTypeUint64, measurement.ExactUint64
	types["tag-1"] = schema.DataTypeUint64
	return long, wide, types
}

func TestWriteGroupTransitionMatrixCumulativeEquivalentRawTypes(t *testing.T) {
	for _, engine := range []string{"sqlite", "postgres"} {
		long, wide, types := transitionMatrixCumulativePair()
		beforeLong, beforeWide := cloneWriteGroup(long), cloneWriteGroup(wide)
		require.NoError(t, ValidateWriteGroupTransitionStructure(long, wide, engine, engine, types))
		require.NoError(t, ValidateWriteGroupTransitionStructure(wide, long, engine, engine, types))
		require.Equal(t, beforeLong, long)
		require.Equal(t, beforeWide, wide)
	}
}

func TestWriteGroupTransitionMatrixRejectsUnsupportedShapes(t *testing.T) {
	for name, mutate := range map[string]func(*WriteGroup, map[string]schema.DataType){
		"unknown profile kind": func(g *WriteGroup, _ map[string]schema.DataType) { g.RecordingProfile.Kind = "private://secret" },
		"unknown reduction": func(g *WriteGroup, _ map[string]schema.DataType) {
			g.RecordingProfile.SampleReduction = "private://secret"
		},
		"invalid interval":      func(g *WriteGroup, _ map[string]schema.DataType) { g.RecordingProfile.Interval = 0 },
		"unknown unit":          func(g *WriteGroup, _ map[string]schema.DataType) { g.RecordingProfile.Unit = "private://secret" },
		"missing type evidence": func(_ *WriteGroup, types map[string]schema.DataType) { delete(types, "tag-1") },
		"unknown type evidence": func(_ *WriteGroup, types map[string]schema.DataType) { types["tag-1"] = "private://secret" },
		"source type differs from evidence": func(g *WriteGroup, _ map[string]schema.DataType) {
			g.Members[0].SourceValueType = schema.DataTypeString
		},
		"wrong exact result": func(g *WriteGroup, _ map[string]schema.DataType) {
			g.Members[0].ResultValueType = measurement.ExactDecimal
		},
		"missing frozen source type": func(g *WriteGroup, _ map[string]schema.DataType) { g.Members[0].SourceValueType = "" },
		"missing frozen result type": func(g *WriteGroup, _ map[string]schema.DataType) { g.Members[0].ResultValueType = "" },
		"unknown role":               func(g *WriteGroup, _ map[string]schema.DataType) { g.Members[0].ValueRole = "private://secret" },
		"missing wide role":          func(g *WriteGroup, _ map[string]schema.DataType) { g.Members[0].ValueRole = "" },
		"boolean mean gauge": func(g *WriteGroup, _ map[string]schema.DataType) {
			g.RecordingProfile.SampleReduction = "mean"
			g.Members[0].ValueRole = WriteGroupValueRoleGauge
		},
		"boolean raw counter": func(g *WriteGroup, _ map[string]schema.DataType) {
			id := "counter-1"
			g.Members[0].ValueRole = WriteGroupValueRoleCounterRaw
			g.Members[0].MeasurementID = &id
		},
		"missing counter measurement": func(g *WriteGroup, types map[string]schema.DataType) {
			g.Members[0].ValueRole = WriteGroupValueRoleCounterRaw
			g.Members[0].SourceValueType = schema.DataTypeInt64
			g.Members[0].ResultValueType = measurement.ExactInt64
			types["tag-1"] = schema.DataTypeInt64
		},
		"blank counter measurement": func(g *WriteGroup, types map[string]schema.DataType) {
			id := " "
			g.Members[0].ValueRole = WriteGroupValueRoleCounterRaw
			g.Members[0].SourceValueType = schema.DataTypeInt64
			g.Members[0].ResultValueType = measurement.ExactInt64
			g.Members[0].MeasurementID = &id
			types["tag-1"] = schema.DataTypeInt64
		},
	} {
		t.Run(name, func(t *testing.T) {
			source, target, types := transitionMatrixPair(WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedWideProfileV1)
			mutate(source, types)
			mutate(target, types)
			beforeSource, beforeTarget, beforeTypes := cloneWriteGroup(source), cloneWriteGroup(target), maps.Clone(types)
			requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), "unsupported-shape")
			require.Equal(t, beforeSource, source)
			require.Equal(t, beforeTarget, target)
			require.Equal(t, beforeTypes, types)
		})
	}
}

func TestWriteGroupTransitionMatrixRejectsUnconfirmedCumulativeShapes(t *testing.T) {
	for name, mutate := range map[string]func(*WriteGroup, map[string]schema.DataType){
		"mean cumulative":           func(g *WriteGroup, _ map[string]schema.DataType) { g.RecordingProfile.SampleReduction = "mean" },
		"missing measurement":       func(g *WriteGroup, _ map[string]schema.DataType) { g.Members[0].MeasurementID = nil },
		"non numeric counter":       func(_ *WriteGroup, types map[string]schema.DataType) { types["tag-1"] = schema.DataTypeString },
		"missing raw type evidence": func(_ *WriteGroup, types map[string]schema.DataType) { delete(types, "tag-1") },
	} {
		t.Run(name, func(t *testing.T) {
			source, _, types := transitionMatrixCumulativePair()
			target := cloneWriteGroup(source)
			mutate(source, types)
			mutate(target, types)
			requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), "unsupported-shape")
		})
	}
}

func TestWriteGroupTransitionMatrixSourceOrderIsNotIdentity(t *testing.T) {
	source, target, types := transitionMatrixPair(WriteGroupLayoutManagedLongV1, WriteGroupLayoutManagedWideProfileV1)
	second := source.Members[0]
	second.PointID, second.TagID, second.TargetColumn = "point-2", "tag-2", "state_2"
	source.Members = append(source.Members, second)
	target.Members = append([]WriteGroupMember{second}, target.Members...)
	types["tag-2"] = schema.DataTypeBool
	require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types))
}
