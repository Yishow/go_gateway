package workspace

import (
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func transitionMatrixExactPair() (*WriteGroup, *WriteGroup, map[string]schema.DataType) {
	source, _, types := transitionMatrixCumulativePair()
	_, member := exactCounterBindingFixture()
	member.CircuitID, member.TargetColumn = "circuit-1", "counter_value"
	source.Members[0] = member
	target := cloneWriteGroup(source)
	target.ID, target.Destination.ConnectorID, target.Destination.ConnectorRevision = "G2", "B", "B1"
	return source, target, types
}

func TestWriteGroupTransitionMatrixExactMeaningIsPreserved(t *testing.T) {
	source, target, types := transitionMatrixExactPair()
	require.NoError(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types))
	for name, mutate := range map[string]func(*WriteGroupMember){
		"definition revision": func(m *WriteGroupMember) { m.ExactRecordingSource.DefinitionRevision = "def-2" },
		"binding revision":    func(m *WriteGroupMember) { m.ExactRecordingSource.SourceBindingRevision = "binding-2" },
		"epoch":               func(m *WriteGroupMember) { m.ExactRecordingSource.SeriesEpoch = "epoch-2" },
		"scale":               func(m *WriteGroupMember) { m.ExactRecordingSource.Scale = "0.2" },
		"unit":                func(m *WriteGroupMember) { m.ExactRecordingSource.Unit = "MWh" },
		"kind":                func(m *WriteGroupMember) { m.ExactRecordingSource.SemanticKind = measurement.SemanticKindSignedCounter },
		"rollover policy":     func(m *WriteGroupMember) { m.ExactRecordingSource.CounterPolicy.RolloverModulus = "100" },
		"rate policy":         func(m *WriteGroupMember) { m.ExactRecordingSource.CounterPolicy.MaxRatePerSecond = "1" },
		"measurement": func(m *WriteGroupMember) {
			id := "counter-2"
			m.MeasurementID = &id
			m.ExactRecordingSource.MeasurementID = id
		},
		"remove exact meaning": func(m *WriteGroupMember) { m.ExactRecordingSource = nil; m.ResultValueType = measurement.ExactUint64 },
	} {
		t.Run(name, func(t *testing.T) {
			source, target, types := transitionMatrixExactPair()
			mutate(&target.Members[0])
			requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), "recording-meaning-change")
		})
	}
}

func TestWriteGroupTransitionMatrixNoEquivalentWideScaledSource(t *testing.T) {
	source, target, types := transitionMatrixExactPair()
	target.LayoutVersion, target.StorageMode = WriteGroupLayoutManagedWideProfileV1, ""
	requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), "unsupported-shape")
	// Coercing both to raw does not make an attached exact proof wide-compatible.
	for _, group := range []*WriteGroup{source, target} {
		group.LayoutVersion, group.StorageMode = WriteGroupLayoutManagedWideProfileV1, ""
		group.Members[0].ResultValueType = measurement.ExactUint64
	}
	requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), "unsupported-shape")
}

func TestWriteGroupTransitionMatrixMalformedExactSourceStaysBlocked(t *testing.T) {
	source, target, types := transitionMatrixExactPair()
	source.Members[0].ExactRecordingSource.Scale = "private://secret"
	target.Members[0].ExactRecordingSource.Scale = "private://secret"
	requireTransitionReason(t, ValidateWriteGroupTransitionStructure(source, target, "sqlite", "sqlite", types), "unsupported-shape")
}
