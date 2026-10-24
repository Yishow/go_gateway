package runtime

import (
	"testing"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func confirmedDefinition(kind measurement.SemanticKind, policy *measurement.CounterPolicy) measurement.MeasurementDefinition {
	return measurement.MeasurementDefinition{
		ID: "meas-1", DefinitionRevision: "def-3", SourceBindingRevision: "bind-2", SeriesEpoch: "epoch-4",
		Unit: "kWh", SemanticKind: kind, CounterPolicy: policy,
	}
}

func TestCounterSourceFromDefinitionOnlyTrustsConfirmedSemantics(t *testing.T) {
	source, ok := CounterSourceFromDefinition(confirmedDefinition(measurement.SemanticKindCounter, nil))
	require.True(t, ok)
	require.Equal(t, CounterSource{
		MeasurementID: "meas-1", DefinitionRevision: "def-3", SourceBindingRevision: "bind-2", SeriesEpoch: "epoch-4",
		Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
	}, source)

	signed, ok := CounterSourceFromDefinition(confirmedDefinition(measurement.SemanticKindSignedCounter, &measurement.CounterPolicy{AllowNegative: true}))
	require.True(t, ok)
	require.Equal(t, aggregation.CounterKindSigned, signed.Kind)

	maxRollOver := 65536.0
	withFloat, ok := CounterSourceFromDefinition(confirmedDefinition(measurement.SemanticKindCounter, &measurement.CounterPolicy{MaxRollOver: &maxRollOver}))
	require.True(t, ok)
	require.Nil(t, withFloat.RolloverModulus, "a float MaxRollOver is never taken as exact evidence")
	require.Nil(t, withFloat.MaxRatePerSecond)

	for name, def := range map[string]measurement.MeasurementDefinition{
		"signed without the confirmation":  confirmedDefinition(measurement.SemanticKindSignedCounter, nil),
		"signed with negatives disallowed": confirmedDefinition(measurement.SemanticKindSignedCounter, &measurement.CounterPolicy{}),
		"a counter that allows negatives":  confirmedDefinition(measurement.SemanticKindCounter, &measurement.CounterPolicy{AllowNegative: true}),
		"a gauge named like a meter":       confirmedDefinition(measurement.SemanticKindGauge, nil),
		"a rate":                           confirmedDefinition(measurement.SemanticKindRate, nil),
		"no unit": func() measurement.MeasurementDefinition {
			d := confirmedDefinition(measurement.SemanticKindCounter, nil)
			d.Unit = " "
			return d
		}(),
		"no series epoch": func() measurement.MeasurementDefinition {
			d := confirmedDefinition(measurement.SemanticKindCounter, nil)
			d.SeriesEpoch = ""
			return d
		}(),
		"no definition revision": func() measurement.MeasurementDefinition {
			d := confirmedDefinition(measurement.SemanticKindCounter, nil)
			d.DefinitionRevision = ""
			return d
		}(),
		"no source binding revision": func() measurement.MeasurementDefinition {
			d := confirmedDefinition(measurement.SemanticKindCounter, nil)
			d.SourceBindingRevision = ""
			return d
		}(),
	} {
		_, ok := CounterSourceFromDefinition(def)
		require.False(t, ok, name)
	}
}

func TestCounterSourceFromDefinitionMapsTheExactRolloverEvidence(t *testing.T) {
	definition := func(policy *measurement.CounterPolicy) measurement.MeasurementDefinition {
		def := confirmedDefinition(measurement.SemanticKindCounter, policy)
		def.NumericEncoding = "uint16"
		return def
	}
	digits := func(value *measurement.ExactValue) string {
		require.NotNil(t, value)
		text, ok := value.CanonicalDigits()
		require.True(t, ok)
		return text
	}

	source, ok := CounterSourceFromDefinition(definition(&measurement.CounterPolicy{RolloverModulus: "9007199254740993", MaxRatePerSecond: "0.125"}))
	require.True(t, ok)
	require.Equal(t, "9007199254740993", digits(source.RolloverModulus), "exact above 2^53, never through a float")
	require.Equal(t, "0.125", digits(source.MaxRatePerSecond))

	modulusOnly, ok := CounterSourceFromDefinition(definition(&measurement.CounterPolicy{RolloverModulus: "65536"}))
	require.True(t, ok)
	require.Equal(t, "65536", digits(modulusOnly.RolloverModulus))
	require.Nil(t, modulusOnly.MaxRatePerSecond, "a rate is never guessed for a modulus alone")

	legacy := 65536.0
	floatOnly, ok := CounterSourceFromDefinition(definition(&measurement.CounterPolicy{MaxRollOver: &legacy}))
	require.True(t, ok)
	require.Nil(t, floatOnly.RolloverModulus, "the legacy float is not exact evidence")
	require.Nil(t, floatOnly.MaxRatePerSecond)

	for name, def := range map[string]measurement.MeasurementDefinition{
		"a zero modulus saved before validation existed": definition(&measurement.CounterPolicy{RolloverModulus: "0"}),
		"an exponent":                             definition(&measurement.CounterPolicy{RolloverModulus: "6.5536e4"}),
		"a rate without the modulus":              definition(&measurement.CounterPolicy{MaxRatePerSecond: "100"}),
		"a negative rate":                         definition(&measurement.CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "-1"}),
		"evidence on a negative-allowing counter": definition(&measurement.CounterPolicy{RolloverModulus: "65536", AllowNegative: true}),
		"evidence on a float source": func() measurement.MeasurementDefinition {
			def := definition(&measurement.CounterPolicy{RolloverModulus: "65536"})
			def.NumericEncoding = "float64"
			return def
		}(),
		"evidence on a signed counter": confirmedDefinition(measurement.SemanticKindSignedCounter, &measurement.CounterPolicy{AllowNegative: true, RolloverModulus: "65536"}),
	} {
		_, ok := CounterSourceFromDefinition(def)
		require.False(t, ok, "damaged or contradictory saved evidence fails closed: %s", name)
	}
}
