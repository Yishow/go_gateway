package runtime

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mixedCounterBoundaryFixture(t *testing.T, scaled bool) *longFixture {
	t.Helper()
	f := newLongFixture(t)
	g := f.config.Group
	g.Members = g.Members[4:]
	f.members = g.Members
	m := &g.Members[0]
	m.Calculation = workspace.WriteGroupMemberCalculationDelta
	m.Cadence = &workspace.WriteGroupMemberCadence{Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
	m.ValueRole = workspace.WriteGroupValueRoleCounterRaw
	m.SourceValueType = schema.DataTypeUint64
	m.ResultValueType = measurement.ExactDecimal
	g.RecordID = "84ec2c79-2fd1-4f94-a514-5078b8b6d2af"
	g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileMixed, ResultShape: workspace.WriteGroupResultShapeMemberCalculationV1}
	g.RowPolicy.IntervalSeconds = 0
	f.config.FirstBucket = time.Time{}
	f.config.StartsAt = boundaryStart
	facts := workspace.CounterBindingFacts{Member: workspace.CounterBindingMemberIdentity{WorkspaceID: g.WorkspaceID, GroupID: g.ID, AppliedRevision: g.Revision, DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}, MeasurementID: *m.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1", SemanticKind: measurement.SemanticKindCounter, Unit: "kWh", SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision, SourceValueType: m.SourceValueType, ResultValueType: measurement.ExactUint64, ScaleMode: workspace.CounterScaleUnscaled, ExactScale: "1"}
	if scaled {
		m.ExactRecordingSource = &workspace.ExactRecordingSource{MeasurementID: *m.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1", SemanticKind: measurement.SemanticKindCounter, Unit: "kWh", OriginalDataType: schema.DataTypeUint64, Scale: "0.1"}
		facts.ScaleMode = workspace.CounterScaleExact
		facts.ExactScale = "0.1"
		facts.ResultValueType = measurement.ExactDecimal
	}
	proof, err := workspace.NewCounterBindingProof(facts)
	require.NoError(t, err)
	f.config.CounterBindings = workspace.CounterBindingProofs{proof.Member.Key(): proof}
	return f
}

func TestMixedBoundarySealedDeltaAndExactScale(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		t.Run(map[bool]string{false: "unscaled", true: "scaled"}[scaled], func(t *testing.T) {
			f := mixedCounterBoundaryFixture(t, scaled)
			b, err := NewGroupBoundary(f.config)
			require.NoError(t, err)
			f.accept(t, b, f.envelope(0, "base", boundaryAt(14), uint64(100)))
			require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
			require.Empty(t, f.sink.rows)
			f.accept(t, b, f.envelope(0, "next", boundaryAt(29), uint64(108)))
			require.NoError(t, b.Tick(t.Context(), boundaryAt(30)))
			require.Len(t, f.sink.rows, 1, "reported=%+v checkpoints=%+v", f.sink.outcomes, b.assembler.CounterCheckpoints())
			result := f.sink.rows[0].Outcome.Members[0]
			require.True(t, result.Usable())
			require.NotNil(t, result.Counter)
			digits, ok := result.Counter.Delta.Value.DecimalDigits()
			require.True(t, ok)
			if scaled {
				require.Equal(t, "0.8", digits)
			} else {
				require.Equal(t, "8", digits)
			}
		})
	}
}

func TestMixedBoundaryRejectsStaleFrozenProofsAndObservationClaims(t *testing.T) {
	for _, field := range []string{"source", "mapping", "measurement", "scale", "definition", "foreign proof"} {
		t.Run(field, func(t *testing.T) {
			f := mixedCounterBoundaryFixture(t, true)
			g := f.config.Group
			m := &g.Members[0]
			switch field {
			case "source":
				m.SourceRevision = "new-source"
			case "mapping":
				m.MappingRevision = "new-mapping"
			case "measurement":
				m.MeasurementID = ptr("foreign")
			case "scale":
				m.ExactRecordingSource.Scale = "0.01"
			case "definition":
				m.ExactRecordingSource.DefinitionRevision = "def-2"
			case "foreign proof":
				g.ID = "foreign"
			}
			_, err := NewGroupBoundary(f.config)
			require.Error(t, err)
		})
	}
	for _, field := range []string{"workspace", "source", "mapping", "measurement", "definition", "binding", "epoch"} {
		t.Run("observation/"+field, func(t *testing.T) {
			f := mixedCounterBoundaryFixture(t, true)
			b, err := NewGroupBoundary(f.config)
			require.NoError(t, err)
			e := f.envelope(0, "foreign", boundaryAt(14), uint64(100))
			switch field {
			case "workspace":
				e.WorkspaceID = "foreign"
			case "source":
				e.SourceRevision = "foreign"
			case "mapping":
				e.MappingRevision = "foreign"
			case "measurement":
				e.MeasurementID = "foreign"
			case "definition":
				e.DefinitionRevision = "foreign"
			case "binding":
				e.SourceBindingRevision = "foreign"
			case "epoch":
				e.SeriesEpoch = "foreign"
			}
			require.Error(t, b.AcceptSample(t.Context(), e))
			require.Empty(t, b.assembler.CounterCheckpoints())
		})
	}
}

func TestMixedLatestExactSourceIsExplicitlyBlockedWithoutReinterpretation(t *testing.T) {
	f := mixedCounterBoundaryFixture(t, true)
	member := &f.config.Group.Members[0]
	member.Calculation = workspace.WriteGroupMemberCalculationLatest
	f.config.CounterBindings = nil
	_, err := NewGroupBoundary(f.config)
	require.ErrorIs(t, err, workspace.ErrCounterBindingProofConflict)
	require.Equal(t, "0.1", member.ExactRecordingSource.Scale)
	require.Equal(t, schema.DataTypeUint64, member.SourceValueType)
	require.Equal(t, measurement.ExactDecimal, member.ResultValueType)
	require.Empty(t, f.sink.rows)
}
