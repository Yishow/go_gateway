package groupdelivery

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// legacySamplePayload is a journal row written before ReceivedAt existed.
const legacySamplePayload = `{"sample_id":"s-1","member_key":"m-1","observed_at":"2026-01-01T00:00:02Z","quality":"good","value":{"type":"int64","encoding":"decimal-string","value":"7"},"source_revision":"src-1","mapping_revision":"map-1"}`

func TestSampleCodecKeepsReceivedAtCompatibly(t *testing.T) {
	t.Run("a journal row without a receive time decodes as unknown", func(t *testing.T) {
		sample, err := decodeSample([]byte(legacySamplePayload))
		require.NoError(t, err)
		require.True(t, sample.ReceivedAt.IsZero(), "unknown stays unknown; it is never replaced by the restart time")
		require.Equal(t, "s-1", sample.SampleID)
	})

	t.Run("the receive time round-trips in UTC with sub-second precision", func(t *testing.T) {
		original := sampleAt("s-2", "m-1", 4, 9)
		original.ReceivedAt = at(4).Add(250*time.Millisecond + 7*time.Microsecond).In(time.FixedZone("UTC+8", 8*3600))
		payload, _, err := encodeSample(original)
		require.NoError(t, err)
		decoded, err := decodeSample(payload)
		require.NoError(t, err)
		require.True(t, original.ReceivedAt.Equal(decoded.ReceivedAt))
		require.Equal(t, time.UTC, decoded.ReceivedAt.Location())
	})

	t.Run("an unknown receive time is omitted, and the identity digest ignores it", func(t *testing.T) {
		withoutReceived := sampleAt("s-3", "m-1", 5, 3)
		payload, digestWithout, err := encodeSample(withoutReceived)
		require.NoError(t, err)
		require.NotContains(t, string(payload), "received_at")

		withReceived := withoutReceived
		withReceived.ReceivedAt = at(5).Add(time.Second)
		payloadWith, digestWith, err := encodeSample(withReceived)
		require.NoError(t, err)
		require.Contains(t, string(payloadWith), "received_at")
		require.Equal(t, digestWithout, digestWith, "a recorded receive time must not turn a resend into a conflict")
	})

	t.Run("a legacy sample re-encodes to the same bytes and the same digest", func(t *testing.T) {
		sample, err := decodeSample([]byte(legacySamplePayload))
		require.NoError(t, err)
		payload, digest, err := encodeSample(sample)
		require.NoError(t, err)
		require.Equal(t, legacySamplePayload, string(payload))
		// Before ReceivedAt existed the digest was the hash of the payload itself.
		legacy := sha256.Sum256([]byte(legacySamplePayload))
		require.Equal(t, hex.EncodeToString(legacy[:]), digest)
	})
}

func TestJournalKeepsTheFirstReceiveTimeOfAResentSample(t *testing.T) {
	store, db, _ := newTestStore(t)
	ledger := store.Ledger(testKey)
	first := sampleAt("s-1", "m-1", 2, 7)
	first.ReceivedAt = at(2).Add(100 * time.Millisecond)
	require.NoError(t, ledger.AppendSample(t.Context(), at(0), first))

	resent := first
	resent.ReceivedAt = at(30) // a later resend of the same acquisition
	require.NoError(t, ledger.AppendSample(t.Context(), at(0), resent), "same identity is a duplicate, not a conflict")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))

	restored, err := ledger.Restore(t.Context())
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.True(t, restored.Samples[0].ReceivedAt.Equal(first.ReceivedAt), "the journaled receive time is the original one")

	changed := first
	changed.Value = sampleAt("s-1", "m-1", 2, 8).Value
	require.ErrorIs(t, ledger.AppendSample(t.Context(), at(0), changed), ErrSampleConflict, "a different value is still a conflict")
}
