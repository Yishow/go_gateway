package runtime

import (
	"testing"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestExactRecordingScaleBoundarySample(t *testing.T) {
	member := boundaryMember{
		key: "source", kind: measurement.ExactUint64, exactScale: "0.1",
		sourceRevision: "source-1", mappingRevision: "mapping-1",
		counter: &aggregation.CounterBinding{
			MeasurementID: "measurement", DefinitionRevision: "def-1",
			SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
		},
	}
	envelope := measurement.SampleEnvelope{
		SampleID: "sample", Quality: schema.QualityGood,
		SourceRevision: "source-1", MappingRevision: "mapping-1",
		MeasurementID: "measurement", DefinitionRevision: "def-1",
		SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
	}
	for _, input := range []struct {
		value uint64
		want  string
	}{{120000, "12000.0"}, {120004, "12000.4"}} {
		envelope.Value = input.value
		sample, err := sampleFrom(member, envelope)
		require.NoError(t, err)
		require.Equal(t, schema.QualityGood, sample.Quality)
		require.Equal(t, measurement.ExactDecimal, sample.Value.Type())
		digits, ok := sample.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, input.want, digits, "the journal input is scaled once without floating point")
	}

	for _, bad := range []any{120004.0, "120004", int64(-1)} {
		envelope.Value = bad
		sample, err := sampleFrom(member, envelope)
		require.NoError(t, err)
		require.Equal(t, schema.QualityBad, sample.Quality)
		require.Empty(t, sample.Value.Type())
	}
	envelope.Value = uint64(120004)
	for _, changed := range []measurement.SampleEnvelope{
		func() measurement.SampleEnvelope { e := envelope; e.SourceRevision = ""; return e }(),
		func() measurement.SampleEnvelope { e := envelope; e.MappingRevision = "mapping-other"; return e }(),
	} {
		_, err := sampleFrom(member, changed)
		require.Error(t, err, "source and mapping proof are required before applying scale")
	}
	envelope.DefinitionRevision = "def-stale"
	_, err := sampleFrom(member, envelope)
	require.Error(t, err)

	envelope.DefinitionRevision = "def-1"
	member.exactScale = ""
	sample, err := sampleFrom(member, envelope)
	require.NoError(t, err)
	require.Equal(t, measurement.ExactUint64, sample.Value.Type(), "unscaled sources keep their original type")
}
