package groupdelivery

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMixedOutboxPayloadRoundTripSeparatesDispatchAndWindow(t *testing.T) {
	for _, long := range []bool{false, true} {
		t.Run(map[bool]string{false: "wide", true: "long"}[long], func(t *testing.T) {
			row := *rowBucket("", t0, 9007199254740993).Row
			row.DueAt = at(60)
			if !long {
				row.BucketStart = time.Time{}
			}
			payload, digest, err := encodeRowPayload(row)
			require.NoError(t, err)
			require.NotEmpty(t, digest)
			var raw map[string]any
			require.NoError(t, json.Unmarshal(payload, &raw))
			require.EqualValues(t, 2, raw["version"])
			require.Equal(t, timestamp(at(60)), raw["due_at"])
			if !long {
				require.Equal(t, "", raw["bucket_start"], "wide has no common window")
			}
			restored, err := DecodeRowPayload(payload)
			require.NoError(t, err)
			require.Equal(t, row, restored)
			if !long {
				require.True(t, restored.BucketStart.IsZero())
			}
			changed := bytes.Replace(payload, []byte(`"due_at":"`+timestamp(at(60))+`"`), []byte(`"due_at":"invalid"`), 1)
			_, err = DecodeRowPayload(changed)
			require.Error(t, err)
		})
	}
}

func TestLegacyOutboxPayloadDoesNotGainMixedDispatch(t *testing.T) {
	row := *rowBucket("", t0, 17).Row
	payload, _, err := encodeRowPayload(row)
	require.NoError(t, err)
	require.NotContains(t, string(payload), "due_at")
	var wire map[string]any
	require.NoError(t, json.Unmarshal(payload, &wire))
	require.EqualValues(t, 1, wire["version"])
	wire["due_at"] = timestamp(at(60))
	forged, err := json.Marshal(wire)
	require.NoError(t, err)
	_, err = DecodeRowPayload(forged)
	require.Error(t, err, "legacy payload cannot silently acquire mixed timestamp semantics")
}
