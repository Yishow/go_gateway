package aggregation

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

// TestCounterCheckpointDecodeRejectsTrailingContent pins that a baseline is read
// only from exactly one JSON object: anything after it, valid or not, means the
// stored bytes are damaged and must not be guessed into a baseline.
func TestCounterCheckpointDecodeRejectsTrailingContent(t *testing.T) {
	_, checkpoint := closeSeries(t, counterBindingFixture(),
		[]CounterReading{counterReading("a", 0, time.Minute, measurement.NewUint64(9007199254740993))})
	encoded, err := EncodeCounterCheckpoint(*checkpoint)
	require.NoError(t, err)

	t.Run("a single object, with or without surrounding whitespace, is accepted", func(t *testing.T) {
		for name, suffix := range map[string]string{"nothing": "", "newline": "\n", "mixed whitespace": " \t\r\n "} {
			decoded, err := DecodeCounterCheckpoint(append(append([]byte{}, encoded...), suffix...))
			require.NoError(t, err, name)
			require.True(t, checkpoint.Reading.Value.Equal(decoded.Reading.Value), name)
		}
		decoded, err := DecodeCounterCheckpoint(append([]byte("\n  "), encoded...))
		require.NoError(t, err, "leading whitespace")
		require.Equal(t, checkpoint.Reading.SampleID, decoded.Reading.SampleID)
	})

	t.Run("damaged or extra content after the object is refused", func(t *testing.T) {
		for name, suffix := range map[string]string{
			"a stray letter":          "x",
			"a lone backslash":        `\`,
			"an unopened close":       "}",
			"a comma":                 ",",
			"an unterminated string":  `"abc`,
			"a second object":         `{"version":1}`,
			"a second scalar":         "1",
			"a second null":           "null",
			"garbage after a newline": "\n\x00",
			"a truncated second":      `{"version":`,
		} {
			t.Run(name, func(t *testing.T) {
				_, err := DecodeCounterCheckpoint(append(append([]byte{}, encoded...), suffix...))
				require.ErrorIs(t, err, ErrCounterCheckpointInvalid)
			})
		}
	})
}
