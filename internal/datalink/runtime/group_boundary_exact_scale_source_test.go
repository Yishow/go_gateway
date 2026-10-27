package runtime

import (
	"testing"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func TestExactRecordingScaleSourceNeedsFrozenProof(t *testing.T) {
	definition := confirmedDefinition(measurement.SemanticKindCounter, nil)
	definition.NumericEncoding = "uint64"
	definition.ExactRecordingScale = "0.1"
	require.Equal(t, measurement.CounterConfirmed, measurement.ConfirmCounter(definition))
	_, ok := CounterSourceFromDefinition(definition)
	require.False(t, ok, "an editable definition cannot silently activate an unscaled counter when exact scale requires frozen proof")
}
