package aggregation

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func scaleBinding(scale string) CounterBinding {
	return CounterBinding{
		WorkspaceID: "ws", GroupID: "group", ProfileRevision: "profile-rev-1", MemberKey: "device/point/tag",
		MeasurementID: "meas-1", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", Kind: CounterKindMonotonic, Unit: "kWh", ExactRecordingScale: scale,
	}
}

func legacyScaleKeyParts() []string {
	return []string{
		counterBaselineKeyVersion, "ws", "group", "profile-rev-1", "device/point/tag", "meas-1",
		"def-1", "bind-1", "epoch-1", string(CounterKindMonotonic), "kWh", "none", "none",
	}
}

func TestCounterBindingExactRecordingScaleKey(t *testing.T) {
	legacyWant, err := measurement.HashStringTuple(legacyScaleKeyParts()...)
	require.NoError(t, err)
	legacy, err := scaleBinding("").Key()
	require.NoError(t, err)
	require.Equal(t, legacyWant, legacy, "an absent scale preserves the exact legacy key")

	first, err := scaleBinding("0.1").Key()
	require.NoError(t, err)
	require.NotEqual(t, legacy, first)
	scaledWant, err := measurement.HashStringTuple(append(legacyScaleKeyParts(), "exact-recording-scale-v1", "0.1")...)
	require.NoError(t, err)
	require.Equal(t, scaledWant, first)

	equivalent, err := scaleBinding("0.10").Key()
	require.NoError(t, err)
	require.Equal(t, first, equivalent, "numerically equal scales share a key")

	changed, err := scaleBinding("0.2").Key()
	require.NoError(t, err)
	require.NotEqual(t, first, changed)

	otherBinding := scaleBinding("0.1")
	otherBinding.DefinitionRevision = "def-2"
	otherKey, err := otherBinding.Key()
	require.NoError(t, err)
	require.NotEqual(t, first, otherKey, "existing binding fields remain bound")
}

func TestCounterBindingExactRecordingScaleValidation(t *testing.T) {
	for _, scale := range []string{"0.1", "0.10", "1", "12000.4"} {
		require.NoError(t, scaleBinding(scale).Validate(), scale)
	}
	require.NoError(t, scaleBinding("").Validate(), "empty scale is legacy")

	for _, scale := range []string{"0", "-0.1", "1e-2", "NaN", "Infinity", " 0.1", "0x10", "abc"} {
		err := scaleBinding(scale).Validate()
		require.Error(t, err, scale)
	}
}

func TestCounterBindingExactRecordingScaleJSONRoundTrip(t *testing.T) {
	binding := scaleBinding("0.10")
	raw, err := json.Marshal(binding)
	require.NoError(t, err)
	var decoded CounterBinding
	require.NoError(t, json.Unmarshal(raw, &decoded))
	require.Equal(t, "0.10", decoded.ExactRecordingScale)
	key, err := decoded.Key()
	require.NoError(t, err)
	original, err := binding.Key()
	require.NoError(t, err)
	require.Equal(t, original, key)
}
